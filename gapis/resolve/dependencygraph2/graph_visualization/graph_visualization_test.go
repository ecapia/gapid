package graph_visualization

import (
	"reflect"
	"testing"
)

const (
	VK_BEGIN_COMMAND_BUFFER = "vkBeginCommandBuffer"
	VK_END_COMMAND_BUFFER   = "vkEndCommandBuffer"
	VK_COMMAND_BUFFER       = "vkCommandBuffer"
	VK_CMD_END_RENDER_PASS  = "vkCmdEndRenderPass"
	VK_RENDER_PASS          = "vkRenderPass"
	VK_CMD_NEXT_SUBPASS     = "vkCmdNextSubpass"
	VK_SUBPASS              = "vkSubpass"
)

func TestGetLabelFromHierarchy1(t *testing.T) {

	commandHierarchyNames, _ := getNameForCommandAndSubCommandHierarchy()

	auxiliarCommands := []string{
		"SetViewPort",
		"SetBarrier",
		"DrawIndex",
		"SetScissor",
		"SetIndexForDraw",
	}

	commands := []string{
		VK_BEGIN_COMMAND_BUFFER,
		auxiliarCommands[0],
		auxiliarCommands[0],
		VK_CMD_BEGIN_RENDER_PASS,
		auxiliarCommands[0],
		auxiliarCommands[1],
		VK_CMD_NEXT_SUBPASS,
		auxiliarCommands[2],
		auxiliarCommands[1],
		VK_CMD_NEXT_SUBPASS,
		auxiliarCommands[0],
		VK_CMD_END_RENDER_PASS,
		auxiliarCommands[2],
		auxiliarCommands[3],
		VK_CMD_BEGIN_RENDER_PASS,
		VK_CMD_END_RENDER_PASS,
		VK_CMD_BEGIN_RENDER_PASS,
		VK_CMD_END_RENDER_PASS,
		VK_END_COMMAND_BUFFER,

		auxiliarCommands[3],
		auxiliarCommands[4],

		VK_BEGIN_COMMAND_BUFFER,
		auxiliarCommands[0],
		auxiliarCommands[1],
		auxiliarCommands[2],
		auxiliarCommands[3],
		VK_END_COMMAND_BUFFER,

		VK_BEGIN_COMMAND_BUFFER,
		auxiliarCommands[4],
		VK_END_COMMAND_BUFFER,
	}
	wantedLabels := []*Label{
		&Label{name: []string{VK_COMMAND_BUFFER}, id: []int{1}},
		&Label{name: []string{VK_COMMAND_BUFFER}, id: []int{1}},
		&Label{name: []string{VK_COMMAND_BUFFER}, id: []int{1}},
		&Label{name: []string{VK_COMMAND_BUFFER, VK_RENDER_PASS}, id: []int{1, 1}},
		&Label{name: []string{VK_COMMAND_BUFFER, VK_RENDER_PASS, VK_SUBPASS}, id: []int{1, 1, 1}},
		&Label{name: []string{VK_COMMAND_BUFFER, VK_RENDER_PASS, VK_SUBPASS}, id: []int{1, 1, 1}},
		&Label{name: []string{VK_COMMAND_BUFFER, VK_RENDER_PASS, VK_SUBPASS}, id: []int{1, 1, 2}},
		&Label{name: []string{VK_COMMAND_BUFFER, VK_RENDER_PASS, VK_SUBPASS}, id: []int{1, 1, 2}},
		&Label{name: []string{VK_COMMAND_BUFFER, VK_RENDER_PASS, VK_SUBPASS}, id: []int{1, 1, 2}},
		&Label{name: []string{VK_COMMAND_BUFFER, VK_RENDER_PASS, VK_SUBPASS}, id: []int{1, 1, 3}},
		&Label{name: []string{VK_COMMAND_BUFFER, VK_RENDER_PASS, VK_SUBPASS}, id: []int{1, 1, 3}},
		&Label{name: []string{VK_COMMAND_BUFFER, VK_RENDER_PASS}, id: []int{1, 1}},
		&Label{name: []string{VK_COMMAND_BUFFER}, id: []int{1}},
		&Label{name: []string{VK_COMMAND_BUFFER}, id: []int{1}},
		&Label{name: []string{VK_COMMAND_BUFFER, VK_RENDER_PASS}, id: []int{1, 2}},
		&Label{name: []string{VK_COMMAND_BUFFER, VK_RENDER_PASS}, id: []int{1, 2}},
		&Label{name: []string{VK_COMMAND_BUFFER, VK_RENDER_PASS}, id: []int{1, 3}},
		&Label{name: []string{VK_COMMAND_BUFFER, VK_RENDER_PASS}, id: []int{1, 3}},
		&Label{name: []string{VK_COMMAND_BUFFER}, id: []int{1}},
		&Label{},
		&Label{},
		&Label{name: []string{VK_COMMAND_BUFFER}, id: []int{2}},
		&Label{name: []string{VK_COMMAND_BUFFER}, id: []int{2}},
		&Label{name: []string{VK_COMMAND_BUFFER}, id: []int{2}},
		&Label{name: []string{VK_COMMAND_BUFFER}, id: []int{2}},
		&Label{name: []string{VK_COMMAND_BUFFER}, id: []int{2}},
		&Label{name: []string{VK_COMMAND_BUFFER}, id: []int{2}},
		&Label{name: []string{VK_COMMAND_BUFFER}, id: []int{3}},
		&Label{name: []string{VK_COMMAND_BUFFER}, id: []int{3}},
		&Label{name: []string{VK_COMMAND_BUFFER}, id: []int{3}},
	}

	currentHierarchy := &Hierarchy{}
	obtainedLabels := []*Label{}
	for _, name := range commands {
		label := getLabelFromHierarchy(name, commandHierarchyNames, currentHierarchy)
		obtainedLabels = append(obtainedLabels, label)
	}
	for i := range wantedLabels {
		if reflect.DeepEqual(wantedLabels[i], obtainedLabels[i]) == false {
			t.Errorf("The label for command %s with id %d is different", commands[i], i)
			t.Errorf("Found %v wanted %v", obtainedLabels[i], wantedLabels[i])
		}
	}

}

