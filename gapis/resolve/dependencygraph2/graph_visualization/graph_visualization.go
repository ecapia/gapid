package graph_visualization

import (
	"context"
	"fmt"
	"github.com/google/gapid/core/math/interval"
	"github.com/google/gapid/gapis/api"
	"github.com/google/gapid/gapis/capture"
	"github.com/google/gapid/gapis/resolve/dependencygraph2"
)

var (
	beginCommands = map[string]int{
		"vkBeginCommandBuffer": 0,
		"vkCmdBeginRenderPass": 1,
		"vkCmdNextSubpass":     2,
	}
	listBeginCommands = []string{
		"vkBeginCommandBuffer",
		"vkCmdBeginRenderPass",
		"vkCmdNextSubpass",
	}
	listNameCommands = []string {
		"vkCommandBuffer",
		"vkRenderPass",
		"vkSubpass",
	}
	endCommands = map[string]int{
		"vkEndCommandBuffer": 0,
		"vkCmdEndRenderPass": 1,
		"vkCmdNextSubpass":   2,
	}
	commandsInsideRenderScope = map[string]struct{}{
		"vkCmdDrawIndexed": struct{}{},
		"vkCmdNextSubpass": struct{}{},
		"vkCmdDraw":        struct{}{},
	}
)
const (
	COMMAND_BEGIN_RENDER_PASS = "vkCmdBeginRenderPass"
	COMMAND_BUFFER            = "commandBuffer"
	MAXIMUM_LIMIT_IN_HIERARCHY = 5
)

type Hierarchy struct {
	idLevels     [MAXIMUM_LIMIT_IN_HIERARCHY]int
	currentId    int
	currentLevel int
}

func getCommandBuffer(command api.Cmd) string {
	parameters := command.CmdParams()
	for _, v := range parameters {
		if v.Name == COMMAND_BUFFER {
			commandBuffer := v.Name + fmt.Sprintf("%d", v.Get()) + "/"
			return commandBuffer
		}
	}
	return ""
}

func getCommandLabel(currentHierarchy *Hierarchy, command api.Cmd) string {
	commandName := command.CmdName()
	isEndCommand := false
	if currentLevel, ok := beginCommands[commandName]; ok && currentLevel <= currentHierarchy.currentLevel {
		currentHierarchy.idLevels[currentLevel] = currentHierarchy.currentId
		currentHierarchy.currentId++
		currentHierarchy.currentLevel = currentLevel + 1
	} else {
		if currentLevel, ok := endCommands[commandName]; ok && currentLevel <= currentHierarchy.currentLevel {
			currentHierarchy.currentLevel = currentLevel + 1
			isEndCommand = true
		}
	}

	label := "\""
	for i := 0; i < currentHierarchy.currentLevel; i++ {
		if i == 0 {
			label += getCommandBuffer(command)
		} else {
			label += fmt.Sprintf("%s%d/", listNameCommands[i], currentHierarchy.idLevels[i])
		}
	}
	if isEndCommand {
		currentHierarchy.currentLevel--
	} else {
		if _, ok := beginCommands[commandName]; ok {
			if commandName == COMMAND_BEGIN_RENDER_PASS {
				currentHierarchy.idLevels[currentHierarchy.currentLevel] = currentHierarchy.currentId
				currentHierarchy.currentId++
				currentHierarchy.currentLevel++
			}
		}
	}
	return label
}

func getSubCommandLabel(cmdNode dependencygraph2.CmdNode) string {
	label := ""
	for i := 1; i < len(cmdNode.Index); i++ {
		label += fmt.Sprintf("/%d", cmdNode.Index[i])
	}
	return label
}

func splitLabelByChar(label *string, splitChar byte) []string {
	splitLabel := []string{}
	prevPos := 0
	for i := 0; i <= len(*label); i++ {
		if i == len(*label) || (*label)[i] == splitChar {
			splitLabel = append(splitLabel, (*label)[prevPos:i])
			prevPos = i + 1
		}
	}
	return splitLabel
}
func getMaxCommonPrefixBetweenSplitLabels(splitLabel1 *[]string, splitLabel2 *[]string) int {
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
	return getMaxCommonPrefixBetweenSplitLabels(&splitLabel1, &splitLabel2)
}

func createGraphFromDependencyGraph(dependencyGraph dependencygraph2.DependencyGraph) (*Graph, error) {

	numberNodes := dependencyGraph.NumNodes()
	graph := createGraph(0)
	currentHierarchy := &Hierarchy{}
	prevNode := &Node{}
	for i := 0; i < numberNodes; i++ {
		dependencyNode := dependencyGraph.GetNode(dependencygraph2.NodeID(i))
		if cmdNode, ok := dependencyNode.(dependencygraph2.CmdNode); ok {
			idCmdNode := cmdNode.Index[0]
			command := dependencyGraph.GetCommand(api.CmdID(idCmdNode))
			commandName := command.CmdName()
			label := getCommandLabel(currentHierarchy, command)
			label += fmt.Sprintf("%s%d", commandName, idCmdNode)
			label += getSubCommandLabel(cmdNode)
			label += "\""
			attr := fmt.Sprintf("\"%v\"", command.CmdParams())

			graph.addNodeByIdAndCommandNameAndAttr(i, label, commandName, attr)

			node := graph.idNodeToNode[i]
			if _, ok1 := commandsInsideRenderScope[prevNode.commandName]; ok1 {
				if _, ok2 := commandsInsideRenderScope[node.commandName]; ok2 {
					if getMaxCommonPrefixBetweenLabels(prevNode.label, node.label) >= 2 {
						graph.addEdgeByNode(node, prevNode)
					}
				}
			}
			if _, ok := commandsInsideRenderScope[node.commandName]; ok {
				prevNode = node
			}
		}
	}

	addDependencyInGraph := func(source, sink dependencygraph2.NodeID) error {
		idSource, idSink := int(source), int(sink)
		if sourceNode, ok1 := graph.idNodeToNode[idSource]; ok1 {
			if sinkNode, ok2 := graph.idNodeToNode[idSink]; ok2 {
				_, ok1 = commandsInsideRenderScope[sourceNode.commandName]
				_, ok2 = commandsInsideRenderScope[sinkNode.commandName]
				if ok1 == false || ok2 == false {
					graph.addEdgeByIdNodes(idSource, idSink)
				}
			}
		}
		return nil
	}

	err := dependencyGraph.ForeachDependency(addDependencyInGraph)
	return graph, err
}

func GetGraphVisualizationFileFromCapture(ctx context.Context, p *capture.Capture) (string, error) {
	config := dependencygraph2.DependencyGraphConfig{}
	dependencyGraph, err := dependencygraph2.BuildDependencyGraph(ctx, config, p, []api.Cmd{}, interval.U64RangeList{})
	if err != nil {
		return "", err
	}

	graph, err := createGraphFromDependencyGraph(dependencyGraph)
	graph.removeNodesWithZeroDegree()
	file := graph.getPbtxtFile()
	return file, err
}
