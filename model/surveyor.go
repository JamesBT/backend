package model

import (
	"TemplateProject/db"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

// CRUD surveyor ============================================================================
func LoginSurveyor(akun string) (Response, error) {
	var res Response

	var usr = User{}
	var loginUsr = User{}

	err := json.Unmarshal([]byte(akun), &usr)
	if err != nil {
		res.Status = 401
		res.Message = "gagal decode json"
		res.Data = err.Error()
		return res, err
	}

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	// cek sudah terdaftar atau belum
	query := "SELECT user_id FROM user WHERE username = ? AND deleted_at IS NULL"
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}

	var userId int
	err = stmt.QueryRow(usr.Username).Scan(&userId)
	if err != nil {
		res.Status = 401
		res.Message = "Pengguna belum terdaftar atau telah dihapus"
		res.Data = err.Error()
		return res, errors.New("pengguna belum terdaftar atau telah dihapus")
	}
	defer stmt.Close()

	// cek apakah password benar atau tidak
	queryinsert := "SELECT u.user_id, u.username, u.nama_lengkap, u.alamat, u.jenis_kelamin, u.tanggal_lahir, u.email, u.nomor_telepon, u.foto_profil, u.ktp, ud.user_kelas_id, ud.status, ud.tipe, ud.first_login, ud.denied_by_admin, s.lokasi, s.availability_surveyor FROM user u JOIN user_detail ud ON u.user_id = ud.user_detail_id JOIN surveyor s ON u.user_id = s.user_id WHERE u.username = ? AND u.password = ?;"
	stmtinsert, err := con.Prepare(queryinsert)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmtinsert.Close()

	var lokasi string
	var availability_surveyor string
	err = stmtinsert.QueryRow(usr.Username, usr.Password).Scan(&loginUsr.Id, &loginUsr.Username, &loginUsr.Nama_lengkap, &loginUsr.Alamat, &loginUsr.Jenis_kelamin, &loginUsr.Tgl_lahir, &loginUsr.Email, &loginUsr.No_telp, &loginUsr.Foto_profil, &loginUsr.Ktp, &loginUsr.Kelas, &loginUsr.Status, &loginUsr.Tipe, &loginUsr.First_login, &loginUsr.Denied_by_admin, &lokasi, &availability_surveyor)
	if err != nil {
		res.Status = 401
		res.Message = "password salah"
		res.Data = err.Error()
		return res, errors.New("password salah")
	}

	// ambil role + privilege
	getRoleQuery := "SELECT ur.role_id, r.nama_role FROM user_role ur JOIN role r ON ur.role_id = r.role_id WHERE ur.user_id = ?;"
	rolestmt, err := con.Prepare(getRoleQuery)
	if err != nil {
		res.Status = 401
		res.Message = "stmt update gagal"
		res.Data = err.Error()
		return res, err
	}
	defer rolestmt.Close()

	var roleId int
	var roleName string
	err = rolestmt.QueryRow(userId).Scan(&roleId, &roleName)
	if err != nil {
		res.Status = 401
		res.Message = "gagal mendapatkan role"
		res.Data = err.Error()
		return res, err
	}

	getPrivilegeQuery := "SELECT pr.privilege_id, p.nama_privilege FROM user_privilege pr JOIN privilege p ON pr.privilege_id = p.privilege_id WHERE pr.user_id = ?;"
	privilegestmt, err := con.Prepare(getPrivilegeQuery)
	if err != nil {
		res.Status = 401
		res.Message = "stmt update gagal"
		res.Data = err.Error()
		return res, err
	}
	defer privilegestmt.Close()

	var privilegeId int
	var privilegeName string
	err = privilegestmt.QueryRow(userId).Scan(&privilegeId, &privilegeName)
	if err != nil {
		res.Status = 401
		res.Message = "gagal mendapatkan privilege"
		res.Data = err.Error()
		return res, err
	}

	// berhasil login => update timestamp terakhir login
	updateQuery := "UPDATE user SET login_timestamp = NOW() WHERE user_id = ?"
	updatestmt, err := con.Prepare(updateQuery)
	if err != nil {
		res.Status = 401
		res.Message = "stmt update gagal"
		res.Data = err.Error()
		return res, err
	}
	defer updatestmt.Close()

	_, err = updatestmt.Exec(userId)
	if err != nil {
		res.Status = 401
		res.Message = "update login_timestamp gagal"
		res.Data = err.Error()
		return res, err
	}

	res.Status = http.StatusOK
	res.Message = "Berhasil login"
	res.Data = map[string]interface{}{
		"id":              loginUsr.Id,
		"username":        loginUsr.Username,
		"nama_lengkap":    loginUsr.Nama_lengkap,
		"alamat":          loginUsr.Alamat,
		"jenis_kelamin":   loginUsr.Jenis_kelamin,
		"tanggal_lahir":   loginUsr.Tgl_lahir,
		"email":           loginUsr.Email,
		"nomor_telepon":   loginUsr.No_telp,
		"foto_profil":     loginUsr.Foto_profil,
		"ktp":             loginUsr.Ktp,
		"status":          loginUsr.Status,
		"tipe":            loginUsr.Tipe,
		"first_login":     loginUsr.First_login,
		"denied_by_admin": loginUsr.Denied_by_admin,
		"role_id":         roleId,
		"role_nama":       roleName,
		"privilege_id":    privilegeId,
		"nama_privilege":  privilegeName,
	}

	defer db.DbClose(con)

	return res, nil
}