func TestGetLabelFromHierarchy2(t *testing.T) {

	commandHierarchyNames, _ := getNameForCommandAndSubCommandHierarchy()

	auxiliarCommands := []string{
		"SetViewPort",
		"SetBarrier",
		"DrawIndex",
		"SetScissor",
		"SetIndexForDraw",
	}

	commands := []string{
		VK_BEGIN_COMMAND_BUFFER,
		auxiliarCommands[0],
		auxiliarCommands[0],
		VK_CMD_BEGIN_RENDER_PASS,
		auxiliarCommands[0],
		auxiliarCommands[1],
		VK_CMD_NEXT_SUBPASS,
		auxiliarCommands[2],
		auxiliarCommands[1],
		VK_CMD_NEXT_SUBPASS,
		auxiliarCommands[0],
		VK_CMD_END_RENDER_PASS,
		auxiliarCommands[2],
		auxiliarCommands[3],
		VK_END_COMMAND_BUFFER,

		auxiliarCommands[3],
		auxiliarCommands[4],

		VK_BEGIN_COMMAND_BUFFER,
		auxiliarCommands[0],
		auxiliarCommands[0],
		VK_CMD_BEGIN_RENDER_PASS,
		auxiliarCommands[0],
		auxiliarCommands[1],
		VK_CMD_NEXT_SUBPASS,
		auxiliarCommands[2],
		auxiliarCommands[1],
		VK_CMD_NEXT_SUBPASS,
		auxiliarCommands[0],
		VK_CMD_END_RENDER_PASS,
		auxiliarCommands[2],
		auxiliarCommands[3],
		VK_END_COMMAND_BUFFER,

		auxiliarCommands[3],
		auxiliarCommands[4],

		VK_BEGIN_COMMAND_BUFFER,
		auxiliarCommands[0],
		auxiliarCommands[0],
		VK_CMD_BEGIN_RENDER_PASS,
		auxiliarCommands[0],
		auxiliarCommands[1],
		VK_CMD_NEXT_SUBPASS,
		auxiliarCommands[2],
		auxiliarCommands[1],
		VK_CMD_NEXT_SUBPASS,
		auxiliarCommands[0],
		VK_CMD_END_RENDER_PASS,
		auxiliarCommands[2],
		auxiliarCommands[3],
		VK_END_COMMAND_BUFFER,

		auxiliarCommands[3],
		auxiliarCommands[4],
	}
	wantedLabels := []*Label{
		&Label{name: []string{VK_COMMAND_BUFFER}, id: []int{1}},
		&Label{name: []string{VK_COMMAND_BUFFER}, id: []int{1}},
		&Label{name: []string{VK_COMMAND_BUFFER}, id: []int{1}},
		&Label{name: []string{VK_COMMAND_BUFFER, VK_RENDER_PASS}, id: []int{1, 1}},
		&Label{name: []string{VK_COMMAND_BUFFER, VK_RENDER_PASS, VK_SUBPASS}, id: []int{1, 1, 1}},
		&Label{name: []string{VK_COMMAND_BUFFER, VK_RENDER_PASS, VK_SUBPASS}, id: []int{1, 1, 1}},
		&Label{name: []string{VK_COMMAND_BUFFER, VK_RENDER_PASS, VK_SUBPASS}, id: []int{1, 1, 2}},
		&Label{name: []string{VK_COMMAND_BUFFER, VK_RENDER_PASS, VK_SUBPASS}, id: []int{1, 1, 2}},
		&Label{name: []string{VK_COMMAND_BUFFER, VK_RENDER_PASS, VK_SUBPASS}, id: []int{1, 1, 2}},
		&Label{name: []string{VK_COMMAND_BUFFER, VK_RENDER_PASS, VK_SUBPASS}, id: []int{1, 1, 3}},
		&Label{name: []string{VK_COMMAND_BUFFER, VK_RENDER_PASS, VK_SUBPASS}, id: []int{1, 1, 3}},
		&Label{name: []string{VK_COMMAND_BUFFER, VK_RENDER_PASS}, id: []int{1, 1}},
		&Label{name: []string{VK_COMMAND_BUFFER}, id: []int{1}},
		&Label{name: []string{VK_COMMAND_BUFFER}, id: []int{1}},
		&Label{name: []string{VK_COMMAND_BUFFER}, id: []int{1}},
		&Label{},
		&Label{},

		&Label{name: []string{VK_COMMAND_BUFFER}, id: []int{2}},
		&Label{name: []string{VK_COMMAND_BUFFER}, id: []int{2}},
		&Label{name: []string{VK_COMMAND_BUFFER}, id: []int{2}},
		&Label{name: []string{VK_COMMAND_BUFFER, VK_RENDER_PASS}, id: []int{2, 1}},
		&Label{name: []string{VK_COMMAND_BUFFER, VK_RENDER_PASS, VK_SUBPASS}, id: []int{2, 1, 1}},
		&Label{name: []string{VK_COMMAND_BUFFER, VK_RENDER_PASS, VK_SUBPASS}, id: []int{2, 1, 1}},
		&Label{name: []string{VK_COMMAND_BUFFER, VK_RENDER_PASS, VK_SUBPASS}, id: []int{2, 1, 2}},
		&Label{name: []string{VK_COMMAND_BUFFER, VK_RENDER_PASS, VK_SUBPASS}, id: []int{2, 1, 2}},
		&Label{name: []string{VK_COMMAND_BUFFER, VK_RENDER_PASS, VK_SUBPASS}, id: []int{2, 1, 2}},
		&Label{name: []string{VK_COMMAND_BUFFER, VK_RENDER_PASS, VK_SUBPASS}, id: []int{2, 1, 3}},
		&Label{name: []string{VK_COMMAND_BUFFER, VK_RENDER_PASS, VK_SUBPASS}, id: []int{2, 1, 3}},
		&Label{name: []string{VK_COMMAND_BUFFER, VK_RENDER_PASS}, id: []int{2, 1}},
		&Label{name: []string{VK_COMMAND_BUFFER}, id: []int{2}},
		&Label{name: []string{VK_COMMAND_BUFFER}, id: []int{2}},
		&Label{name: []string{VK_COMMAND_BUFFER}, id: []int{2}},
		&Label{},
		&Label{},

		&Label{name: []string{VK_COMMAND_BUFFER}, id: []int{3}},
		&Label{name: []string{VK_COMMAND_BUFFER}, id: []int{3}},
		&Label{name: []string{VK_COMMAND_BUFFER}, id: []int{3}},
		&Label{name: []string{VK_COMMAND_BUFFER, VK_RENDER_PASS}, id: []int{3, 1}},
		&Label{name: []string{VK_COMMAND_BUFFER, VK_RENDER_PASS, VK_SUBPASS}, id: []int{3, 1, 1}},
		&Label{name: []string{VK_COMMAND_BUFFER, VK_RENDER_PASS, VK_SUBPASS}, id: []int{3, 1, 1}},
		&Label{name: []string{VK_COMMAND_BUFFER, VK_RENDER_PASS, VK_SUBPASS}, id: []int{3, 1, 2}},
		&Label{name: []string{VK_COMMAND_BUFFER, VK_RENDER_PASS, VK_SUBPASS}, id: []int{3, 1, 2}},
		&Label{name: []string{VK_COMMAND_BUFFER, VK_RENDER_PASS, VK_SUBPASS}, id: []int{3, 1, 2}},
		&Label{name: []string{VK_COMMAND_BUFFER, VK_RENDER_PASS, VK_SUBPASS}, id: []int{3, 1, 3}},
		&Label{name: []string{VK_COMMAND_BUFFER, VK_RENDER_PASS, VK_SUBPASS}, id: []int{3, 1, 3}},
		&Label{name: []string{VK_COMMAND_BUFFER, VK_RENDER_PASS}, id: []int{3, 1}},
		&Label{name: []string{VK_COMMAND_BUFFER}, id: []int{3}},
		&Label{name: []string{VK_COMMAND_BUFFER}, id: []int{3}},
		&Label{name: []string{VK_COMMAND_BUFFER}, id: []int{3}},
		&Label{},
		&Label{},
	}

	currentHierarchy := &Hierarchy{}
	obtainedLabels := []*Label{}
	for _, name := range commands {
		label := getLabelFromHierarchy(name, commandHierarchyNames, currentHierarchy)
		obtainedLabels = append(obtainedLabels, label)
	}
	for i := range wantedLabels {
		if reflect.DeepEqual(wantedLabels[i], obtainedLabels[i]) == false {
			t.Errorf("The label for command %s with id %d is different", commands[i], i)
			t.Errorf("Found %v wanted %v", obtainedLabels[i], wantedLabels[i])
		}
	}

}

