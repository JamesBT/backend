package model

type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type User struct {
	Id               int    `json:"id"`
	Username         string `json:"username"`
	Password         string `json:"password"`
	Nama_lengkap     string `json:"nama_lengkap"`
	Alamat           string `json:"alamat"`
	Jenis_kelamin    string `json:"jenis_kelamin"`
	Tgl_lahir        string `json:"tgl_lahir"`
	Email            string `json:"email"`
	No_telp          string `json:"no_telp"`
	Foto_profil      string `json:"foto_profil"`
	Ktp              string `json:"ktp"`
	Kelas            int    `json:"kelas"`
	Status           string `json:"status"`
	Tipe             int    `json:"tipe"`
	First_login      string `json:"first_login"`
	Denied_by_admin  string `json:"denied_by_admin"`
	UserRole         []Role `json:"user_role"`
	PerusahaanJoined []Perusahaan
}

type Asset struct {
	Id_asset          int        `json:"id_asset"`
	Nama              string     `json:"nama"`
	Id_asset_parent   int        `json:"id_asset_parent"`
	Id_join           string     `json:"id_join"`
	Tipe              string     `json:"tipe"`
	Nomor_legalitas   string     `json:"nomor_legalitas"`
	File_legalitas    string     `json:"file_legalitas"`
	Status_asset      string     `json:"status_asset"`
	Surat_kuasa       string     `json:"surat_kuasa"`
	Surat_legalitas   string     `json:"surat_legalitas"`
	Alamat            string     `json:"alamat"`
	Kondisi           string     `json:"kondisi"`
	Titik_koordinat   string     `json:"titik_koordinat"`
	Batas_koordinat   string     `json:"batas_koordinat"`
	Luas              float64    `json:"luas"`
	Nilai             float64    `json:"nilai"`
	Provinsi          int        `json:"provinsi"`
	Owned_by          int        `json:"owned_by"`
	Id_asset_child    string     `json:"id_asset_child"`
	Status_pengecekan string     `json:"status_pengecekan"`
	Status_verifikasi string     `json:"status_verifikasi"`
	Status_publik     string     `json:"status_publik"`
	Hak_akses         string     `json:"hak_akses"`
	Masa_sewa         string     `json:"masa_sewa"`
	Created_at        string     `json:"created_at"`
	Deleted_at        string     `json:"deleted_at"`
	LinkGambar        []string   `json:"link_gambar"`
	TagsAssets        []Tags     `json:"tags"`
	Usage             []Kegunaan `json:"usage"`
	ChildAssets       []Asset
}

type BusinessField struct {
	Id     int    `json:"id"`
	Nama   string `json:"nama"`
	Detail string `json:"detail"`
}

type Perusahaan struct {
	Id                  int     `json:"id_perusahaan"`
	Status              string  `json:"status"`
	Nama                string  `json:"nama"`
	Username            string  `json:"username"`
	Lokasi              string  `json:"lokasi"`
	Kelas               int     `json:"kelas"`
	Tipe                string  `json:"tipe"`
	Dokumen_kepemilikan string  `json:"dokumen_kepemilikan"`
	Dokumen_perusahaan  string  `json:"dokumen_perusahaan"`
	Modal               float64 `json:"modal"`
	Deskripsi           string  `json:"deskripsi"`
	CreatedAt           string  `json:"created_at"`
	Field               []BusinessField
	UserJoined          []User
}

type Privilege struct {
	Privilege_id   int    `json:"privilege_id"`
	Nama_privilege string `json:"nama_privilege"`
}

type Role struct {
	Role_id   int      `json:"role_id"`
	Nama_role string   `json:"nama_role"`
	Privilege []string `json:"privilege"`
}

type Surveyor struct {
	Surveyor_id           int    `json:"surveyor_id"`
	User_id               int    `json:"user_id"`
	Registered_by         int    `json:"registered_by"`
	Lokasi                string `json:"lokasi"`
	Availability_surveyor string `json:"availability_surveyor"`
}

