package Mysql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"golandproject/Class"
	"golandproject/middleware/Geoip2"
)

var db *sql.DB

func init() {
	var err error
	dsn := "root:SUIbianla123@tcp(127.0.0.1:3306)/project"
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		fmt.Printf("%#v", err)
	}
	err = db.Ping()
	if err != nil {
		fmt.Printf("%#v", err)
	}
}

func AddUser(uid string, password string, email string) bool {
	tx, _ := db.Begin()
	sqlStr := "insert into users(uid,password,email,create_time) values(?,?,?,NOW())"
	_, err := tx.Exec(sqlStr, uid, password, email)
	if err != nil {
		return false
	}
	tx.Commit()
	return true
}

func GetInfo(uid string) (info Class.User) {
	sqlStr := "select * from users where uid=?"
	_ = db.QueryRow(sqlStr, uid).Scan(&info.Uid, &info.Password, &info.Email, &info.Create_time)
	return info
}

func VerifyPassword(uid string, password string) bool {
	sqlStr := "select password from users where uid=?"
	var user Class.User
	_ = db.QueryRow(sqlStr, uid).Scan(&user.Password)
	return user.Password == password
}

func EditPassword(uid string, password string) bool {
	tx, _ := db.Begin()
	sqlStr := "update users set password=? where uid=?"
	_, err := tx.Exec(sqlStr, password, uid)
	if err != nil {
		return false
	}
	tx.Commit()
	return true
}

func AddLoginRecord(uid, ip string) {
	tx, _ := db.Begin()
	//fmt.Printf("%#v %#v", uid, ip)
	sqlStr := "insert into iprecord(uid,ip,time,city) values(?,?,NOW(),?)"
	_, err := tx.Exec(sqlStr, uid, ip, Geoip2.IpQueryCity(ip))
	if err != nil {
		return
	}
	tx.Commit()
}

func GetIpRecord(uid string) (iprecord []Class.IpRecord) {
	rows, err := db.Query("select uid,time,ip,city from iprecord where uid=? ORDER BY time DESC", uid)
	if err != nil {
		return
	}
	for rows.Next() {
		var temp Class.IpRecord
		rows.Scan(&temp.Uid, &temp.Time, &temp.Ip, &temp.City)
		iprecord = append(iprecord, temp)
	}
	return iprecord
}

func FindAllUsers() {
	//tx, err := db.Begin()
	//if err != nil {
	//	tx.Rollback()
	//	fmt.Printf("%#v:%#v", "[ERROR]", err)
	//	return
	//}
	//sqlStr1 := "select * from users"
	//result, err := tx.Exec(sqlStr1)
	//num, _ := result.RowsAffected()
	//if num == 0 {
	//	fmt.Printf("%#v", "影响0行")
	//}
	//defer db.Close()
	//tx.Commit()
	rows, err := db.Query("select * from users")
	if err != nil {
		return
	}
	for rows.Next() {
		//var user Class.User
		//rows.Scan(&user.Uid, &user.Date, &user.U_name, &user.U_password, &user.U_id, &user.Class, &user.Level)
		//rows.Scan(&user.Uid, &user.Password, &user.FirstName, &user.LastName, &user.Sex)
		//fmt.Println(user)
	}
}
