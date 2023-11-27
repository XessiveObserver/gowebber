package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"

	"golang.org/x/crypto/bcrypt"
)

func isAuthenticated(r *http.Request) bool {
	session, err := store.Get(r, "session")
	if err != nil {
		return false
	}

	auth, ok := session.Values["authenticated"].(bool)
	return ok && auth
}

// renderTemplate function loads all templates
func renderTemplate(w http.ResponseWriter, tmpl string, data TemplateData) {
	baseTemplate := "templates/base.html"
	tmpl = fmt.Sprintf("templates/%s", tmpl)

	// Parse the base template and the specified template
	t, err := template.ParseFiles(baseTemplate, tmpl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Execute the template with the provided data
	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// index func displays home page
func index(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
	// Render index.html template
	data := TemplateData{IsAuthenticated: isAuthenticated(r)}
	renderTemplate(w, "index.html", data)
}

// profiles function list all profiles
func profiles(w http.ResponseWriter, r *http.Request) {

	// Check if the method is allowed
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Authentication
	session, err := store.Get(r, "session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if the user is authenticated
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		log.Println("Not authenticated. Redirecting to /login")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Query database for all profiles
	rows, err := DB.Query("SELECT * FROM profiles")
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	var profiles []Profile
	for rows.Next() {
		var profile Profile
		if err := rows.Scan(&profile.ID, &profile.Name, &profile.Age, &profile.Occupation); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		profiles = append(profiles, profile)
	}

	// Pass authentication status to the template
	data := TemplateData{
		IsAuthenticated: isAuthenticated(r),
		Profiles:        profiles,
	}
	// Render profiles.html with profile listing
	renderTemplate(w, "profiles.html", data)
}

// profile function lists a single Profile
func profile(w http.ResponseWriter, r *http.Request) {
	// Pass authentication status to the template

	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Authentication
	session, err := store.Get(r, "session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if the user is authenticated
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		log.Println("Not authenticated. Redirecting to /login")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Parse profile id from URK
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Profile ID not provided", http.StatusBadRequest)
		return
	}

	// Query  database to rerieve profile details by ID
	var profile Profile
	err = DB.QueryRow("SELECT id, name, age, occupation FROM profiles WHERE id = $1", id).Scan(&profile.ID, &profile.Name, &profile.Age, &profile.Occupation)
	if err != nil {
		http.Error(w, "Profile not found", http.StatusNotFound)
		return
	}

	// Render profile.html template with profile details
	data := TemplateData{
		IsAuthenticated: isAuthenticated(r),
		Profile:         profile,
	}
	renderTemplate(w, "profile.html", data)
}

// addprofile function creates a new profile in the database
func addProfile(w http.ResponseWriter, r *http.Request) {
	// Pass authentication status to the template
	data := TemplateData{IsAuthenticated: isAuthenticated(r)}

	// Authentication
	session, err := store.Get(r, "session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if the user is authenticated
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		log.Println("Not authenticated. Redirecting to /login")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	switch r.Method {
	case http.MethodGet:
		// Render addprofile.html template for GET requests
		renderTemplate(w, "addprofile.html", data)
	case http.MethodPost:
		name := r.FormValue("name")
		age := r.FormValue("age")
		occupation := r.FormValue("occupation")

		// Insert new profile into database
		_, err := DB.Exec("INSERT INTO profiles (name, age, occupation) VALUES ($1, $2, $3)", name, age, occupation)
		if err != nil {
			http.Error(w, "Failed to insert into database", http.StatusInternalServerError)
			return
		}

		// Redirect to /profiles after successful addition
		http.Redirect(w, r, "/profiles", http.StatusSeeOther)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

// editprofile function edits a profile
func editProfile(w http.ResponseWriter, r *http.Request) {

	// Authentication
	session, err := store.Get(r, "session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if the user is authenticated
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		log.Println("Not authenticated. Redirecting to /login")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	switch r.Method {
	case http.MethodGet:
		// Check if profile exists using ID
		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, "Profile ID not provided", http.StatusBadRequest)
			return
		}

		// Query the database to retrieve profile details by ID
		var profile Profile
		err := DB.QueryRow("SELECT id, name, age, occupation FROM profiles WHERE id = $1", id).Scan(&profile.ID, &profile.Name, &profile.Age, &profile.Occupation)
		if err != nil {
			http.Error(w, "Profile not found", http.StatusNotFound)
			return
		}
		// Pass authentication status to the template
		data := TemplateData{
			IsAuthenticated: isAuthenticated(r),
			Profile:         profile,
		}

		// Render editprofile.html template for GET requests
		renderTemplate(w, "editprofile.html", data)
	case http.MethodPost:
		// Handle form submission and update profile for POST requests
		id := r.FormValue("id")
		name := r.FormValue("name")
		age := r.FormValue("age")
		occupation := r.FormValue("occupation")

		_, err := DB.Exec("UPDATE profiles SET name = $1, age = $2, occupation = $3 WHERE id = $4", name, age, occupation, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Redirect to /profiles after successful update
		http.Redirect(w, r, "/profiles", http.StatusSeeOther)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

// deleteProfile deletes the profile
func deleteProfile(w http.ResponseWriter, r *http.Request) {

	// Authentication
	session, err := store.Get(r, "session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if the user is authenticated
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		log.Println("Not authenticated. Redirecting to /login")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	switch r.Method {
	case http.MethodGet:
		// Check if profile exists using ID
		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, "Profile ID not provided", http.StatusBadRequest)
			return
		}

		// Query database to retrieve profile details by ID
		var profile Profile
		err := DB.QueryRow("SELECT id, name, age, occupation FROM profiles WHERE id = $1", id).Scan(&profile.ID, &profile.Name, &profile.Age, &profile.Occupation)
		if err != nil {
			http.Error(w, "Profile not found", http.StatusNotFound)
			return
		}

		// Pass authentication status to the template
		data := TemplateData{
			IsAuthenticated: isAuthenticated(r),
			Profile:         profile,
		}

		// Render the deleteprofile.html template for GET requests
		renderTemplate(w, "deleteprofile.html", data)
	case http.MethodPost:
		// Handle form submission and delete profile for POST requests
		id := r.FormValue("id")

		_, err := DB.Exec("DELETE FROM profiles WHERE id = $1", id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Redirect to /profiles after successful deletion
		http.Redirect(w, r, "/profiles", http.StatusSeeOther)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

// Authorization

func register(w http.ResponseWriter, r *http.Request) {
	// Pass authentication status to the template
	data := TemplateData{IsAuthenticated: isAuthenticated(r)}

	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		// Hash the password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			log.Println(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		_, err = DB.Exec("INSERT INTO users(username, password) VALUES($1, $2)", username, string(hashedPassword))
		if err != nil {
			log.Println(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Display registration form
	renderTemplate(w, "register.html", data)
}

func login(w http.ResponseWriter, r *http.Request) {
	// Pass authentication status to the template
	data := TemplateData{IsAuthenticated: isAuthenticated(r)}

	session, err := store.Get(r, "session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		var storedPassword string
		err := DB.QueryRow("SELECT password FROM users WHERE username=$1", username).Scan(&storedPassword)
		if err != nil {
			log.Println(err)
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(password))
		if err != nil {
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			return
		}

		// Set session values
		session.Values["authenticated"] = true
		session.Values["username"] = username
		err = session.Save(r, w)
		if err != nil {
			log.Println(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Redirect to the profiles page after successful login
		http.Redirect(w, r, "/profiles", http.StatusSeeOther)
		return
	}

	// Display login form
	renderTemplate(w, "login.html", data)
}

func logout(w http.ResponseWriter, r *http.Request) {
	// Get the session
	session, err := store.Get(r, "session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Clear all session values
	for key := range session.Values {
		delete(session.Values, key)
	}

	// Save the session
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Redirect to the login page
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
