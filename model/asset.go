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
	"path/filepath"
	"sort"
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

func CreateAssetMultipleFile(filelegalitas *multipart.FileHeader, suratkuasa *multipart.FileHeader,
	nama, surat_legalitas, tipe, usage, tag, nomorlegalitas, status, alamat, kondisi, koordinat, batas_koordinat, luas, nilai, provinsi string, files []*multipart.FileHeader) (Response, error) {

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

	// file gambar

	filePaths := []string{}
	for _, fileHeader := range files {
		filePath, err := saveAssetGambar(strconv.Itoa(int(lastId)), fileHeader)
		if err != nil {
			res.Status = 500
			res.Message = "Error saving file"
			res.Data = err.Error()
			return res, err
		}
		filePaths = append(filePaths, filePath)
	}

	for _, filePath := range filePaths {
		// Prepare the INSERT statement for each filePath
		insertImageQuery := "INSERT INTO asset_gambar (id_asset_gambar, link_gambar) VALUES (?, ?)"
		_, err = con.Exec(insertImageQuery, strconv.Itoa(int(lastId)), filePath)
		if err != nil {
			// Handle the error if the insert fails
			res.Status = 500
			res.Message = "Error inserting image path"
			res.Data = err.Error()
			return res, err
		}
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

	var tempaset Response
	tempaset, _ = GetAssetById(strconv.Itoa(tempid))

	res.Status = http.StatusOK
	res.Message = "Berhasil memasukkan data"
	res.Data = tempaset.Data

	defer db.DbClose(con)
	return res, nil
}

func saveAssetGambar(idAset string, fileHeader *multipart.FileHeader) (string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	uploadDir := "uploads/asset/foto/"
	fileHeader.Filename = idAset + "_" + fileHeader.Filename
	os.MkdirAll(uploadDir, os.ModePerm)
	filePath := filepath.Join(uploadDir, fileHeader.Filename)

	out, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		return "", err
	}

	return filePath, nil
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
	ParentQuery := "SELECT provinsi,IFNULL(id_asset_child,'') FROM asset WHERE id_asset = ?"
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

func CreateAssetChildMultipleGambar(
	filelegalitas *multipart.FileHeader, suratkuasa *multipart.FileHeader,
	parent_id, nama, surat_legalitas, tipe, usage, tag, nomor_legalitas, status,
	alamat, kondisi, koordinat, batas_koordinat, luas, nilai string, files []*multipart.FileHeader) (Response, error) {
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
	ParentQuery := "SELECT provinsi,IFNULL(id_asset_child,'') FROM asset WHERE id_asset = ?"
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

	filePaths := []string{}
	for _, fileHeader := range files {
		filePath, err := saveAssetGambar(strconv.Itoa(int(lastId)), fileHeader)
		if err != nil {
			res.Status = 500
			res.Message = "Error saving file"
			res.Data = err.Error()
			return res, err
		}
		filePaths = append(filePaths, filePath)
	}

	for _, filePath := range filePaths {
		// Prepare the INSERT statement for each filePath
		insertImageQuery := "INSERT INTO asset_gambar (id_asset_gambar, link_gambar) VALUES (?, ?)"
		_, err = con.Exec(insertImageQuery, strconv.Itoa(int(lastId)), filePath)
		if err != nil {
			// Handle the error if the insert fails
			res.Status = 500
			res.Message = "Error inserting image path"
			res.Data = err.Error()
			return res, err
		}
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
	var arrAset []Asset
	assetMap := make(map[int]*Asset) // Map to store assets by their id_asset

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	query := `SELECT a.*, IFNULL(ag.link_gambar,'') as link_gambar
			  FROM asset a
			  LEFT JOIN asset_gambar ag ON a.id_asset = ag.id_asset_gambar
			  WHERE a.deleted_at IS NULL`
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

	var masaSewa, deleteAt []byte
	var linkGambar sql.NullString
	var idJoin, idAssetChild, tipe sql.NullString
	var idAssetParent, idProvinsi sql.NullInt32

	for result.Next() {
		var dtAset Asset
		err = result.Scan(
			&dtAset.Id_asset, &idAssetParent, &idAssetChild, &idJoin, &dtAset.Nama, &tipe, &dtAset.Nomor_legalitas, &dtAset.File_legalitas, &dtAset.Status_asset,
			&dtAset.Surat_kuasa, &dtAset.Surat_legalitas, &dtAset.Alamat, &dtAset.Kondisi, &dtAset.Titik_koordinat,
			&dtAset.Batas_koordinat, &dtAset.Luas, &dtAset.Nilai, &idProvinsi, &dtAset.Status_pengecekan,
			&dtAset.Status_verifikasi, &dtAset.Status_publik, &dtAset.Hak_akses, &masaSewa, &dtAset.Created_at,
			&deleteAt, &linkGambar)
		if err != nil {
			res.Status = 401
			res.Message = "rows scan"
			res.Data = err.Error()
			return res, err
		}

		// Handle potential null values
		if masaSewa != nil {
			masaSewaWaktu, masaSewaErr := time.Parse("2006-01-02 15:04:05", string(masaSewa))
			if masaSewaErr != nil {
				dtAset.Masa_sewa = ""
			} else {
				dtAset.Masa_sewa = masaSewaWaktu.Format("2006-01-02 15:04:05")
			}
		} else {
			dtAset.Masa_sewa = ""
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
		if tipe.Valid {
			dtAset.Tipe = tipe.String
		} else {
			dtAset.Tipe = ""
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
		if idProvinsi.Valid {
			dtAset.Provinsi = int(idProvinsi.Int32)
		} else {
			dtAset.Provinsi = 0
		}

		// Check if the asset already exists in the map
		if asset, exists := assetMap[dtAset.Id_asset]; exists {
			// If asset exists, append the image to its list
			if linkGambar.Valid && linkGambar.String != "" {
				asset.LinkGambar = append(asset.LinkGambar, linkGambar.String)
			}
		} else {
			// If it's a new asset, initialize the LinkGambar slice and add it to the map
			if linkGambar.Valid && linkGambar.String != "" {
				dtAset.LinkGambar = []string{linkGambar.String}
			} else {
				dtAset.LinkGambar = []string{}
			}
			assetMap[dtAset.Id_asset] = &dtAset
		}
	}

	// Convert map to slice
	for _, asset := range assetMap {
		// Ensure that LinkGambar is truly empty if no valid images were found
		if len(asset.LinkGambar) == 1 && asset.LinkGambar[0] == "" {
			asset.LinkGambar = []string{}
		}
		arrAset = append(arrAset, *asset)
	}

	// Sort arrAset by Id_asset
	sort.Slice(arrAset, func(i, j int) bool {
		return arrAset[i].Id_asset < arrAset[j].Id_asset
	})

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
	var idJoin, idAssetChild, tipe sql.NullString
	var idAssetParent, idProvinsi sql.NullInt32
	err = stmt.QueryRow(nId).Scan(&dtAset.Id_asset, &idAssetParent, &idAssetChild, &idJoin, &dtAset.Nama, &tipe, &dtAset.Nomor_legalitas, &dtAset.File_legalitas, &dtAset.Status_asset, &dtAset.Surat_kuasa, &dtAset.Surat_legalitas, &dtAset.Alamat, &dtAset.Kondisi, &dtAset.Titik_koordinat, &dtAset.Batas_koordinat, &dtAset.Luas, &dtAset.Nilai, &idProvinsi, &dtAset.Status_pengecekan, &dtAset.Status_verifikasi, &dtAset.Status_publik, &dtAset.Hak_akses, &masaSewa, &dtAset.Created_at, &deleteAt)
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
	if tipe.Valid {
		dtAset.Tipe = tipe.String
	} else {
		dtAset.Tipe = ""
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

	usageQuery := "SELECT p.id,p.nama FROM asset_penggunaan ap JOIN penggunaan p ON ap.id_penggunaan = p.id WHERE ap.id_asset = ?"
	rowsusage, err := con.Query(usageQuery, nId)
	if err != nil {
		res.Status = 401
		res.Message = "gagal mengambil gambar"
		res.Data = err.Error()
		return res, err
	}
	defer rowsusage.Close()

	var usage []Kegunaan
	for rowsusage.Next() {
		var link Kegunaan
		err := rowsusage.Scan(&link.Id, &link.Nama)
		if err != nil {
			res.Status = 401
			res.Message = "gagal membaca gambar"
			res.Data = err.Error()
			return res, err
		}
		usage = append(usage, link)
	}
	dtAset.Usage = usage

	tagsQuery := "SELECT t.id,t.nama FROM asset_tags at JOIN tags t ON at.id_tags = t.id WHERE at.id_asset = ?"
	tagsRows, err := con.Query(tagsQuery, nId)
	if err != nil {
		res.Status = 401
		res.Message = "gagal mengambil tags"
		res.Data = err.Error()
		return res, err
	}
	defer tagsRows.Close()

	var tags []Tags
	for tagsRows.Next() {
		var tag Tags
		err := tagsRows.Scan(&tag.Id, &tag.Nama)
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
			var idAssetParent, idProvinsi sql.NullInt32
			err = con.QueryRow(usageQuery, id).Scan(&dtAset.Id_asset, &idAssetParent, &idAssetChild, &idJoin, &dtAset.Nama, &dtAset.Tipe, &dtAset.Nomor_legalitas, &dtAset.File_legalitas, &dtAset.Status_asset, &dtAset.Surat_kuasa, &dtAset.Surat_legalitas, &dtAset.Alamat, &dtAset.Kondisi, &dtAset.Titik_koordinat, &dtAset.Batas_koordinat, &dtAset.Luas, &dtAset.Nilai, &idProvinsi, &dtAset.Status_pengecekan, &dtAset.Status_verifikasi, &dtAset.Status_publik, &dtAset.Hak_akses, &masaSewa, &dtAset.Created_at, &deleteAt)
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

	query := "SELECT * FROM asset WHERE id_asset = ? AND deleted_at IS NULL"
	stmt, err := con.Prepare(query)
	if err != nil {
		return dtAset, err
	}
	defer stmt.Close()

	nId, _ := strconv.Atoi(aset_id)
	var masaSewa []byte
	var deleteAt []byte
	var idJoin, idAssetChild, tipe sql.NullString
	var idAssetParent, idProvinsi sql.NullInt32
	err = stmt.QueryRow(nId).Scan(
		&dtAset.Id_asset, &idAssetParent, &idAssetChild, &idJoin, &dtAset.Nama, &tipe, &dtAset.Nomor_legalitas, &dtAset.File_legalitas, &dtAset.Status_asset, &dtAset.Surat_kuasa, &dtAset.Surat_legalitas, &dtAset.Alamat, &dtAset.Kondisi, &dtAset.Titik_koordinat, &dtAset.Batas_koordinat, &dtAset.Luas, &dtAset.Nilai, &idProvinsi, &dtAset.Status_pengecekan, &dtAset.Status_verifikasi, &dtAset.Status_publik, &dtAset.Hak_akses, &masaSewa, &dtAset.Created_at, &deleteAt)
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
	if tipe.Valid {
		dtAset.Tipe = tipe.String
	} else {
		dtAset.Tipe = ""
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
	if idProvinsi.Valid {
		dtAset.Provinsi = int(idProvinsi.Int32)
	} else {
		dtAset.Provinsi = 0
	}

	fmt.Println("ambil gambar")
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

	fmt.Println("ambil usage")
	usageQuery := "SELECT p.id,p.nama FROM asset_penggunaan ap JOIN penggunaan p ON ap.id_penggunaan = p.id WHERE ap.id_asset = ?"
	rowsusage, err := con.Query(usageQuery, nId)
	if err != nil {
		return dtAset, err
	}
	defer rowsusage.Close()

	var usage []Kegunaan
	for rowsusage.Next() {
		var link Kegunaan
		err := rowsusage.Scan(&link.Id, &link.Nama)
		if err != nil {
			return dtAset, err
		}
		usage = append(usage, link)
	}
	dtAset.Usage = usage

	fmt.Println("ambil tag")
	tagsQuery := "SELECT t.id,t.nama FROM asset_tags at JOIN tags t ON at.id_tags = t.id WHERE at.id_asset = ?"
	tagsRows, err := con.Query(tagsQuery, nId)
	if err != nil {
		return dtAset, err
	}
	defer tagsRows.Close()

	var tags []Tags
	for tagsRows.Next() {
		var tag Tags
		err := tagsRows.Scan(&tag.Id, &tag.Nama)
		if err != nil {
			return dtAset, err
		}
		tags = append(tags, tag)
	}
	dtAset.TagsAssets = tags

	fmt.Println("ambil child")
	fmt.Println("id ", dtAset.Id_asset)
	childIds := strings.Split(dtAset.Id_asset_child, ",")
	fmt.Println("ambil child id", childIds)
	for _, childId := range childIds {
		trimmedChildId := strings.TrimSpace(childId)

		var deletedAt sql.NullString
		checkQuery := "SELECT deleted_at FROM asset WHERE id_asset = ?"
		err := con.QueryRow(checkQuery, trimmedChildId).Scan(&deletedAt)
		if err != nil {
			if err == sql.ErrNoRows {
				fmt.Printf("Child asset %s not found, skipping.\n", trimmedChildId)
				continue
			}
			return dtAset, err
		}

		if deletedAt.Valid {
			fmt.Printf("Skipping child asset %s because it has been deleted.\n", trimmedChildId)
			continue
		}

		childAset, err := fetchAssetDetailed(con, trimmedChildId)
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

	fmt.Println("Fetching assets by perusahaan_id:", perusahaan_id)

	query := `
	SELECT a.* 
	FROM transaction_request tr
	JOIN asset a ON tr.id_asset = a.id_asset
	WHERE tr.perusahaan_id = ?
	GROUP BY a.id_asset
	`
	rows, err := con.Query(query, perusahaan_id)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var dtAset Asset
		var masaSewa, deleteAt []byte
		var idJoin, idAssetChild sql.NullString
		var idAssetParent, idProvinsi sql.NullInt32

		err := rows.Scan(&dtAset.Id_asset, &idAssetParent, &idAssetChild, &idJoin, &dtAset.Nama,
			&dtAset.Tipe, &dtAset.Nomor_legalitas, &dtAset.File_legalitas, &dtAset.Status_asset, &dtAset.Surat_kuasa,
			&dtAset.Surat_legalitas, &dtAset.Alamat, &dtAset.Kondisi, &dtAset.Titik_koordinat, &dtAset.Batas_koordinat,
			&dtAset.Luas, &dtAset.Nilai, &idProvinsi, &dtAset.Status_pengecekan, &dtAset.Status_verifikasi,
			&dtAset.Status_publik, &dtAset.Hak_akses, &masaSewa, &dtAset.Created_at, &deleteAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		// Assign nullable fields
		dtAset.Id_asset_parent = ifValidInt(idAssetParent)
		dtAset.Id_asset_child = ifValidString(idAssetChild)
		dtAset.Id_join = ifValidString(idJoin)
		dtAset.Provinsi = ifValidInt(idProvinsi)

		// Fetch associated images
		if dtAset.LinkGambar, err = fetchGambarLinks(con, dtAset.Id_asset); err != nil {
			return nil, err
		}

		// Fetch associated usage
		if dtAset.Usage, err = fetchUsage(con, dtAset.Id_asset); err != nil {
			return nil, err
		}

		// Fetch associated tags
		if dtAset.TagsAssets, err = fetchTags(con, dtAset.Id_asset); err != nil {
			return nil, err
		}

		// Fetch child assets
		if dtAset.ChildAssets, err = fetchChildAssets(con, dtAset.Id_asset); err != nil {
			return nil, err
		}

		assets = append(assets, dtAset)
	}

	if len(assets) == 0 {
		fmt.Println("No assets found for perusahaan_id:", perusahaan_id)
	}

	return assets, nil
}

// Helper function to handle nullable integers
func ifValidInt(value sql.NullInt32) int {
	if value.Valid {
		return int(value.Int32)
	}
	return 0
}

// Helper function to handle nullable strings
func ifValidString(value sql.NullString) string {
	if value.Valid {
		return value.String
	}
	return ""
}

// Fetch related gambar links
func fetchGambarLinks(con *sql.DB, assetID int) ([]string, error) {
	var links []string
	query := "SELECT link_gambar FROM asset_gambar WHERE id_asset_gambar = ?"
	rows, err := con.Query(query, assetID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var link string
		if err := rows.Scan(&link); err != nil {
			return nil, err
		}
		links = append(links, link)
	}
	return links, nil
}

// Fetch associated usage
func fetchUsage(con *sql.DB, assetID int) ([]Kegunaan, error) {
	var usage []Kegunaan
	query := "SELECT p.id, p.nama FROM asset_penggunaan ap JOIN penggunaan p ON ap.id_penggunaan = p.id WHERE ap.id_asset = ?"
	rows, err := con.Query(query, assetID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var u Kegunaan
		if err := rows.Scan(&u.Id, &u.Nama); err != nil {
			return nil, err
		}
		usage = append(usage, u)
	}
	return usage, nil
}

// Fetch associated tags
func fetchTags(con *sql.DB, assetID int) ([]Tags, error) {
	var tags []Tags
	query := "SELECT t.id, t.nama FROM asset_tags at JOIN tags t ON at.id_tags = t.id WHERE at.id_asset = ?"
	rows, err := con.Query(query, assetID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var tag Tags
		if err := rows.Scan(&tag.Id, &tag.Nama); err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}
	return tags, nil
}

// Fetch child assets
func fetchChildAssets(con *sql.DB, parentID int) ([]Asset, error) {
	var childAssets []Asset

	// Fetch the id_asset_child field for the given parent asset.
	query := "SELECT id_asset_child FROM asset WHERE id_asset = ?"
	var idAssetChild sql.NullString

	if err := con.QueryRow(query, parentID).Scan(&idAssetChild); err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("No child assets found for parent ID:", parentID)
			return []Asset{}, nil // Return an empty slice if no children are found.
		}
		return nil, fmt.Errorf("failed to fetch id_asset_child: %w", err)
	}

	if !idAssetChild.Valid || idAssetChild.String == "" {
		fmt.Println("No child asset IDs for parent ID:", parentID)
		return []Asset{}, nil
	}

	// Split the comma-separated child IDs.
	childIDs := strings.Split(idAssetChild.String, ",")

	for _, idStr := range childIDs {
		childID, err := strconv.Atoi(strings.TrimSpace(idStr))
		if err != nil {
			return nil, fmt.Errorf("invalid child ID: %s, error: %w", idStr, err)
		}

		// Fetch detailed information for each child asset.
		child, err := fetchAssetDetailed(con, strconv.Itoa(childID))
		if err != nil {
			return nil, fmt.Errorf("failed to fetch child asset %d: %w", childID, err)
		}

		fmt.Println("Child Asset ID:", childID)
		fmt.Println("Child Asset:", child)

		childAssets = append(childAssets, child)
	}

	fmt.Println("All Child Assets:", childAssets)
	return childAssets, nil
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
		SELECT a.* 
		FROM transaction_request tr
		LEFT JOIN asset a ON tr.id_asset = a.id_asset
		WHERE tr.perusahaan_id = ?
		`
		fmt.Println(perusahaanId)
		rowsAssets, err := con.Query(query, perusahaanId)
		if err != nil {
			return nil, err
		}
		defer rowsAssets.Close()

		for rowsAssets.Next() {
			var dtAset Asset
			var masaSewa, deleteAt []byte
			var idJoin, idAssetChild sql.NullString
			var idAssetParent, idProvinsi sql.NullInt32

			err := rowsAssets.Scan(
				&dtAset.Id_asset, &idAssetParent, &idAssetChild, &idJoin, &dtAset.Nama,
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
			usageQuery := `SELECT p.id,p.nama FROM asset_penggunaan ap
				JOIN penggunaan p ON ap.id_penggunaan = p.id
				WHERE ap.id_asset = ?`
			usageRows, err := con.Query(usageQuery, dtAset.Id_asset)
			if err != nil {
				return nil, err
			}
			defer usageRows.Close()

			for usageRows.Next() {
				var usageName Kegunaan
				err := usageRows.Scan(&usageName.Id, &usageName.Nama)
				if err != nil {
					return nil, err
				}
				dtAset.Usage = append(dtAset.Usage, usageName)
			}

			// Fetch tags
			tagQuery := `SELECT t.id,t.nama FROM asset_tags at
				JOIN tags t ON at.id_tags = t.id
				WHERE at.id_asset = ?`
			tagRows, err := con.Query(tagQuery, dtAset.Id_asset)
			if err != nil {
				return nil, err
			}
			defer tagRows.Close()

			for tagRows.Next() {
				var tagName Tags
				err := tagRows.Scan(&tagName.Id, &tagName.Nama)
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
		err := rows.Scan(&dtAset.Id_asset, &dtAset.Id_asset_parent, &dtAset.Id_asset_child, &dtAset.Nama, &dtAset.Tipe, &dtAset.Nomor_legalitas, &dtAset.File_legalitas, &dtAset.Status_asset, &dtAset.Surat_kuasa, &dtAset.Alamat, &dtAset.Kondisi, &dtAset.Titik_koordinat, &dtAset.Batas_koordinat, &dtAset.Luas, &dtAset.Nilai, &dtAset.Status_pengecekan, &dtAset.Status_verifikasi, &dtAset.Hak_akses, &masaSewa, &dtAset.Created_at, &deleteAt)
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

	var masaSewa []byte
	var deleteAt []byte
	var idAssetParent, idAssetChild, idJoin sql.NullInt32
	err = stmtAsset1.QueryRow(tempjoinAsset.IdAsset1).Scan(&dtAsset1.Id_asset, &idAssetParent, &idAssetChild, &idJoin,
		&dtAsset1.Nama, &dtAsset1.Tipe, &dtAsset1.Nomor_legalitas,
		&dtAsset1.File_legalitas, &dtAsset1.Status_asset, &dtAsset1.Surat_kuasa, &dtAsset1.Surat_legalitas, &dtAsset1.Alamat,
		&dtAsset1.Kondisi, &dtAsset1.Titik_koordinat, &dtAsset1.Batas_koordinat, &dtAsset1.Luas,
		&dtAsset1.Nilai, &dtAsset1.Provinsi, &dtAsset1.Status_pengecekan,
		&dtAsset1.Status_verifikasi, &dtAsset1.Status_publik, &dtAsset1.Hak_akses, &masaSewa,
		&dtAsset1.Created_at, &deleteAt)
	if err != nil {
		res.Status = 401
		res.Message = "exec gagal"
		res.Data = err.Error()
		return res, err
	}
	if masaSewa != nil {
		masaSewaWaktu, masaSewaErr := time.Parse("2006-01-02 15:04:05", string(deleteAt))
		if masaSewaErr != nil {
			dtAsset1.Deleted_at = ""
		} else {
			dtAsset1.Deleted_at = masaSewaWaktu.Format("2006-01-02 15:04:05")
		}
	} else {
		dtAsset1.Deleted_at = ""
	}

	if deleteAt != nil {
		parsedTime, parseErr := time.Parse("2006-01-02 15:04:05", string(deleteAt))
		if parseErr != nil {
			dtAsset1.Deleted_at = ""
		} else {
			dtAsset1.Deleted_at = parsedTime.Format("2006-01-02 15:04:05")
		}
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

	var masaSewa2 []byte
	var deleteAt2 []byte
	var idAssetParent2, idAssetChild2, idJoin2 sql.NullInt32
	err = stmtAsset1.QueryRow(tempjoinAsset.IdAsset2).Scan(&dtAsset2.Id_asset, &idAssetParent2, &idAssetChild2, &idJoin2,
		&dtAsset2.Nama, &dtAsset2.Tipe, &dtAsset2.Nomor_legalitas,
		&dtAsset2.File_legalitas, &dtAsset2.Status_asset, &dtAsset2.Surat_kuasa, &dtAsset2.Surat_legalitas, &dtAsset2.Alamat,
		&dtAsset2.Kondisi, &dtAsset2.Titik_koordinat, &dtAsset2.Batas_koordinat, &dtAsset2.Luas,
		&dtAsset2.Nilai, &dtAsset2.Provinsi, &dtAsset2.Status_pengecekan,
		&dtAsset2.Status_verifikasi, &dtAsset2.Status_publik, &dtAsset2.Hak_akses, &masaSewa2,
		&dtAsset2.Created_at, &deleteAt2)
	if err != nil {
		res.Status = 401
		res.Message = "exec gagal"
		res.Data = err.Error()
		return res, err
	}
	if masaSewa != nil {
		masaSewaWaktu, masaSewaErr := time.Parse("2006-01-02 15:04:05", string(deleteAt))
		if masaSewaErr != nil {
			dtAsset2.Deleted_at = ""
		} else {
			dtAsset2.Deleted_at = masaSewaWaktu.Format("2006-01-02 15:04:05")
		}
	} else {
		dtAsset2.Deleted_at = ""
	}

	if deleteAt != nil {
		parsedTime, parseErr := time.Parse("2006-01-02 15:04:05", string(deleteAt))
		if parseErr != nil {
			dtAsset2.Deleted_at = ""
		} else {
			dtAsset2.Deleted_at = parsedTime.Format("2006-01-02 15:04:05")
		}
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

	tempIdJoin1 := strconv.Itoa(tempjoinAsset.IdAsset1)
	tempIdJoin2 := strconv.Itoa(tempjoinAsset.IdAsset2)
	tempGabungIdJoin := tempIdJoin1 + "," + tempIdJoin2
	var existingAssetId int
	queryCheck := `SELECT id_asset FROM asset WHERE id_join = ?`
	err = con.QueryRow(queryCheck, tempGabungIdJoin).Scan(&existingAssetId)
	if err == nil {
		queryRemoveDeletedAt := `UPDATE asset SET deleted_at = NULL WHERE id_asset = ?`
		_, err = con.Exec(queryRemoveDeletedAt, existingAssetId)
		if err != nil {
			res.Status = 401
			res.Message = "gagal menghapus deleted_at"
			res.Data = err.Error()
			return res, err
		}

		querySetDeletedAt := `UPDATE asset SET deleted_at = NOW() WHERE id_asset = ?`
		_, err = con.Exec(querySetDeletedAt, tempjoinAsset.IdAsset1)
		if err != nil {
			res.Status = 401
			res.Message = "gagal mengatur deleted_at untuk IdAsset1"
			res.Data = err.Error()
			return res, err
		}
		_, err = con.Exec(querySetDeletedAt, tempjoinAsset.IdAsset2)
		if err != nil {
			res.Status = 401
			res.Message = "gagal mengatur deleted_at untuk IdAsset2"
			res.Data = err.Error()
			return res, err
		}

		if dtAsset1.Id_asset_parent == dtAsset2.Id_asset_parent {
			var idAsetChild string
			queryParent := `SELECT id_asset_child FROM asset WHERE id_asset = ?`
			stmtParent, err := con.Prepare(queryParent)
			if err != nil {
				res.Status = 401
				res.Message = "stmt gagal"
				res.Data = err.Error()
				return res, err
			}
			defer stmtParent.Close()

			err = stmtParent.QueryRow(dtAsset1.Id_asset_parent).Scan(&idAsetChild)
			if err != nil {
				res.Status = 401
				res.Message = "exec gagal"
				res.Data = err.Error()
				return res, err
			}

			idAsetChild = idAsetChild + "," + strconv.Itoa(int(existingAssetId))
			queryUpdateParent := `UPDATE asset SET id_asset_child = ? WHERE id_asset = ?`
			stmtUpdateParent, err := con.Prepare(queryUpdateParent)
			if err != nil {
				res.Status = 401
				res.Message = "stmt gagal"
				res.Data = err.Error()
				return res, err
			}
			defer stmtUpdateParent.Close()

			_, err = stmtUpdateParent.Exec(idAsetChild, dtAsset1.Id_asset_parent)
			if err != nil {
				res.Status = 401
				res.Message = "exec gagal"
				res.Data = err.Error()
				return res, err
			}
		}

		// Return success response
		res.Status = http.StatusOK
		res.Message = "Berhasil memperbarui data"
		return res, nil
	}

	tempGabungIdJoin2 := tempIdJoin2 + "," + tempIdJoin1
	var existingAssetId2 int
	queryCheck2 := `SELECT id_asset FROM asset WHERE id_join = ?`
	err = con.QueryRow(queryCheck2, tempGabungIdJoin2).Scan(&existingAssetId2)
	if err == nil {
		queryRemoveDeletedAt := `UPDATE asset SET deleted_at = NULL WHERE id_asset = ?`
		_, err = con.Exec(queryRemoveDeletedAt, existingAssetId)
		if err != nil {
			res.Status = 401
			res.Message = "gagal menghapus deleted_at"
			res.Data = err.Error()
			return res, err
		}

		querySetDeletedAt := `UPDATE asset SET deleted_at = NOW() WHERE id_asset = ?`
		_, err = con.Exec(querySetDeletedAt, tempjoinAsset.IdAsset1)
		if err != nil {
			res.Status = 401
			res.Message = "gagal mengatur deleted_at untuk IdAsset1"
			res.Data = err.Error()
			return res, err
		}
		_, err = con.Exec(querySetDeletedAt, tempjoinAsset.IdAsset2)
		if err != nil {
			res.Status = 401
			res.Message = "gagal mengatur deleted_at untuk IdAsset2"
			res.Data = err.Error()
			return res, err
		}

		if dtAsset1.Id_asset_parent == dtAsset2.Id_asset_parent {
			var idAsetChild string
			queryParent := `SELECT id_asset_child FROM asset WHERE id_asset = ?`
			stmtParent, err := con.Prepare(queryParent)
			if err != nil {
				res.Status = 401
				res.Message = "stmt gagal"
				res.Data = err.Error()
				return res, err
			}
			defer stmtParent.Close()

			err = stmtParent.QueryRow(dtAsset1.Id_asset_parent).Scan(&idAsetChild)
			if err != nil {
				res.Status = 401
				res.Message = "exec gagal"
				res.Data = err.Error()
				return res, err
			}

			idAsetChild = idAsetChild + "," + strconv.Itoa(int(existingAssetId))
			queryUpdateParent := `UPDATE asset SET id_asset_child = ? WHERE id_asset = ?`
			stmtUpdateParent, err := con.Prepare(queryUpdateParent)
			if err != nil {
				res.Status = 401
				res.Message = "stmt gagal"
				res.Data = err.Error()
				return res, err
			}
			defer stmtUpdateParent.Close()

			_, err = stmtUpdateParent.Exec(idAsetChild, dtAsset1.Id_asset_parent)
			if err != nil {
				res.Status = 401
				res.Message = "exec gagal"
				res.Data = err.Error()
				return res, err
			}
		}

		// Return success response
		res.Status = http.StatusOK
		res.Message = "Berhasil memperbarui data"
		return res, nil
	}

	fmt.Println("gabung aset")

	luasBaru := dtAsset1.Luas + dtAsset2.Luas
	nilaiBaru := dtAsset1.Nilai + dtAsset2.Nilai
	if dtAsset1.Provinsi != dtAsset2.Provinsi {
		res.Status = 401
		res.Message = "provinsi tidak sama"
		return res, errors.New(res.Message)
	}
	if dtAsset1.Id_asset_parent != dtAsset2.Id_asset_parent {

		res.Status = 401
		res.Message = "parent tidak sama"
		return res, errors.New(res.Message)
	}
	tempIdJoin := strconv.Itoa(tempjoinAsset.IdAsset1) + "," + strconv.Itoa(tempjoinAsset.IdAsset2)

	var query string
	var result sql.Result

	query = `
		INSERT INTO asset (id_asset_parent,id_join, luas, nilai, provinsi, created_at) 
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

	result, err = stmt.Exec(dtAsset1.Id_asset_parent, tempIdJoin, luasBaru, nilaiBaru, dtAsset1.Provinsi)
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

	if dtAsset1.Id_asset_parent == dtAsset2.Id_asset_parent {
		var idAsetChild string
		queryParent := `SELECT id_asset_child FROM asset WHERE id_asset = ?`
		stmtParent, err := con.Prepare(queryParent)
		if err != nil {
			res.Status = 401
			res.Message = "stmt gagal"
			res.Data = err.Error()
			return res, err
		}
		defer stmtParent.Close()

		err = stmtParent.QueryRow(dtAsset1.Id_asset_parent).Scan(&idAsetChild)
		if err != nil {
			res.Status = 401
			res.Message = "exec gagal"
			res.Data = err.Error()
			return res, err
		}

		idAsetChild = idAsetChild + "," + strconv.Itoa(int(lastId))
		queryUpdateParent := `UPDATE asset SET id_asset_child = ? WHERE id_asset = ?`
		stmtUpdateParent, err := con.Prepare(queryUpdateParent)
		if err != nil {
			res.Status = 401
			res.Message = "stmt gagal"
			res.Data = err.Error()
			return res, err
		}
		defer stmtUpdateParent.Close()

		_, err = stmtUpdateParent.Exec(idAsetChild, dtAsset1.Id_asset_parent)
		if err != nil {
			res.Status = 401
			res.Message = "exec gagal"
			res.Data = err.Error()
			return res, err
		}
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
	type TransactionRequestCustom struct {
		Id_transaksi_jual_sewa int    `json:"id_transaksi_jual_sewa"`
		Perusahaan_id          int    `json:"perusahaan_id"`
		Lokasi_perusahaan      string `json:"lokasi_perusahaan"`
		User_id                int    `json:"user_id"`
		Username               string `json:"username"`
		Nama_lengkap           string `json:"nama_lengkap"`
		Id_asset               int    `json:"id_asset"`
		Nama_aset              string `json:"nama_aset"`
		Alamat_aset            string `json:"alamat"`
		Status                 string `json:"status"`
		Nama_progress          string `json:"nama_progress"`
		Proposal               string `json:"proposal"`
		Tgl_meeting            string `json:"tgl_meeting"`
		Waktu_meeting          string `json:"waktu_meeting"`
		Lokasi_meeting         string `json:"lokasi_meeting"`
		Deskripsi              string `json:"deskripsi"`
		Alasan                 string `json:"alasan"`
		Tgl_dateline           string `json:"tgl_dateline"`
		Created_at             string `json:"created_at"`
	}
	// asset + transaction_request + progress
	var arrAsetTranReq = []TransactionRequestCustom{}

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	query := `
	SELECT tr.id_transaksi_jual_sewa,tr.perusahaan_id,tr.user_id,tr.id_asset,
		a.nama, a.alamat, tr.status, tr.nama_progress,tr.proposal,tr.tgl_meeting,tr.waktu_meeting,
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
		var dtTranReq TransactionRequestCustom
		err := rows.Scan(&dtTranReq.Id_transaksi_jual_sewa, &dtTranReq.Perusahaan_id,
			&dtTranReq.User_id, &dtTranReq.Id_asset, &dtTranReq.Nama_aset, &dtTranReq.Alamat_aset, &dtTranReq.Status, &dtTranReq.Nama_progress,
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
	alamat, kondisi, koordinat, batas_koordinat, luas, nilai, provinsi string) (Response, error) {
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
	tagQuery := "DELETE FROM `asset_tags` WHERE id_asset = ?"
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

	fmt.Println("query update")
	query := `
	UPDATE asset 
	SET nama = ?, tipe = ?, nomor_legalitas = ?, status_asset = ?, surat_legalitas = ?, alamat = ?, kondisi = ?, titik_koordinat = ?,
	batas_koordinat = ?, luas = ?, nilai = ?, provinsi = ?
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
	_, err = stmt.Exec(
		nama, tipe, nomor_legalitas, status, surat_legalitas, alamat, kondisi, koordinat, batas_koordinat, luas, nilai, provinsi, id_asset,
	)
	if err != nil {
		res.Status = 401
		res.Message = "exec gagal"
		res.Data = err.Error()
		return res, err
	}

	fmt.Println("tambah usage")

	// tambah usage + tags
	for _, usageId := range usageIds {
		usageQuery := "INSERT INTO asset_penggunaan (id_asset, id_penggunaan) VALUES (?, ?)"
		fmt.Println("usage id: ", id_asset, usageId)
		_, err = con.Exec(usageQuery, id_asset, usageId)
		if err != nil {
			res.Status = 401
			res.Message = "Gagal menambah penggunaan"
			res.Data = err.Error()
			return res, err
		}
	}
	fmt.Println("tambah tags")

	for _, tagId := range tagIds {
		tagQuery := "INSERT INTO asset_tags (id_asset, id_tags) VALUES (?, ?)"
		_, err = con.Exec(tagQuery, id_asset, tagId)
		if err != nil {
			res.Status = 401
			res.Message = "Gagal menambah tag"
			res.Data = err.Error()
			return res, err
		}
	}
	fmt.Println("tambah file legalitas")

	// tambah filelegalitas
	//source
	srclegalitas, err := filelegalitas.Open()
	if err != nil {
		log.Println(err.Error())
		return res, err
	}
	defer srclegalitas.Close()

	// Destination
	filelegalitas.Filename = id_asset + "_" + filelegalitas.Filename
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

	lastId, err := strconv.Atoi(id_asset)
	if err != nil {
		res.Status = 401
		res.Message = "Gagal konversi string"
		res.Data = err.Error()
		return res, err
	}

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
	suratkuasa.Filename = id_asset + "_" + suratkuasa.Filename
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

	tempaset, err := GetAssetById(id_asset)
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

func UpdateAssetById(filelegalitas *multipart.FileHeader, suratkuasa *multipart.FileHeader,
	id_asset, nama, surat_legalitas, tipe, usage, tag, nomor_legalitas, status,
	alamat, kondisi, koordinat, batas_koordinat, luas, nilai, provinsi string, files []*multipart.FileHeader) (Response, error) {
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

	// hapus file legalitas + surat kuasa kalau ada + gambar aset
	linkasetQueue := `
	SELECT a.file_legalitas, a.surat_kuasa, IFNULL(ag.link_gambar,'')
	FROM asset a
	LEFT JOIN asset_gambar ag ON a.id_asset = ag.id_asset_gambar
	WHERE a.id_asset = ?
`
	var linkfilelegalitas string
	var linksuratkuasa string
	var linkGambar string
	var existingImages []string

	rows, err := con.Query(linkasetQueue, id_asset)
	if err != nil {
		res.Status = 401
		res.Message = "exec gagal"
		res.Data = err.Error()
		return res, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&linkfilelegalitas, &linksuratkuasa, &linkGambar)
		if err != nil {
			res.Status = 401
			res.Message = "exec gagal"
			res.Data = err.Error()
			return res, err
		}

		// Collect existing image links
		if linkGambar != "" {
			existingImages = append(existingImages, linkGambar)
		}
	}
	if linkfilelegalitas != "" {
		err = os.Remove(linkfilelegalitas)
		if err != nil {
			return res, err
		}

	}
	if linksuratkuasa != "" {
		err = os.Remove(linksuratkuasa)
		if err != nil {
			return res, err
		}
	}

	for _, deleteFile := range existingImages {
		err = os.Remove(deleteFile)
		if err != nil {
			return res, err
		}
	}
	deleteImagesQuery := "DELETE FROM asset_gambar WHERE id_asset_gambar = ?"
	_, err = con.Exec(deleteImagesQuery, id_asset)
	if err != nil {
		res.Status = 401
		res.Message = "Delete gambar exec gagal"
		res.Data = "Delete gambar exec gagal dengan id aset " + id_asset
		return res, err
	}

	// hapus tags dan usage di db
	tagQuery := "DELETE FROM `asset_tags` WHERE id_asset = ?"
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

	fmt.Println("query update")
	query := `
	UPDATE asset 
	SET nama = ?, tipe = ?, nomor_legalitas = ?, status_asset = ?, surat_legalitas = ?, alamat = ?, kondisi = ?, titik_koordinat = ?,
	batas_koordinat = ?, luas = ?, nilai = ?, provinsi = ?
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
	_, err = stmt.Exec(
		nama, tipe, nomor_legalitas, status, surat_legalitas, alamat, kondisi, koordinat, batas_koordinat, luas, nilai, provinsi, id_asset,
	)
	if err != nil {
		res.Status = 401
		res.Message = "exec gagal"
		res.Data = err.Error()
		return res, err
	}

	// insert gambar

	// ambil foto aset dan hapus
	filePaths := []string{}
	for _, fileHeader := range files {
		filePath, err := updateAssetGambar(id_asset, fileHeader)
		if err != nil {
			res.Status = 500
			res.Message = "Error saving file"
			res.Data = err.Error()
			return res, err
		}
		filePaths = append(filePaths, filePath)
	}

	for _, filePath := range filePaths {
		// Prepare the INSERT statement for each filePath
		insertImageQuery := "INSERT INTO asset_gambar (id_asset_gambar, link_gambar) VALUES (?, ?)"
		_, err = con.Exec(insertImageQuery, id_asset, filePath)
		if err != nil {
			// Handle the error if the insert fails
			res.Status = 500
			res.Message = "Error inserting image path"
			res.Data = err.Error()
			return res, err
		}
	}
	fmt.Println("tambah usage")

	// tambah usage + tags
	for _, usageId := range usageIds {
		usageQuery := "INSERT INTO asset_penggunaan (id_asset, id_penggunaan) VALUES (?, ?)"
		fmt.Println("usage id: ", id_asset, usageId)
		_, err = con.Exec(usageQuery, id_asset, usageId)
		if err != nil {
			res.Status = 401
			res.Message = "Gagal menambah penggunaan"
			res.Data = err.Error()
			return res, err
		}
	}
	fmt.Println("tambah tags")

	for _, tagId := range tagIds {
		tagQuery := "INSERT INTO asset_tags (id_asset, id_tags) VALUES (?, ?)"
		_, err = con.Exec(tagQuery, id_asset, tagId)
		if err != nil {
			res.Status = 401
			res.Message = "Gagal menambah tag"
			res.Data = err.Error()
			return res, err
		}
	}
	fmt.Println("tambah file legalitas")

	// tambah filelegalitas
	//source
	srclegalitas, err := filelegalitas.Open()
	if err != nil {
		log.Println(err.Error())
		return res, err
	}
	defer srclegalitas.Close()

	// Destination
	filelegalitas.Filename = id_asset + "_" + filelegalitas.Filename
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

	lastId, err := strconv.Atoi(id_asset)
	if err != nil {
		res.Status = 401
		res.Message = "Gagal konversi string"
		res.Data = err.Error()
		return res, err
	}

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
	suratkuasa.Filename = id_asset + "_" + suratkuasa.Filename
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

	tempaset, err := GetAssetById(id_asset)
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

func updateAssetGambar(idAset string, fileHeader *multipart.FileHeader) (string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	uploadDir := "uploads/asset/foto/"
	fileHeader.Filename = idAset + "_" + fileHeader.Filename
	os.MkdirAll(uploadDir, os.ModePerm)
	filePath := filepath.Join(uploadDir, fileHeader.Filename)

	out, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		return "", err
	}

	return filePath, nil
}

func UnjoinAsset(asetId string) (Response, error) {
	var res Response

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	var dtAsset Asset

	queryAsset1 := "SELECT IFNULL(id_asset_parent,0),id_join FROM asset WHERE id_asset = ?"
	stmtAsset1, err := con.Prepare(queryAsset1)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmtAsset1.Close()

	err = stmtAsset1.QueryRow(asetId).Scan(&dtAsset.Id_asset_parent, &dtAsset.Id_join)
	if err != nil {
		res.Status = 401
		res.Message = "exec gagal"
		res.Data = err.Error()
		return res, err
	}

	fmt.Println("update asset")
	// update asset
	query := `
	UPDATE asset 
	SET deleted_at = NOW()
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
	_, err = stmt.Exec(asetId)
	if err != nil {
		res.Status = 401
		res.Message = "exec gagal"
		res.Data = err.Error()
		return res, err
	}

	// update asset join
	if dtAsset.Id_join != "" {
		joinIds := strings.Split(dtAsset.Id_join, ",")
		for _, joinId := range joinIds {
			trimmedJoinId := strings.TrimSpace(joinId)

			queryUpdateJoin := `
			UPDATE asset 
			SET deleted_at = NULL
			WHERE id_asset = ?
			`
			stmtUpdateJoin, err := con.Prepare(queryUpdateJoin)
			if err != nil {
				res.Status = 401
				res.Message = "stmt gagal"
				res.Data = err.Error()
				return res, err
			}
			defer stmtUpdateJoin.Close()

			_, err = stmtUpdateJoin.Exec(trimmedJoinId)
			if err != nil {
				res.Status = 401
				res.Message = "exec gagal untuk id join " + trimmedJoinId
				res.Data = err.Error()
				return res, err
			}
		}
	}

	// update asset parent kalau ada
	if dtAsset.Id_asset_parent != 0 {
		queryParent := "SELECT id_asset_child FROM asset WHERE id_asset = ?"
		var idAssetChild string
		err = con.QueryRow(queryParent, dtAsset.Id_asset_parent).Scan(&idAssetChild)
		if err != nil {
			res.Status = 401
			res.Message = "exec gagal mendapatkan id_asset_child"
			res.Data = err.Error()
			return res, err
		}

		childIds := strings.Split(idAssetChild, ",")
		var updatedChildIds []string
		for _, childId := range childIds {
			if strings.TrimSpace(childId) != asetId {
				updatedChildIds = append(updatedChildIds, strings.TrimSpace(childId))
			}
		}

		newIdAssetChild := strings.Join(updatedChildIds, ",")
		queryUpdateParent := `
		UPDATE asset 
		SET id_asset_child = ?
		WHERE id_asset = ?
		`
		stmtUpdateParent, err := con.Prepare(queryUpdateParent)
		if err != nil {
			res.Status = 401
			res.Message = "stmt gagal update id_asset_child"
			res.Data = err.Error()
			return res, err
		}
		defer stmtUpdateParent.Close()

		_, err = stmtUpdateParent.Exec(newIdAssetChild, dtAsset.Id_asset_parent)
		if err != nil {
			res.Status = 401
			res.Message = "exec gagal update id_asset_child"
			res.Data = err.Error()
			return res, err
		}
	}

	var tempaset Response
	nId, _ := strconv.Atoi(asetId)
	tempaset, _ = GetAssetById(strconv.Itoa(nId))

	res.Status = http.StatusOK
	res.Message = "Berhasil unjoin asset"
	res.Data = tempaset.Data

	defer db.DbClose(con)
	return res, nil
}

func FilterAsset(input string) (Response, error) {
	var res Response

	type InputFilter struct {
		Tipe          string `json:"type"`
		Status        string `json:"status"`
		Tagasset      string `json:"tag"`
		Provinsiasset string `json:"provinsi"`
	}

	var tempInput InputFilter
	err := json.Unmarshal([]byte(input), &tempInput)
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

	query := `SELECT * FROM asset WHERE deleted_at IS NULL`
	params := []interface{}{}
	if tempInput.Tipe != "" {
		tipeList := strings.Split(tempInput.Tipe, ",")
		query += " AND tipe IN (?" + strings.Repeat(",?", len(tipeList)-1) + ")"
		for _, t := range tipeList {
			params = append(params, strings.TrimSpace(t))
		}
	}
	if tempInput.Status != "" {
		statusList := strings.Split(tempInput.Status, ",")
		query += " AND status_asset IN (?" + strings.Repeat(",?", len(statusList)-1) + ")"
		for _, s := range statusList {
			params = append(params, strings.TrimSpace(s))
		}
	}
	if tempInput.Tagasset != "" {
		tagList := strings.Split(tempInput.Tagasset, ",")
		query += `
		AND id_asset IN (
			SELECT id_asset 
			FROM asset_tags 
			WHERE id_tags IN ( ` + strings.Repeat("?,", len(tagList)-1) + "?))"
		for _, tag := range tagList {
			params = append(params, strings.TrimSpace(tag))
		}
	}
	if tempInput.Provinsiasset != "" {
		provinsiList := strings.Split(tempInput.Provinsiasset, ",")
		query += " AND provinsi IN (?" + strings.Repeat(",?", len(provinsiList)-1) + ")"
		for _, p := range provinsiList {
			params = append(params, strings.TrimSpace(p))
		}
	}

	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(params...)
	if err != nil {
		res.Status = 401
		res.Message = "exec gagal"
		res.Data = err.Error()
		return res, err
	}
	defer rows.Close()

	// Process the results
	var assets []Asset
	for rows.Next() {
		var dtAset Asset
		var masaSewa []byte
		var deleteAt []byte
		var idJoin, idAssetChild sql.NullString
		var idAssetParent, idProvinsi sql.NullInt32
		err := rows.Scan(&dtAset.Id_asset, &idAssetParent, &idAssetChild, &idJoin, &dtAset.Nama, &dtAset.Tipe, &dtAset.Nomor_legalitas, &dtAset.File_legalitas, &dtAset.Status_asset, &dtAset.Surat_kuasa, &dtAset.Surat_legalitas, &dtAset.Alamat, &dtAset.Kondisi, &dtAset.Titik_koordinat, &dtAset.Batas_koordinat, &dtAset.Luas, &dtAset.Nilai, &idProvinsi, &dtAset.Status_pengecekan, &dtAset.Status_verifikasi, &dtAset.Status_publik, &dtAset.Hak_akses, &masaSewa, &dtAset.Created_at, &deleteAt) // Add appropriate fields here
		if err != nil {
			res.Status = 401
			res.Message = "scan gagal"
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
		if idProvinsi.Valid {
			dtAset.Provinsi = int(idProvinsi.Int32)
		} else {
			dtAset.Provinsi = 0
		}

		gambarQuery := "SELECT link_gambar FROM asset_gambar WHERE id_asset_gambar = ?"
		rows, err := con.Query(gambarQuery, dtAset.Id_asset)
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

		assets = append(assets, dtAset)
	}

	if err = rows.Err(); err != nil {
		res.Status = 401
		res.Message = "rows error"
		res.Data = err.Error()
		return res, err
	}

	res.Status = http.StatusOK
	res.Message = "Berhasil mengambil data asset berdasarkan filter"
	if len(assets) == 0 {
		res.Data = []Asset{}
	} else {
		res.Data = assets
	}

	defer db.DbClose(con)
	return res, nil
}
