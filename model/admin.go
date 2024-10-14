package model

import (
	"TemplateProject/db"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

func AddAdmin(filefoto *multipart.FileHeader, username, password, nama_lengkap, email, no_telp, role string) (Response, error) {
	var res Response

	type UserCompany struct {
		Id_user   string `json:"user_id"`
		Nama_user string `json:"nama"`
		Username  string `json:"username"`
		Password  string `json:"password"`
		Email     string `json:"email"`
		No_telp   string `json:"no_telp"`
		Id_role   string `json:"role_id"`
	}

	var tempUserCompany UserCompany
	tempUserCompany.Username = username
	tempUserCompany.Password = password
	tempUserCompany.Nama_user = nama_lengkap
	tempUserCompany.Email = email
	tempUserCompany.No_telp = no_telp
	tempUserCompany.Id_role = role

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	// Check if id_role exists
	var roleExists bool
	err = con.QueryRow("SELECT EXISTS(SELECT 1 FROM role WHERE role_id = ?)", tempUserCompany.Id_role).Scan(&roleExists)
	if err != nil {
		res.Status = 401
		res.Message = "gagal mengecek role"
		res.Data = err.Error()
		return res, err
	}
	if !roleExists {
		res.Status = 404
		res.Message = "Role tidak ditemukan"
		res.Data = nil
		return res, nil
	}

	var recordExists bool
	err = con.QueryRow("SELECT EXISTS(SELECT 1 FROM user WHERE email = ?)", tempUserCompany.Email).Scan(&recordExists)
	if err != nil {
		res.Status = 401
		res.Message = "gagal mengecek kombinasi user dan role"
		res.Data = err.Error()
		return res, err
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(tempUserCompany.Password), 10)
	if err != nil {
		res.Status = 401
		res.Message = "hashing gagal"
		res.Data = err.Error()
		return res, err
	}

	// If the record doesn't exist, insert a new one
	queryInsert := `
		INSERT INTO user (username, password, nama_lengkap, email, nomor_telepon) VALUES (?,?,?,?,?)
		`
	stmtInsert, err := con.Prepare(queryInsert)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmtInsert.Close()

	result, err := stmtInsert.Exec(tempUserCompany.Username, string(hashedPass), tempUserCompany.Nama_user, tempUserCompany.Email, tempUserCompany.No_telp)
	if err != nil {
		res.Status = 401
		res.Message = "exec gagal"
		res.Data = err.Error()
		return res, err
	}

	userId, err := result.LastInsertId()
	if err != nil {
		res.Status = 500
		res.Message = "gagal mendapatkan user ID"
		res.Data = err.Error()
		return res, err
	}
	tempUserCompany.Id_user = strconv.Itoa(int(userId))

	queryDetail := `
		INSERT INTO user_detail (user_detail_id, user_kelas_id, status, tipe, status_verifikasi_otp) VALUES (?,6,'V',?,'V')
		`
	stmtDetail, err := con.Prepare(queryDetail)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmtDetail.Close()

	_, err = stmtDetail.Exec(tempUserCompany.Id_user, tempUserCompany.Id_role)
	if err != nil {
		res.Status = 401
		res.Message = "exec gagal"
		res.Data = err.Error()
		return res, err
	}

	query := "INSERT INTO user_role (`user_id`, `role_id`) VALUES (?,?)"
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	_, err = stmt.Exec(tempUserCompany.Id_user, tempUserCompany.Id_role)
	if err != nil {
		res.Status = 401
		res.Message = "exec gagal"
		res.Data = err.Error()
		return res, err
	}

	// masukin foto profil
	filefoto.Filename = tempUserCompany.Id_user + "_" + filefoto.Filename
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

	nId, _ := strconv.Atoi(tempUserCompany.Id_user)
	err = UpdateDataFotoPath("user", "foto_profil", pathFotoFile, "user_id", nId)
	if err != nil {
		return res, err
	}

	res.Status = http.StatusOK
	res.Message = "Berhasil memasukkan admin"
	res.Data = tempUserCompany

	defer db.DbClose(con)
	return res, nil
}

func UpdateAdminRoleById(input string) (Response, error) {
	var res Response
	type TempAdmin struct {
		Id      int    `json:"id"`
		Role_id string `json:"role"`
	}

	var dtTempAdmin TempAdmin
	err := json.Unmarshal([]byte(input), &dtTempAdmin)
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

	query := "UPDATE user_role SET `role_id` = ? WHERE user_id = ? "
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(dtTempAdmin.Role_id, dtTempAdmin.Id)
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

func UpdateAdminById(input string) (Response, error) {
	var res Response
	type TempAdmin struct {
		Id           int    `json:"id"`
		Username     string `json:"username"`
		Password     string `json:"password"`
		Nama_lengkap string `json:"nama_lengkap"`
		Email        string `json:"email"`
		No_telp      string `json:"no_telp"`
	}

	var dtTempAdmin TempAdmin
	err := json.Unmarshal([]byte(input), &dtTempAdmin)
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

	query := "UPDATE user SET `username` = ?, `password` = ?, `nama_lengkap` = ?, `email`= ? , `nomor_telepon` = ? WHERE `user_id` = ? "
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(dtTempAdmin.Username, dtTempAdmin.Password, dtTempAdmin.Nama_lengkap, dtTempAdmin.Email, dtTempAdmin.No_telp, dtTempAdmin.Id)
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

func GetAllAdmin() (Response, error) {
	var res Response
	type TempUser struct {
		Id           int    `json:"id"`
		Username     string `json:"username"`
		Password     string `json:"password"`
		Nama_lengkap string `json:"nama_lengkap"`
		Email        string `json:"email"`
		No_telp      string `json:"no_telp"`
		Foto_profil  string `json:"foto_profil"`
		UserRole     string `json:"user_role"`
	}
	var arrAdmin = []TempUser{}

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	query := `
	SELECT u.user_id, u.username, u.password, u.nama_lengkap, 
	u.email, u.nomor_telepon, u.foto_profil, r.nama_role
	FROM user u
	LEFT JOIN user_role ur ON u.user_id = ur.user_id
	LEFT JOIN role r ON ur.role_id = r.role_id
	WHERE r.admin_role='Y' AND deleted_at IS NULL
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
		res.Message = "exec gagal"
		res.Data = err.Error()
		return res, err
	}
	defer rows.Close()

	for rows.Next() {
		var dtAdmin TempUser
		err = rows.Scan(&dtAdmin.Id, &dtAdmin.Username, &dtAdmin.Password, &dtAdmin.Nama_lengkap, &dtAdmin.Email,
			&dtAdmin.No_telp, &dtAdmin.Foto_profil, &dtAdmin.UserRole)
		if err != nil {
			res.Status = 401
			res.Message = "Failed to scan company details"
			res.Data = err.Error()
			return res, err
		}
		arrAdmin = append(arrAdmin, dtAdmin)
	}
	if err = rows.Err(); err != nil {
		res.Status = 401
		res.Message = "rows error"
		res.Data = err.Error()
		return res, err
	}

	if len(arrAdmin) == 0 {
		res.Status = 404
		res.Message = "Data tidak ditemukan"
		res.Data = nil
		return res, nil
	}

	res.Status = http.StatusOK
	res.Message = "Berhasil mengambil data"
	res.Data = arrAdmin

	defer db.DbClose(con)
	return res, nil
}

func GetAdminById(id_user string) (Response, error) {
	type TempRole struct {
		Role_id   int         `json:"role_id"`
		Nama_role string      `json:"nama_role"`
		Privilege []Privilege `json:"privilege"`
	}
	type User struct {
		Id               int      `json:"id"`
		Username         string   `json:"username"`
		Password         string   `json:"password"`
		Nama_lengkap     string   `json:"nama_lengkap"`
		Alamat           string   `json:"alamat"`
		Jenis_kelamin    string   `json:"jenis_kelamin"`
		Tgl_lahir        string   `json:"tgl_lahir"`
		Email            string   `json:"email"`
		No_telp          string   `json:"no_telp"`
		Foto_profil      string   `json:"foto_profil"`
		Ktp              string   `json:"ktp"`
		Kelas            int      `json:"kelas"`
		Status           string   `json:"status"`
		Tipe             int      `json:"tipe"`
		First_login      string   `json:"first_login"`
		Denied_by_admin  string   `json:"denied_by_admin"`
		UserRole         TempRole `json:"user_role"`
		PerusahaanJoined []Perusahaan
	}
	var res Response
	var usr User
	var dtRole TempRole
	var privileges []Privilege

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka koneksi"
		res.Data = err.Error()
		return res, err
	}

	query := `
    SELECT u.user_id, u.username, u.password, u.nama_lengkap, u.email, u.nomor_telepon, u.foto_profil, 
           r.role_id, r.nama_role, IFNULL(p.privilege_id,0), IFNULL(p.nama_privilege, '')  
    FROM user u
    LEFT JOIN user_role ur ON u.user_id = ur.user_id
    LEFT JOIN role r ON ur.role_id = r.role_id
    LEFT JOIN role_privilege_user rp ON r.role_id = rp.id_role
    LEFT JOIN privilege p ON rp.id_privilege = p.privilege_id
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

	rows, err := stmt.Query(id_user)
	if err != nil {
		res.Status = 401
		res.Message = "exec gagal"
		res.Data = err.Error()
		return res, err
	}
	defer rows.Close()

	isFirstRow := true
	var tempPrivilege Privilege

	for rows.Next() {
		err = rows.Scan(
			&usr.Id, &usr.Username, &usr.Password, &usr.Nama_lengkap,
			&usr.Email, &usr.No_telp, &usr.Foto_profil,
			&dtRole.Role_id, &dtRole.Nama_role, &tempPrivilege.Privilege_id, &tempPrivilege.Nama_privilege,
		)
		if err != nil {
			res.Status = 401
			res.Message = "scan gagal"
			res.Data = err.Error()
			return res, err
		}

		if tempPrivilege.Nama_privilege != "" {
			privileges = append(privileges, tempPrivilege)
		}

		if isFirstRow {
			isFirstRow = false
		}
	}

	// Populate role and privileges
	dtRole.Privilege = privileges
	usr.UserRole = dtRole

	res.Status = http.StatusOK
	res.Message = "Berhasil mengambil data"
	res.Data = usr

	defer db.DbClose(con)

	return res, nil
}

func LoginAdmin(akun string) (Response, error) {
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

	// Cek apakah user terdaftar dan memiliki peran admin
	query := `
	SELECT u.user_id 
	FROM user u
	LEFT JOIN user_role ur ON u.user_id = ur.user_id
	LEFT JOIN role r ON ur.role_id = r.role_id 
	WHERE u.username = ? AND u.deleted_at IS NULL AND r.admin_role = 'Y'
	`
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

	fmt.Println("user id: ", userId)

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

	// Cek apakah password benar
	queryCheck := "SELECT u.user_id, u.username, u.nama_lengkap, u.alamat, u.jenis_kelamin, IFNULL(u.tanggal_lahir,''), u.email, u.nomor_telepon, u.foto_profil, u.ktp, ud.user_kelas_id, ud.status, ud.tipe FROM user u JOIN user_detail ud ON u.user_id = ud.user_detail_id WHERE u.user_id = ?;"
	stmtCheck, err := con.Prepare(queryCheck)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmtCheck.Close()

	err = stmtCheck.QueryRow(userId).Scan(&loginUsr.Id, &loginUsr.Username, &loginUsr.Nama_lengkap, &loginUsr.Alamat, &loginUsr.Jenis_kelamin, &loginUsr.Tgl_lahir, &loginUsr.Email, &loginUsr.No_telp, &loginUsr.Foto_profil, &loginUsr.Ktp, &loginUsr.Kelas, &loginUsr.Status, &loginUsr.Tipe)
	if err != nil {
		res.Status = 401
		res.Message = "password salah"
		res.Data = err.Error()
		return res, errors.New("password salah")
	}

	// Ambil role + privilege
	getRoleQuery := "SELECT ur.role_id, r.nama_role FROM user_role ur JOIN role r ON ur.role_id = r.role_id WHERE ur.user_id = ?;"
	roleStmt, err := con.Prepare(getRoleQuery)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer roleStmt.Close()

	var roleId int
	var roleName string
	err = roleStmt.QueryRow(userId).Scan(&roleId, &roleName)
	if err != nil {
		res.Status = 401
		res.Message = "gagal mendapatkan role"
		res.Data = err.Error()
		return res, err
	}

	// Ambil semua privilege dari role yang diambil
	privilegesQuery := `
	SELECT p.privilege_id, p.nama_privilege
	FROM role_privilege_user rpu
	JOIN privilege p ON rpu.id_privilege = p.privilege_id
	WHERE rpu.id_role = ?
	`
	privStmt, err := con.Prepare(privilegesQuery)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer privStmt.Close()

	rows, err := privStmt.Query(roleId)
	if err != nil {
		res.Status = 401
		res.Message = "query gagal"
		res.Data = err.Error()
		return res, err
	}
	defer rows.Close()

	fmt.Println("role: ", roleId)

	// Array to hold all privileges
	var privileges []map[string]interface{}
	for rows.Next() {
		var privilegeId int
		var privilegeName string
		err := rows.Scan(&privilegeId, &privilegeName)
		if err != nil {
			res.Status = 401
			res.Message = "gagal mendapatkan privilege"
			res.Data = err.Error()
			return res, err
		}
		privileges = append(privileges, map[string]interface{}{
			"privilege_id":   privilegeId,
			"nama_privilege": privilegeName,
		})
	}

	// Berhasil login => update timestamp terakhir login
	updateQuery := "UPDATE user SET login_timestamp = NOW() WHERE user_id = ?"
	updateStmt, err := con.Prepare(updateQuery)
	if err != nil {
		res.Status = 401
		res.Message = "stmt update gagal"
		res.Data = err.Error()
		return res, err
	}
	defer updateStmt.Close()

	_, err = updateStmt.Exec(userId)
	if err != nil {
		res.Status = 401
		res.Message = "update login_timestamp gagal"
		res.Data = err.Error()
		return res, err
	}

	// Response success
	res.Status = http.StatusOK
	res.Message = "Berhasil login"
	res.Data = map[string]interface{}{
		"id":            loginUsr.Id,
		"username":      loginUsr.Username,
		"nama_lengkap":  loginUsr.Nama_lengkap,
		"alamat":        loginUsr.Alamat,
		"jenis_kelamin": loginUsr.Jenis_kelamin,
		"tanggal_lahir": loginUsr.Tgl_lahir,
		"email":         loginUsr.Email,
		"nomor_telepon": loginUsr.No_telp,
		"foto_profil":   loginUsr.Foto_profil,
		"ktp":           loginUsr.Ktp,
		"status":        loginUsr.Status,
		"tipe":          loginUsr.Tipe,
		"role_id":       roleId,
		"role_nama":     roleName,
		"privileges":    privileges, // Add privileges to the response
	}

	defer db.DbClose(con)
	return res, nil
}

func DeleteAdmin(id_user string) (Response, error) {
	var res Response

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	checkAdminQuery := `
	SELECT u.user_id
	FROM user u
	LEFT JOIN user_role ur ON u.user_id = ur.user_id
	LEFT JOIN role r ON ur.role_id = r.role_id
	WHERE u.user_id = ? AND r.admin_role = 'Y' AND u.deleted_at IS NULL
	`
	var userID string
	err = con.QueryRow(checkAdminQuery, id_user).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			res.Status = 403
			res.Message = "User bukan admin"
			res.Data = nil
			return res, errors.New(res.Message)
		}
		res.Status = 401
		res.Message = "Gagal memeriksa role user"
		res.Data = err.Error()
		return res, err
	}

	query := `
	UPDATE user SET deleted_at = NOW() 
	WHERE user_id = ?`
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
