package model

import (
	"TemplateProject/db"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"mime/multipart"
	"net/http"
	"net/mail"
	"net/smtp"
	"os"
	"strconv"
	"strings"
	"time"
)

// CRUD user ============================================================================
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
	// queryinsert := "SELECT user_id, username, nama_lengkap, alamat, jenis_kelamin, tanggal_lahir, email, nomor_telepon, foto_profil, ktp FROM user WHERE username = ? AND password = ?"
	queryinsert := "SELECT u.user_id, u.username, u.nama_lengkap, u.alamat, u.jenis_kelamin, u.tanggal_lahir, u.email, u.nomor_telepon, u.foto_profil, u.ktp, ud.user_kelas_id, ud.status, ud.tipe, ud.first_login, ud.denied_by_admin FROM user u JOIN user_detail ud ON u.user_id = ud.user_detail_id WHERE u.username = ? AND u.password = ?;"
	stmtinsert, err := con.Prepare(queryinsert)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmtinsert.Close()

	err = stmtinsert.QueryRow(usr.Username, usr.Password).Scan(&loginUsr.Id, &loginUsr.Username, &loginUsr.Nama_lengkap, &loginUsr.Alamat, &loginUsr.Jenis_kelamin, &loginUsr.Tgl_lahir, &loginUsr.Email, &loginUsr.No_telp, &loginUsr.Foto_profil, &loginUsr.Ktp, &loginUsr.Kelas, &loginUsr.Status, &loginUsr.Tipe, &loginUsr.First_login, &loginUsr.Denied_by_admin)
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

	// ngecek email
	_, err = mail.ParseAddress(usr.Email)
	if err != nil {
		res.Status = 500
		res.Message = "Invalid email"
		res.Data = err.Error()
		return res, err
	} else {
		fmt.Println("email valid")
	}

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
	fmt.Println(usr)
	result, err := insertstmt.Exec(usr.Username, usr.Password, usr.Nama_lengkap, usr.Email, usr.No_telp)
	if err != nil {
		res.Status = 401
		res.Message = "insert user gagal"
		res.Data = err.Error()
		return res, err
	}
	defer insertstmt.Close()

	userId, err = result.LastInsertId()
	if err != nil {
		res.Status = 500
		res.Message = "gagal mendapatkan user ID"
		res.Data = err.Error()
		return res, err
	}
	usr.Id = int(userId)

	// random number generator untuk buat kode otp 4 digit 1000-9999
	randomizer := rand.New(rand.NewSource(time.Now().UTC().UnixNano()))

	randomnumber := randomizer.Intn(9000) + 1000

	// tambah ke user detail
	insertdetailquery := "INSERT INTO user_detail (user_detail_id,user_kelas_id,status,tipe,kode_otp) VALUES (?,?,?,?,?)"
	insertdetailstmt, err := con.Prepare(insertdetailquery)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer insertstmt.Close()

	_, err = insertdetailstmt.Exec(usr.Id, 1, 1, 8, randomnumber)
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

	_, err = insertrolestmt.Exec(usr.Id, 8)
	if err != nil {
		res.Status = 401
		res.Message = "insert user role gagal"
		res.Data = err.Error()
		return res, errors.New("insert user detail gagal")
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
		return res, errors.New("insert user detail gagal")
	}
	defer stmt.Close()

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
	// kirim email untuk kode otp
	to := []string{usr.Email}
	cc := []string{usr.Email}
	subject := "Aset Manajemen: Kode Verifikasi (OTP) untuk Verifikasi Identitas"
	message := "Hai " + usr.Username + "\n\nKode verifikasi (OTP) Aset Manajemen kamu:\n " + strconv.Itoa(randomnumber)
	err = sendMail(to, cc, subject, message)
	if err != nil {
		res.Status = 401
		res.Message = "gagal kirim email verifikasi kode otp"
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

func GetAllUser() (Response, error) {
	var res Response
	var arrUser = []User{}

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	query := `
	SELECT u.user_id, u.username, u.password, u.nama_lengkap, u.alamat, u.jenis_kelamin, 
		u.tanggal_lahir, u.email, u.nomor_telepon, u.foto_profil, u.ktp, ud.user_kelas_id, 
		ud.status, ud.tipe, ud.first_login, ud.denied_by_admin, 
		ur.role_id, r.nama_role, up.privilege_id, p.nama_privilege
	FROM user u 
	INNER JOIN user_detail ud ON u.user_id = ud.user_detail_id
	JOIN user_role ur ON u.user_id = ur.user_id
	JOIN user_privilege up ON u.user_id = up.user_id
	JOIN role r ON ur.role_id = r.role_id
	JOIN privilege p ON up.privilege_id = p.privilege_id
	`
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	result, err := stmt.Query()
	if err != nil {
		res.Status = 401
		res.Message = "exec gagal"
		res.Data = err.Error()
		return res, err
	}
	defer result.Close()

	userMap := make(map[int]*User)

	for result.Next() {
		var dtUser User
		var roleId int
		var roleName string
		var privilegeId int
		var privilegeName string
		err = result.Scan(&dtUser.Id, &dtUser.Username, &dtUser.Password, &dtUser.Nama_lengkap, &dtUser.Alamat, &dtUser.Jenis_kelamin, &dtUser.Tgl_lahir, &dtUser.Email, &dtUser.No_telp, &dtUser.Foto_profil, &dtUser.Ktp, &dtUser.Kelas, &dtUser.Status, &dtUser.Tipe, &dtUser.First_login, &dtUser.Denied_by_admin, &roleId, &roleName, &privilegeId, &privilegeName)
		if err != nil {
			res.Status = 401
			res.Message = "rows scan"
			res.Data = err.Error()
			return res, err
		}

		if existingUser, ok := userMap[dtUser.Id]; ok {
			// Check if the role already exists in the user's role list
			roleExists := false
			for _, r := range existingUser.UserRole {
				if r.Role_id == roleId {
					roleExists = true
					break
				}
			}
			if !roleExists {
				existingUser.UserRole = append(existingUser.UserRole, Role{Role_id: roleId, Nama_role: roleName})
			}

			// Check if the privilege already exists in the user's privilege list
			privilegeExists := false
			for _, p := range existingUser.UserPrivilege {
				if p.Privilege_id == privilegeId {
					privilegeExists = true
					break
				}
			}
			if !privilegeExists {
				existingUser.UserPrivilege = append(existingUser.UserPrivilege, Privilege{Privilege_id: privilegeId, Nama_privilege: privilegeName})
			}
		} else {
			// If the user is not already in the map, add them with their initial role and privilege
			dtUser.UserRole = []Role{{Role_id: roleId, Nama_role: roleName}}
			dtUser.UserPrivilege = []Privilege{{Privilege_id: privilegeId, Nama_privilege: privilegeName}}
			userMap[dtUser.Id] = &dtUser
		}

	}

	for _, user := range userMap {
		arrUser = append(arrUser, *user)
	}

	res.Status = http.StatusOK
	res.Message = "Berhasil mengambil data"
	res.Data = arrUser

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

func GetUserByUsername(username string) (Response, error) {
	var res Response
	var dtUsers = []User{}

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	query := "SELECT user_id, username, nama_lengkap,alamat,jenis_kelamin,tanggal_lahir,email,nomor_telepon,foto_profil,ktp FROM user WHERE username LIKE ?"
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	rows, err := stmt.Query("%" + username + "%")
	if err != nil {
		res.Status = 401
		res.Message = "exec gagal"
		res.Data = err.Error()
		return res, err
	}
	defer rows.Close()
	for rows.Next() {
		var dtUser User
		err := rows.Scan(&dtUser.Id, &dtUser.Username, &dtUser.Nama_lengkap, &dtUser.Alamat, &dtUser.Jenis_kelamin, &dtUser.Tgl_lahir, &dtUser.Email, &dtUser.No_telp, &dtUser.Foto_profil, &dtUser.Ktp)
		if err != nil {
			res.Status = 401
			res.Message = "scan gagal"
			res.Data = err.Error()
			return res, err
		}
		dtUsers = append(dtUsers, dtUser)
	}

	if err = rows.Err(); err != nil {
		res.Status = 401
		res.Message = "rows error"
		res.Data = err.Error()
		return res, err
	}

	if len(dtUsers) == 0 {
		res.Status = 404
		res.Message = "Data tidak ditemukan"
		res.Data = nil
		return res, nil
	}

	res.Status = http.StatusOK
	res.Message = "Berhasil mengambil data"
	res.Data = dtUsers

	defer db.DbClose(con)
	return res, nil
}

func UpdateUser(filefoto *multipart.FileHeader, userid, username, nama_lengkap, alamat, jenis_kelamin, tanggal_lahir, email, no_telp string) (Response, error) {
	var res Response

	userId, _ := strconv.Atoi(userid)

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	query := "UPDATE user SET username = ?, nama_lengkap = ?, alamat = ?, jenis_kelamin = ?, tanggal_lahir = ?, email = ?, nomor_telepon = ?,updated_at = NOW() WHERE user_id = ? "
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(username, nama_lengkap, alamat, jenis_kelamin, tanggal_lahir, email, no_telp, userId)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}

	// tambah file foto profile dan ktp
	// foto profil ======================================================
	fmt.Println(filefoto.Header.Get("Content-type"))
	// tipe := filefoto.Header.Get("Content-type")

	tipeGambar := ".png"
	// if tipe == "image/png" {
	// 	tipeGambar = ".png"
	// } else if tipe == "image/jpg" {
	// 	tipeGambar = ".jpg"
	// } else if tipe == "image/jpeg" {
	// 	tipeGambar = ".jpg"
	// }

	filefoto.Filename = userid + tipeGambar
	pathFotoFile := "uploads/user/foto_profil/" + filefoto.Filename
	//source
	srcfoto, err := filefoto.Open()
	if err != nil {
		log.Println(err.Error())
		return res, err
	}
	defer srcfoto.Close()

	// Destination
	dstfoto, err := os.Create("uploads/user/foto_profil/" + filefoto.Filename)
	if err != nil {
		log.Println("2")
		fmt.Print("2")
		return res, err
	}

	// Copy
	if _, err = io.Copy(dstfoto, srcfoto); err != nil {
		log.Println(err.Error())
		log.Println("3")
		fmt.Print("3")
		return res, err
	}
	dstfoto.Close()

	err = UpdateDataFotoPath("user", "foto_profil", pathFotoFile, "user_id", userId)
	if err != nil {
		return res, err
	}

	res.Status = http.StatusOK
	res.Message = "Berhasil mengupdate data"
	res.Data = result

	defer db.DbClose(con)
	return res, nil
}

func UpdateUserFull(filefoto *multipart.FileHeader, filektp *multipart.FileHeader, userid, username, nama_lengkap, alamat, jenis_kelamin, tanggal_lahir, email, no_telp string) (Response, error) {
	var res Response

	userId, _ := strconv.Atoi(userid)

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}
	var usrStatus int
	statusQuery := "SELECT status FROM user_detail WHERE user_detail_id = ?"
	statusstmt, err := con.Prepare(statusQuery)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer statusstmt.Close()
	err = statusstmt.QueryRow(userid).Scan(&usrStatus)
	if err != nil {
		res.Status = 401
		res.Message = "exec gagal"
		res.Data = err.Error()
		return res, err
	}
	if usrStatus == 0 {
		res.Status = 403
		res.Message = "Akses ditolak: Pengguna tidak diizinkan untuk memperbarui data"
		res.Data = nil
		return res, errors.New("pengguna tidak diizinkan untuk memperbarui data")
	}

	query := "UPDATE user SET username = ?, nama_lengkap = ?, alamat = ?, jenis_kelamin = ?, tanggal_lahir = ?, email = ?, nomor_telepon = ?,updated_at = NOW() WHERE user_id = ? "
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(username, nama_lengkap, alamat, jenis_kelamin, tanggal_lahir, email, no_telp, userId)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}

	// tambah file foto profile dan ktp
	// foto profil ======================================================
	fmt.Println(filefoto.Header.Get("Content-type"))
	// tipe := filefoto.Header.Get("Content-type")

	tipeGambar := ".png"
	// if tipe == "image/png" {
	// 	tipeGambar = ".png"
	// } else if tipe == "image/jpg" {
	// 	tipeGambar = ".jpg"
	// } else if tipe == "image/jpeg" {
	// 	tipeGambar = ".jpg"
	// }

	filefoto.Filename = userid + tipeGambar
	pathFotoFile := "uploads/user/foto_profil/" + filefoto.Filename
	//source
	srcfoto, err := filefoto.Open()
	if err != nil {
		log.Println(err.Error())
		return res, err
	}
	defer srcfoto.Close()

	// Destination
	dstfoto, err := os.Create("uploads/user/foto_profil/" + filefoto.Filename)
	if err != nil {
		log.Println("2")
		fmt.Print("2")
		return res, err
	}

	// Copy
	if _, err = io.Copy(dstfoto, srcfoto); err != nil {
		log.Println(err.Error())
		log.Println("3")
		fmt.Print("3")
		return res, err
	}
	dstfoto.Close()

	err = UpdateDataFotoPath("user", "foto_profil", pathFotoFile, "user_id", userId)
	if err != nil {
		return res, err
	}

	// ktp ======================================================

	filektp.Filename = userid + tipeGambar
	pathKtpFile := "uploads/user/ktp/" + filefoto.Filename
	//source
	srcktp, err := filektp.Open()
	if err != nil {
		log.Println(err.Error())
		return res, err
	}
	defer srcktp.Close()

	// Destination
	dstktp, err := os.Create("uploads/user/ktp/" + filefoto.Filename)
	if err != nil {
		log.Println("2")
		fmt.Print("2")
		return res, err
	}

	// Copy
	if _, err = io.Copy(dstktp, srcktp); err != nil {
		log.Println(err.Error())
		log.Println("3")
		fmt.Print("3")
		return res, err
	}
	dstktp.Close()

	err = UpdateDataFotoPath("user", "ktp", pathKtpFile, "user_id", userId)
	if err != nil {
		return res, err
	}

	res.Status = http.StatusOK
	res.Message = "Berhasil mengupdate data"
	res.Data = result

	defer db.DbClose(con)
	return res, nil
}

