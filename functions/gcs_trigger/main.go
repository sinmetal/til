package gcs_trigger

import (
	"context"
	"fmt"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/cloudevents/sdk-go/v2/event"
)

func init() {
	// Register a CloudEvent function with the Functions Framework
	functions.CloudEvent("GCSTriggerFunction", gcsTriggerEventFunction)
}

// Function gcsTriggerEventFunction accepts and handles a CloudEvent object
func gcsTriggerEventFunction(ctx context.Context, e event.Event) error {
	// Your code here
	// Access the CloudEvent data payload via e.Data() or e.DataAs(...)
	fmt.Printf("%s\n", string(e.Data()))

	// Return nil if no error occurred
	return nil
}
