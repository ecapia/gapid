package graph_visualization

import (
	"reflect"
	"sort"
	"testing"
)

func TestDfs1(t *testing.T) {

	numberNodes := 10
	graph := createGraph(numberNodes)
	graph.addEdgeByIdNodes(6, 2)
	graph.addEdgeByIdNodes(6, 3)
	graph.addEdgeByIdNodes(6, 7)
	graph.addEdgeByIdNodes(3, 7)
	graph.addEdgeByIdNodes(7, 9)

	graph.addEdgeByIdNodes(5, 4)
	graph.addEdgeByIdNodes(5, 0)
	graph.addEdgeByIdNodes(0, 4)

	graph.idNodeToNode[6].name = VK_QUEUE_PRESENT
	graph.idNodeToNode[5].name = VK_QUEUE_PRESENT
	graph.idNodeToNode[8].name = VK_QUEUE_PRESENT
	graph.idNodeToNode[1].name = VK_QUEUE_PRESENT

	graph.makeAdjacentList()
	graph.sortAdjacentList()
	for _, node := range graph.nodes {
		node.isReal = true
	}

	wantedComponents := [][]*Node{
		[]*Node{graph.nodes[1]},
		[]*Node{graph.nodes[0], graph.nodes[4], graph.nodes[5]},
		[]*Node{graph.nodes[2], graph.nodes[3], graph.nodes[6], graph.nodes[7], graph.nodes[9]},
		[]*Node{graph.nodes[8]},
	}

	obtainedComponents := [][]*Node{}
	visited := make([]bool, numberNodes)
	for _, node := range graph.nodes {
		if node.name == VK_QUEUE_PRESENT && visited[node.id] == false && node.isReal {
			visitedNodes := []*Node{}
			graph.dfs(node, &visited, &visitedNodes)
			sort.Sort(NodeSorter(visitedNodes))
			obtainedComponents = append(obtainedComponents, visitedNodes)
		}
	}

	for numberComponent := range obtainedComponents {
		if reflect.DeepEqual(obtainedComponents[numberComponent], wantedComponents[numberComponent]) == false {
			t.Errorf("The component %d is different\n", numberComponent)
			t.Errorf("Wanted %v , obtained %v\n", wantedComponents[numberComponent], obtainedComponents[numberComponent])
		}
	}

}

func TestDfs2(t *testing.T) {
	numberNodes := 100000
	graph := createGraph(numberNodes)
	for i := 0; i+1 < numberNodes; i++ {
		graph.addEdgeByIdNodes(i, i+1)
	}
	graph.idNodeToNode[0].name = VK_QUEUE_PRESENT

	graph.makeAdjacentList()
	graph.sortAdjacentList()
	wantedComponents := [][]*Node{[]*Node{}}
	for _, node := range graph.nodes {
		node.isReal = true
		wantedComponents[0] = append(wantedComponents[0], node)
	}

	obtainedComponents := [][]*Node{}
	visited := make([]bool, numberNodes)
	for _, node := range graph.nodes {
		if node.name == VK_QUEUE_PRESENT && visited[node.id] == false && node.isReal {
			visitedNodes := []*Node{}
			graph.dfs(node, &visited, &visitedNodes)
			sort.Sort(NodeSorter(visitedNodes))
			obtainedComponents = append(obtainedComponents, visitedNodes)
		}
	}
	for numberComponent := range obtainedComponents {
		if reflect.DeepEqual(obtainedComponents[numberComponent], wantedComponents[numberComponent]) == false {
			t.Errorf("The component %d is different\n", numberComponent)
			t.Errorf("Wanted %v , obtained %v\n", wantedComponents[numberComponent], obtainedComponents[numberComponent])
		}
	}
}

