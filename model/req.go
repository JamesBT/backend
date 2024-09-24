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
)

// CRUD survey_request ============================================================================
func CreateSurveyReq(surat *multipart.FileHeader, user_id, id_asset, dateline string) (Response, error) {
	var res Response
	var dtSurveyReq = SurveyRequest{}

	dtSurveyReq.User_id, _ = strconv.Atoi(user_id)
	dtSurveyReq.Id_asset, _ = strconv.Atoi(id_asset)
	dtSurveyReq.Dateline = dateline

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
	fmt.Println("ambil data tags dan usage")

	// ambil data aset dan update bagian lama
	var nilai, kondisi, batas_koordinat, titik_koordinat string
	var luas float64
	var tags []string
	var usages []string

	queryAsset := "SELECT id_penggunaan FROM asset_penggunaan WHERE id_asset = ?"
	rowtags, err := con.Query(queryAsset, dtSurveyReq.Id_asset)
	if err != nil {
		res.Status = 401
		res.Message = "query gagal"
		res.Data = err.Error()
		return res, err
	}
	defer rowtags.Close()
	for rowtags.Next() {
		var usage string
		if err := rowtags.Scan(&usage); err != nil {
			return res, err
		}
		usages = append(usages, usage)
	}

	queryTags := "SELECT id_tags FROM asset_tags WHERE id_asset = ?"
	rows, err := con.Query(queryTags, dtSurveyReq.Id_asset)
	if err != nil {
		res.Status = 401
		res.Message = "query gagal"
		res.Data = err.Error()
		return res, err
	}
	defer rows.Close()

	for rows.Next() {
		var tag string
		if err := rows.Scan(&tag); err != nil {
			return res, err
		}
		tags = append(tags, tag)
	}
	fmt.Println("update")

	tagsString := strings.Join(tags, ", ")
	usagesString := strings.Join(usages, ", ")

	queryGETasset := "SELECT luas, nilai, kondisi, batas_koordinat, titik_koordinat FROM asset WHERE id_asset = ?"
	err = con.QueryRow(queryGETasset, id_asset).Scan(&luas, &nilai, &kondisi, &batas_koordinat, &titik_koordinat)
	if err != nil {
		res.Status = 401
		res.Message = "query gagal"
		res.Data = err.Error()
		return res, err
	}

	updateQuery := `
		UPDATE survey_request 
		SET usage_old = ?, luas_old = ?, nilai_old = ?, kondisi_old = ?, batas_koordinat_old = ?, titik_koordinat_old = ?, tags_old = ? 
		WHERE id_transaksi_jual_sewa = ?
	`

	_, err = con.Exec(updateQuery, usagesString, luas, nilai, kondisi, batas_koordinat, titik_koordinat, tagsString, lastId)
	if err != nil {
		return res, err
	}
	fmt.Println("masukin file")
	dtSurveyReq.Id_transaksi_jual_sewa = int(lastId)
	dtSurveyReq.Status_request = "O"

	nId := strconv.Itoa(int(lastId))
	// masukin file
	surat.Filename = nId + "_" + surat.Filename
	pathFotoFile := "uploads/survey_req/surat/" + surat.Filename
	//source
	srcfoto, err := surat.Open()
	if err != nil {
		log.Println(err.Error())
		return res, err
	}
	defer srcfoto.Close()

	// Destination
	dstfoto, err := os.Create("uploads/survey_req/surat/" + surat.Filename)
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

	err = UpdateDataFotoPath("survey_request", "surat_penugasan", pathFotoFile, "id_transaksi_jual_sewa", int(lastId))
	if err != nil {
		return res, err
	}

	res.Status = http.StatusOK
	res.Message = "Berhasil memasukkan data"
	res.Data = dtSurveyReq

	defer db.DbClose(con)
	return res, nil
}

