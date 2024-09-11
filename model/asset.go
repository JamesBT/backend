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
	"strings"
	"time"
)

// CRUD aset ============================================================================
func CreateAsset(filelegalitas *multipart.FileHeader, suratkuasa *multipart.FileHeader,
	gambar_asset *multipart.FileHeader, nama, surat_legalitas, tipe, usage, tag, nomorlegalitas, status, alamat, kondisi, koordinat, batas_koordinat, luas, nilai, provinsi string) (Response, error) {

	var res Response

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	var provinsiExists bool
	provQuery := "SELECT EXISTS(SELECT 1 FROM provinsi WHERE id_provinsi = ?)"
	err = con.QueryRow(provQuery, provinsi).Scan(&provinsiExists)
	if err != nil || !provinsiExists {
		res.Status = 401
		res.Message = "Provinsi tidak valid"
		res.Data = "Provinsi ID tidak ditemukan"
		return res, err
	}

	usageIds := strings.Split(usage, ",")
	for _, id := range usageIds {
		var usageExists bool
		fmt.Println("usage", id)
		usageQuery := "SELECT EXISTS(SELECT 1 FROM penggunaan WHERE id = ?)"
		err = con.QueryRow(usageQuery, id).Scan(&usageExists)
		if err != nil || !usageExists {
			res.Status = 401
			res.Message = "Penggunaan tidak valid"
			res.Data = "Penggunaan ID " + id + " tidak ditemukan"
			return res, err
		}
	}

	tagIds := strings.Split(tag, ",")
	for _, id2 := range tagIds {
		var tagExists bool
		fmt.Println("tag", id2)
		tagQuery := "SELECT EXISTS(SELECT 1 FROM tags WHERE id = ?)"
		err = con.QueryRow(tagQuery, id2).Scan(&tagExists)
		if err != nil || !tagExists {
			res.Status = 401
			res.Message = "Tag tidak valid"
			res.Data = "Tag ID " + id2 + " tidak ditemukan"
			return res, err
		}
	}

	// query := `
	// INSERT INTO asset (perusahaan_id, nama, tipe, nomor_legalitas, status_asset, surat_legalitas, alamat,
	// kondisi, titik_koordinat, batas_koordinat, luas, nilai, provinsi, created_at)
	// VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,NOW())
	// `
	query := `
	INSERT INTO asset (nama, tipe, nomor_legalitas, status_asset, surat_legalitas, alamat, 
	kondisi, titik_koordinat, batas_koordinat, luas, nilai, provinsi, created_at) 
	VALUES (?,?,?,?,?,?,?,?,?,?,?,?,NOW())
	`
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	// result, err := stmt.Exec(
	// 	perusahaan_id, nama, tipe, nomorlegalitas, status, surat_legalitas, alamat, kondisi, koordinat, batas_koordinat,
	// 	luas, nilai, provinsi)
	result, err := stmt.Exec(
		nama, tipe, nomorlegalitas, status, surat_legalitas, alamat, kondisi, koordinat, batas_koordinat,
		luas, nilai, provinsi)
	if err != nil {
		res.Status = 401
		res.Message = "exec gagal"
		res.Data = err.Error()
		return res, err
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		log.Println(err.Error())
		return res, err
	}

	// tambah usage + tags
	for _, usageId := range usageIds {
		usageQuery := "INSERT INTO asset_penggunaan (id_asset, id_penggunaan) VALUES (?, ?)"
		_, err = con.Exec(usageQuery, lastId, usageId)
		if err != nil {
			res.Status = 401
			res.Message = "Gagal menambah penggunaan"
			res.Data = err.Error()
			return res, err
		}
	}

	// Insert into asset_tags (id_asset, id_tags)
	for _, tagId := range tagIds {
		tagQuery := "INSERT INTO asset_tags (id_asset, id_tags) VALUES (?, ?)"
		_, err = con.Exec(tagQuery, lastId, tagId)
		if err != nil {
			res.Status = 401
			res.Message = "Gagal menambah tag"
			res.Data = err.Error()
			return res, err
		}
	}

	// tambah filelegalitas
	//source
	srclegalitas, err := filelegalitas.Open()
	if err != nil {
		log.Println(err.Error())
		return res, err
	}
	defer srclegalitas.Close()

	tempid := int(lastId)
	_tempid := strconv.Itoa(tempid)
	// Destination
	filelegalitas.Filename = _tempid + "_" + filelegalitas.Filename
	pathFileLegalitas := "uploads/asset/file_legalitas/" + filelegalitas.Filename
	dstlegalitas, err := os.Create("uploads/asset/file_legalitas/" + filelegalitas.Filename)
	if err != nil {
		log.Println("2")
		fmt.Print("2")
		return res, err
	}

	// Copy
	if _, err = io.Copy(dstlegalitas, srclegalitas); err != nil {
		log.Println(err.Error())
		log.Println("3")
		fmt.Print("3")
		return res, err
	}
	dstlegalitas.Close()

	err = UpdateDataFotoPath("asset", "file_legalitas", pathFileLegalitas, "id_asset", int(lastId))
	if err != nil {
		return res, err
	}

	// tambah suratkuasa
	//source
	srcsuratkuasa, err := suratkuasa.Open()
	if err != nil {
		log.Println(err.Error())
		return res, err
	}
	defer srcsuratkuasa.Close()

	// Destination
	suratkuasa.Filename = _tempid + "_" + "_" + suratkuasa.Filename
	pathFileSuratKuasa := "uploads/asset/surat_kuasa/" + suratkuasa.Filename
	dstsuratkuasa, err := os.Create("uploads/asset/surat_kuasa/" + suratkuasa.Filename)
	if err != nil {
		log.Println("2")
		fmt.Print("2")
		return res, err
	}

	// Copy
	if _, err = io.Copy(dstsuratkuasa, srcsuratkuasa); err != nil {
		log.Println(err.Error())
		log.Println("3")
		fmt.Print("3")
		return res, err
	}
	dstsuratkuasa.Close()

	err = UpdateDataFotoPath("asset", "surat_kuasa", pathFileSuratKuasa, "id_asset", int(lastId))
	if err != nil {
		return res, err
	}

	// gambar
	srcgambar, err := gambar_asset.Open()
	if err != nil {
		log.Println(err.Error())
		return res, err
	}
	defer srcgambar.Close()

	// Destination
	gambar_asset.Filename = _tempid + "_" + "_" + gambar_asset.Filename
	pathFileGambar := "uploads/asset/foto/" + gambar_asset.Filename
	dstgambar, err := os.Create("uploads/asset/foto/" + gambar_asset.Filename)
	if err != nil {
		log.Println("2")
		fmt.Print("2")
		return res, err
	}

	// Copy
	if _, err = io.Copy(dstgambar, srcgambar); err != nil {
		log.Println(err.Error())
		log.Println("3")
		fmt.Print("3")
		return res, err
	}
	dstgambar.Close()

	queryGambar := `
	INSERT INTO asset_gambar (id_asset_gambar, link_gambar) VALUES (?,?)
	`
	stmtGambar, err := con.Prepare(queryGambar)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmtGambar.Close()

	_, err = stmtGambar.Exec(int(lastId), pathFileGambar)
	if err != nil {
		res.Status = 401
		res.Message = "exec gagal"
		res.Data = err.Error()
		return res, err
	}

	var tempaset Response
	tempaset, _ = GetAssetById(strconv.Itoa(tempid))

	res.Status = http.StatusOK
	res.Message = "Berhasil memasukkan data"
	res.Data = tempaset.Data

	defer db.DbClose(con)
	return res, nil
}

