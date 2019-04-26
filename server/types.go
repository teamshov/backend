package main

type Building struct {
	id        string
	floors    map[string]*Floor
	emergency bool
}

type Floor struct {
	id       string
	graph    *Graph
	devices  *IOTDevice
	building *Building
}

type Graph struct {
	nodes        map[uint]*Node
	exits        []uint
	devices      []*IOTDevice
	coloredNodes map[string]*ColorPath

	building *Building
	floor    *Floor

	red int
}

type Node struct {
	cNodes []uint
	index  uint
	x      float64
	y      float64
	isExit bool
}

type IOTDevice struct {
	deviceType  string
	deviceID    string
	node        *Node
	dangerLevel float64
	x           float64
	y           float64
}

type ColorPath struct {
	origin     *Node
	deviceType string
	deviceID   string
	path       []*Node
}
