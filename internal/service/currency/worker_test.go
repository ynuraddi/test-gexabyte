package currency

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWorker(t *testing.T) {
	tasks := []doTaskFunc{
		func() interface{} {
			return 1
		},
		func() interface{} {
			return 1
		},
		func() interface{} {
			return 1
		},
	}

	service := Currency{}
	out := service.taskResultStream(context.Background(), tasks...)
	for i := 0; i < 3; i++ {
		<-out
	}
	res, isOpen := <-out
	assert.Equal(t, nil, res)
	assert.False(t, isOpen)
}
