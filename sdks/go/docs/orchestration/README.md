# Orchestration

### Available Operations

* [CancelEvent](#cancelevent) - Cancel a running workflow
* [CreateWorkflow](#createworkflow) - Create workflow
* [GetInstance](#getinstance) - Get a workflow instance by id
* [GetInstanceHistory](#getinstancehistory) - Get a workflow instance history by id
* [GetInstanceStageHistory](#getinstancestagehistory) - Get a workflow instance stage history
* [GetWorkflow](#getworkflow) - Get a flow by id
* [ListInstances](#listinstances) - List instances of a workflow
* [ListWorkflows](#listworkflows) - List registered workflows
* [OrchestrationgetServerInfo](#orchestrationgetserverinfo) - Get server info
* [RunWorkflow](#runworkflow) - Run workflow
* [SendEvent](#sendevent) - Send an event to a running workflow

## CancelEvent

Cancel a running workflow

### Example Usage

```go
package main

import(
	"context"
	"log"
	"github.com/formancehq/formance-sdk-go"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
)

func main() {
    s := formance.New(
        formance.WithSecurity(shared.Security{
            Authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
        }),
    )

    ctx := context.Background()
    res, err := s.Orchestration.CancelEvent(ctx, operations.CancelEventRequest{
        InstanceID: "dolorum",
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.StatusCode == http.StatusOK {
        // handle response
    }
}
```

## CreateWorkflow

Create a workflow

### Example Usage

```go
package main

import(
	"context"
	"log"
	"github.com/formancehq/formance-sdk-go"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
)

func main() {
    s := formance.New(
        formance.WithSecurity(shared.Security{
            Authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
        }),
    )

    ctx := context.Background()
    res, err := s.Orchestration.CreateWorkflow(ctx, shared.CreateWorkflowRequest{
        Name: formance.String("Eddie Prosacco"),
        Stages: []map[string]interface{}{
            map[string]interface{}{
                "non": "eligendi",
                "sint": "aliquid",
            },
            map[string]interface{}{
                "necessitatibus": "sint",
                "officia": "dolor",
                "debitis": "a",
            },
            map[string]interface{}{
                "in": "in",
                "illum": "maiores",
                "rerum": "dicta",
            },
            map[string]interface{}{
                "cumque": "facere",
                "ea": "aliquid",
            },
        },
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.CreateWorkflowResponse != nil {
        // handle response
    }
}
```

## GetInstance

Get a workflow instance by id

### Example Usage

```go
package main

import(
	"context"
	"log"
	"github.com/formancehq/formance-sdk-go"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
)

func main() {
    s := formance.New(
        formance.WithSecurity(shared.Security{
            Authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
        }),
    )

    ctx := context.Background()
    res, err := s.Orchestration.GetInstance(ctx, operations.GetInstanceRequest{
        InstanceID: "laborum",
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.GetWorkflowInstanceResponse != nil {
        // handle response
    }
}
```

## GetInstanceHistory

Get a workflow instance history by id

### Example Usage

```go
package main

import(
	"context"
	"log"
	"github.com/formancehq/formance-sdk-go"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
)

func main() {
    s := formance.New(
        formance.WithSecurity(shared.Security{
            Authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
        }),
    )

    ctx := context.Background()
    res, err := s.Orchestration.GetInstanceHistory(ctx, operations.GetInstanceHistoryRequest{
        InstanceID: "accusamus",
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.GetWorkflowInstanceHistoryResponse != nil {
        // handle response
    }
}
```

## GetInstanceStageHistory

Get a workflow instance stage history

### Example Usage

```go
package main

import(
	"context"
	"log"
	"github.com/formancehq/formance-sdk-go"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
)

func main() {
    s := formance.New(
        formance.WithSecurity(shared.Security{
            Authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
        }),
    )

    ctx := context.Background()
    res, err := s.Orchestration.GetInstanceStageHistory(ctx, operations.GetInstanceStageHistoryRequest{
        InstanceID: "non",
        Number: 581273,
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.GetWorkflowInstanceHistoryStageResponse != nil {
        // handle response
    }
}
```

## GetWorkflow

Get a flow by id

### Example Usage

```go
package main

import(
	"context"
	"log"
	"github.com/formancehq/formance-sdk-go"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
)

func main() {
    s := formance.New(
        formance.WithSecurity(shared.Security{
            Authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
        }),
    )

    ctx := context.Background()
    res, err := s.Orchestration.GetWorkflow(ctx, operations.GetWorkflowRequest{
        FlowID: "enim",
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.GetWorkflowResponse != nil {
        // handle response
    }
}
```

## ListInstances

List instances of a workflow

### Example Usage

```go
package main

import(
	"context"
	"log"
	"github.com/formancehq/formance-sdk-go"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
)

func main() {
    s := formance.New(
        formance.WithSecurity(shared.Security{
            Authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
        }),
    )

    ctx := context.Background()
    res, err := s.Orchestration.ListInstances(ctx, operations.ListInstancesRequest{
        Running: formance.Bool(false),
        WorkflowID: formance.String("accusamus"),
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.ListRunsResponse != nil {
        // handle response
    }
}
```

## ListWorkflows

List registered workflows

### Example Usage

```go
package main

import(
	"context"
	"log"
	"github.com/formancehq/formance-sdk-go"
)

func main() {
    s := formance.New(
        formance.WithSecurity(shared.Security{
            Authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
        }),
    )

    ctx := context.Background()
    res, err := s.Orchestration.ListWorkflows(ctx)
    if err != nil {
        log.Fatal(err)
    }

    if res.ListWorkflowsResponse != nil {
        // handle response
    }
}
```

## OrchestrationgetServerInfo

Get server info

### Example Usage

```go
package main

import(
	"context"
	"log"
	"github.com/formancehq/formance-sdk-go"
)

func main() {
    s := formance.New(
        formance.WithSecurity(shared.Security{
            Authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
        }),
    )

    ctx := context.Background()
    res, err := s.Orchestration.OrchestrationgetServerInfo(ctx)
    if err != nil {
        log.Fatal(err)
    }

    if res.ServerInfo != nil {
        // handle response
    }
}
```

## RunWorkflow

Run workflow

### Example Usage

```go
package main

import(
	"context"
	"log"
	"github.com/formancehq/formance-sdk-go"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
)

func main() {
    s := formance.New(
        formance.WithSecurity(shared.Security{
            Authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
        }),
    )

    ctx := context.Background()
    res, err := s.Orchestration.RunWorkflow(ctx, operations.RunWorkflowRequest{
        RequestBody: map[string]string{
            "quidem": "provident",
            "nam": "id",
            "blanditiis": "deleniti",
            "sapiente": "amet",
        },
        Wait: formance.Bool(false),
        WorkflowID: "deserunt",
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.RunWorkflowResponse != nil {
        // handle response
    }
}
```

## SendEvent

Send an event to a running workflow

### Example Usage

```go
package main

import(
	"context"
	"log"
	"github.com/formancehq/formance-sdk-go"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
)

func main() {
    s := formance.New(
        formance.WithSecurity(shared.Security{
            Authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
        }),
    )

    ctx := context.Background()
    res, err := s.Orchestration.SendEvent(ctx, operations.SendEventRequest{
        RequestBody: &operations.SendEventRequestBody{
            Name: "Wilma Mosciski",
        },
        InstanceID: "perferendis",
    })
    if err != nil {
        log.Fatal(err)
    }

    if res.StatusCode == http.StatusOK {
        // handle response
    }
}
```
