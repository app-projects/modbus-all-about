package main

import (
	"mytask/source/runner"
	"time"
	"fmt"
)

func task1(detal int) int {
	fmt.Println("this is a task process task1")
	return 0
}
func task2(detal int) int {
	fmt.Println("this is a task process task2")
	return 0
}
func task3(detal int) int {
	fmt.Println("this is a task process task3")
	return 0
}

func main() {
	engine := runner.NewTaskContext(3 * time.Duration(time.Second))
	engine.AddTask(task1, task2, task3)

	err := engine.Start()

	switch err {
	case runner.ERR_TIMEOUT:
		fmt.Println(runner.ERR_TIMEOUT.Error())
	case runner.ERR_INTER:
		fmt.Println(runner.ERR_INTER.Error())
	}

}
