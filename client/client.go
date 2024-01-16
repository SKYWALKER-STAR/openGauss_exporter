package main

import (
	"os"
	_ "fmt"
	"metrics"
)

func main() {
	counter := metrics.CreateMetric("unter","Test metric","Test Help Metric Help")

	if counter == nil {
		os.Exit(-1)
	}
}
