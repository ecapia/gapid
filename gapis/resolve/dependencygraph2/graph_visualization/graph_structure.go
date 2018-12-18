package graph_visualization

import (
	"bytes"
	"context"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/google/gapid/core/data/protoconv"
	"github.com/google/gapid/gapis/resolve/dependencygraph2"
	"github.com/google/gapid/gapis/resolve/dependencygraph2/graph_visualization/protobuf"
	"sort"
)

const (
	QUEUE_PRESENT        = "vkQueuePresentKHR"
	SUPER_COMMAND_BUFFER = "SuperBuffer"
	SUPER_QUEUE_SUBMIT   = "SuperQueue"
	UNUSED_COMMANDS      = "UnusedCommands"
	MAX_NODES_BY_LEVEL   = 5
)

type Node struct {
	incomingIdNodesToIdEdge  map[int]int
	outcomingIdNodesToIdEdge map[int]int
	id                       int
	label                    string
	name                     string
	nameFrame                string
	attributes               []string
	isReal                   bool
}

type Edge struct {
	source, sink *Node
	id           int
	label        string
}

type Graph struct {
	idNodeToNode map[int]*Node
	idEdgeToEdge map[int]*Edge
	maxIdNode    int
	maxIdEdge    int
	numberNodes  int
	numberEdges  int
}

func createGraph(numberNodes int) *Graph {
	newGraph := &Graph{idNodeToNode: map[int]*Node{}, idEdgeToEdge: map[int]*Edge{}}
	for i := 0; i < numberNodes; i++ {
		newNode := &Node{incomingIdNodesToIdEdge: map[int]int{}, outcomingIdNodesToIdEdge: map[int]int{}, id: newGraph.maxIdNode}
		newGraph.idNodeToNode[newNode.id] = newNode
		newGraph.numberNodes++
		newGraph.maxIdNode++
	}
	return newGraph
}

func (g *Graph) addNodeByDefault(label string) int {
	id := g.maxIdNode
	newNode := &Node{incomingIdNodesToIdEdge: map[int]int{}, outcomingIdNodesToIdEdge: map[int]int{}, id: id, label: label}
	g.idNodeToNode[id] = newNode
	g.numberNodes++
	g.maxIdNode++
	return id
}

func (g *Graph) addNodeById(id int, label string) bool {
	_, ok := g.idNodeToNode[id]
	if ok == true {
		return false
	}

	newNode := &Node{incomingIdNodesToIdEdge: map[int]int{}, outcomingIdNodesToIdEdge: map[int]int{}, id: id, label: label}
	g.idNodeToNode[id] = newNode
	g.numberNodes++
	g.maxIdNode++
	if g.maxIdNode <= id {
		g.maxIdNode = id + 1
	}
	return true
}

func (g *Graph) addNodeByIdAndNameAndAttrAndIsReal(id int, label string, name string, attributes []string, isReal bool) bool {
	_, ok := g.idNodeToNode[id]
	if ok == true {
		return false
	}

	newNode := &Node{incomingIdNodesToIdEdge: map[int]int{}, outcomingIdNodesToIdEdge: map[int]int{}, id: id, label: label,
		name: name, attributes: attributes, isReal: isReal}
	g.idNodeToNode[id] = newNode
	g.numberNodes++
	g.maxIdNode++
	if g.maxIdNode <= id {
		g.maxIdNode = id + 1
	}
	return true
}

func (g *Graph) addEdge(newEdge *Edge) bool {
	source, sink := newEdge.source, newEdge.sink
	if _, ok := source.outcomingIdNodesToIdEdge[sink.id]; ok {
		return false
	}

	id := g.maxIdEdge
	g.idEdgeToEdge[id] = newEdge
	g.numberEdges++
	g.maxIdEdge++

	source.outcomingIdNodesToIdEdge[sink.id] = id
	sink.incomingIdNodesToIdEdge[source.id] = id
	return true
}

func (g *Graph) addEdgeByNodes(source, sink *Node) {
	id := g.maxIdEdge
	newEdge := &Edge{source: source, sink: sink, id: id}
	g.addEdge(newEdge)
}

func (g *Graph) addEdgeByIdNodes(idSource, idSink int) (int, bool) {
	source, ok := g.idNodeToNode[idSource]
	if ok == false {
		return 0, false
	}
	sink, ok := g.idNodeToNode[idSink]
	if ok == false {
		return 0, false
	}
	id := g.maxIdEdge
	newEdge := &Edge{source: source, sink: sink, id: id}
	g.addEdge(newEdge)
	return id, true
}

