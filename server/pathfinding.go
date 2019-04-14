package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"math"
	"net/http"

	"github.com/labstack/echo/v4"
)

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

type Graph struct {
	nodes map[uint]*Node
	exits []uint
}

func createGraph(njson []interface{}) (*Graph, error) {
	var g *Graph
	g = new(Graph)
	g.nodes = make(map[uint]*Node)
	g.exits = make([]uint, 0)

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

func GraphRoutine() {

	//load the graph
	doc := DBGet("graphs", "eb2_L1")

	nodes := doc["nodes"].([]interface{})
	graph, _ = createGraph(nodes)

}

func InitPathfindingService(e *echo.Echo) {
	GraphRoutine()
	e.PUT("/getpath", getPath)
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
	minCost := math.MaxFloat64
	var minResult []float64

	for _, eindex := range graph.exits {
		result, cost, err := AStar(graph, minNode, graph.nodes[eindex])
		if cost < minCost {
			minCost = cost
			minResult = result
		}

		if err != nil {
			panic(err)
		}
	}

	json, _ := json.Marshal(minResult)
	return c.JSON(http.StatusOK, string(json))
}

type Path struct {
	cost   float64
	acculCost float64
	parent *Path
	node   *Node
}

type Vec2 struct {
	X float64
	Y float64
}

func AStar(graph *Graph, start *Node, target *Node) ([]float64, float64, error) {
	parentPath := new(Path)
	parentPath.node = start

	paths := make(map[*Path]bool)
	visted := make(map[uint]bool)
	//f(n) = g(n) + h(n)
	getCost := func(n *Node) (float64, float64) {
		h := getDist(n, target)
		g := getDist(n, parentPath.node)

		return g + h, g
	}

	for parentPath.node != target {
		pnode := parentPath.node
		var minPath *Path
		minCost := math.MaxFloat64
		for _, e := range pnode.cNodes {
			var g float64
			if _, ok := visted[e]; ok {
				continue
			}
			visted[e] = true

			path := new(Path)
			path.parent = parentPath
			path.node = graph.nodes[e]
			path.cost, g = getCost(path.node)
			path.acculCost = parentPath.acculCost + g
			paths[path] = true
		}

		if len(paths) == 0 {
			return nil, 0, errors.New("failed to find a path to exit")
		}

		for p := range paths {
			if p.cost < minCost {
				minPath = p
				minCost = p.cost
			}
		}

		delete(paths, minPath)
		parentPath = minPath
	}
	cost := parentPath.acculCost
	pathsindices := make([]float64, 0)
	for {
		n := parentPath.node
		pathsindices = append(pathsindices, n.x)
		pathsindices = append(pathsindices, n.y)

		if parentPath.parent == nil {
			break
		}

		parentPath = parentPath.parent
	}

	return pathsindices, cost, nil
}
