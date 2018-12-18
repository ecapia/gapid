package graph_visualization

import (
	"bytes"
	"fmt"
	"sort"
)

type Label struct {
	name []string
	id   []int
}

func (label *Label) pushBack(name string, id int) {
	label.name = append(label.name, name)
	label.id = append(label.id, id)
}
func (label *Label) pushFront(name string, id int) {
	temporal := &Label{name: []string{name}, id: []int{id}}
	temporal.pushBackLabel(label)
	label.name = temporal.name
	label.id = temporal.id
}

func (label *Label) pushBackLabel(labelToPush *Label) {
	label.name = append(label.name, labelToPush.name...)
	label.id = append(label.id, labelToPush.id...)
}

func (label *Label) update(index int, name string, id int) bool {
	if index >= len(label.id) {
		return false
	}
	label.name[index] = name
	label.id[index] = id
	return true
}

func getMaxCommonPrefix(label1 *Label, label2 *Label) int {
	size := len(label1.id)
	if len(label2.id) < size {
		size = len(label2.id)
	}
	for i := 0; i < size; i++ {
		if label1.name[i] != label2.name[i] || label1.id[i] != label2.id[i] {
			return i
		}
	}
	return size
}

func (label *Label) getLabelAsString() string {
	var output bytes.Buffer
	size := len(label.name)
	for i := range label.name {
		if i+1 < size {
			output.WriteString(label.name[i] + fmt.Sprintf("_%d/", label.id[i]))
		} else {
			output.WriteString(label.name[i] + fmt.Sprintf("_%d", label.id[i]))
		}
	}
	return output.String()
}

type Node struct {
	incomingIdNodesToIdEdge       map[int]int
	outcomingIdNodesToIdEdge      map[int]int
	incomingNodes, outcomingNodes []*Node
	id                            int
	name                          string
	nameFrame                     string
	attributes                    []string
	isReal                        bool
	label                         *Label
	color                         string
}

type Edge struct {
	source, sink *Node
	id           int
	label        string
}

type Graph struct {
	idNodeToNode map[int]*Node
	idEdgeToEdge map[int]*Edge
	nodes        []*Node
	edges        []*Edge
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

func (g *Graph) addNodeByDefault(name string) int {
	id := g.maxIdNode
	newNode := &Node{incomingIdNodesToIdEdge: map[int]int{}, outcomingIdNodesToIdEdge: map[int]int{}, id: id, name: name}
	g.idNodeToNode[id] = newNode
	g.numberNodes++
	g.maxIdNode++
	return id
}

func (g *Graph) addNodeById(id int, name string) bool {
	_, ok := g.idNodeToNode[id]
	if ok == true {
		return false
	}

	newNode := &Node{incomingIdNodesToIdEdge: map[int]int{}, outcomingIdNodesToIdEdge: map[int]int{}, id: id, name: name}
	g.idNodeToNode[id] = newNode
	g.numberNodes++
	g.maxIdNode++
	if g.maxIdNode <= id {
		g.maxIdNode = id + 1
	}
	return true
}

func (g *Graph) addNodeByIdAndNameAndAttrAndIsReal(id int, label *Label, name string, attributes []string, isReal bool) bool {
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

func (g *Graph) getNodesFromKeysInMap(input map[int]int) []*Node {
	nodes := []*Node{}
	for key := range input {
		nodes = append(nodes, g.idNodeToNode[key])
	}
	return nodes
}

func (g *Graph) makeAdjacentList() {
	g.nodes = make([]*Node, 0)
	for _, node := range g.idNodeToNode {
		g.nodes = append(g.nodes, node)
	}
	for _, node := range g.nodes {
		node.incomingNodes = g.getNodesFromKeysInMap(node.incomingIdNodesToIdEdge)
		node.outcomingNodes = g.getNodesFromKeysInMap(node.outcomingIdNodesToIdEdge)
	}
}

func (g *Graph) sortAdjacentList() {
	sort.Sort(NodeSorter(g.nodes))
	for _, node := range g.nodes {
		sort.Sort(NodeSorter(node.incomingNodes))
		sort.Sort(NodeSorter(node.outcomingNodes))
	}
}