func (g *Graph) removeEdgeById(id int) bool {
	edge := g.idEdgeToEdge[id]
	source, sink := edge.source, edge.sink
	delete(source.outcomingIdNodesToIdEdge, sink.id)
	delete(sink.incomingIdNodesToIdEdge, source.id)

	delete(g.idEdgeToEdge, id)
	g.numberEdges--
	return true
}

func (g *Graph) removeNodeById(id int) bool {

	node, ok := g.idNodeToNode[id]
	if ok == false {
		return false
	}

	in, out := node.incomingIdNodesToIdEdge, node.outcomingIdNodesToIdEdge
	for _, value := range in {
		g.removeEdgeById(value)
	}
	for _, value := range out {
		g.removeEdgeById(value)
	}
	delete(g.idNodeToNode, id)
	g.numberNodes--
	return true
}

func (g *Graph) removeNodesWithZeroDegree() {

	for id, node := range g.idNodeToNode {
		if (len(node.incomingIdNodesToIdEdge) + len(node.outcomingIdNodesToIdEdge)) == 0 {
			g.removeNodeById(id)
		}
	}
}

func (g *Graph) joinEdgesThroughtNode(idNode int) bool {

	node, ok := g.idNodeToNode[idNode]
	if ok == false {
		return false
	}
	for idSource := range node.incomingIdNodesToIdEdge {
		for idSink := range node.outcomingIdNodesToIdEdge {
			g.addEdgeByIdNodes(idSource, idSink)
		}
	}
	return true
}

func (g *Graph) removeNodeByIdKeepingEdges(id int) bool {
	if g.joinEdgesThroughtNode(id) == false {
		return false
	}
	if g.removeNodeById(id) == false {
		return false
	}
	return true
}
func getKeysSortedFromMap(input map[int]int) []int {
	sortedKeys := []int{}
	for key := range input {
		sortedKeys = append(sortedKeys, key)
	}
	sort.Ints(sortedKeys)
	return sortedKeys
}

func (g *Graph) dfs(node *Node, visited *[]bool, numberFrame *int, queueSubmitToPos, commandBufferToPos *map[string]int, visitedIdNodes *[]*Node) {
	*visitedIdNodes = append(*visitedIdNodes, node)
	nameTopLevel := splitLabelByChar(&node.label, '/')[0]
	if node.label[0] == 'c' { //Needs to improve this part
		size := len(*commandBufferToPos)
		if _, ok := (*commandBufferToPos)[nameTopLevel]; ok == false {
			(*commandBufferToPos)[nameTopLevel] = size
		}
	} else {
		size := len(*queueSubmitToPos)
		if _, ok := (*queueSubmitToPos)[nameTopLevel]; ok == false {
			(*queueSubmitToPos)[nameTopLevel] = size
		}
	}

	(*visited)[node.id] = true
	node.nameFrame = QUEUE_PRESENT + fmt.Sprintf("%d", *numberFrame)
	idNeighbors := getKeysSortedFromMap(node.outcomingIdNodesToIdEdge)
	for _, idNeighbor := range idNeighbors {
		neighbor := g.idNodeToNode[idNeighbor]
		if (*visited)[neighbor.id] == false {
			g.dfs(neighbor, visited, numberFrame, queueSubmitToPos, commandBufferToPos, visitedIdNodes)
		}
	}
}

func (g *Graph) joinNodesByFrame() {
	visited := make([]bool, g.maxIdNode)
	numberFrame := 1
	for i := 0; i < g.maxIdNode; i++ {
		if node, ok := g.idNodeToNode[i]; ok && node.name == QUEUE_PRESENT && visited[node.id] == false {
			commandBufferToPos := map[string]int{}
			queueSubmitToPos := map[string]int{}
			visitedIdNodes := []*Node{}
			g.dfs(node, &visited, &numberFrame, &queueSubmitToPos, &commandBufferToPos, &visitedIdNodes)

			sizeSuperCommandBuffer := (len(commandBufferToPos) + MAX_NODES_BY_LEVEL - 1) / MAX_NODES_BY_LEVEL
			sizeSuperQueueSubmit := (len(queueSubmitToPos) + MAX_NODES_BY_LEVEL - 1) / MAX_NODES_BY_LEVEL
			for _, node := range visitedIdNodes {
				nameTopLevel := splitLabelByChar(&node.label, '/')[0]
				pos := 0
				if node.label[0] == 'c' { // needs to improve this part
					pos = commandBufferToPos[nameTopLevel]
					node.label = SUPER_COMMAND_BUFFER + fmt.Sprintf("_%d/", pos/sizeSuperCommandBuffer) + node.label
				} else {
					pos = queueSubmitToPos[nameTopLevel]
					node.label = SUPER_QUEUE_SUBMIT + fmt.Sprintf("_%d/", pos/sizeSuperQueueSubmit) + node.label
				}
			}
			numberFrame++
		}
	}
	for _, node := range g.idNodeToNode {
		if node.nameFrame != "" {
			node.label = node.nameFrame + "/" + node.label
		}
	}
	for _, node := range g.idNodeToNode {
		if visited[node.id] == false {
			node.label = UNUSED_COMMANDS + "/" + node.label
		}
	}
}

