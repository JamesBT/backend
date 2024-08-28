package model

import (
	"TemplateProject/db"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
)

// CRUD perusahaan ============================================================================
func CreatePerusahaan(file_kepemilikan *multipart.FileHeader, file_perusahaan *multipart.FileHeader, userid, nama, username, lokasi, tipe, modal, deskripsi string) (Response, error) {
	var res Response
	var dtPerusahaan = Perusahaan{}

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	query := "INSERT INTO perusahaan (name, username, lokasi, tipe, modal_awal, deskripsi) VALUES (?,?,?,?,?,?)"
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(nama, username, lokasi, tipe, modal, deskripsi)
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
	fmt.Println("lastid", lastId)
	// insert file kepemilikan
	templastid := strconv.Itoa(int(lastId))
	file_kepemilikan.Filename = templastid + "_" + file_kepemilikan.Filename
	pathFile := "uploads/perusahaan/file_kepemilikan/" + file_kepemilikan.Filename
	//source
	srcfoto, err := file_kepemilikan.Open()
	if err != nil {
		log.Println(err.Error())
		return res, err
	}
	defer srcfoto.Close()

	// Destination
	dstfoto, err := os.Create("uploads/perusahaan/file_kepemilikan/" + file_kepemilikan.Filename)
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

	err = UpdateDataFotoPath("perusahaan", "dokumen_kepemilikan", pathFile, "perusahaan_id", int(lastId))
	if err != nil {
		return res, err
	}

	// insert file perusahaan
	file_perusahaan.Filename = templastid + "_" + file_perusahaan.Filename
	pathFilePerusahaan := "uploads/perusahaan/file_perusahaan/" + file_perusahaan.Filename
	//source
	srcfotoPerusahaan, err := file_perusahaan.Open()
	if err != nil {
		log.Println(err.Error())
		return res, err
	}
	defer srcfoto.Close()

	// Destination
	dstfotoPerusahaan, err := os.Create("uploads/perusahaan/file_perusahaan/" + file_perusahaan.Filename)
	if err != nil {
		log.Println("2")
		fmt.Print("2")
		return res, err
	}

	// Copy
	if _, err = io.Copy(dstfotoPerusahaan, srcfotoPerusahaan); err != nil {
		log.Println(err.Error())
		log.Println("3")
		fmt.Print("3")
		return res, err
	}
	dstfoto.Close()

	err = UpdateDataFotoPath("perusahaan", "dokumen_perusahaan", pathFilePerusahaan, "perusahaan_id", int(lastId))
	if err != nil {
		return res, err
	}

	dtPerusahaan.Id = int(lastId)
	dtPerusahaan.Nama = nama
	dtPerusahaan.Username = username
	dtPerusahaan.Lokasi = lokasi
	dtPerusahaan.Tipe = tipe
	dtPerusahaan.Dokumen_kepemilikan = pathFile
	dtPerusahaan.Dokumen_perusahaan = pathFilePerusahaan
	dtPerusahaan.Modal = modal
	dtPerusahaan.Deskripsi = deskripsi

	// masukin ke user_perusahaan
	queryuser := "INSERT INTO user_perusahaan (id_user,id_perusahaan) VALUES (?,?)"
	stmtuser, err := con.Prepare(queryuser)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmtuser.Close()

	_, err = stmtuser.Exec(userid, dtPerusahaan.Id)
	if err != nil {
		res.Status = 401
		res.Message = "exec gagal"
		res.Data = err.Error()
		return res, err
	}

	res.Status = http.StatusOK
	res.Message = "Berhasil memasukkan data"
	res.Data = dtPerusahaan

	defer db.DbClose(con)
	return res, nil
}

func GetAllPerusahaanUnverified() (Response, error) {
	var res Response
	var arrPerusahaan = []Perusahaan{}
	var dtPerusahaan Perusahaan

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	query := "SELECT perusahaan_id, status, name, username, lokasi, tipe, dokumen_kepemilikan, dokumen_perusahaan, modal_awal, deskripsi, created_at FROM perusahaan WHERE status = 'N'"
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
		err = result.Scan(&dtPerusahaan.Id, &dtPerusahaan.Status, &dtPerusahaan.Nama, &dtPerusahaan.Username, &dtPerusahaan.Lokasi, &dtPerusahaan.Tipe, &dtPerusahaan.Dokumen_kepemilikan, &dtPerusahaan.Dokumen_perusahaan, &dtPerusahaan.Modal, &dtPerusahaan.Deskripsi, &dtPerusahaan.CreatedAt)
		if err != nil {
			res.Status = 401
			res.Message = "rows scan"
			res.Data = err.Error()
			return res, err
		}
		arrPerusahaan = append(arrPerusahaan, dtPerusahaan)
	}

	res.Status = http.StatusOK
	res.Message = "Berhasil mengambil data"
	res.Data = arrPerusahaan

	defer db.DbClose(con)
	return res, nil
}