func SignUpSurveyor(akun string) (Response, error) {
	var res Response
	var usr = RegisSurveyor{}

	err := json.Unmarshal([]byte(akun), &usr)
	if err != nil {
		res.Status = 401
		res.Message = "gagal decode json"
		res.Data = err.Error()
		return res, err
	}

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	// cek sudah terdaftar atau belum
	query := "SELECT user_id FROM user WHERE `username` = ?"
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}

	var userId int64
	err = stmt.QueryRow(usr.Username).Scan(&userId)
	if err != nil && err != sql.ErrNoRows {
		res.Status = 500
		res.Message = "Query execution failed"
		res.Data = err.Error()
		return res, err
	}
	// fmt.Println("username:", usr.Username, "id: ", userId)
	if err == nil {
		// User sudah kedaftar
		res.Status = 401
		res.Message = "User already registered"
		res.Data = "User ID: " + fmt.Sprint(userId)
		return res, errors.New("user already registered")
	}
	defer stmt.Close()

	// // cek registered by
	// registeredquery := "SELECT user_id FROM user WHERE user_id = ?"
	// registeredstmt, err := con.Prepare(registeredquery)
	// if err != nil {
	// 	res.Status = 401
	// 	res.Message = "stmt gagal"
	// 	res.Data = err.Error()
	// 	return res, err
	// }

	// var registereduserId int64
	// err = registeredstmt.QueryRow(usr.Registered_by).Scan(&registereduserId)
	// if err == nil {
	// 	res.Status = 401
	// 	res.Message = "user id not found"
	// 	res.Data = "User ID: " + fmt.Sprint(userId)
	// 	return res, errors.New("user already registered")
	// } else if err != sql.ErrNoRows {
	// 	res.Status = 500
	// 	res.Message = "Query execution failed"
	// 	res.Data = err.Error()
	// 	return res, err
	// }
	// defer registeredstmt.Close()

	// masukkan ke db
	insertquery := "INSERT INTO user (username,password,nama_lengkap,email,nomor_telepon,tanggal_lahir) VALUES (?,?,?,?,?,NOW())"
	insertstmt, err := con.Prepare(insertquery)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer insertstmt.Close()

	result, err := insertstmt.Exec(usr.Username, usr.Password, usr.Nama_lengkap, usr.Email, usr.No_telp)
	if err != nil {
		res.Status = 401
		res.Message = "insert user gagal"
		res.Data = err.Error()
		return res, errors.New("insert user gagal")
	}
	defer stmt.Close()

	userId, err = result.LastInsertId()
	if err != nil {
		res.Status = 500
		res.Message = "gagal mendapatkan user ID"
		res.Data = err.Error()
		return res, err
	}
	usr.Id = int(userId)

	// tambah ke user detail
	insertdetailquery := "INSERT INTO user_detail (user_detail_id,user_kelas_id,status,tipe) VALUES (?,?,?,?)"
	insertdetailstmt, err := con.Prepare(insertdetailquery)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer insertstmt.Close()

	_, err = insertdetailstmt.Exec(usr.Id, 1, 1, 7)
	if err != nil {
		res.Status = 401
		res.Message = "insert user detail gagal"
		res.Data = err.Error()
		return res, errors.New("insert user detail gagal")
	}
	defer stmt.Close()

	// tambah ke user role dan user privilege
	insertrolequery := "INSERT INTO user_role (user_id,role_id) VALUES (?,?)"
	insertrolestmt, err := con.Prepare(insertrolequery)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer insertrolestmt.Close()

	_, err = insertrolestmt.Exec(usr.Id, 7)
	if err != nil {
		res.Status = 401
		res.Message = "insert user role gagal"
		res.Data = err.Error()
		return res, errors.New("insert user role gagal")
	}
	defer stmt.Close()

	insertprivquery := "INSERT INTO user_privilege (user_id,privilege_id) VALUES (?,?)"
	insertprivstmt, err := con.Prepare(insertprivquery)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer insertprivstmt.Close()

	_, err = insertprivstmt.Exec(usr.Id, 17)
	if err != nil {
		res.Status = 401
		res.Message = "insert user privilege gagal"
		res.Data = err.Error()
		return res, errors.New("insert user privilege gagal")
	}
	defer stmt.Close()

	// tambah ke surveyor
	insertsurvquery := "INSERT INTO surveyor (user_id,registered_by,lokasi) VALUES (?,?,?)"
	insertsurvstmt, err := con.Prepare(insertsurvquery)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer insertsurvstmt.Close()

	_, err = insertsurvstmt.Exec(usr.Id, usr.Registered_by, "")
	if err != nil {
		res.Status = 401
		res.Message = "insert surveyor gagal"
		res.Data = err.Error()
		return res, errors.New("insert surveyor gagal")
	}

	// set waktu login dan created_at login => update timestamp terakhir login
	updateQuery := "UPDATE user SET login_timestamp = NOW(), created_at = NOW() WHERE user_id = ?"
	updatestmt, err := con.Prepare(updateQuery)
	if err != nil {
		res.Status = 401
		res.Message = "stmt update gagal"
		res.Data = err.Error()
		return res, err
	}
	defer updatestmt.Close()

	_, err = updatestmt.Exec(usr.Id)
	if err != nil {
		res.Status = 401
		res.Message = "update login_timestamp gagal"
		res.Data = err.Error()
		return res, err
	}

	// hilangkan password buat global variabel
	usr.Password = ""
	res.Status = http.StatusOK
	res.Message = "Berhasil buat user"
	res.Data = usr

	defer db.DbClose(con)

	return res, nil
}

