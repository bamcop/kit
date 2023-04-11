package simple_task

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/smallnest/chanx"
	"golang.org/x/exp/slog"
)

type Task interface {
	ID() string
	LoadTask() []Task
	Handler(Task) error
	OnSucc(Task)
	OnFail(Task)
}

type TaskRunner[T Task] struct {
	sync.Mutex
	fn       func() Task
	ch       *chanx.UnboundedChan[Task]
	tasks    []Task         // 所有任务
	retry    int            // 重试次数
	ncpu     int            // 线程数
	failMap  map[string]int //
	succList []Task         // 成功列表
	failList []Task         // 失败列表
}

func NewTaskRunner(ncpu int, fn func() Task) TaskRunner[Task] {
	return TaskRunner[Task]{
		fn:      fn,
		ch:      chanx.NewUnboundedChan[Task](ncpu),
		retry:   0,
		ncpu:    ncpu,
		failMap: map[string]int{},
	}
}

func (t *TaskRunner[T]) Retry(n int) *TaskRunner[T] {
	t.retry = n

	return t
}

func (t *TaskRunner[T]) Run() []Task {
	t.tasks = t.fn().LoadTask()
	for _, task := range t.tasks {
		t.ch.In <- task
		t.failMap[task.ID()] = 0
	}

	// 检查任务是否调度完成
	timer := time.NewTimer(time.Second)
	go func() {
		for {
			<-timer.C

			if !t.hasRetryTask() {
				close(t.ch.In)
			}
		}
	}()

	// 消费
	done := make(chan struct{}, t.ncpu)
	for i := 1; i <= t.ncpu; i++ {
		go t.consumer(i, t.ch.Out, done)
	}

	// 等待完成
	for i := 1; i <= t.ncpu; i++ {
		<-done
		slog.Info(fmt.Sprintf("goroutine  %d done", i))
	}

	return t.failList
}

func (t *TaskRunner[T]) hasRetryTask() bool {
	t.Lock()
	defer t.Unlock()

	for _, n := range t.failMap {
		if n < t.retry {
			return true
		}
	}

	return false
}

func (t *TaskRunner[T]) couldRetry(item Task) bool {
	t.Lock()
	defer t.Unlock()

	if t.failMap[item.ID()] < t.retry {
		return true
	}
	return false
}

func (t *TaskRunner[T]) consumer(goid int, ch <-chan Task, done chan struct{}) {
	for task := range ch {
		if err := t.handle(goid, task); err != nil {
			slog.LogAttrs(context.Background(), slog.LevelError, err.Error(), slog.String("id", task.ID()))

			if t.couldRetry(task) {
				t.ch.In <- task

				t.Lock()
				t.failMap[task.ID()] = t.failMap[task.ID()] + 1
				t.Unlock()
			} else {
				t.failList = append(t.failList, task)
				t.fn().OnFail(task)
			}
		}
	}

	done <- struct{}{}
}

func (t *TaskRunner[T]) handle(goid int, item Task) (err error) {
	defer func() {
		if r := recover(); r != nil {
			v, ok := r.(error)
			if ok {
				err = v
			} else {
				err = fmt.Errorf("%+v", r)
			}
		}
	}()

	if err = t.fn().Handler(item); err != nil {
		return err
	}

	slog.LogAttrs(
		context.Background(), slog.LevelInfo, "task done",
		slog.String("id", item.ID()),
		slog.Int("count", len(t.succList)+1),
	)
	t.fn().OnSucc(item)

	t.Lock()
	t.succList = append(t.succList, item)
	delete(t.failMap, item.ID())
	t.Unlock()

	return nil
}
