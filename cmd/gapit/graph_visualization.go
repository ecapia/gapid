// Copyright (C) 2018 Google Inc.
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

package main

import (
	"context"
	"flag"
	"github.com/google/gapid/core/app"
	"github.com/google/gapid/core/log"
	"io/ioutil"
	"path/filepath"
)

type graph_visualizationVerb struct{ GraphVisualizationFlags }

func init() {
	verb := &graph_visualizationVerb{}
	app.AddVerb(&app.Verb{
		Name:      "graph_visualization",
		ShortHelp: "Get and Write graph-visualization file from capture",
		Action:    verb,
	})
}
func getLastLevelNameFromCapturePath(capturePath string) string {
	for i := len(capturePath) - 1; i >= 0; i-- {
		if capturePath[i] == '/' {
			return capturePath[i+1:]
		}
	}
	return capturePath
}
func (verb *graph_visualizationVerb) Run(ctx context.Context, flags flag.FlagSet) error {

	if flags.NArg() != 1 {
		app.Usage(ctx, "Exactly two parameters expected:  trace file path and output format, got %d", flags.NArg())
		return nil
	}

	capturePath, err := filepath.Abs(flags.Arg(0))
	if err != nil {
		return log.Errf(ctx, err, "Finding file: %v", flags.Arg(0))
	}

	client, err := getGapis(ctx, verb.Gapis, verb.Gapir)
	if err != nil {
		return log.Err(ctx, err, "Failed to connect to the GAPIS server")
	}
	defer client.Close()

	capture, err := client.LoadCapture(ctx, capturePath)
	if err != nil {
		return log.Errf(ctx, err, "LoadCapture(%v)", capturePath)
	}
	format := verb.Format
	graphVisualizationFile, err := client.GetGraphVisualizationFile(ctx, capture, format)
	if err != nil {
		return log.Errf(ctx, err, "ExportCapture(%v)", capture)
	}

	graphVisualizationName := verb.Out
	if graphVisualizationName == "" {
		graphVisualizationName = getLastLevelNameFromCapturePath(capturePath)
	}
	graphVisualizationName += "." + format

	if err := ioutil.WriteFile(graphVisualizationName, []byte(graphVisualizationFile), 0666); err != nil {
		return log.Errf(ctx, err, "Writing file: %v", graphVisualizationName)
	}
	return nil
}
