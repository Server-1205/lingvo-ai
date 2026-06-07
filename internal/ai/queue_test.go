package ai

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestQueue_EnqueueNormal(t *testing.T) {
	q := NewQueue()

	req := &AIRequest{
		Priority: PriorityNormal,
		Prompt:   "say hello",
		ResultCh: make(chan string, 1),
		ErrorCh:  make(chan error, 1),
		Ctx:      context.Background(),
	}

	q.Enqueue(req)
	assert.Len(t, q.normalCh, 1, "normal channel should have 1 item")
	assert.Len(t, q.premiumCh, 0, "premium channel should be empty")
}

func TestQueue_EnqueuePremium(t *testing.T) {
	q := NewQueue()

	req := &AIRequest{
		Priority: PriorityPremium,
		Prompt:   "premium req",
		ResultCh: make(chan string, 1),
		ErrorCh:  make(chan error, 1),
		Ctx:      context.Background(),
	}

	q.Enqueue(req)
	assert.Len(t, q.premiumCh, 1, "premium channel should have 1 item")
	assert.Len(t, q.normalCh, 0, "normal channel should be empty")
}

func TestQueue_EnqueueBothChannels(t *testing.T) {
	q := NewQueue()

	normalReq := &AIRequest{
		Priority: PriorityNormal,
		Prompt:   "normal",
		ResultCh: make(chan string, 1),
		ErrorCh:  make(chan error, 1),
		Ctx:      context.Background(),
	}
	premiumReq := &AIRequest{
		Priority: PriorityPremium,
		Prompt:   "premium",
		ResultCh: make(chan string, 1),
		ErrorCh:  make(chan error, 1),
		Ctx:      context.Background(),
	}

	q.Enqueue(normalReq)
	q.Enqueue(premiumReq)

	assert.Len(t, q.normalCh, 1)
	assert.Len(t, q.premiumCh, 1)
}

func TestQueue_StartAndStopWorker(t *testing.T) {
	sugar := newTestSugar()
	q := NewQueue()

	q.StartWorker(context.Background(), nil, sugar)
	assert.True(t, q.IsEnabled())

	q.StopWorker()
	assert.False(t, q.IsEnabled())
}

func TestQueue_StopWorkerIdempotent(t *testing.T) {
	sugar := newTestSugar()
	q := NewQueue()

	q.StartWorker(context.Background(), nil, sugar)
	q.StopWorker()
	q.StopWorker()
	assert.False(t, q.IsEnabled())
}

func TestQueue_StartWorkerTwiceDoesNotCrash(t *testing.T) {
	sugar := newTestSugar()
	q := NewQueue()

	q.StartWorker(context.Background(), nil, sugar)
	q.StartWorker(context.Background(), nil, sugar)

	q.StopWorker()
	assert.False(t, q.IsEnabled())
}

func TestQueue_PremiumChannelReadFirst(t *testing.T) {
	q := NewQueue()

	premiumReq := &AIRequest{
		Priority: PriorityPremium,
		Prompt:   "premium",
		ResultCh: make(chan string, 1),
		ErrorCh:  make(chan error, 1),
		Ctx:      context.Background(),
	}
	normalReq := &AIRequest{
		Priority: PriorityNormal,
		Prompt:   "normal",
		ResultCh: make(chan string, 1),
		ErrorCh:  make(chan error, 1),
		Ctx:      context.Background(),
	}

	q.Enqueue(premiumReq)
	q.Enqueue(normalReq)

	select {
	case got := <-q.premiumCh:
		assert.Equal(t, "premium", got.Prompt)
	case <-time.After(time.Second):
		t.Fatal("timeout reading from premium channel")
	}

	select {
	case got := <-q.normalCh:
		assert.Equal(t, "normal", got.Prompt)
	case <-time.After(time.Second):
		t.Fatal("timeout reading from normal channel")
	}
}



func TestQueue_EmptyQueueDoesNotBlock(t *testing.T) {
	sugar := newTestSugar()
	q := NewQueue()

	q.StartWorker(context.Background(), nil, sugar)

	time.Sleep(50 * time.Millisecond)

	q.StopWorker()
	assert.False(t, q.IsEnabled())
}

func TestQueue_NewQueueNotEnabled(t *testing.T) {
	q := NewQueue()
	assert.False(t, q.IsEnabled())
}
