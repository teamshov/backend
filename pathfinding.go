package main

import (
	"encoding/json"
	"io/ioutil"
	"math"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Node struct {
	cNodes []int
	index  int
	x      float64
	y      float64
	ntype  string
}

func createNode(e map[string]interface{}) *Node {
	var n *Node
	n = new(Node)

	cnodes := e["cnodes"].([]interface{})
	n.cNodes = make([]int, len(cnodes))
	for i, elem := range cnodes {
		n.cNodes[i] = int(elem.(float64))
	}

	n.index = int(e["index"].(float64))

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

type Graph struct {
	nodes map[int]*Node
	exits map[int]*Node
}

func createGraph(njson []interface{}) (*Graph, error) {
	var g *Graph
	g = new(Graph)
	g.nodes = make(map[int]*Node)
	g.exits = make(map[int]*Node)

	for _, e := range njson {
		n := createNode(e.(map[string]interface{}))

		g.nodes[n.index] = n
		if n.ntype == "exit" {
			g.exits[n.index] = n
		}
	}

	return g, nil
}

var (
	graph *Graph
)

func GraphRoutine() {

	//load the graph
	doc := DBGet("graphs", "eb2_L1")

	nodes := doc["nodes"].([]interface{})
	graph, _ = createGraph(nodes)

}

func InitPathfindingService(e *echo.Echo) {
	GraphRoutine()
	e.GET("/getpath", getPath)
}

func getPath(c echo.Context) error {
	var input map[string]interface{}
	//var paths map[string]interface{}

	body, _ := ioutil.ReadAll(c.Request().Body)
	json.Unmarshal(body, &input)

	x := input["x"].(float64)
	y := input["y"].(float64)
	//uid := input["uid"]

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

	//add pathfinding here

	return c.JSON(http.StatusOK, minNode.index)
}
