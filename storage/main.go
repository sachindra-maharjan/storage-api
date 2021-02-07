package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"github.com/gorilla/mux"
	"google.golang.org/api/iterator"
)

//Response API response
type Response struct {
	ProjectID  string `json:"projectID,omitempty"`
	BucketName string `json:"bucketName,omitempty"`
	Message    string `json:"message,omitempty"`
}

//Status of the response
type Status struct {
	Status string `json:"status,omitempty"`
}

func status(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(Status{Status: "Storage service is running!!!"})
}

func createBucket(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectID := vars["projectID"]
	bucketName := vars["bucketName"]

	err := createBucketClassLocation(projectID, bucketName)
	if err != nil {
		json.NewEncoder(w).Encode(Response{ProjectID: projectID, BucketName: bucketName, Message: err.Error()})
	} else {
		json.NewEncoder(w).Encode(Response{ProjectID: projectID, BucketName: bucketName, Message: "Bucket created"})
	}
}

func getBuckets(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectID := vars["projectID"]
	buckets, err := list(projectID)

	if err != nil {
		json.NewEncoder(w).Encode(Response{ProjectID: projectID, Message: err.Error()})
	} else {
		json.NewEncoder(w).Encode(buckets)
	}
}

func deleteBucket(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bucketName := vars["bucketName"]
	err := delete(bucketName)
	if err != nil {
		json.NewEncoder(w).Encode(Response{BucketName: bucketName, Message: err.Error()})
	} else {
		json.NewEncoder(w).Encode(Response{BucketName: bucketName, Message: "Bucket deleted"})
	}
}

func uploadFile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bucketName := vars["bucketName"]
	objectName := r.URL.Query().Get("object")
	var message = "File uploaded"

	//Parse multipart form, 10 << 20 specific maximum size of 10 MB
	r.ParseMultipartForm(10 << 20)

	file, handler, err := r.FormFile("file")
	if err != nil {
		json.NewEncoder(w).Encode(Response{BucketName: bucketName, Message: err.Error()})
		return
	}
	defer file.Close()

	fmt.Printf("Uploaded file: %+v\n", handler.Filename)
	fmt.Printf("File size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		json.NewEncoder(w).Encode(Response{BucketName: bucketName, Message: err.Error()})
		return
	}

	object := handler.Filename
	if len(strings.TrimSpace(object)) != 0 {
		object = objectName + "/" + handler.Filename
	}

	err = uploadObject(fileBytes, bucketName, object)
	if err != nil {
		json.NewEncoder(w).Encode(Response{BucketName: bucketName, Message: err.Error()})
		return
	}
	json.NewEncoder(w).Encode(Response{BucketName: bucketName, Message: message})
}

func downloadFile(w http.ResponseWriter, r *http.Request) {
	fmt.Print("File downloaded")
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", status).Methods("GET")
	router.HandleFunc("/bucket/{projectID}/{bucketName}", createBucket).Methods("POST")
	router.HandleFunc("/bucket/list/{projectID}", getBuckets).Methods("GET")
	router.HandleFunc("/bucket/{bucketName}", deleteBucket).Methods("DELETE")
	router.HandleFunc("/bucket/object/upload/{bucketName}", uploadFile).Methods("POST")
	router.HandleFunc("/bucket/object/download/{bucketName}/{object}", downloadFile).Methods("GET")

	// [START setting_port]
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatal(err)
	}
}

//Bucket Operations
func createBucketClassLocation(projectID, bucketName string) error {
	ctx := context.Background()

	if projectID == "" {
		return fmt.Errorf("Project ID is required")
	}
	if bucketName == "" {
		return fmt.Errorf("BucketName is required")
	}

	client, err := storage.NewClient(ctx)

	if err != nil {
		return fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	storageClassAndLocation := &storage.BucketAttrs{
		StorageClass: "STANDARD",
		Location:     "US",
		LocationType: "multi-region",
	}

	bucket := client.Bucket(bucketName)
	if err := bucket.Create(ctx, projectID, storageClassAndLocation); err != nil {
		return fmt.Errorf("Bucket(%q).Create: %v", bucketName, err)
	}
	fmt.Printf("Created bucket %v in %v with storage class %v\n", bucketName, storageClassAndLocation.Location, storageClassAndLocation.StorageClass)
	return nil
}

func list(projectID string) ([]string, error) {
	ctx := context.Background()

	if projectID == "" {
		return nil, fmt.Errorf("Project ID is required")
	}
	var buckets []string
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	it := client.Buckets(ctx, projectID)

	for {
		battrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		buckets = append(buckets, battrs.Name)
	}
	return buckets, nil
}

func delete(bucketName string) error {
	ctx := context.Background()

	if bucketName == "" {
		return fmt.Errorf("BucketName is required")
	}

	client, err := storage.NewClient(ctx)

	if err != nil {
		return fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	if err := client.Bucket(bucketName).Delete(ctx); err != nil {
		return fmt.Errorf("Bucket(%q).Delete: %v", bucketName, err)
	}
	return nil
}

func uploadObject(buf []byte, bucketName, object string) error {
	ctx := context.Background()

	if bucketName == "" {
		return fmt.Errorf("BucketName is required")
	}

	client, err := storage.NewClient(ctx)

	if err != nil {
		return fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	wc := client.Bucket(bucketName).Object(object).NewWriter(ctx)
	if _, err := io.Copy(wc, bytes.NewReader(buf)); err != nil {
		return fmt.Errorf("io.Copy: %v", err)
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %v", err)
	}
	fmt.Printf("Blob %v uploaded", object)
	return nil
}
