package graph_visualization

import (
	"bytes"
	"fmt"
	"sort"
)

const (
	QUEUE_PRESENT      = "vkQueuePresentKHR"
	SUPER_COMMAND      = "SUPER"
	UNUSED_COMMAND     = "UNUSED"
	LIMIT_BY_HIERARCHY = 5
)

type NodeSorter []*Node

func (input NodeSorter) Len() int {
	return len(input)
}
func (input NodeSorter) Swap(i, j int) {
	input[i], input[j] = input[j], input[i]
}
func (input NodeSorter) Less(i, j int) bool {
	return input[i].id < input[j].id
}

func (g *Graph) assignColors() {
	for _, node := range g.nodes {
		node.color = node.label.name[0] + fmt.Sprintf("%d", node.label.id[0])
	}
}

func (g *Graph) getNumberNodesInTopLevel() int {
	uniquesNamesInTopLevel := map[string]int{}
	for _, node := range g.nodes {
		nameTopLevel := node.label.name[0] + fmt.Sprintf("%d", node.label.id[0])
		uniquesNamesInTopLevel[nameTopLevel] = 1
	}
	return len(uniquesNamesInTopLevel)
}

type Chunk struct {
	idToPosition       map[int]int
	positionToLabelIds [][]int
	done               bool
}

func assignLabelIdsToChunk(begin, end int, currentLabelIds *[]int, positionToLabelIds *[][]int) {
	if (end - begin + 1) <= LIMIT_BY_HIERARCHY {
		for i := begin; i <= end; i++ {
			(*positionToLabelIds)[i] = make([]int, len(*currentLabelIds))
			copy((*positionToLabelIds)[i], *currentLabelIds)
		}
	} else {
		*currentLabelIds = append(*currentLabelIds, 1)
		size := (end - begin) / LIMIT_BY_HIERARCHY
		newBegin := begin
		newEnd := newBegin + size
		id := 1
		for newBegin <= end {
			if newEnd > end {
				newEnd = end
			}
			(*currentLabelIds)[len(*currentLabelIds)-1] = id
			assignLabelIdsToChunk(newBegin, newEnd, currentLabelIds, positionToLabelIds)
			id++
			newBegin = newEnd + 1
			newEnd = newBegin + size
		}
		*currentLabelIds = (*currentLabelIds)[:len(*currentLabelIds)-1]
	}
}

func makeChunks(nodes *[]*Node) {
	nameToChunk := map[string]*Chunk{}
	for _, node := range *nodes {
		var currentName bytes.Buffer
		label := node.label
		for i, name := range label.name {
			id := label.id[i]
			currentName.WriteString(name)
			if _, ok := nameToChunk[currentName.String()]; ok == false {
				nameToChunk[currentName.String()] = &Chunk{idToPosition: map[int]int{}}
			}
			currentChunk := nameToChunk[currentName.String()]
			if _, ok := currentChunk.idToPosition[id]; ok == false {
				size := len(currentChunk.idToPosition)
				currentChunk.idToPosition[id] = size
			}
			currentName.WriteString(fmt.Sprintf("%d/", id))
		}
	}

	for _, node := range *nodes {
		var currentName bytes.Buffer
		label := node.label
		newLabel := &Label{}
		for i, name := range label.name {
			id := label.id[i]
			currentName.WriteString(name)
			currentChunk := nameToChunk[currentName.String()]

			size := len(currentChunk.idToPosition)
			if size > 1 {
				if currentChunk.done == false {
					fmt.Println("CreatingChunk of size ", size, " with name ", currentName.String())
					currentChunk.done = true
					currentChunk.positionToLabelIds = make([][]int, size)
					currentLabelIds := make([]int, 1)
					assignLabelIdsToChunk(0, size-1, &currentLabelIds, &currentChunk.positionToLabelIds)
				}
				pos := currentChunk.idToPosition[id]
				newIds := currentChunk.positionToLabelIds[pos]
				for _, newId := range newIds {
					newLabel.pushBack(SUPER_COMMAND+name, newId)
				}
			}
			newLabel.pushBack(name, id)
			currentName.WriteString(fmt.Sprintf("%d/", id))
		}
		node.label = newLabel
	}
}

func (g *Graph) dfs(curr *Node, visited *[]bool, nodesVisited *[]*Node) {
	*nodesVisited = append(*nodesVisited, curr)
	(*visited)[curr.id] = true

	for _, neighbor := range curr.outcomingNodes {
		if (*visited)[neighbor.id] == false && neighbor.isReal {
			g.dfs(neighbor, visited, nodesVisited)
		}
	}
}

func (g *Graph) makeChunksByFrame() {

	visited := make([]bool, g.maxIdNode)
	numFrame := 1
	for _, currNode := range g.nodes {
		if currNode.name == QUEUE_PRESENT && visited[currNode.id] == false && currNode.isReal {
			nodesVisited := []*Node{}
			fmt.Printf("Starting Frame%d\n", numFrame)
			g.dfs(currNode, &visited, &nodesVisited)

			sort.Sort(NodeSorter(nodesVisited))
			makeChunks(&nodesVisited)

			for _, node := range nodesVisited {
				node.label.pushFront("Frame", numFrame)
			}
			numFrame++
		}
	}
	for _, currNode := range g.nodes {
		if visited[currNode.id] == false {
			currNode.label.pushFront(UNUSED_COMMAND, 0)
			currNode.color = ""
		}
	}
}
