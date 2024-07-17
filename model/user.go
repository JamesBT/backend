package model

import (
	"TemplateProject/db"
	"encoding/json"
	"errors"
	"net/http"
)

func Login(akun string) (Response, error) {
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
		res.Message = "Pengguna belum terdaftar"
		res.Data = err.Error()
		return res, errors.New("pengguna belum terdaftar")
	}
	defer stmt.Close()

	// cek apakah password benar atau tidak
	query = "SELECT user_id FROM user WHERE username = ? AND password = ?"
	stmt, err = con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}

	err = stmt.QueryRow(usr.Username, usr.Password).Scan(&userId)
	if err != nil {
		res.Status = 401
		res.Message = "username/password salah"
		res.Data = err.Error()
		return res, errors.New("username/password salah")
	}
	defer stmt.Close()

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
	res.Data = map[string]int{"id": userId}

	defer db.DbClose(con)

	return res, nil
}

func SignUp(akun string) (Response, error) {
	var res Response

	return res, nil
}

func ForgotPass(email string) (Response, error) {
	var res Response

	return res, nil
}

func ChangePass() (Response, error) {
	var res Response

	return res, nil
}
