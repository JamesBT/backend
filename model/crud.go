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
)

// ini untuk crud 10 tabel aja
// CRUD aset ============================================================================
func CreateAsset(filelegalitas *multipart.FileHeader, suratkuasa *multipart.FileHeader, nama, perusahaan_id, tipe, nomorlegalitas, status, alamat, kondisi, koordinat, batas_koordinat, luas, nilai string) (Response, error) {
	var res Response

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	query := "INSERT INTO asset (perusahaan_id, nama, tipe, nomor_legalitas, status_asset, alamat, kondisi, titik_koordinat, batas_koordinat, luas, nilai, created_at) VALUES (?,?,?,?,?,?,?,?,?,?,?,NOW())"
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(perusahaan_id, nama, tipe, nomorlegalitas, status, alamat, kondisi, koordinat, batas_koordinat, luas, nilai)
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
	filelegalitas.Filename = _tempid + "_" + perusahaan_id + "_" + filelegalitas.Filename
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
	suratkuasa.Filename = _tempid + "_" + perusahaan_id + "_" + suratkuasa.Filename
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

	tempaset, _ := GetAssetById(string(tempid))

	res.Status = http.StatusOK
	res.Message = "Berhasil memasukkan data"
	res.Data = tempaset

	defer db.DbClose(con)
	return res, nil
}

