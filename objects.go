package goeng

type Dict struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type Word struct {
	ID     int64  `json:"id"`
	First  string `json:"first"`
	Second string `json:"second"`
}

type User struct {
	ID           int64  `json:"id"`
	userName     string `json:"user_name"`
	passwordHash string `json:"password_hash"`
	salt         string `json:"salt"`
}

type DtoUser struct {
	userName string `json:"user_name"`
	password string `json:"password"`
}
