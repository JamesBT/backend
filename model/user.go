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
	"os"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
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

	var tempDBPass string
	querypass := `SELECT password FROM user WHERE user_id = ?;`
	stmtpass, err := con.Prepare(querypass)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmtpass.Close()

	err = stmtpass.QueryRow(userId).Scan(&tempDBPass)
	if err != nil {
		res.Status = 401
		res.Message = "query password gagal"
		res.Data = err.Error()
		return res, errors.New("query password gagal")
	}

	// cek pass sama atau tidak
	err = bcrypt.CompareHashAndPassword([]byte(tempDBPass), []byte(usr.Password))
	if err != nil {
		res.Status = 404
		res.Message = "password salah"
		res.Data = err.Error()
		return res, err
	}

	// cek apakah password benar atau tidak
	// queryinsert := "SELECT user_id, username, nama_lengkap, alamat, jenis_kelamin, tanggal_lahir, email, nomor_telepon, foto_profil, ktp FROM user WHERE username = ? AND password = ?"
	queryinsert := `SELECT u.user_id, u.username, u.nama_lengkap, u.alamat, u.jenis_kelamin, IFNULL(u.tanggal_lahir,""), u.email, u.nomor_telepon, u.foto_profil, u.ktp, ud.user_kelas_id, ud.status, ud.tipe, ud.first_login, ud.denied_by_admin FROM user u JOIN user_detail ud ON u.user_id = ud.user_detail_id WHERE u.user_id = ?;`
	stmtinsert, err := con.Prepare(queryinsert)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmtinsert.Close()

	err = stmtinsert.QueryRow(userId).Scan(&loginUsr.Id, &loginUsr.Username, &loginUsr.Nama_lengkap, &loginUsr.Alamat, &loginUsr.Jenis_kelamin, &loginUsr.Tgl_lahir, &loginUsr.Email, &loginUsr.No_telp, &loginUsr.Foto_profil, &loginUsr.Ktp, &loginUsr.Kelas, &loginUsr.Status, &loginUsr.Tipe, &loginUsr.First_login, &loginUsr.Denied_by_admin)
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
	}

	// cek email sudah terpakai 3 kali atau tidak
	cekEmailQuery := `
		SELECT email FROM user WHERE email = ?
	`
	stmtEmailQuery, err := con.Prepare(cekEmailQuery)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}

	var userEmail string
	err = stmtEmailQuery.QueryRow(usr.Email).Scan(&userEmail)
	if err == nil {
		if userEmail != "" {
			res.Status = 401
			res.Message = "Email already registered"
			res.Data = "User email: " + fmt.Sprint(usr.Email)
			return res, errors.New("user already registered")
		} else {
			res.Status = 500
			res.Message = "Query execution failed"
			res.Data = errors.New("query exec failed")
			return res, err
		}
	} else if err != sql.ErrNoRows {
		res.Status = 500
		res.Message = "Query execution failed"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(usr.Password), 10)
	if err != nil {
		res.Status = 401
		res.Message = "hashing gagal"
		res.Data = err.Error()
		return res, err
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
	result, err := insertstmt.Exec(usr.Username, string(hashedPass), usr.Nama_lengkap, usr.Email, usr.No_telp)
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
	insertdetailquery := "INSERT INTO user_detail (user_detail_id,user_kelas_id,tipe,kode_otp) VALUES (?,?,?,?)"
	insertdetailstmt, err := con.Prepare(insertdetailquery)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer insertstmt.Close()

	_, err = insertdetailstmt.Exec(usr.Id, 1, 8, randomnumber)
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
	LEFT JOIN user_role ur ON u.user_id = ur.user_id
	LEFT JOIN role r ON ur.role_id = r.role_id
	LEFT JOIN privilege p ON up.privilege_id = p.privilege_id
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
		err = result.Scan(&dtUser.Id, &dtUser.Username, &dtUser.Password, &dtUser.Nama_lengkap, &dtUser.Alamat, &dtUser.Jenis_kelamin, &dtUser.Tgl_lahir, &dtUser.Email, &dtUser.No_telp, &dtUser.Foto_profil, &dtUser.Ktp, &dtUser.Kelas, &dtUser.Status, &dtUser.Tipe, &dtUser.First_login, &dtUser.Denied_by_admin, &roleId, &roleName)
		if err != nil {
			res.Status = 401
			res.Message = "rows scan"
			res.Data = err.Error()
			return res, err
		}

		if existingUser, ok := userMap[dtUser.Id]; ok {
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
		} else {
			dtUser.UserRole = []Role{{Role_id: roleId, Nama_role: roleName}}
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

	query := `
	SELECT u.user_id, u.username, u.password, u.nama_lengkap, u.alamat, u.jenis_kelamin, IFNULL(u.tanggal_lahir,""), 
	u.email, u.nomor_telepon, u.foto_profil, u.ktp, ur.role_id, r.nama_role 
	FROM user u
	LEFT JOIN user_role ur ON u.user_id = ur.user_id
	LEFT JOIN role r ON ur.role_id = r.role_id
	WHERE u.user_id = ?
	`
	stmt, err := con.Prepare(query)

	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	nId, _ := strconv.Atoi(id_user)
	var tempDtRole Role
	err = stmt.QueryRow(nId).Scan(&usr.Id, &usr.Username, &usr.Password, &usr.Nama_lengkap, &usr.Alamat, &usr.Jenis_kelamin, &usr.Tgl_lahir, &usr.Email, &usr.No_telp, &usr.Foto_profil, &usr.Ktp, &tempDtRole.Role_id, &tempDtRole.Nama_role)
	if err != nil {
		res.Status = 401
		res.Message = "rows gagal"
		res.Data = err.Error()
		return res, err
	}

	usr.UserRole = append(usr.UserRole, tempDtRole)

	res.Status = http.StatusOK
	res.Message = "Berhasil mengambil data"
	res.Data = usr

	defer db.DbClose(con)

	return res, nil
}

func GetUserDetailedById(id_user string) (Response, error) {
	var res Response
	var usr User
	var perusahaanList []Perusahaan

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka koneksi"
		res.Data = err.Error()
		return res, err
	}

	query := "SELECT user_id, username, nama_lengkap, password, email, nomor_telepon,foto_profil,ktp FROM user WHERE user_id = ?"
	stmt, err := con.Prepare(query)

	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	nId, _ := strconv.Atoi(id_user)
	err = stmt.QueryRow(nId).Scan(&usr.Id, &usr.Username, &usr.Nama_lengkap, &usr.Password, &usr.Email, &usr.No_telp, &usr.Foto_profil, &usr.Ktp)
	if err != nil {
		res.Status = 401
		res.Message = "rows gagal"
		res.Data = err.Error()
		return res, err
	}

	query = `
		SELECT p.perusahaan_id, p.name, p.username, p.lokasi, p.tipe, p.modal_awal, p.deskripsi, p.created_at 
		FROM perusahaan p
		LEFT JOIN user_perusahaan up ON p.perusahaan_id = up.id_perusahaan
		WHERE up.id_user = ?
	`
	stmt, err = con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "Failed to prepare statement"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(nId)
	if err != nil {
		res.Status = 401
		res.Message = "Failed to execute query"
		res.Data = err.Error()
		return res, err
	}
	defer rows.Close()

	for rows.Next() {
		var perusahaan Perusahaan
		err = rows.Scan(&perusahaan.Id, &perusahaan.Nama, &perusahaan.Username, &perusahaan.Lokasi, &perusahaan.Tipe, &perusahaan.Modal, &perusahaan.Deskripsi, &perusahaan.CreatedAt)
		if err != nil {
			res.Status = 401
			res.Message = "Failed to scan company details"
			res.Data = err.Error()
			return res, err
		}
		perusahaanList = append(perusahaanList, perusahaan)
	}

	usr.PerusahaanJoined = perusahaanList

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

	query := "UPDATE user SET username = ?, nama_lengkap = ?, `alamat` = ?, jenis_kelamin = ?, tanggal_lahir = ?, email = ?, nomor_telepon = ?,updated_at = NOW() WHERE user_id = ? "
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
	// tipe := filefoto.Header.Get("Content-type")

	filefoto.Filename = userid + "_" + filefoto.Filename
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

func UpdateUserById(filefoto *multipart.FileHeader, id, username, password, nama_lengkap, no_telp string) (Response, error) {
	var res Response

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	var tempPass string
	queryPass := "SELECT password FROM user WHERE user_id = ?"
	stmtPass, err := con.Prepare(queryPass)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmtPass.Close()
	err = stmtPass.QueryRow(id).Scan(&tempPass)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}

	var hashedPass string
	if tempPass != password {
		tempHashedPass, err := bcrypt.GenerateFromPassword([]byte(password), 10)
		if err != nil {
			res.Status = 401
			res.Message = "stmt gagal"
			res.Data = err.Error()
			return res, err
		}
		hashedPass = string(tempHashedPass)
	} else {
		hashedPass = password
	}

	query := "UPDATE user SET username = ?, password = ?, nama_lengkap = ?, nomor_telepon = ?, updated_at = NOW() WHERE user_id = ? "
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(username, hashedPass, nama_lengkap, no_telp, id)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}

	filefoto.Filename = id + "_" + filefoto.Filename
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

	nId, _ := strconv.Atoi(id)
	err = UpdateDataFotoPath("user", "foto_profil", pathFotoFile, "user_id", nId)
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

