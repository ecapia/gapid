
package graph_visualization

import (
	"testing"
	"sort"
	"reflect"
)

func getSortedKeys(input map[int]int) []int{
	sortedKeys := []int{}
	for key := range input {
		sortedKeys = append(sortedKeys , key)
	}
	sort.Ints(sortedKeys)
	return sortedKeys
}

func getSortedKeysForNodes(input map[int]*Node) []int{
	sortedKeys :=[]int{}
	for key := range input {
		sortedKeys = append(sortedKeys , key)
	}
	sort.Ints(sortedKeys)
	return sortedKeys
}


func areEqualGraphs(t *testing.T, wantedGraph *Graph, testedGraph *Graph) bool {

	if wantedGraph.numberNodes != testedGraph.numberNodes {
		t.Errorf("The numbers of nodes are different %v != %v\n",wantedGraph.numberNodes,testedGraph.numberNodes)
	}
	if wantedGraph.numberEdges != testedGraph.numberEdges {
		t.Errorf("The numbers of edges are different %v != %v\n",wantedGraph.numberEdges,testedGraph.numberEdges)
	}

	wantedSortedIdNodes := getSortedKeysForNodes(wantedGraph.idNodeToNode)
	testedSortedIdNodes := getSortedKeysForNodes(testedGraph.idNodeToNode)
	if reflect.DeepEqual(wantedSortedIdNodes , testedSortedIdNodes) == false {
		t.Errorf("The id nodes are different in the graphs\n" )
	}

	for _ , id := range wantedSortedIdNodes {
		wantedNode := wantedGraph.idNodeToNode[id]
		testedNode := testedGraph.idNodeToNode[id]
		if reflect.DeepEqual(wantedNode.label , testedNode.label) == false {
			t.Errorf("The labels from nodes with id %d are different %v != %v\n",id,wantedNode.label,testedNode.label)
		}
		if reflect.DeepEqual(wantedNode.name , testedNode.name) == false {
			t.Errorf("The name from nodes with id %d are different %v != %v\n",id,wantedNode.name,testedNode.name)
		}
		wantedIncomingIdNodesSorted := getSortedKeys(wantedNode.incomingIdNodesToIdEdge)
		testedIncomingIdNodesSorted := getSortedKeys(testedNode.incomingIdNodesToIdEdge)
		if reflect.DeepEqual(wantedIncomingIdNodesSorted , testedIncomingIdNodesSorted) == false{
			t.Errorf("The incoming id Nodes are different for Nodes with id %d\n",id)
		}

		wantedOutcomingIdNodesSorted := getSortedKeys(wantedNode.outcomingIdNodesToIdEdge)
		testedOutcomingIdNodesSorted := getSortedKeys(testedNode.outcomingIdNodesToIdEdge)
		if reflect.DeepEqual(wantedOutcomingIdNodesSorted , testedOutcomingIdNodesSorted) == false {
			t.Errorf("The outcoming id Nodes are different for Nodes with id %d\n",id)
		}
	}
	return true
}



func TestGraph1(t *testing.T) {

	wantedGraph := createGraph(0)
	wantedGraph.addNodeById(0 , "A")
	wantedGraph.addNodeById(1 , "B")
	wantedGraph.addNodeById(2 , "C")
	wantedGraph.addNodeById(3 , "D")
	wantedGraph.addNodeById(4 , "E")
	wantedGraph.addNodeById(5 , "F")
	wantedGraph.addNodeById(6 , "G")
	wantedGraph.addNodeById(7 , "H")
	wantedGraph.addNodeById(8 , "I")
	wantedGraph.addNodeById(9 , "J")

	testedGraph := createGraph(0)
	testedGraph.addNodeById(0 , "A")
	testedGraph.addNodeById(1 , "B")
	testedGraph.addNodeById(2 , "C")
	testedGraph.addNodeById(3 , "D")
	testedGraph.addNodeById(4 , "E")
	testedGraph.addNodeById(5 , "F")
	testedGraph.addNodeById(6 , "G")
	testedGraph.addNodeById(7 , "H")
	testedGraph.addNodeById(8 , "I")
	testedGraph.addNodeById(9 , "J")

	testedGraph.addNodeById(10 , "K")
	testedGraph.addNodeById(11, "L")
	testedGraph.addNodeById(12, "M")
	testedGraph.removeNodeById(10)
	testedGraph.removeNodeById(11)
	testedGraph.removeNodeById(12)

	if areEqualGraphs(t , wantedGraph , testedGraph) == false {
		t.Errorf("The graphs are different\n")
	}

	testedGraph.addNodeById(10 , "K")
	testedGraph.addNodeById(11, "L")
	testedGraph.addNodeById(12, "M")
	testedGraph.addEdgeByIdNodes(10 , 1)
	testedGraph.addEdgeByIdNodes(10 , 11)
	testedGraph.addEdgeByIdNodes(10 , 4)
	testedGraph.addEdgeByIdNodes(11 , 12)
	testedGraph.addEdgeByIdNodes(2 , 12)
	testedGraph.addEdgeByIdNodes(4 , 11)
	testedGraph.addEdgeByIdNodes(10 , 12)
	testedGraph.removeNodeById(10)
	testedGraph.removeNodeById(11)
	testedGraph.removeNodeById(12)


	if areEqualGraphs(t , wantedGraph , testedGraph) == false {
		t.Errorf("The graphs are different\n")
	}

}