func GetAllUserUnverified() (Response, error) {
	var res Response
	var arrUser = []User{}

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	query := "SELECT u.user_id, u.username, u.password, u.nama_lengkap, u.alamat, u.jenis_kelamin, u.tanggal_lahir, u.email, u.nomor_telepon, u.foto_profil, u.ktp FROM user u INNER JOIN user_detail ud ON u.user_id = ud.user_detail_id WHERE ud.status = 'N'"
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	result, err := stmt.Query()
	if err != nil {
		res.Status = 401
		res.Message = "exec gagal"
		res.Data = err.Error()
		return res, err
	}
	defer result.Close()
	fmt.Println("")
	for result.Next() {
		var dtUser User
		err = result.Scan(&dtUser.Id, &dtUser.Username, &dtUser.Password, &dtUser.Nama_lengkap, &dtUser.Alamat, &dtUser.Jenis_kelamin, &dtUser.Tgl_lahir, &dtUser.Email, &dtUser.No_telp, &dtUser.Foto_profil, &dtUser.Ktp)
		if err != nil {
			res.Status = 401
			res.Message = "rows scan"
			res.Data = err.Error()
			return res, err
		}
		fmt.Println("data user:", dtUser)
		arrUser = append(arrUser, dtUser)
	}

	res.Status = http.StatusOK
	res.Message = "Berhasil mengambil data"
	res.Data = arrUser

	defer db.DbClose(con)
	return res, nil
}