func CreateAssetChild(filelegalitas *multipart.FileHeader, suratkuasa *multipart.FileHeader, parent_id, nama, perusahaan_id, tipe, nomorlegalitas, status, alamat, kondisi, koordinat, batas_koordinat, luas, nilai string) (Response, error) {
	var res Response

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	query := "INSERT INTO asset (id_asset_parent, perusahaan_id, nama, tipe, nomor_legalitas, status_asset, alamat, kondisi, titik_koordinat, batas_koordinat, luas, nilai, created_at) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,NOW())"
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(parent_id, perusahaan_id, nama, tipe, nomorlegalitas, status, alamat, kondisi, koordinat, batas_koordinat, luas, nilai)
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
	// tambah filelegalitas
	//source
	srclegalitas, err := filelegalitas.Open()
	if err != nil {
		log.Println(err.Error())
		return res, err
	}
	defer srclegalitas.Close()

	// Destination

	tempid := int(lastId)
	_tempid := strconv.Itoa(tempid)
	filelegalitas.Filename = _tempid + "_" + perusahaan_id + "_" + filelegalitas.Filename
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
	suratkuasa.Filename = _tempid + "_" + perusahaan_id + "_" + suratkuasa.Filename
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
	queryupdate := "UPDATE asset SET id_asset_child = ? WHERE id_asset = ?"
	stmtupdate, err := con.Prepare(queryupdate)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmtupdate.Close()

	_, err = stmtupdate.Exec(templastid, parent_id)
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
	var masaSewa sql.NullTime
	var deleteAt sql.NullTime
	var idAssetParent, idAssetChild, idJoin, idPerusahaan sql.NullInt32

	for result.Next() {
		err = result.Scan(&dtAset.Id_asset, &idAssetParent, &idAssetChild, &idJoin, &idPerusahaan, &dtAset.Nama, &dtAset.Tipe, &dtAset.Nomor_legalitas, &dtAset.File_legalitas, &dtAset.Status_asset, &dtAset.Surat_kuasa, &dtAset.Alamat, &dtAset.Kondisi, &dtAset.Titik_koordinat, &dtAset.Batas_koordinat, &dtAset.Luas, &dtAset.Nilai, &dtAset.Usage, &dtAset.Status_pengecekan, &dtAset.Status_verifikasi, &dtAset.Status_publik, &dtAset.Hak_akses, &masaSewa, &dtAset.Created_at, &deleteAt)
		if err != nil {
			res.Status = 401
			res.Message = "rows scan"
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
		if idAssetParent.Valid {
			dtAset.Id_asset_parent = int(idAssetParent.Int32)
		} else {
			dtAset.Id_asset_parent = 0
		}
		if idAssetChild.Valid {
			dtAset.Id_asset_child = strconv.Itoa(int(idAssetChild.Int32))
		} else {
			dtAset.Id_asset_child = ""
		}
		if idJoin.Valid {
			dtAset.Id_join = int(idJoin.Int32)
		} else {
			dtAset.Id_join = 0
		}
		if idPerusahaan.Valid {
			dtAset.Id_perusahaan = int(idPerusahaan.Int32)
		} else {
			dtAset.Id_perusahaan = 0
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
	var masaSewa sql.NullTime
	var deleteAt sql.NullTime
	var idAssetParent, idAssetChild, idJoin sql.NullInt32
	err = stmt.QueryRow(nId).Scan(&dtAset.Id_asset, &idAssetParent, &idAssetChild, &idJoin, &dtAset.Id_perusahaan, &dtAset.Nama, &dtAset.Tipe, &dtAset.Nomor_legalitas, &dtAset.File_legalitas, &dtAset.Status_asset, &dtAset.Surat_kuasa, &dtAset.Alamat, &dtAset.Kondisi, &dtAset.Titik_koordinat, &dtAset.Batas_koordinat, &dtAset.Luas, &dtAset.Nilai, &dtAset.Usage, &dtAset.Status_pengecekan, &dtAset.Status_verifikasi, &dtAset.Status_publik, &dtAset.Hak_akses, &masaSewa, &dtAset.Created_at, &deleteAt)
	if err != nil {
		res.Status = 401
		res.Message = "exec gagal"
		res.Data = err.Error()
		return res, err
	}
	fmt.Println("ambil berhasil")

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
	if idAssetParent.Valid {
		dtAset.Id_asset_parent = int(idAssetParent.Int32)
	} else {
		dtAset.Id_asset_parent = 0
	}
	if idAssetChild.Valid {
		dtAset.Id_asset_child = strconv.Itoa(int(idAssetChild.Int32))
	} else {
		dtAset.Id_asset_child = ""
	}
	if idJoin.Valid {
		dtAset.Id_join = int(idJoin.Int32)
	} else {
		dtAset.Id_join = 0
	}

	res.Status = http.StatusOK
	res.Message = "Berhasil mengambil data"
	res.Data = dtAset

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
	fmt.Println(aset_id)
	query := "SELECT id_asset, id_asset_parent, id_asset_child, id_join, perusahaan_id, nama, tipe, nomor_legalitas, file_legalitas, status_asset, surat_kuasa, alamat, kondisi, titik_koordinat, batas_koordinat, luas, nilai, `usage`, status_pengecekan, status_verifikasi, status_publik, hak_akses, masa_sewa, created_at, deleted_at FROM asset WHERE id_asset = ?"
	stmt, err := con.Prepare(query)
	if err != nil {
		return dtAset, err
	}
	defer stmt.Close()

	nId, _ := strconv.Atoi(aset_id)
	var masaSewa sql.NullTime
	var deleteAt sql.NullTime
	var idAssetParent, idAssetChild, idJoin sql.NullInt32

	err = stmt.QueryRow(nId).Scan(
		&dtAset.Id_asset, &idAssetParent, &idAssetChild, &idJoin, &dtAset.Id_perusahaan,
		&dtAset.Nama, &dtAset.Tipe, &dtAset.Nomor_legalitas, &dtAset.File_legalitas, &dtAset.Status_asset,
		&dtAset.Surat_kuasa, &dtAset.Alamat, &dtAset.Kondisi, &dtAset.Titik_koordinat, &dtAset.Batas_koordinat,
		&dtAset.Luas, &dtAset.Nilai, &dtAset.Usage, &dtAset.Status_pengecekan, &dtAset.Status_verifikasi, &dtAset.Status_publik, &dtAset.Hak_akses,
		&masaSewa, &dtAset.Created_at, &deleteAt,
	)
	if err != nil {
		return dtAset, err
	}

	// Format nullable times
	if masaSewa.Valid {
		dtAset.Masa_sewa = masaSewa.Time.Format("2006-01-02")
	} else {
		dtAset.Masa_sewa = ""
	}
	if deleteAt.Valid {
		dtAset.Deleted_at = deleteAt.Time.Format("2006-01-02")
	} else {
		dtAset.Deleted_at = ""
	}
	if idAssetParent.Valid {
		dtAset.Id_asset_parent = int(idAssetParent.Int32)
	} else {
		dtAset.Id_asset_parent = 0
	}
	if idAssetChild.Valid {
		dtAset.Id_asset_child = strconv.Itoa(int(idAssetChild.Int32))
	} else {
		dtAset.Id_asset_child = ""
	}
	if idJoin.Valid {
		dtAset.Id_join = int(idJoin.Int32)
	} else {
		dtAset.Id_join = 0
	}

	childQuery := "SELECT id_asset FROM asset WHERE id_asset_parent = ?"
	rows, err := con.Query(childQuery, dtAset.Id_asset)
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

	// untuk gambar
	imageQuery := "SELECT link_gambar FROM asset_gambar WHERE id_asset_gambar = ?"
	imageRows, err := con.Query(imageQuery, dtAset.Id_asset)
	if err != nil {
		return dtAset, err
	}
	defer imageRows.Close()

	for imageRows.Next() {
		var linkGambar string
		if err := imageRows.Scan(&linkGambar); err != nil {
			return dtAset, err
		}
		dtAset.LinkGambar = append(dtAset.LinkGambar, linkGambar)
	}

	// untuk tags
	tagQuery := `SELECT t.nama FROM asset_tags at
		JOIN tags t ON at.id_tags = t.id
		WHERE at.id_asset = ?`
	tagRows, err := con.Query(tagQuery, dtAset.Id_asset)
	if err != nil {
		return dtAset, err
	}
	defer tagRows.Close()

	for tagRows.Next() {
		var tagName string
		if err := tagRows.Scan(&tagName); err != nil {
			return dtAset, err
		}
		dtAset.TagsAssets = append(dtAset.TagsAssets, tagName)
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

/**
func GetAllPerusahaan() (Response, error) {
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

	query := "SELECT * FROM perusahaan"
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
		err = result.Scan(&dtPerusahaan.Perusahaan_id, &dtPerusahaan.User_id, &dtPerusahaan.Sertifikat_perusahaan)
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

func GetPerusahaanById(perusahaan_id string) (Response, error) {
	var res Response
	var dtPerusahaan Perusahaan

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	query := "SELECT * FROM perusahaan WHERE perusahaan_id = ?"
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()
	nId, _ := strconv.Atoi(perusahaan_id)
	err = stmt.QueryRow(nId).Scan(&dtPerusahaan.Perusahaan_id, &dtPerusahaan.User_id, &dtPerusahaan.Sertifikat_perusahaan)
	if err != nil {
		res.Status = 401
		res.Message = "exec gagal"
		res.Data = err.Error()
		return res, err
	}

	res.Status = http.StatusOK
	res.Message = "Berhasil mengambil data"
	res.Data = dtPerusahaan

	defer db.DbClose(con)
	return res, nil
}

func UpdatePerusahaanById(perusahaan string) (Response, error) {
	var res Response

	var dtPerusahaan = Perusahaan{}

	err := json.Unmarshal([]byte(perusahaan), &dtPerusahaan)
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

	query := "UPDATE perusahaan SET sertifikat_perusahaan = ? WHERE perusahaan_id = ?"
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(dtPerusahaan.Sertifikat_perusahaan, dtPerusahaan.Perusahaan_id)
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

func DeletePerusahaanById(perusahaan string) (Response, error) {
	var res Response

	var dtPerusahaan = Perusahaan{}

	err := json.Unmarshal([]byte(perusahaan), &dtPerusahaan)
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

	query := "DELETE FROM perusahaan WHERE perusahaan_id = ?"
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(dtPerusahaan.Perusahaan_id)
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
**/
// CRUD privilege ============================================================================
func CreatePrivilege(hakakses string) (Response, error) {
	var res Response
	var dtPrivilege = Privilege{}

	err := json.Unmarshal([]byte(hakakses), &dtPrivilege)
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

	query := "INSERT INTO privilege (nama_privilege) VALUES (?)"
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(dtPrivilege.Nama_privilege)
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
	dtPrivilege.Privilege_id = int(lastId)

	res.Status = http.StatusOK
	res.Message = "Berhasil memasukkan data"
	res.Data = dtPrivilege

	defer db.DbClose(con)
	return res, nil
}

func GetAllPrivilege() (Response, error) {
	var res Response
	var arrPrivilege = []Privilege{}
	var dtPrivilege Privilege

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	query := "SELECT * FROM privilege"
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
		err = result.Scan(&dtPrivilege.Privilege_id, &dtPrivilege.Nama_privilege)
		if err != nil {
			res.Status = 401
			res.Message = "rows scan"
			res.Data = err.Error()
			return res, err
		}
		arrPrivilege = append(arrPrivilege, dtPrivilege)
	}

	res.Status = http.StatusOK
	res.Message = "Berhasil mengambil data"
	res.Data = arrPrivilege

	defer db.DbClose(con)
	return res, nil
}

func GetPrivilegeById(hak_akses_id string) (Response, error) {
	var res Response
	var dtPrivilege Privilege

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	query := "SELECT * FROM privilege WHERE privilege_id = ?"
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()
	nId, _ := strconv.Atoi(hak_akses_id)
	err = stmt.QueryRow(nId).Scan(&dtPrivilege.Privilege_id, &dtPrivilege.Nama_privilege)
	if err != nil {
		res.Status = 401
		res.Message = "exec gagal"
		res.Data = err.Error()
		return res, err
	}

	res.Status = http.StatusOK
	res.Message = "Berhasil mengambil data"
	res.Data = dtPrivilege

	defer db.DbClose(con)
	return res, nil
}

func GetPrivilegeByName(nama_privilege string) (Response, error) {
	var res Response
	var dtPrivileges = []Privilege{}

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	query := "SELECT * FROM asset WHERE nama_privilege LIKE ?"
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	rows, err := stmt.Query("%" + nama_privilege + "%")
	if err != nil {
		res.Status = 401
		res.Message = "exec gagal"
		res.Data = err.Error()
		return res, err
	}
	defer rows.Close()
	for rows.Next() {
		var dtPrivilege Privilege
		err := rows.Scan(&dtPrivilege.Privilege_id, &dtPrivilege.Nama_privilege)
		if err != nil {
			res.Status = 401
			res.Message = "scan gagal"
			res.Data = err.Error()
			return res, err
		}
		dtPrivileges = append(dtPrivileges, dtPrivilege)
	}

	if err = rows.Err(); err != nil {
		res.Status = 401
		res.Message = "rows error"
		res.Data = err.Error()
		return res, err
	}

	if len(dtPrivileges) == 0 {
		res.Status = 404
		res.Message = "Data tidak ditemukan"
		res.Data = nil
		return res, nil
	}

	res.Status = http.StatusOK
	res.Message = "Berhasil mengambil data"
	res.Data = dtPrivileges

	defer db.DbClose(con)
	return res, nil
}

func UpdatePrivilegeById(hakakses string) (Response, error) {
	var res Response

	var dtPrivilege = Privilege{}

	err := json.Unmarshal([]byte(hakakses), &dtPrivilege)
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

	query := "UPDATE privilege SET nama_privilege = ? WHERE privilege_id = ?"
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(dtPrivilege.Nama_privilege, dtPrivilege.Privilege_id)
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

func DeletePrivilegeById(hakakses string) (Response, error) {
	var res Response

	var dtPrivilege = Privilege{}

	err := json.Unmarshal([]byte(hakakses), &dtPrivilege)
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

	query := "DELETE FROM privilege WHERE privilege_id = ?"
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(dtPrivilege.Privilege_id)
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

func GetAllPerusahaanDetailed() (Response, error) {
	var res Response
	var dtPerusahaan = []UserPerusahaan{}

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	query := "SELECT p.perusahaan_id, p.name, COUNT(DISTINCT up.id_user) AS user_count, COUNT(DISTINCT tr.id_transaksi_jual_sewa) AS transaction_count FROM  perusahaan p LEFT JOIN user_perusahaan up ON p.perusahaan_id = up.id_perusahaan LEFT JOIN transaction_request tr ON p.perusahaan_id = tr.perusahaan_id GROUP BY p.perusahaan_id;"
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
		var _dtUserPerusahaan UserPerusahaan
		err := rows.Scan(&_dtUserPerusahaan.Perusahaan_id, &_dtUserPerusahaan.Name, &_dtUserPerusahaan.UserCount, &_dtUserPerusahaan.TransactionCount)
		if err != nil {
			res.Status = 401
			res.Message = "scan gagal"
			res.Data = err.Error()
			return res, err
		}

		dtPerusahaan = append(dtPerusahaan, _dtUserPerusahaan)
	}

	if err = rows.Err(); err != nil {
		res.Status = 401
		res.Message = "rows error"
		res.Data = err.Error()
		return res, err
	}

	if len(dtPerusahaan) == 0 {
		res.Status = 404
		res.Message = "Data tidak ditemukan"
		res.Data = nil
		return res, nil
	}

	res.Status = http.StatusOK
	res.Message = "Berhasil mengambil data"
	res.Data = dtPerusahaan

	defer db.DbClose(con)
	return res, nil
}

// CRUD role ============================================================================
func CreateRole(peran string) (Response, error) {
	var res Response
	var dtRole = Role{}

	err := json.Unmarshal([]byte(peran), &dtRole)
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

	query := "INSERT INTO role (nama_role) VALUES (?)"
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(dtRole.Nama_role)
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
	dtRole.Role_id = int(lastId)

	res.Status = http.StatusOK
	res.Message = "Berhasil memasukkan data"
	res.Data = dtRole

	defer db.DbClose(con)
	return res, nil
}

func GetAllRole() (Response, error) {
	var res Response
	var arrRole = []Role{}
	var dtRole Role

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	query := "SELECT * FROM role"
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
		err = result.Scan(&dtRole.Role_id, &dtRole.Nama_role)
		if err != nil {
			res.Status = 401
			res.Message = "rows scan"
			res.Data = err.Error()
			return res, err
		}
		arrRole = append(arrRole, dtRole)
	}

	res.Status = http.StatusOK
	res.Message = "Berhasil mengambil data"
	res.Data = arrRole

	defer db.DbClose(con)
	return res, nil
}

func GetRoleById(role_id string) (Response, error) {
	var res Response
	var dtRole Role

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	query := "SELECT * FROM role WHERE role_id = ?"
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()
	nId, _ := strconv.Atoi(role_id)
	err = stmt.QueryRow(nId).Scan(&dtRole.Role_id, &dtRole.Nama_role)
	if err != nil {
		res.Status = 401
		res.Message = "exec gagal"
		res.Data = err.Error()
		return res, err
	}

	res.Status = http.StatusOK
	res.Message = "Berhasil mengambil data"
	res.Data = dtRole

	defer db.DbClose(con)
	return res, nil
}

func GetRoleByName(nama_role string) (Response, error) {
	var res Response
	var dtRoles = []Role{}

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	query := "SELECT * FROM role WHERE nama_role LIKE ?"
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	rows, err := stmt.Query("%" + nama_role + "%")
	if err != nil {
		res.Status = 401
		res.Message = "exec gagal"
		res.Data = err.Error()
		return res, err
	}
	defer rows.Close()
	for rows.Next() {
		var dtRole Role
		err := rows.Scan(&dtRole.Role_id, &dtRole.Nama_role)
		if err != nil {
			res.Status = 401
			res.Message = "scan gagal"
			res.Data = err.Error()
			return res, err
		}
		dtRoles = append(dtRoles, dtRole)
	}

	if err = rows.Err(); err != nil {
		res.Status = 401
		res.Message = "rows error"
		res.Data = err.Error()
		return res, err
	}

	if len(dtRoles) == 0 {
		res.Status = 404
		res.Message = "Data tidak ditemukan"
		res.Data = nil
		return res, nil
	}

	res.Status = http.StatusOK
	res.Message = "Berhasil mengambil data"
	res.Data = dtRoles

	defer db.DbClose(con)
	return res, nil
}

func UpdateRoleById(peran string) (Response, error) {
	var res Response

	var dtRole = Role{}

	err := json.Unmarshal([]byte(peran), &dtRole)
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

	query := "UPDATE role SET nama_role = ? WHERE role_id = ?"
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(dtRole.Nama_role, dtRole.Role_id)
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

func DeleteRoleById(peran string) (Response, error) {
	var res Response

	var dtRole = Role{}

	err := json.Unmarshal([]byte(peran), &dtRole)
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

	query := "DELETE FROM role WHERE role_id = ?"
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(dtRole.Role_id)
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

// CRUD surveyor ============================================================================
func LoginSurveyor(akun string) (Response, error) {
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

	// cek apakah password benar atau tidak
	queryinsert := "SELECT u.user_id, u.username, u.nama_lengkap, u.alamat, u.jenis_kelamin, u.tanggal_lahir, u.email, u.nomor_telepon, u.foto_profil, u.ktp, ud.user_kelas_id, ud.status, ud.tipe, ud.first_login, ud.denied_by_admin, s.lokasi, s.availability_surveyor FROM user u JOIN user_detail ud ON u.user_id = ud.user_detail_id JOIN surveyor s ON u.user_id = s.user_id WHERE u.username = ? AND u.password = ?;"
	stmtinsert, err := con.Prepare(queryinsert)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmtinsert.Close()

	var lokasi string
	var availability_surveyor string
	err = stmtinsert.QueryRow(usr.Username, usr.Password).Scan(&loginUsr.Id, &loginUsr.Username, &loginUsr.Nama_lengkap, &loginUsr.Alamat, &loginUsr.Jenis_kelamin, &loginUsr.Tgl_lahir, &loginUsr.Email, &loginUsr.No_telp, &loginUsr.Foto_profil, &loginUsr.Ktp, &loginUsr.Kelas, &loginUsr.Status, &loginUsr.Tipe, &loginUsr.First_login, &loginUsr.Denied_by_admin, &lokasi, &availability_surveyor)
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

	getPrivilegeQuery := "SELECT pr.privilege_id, p.nama_privilege FROM user_privilege pr JOIN privilege p ON pr.privilege_id = p.privilege_id WHERE pr.user_id = ?;"
	privilegestmt, err := con.Prepare(getPrivilegeQuery)
	if err != nil {
		res.Status = 401
		res.Message = "stmt update gagal"
		res.Data = err.Error()
		return res, err
	}
	defer privilegestmt.Close()

	var privilegeId int
	var privilegeName string
	err = privilegestmt.QueryRow(userId).Scan(&privilegeId, &privilegeName)
	if err != nil {
		res.Status = 401
		res.Message = "gagal mendapatkan privilege"
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
		"privilege_id":    privilegeId,
		"nama_privilege":  privilegeName,
	}

	defer db.DbClose(con)

	return res, nil
}

func SignUpSurveyor(akun string) (Response, error) {
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

	result, err := insertstmt.Exec(usr.Username, usr.Password, usr.Nama_lengkap, usr.Email, usr.No_telp)
	if err != nil {
		res.Status = 401
		res.Message = "insert user gagal"
		res.Data = err.Error()
		return res, errors.New("insert user gagal")
	}
	defer stmt.Close()

	userId, err = result.LastInsertId()
	if err != nil {
		res.Status = 500
		res.Message = "gagal mendapatkan user ID"
		res.Data = err.Error()
		return res, err
	}
	usr.Id = int(userId)

	// tambah ke user detail
	insertdetailquery := "INSERT INTO user_detail (user_detail_id,user_kelas_id,status,tipe) VALUES (?,?,?,?)"
	insertdetailstmt, err := con.Prepare(insertdetailquery)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer insertstmt.Close()

	_, err = insertdetailstmt.Exec(usr.Id, 1, 1, 7)
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

	_, err = insertrolestmt.Exec(usr.Id, 7)
	if err != nil {
		res.Status = 401
		res.Message = "insert user role gagal"
		res.Data = err.Error()
		return res, errors.New("insert user role gagal")
	}
	defer stmt.Close()

	insertprivquery := "INSERT INTO user_privilege (user_id,privilege_id) VALUES (?,?)"
	insertprivstmt, err := con.Prepare(insertprivquery)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer insertprivstmt.Close()

	_, err = insertprivstmt.Exec(usr.Id, 17)
	if err != nil {
		res.Status = 401
		res.Message = "insert user privilege gagal"
		res.Data = err.Error()
		return res, errors.New("insert user privilege gagal")
	}
	defer stmt.Close()

	// tambah ke surveyor
	insertsurvquery := "INSERT INTO surveyor (user_id,lokasi) VALUES (?,?)"
	insertsurvstmt, err := con.Prepare(insertsurvquery)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer insertsurvstmt.Close()

	_, err = insertsurvstmt.Exec(usr.Id, "")
	if err != nil {
		res.Status = 401
		res.Message = "insert surveyor gagal"
		res.Data = err.Error()
		return res, errors.New("insert surveyor gagal")
	}

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

	// hilangkan password buat global variabel
	usr.Password = ""
	res.Status = http.StatusOK
	res.Message = "Berhasil buat user"
	res.Data = usr

	defer db.DbClose(con)

	return res, nil
}

func GetAllSurveyor() (Response, error) {
	var res Response
	var dtUserSurveyor = []UserSurveyor{}

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	query := "SELECT u.user_id,u.username,u.nama_lengkap,u.alamat,u.jenis_kelamin,u.tanggal_lahir,u.email,u.nomor_telepon,u.foto_profil,u.ktp,s.suveyor_id,s.lokasi,s.availability_surveyor FROM user u JOIN surveyor s ON u.user_id = s.user_id"
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
		var _dtUserSurveyor UserSurveyor
		err := rows.Scan(&_dtUserSurveyor.User_id, &_dtUserSurveyor.Username, &_dtUserSurveyor.Nama_lengkap, &_dtUserSurveyor.Alamat, &_dtUserSurveyor.Jenis_kelamin, &_dtUserSurveyor.Tgl_lahir, &_dtUserSurveyor.Email, &_dtUserSurveyor.No_telp, &_dtUserSurveyor.Foto_profil, &_dtUserSurveyor.Ktp, &_dtUserSurveyor.Surveyor_id, &_dtUserSurveyor.Lokasi, &_dtUserSurveyor.Availability_surveyor)
		if err != nil {
			res.Status = 401
			res.Message = "scan gagal"
			res.Data = err.Error()
			return res, err
		}

		dtUserSurveyor = append(dtUserSurveyor, _dtUserSurveyor)
	}

	if err = rows.Err(); err != nil {
		res.Status = 401
		res.Message = "rows error"
		res.Data = err.Error()
		return res, err
	}

	if len(dtUserSurveyor) == 0 {
		res.Status = 404
		res.Message = "Data tidak ditemukan"
		res.Data = nil
		return res, nil
	}

	res.Status = http.StatusOK
	res.Message = "Berhasil mengambil data"
	res.Data = dtUserSurveyor

	defer db.DbClose(con)
	return res, nil
}

func GetSurveyorById(surveyor_id string) (Response, error) {
	var res Response
	var dtSurveyor Surveyor

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	query := "SELECT * FROM surveyor WHERE surveyor_id = ?"
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()
	nId, _ := strconv.Atoi(surveyor_id)
	err = stmt.QueryRow(nId).Scan(&dtSurveyor.Surveyor_id, &dtSurveyor.User_id, &dtSurveyor.Lokasi, &dtSurveyor.Availability_surveyor)
	if err != nil {
		res.Status = 401
		res.Message = "exec gagal"
		res.Data = err.Error()
		return res, err
	}

	res.Status = http.StatusOK
	res.Message = "Berhasil mengambil data"
	res.Data = dtSurveyor

	defer db.DbClose(con)
	return res, nil
}

func GetSurveyorByName(nama_surveyor string) (Response, error) {
	var res Response
	var dtUserSurveyor = []UserSurveyor{}

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	query := "SELECT u.user_id,u.username,u.nama_lengkap,u.alamat,u.jenis_kelamin,u.tanggal_lahir,u.email,u.nomor_telepon,u.foto_profil,u.ktp,s.suveyor_id,s.lokasi,s.availability_surveyor,COUNT(CASE WHEN sr.status_request = 'O' THEN 1 END) AS surveyonprogress,COUNT(sr.id_transaksi_jual_sewa) AS totalsurvey FROM user u JOIN surveyor s ON u.user_id = s.user_id LEFT JOIN survey_request sr ON u.user_id = sr.user_id WHERE u.username LIKE ? GROUP BY u.user_id"
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	rows, err := stmt.Query("%" + nama_surveyor + "%")
	if err != nil {
		res.Status = 401
		res.Message = "exec gagal"
		res.Data = err.Error()
		return res, err
	}
	defer rows.Close()

	for rows.Next() {
		var _dtUserSurveyor UserSurveyor
		err := rows.Scan(&_dtUserSurveyor.User_id, &_dtUserSurveyor.Username, &_dtUserSurveyor.Nama_lengkap, &_dtUserSurveyor.Alamat, &_dtUserSurveyor.Jenis_kelamin, &_dtUserSurveyor.Tgl_lahir, &_dtUserSurveyor.Email, &_dtUserSurveyor.No_telp, &_dtUserSurveyor.Foto_profil, &_dtUserSurveyor.Ktp, &_dtUserSurveyor.Surveyor_id, &_dtUserSurveyor.Lokasi, &_dtUserSurveyor.Availability_surveyor, &_dtUserSurveyor.SurveyOnProgress, &_dtUserSurveyor.TotalSurvey)
		if err != nil {
			res.Status = 401
			res.Message = "scan gagal"
			res.Data = err.Error()
			return res, err
		}

		dtUserSurveyor = append(dtUserSurveyor, _dtUserSurveyor)
	}

	if err = rows.Err(); err != nil {
		res.Status = 401
		res.Message = "rows error"
		res.Data = err.Error()
		return res, err
	}

	if len(dtUserSurveyor) == 0 {
		res.Status = 404
		res.Message = "Data tidak ditemukan"
		res.Data = nil
		return res, nil
	}

	res.Status = http.StatusOK
	res.Message = "Berhasil mengambil data"
	res.Data = dtUserSurveyor

	defer db.DbClose(con)
	return res, nil
}

func UpdateSurveyorById(data_surveyor string) (Response, error) {
	var res Response

	var userSurveyor UserSurveyor

	err := json.Unmarshal([]byte(data_surveyor), &userSurveyor)
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

	query := "UPDATE user u JOIN surveyor s ON u.user_id = s.user_id SET u.username = ?, u.password = ?, u.email = ?, u.nomor_telepon = ?, u.updated_at = NOW() WHERE s.suveyor_id = ? "
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(userSurveyor.Username, userSurveyor.Password, userSurveyor.Email, userSurveyor.No_telp, userSurveyor.Surveyor_id)
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

func DeleteSurveyorById(inspektur string) (Response, error) {
	var res Response

	var dtSurveyor = Surveyor{}

	err := json.Unmarshal([]byte(inspektur), &dtSurveyor)
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

	query := "DELETE FROM surveyor WHERE surveyor_id = ?"
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(dtSurveyor.Surveyor_id)
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

func GetAllSurveyorDetailed() (Response, error) {
	var res Response
	var dtUserSurveyor = []UserSurveyor{}

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	query := "SELECT u.user_id,u.username,u.nama_lengkap,u.alamat,u.jenis_kelamin,u.tanggal_lahir,u.email,u.nomor_telepon,u.foto_profil,u.ktp,s.suveyor_id,s.lokasi,s.availability_surveyor,COUNT(CASE WHEN sr.status_request = 'O' THEN 1 END) AS surveyonprogress,COUNT(sr.id_transaksi_jual_sewa) AS totalsurvey FROM user u JOIN surveyor s ON u.user_id = s.user_id LEFT JOIN survey_request sr ON u.user_id = sr.user_id GROUP BY u.user_id"
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
		var _dtUserSurveyor UserSurveyor
		err := rows.Scan(&_dtUserSurveyor.User_id, &_dtUserSurveyor.Username, &_dtUserSurveyor.Nama_lengkap, &_dtUserSurveyor.Alamat, &_dtUserSurveyor.Jenis_kelamin, &_dtUserSurveyor.Tgl_lahir, &_dtUserSurveyor.Email, &_dtUserSurveyor.No_telp, &_dtUserSurveyor.Foto_profil, &_dtUserSurveyor.Ktp, &_dtUserSurveyor.Surveyor_id, &_dtUserSurveyor.Lokasi, &_dtUserSurveyor.Availability_surveyor, &_dtUserSurveyor.SurveyOnProgress, &_dtUserSurveyor.TotalSurvey)
		if err != nil {
			res.Status = 401
			res.Message = "scan gagal"
			res.Data = err.Error()
			return res, err
		}

		dtUserSurveyor = append(dtUserSurveyor, _dtUserSurveyor)
	}

	if err = rows.Err(); err != nil {
		res.Status = 401
		res.Message = "rows error"
		res.Data = err.Error()
		return res, err
	}

	if len(dtUserSurveyor) == 0 {
		res.Status = 404
		res.Message = "Data tidak ditemukan"
		res.Data = nil
		return res, nil
	}

	res.Status = http.StatusOK
	res.Message = "Berhasil mengambil data"
	res.Data = dtUserSurveyor

	defer db.DbClose(con)
	return res, nil
}

// CRUD survey_request ============================================================================
func CreateSurveyReq(surveyreq string) (Response, error) {
	var res Response
	var dtSurveyReq = SurveyRequest{}

	err := json.Unmarshal([]byte(surveyreq), &dtSurveyReq)
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

	query := "INSERT INTO survey_request (user_id, id_asset, dateline) VALUES (?,?,?)"
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(dtSurveyReq.User_id, dtSurveyReq.Id_asset, dtSurveyReq.Dateline)
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
	dtSurveyReq.Id_transaksi_jual_sewa = int(lastId)
	dtSurveyReq.Status_request = "O"
	res.Status = http.StatusOK
	res.Message = "Berhasil memasukkan data"
	res.Data = dtSurveyReq

	defer db.DbClose(con)
	return res, nil
}

func GetAllSurveyReq() (Response, error) {
	var res Response
	var arrSurveyReq = []SurveyRequest{}
	var dtSurveyReq SurveyRequest

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	query := "SELECT * FROM survey_request"
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
		err = result.Scan(&dtSurveyReq.Id_transaksi_jual_sewa, &dtSurveyReq.User_id, &dtSurveyReq.Id_asset, &dtSurveyReq.Dateline, &dtSurveyReq.Status_request)
		if err != nil {
			res.Status = 401
			res.Message = "rows scan"
			res.Data = err.Error()
			return res, err
		}
		arrSurveyReq = append(arrSurveyReq, dtSurveyReq)
	}

	res.Status = http.StatusOK
	res.Message = "Berhasil mengambil data"
	res.Data = arrSurveyReq

	defer db.DbClose(con)
	return res, nil
}

func GetSurveyReqById(surveyreq_id string) (Response, error) {
	var res Response
	var dtSurveyReq SurveyRequest

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	query := "SELECT * FROM survey_request WHERE id_transaksi_jual_sewa = ?"
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()
	nId, _ := strconv.Atoi(surveyreq_id)
	err = stmt.QueryRow(nId).Scan(&dtSurveyReq.Id_transaksi_jual_sewa, &dtSurveyReq.User_id, &dtSurveyReq.Id_asset, &dtSurveyReq.Dateline, &dtSurveyReq.Status_request)
	if err != nil {
		res.Status = 401
		res.Message = "exec gagal"
		res.Data = err.Error()
		return res, err
	}

	res.Status = http.StatusOK
	res.Message = "Berhasil mengambil data"
	res.Data = dtSurveyReq

	defer db.DbClose(con)
	return res, nil
}

func GetSurveyReqByAsetId(aset_id string) (Response, error) {
	var res Response
	var dtSurveyReq SurveyRequest

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	query := "SELECT * FROM survey_request WHERE id_asset = ?"
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()
	nId, _ := strconv.Atoi(aset_id)
	err = stmt.QueryRow(nId).Scan(&dtSurveyReq.Id_transaksi_jual_sewa, &dtSurveyReq.User_id, &dtSurveyReq.Id_asset, &dtSurveyReq.Dateline, &dtSurveyReq.Status_request, &dtSurveyReq.Data_lengkap, &dtSurveyReq.Usage_old, &dtSurveyReq.Usage_new, &dtSurveyReq.Luas_old, &dtSurveyReq.Luas_new, &dtSurveyReq.Nilai_old, &dtSurveyReq.Nilai_new, &dtSurveyReq.Kondisi_old, &dtSurveyReq.Kondisi_new, &dtSurveyReq.Batas_koordinat_old, &dtSurveyReq.Batas_koordinat_new, &dtSurveyReq.Tags_old, &dtSurveyReq.Tags_new)
	if err != nil {
		res.Status = 401
		res.Message = "exec gagal"
		res.Data = err.Error()
		return res, err
	}

	res.Status = http.StatusOK
	res.Message = "Berhasil mengambil data"
	res.Data = dtSurveyReq

	defer db.DbClose(con)
	return res, nil
}

func UpdateSurveyReqById(surveyreq string) (Response, error) {
	var res Response

	var dtSurveyReq = SurveyRequest{}

	err := json.Unmarshal([]byte(surveyreq), &dtSurveyReq)
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

	query := "UPDATE surveyor SET dateline = ?, status_request = ? WHERE id_transaksi_jual_sewa = ?"
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(dtSurveyReq.Dateline, dtSurveyReq.Status_request, dtSurveyReq.Id_transaksi_jual_sewa)
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

func DeleteSurveyReqById(surveyreq string) (Response, error) {
	var res Response

	var dtSurveyReq = SurveyRequest{}

	err := json.Unmarshal([]byte(surveyreq), &dtSurveyReq)
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

	query := "DELETE FROM survey_request WHERE id_transaksi_jual_sewa = ?"
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(dtSurveyReq.Id_transaksi_jual_sewa)
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

// CRUD transaction_request ============================================================================
func CreateTranReq(tranreq string) (Response, error) {
	var res Response
	var dtTranReq = TransactionRequest{}

	err := json.Unmarshal([]byte(tranreq), &dtTranReq)
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

	query := "INSERT INTO transaction_request (user_id, id_asset, tipe, masa_sewa, meeting_log) VALUES (?,?,?,?,?)"
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(dtTranReq.User_id, dtTranReq.Id_asset, dtTranReq.Tipe, dtTranReq.Masa_sewa, dtTranReq.Meeting_log)
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
	dtTranReq.Id_transaksi_jual_sewa = int(lastId)

	res.Status = http.StatusOK
	res.Message = "Berhasil memasukkan data"
	res.Data = dtTranReq

	defer db.DbClose(con)
	return res, nil
}

func GetAllTranReq() (Response, error) {
	var res Response
	var arrTranReq = []TransactionRequest{}
	var dtTranReq TransactionRequest

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	query := "SELECT * FROM transaction_request"
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
		err = result.Scan(&dtTranReq.Id_transaksi_jual_sewa, &dtTranReq.User_id, &dtTranReq.Id_asset, &dtTranReq.Tipe, &dtTranReq.Masa_sewa, &dtTranReq.Meeting_log)
		if err != nil {
			res.Status = 401
			res.Message = "rows scan"
			res.Data = err.Error()
			return res, err
		}
		arrTranReq = append(arrTranReq, dtTranReq)
	}

	res.Status = http.StatusOK
	res.Message = "Berhasil mengambil data"
	res.Data = arrTranReq

	defer db.DbClose(con)
	return res, nil
}

func GetTranReqById(tranreq_id string) (Response, error) {
	var res Response
	var dtTranReq TransactionRequest

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	query := "SELECT * FROM transaction_request WHERE id_transaksi_jual_sewa = ?"
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()
	nId, _ := strconv.Atoi(tranreq_id)
	err = stmt.QueryRow(nId).Scan(&dtTranReq.Id_transaksi_jual_sewa, &dtTranReq.User_id, &dtTranReq.Id_asset, &dtTranReq.Tipe, &dtTranReq.Masa_sewa, &dtTranReq.Meeting_log)
	if err != nil {
		res.Status = 401
		res.Message = "exec gagal"
		res.Data = err.Error()
		return res, err
	}

	res.Status = http.StatusOK
	res.Message = "Berhasil mengambil data"
	res.Data = dtTranReq

	defer db.DbClose(con)
	return res, nil
}

func GetTranReqByTipe(nama_tipe string) (Response, error) {
	var res Response
	var dtTranReqs = []TransactionRequest{}

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	query := "SELECT * FROM transaction_request WHERE tipe LIKE ?"
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	rows, err := stmt.Query("%" + nama_tipe + "%")
	if err != nil {
		res.Status = 401
		res.Message = "exec gagal"
		res.Data = err.Error()
		return res, err
	}
	defer rows.Close()
	for rows.Next() {
		var dtTranReq TransactionRequest
		err := rows.Scan(&dtTranReq.Id_transaksi_jual_sewa, &dtTranReq.User_id, &dtTranReq.Id_asset, &dtTranReq.Tipe, &dtTranReq.Masa_sewa, &dtTranReq.Meeting_log)
		if err != nil {
			res.Status = 401
			res.Message = "scan gagal"
			res.Data = err.Error()
			return res, err
		}
		dtTranReqs = append(dtTranReqs, dtTranReq)
	}

	if err = rows.Err(); err != nil {
		res.Status = 401
		res.Message = "rows error"
		res.Data = err.Error()
		return res, err
	}

	if len(dtTranReqs) == 0 {
		res.Status = 404
		res.Message = "Data tidak ditemukan"
		res.Data = nil
		return res, nil
	}

	res.Status = http.StatusOK
	res.Message = "Berhasil mengambil data"
	res.Data = dtTranReqs

	defer db.DbClose(con)
	return res, nil
}

func UpdateTranReqById(tranreq string) (Response, error) {
	var res Response

	var dtTranReq = TransactionRequest{}

	err := json.Unmarshal([]byte(tranreq), &dtTranReq)
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

	query := "UPDATE transaction_request SET tipe = ?, masa_sewa = ?, meeting_log = ? WHERE id_transaksi_jual_sewa = ?"
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(dtTranReq.Tipe, dtTranReq.Masa_sewa, dtTranReq.Meeting_log, dtTranReq.Id_transaksi_jual_sewa)
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

func DeleteTranReqById(tranreq string) (Response, error) {
	var res Response

	var dtTranReq = SurveyRequest{}

	err := json.Unmarshal([]byte(tranreq), &dtTranReq)
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

	query := "DELETE FROM transaction_request WHERE id_transaksi_jual_sewa = ?"
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(dtTranReq.Id_transaksi_jual_sewa)
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

func GetTranReqByUserId(user_id string) (Response, error) {
	var res Response
	// var dtTranReq TransactionRequest
	var dtTranReq = []TransactionRequest{}

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	query := "SELECT * FROM transaction_request WHERE user_id = ?"
	rows, err := con.Query(query, user_id)
	if err != nil {
		res.Status = 401
		res.Message = "Failed to execute query"
		res.Data = err.Error()
		return res, err
	}
	defer rows.Close()

	for rows.Next() {
		var tranReq TransactionRequest
		err = rows.Scan(&tranReq.Id_transaksi_jual_sewa, &tranReq.User_id, &tranReq.Id_asset, &tranReq.Tipe, &tranReq.Masa_sewa, &tranReq.Meeting_log)
		if err != nil {
			res.Status = 401
			res.Message = "Failed to scan row"
			res.Data = err.Error()
			return res, err
		}
		dtTranReq = append(dtTranReq, tranReq)
	}
	if err = rows.Err(); err != nil {
		res.Status = 401
		res.Message = "Failed during row iteration"
		res.Data = err.Error()
		return res, err
	}

	res.Status = http.StatusOK
	res.Message = "Berhasil mengambil data"
	res.Data = dtTranReq

	defer db.DbClose(con)
	return res, nil
}

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

	// cek apakah password benar atau tidak
	// queryinsert := "SELECT user_id, username, nama_lengkap, alamat, jenis_kelamin, tanggal_lahir, email, nomor_telepon, foto_profil, ktp FROM user WHERE username = ? AND password = ?"
	queryinsert := "SELECT u.user_id, u.username, u.nama_lengkap, u.alamat, u.jenis_kelamin, u.tanggal_lahir, u.email, u.nomor_telepon, u.foto_profil, u.ktp, ud.user_kelas_id, ud.status, ud.tipe, ud.first_login, ud.denied_by_admin FROM user u JOIN user_detail ud ON u.user_id = ud.user_detail_id WHERE u.username = ? AND u.password = ?;"
	stmtinsert, err := con.Prepare(queryinsert)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmtinsert.Close()

	err = stmtinsert.QueryRow(usr.Username, usr.Password).Scan(&loginUsr.Id, &loginUsr.Username, &loginUsr.Nama_lengkap, &loginUsr.Alamat, &loginUsr.Jenis_kelamin, &loginUsr.Tgl_lahir, &loginUsr.Email, &loginUsr.No_telp, &loginUsr.Foto_profil, &loginUsr.Ktp, &loginUsr.Kelas, &loginUsr.Status, &loginUsr.Tipe, &loginUsr.First_login, &loginUsr.Denied_by_admin)
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

	getPrivilegeQuery := "SELECT pr.privilege_id, p.nama_privilege FROM user_privilege pr JOIN privilege p ON pr.privilege_id = p.privilege_id WHERE pr.user_id = ?;"
	privilegestmt, err := con.Prepare(getPrivilegeQuery)
	if err != nil {
		res.Status = 401
		res.Message = "stmt update gagal"
		res.Data = err.Error()
		return res, err
	}
	defer privilegestmt.Close()

	var privilegeId int
	var privilegeName string
	err = privilegestmt.QueryRow(userId).Scan(&privilegeId, &privilegeName)
	if err != nil {
		res.Status = 401
		res.Message = "gagal mendapatkan privilege"
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
		"privilege_id":    privilegeId,
		"nama_privilege":  privilegeName,
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

	result, err := insertstmt.Exec(usr.Username, usr.Password, usr.Nama_lengkap, usr.Email, usr.No_telp)
	if err != nil {
		res.Status = 401
		res.Message = "insert user gagal"
		res.Data = err.Error()
		return res, errors.New("insert user gagal")
	}
	defer stmt.Close()

	userId, err = result.LastInsertId()
	if err != nil {
		res.Status = 500
		res.Message = "gagal mendapatkan user ID"
		res.Data = err.Error()
		return res, err
	}
	usr.Id = int(userId)

	// tambah ke user detail
	insertdetailquery := "INSERT INTO user_detail (user_detail_id,user_kelas_id,status,tipe) VALUES (?,?,?,?)"
	insertdetailstmt, err := con.Prepare(insertdetailquery)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer insertstmt.Close()

	_, err = insertdetailstmt.Exec(usr.Id, 1, 1, 8)
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

	insertprivquery := "INSERT INTO user_privilege (user_id,privilege_id) VALUES (?,?)"
	insertprivstmt, err := con.Prepare(insertprivquery)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer insertprivstmt.Close()

	_, err = insertprivstmt.Exec(usr.Id, 17)
	if err != nil {
		res.Status = 401
		res.Message = "insert user privilege gagal"
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
	var dtUser User

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	query := "SELECT user_id, username, password,nama_lengkap,alamat,jenis_kelamin,tanggal_lahir,email,nomor_telepon,foto_profil,ktp FROM user"
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

	query := "SELECT user_id, username, nama_lengkap, alamat, jenis_kelamin, tanggal_lahir, email, nomor_telepon, foto_profil, ktp FROM user WHERE user_id = ?"
	stmt, err := con.Prepare(query)

	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	nId, _ := strconv.Atoi(id_user)
	err = stmt.QueryRow(nId).Scan(&usr.Id, &usr.Username, &usr.Nama_lengkap, &usr.Alamat, &usr.Jenis_kelamin, &usr.Tgl_lahir, &usr.Email, &usr.No_telp, &usr.Foto_profil, &usr.Ktp)
	if err != nil {
		res.Status = 401
		res.Message = "rows gagal"
		res.Data = err.Error()
		return res, err
	}

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
	fmt.Println(filefoto.Header.Get("Content-type"))
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
	fmt.Println(filefoto.Header.Get("Content-type"))
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
	fmt.Println("")
	for result.Next() {
		var dtUser User
		err = result.Scan(&dtUser.Id, &dtUser.Username, &dtUser.Password, &dtUser.Nama_lengkap, &dtUser.Alamat, &dtUser.Jenis_kelamin, &dtUser.Tgl_lahir, &dtUser.Email, &dtUser.No_telp, &dtUser.Foto_profil, &dtUser.Ktp)
		if err != nil {
			res.Status = 401
			res.Message = "rows scan"
			res.Data = err.Error()
			return res, err
		}
		fmt.Println("data user:", dtUser)
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

// CRUD user_privilege ============================================================================
func CreateUserPriv(userPriv string) (Response, error) {
	var res Response
	var dtUserPriv = UserPrivilege{}

	err := json.Unmarshal([]byte(userPriv), &dtUserPriv)
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

	query := "INSERT INTO user_privilege (privilege_id, user_id) VALUES (?,?)"
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(dtUserPriv.Privilege_id, dtUserPriv.User_id)
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
	dtUserPriv.User_privilege_id = int(lastId)

	res.Status = http.StatusOK
	res.Message = "Berhasil memasukkan data"
	res.Data = dtUserPriv

	defer db.DbClose(con)
	return res, nil
}

func GetAllUserPriv() (Response, error) {
	var res Response
	var arrUserPriv = []UserPrivilege{}
	var dtUserPriv UserPrivilege

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	query := "SELECT * FROM user_privilege"
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
		err = result.Scan(&dtUserPriv.User_privilege_id, &dtUserPriv.Privilege_id, &dtUserPriv.User_id)
		if err != nil {
			res.Status = 401
			res.Message = "rows scan"
			res.Data = err.Error()
			return res, err
		}
		arrUserPriv = append(arrUserPriv, dtUserPriv)
	}

	res.Status = http.StatusOK
	res.Message = "Berhasil mengambil data"
	res.Data = arrUserPriv

	defer db.DbClose(con)
	return res, nil
}

func GetUserPrivById(user_priv_id string) (Response, error) {
	var res Response
	var dtUserPriv UserPrivilege

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	query := "SELECT * FROM user_privilege WHERE user_privilege_id = ?"
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()
	nId, _ := strconv.Atoi(user_priv_id)
	err = stmt.QueryRow(nId).Scan(&dtUserPriv.User_privilege_id, &dtUserPriv.User_id)
	if err != nil {
		res.Status = 401
		res.Message = "exec gagal"
		res.Data = err.Error()
		return res, err
	}

	res.Status = http.StatusOK
	res.Message = "Berhasil mengambil data"
	res.Data = dtUserPriv

	defer db.DbClose(con)
	return res, nil
}

func GetUserPrivByUserId(user_id string) (Response, error) {
	var res Response
	var dtUserPriv UserPrivilege

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	query := "SELECT * FROM user_privilege WHERE user_id = ?"
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()
	nId, _ := strconv.Atoi(user_id)
	err = stmt.QueryRow(nId).Scan(&dtUserPriv.User_privilege_id, &dtUserPriv.User_id)
	if err != nil {
		res.Status = 401
		res.Message = "exec gagal"
		res.Data = err.Error()
		return res, err
	}

	res.Status = http.StatusOK
	res.Message = "Berhasil mengambil data"
	res.Data = dtUserPriv

	defer db.DbClose(con)
	return res, nil
}

func GetUserPrivDetailByUserId(user_id string) (Response, error) {
	var res Response
	var privileges []map[string]interface{}

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	query := "SELECT up.user_id, up.privilege_id, p.nama_privilege FROM user_privilege up JOIN privilege p ON up.privilege_id = p.privilege_id WHERE up.user_id = ?"

	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	nId, _ := strconv.Atoi(user_id)
	rows, err := stmt.Query(nId)
	// err = stmt.QueryRow(nId).Scan(&temp_user_id, &temp_privilege_id, &temp_nama_privilege)
	if err != nil {
		res.Status = 401
		res.Message = "exec gagal"
		res.Data = err.Error()
		return res, err
	}
	defer rows.Close()

	for rows.Next() {
		var temp_privilege_id, temp_user_id int
		var temp_nama_privilege string

		err := rows.Scan(&temp_user_id, &temp_privilege_id, &temp_nama_privilege)
		if err != nil {
			res.Status = 401
			res.Message = "scan gagal"
			res.Data = err.Error()
			return res, err
		}

		privilege := map[string]interface{}{
			"user_id":        temp_user_id,
			"privilege_id":   temp_privilege_id,
			"nama_privilege": temp_nama_privilege,
		}
		privileges = append(privileges, privilege)
	}

	if len(privileges) == 0 {
		res.Status = 404
		res.Message = "Data tidak ditemukan"
		res.Data = nil
		return res, nil
	}

	res.Status = http.StatusOK
	res.Message = "Berhasil mengambil data"
	res.Data = privileges

	defer db.DbClose(con)
	return res, nil
}

func UpdateUserPriv(userPriv string) (Response, error) {
	var res Response

	var dtUserPriv = UserPrivilege{}

	err := json.Unmarshal([]byte(userPriv), &dtUserPriv)
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

	query := "UPDATE user_privilege SET privilege_id = ?, user_id = ? WHERE user_privilege_id = ?"
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(dtUserPriv.Privilege_id, dtUserPriv.User_id, dtUserPriv.User_privilege_id)
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

func DeleteUserPriv(userPriv string) (Response, error) {
	var res Response

	var dtUserPriv = UserPrivilege{}

	err := json.Unmarshal([]byte(userPriv), &dtUserPriv)
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

	query := "DELETE FROM user_privilege WHERE user_privilege_id = ?"
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(dtUserPriv.User_privilege_id)
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

// CRUD user_role

func CreateUserRole(userRole string) (Response, error) {
	var res Response
	var dtUserRole = UserRole{}

	err := json.Unmarshal([]byte(userRole), &dtUserRole)
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

	query := "INSERT INTO user_role (user_id, role_id) VALUES (?,?)"
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(dtUserRole.User_id, dtUserRole.Role_id)
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
	dtUserRole.User_role_id = int(lastId)

	res.Status = http.StatusOK
	res.Message = "Berhasil memasukkan data"
	res.Data = dtUserRole

	defer db.DbClose(con)
	return res, nil
}

func GetAllUserRole() (Response, error) {
	var res Response
	var arrUserRole = []UserRole{}
	var dtUserRole UserRole

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	query := "SELECT * FROM user_role"
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
		err = result.Scan(&dtUserRole.User_role_id, &dtUserRole.User_id, &dtUserRole.Role_id)
		if err != nil {
			res.Status = 401
			res.Message = "rows scan"
			res.Data = err.Error()
			return res, err
		}
		arrUserRole = append(arrUserRole, dtUserRole)
	}

	res.Status = http.StatusOK
	res.Message = "Berhasil mengambil data"
	res.Data = arrUserRole

	defer db.DbClose(con)
	return res, nil
}

func GetUserRoleById(user_role_id string) (Response, error) {
	var res Response
	var dtUserRole UserRole

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	query := "SELECT * FROM user_role WHERE user_role_id = ?"
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()
	nId, _ := strconv.Atoi(user_role_id)
	err = stmt.QueryRow(nId).Scan(&dtUserRole.User_role_id, &dtUserRole.User_id, &dtUserRole.Role_id)
	if err != nil {
		res.Status = 401
		res.Message = "exec gagal"
		res.Data = err.Error()
		return res, err
	}

	res.Status = http.StatusOK
	res.Message = "Berhasil mengambil data"
	res.Data = dtUserRole

	defer db.DbClose(con)
	return res, nil
}

func GetUserRoleByUserId(user_id string) (Response, error) {
	var res Response
	var dtUserRole UserRole

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	query := "SELECT * FROM user_role WHERE user_id = ?"
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()
	nId, _ := strconv.Atoi(user_id)
	err = stmt.QueryRow(nId).Scan(&dtUserRole.User_role_id, &dtUserRole.User_id, &dtUserRole.Role_id)
	if err != nil {
		res.Status = 401
		res.Message = "exec gagal"
		res.Data = err.Error()
		return res, err
	}

	res.Status = http.StatusOK
	res.Message = "Berhasil mengambil data"
	res.Data = dtUserRole

	defer db.DbClose(con)
	return res, nil
}

func GetUserRoleDetailByUserId(user_id string) (Response, error) {
	var res Response
	var roles []map[string]interface{}

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	query := "SELECT ur.user_id, ur.role_id, r.nama_role FROM user_role ur JOIN role r ON ur.role_id = r.role_id WHERE ur.user_id = ?"
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	nId, _ := strconv.Atoi(user_id)
	rows, err := stmt.Query(nId)
	// err = stmt.QueryRow(nId).Scan(&temp_user_id, &temp_role_id, &temp_nama_role)
	if err != nil {
		res.Status = 401
		res.Message = "exec gagal"
		res.Data = err.Error()
		return res, err
	}
	defer rows.Close()

	for rows.Next() {
		var temp_user_id, temp_role_id int
		var temp_nama_role string
		err := rows.Scan(&temp_user_id, &temp_role_id, &temp_nama_role)
		if err != nil {
			res.Status = 401
			res.Message = "scan gagal"
			res.Data = err.Error()
			return res, err
		}

		role := map[string]interface{}{
			"user_id":   temp_user_id,
			"role_id":   temp_role_id,
			"nama_role": temp_nama_role,
		}
		roles = append(roles, role)
	}

	if len(roles) == 0 {
		res.Status = 404
		res.Message = "Data tidak ditemukan"
		res.Data = nil
		return res, nil
	}

	res.Status = http.StatusOK
	res.Message = "Berhasil mengambil data"
	res.Data = roles

	defer db.DbClose(con)
	return res, nil
}

func UpdateUserRole(userRole string) (Response, error) {
	var res Response

	var dtUserRole = UserRole{}

	err := json.Unmarshal([]byte(userRole), &dtUserRole)
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

	query := "UPDATE user_role SET user_id = ?, role_id = ? WHERE user_role_id = ?"
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(dtUserRole.User_id, dtUserRole.Role_id, dtUserRole.User_role_id)
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

func DeleteUserRole(userRole string) (Response, error) {
	var res Response

	var dtUserRole = UserRole{}

	err := json.Unmarshal([]byte(userRole), &dtUserRole)
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

	query := "DELETE FROM user_role WHERE user_role_id = ?"
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(dtUserRole.User_role_id)
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

// fungsi tambahan
func UploadFile(file *multipart.FileHeader, id string, kolom_id string, folder string) (Response, error) {
	var res Response

	log.Println("Upload File")
	nId, _ := strconv.Atoi(id)
	// file.Filename =
	pathFile := "uploads/user/" + file.Filename
	//source
	src, err := file.Open()
	if err != nil {
		log.Println(err.Error())
		log.Println("1")
		fmt.Print("1")
		return res, err
	}
	defer src.Close()

	// Destination
	dst, err := os.Create("uploads/" + folder + "/" + file.Filename)
	if err != nil {
		log.Println("2")
		fmt.Print("2")
		return res, err
	}

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		log.Println(err.Error())
		log.Println("3")
		fmt.Print("3")
		return res, err
	}
	dst.Close()

	err = UpdateDataFotoPath(folder, "foto", pathFile, kolom_id, nId)
	if err != nil {
		return res, err
	}

	res.Status = http.StatusOK
	res.Message = "Sukses Upload File"
	res.Data = file.Filename

	return res, nil
}

func UpdateDataFotoPath(tabel string, kolom string, path string, kolom_id string, id int) error {
	log.Println("mengubah status foto di DB")
	fmt.Println("mengubah status foto di DB")
	// Open DB connection
	con, err := db.DbConnection()
	if err != nil {
		log.Println("error: " + err.Error())
		return err
	}
	defer db.DbClose(con) // Ensure the connection is closed

	// Build the SQL query
	query := fmt.Sprintf("UPDATE %s SET %s='%s' WHERE %s = %d", tabel, kolom, path, kolom_id, id)
	fmt.Println(query)
	// Execute the query
	_, err = con.Exec(query) // Use Exec instead of Query since this is an UPDATE operation
	if err != nil {
		log.Println("error executing query: " + err.Error())
		return err
	}

	fmt.Println("status foto di edit")
	return nil
}

func VerifyUserAccept(input string) (Response, error) {
	var res Response

	type temp_verif_user_acc struct {
		UserID int `json:"userid"`
		Kelas  int `json:"kelas"`
	}
	var requestacc temp_verif_user_acc
	err := json.Unmarshal([]byte(input), &requestacc)
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

	query := "UPDATE user_detail SET user_kelas_id=?,status='V' WHERE user_detail_id = ?"
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(requestacc.Kelas, requestacc.UserID)
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

func VerifyPerusahaanAccept(input string) (Response, error) {
	var res Response

	type temp_verif_perusahaan_acc struct {
		PerusahaanId int `json:"perusahaan_id"`
		Kelas        int `json:"kelas"`
	}
	var requestacc temp_verif_perusahaan_acc
	err := json.Unmarshal([]byte(input), &requestacc)
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

	query := "UPDATE perusahaan SET kelas=?,status='V' WHERE perusahaan_id = ?"
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(requestacc.Kelas, requestacc.PerusahaanId)
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

func VerifyUserDecline(input string) (Response, error) {
	var res Response

	type temp_verif_user_deny struct {
		UserID int `json:"userid"`
	}
	var requestacc temp_verif_user_deny
	err := json.Unmarshal([]byte(input), &requestacc)
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

	query := "UPDATE user_detail SET status='N',denied_by_admin='Y' WHERE user_detail_id = ?"
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(requestacc.UserID)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}

	// kirim notif (masih mendatang)

	res.Status = http.StatusOK
	res.Message = "Berhasil mengupdate data"
	res.Data = result

	defer db.DbClose(con)
	return res, nil
}

func VerifyPerusahaanDecline(input string) (Response, error) {
	var res Response

	type temp_verif_perusahaan_deny struct {
		PerusahaanId int `json:"perusahaan_id"`
	}
	var requestacc temp_verif_perusahaan_deny
	err := json.Unmarshal([]byte(input), &requestacc)
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

	query := "UPDATE perusahaan SET status='N',denied_by_admin='Y' WHERE perusahaan_id = ?"
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(requestacc.PerusahaanId)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}

	// kirim notif (masih mendatang)

	res.Status = http.StatusOK
	res.Message = "Berhasil mengupdate data"
	res.Data = result

	defer db.DbClose(con)
	return res, nil
}

func VerifyAssetAccept(input string) (Response, error) {
	var res Response
	var dtSurveyReq SurveyRequest

	type temp_verif_asset_acc struct {
		SurveryReqId int `json:"surveyreq_id"`
	}

	var requestacc temp_verif_asset_acc
	err := json.Unmarshal([]byte(input), &requestacc)
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

	query := "UPDATE survey_request SET status_request='F' WHERE id_transaksi_jual_sewa = ?"
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	_, err = stmt.Exec(requestacc.SurveryReqId)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}

	selectquery := "SELECT * FROM survey_request WHERE id_transaksi_jual_sewa = ?"
	selectstmt, err := con.Prepare(selectquery)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer selectstmt.Close()
	var temp_tags_old sql.NullString
	var temp_tags_new sql.NullString

	err = selectstmt.QueryRow(requestacc.SurveryReqId).Scan(&dtSurveyReq.Id_transaksi_jual_sewa, &dtSurveyReq.User_id, &dtSurveyReq.Id_asset, &dtSurveyReq.Dateline, &dtSurveyReq.Status_request, &dtSurveyReq.Data_lengkap, &dtSurveyReq.Usage_old, &dtSurveyReq.Usage_new, &dtSurveyReq.Luas_old, &dtSurveyReq.Luas_new, &dtSurveyReq.Nilai_old, &dtSurveyReq.Nilai_new, &dtSurveyReq.Kondisi_old, &dtSurveyReq.Kondisi_new, &dtSurveyReq.Batas_koordinat_old, &dtSurveyReq.Batas_koordinat_new, &temp_tags_old, &temp_tags_new)
	if err != nil {
		res.Status = 401
		res.Message = "rows gagal"
		res.Data = err.Error()
		return res, err
	}

	if temp_tags_old.Valid {
		dtSurveyReq.Tags_old = temp_tags_old.String
	} else {
		dtSurveyReq.Tags_old = ""
	}

	if temp_tags_new.Valid {
		dtSurveyReq.Tags_new = temp_tags_new.String
	} else {
		dtSurveyReq.Tags_new = ""
	}

	// update data asset dengan yang baru
	updatequery := "UPDATE asset SET `kondisi`= ?,`batas_koordinat`= ?,`luas`= ?,`nilai`= ?,`usage`= ? WHERE `id_asset`= ?"
	updatestmt, err := con.Prepare(updatequery)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer updatestmt.Close()

	_, err = updatestmt.Exec(dtSurveyReq.Kondisi_new, dtSurveyReq.Batas_koordinat_new, dtSurveyReq.Luas_new, dtSurveyReq.Nilai_new, dtSurveyReq.Usage_new, dtSurveyReq.Id_asset)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}

	update2query := "UPDATE asset_tags SET id_tags= ? WHERE id_asset= ?"
	update2stmt, err := con.Prepare(update2query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer update2stmt.Close()

	_, err = update2stmt.Exec(dtSurveyReq.Tags_new, dtSurveyReq.Id_asset)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}

	tempaset, _ := GetAssetById(string(dtSurveyReq.Id_asset))

	res.Status = http.StatusOK
	res.Message = "Berhasil mengupdate data"
	res.Data = tempaset

	defer db.DbClose(con)
	return res, nil
}

func ReassignAsset(input string) (Response, error) {
	var res Response

	type temp_verif_asset_acc struct {
		SurveyReqId int    `json:"surveyreq_id"`
		SurveyorId  int    `json:"surveyor_id"`
		Dateline    string `json:"dateline"`
	}

	var requestacc temp_verif_asset_acc
	err := json.Unmarshal([]byte(input), &requestacc)
	if err != nil {
		res.Status = 401
		res.Message = "gagal unmarshal JSON"
		res.Data = err.Error()
		return res, err
	}
	fmt.Println(requestacc)
	fmt.Println(requestacc.SurveyorId)

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	// ambil user id dari surveyor
	surveyorquery := "SELECT user_id FROM surveyor WHERE `suveyor_id` = ?"
	surveyorstmt, err := con.Prepare(surveyorquery)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer surveyorstmt.Close()

	var _tempuserid int
	err = surveyorstmt.QueryRow(requestacc.SurveyorId).Scan(&_tempuserid)
	if err != nil {
		res.Status = 401
		res.Message = "gagal mengambil user_id"
		res.Data = err.Error()
		return res, err
	}
	fmt.Println(requestacc.SurveyorId)
	fmt.Println(_tempuserid)

	query := "UPDATE survey_request SET `user_id`=?,`dateline`=?,`status_request`='R' WHERE id_transaksi_jual_sewa = ?"
	stmt2, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt2.Close()

	result, err := stmt2.Exec(_tempuserid, requestacc.Dateline, requestacc.SurveyReqId)
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

// BELUM SELESAI
func ForgotPass(email string) (Response, error) {
	var res Response

	// asd

	return res, nil
}

func ChangePass() (Response, error) {
	var res Response

	return res, nil
}
