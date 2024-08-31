package dto

import "golang.org/x/crypto/bcrypt"

type User struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (u *User) HashAndSalt() *User {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	u.Password = string(hash)
	return u
}

func NewUser() *User {
	return &User{}
}