func GetAllSurveyor() (Response, error) {
	var res Response
	var dtUserSurveyor = []UserSurveyor{}

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	query := "SELECT u.user_id,u.username,u.password,u.nama_lengkap,u.alamat,u.jenis_kelamin,u.tanggal_lahir,u.email,u.nomor_telepon,u.foto_profil,u.ktp,s.suveyor_id,s.registered_by,s.lokasi,s.availability_surveyor FROM user u JOIN surveyor s ON u.user_id = s.user_id"
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		res.Status = 401
		res.Message = "exec gagal"
		res.Data = err.Error()
		return res, err
	}
	defer rows.Close()

	for rows.Next() {
		var _dtUserSurveyor UserSurveyor
		err := rows.Scan(&_dtUserSurveyor.User_id, &_dtUserSurveyor.Username, &_dtUserSurveyor.Password, &_dtUserSurveyor.Nama_lengkap, &_dtUserSurveyor.Alamat, &_dtUserSurveyor.Jenis_kelamin, &_dtUserSurveyor.Tgl_lahir, &_dtUserSurveyor.Email, &_dtUserSurveyor.No_telp, &_dtUserSurveyor.Foto_profil, &_dtUserSurveyor.Ktp, &_dtUserSurveyor.Surveyor_id, &_dtUserSurveyor.Registered_by, &_dtUserSurveyor.Lokasi, &_dtUserSurveyor.Availability_surveyor)
		if err != nil {
			res.Status = 401
			res.Message = "scan gagal"
			res.Data = err.Error()
			return res, err
		}

		dtUserSurveyor = append(dtUserSurveyor, _dtUserSurveyor)
	}

	if err = rows.Err(); err != nil {
		res.Status = 401
		res.Message = "rows error"
		res.Data = err.Error()
		return res, err
	}

	if len(dtUserSurveyor) == 0 {
		res.Status = 404
		res.Message = "Data tidak ditemukan"
		res.Data = nil
		return res, nil
	}

	res.Status = http.StatusOK
	res.Message = "Berhasil mengambil data"
	res.Data = dtUserSurveyor

	defer db.DbClose(con)
	return res, nil
}

