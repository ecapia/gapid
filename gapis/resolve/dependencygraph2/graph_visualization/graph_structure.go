package graph_visualization

import (
	"fmt"
	"os"
	"strconv"
)

type Node struct {
	incomingIdNodesToIdEdge  map[int]int
	outcomingIdNodesToIdEdge map[int]int
	id                       int
	label                    string
	commandName              string
	idCommandType            int
	attributes               string
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

func (g *Graph) addNodeByIdAndIdCommandType(id int, label string, idCommandType int) bool {
	_, ok := g.idNodeToNode[id]
	if ok == true {
		return false
	}

	newNode := &Node{incomingIdNodesToIdEdge: map[int]int{}, outcomingIdNodesToIdEdge: map[int]int{},
		id: id, idCommandType: idCommandType, label: label}
	g.idNodeToNode[id] = newNode
	g.numberNodes++
	g.maxIdNode++
	if g.maxIdNode <= id {
		g.maxIdNode = id + 1
	}
	return true
}

func (g *Graph) addNodeByIdAndIdCommandTypeAndAttr(id int, label string, idCommandType int, attributes string) bool {
	_, ok := g.idNodeToNode[id]
	if ok == true {
		return false
	}

	newNode := &Node{incomingIdNodesToIdEdge: map[int]int{}, outcomingIdNodesToIdEdge: map[int]int{},
		id: id, idCommandType: idCommandType, label: label, attributes: attributes}
	g.idNodeToNode[id] = newNode
	g.numberNodes++
	g.maxIdNode++
	if g.maxIdNode <= id {
		g.maxIdNode = id + 1
	}
	return true
}

func (g *Graph) addNodeByIdAndCommandNameAndAttr(id int, label string, commandName string, attributes string) bool {
	_, ok := g.idNodeToNode[id]
	if ok == true {
		return false
	}

	newNode := &Node{incomingIdNodesToIdEdge: map[int]int{}, outcomingIdNodesToIdEdge: map[int]int{},
		id: id, label: label, commandName: commandName, attributes: attributes}
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

func (g *Graph) addEdgeByNode(source, sink *Node) {
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

func (g *Graph) removeNodeKeepingEdges(idNode int) bool {
	if g.joinEdgesThroughtNode(idNode) == false {
		return false
	}
	if g.removeNodeById(idNode) == false {
		return false
	}
	return true
}

func (g *Graph) dfs(curr *Node, time, mini, idInSCC, s *[]int, currSCC, counter *int) {
	*s = append(*s, curr.id)
	(*time)[curr.id] = *counter
	(*mini)[curr.id] = *counter
	(*counter)++

	for idNext := range curr.outcomingIdNodesToIdEdge {
		next := g.idNodeToNode[idNext]
		if (*time)[next.id] == 0 {
			g.dfs(next, time, mini, idInSCC, s, currSCC, counter)
		}
		if (*time)[next.id] != -1 {
			if (*mini)[next.id] < (*mini)[curr.id] {
				(*mini)[curr.id] = (*mini)[next.id]
			}
		}
	}

	if (*mini)[curr.id] == (*time)[curr.id] {
		for {
			tmp := (*s)[len(*s)-1]
			(*time)[tmp] = -1
			*s = (*s)[:len(*s)-1]
			(*idInSCC)[tmp] = *currSCC
			if tmp == curr.id {
				break
			}
		}
		(*currSCC)++
	}
}

func (g *Graph) getSCC() []int {
	currSCC := 0
	counter := 1
	time := make([]int, g.maxIdNode)
	mini := make([]int, g.maxIdNode)
	idInSCC := make([]int, g.maxIdNode)
	s := make([]int, 0)

	for _, curr := range g.idNodeToNode {
		if time[curr.id] == 0 {
			g.dfs(curr, &time, &mini, &idInSCC, &s, &currSCC, &counter)
		}
	}
	return idInSCC
}

func (g *Graph) makeSccCompressionByIdCommandType() {
	newGraph := createGraph(0)
	for _, curr := range g.idNodeToNode {
		newGraph.addNodeById(curr.idCommandType, "")
	}

	for _, curr := range g.idNodeToNode {
		for idNext := range curr.outcomingIdNodesToIdEdge {
			next := g.idNodeToNode[idNext]
			newGraph.addEdgeByIdNodes(curr.idCommandType, next.idCommandType)
		}
	}
	idInSCC := newGraph.getSCC()
	for _, curr := range g.idNodeToNode {
		scc := idInSCC[curr.idCommandType]
		curr.label = curr.label + "/" + fmt.Sprintf("scc%d", scc)
	}
}

func printMap(m map[int]int) {
	for k, v := range m {
		fmt.Println("idNeighbor = ", k, " idEdge = ", v)
	}
}

func (g *Graph) printEdges() {
	for _, e := range g.idEdgeToEdge {
		fmt.Println(" ( ", e.source.id, ",", e.sink.id, ")", "idEdge = ", e.id)
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

func (g *Graph) writeEdges(f *os.File) {
	for _, e := range g.idEdgeToEdge {
		line := strconv.Itoa(e.source.id) + " -> " + strconv.Itoa(e.sink.id) + ";\n"
		f.WriteString(line)
	}
}

func (g *Graph) writeNodes(f *os.File) {
	for _, n := range g.idNodeToNode {
		line := strconv.Itoa(n.id) + "[label=" + n.label + "]" + ";\n"
		f.WriteString(line)
	}
}

func (g *Graph) writeDigraph(filename string) {
	_, e := os.Stat(filename)
	if os.IsNotExist(e) {
		os.Create(filename)
	}

	f, err := os.OpenFile(filename, os.O_RDWR, 0600)
	if err != nil {
		panic(err)
	}
	f.WriteString("digraph g {\n")
	g.writeNodes(f)
	g.writeEdges(f)
	f.WriteString("}\n")
}

func (g *Graph) getEdgesAsString() string {
	output := ""
	for _, e := range g.idEdgeToEdge {
		line := strconv.Itoa(e.source.id) + " -> " + strconv.Itoa(e.sink.id) + ";\n"
		output += line
	}
	return output
}

func (g *Graph) getNodesAsString() string {
	output := ""
	for _, n := range g.idNodeToNode {
		line := strconv.Itoa(n.id) + "[label=" + n.label + "]" + ";\n"
		output += line
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

	for _, node := range g.idNodeToNode {
		line := "node {\n"
		line += "name: " + node.label + "\n"
		line += "op: " + node.label + "\n"
		for idNeighbor := range node.incomingIdNodesToIdEdge {
			nodeNeighbor := g.idNodeToNode[idNeighbor]
			line += "input: " + nodeNeighbor.label + "\n"
		}
		line += "attr {\n"
		line += "key: " + "\"" + node.commandName + "\"\n"
		line += "}\n"

		line += "}\n"
		output += line
	}
	return output
}
