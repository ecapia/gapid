package graph_visualization

import (
	"bytes"
	"context"
	"fmt"
	"github.com/google/gapid/gapis/api"
	"github.com/google/gapid/gapis/api/vulkan"
	"github.com/google/gapid/gapis/resolve/dependencygraph2"
	"github.com/google/gapid/gapis/service/path"
)

var (
	commandsInsideRenderScope = map[string]struct{}{
		"vkCmdDrawIndexed": struct{}{},
		"vkCmdNextSubpass": struct{}{},
		"vkCmdDraw":        struct{}{},
	}
)

const (
	COMMAND_QUEUE_SUBMIT       = "vkQueueSubmit"
	COMMAND_BEGIN_DEBUG_MARKER = "vkCmdDebugMarkerBeginEXT"
	COMMAND_END_DEBUG_MARKER   = "vkCmdDebugMarkerEndEXT"
	COMMAND_DEBUG_MARKER       = "vkCmdDebugMarker"
	COMMAND_BEGIN_RENDER_PASS  = "vkCmdBeginRenderPass"
	COMMAND_BUFFER             = "commandBuffer"
	EMPTY_NAME                 = "Empty"
	MAXIMUM_LEVEL_IN_HIERARCHY = 10
)

type HierarchyNames struct {
	beginNames     map[string]int
	endNames       map[string]int
	listBeginNames []string
	listNames      []string
}

func (currentHierarchy *HierarchyNames) add(beginName, endName, name string) {
	size := len(currentHierarchy.listNames)
	currentHierarchy.beginNames[beginName] = size
	currentHierarchy.endNames[endName] = size
	currentHierarchy.listBeginNames = append(currentHierarchy.listBeginNames, beginName)
	currentHierarchy.listNames = append(currentHierarchy.listNames, name)
}

func getNameForCommandAndSubCommandHierarchy() (*HierarchyNames, *HierarchyNames) {
	commandHierarchyNames := &HierarchyNames{beginNames: map[string]int{}, endNames: map[string]int{},
		listBeginNames: []string{}, listNames: []string{}}
	commandHierarchyNames.add("vkBeginCommandBuffer", "vkEndCommandBuffer", "vkCommandBuffer")
	commandHierarchyNames.add("vkCmdBeginRenderPass", "vkCmdEndRenderPass", "vkRenderPass")
	commandHierarchyNames.add("vkCmdNextSubpass", "vkCmdNextSubpass", "vkSubpass")

	subCommandHierarchyNames := &HierarchyNames{beginNames: map[string]int{}, endNames: map[string]int{},
		listBeginNames: []string{}, listNames: []string{}}
	subCommandHierarchyNames.add("vkCmdBeginRenderPass", "vkCmdEndRenderPass", "vkRenderPass")
	subCommandHierarchyNames.add("vkCmdNextSubpass", "vkCmdNextSubpass", "vkSubpass")
	return commandHierarchyNames, subCommandHierarchyNames
}

type Hierarchy struct {
	idLevels              [MAXIMUM_LEVEL_IN_HIERARCHY]int
	currentLevel          int
	numberOfCommandBuffer int
}

func (h *Hierarchy) setZerosFrom(from int) {
	for i := from; i < MAXIMUM_LEVEL_IN_HIERARCHY; i++ {
		if i >= 0 {
			h.idLevels[i] = 0
		}
	}
}
func addDebugMarkersToNodes(nodes []*Node) {
	posBeginDebugMarker := -1
	labelBeginDebugMarker := &Label{}
	for pos, node := range nodes {
		if node.name == COMMAND_BEGIN_DEBUG_MARKER {
			posBeginDebugMarker = pos
			labelBeginDebugMarker = node.label
		} else {
			if posBeginDebugMarker >= 0 && node.name == COMMAND_END_DEBUG_MARKER {
				labelEndDebugMarker := node.label
				if len(labelBeginDebugMarker.id) == len(labelEndDebugMarker.id) &&
					getMaxCommonPrefix(labelBeginDebugMarker, labelEndDebugMarker) == len(labelEndDebugMarker.id)-1 {
					for i := posBeginDebugMarker; i <= pos; i++ {
						nodeInDebugMarker := nodes[i]
						size := len(nodeInDebugMarker.label.id)
						nodeInDebugMarker.label.pushBack(nodeInDebugMarker.label.name[size-1], nodeInDebugMarker.label.id[size-1])
						nodeInDebugMarker.label.update(size-1, COMMAND_DEBUG_MARKER, 0)
					}
				}
				posBeginDebugMarker = -1
			}
		}
	}
}
func getCommandBuffer(command api.Cmd) (string, vulkan.VkCommandBuffer) {
	parameters := command.CmdParams()
	for _, parameter := range parameters {
		if parameter.Name == COMMAND_BUFFER {
			id := parameter.Get().(vulkan.VkCommandBuffer)
			return parameter.Name, id
		}
	}
	return "", 0
}