func UpdateUserByIdTanpaFoto(id, username, password, nama_lengkap, no_telp string) (Response, error) {
	var res Response

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	var tempPass string
	queryPass := "SELECT password FROM user WHERE user_id = ?"
	stmtPass, err := con.Prepare(queryPass)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmtPass.Close()
	err = stmtPass.QueryRow(id).Scan(&tempPass)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}

	var hashedPass string
	if tempPass != password {
		tempHashedPass, err := bcrypt.GenerateFromPassword([]byte(password), 10)
		if err != nil {
			res.Status = 401
			res.Message = "stmt gagal"
			res.Data = err.Error()
			return res, err
		}
		hashedPass = string(tempHashedPass)
	} else {
		hashedPass = password
	}

	query := "UPDATE user SET username = ?, password = ?, nama_lengkap = ?, nomor_telepon = ?, updated_at = NOW() WHERE user_id = ? "
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(username, hashedPass, nama_lengkap, no_telp, id)
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
	for result.Next() {
		var dtUser User
		err = result.Scan(&dtUser.Id, &dtUser.Username, &dtUser.Password, &dtUser.Nama_lengkap, &dtUser.Alamat, &dtUser.Jenis_kelamin, &dtUser.Tgl_lahir, &dtUser.Email, &dtUser.No_telp, &dtUser.Foto_profil, &dtUser.Ktp)
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

func AdminUserManagement() (Response, error) {
	var res Response
	type TempAdminUser struct {
		IdUser          int    `json:"id"`
		Nama            string `json:"nama"`
		Foto_profil     string `json:"foto_profil"`
		TotalPerusahaan int    `json:"totalPerusahaan"`
		TotalTransaksi  int    `json:"totalTransaksi"`
	}
	var arrUser = []TempAdminUser{}

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	query := `
	WITH TransactionCounts AS (
		SELECT 
			tr.user_id, 
			COUNT(DISTINCT tr.id_asset) AS appearance_count
		FROM 
			transaction_request tr
		WHERE 
			tr.status IN ('A', 'D')
		GROUP BY 
			tr.user_id
	)

	SELECT 
		u.nama_lengkap, 
		u.user_id, 
		u.foto_profil,
		COUNT(DISTINCT CASE WHEN p.status = 'A' THEN up.id_perusahaan END) AS totalPerusahaan, 
		COALESCE(tc.appearance_count, 0) AS totalTransaksi
	FROM 
		user u 
	INNER JOIN 
		user_detail ud ON u.user_id = ud.user_detail_id
	LEFT JOIN 
		user_perusahaan up ON u.user_id = up.id_user
	LEFT JOIN 
		TransactionCounts tc ON u.user_id = tc.user_id
	LEFT JOIN 
		user_role ur ON u.user_id = ur.user_id
	LEFT JOIN 
		role r ON ur.role_id = r.role_id
	LEFT JOIN 
		perusahaan p ON up.id_perusahaan = p.perusahaan_id
	WHERE 
		ud.status = 'V' AND r.admin_role = 'N' 
	GROUP BY 
		u.user_id, 
		u.nama_lengkap, 
		u.foto_profil;

	`

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
		res.Message = "Failed to execute query"
		res.Data = err.Error()
		return res, err
	}
	defer rows.Close()

	for rows.Next() {
		var tempUser TempAdminUser
		var totalTransaksi sql.NullInt64
		err = rows.Scan(&tempUser.Nama, &tempUser.IdUser, &tempUser.Foto_profil, &tempUser.TotalPerusahaan, &totalTransaksi)
		if err != nil {
			res.Status = 401
			res.Message = "Failed to scan row"
			res.Data = err.Error()
			return res, err
		}

		if totalTransaksi.Valid {
			tempUser.TotalTransaksi = int(totalTransaksi.Int64)
		} else {
			tempUser.TotalTransaksi = 0
		}
		arrUser = append(arrUser, tempUser)
	}

	if err = rows.Err(); err != nil {
		res.Status = 401
		res.Message = "Row iteration failed"
		res.Data = err.Error()
		return res, err
	}

	res.Status = http.StatusOK
	res.Message = "Berhasil mengambil data"
	res.Data = arrUser

	defer db.DbClose(con)
	return res, nil
}

func CobaHashing(input string) (Response, error) {
	var res Response
	hasilHash, err := bcrypt.GenerateFromPassword([]byte(input), 10)
	if err != nil {
		res.Status = 401
		res.Message = "Error"
		return res, err
	}

	res.Status = 200
	res.Message = "hashing berhasil"
	res.Data = string(hasilHash)
	return res, nil
}

func SamainPassword(input, id string) (Response, error) {
	var res Response

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	query := "SELECT password FROM user WHERE user_id = ?"
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	var dtPass string
	err = stmt.QueryRow(id).Scan(&dtPass)
	if err != nil {
		res.Status = 404
		res.Message = "password not found"
		res.Data = err.Error()
		return res, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(dtPass), []byte(input))
	if err != nil {
		res.Status = 404
		res.Message = "password salah"
		res.Data = err.Error()
		return res, err
	}

	res.Status = http.StatusOK
	res.Message = "Berhasil mengcek password"

	defer db.DbClose(con)

	return res, nil
}
