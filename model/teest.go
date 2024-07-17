package model

import (
	"TemplateProject/db"
	"encoding/json"
	"net/http"
	"strconv"
)

func KirimAngka(jumlah_angka string) (Response, error) {
	var res Response
	var arrInt = []int{}
	jumlah_angka_int, err := strconv.Atoi(jumlah_angka)

	if err != nil {
		res.Status = 401
		res.Message = "parameter tidak valid"
		res.Data = err.Error()
		return res, err
	}

	for i := 0; i < jumlah_angka_int; i++ {
		arrInt = append(arrInt, i)
	}

	res.Status = http.StatusOK
	res.Message = "Array Integer sampai " + jumlah_angka
	res.Data = arrInt
	return res, nil
}

func GetAllBarang() (Response, error) {
	var res Response
	var obj Barang
	var arrObj = []Barang{}
	con, err := db.DbConnection()

	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka koneksi"
		res.Data = err.Error()
		return res, err
	}

	query := "SELECT * FROM barang"
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
		res.Message = "rows gagal"
		res.Data = err.Error()
		return res, err
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&obj.Id, &obj.Nama, &obj.Harga)
		if err != nil {
			res.Status = 401
			res.Message = "rows scan"
			res.Data = err.Error()
			return res, err
		}
		arrObj = append(arrObj, obj)
	}
	res.Status = http.StatusOK
	res.Message = "Berhasil mengambil data"
	res.Data = arrObj

	defer db.DbClose(con)

	return res, nil
}

func GetBarangById(id_barang string) (Response, error) {
	var res Response
	var obj Barang
	// var arrObj = []Barang{}
	con, err := db.DbConnection()

	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka koneksi"
		res.Data = err.Error()
		return res, err
	}

	query := "SELECT id,nama,harga FROM barang WHERE id = ?"
	stmt, err := con.Prepare(query)

	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()
	nId, _ := strconv.Atoi(id_barang)
	err = stmt.QueryRow(nId).Scan(&obj.Id, &obj.Nama, &obj.Harga)
	if err != nil {
		res.Status = 401
		res.Message = "rows gagal"
		res.Data = err.Error()
		return res, err
	}
	// defer rows.Close()
	// for rows.Next() {
	// 	err = rows.Scan(&obj.Id, &obj.Nama, &obj.Harga)
	// 	if err != nil {
	// 		res.Status = 401
	// 		res.Message = "rows scan"
	// 		res.Data = err.Error()
	// 		return res, err
	// 	}
	// 	arrObj = append(arrObj, obj)
	// }
	res.Status = http.StatusOK
	res.Message = "Berhasil mengambil data"
	res.Data = obj

	defer db.DbClose(con)

	return res, nil
}

func InsertBarang(barang string) (Response, error) {
	var res Response
	// var obj Barang
	var arrObj = []Barang{}

	err := json.Unmarshal([]byte(barang), &arrObj)

	if err != nil {
		res.Status = 401
		res.Message = "gagal decode json"
		res.Data = err.Error()
		return res, err
	}
	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka koneksi database"
		res.Data = err.Error()
		return res, err
	}

	query := "INSERT INTO barang (nama,harga) VALUES (?,?)"
	stmt, err := con.Prepare(query)

	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()
	for i, x := range arrObj {
		result, err := stmt.Exec(x.Nama, x.Harga)
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
		arrObj[i].Id = int(lastId)
	}
	// defer rows.Close()
	// for rows.Next() {
	// 	err = rows.Scan(&obj.Id, &obj.Nama, &obj.Harga)
	// 	if err != nil {
	// 		res.Status = 401
	// 		res.Message = "rows scan"
	// 		res.Data = err.Error()
	// 		return res, err
	// 	}
	// 	arrObj = append(arrObj, obj)
	// }
	res.Status = http.StatusOK
	res.Message = "Berhasil mengambil data"
	res.Data = arrObj

	defer db.DbClose(con)

	return res, nil
}