type SurveyRequest struct {
	Id_transaksi_jual_sewa int        `json:"id_transaksi_jual_sewa"`
	User_id                int        `json:"user_id"`
	Id_asset               int        `json:"id_asset"`
	Nama_asset             string     `json:"nama_aset"`
	Lokasi_asset           string     `json:"lokasi_asset"`
	Tipe_asset             string     `json:"tipe_asset"`
	Created_at             string     `json:"created_at"`
	Surat_penugasan        string     `json:"surat_penugasan"`
	Dateline               string     `json:"dateline"`
	Status_request         string     `json:"status_request"`
	Status_verifikasi      string     `json:"status_verifikasi"`
	Status_submitted       string     `json:"status_submit"`
	Data_lengkap           string     `json:"data_lengkap"`
	Usage_old              []Kegunaan `json:"usage_old"`
	Usage_new              []Kegunaan `json:"usage_new"`
	Luas_old               string     `json:"luas_old"`
	Luas_new               string     `json:"luas_new"`
	Nilai_old              float64    `json:"nilai_old"`
	Nilai_new              float64    `json:"nilai_new"`
	Kondisi_old            string     `json:"kondisi_old"`
	Kondisi_new            string     `json:"kondisi_new"`
	Titik_koordinat_old    string     `json:"titik_koordinat_old"`
	Titik_koordinat_new    string     `json:"titik_koordinat_new"`
	Batas_koordinat_old    string     `json:"batas_koordinat_old"`
	Batas_koordinat_new    string     `json:"batas_koordinat_new"`
	Tags_old               []Tags     `json:"tags_old"`
	Tags_new               []Tags     `json:"tags_new"`
}

type TransactionRequest struct {
	Id_transaksi_jual_sewa int    `json:"id_transaksi_jual_sewa"`
	Perusahaan_id          int    `json:"perusahaan_id"`
	Lokasi_perusahaan      string `json:"lokasi_perusahaan"`
	User_id                int    `json:"user_id"`
	Username               string `json:"username"`
	Nama_lengkap           string `json:"nama_lengkap"`
	Id_asset               int    `json:"id_asset"`
	Nama_aset              string `json:"nama_aset"`
	Status                 string `json:"status"`
	Nama_progress          string `json:"nama_progress"`
	Proposal               string `json:"proposal"`
	Tgl_meeting            string `json:"tgl_meeting"`
	Waktu_meeting          string `json:"waktu_meeting"`
	Lokasi_meeting         string `json:"lokasi_meeting"`
	Deskripsi              string `json:"deskripsi"`
	Alasan                 string `json:"alasan"`
	Tgl_dateline           string `json:"tgl_dateline"`
	Created_at             string `json:"created_at"`
}

type UserRole struct {
	User_role_id int `json:"user_role_id"`
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
	Name             string `json:"nama"`
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

type Notification struct {
	Notification_id        int    `json:"notification_id"`
	User_id_sender         int    `json:"user_id_sender"`
	User_id_receiver       int    `json:"user_id_receiver"`
	Perusahaan_id_receiver int    `json:"perusahaan_id_receiver"`
	Created_at             string `json:"created_at"`
	Title                  string `json:"notification_title"`
	Detail                 string `json:"notification_detail"`
}

type Kelas struct {
	Id             int     `json:"id"`
	Nama           string  `json:"nama"`
	Modal_minimal  float64 `json:"modal_minimal"`
	Modal_maksimal float64 `json:"modal_maksimal"`
}

type Kegunaan struct {
	Id   int    `json:"id"`
	Nama string `json:"nama"`
}

type Tags struct {
	Id     int    `json:"id"`
	Nama   string `json:"nama"`
	Detail string `json:"detail"`
}

type Provinsi struct {
	Id   int    `json:"id"`
	Nama string `json:"nama"`
}

type Progress struct {
	Id                    int    `json:"id"`
	User_id               int    `json:"user_id"`
	Perusahaan_id         int    `json:"perusahaan_id"`
	Id_asset              int    `json:"id_asset"`
	Nama_asset            string `json:"nama_asset"`
	Status                string `json:"status"`
	Data_lengkap          string `json:"data_lengkap"`
	Nama                  string `json:"nama"`
	Proposal              string `json:"proposal"`
	Tanggal_meeting       string `json:"tgl_meeting"`
	Waktu_meeting         string `json:"waktu_meeting"`
	Tempat_meeting        string `json:"tempat_meeting"`
	Waktu_mulai_meeting   string `json:"waktu_mulai"`
	Waktu_selesai_meeting string `json:"waktu_selesai"`
	Notes                 string `json:"notes"`
	Dokumen               string `json:"dokumen"`
	Tipe_dokumen          string `json:"tipe"`
}

type GrupAsset struct {
	Id_asset       int        `json:"id_asset"`
	Asset_name     string     `json:"nama_asset"`
	User_id        int        `json:"user_id"`
	Perusahaan_id  int        `json:"perusahaan_id"`
	Semua_progress []Progress `json:"progress"`
}
