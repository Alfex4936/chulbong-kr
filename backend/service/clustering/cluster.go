package main

import (
	"fmt"
	"math"
)

type Point struct {
	Latitude  float64
	Longitude float64
	ClusterID int // 0: unvisited, -1: noise, >0: cluster ID
	Address   string
}

type Grid struct {
	CellSize float64
	Cells    map[int]map[int][]*Point
}

func NewGrid(cellSize float64) *Grid {
	return &Grid{
		CellSize: cellSize,
		Cells:    make(map[int]map[int][]*Point),
	}
}

func (g *Grid) Insert(p *Point) {
	xIdx := int(p.Longitude / g.CellSize)
	yIdx := int(p.Latitude / g.CellSize)

	if _, ok := g.Cells[xIdx]; !ok {
		g.Cells[xIdx] = make(map[int][]*Point)
	}
	g.Cells[xIdx][yIdx] = append(g.Cells[xIdx][yIdx], p)
}

func (g *Grid) GetNeighbors(p *Point, eps float64) []*Point {
	epsDeg := eps / 111000.0 // Convert meters to degrees
	cellRadius := int(math.Ceil(epsDeg / g.CellSize))

	xIdx := int(p.Longitude / g.CellSize)
	yIdx := int(p.Latitude / g.CellSize)

	neighbors := []*Point{}

	for dx := -cellRadius; dx <= cellRadius; dx++ {
		for dy := -cellRadius; dy <= cellRadius; dy++ {
			nx := xIdx + dx
			ny := yIdx + dy

			if cell, ok := g.Cells[nx][ny]; ok {
				for _, np := range cell {
					if distance(p.Latitude, p.Longitude, np.Latitude, np.Longitude) <= eps {
						neighbors = append(neighbors, np)
					}
				}
			}
		}
	}
	return neighbors
}

// Optimized Haversine formula for small distances
func distance(lat1, lon1, lat2, lon2 float64) float64 {
	const (
		rad = math.Pi / 180
		r   = 6371e3 // Earth radius in meters
	)
	dLat := (lat2 - lat1) * rad
	dLon := (lon2 - lon1) * rad
	lat1Rad := lat1 * rad
	lat2Rad := lat2 * rad

	sinDLat := math.Sin(dLat / 2)
	sinDLon := math.Sin(dLon / 2)

	a := sinDLat*sinDLat + math.Cos(lat1Rad)*math.Cos(lat2Rad)*sinDLon*sinDLon
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return r * c // Distance in meters
}

// expandCluster expands the cluster with id clusterID by recursively adding all density-reachable points
func expandCluster(grid *Grid, p *Point, neighbors []*Point, clusterID int, eps float64, minPts int) {
	p.ClusterID = clusterID

	queue := make([]*Point, 0, len(neighbors))
	queue = append(queue, neighbors...)

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if current.ClusterID == 0 { // Unvisited
			current.ClusterID = clusterID
			nNeighbors := grid.GetNeighbors(current, eps)
			if len(nNeighbors) >= minPts {
				queue = append(queue, nNeighbors...)
			}
		} else if current.ClusterID == -1 { // Noise
			current.ClusterID = clusterID
		}
	}
}

// DBSCAN performs DBSCAN clustering on the points
func DBSCAN(points []*Point, eps float64, minPts int) {
	clusterID := 0
	epsDeg := eps / 111000.0 // Convert meters to degrees
	cellSize := epsDeg       // Use epsDeg as cell size

	grid := NewGrid(cellSize)

	for _, p := range points {
		grid.Insert(p)
	}

	for _, p := range points {
		if p.ClusterID != 0 {
			continue
		}
		neighbors := grid.GetNeighbors(p, eps)
		if len(neighbors) < minPts {
			p.ClusterID = -1 // Mark as noise
		} else {
			clusterID++
			expandCluster(grid, p, neighbors, clusterID, eps, minPts)
		}
	}
}

func main() {
	points := []*Point{
		{Latitude: 37.55808862059195, Longitude: 126.95976545165765, Address: "서울 서대문구 북아현동 884"},
		{Latitude: 37.568166, Longitude: 126.974102, Address: "서울 중구 정동 1-76"},
		{Latitude: 37.568661, Longitude: 126.972375, Address: "서울 종로구 신문로2가 171"},
		{Latitude: 37.56885, Longitude: 126.972064, Address: "서울 종로구 신문로2가 171"},
		{Latitude: 37.56589411615361, Longitude: 126.96930309974685, Address: "서울 중구 순화동 1-1"},
		{Latitude: 37.57838984677184, Longitude: 126.98853202207196, Address: "서울 종로구 원서동 181"},
		{Latitude: 37.57318309415514, Longitude: 126.95501424473001, Address: "서울 서대문구 현저동 101"},
		{Latitude: 37.5541479820707, Longitude: 126.98370331932351, Address: "서울 중구 회현동1가 산 1-2"},
		{Latitude: 37.58411863798303, Longitude: 126.97246285644356, Address: "서울 종로구 궁정동 17-3"},
		{Latitude: 36.33937565888829, Longitude: 127.41575408006757, Address: "대전 중구 선화동 223"},
		{Latitude: 36.346176003613984, Longitude: 127.41482385609581, Address: "대전 대덕구 오정동 496-1"},
	}

	eps := 5000.0 // Epsilon in meters
	minPts := 2   // Minimum number of points to form a dense region

	// Will group points within a region of 500 meters
	DBSCAN(points, eps, minPts)

	for _, p := range points {
		fmt.Printf("Point at (%f, %f), Address: %s, Cluster ID: %d\n", p.Latitude, p.Longitude, p.Address, p.ClusterID)
	}
}
