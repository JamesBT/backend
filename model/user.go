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

func Login(akun string) (Response, error) {
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
	query = "SELECT u.user_id,u.username,u.nama_lengkap,u.alamat,u.jenis_kelamin,u.tanggal_lahir,u.email,u.nomor_telepon,u.foto_profil,u.ktp,ud.status,ud.tipe,ud.first_login,ud.denied_by_admin FROM user u JOIN user_detail ud WHERE u.username = ? AND u.password = ?"
	stmt, err = con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}

	err = stmt.QueryRow(usr.Username, usr.Password).Scan(&loginUsr.Id, &loginUsr.Username, &loginUsr.Nama_lengkap, &loginUsr.Alamat, &loginUsr.Jenis_kelamin, &loginUsr.Tgl_lahir, &loginUsr.Email, &loginUsr.No_telp, &loginUsr.Foto_profil, &loginUsr.Ktp, &loginUsr.Status, &loginUsr.Tipe, &loginUsr.First_login, &loginUsr.Denied_by_admin)
	if err != nil {
		res.Status = 401
		res.Message = "password salah"
		res.Data = err.Error()
		return res, errors.New("password salah")
	}
	defer stmt.Close()

	// ambil role + privilege
	getRoleQuery := "SELECT ur.role_id, r.nama_role FROM  user_role ur JOIN role r ON ur.role_id = r.role_id WHERE ur.user_id = ?;"
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

func SignUp(akun string) (Response, error) {
	var res Response

	var usr = User{}

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
	query := "SELECT user_id FROM user WHERE username = ?"
	// query := "INSERT INTO user (username,password,nama_lengkap,email,nomor_telepon) VALUES (?,?,?,?,?)"
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}

	var userId int64
	err = stmt.QueryRow(usr.Username).Scan(&userId)
	if err == nil {
		res.Status = 401
		res.Message = "User already registered"
		res.Data = "User ID: " + fmt.Sprint(userId)
		return res, errors.New("user already registered")
	} else if err != sql.ErrNoRows {
		res.Status = 500
		res.Message = "Query execution failed"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	// masukkan ke db
	insertquery := "INSERT INTO user (username,password,nama_lengkap,email,nomor_telepon) VALUES (?,?,?,?,?)"
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

func GetUserById(id_user string) (Response, error) {
	var res Response

	var usr User

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka koneksi"
		res.Data = err.Error()
		return res, err
	}

	query := "SELECT user_id, username, nama_lengkap, alamat, jenis_kelamin, tanggal_lahir, email, nomor_telepon, foto_profil, ktp FROM user WHERE user_id = ?"
	stmt, err := con.Prepare(query)

	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	nId, _ := strconv.Atoi(id_user)
	err = stmt.QueryRow(nId).Scan(&usr.Id, &usr.Username, &usr.Nama_lengkap, &usr.Alamat, &usr.Jenis_kelamin, &usr.Tgl_lahir, &usr.Email, &usr.No_telp, &usr.Foto_profil, &usr.Ktp)
	if err != nil {
		res.Status = 401
		res.Message = "rows gagal"
		res.Data = err.Error()
		return res, err
	}

	res.Status = http.StatusOK
	res.Message = "Berhasil mengambil data"
	res.Data = usr

	defer db.DbClose(con)

	return res, nil
}

func ForgotPass(email string) (Response, error) {
	var res Response

	// asd

	return res, nil
}

func ChangePass() (Response, error) {
	var res Response

	return res, nil
}
