package model

import (
	"TemplateProject/db"
	"net/http"
)

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
