package model

type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type User struct {
	Id              int         `json:"id"`
	Username        string      `json:"username"`
	Password        string      `json:"password"`
	Nama_lengkap    string      `json:"nama_lengkap"`
	Alamat          string      `json:"alamat"`
	Jenis_kelamin   string      `json:"jenis_kelamin"`
	Tgl_lahir       string      `json:"tgl_lahir"`
	Email           string      `json:"email"`
	No_telp         string      `json:"no_telp"`
	Foto_profil     string      `json:"foto_profil"`
	Ktp             string      `json:"ktp"`
	Kelas           int         `json:"kelas"`
	Status          string      `json:"status"`
	Tipe            int         `json:"tipe"`
	First_login     string      `json:"first_login"`
	Denied_by_admin string      `json:"denied_by_admin"`
	UserRole        []Role      `json:"user_role"`
	UserPrivilege   []Privilege `json:"user_privilege"`
}

type Asset struct {
	Id_asset          int      `json:"id_asset"`
	Nama              string   `json:"nama"`
	Id_asset_parent   int      `json:"id_asset_parent"`
	Id_perusahaan     int      `json:"id_perusahaan"`
	Id_join           string   `json:"id_join"`
	Tipe              string   `json:"tipe"`
	Nomor_legalitas   string   `json:"nomor_legalitas"`
	File_legalitas    string   `json:"file_legalitas"`
	Status_asset      string   `json:"status_asset"`
	Surat_kuasa       string   `json:"surat_kuasa"`
	Alamat            string   `json:"alamat"`
	Kondisi           string   `json:"kondisi"`
	Titik_koordinat   string   `json:"titik_koordinat"`
	Batas_koordinat   string   `json:"batas_koordinat"`
	Luas              float32  `json:"luas"`
	Nilai             float32  `json:"nilai"`
	Provinsi          string   `json:"provinsi"`
	Usage             string   `json:"usage"`
	Owned_by          int      `json:"owned_by"`
	Id_asset_child    string   `json:"id_asset_child"`
	Status_pengecekan string   `json:"status_pengecekan"`
	Status_verifikasi string   `json:"status_verifikasi"`
	Status_publik     string   `json:"status_publik"`
	Hak_akses         string   `json:"hak_akses"`
	Masa_sewa         string   `json:"masa_sewa"`
	Created_at        string   `json:"created_at"`
	Deleted_at        string   `json:"deleted_at"`
	LinkGambar        []string `json:"link_gambar"`
	TagsAssets        []string `json:"tags"`
	ChildAssets       []Asset
}

type Perusahaan struct {
	Id                  int    `json:"id_perusahaan"`
	Status              string `json:"status"`
	Nama                string `json:"nama"`
	Username            string `json:"username"`
	Lokasi              string `json:"lokasi"`
	Kelas               int    `json:"kelas"`
	Tipe                string `json:"tipe"`
	Dokumen_kepemilikan string `json:"dokumen_kepemilikan"`
	Dokumen_perusahaan  string `json:"dokumen_perusahaan"`
	Modal               string `json:"modal"`
	Deskripsi           string `json:"deskripsi"`
	CreatedAt           string `json:"created_at"`
	UserJoined          []User
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
	Registered_by         int    `json:"registered_by"`
	Lokasi                string `json:"lokasi"`
	Availability_surveyor string `json:"availability_surveyor"`
}

type SurveyRequest struct {
	Id_transaksi_jual_sewa int    `json:"id_transaksi_jual_sewa"`
	User_id                int    `json:"user_id"`
	Id_asset               int    `json:"id_asset"`
	Nama_asset             string `json:"nama_aset"`
	Created_at             string `json:"created_at"`
	Dateline               string `json:"dateline"`
	Status_request         string `json:"status_request"`
	Status_verifikasi      string `json:"status_verifikasi"`
	Data_lengkap           string `json:"data_lengkap"`
	Usage_old              string `json:"usage_old"`
	Usage_new              string `json:"usage_new"`
	Luas_old               string `json:"luas_old"`
	Luas_new               string `json:"luas_new"`
	Nilai_old              string `json:"nilai_old"`
	Nilai_new              string `json:"nilai_new"`
	Kondisi_old            string `json:"kondisi_old"`
	Kondisi_new            string `json:"kondisi_new"`
	Batas_koordinat_old    string `json:"batas_koordinat_old"`
	Batas_koordinat_new    string `json:"batas_koordinat_new"`
	Tags_old               string `json:"tags_old"`
	Tags_new               string `json:"tags_new"`
}

type TransactionRequest struct {
	Id_transaksi_jual_sewa int    `json:"id_transaksi_jual_sewa"`
	Perusahaan_id          int    `json:"perusahaan_id"`
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

type UserSurveyor struct {
	User_id               int    `json:"user_id"`
	Username              string `json:"username"`
	Password              string `json:"password"`
	Nama_lengkap          string `json:"nama_lengkap"`
	Alamat                string `json:"alamat"`
	Jenis_kelamin         string `json:"jenis_kelamin"`
	Tgl_lahir             string `json:"tgl_lahir"`
	Email                 string `json:"email"`
	No_telp               string `json:"no_telp"`
	Foto_profil           string `json:"foto_profil"`
	Ktp                   string `json:"ktp"`
	Surveyor_id           int    `json:"surveyor_id"`
	Registered_by         int    `json:"registered_by"`
	Lokasi                string `json:"lokasi"`
	Availability_surveyor string `json:"availability_surveyor"`
	SurveyOnProgress      int    `json:"surveyonprogress"`
	FinishedSurvey        int    `json:"finished_survey"`
	TotalSurvey           int    `json:"totalsurvey"`
	Survey_Request        []SurveyRequest
}

type UserPerusahaan struct {
	Perusahaan_id    int    `json:"perusahaan_id"`
	Name             string `json:"username"`
	UserCount        string `json:"usercount"`
	TransactionCount string `json:"transactioncount"`
}

type RegisSurveyor struct {
	Id              int         `json:"id"`
	Registered_by   int         `json:"registered_by"`
	Username        string      `json:"username"`
	Password        string      `json:"password"`
	Nama_lengkap    string      `json:"nama_lengkap"`
	Alamat          string      `json:"alamat"`
	Jenis_kelamin   string      `json:"jenis_kelamin"`
	Tgl_lahir       string      `json:"tgl_lahir"`
	Email           string      `json:"email"`
	No_telp         string      `json:"no_telp"`
	Foto_profil     string      `json:"foto_profil"`
	Ktp             string      `json:"ktp"`
	Kelas           int         `json:"kelas"`
	Status          string      `json:"status"`
	Tipe            int         `json:"tipe"`
	First_login     string      `json:"first_login"`
	Denied_by_admin string      `json:"denied_by_admin"`
	UserRole        []Role      `json:"user_role"`
	UserPrivilege   []Privilege `json:"user_privilege"`
}
