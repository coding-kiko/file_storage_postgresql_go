package main

import (
	"encoding/json"
	"io/ioutil"
)

type Users struct {
	Users []User `json:"users"`
}

type User struct {
	Name  string `json:"email"`
	Pwd   string `json:"pwd"`
	Token string `json:"jwt"`
}

var (
	filename = "./cache/users.json"
)

func main() {
	var users Users
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	json.Unmarshal(data, &users)
	new := User{Name: "Cuatro", Pwd: "Cuatro"}
	users.Users = append(users.Users, new)

	replace := Users{Users: users.Users}
	data, _ = json.MarshalIndent(replace, "", "\t")
	_ = ioutil.WriteFile(filename, data, 0644)

}
