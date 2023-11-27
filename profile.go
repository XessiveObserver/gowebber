package main

// User Profiles
type Profile struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Age        int    `json:"age"`
	Occupation string `json:"occupation"`
}

// User account
type User struct {
	ID       int
	Username string
	Password string
}

type TemplateData struct {
	IsAuthenticated bool
	Profiles        []Profile
	Profile         Profile
	// Add other data you want to pass to templates
}