func GetPerusahaanDetailById(id_perusahaan string) (Response, error) {
	var res Response
	var dtPerusahaan Perusahaan
	var dtUser []User

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka koneksi"
		res.Data = err.Error()
		return res, err
	}

	query := `
		SELECT p.perusahaan_id, p.username, p.lokasi, p.tipe, p.kelas, p.dokumen_kepemilikan, p.dokumen_perusahaan,
			p.modal_awal,p.deskripsi, u.user_id, u.username 
		FROM perusahaan p
		LEFT JOIN user_perusahaan up ON p.perusahaan_id = up.id_perusahaan
		LEFT JOIN user u on up.id_user = u.user_id
		WHERE p.perusahaan_id = ?
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
	rows, err := stmt.Query(nId)
	if err != nil {
		res.Status = 401
		res.Message = "rows gagal"
		res.Data = err.Error()
		return res, err
	}
	defer rows.Close()

	for rows.Next() {
		var usr User
		err = rows.Scan(&dtPerusahaan.Id, &dtPerusahaan.Username, &dtPerusahaan.Lokasi, &dtPerusahaan.Tipe,
			&dtPerusahaan.Kelas,
			&dtPerusahaan.Dokumen_kepemilikan, &dtPerusahaan.Dokumen_perusahaan, &dtPerusahaan.Modal,
			&dtPerusahaan.Deskripsi, &usr.Id, &usr.Username)
		if err != nil {
			res.Status = 401
			res.Message = "Failed to scan row"
			res.Data = err.Error()
			return res, err
		}
		if usr.Id != 0 {
			dtUser = append(dtUser, usr)
		}
	}

	dtPerusahaan.UserJoined = dtUser

	res.Status = http.StatusOK
	res.Message = "Berhasil mengambil data"
	res.Data = dtPerusahaan

	defer db.DbClose(con)

	return res, nil
}

func GetPerusahaanByUserId(user_id string) (Response, error) {
	var res Response
	var arrPerusahaan = []Perusahaan{}

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	query := `
		SELECT p.perusahaan_id, p.status, p.name, p.username, p.lokasi, p.tipe, 
		       p.dokumen_kepemilikan, p.dokumen_perusahaan, p.modal_awal, p.deskripsi, p.created_at
		FROM perusahaan p
		JOIN user_perusahaan up ON p.perusahaan_id = up.id_perusahaan
		WHERE up.id_user = ? AND p.status = 'N'
	`
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(user_id)
	if err != nil {
		res.Status = 401
		res.Message = "Failed to execute query"
		res.Data = err.Error()
		return res, err
	}
	defer rows.Close()
	for rows.Next() {
		var dtPerusahaan Perusahaan
		err = rows.Scan(
			&dtPerusahaan.Id, &dtPerusahaan.Status, &dtPerusahaan.Nama, &dtPerusahaan.Username,
			&dtPerusahaan.Lokasi, &dtPerusahaan.Tipe, &dtPerusahaan.Dokumen_kepemilikan,
			&dtPerusahaan.Dokumen_perusahaan, &dtPerusahaan.Modal, &dtPerusahaan.Deskripsi,
			&dtPerusahaan.CreatedAt,
		)
		if err != nil {
			res.Status = 401
			res.Message = "Failed to scan row"
			res.Data = err.Error()
			return res, err
		}
		arrPerusahaan = append(arrPerusahaan, dtPerusahaan)
	}
	if err = rows.Err(); err != nil {
		res.Status = 401
		res.Message = "Error during row iteration"
		res.Data = err.Error()
		return res, err
	}

	res.Status = http.StatusOK
	res.Message = "Berhasil mengambil data"
	res.Data = arrPerusahaan

	defer db.DbClose(con)
	return res, nil
}

func HomeUserPerusahaan(perusahaan_id string) (Response, error) {
	var res Response

	// total transaction request + semua aset
	type HomeUserPerusahaan struct {
		TotalRequest int     `json:"total_request"`
		SemuaAsset   []Asset `json:"asset"`
	}
	var homeuser HomeUserPerusahaan

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka koneksi"
		res.Data = err.Error()
		return res, err
	}

	query := `
		SELECT (
			SELECT COUNT(*)
			FROM transaction_request
			WHERE perusahaan_id = ?
		) AS total_request, a.id_asset, a.nama, a.alamat, a.status_asset
		FROM asset a
		WHERE a.perusahaan_id = ?
	`
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	nId, _ := strconv.Atoi(perusahaan_id)
	rows, err := stmt.Query(nId, nId)
	if err != nil {
		res.Status = 401
		res.Message = "rows gagal"
		res.Data = err.Error()
		return res, err
	}
	defer rows.Close()

	for rows.Next() {
		var asset Asset
		err = rows.Scan(&homeuser.TotalRequest, &asset.Id_asset, &asset.Nama, &asset.Alamat, &asset.Status_asset)
		if err != nil {
			res.Status = 401
			res.Message = "Failed to scan row"
			res.Data = err.Error()
			return res, err
		}
		fmt.Println("1")
		homeuser.SemuaAsset = append(homeuser.SemuaAsset, asset)
	}

	res.Status = http.StatusOK
	res.Message = "Berhasil mengambil data"
	res.Data = homeuser

	defer db.DbClose(con)
	return res, nil
}
