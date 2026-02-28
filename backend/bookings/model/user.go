package model

type User struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewUser(username string, password string) User {
	return User{
		Username: username,
		Password: password,
	}
}

func (User) TableName() string {
	return "usertable"
}
