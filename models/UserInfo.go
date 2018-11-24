package models

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

//CREATE TABLE `user` ( `id` int(11) NOT NULL AUTO_INCREMENT, `age` int(11) DEFAULT NULL, `name` varchar(255) DEFAULT NULL, `enabled` int(1) NOT NULL DEFAULT '1', PRIMARY KEY (`id`) ) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8;
type User struct {
	Id int
	Age int16
	Name string
}

var(
	db orm.Ormer
)


func init() {

	// open the debug or not
	// the debug statement will print the sql sentence
	orm.Debug = true

	orm.RegisterDriver("mysql", orm.DRMySQL)

	orm.RegisterModel(new(User)) // user_info UserInfo

	orm.RegisterDataBase("default", "mysql",
		"root:1234@tcp(127.0.0.1:3306)/test?charset=utf8", 30)

	db = orm.NewOrm()

}

// add a object
func AddUser(user_info *User)(int64, error) {
	id, err := db.Insert(user_info)
	return id, err
}

// read a object
func ReadUser(users *[]User) {

	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Select("*").From("user")
	sql := qb.String()
	// put the result to the users
	db.Raw(sql).QueryRows(users)

}









