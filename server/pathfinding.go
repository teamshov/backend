package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"math"
	"net/http"

	"github.com/labstack/echo/v4"
)

func InitPathfindingService(e *echo.Echo) {
	InitGraph()
	e.PUT("/getpath", getPath)
}

func getPathXY(x float64, y float64) string {

	node := getNearestNode(x, y)

	//add pathfinding here
	minCost := math.MaxFloat64
	var minResult []float64

	for _, eindex := range graph.exits {
		result, cost, err := AStar(graph, node, graph.nodes[eindex])
		if cost < minCost {
			minCost = cost
			minResult = result
		}

		if err != nil {
			panic(err)
		}
	}

	json, _ := json.Marshal(minResult)
	return string(json)
}

func getPath(c echo.Context) error {
	var input map[string]interface{}
	//var paths map[string]interface{}

	body, _ := ioutil.ReadAll(c.Request().Body)
	json.Unmarshal(body, &input)

	x := input["x"].(float64)
	y := input["y"].(float64)
	//uid := input["uid"]

	return c.JSON(http.StatusOK, getPathXY(x, y))
}

type Path struct {
	cost      float64
	acculCost float64
	parent    *Path
	node      *Node
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