func GetSurveyorById(surveyor_id string) (Response, error) {
	var res Response
	var dtSurveyor Surveyor

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	query := "SELECT * FROM surveyor WHERE surveyor_id = ?"
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()
	nId, _ := strconv.Atoi(surveyor_id)
	err = stmt.QueryRow(nId).Scan(&dtSurveyor.Surveyor_id, &dtSurveyor.User_id, &dtSurveyor.Lokasi, &dtSurveyor.Availability_surveyor)
	if err != nil {
		res.Status = 401
		res.Message = "exec gagal"
		res.Data = err.Error()
		return res, err
	}

	res.Status = http.StatusOK
	res.Message = "Berhasil mengambil data"
	res.Data = dtSurveyor

	defer db.DbClose(con)
	return res, nil
}

func GetSurveyorByName(nama_surveyor string) (Response, error) {
	var res Response
	var dtUserSurveyor = []UserSurveyor{}

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	query := "SELECT u.user_id,u.username,u.nama_lengkap,u.alamat,u.jenis_kelamin,u.tanggal_lahir,u.email,u.nomor_telepon,u.foto_profil,u.ktp,s.suveyor_id,s.lokasi,s.availability_surveyor,COUNT(CASE WHEN sr.status_request = 'O' THEN 1 END) AS surveyonprogress,COUNT(sr.id_transaksi_jual_sewa) AS totalsurvey FROM user u JOIN surveyor s ON u.user_id = s.user_id LEFT JOIN survey_request sr ON u.user_id = sr.user_id WHERE u.username LIKE ? GROUP BY u.user_id"
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	rows, err := stmt.Query("%" + nama_surveyor + "%")
	if err != nil {
		res.Status = 401
		res.Message = "exec gagal"
		res.Data = err.Error()
		return res, err
	}
	defer rows.Close()

	for rows.Next() {
		var _dtUserSurveyor UserSurveyor
		err := rows.Scan(&_dtUserSurveyor.User_id, &_dtUserSurveyor.Username, &_dtUserSurveyor.Nama_lengkap, &_dtUserSurveyor.Alamat, &_dtUserSurveyor.Jenis_kelamin, &_dtUserSurveyor.Tgl_lahir, &_dtUserSurveyor.Email, &_dtUserSurveyor.No_telp, &_dtUserSurveyor.Foto_profil, &_dtUserSurveyor.Ktp, &_dtUserSurveyor.Surveyor_id, &_dtUserSurveyor.Lokasi, &_dtUserSurveyor.Availability_surveyor, &_dtUserSurveyor.SurveyOnProgress, &_dtUserSurveyor.TotalSurvey)
		if err != nil {
			res.Status = 401
			res.Message = "scan gagal"
			res.Data = err.Error()
			return res, err
		}

		dtUserSurveyor = append(dtUserSurveyor, _dtUserSurveyor)
	}

	if err = rows.Err(); err != nil {
		res.Status = 401
		res.Message = "rows error"
		res.Data = err.Error()
		return res, err
	}

	if len(dtUserSurveyor) == 0 {
		res.Status = 404
		res.Message = "Data tidak ditemukan"
		res.Data = nil
		return res, nil
	}

	res.Status = http.StatusOK
	res.Message = "Berhasil mengambil data"
	res.Data = dtUserSurveyor

	defer db.DbClose(con)
	return res, nil
}

