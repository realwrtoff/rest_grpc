package model

type Account struct {
	Username     string `json:"username"`
	Password string `json:"password"`
}

type Token struct {
	Username  string `json:"username"`
	Token string `json:"token"`
}

func (Account) TableName() string {
	return "account"
}

func (Token) TableName() string {
	return "token"
}