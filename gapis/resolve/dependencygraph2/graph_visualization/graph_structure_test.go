package graph_visualization

import (
	"reflect"
	"sort"
	"testing"
)

func getSortedKeys(input map[int]int) []int {
	sortedKeys := []int{}
	for key := range input {
		sortedKeys = append(sortedKeys, key)
	}
	sort.Ints(sortedKeys)
	return sortedKeys
}

func getSortedKeysForNodes(input map[int]*Node) []int {
	sortedKeys := []int{}
	for key := range input {
		sortedKeys = append(sortedKeys, key)
	}
	sort.Ints(sortedKeys)
	return sortedKeys
}

func areEqualGraphs(t *testing.T, wantedGraph *Graph, obtainedGraph *Graph) bool {

	if wantedGraph.numberNodes != obtainedGraph.numberNodes {
		t.Errorf("The numbers of nodes are different %v != %v\n", wantedGraph.numberNodes, obtainedGraph.numberNodes)
	}
	if wantedGraph.numberEdges != obtainedGraph.numberEdges {
		t.Errorf("The numbers of edges are different %v != %v\n", wantedGraph.numberEdges, obtainedGraph.numberEdges)
	}

	wantedSortedIdNodes := getSortedKeysForNodes(wantedGraph.idNodeToNode)
	obtainedSortedIdNodes := getSortedKeysForNodes(obtainedGraph.idNodeToNode)
	if reflect.DeepEqual(wantedSortedIdNodes, obtainedSortedIdNodes) == false {
		t.Errorf("The id nodes are different in the graphs\n")
	}

	for _, id := range wantedSortedIdNodes {
		wantedNode := wantedGraph.idNodeToNode[id]
		obtainedNode := obtainedGraph.idNodeToNode[id]
		if reflect.DeepEqual(wantedNode.label, obtainedNode.label) == false {
			t.Errorf("The labels from nodes with id %d are different %v != %v\n", id, wantedNode.label, obtainedNode.label)
		}
		if reflect.DeepEqual(wantedNode.name, obtainedNode.name) == false {
			t.Errorf("The name from nodes with id %d are different %v != %v\n", id, wantedNode.name, obtainedNode.name)
		}
		wantedIncomingIdNodesSorted := getSortedKeys(wantedNode.incomingIdNodesToIdEdge)
		obtainedIncomingIdNodesSorted := getSortedKeys(obtainedNode.incomingIdNodesToIdEdge)
		if reflect.DeepEqual(wantedIncomingIdNodesSorted, obtainedIncomingIdNodesSorted) == false {
			t.Errorf("The incoming id Nodes are different for Nodes with id %d\n", id)
		}

		wantedOutcomingIdNodesSorted := getSortedKeys(wantedNode.outcomingIdNodesToIdEdge)
		obtainedOutcomingIdNodesSorted := getSortedKeys(obtainedNode.outcomingIdNodesToIdEdge)
		if reflect.DeepEqual(wantedOutcomingIdNodesSorted, obtainedOutcomingIdNodesSorted) == false {
			t.Errorf("The outcoming id Nodes are different for Nodes with id %d\n", id)
		}
	}
	return true
}

