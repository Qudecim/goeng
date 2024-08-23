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
