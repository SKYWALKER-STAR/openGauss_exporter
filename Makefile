
all:
	go build client/gaussdb_exporter.go

clean:
	rm -rf client/gaussdb_exporter.go
