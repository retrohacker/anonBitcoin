package main

import(
	"net/http"
)

func main() {
	/*
	exec.LookPath("/home/crackerz/bin")
	bitcoind := exec.Command("bitcoind","getaccountaddress","")
	output,err := bitcoind.Output()

	if err==nil {
		fmt.Println(string(output))
	} else {
		fmt.Println("Error: ",err.Error())
		fmt.Println(bitcoind)
	}
	*/
	http.HandleFunc("/",home)
	http.HandleFunc("/register",register)
	http.HandleFunc("/register/submit",processRegistration)
	http.HandleFunc("/login",login)
	http.HandleFunc("/login/submit",processLogin)
	http.ListenAndServe(":8080",nil)
}
