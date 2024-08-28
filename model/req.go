package model

import (
	"TemplateProject/db"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

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

	// ambil data aset dan update bagian lama
	var usage, luas, nilai, kondisi, batas_koordinat string
	var tags []string

	queryAsset := "SELECT `usage`, luas, nilai, kondisi, batas_koordinat FROM asset WHERE id_asset = ? LIMIT 1"
	err = con.QueryRow(queryAsset, dtSurveyReq.Id_asset).Scan(&usage, &luas, &nilai, &kondisi, &batas_koordinat)
	if err != nil {
		res.Status = 401
		res.Message = "query gagal"
		res.Data = err.Error()
		return res, err
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

	tagsString := strings.Join(tags, ", ")

	updateQuery := `
		UPDATE survey_request 
		SET usage_old = ?, luas_old = ?, nilai_old = ?, kondisi_old = ?, batas_koordinat_old = ?, tags_old = ? 
		WHERE id_transaksi_jual_sewa = ?
	`
	_, err = con.Exec(updateQuery, usage, luas, nilai, kondisi, batas_koordinat, tagsString, lastId)
	if err != nil {
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
			sr.status_request,sr.status_verifikasi,sr.dateline
		FROM survey_request sr
		JOIN asset a ON sr.id_asset = a.id_asset
		JOIN user u ON sr.user_id = u.user_id
		JOIN surveyor s ON sr.user_id = s.user_id
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
			&dtSurveyReq.Status_verifikasi, &dtSurveyReq.Dateline,
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

	query := `
	SELECT sr.*,a.nama 
	FROM survey_request sr
	JOIN asset a ON sr.id_asset = a.id_asset 
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
	nId, _ := strconv.Atoi(surveyreq_id)
	err = stmt.QueryRow(nId).Scan(&dtSurveyReq.Id_transaksi_jual_sewa, &dtSurveyReq.User_id, &dtSurveyReq.Id_asset, &dtSurveyReq.Created_at, &dtSurveyReq.Dateline, &dtSurveyReq.Status_request, &dtSurveyReq.Status_verifikasi, &dtSurveyReq.Data_lengkap, &dtSurveyReq.Usage_old, &dtSurveyReq.Usage_new, &dtSurveyReq.Luas_old, &dtSurveyReq.Luas_new, &dtSurveyReq.Nilai_old, &dtSurveyReq.Nilai_new, &dtSurveyReq.Kondisi_old, &dtSurveyReq.Kondisi_new, &dtSurveyReq.Batas_koordinat_old, &dtSurveyReq.Batas_koordinat_new, &dtSurveyReq.Tags_old, &dtSurveyReq.Tags_new, &dtSurveyReq.Nama_asset)
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

	query := `
	SELECT sr.*,a.nama 
	FROM survey_request sr
	JOIN asset a ON sr.id_asset = a.id_asset 
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
	nId, _ := strconv.Atoi(aset_id)
	err = stmt.QueryRow(nId).Scan(&dtSurveyReq.Id_transaksi_jual_sewa, &dtSurveyReq.User_id, &dtSurveyReq.Id_asset, &dtSurveyReq.Created_at, &dtSurveyReq.Dateline, &dtSurveyReq.Status_request, &dtSurveyReq.Status_verifikasi, &dtSurveyReq.Data_lengkap, &dtSurveyReq.Usage_old, &dtSurveyReq.Usage_new, &dtSurveyReq.Luas_old, &dtSurveyReq.Luas_new, &dtSurveyReq.Nilai_old, &dtSurveyReq.Nilai_new, &dtSurveyReq.Kondisi_old, &dtSurveyReq.Kondisi_new, &dtSurveyReq.Batas_koordinat_old, &dtSurveyReq.Batas_koordinat_new, &dtSurveyReq.Tags_old, &dtSurveyReq.Tags_new, &dtSurveyReq.Nama_asset)
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

// func GetAllSurveyReqByPerusahaanId(perusahaan_id string) (Response, error) {
// 	var res Response
// 	type SurveyorAssignment struct {
// 		Id_transaksi_jual_sewa int    `json:"id_transaksi_jual_sewa"`
// 		User_id                int    `json:"user_id"`
// 		Id_asset               int    `json:"id_asset"`
// 		Created_at             string `json:"created_at"`
// 		Dateline               string `json:"dateline"`
// 		Status_request         string `json:"status_request"`
// 		Status_verifikasi      string `json:"status_verifikasi"`
// 		Asset_nama             string `json:"asset_nama"`
// 		Asset_alamat           string `json:"asset_alamat"`
// 		Asset_titikkoordinat   string `json:"asset_titikkoordinat"`
// 	}
// 	type tempSurvAssignment struct {
// 		OngoingAssignment  []SurveyorAssignment `json:"ongoing_assignment"`
// 		FinishedAssignment []SurveyorAssignment `json:"finished_assignment"`
// 	}
// 	var arrSurveyReq = []SurveyorAssignment{}
// 	var dtSurveyReq SurveyorAssignment

// 	con, err := db.DbConnection()
// 	if err != nil {
// 		res.Status = 401
// 		res.Message = "gagal membuka database"
// 		res.Data = err.Error()
// 		return res, err
// 	}

// 	query := `
// 		SELECT sr.id_transaksi_jual_sewa, sr.user_id, sr.id_asset, sr.created_at, sr.status_request, sr.status_verifikasi, sr.dateline, a.nama, a.alamat, a.titik_koordinat
// 		FROM survey_request sr
// 		JOIN asset a ON sr.id_asset = a.id_asset
// 		WHERE sr.user_id = ? AND (sr.status_request = 'O' OR sr.status_request = 'R')
// 	`
// 	stmt, err := con.Prepare(query)
// 	if err != nil {
// 		res.Status = 401
// 		res.Message = "stmt gagal"
// 		res.Data = err.Error()
// 		return res, err
// 	}
// 	defer stmt.Close()

// 	nId, _ := strconv.Atoi(user_id)
// 	result, err := stmt.Query(nId)
// 	if err != nil {
// 		res.Status = 401
// 		res.Message = "exec gagal"
// 		res.Data = err.Error()
// 		return res, err
// 	}
// 	defer result.Close()
// 	for result.Next() {
// 		err = result.Scan(&dtSurveyReq.Id_transaksi_jual_sewa, &dtSurveyReq.User_id, &dtSurveyReq.Id_asset, &dtSurveyReq.Created_at,
// 			&dtSurveyReq.Status_request, &dtSurveyReq.Status_verifikasi, &dtSurveyReq.Dateline, &dtSurveyReq.Asset_nama, &dtSurveyReq.Asset_alamat, &dtSurveyReq.Asset_titikkoordinat)
// 		if err != nil {
// 			res.Status = 401
// 			res.Message = "rows scan"
// 			res.Data = err.Error()
// 			return res, err
// 		}
// 		arrSurveyReq = append(arrSurveyReq, dtSurveyReq)
// 	}
// 	var survey_assignment tempSurvAssignment
// 	survey_assignment.OngoingAssignment = arrSurveyReq
// 	fmt.Println("ambil finished survey request")
// 	// finished assignment
// 	var arrSurveyReqFinished = []SurveyorAssignment{}
// 	var dtSurveyReqFinished SurveyorAssignment
// 	queryfinished := `
// 		SELECT sr.id_transaksi_jual_sewa, sr.user_id, sr.id_asset, sr.created_at, sr.status_request, sr.status_verifikasi, sr.dateline, a.nama, a.alamat, a.titik_koordinat
// 		FROM survey_request sr
// 		JOIN asset a ON sr.id_asset = a.id_asset
// 		WHERE sr.user_id = ? AND sr.status_request = 'F'
// 	`
// 	stmtfinished, err := con.Prepare(queryfinished)
// 	if err != nil {
// 		res.Status = 401
// 		res.Message = "stmt gagal"
// 		res.Data = err.Error()
// 		return res, err
// 	}
// 	defer stmtfinished.Close()

// 	resultfinished, err := stmtfinished.Query(nId)
// 	if err != nil {
// 		res.Status = 401
// 		res.Message = "exec gagal"
// 		res.Data = err.Error()
// 		return res, err
// 	}
// 	defer resultfinished.Close()
// 	for resultfinished.Next() {
// 		err = resultfinished.Scan(&dtSurveyReq.Id_transaksi_jual_sewa, &dtSurveyReq.User_id, &dtSurveyReq.Id_asset, &dtSurveyReq.Created_at,
// 			&dtSurveyReq.Status_request, &dtSurveyReq.Status_verifikasi, &dtSurveyReq.Dateline, &dtSurveyReq.Asset_nama, &dtSurveyReq.Asset_alamat, &dtSurveyReq.Asset_titikkoordinat)
// 		if err != nil {
// 			res.Status = 401
// 			res.Message = "rows scan"
// 			res.Data = err.Error()
// 			return res, err
// 		}
// 		arrSurveyReqFinished = append(arrSurveyReqFinished, dtSurveyReqFinished)
// 	}

// 	survey_assignment.FinishedAssignment = arrSurveyReqFinished

// 	res.Status = http.StatusOK
// 	res.Message = "Berhasil mengambil data"
// 	res.Data = survey_assignment

// 	defer db.DbClose(con)
// 	return res, nil
// }

func GetAllSurveyReqByUserId(user_id string) (Response, error) {
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
	var survey_assignment tempSurvAssignment
	survey_assignment.OngoingAssignment = arrSurveyReq
	fmt.Println("ambil finished survey request")
	// finished assignment
	var arrSurveyReqFinished = []SurveyorAssignment{}
	var dtSurveyReqFinished SurveyorAssignment
	queryfinished := `
		SELECT sr.id_transaksi_jual_sewa, sr.user_id, sr.id_asset, sr.created_at, sr.status_request, sr.status_verifikasi, sr.dateline, a.nama, a.alamat, a.titik_koordinat 
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
			&dtSurveyReq.Status_request, &dtSurveyReq.Status_verifikasi, &dtSurveyReq.Dateline, &dtSurveyReq.Asset_nama, &dtSurveyReq.Asset_alamat, &dtSurveyReq.Asset_titikkoordinat)
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

	query := "SELECT * FROM transaction_request WHERE perusahaan_id = ?"
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
