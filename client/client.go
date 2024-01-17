package main

import (
	"os"
	"fmt"
	"metrics"
	"exporter"
)

func main() {
	counter := metrics.CreateMetric("Counter","Test metric","Test Help Metric Help")

	if counter == nil {
		os.Exit(-1)
	}

	db := exporter.CreateInstance("127.0.0.1",5433,"ming","postgres","")

	db.Connect()
	obj := db.GetConn()

	fmt.Println(obj)
}
