package main

func initBuilding(buildingID string) *Building {
	building := &Building{id: buildingID}
	building.floors = make(map[string]*Floor)

	bdoc, _ := DBGet("buildings", buildingID)

	floors := bdoc["floors"].(map[string]interface{})
	for k, v := range floors {
		floor := building.initFloor(k, v.(map[string]interface{}))
		floor.building = building
		building.floors[k] = floor
	}

	return building
}

func (b *Building) initFloor(floorID string, floordoc map[string]interface{}) *Floor {
	floor := &Floor{id: floorID, building: b}

	graphID := floordoc["graph"].(string)
	floor.graph = floor.InitGraph(graphID)

	return floor
}

var g *Graph
func (f *Floor) InitGraph(graphID string) (*Graph) {
	//load the graph
	doc, _ := DBGet("graphs", graphID)

	nodes := doc["nodes"].([]interface{})
	graph, _ := createGraph(nodes)

	graph.loadDevicesFromDB("beacons")
	graph.loadDevicesFromDB("esp32")

	graph.floor = f
	graph.building = f.building
	go graph.GraphRoutine()
	g = graph

	return graph;
}

func (building *Building) SetEmergency(e bool) {
	building.emergency = e
}

func getMaxExitOccupation() int {
	return 30
}

func getMaxDangerLevel() float64 {
	return 0.25
}
