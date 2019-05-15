package main

import (
	"./metric"
	"./fast"
)

func main(){

	metric.RunMetricServer()

	fast.RunFastServer()
	//standard.RunStandardServer()
}
