package graph_visualization

import (
	"bytes"
	"context"
	"fmt"
	"github.com/google/gapid/gapis/api"
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
	COMMAND_BEGIN_DEBUG_MARKER = "vkCmdDebugMarkerBeginEXT"
	COMMAND_END_DEBUG_MARKER   = "vkCmdDebugMarkerEndEXT"
	COMMAND_DEBUG_MARKER       = "vkCmdDebugMarker"
	COMMAND_BEGIN_RENDER_PASS  = "vkCmdBeginRenderPass"
	COMMAND_BUFFER             = "commandBuffer"
	EMPTY                      = "Empty"
	MAXIMUM_LEVEL_IN_HIERARCHY = 5
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
func splitLabelByChar(label *string, splitChar byte) []string {
	output := []string{}
	prevPos := 0
	for i := 0; i <= len(*label); i++ {
		if i == len(*label) || (*label)[i] == splitChar {
			output = append(output, (*label)[prevPos:i])
			prevPos = i + 1
		}
	}
	return output
}

func mergeSplitLabel(splitLabel []string) string {
	output := ""
	for _, val := range splitLabel {
		output += val
		output += "/"
	}
	return output
}

func getMaxCommonPrefix(splitLabel1 *[]string, splitLabel2 *[]string) int {
	size := len(*splitLabel1)
	if len(*splitLabel2) < size {
		size = len(*splitLabel2)
	}
	for i := 0; i < size; i++ {
		if (*splitLabel1)[i] != (*splitLabel2)[i] {
			return i
		}
	}
	return size
}

func getMaxCommonPrefixBetweenLabels(label1, label2 string) int {
	splitLabel1 := splitLabelByChar(&label1, '/')
	splitLabel2 := splitLabelByChar(&label2, '/')
	return getMaxCommonPrefix(&splitLabel1, &splitLabel2)
}

type Hierarchy struct {
	idLevels     [MAXIMUM_LEVEL_IN_HIERARCHY]int
	currentLevel int
}

func (h *Hierarchy) SetZeroFrom(from int) {
	for i := from; i < MAXIMUM_LEVEL_IN_HIERARCHY; i++ {
		if i >= 0 {
			h.idLevels[i] = 0
		}
	}
}

func addDebugMarkersToGraph(graph *Graph, idNodes []int) {
	posBeginDebugMarker := 0
	labelBeginDebugMarker := []string{}
	for pos, idNode := range idNodes {
		node := graph.idNodeToNode[idNode]
		if node.name == COMMAND_BEGIN_DEBUG_MARKER {
			posBeginDebugMarker = pos
			labelBeginDebugMarker = splitLabelByChar(&node.label, '/')
		} else {
			if node.name == COMMAND_END_DEBUG_MARKER {
				labelEndDebugMarker := splitLabelByChar(&node.label, '/')
				if getMaxCommonPrefix(&labelBeginDebugMarker, &labelEndDebugMarker) == len(labelEndDebugMarker)-1 {
					for i := posBeginDebugMarker; i <= pos; i++ {
						node = graph.idNodeToNode[idNodes[i]]
						splitLabel := splitLabelByChar(&node.label, '/')
						lastLabel := splitLabel[len(splitLabel)-1]
						splitLabel[len(splitLabel)-1] = COMMAND_DEBUG_MARKER
						node.label = mergeSplitLabel(splitLabel) + lastLabel
					}
				}
			}
		}
	}
}

func getCommandBuffer(command api.Cmd) string {
	parameters := command.CmdParams()
	for _, v := range parameters {
		if v.Name == COMMAND_BUFFER {
			commandBuffer := v.Name + fmt.Sprintf("%d", v.Get())
			return commandBuffer
		}
	}
	return ""
}
func getCommandLabel(command api.Cmd, idCommandNode uint64, commandHierarchyNames *HierarchyNames, labelToHierarchy *map[string]*Hierarchy) string {
	commmandName := command.CmdName()
	var label bytes.Buffer
	if commandBuffer := getCommandBuffer(command); commandBuffer != "" {
		if _, ok := (*labelToHierarchy)[commandBuffer]; ok == false {
			(*labelToHierarchy)[commandBuffer] = &Hierarchy{}
		}
		currentHierarchy := (*labelToHierarchy)[commandBuffer]
		label.WriteString(commandBuffer + "/")
		label.WriteString(getLabelFromHierarchy(commmandName, commandHierarchyNames, currentHierarchy))
		label.WriteString(fmt.Sprintf("%d_%s", idCommandNode, commmandName))
	} else {
		label.WriteString(fmt.Sprintf("%d_%s", idCommandNode, commmandName))
	}
	return label.String()
}

func getLabelFromHierarchy(name string, hierarchyNames *HierarchyNames, currentHierarchy *Hierarchy) string {
	if currentLevel, ok := hierarchyNames.beginNames[name]; ok && currentLevel <= currentHierarchy.currentLevel {
		currentHierarchy.idLevels[currentLevel]++
		currentHierarchy.currentLevel = currentLevel + 1
	} else {
		if currentLevel, ok := hierarchyNames.endNames[name]; ok && currentLevel <= currentHierarchy.currentLevel {
			currentHierarchy.currentLevel = currentLevel + 1
		}
	}
	var label bytes.Buffer
	for i := 0; i < currentHierarchy.currentLevel; i++ {
		label.WriteString(fmt.Sprintf("%d_%s/", currentHierarchy.idLevels[i], hierarchyNames.listNames[i]))
	}
	if _, ok := hierarchyNames.beginNames[name]; ok {
		if name == COMMAND_BEGIN_RENDER_PASS {
			currentHierarchy.idLevels[currentHierarchy.currentLevel]++
			currentHierarchy.currentLevel++
		}
	}
	currentHierarchy.SetZeroFrom(currentHierarchy.currentLevel)
	return label.String()
}

func makeChainByRenderScope(graph *Graph, prevNode *Node, currNode *Node) {
	if _, ok1 := commandsInsideRenderScope[prevNode.name]; ok1 {
		if _, ok2 := commandsInsideRenderScope[currNode.name]; ok2 {
			if getMaxCommonPrefixBetweenLabels(prevNode.label, currNode.label) >= 2 {
				graph.addEdgeByNodes(currNode, prevNode)
			}
		}
	}
}

func getSubCommandLabel(commandNode dependencygraph2.CmdNode, commmandName string, subCommandToLabel *map[string]string) (string, string) {
	subCommandName := commmandName
	label := commmandName
	for i := 1; i < len(commandNode.Index); i++ {
		subCommandName += fmt.Sprintf("/%d", commandNode.Index[i])
		if i+1 < len(commandNode.Index) {
			if name, ok := (*subCommandToLabel)[subCommandName]; ok {
				label += "/" + name
			} else {
				label += fmt.Sprintf("/%d", commandNode.Index[i])
			}
		}
	}
	return label, subCommandName
}

func createGraphFromDependencyGraph(dependencyGraph dependencygraph2.DependencyGraph) (*Graph, error) {

	graph := createGraph(0)
	numberNodes := dependencyGraph.NumNodes()
	commandHierarchyNames, subCommandHierarchyNames := getNameForCommandAndSubCommandHierarchy()
	subCommandToLabel := map[string]string{}
	labelToHierarchy := map[string]*Hierarchy{}
	prevNode := &Node{}
	validIdNodes := []int{}

	for i := 0; i < numberNodes; i++ {
		dependencyNode := dependencyGraph.GetNode(dependencygraph2.NodeID(i))
		if commandNode, ok := dependencyNode.(dependencygraph2.CmdNode); ok {
			idCommandNode := commandNode.Index[0]
			command := dependencyGraph.GetCommand(api.CmdID(idCommandNode))
			var name, label bytes.Buffer

			isReal := api.CmdID(idCommandNode).IsReal()
			if len(commandNode.Index) == 1 {
				label.WriteString(getCommandLabel(command, idCommandNode, commandHierarchyNames, &labelToHierarchy))
				name.WriteString(command.CmdName())
			} else {
				if len(commandNode.Index) > 1 {
					commmandName := fmt.Sprintf("%d_%s", idCommandNode, command.CmdName())
					subCommandLabel, tmpLabel := getSubCommandLabel(commandNode, commmandName, &subCommandToLabel)
					if _, ok := labelToHierarchy[subCommandLabel]; ok == false {
						labelToHierarchy[subCommandLabel] = &Hierarchy{}
					}

					currentHierarchy := labelToHierarchy[subCommandLabel]
					dependencyNodeAccesses := dependencyGraph.GetNodeAccesses(dependencygraph2.NodeID(i))
					subCommandName := EMPTY
					if len(dependencyNodeAccesses.InitCmdNodes) > 0 {
						subCommandName = graph.idNodeToNode[int(dependencyNodeAccesses.InitCmdNodes[0])].name
					}
					hierarchyLabel := getLabelFromHierarchy(subCommandName, subCommandHierarchyNames, currentHierarchy)
					hierarchyLabel += fmt.Sprintf("%d_", commandNode.Index[len(commandNode.Index)-1]) + subCommandName
					subCommandToLabel[tmpLabel] = hierarchyLabel

					label.WriteString(subCommandLabel + "/")
					label.WriteString(hierarchyLabel)
					name.WriteString(subCommandName)
				}
			}
			attributes := []string{}
			for _, property := range command.CmdParams() {
				attributes = append(attributes, property.Name+fmt.Sprintf("%d", property.Get()))
			}
			graph.addNodeByIdAndNameAndAttrAndIsReal(i, label.String(), name.String(), attributes, isReal)
			validIdNodes = append(validIdNodes, i)

			currNode := graph.idNodeToNode[i]
			makeChainByRenderScope(graph, prevNode, currNode)
			if _, ok := commandsInsideRenderScope[currNode.name]; ok {
				prevNode = currNode
			}
		}
	}

	addDebugMarkersToGraph(graph, validIdNodes)

	addDependencyToGraph := func(source, sink dependencygraph2.NodeID) error {
		idSource, idSink := int(source), int(sink)
		if sourceNode, ok1 := graph.idNodeToNode[idSource]; ok1 {
			if sinkNode, ok2 := graph.idNodeToNode[idSink]; ok2 {
				graph.addEdgeByNodes(sourceNode, sinkNode)
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
	dependencyGraph, err := dependencygraph2.GetDependencyGraph(
		ctx, p, config)
	if err != nil {
		return "", err
	}

	graph, err := createGraphFromDependencyGraph(dependencyGraph)
	graph.joinNodesByFrame()
	file := ""
	if format == "pbtxt" {
		file = graph.getPbtxtFileFaster()
	} else if format == "proto" {
		file = graph.getProtoFile()
	} else if format == "dot" {
		file = graph.getDotFile()
	} else {
		return "", fmt.Errorf("Invalid format (Format supported: proto,pbtxt and dot)")
	}

	return file, err
}