func getCommandLabel(command api.Cmd, idCommandNode uint64, commandHierarchyNames *HierarchyNames, commandBufferIdToHierarchy *map[vulkan.VkCommandBuffer]*Hierarchy) *Label {
	label := &Label{}
	commandName := command.CmdName()
	if command, idCommandBuffer := getCommandBuffer(command); command != "" {
		if _, ok := (*commandBufferIdToHierarchy)[idCommandBuffer]; ok == false {
			size := len(*commandBufferIdToHierarchy)
			(*commandBufferIdToHierarchy)[idCommandBuffer] = &Hierarchy{numberOfCommandBuffer: size + 1}
		}
		currentHierarchy := (*commandBufferIdToHierarchy)[idCommandBuffer]
		label.pushBack(command, currentHierarchy.numberOfCommandBuffer)
		label.pushBackLabel(getLabelFromHierarchy(commandName, commandHierarchyNames, currentHierarchy))
		label.pushBack(commandName, int(idCommandNode))
	} else {
		label.pushBack(commandName, int(idCommandNode))
	}
	return label
}

func getLabelFromHierarchy(name string, hierarchyNames *HierarchyNames, currentHierarchy *Hierarchy) *Label {
	label := &Label{}
	if currentLevel, ok := hierarchyNames.beginNames[name]; ok && currentLevel <= currentHierarchy.currentLevel {
		currentHierarchy.idLevels[currentLevel]++
		currentHierarchy.currentLevel = currentLevel + 1
	} else {
		if currentLevel, ok := hierarchyNames.endNames[name]; ok && currentLevel <= currentHierarchy.currentLevel {
			currentHierarchy.currentLevel = currentLevel + 1
		}
	}
	for i := 0; i < currentHierarchy.currentLevel; i++ {
		label.pushBack(hierarchyNames.listNames[i], currentHierarchy.idLevels[i])
	}
	if _, ok := hierarchyNames.beginNames[name]; ok {
		if name == COMMAND_BEGIN_RENDER_PASS {
			currentHierarchy.idLevels[currentHierarchy.currentLevel]++
			currentHierarchy.currentLevel++
		}
	}
	currentHierarchy.setZerosFrom(currentHierarchy.currentLevel)
	return label
}

func getSubCommandLabel(commandNode dependencygraph2.CmdNode, commandName string,
	subCommandNameToLabel *map[string]*Label) (*Label, string) {

	currentSubCommandName := commandName + fmt.Sprintf("_%d", commandNode.Index[0])
	subCommandName := currentSubCommandName
	label := &Label{}
	label.pushBack(commandName, int(commandNode.Index[0]))
	for curr := 1; curr < len(commandNode.Index); curr++ {
		currentSubCommandName += fmt.Sprintf("/_%d", commandNode.Index[curr])
		if curr+1 < len(commandNode.Index) {
			if labelFromHierarchy, ok := (*subCommandNameToLabel)[currentSubCommandName]; ok {
				label.pushBackLabel(labelFromHierarchy)
			} else {
				label.pushBack("", int(commandNode.Index[curr]))
			}
			subCommandName += fmt.Sprintf("/_%d", commandNode.Index[curr])
		}
	}
	return label, subCommandName
}

