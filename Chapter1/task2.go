package main

import (
	"fmt"
	"sync"
	"time"
)

func Pointer(nums *int) {
	*nums += 10
}

func SliceExecute(nums *[]int) {
	for i := range *nums {
		(*nums)[i] *= 2
	}
}

func printOdd(ch chan struct{}, next chan struct{}, done chan struct{}) {
	defer close(done)
	for i := 1; i <= 10; i += 2 {
		<-ch // 等待轮到自己
		fmt.Println(i)
		next <- struct{}{} // 通知下一个协程
	}
}

func printEven(ch chan struct{}, next chan struct{}, done chan struct{}) {
	defer close(done)
	for i := 2; i <= 10; i += 2 {
		<-ch // 等待轮到自己
		fmt.Println(i)
		if i < 10 {
			next <- struct{}{} // 通知下一个协程
		}
	}
}

// Task 表示一个可执行的任务
type Task struct {
	Name      string        // 任务名称
	Execute   func()        // 任务执行函数
	TimeTaken time.Duration // 执行耗时
}

// TaskScheduler 任务调度器
type TaskScheduler struct {
	maxWorkers int    // 最大工作协程数
	tasks      []Task // 待执行任务列表
	results    []Task // 已完成任务列表
	wg         sync.WaitGroup
	mu         sync.Mutex // 保护 results
}

// NewTaskScheduler 创建新的任务调度器
func NewTaskScheduler(maxWorkers int) *TaskScheduler {
	return &TaskScheduler{
		maxWorkers: maxWorkers,
		tasks:      make([]Task, 0),
		results:    make([]Task, 0),
	}
}

// AddTask 添加任务到调度器
func (s *TaskScheduler) AddTask(name string, task func()) {
	s.tasks = append(s.tasks, Task{
		Name:    name,
		Execute: task,
	})
}

// Run 执行所有任务
func (s *TaskScheduler) Run() {
	taskCh := make(chan Task, len(s.tasks))

	// 启动工作协程
	for i := 0; i < s.maxWorkers; i++ {
		s.wg.Add(1)
		go s.worker(taskCh)
	}

	// 发送任务到通道
	go func() {
		for _, task := range s.tasks {
			taskCh <- task
		}
		close(taskCh)
	}()

	// 等待所有工作协程完成
	s.wg.Wait()
}

// worker 工作协程，处理任务并记录时间
func (s *TaskScheduler) worker(taskCh chan Task) {
	defer s.wg.Done()

	for task := range taskCh {
		start := time.Now()
		task.Execute()
		task.TimeTaken = time.Since(start)

		// 记录结果（加锁保护共享资源）
		s.mu.Lock()
		s.results = append(s.results, task)
		s.mu.Unlock()
	}
}

// GetResults 获取所有任务的执行结果
func (s *TaskScheduler) GetResults() []Task {
	return s.results
}

func main() {
	//question1
	nums := 10
	Pointer(&nums)
	fmt.Println(nums)

	//question2
	numsSlice := []int{1, 2, 3, 4, 5}
	SliceExecute(&numsSlice)
	fmt.Println(numsSlice)

	//question3
	oddCh := make(chan struct{})
	evenCh := make(chan struct{})
	doneOdd := make(chan struct{})
	doneEven := make(chan struct{})

	go printOdd(oddCh, evenCh, doneOdd)
	go printEven(evenCh, oddCh, doneEven)

	oddCh <- struct{}{}

	<-doneOdd
	<-doneEven

	//question4
	scheduler := NewTaskScheduler(3)

	// 添加示例任务
	scheduler.AddTask("Task1", func() {
		time.Sleep(2 * time.Second)
	})
	scheduler.AddTask("Task2", func() {
		time.Sleep(1 * time.Second)
	})
	scheduler.AddTask("Task3", func() {
		time.Sleep(3 * time.Second)
	})

	// 执行任务
	scheduler.Run()

	// 输出结果
	fmt.Println("任务执行结果：")
	for _, task := range scheduler.GetResults() {
		fmt.Printf("%s 执行耗时: %v\n", task.Name, task.TimeTaken)
	}
}
