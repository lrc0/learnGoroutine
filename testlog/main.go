// package main

// import (
// 	"fmt"
// 	"io/ioutil"
// 	"os"
// 	// "time"
// 	"net/url"

// 	"github.com/astaxie/beego/logs"
// )

// func Logs() {
// 	log := logs.NewLogger(10000)
// 	log.SetLogger("console", "")
// 	log.SetLevel(logs.LevelDebug)
// 	log.EnableFuncCallDepth(true)
// 	// msg, err := ioutil.ReadFile(os.Args[1])
// 	// if err != nil {
// 	//     log.Error("err: ", err)
// 	//     return
// 	// }
// 	msg := `我是日志信息，我不是日志信息，你是日志信息`
// 	// log.Informational("logs: %s", string(msg))
// 	if err := ioutil.WriteFile("/tmp/test.log", []byte(msg), os.ModePerm); err != nil {
// 		log.Error(err.Error())
// 		return
// 	}

// }

// func main() {
// 	// for i := 0; i < 2; i++ {
// 	// 	go func() {
// 	// 		for {
// 	// 			Logs()
// 	// 		}
// 	// 	}()
// 	// }
// 	// Logs()

// 	// for {
// 	// 	fmt.Println("小牛扣扣使劲揪,小妞扣扣对准扣眼扣,小牛和小妞,谁学会了扣纽扣？")
// 	// 	time.Sleep(time.Second * 1)
// 	// }
// 	x := 3
// 	y := 2
// 	x *= y
// 	fmt.Println(x)
// 	m := url.Values{"lang": {"en"}}
// 	m.Add("item", "1")
// 	m.Add("item", "2")

// 	fmt.Println(m.Get("lang")) // "en"
// 	fmt.Println(m.Get("q"))    // ""
// 	fmt.Println(m.Get("item")) // "1"      (first value)
// 	fmt.Println(m["item"])     // "[1 2]"  (direct map access)

// 	fmt.Println(m(nil).Get("item")) // ""
// 	m.Add("item", "3")
// }

package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		fmt.Println(err)
	}
	done := make(chan int)
	go func() {
		io.Copy(os.Stdout, conn)
		done <- 1
	}()
	mustCopy(conn, os.Stdin)
	conn.Close()
	<-done
}

func mustCopy(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		fmt.Println(err)
	}
}