func TestMakeChunks1(t *testing.T) {
	numberNodes := 30
	graph := createGraph(numberNodes)
	auxiliarCommands := []string{
		"SetBarrier",
		"DrawIndex",
		"SetScissor",
		"SetIndexForDraw",
	}
	labels := []*Label{
		&Label{name: []string{auxiliarCommands[0]}, id: []int{0}},
		&Label{name: []string{COMMAND_BUFFER, VK_BEGIN_COMMAND_BUFFER}, id: []int{1, 1}},
		&Label{name: []string{COMMAND_BUFFER, VK_RENDER_PASS, VK_CMD_BEGIN_RENDER_PASS}, id: []int{1, 1, 2}},
		&Label{name: []string{COMMAND_BUFFER, VK_RENDER_PASS, VK_SUBPASS, auxiliarCommands[1]}, id: []int{1, 1, 1, 3}},
		&Label{name: []string{COMMAND_BUFFER, VK_RENDER_PASS, VK_SUBPASS, auxiliarCommands[0]}, id: []int{1, 1, 1, 4}},
		&Label{name: []string{COMMAND_BUFFER, VK_RENDER_PASS, VK_SUBPASS, auxiliarCommands[0]}, id: []int{1, 1, 2, 5}},
		&Label{name: []string{COMMAND_BUFFER, VK_RENDER_PASS, VK_SUBPASS, auxiliarCommands[0]}, id: []int{1, 1, 2, 6}},
		&Label{name: []string{COMMAND_BUFFER, VK_RENDER_PASS, VK_SUBPASS, auxiliarCommands[0]}, id: []int{1, 1, 2, 7}},
		&Label{name: []string{COMMAND_BUFFER, VK_RENDER_PASS, VK_SUBPASS, auxiliarCommands[0]}, id: []int{1, 1, 2, 8}},
		&Label{name: []string{COMMAND_BUFFER, VK_RENDER_PASS, VK_SUBPASS, auxiliarCommands[0]}, id: []int{1, 1, 2, 9}},
		&Label{name: []string{COMMAND_BUFFER, VK_RENDER_PASS, VK_SUBPASS, auxiliarCommands[0]}, id: []int{1, 1, 2, 10}},
		&Label{name: []string{COMMAND_BUFFER, VK_RENDER_PASS, VK_SUBPASS, auxiliarCommands[0]}, id: []int{1, 1, 2, 11}},
		&Label{name: []string{COMMAND_BUFFER, VK_RENDER_PASS, VK_SUBPASS, auxiliarCommands[0]}, id: []int{1, 1, 2, 12}},
		&Label{name: []string{COMMAND_BUFFER, VK_RENDER_PASS, VK_SUBPASS, auxiliarCommands[0]}, id: []int{1, 1, 2, 13}},
		&Label{name: []string{COMMAND_BUFFER, VK_RENDER_PASS, VK_SUBPASS, auxiliarCommands[0]}, id: []int{1, 1, 2, 14}},
		&Label{name: []string{auxiliarCommands[0]}, id: []int{15}},
		&Label{name: []string{auxiliarCommands[0]}, id: []int{16}},
		&Label{name: []string{auxiliarCommands[0]}, id: []int{17}},
		&Label{name: []string{auxiliarCommands[0]}, id: []int{18}},
		&Label{name: []string{auxiliarCommands[0]}, id: []int{19}},
		&Label{name: []string{auxiliarCommands[0]}, id: []int{20}},
		&Label{name: []string{auxiliarCommands[0]}, id: []int{21}},
		&Label{name: []string{COMMAND_BUFFER, VK_RENDER_PASS, VK_SUBPASS, auxiliarCommands[0]}, id: []int{1, 2, 1, 22}},
		&Label{name: []string{COMMAND_BUFFER, VK_RENDER_PASS, VK_SUBPASS, auxiliarCommands[0]}, id: []int{1, 2, 2, 23}},
		&Label{name: []string{COMMAND_BUFFER, VK_RENDER_PASS, VK_SUBPASS, auxiliarCommands[0]}, id: []int{1, 2, 3, 24}},
		&Label{name: []string{COMMAND_BUFFER, VK_RENDER_PASS, VK_SUBPASS, auxiliarCommands[0]}, id: []int{1, 2, 4, 25}},
		&Label{name: []string{COMMAND_BUFFER, VK_RENDER_PASS, VK_SUBPASS, auxiliarCommands[0]}, id: []int{1, 2, 5, 26}},
		&Label{name: []string{COMMAND_BUFFER, VK_RENDER_PASS, VK_SUBPASS, auxiliarCommands[0]}, id: []int{1, 2, 6, 27}},
		&Label{name: []string{auxiliarCommands[0]}, id: []int{28}},
		&Label{name: []string{auxiliarCommands[0]}, id: []int{29}},
	}
	graph.makeAdjacentList()
	graph.sortAdjacentList()
	for i := 0; i < numberNodes; i++ {
		graph.nodes[i].label = labels[i]
	}
	makeChunks(&graph.nodes)

	wantedLabels := []*Label{
		&Label{name: []string{SUPER + auxiliarCommands[0], SUPER + auxiliarCommands[0], auxiliarCommands[0]}, id: []int{0, 1, 0}},
		&Label{name: []string{COMMAND_BUFFER, VK_BEGIN_COMMAND_BUFFER}, id: []int{1, 1}},
		&Label{name: []string{COMMAND_BUFFER, SUPER + VK_RENDER_PASS, VK_RENDER_PASS, VK_CMD_BEGIN_RENDER_PASS}, id: []int{1, 0, 1, 2}},
		&Label{name: []string{COMMAND_BUFFER, SUPER + VK_RENDER_PASS, VK_RENDER_PASS, SUPER + VK_SUBPASS, VK_SUBPASS, auxiliarCommands[1]},
			id: []int{1, 0, 1, 0, 1, 3}},
		&Label{name: []string{COMMAND_BUFFER, SUPER + VK_RENDER_PASS, VK_RENDER_PASS, SUPER + VK_SUBPASS, VK_SUBPASS, auxiliarCommands[0]},
			id: []int{1, 0, 1, 0, 1, 4}},
		&Label{name: []string{COMMAND_BUFFER, SUPER + VK_RENDER_PASS, VK_RENDER_PASS, SUPER + VK_SUBPASS, VK_SUBPASS, SUPER + auxiliarCommands[0],
			SUPER + auxiliarCommands[0], auxiliarCommands[0]}, id: []int{1, 0, 1, 0, 2, 0, 1, 5}},
		&Label{name: []string{COMMAND_BUFFER, SUPER + VK_RENDER_PASS, VK_RENDER_PASS, SUPER + VK_SUBPASS, VK_SUBPASS, SUPER + auxiliarCommands[0],
			SUPER + auxiliarCommands[0], auxiliarCommands[0]}, id: []int{1, 0, 1, 0, 2, 0, 1, 6}},
		&Label{name: []string{COMMAND_BUFFER, SUPER + VK_RENDER_PASS, VK_RENDER_PASS, SUPER + VK_SUBPASS, VK_SUBPASS, SUPER + auxiliarCommands[0],
			SUPER + auxiliarCommands[0], auxiliarCommands[0]}, id: []int{1, 0, 1, 0, 2, 0, 2, 7}},
		&Label{name: []string{COMMAND_BUFFER, SUPER + VK_RENDER_PASS, VK_RENDER_PASS, SUPER + VK_SUBPASS, VK_SUBPASS, SUPER + auxiliarCommands[0],
			SUPER + auxiliarCommands[0], auxiliarCommands[0]}, id: []int{1, 0, 1, 0, 2, 0, 2, 8}},
		&Label{name: []string{COMMAND_BUFFER, SUPER + VK_RENDER_PASS, VK_RENDER_PASS, SUPER + VK_SUBPASS, VK_SUBPASS, SUPER + auxiliarCommands[0],
			SUPER + auxiliarCommands[0], auxiliarCommands[0]}, id: []int{1, 0, 1, 0, 2, 0, 3, 9}},
		&Label{name: []string{COMMAND_BUFFER, SUPER + VK_RENDER_PASS, VK_RENDER_PASS, SUPER + VK_SUBPASS, VK_SUBPASS, SUPER + auxiliarCommands[0],
			SUPER + auxiliarCommands[0], auxiliarCommands[0]}, id: []int{1, 0, 1, 0, 2, 0, 3, 10}},
		&Label{name: []string{COMMAND_BUFFER, SUPER + VK_RENDER_PASS, VK_RENDER_PASS, SUPER + VK_SUBPASS, VK_SUBPASS, SUPER + auxiliarCommands[0],
			SUPER + auxiliarCommands[0], auxiliarCommands[0]}, id: []int{1, 0, 1, 0, 2, 0, 4, 11}},
		&Label{name: []string{COMMAND_BUFFER, SUPER + VK_RENDER_PASS, VK_RENDER_PASS, SUPER + VK_SUBPASS, VK_SUBPASS, SUPER + auxiliarCommands[0],
			SUPER + auxiliarCommands[0], auxiliarCommands[0]}, id: []int{1, 0, 1, 0, 2, 0, 4, 12}},
		&Label{name: []string{COMMAND_BUFFER, SUPER + VK_RENDER_PASS, VK_RENDER_PASS, SUPER + VK_SUBPASS, VK_SUBPASS, SUPER + auxiliarCommands[0],
			SUPER + auxiliarCommands[0], auxiliarCommands[0]}, id: []int{1, 0, 1, 0, 2, 0, 5, 13}},
		&Label{name: []string{COMMAND_BUFFER, SUPER + VK_RENDER_PASS, VK_RENDER_PASS, SUPER + VK_SUBPASS, VK_SUBPASS, SUPER + auxiliarCommands[0],
			SUPER + auxiliarCommands[0], auxiliarCommands[0]}, id: []int{1, 0, 1, 0, 2, 0, 5, 14}},
		&Label{name: []string{SUPER + auxiliarCommands[0], SUPER + auxiliarCommands[0], auxiliarCommands[0]}, id: []int{0, 1, 15}},
		&Label{name: []string{SUPER + auxiliarCommands[0], SUPER + auxiliarCommands[0], auxiliarCommands[0]}, id: []int{0, 2, 16}},
		&Label{name: []string{SUPER + auxiliarCommands[0], SUPER + auxiliarCommands[0], auxiliarCommands[0]}, id: []int{0, 2, 17}},
		&Label{name: []string{SUPER + auxiliarCommands[0], SUPER + auxiliarCommands[0], auxiliarCommands[0]}, id: []int{0, 3, 18}},
		&Label{name: []string{SUPER + auxiliarCommands[0], SUPER + auxiliarCommands[0], auxiliarCommands[0]}, id: []int{0, 3, 19}},
		&Label{name: []string{SUPER + auxiliarCommands[0], SUPER + auxiliarCommands[0], auxiliarCommands[0]}, id: []int{0, 4, 20}},
		&Label{name: []string{SUPER + auxiliarCommands[0], SUPER + auxiliarCommands[0], auxiliarCommands[0]}, id: []int{0, 4, 21}},
		&Label{name: []string{COMMAND_BUFFER, SUPER + VK_RENDER_PASS, VK_RENDER_PASS, SUPER + VK_SUBPASS, SUPER + VK_SUBPASS, VK_SUBPASS,
			auxiliarCommands[0]}, id: []int{1, 0, 2, 0, 1, 1, 22}},
		&Label{name: []string{COMMAND_BUFFER, SUPER + VK_RENDER_PASS, VK_RENDER_PASS, SUPER + VK_SUBPASS, SUPER + VK_SUBPASS, VK_SUBPASS,
			auxiliarCommands[0]}, id: []int{1, 0, 2, 0, 1, 2, 23}},
		&Label{name: []string{COMMAND_BUFFER, SUPER + VK_RENDER_PASS, VK_RENDER_PASS, SUPER + VK_SUBPASS, SUPER + VK_SUBPASS, VK_SUBPASS,
			auxiliarCommands[0]}, id: []int{1, 0, 2, 0, 2, 3, 24}},
		&Label{name: []string{COMMAND_BUFFER, SUPER + VK_RENDER_PASS, VK_RENDER_PASS, SUPER + VK_SUBPASS, SUPER + VK_SUBPASS, VK_SUBPASS,
			auxiliarCommands[0]}, id: []int{1, 0, 2, 0, 2, 4, 25}},
		&Label{name: []string{COMMAND_BUFFER, SUPER + VK_RENDER_PASS, VK_RENDER_PASS, SUPER + VK_SUBPASS, SUPER + VK_SUBPASS, VK_SUBPASS,
			auxiliarCommands[0]}, id: []int{1, 0, 2, 0, 3, 5, 26}},
		&Label{name: []string{COMMAND_BUFFER, SUPER + VK_RENDER_PASS, VK_RENDER_PASS, SUPER + VK_SUBPASS, SUPER + VK_SUBPASS, VK_SUBPASS,
			auxiliarCommands[0]}, id: []int{1, 0, 2, 0, 3, 6, 27}},
		&Label{name: []string{SUPER + auxiliarCommands[0], SUPER + auxiliarCommands[0], auxiliarCommands[0]}, id: []int{0, 5, 28}},
		&Label{name: []string{SUPER + auxiliarCommands[0], SUPER + auxiliarCommands[0], auxiliarCommands[0]}, id: []int{0, 5, 29}},
	}
	for i := range wantedLabels {
		if reflect.DeepEqual(wantedLabels[i], graph.nodes[i].label) == false {
			t.Errorf("The label for the node %d is different\n", i)
			t.Errorf("Obtained %v Wanted %v\n", graph.nodes[i].label, wantedLabels[i])
		}
	}

}