func TestGraph1(t *testing.T) {

	wantedGraph := createGraph(0)
	wantedGraph.addNodeById(0, "A")
	wantedGraph.addNodeById(1, "B")
	wantedGraph.addNodeById(2, "C")
	wantedGraph.addNodeById(3, "D")
	wantedGraph.addNodeById(4, "E")
	wantedGraph.addNodeById(5, "F")
	wantedGraph.addNodeById(6, "G")
	wantedGraph.addNodeById(7, "H")
	wantedGraph.addNodeById(8, "I")
	wantedGraph.addNodeById(9, "J")

	obtainedGraph := createGraph(0)
	obtainedGraph.addNodeById(0, "A")
	obtainedGraph.addNodeById(1, "B")
	obtainedGraph.addNodeById(2, "C")
	obtainedGraph.addNodeById(3, "D")
	obtainedGraph.addNodeById(4, "E")
	obtainedGraph.addNodeById(5, "F")
	obtainedGraph.addNodeById(6, "G")
	obtainedGraph.addNodeById(7, "H")
	obtainedGraph.addNodeById(8, "I")
	obtainedGraph.addNodeById(9, "J")

	obtainedGraph.addNodeById(10, "K")
	obtainedGraph.addNodeById(11, "L")
	obtainedGraph.addNodeById(12, "M")
	obtainedGraph.removeNodeById(10)
	obtainedGraph.removeNodeById(11)
	obtainedGraph.removeNodeById(12)

	if areEqualGraphs(t, wantedGraph, obtainedGraph) == false {
		t.Errorf("The graphs are different\n")
	}

	obtainedGraph.addNodeById(10, "K")
	obtainedGraph.addNodeById(11, "L")
	obtainedGraph.addNodeById(12, "M")
	obtainedGraph.addEdgeByIdNodes(10, 1)
	obtainedGraph.addEdgeByIdNodes(10, 11)
	obtainedGraph.addEdgeByIdNodes(10, 4)
	obtainedGraph.addEdgeByIdNodes(11, 12)
	obtainedGraph.addEdgeByIdNodes(2, 12)
	obtainedGraph.addEdgeByIdNodes(4, 11)
	obtainedGraph.addEdgeByIdNodes(10, 12)
	obtainedGraph.removeNodeById(10)
	obtainedGraph.removeNodeById(11)
	obtainedGraph.removeNodeById(12)

	if areEqualGraphs(t, wantedGraph, obtainedGraph) == false {
		t.Errorf("The graphs are different\n")
	}

}

