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

// NOTIFICATION
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

	if kirimnotif.User_id_receiver != 0 && kirimnotif.Perusahaan_id_receiver != 0 {
		query = "INSERT INTO notification (user_id_sender, user_id_receiver, perusahaan_id_receiver, created_at, notification_title, notification_detail) VALUES (?,?,?,NOW(),?,?)"
		stmt, err := con.Prepare(query)
		if err != nil {
			res.Status = 401
			res.Message = "stmt gagal"
			res.Data = err.Error()
			return res, err
		}
		defer stmt.Close()

		result, err = stmt.Exec(kirimnotif.User_id_sender, kirimnotif.User_id_receiver, kirimnotif.Perusahaan_id_receiver, kirimnotif.Title, kirimnotif.Detail)
		if err != nil {
			res.Status = 401
			res.Message = "exec gagal"
			res.Data = err.Error()
			return res, err
		}
	} else if kirimnotif.User_id_receiver != 0 {
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