func UpdateSurveyorById(data_surveyor string) (Response, error) {
	var res Response

	var userSurveyor UserSurveyor

	err := json.Unmarshal([]byte(data_surveyor), &userSurveyor)
	if err != nil {
		res.Status = 401
		res.Message = "gagal decode json"
		res.Data = err.Error()
		return res, err
	}

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	query := "UPDATE user u JOIN surveyor s ON u.user_id = s.user_id SET u.username = ?, u.password = ?, u.email = ?, u.nomor_telepon = ?, u.updated_at = NOW() WHERE s.suveyor_id = ? "
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(userSurveyor.Username, userSurveyor.Password, userSurveyor.Email, userSurveyor.No_telp, userSurveyor.Surveyor_id)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}

	res.Status = http.StatusOK
	res.Message = "Berhasil mengupdate data"
	res.Data = result

	defer db.DbClose(con)
	return res, nil
}

func DeleteSurveyorById(inspektur string) (Response, error) {
	var res Response

	var dtSurveyor = Surveyor{}

	err := json.Unmarshal([]byte(inspektur), &dtSurveyor)
	if err != nil {
		res.Status = 401
		res.Message = "gagal decode json"
		res.Data = err.Error()
		return res, err
	}

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	query := "DELETE FROM surveyor WHERE surveyor_id = ?"
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(dtSurveyor.Surveyor_id)
	if err != nil {
		res.Status = 401
		res.Message = "exec gagal"
		res.Data = err.Error()
		return res, err
	}

	res.Status = http.StatusOK
	res.Message = "Berhasil menghapus data"
	res.Data = result

	defer db.DbClose(con)

	return res, nil
}

func GetAllSurveyorDetailed() (Response, error) {
	var res Response
	var dtUserSurveyor = []UserSurveyor{}

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	query := "SELECT u.user_id,u.username,u.password,u.nama_lengkap,u.alamat,u.jenis_kelamin,u.tanggal_lahir,u.email,u.nomor_telepon,u.foto_profil,u.ktp,s.suveyor_id,s.registered_by,s.lokasi,s.availability_surveyor,COUNT(CASE WHEN sr.status_request = 'O' THEN 1 END) AS surveyonprogress,COUNT(sr.id_transaksi_jual_sewa) AS totalsurvey FROM user u JOIN surveyor s ON u.user_id = s.user_id LEFT JOIN survey_request sr ON u.user_id = sr.user_id GROUP BY u.user_id"
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		res.Status = 401
		res.Message = "exec gagal"
		res.Data = err.Error()
		return res, err
	}
	defer rows.Close()
	for rows.Next() {
		var _dtUserSurveyor UserSurveyor
		err := rows.Scan(&_dtUserSurveyor.User_id, &_dtUserSurveyor.Username, &_dtUserSurveyor.Password, &_dtUserSurveyor.Nama_lengkap, &_dtUserSurveyor.Alamat, &_dtUserSurveyor.Jenis_kelamin, &_dtUserSurveyor.Tgl_lahir, &_dtUserSurveyor.Email, &_dtUserSurveyor.No_telp, &_dtUserSurveyor.Foto_profil, &_dtUserSurveyor.Ktp, &_dtUserSurveyor.Surveyor_id, &_dtUserSurveyor.Registered_by, &_dtUserSurveyor.Lokasi, &_dtUserSurveyor.Availability_surveyor, &_dtUserSurveyor.SurveyOnProgress, &_dtUserSurveyor.TotalSurvey)
		if err != nil {
			res.Status = 401
			res.Message = "scan gagal"
			res.Data = err.Error()
			return res, err
		}

		dtUserSurveyor = append(dtUserSurveyor, _dtUserSurveyor)
	}

	if err = rows.Err(); err != nil {
		res.Status = 401
		res.Message = "rows error"
		res.Data = err.Error()
		return res, err
	}

	if len(dtUserSurveyor) == 0 {
		res.Status = 404
		res.Message = "Data tidak ditemukan"
		res.Data = nil
		return res, nil
	}

	res.Status = http.StatusOK
	res.Message = "Berhasil mengambil data"
	res.Data = dtUserSurveyor

	defer db.DbClose(con)
	return res, nil
}