func CreateAssetChild(
	filelegalitas *multipart.FileHeader, suratkuasa *multipart.FileHeader, gambar_asset *multipart.FileHeader,
	parent_id, nama, surat_legalitas, tipe, usage, tag, nomor_legalitas, status,
	alamat, kondisi, koordinat, batas_koordinat, luas, nilai string) (Response, error) {
	var res Response

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	usageIds := strings.Split(usage, ",")
	for _, id := range usageIds {
		var usageExists bool
		fmt.Println("usage", id)
		usageQuery := "SELECT EXISTS(SELECT 1 FROM penggunaan WHERE id = ?)"
		err = con.QueryRow(usageQuery, id).Scan(&usageExists)
		if err != nil || !usageExists {
			res.Status = 401
			res.Message = "Penggunaan tidak valid"
			res.Data = "Penggunaan ID " + id + " tidak ditemukan"
			return res, err
		}
	}

	tagIds := strings.Split(tag, ",")
	for _, id2 := range tagIds {
		var tagExists bool
		fmt.Println("tag", id2)
		tagQuery := "SELECT EXISTS(SELECT 1 FROM tags WHERE id = ?)"
		err = con.QueryRow(tagQuery, id2).Scan(&tagExists)
		if err != nil || !tagExists {
			res.Status = 401
			res.Message = "Tag tidak valid"
			res.Data = "Tag ID " + id2 + " tidak ditemukan"
			return res, err
		}
	}

	// ambil parent aset + dt provinsi
	var dtProvinsi string
	var dtAsetChild string
	ParentQuery := "SELECT provinsi,id_asset_child FROM asset WHERE id_asset = ?"
	err = con.QueryRow(ParentQuery, parent_id).Scan(&dtProvinsi, &dtAsetChild)
	if err != nil || (dtProvinsi == "") {
		res.Status = 401
		res.Message = "provinsi tidak valid"
		res.Data = "Aset ID " + parent_id + " tidak ditemukan"
		return res, err
	}

	// query := `
	// INSERT INTO asset (id_asset_parent,perusahaan_id, nama, tipe, nomor_legalitas, status_asset, surat_legalitas, alamat,
	// kondisi, titik_koordinat, batas_koordinat, luas, nilai, provinsi, created_at)
	// VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,NOW())
	// `
	query := `
	INSERT INTO asset (id_asset_parent, nama, tipe, nomor_legalitas, status_asset, surat_legalitas, alamat, 
	kondisi, titik_koordinat, batas_koordinat, luas, nilai, provinsi, created_at) 
	VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,NOW())
	`
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	// result, err := stmt.Exec(
	// 	parent_id, perusahaan_id, nama, tipe, nomor_legalitas, status, surat_legalitas, alamat, kondisi, koordinat, batas_koordinat,
	// 	luas, nilai, provinsi)
	result, err := stmt.Exec(
		parent_id, nama, tipe, nomor_legalitas, status, surat_legalitas, alamat, kondisi, koordinat, batas_koordinat,
		luas, nilai, dtProvinsi)
	if err != nil {
		res.Status = 401
		res.Message = "exec gagal"
		res.Data = err.Error()
		return res, err
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		log.Println(err.Error())
		return res, err
	}

	// tambah usage + tags
	for _, usageId := range usageIds {
		usageQuery := "INSERT INTO asset_penggunaan (id_asset, id_penggunaan) VALUES (?, ?)"
		_, err = con.Exec(usageQuery, lastId, usageId)
		if err != nil {
			res.Status = 401
			res.Message = "Gagal menambah penggunaan"
			res.Data = err.Error()
			return res, err
		}
	}

	for _, tagId := range tagIds {
		tagQuery := "INSERT INTO asset_tags (id_asset, id_tags) VALUES (?, ?)"
		_, err = con.Exec(tagQuery, lastId, tagId)
		if err != nil {
			res.Status = 401
			res.Message = "Gagal menambah tag"
			res.Data = err.Error()
			return res, err
		}
	}

	// tambah filelegalitas
	//source
	srclegalitas, err := filelegalitas.Open()
	if err != nil {
		log.Println(err.Error())
		return res, err
	}
	defer srclegalitas.Close()

	tempid := int(lastId)
	_tempid := strconv.Itoa(tempid)
	// Destination
	filelegalitas.Filename = _tempid + "_" + filelegalitas.Filename
	pathFileLegalitas := "uploads/asset/file_legalitas/" + filelegalitas.Filename
	dstlegalitas, err := os.Create("uploads/asset/file_legalitas/" + filelegalitas.Filename)
	if err != nil {
		log.Println("2")
		fmt.Print("2")
		return res, err
	}

	// Copy
	if _, err = io.Copy(dstlegalitas, srclegalitas); err != nil {
		log.Println(err.Error())
		log.Println("3")
		fmt.Print("3")
		return res, err
	}
	dstlegalitas.Close()

	err = UpdateDataFotoPath("asset", "file_legalitas", pathFileLegalitas, "id_asset", int(lastId))
	if err != nil {
		return res, err
	}

	// tambah suratkuasa
	//source
	srcsuratkuasa, err := suratkuasa.Open()
	if err != nil {
		log.Println(err.Error())
		return res, err
	}
	defer srcsuratkuasa.Close()

	// Destination
	suratkuasa.Filename = _tempid + "_" + suratkuasa.Filename
	pathFileSuratKuasa := "uploads/asset/surat_kuasa/" + suratkuasa.Filename
	dstsuratkuasa, err := os.Create("uploads/asset/surat_kuasa/" + suratkuasa.Filename)
	if err != nil {
		log.Println("2")
		fmt.Print("2")
		return res, err
	}

	// Copy
	if _, err = io.Copy(dstsuratkuasa, srcsuratkuasa); err != nil {
		log.Println(err.Error())
		log.Println("3")
		fmt.Print("3")
		return res, err
	}
	dstsuratkuasa.Close()

	err = UpdateDataFotoPath("asset", "surat_kuasa", pathFileSuratKuasa, "id_asset", int(lastId))
	if err != nil {
		return res, err
	}

	// gambar
	srcgambar, err := gambar_asset.Open()
	if err != nil {
		log.Println(err.Error())
		return res, err
	}
	defer srcgambar.Close()

	// Destination
	gambar_asset.Filename = _tempid + "_" + gambar_asset.Filename
	pathFileGambar := "uploads/asset/foto/" + gambar_asset.Filename
	dstgambar, err := os.Create("uploads/asset/foto/" + gambar_asset.Filename)
	if err != nil {
		log.Println("2")
		fmt.Print("2")
		return res, err
	}

	// Copy
	if _, err = io.Copy(dstgambar, srcgambar); err != nil {
		log.Println(err.Error())
		log.Println("3")
		fmt.Print("3")
		return res, err
	}
	dstgambar.Close()

	queryGambar := `
	INSERT INTO asset_gambar (id_asset_gambar, link_gambar) VALUES (?,?)
	`
	stmtGambar, err := con.Prepare(queryGambar)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmtGambar.Close()

	_, err = stmtGambar.Exec(int(lastId), pathFileGambar)
	if err != nil {
		res.Status = 401
		res.Message = "exec gagal"
		res.Data = err.Error()
		return res, err
	}

	templastid := int(lastId)
	// update parent jadi punya child
	if dtAsetChild != "" {
		dtAsetChild = dtAsetChild + "," + strconv.Itoa(templastid)
	} else {
		dtAsetChild = strconv.Itoa(templastid)
	}
	queryupdate := "UPDATE asset SET id_asset_child = ? WHERE id_asset = ?"
	stmtupdate, err := con.Prepare(queryupdate)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmtupdate.Close()

	_, err = stmtupdate.Exec(dtAsetChild, parent_id)
	if err != nil {
		res.Status = 401
		res.Message = "exec gagal"
		res.Data = err.Error()
		return res, err
	}

	tempaset, err := GetAssetById(_tempid)
	if err != nil {
		res.Status = 401
		res.Message = "get aset by id gagal"
		res.Data = err.Error()
		return res, err
	}

	res.Status = http.StatusOK
	res.Message = "Berhasil memasukkan data"
	res.Data = tempaset

	defer db.DbClose(con)
	return res, nil
}

