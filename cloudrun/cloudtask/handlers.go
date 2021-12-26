package main

import (
	"fmt"
	"net/http"
	"time"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"
	taskspb "google.golang.org/genproto/googleapis/cloud/tasks/v2"
	"google.golang.org/protobuf/types/known/durationpb"
)

type Handlers struct {
	tasks *cloudtasks.Client
}

func (h *Handlers) AddTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	r.FormValue("")

	// Build the Task queue path.
	queuePath := fmt.Sprintf("projects/%s/locations/%s/queues/%s", "sinmetal-playground-20211225", "asia-northeast1", "test")

	// Build the Task payload.
	// https://godoc.org/google.golang.org/genproto/googleapis/cloud/tasks/v2#CreateTaskRequest
	req := &taskspb.CreateTaskRequest{
		Parent: queuePath,
		Task: &taskspb.Task{
			// https://godoc.org/google.golang.org/genproto/googleapis/cloud/tasks/v2#HttpRequest
			MessageType: &taskspb.Task_HttpRequest{
				HttpRequest: &taskspb.HttpRequest{
					HttpMethod: taskspb.HttpMethod_POST,
					Url:        "https://appserver-2jsu5stp3a-an.a.run.app/tasks/serve",
					AuthorizationHeader: &taskspb.HttpRequest_OidcToken{
						OidcToken: &taskspb.OidcToken{
							ServiceAccountEmail: "run-default@sinmetal-playground-20211225.iam.gserviceaccount.com",
						},
					},
				},
			},
			DispatchDeadline: &durationpb.Duration{
				Seconds: 1800,
				Nanos:   0,
			},
		},
	}

	task, err := h.tasks.CreateTask(ctx, req)
	if err != nil {
		fmt.Printf("failed cloudtasks.CreateTask: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(task.Name))
	if err != nil {
		fmt.Println(err)
	}
}

func (h *Handlers) ServeTask(w http.ResponseWriter, r *http.Request) {
	taskName := r.Header.Get("X-CloudTasks-TaskName")
	taskExecutionCount := r.Header.Get("X-CloudTasks-TaskExecutionCount")
	taskRetryCount := r.Header.Get("X-CloudTasks-TaskRetryCount")
	for i := 0; i < 4000; i++ {
		fmt.Printf("timeCount-v2:%s:%s:%s:%d\n", taskName, taskExecutionCount, taskRetryCount, i)
		time.Sleep(1 * time.Second)
	}
	fmt.Printf("timeCount-v2:%s:%s:%s:%d\n", taskName, taskExecutionCount, taskRetryCount, 9999)
	w.WriteHeader(http.StatusOK)
}
