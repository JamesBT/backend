package model

import (
	"TemplateProject/db"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

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
		SenderId int `json:"senderid"`
		UserID   int `json:"userid"`
		Kelas    int `json:"kelas"`
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

	title := "You're accepted!"
	detail := "Your request as user has been accepted"

	testingString := fmt.Sprintf(`
	{
		"user_id_sender": %d,
		"user_id_receiver": %d,
		"perusahaan_id_receiver": 0,
		"notification_title": "%s",
		"notification_detail": "%s"
	}
	`, requestacc.SenderId, requestacc.UserID, title, detail)

	res, err = CreateNotification(testingString)
	if err != nil {
		res.Status = 401
		res.Message = "kirim notifikasi gagal"
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
		SenderId     int    `json:"senderId"`
		UserId       int    `json:"userId"`
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

	selectUserQuery := `
	SELECT id_user
	FROM user_perusahaan
	WHERE id_perusahaan = ?`
	stmtSelect, err := con.Prepare(selectUserQuery)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmtSelect.Close()

	err = stmtSelect.QueryRow(requestacc.PerusahaanId).Scan(&requestacc.UserId)
	if err != nil {
		res.Status = 401
		res.Message = "exec gagal"
		res.Data = err.Error()
		return res, err
	}

	insertPrivPerusahaan := `
		INSERT INTO role_privilege_all (id_perusahaan,id_role,id_privilege) VALUES (?,?,?)
	`
	stmtInsertPriv, err := con.Prepare(insertPrivPerusahaan)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmtInsertPriv.Close()

	_, err = stmtInsertPriv.Exec(requestacc.PerusahaanId, 5, 25)
	if err != nil {
		res.Status = 401
		res.Message = "exec gagal"
		res.Data = err.Error()
		return res, err
	}

	_, err = stmtInsertPriv.Exec(requestacc.PerusahaanId, 5, 28)
	if err != nil {
		res.Status = 401
		res.Message = "exec gagal"
		res.Data = err.Error()
		return res, err
	}

	title := "Your company request has been accepted!"
	detail := ""

	testingString := fmt.Sprintf(`
	{
		"user_id_sender": %d,
		"user_id_receiver": %d,
		"perusahaan_id_receiver": %d,
		"notification_title": "%s",
		"notification_detail": "%s"
	}
	`, requestacc.SenderId, requestacc.UserId, requestacc.PerusahaanId, title, detail)

	res, err = CreateNotification(testingString)
	if err != nil {
		res.Status = 401
		res.Message = "kirim notifikasi gagal"
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
		UserID   int `json:"userid"`
		SenderId int `json:"senderId"`
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

	title := "You're declined!"
	detail := "Your request as user has been declined"

	testingString := fmt.Sprintf(`
	{
		"user_id_sender": %d,
		"user_id_receiver": %d,
		"perusahaan_id_receiver": 0,
		"notification_title": "%s",
		"notification_detail": "%s"
	}
	`, requestacc.SenderId, requestacc.UserID, title, detail)

	res, err = CreateNotification(testingString)
	if err != nil {
		res.Status = 401
		res.Message = "kirim notifikasi gagal"
		res.Data = err.Error()
		return res, err
	}

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
		SenderId     int    `json:"senderId"`
		UserId       int    `json:"userId"`
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
	result, err := stmt.Exec(requestacc.Alasan, requestacc.PerusahaanId)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}

	// kirim notif (masih mendatang)
	selectUserQuery := `
	SELECT id_user
	FROM user_perusahaan
	WHERE id_perusahaan = ?`
	stmtSelect, err := con.Prepare(selectUserQuery)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmtSelect.Close()

	err = stmtSelect.QueryRow(requestacc.PerusahaanId).Scan(&requestacc.UserId)
	if err != nil {
		res.Status = 401
		res.Message = "exec gagal"
		res.Data = err.Error()
		return res, err
	}

	title := "Your company request has been declined!"
	detail := ""

	testingString := fmt.Sprintf(`
	{
		"user_id_sender": %d,
		"user_id_receiver": %d,
		"perusahaan_id_receiver": 0,
		"notification_title": "%s",
		"notification_detail": "%s"
	}
	`, requestacc.SenderId, requestacc.UserId, title, detail)

	res, err = CreateNotification(testingString)
	if err != nil {
		res.Status = 401
		res.Message = "kirim notifikasi gagal"
		res.Data = err.Error()
		return res, err
	}

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
		SenderId     int `json:"senderId"`
		UserId       int `json:"userId"`
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

	queryNamaVerify := `SELECT nama_lengkap FROM user WHERE user_id = ?`
	stmtNamaVerify, err := con.Prepare(queryNamaVerify)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmtNamaVerify.Close()

	var namaVerify string
	err = stmtNamaVerify.QueryRow(requestacc.SenderId).Scan(&namaVerify)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}

	query := "UPDATE survey_request SET id_user_verify = ?, nama_user_verify = ?,status_request='F',status_verifikasi='V',data_lengkap='Y',status_submitted='Y' WHERE id_transaksi_jual_sewa = ?"
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	_, err = stmt.Exec(requestacc.SenderId, namaVerify, requestacc.SurveryReqId)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}

	selectquery := "SELECT id_asset,usage_new,luas_new,nilai_new,kondisi_new,titik_koordinat_new,batas_koordinat_new,tags_new,gambar_new,user_id FROM survey_request WHERE id_transaksi_jual_sewa = ?"
	selectstmt, err := con.Prepare(selectquery)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer selectstmt.Close()
	var usage_new, luas_new, nilai_new, kondisi_new, titik_koordinat_new, batas_koordinat_new, tags_new, gambar_new sql.NullString

	err = selectstmt.QueryRow(requestacc.SurveryReqId).Scan(
		&dtSurveyReq.Id_asset, &usage_new, &luas_new, &nilai_new, &kondisi_new, &titik_koordinat_new, &batas_koordinat_new, &tags_new, &gambar_new, &requestacc.UserId,
	)
	if err != nil {
		res.Status = 401
		res.Message = "rows gagal"
		res.Data = err.Error()
		return res, err
	}

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

	// update gambar (copy, update db)
	oldImagesQuery := "SELECT link_gambar FROM asset_gambar WHERE id_asset_gambar = ?"
	rows, err := con.Query(oldImagesQuery, dtSurveyReq.Id_asset)
	if err != nil {
		res.Status = 401
		res.Message = "query gagal"
		res.Data = err.Error()
		return res, err
	}
	defer rows.Close()

	var oldImagePaths []string
	for rows.Next() {
		var oldImage string
		if err := rows.Scan(&oldImage); err != nil {
			res.Status = 401
			res.Message = "scan gagal"
			res.Data = err.Error()
			return res, err
		}

		oldImagePaths = append(oldImagePaths, oldImage)

		// Ensure the destination directory exists before moving the file
		historyDir := fmt.Sprintf("uploads/asset/history/%d_%d", requestacc.SurveryReqId, dtSurveyReq.Id_asset)
		err = os.MkdirAll(historyDir, os.ModePerm)
		if err != nil {
			res.Status = 401
			res.Message = "mkdir gagal"
			res.Data = err.Error()
			return res, err
		}

		// Move old image to asset/history/(id_survey_request)_(id_asset)_(filename)
		historyPath := fmt.Sprintf("%s/%s", historyDir, filepath.Base(oldImage))
		err = os.Rename(oldImage, historyPath)
		if err != nil {
			res.Status = 401
			res.Message = "rename gagal"
			res.Data = err.Error()
			return res, err
		}

		// Update survey_request.gambar_old to new history path
		_, err = con.Exec("UPDATE survey_request SET gambar_old = ? WHERE id_transaksi_jual_sewa = ?", historyPath, requestacc.SurveryReqId)
		if err != nil {
			res.Status = 401
			res.Message = "query gagal"
			res.Data = err.Error()
			return res, err
		}
	}

	// Process and insert new images: Copy from survey_req/gambar to asset/gambar, rename, and update
	var updatedGambarNew []string
	if gambar_new.Valid {
		newImages := strings.Split(gambar_new.String, ",")
		for _, newImage := range newImages {
			newImage = strings.TrimSpace(newImage)
			if newImage != "" {
				filename := filepath.Base(newImage)
				parts := strings.SplitN(filename, "_", 2)
				if len(parts) == 2 {
					newFilename := parts[1]

					// Copy image from survey_req/gambar to asset/gambar and rename
					srcPath := fmt.Sprintf("uploads/survey_req/gambar/%s", filepath.Base(newImage))
					destPath := fmt.Sprintf("uploads/asset/foto/%s", newFilename)

					input, err := os.ReadFile(srcPath)
					if err != nil {
						res.Status = 401
						res.Message = "read file gagal"
						res.Data = err.Error()
						return res, err
					}
					err = os.WriteFile(destPath, input, 0644)
					if err != nil {
						res.Status = 401
						res.Message = "write file gagal"
						res.Data = err.Error()
						return res, err
					}

					// Insert the new image path into asset_gambar
					insertImageQuery := "INSERT INTO asset_gambar (id_asset_gambar, link_gambar) VALUES (?, ?)"
					_, err = con.Exec(insertImageQuery, dtSurveyReq.Id_asset, destPath)
					if err != nil {
						res.Status = 401
						res.Message = "insert gambar gagal"
						res.Data = err.Error()
						return res, err
					}

					updatedGambarNew = append(updatedGambarNew, destPath)
				}

			}
		}
	}

	title := "Your assignment has been accepted!"
	detail := ""

	testingString := fmt.Sprintf(`
	{
		"user_id_sender": %d,
		"user_id_receiver": %d,
		"perusahaan_id_receiver": 0,
		"notification_title": "%s",
		"notification_detail": "%s"
	}
	`, requestacc.SenderId, requestacc.UserId, title, detail)

	res, err = CreateNotification(testingString)
	if err != nil {
		res.Status = 401
		res.Message = "kirim notifikasi gagal"
		res.Data = err.Error()
		return res, err
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
		SurveyorId  int    `json:"user_id"`
		Dateline    string `json:"dateline"`
		SenderId    int    `json:"senderId"`
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

	// ambil user id dari surveyor
	surveyorquery := "SELECT user_id FROM surveyor WHERE `user_id` = ? "
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

	query := "UPDATE survey_request SET `user_id`=?,`dateline`=?,`status_request`='O',`status_submitted`='N',`status_verifikasi`='R' WHERE id_transaksi_jual_sewa = ?"
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

	assetquery := `
	SELECT a.nama 
	FROM survey_request sr
	LEFT JOIN asset a ON sr.id_asset = a.id_asset
	WHERE sr.id_transaksi_jual_sewa = ? `
	assetstmt, err := con.Prepare(assetquery)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer assetstmt.Close()

	var tempNamaAsset string
	err = assetstmt.QueryRow(requestacc.SurveyReqId).Scan(&tempNamaAsset)
	if err != nil {
		res.Status = 401
		res.Message = "gagal mengambil user_id"
		res.Data = err.Error()
		return res, err
	}

	title := "You’ve got a new assignment!"
	detail := fmt.Sprintf("Please do a survey on : %s", tempNamaAsset)

	testingString := fmt.Sprintf(`
	{
		"user_id_sender": %d,
		"user_id_receiver": %d,
		"perusahaan_id_receiver": 0,
		"notification_title": "%s",
		"notification_detail": "%s"
	}
	`, requestacc.SenderId, _tempuserid, title, detail)

	res, err = CreateNotification(testingString)
	if err != nil {
		res.Status = 401
		res.Message = "kirim notifikasi gagal"
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
		Id             int    `json:"id"`
		UserId         int    `json:"userId"`
		SenderId       int    `json:"senderId"`
		PerusahaanId   int    `json:"perusahaanId"`
		PerusahaanNama string `json:"perusahaanNama"`
		AssetId        int    `json:"assetId"`
		AssetNama      string `json:"assetNama"`
		NamaProgress   string `json:"namaProgress"`
		Proposal       string `json:"proposal"`
		Alasan         string `json:"alasan"`
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
	SELECT tr.id_asset,tr.user_id,tr.perusahaan_id,tr.nama_progress, tr.proposal, a.nama, p.name
	FROM transaction_request tr
	LEFT JOIN asset a ON tr.id_asset = a.id_asset
	LEFT JOIN perusahaan p ON tr.perusahaan_id = p.perusahaan_id
	WHERE tr.id_transaksi_jual_sewa = ?
	`
	err = con.QueryRow(getProgressAndAssetQuery, tempTranReq.Id).Scan(
		&tempTranReq.AssetId, &tempTranReq.UserId, &tempTranReq.PerusahaanId,
		&tempTranReq.NamaProgress, &tempTranReq.Proposal, &tempTranReq.AssetNama, &tempTranReq.PerusahaanNama)
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

	title := fmt.Sprintf("Your transaction request for %s has been ACCEPTED", tempTranReq.PerusahaanNama)
	detail := ""

	testingString := fmt.Sprintf(`
	{
		"user_id_sender": %d,
		"user_id_receiver": %d,
		"perusahaan_id_receiver": %d,
		"notification_title": "%s",
		"notification_detail": "%s"
	}
	`, tempTranReq.SenderId, tempTranReq.UserId, tempTranReq.PerusahaanId, title, detail)

	res, err = CreateNotification(testingString)
	if err != nil {
		res.Status = 401
		res.Message = "kirim notifikasi gagal"
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
		Id             int    `json:"id"`
		UserId         int    `json:"userId"`
		SenderId       int    `json:"senderId"`
		PerusahaanId   int    `json:"perusahaanId"`
		PerusahaanNama string `json:"perusahaanNama"`
		AssetId        int    `json:"assetId"`
		AssetNama      string `json:"assetNama"`
		NamaProgress   string `json:"namaProgress"`
		Proposal       string `json:"proposal"`
		Alasan         string `json:"alasan"`
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
	SELECT tr.id_asset,tr.user_id,tr.perusahaan_id,tr.nama_progress, tr.proposal, a.nama, tr.nama_progress, p.name 
	FROM transaction_request tr
	LEFT JOIN asset a ON tr.id_asset = a.id_asset
	LEFT JOIN perusahaan p ON tr.perusahaan_id = p.perusahaan_id
	WHERE tr.id_transaksi_jual_sewa = ?
	`
	err = con.QueryRow(getProgressAndAssetQuery, tempTranReq.Id).Scan(
		&tempTranReq.AssetId, &tempTranReq.UserId, &tempTranReq.PerusahaanId,
		&tempTranReq.NamaProgress, &tempTranReq.Proposal, &tempTranReq.AssetNama,
		&tempTranReq.NamaProgress, &tempTranReq.PerusahaanNama)
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

	title := fmt.Sprintf("Your transaction request for %s on company %s has been DECLINED", tempTranReq.NamaProgress, tempTranReq.PerusahaanNama)

	testingString := fmt.Sprintf(`
	{
		"user_id_sender": %d,
		"user_id_receiver": %d,
		"perusahaan_id_receiver": %d,
		"notification_title": "%s",
		"notification_detail": ""
	}
	`, tempTranReq.SenderId, tempTranReq.UserId, tempTranReq.PerusahaanId, title)

	res, err = CreateNotification(testingString)
	if err != nil {
		res.Status = 401
		res.Message = "kirim notifikasi gagal"
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
	GROUP BY p.perusahaan_id
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
