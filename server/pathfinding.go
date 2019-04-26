package main

import (
	"errors"
	"math"
)

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

func (graph *Graph) AStar(start *Node, target *Node) ([]float64, float64, error) {
	parentPath := new(Path)
	parentPath.node = start

	paths := make(map[*Path]bool)
	visted := make(map[uint]bool)
	//f(n) = g(n) + h(n)
	getCost := func(n *Node) (float64, float64) {
		h := getDist(n, target)
		g := getDist(n, parentPath.node)
		s := graph.sampleDangerLeveL(n.x, n.y)
		d := g * (1/(1-math.Sqrt(s)) - 1)

		return g + d + h, g + d
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

func (graph *Graph) AStarNodes(start *Node, target *Node) ([]*Node, float64, error) {
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
	nodes := make([]*Node, 0)
	for {
		nodes = append(nodes, parentPath.node)

		if parentPath.parent == nil {
			break
		}

		parentPath = parentPath.parent
	}

	return nodes, cost, nil
}
