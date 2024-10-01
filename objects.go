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
	UserName     string `json:"user_name"`
	PasswordHash string `json:"password_hash"`
	Salt         string `json:"salt"`
}

type DtoUser struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

type DtoError struct {
	Code    int    `json:"code"`
	Field   string `json:"field"`
	Message string `json:"message"`
}

type DtoSuccess struct {
	Success bool `json:"success"`
}