func GetAllAsset() (Response, error) {
	var res Response
	var arrAset = []Asset{}
	var dtAset Asset

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	query := "SELECT * FROM asset"
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
	var masaSewa []byte
	var deleteAt []byte
	var idJoin, idAssetChild sql.NullString
	var idAssetParent, idPerusahaan, idProvinsi sql.NullInt32

	for result.Next() {
		err = result.Scan(&dtAset.Id_asset, &idAssetParent, &idAssetChild, &idJoin, &idPerusahaan, &dtAset.Nama, &dtAset.Tipe, &dtAset.Nomor_legalitas, &dtAset.File_legalitas, &dtAset.Status_asset, &dtAset.Surat_kuasa, &dtAset.Surat_legalitas, &dtAset.Alamat, &dtAset.Kondisi, &dtAset.Titik_koordinat, &dtAset.Batas_koordinat, &dtAset.Luas, &dtAset.Nilai, &idProvinsi, &dtAset.Status_pengecekan, &dtAset.Status_verifikasi, &dtAset.Status_publik, &dtAset.Hak_akses, &masaSewa, &dtAset.Created_at, &deleteAt)
		if err != nil {
			res.Status = 401
			res.Message = "rows scan"
			res.Data = err.Error()
			return res, err
		}

		if masaSewa != nil {
			masaSewaWaktu, masaSewaErr := time.Parse("2006-01-02 15:04:05", string(deleteAt))
			if masaSewaErr != nil {
				dtAset.Deleted_at = ""
			} else {
				dtAset.Deleted_at = masaSewaWaktu.Format("2006-01-02 15:04:05")
			}
		} else {
			dtAset.Deleted_at = ""
		}

		if deleteAt != nil {
			parsedTime, parseErr := time.Parse("2006-01-02 15:04:05", string(deleteAt))
			if parseErr != nil {
				dtAset.Deleted_at = ""
			} else {
				dtAset.Deleted_at = parsedTime.Format("2006-01-02 15:04:05")
			}
		} else {
			dtAset.Deleted_at = ""
		}
		if idAssetParent.Valid {
			dtAset.Id_asset_parent = int(idAssetParent.Int32)
		} else {
			dtAset.Id_asset_parent = 0
		}
		if idAssetChild.Valid {
			dtAset.Id_asset_child = idAssetChild.String
		} else {
			dtAset.Id_asset_child = ""
		}
		if idJoin.Valid {
			dtAset.Id_join = idJoin.String
		} else {
			dtAset.Id_join = "0"
		}
		if idPerusahaan.Valid {
			dtAset.Id_perusahaan = int(idPerusahaan.Int32)
		} else {
			dtAset.Id_perusahaan = 0
		}
		if idProvinsi.Valid {
			dtAset.Provinsi = int(idProvinsi.Int32)
		} else {
			dtAset.Provinsi = 0
		}
		arrAset = append(arrAset, dtAset)
	}

	res.Status = http.StatusOK
	res.Message = "Berhasil mengambil data"
	res.Data = arrAset

	defer db.DbClose(con)
	return res, nil
}

func GetAssetById(aset_id string) (Response, error) {
	var res Response
	var dtAset Asset

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	query := "SELECT * FROM asset WHERE id_asset = ?"
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	nId, _ := strconv.Atoi(aset_id)
	var masaSewa []byte
	var deleteAt []byte
	var idJoin, idAssetChild sql.NullString
	var idAssetParent, idPerusahaan, idProvinsi sql.NullInt32
	err = stmt.QueryRow(nId).Scan(&dtAset.Id_asset, &idAssetParent, &idAssetChild, &idJoin, &idPerusahaan, &dtAset.Nama, &dtAset.Tipe, &dtAset.Nomor_legalitas, &dtAset.File_legalitas, &dtAset.Status_asset, &dtAset.Surat_kuasa, &dtAset.Surat_legalitas, &dtAset.Alamat, &dtAset.Kondisi, &dtAset.Titik_koordinat, &dtAset.Batas_koordinat, &dtAset.Luas, &dtAset.Nilai, &idProvinsi, &dtAset.Status_pengecekan, &dtAset.Status_verifikasi, &dtAset.Status_publik, &dtAset.Hak_akses, &masaSewa, &dtAset.Created_at, &deleteAt)
	if err != nil {
		res.Status = 401
		res.Message = "exec gagal"
		res.Data = err.Error()
		return res, err
	}
	fmt.Println("ambil berhasil")

	if masaSewa != nil {
		masaSewaWaktu, masaSewaErr := time.Parse("2006-01-02 15:04:05", string(deleteAt))
		if masaSewaErr != nil {
			dtAset.Deleted_at = ""
		} else {
			dtAset.Deleted_at = masaSewaWaktu.Format("2006-01-02 15:04:05")
		}
	} else {
		dtAset.Deleted_at = ""
	}

	if deleteAt != nil {
		parsedTime, parseErr := time.Parse("2006-01-02 15:04:05", string(deleteAt))
		if parseErr != nil {
			dtAset.Deleted_at = ""
		} else {
			dtAset.Deleted_at = parsedTime.Format("2006-01-02 15:04:05")
		}
	} else {
		dtAset.Deleted_at = ""
	}
	if idAssetParent.Valid {
		dtAset.Id_asset_parent = int(idAssetParent.Int32)
	} else {
		dtAset.Id_asset_parent = 0
	}
	if idAssetChild.Valid {
		dtAset.Id_asset_child = idAssetChild.String
	} else {
		dtAset.Id_asset_child = ""
	}
	if idJoin.Valid {
		dtAset.Id_join = idJoin.String
	} else {
		dtAset.Id_join = "0"
	}
	if idPerusahaan.Valid {
		dtAset.Id_perusahaan = int(idPerusahaan.Int32)
	} else {
		dtAset.Id_perusahaan = 0
	}
	if idProvinsi.Valid {
		dtAset.Provinsi = int(idProvinsi.Int32)
	} else {
		dtAset.Provinsi = 0
	}

	gambarQuery := "SELECT link_gambar FROM asset_gambar WHERE id_asset_gambar = ?"
	rows, err := con.Query(gambarQuery, nId)
	if err != nil {
		res.Status = 401
		res.Message = "gagal mengambil gambar"
		res.Data = err.Error()
		return res, err
	}
	defer rows.Close()

	var gambarLinks []string
	for rows.Next() {
		var link string
		err := rows.Scan(&link)
		if err != nil {
			res.Status = 401
			res.Message = "gagal membaca gambar"
			res.Data = err.Error()
			return res, err
		}
		gambarLinks = append(gambarLinks, link)
	}
	dtAset.LinkGambar = gambarLinks

	usageQuery := "SELECT p.nama FROM asset_penggunaan ap JOIN penggunaan p ON ap.id_penggunaan = p.id WHERE ap.id_asset = ?"
	rowsusage, err := con.Query(usageQuery, nId)
	if err != nil {
		res.Status = 401
		res.Message = "gagal mengambil gambar"
		res.Data = err.Error()
		return res, err
	}
	defer rowsusage.Close()

	var usage []string
	for rowsusage.Next() {
		var link string
		err := rowsusage.Scan(&link)
		if err != nil {
			res.Status = 401
			res.Message = "gagal membaca gambar"
			res.Data = err.Error()
			return res, err
		}
		usage = append(usage, link)
	}
	dtAset.LinkGambar = usage

	tagsQuery := "SELECT t.nama FROM asset_tags at JOIN tags t ON at.id_tags = t.id WHERE at.id_asset = ?"
	tagsRows, err := con.Query(tagsQuery, nId)
	if err != nil {
		res.Status = 401
		res.Message = "gagal mengambil tags"
		res.Data = err.Error()
		return res, err
	}
	defer tagsRows.Close()

	var tags []string
	for tagsRows.Next() {
		var tag string
		err := tagsRows.Scan(&tag)
		if err != nil {
			res.Status = 401
			res.Message = "gagal membaca tags"
			res.Data = err.Error()
			return res, err
		}
		tags = append(tags, tag)
	}
	dtAset.TagsAssets = tags

	res.Status = http.StatusOK
	res.Message = "Berhasil mengambil data"
	res.Data = dtAset

	defer db.DbClose(con)
	return res, nil
}

