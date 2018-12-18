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
		lines := fmt.Sprintf("%d", node.id) + "[label=" + "\"" + node.label.getLabelAsString() + "\"" + "]" + ";\n"
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
	var output bytes.Buffer
	totalRealNodes := 0
	totalRealEdges := 0
	for _, node := range g.nodes {
		if node.isReal == false {
			continue
		}
		totalRealNodes++
		output.WriteString("node {\n")
		output.WriteString("\tname: \"" + node.label.getLabelAsString() + "\"\n")
		output.WriteString("\top: \"" + node.label.name[len(node.label.name)-1] + fmt.Sprintf("%d", node.id) + "\"\n")

		for _, nodeNeighbor := range node.incomingNodes {
			if nodeNeighbor.isReal == false {
				continue
			}
			totalRealEdges++
			output.WriteString("\tinput: \"" + nodeNeighbor.label.getLabelAsString() + "\"\n")
		}
		output.WriteString("\tdevice: \"" + node.color + "\"\n")

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
	fmt.Println("Displayed Nodes = ", totalRealNodes, " Displayed Edges = ", totalRealEdges)
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
	protoNode.Name = node.label.getLabelAsString()
	protoNode.Op = node.label.getLabelAsString()

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
		protoNode.Input[i] = g.idNodeToNode[idNeighbor].label.getLabelAsString()
	}
	return protoNode
}