func GetSurveyorByUserId(user_id string) (Response, error) {
	var res Response
	var dtSurveyor UserSurveyor
	fmt.Println("get user by user id")
	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	query := `
		SELECT u.username,u.password,u.nama_lengkap,u.alamat,u.jenis_kelamin,u.tanggal_lahir,u.email,u.nomor_telepon,u.foto_profil,u.ktp,s.lokasi,COUNT(CASE WHEN sr.status_request = 'O' THEN 1 END), 
			COUNT(CASE WHEN sr.status_request = 'R' THEN 1 END), COUNT(CASE WHEN sr.status_request = 'F' THEN 1 END) AS requests,
			s.suveyor_id, s.registered_by
		FROM surveyor s 
		JOIN user u ON s.user_id = u.user_id 
		LEFT JOIN survey_request sr ON u.user_id = sr.user_id 
		WHERE s.user_id = ?`
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	var _ongoingReq1 int
	var _ongoingReq2 int
	var ongoingReq int
	var finishedReq int
	nId, _ := strconv.Atoi(user_id)
	dtSurveyor.User_id = nId
	err = stmt.QueryRow(nId).Scan(&dtSurveyor.Username, &dtSurveyor.Password, &dtSurveyor.Nama_lengkap, &dtSurveyor.Alamat, &dtSurveyor.Jenis_kelamin, &dtSurveyor.Tgl_lahir,
		&dtSurveyor.Email, &dtSurveyor.No_telp, &dtSurveyor.Foto_profil, &dtSurveyor.Ktp, &dtSurveyor.Lokasi, &_ongoingReq1, &_ongoingReq2, &finishedReq, &dtSurveyor.Surveyor_id, &dtSurveyor.Registered_by)
	if err != nil {
		res.Status = 401
		res.Message = "exec gagal"
		res.Data = err.Error()
		return res, err
	}

	ongoingReq = _ongoingReq1 + _ongoingReq2
	dtSurveyor.FinishedSurvey = ongoingReq

	// looping untuk ambil semua survey request
	querysurreq := "SELECT * FROM survey_request WHERE user_id = ?"
	stmtsurreq, err := con.Prepare(querysurreq)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmtsurreq.Close()

	rows, err := stmtsurreq.Query(nId)
	if err != nil {
		res.Status = 401
		res.Message = "exec gagal"
		res.Data = err.Error()
		return res, err
	}

	defer rows.Close()
	for rows.Next() {
		var _surveyReq SurveyRequest
		err := rows.Scan(&_surveyReq.Id_transaksi_jual_sewa, &_surveyReq.User_id, &_surveyReq.Id_asset, &_surveyReq.Created_at, &_surveyReq.Dateline, &_surveyReq.Status_request, &_surveyReq.Status_verifikasi, &_surveyReq.Data_lengkap, &_surveyReq.Usage_old, &_surveyReq.Usage_new, &_surveyReq.Luas_old, &_surveyReq.Luas_new, &_surveyReq.Nilai_old, &_surveyReq.Nilai_new, &_surveyReq.Kondisi_old, &_surveyReq.Kondisi_new, &_surveyReq.Batas_koordinat_old, &_surveyReq.Batas_koordinat_new, &_surveyReq.Tags_old, &_surveyReq.Tags_new)
		if err != nil {
			res.Status = 401
			res.Message = "scan gagal"
			res.Data = err.Error()
			return res, err
		}

		dtSurveyor.Survey_Request = append(dtSurveyor.Survey_Request, _surveyReq)
	}

	res.Status = http.StatusOK
	res.Message = "Berhasil mengambil data"
	res.Data = dtSurveyor

	defer db.DbClose(con)

	return res, nil
}