func GetAssetChildByParentId(aset_id string) (Response, error) {
	var res Response

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	query := "SELECT id_asset_child FROM asset WHERE id_asset = ?"
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	nId, _ := strconv.Atoi(aset_id)
	var idAssetChild sql.NullString
	var idChild string
	err = stmt.QueryRow(nId).Scan(&idAssetChild)
	if err != nil {
		res.Status = 401
		res.Message = "exec gagal"
		res.Data = err.Error()
		return res, err
	}
	fmt.Println("ambil berhasil")

	var arrAset []Asset
	if idAssetChild.Valid {
		idChild = idAssetChild.String
	} else {
		idChild = ""
	}

	// ambil aset pisah berdasarkan ,

	if idChild != "" {
		childIds := strings.Split(idChild, ",")
		for _, id := range childIds {
			var dtAset Asset
			usageQuery := "SELECT * FROM asset WHERE id_asset = ?"
			var masaSewa []byte
			var deleteAt []byte
			var idJoin, idAssetChild sql.NullString
			var idAssetParent, idPerusahaan, idProvinsi sql.NullInt32
			err = con.QueryRow(usageQuery, id).Scan(&dtAset.Id_asset, &idAssetParent, &idAssetChild, &idJoin, &idPerusahaan, &dtAset.Nama, &dtAset.Tipe, &dtAset.Nomor_legalitas, &dtAset.File_legalitas, &dtAset.Status_asset, &dtAset.Surat_kuasa, &dtAset.Surat_legalitas, &dtAset.Alamat, &dtAset.Kondisi, &dtAset.Titik_koordinat, &dtAset.Batas_koordinat, &dtAset.Luas, &dtAset.Nilai, &idProvinsi, &dtAset.Status_pengecekan, &dtAset.Status_verifikasi, &dtAset.Status_publik, &dtAset.Hak_akses, &masaSewa, &dtAset.Created_at, &deleteAt)
			if err != nil {
				res.Status = 401
				res.Message = "exec error "
				res.Data = err
				return res, err
			}
			if masaSewa != nil {
				masaSewaWaktu, masaSewaErr := time.Parse("2006-01-02 15:04:05", string(deleteAt))
				if masaSewaErr != nil {
					dtAset.Deleted_at = ""
				} else {
					dtAset.Deleted_at = masaSewaWaktu.Format("2006-01-02 15:04:05")
				}
			} else {
				dtAset.Deleted_at = ""
			}
			if deleteAt != nil {
				parsedTime, parseErr := time.Parse("2006-01-02 15:04:05", string(deleteAt))
				if parseErr != nil {
					dtAset.Deleted_at = ""
				} else {
					dtAset.Deleted_at = parsedTime.Format("2006-01-02 15:04:05")
				}
			} else {
				dtAset.Deleted_at = ""
			}
			if idAssetParent.Valid {
				dtAset.Id_asset_parent = int(idAssetParent.Int32)
			} else {
				dtAset.Id_asset_parent = 0
			}
			if idAssetChild.Valid {
				dtAset.Id_asset_child = idAssetChild.String
			} else {
				dtAset.Id_asset_child = ""
			}
			if idJoin.Valid {
				dtAset.Id_join = idJoin.String
			} else {
				dtAset.Id_join = "0"
			}
			if idPerusahaan.Valid {
				dtAset.Id_perusahaan = int(idPerusahaan.Int32)
			} else {
				dtAset.Id_perusahaan = 0
			}
			if idProvinsi.Valid {
				dtAset.Provinsi = int(idProvinsi.Int32)
			} else {
				dtAset.Provinsi = 0
			}
			arrAset = append(arrAset, dtAset)
		}
	} else {
		res.Status = 401
		res.Message = "tidak ada child dari aset parent id: " + aset_id
		return res, errors.New(res.Message)
	}

	res.Status = http.StatusOK
	res.Message = "Berhasil mengambil data"
	res.Data = arrAset

	defer db.DbClose(con)
	return res, nil
}

func GetAssetDetailedById(aset_id string) (Response, error) {
	var res Response

	fmt.Println("get aset detailed by id")
	// ambil data parent
	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	dtAset, err := fetchAssetDetailed(con, aset_id)
	if err != nil {
		res.Status = 401
		res.Message = "Failed to fetch asset details"
		res.Data = err.Error()
		return res, err
	}

	res.Status = http.StatusOK
	res.Message = "Berhasil mengambil data"
	res.Data = dtAset

	defer db.DbClose(con)
	return res, nil
}

func fetchAssetDetailed(con *sql.DB, aset_id string) (Asset, error) {
	var dtAset Asset

	query := "SELECT * FROM asset WHERE id_asset = ?"
	stmt, err := con.Prepare(query)
	if err != nil {
		return dtAset, err
	}
	defer stmt.Close()

	nId, _ := strconv.Atoi(aset_id)
	var masaSewa []byte
	var deleteAt []byte
	var idJoin, idAssetChild sql.NullString
	var idAssetParent, idPerusahaan, idProvinsi sql.NullInt32
	err = stmt.QueryRow(nId).Scan(&dtAset.Id_asset, &idAssetParent, &idAssetChild, &idJoin, &idPerusahaan, &dtAset.Nama, &dtAset.Tipe, &dtAset.Nomor_legalitas, &dtAset.File_legalitas, &dtAset.Status_asset, &dtAset.Surat_kuasa, &dtAset.Surat_legalitas, &dtAset.Alamat, &dtAset.Kondisi, &dtAset.Titik_koordinat, &dtAset.Batas_koordinat, &dtAset.Luas, &dtAset.Nilai, &idProvinsi, &dtAset.Status_pengecekan, &dtAset.Status_verifikasi, &dtAset.Status_publik, &dtAset.Hak_akses, &masaSewa, &dtAset.Created_at, &deleteAt)
	if err != nil {
		return dtAset, err
	}

	if masaSewa != nil {
		masaSewaWaktu, masaSewaErr := time.Parse("2006-01-02 15:04:05", string(deleteAt))
		if masaSewaErr != nil {
			dtAset.Deleted_at = ""
		} else {
			dtAset.Deleted_at = masaSewaWaktu.Format("2006-01-02 15:04:05")
		}
	} else {
		dtAset.Deleted_at = ""
	}

	if deleteAt != nil {
		parsedTime, parseErr := time.Parse("2006-01-02 15:04:05", string(deleteAt))
		if parseErr != nil {
			dtAset.Deleted_at = ""
		} else {
			dtAset.Deleted_at = parsedTime.Format("2006-01-02 15:04:05")
		}
	} else {
		dtAset.Deleted_at = ""
	}
	if idAssetParent.Valid {
		dtAset.Id_asset_parent = int(idAssetParent.Int32)
	} else {
		dtAset.Id_asset_parent = 0
	}
	if idAssetChild.Valid {
		dtAset.Id_asset_child = idAssetChild.String
	} else {
		dtAset.Id_asset_child = ""
	}
	if idJoin.Valid {
		dtAset.Id_join = idJoin.String
	} else {
		dtAset.Id_join = "0"
	}
	if idPerusahaan.Valid {
		dtAset.Id_perusahaan = int(idPerusahaan.Int32)
	} else {
		dtAset.Id_perusahaan = 0
	}
	if idProvinsi.Valid {
		dtAset.Provinsi = int(idProvinsi.Int32)
	} else {
		dtAset.Provinsi = 0
	}

	gambarQuery := "SELECT link_gambar FROM asset_gambar WHERE id_asset_gambar = ?"
	rows, err := con.Query(gambarQuery, nId)
	if err != nil {
		return dtAset, err
	}
	defer rows.Close()

	var gambarLinks []string
	for rows.Next() {
		var link string
		err := rows.Scan(&link)
		if err != nil {
			return dtAset, err
		}
		gambarLinks = append(gambarLinks, link)
	}
	dtAset.LinkGambar = gambarLinks

	usageQuery := "SELECT p.nama FROM asset_penggunaan ap JOIN penggunaan p ON ap.id_penggunaan = p.id WHERE ap.id_asset = ?"
	rowsusage, err := con.Query(usageQuery, nId)
	if err != nil {
		return dtAset, err
	}
	defer rowsusage.Close()

	var usage []string
	for rowsusage.Next() {
		var link string
		err := rowsusage.Scan(&link)
		if err != nil {
			return dtAset, err
		}
		usage = append(usage, link)
	}
	dtAset.Usage = usage

	tagsQuery := "SELECT t.nama FROM asset_tags at JOIN tags t ON at.id_tags = t.id WHERE at.id_asset = ?"
	tagsRows, err := con.Query(tagsQuery, nId)
	if err != nil {
		return dtAset, err
	}
	defer tagsRows.Close()

	var tags []string
	for tagsRows.Next() {
		var tag string
		err := tagsRows.Scan(&tag)
		if err != nil {
			return dtAset, err
		}
		tags = append(tags, tag)
	}
	dtAset.TagsAssets = tags

	childQuery := "SELECT id_asset FROM asset WHERE id_asset_parent = ?"
	rows, err = con.Query(childQuery, dtAset.Id_asset)
	if err != nil {
		return dtAset, err
	}
	defer rows.Close()

	for rows.Next() {
		var childId int
		if err := rows.Scan(&childId); err != nil {
			return dtAset, err
		}

		childAset, err := fetchAssetDetailed(con, strconv.Itoa(childId))
		if err != nil {
			return dtAset, err
		}
		dtAset.ChildAssets = append(dtAset.ChildAssets, childAset)
	}

	return dtAset, nil
}

