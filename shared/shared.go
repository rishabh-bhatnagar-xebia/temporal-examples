package shared

import "fmt"

var QueueNameBasic = "QUEUE_BASIC"
var QueueNameAsyncV1 = "QUEUE_ASYNC_V1"
var QueueNameAsyncV2 = "QUEUE_ASYNC_V2"
var HttpWorkflowTypeParamName = "type"

type Counter struct {
	count int
}

func (c *Counter) Get() string {
	c.count += 1
	return fmt.Sprint(c.count)
}

var C *Counter = nil

func NewCounter() *Counter {
	if C != nil {
		return C
	}
	c := 0
	C = &Counter{c}
	return C
}
