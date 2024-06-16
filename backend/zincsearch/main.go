package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type MarkerData struct {
	MarkerID int32  `json:"markerId"`
	Address  string `json:"address"`
}

type BulkData struct {
	Index   string       `json:"index"`
	Records []MarkerData `json:"records"`
}

type PropertyDetail struct {
	Type          string `json:"type"`
	Index         bool   `json:"index"`
	Store         bool   `json:"store"`
	Sortable      bool   `json:"sortable"`
	Aggregatable  bool   `json:"aggregatable"`
	Highlightable bool   `json:"highlightable"`
}

type Mapping struct {
	Properties map[string]PropertyDetail `json:"properties"`
}

type IndexerData struct {
	Name         string  `json:"name"`
	StorageType  string  `json:"storage_type"`
	ShardNum     int     `json:"shard_num"`
	MappingField Mapping `json:"mappings"`
}

var (
	zincApi      string
	zincUser     string
	zincPassword string
)

func main() {
	godotenv.Load()
	zincApi = os.Getenv("ZINCSEARCH_URL")
	zincUser = os.Getenv("ZINCSEARCH_USER")
	zincPassword = os.Getenv("ZINCSEARCH_PASSWORD")

	log.Println("Starting indexer!")
	indexerData, err := createIndexerFromJsonFile("./index.json")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Deleting index if exists...")
	deleted := deleteIndexOnZincSearch("markers")
	if deleted != nil {
		log.Println("Index doesn't exist. Creating...")
	}

	sent := createIndexOnZincSearch(indexerData)
	if sent != nil {
		log.Fatal(sent)
	}

	log.Println("Index created successfully.")
	log.Println("Start indexing, this might take a few minutes...")
	startTime := time.Now()

	records, _ := getMarkersFromJson("markers.json")

	sendBulkToZincSearch(records)

	duration := time.Since(startTime)
	log.Printf("Finished indexing. Time taken: %.2f seconds", duration.Seconds())
}

// https://zincsearch-docs.zinc.dev/api/document/bulkv2/
func sendBulkToZincSearch(records []MarkerData) {
	bulkData := BulkData{
		Index:   "markers",
		Records: records,
	}

	jsonData, err := json.Marshal(bulkData)
	if err != nil {
		log.Println(err)
		return
	}

	req, err := http.NewRequest("POST", zincApi+"/api/_bulkv2", bytes.NewReader(jsonData))
	if err != nil {
		log.Println(err)
		return
	}
	req.SetBasicAuth(zincUser, zincPassword)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()

	_, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}
}

func createIndexerFromJsonFile(filepath string) (IndexerData, error) {
	var indexerData IndexerData

	file, err := os.Open(filepath)
	if err != nil {
		return indexerData, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&indexerData)
	if err != nil {
		return indexerData, err
	}

	return indexerData, nil
}

func getMarkersFromJson(filepath string) ([]MarkerData, error) {
	var markerData []MarkerData

	file, err := os.Open(filepath)
	if err != nil {
		return markerData, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&markerData)
	if err != nil {
		return markerData, err
	}

	return markerData, nil
}

func createIndexOnZincSearch(indexerData IndexerData) error {
	jsonData, err := json.Marshal(indexerData)
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest("POST", zincApi+"/api/index", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(zincUser, zincPassword)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("failed to create indexer, status code: %d", resp.StatusCode)
	}

	return nil
}

func deleteIndexOnZincSearch(indexName string) error {
	req, err := http.NewRequest("DELETE", zincApi+"/api/index/"+indexName, nil)
	if err != nil {
		return err
	}

	req.SetBasicAuth(zincUser, zincPassword)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to delete indexer, status code: %d", resp.StatusCode)
	}

	log.Println("Index deleted successfully")
	return nil
}
