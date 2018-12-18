package graph_visualization

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/google/gapid/core/data/protoconv"
	"github.com/google/gapid/gapis/resolve/dependencygraph2"
	"sort"
)

type Node struct {
	incomingIdNodesToIdEdge  map[int]int
	outcomingIdNodesToIdEdge map[int]int
	id                       int
	label                    string
	name                     string
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

func printMap(m map[int]int) {
	for k, v := range m {
		fmt.Println("idNeighbor = ", k, " idEdge = ", v)
	}
}

func (g *Graph) printEdges() {
	for _, edge := range g.idEdgeToEdge {
		fmt.Println(" ( ", edge.source.id, ",", edge.sink.id, ")", "idEdge = ", edge.id)
	}
}

func (g *Graph) printNodes() {
	for id, node := range g.idNodeToNode {
		fmt.Println("node = ", id)
		fmt.Println("in = ")
		printMap(node.incomingIdNodesToIdEdge)
		fmt.Println("out = ")
		printMap(node.outcomingIdNodesToIdEdge)
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

func (g *Graph) getPbtxtFile() string {
	output := ""
	orderedIdNodes := []int{}
	for id := range g.idNodeToNode {
		orderedIdNodes = append(orderedIdNodes, id)
	}
	sort.Ints(orderedIdNodes)
	for _, idNode := range orderedIdNodes {
		node := g.idNodeToNode[idNode]
		if node.isReal == false {
			continue
		}
		lines := "node {\n"
		lines += "\tname: " + "\"" + node.label + "\"" + "\n"
		lines += "\top: " + "\"" + node.label + "\"" + "\n"

		orderedIdEdges := []int{}
		for idNeighbor := range node.incomingIdNodesToIdEdge {
			orderedIdEdges = append(orderedIdEdges, idNeighbor)
		}
		sort.Ints(orderedIdEdges)

		for _, idNeighbor := range orderedIdEdges {
			nodeNeighbor := g.idNodeToNode[idNeighbor]
			if nodeNeighbor.isReal == false {
				continue
			}
			lines += "\tinput: " + "\"" + nodeNeighbor.label + "\"" + "\n"
		}

		for i, val := range node.attributes {
			lines += "\t\tattr {\n"
			lines += "\t\t\tkey: " + "Param" + fmt.Sprintf("%d", i+1) + "\n"
			lines += "\t\t\tvalue {\n"
			lines += "\t\t\t\t\t" + val + "  \n"
			lines += "\t\t\t}\n"
			lines += "\t\t}\n"
		}

		lines += "}\n"
		output += lines
	}
	return output
}

func getProtoFile(ctx context.Context, dependencyGraph dependencygraph2.DependencyGraph) string {
	msg, err := protoconv.ToProto(ctx, dependencyGraph)
	if err != nil {
		panic(msg)
	}
	output := proto.MarshalTextString(msg)
	return output
}