func UpdateUserBySurveyorId(input string) (Response, error) {
	var res Response

	type TempUpdateSurveyorId struct {
		UserId     int    `json:"user_id"`
		SurveyorId int    `json:"surveyor_id"`
		Nama       string `json:"nama"`
		Username   string `json:"username"`
		Password   string `json:"password"`
		Email      string `json:"email"`
		NoTelp     string `json:"notelp"`
	}

	var userSurveyor TempUpdateSurveyorId
	err := json.Unmarshal([]byte(input), &userSurveyor)
	if err != nil {
		res.Status = 401
		res.Message = "gagal decode json"
		res.Data = err.Error()
		return res, err
	}

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	queryUserId := "SELECT user_id FROM surveyor WHERE suveyor_id = ?"
	stmtUserId, err := con.Prepare(queryUserId)
	if err != nil {
		res.Status = 401
		res.Message = "gagal mendapatkan user_id"
		res.Data = err.Error()
		return res, err
	}
	defer stmtUserId.Close()

	err = stmtUserId.QueryRow(userSurveyor.SurveyorId).Scan(&userSurveyor.UserId)
	if err != nil {
		if err == sql.ErrNoRows {
			res.Status = 404
			res.Message = "Surveyor not found"
		} else {
			res.Status = 401
			res.Message = "Failed to execute statement"
		}
		res.Data = err.Error()
		return res, err
	}

	fmt.Println("userid: ", userSurveyor.UserId)
	query := "UPDATE user SET username = ?, password = ?, nama_lengkap = ?, email = ?, nomor_telepon = ?,updated_at = NOW() WHERE user_id = ? "
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	_, err = stmt.Exec(userSurveyor.Username, userSurveyor.Password, userSurveyor.Nama, userSurveyor.Email, userSurveyor.NoTelp, userSurveyor.UserId)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}

	var tempuser Response
	tempuser, _ = GetSurveyorByUserId(strconv.Itoa(userSurveyor.UserId))

	res.Status = http.StatusOK
	res.Message = "Berhasil mengupdate data"
	res.Data = tempuser.Data

	defer db.DbClose(con)
	return res, nil
}

func UpdateSurveyorByUserId(input string) (Response, error) {
	var res Response

	type TempUpdateSurveyorId struct {
		UserId   int    `json:"user_id"`
		Nama     string `json:"nama"`
		Username string `json:"username"`
		Password string `json:"password"`
		Email    string `json:"email"`
		NoTelp   string `json:"notelp"`
	}

	var userSurveyor TempUpdateSurveyorId
	err := json.Unmarshal([]byte(input), &userSurveyor)
	if err != nil {
		res.Status = 401
		res.Message = "gagal decode json"
		res.Data = err.Error()
		return res, err
	}

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	query := "UPDATE user SET username = ?, password = ?, nama_lengkap = ?, email = ?, nomor_telepon = ?,updated_at = NOW() WHERE user_id = ? "
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	_, err = stmt.Exec(userSurveyor.Username, userSurveyor.Password, userSurveyor.Nama, userSurveyor.Email, userSurveyor.NoTelp, userSurveyor.UserId)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}

	var tempuser Response
	tempuser, _ = GetSurveyorByUserId(strconv.Itoa(userSurveyor.UserId))

	res.Status = http.StatusOK
	res.Message = "Berhasil mengupdate data"
	res.Data = tempuser.Data

	defer db.DbClose(con)
	return res, nil
}
