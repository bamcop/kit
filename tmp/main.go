package main

import (
	"fmt"
	"os"

	"github.com/bamcop/kit/data"
	"github.com/bamcop/kit/task/simple_task"
	"golang.org/x/exp/slog"
)

type SSR string

func (S SSR) ID() string {
	return string(S)
}

func (S SSR) LoadTask() []simple_task.Task {
	return []simple_task.Task{
		SSR("A"),
		SSR("B"),
		SSR("C"),
	}
}

func (S SSR) Handler(task simple_task.Task) error {
	fmt.Println("==>", task.ID())
	if task.ID() == "B" {
		panic(1)
	}

	return nil
}

func (S SSR) OnSucc(task simple_task.Task) {
	//TODO implement me
	//panic("implement me")
}

func (S SSR) OnFail(task simple_task.Task) {
	//TODO implement me
	//panic("implement me")
}

func main() {
	h := slog.HandlerOptions{
		AddSource: true,
	}.NewJSONHandler(os.Stdout)
	logger := slog.New(h)
	slog.SetDefault(logger)

	r := simple_task.NewTaskRunner(2, func() simple_task.Task {
		return SSR("1")
	})

	fails := r.Retry(3).Run()
	if len(fails) != 0 {
		fmt.Println(string(data.MarshalIndent(fails).Unwrap()))
	}
}
