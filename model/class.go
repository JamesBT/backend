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
	Ktp           string `json:"ktp"`
}

type Asset struct {
	Id_asset_parent   int     `json:"id_asset_parent"`
	Nama              string  `json:"nama"`
	Nama_legalitas    string  `json:"nama_legalitas"`
	Nomor_legalitas   string  `json:"nomor_legalitas"`
	Tipe              string  `json:"tipe"`
	Nilai             int     `json:"nilai"`
	Luas              float32 `json:"luas"`
	Titik_koordinat   string  `json:"titik_koordinat"`
	Batas_koordinat   string  `json:"batas_koordinat"`
	Kondisi           string  `json:"kondisi"`
	Id_asset_child    string  `json:"id_asset_child"`
	Alamat            string  `json:"alamat"`
	Status_pengecekan string  `json:"status_pengecekan"`
	Status_verifikasi string  `json:"status_verifikasi"`
	Hak_akses         string  `json:"hak_akses"`
	Status_asset      string  `json:"status_asset"`
	Masa_sewa         string  `json:"masa_sewa"`
}

type Perusahaan struct {
	Perusahaan_id         int    `json:"perusahaan_id"`
	User_id               int    `json:"user_id"`
	Sertifikat_perusahaan string `json:"sertifikat_perusahaan"`
}

type Privilege struct {
	Privilege_id   int    `json:"privilege_id"`
	Nama_privilege string `json:"nama_privilege"`
}

type Role struct {
	Role_id   int    `json:"role_id"`
	Nama_role string `json:"nama_role"`
}

type Surveyor struct {
	Surveyor_id           int    `json:"surveyor_id"`
	User_id               int    `json:"user_id"`
	Lokasi                string `json:"lokasi"`
	Availability_surveyor int    `json:"availability_surveyor"`
}

type SurveyRequest struct {
	Id_transaksi_jual_sewa int    `json:"id_transaksi_jual_sewa"`
	User_id                int    `json:"user_id"`
	Id_asset               int    `json:"id_asset"`
	Dateline               string `json:"dateline"`
	Status_request         string `json:"status_request"`
}

type TransactionRequest struct {
	Id_transaksi_jual_sewa int    `json:"id_transaksi_jual_sewa"`
	User_id                int    `json:"user_id"`
	Id_asset               int    `json:"id_asset"`
	Tipe                   string `json:"tipe"`
	Masa_sewa              string `json:"masa_sewa"`
	Meeting_log            string `json:"meeting_log"`
}

type UserPrivilege struct {
	User_privilege_id int `json:"user_privilege_id"`
	Privilege_id      int `json:"privilege_id"`
	User_id           int `json:"user_id"`
}

type UserRole struct {
	User_role_id int `json:"user_privilege_id"`
	User_id      int `json:"user_id"`
	Role_id      int `json:"role_id"`
}
