package models

// API user credentials
// It is used to sign in
//
// swagger:model user
type User struct {
	// User's login
	//
	// required: true
	Username string `json:"username"`
	// User's password
	//
	// required: true
	Password string `json:"password"`
}
