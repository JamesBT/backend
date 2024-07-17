package model

type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type Barang struct {
	Id    int    `json:"id"`
	Nama  string `json:"nama"`
	Harga int    `json:"harga"`
}

type User struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}
