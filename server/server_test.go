package main

import (
	"math/rand"
	"testing"
)

func TestPathfinding(t *testing.T) {
	//load the graph
	doc := DBGet("graphs", "eb2_L1")

	nodes := doc["nodes"].([]interface{})
	graph, _ = createGraph(nodes)

	for i := 0; i < 10; i++ {
		getPathXY(rand.Float64()*40, rand.Float64()*40)
	}
}

func BenchmarkPathfinding(b *testing.B) {
	//load the graph
	doc := DBGet("graphs", "eb2_L1")

	nodes := doc["nodes"].([]interface{})
	graph, _ = createGraph(nodes)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		getPathXY(rand.Float64()*40, rand.Float64()*40)
	}
}

func BenchmarkPathfindingParrallel(b *testing.B) {
	//load the graph
	doc := DBGet("graphs", "eb2_L1")

	nodes := doc["nodes"].([]interface{})
	graph, _ = createGraph(nodes)
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			getPathXY(rand.Float64()*40, rand.Float64()*40)
		}
	})
}
