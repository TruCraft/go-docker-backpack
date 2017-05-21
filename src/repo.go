package main

import (
	"database/sql"
	"fmt"
	"time"
	"os"
	"io/ioutil"

	_ "github.com/go-sql-driver/mysql"
	"strings"
)

type DBInfo struct {
	Name string `json:"name"`
	User string `json:"user"`
	Pass string `json:"pass"`
	Host string `json:"host"`
	Port int    `json:"port"`
}

var db_table string = "tf2_backpacks"
var db_host string = "127.0.0.1"
var db_name string = "trucraft"
var db_user string = "trucraft_service"
var db_info DBInfo

func InitDBInfo() {
	if db_info == (DBInfo{}) {
		getTableName()
		if passfile, ok := os.LookupEnv("MYSQL_PASSWORD_FILE"); ok {
			if env_db_host, ok := os.LookupEnv("MYSQL_DB_HOST"); ok {
				db_host = env_db_host
			}

			if env_db_name, ok := os.LookupEnv("MYSQL_DB_NAME"); ok {
				db_name = env_db_name
			}

			if env_db_user, ok := os.LookupEnv("MYSQL_DB_USER"); ok {
				db_user = env_db_user
			}

			pass, err := ioutil.ReadFile(passfile)
			checkErr(err)
			db_info = DBInfo{
				db_name,
				db_user,
				strings.TrimSpace(string(pass)),
				db_host,
				3306,
			}
		} else {
			panic("Environment variable \"MYSQL_PASSWORD_FILE\" not found")
		}
	}
}

// Get the {limit} most recent backpacks from db for {username}
func RepoFindBackpack(username string, limit int) Backpacks {
	InitDBInfo()
	db, err := sql.Open("mysql", getDataSourceName())
	checkErr(err)

	query := fmt.Sprintf("SELECT * FROM %s WHERE username = '%s' ORDER BY timestamp DESC LIMIT %d", db_table, username, limit)

	rows, err := db.Query(query)
	checkErr(err)

	var backpacks Backpacks

	for rows.Next() {
		var id int64
		var username string
		var timestamp string
		var backpack_json string
		err = rows.Scan(&id, &username, &timestamp, &backpack_json)
		checkErr(err)
		backpacks = append(backpacks, Backpack{id, username, timestamp, backpack_json})
	}
	return backpacks
}

func getTableName() {
	if env_db_table, ok := os.LookupEnv("MYSQL_DB_TABLE"); ok {
		db_table = env_db_table
	}
}

func RepoAddBackpack(backpack Backpack) Backpack {
	InitDBInfo()
	db, err := sql.Open("mysql", getDataSourceName())
	checkErr(err)

	query := fmt.Sprintf("INSERT INTO %s SET username=?, timestamp=?, backpack_json=?", db_table)

	stmt, err := db.Prepare(query)
	checkErr(err)

	backpack.Timestamp = time.Now().Format("2006-01-02 15:04:05")
	res, err := stmt.Exec(backpack.Username, backpack.Timestamp, backpack.Items)
	checkErr(err)

	if id, err := res.LastInsertId(); err == nil {
		backpack.Id = id
	}

	//affect, err := res.RowsAffected()
	//checkErr(err)

	return backpack
}

func getDataSourceName() (string) {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8", db_info.User, db_info.Pass, db_info.Host, db_info.Port, db_info.Name)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
