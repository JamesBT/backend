package model

import (
	"TemplateProject/db"
	"net/http"
)

func GetAllBusinessField() (Response, error) {
	var res Response
	var arrField = []BusinessField{}

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	query := `
	SELECT *
	FROM business_field
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
		var dtField BusinessField
		err = result.Scan(&dtField.Id, &dtField.Nama, &dtField.Detail)
		if err != nil {
			res.Status = 401
			res.Message = "rows scan"
			res.Data = err.Error()
			return res, err
		}

		arrField = append(arrField, dtField)
	}
	res.Status = http.StatusOK
	res.Message = "Berhasil mengambil data"
	res.Data = arrField

	defer db.DbClose(con)
	return res, nil
}
