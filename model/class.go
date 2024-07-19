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
	Id            int    `json:"id"`
	Username      string `json:"username"`
	Password      string `json:"password"`
	Nama_lengkap  string `json:"nama_lengkap"`
	Alamat        string `json:"alamat"`
	Jenis_kelamin string `json:"jenis_kelamin"`
	Tgl_lahir     string `json:"tgl_lahir"`
	Email         string `json:"email"`
	No_telp       string `json:"no_telp"`
	Foto_profil   string `json:"foto_profil"`
}
