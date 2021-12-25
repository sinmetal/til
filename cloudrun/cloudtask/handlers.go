package main

import (
	"fmt"
	"net/http"
	"time"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"
	taskspb "google.golang.org/genproto/googleapis/cloud/tasks/v2"
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
		},
	}

	_, err := h.tasks.CreateTask(ctx, req)
	if err != nil {
		fmt.Printf("failed cloudtasks.CreateTask: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *Handlers) ServeTask(w http.ResponseWriter, r *http.Request) {
	taskName := r.Header.Get("X-CloudTasks-TaskName")
	taskExecutionCount := r.Header.Get("X-CloudTasks-TaskExecutionCount")
	for i := 0; i < 4000; i++ {
		fmt.Printf("%s:%d:%d\n", taskName, taskExecutionCount, i)
		time.Sleep(1 * time.Second)
	}
}
