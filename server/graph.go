package main

import "math"

type Node struct {
	cNodes []uint
	index  uint
	x      float64
	y      float64
	ntype  string
}

func createNode(e map[string]interface{}) *Node {
	var n *Node
	n = new(Node)

	cnodes := e["cnodes"].([]interface{})
	n.cNodes = make([]uint, len(cnodes))
	for i, elem := range cnodes {
		n.cNodes[i] = uint(elem.(float64))
	}

	n.index = uint(e["index"].(float64))

	pos := e["pos"].(map[string]interface{})
	n.x = pos["x"].(float64)
	n.y = pos["y"].(float64)
	n.ntype = e["type"].(string)

	return n
}

func (n *Node) getDist(x float64, y float64) float64 {
	xx := math.Pow((x - n.x), 2)
	yy := math.Pow((y - n.y), 2)

	return math.Sqrt(xx + yy)
}

func getDist(n1 *Node, n2 *Node) float64 {
	xx := math.Pow(n1.x-n2.x, 2)
	yy := math.Pow(n1.y-n2.y, 2)

	return math.Sqrt(xx + yy)
}

func getNearestNode(x float64, y float64) *Node {
	var minNode *Node
	minDist := math.MaxFloat64
	//get the node
	for _, v := range graph.nodes {
		dist := v.getDist(x, y)

		if dist < minDist || minNode == nil {
			minNode = v
			minDist = dist
		}

	}

	return minNode
}

type Graph struct {
	nodes   map[uint]*Node
	exits   []uint
	devices map[string]map[string]*Node
}

func createGraph(njson []interface{}) (*Graph, error) {
	var g *Graph
	g = new(Graph)
	g.nodes = make(map[uint]*Node)
	g.exits = make([]uint, 0)
	g.devices = make(map[string]map[string]*Node)

	for _, e := range njson {
		n := createNode(e.(map[string]interface{}))

		g.nodes[n.index] = n
		if n.ntype == "exit" {
			g.exits = append(g.exits, n.index)
		}
	}

	return g, nil
}

var (
	graph *Graph
)

func InitGraph() {
	//load the graph
	doc, _ := DBGet("graphs", "eb2_L1")

	nodes := doc["nodes"].([]interface{})
	graph, _ = createGraph(nodes)

	graph.loadDevicesFromDB("beacons")
	graph.loadDevicesFromDB("esp32")
}

func GraphRoutine() {

}

func (g *Graph) loadDevicesFromDB(db string) {
	dlist, err := DBAll(db)
	if err != nil {
		panic(err)
	}

	g.devices[db] = make(map[string]*Node)

	for _, v := range dlist {
		piedoc, _ := DBGet("pies", v)

		x, xok := piedoc["xpos"].(float64)
		y, yok := piedoc["ypos"].(float64)

		if !xok || !yok {
			continue
		}

		n := getNearestNode(x, y)

		g.devices[db][v] = n
	}
}
