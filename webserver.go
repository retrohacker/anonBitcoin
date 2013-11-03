package main

import( "net/http"
	"html/template"
	"fmt"
)

func home(w http.ResponseWriter, r *http.Request) {
	t,_:=template.ParseFiles("html/home.html")
	t.Execute(w,nil)
}

func register(w http.ResponseWriter, r *http.Request) {
	t,_:=template.ParseFiles("html/register.html")
	t.Execute(w,nil)
}

type sregister struct {
	Success string
	Username string
}

func processRegistration(w http.ResponseWriter, r *http.Request) {
	t,_:=template.ParseFiles("html/processRegistration.html")
	uname:=r.FormValue("name")
	pwd:=r.FormValue("password")

	if len(uname)==0 {
		fmt.Fprintf(w,"No User")
		return;
	}

	if !createUser(uname,pwd) {
		fmt.Fprintf(w,"Registration with SQL failed")
	}

	tfill:=&sregister{
		Success:"REGISTERED!",
		Username:uname,
	}

	err:=t.Execute(w,tfill)
	if err!=nil {
		fmt.Fprintf(w,err.Error())
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	t,_:=template.ParseFiles("html/login.html")
	t.Execute(w,nil)
}

type slogin struct {
	Success string
	Username string
}

func processLogin(w http.ResponseWriter, r *http.Request) {
	t,_:=template.ParseFiles("html/processLogin.html")
	uname:=r.FormValue("name")
	pwd:=r.FormValue("password")

	if len(uname)==0 {
		fmt.Fprintf(w,"No User")
	}

	var success slogin
	err,uid:=confirmUser(uname,pwd)

	if !err {
		success=slogin {
			Success:"FAILED!",
			Username:uname,
		}
		t.Execute(w,success)
		return
	}

	//Create session
	_,err=updateSessionID(uid)
	if !err {
		fmt.Println("Could not create session!")
		return
	}
	success=slogin {
		Success:"SUCCESS!",
		Username:uname,
	}
	t.Execute(w,success)
}