func TestGraph2(t *testing.T) {
	labelA := &Label{name: []string{"A"}, id: []int{1}}
	labelB := &Label{name: []string{"B"}, id: []int{2}}
	labelC := &Label{name: []string{"C"}, id: []int{3}}
	labelD := &Label{name: []string{"D"}, id: []int{4}}
	labelE := &Label{name: []string{"E"}, id: []int{5}}
	labelF := &Label{name: []string{"F"}, id: []int{6}}
	labelG := &Label{name: []string{"G"}, id: []int{7}}

	wantedGraph := createGraph(0)
	obtainedGraph := createGraph(0)
	attributes := []string{}
	obtainedGraph.addNodeByIdAndNameAndAttrAndIsReal(0, labelA, "vkCommandBuffer0", attributes, true)
	obtainedGraph.addNodeByIdAndNameAndAttrAndIsReal(1, labelA, "vkCommandBuffer1", attributes, true)
	obtainedGraph.addNodeByIdAndNameAndAttrAndIsReal(2, labelA, "vkCommandBuffer2", attributes, true)
	obtainedGraph.addNodeByIdAndNameAndAttrAndIsReal(3, labelA, "vkCommandBuffer3", attributes, true)
	obtainedGraph.addNodeByIdAndNameAndAttrAndIsReal(4, labelA, "vkCommandBuffer4", attributes, true)
	obtainedGraph.addNodeByIdAndNameAndAttrAndIsReal(5, labelA, "vkCommandBuffer5", attributes, true)
	obtainedGraph.addNodeByIdAndNameAndAttrAndIsReal(6, labelA, "vkCommandBuffer6", attributes, true)
	obtainedGraph.addNodeByIdAndNameAndAttrAndIsReal(7, labelA, "vkCommandBuffer7", attributes, true)
	obtainedGraph.removeNodesWithZeroDegree()

	if areEqualGraphs(t, wantedGraph, obtainedGraph) == false {
		t.Errorf("The graphs are different\n")
	}

	wantedGraph.addNodeByIdAndNameAndAttrAndIsReal(0, labelA, "vkCommandBuffer0", attributes, true)
	wantedGraph.addNodeByIdAndNameAndAttrAndIsReal(2, labelB, "vkCommandBuffer2", attributes, true)
	wantedGraph.addNodeByIdAndNameAndAttrAndIsReal(6, labelC, "vkCommandBuffer6", attributes, true)
	wantedGraph.addNodeByIdAndNameAndAttrAndIsReal(3, labelD, "vkCommandBuffer3", attributes, true)
	wantedGraph.addNodeByIdAndNameAndAttrAndIsReal(4, labelE, "vkCommandBuffer4", attributes, true)
	wantedGraph.addNodeByIdAndNameAndAttrAndIsReal(5, labelF, "vkCommandBuffer5", attributes, true)
	wantedGraph.addEdgeByIdNodes(0, 3)
	wantedGraph.addEdgeByIdNodes(0, 4)
	wantedGraph.addEdgeByIdNodes(0, 5)
	wantedGraph.addEdgeByIdNodes(2, 3)
	wantedGraph.addEdgeByIdNodes(2, 4)
	wantedGraph.addEdgeByIdNodes(2, 5)
	wantedGraph.addEdgeByIdNodes(6, 3)
	wantedGraph.addEdgeByIdNodes(6, 4)
	wantedGraph.addEdgeByIdNodes(6, 5)

	obtainedGraph.addNodeByIdAndNameAndAttrAndIsReal(0, labelA, "vkCommandBuffer0", attributes, true)
	obtainedGraph.addNodeByIdAndNameAndAttrAndIsReal(2, labelB, "vkCommandBuffer2", attributes, true)
	obtainedGraph.addNodeByIdAndNameAndAttrAndIsReal(6, labelC, "vkCommandBuffer6", attributes, true)
	obtainedGraph.addNodeByIdAndNameAndAttrAndIsReal(3, labelD, "vkCommandBuffer3", attributes, true)
	obtainedGraph.addNodeByIdAndNameAndAttrAndIsReal(4, labelE, "vkCommandBuffer4", attributes, true)
	obtainedGraph.addNodeByIdAndNameAndAttrAndIsReal(5, labelF, "vkCommandBuffer5", attributes, true)
	obtainedGraph.addNodeByIdAndNameAndAttrAndIsReal(1, labelG, "vkCommandBuffer1", attributes, true)
	obtainedGraph.addEdgeByIdNodes(0, 1)
	obtainedGraph.addEdgeByIdNodes(2, 1)
	obtainedGraph.addEdgeByIdNodes(6, 1)
	obtainedGraph.addEdgeByIdNodes(1, 3)
	obtainedGraph.addEdgeByIdNodes(1, 4)
	obtainedGraph.addEdgeByIdNodes(1, 5)
	obtainedGraph.removeNodeByIdKeepingEdges(1)

	if areEqualGraphs(t, wantedGraph, obtainedGraph) == false {
		t.Errorf("The graphs are different\n")
	}
}

func TestGraph3(t *testing.T) {

	wantedGraph := createGraph(123456)
	obtainedGraph := createGraph(123455)
	obtainedGraph.addNodeByDefault("")
	if areEqualGraphs(t, wantedGraph, obtainedGraph) == false {
		t.Errorf("The graphs are different\n")
	}

	wantedGraph.removeNodesWithZeroDegree()
	obtainedGraph.removeNodesWithZeroDegree()
	if areEqualGraphs(t, wantedGraph, obtainedGraph) == false {
		t.Errorf("The graphs are different\n")
	}
	wantedGraph.addNodeById(123456, "")
	wantedGraph.addNodeById(10, "")
	wantedGraph.addEdgeByIdNodes(123456, 10)

	obtainedGraph.addNodeById(123456, "")
	obtainedGraph.addNodeById(10, "")
	obtainedGraph.addEdgeByIdNodes(10, 123456)
	obtainedGraph.removeEdgeById(0)
	obtainedGraph.addEdgeByIdNodes(123456, 10)

	if areEqualGraphs(t, wantedGraph, obtainedGraph) == false {
		t.Errorf("The graphs are different\n")
	}
}