func GetAllSurveyReq() (Response, error) {
	var res Response
	type SurveyorAssignment struct {
		Id_transaksi_jual_sewa int    `json:"id_transaksi_jual_sewa"`
		User_id                int    `json:"user_id"`
		User_nama              string `json:"nama_lengkap"`
		Id_asset               int    `json:"id_asset"`
		Created_at             string `json:"created_at"`
		Surat_penugasan        string `json:"surat"`
		Dateline               string `json:"dateline"`
		Status_request         string `json:"status_request"`
		Status_verifikasi      string `json:"status_verifikasi"`
		Asset_nama             string `json:"asset_nama"`
		Asset_alamat           string `json:"asset_alamat"`
		Asset_titikkoordinat   string `json:"asset_titikkoordinat"`
	}
	var arrSurveyReq = []SurveyorAssignment{}
	var dtSurveyReq SurveyorAssignment

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	query := `
		SELECT sr.id_transaksi_jual_sewa, sr.user_id, u.nama_lengkap, sr.id_asset, sr.created_at, sr.surat_penugasan, sr.status_request, sr.status_verifikasi, sr.dateline, a.nama, a.alamat, a.titik_koordinat 
		FROM survey_request sr
		LEFT JOIN asset a ON sr.id_asset = a.id_asset
		LEFT JOIN user u ON sr.user_id = u.user_id
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
		err = result.Scan(&dtSurveyReq.Id_transaksi_jual_sewa, &dtSurveyReq.User_id, &dtSurveyReq.User_nama, &dtSurveyReq.Id_asset, &dtSurveyReq.Created_at, &dtSurveyReq.Surat_penugasan,
			&dtSurveyReq.Status_request, &dtSurveyReq.Status_verifikasi, &dtSurveyReq.Dateline, &dtSurveyReq.Asset_nama, &dtSurveyReq.Asset_alamat, &dtSurveyReq.Asset_titikkoordinat)
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

func GetAllSurveyReqDetailed() (Response, error) {
	var res Response
	type tempAllSurveyReqDetail struct {
		Id_transaksi_jual_sewa int    `json:"id_transaksi_jual_sewa"`
		Data_lengkap           string `json:"data_lengkap"`
		Nama_asset             string `json:"nama_asset"`
		Alamat                 string `json:"alamat"`
		Id_Surveyor            int    `json:"id_surveyor"`
		Nama_lengkap           string `json:"nama_lengkap"`
		Status_request         string `json:"status_request"`
		Status_verifikasi      string `json:"status_verifikasi"`
		Surat_penugasan        string `json:"surat"`
		Dateline               string `json:"dateline"`
	}
	var arrSurveyReq = []tempAllSurveyReqDetail{}

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	query := `
		SELECT sr.id_transaksi_jual_sewa, sr.data_lengkap,a.nama,a.alamat,
			s.suveyor_id,u.nama_lengkap,
			sr.status_request,sr.status_verifikasi,sr.surat_penugasan,sr.dateline
		FROM survey_request sr
		LEFT JOIN asset a ON sr.id_asset = a.id_asset
		LEFT JOIN user u ON sr.user_id = u.user_id
		LEFT JOIN surveyor s ON sr.user_id = s.user_id
		WHERE sr.status_submitted = 'Y'
		ORDER BY sr.dateline
	`
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

	for rows.Next() {
		var dtSurveyReq tempAllSurveyReqDetail
		err := rows.Scan(
			&dtSurveyReq.Id_transaksi_jual_sewa, &dtSurveyReq.Data_lengkap, &dtSurveyReq.Nama_asset, &dtSurveyReq.Alamat,
			&dtSurveyReq.Id_Surveyor, &dtSurveyReq.Nama_lengkap, &dtSurveyReq.Status_request,
			&dtSurveyReq.Status_verifikasi, &dtSurveyReq.Surat_penugasan, &dtSurveyReq.Dateline,
		)
		if err != nil {
			res.Status = 401
			res.Message = "rows scan gagal"
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

	// Modify the query to include usage_old, usage_new
	query := `
	SELECT sr.id_transaksi_jual_sewa, sr.user_id, sr.id_asset, sr.created_at, sr.dateline,
       sr.surat_penugasan, sr.status_request, sr.status_verifikasi, sr.data_lengkap, 
       sr.usage_old, sr.usage_new, sr.luas_old, sr.luas_new, sr.nilai_old, sr.nilai_new,
       sr.kondisi_old, sr.kondisi_new, sr.titik_koordinat_old, sr.titik_koordinat_new, 
       sr.batas_koordinat_old, sr.batas_koordinat_new, sr.tags_old, sr.tags_new, 
       a.nama, a.alamat, a.tipe
	FROM survey_request sr
	LEFT JOIN asset a ON sr.id_asset = a.id_asset
	WHERE sr.id_transaksi_jual_sewa = ?
	`

	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	// Define sql.NullString to handle NULL values from the database
	var usageOld, usageNew, tagsOld, tagsNew sql.NullString
	nId, _ := strconv.Atoi(surveyreq_id)
	err = stmt.QueryRow(nId).Scan(
		&dtSurveyReq.Id_transaksi_jual_sewa, &dtSurveyReq.User_id, &dtSurveyReq.Id_asset, &dtSurveyReq.Created_at,
		&dtSurveyReq.Dateline, &dtSurveyReq.Surat_penugasan, &dtSurveyReq.Status_request, &dtSurveyReq.Status_verifikasi,
		&dtSurveyReq.Data_lengkap, &usageOld, &usageNew, &dtSurveyReq.Luas_old, &dtSurveyReq.Luas_new,
		&dtSurveyReq.Nilai_old, &dtSurveyReq.Nilai_new, &dtSurveyReq.Kondisi_old, &dtSurveyReq.Kondisi_new,
		&dtSurveyReq.Titik_koordinat_old, &dtSurveyReq.Titik_koordinat_new, &dtSurveyReq.Batas_koordinat_old,
		&dtSurveyReq.Batas_koordinat_new, &tagsOld, &tagsNew, &dtSurveyReq.Nama_asset, &dtSurveyReq.Lokasi_asset, &dtSurveyReq.Tipe_asset)

	if err != nil {
		res.Status = 401
		res.Message = "exec gagal"
		res.Data = err.Error()
		return res, err
	}

	// Fetch names for old and new usage
	if usageOld.Valid {
		dtSurveyReq.Usage_old, err = fetchUsageNames(con, usageOld.String)
		if err != nil {
			res.Status = 401
			res.Message = fmt.Sprintf("Error fetching Usage_old: %v", err)
			res.Data = nil
			return res, err
		}
	} else {
		dtSurveyReq.Usage_old = []Kegunaan{}
	}

	if usageNew.Valid {
		dtSurveyReq.Usage_new, err = fetchUsageNames(con, usageNew.String)
		if err != nil {
			res.Status = 401
			res.Message = fmt.Sprintf("Error fetching Usage_new: %v", err)
			res.Data = nil
			return res, err
		}
	} else {
		dtSurveyReq.Usage_new = []Kegunaan{}
	}

	// Fetch names for old and new tags
	if tagsOld.Valid {
		dtSurveyReq.Tags_old, err = fetchTagNames(con, tagsOld.String)
		if err != nil {
			res.Status = 401
			res.Message = fmt.Sprintf("Error fetching Tags_old: %v", err)
			res.Data = nil
			return res, err
		}
	} else {
		dtSurveyReq.Tags_old = []Tags{}
	}

	if tagsNew.Valid {
		dtSurveyReq.Tags_new, err = fetchTagNames(con, tagsNew.String)
		if err != nil {
			res.Status = 401
			res.Message = fmt.Sprintf("Error fetching Tags_new: %v", err)
			res.Data = nil
			return res, err
		}
	} else {
		dtSurveyReq.Tags_new = []Tags{}
	}

	res.Status = http.StatusOK
	res.Message = "Berhasil mengambil data"
	res.Data = dtSurveyReq

	defer db.DbClose(con)
	return res, nil
}

func fetchUsageNames(con *sql.DB, usageIds string) ([]Kegunaan, error) {
	if strings.TrimSpace(usageIds) == "" {
		return []Kegunaan{}, nil
	}

	idList := strings.Split(usageIds, ",")
	var kegunaan []Kegunaan

	query := `
		SELECT id, nama
		FROM penggunaan
		WHERE id IN (` + strings.Join(idList, ",") + `)
	`

	rows, err := con.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var usage Kegunaan
		if err := rows.Scan(&usage.Id, &usage.Nama); err != nil {
			return nil, err
		}
		kegunaan = append(kegunaan, usage)
	}

	return kegunaan, nil
}

func fetchTagNames(con *sql.DB, tagIds string) ([]Tags, error) {
	if strings.TrimSpace(tagIds) == "" {
		return []Tags{}, nil
	}

	// Split the comma-separated string of tag IDs
	idList := strings.Split(tagIds, ",")
	var tags []Tags

	// Query to fetch names for the tag IDs
	query := `
		SELECT id, nama
		FROM tags
		WHERE id IN (` + strings.Join(idList, ",") + `)
	`

	rows, err := con.Query(query)
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

	query := `
	SELECT sr.id_transaksi_jual_sewa, sr.user_id, sr.id_asset, sr.created_at, sr.dateline,
       sr.surat_penugasan, sr.status_request, sr.status_verifikasi, sr.data_lengkap, 
       sr.usage_old, sr.usage_new, sr.luas_old, sr.luas_new, sr.nilai_old, sr.nilai_new,
       sr.kondisi_old, sr.kondisi_new, sr.titik_koordinat_old, sr.titik_koordinat_new, 
       sr.batas_koordinat_old, sr.batas_koordinat_new, sr.tags_old, sr.tags_new, 
       a.nama, a.alamat, a.tipe
	FROM survey_request sr
	LEFT JOIN asset a ON sr.id_asset = a.id_asset
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

	// Define sql.NullString to handle NULL values from the database
	var usageOld, usageNew, tagsOld, tagsNew sql.NullString
	nId, _ := strconv.Atoi(aset_id)
	err = stmt.QueryRow(nId).Scan(
		&dtSurveyReq.Id_transaksi_jual_sewa, &dtSurveyReq.User_id, &dtSurveyReq.Id_asset, &dtSurveyReq.Created_at,
		&dtSurveyReq.Dateline, &dtSurveyReq.Surat_penugasan, &dtSurveyReq.Status_request, &dtSurveyReq.Status_verifikasi,
		&dtSurveyReq.Data_lengkap, &usageOld, &usageNew, &dtSurveyReq.Luas_old, &dtSurveyReq.Luas_new,
		&dtSurveyReq.Nilai_old, &dtSurveyReq.Nilai_new, &dtSurveyReq.Kondisi_old, &dtSurveyReq.Kondisi_new,
		&dtSurveyReq.Titik_koordinat_old, &dtSurveyReq.Titik_koordinat_new, &dtSurveyReq.Batas_koordinat_old,
		&dtSurveyReq.Batas_koordinat_new, &tagsOld, &tagsNew, &dtSurveyReq.Nama_asset, &dtSurveyReq.Lokasi_asset, &dtSurveyReq.Tipe_asset)

	if err != nil {
		res.Status = 401
		res.Message = "exec gagal"
		res.Data = err.Error()
		return res, err
	}

	// Fetch names for old and new usage
	if usageOld.Valid {
		dtSurveyReq.Usage_old, err = fetchUsageNames(con, usageOld.String)
		if err != nil {
			res.Status = 401
			res.Message = fmt.Sprintf("Error fetching Usage_old: %v", err)
			res.Data = nil
			return res, err
		}
	} else {
		dtSurveyReq.Usage_old = []Kegunaan{}
	}

	if usageNew.Valid {
		dtSurveyReq.Usage_new, err = fetchUsageNames(con, usageNew.String)
		if err != nil {
			res.Status = 401
			res.Message = fmt.Sprintf("Error fetching Usage_new: %v", err)
			res.Data = nil
			return res, err
		}
	} else {
		dtSurveyReq.Usage_new = []Kegunaan{}
	}

	// Fetch names for old and new tags
	if tagsOld.Valid {
		dtSurveyReq.Tags_old, err = fetchTagNames(con, tagsOld.String)
		if err != nil {
			res.Status = 401
			res.Message = fmt.Sprintf("Error fetching Tags_old: %v", err)
			res.Data = nil
			return res, err
		}
	} else {
		dtSurveyReq.Tags_old = []Tags{}
	}

	if tagsNew.Valid {
		dtSurveyReq.Tags_new, err = fetchTagNames(con, tagsNew.String)
		if err != nil {
			res.Status = 401
			res.Message = fmt.Sprintf("Error fetching Tags_new: %v", err)
			res.Data = nil
			return res, err
		}
	} else {
		dtSurveyReq.Tags_new = []Tags{}
	}

	res.Status = http.StatusOK
	res.Message = "Berhasil mengambil data"
	res.Data = dtSurveyReq

	defer db.DbClose(con)
	return res, nil
}