func createGraphFromDependencyGraph(dependencyGraph dependencygraph2.DependencyGraph) (*Graph, error) {

	graph := createGraph(0)
	numberNodes := dependencyGraph.NumNodes()
	commandHierarchyNames, subCommandHierarchyNames := getNameForCommandAndSubCommandHierarchy()
	subCommandNameToLabel := map[string]*Label{}
	subCommandNameToHierarchy := map[string]*Hierarchy{}
	commandBufferIdToHierarchy := map[vulkan.VkCommandBuffer]*Hierarchy{}
	validNodes := []*Node{}
	for i := 0; i < numberNodes; i++ {
		dependencyNode := dependencyGraph.GetNode(dependencygraph2.NodeID(i))
		if commandNode, ok := dependencyNode.(dependencygraph2.CmdNode); ok {
			idCommandNode := commandNode.Index[0]
			command := dependencyGraph.GetCommand(api.CmdID(idCommandNode))
			commandName := command.CmdName()
			var name bytes.Buffer
			label := &Label{}
			isReal := api.CmdID(idCommandNode).IsReal()

			if len(commandNode.Index) == 1 {
				label = getCommandLabel(command, idCommandNode, commandHierarchyNames, &commandBufferIdToHierarchy)
				name.WriteString(commandName)
			} else {
				if len(commandNode.Index) > 1 {
					subCommandName := ""
					label, subCommandName = getSubCommandLabel(commandNode, commandName, &subCommandNameToLabel)
					if _, ok := subCommandNameToHierarchy[subCommandName]; ok == false {
						subCommandNameToHierarchy[subCommandName] = &Hierarchy{}
					}

					currentHierarchy := subCommandNameToHierarchy[subCommandName]
					subCommandName += fmt.Sprintf("/_%d", commandNode.Index[len(commandNode.Index)-1])

					nameLastLevel := EMPTY_NAME
					dependencyNodeAccesses := dependencyGraph.GetNodeAccesses(dependencygraph2.NodeID(i))
					if len(dependencyNodeAccesses.InitCmdNodes) > 0 {
						nameLastLevel = graph.idNodeToNode[int(dependencyNodeAccesses.InitCmdNodes[0])].name
					}

					labelFromHierarchy := getLabelFromHierarchy(nameLastLevel, subCommandHierarchyNames, currentHierarchy)
					labelFromHierarchy.pushBack(nameLastLevel, int(commandNode.Index[len(commandNode.Index)-1]))

					subCommandNameToLabel[subCommandName] = labelFromHierarchy
					label.pushBackLabel(labelFromHierarchy)
					name.WriteString(nameLastLevel)
				}
			}
			attributes := []string{}
			for _, parameter := range command.CmdParams() {
				attributes = append(attributes, parameter.Name+fmt.Sprintf("%d", parameter.Get()))
			}

			graph.addNodeByIdAndNameAndAttrAndIsReal(i, label, name.String(), attributes, isReal)
			validNodes = append(validNodes, graph.idNodeToNode[i])
		}
	}

	addDebugMarkersToNodes(validNodes)

	addDependencyToGraph := func(source, sink dependencygraph2.NodeID) error {
		idSource, idSink := int(source), int(sink)
		if sourceNode, ok := graph.idNodeToNode[idSource]; ok {
			if sinkNode, ok := graph.idNodeToNode[idSink]; ok {
				graph.addEdgeByNodes(sourceNode, sinkNode)

				if sourceNode.name == COMMAND_BEGIN_DEBUG_MARKER || sourceNode.name == COMMAND_END_DEBUG_MARKER {
					if sinkNode.name == COMMAND_QUEUE_SUBMIT {
						graph.addEdgeByNodes(sinkNode, sourceNode)
					}
				}
			}
		}
		return nil
	}

	err := dependencyGraph.ForeachDependency(addDependencyToGraph)
	return graph, err
}

func GetGraphVisualizationFileFromCapture(ctx context.Context, p *path.Capture, format string) (string, error) {
	config := dependencygraph2.DependencyGraphConfig{
		SaveNodeAccesses:       true,
		IncludeInitialCommands: true,
	}
	dependencyGraph, err := dependencygraph2.GetDependencyGraph(ctx, p, config)
	if err != nil {
		return "", err
	}

	graph, err := createGraphFromDependencyGraph(dependencyGraph)
	graph.makeAdjacentList()
	graph.sortAdjacentList()
	graph.assignColors()
	graph.makeChunksByFrame()
	fmt.Println("Number of nodes in the top level = ", graph.getNumberNodesInTopLevel())
	fmt.Println("Total number of nodes and edges = ", graph.numberNodes, graph.numberEdges)

	output := ""
	if format == "pbtxt" {
		output = graph.getPbtxtFile()
	} else if format == "proto" {
		output = graph.getProtoFile()
	} else if format == "dot" {
		output = graph.getDotFile()
	} else {
		return "", fmt.Errorf("Invalid format (Format supported: proto,pbtxt and dot)")
	}

	return output, err
}
