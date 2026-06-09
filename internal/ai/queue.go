package ai

import (
	"context"
	"sync"

	"go.uber.org/zap"
)

type Priority int

const (
	PriorityNormal  Priority = 0
	PriorityPremium Priority = 1
)

type AIRequest struct {
	Priority Priority
	Prompt   string
	ResultCh chan string
	ErrorCh  chan error
	Ctx      context.Context
}

type Queue struct {
	premiumCh chan *AIRequest
	normalCh  chan *AIRequest
	stopCh    chan struct{}
	wg        sync.WaitGroup
	enabled   bool
}

func NewQueue() *Queue {
	return &Queue{
		premiumCh: make(chan *AIRequest, 100),
		normalCh:  make(chan *AIRequest, 100),
		stopCh:    make(chan struct{}),
	}
}

func (q *Queue) StartWorker(ctx context.Context, gemini *Client, sugar *zap.SugaredLogger) {
	if q.enabled {
		sugar.Warn("ai queue worker already running")
		return
	}
	q.enabled = true

	q.wg.Add(1)
	go func() {
		defer q.wg.Done()
		sugar.Infow("ai queue worker started")

		for {
			select {
			case <-q.stopCh:
				sugar.Infow("ai queue worker stopped")
				return
			case req := <-q.premiumCh:
				q.processRequest(ctx, gemini, sugar, req)
			default:
				select {
				case <-q.stopCh:
					sugar.Infow("ai queue worker stopped")
					return
				case req := <-q.premiumCh:
					q.processRequest(ctx, gemini, sugar, req)
				case req := <-q.normalCh:
					q.processRequest(ctx, gemini, sugar, req)
				}
			}
		}
	}()
}

func (q *Queue) processRequest(ctx context.Context, gemini *Client, sugar *zap.SugaredLogger, req *AIRequest) {
	select {
	case <-req.Ctx.Done():
		sugar.Debugw("ai queue: request cancelled before processing")
		return
	default:
	}

	result, err := gemini.Generate(req.Ctx, req.Prompt)
	if err != nil {
		sugar.Errorw("ai queue: request failed", "error", err)
		req.ErrorCh <- err
		return
	}

	req.ResultCh <- result
}

func (q *Queue) Enqueue(req *AIRequest) {
	if req.Priority == PriorityPremium {
		q.premiumCh <- req
	} else {
		q.normalCh <- req
	}
}

func (q *Queue) StopWorker() {
	if !q.enabled {
		return
	}
	close(q.stopCh)
	q.wg.Wait()
	q.enabled = false
}

func (q *Queue) IsEnabled() bool {
	return q.enabled
}
