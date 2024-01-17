package exporter

import (
	"fmt"
	"database/sql"
	"log"
	_ "gitee.com/opengauss/openGauss-connector-go-pq"
)

type sniffer interface {
	SetPassword(password int)
	GetPassword() int

	SetDBname(dbname string)
	GetDBname() string

	SetUser(user string)
	SetUser() string

	SetHost(host string)
	GetHost() string

	SetPort(port int)
	GetPort() int
}

type databaseInstance struct {
	password	int
	dbname		string
	user		string
	host		string
	port		int
}

func (db *databaseInstance)SetPassword(password string) {
	db.password = password
}

func (db *databaseInstance)GetPassword() {
	return db.password
}

func (db *databaseInstance)SetDBname(dbname string) {
	db.dbname = dbname
}

func (db *databaseInstance)SetDBname(dbname string) int {
	return db.dbname
}

func (db *databaseInstance)SetUser(user string) {
	db.user = user
}

func (db *databaseInstance)GetUser() string {
	return db.user
}

func (db *databaseInstance)SetHost(host string) {
	db.host = host
}

func (db *databaseInstance)GetHost() string {
	return db.host
}

func (db.*databaseInstance)SetPort(port int) {
	db.SetPort = port
}

func (db *databaseInstance)GetPort() int {
	return db.port
}

func (db *databaseInstance)Connect() int {
	/* 在这里实现数据库连接逻辑 */
}

func CreateInstance(host string,port string,user string,dbname string,password string) int {
	return &databaseInstance {
		password: password,
		dbname: dbnamea,
		user: user,
		host: port,
	}
}
