package main

import (
	"fmt"
	"net/http"
	// "runtime/pprof"
	"time"
)

//Job .
type Job interface {
	Do(ch chan bool)
}

//Worker .
type Worker struct {
	//控制任务队列的并发量
	WorkPool chan chan Job
	//任务缓存队列
	JobChan chan Job
	//退出任务信号
	Quit chan bool
}

//全局worker
var worker = NewWorker()

//NewWorker 初始化一个worker
func NewWorker() *Worker {
	return &Worker{
		WorkPool: make(chan chan Job, 1), //保证只有1个队列
		JobChan:  make(chan Job, 10),     //队列大小自定，此处定为10个
		Quit:     make(chan bool),
	}
}

//Start 启动创建工作
func (w *Worker) Start() {
	go func() {
		//注册Job队列到pool里面
		w.WorkPool <- w.JobChan
		for {
			select {
			case job := <-w.JobChan:
				job.Do(w.Quit)
			}
		}
	}()
}

//Stop 停止工作，每次调用结束一个task
func (w *Worker) Stop() {
	go func() {
		w.Quit <- true
	}()
}

// Task .
type Task struct {
}

//Do .
func (t *Task) Do(ch chan bool) {
	for {
		select {

		case <-ch:
			fmt.Println("quit")
			return

		default:
			fmt.Println(time.Now().String())
			time.Sleep(1 * time.Second)
		}
	}
}

func main() {

	http.HandleFunc("/start", start)
	http.HandleFunc("/stop", stop)

	http.ListenAndServe(":8888", nil)

}

func start(w http.ResponseWriter, r *http.Request) {
	t := new(Task)
	worker.JobChan <- t
	fmt.Fprintln(w, fmt.Sprintf("======>> 任务数量: %+v \n", len(worker.JobChan)))
	// p := pprof.Lookup("goroutine")
	// p.WriteTo(w, 1)
	go worker.Start()
}

func stop(w http.ResponseWriter, r *http.Request) {
	worker.Stop()
}
