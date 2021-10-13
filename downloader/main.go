package main

import (
	"bufio"
	"fmt"
	"gopkg.in/logger.v1"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
)

type task struct {
	url      string
	filename string
	retry    int
}

type context struct {
	locker   *sync.RWMutex
	taskChan chan *task
	taskMap  map[string]int
	folder   string
}

const (
	ready = iota
	done
	fail

	maxRetry      = 3
	maxDownloader = 10
)

func main() {
	ctx := initContext()
	start(ctx)
	select {}
}

func start(ctx *context) {
	client := new(http.Client)

	for i := 0; i <= maxDownloader; i++ {
		go func() {
			task := <-ctx.taskChan
			task.download(ctx, client)
		}()
	}

	go func() {
		for {
			ctx.taskChan <- generateTask()
		}
	}()
}

func initContext() *context {
	err := os.Mkdir("download", 0777)
	if os.IsExist(err) {
		goto l
	}
	if err != nil {
		log.Error(err)
		return nil
	}
l:
	return &context{
		locker:   new(sync.RWMutex),
		taskChan: make(chan *task, maxDownloader),
		taskMap:  make(map[string]int),
		folder:   "download",
	}
}

func generateTask() *task {
	var t task
	var filename, url string
	fmt.Println("enter download url: ")
	fmt.Scanln(&filename, &url)

	t.filename = filename
	t.url = parseURL(url)
	t.retry = 0
	return &t
}

func (t *task) download(ctx *context, client *http.Client) {

	ctx.locker.RLock()
	if ctx.taskMap[t.url] == done {
		ctx.locker.RUnlock()
		return
	}
	ctx.locker.RUnlock()

	defer t.downloadRetry(ctx)

	req, err := http.NewRequest("GET", t.url, nil)
	if err != nil {
		log.Error(err)
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Error(err)
		return
	}
	defer resp.Body.Close()

	file, err := os.OpenFile(ctx.folder+"/"+t.filename, os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		log.Error(err)
		return
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		log.Error(err)
		return
	}

	newReader := bufio.NewReaderSize(resp.Body, 4096)

	var total int64
	for {
		n, err := io.Copy(file, newReader)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Error(err)
			return
		}
		total += n

		//进度条
		per := (float64(total) / float64(fileInfo.Size())) * 100
		fmt.Printf("\r[%s] %s%s", bar(per, 100), fmt.Sprintf("%.2f", per), fmt.Sprintf("%s", "% "))
	}

	ctx.locker.Lock()
	ctx.taskMap[t.url] = done
	ctx.locker.Unlock()

	fmt.Println()
}

func (t *task) downloadRetry(ctx *context) {
	ctx.locker.RLock()
	if ctx.taskMap[t.url] == done {
		ctx.locker.RUnlock()
		return
	}
	ctx.locker.RUnlock()

	if t.retry++; t.retry < maxRetry {
		go func() {
			ctx.taskChan <- t
		}()
	} else {
		ctx.locker.Lock()
		ctx.taskMap[t.url] = fail
		ctx.locker.Unlock()
	}
}

func bar(count, size float64) string {
	str := ""
	for i := float64(0); i < size; i++ {
		if i < count {
			str += "="
		} else {
			str += " "
		}
	}
	return str
}

func parseURL(url string) string {
	slice := strings.Split(url, "@")
	target := strings.Split(slice[1], "/")
	return "http://" + target[0]
}
