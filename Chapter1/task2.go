package main

import (
	"fmt"
	"math"
	"sync"
	"sync/atomic"
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


type Shape interface {
    Area() float64
    Perimeter() float64
}

type Shapes interface {
	Area() float64
	Perimeter() float64
}

type Rectangle struct {
	Width float64
	Height float64
}

func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}


func (r Rectangle) Perimeter() float64 {
    return 2 * (r.Width + r.Height)
}


type Circle struct {
    Radius float64
}

func (c Circle) Area() float64 {
    return math.Pi * math.Pow(c.Radius, 2)
}

func (c Circle) Perimeter() float64 {
    return 2 * math.Pi * c.Radius
}

type Person struct {
    Name string
    Age  int
}


type Employee struct {
    Person
    EmployeeID string
}

func (e Employee) PrintInfo() {
    fmt.Printf("员工ID: %s\n", e.EmployeeID)
    fmt.Printf("姓名: %s\n", e.Name)
    fmt.Printf("年龄: %d\n", e.Age)
}


type Counter struct {
    value int
    mutex sync.Mutex
}

// Increment 安全地增加计数器值
func (c *Counter) Increment() {
    c.mutex.Lock()
    defer c.mutex.Unlock()
    c.value++
}

// Value 返回当前计数器值
func (c *Counter) Value() int {
    c.mutex.Lock()
    defer c.mutex.Unlock()
    return c.value
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
		time.Sleep(200 * time.Millisecond)
	})
	scheduler.AddTask("Task2", func() {
		time.Sleep(500 * time.Millisecond)
	})
	scheduler.AddTask("Task3", func() {
		time.Sleep(1 * time.Second)
	})

	// 执行任务
	scheduler.Run()

	// 输出结果
	fmt.Println("任务执行结果：")
	for _, task := range scheduler.GetResults() {
		fmt.Printf("%s 执行耗时: %v\n", task.Name, task.TimeTaken)
	}

	//question5
	rectangle := Rectangle{Width: 5, Height: 10}
    fmt.Printf("矩形面积: %.2f\n", rectangle.Area())
    fmt.Printf("矩形周长: %.2f\n", rectangle.Perimeter())

    circle := Circle{Radius: 3}
    fmt.Printf("圆形面积: %.2f\n", circle.Area())
    fmt.Printf("圆形周长: %.2f\n", circle.Perimeter())

	//question6
	employee := Employee{
		Person: Person{
			Name: "John",
			Age:  30,
		},
		EmployeeID: "123456",
	}
	employee.PrintInfo()

	//question7
	ch := make(chan int)
    var wg sync.WaitGroup

    wg.Add(1)
    go func() {
        defer wg.Done()
        for i := 1; i <= 10; i++ {
            ch <- i
        }
        close(ch)
    }()

    wg.Add(1)
    go func() {
        defer wg.Done()
        for num := range ch {
            fmt.Println("接收到:", num)
        }
    }()

    wg.Wait()
    fmt.Println("finished")

	//question8
	ch2 := make(chan int, 50)
    var wg2 sync.WaitGroup

    // 生产者协程
    wg2.Add(1)
    go func() {
        defer wg2.Done()
        for i := 1; i <= 100; i++ {
            ch2 <- i // 发送数据（缓冲区未满时不会阻塞）
            fmt.Printf("生产者发送: %d\n", i)
        }
        close(ch2) // 生产完毕后关闭通道
    }()

    // 消费者协程
    wg2.Add(1)
    go func() {
        defer wg2.Done()
        for num := range ch2 { // 自动处理通道关闭
            fmt.Printf("消费者接收: %d\n", num)
        }
    }()

    wg2.Wait()

	//question9
	var (
        counter Counter
        wg3      sync.WaitGroup
        workers = 10
        ops     = 1000
    )

    wg3.Add(workers)
    for i := 0; i < workers; i++ {
        go func() {
            defer wg3.Done()
            for j := 0; j < ops; j++ {
                counter.Increment()
            }
        }()
    }

    wg3.Wait()
    fmt.Printf("最终计数器值: %d\n", counter.Value())

	//question10
	var (
        counter2 int64
        wg4      sync.WaitGroup
        workers2 = 10
        ops2     = 1000
    )

    wg4.Add(workers2)
    for i := 0; i < workers2; i++ {
        go func() {
            defer wg4.Done()
            for j := 0; j < ops2; j++ {
                atomic.AddInt64(&counter2, 1)
            }
        }()
    }

    wg4.Wait()
    fmt.Printf("final counter value: %d\n", atomic.LoadInt64(&counter2))
}