func TestGetLabelFromHierarchy3(t *testing.T) {

	commandHierarchyNames, _ := getNameForCommandAndSubCommandHierarchy()

	auxiliarCommands := []string{
		"SetViewPort",
		"SetBarrier",
		"DrawIndex",
		"SetScissor",
		"SetIndexForDraw",
	}

	commands := []string{
		auxiliarCommands[0],
		auxiliarCommands[1],
		VK_BEGIN_COMMAND_BUFFER,
		auxiliarCommands[0],
		auxiliarCommands[1],
		VK_CMD_BEGIN_RENDER_PASS,
		auxiliarCommands[1],
		auxiliarCommands[2],
		auxiliarCommands[3],
		VK_CMD_NEXT_SUBPASS,
		auxiliarCommands[1],
		auxiliarCommands[2],
		VK_CMD_NEXT_SUBPASS,
		VK_CMD_NEXT_SUBPASS,
		auxiliarCommands[1],
		auxiliarCommands[2],
		VK_CMD_NEXT_SUBPASS,
		VK_CMD_END_RENDER_PASS,
		auxiliarCommands[0],
		VK_CMD_BEGIN_RENDER_PASS,
		VK_CMD_NEXT_SUBPASS,
		VK_CMD_NEXT_SUBPASS,
		VK_CMD_NEXT_SUBPASS,
		VK_CMD_NEXT_SUBPASS,
		VK_CMD_END_RENDER_PASS,
		auxiliarCommands[0],
		VK_CMD_BEGIN_RENDER_PASS,
		VK_CMD_END_RENDER_PASS,
		auxiliarCommands[0],
		VK_CMD_BEGIN_RENDER_PASS,
		VK_CMD_END_RENDER_PASS,
		auxiliarCommands[3],
		VK_CMD_BEGIN_RENDER_PASS,
		VK_CMD_END_RENDER_PASS,
		VK_CMD_BEGIN_RENDER_PASS,
		VK_CMD_END_RENDER_PASS,
		VK_END_COMMAND_BUFFER,
		auxiliarCommands[2],
		auxiliarCommands[1],
		VK_BEGIN_COMMAND_BUFFER,
		VK_END_COMMAND_BUFFER,
		auxiliarCommands[1],
	}
	wantedLabels := []*Label{
		&Label{},
		&Label{},
		&Label{name: []string{VK_COMMAND_BUFFER}, id: []int{1}},
		&Label{name: []string{VK_COMMAND_BUFFER}, id: []int{1}},
		&Label{name: []string{VK_COMMAND_BUFFER}, id: []int{1}},
		&Label{name: []string{VK_COMMAND_BUFFER, VK_RENDER_PASS}, id: []int{1, 1}},
		&Label{name: []string{VK_COMMAND_BUFFER, VK_RENDER_PASS, VK_SUBPASS}, id: []int{1, 1, 1}},
		&Label{name: []string{VK_COMMAND_BUFFER, VK_RENDER_PASS, VK_SUBPASS}, id: []int{1, 1, 1}},
		&Label{name: []string{VK_COMMAND_BUFFER, VK_RENDER_PASS, VK_SUBPASS}, id: []int{1, 1, 1}},
		&Label{name: []string{VK_COMMAND_BUFFER, VK_RENDER_PASS, VK_SUBPASS}, id: []int{1, 1, 2}},
		&Label{name: []string{VK_COMMAND_BUFFER, VK_RENDER_PASS, VK_SUBPASS}, id: []int{1, 1, 2}},
		&Label{name: []string{VK_COMMAND_BUFFER, VK_RENDER_PASS, VK_SUBPASS}, id: []int{1, 1, 2}},
		&Label{name: []string{VK_COMMAND_BUFFER, VK_RENDER_PASS, VK_SUBPASS}, id: []int{1, 1, 3}},
		&Label{name: []string{VK_COMMAND_BUFFER, VK_RENDER_PASS, VK_SUBPASS}, id: []int{1, 1, 4}},
		&Label{name: []string{VK_COMMAND_BUFFER, VK_RENDER_PASS, VK_SUBPASS}, id: []int{1, 1, 4}},
		&Label{name: []string{VK_COMMAND_BUFFER, VK_RENDER_PASS, VK_SUBPASS}, id: []int{1, 1, 4}},
		&Label{name: []string{VK_COMMAND_BUFFER, VK_RENDER_PASS, VK_SUBPASS}, id: []int{1, 1, 5}},
		&Label{name: []string{VK_COMMAND_BUFFER, VK_RENDER_PASS}, id: []int{1, 1}},
		&Label{name: []string{VK_COMMAND_BUFFER}, id: []int{1}},
		&Label{name: []string{VK_COMMAND_BUFFER, VK_RENDER_PASS}, id: []int{1, 2}},
		&Label{name: []string{VK_COMMAND_BUFFER, VK_RENDER_PASS, VK_SUBPASS}, id: []int{1, 2, 2}},
		&Label{name: []string{VK_COMMAND_BUFFER, VK_RENDER_PASS, VK_SUBPASS}, id: []int{1, 2, 3}},
		&Label{name: []string{VK_COMMAND_BUFFER, VK_RENDER_PASS, VK_SUBPASS}, id: []int{1, 2, 4}},
		&Label{name: []string{VK_COMMAND_BUFFER, VK_RENDER_PASS, VK_SUBPASS}, id: []int{1, 2, 5}},
		&Label{name: []string{VK_COMMAND_BUFFER, VK_RENDER_PASS}, id: []int{1, 2}},
		&Label{name: []string{VK_COMMAND_BUFFER}, id: []int{1}},
		&Label{name: []string{VK_COMMAND_BUFFER, VK_RENDER_PASS}, id: []int{1, 3}},
		&Label{name: []string{VK_COMMAND_BUFFER, VK_RENDER_PASS}, id: []int{1, 3}},
		&Label{name: []string{VK_COMMAND_BUFFER}, id: []int{1}},
		&Label{name: []string{VK_COMMAND_BUFFER, VK_RENDER_PASS}, id: []int{1, 4}},
		&Label{name: []string{VK_COMMAND_BUFFER, VK_RENDER_PASS}, id: []int{1, 4}},
		&Label{name: []string{VK_COMMAND_BUFFER}, id: []int{1}},
		&Label{name: []string{VK_COMMAND_BUFFER, VK_RENDER_PASS}, id: []int{1, 5}},
		&Label{name: []string{VK_COMMAND_BUFFER, VK_RENDER_PASS}, id: []int{1, 5}},
		&Label{name: []string{VK_COMMAND_BUFFER, VK_RENDER_PASS}, id: []int{1, 6}},
		&Label{name: []string{VK_COMMAND_BUFFER, VK_RENDER_PASS}, id: []int{1, 6}},
		&Label{name: []string{VK_COMMAND_BUFFER}, id: []int{1}},
		&Label{},
		&Label{},
		&Label{name: []string{VK_COMMAND_BUFFER}, id: []int{2}},
		&Label{name: []string{VK_COMMAND_BUFFER}, id: []int{2}},
		&Label{},
	}

	currentHierarchy := &Hierarchy{}
	obtainedLabels := []*Label{}
	for _, name := range commands {
		label := getLabelFromHierarchy(name, commandHierarchyNames, currentHierarchy)
		obtainedLabels = append(obtainedLabels, label)
	}
	for i := range wantedLabels {
		if reflect.DeepEqual(wantedLabels[i], obtainedLabels[i]) == false {
			t.Errorf("The label for command %s with id %d is different", commands[i], i)
			t.Errorf("Found %v wanted %v", obtainedLabels[i], wantedLabels[i])
		}
	}

}
