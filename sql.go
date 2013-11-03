package main

import(
	"database/sql"
	_ "github.com/lib/pq"
	"fmt"
	"strings"
	"strconv"
)

var db *sql.DB
const (
	userTable = "users"
	sessionTable = "sessions"
)

func init() {
	var err error
	db, err = sql.Open("postgres","user=psql dbname=anonbitcoin password=psql")
	if err!=nil {
		fmt.Println(err.Error())
	}
	query:="CREATE TABLE IF NOT EXISTS "+userTable+"(userid serial PRIMARY KEY, username varchar(50) UNIQUE NOT NULL, password varchar(50) NOT NULL, balance decimal)"
	_,err = db.Exec(query)
	if err!=nil {
		fmt.Printf("User Table:" + err.Error())
	}
	query="CREATE TABLE IF NOT EXISTS "+sessionTable+"(sessionid serial PRIMARY KEY,  userid integer UNIQUE)"
	_,err = db.Exec(query)
	if err!=nil {
		fmt.Printf("Session Table:" + err.Error())
	}
}

func createUser(username, password string) bool {
	fmt.Println("Creating user "+username+"...")
	_,err:=db.Exec("INSERT INTO "+userTable+"(username,password,balance) VALUES ('"+username+"','"+password+"',0)")
	if err!=nil {
		fmt.Println(err.Error())
	}
	fmt.Println("User "+username+" created!")
	return true
}

func confirmUser(username, password string) (bool,int) {
	fmt.Println("Logging in user "+username+"...")
	rows,err:=db.Query("SELECT password,userid FROM "+userTable+" WHERE username='"+username+"'");
	if err!=nil {
		fmt.Println(err.Error())
		return false,0
	}
	if !rows.Next() {
		fmt.Println("NO RESULTS FOUND!")
		return false,0
	}

	var pwd string
	var uid int
	err=rows.Scan(&pwd,&uid)

	if err!=nil {
		fmt.Println(err.Error())
		return false,0
	}

	if !strings.EqualFold(pwd,password) {
		fmt.Println("Incorrect Password!")
		return false,0
	}

	fmt.Println("User "+username+" logged in!")
	return true,uid
}

func updateSessionID(userid int) (int,bool) {
	uid:=strconv.Itoa(userid)
	db.Exec("DELETE FROM "+sessionTable+" WHERE userid="+uid)
	result,err:=db.Exec("INSERT INTO "+sessionTable+"(userid) VALUES("+uid+")")
	fmt.Println(result)
	if err!=nil {
		fmt.Println(err.Error())
		return -1,false
	}
	return 0,true
}

type sqlUser struct {
	Userid int
	Username string
	Balance float32
}

func getSession(sessionid int) *sqlUser {
	fmt.Println("Looking up users session...")
	results,err:=db.Query("SELECT (sessions.userid,users.username,users.balance) FROM sessions INNER JOIN sessions.userid=users.userid WHERE sessions.sessionid="+strconv.Itoa(sessionid))
	if err!=nil {
		fmt.Println(err.Error())
		return nil
	}

	if !results.Next() {
		fmt.Println("ID Not Found")
		return nil
	}

	var uid int
	var uname string
	var bal float32

	err=results.Scan(&uid,&uname,&bal)
	if err!=nil {
		fmt.Println(err.Error())
		return nil
	}
	user:=sqlUser{
		Userid:uid,
		Username:uname,
		Balance:bal,
	}

	return &user;
}
