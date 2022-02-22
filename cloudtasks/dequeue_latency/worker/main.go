package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"
	"cloud.google.com/go/compute/metadata"
	"github.com/google/uuid"
	taskspb "google.golang.org/genproto/googleapis/cloud/tasks/v2"
)

var tasksClient *cloudtasks.Client

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World")
}

func ReceiveTaskHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	fmt.Printf("__latency__,receiveTask,%s,%d\n", string(body), time.Now().UnixNano())
}

func AddTaskHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ip, err := metadata.ExternalIP()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	pID, err := metadata.ProjectID()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	queuePath := fmt.Sprintf("projects/%s/locations/%s/queues/%s", pID, "asia-northeast1", "latency")
	id := uuid.New().String()
	fmt.Printf("__latency__,addTask,%s,%d\n", id, time.Now().UnixNano())
	taskName := fmt.Sprintf("%s/tasks/%s", queuePath, id)
	// Build the Task payload.
	// https://godoc.org/google.golang.org/genproto/googleapis/cloud/tasks/v2#CreateTaskRequest
	req := &taskspb.CreateTaskRequest{
		Parent: queuePath,
		Task: &taskspb.Task{
			Name: taskName,
			// https://godoc.org/google.golang.org/genproto/googleapis/cloud/tasks/v2#HttpRequest
			MessageType: &taskspb.Task_HttpRequest{
				HttpRequest: &taskspb.HttpRequest{
					HttpMethod: taskspb.HttpMethod_POST,
					Url:        fmt.Sprintf("http://%s:8080/tasks/receive", ip),
					Body:       []byte(id),
				},
			},
		},
	}

	// Add a payload message if one is present.
	req.Task.GetHttpRequest().Body = []byte(id)

	createdTask, err := tasksClient.CreateTask(ctx, req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Write([]byte(createdTask.GetName()))
	fmt.Printf("__latency__,createdTask,%s,%d\n", id, time.Now().UnixNano())
}

func main() {
	ctx := context.Background()

	var err error
	tasksClient, err = cloudtasks.NewClient(ctx)
	if err != nil {
		panic(err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("Start Listen PORT:%s\n", port)

	http.HandleFunc("/tasks/receive", ReceiveTaskHandler)
	http.HandleFunc("/tasks/add", AddTaskHandler)
	http.HandleFunc("/", handler)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil); err != nil {
		fmt.Println(err.Error())
	}
}
