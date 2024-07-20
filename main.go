package main

import (
	"fmt"
	"log"
	"net/http"

	db "user_admin/DB"
	handlers "user_admin/Handlers"
)

func main() {

	//User
	http.HandleFunc("/", handlers.RootHandler)
	http.HandleFunc("/signup", handlers.SignupHandler)
	http.HandleFunc("/home", handlers.HomeHandler)
	http.HandleFunc("/logout", handlers.LogoutHandler)

	//Admin
	http.HandleFunc("/admin", handlers.Admin)
	http.HandleFunc("/adminAddUser", handlers.AdminAddUser)
	http.HandleFunc("/adminUserUpdate", handlers.AdminUserUpdate)
	http.HandleFunc("/adminUserDelete", handlers.AdminUserDelete)
	http.HandleFunc("/adminLogout", handlers.AdminLogout)
	http.HandleFunc("/adminSearchUser", handlers.AdminSearchUser)

	fmt.Printf("Starting server at port 8080\n")
	db.Init()
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}

}
