package main

import (
	"fmt"
	"math"
	"sync"

	"github.com/dhconnelly/rtreego"
)

// Point represents a geospatial point
type Point struct {
	Latitude  float64
	Longitude float64
	ClusterID int // Used to identify which cluster this point belongs to
	Address   string
}

// RTreePoint wraps a Point to implement the rtreego.Spatial interface
type RTreePoint struct {
	Point
}

func (p RTreePoint) Bounds() rtreego.Rect {
	return rtreego.Point{p.Latitude, p.Longitude}.ToRect(0.00001)
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

// regionQuery returns the indices of all points within eps distance of point p using an R-tree for efficient querying
func regionQuery(tree *rtreego.Rtree, points []Point, p Point, eps float64) []int {
	epsDeg := eps / 111000 // Approximate conversion from meters to degrees
	searchRect := rtreego.Point{p.Latitude, p.Longitude}.ToRect(epsDeg)
	results := tree.SearchIntersect(searchRect)

	var neighbors []int
	for _, item := range results {
		rtp := item.(RTreePoint)
		if distance(p.Latitude, p.Longitude, rtp.Latitude, rtp.Longitude) < eps {
			for idx, pt := range points {
				if pt.Latitude == rtp.Latitude && pt.Longitude == rtp.Longitude {
					neighbors = append(neighbors, idx)
					break
				}
			}
		}
	}
	return neighbors
}

// expandCluster expands the cluster with id clusterID by recursively adding all density-reachable points
func expandCluster(tree *rtreego.Rtree, points []Point, pIndex int, neighbors []int, clusterID int, eps float64, minPts int, wg *sync.WaitGroup, mu *sync.Mutex) {
	defer wg.Done()

	points[pIndex].ClusterID = clusterID
	i := 0
	for i < len(neighbors) {
		nIndex := neighbors[i]
		if points[nIndex].ClusterID == 0 { // Not visited
			points[nIndex].ClusterID = clusterID
			nNeighbors := regionQuery(tree, points, points[nIndex], eps)
			if len(nNeighbors) >= minPts {
				neighbors = append(neighbors, nNeighbors...)
			}
		} else if points[nIndex].ClusterID == -1 { // Change noise to border point
			points[nIndex].ClusterID = clusterID
		}
		i++
	}

	mu.Lock()
	for _, neighborIdx := range neighbors {
		if points[neighborIdx].ClusterID == 0 {
			points[neighborIdx].ClusterID = clusterID
		}
	}
	mu.Unlock()
}

// DBSCAN performs DBSCAN clustering on the points
func DBSCAN(points []Point, eps float64, minPts int) []Point {
	clusterID := 0
	tree := rtreego.NewTree(2, 25, 50) // Create an R-tree for 2D points

	for _, p := range points {
		tree.Insert(RTreePoint{p})
	}

	var wg sync.WaitGroup
	var mu sync.Mutex

	for i := range points {
		if points[i].ClusterID != 0 {
			continue
		}
		neighbors := regionQuery(tree, points, points[i], eps)
		if len(neighbors) < minPts {
			points[i].ClusterID = -1
		} else {
			clusterID++
			wg.Add(1)
			go expandCluster(tree, points, i, neighbors, clusterID, eps, minPts, &wg, &mu)
		}
	}

	wg.Wait()
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
