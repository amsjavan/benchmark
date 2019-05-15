package fast

import (
	"fmt"
	"github.com/valyala/fasthttp"
	"sync"
	"sync/atomic"
	"time"
)

var ops uint64
var client *fasthttp.Client
func init(){
	client = &fasthttp.Client{
		MaxConnsPerHost: 1000,
		MaxIdleConnDuration: 5 * time.Second}
}

func doRequest(url string, wait *sync.WaitGroup) {
	if wait != nil {
		defer wait.Done()
	}

	req := fasthttp.AcquireRequest()
	req.SetRequestURI(url)

	resp := fasthttp.AcquireResponse()
	err := client.Do(req, resp)
	if err != nil {
		fmt.Println(err)
	}

	//bodyBytes := resp.Body()
	//println(string(bodyBytes))
	// User-Agent: fasthttp
	// Body:
	atomic.AddUint64(&ops, 1)

}


func batchRequest(n int, waitGroup *sync.WaitGroup){
	defer waitGroup.Done()
	for i:=0 ; i < n; i++ {
		doRequest("http://localhost:8080", nil)
	}
}

//Good
func RunPool(){
	var waitGroup sync.WaitGroup


	for i:=0; i <5000; i++ {
		waitGroup.Add(1)
		go batchRequest(100, &waitGroup)

	}

	start := time.Now()

	waitGroup.Wait()

	fmt.Println(time.Since(start))

}

//Bad
func Run() {
	start := time.Now()
	doRequest("http://localhost:8080", nil)
	fmt.Println(time.Since(start))
	//var waitGroup sync.WaitGroup
	//
	//
	//for i:=0; i < 5000; i++ {
	//	waitGroup.Add(1)
	//	go 	doRequest("http://localhost:8080", &waitGroup)
	//}
	//start := time.Now()
	//
	//waitGroup.Wait()
	//
	//fmt.Println(time.Since(start))
}