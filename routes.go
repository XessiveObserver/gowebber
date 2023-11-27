package main

import (
	"fmt"
	"log"
	"net/http"
)

func routing() {
	// Static file routing
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/", index)
	http.HandleFunc("/profiles", profiles)
	http.HandleFunc("/profile", profile)
	http.HandleFunc("/addprofile", addProfile)
	http.HandleFunc("/editprofile", editProfile)
	http.HandleFunc("/deleteprofile", deleteProfile)

	http.HandleFunc("/register", register)
	http.HandleFunc("/login", login)
	http.HandleFunc("/logout", logout)

	fmt.Println("Server Running on port: 8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
