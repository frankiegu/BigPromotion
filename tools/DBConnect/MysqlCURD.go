package main

import (
	"database/sql"
	"fmt"
	_ "github.com/GO-SQL-Driver/MySQL"
	"log"
)

func initMysql() *sql.DB {
	db, err := sql.Open("mysql", "root:1234@tcp(127.0.0.1:3306)/test?charset=utf8")
	if err != nil {
		log.Println(err)
	}
	db.SetMaxOpenConns(30)
	db.SetMaxIdleConns(1)
	db.Ping()
	return db
}

func query() {
	db := initMysql()
	defer db.Close()

	row, err := db.Query("select * from user")
	if err != nil {
		log.Fatal(err)
	}

	for row.Next() {
		//var m = make(map[string]interface{})
		var user_id int
		var user_name string
		var user_age int
		var user_sex int
		err = row.Scan(&user_id, &user_name, &user_age, &user_sex)
		fmt.Println("user_id: ", user_id)
		fmt.Println("user_name: ", user_name)
		fmt.Println("user_age: ", user_age)
		fmt.Println("user_sex: ", user_sex)
	}
}

func insert() {
	db := initMysql()
	defer db.Close()

	db.Begin() //transaction begin
	stmt, _ := db.Prepare(`insert user (user_name, user_age, user_sex) values (?,?,?)`)
	result, _ := stmt.Exec("tony", 20, 1)
	fmt.Println(result)
	stmt.Close()


}

func delete() {
	db := initMysql()
	defer db.Close()

	result, _ := db.Exec("delete from user where user_id = ?", 2)
	c, _ := result.RowsAffected()
	log.Println("delete affected rows: ",c)
}

func update() {
	db := initMysql()
	defer db.Close()

	result, _ := db.Exec("update user set user_name = ? where user_id = ?", "kk", 1)
	c, _ := result.RowsAffected()
	log.Println("update affected rows: ", c)
}

func main() {
	initMysql()
	insert()
	query()
	delete()
	update()
}
