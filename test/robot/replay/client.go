// Copyright (C) 2017 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package replay

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/gapid/core/app/crash"
	"github.com/google/gapid/core/app/layout"
	"github.com/google/gapid/core/event/task"
	"github.com/google/gapid/core/log"
	"github.com/google/gapid/core/os/device"
	"github.com/google/gapid/core/os/device/bind"
	"github.com/google/gapid/core/os/device/host"
	"github.com/google/gapid/core/os/file"
	"github.com/google/gapid/core/os/shell"
	"github.com/google/gapid/test/robot/job"
	"github.com/google/gapid/test/robot/job/worker"
	"github.com/google/gapid/test/robot/stash"
)

const (
	replayTimeout = time.Hour
	// this string is returned when GAPIT fails to connect to the GAPIS, particularly due to ETXTBSY
	// look at https://github.com/google/gapid/pull/933 for more information
	retryString = "Failed to connect to the GAPIS server"
)

type client struct {
	store   *stash.Client
	manager Manager
	tempDir file.Path
}

// Run starts new replay client if any hardware is available.
func Run(ctx context.Context, store *stash.Client, manager Manager, tempDir file.Path) error {
	c := &client{store: store, manager: manager, tempDir: tempDir}
	job.OnDeviceAdded(ctx, c.onDeviceAdded)
	host := host.Instance(ctx)
	return manager.Register(ctx, host, host, c.replay)
}

func (c *client) onDeviceAdded(ctx context.Context, host *device.Instance, target bind.Device) {
	replayOnTarget := func(ctx context.Context, t *Task) error {
		job.LockDevice(ctx, target)
		defer func() {
			// HACK: kill gapid.apk manually for now as subsequent reports/replays may freeze the app.
			// Remove when https://github.com/google/gapid/issues/1666 is fixed.
			target.Shell("am", "force-stop", "com.google.android.gapid.aarch64").Run(ctx)
			target.Shell("am", "force-stop", "com.google.android.gapid.armeabi").Run(ctx)
		}()
		defer job.UnlockDevice(ctx, target)
		if target.Status() != bind.Status_Online {
			log.I(ctx, "Trying to replay %s on %s not started, device status %s",
				t.Input.Trace, target.Instance().GetSerial(), target.Status().String())
			return nil
		}
		return c.replay(ctx, t)
	}
	crash.Go(func() {
		if err := c.manager.Register(ctx, host, target.Instance(), replayOnTarget); err != nil {
			log.E(ctx, "Error running replay client: %v", err)
		}
	})
}

func (c *client) replay(ctx context.Context, t *Task) error {
	if err := c.manager.Update(ctx, t.Action, job.Running, nil); err != nil {
		return err
	}
	var output *Output
	err := worker.RetryFunction(ctx, 4, time.Millisecond*100, func() (err error) {
		ctx, cancel := task.WithTimeout(ctx, replayTimeout)
		defer cancel()
		output, err = doReplay(ctx, t.Action, t.Input, c.store, c.tempDir)
		return err
	})

	status := job.Succeeded
	if err != nil {
		status = job.Failed
		log.E(ctx, "Error running replay: %v", err)
	} else if output.Err != "" {
		status = job.Failed
		log.E(ctx, "Error during replay: %v", output.Err)
	}

	return c.manager.Update(ctx, t.Action, status, output)
}

// doReplay extracts input files and runs `gapit video` on them, capturing the output. The output object will
// be partially filled in the event of an upload error from store in order to allow examination of the logs.
func doReplay(ctx context.Context, action string, in *Input, store *stash.Client, tempDir file.Path) (*Output, error) {
	tracefile := tempDir.Join(action + ".gfxtrace")
	videofile := tempDir.Join(action + "_replay.mp4")

	extractedDir := tempDir.Join(action + "_tools")
	extractedLayout, err := layout.NewPkgLayout(extractedDir, true)
	if err != nil {
		return nil, err
	}

	gapit, err := extractedLayout.Gapit(ctx)
	if err != nil {
		return nil, err
	}
	gapir, err := extractedLayout.Gapir(ctx)
	if err != nil {
		return nil, err
	}
	gapis, err := extractedLayout.Gapis(ctx)
	if err != nil {
		return nil, err
	}
	vscLib, err := extractedLayout.Json(ctx, layout.LibVirtualSwapChain)
	if err != nil {
		return nil, err
	}
	vscJson, err := extractedLayout.Json(ctx, layout.LibVirtualSwapChain)
	if err != nil {
		return nil, err
	}

	defer func() {
		file.Remove(tracefile)
		file.Remove(videofile)
		file.RemoveAll(extractedDir)
	}()

	for _, file := range []struct {
		in  string
		out file.Path
	}{
		{in.Trace, tracefile},
		{in.Gapit, gapit},
		{in.Gapis, gapis},
		{in.Gapir, gapir},
		{in.VirtualSwapChainLib, vscLib},
		{in.VirtualSwapChainJson, vscJson},
	} {
		if err := store.GetFile(ctx, file.in, file.out); err != nil {
			return nil, err
		}
	}

	if in.GetToolingLayout() != nil {
		gapidAPK, err := extractedLayout.GapidApk(ctx, in.GetToolingLayout().GetGapidAbi())
		if err != nil {
			return nil, err
		}
		if err := store.GetFile(ctx, in.GapidApk, gapidAPK); err != nil {
			return nil, err
		}
	}

	params := []string{
		"video",
		"-gapir-device", in.GetGapirDevice(),
		"-frames-minimum", "10",
		"-type", "sxs",
		"-out", videofile.System(),
		tracefile.System(),
	}
	cmd := shell.Command(gapit.System(), params...)
	outBuf := &bytes.Buffer{}
	errBuf := &bytes.Buffer{}
	outputObj := &Output{}
	errs := []string{}
	log.I(ctx, "Running replay action %s", cmd)
	if err := cmd.Capture(outBuf, errBuf).Run(ctx); err != nil {
		if err := worker.NeedsRetry(err.Error()); err != nil {
			return nil, err
		}
		errs = append(errs, err.Error())
	}

	errs = append(errs, strings.TrimSpace(errBuf.String()))
	outputObj.Err = strings.Join(errs, "\n")
	if err := worker.NeedsRetry(outputObj.Err, retryString); err != nil {
		return nil, err
	}

	output := fmt.Sprintf("%s\n\n%s\n%s\n\n", cmd, strings.TrimSpace(outBuf.String()), strings.TrimSpace(errBuf.String()))
	log.I(ctx, output)
	logID, err := store.UploadString(ctx, stash.Upload{Name: []string{"replay.log"}, Type: []string{"text/plain"}}, output)
	if err != nil {
		return outputObj, err
	}
	outputObj.Log = logID
	videoID, err := store.UploadFile(ctx, videofile)
	if err != nil {
		return outputObj, err
	}
	outputObj.Video = videoID
	return outputObj, nil
}
