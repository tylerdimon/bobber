package main

import (
	"encoding/json"
	"fmt"
	"go.etcd.io/bbolt"
	"net/http"
	"time"
)

type RequestData struct {
	Method    string
	URL       string
	Host      string
	Path      string
	Timestamp string
	Body      string
	Headers   string
}

func (r RequestData) String() string {
	return fmt.Sprintf("Timestamp: %v\nMethod: %v\nURL: %v\nHost: %v\nPath: %v\nHeaders: %v\nBody: %v",
		time.Now().Format(time.RFC3339), r.Method, r.URL, r.Host, r.Path, r.Headers, r.Body)
}

func all() ([]RequestData, error) {
	var requests []RequestData
	// Open a view transaction to read data
	err := db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("Requests"))
		if b == nil {
			return fmt.Errorf("Bucket not found")
		}

		return b.ForEach(func(k, v []byte) error {
			var reqData RequestData
			if err := json.Unmarshal(v, &reqData); err != nil {
				return err
			}
			requests = append(requests, reqData)
			return nil
		})
	})
	return requests, err
}

func GetAllRequests(w http.ResponseWriter, r *http.Request) {
	requests, err := all()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var strings []string
	for _, req := range requests {
		strings = append(strings, req.String())
	}

	jsonData, err := json.Marshal(strings)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func DeleteAllRequestsHandler(w http.ResponseWriter, r *http.Request) {
	// Clear data in BoltDB
	err := db.Update(func(tx *bbolt.Tx) error {
		err := tx.DeleteBucket([]byte("Requests"))
		if err != nil {
			return err
		}
		_, err = tx.CreateBucket([]byte("Requests"))
		return err
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "All requests cleared")
}

func SaveRequest(data RequestData) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("Requests"))
		return b.Put([]byte(data.Timestamp), jsonData)
	})
	return err
}