func (g *Graph) getEdgesAsString() string {
	output := ""
	for _, edge := range g.idEdgeToEdge {
		lines := fmt.Sprintf("%d", edge.source.id) + " -> " + fmt.Sprintf("%d", edge.sink.id) + ";\n"
		output += lines
	}
	return output
}

func (g *Graph) getNodesAsString() string {
	output := ""
	for _, node := range g.idNodeToNode {
		lines := fmt.Sprintf("%d", node.id) + "[label=" + "\"" + node.label + "\"" + "]" + ";\n"
		output += lines
	}
	return output
}

func (g *Graph) getDotFile() string {
	output := "digraph g {\n"
	output += g.getNodesAsString()
	output += g.getEdgesAsString()
	output += "}\n"
	return output
}

func (g *Graph) getNumberNodesInTopLevel() int {
	uniquesNamesInTopLevel := map[string]int{}
	for _, node := range g.idNodeToNode {
		nameTopLevel := splitLabelByChar(&node.label, '/')[0]
		uniquesNamesInTopLevel[nameTopLevel] = 1
	}
	return len(uniquesNamesInTopLevel)
}

func (g *Graph) getPbtxtFile() string {
	var output bytes.Buffer
	validIdNodes := []int{}
	for id, node := range g.idNodeToNode {
		if node.isReal {
			validIdNodes = append(validIdNodes, id)
		}
	}
	sort.Ints(validIdNodes)

	for _, idNode := range validIdNodes {
		node := g.idNodeToNode[idNode]
		output.WriteString("node {\n")
		output.WriteString("\tname: \"" + node.label + "\"\n")
		output.WriteString("\top: \"" + node.label + "\"\n")

		validIdNeighbors := []int{}
		for idNeighbor := range node.incomingIdNodesToIdEdge {
			if g.idNodeToNode[idNeighbor].isReal {
				validIdNeighbors = append(validIdNeighbors, idNeighbor)
			}
		}
		sort.Ints(validIdNeighbors)
		for _, idNeighbor := range validIdNeighbors {
			nodeNeighbor := g.idNodeToNode[idNeighbor]
			output.WriteString("\tinput: \"" + nodeNeighbor.label + "\"\n")
		}
		for i, attribute := range node.attributes {
			output.WriteString("\t\tattr {\n")
			output.WriteString("\t\t\tkey: " + "Param" + fmt.Sprintf("%d", i+1) + "\n")
			output.WriteString("\t\t\tvalue {\n")
			output.WriteString("\t\t\t\t\t: " + attribute + "  \n")
			output.WriteString("\t\t\t}\n")
			output.WriteString("\t\t}\n")
		}
		output.WriteString("}\n")
	}
	return output.String()
}

func getProtoFileFromDependencyGraph(ctx context.Context, g dependencygraph2.DependencyGraph) string {
	msg, err := protoconv.ToProto(ctx, g)
	if err != nil {
		panic(msg)
	}
	output := proto.MarshalTextString(msg)
	return output
}
func (g *Graph) getProtoFile() string {
	validIdNodes := []int{}
	for id, node := range g.idNodeToNode {
		if node.isReal {
			validIdNodes = append(validIdNodes, id)
		}
	}
	sort.Ints(validIdNodes)
	numberValidNodes := len(validIdNodes)

	protoGraph := &protobuf.Graph{}
	protoGraph.Nodes = make([]*protobuf.Node, numberValidNodes)

	for i, idNode := range validIdNodes {
		tmp := g.getNodeProto(idNode)
		protoGraph.Nodes[i] = tmp
	}

	output := proto.MarshalTextString(protoGraph)
	return output
}

func (g *Graph) getNodeProto(idNode int) *protobuf.Node {
	protoNode := &protobuf.Node{}
	node := g.idNodeToNode[idNode]
	protoNode.Name = node.label
	protoNode.Op = node.label

	validIdNeighbors := []int{}
	for idNeighbor := range node.incomingIdNodesToIdEdge {
		if g.idNodeToNode[idNeighbor].isReal {
			validIdNeighbors = append(validIdNeighbors, idNeighbor)
		}
	}
	sort.Ints(validIdNeighbors)
	numberValidNeighbors := len(validIdNeighbors)
	protoNode.Input = make([]string, numberValidNeighbors)
	for i, idNeighbor := range validIdNeighbors {
		protoNode.Input[i] = g.idNodeToNode[idNeighbor].label
	}
	return protoNode
}