func UbahVisibilitasAset(aset_id, input string) (Response, error) {
	var res Response

	type temp_visibilitas_asset_acc struct {
		Visibilitas string `json:"visibilitas"`
	}

	var visibilitasAsset temp_visibilitas_asset_acc
	err := json.Unmarshal([]byte(input), &visibilitasAsset)
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

	query := "UPDATE asset SET `status_publik`= ? WHERE `id_asset` = ?"
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	nId, _ := strconv.Atoi(aset_id)
	_, err = stmt.Exec(visibilitasAsset.Visibilitas, nId)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}

	tempaset, _ := GetAssetById(aset_id)

	res.Status = http.StatusOK
	res.Message = "Berhasil mengupdate data"
	res.Data = tempaset

	defer db.DbClose(con)
	return res, nil
}

func GetAssetDetailedByPerusahaanId(perusahaan_id string) (Response, error) {
	var res Response

	fmt.Println("get aset detailed by id")
	// ambil data parent
	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	dtAset, err := fetchAssetsByPerusahaanId(con, perusahaan_id)
	if err != nil {
		res.Status = 401
		res.Message = "Failed to fetch asset details"
		res.Data = err.Error()
		return res, err
	}

	res.Status = http.StatusOK
	res.Message = "Berhasil mengambil data"
	res.Data = dtAset

	defer db.DbClose(con)
	return res, nil
}

func fetchAssetsByPerusahaanId(con *sql.DB, perusahaan_id string) ([]Asset, error) {
	var assets []Asset

	query := "SELECT * FROM asset WHERE perusahaan_id = ?"
	rows, err := con.Query(query, perusahaan_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var dtAset Asset
		var masaSewa []byte
		var deleteAt []byte
		var idJoin, idAssetChild sql.NullString
		var idAssetParent, idPerusahaan, idProvinsi sql.NullInt32

		err := rows.Scan(&dtAset.Id_asset, &idAssetParent, &idAssetChild, &idJoin, &idPerusahaan, &dtAset.Nama,
			&dtAset.Tipe, &dtAset.Nomor_legalitas, &dtAset.File_legalitas, &dtAset.Status_asset, &dtAset.Surat_kuasa,
			&dtAset.Surat_legalitas, &dtAset.Alamat, &dtAset.Kondisi, &dtAset.Titik_koordinat, &dtAset.Batas_koordinat,
			&dtAset.Luas, &dtAset.Nilai, &idProvinsi, &dtAset.Status_pengecekan, &dtAset.Status_verifikasi,
			&dtAset.Status_publik, &dtAset.Hak_akses, &masaSewa, &dtAset.Created_at, &deleteAt)
		if err != nil {
			return nil, err
		}

		if idAssetParent.Valid {
			dtAset.Id_asset_parent = int(idAssetParent.Int32)
		} else {
			dtAset.Id_asset_parent = 0
		}
		if idAssetChild.Valid {
			dtAset.Id_asset_child = idAssetChild.String
		} else {
			dtAset.Id_asset_child = ""
		}
		if idJoin.Valid {
			dtAset.Id_join = idJoin.String
		} else {
			dtAset.Id_join = "0"
		}
		if idPerusahaan.Valid {
			dtAset.Id_perusahaan = int(idPerusahaan.Int32)
		} else {
			dtAset.Id_perusahaan = 0
		}
		if idProvinsi.Valid {
			dtAset.Provinsi = int(idProvinsi.Int32)
		} else {
			dtAset.Provinsi = 0
		}

		// Fetch gambar
		gambarQuery := "SELECT link_gambar FROM asset_gambar WHERE id_asset_gambar = ?"
		gambarRows, err := con.Query(gambarQuery, dtAset.Id_asset)
		if err != nil {
			return nil, err
		}
		defer gambarRows.Close()

		var gambarLinks []string
		for gambarRows.Next() {
			var link string
			err := gambarRows.Scan(&link)
			if err != nil {
				return nil, err
			}
			gambarLinks = append(gambarLinks, link)
		}
		dtAset.LinkGambar = gambarLinks

		// Fetch usage
		usageQuery := "SELECT p.nama FROM asset_penggunaan ap JOIN penggunaan p ON ap.id_penggunaan = p.id WHERE ap.id_asset = ?"
		usageRows, err := con.Query(usageQuery, dtAset.Id_asset)
		if err != nil {
			return nil, err
		}
		defer usageRows.Close()

		var usage []string
		for usageRows.Next() {
			var name string
			err := usageRows.Scan(&name)
			if err != nil {
				return nil, err
			}
			usage = append(usage, name)
		}
		dtAset.Usage = usage

		// Fetch tags
		tagsQuery := "SELECT t.nama FROM asset_tags at JOIN tags t ON at.id_tags = t.id WHERE at.id_asset = ?"
		tagsRows, err := con.Query(tagsQuery, dtAset.Id_asset)
		if err != nil {
			return nil, err
		}
		defer tagsRows.Close()

		var tags []string
		for tagsRows.Next() {
			var tag string
			err := tagsRows.Scan(&tag)
			if err != nil {
				return nil, err
			}
			tags = append(tags, tag)
		}
		dtAset.TagsAssets = tags

		// Fetch child assets
		childQuery := "SELECT id_asset FROM asset WHERE id_asset_parent = ?"
		childRows, err := con.Query(childQuery, dtAset.Id_asset)
		if err != nil {
			return nil, err
		}
		defer childRows.Close()

		var childAssets []Asset
		for childRows.Next() {
			var childId int
			err := childRows.Scan(&childId)
			if err != nil {
				return nil, err
			}
			childAset, err := fetchAssetDetailed(con, strconv.Itoa(childId))
			if err != nil {
				return nil, err
			}
			childAssets = append(childAssets, childAset)
		}
		dtAset.ChildAssets = childAssets

		assets = append(assets, dtAset)
	}

	return assets, nil
}

func GetAssetDetailedByUserId(perusahaan_id string) (Response, error) {
	var res Response

	fmt.Println("get aset detailed by id")
	// ambil data parent
	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	dtAset, err := fetchAssetDetailedByUserId(con, perusahaan_id)
	if err != nil {
		res.Status = 401
		res.Message = "Failed to fetch asset details"
		res.Data = err.Error()
		return res, err
	}

	res.Status = http.StatusOK
	res.Message = "Berhasil mengambil data"
	res.Data = dtAset

	defer db.DbClose(con)
	return res, nil
}

