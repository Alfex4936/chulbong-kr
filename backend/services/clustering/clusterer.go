package main

import (
	"fmt"
	"math"
)

// Point represents a geospatial point
type Point struct {
	Latitude  float64
	Longitude float64
	ClusterID int // Used to identify which cluster this point belongs to
	Address   string
}

// Haversine formula to calculate geographic distance between points
func distance(lat1, lon1, lat2, lon2 float64) float64 {
	rad := math.Pi / 180
	r := 6371e3 // Earth radius in meters
	dLat := (lat2 - lat1) * rad
	dLon := (lon2 - lon1) * rad
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*rad)*math.Cos(lat2*rad)*
			math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return r * c // Distance in meters
}

func regionQuery(points []Point, p Point, eps float64) []int {
	var neighbors []int
	for idx, pt := range points {
		if distance(p.Latitude, p.Longitude, pt.Latitude, pt.Longitude) < eps {
			neighbors = append(neighbors, idx)
		}
	}
	return neighbors
}

func expandCluster(points []Point, pIndex int, neighbors []int, clusterID int, eps float64, minPts int) {
	points[pIndex].ClusterID = clusterID
	i := 0
	for i < len(neighbors) {
		nIndex := neighbors[i]
		if points[nIndex].ClusterID == 0 { // Not visited
			points[nIndex].ClusterID = clusterID
			nNeighbors := regionQuery(points, points[nIndex], eps)
			if len(nNeighbors) >= minPts {
				neighbors = append(neighbors, nNeighbors...)
			}
		} else if points[nIndex].ClusterID == -1 { // Change noise to border point
			points[nIndex].ClusterID = clusterID
		}
		i++
	}
}

func DBSCAN(points []Point, eps float64, minPts int) []Point {
	clusterID := 0
	for i := range points {
		if points[i].ClusterID != 0 {
			continue
		}
		neighbors := regionQuery(points, points[i], eps)
		if len(neighbors) < minPts {
			points[i].ClusterID = -1
		} else {
			clusterID++
			expandCluster(points, i, neighbors, clusterID, eps, minPts)
		}
	}
	return points
}

func main() {
	points := []Point{
		{Latitude: 37.568166, Longitude: 126.974102, Address: "서울 중구 정동 1-76"},
		{Latitude: 37.568661, Longitude: 126.972375, Address: "서울 종로구 신문로2가 171"},
		{Latitude: 37.56885, Longitude: 126.972064, Address: "서울 종로구 신문로2가 171"},
		{Latitude: 37.56589411615361, Longitude: 126.96930309974685, Address: "서울 중구 순화동 1-1"},

		{Latitude: 37.55808862059195, Longitude: 126.95976545165765, Address: "서울 서대문구 북아현동 884"},
		{Latitude: 37.57838984677184, Longitude: 126.98853202207196, Address: "서울 종로구 원서동 181"},
		{Latitude: 37.57318309415514, Longitude: 126.95501424473001, Address: "서울 서대문구 현저동 101"},
		{Latitude: 37.5541479820707, Longitude: 126.98370331932351, Address: "서울 중구 회현동1가 산 1-2"},
		{Latitude: 37.58411863798303, Longitude: 126.97246285644356, Address: "서울 종로구 궁정동 17-3"},
	}

	eps := 500.0 // Epsilon in meters
	minPts := 2  // Minimum number of points to form a dense region

	clusteredPoints := DBSCAN(points, eps, minPts)
	for _, p := range clusteredPoints {
		fmt.Printf("Point at (%f, %f), Address: %s, Cluster ID: %d\n", p.Latitude, p.Longitude, p.Address, p.ClusterID)
	}
}
