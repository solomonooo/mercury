//import "github.com/solomonooo/mercury"
//author : solomonooo
//create time : 2016-09-08

package plugin

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/solomonooo/mercury"
	"time"
)

type Mysql struct {
	host     string
	port     int
	user     string
	password string
	database string

	dbAddr string
	conn   *sql.DB
}

type MysqlResult struct {
	rows *sql.Rows
}

func NewMysqlConn(host string, port int, user string, pwd string, database string) (*Mysql, error) {
	var mysql Mysql
	mysql.host = host
	mysql.port = port
	mysql.user = user
	mysql.password = pwd
	mysql.database = database

	var err error
	mysql.dbAddr = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", user, pwd, host, port, database)
	mysql.conn, err = sql.Open("mysql", mysql.dbAddr)
	if err != nil {
		mercury.Error("connect to db %s:%d/%s failed, err[%s]", host, port, database, err.Error())
		return nil, err
	}
	mysql.conn.SetMaxOpenConns(1000)
	mysql.conn.SetMaxIdleConns(500)
	err = mysql.conn.Ping()
	return &mysql, err
}

func (mysql *Mysql) Query(queryStr string) (*MysqlResult, error) {
	var result MysqlResult
	var err error
	start := time.Now().UnixNano()
	result.rows, err = mysql.conn.Query(queryStr)
	if err != nil {
		mercury.Error("query from db %s:%d/%s failed, sql[%s], err[%s]", mysql.host, mysql.port, mysql.database, queryStr, err.Error())
	}
	end := time.Now().UnixNano()
	mercury.Info("mysql query cost[%dms], sql[%s]", (end-start)/1000000, queryStr)
	return &result, err
}

func (result *MysqlResult) Get(args ...interface{}) bool {
	if false == result.rows.Next() {
		return false
	}

	err := result.rows.Scan(args...)
	if err != nil {
		mercury.Error("parse mysql result failed, err[%s]", err.Error())
		return false
	}
	return true
}

func (result *MysqlResult) Close() error {
	return result.rows.Close()
}

func (mysql *Mysql) Close() error {
	return mysql.conn.Close()
}