func TestGraph2(t *testing.T){
	labelA := &Label{name : []string{"A"} , id : []int{1}}
	labelB := &Label{name : []string{"B"} , id : []int{2}}
	labelC := &Label{name : []string{"C"} , id : []int{3}}
	labelD := &Label{name : []string{"D"} , id : []int{4}}
	labelE := &Label{name : []string{"E"} , id : []int{5}}
	labelF := &Label{name : []string{"F"} , id : []int{6}}
	labelG := &Label{name : []string{"G"} , id : []int{7}}

	wantedGraph := createGraph(0)
	testedGraph := createGraph(0)
	attributes := []string{}
	testedGraph.addNodeByIdAndNameAndAttrAndIsReal(0,labelA,"vkCommandBuffer0",attributes,true)
	testedGraph.addNodeByIdAndNameAndAttrAndIsReal(1,labelA,"vkCommandBuffer1",attributes,true)
	testedGraph.addNodeByIdAndNameAndAttrAndIsReal(2,labelA,"vkCommandBuffer2",attributes,true)
	testedGraph.addNodeByIdAndNameAndAttrAndIsReal(3,labelA,"vkCommandBuffer3",attributes,true)
	testedGraph.addNodeByIdAndNameAndAttrAndIsReal(4,labelA,"vkCommandBuffer4",attributes,true)
	testedGraph.addNodeByIdAndNameAndAttrAndIsReal(5,labelA,"vkCommandBuffer5",attributes,true)
	testedGraph.addNodeByIdAndNameAndAttrAndIsReal(6,labelA,"vkCommandBuffer6",attributes,true)
	testedGraph.addNodeByIdAndNameAndAttrAndIsReal(7,labelA,"vkCommandBuffer7",attributes,true)
	testedGraph.removeNodesWithZeroDegree()

	if areEqualGraphs(t , wantedGraph , testedGraph) == false {
		t.Errorf("The graphs are different\n")
	}

	wantedGraph.addNodeByIdAndNameAndAttrAndIsReal(0,labelA,"vkCommandBuffer0",attributes,true)
	wantedGraph.addNodeByIdAndNameAndAttrAndIsReal(2,labelB,"vkCommandBuffer2",attributes,true)
	wantedGraph.addNodeByIdAndNameAndAttrAndIsReal(6,labelC,"vkCommandBuffer6",attributes,true)
	wantedGraph.addNodeByIdAndNameAndAttrAndIsReal(3,labelD,"vkCommandBuffer3",attributes,true)
	wantedGraph.addNodeByIdAndNameAndAttrAndIsReal(4,labelE,"vkCommandBuffer4",attributes,true)
	wantedGraph.addNodeByIdAndNameAndAttrAndIsReal(5,labelF,"vkCommandBuffer5",attributes,true)
	wantedGraph.addEdgeByIdNodes(0 , 3)
	wantedGraph.addEdgeByIdNodes(0 , 4)
	wantedGraph.addEdgeByIdNodes(0 , 5)
	wantedGraph.addEdgeByIdNodes(2 , 3)
	wantedGraph.addEdgeByIdNodes(2 , 4)
	wantedGraph.addEdgeByIdNodes(2 , 5)
	wantedGraph.addEdgeByIdNodes(6 , 3)
	wantedGraph.addEdgeByIdNodes(6 , 4)
	wantedGraph.addEdgeByIdNodes(6 , 5)

	testedGraph.addNodeByIdAndNameAndAttrAndIsReal(0,labelA,"vkCommandBuffer0",attributes,true)
	testedGraph.addNodeByIdAndNameAndAttrAndIsReal(2,labelB,"vkCommandBuffer2",attributes,true)
	testedGraph.addNodeByIdAndNameAndAttrAndIsReal(6,labelC,"vkCommandBuffer6",attributes,true)
	testedGraph.addNodeByIdAndNameAndAttrAndIsReal(3,labelD,"vkCommandBuffer3",attributes,true)
	testedGraph.addNodeByIdAndNameAndAttrAndIsReal(4,labelE,"vkCommandBuffer4",attributes,true)
	testedGraph.addNodeByIdAndNameAndAttrAndIsReal(5,labelF,"vkCommandBuffer5",attributes,true)
	testedGraph.addNodeByIdAndNameAndAttrAndIsReal(1,labelG,"vkCommandBuffer1",attributes,true)
	testedGraph.addEdgeByIdNodes(0,1)
	testedGraph.addEdgeByIdNodes(2,1)
	testedGraph.addEdgeByIdNodes(6,1)
	testedGraph.addEdgeByIdNodes(1,3)
	testedGraph.addEdgeByIdNodes(1,4)
	testedGraph.addEdgeByIdNodes(1,5)
	testedGraph.removeNodeByIdKeepingEdges(1)

	if areEqualGraphs(t , wantedGraph , testedGraph) == false {
		t.Errorf("The graphs are different\n")
	}
}

func TestGraph3(t *testing.T) {

	wantedGraph := createGraph(123456)
	testedGraph := createGraph(123455)
	testedGraph.addNodeByDefault("")
	if areEqualGraphs(t , wantedGraph , testedGraph) == false {
		t.Errorf("The graphs are different\n")
	}

	wantedGraph.removeNodesWithZeroDegree()
	testedGraph.removeNodesWithZeroDegree()
	if areEqualGraphs(t , wantedGraph , testedGraph) == false {
		t.Errorf("The graphs are different\n")
	}
	wantedGraph.addNodeById(123456,"")
	wantedGraph.addNodeById(10,"")
	wantedGraph.addEdgeByIdNodes(123456, 10)

	testedGraph.addNodeById(123456,"")
	testedGraph.addNodeById(10,"")
	testedGraph.addEdgeByIdNodes(10,123456)
	testedGraph.removeEdgeById(0)
	testedGraph.addEdgeByIdNodes(123456,10)

	if areEqualGraphs(t , wantedGraph , testedGraph) == false {
		t.Errorf("The graphs are different\n")
	}
}
