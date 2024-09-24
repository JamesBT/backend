package model

import (
	"TemplateProject/db"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
)

func CreateKelas(input string) (Response, error) {
	var res Response

	type InputRole struct {
		Class_name string `json:"nama"`
		Minimum    string `json:"minimum"`
		Maximum    string `json:"maximum"`
	}
	var dtKelas InputRole
	err := json.Unmarshal([]byte(input), &dtKelas)
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

	query := "INSERT INTO kelas (kelas_nama,kelas_modal_minimal,kelas_modal_maksimal) VALUES (?,?,?)"
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(dtKelas.Class_name, dtKelas.Minimum, dtKelas.Maximum)
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

	res.Status = http.StatusOK
	res.Message = "Berhasil membuat kelas"
	res.Data = map[string]interface{}{
		"id_role":        lastId,
		"nama kelas":     dtKelas.Class_name,
		"modal minimal":  dtKelas.Minimum,
		"modal maksimum": dtKelas.Maximum,
	}

	defer db.DbClose(con)
	return res, nil
}

func UpdateKelas(input string) (Response, error) {
	var res Response

	type InputRole struct {
		Class_id   int    `json:"id"`
		Class_name string `json:"nama"`
		Minimum    string `json:"minimum"`
		Maximum    string `json:"maximum"`
	}
	var dtKelas InputRole
	err := json.Unmarshal([]byte(input), &dtKelas)
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

	query := "UPDATE kelas SET kelas_nama = ?, kelas_modal_minimal = ?, kelas_modal_maksimal = ? WHERE kelas_id = ?"
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	_, err = stmt.Exec(dtKelas.Class_name, dtKelas.Minimum, dtKelas.Maximum, dtKelas.Class_id)
	if err != nil {
		res.Status = 401
		res.Message = "exec gagal"
		res.Data = err.Error()
		return res, err
	}

	res.Status = http.StatusOK
	res.Message = "Berhasil update kelas"
	res.Data = map[string]interface{}{
		"id_role":        dtKelas.Class_id,
		"nama kelas":     dtKelas.Class_name,
		"modal minimal":  dtKelas.Minimum,
		"modal maksimum": dtKelas.Maximum,
	}

	defer db.DbClose(con)
	return res, nil
}

func GetAllKelas() (Response, error) {
	var res Response
	var arrKelas = []Kelas{}

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	query := `
	SELECT *
	FROM kelas 
	WHERE kelas_id != 6
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

	for result.Next() {
		var dtKelas Kelas
		err = result.Scan(&dtKelas.Id, &dtKelas.Nama, &dtKelas.Modal_minimal, &dtKelas.Modal_maksimal)
		if err != nil {
			res.Status = 401
			res.Message = "rows scan"
			res.Data = err.Error()
			return res, err
		}

		arrKelas = append(arrKelas, dtKelas)
	}
	res.Status = http.StatusOK
	res.Message = "Berhasil mengambil data"
	res.Data = arrKelas

	defer db.DbClose(con)
	return res, nil
}

func GetKelasById(kelas_id string) (Response, error) {
	var res Response
	var dtKelas Kelas

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	query := `
	SELECT *
	FROM kelas
	WHERE kelas_id = ?
	`
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(kelas_id).Scan(&dtKelas.Id, &dtKelas.Nama, &dtKelas.Modal_minimal, &dtKelas.Modal_maksimal)
	if err != nil {
		res.Status = 401
		res.Message = "exec gagal"
		res.Data = err.Error()
		return res, err
	}

	res.Status = http.StatusOK
	res.Message = "Berhasil mengambil data"
	res.Data = dtKelas

	defer db.DbClose(con)
	return res, nil
}

func DeleteKelasById(id string) (Response, error) {
	var res Response
	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	checkAdminQuery := `
	SELECT kelas_id
	FROM kelas
	WHERE kelas_id = ?
	`
	var kelasId string
	err = con.QueryRow(checkAdminQuery, id).Scan(&kelasId)
	if err != nil {
		if err == sql.ErrNoRows {
		} else {
			res.Status = 401
			res.Message = "Gagal memeriksa kelas"
			res.Data = err.Error()
			return res, err
		}
	} else {
		res.Status = 403
		res.Message = "Kelas masih terpakai di user/perusahaan lain"
		res.Data = nil
		return res, errors.New(res.Message)
	}

	query := "DELETE FROM kelas WHERE kelas_id = ?"
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(id)
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
