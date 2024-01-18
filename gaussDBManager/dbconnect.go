package dbmanager

import (
	"fmt"
	"database/sql"
	"log"
	_ "gitee.com/opengauss/openGauss-connector-go-pq"
)

type sniffer interface {
	SetPassword(password string)
	GetPassword() string

	SetDBname(dbname string)
	GetDBname() string

	SetUser(user string)
	GetUser() string

	SetHost(host string)
	GetHost() string

	SetPort(port int)
	GetPort() int

	Connect() int
	SetConn(conn *sql.DB) 
	GetConn() *sql.DB
}

type databaseInstance struct {
	password	string
	dbname		string
	user		string
	host		string
	port		int
	conn		*sql.DB
}

func (db *databaseInstance)SetPassword(password string) {
	db.password = password
}

func (db *databaseInstance)GetPassword() string {
	return db.password
}

func (db *databaseInstance)SetDBname(dbname string) {
	db.dbname = dbname
}

func (db *databaseInstance)GetDBname() string {
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

func (db *databaseInstance)SetPort(port int) {
	db.port = port
}

func (db *databaseInstance)GetPort() int {
	return db.port
}

func (db *databaseInstance)SetConn(conn *sql.DB) {
	db.conn = conn
}

func (db *databaseInstance)Connect() int {
	/* 数据库连接逻辑 */
	var connStr string

	connStr = fmt.Sprintf("user=%s dbname=%s password=%s host=%s port=%d",db.user,db.dbname,db.password,db.host,db.port)

	conn,err := sql.Open("opengauss",connStr)
	if err != nil {
		log.Fatal(err)
	}

	db.SetConn(conn)
	return 0
}

func (db *databaseInstance)GetConn() *sql.DB {
	return db.conn
}

func CreateInstance(host string,port int,user string,dbname string,password string) sniffer {
	return &databaseInstance {
		password: password,
		dbname: dbname,
		user: user,
		host: host,
		port: port,
	}
}
