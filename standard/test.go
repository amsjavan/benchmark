package standard

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

var ops uint64

var client = &http.Client{}

func init() {
	http.DefaultTransport.(*http.Transport).MaxIdleConnsPerHost = 100
	http.DefaultTransport.(*http.Transport).MaxConnsPerHost = 10000
}

func io(waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()
	rsp, err := client.Get("http://localhost:9090")

	if err != nil {
		fmt.Println(err)
	} else {
		_, _ = ioutil.ReadAll(rsp.Body)
		//fmt.Println(string(r))
	}
	_ = rsp.Body.Close()
	atomic.AddUint64(&ops, 1)

}

func StressTest(n int){
	//ch := make(chan int, 10000)
	var waitGroup sync.WaitGroup
	for i:=0 ; i < n; i++ {
		waitGroup.Add(1)
		go io(&waitGroup)
	}
	start := time.Now()

	//for i:=0 ; i < n; i++ {
	//	fmt.Println(<- ch)
	//}
	waitGroup.Wait()
	fmt.Println("==", ops)
	fmt.Println(time.Since(start))
}
