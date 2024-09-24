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
	"net/smtp"
	"os"
	"strconv"
	"strings"
)

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

func VerifyOTP(input string) (Response, error) {
	var res Response
	type temp_verif_user_acc struct {
		UserID   int `json:"userid"`
		Kode_OTP int `json:"kode_otp"`
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

	query := "SELECT id_user_detail FROM user_detail WHERE user_detail_id = ? AND kode_otp = ? "
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	var user_id int
	err = stmt.QueryRow(requestacc.UserID, requestacc.Kode_OTP).Scan(&user_id)
	if err != nil {
		if err == sql.ErrNoRows {
			res.Status = 401
			res.Message = "Kode OTP tidak valid"
			res.Data = nil
			return res, nil
		}
		res.Status = 401
		res.Message = "exec gagal"
		res.Data = err.Error()
		return res, err
	}

	updatequery := "UPDATE user_detail SET status_verifikasi_otp = 'V' where user_detail_id = ?"
	updatestmt, err := con.Prepare(updatequery)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer updatestmt.Close()

	_, err = updatestmt.Exec(requestacc.UserID)
	if err != nil {
		res.Status = 401
		res.Message = "gagal mengupdate data"
		res.Data = err.Error()
		return res, err
	}

	res.Status = http.StatusOK
	res.Message = "Berhasil verifikasi OTP"
	res.Data = user_id

	defer db.DbClose(con)
	return res, nil
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
		PerusahaanId int    `json:"perusahaan_id"`
		Kelas        int    `json:"kelas"`
		Field        string `json:"business_field"`
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

	query := "UPDATE perusahaan SET kelas=?,status='A' WHERE perusahaan_id = ?"
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

	// Insert into perusahaan_business table
	insertQuery := "INSERT INTO perusahaan_business (id_perusahaan, id_business) VALUES (?, ?)"
	fields := strings.Split(requestacc.Field, ",")
	for _, field := range fields {
		insertStmt, err := con.Prepare(insertQuery)
		if err != nil {
			res.Status = 401
			res.Message = "Stmt gagal"
			res.Data = err.Error()
			return res, err
		}
		defer insertStmt.Close()

		_, err = insertStmt.Exec(requestacc.PerusahaanId, field)
		if err != nil {
			res.Status = 401
			res.Message = "Gagal memasukkan ke perusahaan_business"
			res.Data = err.Error()
			return res, err
		}
	}

	// update user yang apply jadi admin
	queryuser := "UPDATE user_perusahaan SET `id_role`='5' WHERE id_perusahaan = ?"
	stmtuser, err := con.Prepare(queryuser)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmtuser.Close()

	_, err = stmtuser.Exec(requestacc.PerusahaanId)
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
		PerusahaanId int    `json:"perusahaan_id"`
		Alasan       string `json:"decline_message"`
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

	query := "UPDATE perusahaan SET status='D',denied_by_admin='Y',`alasan`=? WHERE perusahaan_id = ?"
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()
	fmt.Println(requestacc.Alasan)
	result, err := stmt.Exec(requestacc.Alasan, requestacc.PerusahaanId)
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

func GetVerifyPerusahaanDetailedById(id_perusahaan string) (Response, error) {
	var res Response
	type DetailVerifPerusahaan struct {
		Perusahaan_id       int             `json:"perusahaan_id"`
		Nama_perusahaan     string          `json:"perusahaan_nama"`
		Nama_user           string          `json:"user_nama"`
		Username_user       string          `json:"user_username"`
		Created_at          string          `json:"created_at"`
		Username_perusahaan string          `json:"perusahaan_username"`
		Lokasi              string          `json:"lokasi"`
		Tipe                string          `json:"tipe"`
		Status              string          `json:"status"`
		Kelas               int             `json:"kelas"`
		Dokumen_kepemilikan string          `json:"dokumen_kepemilikan"`
		Dokumen_perusahaan  string          `json:"dokumen_perusahaan"`
		Modal_awal          string          `json:"modal"`
		Deskripsi           string          `json:"deskripsi"`
		Alasan              string          `json:"alasan"`
		Field               []BusinessField `json:"field"`
	}
	var tempVerifPerusahaan DetailVerifPerusahaan

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	// nama perusahaan + nama user + username(user) + tanggal waktu + username (perusahaan)
	// lokasi + tipe + dokumen_kepemilikan + dokumen_perusahaan + modal awal + deskripsi
	query := `
		SELECT p.perusahaan_id, p.name, u.nama_lengkap, u.username, p.created_at, p.username,
		p.lokasi, p.tipe, p.status, IFNULL(p.kelas,0), p.dokumen_kepemilikan, p.dokumen_perusahaan, p.modal_awal, p.deskripsi, p.alasan
		FROM perusahaan p
		LEFT JOIN user_perusahaan up ON p.perusahaan_id = up.id_perusahaan
		LEFT JOIN user u ON up.id_user = u.user_id
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
	err = stmt.QueryRow(nId).Scan(
		&tempVerifPerusahaan.Perusahaan_id,
		&tempVerifPerusahaan.Nama_perusahaan,
		&tempVerifPerusahaan.Nama_user,
		&tempVerifPerusahaan.Username_user,
		&tempVerifPerusahaan.Created_at,
		&tempVerifPerusahaan.Username_perusahaan,
		&tempVerifPerusahaan.Lokasi,
		&tempVerifPerusahaan.Tipe,
		&tempVerifPerusahaan.Status,
		&tempVerifPerusahaan.Kelas,
		&tempVerifPerusahaan.Dokumen_kepemilikan,
		&tempVerifPerusahaan.Dokumen_perusahaan,
		&tempVerifPerusahaan.Modal_awal,
		&tempVerifPerusahaan.Deskripsi,
		&tempVerifPerusahaan.Alasan,
	)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}

	queryPerusahaan := `
	SELECT bf.* 
	FROM business_field bf
	LEFT JOIN perusahaan_business pb ON bf.id = pb.id_business
	WHERE pb.id_perusahaan = ?
	`
	stmtPerusahaan, err := con.Prepare(queryPerusahaan)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmtPerusahaan.Close()

	resultPerusahaan, err := stmtPerusahaan.Query(nId)
	if err != nil {
		res.Status = 401
		res.Message = "exec gagal"
		res.Data = err.Error()
		return res, err
	}
	defer resultPerusahaan.Close()

	var arrBusiness = []BusinessField{}
	for resultPerusahaan.Next() {
		var dtBusiness BusinessField
		err = resultPerusahaan.Scan(&dtBusiness.Id, &dtBusiness.Nama, &dtBusiness.Detail)
		if err != nil {
			res.Status = 401
			res.Message = "rows scan"
			res.Data = err.Error()
			return res, err
		}
		arrBusiness = append(arrBusiness, dtBusiness)
	}
	tempVerifPerusahaan.Field = arrBusiness

	res.Status = http.StatusOK
	res.Message = "Berhasil mengambil data"
	res.Data = tempVerifPerusahaan

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

	query := "UPDATE survey_request SET status_request='F',status_verifikasi='V' WHERE id_transaksi_jual_sewa = ?"
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

	selectquery := "SELECT id_asset,usage_new,luas_new,nilai_new,kondisi_new,titik_koordinat_new,batas_koordinat_new,tags_new FROM survey_request WHERE id_transaksi_jual_sewa = ?"
	selectstmt, err := con.Prepare(selectquery)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer selectstmt.Close()
	var usage_new, luas_new, nilai_new, kondisi_new, titik_koordinat_new, batas_koordinat_new, tags_new sql.NullString

	err = selectstmt.QueryRow(requestacc.SurveryReqId).Scan(
		&dtSurveyReq.Id_asset, &usage_new, &luas_new, &nilai_new, &kondisi_new, &titik_koordinat_new, &batas_koordinat_new, &tags_new,
	)
	if err != nil {
		res.Status = 401
		res.Message = "rows gagal"
		res.Data = err.Error()
		return res, err
	}
	fmt.Println("usage_new", usage_new)
	fmt.Println("luas_new", luas_new)
	fmt.Println("nilai_new", nilai_new)
	fmt.Println("kondisi", kondisi_new)
	fmt.Println("titik_koordinat_new", titik_koordinat_new)
	fmt.Println("batas_koordinat_new", batas_koordinat_new)
	fmt.Println("tags_new", tags_new)
	fmt.Println("id", dtSurveyReq.Id_asset)

	// update data asset dengan yang baru
	updatequery := "UPDATE asset SET `kondisi`= ?,`titik_koordinat`= ?,`batas_koordinat`= ?,`luas`= ?,`nilai`= ? WHERE `id_asset`= ?"
	updatestmt, err := con.Prepare(updatequery)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer updatestmt.Close()

	_, err = updatestmt.Exec(kondisi_new, titik_koordinat_new, batas_koordinat_new, luas_new, nilai_new, dtSurveyReq.Id_asset)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}

	// delete baru insert tag asset
	deleteTagsQuery := "DELETE FROM asset_tags WHERE id_asset = ?"
	deleteStmt, err := con.Prepare(deleteTagsQuery)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer deleteStmt.Close()

	_, err = deleteStmt.Exec(dtSurveyReq.Id_asset)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}

	deleteTagsQuery2 := "DELETE FROM asset_penggunaan WHERE id_asset = ?"
	deleteStmt2, err := con.Prepare(deleteTagsQuery2)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer deleteStmt2.Close()

	_, err = deleteStmt2.Exec(dtSurveyReq.Id_asset)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}

	// Insert new usage values
	if usage_new.Valid {
		usageList := strings.Split(usage_new.String, ",")
		for _, usage := range usageList {
			usage = strings.TrimSpace(usage)
			if usage == "" {
				continue
			}
			insertUsageQuery := "INSERT INTO asset_penggunaan (id_asset, id_penggunaan) VALUES (?, ?)"
			insertStmt, err := con.Prepare(insertUsageQuery)
			if err != nil {
				res.Status = 401
				res.Message = "stmt gagal"
				res.Data = err.Error()
				return res, err
			}
			defer insertStmt.Close()

			_, err = insertStmt.Exec(dtSurveyReq.Id_asset, usage)
			if err != nil {
				res.Status = 401
				res.Message = "stmt gagal"
				res.Data = err.Error()
				return res, err
			}
		}
	}

	// Insert new tags
	if tags_new.Valid {
		tagList := strings.Split(tags_new.String, ",")
		for _, tag := range tagList {
			tag = strings.TrimSpace(tag)
			if tag == "" {
				continue
			}
			insertTagQuery := "INSERT INTO asset_tags (id_asset, id_tags) VALUES (?, ?)"
			insertTagStmt, err := con.Prepare(insertTagQuery)
			if err != nil {
				res.Status = 401
				res.Message = "stmt gagal"
				res.Data = err.Error()
				return res, err
			}
			defer insertTagStmt.Close()

			_, err = insertTagStmt.Exec(dtSurveyReq.Id_asset, tag)
			if err != nil {
				res.Status = 401
				res.Message = "stmt gagal"
				res.Data = err.Error()
				return res, err
			}
		}
	}

	tempaset, _ := GetAssetById(strconv.Itoa(dtSurveyReq.Id_asset))

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

func AcceptTransaction(input string) (Response, error) {
	var res Response
	type TranReqStatus struct {
		Id           int    `json:"id"`
		UserId       int    `json:"userId"`
		PerusahaanId int    `json:"perusahaanId"`
		AssetId      int    `json:"assetId"`
		AssetNama    string `json:"assetNama"`
		NamaProgress string `json:"namaProgress"`
		Proposal     string `json:"proposal"`
		Alasan       string `json:"alasan"`
	}

	var tempTranReq TranReqStatus
	err := json.Unmarshal([]byte(input), &tempTranReq)
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

	// ambil user id dari surveyor
	query := `
	UPDATE transaction_request SET status = 'A', alasan = ?
	WHERE id_transaksi_jual_sewa = ?
	`
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	_, err = stmt.Exec(tempTranReq.Alasan, tempTranReq.Id)
	if err != nil {
		res.Status = 401
		res.Message = "Gagal mengupdate status"
		res.Data = err.Error()
		return res, err
	}

	getProgressAndAssetQuery := `
	SELECT tr.id_asset,tr.user_id,tr.perusahaan_id,tr.nama_progress, tr.proposal, a.nama 
	FROM transaction_request tr
	LEFT JOIN asset a ON tr.id_asset = a.id_asset
	WHERE tr.id_transaksi_jual_sewa = ?
	`
	err = con.QueryRow(getProgressAndAssetQuery, tempTranReq.Id).Scan(&tempTranReq.AssetId, &tempTranReq.UserId, &tempTranReq.PerusahaanId, &tempTranReq.NamaProgress, &tempTranReq.Proposal, &tempTranReq.AssetNama)
	if err != nil {
		res.Status = 401
		res.Message = "Gagal mengambil nama_progress, proposal, dan nama asset"
		res.Data = err.Error()
		return res, err
	}

	insertProgressQuery := `
	INSERT INTO progress (user_id, perusahaan_id, id_asset, nama, proposal) 
	VALUES (?, ?, ?, ?, ?)
	`
	stmtProgress, err := con.Prepare(insertProgressQuery)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmtProgress.Close()

	_, err = stmtProgress.Exec(tempTranReq.UserId, tempTranReq.PerusahaanId, tempTranReq.AssetId, tempTranReq.NamaProgress, tempTranReq.Proposal)
	if err != nil {
		res.Status = 401
		res.Message = "Gagal memasukkan data ke progress"
		res.Data = err.Error()
		return res, err
	}

	res.Status = http.StatusOK
	res.Message = "Berhasil mengupdate data"
	res.Data = nil

	defer db.DbClose(con)
	return res, nil
}

func DeclineTransaction(input string) (Response, error) {
	var res Response
	type TranReqStatus struct {
		Id           int    `json:"id"`
		UserId       int    `json:"userId"`
		PerusahaanId int    `json:"perusahaanId"`
		AssetId      int    `json:"assetId"`
		AssetNama    string `json:"assetNama"`
		NamaProgress string `json:"namaProgress"`
		Proposal     string `json:"proposal"`
		Alasan       string `json:"alasan"`
	}

	var tempTranReq TranReqStatus
	err := json.Unmarshal([]byte(input), &tempTranReq)
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

	// ambil user id dari surveyor
	query := `
	UPDATE transaction_request SET status = 'D', alasan = ?
	WHERE id_transaksi_jual_sewa = ?
	`
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	_, err = stmt.Exec(tempTranReq.Alasan, tempTranReq.Id)
	if err != nil {
		res.Status = 401
		res.Message = "Gagal mengupdate status"
		res.Data = err.Error()
		return res, err
	}

	getProgressAndAssetQuery := `
	SELECT tr.id_asset,tr.user_id,tr.perusahaan_id,tr.nama_progress, tr.proposal, a.nama 
	FROM transaction_request tr
	LEFT JOIN asset a ON tr.id_asset = a.id_asset
	WHERE tr.id_transaksi_jual_sewa = ?
	`
	err = con.QueryRow(getProgressAndAssetQuery, tempTranReq.Id).Scan(&tempTranReq.AssetId, &tempTranReq.UserId, &tempTranReq.PerusahaanId, &tempTranReq.NamaProgress, &tempTranReq.Proposal, &tempTranReq.AssetNama)
	if err != nil {
		res.Status = 401
		res.Message = "Gagal mengambil nama_progress, proposal, dan nama asset"
		res.Data = err.Error()
		return res, err
	}

	insertProgressQuery := `
	INSERT INTO progress (user_id, perusahaan_id, id_asset, nama, proposal, status) 
	VALUES (?, ?, ?, ?, ?, 'D')
	`
	stmtProgress, err := con.Prepare(insertProgressQuery)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmtProgress.Close()

	_, err = stmtProgress.Exec(tempTranReq.UserId, tempTranReq.PerusahaanId, tempTranReq.AssetId, tempTranReq.NamaProgress, tempTranReq.Proposal)
	if err != nil {
		res.Status = 401
		res.Message = "Gagal memasukkan data ke progress"
		res.Data = err.Error()
		return res, err
	}

	res.Status = http.StatusOK
	res.Message = "Berhasil mengupdate data"
	res.Data = nil

	defer db.DbClose(con)
	return res, nil
}

func GetAllVerify() (Response, error) {
	var res Response

	type TempVerifyPerusahaan struct {
		Id                  int    `json:"id_perusahaan"`
		Status              string `json:"status"`
		Nama                string `json:"nama"`
		Username            string `json:"username"`
		NamaUser            string `json:"namauser"`
		NamaLengkapUser     string `json:"namalengkapuser"`
		Lokasi              string `json:"lokasi"`
		Kelas               int    `json:"kelas"`
		Tipe                string `json:"tipe"`
		Dokumen_kepemilikan string `json:"dokumen_kepemilikan"`
		Dokumen_perusahaan  string `json:"dokumen_perusahaan"`
		Modal               string `json:"modal"`
		Deskripsi           string `json:"deskripsi"`
		CreatedAt           string `json:"created_at"`
		UserJoined          []User
	}
	type AllVerify struct {
		Users      []User                 `json:"users"`
		Perusahaan []TempVerifyPerusahaan `json:"perusahaan"`
	}
	var allVerify AllVerify
	var arrUsers = []User{}
	var arrPerusahaan = []TempVerifyPerusahaan{}

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}
	defer db.DbClose(con)

	// Query for perusahaan data
	queryPerusahaan := `
	SELECT p.perusahaan_id, p.status, p.name, IFNULL(p.username,""), IFNULL(u.username,""), IFNULL(u.nama_lengkap,""), p.lokasi, p.tipe, IFNULL(p.kelas,0), p.dokumen_kepemilikan, 
	p.dokumen_perusahaan, p.modal_awal, p.deskripsi, p.created_at 
	FROM perusahaan p
	LEFT JOIN user_perusahaan up ON p.perusahaan_id = up.id_perusahaan
	LEFT JOIN user u ON up.id_user = u.user_id
	ORDER BY p.created_at DESC`
	stmtPerusahaan, err := con.Prepare(queryPerusahaan)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmtPerusahaan.Close()

	resultPerusahaan, err := stmtPerusahaan.Query()
	if err != nil {
		res.Status = 401
		res.Message = "exec gagal"
		res.Data = err.Error()
		return res, err
	}
	defer resultPerusahaan.Close()

	for resultPerusahaan.Next() {
		var dtPerusahaan TempVerifyPerusahaan
		err = resultPerusahaan.Scan(&dtPerusahaan.Id, &dtPerusahaan.Status, &dtPerusahaan.Nama, &dtPerusahaan.Username, &dtPerusahaan.NamaUser, &dtPerusahaan.NamaLengkapUser, &dtPerusahaan.Lokasi, &dtPerusahaan.Tipe, &dtPerusahaan.Kelas, &dtPerusahaan.Dokumen_kepemilikan, &dtPerusahaan.Dokumen_perusahaan, &dtPerusahaan.Modal, &dtPerusahaan.Deskripsi, &dtPerusahaan.CreatedAt)
		if err != nil {
			res.Status = 401
			res.Message = "rows scan"
			res.Data = err.Error()
			return res, err
		}
		arrPerusahaan = append(arrPerusahaan, dtPerusahaan)
	}
	allVerify.Perusahaan = arrPerusahaan

	// Query for user data
	queryUser := `
	SELECT u.user_id, u.username, u.password, u.nama_lengkap, u.alamat, u.jenis_kelamin, 
	IFNULL(u.tanggal_lahir,""), u.email, u.nomor_telepon, u.foto_profil, u.ktp, ud.status
	FROM user u
	LEFT JOIN user_detail ud ON u.user_id = ud.user_detail_id
	ORDER BY u.created_at DESC`
	stmtUser, err := con.Prepare(queryUser)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmtUser.Close()

	resultUser, err := stmtUser.Query()
	if err != nil {
		res.Status = 401
		res.Message = "exec gagal"
		res.Data = err.Error()
		return res, err
	}
	defer resultUser.Close()

	for resultUser.Next() {
		var dtUser User
		err = resultUser.Scan(&dtUser.Id, &dtUser.Username, &dtUser.Password, &dtUser.Nama_lengkap,
			&dtUser.Alamat, &dtUser.Jenis_kelamin, &dtUser.Tgl_lahir, &dtUser.Email, &dtUser.No_telp, &dtUser.Foto_profil,
			&dtUser.Ktp, &dtUser.Status,
		)
		if err != nil {
			res.Status = 401
			res.Message = "rows scan"
			res.Data = err.Error()
			return res, err
		}
		arrUsers = append(arrUsers, dtUser)
	}
	allVerify.Users = arrUsers

	// Return the combined results
	res.Status = http.StatusOK
	res.Message = "Berhasil mengambil data"
	res.Data = allVerify

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

const CONFIG_SMTP_HOST = "smtp.gmail.com"
const CONFIG_SMTP_PORT = 587
const CONFIG_SENDER_NAME = "LEAP - Testing Kirim Email"
const CONFIG_AUTH_EMAIL = "c14210026@john.petra.ac.id"
const CONFIG_AUTH_PASSWORD = "alzx sjan ikkr ipsm"

func sendMail(to []string, cc []string, subject, message string) error {
	body := "From: " + CONFIG_SENDER_NAME + "\n" +
		"To: " + strings.Join(to, ",") + "\n" +
		"Cc: " + strings.Join(cc, ",") + "\n" +
		"Subject: " + subject + "\n\n" +
		message

	auth := smtp.PlainAuth("", CONFIG_AUTH_EMAIL, CONFIG_AUTH_PASSWORD, CONFIG_SMTP_HOST)
	smtpAddr := fmt.Sprintf("%s:%d", CONFIG_SMTP_HOST, CONFIG_SMTP_PORT)

	err := smtp.SendMail(smtpAddr, auth, CONFIG_AUTH_EMAIL, append(to, cc...), []byte(body))
	if err != nil {
		return err
	}

	return nil
}

func CreateNotification(input string) (Response, error) {
	var res Response
	var kirimnotif Notification
	err := json.Unmarshal([]byte(input), &kirimnotif)
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

	DeleteNotification()
	var query string
	var result sql.Result
	if kirimnotif.User_id_receiver != 0 {
		query = "INSERT INTO notification (user_id_sender, user_id_receiver, created_at, notification_title, notification_detail) VALUES (?,?,NOW(),?,?)"
		stmt, err := con.Prepare(query)
		if err != nil {
			res.Status = 401
			res.Message = "stmt gagal"
			res.Data = err.Error()
			return res, err
		}
		defer stmt.Close()

		result, err = stmt.Exec(kirimnotif.User_id_sender, kirimnotif.User_id_receiver, kirimnotif.Title, kirimnotif.Detail)
		if err != nil {
			res.Status = 401
			res.Message = "exec gagal"
			res.Data = err.Error()
			return res, err
		}
	} else if kirimnotif.Perusahaan_id_receiver != 0 {
		query = "INSERT INTO notification (user_id_sender, perusahaan_id_receiver, created_at, notification_title, notification_detail) VALUES (?,?,NOW(),?,?)"
		stmt, err := con.Prepare(query)
		if err != nil {
			res.Status = 401
			res.Message = "stmt gagal"
			res.Data = err.Error()
			return res, err
		}
		defer stmt.Close()

		result, err = stmt.Exec(kirimnotif.User_id_sender, kirimnotif.Perusahaan_id_receiver, kirimnotif.Title, kirimnotif.Detail)
		if err != nil {
			res.Status = 401
			res.Message = "exec gagal"
			res.Data = err.Error()
			return res, err
		}
	} else {
		res.Status = 401
		res.Message = "kedua parameter kosong"
		return res, errors.New(res.Message)
	}

	notifId, err := result.LastInsertId()
	if err != nil {
		res.Status = 500
		res.Message = "gagal mendapatkan user ID"
		res.Data = err.Error()
		return res, err
	}
	kirimnotif.Notification_id = int(notifId)

	var tempNotif Response
	tempNotif, _ = GetNotificationById(strconv.Itoa(kirimnotif.Notification_id))
	res.Status = http.StatusOK
	res.Message = "Berhasil kirim notifikasi"
	res.Data = tempNotif.Data

	defer db.DbClose(con)
	return res, nil
}

func GetNotificationById(id string) (Response, error) {
	var res Response
	var dtNotif Notification

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	DeleteNotification()

	query := `
	SELECT notification_id, user_id_sender, IFNULL(user_id_receiver,0), IFNULL(perusahaan_id_receiver,0), created_at, notification_title, notification_detail
	FROM notification
	WHERE notification_id = ?
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
	err = stmt.QueryRow(nId).Scan(
		&dtNotif.Notification_id, &dtNotif.User_id_sender, &dtNotif.User_id_receiver, &dtNotif.Perusahaan_id_receiver,
		&dtNotif.Created_at, &dtNotif.Title, &dtNotif.Detail,
	)
	if err != nil {
		res.Status = 401
		res.Message = "exec gagal"
		res.Data = err.Error()
		return res, err
	}

	res.Status = http.StatusOK
	res.Message = "Berhasil mengambil data"
	res.Data = dtNotif

	defer db.DbClose(con)
	return res, nil
}

func GetNotificationByUserIdReceiver(user_id string) (Response, error) {
	var res Response
	var arrNotification []Notification

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	DeleteNotification()

	query := `
	SELECT notification_id, user_id_sender, IFNULL(user_id_receiver,0), IFNULL(perusahaan_id_receiver,0), created_at, notification_title, notification_detail 
	FROM notification
	WHERE user_id_receiver = ?
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
	rows, err := stmt.Query(nId)
	if err != nil {
		res.Status = 401
		res.Message = "query gagal"
		res.Data = err.Error()
		return res, err
	}
	defer rows.Close()

	for rows.Next() {
		var dtNotif Notification
		err := rows.Scan(
			&dtNotif.Notification_id, &dtNotif.User_id_sender,
			&dtNotif.User_id_receiver, &dtNotif.Perusahaan_id_receiver, &dtNotif.Created_at,
			&dtNotif.Title, &dtNotif.Detail,
		)
		if err != nil {
			res.Status = 401
			res.Message = "exec gagal"
			res.Data = err.Error()
			return res, err
		}

		arrNotification = append(arrNotification, dtNotif)
	}

	res.Status = http.StatusOK
	res.Message = "Berhasil mengambil data"
	res.Data = arrNotification

	defer db.DbClose(con)
	return res, nil
}

func GetNotificationByPerusahaanIdReceiver(perusahaan_id string) (Response, error) {
	var res Response
	var arrNotification []Notification

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	DeleteNotification()

	query := `
	SELECT notification_id, user_id_sender, IFNULL(user_id_receiver,0), IFNULL(perusahaan_id_receiver,0), created_at, notification_title, notification_detail 
	FROM notification
	WHERE perusahaan_id_receiver = ?
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
	rows, err := stmt.Query(nId)
	if err != nil {
		res.Status = 401
		res.Message = "query gagal"
		res.Data = err.Error()
		return res, err
	}
	defer rows.Close()

	for rows.Next() {
		var dtNotif Notification
		err := rows.Scan(
			&dtNotif.Notification_id, &dtNotif.User_id_sender,
			&dtNotif.User_id_receiver, &dtNotif.Perusahaan_id_receiver, &dtNotif.Created_at,
			&dtNotif.Title, &dtNotif.Detail,
		)
		if err != nil {
			res.Status = 401
			res.Message = "exec gagal"
			res.Data = err.Error()
			return res, err
		}

		arrNotification = append(arrNotification, dtNotif)
	}

	if len(arrNotification) == 0 {
		res.Status = 401
		res.Message = "data kosong"
		return res, errors.New(res.Message)
	}

	res.Status = http.StatusOK
	res.Message = "Berhasil mengambil data"
	res.Data = arrNotification

	defer db.DbClose(con)
	return res, nil
}

func DeleteNotification() (Response, error) {
	var res Response
	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	query := "DELETE FROM notification WHERE created_at < NOW() - INTERVAL 6 MONTH"
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	_, err = stmt.Exec()
	if err != nil {
		res.Status = 401
		res.Message = "exec gagal"
		res.Data = err.Error()
		return res, err
	}

	res.Status = http.StatusOK
	res.Message = "Berhasil menghapus notifikasi lama"

	return res, nil
}

func GetAllUsage() (Response, error) {
	var res Response
	var arrUsage = []Kegunaan{}

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	query := `
	SELECT *
	FROM penggunaan
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
	defer rows.Close()
	for rows.Next() {
		var dtUsage Kegunaan
		err := rows.Scan(&dtUsage.Id, &dtUsage.Nama)
		if err != nil {
			res.Status = 401
			res.Message = "scan gagal"
			res.Data = err.Error()
			return res, err
		}

		arrUsage = append(arrUsage, dtUsage)
	}

	if err = rows.Err(); err != nil {
		res.Status = 401
		res.Message = "rows error"
		res.Data = err.Error()
		return res, err
	}

	if len(arrUsage) == 0 {
		res.Status = 404
		res.Message = "Data tidak ditemukan"
		res.Data = nil
		return res, nil
	}

	res.Status = http.StatusOK
	res.Message = "Berhasil mengambil data"
	res.Data = arrUsage

	defer db.DbClose(con)
	return res, nil
}

func GetAllTags() (Response, error) {
	var res Response
	var arrTags = []Tags{}

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	query := `
	SELECT *
	FROM tags
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
	defer rows.Close()
	for rows.Next() {
		var dtTags Tags
		err := rows.Scan(&dtTags.Id, &dtTags.Nama, &dtTags.Detail)
		if err != nil {
			res.Status = 401
			res.Message = "scan gagal"
			res.Data = err.Error()
			return res, err
		}

		arrTags = append(arrTags, dtTags)
	}

	if err = rows.Err(); err != nil {
		res.Status = 401
		res.Message = "rows error"
		res.Data = err.Error()
		return res, err
	}

	if len(arrTags) == 0 {
		res.Status = 404
		res.Message = "Data tidak ditemukan"
		res.Data = nil
		return res, nil
	}

	res.Status = http.StatusOK
	res.Message = "Berhasil mengambil data"
	res.Data = arrTags

	defer db.DbClose(con)
	return res, nil
}

func GetTagsUsed() (Response, error) {
	var res Response
	var arrTags = []Tags{}

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	query := `
	SELECT at.id_tags, t.nama, t.detail
	FROM asset_tags at
	LEFT JOIN tags t ON at.id_tags = t.id
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
	defer rows.Close()

	addedTags := make(map[int]bool)

	for rows.Next() {
		var dtTags Tags
		err := rows.Scan(&dtTags.Id, &dtTags.Nama, &dtTags.Detail)
		if err != nil {
			res.Status = 401
			res.Message = "scan gagal"
			res.Data = err.Error()
			return res, err
		}

		if _, exists := addedTags[dtTags.Id]; !exists {
			arrTags = append(arrTags, dtTags)
			addedTags[dtTags.Id] = true
		}
	}

	if err = rows.Err(); err != nil {
		res.Status = 401
		res.Message = "rows error"
		res.Data = err.Error()
		return res, err
	}

	if len(arrTags) == 0 {
		res.Status = 404
		res.Message = "Data tidak ditemukan"
		res.Data = nil
		return res, nil
	}

	res.Status = http.StatusOK
	res.Message = "Berhasil mengambil data"
	res.Data = arrTags

	defer db.DbClose(con)
	return res, nil
}

func GetAllProvinsi() (Response, error) {
	var res Response
	var arrProvinsi = []Provinsi{}

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	query := `
	SELECT *
	FROM provinsi
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
	defer rows.Close()

	addedProvinsi := make(map[int]bool)

	for rows.Next() {
		var dtProvinsi Provinsi
		err := rows.Scan(&dtProvinsi.Id, &dtProvinsi.Nama)
		if err != nil {
			res.Status = 401
			res.Message = "scan gagal"
			res.Data = err.Error()
			return res, err
		}

		if _, exists := addedProvinsi[dtProvinsi.Id]; !exists {
			// If not added, append it to the array and mark it as added
			arrProvinsi = append(arrProvinsi, dtProvinsi)
			addedProvinsi[dtProvinsi.Id] = true
		}
	}

	if err = rows.Err(); err != nil {
		res.Status = 401
		res.Message = "rows error"
		res.Data = err.Error()
		return res, err
	}

	if len(arrProvinsi) == 0 {
		res.Status = 404
		res.Message = "Data tidak ditemukan"
		res.Data = nil
		return res, nil
	}

	res.Status = http.StatusOK
	res.Message = "Berhasil mengambil data"
	res.Data = arrProvinsi

	defer db.DbClose(con)
	return res, nil
}

func GetProvinsiUsed() (Response, error) {
	var res Response
	var arrProvinsi = []Provinsi{}

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	query := `
	SELECT DISTINCT a.provinsi, p.nama
	FROM asset a
	LEFT JOIN provinsi p ON a.provinsi = p.id_provinsi
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
	defer rows.Close()
	for rows.Next() {
		var dtProvinsi Provinsi
		err := rows.Scan(&dtProvinsi.Id, &dtProvinsi.Nama)
		if err != nil {
			res.Status = 401
			res.Message = "scan gagal"
			res.Data = err.Error()
			return res, err
		}

		arrProvinsi = append(arrProvinsi, dtProvinsi)
	}

	if err = rows.Err(); err != nil {
		res.Status = 401
		res.Message = "rows error"
		res.Data = err.Error()
		return res, err
	}

	if len(arrProvinsi) == 0 {
		res.Status = 404
		res.Message = "Data tidak ditemukan"
		res.Data = nil
		return res, nil
	}

	res.Status = http.StatusOK
	res.Message = "Berhasil mengambil data"
	res.Data = arrProvinsi

	defer db.DbClose(con)
	return res, nil
}