func GetSurveyReqByAsetName(aset_name string) (Response, error) {
	var res Response
	var arrSurveyReq []SurveyRequest

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	query := `
	SELECT sr.*,a.nama 
	FROM survey_request sr
	JOIN asset a ON sr.id_asset = a.id_asset 
	WHERE a.nama LIKE ?
	`
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	likePattern := "%" + aset_name + "%"
	rows, err := stmt.Query(likePattern)
	if err != nil {
		res.Status = 401
		res.Message = "query gagal"
		res.Data = err.Error()
		return res, err
	}
	defer rows.Close()

	for rows.Next() {
		var dtSurveyReq SurveyRequest
		err := rows.Scan(
			&dtSurveyReq.Id_transaksi_jual_sewa, &dtSurveyReq.User_id, &dtSurveyReq.Id_asset,
			&dtSurveyReq.Created_at, &dtSurveyReq.Dateline, &dtSurveyReq.Status_request,
			&dtSurveyReq.Status_verifikasi, &dtSurveyReq.Data_lengkap, &dtSurveyReq.Usage_old,
			&dtSurveyReq.Usage_new, &dtSurveyReq.Luas_old, &dtSurveyReq.Luas_new,
			&dtSurveyReq.Nilai_old, &dtSurveyReq.Nilai_new, &dtSurveyReq.Kondisi_old,
			&dtSurveyReq.Kondisi_new, &dtSurveyReq.Batas_koordinat_old, &dtSurveyReq.Batas_koordinat_new,
			&dtSurveyReq.Tags_old, &dtSurveyReq.Tags_new, &dtSurveyReq.Nama_asset,
		)
		if err != nil {
			res.Status = 401
			res.Message = "exec gagal"
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

func GetAllSurveyReqByUserId(user_id string) (Response, error) {
	var res Response
	type SurveyorAssignment struct {
		Id_transaksi_jual_sewa int    `json:"id_transaksi_jual_sewa"`
		User_id                int    `json:"user_id"`
		Id_asset               int    `json:"id_asset"`
		Created_at             string `json:"created_at"`
		Surat_penugasan        string `json:"surat_penugasan"`
		Dateline               string `json:"dateline"`
		Status_request         string `json:"status_request"`
		Status_verifikasi      string `json:"status_verifikasi"`
		Asset_nama             string `json:"asset_nama"`
		Asset_alamat           string `json:"asset_alamat"`
		Asset_titikkoordinat   string `json:"asset_titikkoordinat"`
	}
	type tempSurvAssignment struct {
		OngoingAssignment  []SurveyorAssignment `json:"ongoing_assignment"`
		FinishedAssignment []SurveyorAssignment `json:"finished_assignment"`
	}
	var arrSurveyReq = []SurveyorAssignment{}
	var dtSurveyReq SurveyorAssignment

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	query := `
		SELECT sr.id_transaksi_jual_sewa, sr.user_id, sr.id_asset, sr.created_at, sr.status_request, sr.status_verifikasi, sr.surat_penugasan, sr.dateline, a.nama, a.alamat, a.titik_koordinat 
		FROM survey_request sr
		JOIN asset a ON sr.id_asset = a.id_asset
		WHERE sr.user_id = ? AND (sr.status_request = 'O' OR sr.status_request = 'R')
	`
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	nId, _ := strconv.Atoi(user_id)
	result, err := stmt.Query(nId)
	if err != nil {
		res.Status = 401
		res.Message = "exec gagal"
		res.Data = err.Error()
		return res, err
	}
	defer result.Close()
	for result.Next() {
		err = result.Scan(&dtSurveyReq.Id_transaksi_jual_sewa, &dtSurveyReq.User_id, &dtSurveyReq.Id_asset, &dtSurveyReq.Created_at,
			&dtSurveyReq.Status_request, &dtSurveyReq.Status_verifikasi, &dtSurveyReq.Surat_penugasan, &dtSurveyReq.Dateline, &dtSurveyReq.Asset_nama, &dtSurveyReq.Asset_alamat, &dtSurveyReq.Asset_titikkoordinat)
		if err != nil {
			res.Status = 401
			res.Message = "rows scan"
			res.Data = err.Error()
			return res, err
		}
		arrSurveyReq = append(arrSurveyReq, dtSurveyReq)
	}
	var survey_assignment tempSurvAssignment
	survey_assignment.OngoingAssignment = arrSurveyReq
	fmt.Println("ambil finished survey request")
	// finished assignment
	var arrSurveyReqFinished = []SurveyorAssignment{}
	var dtSurveyReqFinished SurveyorAssignment
	queryfinished := `
		SELECT sr.id_transaksi_jual_sewa, sr.user_id, sr.id_asset, sr.created_at, sr.status_request, sr.status_verifikasi, sr.surat_penugasan, sr.dateline, a.nama, a.alamat, a.titik_koordinat 
		FROM survey_request sr
		JOIN asset a ON sr.id_asset = a.id_asset
		WHERE sr.user_id = ? AND sr.status_request = 'F'
	`
	stmtfinished, err := con.Prepare(queryfinished)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmtfinished.Close()

	resultfinished, err := stmtfinished.Query(nId)
	if err != nil {
		res.Status = 401
		res.Message = "exec gagal"
		res.Data = err.Error()
		return res, err
	}
	defer resultfinished.Close()
	for resultfinished.Next() {
		err = resultfinished.Scan(&dtSurveyReq.Id_transaksi_jual_sewa, &dtSurveyReq.User_id, &dtSurveyReq.Id_asset, &dtSurveyReq.Created_at,
			&dtSurveyReq.Status_request, &dtSurveyReq.Status_verifikasi, &dtSurveyReq.Surat_penugasan, &dtSurveyReq.Dateline, &dtSurveyReq.Asset_nama, &dtSurveyReq.Asset_alamat, &dtSurveyReq.Asset_titikkoordinat)
		if err != nil {
			res.Status = 401
			res.Message = "rows scan"
			res.Data = err.Error()
			return res, err
		}
		arrSurveyReqFinished = append(arrSurveyReqFinished, dtSurveyReqFinished)
	}

	survey_assignment.FinishedAssignment = arrSurveyReqFinished

	res.Status = http.StatusOK
	res.Message = "Berhasil mengambil data"
	res.Data = survey_assignment

	defer db.DbClose(con)
	return res, nil
}

func GetAllOngoingSurveyReqByUserId(user_id string) (Response, error) {
	var res Response
	type SurveyorAssignment struct {
		Id_transaksi_jual_sewa int    `json:"id_transaksi_jual_sewa"`
		User_id                int    `json:"user_id"`
		Id_asset               int    `json:"id_asset"`
		Created_at             string `json:"created_at"`
		Dateline               string `json:"dateline"`
		Status_request         string `json:"status_request"`
		Status_verifikasi      string `json:"status_verifikasi"`
		Asset_nama             string `json:"asset_nama"`
		Asset_alamat           string `json:"asset_alamat"`
		Asset_titikkoordinat   string `json:"asset_titikkoordinat"`
	}
	var arrSurveyReq = []SurveyorAssignment{}
	var dtSurveyReq SurveyorAssignment

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	query := `
		SELECT sr.id_transaksi_jual_sewa, sr.user_id, sr.id_asset, sr.created_at, sr.status_request, sr.status_verifikasi, sr.dateline, a.nama, a.alamat, a.titik_koordinat 
		FROM survey_request sr
		JOIN asset a ON sr.id_asset = a.id_asset
		WHERE sr.user_id = ? AND (sr.status_request = 'O' OR sr.status_request = 'R')
	`
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	nId, _ := strconv.Atoi(user_id)
	result, err := stmt.Query(nId)
	if err != nil {
		res.Status = 401
		res.Message = "exec gagal"
		res.Data = err.Error()
		return res, err
	}
	defer result.Close()
	for result.Next() {
		err = result.Scan(&dtSurveyReq.Id_transaksi_jual_sewa, &dtSurveyReq.User_id, &dtSurveyReq.Id_asset, &dtSurveyReq.Created_at,
			&dtSurveyReq.Status_request, &dtSurveyReq.Status_verifikasi, &dtSurveyReq.Dateline, &dtSurveyReq.Asset_nama, &dtSurveyReq.Asset_alamat, &dtSurveyReq.Asset_titikkoordinat)
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

func GetAllFinishedSurveyReqByUserId(user_id string) (Response, error) {
	var res Response
	type SurveyorAssignment struct {
		Id_transaksi_jual_sewa int    `json:"id_transaksi_jual_sewa"`
		User_id                int    `json:"user_id"`
		Id_asset               int    `json:"id_asset"`
		Created_at             string `json:"created_at"`
		Dateline               string `json:"dateline"`
		Status_request         string `json:"status_request"`
		Status_verifikasi      string `json:"status_verifikasi"`
		Asset_nama             string `json:"asset_nama"`
		Asset_alamat           string `json:"asset_alamat"`
		Asset_titikkoordinat   string `json:"asset_titikkoordinat"`
	}
	var arrSurveyReq = []SurveyorAssignment{}
	var dtSurveyReq SurveyorAssignment

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	query := `
		SELECT sr.id_transaksi_jual_sewa, sr.user_id, sr.id_asset, sr.created_at, sr.status_request, sr.status_verifikasi, sr.dateline, a.nama, a.alamat, a.titik_koordinat 
		FROM survey_request sr
		JOIN asset a ON sr.id_asset = a.id_asset
		WHERE sr.user_id = ? AND sr.status_request = 'F'
	`
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	nId, _ := strconv.Atoi(user_id)
	result, err := stmt.Query(nId)
	if err != nil {
		res.Status = 401
		res.Message = "exec gagal"
		res.Data = err.Error()
		return res, err
	}
	defer result.Close()
	for result.Next() {
		err = result.Scan(&dtSurveyReq.Id_transaksi_jual_sewa, &dtSurveyReq.User_id, &dtSurveyReq.Id_asset, &dtSurveyReq.Created_at,
			&dtSurveyReq.Status_request, &dtSurveyReq.Status_verifikasi, &dtSurveyReq.Dateline, &dtSurveyReq.Asset_nama, &dtSurveyReq.Asset_alamat, &dtSurveyReq.Asset_titikkoordinat)
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

func SubmitSurveyReqById(surveyreq string) (Response, error) {
	var res Response
	type TempSubmitSurveyReq struct {
		Id              int     `json:"id"`
		Usage           string  `json:"usage"`
		Luas            float64 `json:"luas"`
		Nilai           float64 `json:"nilai"`
		Kondisi         string  `json:"kondisi"`
		Titik_koordinat string  `json:"titik_koordinat"`
		Batas_koordinat string  `json:"batas_koordinat"`
		Tags            string  `json:"tags"`
	}

	var tempsubmitsurveyreq TempSubmitSurveyReq
	var datalengkap string

	err := json.Unmarshal([]byte(surveyreq), &tempsubmitsurveyreq)
	if err != nil {
		res.Status = 401
		res.Message = "gagal decode json"
		res.Data = err.Error()
		return res, err
	}

	if tempsubmitsurveyreq.Usage != "" && tempsubmitsurveyreq.Luas > 0 &&
		tempsubmitsurveyreq.Nilai > 0 && tempsubmitsurveyreq.Kondisi != "" &&
		tempsubmitsurveyreq.Titik_koordinat != "" && tempsubmitsurveyreq.Batas_koordinat != "" &&
		tempsubmitsurveyreq.Tags != "" {
		datalengkap = "Y"
	} else {
		datalengkap = "N"
	}
	fmt.Println("usage", tempsubmitsurveyreq.Usage)
	fmt.Println("Luas", tempsubmitsurveyreq.Luas)
	fmt.Println("Nilai", tempsubmitsurveyreq.Nilai)
	fmt.Println("Kondisi", tempsubmitsurveyreq.Kondisi)
	fmt.Println("Titik_koordinat", tempsubmitsurveyreq.Titik_koordinat)
	fmt.Println("Batas_koordinat", tempsubmitsurveyreq.Batas_koordinat)
	fmt.Println("Tags", tempsubmitsurveyreq.Tags)

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	query := "UPDATE survey_request SET `usage_new` = ?, luas_new = ?, nilai_new = ?, kondisi_new = ?, titik_koordinat_new = ?, batas_koordinat_new = ?, tags_new = ?,data_lengkap = ?,status_submitted='Y' WHERE id_transaksi_jual_sewa = ?"
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(
		tempsubmitsurveyreq.Usage,
		tempsubmitsurveyreq.Luas,
		tempsubmitsurveyreq.Nilai,
		tempsubmitsurveyreq.Kondisi,
		tempsubmitsurveyreq.Titik_koordinat,
		tempsubmitsurveyreq.Batas_koordinat,
		tempsubmitsurveyreq.Tags,
		datalengkap,
		tempsubmitsurveyreq.Id,
	)

	fmt.Println("datalengkap:", datalengkap)
	if err != nil {
		res.Status = 401
		res.Message = "Failed to execute statement"
		res.Data = err.Error()
		return res, err
	}

	// Check if any rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		res.Status = 401
		res.Message = "Failed to retrieve rows affected"
		res.Data = err.Error()
		return res, err
	}

	if rowsAffected == 0 {
		res.Status = 404
		res.Message = "No records updated"
		res.Data = nil
		return res, nil
	}

	res.Status = http.StatusOK
	res.Message = "Berhasil mengupdate data"
	res.Data = result

	defer db.DbClose(con)
	return res, nil
}

func GetSurveyReqHistoryByUserId(user_id string) (Response, error) {
	var res Response

	return res, nil
}

// CRUD transaction_request ============================================================================
func CreateTranReq(proposal *multipart.FileHeader, id_asset, user_id, perusahaan_id, nama_progress, tgl_meeting, waktu_meeting, lokasi_meeting, deskripsi string) (Response, error) {
	var res Response
	var dtTranReq = TransactionRequest{}
	dtTranReq.Perusahaan_id, _ = strconv.Atoi(perusahaan_id)
	dtTranReq.User_id, _ = strconv.Atoi(user_id)
	dtTranReq.Id_asset, _ = strconv.Atoi(id_asset)
	dtTranReq.Nama_progress = nama_progress
	dtTranReq.Tgl_meeting = tgl_meeting
	dtTranReq.Waktu_meeting = waktu_meeting
	dtTranReq.Lokasi_meeting = lokasi_meeting
	dtTranReq.Deskripsi = deskripsi

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	query := "INSERT INTO transaction_request (id_asset, user_id, perusahaan_id, nama_progress, tgl_meeting, waktu_meeting, lokasi_meeting, deskripsi) VALUES (?,?,?,?,?,?,?,?)"
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(dtTranReq.Id_asset, dtTranReq.User_id, dtTranReq.Perusahaan_id, dtTranReq.Nama_progress, dtTranReq.Tgl_meeting, dtTranReq.Waktu_meeting, dtTranReq.Lokasi_meeting, dtTranReq.Deskripsi)
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

	// insert file proposal
	proposal.Filename = strconv.Itoa(dtTranReq.Id_transaksi_jual_sewa) + "_" + proposal.Filename
	pathFotoFile := "uploads/transaction/" + proposal.Filename
	//source
	srcfoto, err := proposal.Open()
	if err != nil {
		log.Println(err.Error())
		return res, err
	}
	defer srcfoto.Close()

	// Destination
	dstfoto, err := os.Create("uploads/transaction/" + proposal.Filename)
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

	err = UpdateDataFotoPath("transaction_request", "proposal", pathFotoFile, "id_transaksi_jual_sewa", dtTranReq.Id_transaksi_jual_sewa)
	if err != nil {
		return res, err
	}
	dtTranReq.Proposal = pathFotoFile

	res.Status = http.StatusOK
	res.Message = "Berhasil memasukkan data"
	res.Data = dtTranReq

	defer db.DbClose(con)
	return res, nil
}

func GetAllTranReq() (Response, error) {
	var res Response
	var arrTranReq = []TransactionRequest{}

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	query := `
	SELECT tr.id_transaksi_jual_sewa, tr.perusahaan_id, tr.user_id, u.username, u.nama_lengkap, tr.id_asset, 
		a.nama, tr.status, tr.nama_progress, tr.proposal, IFNULL(tr.tgl_meeting, ''), IFNULL(tr.waktu_meeting, ''),tr.lokasi_meeting, 
		tr.deskripsi, tr.alasan, IFNULL(tr.tgl_dateline, ''), tr.created_at 
	FROM transaction_request tr
	LEFT JOIN asset a ON tr.id_asset = a.id_asset
	LEFT JOIN user u ON tr.user_id = u.user_id
	ORDER BY tr.created_at DESC
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
		var dtTranReq TransactionRequest
		err = result.Scan(&dtTranReq.Id_transaksi_jual_sewa, &dtTranReq.Perusahaan_id,
			&dtTranReq.User_id, &dtTranReq.Username, &dtTranReq.Nama_lengkap, &dtTranReq.Id_asset, &dtTranReq.Nama_aset, &dtTranReq.Status, &dtTranReq.Nama_progress,
			&dtTranReq.Proposal, &dtTranReq.Tgl_meeting, &dtTranReq.Waktu_meeting, &dtTranReq.Lokasi_meeting,
			&dtTranReq.Deskripsi, &dtTranReq.Alasan, &dtTranReq.Tgl_dateline,
			&dtTranReq.Created_at)
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

	query := `
	SELECT tr.id_transaksi_jual_sewa, tr.perusahaan_id, tr.user_id, u.username, u.nama_lengkap, tr.id_asset, 
		a.nama, tr.status, tr.nama_progress, tr.proposal, IFNULL(tr.tgl_meeting, ''), IFNULL(tr.waktu_meeting, ''),tr.lokasi_meeting, 
		tr.deskripsi, tr.alasan, IFNULL(tr.tgl_dateline, ''), tr.created_at 
	FROM transaction_request tr
	LEFT JOIN asset a ON tr.id_asset = a.id_asset
	LEFT JOIN user u ON tr.user_id = u.user_id
	WHERE id_transaksi_jual_sewa = ?
	ORDER BY tr.created_at DESC
	`
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	nId, _ := strconv.Atoi(tranreq_id)
	err = stmt.QueryRow(nId).Scan(&dtTranReq.Id_transaksi_jual_sewa, &dtTranReq.Perusahaan_id,
		&dtTranReq.User_id, &dtTranReq.Username, &dtTranReq.Nama_lengkap, &dtTranReq.Id_asset, &dtTranReq.Nama_aset, &dtTranReq.Status, &dtTranReq.Nama_progress,
		&dtTranReq.Proposal, &dtTranReq.Tgl_meeting, &dtTranReq.Waktu_meeting, &dtTranReq.Lokasi_meeting,
		&dtTranReq.Deskripsi, &dtTranReq.Alasan, &dtTranReq.Tgl_dateline,
		&dtTranReq.Created_at)
	if err != nil {
		if err == sql.ErrNoRows {
			res.Status = 404
			res.Message = "Data tidak ditemukan"
		} else {
			res.Status = 401
			res.Message = "Gagal scan row"
		}
		res.Data = err.Error()
		return res, err
	}

	res.Status = http.StatusOK
	res.Message = "Berhasil mengambil data"
	res.Data = dtTranReq

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

	query := `
	SELECT tr.id_transaksi_jual_sewa, tr.perusahaan_id, p.lokasi, tr.user_id, u.username, u.nama_lengkap, tr.id_asset, 
		a.nama, tr.status, tr.nama_progress, tr.proposal, IFNULL(tr.tgl_meeting, ''), tr.lokasi_meeting, 
		tr.deskripsi, tr.alasan, IFNULL(tr.tgl_dateline, ''), tr.created_at 
	FROM transaction_request tr
	LEFT JOIN perusahaan p ON tr.perusahaan_id = p.perusahaan_id
	LEFT JOIN user u ON tr.user_id = u.user_id
	LEFT JOIN asset a ON tr.id_asset = a.id_asset
	WHERE tr.user_id = ?
	`
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
		err = rows.Scan(&tranReq.Id_transaksi_jual_sewa, &tranReq.Perusahaan_id, &tranReq.Lokasi_perusahaan,
			&tranReq.User_id, &tranReq.Username, &tranReq.Nama_lengkap, &tranReq.Id_asset, &tranReq.Nama_aset, &tranReq.Status, &tranReq.Nama_progress,
			&tranReq.Proposal, &tranReq.Tgl_meeting, &tranReq.Lokasi_meeting,
			&tranReq.Deskripsi, &tranReq.Alasan, &tranReq.Tgl_dateline,
			&tranReq.Created_at)
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

func GetTranReqByPerusahaanId(perusahaan_id string) (Response, error) {
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

	query := `
	SELECT id_transaksi_jual_sewa, perusahaan_id, user_id, id_asset, 
		status, nama_progress, proposal, IFNULL(tgl_meeting, ''), lokasi_meeting, 
		deskripsi, alasan, IFNULL(tgl_dateline, ''), created_at 
	FROM transaction_request
	WHERE perusahaan_id = ?
	`
	rows, err := con.Query(query, perusahaan_id)
	if err != nil {
		res.Status = 401
		res.Message = "Failed to execute query"
		res.Data = err.Error()
		return res, err
	}
	defer rows.Close()

	for rows.Next() {
		var tranReq TransactionRequest
		err = rows.Scan(&tranReq.Id_transaksi_jual_sewa, &tranReq.Perusahaan_id,
			&tranReq.User_id, &tranReq.Id_asset, &tranReq.Status, &tranReq.Nama_progress,
			&tranReq.Proposal, &tranReq.Tgl_meeting, &tranReq.Lokasi_meeting,
			&tranReq.Deskripsi, &tranReq.Alasan, &tranReq.Tgl_dateline,
			&tranReq.Created_at)
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

func GetAllUserTransaction() (Response, error) {
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

	query := `
	SELECT tr.id_transaksi_jual_sewa, tr.perusahaan_id, p.lokasi, tr.user_id, u.username, u.nama_lengkap, tr.id_asset, 
		a.nama, tr.status, tr.nama_progress, tr.proposal, IFNULL(tr.tgl_meeting, ''), tr.lokasi_meeting, 
		tr.deskripsi, tr.alasan, IFNULL(tr.tgl_dateline, ''), tr.created_at 
	FROM transaction_request tr
	LEFT JOIN perusahaan p ON tr.perusahaan_id = p.perusahaan_id
	LEFT JOIN user u ON tr.user_id = u.user_id
	LEFT JOIN asset a ON tr.id_asset = a.id_asset
	WHERE tr.status = 'A'
	`
	rows, err := con.Query(query)
	if err != nil {
		res.Status = 401
		res.Message = "Failed to execute query"
		res.Data = err.Error()
		return res, err
	}
	defer rows.Close()

	for rows.Next() {
		var tranReq TransactionRequest
		err = rows.Scan(&tranReq.Id_transaksi_jual_sewa, &tranReq.Perusahaan_id, &tranReq.Lokasi_perusahaan,
			&tranReq.User_id, &tranReq.Username, &tranReq.Nama_lengkap, &tranReq.Id_asset, &tranReq.Nama_aset, &tranReq.Status, &tranReq.Nama_progress,
			&tranReq.Proposal, &tranReq.Tgl_meeting, &tranReq.Lokasi_meeting,
			&tranReq.Deskripsi, &tranReq.Alasan, &tranReq.Tgl_dateline,
			&tranReq.Created_at)
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

func SendProposal(proposal *multipart.FileHeader, id_asset, user_id, perusahaan_id, deskripsi string) (Response, error) {
	var res Response
	var dtTranReq = TransactionRequest{}
	dtTranReq.Perusahaan_id, _ = strconv.Atoi(perusahaan_id)
	dtTranReq.User_id, _ = strconv.Atoi(user_id)
	dtTranReq.Id_asset, _ = strconv.Atoi(id_asset)
	dtTranReq.Deskripsi = deskripsi

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	query := "INSERT INTO transaction_request (id_asset, user_id, perusahaan_id, deskripsi) VALUES (?,?,?,?)"
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(dtTranReq.Id_asset, dtTranReq.User_id, dtTranReq.Perusahaan_id, dtTranReq.Deskripsi)
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

	// insert file proposal
	proposal.Filename = strconv.Itoa(dtTranReq.Id_transaksi_jual_sewa) + "_" + proposal.Filename
	pathFotoFile := "uploads/transaction/" + proposal.Filename
	//source
	srcfoto, err := proposal.Open()
	if err != nil {
		log.Println(err.Error())
		return res, err
	}
	defer srcfoto.Close()

	// Destination
	dstfoto, err := os.Create("uploads/transaction/" + proposal.Filename)
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

	err = UpdateDataFotoPath("transaction_request", "proposal", pathFotoFile, "id_transaksi_jual_sewa", dtTranReq.Id_transaksi_jual_sewa)
	if err != nil {
		return res, err
	}
	dtTranReq.Proposal = pathFotoFile

	res.Status = http.StatusOK
	res.Message = "Berhasil memasukkan data"
	res.Data = dtTranReq

	defer db.DbClose(con)
	return res, nil
}

// func UserManagementGetMeetingByUserId(user_id string) (Response, error) {
// 	var res Response
// 	var assetsMap = make(map[int][]Progress)
// 	var grupAsset []GrupAsset

// 	con, err := db.DbConnection()
// 	if err != nil {
// 		res.Status = 401
// 		res.Message = "gagal membuka database"
// 		res.Data = err.Error()
// 		return res, err
// 	}

// 	query := `
// 	SELECT p.id,p.user_id,p.id_asset,p.perusahaan_id,p.status,p.nama,p.proposal,IFNULL(p.tanggal_meeting,""),IFNULL(p.waktu_meeting,""),
// 	IFNULL(p.tempat_meeting,""),IFNULL(p.waktu_mulai_meeting,""),IFNULL(p.waktu_selesai_meeting,""),IFNULL(p.notes,""),IFNULL(p.file,""),IFNULL(p.tipe_file,""),a.nama
// 	FROM progress p
// 	LEFT JOIN asset a ON p.id_asset = a.id_asset
// 	WHERE p.user_id = ?
// 	`

// 	nId, _ := strconv.Atoi(user_id)
// 	rows, err := con.Query(query, nId)
// 	if err != nil {
// 		res.Status = 401
// 		res.Message = "Failed to execute query"
// 		res.Data = err.Error()
// 		return res, err
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		var progress Progress
// 		err = rows.Scan(&progress.Id, &progress.User_id, &progress.Id_asset, &progress.Perusahaan_id, &progress.Status, &progress.Nama, &progress.Proposal, &progress.Tanggal_meeting, &progress.Waktu_meeting,
// 			&progress.Tempat_meeting, &progress.Waktu_mulai_meeting, &progress.Waktu_selesai_meeting, &progress.Notes, &progress.Dokumen, &progress.Tipe_dokumen,
// 			&progress.Nama_asset)
// 		if err != nil {
// 			res.Status = 401
// 			res.Message = "Failed to scan row"
// 			res.Data = err.Error()
// 			return res, err
// 		}
// 		assetsMap[progress.Id_asset] = append(assetsMap[progress.Id_asset], progress)
// 	}
// 	if err = rows.Err(); err != nil {
// 		res.Status = 401
// 		res.Message = "Failed during row iteration"
// 		res.Data = err.Error()
// 		return res, err
// 	}

// 	for assetId, progressList := range assetsMap {
// 		grupAsset = append(grupAsset, GrupAsset{
// 			Id_asset:       assetId,
// 			Asset_name:     progressList[0].Nama_asset,
// 			Semua_progress: progressList,
// 		})
// 	}
// 	if err = rows.Err(); err != nil {
// 		res.Status = 401
// 		res.Message = "Failed during row iteration"
// 		res.Data = err.Error()
// 		return res, err
// 	}

// 	res.Status = http.StatusOK
// 	res.Message = "Berhasil mengambil data"
// 	res.Data = grupAsset

// 	defer db.DbClose(con)
// 	return res, nil
// }

func UserManagementGetMeetingByUserId(user_id string) (Response, error) {
	var res Response
	var assetsMap = make(map[int][]Progress)
	var grupAsset []GrupAsset

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	query := `
	SELECT p.id, p.user_id, p.id_asset, p.perusahaan_id, p.status, p.nama, p.proposal, 
	IFNULL(p.tanggal_meeting,""), IFNULL(p.waktu_meeting,""), IFNULL(p.tempat_meeting,""), 
	IFNULL(p.waktu_mulai_meeting,""), IFNULL(p.waktu_selesai_meeting,""), IFNULL(p.notes,""), 
	IFNULL(p.file,""), IFNULL(p.tipe_file,""), a.nama
	FROM progress p
	LEFT JOIN asset a ON p.id_asset = a.id_asset
	WHERE p.user_id = ?
	`

	nId, _ := strconv.Atoi(user_id)
	rows, err := con.Query(query, nId)
	if err != nil {
		res.Status = 401
		res.Message = "Failed to execute query"
		res.Data = err.Error()
		return res, err
	}
	defer rows.Close()

	for rows.Next() {
		var progress Progress
		err = rows.Scan(&progress.Id, &progress.User_id, &progress.Id_asset, &progress.Perusahaan_id, &progress.Status, &progress.Nama, &progress.Proposal,
			&progress.Tanggal_meeting, &progress.Waktu_meeting, &progress.Tempat_meeting, &progress.Waktu_mulai_meeting, &progress.Waktu_selesai_meeting, &progress.Notes,
			&progress.Dokumen, &progress.Tipe_dokumen, &progress.Nama_asset)
		if err != nil {
			res.Status = 401
			res.Message = "Failed to scan row"
			res.Data = err.Error()
			return res, err
		}
		assetsMap[progress.Id_asset] = append(assetsMap[progress.Id_asset], progress)
	}
	if err = rows.Err(); err != nil {
		res.Status = 401
		res.Message = "Failed during row iteration"
		res.Data = err.Error()
		return res, err
	}

	for assetId, progressList := range assetsMap {
		grupAsset = append(grupAsset, GrupAsset{
			Id_asset:       assetId,
			Asset_name:     progressList[0].Nama_asset,
			User_id:        progressList[0].User_id,
			Perusahaan_id:  progressList[0].Perusahaan_id,
			Semua_progress: []Progress{},
		})
		for _, prog := range progressList {
			grupAsset[len(grupAsset)-1].Semua_progress = append(grupAsset[len(grupAsset)-1].Semua_progress, Progress{
				Id:                    prog.Id,
				Status:                prog.Status,
				Nama:                  prog.Nama,
				Proposal:              prog.Proposal,
				Tanggal_meeting:       prog.Tanggal_meeting,
				Waktu_meeting:         prog.Waktu_meeting,
				Tempat_meeting:        prog.Tempat_meeting,
				Waktu_mulai_meeting:   prog.Waktu_mulai_meeting,
				Waktu_selesai_meeting: prog.Waktu_selesai_meeting,
				Notes:                 prog.Notes,
				Dokumen:               prog.Dokumen,
				Tipe_dokumen:          prog.Tipe_dokumen,
			})
		}
	}

	res.Status = http.StatusOK
	res.Message = "Berhasil mengambil data"
	res.Data = grupAsset

	defer db.DbClose(con)
	return res, nil
}

func UserManagementGetMeetingByPerusahaanId(perusahaan_id string) (Response, error) {
	var res Response
	var assetsMap = make(map[int][]Progress)
	var grupAsset []GrupAsset

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	query := `
	SELECT p.id, p.user_id, p.id_asset, p.perusahaan_id, p.status, p.nama, p.proposal, 
	IFNULL(p.tanggal_meeting,""), IFNULL(p.waktu_meeting,""), IFNULL(p.tempat_meeting,""), 
	IFNULL(p.waktu_mulai_meeting,""), IFNULL(p.waktu_selesai_meeting,""), IFNULL(p.notes,""), 
	IFNULL(p.file,""), IFNULL(p.tipe_file,""), a.nama
	FROM progress p
	LEFT JOIN asset a ON p.id_asset = a.id_asset
	WHERE p.perusahaan_id = ?
	`

	nId, _ := strconv.Atoi(perusahaan_id)
	rows, err := con.Query(query, nId)
	if err != nil {
		res.Status = 401
		res.Message = "Failed to execute query"
		res.Data = err.Error()
		return res, err
	}
	defer rows.Close()

	for rows.Next() {
		var progress Progress
		err = rows.Scan(&progress.Id, &progress.User_id, &progress.Id_asset, &progress.Perusahaan_id, &progress.Status, &progress.Nama, &progress.Proposal,
			&progress.Tanggal_meeting, &progress.Waktu_meeting, &progress.Tempat_meeting, &progress.Waktu_mulai_meeting, &progress.Waktu_selesai_meeting, &progress.Notes,
			&progress.Dokumen, &progress.Tipe_dokumen, &progress.Nama_asset)
		if err != nil {
			res.Status = 401
			res.Message = "Failed to scan row"
			res.Data = err.Error()
			return res, err
		}
		assetsMap[progress.Id_asset] = append(assetsMap[progress.Id_asset], progress)
	}
	if err = rows.Err(); err != nil {
		res.Status = 401
		res.Message = "Failed during row iteration"
		res.Data = err.Error()
		return res, err
	}

	for assetId, progressList := range assetsMap {
		grupAsset = append(grupAsset, GrupAsset{
			Id_asset:       assetId,
			Asset_name:     progressList[0].Nama_asset,
			User_id:        progressList[0].User_id,
			Perusahaan_id:  progressList[0].Perusahaan_id,
			Semua_progress: []Progress{},
		})
		for _, prog := range progressList {
			grupAsset[len(grupAsset)-1].Semua_progress = append(grupAsset[len(grupAsset)-1].Semua_progress, Progress{
				Id:                    prog.Id,
				Status:                prog.Status,
				Nama:                  prog.Nama,
				Proposal:              prog.Proposal,
				Tanggal_meeting:       prog.Tanggal_meeting,
				Waktu_meeting:         prog.Waktu_meeting,
				Tempat_meeting:        prog.Tempat_meeting,
				Waktu_mulai_meeting:   prog.Waktu_mulai_meeting,
				Waktu_selesai_meeting: prog.Waktu_selesai_meeting,
				Notes:                 prog.Notes,
				Dokumen:               prog.Dokumen,
				Tipe_dokumen:          prog.Tipe_dokumen,
			})
		}
	}

	res.Status = http.StatusOK
	res.Message = "Berhasil mengambil data"
	res.Data = grupAsset

	defer db.DbClose(con)
	return res, nil
}

func CreateMeeting(dokumen *multipart.FileHeader, id, tanggal_meeting, waktu_meeting, tempat_meeting, waktu_mulai_meeting, waktu_selesai_meeting, notes, tipe_dokumen string) (Response, error) {
	var res Response
	var dtMeeting Progress

	if id == "" && tanggal_meeting == "" && waktu_meeting == "" && tempat_meeting == "" && waktu_mulai_meeting == "" && waktu_selesai_meeting == "" && notes == "" {
		res.Status = 400
		res.Message = "Required fields are missing"
		res.Data = "data tidak lengkap"
		return res, errors.New("required fields are missing")
	}
	if tipe_dokumen != "L" && tipe_dokumen != "C" && tipe_dokumen != "A" {
		tipe_dokumen = ""
		dtMeeting.Tipe_dokumen = ""
	}

	dtMeeting.Id, _ = strconv.Atoi(id)
	dtMeeting.Tanggal_meeting = tanggal_meeting
	dtMeeting.Waktu_meeting = waktu_meeting
	dtMeeting.Tempat_meeting = tempat_meeting
	dtMeeting.Waktu_mulai_meeting = waktu_mulai_meeting
	dtMeeting.Waktu_selesai_meeting = waktu_selesai_meeting
	dtMeeting.Notes = notes
	dtMeeting.Tipe_dokumen = tipe_dokumen

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	query := `
		UPDATE progress SET tanggal_meeting = ?, waktu_meeting = ?, tempat_meeting = ?, 
		waktu_mulai_meeting=?, waktu_selesai_meeting = ?, notes = ?, tipe_file = ?, data_lengkap = 'Y'
		WHERE id = ?`
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		dtMeeting.Tanggal_meeting,
		dtMeeting.Waktu_meeting,
		dtMeeting.Tempat_meeting,
		dtMeeting.Waktu_mulai_meeting,
		dtMeeting.Waktu_selesai_meeting,
		dtMeeting.Notes,
		dtMeeting.Tipe_dokumen,
		dtMeeting.Id,
	)
	if err != nil {
		res.Status = 401
		res.Message = "exec gagal"
		res.Data = err.Error()
		return res, err
	}

	// tambah file
	//  ======================================================

	dokumen.Filename = id + "_" + dokumen.Filename
	pathFotoFile := "uploads/progress/" + dokumen.Filename
	//source
	srcfoto, err := dokumen.Open()
	if err != nil {
		log.Println(err.Error())
		return res, err
	}
	defer srcfoto.Close()

	// Destination
	dstfoto, err := os.Create("uploads/progress/" + dokumen.Filename)
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

	err = UpdateDataFotoPath("progress", "file", pathFotoFile, "id", dtMeeting.Id)
	if err != nil {
		return res, err
	}

	res.Status = http.StatusOK
	res.Message = "Berhasil memasukkan data"
	res.Data = dtMeeting

	defer db.DbClose(con)

	return res, nil
}

func CreateMeetingWithoutDocument(id, tanggal_meeting, waktu_meeting, tempat_meeting, waktu_mulai_meeting, waktu_selesai_meeting, notes string) (Response, error) {
	var res Response
	var dtMeeting Progress

	dtMeeting.Id, _ = strconv.Atoi(id)
	dtMeeting.Tanggal_meeting = tanggal_meeting
	dtMeeting.Waktu_meeting = waktu_meeting
	dtMeeting.Tempat_meeting = tempat_meeting
	dtMeeting.Waktu_mulai_meeting = waktu_mulai_meeting
	dtMeeting.Waktu_selesai_meeting = waktu_selesai_meeting
	dtMeeting.Notes = notes

	if id == "" && tanggal_meeting == "" && waktu_meeting == "" && tempat_meeting == "" && waktu_mulai_meeting == "" && waktu_selesai_meeting == "" && notes == "" {
		res.Status = 400
		res.Message = "Required fields are missing"
		res.Data = "data tidak lengkap"
		return res, errors.New("required fields are missing")
	}

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	query := `
		UPDATE progress SET tanggal_meeting = ?, waktu_meeting = ?, tempat_meeting = ?, 
		waktu_mulai_meeting=?, waktu_selesai_meeting = ?, notes = ?, data_lengkap = 'Y'
		WHERE id = ?`
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		dtMeeting.Tanggal_meeting,
		dtMeeting.Waktu_meeting,
		dtMeeting.Tempat_meeting,
		dtMeeting.Waktu_mulai_meeting,
		dtMeeting.Waktu_selesai_meeting,
		dtMeeting.Notes,
		dtMeeting.Id,
	)
	if err != nil {
		res.Status = 401
		res.Message = "exec gagal"
		res.Data = err.Error()
		return res, err
	}

	res.Status = http.StatusOK
	res.Message = "Berhasil memasukkan data"
	res.Data = dtMeeting

	defer db.DbClose(con)

	return res, nil
}

func GetAllProgress() (Response, error) {
	var res Response
	var arrProgress = []Progress{}

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	query := `
	SELECT id, user_id, perusahaan_id, id_asset, nama, proposal, status, data_lengkap, IFNULL(tanggal_meeting,""), IFNULL(waktu_meeting,""), 
	IFNULL(tempat_meeting,""), IFNULL(waktu_mulai_meeting,""), IFNULL(waktu_selesai_meeting,""), IFNULL(notes,""), IFNULL(file,""), IFNULL(tipe_file,"")
	FROM progress
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
		var dtProgress Progress
		err = result.Scan(
			&dtProgress.Id,
			&dtProgress.User_id,
			&dtProgress.Perusahaan_id,
			&dtProgress.Id_asset,
			&dtProgress.Nama_asset,
			&dtProgress.Proposal,
			&dtProgress.Status,
			&dtProgress.Data_lengkap,
			&dtProgress.Tanggal_meeting,
			&dtProgress.Waktu_meeting,
			&dtProgress.Tempat_meeting,
			&dtProgress.Waktu_mulai_meeting,
			&dtProgress.Waktu_selesai_meeting,
			&dtProgress.Notes,
			&dtProgress.Dokumen,
			&dtProgress.Tipe_dokumen,
		)
		if err != nil {
			res.Status = 401
			res.Message = "rows scan"
			res.Data = err.Error()
			return res, err
		}
		arrProgress = append(arrProgress, dtProgress)
	}

	res.Status = http.StatusOK
	res.Message = "Berhasil mengambil data"
	res.Data = arrProgress

	defer db.DbClose(con)
	return res, nil
}

func GetProgressByUserId(user_id string) (Response, error) {
	var res Response
	var arrProgress = []Progress{}

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	query := `
	SELECT id, user_id, perusahaan_id, id_asset, nama, proposal, status, data_lengkap, IFNULL(tanggal_meeting,""), IFNULL(waktu_meeting,""), 
	IFNULL(tempat_meeting,""), IFNULL(waktu_mulai_meeting,""), IFNULL(waktu_selesai_meeting,""), IFNULL(notes,""), IFNULL(file,""), IFNULL(tipe_file,"")
	FROM progress
	WHERE user_id = ?
	`
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	nId, _ := strconv.Atoi(user_id)
	result, err := stmt.Query(nId)
	if err != nil {
		res.Status = 401
		res.Message = "exec gagal"
		res.Data = err.Error()
		return res, err
	}
	defer result.Close()
	for result.Next() {
		var dtProgress Progress
		err = result.Scan(
			&dtProgress.Id,
			&dtProgress.User_id,
			&dtProgress.Perusahaan_id,
			&dtProgress.Id_asset,
			&dtProgress.Nama,
			&dtProgress.Proposal,
			&dtProgress.Status,
			&dtProgress.Data_lengkap,
			&dtProgress.Tanggal_meeting,
			&dtProgress.Waktu_meeting,
			&dtProgress.Tempat_meeting,
			&dtProgress.Waktu_mulai_meeting,
			&dtProgress.Waktu_selesai_meeting,
			&dtProgress.Notes,
			&dtProgress.Dokumen,
			&dtProgress.Tipe_dokumen,
		)
		if err != nil {
			res.Status = 401
			res.Message = "rows scan"
			res.Data = err.Error()
			return res, err
		}
		arrProgress = append(arrProgress, dtProgress)
	}

	if len(arrProgress) == 0 {
		res.Status = 404
		res.Message = "No progress data found"
		return res, fmt.Errorf("no progress data found for user_id: %s", user_id)
	}

	res.Status = http.StatusOK
	res.Message = "Berhasil mengambil data"
	res.Data = arrProgress

	defer db.DbClose(con)
	return res, nil
}

func GetProgressNotDoneByUserId(user_id string) (Response, error) {
	var res Response
	var arrProgress = []Progress{}

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	query := `
	SELECT id, user_id, perusahaan_id, id_asset, nama, proposal, status, data_lengkap, IFNULL(tanggal_meeting,""), IFNULL(waktu_meeting,""), 
	IFNULL(tempat_meeting,""), IFNULL(waktu_mulai_meeting,""), IFNULL(waktu_selesai_meeting,""), IFNULL(notes,""), IFNULL(file,""), IFNULL(tipe_file,"")
	FROM progress
	WHERE user_id = ? AND data_lengkap = 'N' AND status = 'A'
	`
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	nId, _ := strconv.Atoi(user_id)
	result, err := stmt.Query(nId)
	if err != nil {
		res.Status = 401
		res.Message = "exec gagal"
		res.Data = err.Error()
		return res, err
	}
	defer result.Close()
	for result.Next() {
		var dtProgress Progress
		err = result.Scan(
			&dtProgress.Id,
			&dtProgress.User_id,
			&dtProgress.Perusahaan_id,
			&dtProgress.Id_asset,
			&dtProgress.Nama,
			&dtProgress.Proposal,
			&dtProgress.Status,
			&dtProgress.Data_lengkap,
			&dtProgress.Tanggal_meeting,
			&dtProgress.Waktu_meeting,
			&dtProgress.Tempat_meeting,
			&dtProgress.Waktu_mulai_meeting,
			&dtProgress.Waktu_selesai_meeting,
			&dtProgress.Notes,
			&dtProgress.Dokumen,
			&dtProgress.Tipe_dokumen,
		)
		if err != nil {
			res.Status = 401
			res.Message = "rows scan"
			res.Data = err.Error()
			return res, err
		}
		arrProgress = append(arrProgress, dtProgress)
	}

	if len(arrProgress) == 0 {
		res.Status = 404
		res.Message = "No progress data found"
		return res, fmt.Errorf("no progress data found for user_id: %s", user_id)
	}

	res.Status = http.StatusOK
	res.Message = "Berhasil mengambil data"
	res.Data = arrProgress

	defer db.DbClose(con)
	return res, nil
}

func GetProgressByUserAsetId(user_id, aset_id string) (Response, error) {
	var res Response
	var arrProgress = []Progress{}
	fmt.Println(user_id, aset_id)

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	query := `
	SELECT id, user_id, perusahaan_id, id_asset, nama, proposal, status, data_lengkap, IFNULL(tanggal_meeting,""), IFNULL(waktu_meeting,""), 
	IFNULL(tempat_meeting,""), IFNULL(waktu_mulai_meeting,""), IFNULL(waktu_selesai_meeting,""), IFNULL(notes,""), IFNULL(file,""), IFNULL(tipe_file,"")
	FROM progress
	WHERE user_id = ? AND id_asset = ? AND status = 'A' AND data_lengkap = 'N'
	`
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	nId, _ := strconv.Atoi(user_id)
	nId2, _ := strconv.Atoi(aset_id)
	result, err := stmt.Query(nId, nId2)
	if err != nil {
		res.Status = 401
		res.Message = "exec gagal"
		res.Data = err.Error()
		return res, err
	}
	defer result.Close()
	for result.Next() {
		var dtProgress Progress
		err = result.Scan(
			&dtProgress.Id,
			&dtProgress.User_id,
			&dtProgress.Perusahaan_id,
			&dtProgress.Id_asset,
			&dtProgress.Nama,
			&dtProgress.Proposal,
			&dtProgress.Status,
			&dtProgress.Data_lengkap,
			&dtProgress.Tanggal_meeting,
			&dtProgress.Waktu_meeting,
			&dtProgress.Tempat_meeting,
			&dtProgress.Waktu_mulai_meeting,
			&dtProgress.Waktu_selesai_meeting,
			&dtProgress.Notes,
			&dtProgress.Dokumen,
			&dtProgress.Tipe_dokumen,
		)
		if err != nil {
			res.Status = 401
			res.Message = "rows scan"
			res.Data = err.Error()
			return res, err
		}
		arrProgress = append(arrProgress, dtProgress)
	}

	if len(arrProgress) == 0 {
		res.Status = 404
		res.Message = "No progress data found"
		return res, fmt.Errorf("no progress data found for user_id: %s, aset_id: %s", user_id, aset_id)
	}

	res.Status = http.StatusOK
	res.Message = "Berhasil mengambil data"
	res.Data = arrProgress

	defer db.DbClose(con)
	return res, nil
}

func GetProgressById(id string) (Response, error) {
	var res Response
	var arrProgress = []Progress{}

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	query := `
	SELECT id, user_id, perusahaan_id, id_asset, nama, proposal, status, data_lengkap, IFNULL(tanggal_meeting,""), IFNULL(waktu_meeting,""), 
	IFNULL(tempat_meeting,""), IFNULL(waktu_mulai_meeting,""), IFNULL(waktu_selesai_meeting,""), IFNULL(notes,""), IFNULL(file,""), IFNULL(tipe_file,"")
	FROM progress
	WHERE id = ?
	`
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	nId, _ := strconv.Atoi(id)
	result, err := stmt.Query(nId)
	if err != nil {
		res.Status = 401
		res.Message = "exec gagal"
		res.Data = err.Error()
		return res, err
	}
	defer result.Close()
	for result.Next() {
		var dtProgress Progress
		err = result.Scan(
			&dtProgress.Id,
			&dtProgress.User_id,
			&dtProgress.Perusahaan_id,
			&dtProgress.Id_asset,
			&dtProgress.Nama,
			&dtProgress.Proposal,
			&dtProgress.Status,
			&dtProgress.Data_lengkap,
			&dtProgress.Tanggal_meeting,
			&dtProgress.Waktu_meeting,
			&dtProgress.Tempat_meeting,
			&dtProgress.Waktu_mulai_meeting,
			&dtProgress.Waktu_selesai_meeting,
			&dtProgress.Notes,
			&dtProgress.Dokumen,
			&dtProgress.Tipe_dokumen,
		)
		if err != nil {
			res.Status = 401
			res.Message = "rows scan"
			res.Data = err.Error()
			return res, err
		}
		arrProgress = append(arrProgress, dtProgress)
	}

	res.Status = http.StatusOK
	res.Message = "Berhasil mengambil data"
	res.Data = arrProgress

	defer db.DbClose(con)
	return res, nil
}
