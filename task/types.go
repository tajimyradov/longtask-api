package task

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Task interface {
	ID() string
	Type() string
	Payload() any
	Run(ctx context.Context) (any, error)
}

type LongTask struct {
	id      string
	payload map[string]interface{}
}

func NewLongTask(payload map[string]interface{}) *LongTask {
	return &LongTask{
		id:      uuid.NewString(),
		payload: payload,
	}
}

func (t *LongTask) ID() string   { return t.id }
func (t *LongTask) Type() string { return "long_task" }
func (t *LongTask) Payload() any { return t.payload }

func (t *LongTask) Run(ctx context.Context) (any, error) {
	select {
	case <-time.After(3 * time.Minute):
		return map[string]string{"message": "done"}, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}
