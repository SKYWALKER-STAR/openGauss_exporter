package main

import (
	_ "os"
	_ "fmt"
	"log"
	"net/http"
	"exporter"
	"dbmanager"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	db := dbmanager.CreateInstance("127.0.0.1",5433,"ming","postgres","")

	db.Connect()
	//obj := db.GetConn()

	// 这是一个Metrics
	m := exporter.CreateTemplateMetrics("Test","This is test help")
	m.Register()
	m.SetCpuTemp(129)

	http.Handle("/metrics",promhttp.HandlerFor(m.GetRegistry(),promhttp.HandlerOpts{Registry: m.GetRegistry()}))

	log.Fatal(http.ListenAndServe(":9876",nil))
}