func fetchAssetDetailedByUserId(con *sql.DB, user_id string) ([]Asset, error) {
	var arrAset []Asset

	queryPerusahaan := `
		SELECT id_perusahaan 
		FROM user_perusahaan
		WHERE id_user = ?
	`
	rows, err := con.Query(queryPerusahaan, user_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var perusahaanIds []int
	for rows.Next() {
		var perusahaanId int
		if err := rows.Scan(&perusahaanId); err != nil {
			return nil, err
		}
		perusahaanIds = append(perusahaanIds, perusahaanId)
	}

	if len(perusahaanIds) == 0 {
		return arrAset, nil
	}

	for _, perusahaanId := range perusahaanIds {
		query := `
			SELECT *
			FROM asset 
			WHERE perusahaan_id = ?
		`

		rowsAssets, err := con.Query(query, perusahaanId)
		if err != nil {
			return nil, err
		}
		defer rowsAssets.Close()

		for rowsAssets.Next() {
			var dtAset Asset
			var masaSewa, deleteAt []byte
			var idJoin, idAssetChild sql.NullString
			var idAssetParent, idPerusahaan, idProvinsi sql.NullInt32

			err := rowsAssets.Scan(
				&dtAset.Id_asset, &idAssetParent, &idAssetChild, &idJoin, &idPerusahaan, &dtAset.Nama,
				&dtAset.Tipe, &dtAset.Nomor_legalitas, &dtAset.File_legalitas, &dtAset.Status_asset, &dtAset.Surat_kuasa,
				&dtAset.Surat_legalitas, &dtAset.Alamat, &dtAset.Kondisi, &dtAset.Titik_koordinat, &dtAset.Batas_koordinat,
				&dtAset.Luas, &dtAset.Nilai, &idProvinsi, &dtAset.Status_pengecekan, &dtAset.Status_verifikasi,
				&dtAset.Status_publik, &dtAset.Hak_akses, &masaSewa, &dtAset.Created_at, &deleteAt,
			)
			if err != nil {
				return nil, err
			}

			// Handle nullable fields
			if idAssetParent.Valid {
				dtAset.Id_asset_parent = int(idAssetParent.Int32)
			} else {
				dtAset.Id_asset_parent = 0
			}
			if idAssetChild.Valid {
				dtAset.Id_asset_child = idAssetChild.String
			} else {
				dtAset.Id_asset_child = ""
			}
			if idJoin.Valid {
				dtAset.Id_join = idJoin.String
			} else {
				dtAset.Id_join = "0"
			}
			if idPerusahaan.Valid {
				dtAset.Id_perusahaan = int(idPerusahaan.Int32)
			} else {
				dtAset.Id_perusahaan = 0
			}
			if idProvinsi.Valid {
				dtAset.Provinsi = int(idProvinsi.Int32)
			} else {
				dtAset.Provinsi = 0
			}

			// Fetch child assets
			childQuery := "SELECT id_asset FROM asset WHERE id_asset_parent = ?"
			childRows, err := con.Query(childQuery, dtAset.Id_asset)
			if err != nil {
				return nil, err
			}
			defer childRows.Close()

			for childRows.Next() {
				var childId int
				err := childRows.Scan(&childId)
				if err != nil {
					return nil, err
				}
				childAset, err := fetchAssetDetailed(con, strconv.Itoa(childId))
				if err != nil {
					return nil, err
				}
				dtAset.ChildAssets = append(dtAset.ChildAssets, childAset)
			}

			// Fetch images
			imageQuery := "SELECT link_gambar FROM asset_gambar WHERE id_asset_gambar = ?"
			imageRows, err := con.Query(imageQuery, dtAset.Id_asset)
			if err != nil {
				return nil, err
			}
			defer imageRows.Close()

			for imageRows.Next() {
				var linkGambar string
				err := imageRows.Scan(&linkGambar)
				if err != nil {
					return nil, err
				}
				dtAset.LinkGambar = append(dtAset.LinkGambar, linkGambar)
			}

			// Fetch usages
			usageQuery := `SELECT p.nama FROM asset_penggunaan ap
				JOIN penggunaan p ON ap.id_penggunaan = p.id
				WHERE ap.id_asset = ?`
			usageRows, err := con.Query(usageQuery, dtAset.Id_asset)
			if err != nil {
				return nil, err
			}
			defer usageRows.Close()

			for usageRows.Next() {
				var usageName string
				err := usageRows.Scan(&usageName)
				if err != nil {
					return nil, err
				}
				dtAset.Usage = append(dtAset.Usage, usageName)
			}

			// Fetch tags
			tagQuery := `SELECT t.nama FROM asset_tags at
				JOIN tags t ON at.id_tags = t.id
				WHERE at.id_asset = ?`
			tagRows, err := con.Query(tagQuery, dtAset.Id_asset)
			if err != nil {
				return nil, err
			}
			defer tagRows.Close()

			for tagRows.Next() {
				var tagName string
				err := tagRows.Scan(&tagName)
				if err != nil {
					return nil, err
				}
				dtAset.TagsAssets = append(dtAset.TagsAssets, tagName)
			}

			// Append asset to the list
			arrAset = append(arrAset, dtAset)
		}
	}

	return arrAset, nil
}

func GetAssetByName(nama_aset string) (Response, error) {
	var res Response
	var dtAsets = []Asset{}

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	query := "SELECT * FROM asset WHERE nama LIKE ?"
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	rows, err := stmt.Query("%" + nama_aset + "%")
	if err != nil {
		res.Status = 401
		res.Message = "exec gagal"
		res.Data = err.Error()
		return res, err
	}
	defer rows.Close()

	for rows.Next() {
		var dtAset Asset
		var masaSewa sql.NullTime
		var deleteAt sql.NullTime
		err := rows.Scan(&dtAset.Id_asset, &dtAset.Id_asset_parent, &dtAset.Id_asset_child, &dtAset.Id_perusahaan, &dtAset.Nama, &dtAset.Tipe, &dtAset.Nomor_legalitas, &dtAset.File_legalitas, &dtAset.Status_asset, &dtAset.Surat_kuasa, &dtAset.Alamat, &dtAset.Kondisi, &dtAset.Titik_koordinat, &dtAset.Batas_koordinat, &dtAset.Luas, &dtAset.Nilai, &dtAset.Status_pengecekan, &dtAset.Status_verifikasi, &dtAset.Hak_akses, &masaSewa, &dtAset.Created_at, &deleteAt)
		if err != nil {
			res.Status = 401
			res.Message = "scan gagal"
			res.Data = err.Error()
			return res, err
		}

		if masaSewa.Valid {
			dtAset.Masa_sewa = masaSewa.Time.Format("2024-08-08")
		} else {
			dtAset.Masa_sewa = ""
		}
		if deleteAt.Valid {
			dtAset.Deleted_at = deleteAt.Time.Format("2024-08-08")
		} else {
			dtAset.Deleted_at = ""
		}
		dtAsets = append(dtAsets, dtAset)
	}

	if err = rows.Err(); err != nil {
		res.Status = 401
		res.Message = "rows error"
		res.Data = err.Error()
		return res, err
	}

	if len(dtAsets) == 0 {
		res.Status = 404
		res.Message = "Data tidak ditemukan"
		res.Data = nil
		return res, nil
	}

	res.Status = http.StatusOK
	res.Message = "Berhasil mengambil data"
	res.Data = dtAsets

	defer db.DbClose(con)
	return res, nil
}

func DeleteAssetById(aset string) (Response, error) {
	var res Response

	var dtAset = Asset{}

	err := json.Unmarshal([]byte(aset), &dtAset)
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

	query := "DELETE FROM asset WHERE id_asset_parent = ?"
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(dtAset.Id_asset_parent)
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

func JoinAsset(input string) (Response, error) {
	var res Response

	type TempJoinAsset struct {
		IdAsset1 int `json:"id_asset_1"`
		IdAsset2 int `json:"id_asset_2"`
	}
	var tempjoinAsset TempJoinAsset
	err := json.Unmarshal([]byte(input), &tempjoinAsset)
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

	fmt.Println(tempjoinAsset.IdAsset1)
	fmt.Println(tempjoinAsset.IdAsset2)
	var dtAsset1 Asset
	var dtAsset2 Asset

	queryAsset1 := "SELECT * FROM asset WHERE id_asset = ?"
	stmtAsset1, err := con.Prepare(queryAsset1)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmtAsset1.Close()

	var masaSewa sql.NullTime
	var deleteAt sql.NullTime
	var idAssetParent, idAssetChild, idJoin sql.NullInt32
	err = stmtAsset1.QueryRow(tempjoinAsset.IdAsset1).Scan(&dtAsset1.Id_asset, &idAssetParent, &idAssetChild, &idJoin,
		&dtAsset1.Id_perusahaan, &dtAsset1.Nama, &dtAsset1.Tipe, &dtAsset1.Nomor_legalitas,
		&dtAsset1.File_legalitas, &dtAsset1.Status_asset, &dtAsset1.Surat_kuasa, &dtAsset1.Alamat,
		&dtAsset1.Kondisi, &dtAsset1.Titik_koordinat, &dtAsset1.Batas_koordinat, &dtAsset1.Luas,
		&dtAsset1.Nilai, &dtAsset1.Provinsi, &dtAsset1.Usage, &dtAsset1.Status_pengecekan,
		&dtAsset1.Status_verifikasi, &dtAsset1.Status_publik, &dtAsset1.Hak_akses, &masaSewa,
		&dtAsset1.Created_at, &deleteAt)
	if err != nil {
		res.Status = 401
		res.Message = "exec gagal"
		res.Data = err.Error()
		return res, err
	}
	if masaSewa.Valid {
		dtAsset1.Masa_sewa = masaSewa.Time.Format("2024-08-08")
	} else {
		dtAsset1.Masa_sewa = ""
	}
	if deleteAt.Valid {
		dtAsset1.Deleted_at = deleteAt.Time.Format("2024-08-08")
	} else {
		dtAsset1.Deleted_at = ""
	}
	if idAssetParent.Valid {
		dtAsset1.Id_asset_parent = int(idAssetParent.Int32)
	} else {
		dtAsset1.Id_asset_parent = 0
	}
	if idAssetChild.Valid {
		dtAsset1.Id_asset_child = strconv.Itoa(int(idAssetChild.Int32))
	} else {
		dtAsset1.Id_asset_child = ""
	}
	if idJoin.Valid {
		dtAsset1.Id_join = strconv.Itoa(int(idJoin.Int32))
	} else {
		dtAsset1.Id_join = "0"
	}

	queryAsset2 := "SELECT * FROM asset WHERE id_asset = ?"
	stmtAsset2, err := con.Prepare(queryAsset2)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmtAsset2.Close()

	var masaSewa2 sql.NullTime
	var deleteAt2 sql.NullTime
	var idAssetParent2, idAssetChild2, idJoin2 sql.NullInt32
	err = stmtAsset1.QueryRow(tempjoinAsset.IdAsset2).Scan(&dtAsset2.Id_asset, &idAssetParent2, &idAssetChild2, &idJoin2,
		&dtAsset2.Id_perusahaan, &dtAsset2.Nama, &dtAsset2.Tipe, &dtAsset2.Nomor_legalitas,
		&dtAsset2.File_legalitas, &dtAsset2.Status_asset, &dtAsset2.Surat_kuasa, &dtAsset2.Alamat,
		&dtAsset2.Kondisi, &dtAsset2.Titik_koordinat, &dtAsset2.Batas_koordinat, &dtAsset2.Luas,
		&dtAsset2.Nilai, &dtAsset2.Provinsi, &dtAsset2.Usage, &dtAsset2.Status_pengecekan,
		&dtAsset2.Status_verifikasi, &dtAsset2.Status_publik, &dtAsset2.Hak_akses, &masaSewa2,
		&dtAsset2.Created_at, &deleteAt2)
	if err != nil {
		res.Status = 401
		res.Message = "exec gagal"
		res.Data = err.Error()
		return res, err
	}
	if masaSewa2.Valid {
		dtAsset2.Masa_sewa = masaSewa.Time.Format("2024-08-08")
	} else {
		dtAsset2.Masa_sewa = ""
	}
	if deleteAt2.Valid {
		dtAsset2.Deleted_at = deleteAt.Time.Format("2024-08-08")
	} else {
		dtAsset2.Deleted_at = ""
	}
	if idAssetParent2.Valid {
		dtAsset2.Id_asset_parent = int(idAssetParent.Int32)
	} else {
		dtAsset2.Id_asset_parent = 0
	}
	if idAssetChild2.Valid {
		dtAsset2.Id_asset_child = strconv.Itoa(int(idAssetChild.Int32))
	} else {
		dtAsset2.Id_asset_child = ""
	}
	if idJoin2.Valid {
		dtAsset2.Id_join = strconv.Itoa(int(idJoin.Int32))
	} else {
		dtAsset2.Id_join = "0"
	}

	luasBaru := dtAsset1.Luas + dtAsset2.Luas
	nilaiBaru := dtAsset1.Nilai + dtAsset2.Nilai
	tempIdPerusahaan1 := dtAsset1.Id_perusahaan
	tempIdPerusahaan2 := dtAsset2.Id_perusahaan
	if tempIdPerusahaan1 != tempIdPerusahaan2 {
		res.Status = 401
		res.Message = "id perusahaan tidak sama"
		return res, errors.New(res.Message)
	}
	if dtAsset1.Provinsi != dtAsset2.Provinsi {
		res.Status = 401
		res.Message = "provinsi tidak sama"
		return res, errors.New(res.Message)
	}

	tempIdJoin := strconv.Itoa(tempjoinAsset.IdAsset1) + "," + strconv.Itoa(tempjoinAsset.IdAsset2)

	query := `
	INSERT INTO asset (id_join, perusahaan_id, luas, nilai, provinsi, created_at) 
	VALUES (?,?,?,?,?,NOW())
	`
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(tempIdJoin, tempIdPerusahaan1, luasBaru, nilaiBaru, dtAsset1.Provinsi)
	if err != nil {
		res.Status = 401
		res.Message = "exec gagal"
		res.Data = err.Error()
		return res, err
	}
	lastId, err := result.LastInsertId()
	if err != nil {
		log.Println(err.Error())
		return res, err
	}

	queryUpdateAsset1 := "UPDATE asset SET deleted_at = NOW() WHERE id_asset = ?"
	stmtUpdateAsset1, err := con.Prepare(queryUpdateAsset1)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmtUpdateAsset1.Close()
	_, err = stmtUpdateAsset1.Exec(tempjoinAsset.IdAsset1)
	if err != nil {
		res.Status = 401
		res.Message = "exec gagal"
		res.Data = err.Error()
		return res, err
	}

	queryUpdateAsset2 := "UPDATE asset SET deleted_at = NOW() WHERE id_asset = ?"
	stmtUpdateAsset2, err := con.Prepare(queryUpdateAsset2)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmtUpdateAsset2.Close()
	_, err = stmtUpdateAsset2.Exec(tempjoinAsset.IdAsset2)
	if err != nil {
		res.Status = 401
		res.Message = "exec gagal"
		res.Data = err.Error()
		return res, err
	}

	var tempaset Response
	tempaset, _ = GetAssetById(strconv.Itoa(int(lastId)))

	res.Status = http.StatusOK
	res.Message = "Berhasil memasukkan data"
	res.Data = tempaset.Data

	defer db.DbClose(con)
	return res, nil
}

func GetAssetRentedByUserId(userId string) (Response, error) {
	var res Response
	// asset + transaction_request + progress
	var arrAsetTranReq = []TransactionRequest{}

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	query := `
	SELECT tr.id_transaksi_jual_sewa,tr.perusahaan_id,tr.user_id,tr.id_asset,
		a.nama,tr.status, tr.nama_progress,tr.proposal,tr.tgl_meeting,tr.waktu_meeting,
		tr.lokasi_meeting,tr.deskripsi,tr.alasan,IFNULL(tr.tgl_dateline,""),tr.created_at
	FROM transaction_request tr
	LEFT JOIN asset a ON tr.id_asset = a.id_asset
	WHERE tr.user_id = ? AND tr.status = 'A'
	`
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	nId, _ := strconv.Atoi(userId)
	rows, err := stmt.Query(nId)
	if err != nil {
		res.Status = 401
		res.Message = "exec gagal"
		res.Data = err.Error()
		return res, err
	}
	defer rows.Close()
	for rows.Next() {
		var dtTranReq TransactionRequest
		err := rows.Scan(&dtTranReq.Id_transaksi_jual_sewa, &dtTranReq.Perusahaan_id,
			&dtTranReq.User_id, &dtTranReq.Id_asset, &dtTranReq.Nama_aset, &dtTranReq.Status, &dtTranReq.Nama_progress,
			&dtTranReq.Proposal, &dtTranReq.Tgl_meeting, &dtTranReq.Waktu_meeting, &dtTranReq.Lokasi_meeting, &dtTranReq.Deskripsi,
			&dtTranReq.Alasan, &dtTranReq.Tgl_dateline, &dtTranReq.Created_at)
		if err != nil {
			res.Status = 401
			res.Message = "scan gagal"
			res.Data = err.Error()
			return res, err
		}

		arrAsetTranReq = append(arrAsetTranReq, dtTranReq)
	}

	if err = rows.Err(); err != nil {
		res.Status = 401
		res.Message = "rows error"
		res.Data = err.Error()
		return res, err
	}

	if len(arrAsetTranReq) == 0 {
		res.Status = 404
		res.Message = "Data tidak ditemukan"
		res.Data = nil
		return res, nil
	}

	res.Status = http.StatusOK
	res.Message = "Berhasil mengambil data"
	res.Data = arrAsetTranReq

	defer db.DbClose(con)

	return res, nil
}

func GetAssetSurveyHistoryByAssetId(assetId string) (Response, error) {
	var res Response
	type AssetSurveyHistory struct {
		UpdatedOn    string  `json:"updatedon"`
		SurveyorName string  `json:"surveyorname"`
		ValueName    string  `json:"value_name"`
		ValueOld     float64 `json:"value_old"`
		ValueNew     float64 `json:"value_new"`
		KondisiOld   string  `json:"kondisi_old"`
		KondisiNew   string  `json:"kondisi_new"`
	}
	var arrAsetTranReq = []AssetSurveyHistory{}

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	query := `
	SELECT sr.created_at,u.nama_lengkap,sr.nilai_old,sr.nilai_new,sr.kondisi_old, sr.kondisi_new
	FROM survey_request sr
	LEFT JOIN user u ON sr.user_id = u.user_id
	WHERE sr.id_asset = ?
	`
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	nId, _ := strconv.Atoi(assetId)
	rows, err := stmt.Query(nId)
	if err != nil {
		res.Status = 401
		res.Message = "exec gagal"
		res.Data = err.Error()
		return res, err
	}
	defer rows.Close()
	for rows.Next() {
		var dtTranReq AssetSurveyHistory
		err := rows.Scan(&dtTranReq.UpdatedOn, &dtTranReq.SurveyorName,
			&dtTranReq.ValueOld, &dtTranReq.ValueNew, &dtTranReq.KondisiOld, &dtTranReq.KondisiNew)
		if err != nil {
			res.Status = 401
			res.Message = "scan gagal"
			res.Data = err.Error()
			return res, err
		}

		if dtTranReq.ValueOld == dtTranReq.ValueNew {
			dtTranReq.ValueName = "Same"
		} else if dtTranReq.ValueOld > dtTranReq.ValueNew {
			dtTranReq.ValueName = "Decreasing"
		} else {
			dtTranReq.ValueName = "Increasing"
		}

		arrAsetTranReq = append(arrAsetTranReq, dtTranReq)
	}

	if err = rows.Err(); err != nil {
		res.Status = 401
		res.Message = "rows error"
		res.Data = err.Error()
		return res, err
	}

	if len(arrAsetTranReq) == 0 {
		res.Status = 404
		res.Message = "Data tidak ditemukan"
		res.Data = nil
		return res, nil
	}

	res.Status = http.StatusOK
	res.Message = "Berhasil mengambil data"
	res.Data = arrAsetTranReq

	defer db.DbClose(con)
	return res, nil
}

func UpdateAssetByIdWithoutGambar(filelegalitas *multipart.FileHeader, suratkuasa *multipart.FileHeader,
	id_asset, nama, surat_legalitas, tipe, usage, tag, nomor_legalitas, status,
	alamat, kondisi, koordinat, batas_koordinat, luas, nilai string) (Response, error) {
	var res Response
	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	// hapus file legalitas + surat kuasa kalau ada
	linkasetQueue := `
		SELECT file_legalitas, surat_kuasa
		FROM asset
		WHERE id_asset = ?
	`
	var linkfilelegalitas string
	var linksuratkuasa string
	err = con.QueryRow(linkasetQueue, id_asset).Scan(&linkfilelegalitas, &linksuratkuasa)
	if err != nil {
		res.Status = 401
		res.Message = "exec gagal"
		res.Data = err
		return res, err
	}
	err = os.Remove(linkfilelegalitas)
	if err != nil {
		return res, err
	}
	err = os.Remove(linksuratkuasa)
	if err != nil {
		return res, err
	}

	usageIds := strings.Split(usage, ",")
	for _, id := range usageIds {
		var usageExists bool
		fmt.Println("usage", id)
		usageQuery := "SELECT EXISTS(SELECT 1 FROM penggunaan WHERE id = ?)"
		err = con.QueryRow(usageQuery, id).Scan(&usageExists)
		if err != nil || !usageExists {
			res.Status = 401
			res.Message = "Penggunaan tidak valid"
			res.Data = "Penggunaan ID " + id + " tidak ditemukan"
			return res, err
		}
	}

	tagIds := strings.Split(tag, ",")
	for _, id2 := range tagIds {
		var tagExists bool
		fmt.Println("tag", id2)
		tagQuery := "SELECT EXISTS(SELECT 1 FROM tags WHERE id = ?)"
		err = con.QueryRow(tagQuery, id2).Scan(&tagExists)
		if err != nil || !tagExists {
			res.Status = 401
			res.Message = "Tag tidak valid"
			res.Data = "Tag ID " + id2 + " tidak ditemukan"
			return res, err
		}
	}

	// hapus tags dan usage di db
	tagQuery := "DELETE FROM `asset_tag` WHERE id_asset = ?"
	_, err = con.Exec(tagQuery, id_asset)
	if err != nil {
		res.Status = 401
		res.Message = "Delete tag exec gagal"
		res.Data = "Delete tag exec gagal dengan id aset " + id_asset
		return res, err
	}
	usageQuery := "DELETE FROM `asset_penggunaan` WHERE id_asset = ?"
	_, err = con.Exec(usageQuery, id_asset)
	if err != nil {
		res.Status = 401
		res.Message = "Delete usage exec gagal"
		res.Data = "Delete usage exec gagal dengan id aset " + id_asset
		return res, err
	}

	query := `
	UPDATE asset 
	SET nama = ?, tipe = ?, nomor_legalitas = ?, status_asset = ?, surat_legalitas = ?, alamat = ?, kondisi = ?, titik_koordinat = ?,
	batas_koordinat = ?, luas = ?, nilai = ?
	WHERE id_asset = ?
	`
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()
	result, err := stmt.Exec(
		nama, tipe, nomor_legalitas, status, surat_legalitas, alamat, kondisi, koordinat, batas_koordinat, luas, nilai, id_asset,
	)
	if err != nil {
		res.Status = 401
		res.Message = "exec gagal"
		res.Data = err.Error()
		return res, err
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		log.Println(err.Error())
		return res, err
	}

	// tambah usage + tags
	for _, usageId := range usageIds {
		usageQuery := "INSERT INTO asset_penggunaan (id_asset, id_penggunaan) VALUES (?, ?)"
		_, err = con.Exec(usageQuery, lastId, usageId)
		if err != nil {
			res.Status = 401
			res.Message = "Gagal menambah penggunaan"
			res.Data = err.Error()
			return res, err
		}
	}

	for _, tagId := range tagIds {
		tagQuery := "INSERT INTO asset_tags (id_asset, id_tags) VALUES (?, ?)"
		_, err = con.Exec(tagQuery, lastId, tagId)
		if err != nil {
			res.Status = 401
			res.Message = "Gagal menambah tag"
			res.Data = err.Error()
			return res, err
		}
	}

	// tambah filelegalitas
	//source
	srclegalitas, err := filelegalitas.Open()
	if err != nil {
		log.Println(err.Error())
		return res, err
	}
	defer srclegalitas.Close()

	tempid := int(lastId)
	_tempid := strconv.Itoa(tempid)
	// Destination
	filelegalitas.Filename = _tempid + "_" + filelegalitas.Filename
	pathFileLegalitas := "uploads/asset/file_legalitas/" + filelegalitas.Filename
	dstlegalitas, err := os.Create("uploads/asset/file_legalitas/" + filelegalitas.Filename)
	if err != nil {
		log.Println("2")
		fmt.Print("2")
		return res, err
	}

	// Copy
	if _, err = io.Copy(dstlegalitas, srclegalitas); err != nil {
		log.Println(err.Error())
		log.Println("3")
		fmt.Print("3")
		return res, err
	}
	dstlegalitas.Close()

	err = UpdateDataFotoPath("asset", "file_legalitas", pathFileLegalitas, "id_asset", int(lastId))
	if err != nil {
		return res, err
	}

	// tambah suratkuasa
	//source
	srcsuratkuasa, err := suratkuasa.Open()
	if err != nil {
		log.Println(err.Error())
		return res, err
	}
	defer srcsuratkuasa.Close()

	// Destination
	suratkuasa.Filename = _tempid + "_" + suratkuasa.Filename
	pathFileSuratKuasa := "uploads/asset/surat_kuasa/" + suratkuasa.Filename
	dstsuratkuasa, err := os.Create("uploads/asset/surat_kuasa/" + suratkuasa.Filename)
	if err != nil {
		log.Println("2")
		fmt.Print("2")
		return res, err
	}

	// Copy
	if _, err = io.Copy(dstsuratkuasa, srcsuratkuasa); err != nil {
		log.Println(err.Error())
		log.Println("3")
		fmt.Print("3")
		return res, err
	}
	dstsuratkuasa.Close()

	err = UpdateDataFotoPath("asset", "surat_kuasa", pathFileSuratKuasa, "id_asset", int(lastId))
	if err != nil {
		return res, err
	}

	tempaset, err := GetAssetById(_tempid)
	if err != nil {
		res.Status = 401
		res.Message = "get aset by id gagal"
		res.Data = err.Error()
		return res, err
	}

	res.Status = http.StatusOK
	res.Message = "Berhasil memasukkan data"
	res.Data = tempaset

	defer db.DbClose(con)
	return res, nil
}

func UnjoinAsset(input string) (Response, error) {
	var res Response
	// del

	return res, nil
}
