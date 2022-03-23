package auth

type Users struct {
	Users []User `json:"users"`
}

type User struct {
	Email string `json:"email"`
	Pwd   string `json:"pwd"`
	Token string `json:"jwt"`
}

type Credentials struct {
	Email string `json:"email"`
	Pwd   string `json:"pwd"`
}
