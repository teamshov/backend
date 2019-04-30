package main

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"time"
)

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

	ntype := e["type"].(string)
	if ntype == "exit" {
		n.isExit = true
	}

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

func createGraph(njson []interface{}) (*Graph, error) {
	var g *Graph
	g = new(Graph)
	g.nodes = make(map[uint]*Node)
	g.exits = make([]uint, 0)
	g.devices = make([]*IOTDevice, 0)
	g.coloredNodes = make(map[string]*ColorPath)

	for _, e := range njson {
		n := createNode(e.(map[string]interface{}))

		g.nodes[n.index] = n
		if n.isExit {
			g.exits = append(g.exits, n.index)
		}
	}

	return g, nil
}

func (g *Graph) loadDevicesFromDB(db string) {
	dlist, err := DBAll(db)
	if err != nil {
		panic(err)
	}

	for _, v := range dlist {
		doc, err := DBGet(db, v)
		if err != nil {
			fmt.Println(err)
			continue
		}

		x, xok := doc["xpos"].(float64)
		y, yok := doc["ypos"].(float64)

		if !xok || !yok {
			continue
		}

		n := g.getNearestNode(x, y)

		g.devices = append(g.devices, &IOTDevice{deviceType: db, deviceID: v, node: n, x: x, y: y})

		_, cok := doc["color"].(bool)
		if cok {
			g.coloredNodes[v] = g.createColorPath(n, db, v)
		}
	}
}

func (graph *Graph) createColorPath(node *Node, deviceType string, deviceID string) *ColorPath {
	//add pathfinding here
	minCost := math.MaxFloat64
	var minResult []*Node

	for _, eindex := range graph.exits {
		result, cost, err := graph.AStarNodes(node, graph.nodes[eindex])
		if cost < minCost {
			minCost = cost
			minResult = result
		}

		if err != nil {
			fmt.Println(err)
		}
	}

	colorpath := &ColorPath{
		origin:     node,
		deviceType: deviceType,
		deviceID:   deviceID,
		path:       minResult,
	}

	return colorpath
}

func (graph *Graph) GraphRoutine() {
	for {
		graph.red = 0
		for _, v := range graph.devices {
			v.updateDangerLevel()
		}

		for _, v := range graph.coloredNodes {
			color := graph.getColorOfPath(v)
			PublishColor(v.deviceType, v.deviceID, color)
		}

		graph.building.SetEmergency(graph.red > 0)

		time.Sleep(2 * time.Second)
	}
}

func (graph *Graph) getNearestNode(x float64, y float64) *Node {
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

func (graph *Graph) getColorOfPath(colorpath *ColorPath) string {
	var isCongested bool
	var isDangerous bool

	for _, e := range colorpath.path {
		if e.isExit {
			congestion, _ := RedisGetInt(fmt.Sprintf("%s:%s:exit:%i", graph.building.id, graph.floor.id, e.index))

			if congestion > getMaxExitOccupation() {
				isCongested = true
			}
		}
		dl := graph.sampleDangerLeveL(e.x, e.y)
		if dl > getMaxDangerLevel() {
			isDangerous = true
			break
		}
	}

	if isDangerous {
		graph.red++
		return "red"
	} else if isCongested {
		return "yellow"
	} else {
		return "green"
	}
}

func (d *IOTDevice) updateDangerLevel() {
	topic := fmt.Sprintf("%s:%s:dangerlevel", d.deviceType, d.deviceID)
	dls, err := RedisGetString(topic)
	if err != nil {
		//fmt.Printf("%v %v %f\n", err, topic, dls)
	}

	dl, err := strconv.ParseFloat(dls, 64)
	if err != nil {
		dl = 0.0
	}
	d.dangerLevel = dl
}

func (g *Graph) sampleDangerLeveL(x float64, y float64) float64 {
	var max float64
	for _, e := range g.devices {
		dist := math.Sqrt(math.Pow(x-e.x, 2) + math.Pow(y-e.y, 2))
		d := e.dangerLevel * math.Pow(math.E, -math.Pow(dist, 2)/math.Pow(5, 2))
		if d > max {
			max = d
		}
	}

	return max
}

func (graph *Graph) getPathXY(x float64, y float64) string {

	node := graph.getNearestNode(x, y)

	//add pathfinding here
	minCost := math.MaxFloat64
	var minResult []float64

	for _, eindex := range graph.exits {
		result, cost, err := graph.AStar(node, graph.nodes[eindex])
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
