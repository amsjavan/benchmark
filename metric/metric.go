package metric

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func print(){
	for {
		rsp, err := http.Get("http://localhost:8080/metrics")

		if err != nil {
			fmt.Println(err)
		} else {
			r, _ := ioutil.ReadAll(rsp.Body)
			str := string(r)
			s := strings.Split(str, "\n")
			fmt.Println(s[11])
		}
		time.Sleep(200 * time.Millisecond)
	}
}

func RunMetricServer(){
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		_ = http.ListenAndServe(":8080", nil)
	}()
}