func GetUserKTP(id_user string) (Response, error) {
	var res Response

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	query := "SELECT ktp FROM user WHERE user_id = ?"
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	var ktpPath string
	err = stmt.QueryRow(id_user).Scan(&ktpPath)
	if err != nil {
		res.Status = 404
		res.Message = "KTP not found"
		res.Data = err.Error()
		return res, err
	}

	res.Status = http.StatusOK
	res.Message = "Berhasil mengambil data ktp"
	res.Data = ktpPath

	defer db.DbClose(con)

	return res, nil
}

func GetUserFoto(id_user string) (Response, error) {
	var res Response

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	query := "SELECT foto_profil FROM user WHERE user_id = ?"
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	var fotoPath string
	err = stmt.QueryRow(id_user).Scan(&fotoPath)
	if err != nil {
		res.Status = 404
		res.Message = "Foto profil not found"
		res.Data = err.Error()
		return res, err
	}

	res.Status = http.StatusOK
	res.Message = "Berhasil mengambil data ktp"
	res.Data = fotoPath

	defer db.DbClose(con)

	return res, nil
}

func DeleteUserById(id_user string) (Response, error) {
	var res Response

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	query := "UPDATE user SET deleted_at = NOW() WHERE user_id = ?"
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(id_user)
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

func GetAllUserByPerusahaanId(id_perusahaan string) (Response, error) {
	var res Response

	var arrUser = []User{}
	var dtUser User

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	query := `
		SELECT u.user_id, u.username, u.nama_lengkap, u.alamat, u.foto_profil 
		FROM user_perusahaan up
		JOIN user u ON up.id_user = u.user_id
		WHERE up.id_perusahaan = ?
	`
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	nId, _ := strconv.Atoi(id_perusahaan)
	result, err := stmt.Query(nId)
	if err != nil {
		res.Status = 401
		res.Message = "exec gagal"
		res.Data = err.Error()
		return res, err
	}
	defer result.Close()
	for result.Next() {
		err = result.Scan(&dtUser.Id, &dtUser.Username, &dtUser.Nama_lengkap, &dtUser.Alamat, &dtUser.Foto_profil)
		if err != nil {
			res.Status = 401
			res.Message = "rows scan"
			res.Data = err.Error()
			return res, err
		}
		arrUser = append(arrUser, dtUser)
	}

	res.Status = http.StatusOK
	res.Message = "Berhasil mengambil data"
	res.Data = arrUser

	defer db.DbClose(con)
	return res, nil
}

// CRUD user_privilege ============================================================================
func CreateUserPriv(userPriv string) (Response, error) {
	var res Response
	var dtUserPriv = UserPrivilege{}

	err := json.Unmarshal([]byte(userPriv), &dtUserPriv)
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

	query := "INSERT INTO user_privilege (privilege_id, user_id) VALUES (?,?)"
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(dtUserPriv.Privilege_id, dtUserPriv.User_id)
	if err != nil {
		res.Status = 401
		res.Message = "exec gagal"
		res.Data = err.Error()
		return res, err
	}
	lastId, err := result.LastInsertId()
	if err != nil {
		res.Status = 401
		res.Message = "Last Id gagal"
		res.Data = err.Error()
		return res, err
	}
	dtUserPriv.User_privilege_id = int(lastId)

	res.Status = http.StatusOK
	res.Message = "Berhasil memasukkan data"
	res.Data = dtUserPriv

	defer db.DbClose(con)
	return res, nil
}

func GetAllUserPriv() (Response, error) {
	var res Response
	var arrUserPriv = []UserPrivilege{}
	var dtUserPriv UserPrivilege

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	query := "SELECT * FROM user_privilege"
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	result, err := stmt.Query()
	if err != nil {
		res.Status = 401
		res.Message = "exec gagal"
		res.Data = err.Error()
		return res, err
	}
	defer result.Close()
	for result.Next() {
		err = result.Scan(&dtUserPriv.User_privilege_id, &dtUserPriv.Privilege_id, &dtUserPriv.User_id)
		if err != nil {
			res.Status = 401
			res.Message = "rows scan"
			res.Data = err.Error()
			return res, err
		}
		arrUserPriv = append(arrUserPriv, dtUserPriv)
	}

	res.Status = http.StatusOK
	res.Message = "Berhasil mengambil data"
	res.Data = arrUserPriv

	defer db.DbClose(con)
	return res, nil
}

func GetUserPrivById(user_priv_id string) (Response, error) {
	var res Response
	var dtUserPriv UserPrivilege

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	query := "SELECT * FROM user_privilege WHERE user_privilege_id = ?"
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()
	nId, _ := strconv.Atoi(user_priv_id)
	err = stmt.QueryRow(nId).Scan(&dtUserPriv.User_privilege_id, &dtUserPriv.User_id)
	if err != nil {
		res.Status = 401
		res.Message = "exec gagal"
		res.Data = err.Error()
		return res, err
	}

	res.Status = http.StatusOK
	res.Message = "Berhasil mengambil data"
	res.Data = dtUserPriv

	defer db.DbClose(con)
	return res, nil
}

func GetUserPrivByUserId(user_id string) (Response, error) {
	var res Response
	var dtUserPriv UserPrivilege

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	query := "SELECT * FROM user_privilege WHERE user_id = ?"
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()
	nId, _ := strconv.Atoi(user_id)
	err = stmt.QueryRow(nId).Scan(&dtUserPriv.User_privilege_id, &dtUserPriv.User_id)
	if err != nil {
		res.Status = 401
		res.Message = "exec gagal"
		res.Data = err.Error()
		return res, err
	}

	res.Status = http.StatusOK
	res.Message = "Berhasil mengambil data"
	res.Data = dtUserPriv

	defer db.DbClose(con)
	return res, nil
}

func GetUserPrivDetailByUserId(user_id string) (Response, error) {
	var res Response
	var privileges []map[string]interface{}

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	query := "SELECT up.user_id, up.privilege_id, p.nama_privilege FROM user_privilege up JOIN privilege p ON up.privilege_id = p.privilege_id WHERE up.user_id = ?"

	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	nId, _ := strconv.Atoi(user_id)
	rows, err := stmt.Query(nId)
	// err = stmt.QueryRow(nId).Scan(&temp_user_id, &temp_privilege_id, &temp_nama_privilege)
	if err != nil {
		res.Status = 401
		res.Message = "exec gagal"
		res.Data = err.Error()
		return res, err
	}
	defer rows.Close()

	for rows.Next() {
		var temp_privilege_id, temp_user_id int
		var temp_nama_privilege string

		err := rows.Scan(&temp_user_id, &temp_privilege_id, &temp_nama_privilege)
		if err != nil {
			res.Status = 401
			res.Message = "scan gagal"
			res.Data = err.Error()
			return res, err
		}

		privilege := map[string]interface{}{
			"user_id":        temp_user_id,
			"privilege_id":   temp_privilege_id,
			"nama_privilege": temp_nama_privilege,
		}
		privileges = append(privileges, privilege)
	}

	if len(privileges) == 0 {
		res.Status = 404
		res.Message = "Data tidak ditemukan"
		res.Data = nil
		return res, nil
	}

	res.Status = http.StatusOK
	res.Message = "Berhasil mengambil data"
	res.Data = privileges

	defer db.DbClose(con)
	return res, nil
}

func UpdateUserPriv(userPriv string) (Response, error) {
	var res Response

	var dtUserPriv = UserPrivilege{}

	err := json.Unmarshal([]byte(userPriv), &dtUserPriv)
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

	query := "UPDATE user_privilege SET privilege_id = ?, user_id = ? WHERE user_privilege_id = ?"
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(dtUserPriv.Privilege_id, dtUserPriv.User_id, dtUserPriv.User_privilege_id)
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

func DeleteUserPriv(userPriv string) (Response, error) {
	var res Response

	var dtUserPriv = UserPrivilege{}

	err := json.Unmarshal([]byte(userPriv), &dtUserPriv)
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

	query := "DELETE FROM user_privilege WHERE user_privilege_id = ?"
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(dtUserPriv.User_privilege_id)
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

// CRUD user_role

func CreateUserRole(userRole string) (Response, error) {
	var res Response
	var dtUserRole = UserRole{}

	err := json.Unmarshal([]byte(userRole), &dtUserRole)
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

	query := "INSERT INTO user_role (user_id, role_id) VALUES (?,?)"
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(dtUserRole.User_id, dtUserRole.Role_id)
	if err != nil {
		res.Status = 401
		res.Message = "exec gagal"
		res.Data = err.Error()
		return res, err
	}
	lastId, err := result.LastInsertId()
	if err != nil {
		res.Status = 401
		res.Message = "Last Id gagal"
		res.Data = err.Error()
		return res, err
	}
	dtUserRole.User_role_id = int(lastId)

	res.Status = http.StatusOK
	res.Message = "Berhasil memasukkan data"
	res.Data = dtUserRole

	defer db.DbClose(con)
	return res, nil
}

func GetAllUserRole() (Response, error) {
	var res Response
	var arrUserRole = []UserRole{}
	var dtUserRole UserRole

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	query := "SELECT * FROM user_role"
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	result, err := stmt.Query()
	if err != nil {
		res.Status = 401
		res.Message = "exec gagal"
		res.Data = err.Error()
		return res, err
	}
	defer result.Close()
	for result.Next() {
		err = result.Scan(&dtUserRole.User_role_id, &dtUserRole.User_id, &dtUserRole.Role_id)
		if err != nil {
			res.Status = 401
			res.Message = "rows scan"
			res.Data = err.Error()
			return res, err
		}
		arrUserRole = append(arrUserRole, dtUserRole)
	}

	res.Status = http.StatusOK
	res.Message = "Berhasil mengambil data"
	res.Data = arrUserRole

	defer db.DbClose(con)
	return res, nil
}

func GetUserRoleById(user_role_id string) (Response, error) {
	var res Response
	var dtUserRole UserRole

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	query := "SELECT * FROM user_role WHERE user_role_id = ?"
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()
	nId, _ := strconv.Atoi(user_role_id)
	err = stmt.QueryRow(nId).Scan(&dtUserRole.User_role_id, &dtUserRole.User_id, &dtUserRole.Role_id)
	if err != nil {
		res.Status = 401
		res.Message = "exec gagal"
		res.Data = err.Error()
		return res, err
	}

	res.Status = http.StatusOK
	res.Message = "Berhasil mengambil data"
	res.Data = dtUserRole

	defer db.DbClose(con)
	return res, nil
}

func GetUserRoleByUserId(user_id string) (Response, error) {
	var res Response
	var dtUserRole UserRole

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	query := "SELECT * FROM user_role WHERE user_id = ?"
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()
	nId, _ := strconv.Atoi(user_id)
	err = stmt.QueryRow(nId).Scan(&dtUserRole.User_role_id, &dtUserRole.User_id, &dtUserRole.Role_id)
	if err != nil {
		res.Status = 401
		res.Message = "exec gagal"
		res.Data = err.Error()
		return res, err
	}

	res.Status = http.StatusOK
	res.Message = "Berhasil mengambil data"
	res.Data = dtUserRole

	defer db.DbClose(con)
	return res, nil
}

func GetUserRoleDetailByUserId(user_id string) (Response, error) {
	var res Response
	var roles []map[string]interface{}

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	query := "SELECT ur.user_id, ur.role_id, r.nama_role FROM user_role ur JOIN role r ON ur.role_id = r.role_id WHERE ur.user_id = ?"
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	nId, _ := strconv.Atoi(user_id)
	rows, err := stmt.Query(nId)
	// err = stmt.QueryRow(nId).Scan(&temp_user_id, &temp_role_id, &temp_nama_role)
	if err != nil {
		res.Status = 401
		res.Message = "exec gagal"
		res.Data = err.Error()
		return res, err
	}
	defer rows.Close()

	for rows.Next() {
		var temp_user_id, temp_role_id int
		var temp_nama_role string
		err := rows.Scan(&temp_user_id, &temp_role_id, &temp_nama_role)
		if err != nil {
			res.Status = 401
			res.Message = "scan gagal"
			res.Data = err.Error()
			return res, err
		}

		role := map[string]interface{}{
			"user_id":   temp_user_id,
			"role_id":   temp_role_id,
			"nama_role": temp_nama_role,
		}
		roles = append(roles, role)
	}

	if len(roles) == 0 {
		res.Status = 404
		res.Message = "Data tidak ditemukan"
		res.Data = nil
		return res, nil
	}

	res.Status = http.StatusOK
	res.Message = "Berhasil mengambil data"
	res.Data = roles

	defer db.DbClose(con)
	return res, nil
}

func UpdateUserRole(userRole string) (Response, error) {
	var res Response

	var dtUserRole = UserRole{}

	err := json.Unmarshal([]byte(userRole), &dtUserRole)
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

	query := "UPDATE user_role SET user_id = ?, role_id = ? WHERE user_role_id = ?"
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(dtUserRole.User_id, dtUserRole.Role_id, dtUserRole.User_role_id)
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

func DeleteUserRole(userRole string) (Response, error) {
	var res Response

	var dtUserRole = UserRole{}

	err := json.Unmarshal([]byte(userRole), &dtUserRole)
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

	query := "DELETE FROM user_role WHERE user_role_id = ?"
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(dtUserRole.User_role_id)
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

// fungsi tambahan
func UploadFile(file *multipart.FileHeader, id string, kolom_id string, folder string) (Response, error) {
	var res Response

	log.Println("Upload File")
	nId, _ := strconv.Atoi(id)
	// file.Filename =
	pathFile := "uploads/user/" + file.Filename
	//source
	src, err := file.Open()
	if err != nil {
		log.Println(err.Error())
		log.Println("1")
		fmt.Print("1")
		return res, err
	}
	defer src.Close()

	// Destination
	dst, err := os.Create("uploads/" + folder + "/" + file.Filename)
	if err != nil {
		log.Println("2")
		fmt.Print("2")
		return res, err
	}

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		log.Println(err.Error())
		log.Println("3")
		fmt.Print("3")
		return res, err
	}
	dst.Close()

	err = UpdateDataFotoPath(folder, "foto", pathFile, kolom_id, nId)
	if err != nil {
		return res, err
	}

	res.Status = http.StatusOK
	res.Message = "Sukses Upload File"
	res.Data = file.Filename

	return res, nil
}

func UpdateDataFotoPath(tabel string, kolom string, path string, kolom_id string, id int) error {
	log.Println("mengubah status foto di DB")
	fmt.Println("mengubah status foto di DB")
	// Open DB connection
	con, err := db.DbConnection()
	if err != nil {
		log.Println("error: " + err.Error())
		return err
	}
	defer db.DbClose(con) // Ensure the connection is closed

	// Build the SQL query
	query := fmt.Sprintf("UPDATE %s SET %s='%s' WHERE %s = %d", tabel, kolom, path, kolom_id, id)
	fmt.Println(query)
	// Execute the query
	_, err = con.Exec(query) // Use Exec instead of Query since this is an UPDATE operation
	if err != nil {
		log.Println("error executing query: " + err.Error())
		return err
	}

	fmt.Println("status foto di edit")
	return nil
}

func VerifyOTP(input string) (Response, error) {
	var res Response
	type temp_verif_user_acc struct {
		UserID   int `json:"userid"`
		Kode_OTP int `json:"kode_otp"`
	}
	var requestacc temp_verif_user_acc
	err := json.Unmarshal([]byte(input), &requestacc)
	if err != nil {
		res.Status = 401
		res.Message = "gagal unmarshal JSON"
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

	query := "SELECT id_user_detail FROM user_detail WHERE user_detail_id = ? AND kode_otp = ? "
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	var user_id int
	err = stmt.QueryRow(requestacc.UserID, requestacc.Kode_OTP).Scan(&user_id)
	if err != nil {
		if err == sql.ErrNoRows {
			res.Status = 401
			res.Message = "Kode OTP tidak valid"
			res.Data = nil
			return res, nil
		}
		res.Status = 401
		res.Message = "exec gagal"
		res.Data = err.Error()
		return res, err
	}

	updatequery := "UPDATE user_detail SET status_verifikasi_otp = 'V' where user_detail_id = ?"
	updatestmt, err := con.Prepare(updatequery)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer updatestmt.Close()

	_, err = updatestmt.Exec(requestacc.UserID)
	if err != nil {
		res.Status = 401
		res.Message = "gagal mengupdate data"
		res.Data = err.Error()
		return res, err
	}

	res.Status = http.StatusOK
	res.Message = "Berhasil verifikasi OTP"
	res.Data = user_id

	defer db.DbClose(con)
	return res, nil
}

func VerifyUserAccept(input string) (Response, error) {
	var res Response

	type temp_verif_user_acc struct {
		UserID int `json:"userid"`
		Kelas  int `json:"kelas"`
	}
	var requestacc temp_verif_user_acc
	err := json.Unmarshal([]byte(input), &requestacc)
	if err != nil {
		res.Status = 401
		res.Message = "gagal unmarshal JSON"
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

	query := "UPDATE user_detail SET user_kelas_id=?,status='V' WHERE user_detail_id = ?"
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(requestacc.Kelas, requestacc.UserID)
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

func VerifyPerusahaanAccept(input string) (Response, error) {
	var res Response

	type temp_verif_perusahaan_acc struct {
		PerusahaanId int `json:"perusahaan_id"`
		Kelas        int `json:"kelas"`
	}
	var requestacc temp_verif_perusahaan_acc
	err := json.Unmarshal([]byte(input), &requestacc)
	if err != nil {
		res.Status = 401
		res.Message = "gagal unmarshal JSON"
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

	query := "UPDATE perusahaan SET kelas=?,status='V' WHERE perusahaan_id = ?"
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(requestacc.Kelas, requestacc.PerusahaanId)
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

func VerifyUserDecline(input string) (Response, error) {
	var res Response

	type temp_verif_user_deny struct {
		UserID int `json:"userid"`
	}
	var requestacc temp_verif_user_deny
	err := json.Unmarshal([]byte(input), &requestacc)
	if err != nil {
		res.Status = 401
		res.Message = "gagal unmarshal JSON"
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

	query := "UPDATE user_detail SET status='N',denied_by_admin='Y' WHERE user_detail_id = ?"
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(requestacc.UserID)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}

	// kirim notif (masih mendatang)

	res.Status = http.StatusOK
	res.Message = "Berhasil mengupdate data"
	res.Data = result

	defer db.DbClose(con)
	return res, nil
}

func VerifyPerusahaanDecline(input string) (Response, error) {
	var res Response

	type temp_verif_perusahaan_deny struct {
		PerusahaanId int `json:"perusahaan_id"`
	}
	var requestacc temp_verif_perusahaan_deny
	err := json.Unmarshal([]byte(input), &requestacc)
	if err != nil {
		res.Status = 401
		res.Message = "gagal unmarshal JSON"
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

	query := "UPDATE perusahaan SET status='N',denied_by_admin='Y' WHERE perusahaan_id = ?"
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(requestacc.PerusahaanId)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}

	// kirim notif (masih mendatang)

	res.Status = http.StatusOK
	res.Message = "Berhasil mengupdate data"
	res.Data = result

	defer db.DbClose(con)
	return res, nil
}

func VerifyAssetAccept(input string) (Response, error) {
	var res Response
	var dtSurveyReq SurveyRequest

	type temp_verif_asset_acc struct {
		SurveryReqId int `json:"surveyreq_id"`
	}

	var requestacc temp_verif_asset_acc
	err := json.Unmarshal([]byte(input), &requestacc)
	if err != nil {
		res.Status = 401
		res.Message = "gagal unmarshal JSON"
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

	query := "UPDATE survey_request SET status_request='F' WHERE id_transaksi_jual_sewa = ?"
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	_, err = stmt.Exec(requestacc.SurveryReqId)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}

	selectquery := "SELECT * FROM survey_request WHERE id_transaksi_jual_sewa = ?"
	selectstmt, err := con.Prepare(selectquery)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer selectstmt.Close()
	var temp_tags_old sql.NullString
	var temp_tags_new sql.NullString

	err = selectstmt.QueryRow(requestacc.SurveryReqId).Scan(&dtSurveyReq.Id_transaksi_jual_sewa, &dtSurveyReq.User_id, &dtSurveyReq.Id_asset, &dtSurveyReq.Dateline, &dtSurveyReq.Status_request, &dtSurveyReq.Data_lengkap, &dtSurveyReq.Usage_old, &dtSurveyReq.Usage_new, &dtSurveyReq.Luas_old, &dtSurveyReq.Luas_new, &dtSurveyReq.Nilai_old, &dtSurveyReq.Nilai_new, &dtSurveyReq.Kondisi_old, &dtSurveyReq.Kondisi_new, &dtSurveyReq.Batas_koordinat_old, &dtSurveyReq.Batas_koordinat_new, &temp_tags_old, &temp_tags_new)
	if err != nil {
		res.Status = 401
		res.Message = "rows gagal"
		res.Data = err.Error()
		return res, err
	}

	if temp_tags_old.Valid {
		dtSurveyReq.Tags_old = temp_tags_old.String
	} else {
		dtSurveyReq.Tags_old = ""
	}

	if temp_tags_new.Valid {
		dtSurveyReq.Tags_new = temp_tags_new.String
	} else {
		dtSurveyReq.Tags_new = ""
	}

	// update data asset dengan yang baru
	updatequery := "UPDATE asset SET `kondisi`= ?,`batas_koordinat`= ?,`luas`= ?,`nilai`= ?,`usage`= ? WHERE `id_asset`= ?"
	updatestmt, err := con.Prepare(updatequery)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer updatestmt.Close()

	_, err = updatestmt.Exec(dtSurveyReq.Kondisi_new, dtSurveyReq.Batas_koordinat_new, dtSurveyReq.Luas_new, dtSurveyReq.Nilai_new, dtSurveyReq.Usage_new, dtSurveyReq.Id_asset)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}

	update2query := "UPDATE asset_tags SET id_tags= ? WHERE id_asset= ?"
	update2stmt, err := con.Prepare(update2query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer update2stmt.Close()

	_, err = update2stmt.Exec(dtSurveyReq.Tags_new, dtSurveyReq.Id_asset)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}

	tempaset, _ := GetAssetById(string(dtSurveyReq.Id_asset))

	res.Status = http.StatusOK
	res.Message = "Berhasil mengupdate data"
	res.Data = tempaset

	defer db.DbClose(con)
	return res, nil
}

func ReassignAsset(input string) (Response, error) {
	var res Response

	type temp_verif_asset_acc struct {
		SurveyReqId int    `json:"surveyreq_id"`
		SurveyorId  int    `json:"surveyor_id"`
		Dateline    string `json:"dateline"`
	}

	var requestacc temp_verif_asset_acc
	err := json.Unmarshal([]byte(input), &requestacc)
	if err != nil {
		res.Status = 401
		res.Message = "gagal unmarshal JSON"
		res.Data = err.Error()
		return res, err
	}
	fmt.Println(requestacc)
	fmt.Println(requestacc.SurveyorId)

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	// ambil user id dari surveyor
	surveyorquery := "SELECT user_id FROM surveyor WHERE `suveyor_id` = ?"
	surveyorstmt, err := con.Prepare(surveyorquery)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer surveyorstmt.Close()

	var _tempuserid int
	err = surveyorstmt.QueryRow(requestacc.SurveyorId).Scan(&_tempuserid)
	if err != nil {
		res.Status = 401
		res.Message = "gagal mengambil user_id"
		res.Data = err.Error()
		return res, err
	}
	fmt.Println(requestacc.SurveyorId)
	fmt.Println(_tempuserid)

	query := "UPDATE survey_request SET `user_id`=?,`dateline`=?,`status_request`='R' WHERE id_transaksi_jual_sewa = ?"
	stmt2, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt2.Close()

	result, err := stmt2.Exec(_tempuserid, requestacc.Dateline, requestacc.SurveyReqId)
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

// BELUM SELESAI
func ForgotPass(email string) (Response, error) {
	var res Response

	// asd

	return res, nil
}

func ChangePass() (Response, error) {
	var res Response

	return res, nil
}

const CONFIG_SMTP_HOST = "smtp.gmail.com"
const CONFIG_SMTP_PORT = 587
const CONFIG_SENDER_NAME = "LEAP - Testing Kirim Email"
const CONFIG_AUTH_EMAIL = "c14210026@john.petra.ac.id"
const CONFIG_AUTH_PASSWORD = "alzx sjan ikkr ipsm"

func sendMail(to []string, cc []string, subject, message string) error {
	body := "From: " + CONFIG_SENDER_NAME + "\n" +
		"To: " + strings.Join(to, ",") + "\n" +
		"Cc: " + strings.Join(cc, ",") + "\n" +
		"Subject: " + subject + "\n\n" +
		message

	auth := smtp.PlainAuth("", CONFIG_AUTH_EMAIL, CONFIG_AUTH_PASSWORD, CONFIG_SMTP_HOST)
	smtpAddr := fmt.Sprintf("%s:%d", CONFIG_SMTP_HOST, CONFIG_SMTP_PORT)

	err := smtp.SendMail(smtpAddr, auth, CONFIG_AUTH_EMAIL, append(to, cc...), []byte(body))
	if err != nil {
		return err
	}

	return nil
}
