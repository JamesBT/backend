package model

import (
	"TemplateProject/db"
	"database/sql"
	"encoding/json"
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
func CreateAsset(filelegalitas *multipart.FileHeader, suratkuasa *multipart.FileHeader, nama, perusahaan_id, tipe, nomorlegalitas, status, alamat, kondisi, koordinat, batas_koordinat, luas, nilai, provinsi string) (Response, error) {
	var res Response

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	query := "INSERT INTO asset (perusahaan_id, nama, tipe, nomor_legalitas, status_asset, alamat, kondisi, titik_koordinat, batas_koordinat, luas, nilai, provinsi, created_at) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,NOW())"
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(perusahaan_id, nama, tipe, nomorlegalitas, status, alamat, kondisi, koordinat, batas_koordinat, luas, nilai, provinsi)
	if err != nil {
		res.Status = 401
		res.Message = "exec gagal"
		res.Data = err.Error()
		return res, err
	}

	fmt.Println("SELESAI QUEUE MASUKIN")

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

	var tempaset Response
	tempaset, _ = GetAssetById(strconv.Itoa(tempid))

	res.Status = http.StatusOK
	res.Message = "Berhasil memasukkan data"
	res.Data = tempaset.Data

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
		err = result.Scan(&dtAset.Id_asset, &idAssetParent, &idAssetChild, &idJoin, &idPerusahaan, &dtAset.Nama, &dtAset.Tipe, &dtAset.Nomor_legalitas, &dtAset.File_legalitas, &dtAset.Status_asset, &dtAset.Surat_kuasa, &dtAset.Alamat, &dtAset.Kondisi, &dtAset.Titik_koordinat, &dtAset.Batas_koordinat, &dtAset.Luas, &dtAset.Nilai, &dtAset.Provinsi, &dtAset.Usage, &dtAset.Status_pengecekan, &dtAset.Status_verifikasi, &dtAset.Status_publik, &dtAset.Hak_akses, &masaSewa, &dtAset.Created_at, &deleteAt)
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
			dtAset.Id_join = strconv.Itoa(int(idJoin.Int32))
		} else {
			dtAset.Id_join = "0"
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
	err = stmt.QueryRow(nId).Scan(&dtAset.Id_asset, &idAssetParent, &idAssetChild, &idJoin, &dtAset.Id_perusahaan, &dtAset.Nama, &dtAset.Tipe, &dtAset.Nomor_legalitas, &dtAset.File_legalitas, &dtAset.Status_asset, &dtAset.Surat_kuasa, &dtAset.Alamat, &dtAset.Kondisi, &dtAset.Titik_koordinat, &dtAset.Batas_koordinat, &dtAset.Luas, &dtAset.Nilai, &dtAset.Provinsi, &dtAset.Usage, &dtAset.Status_pengecekan, &dtAset.Status_verifikasi, &dtAset.Status_publik, &dtAset.Hak_akses, &masaSewa, &dtAset.Created_at, &deleteAt)
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
		dtAset.Id_join = strconv.Itoa(int(idJoin.Int32))
	} else {
		dtAset.Id_join = "0"
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
		dtAset.Id_join = strconv.Itoa(int(idJoin.Int32))
	} else {
		dtAset.Id_join = "0"
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

func JoinAsset(id_asset_1, id_asset_2 string) (Response, error) {
	var res Response

	return res, nil
}
