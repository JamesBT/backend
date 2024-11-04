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
	"strings"
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
	tempModal, _ := strconv.Atoi(modal)
	dtPerusahaan.Modal = float64(tempModal)
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

	query := "SELECT perusahaan_id, status, name, username, lokasi, tipe, dokumen_kepemilikan, dokumen_perusahaan, modal_awal, deskripsi, created_at FROM perusahaan WHERE status = 'W'"
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
	type TempUser struct {
		Id               int    `json:"id"`
		Username         string `json:"username"`
		Password         string `json:"password"`
		Nama_lengkap     string `json:"nama_lengkap"`
		Alamat           string `json:"alamat"`
		Jenis_kelamin    string `json:"jenis_kelamin"`
		Tgl_lahir        string `json:"tgl_lahir"`
		Email            string `json:"email"`
		No_telp          string `json:"no_telp"`
		Foto_profil      string `json:"foto_profil"`
		Ktp              string `json:"ktp"`
		Kelas            int    `json:"kelas"`
		Status           string `json:"status"`
		Tipe             int    `json:"tipe"`
		First_login      string `json:"first_login"`
		Denied_by_admin  string `json:"denied_by_admin"`
		UserRole         string `json:"user_role"`
		PerusahaanJoined []Perusahaan
	}
	type TempPerusahaan struct {
		Id                  int     `json:"id_perusahaan"`
		Status              string  `json:"status"`
		Nama                string  `json:"nama"`
		Username            string  `json:"username"`
		Lokasi              string  `json:"lokasi"`
		Kelas               int     `json:"kelas"`
		Tipe                string  `json:"tipe"`
		Dokumen_kepemilikan string  `json:"dokumen_kepemilikan"`
		Dokumen_perusahaan  string  `json:"dokumen_perusahaan"`
		Modal               float64 `json:"modal"`
		Deskripsi           string  `json:"deskripsi"`
		CreatedAt           string  `json:"created_at"`
		Field               []BusinessField
		UserJoined          []TempUser
		LinkedPerusahaan    []TempPerusahaan
		ArchivedAsset       []Asset
	}

	var dtPerusahaan TempPerusahaan
	var dtUser []TempUser
	var linkedPerusahaan []TempPerusahaan
	var archivedAssets []Asset

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka koneksi"
		res.Data = err.Error()
		return res, err
	}

	query := `
		SELECT p.perusahaan_id, p.status, p.name, p.username, p.lokasi, p.tipe, IFNULL(p.kelas,0), p.dokumen_kepemilikan, p.dokumen_perusahaan,
			p.modal_awal,p.deskripsi, p.created_at, p.id_parent, p.id_child, p.archive_asset,
			u.user_id, u.username, u.nama_lengkap, IFNULL(r.nama_role,""), u.foto_profil
		FROM perusahaan p
		LEFT JOIN user_perusahaan up ON p.perusahaan_id = up.id_perusahaan
		LEFT JOIN user u on up.id_user = u.user_id
		LEFT JOIN role r ON up.id_role = r.role_id 
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

	var idParent, idChild, idArchiveAsset sql.NullString

	for rows.Next() {
		var usr TempUser
		err = rows.Scan(
			&dtPerusahaan.Id, &dtPerusahaan.Status, &dtPerusahaan.Nama,
			&dtPerusahaan.Username, &dtPerusahaan.Lokasi, &dtPerusahaan.Tipe,
			&dtPerusahaan.Kelas, &dtPerusahaan.Dokumen_kepemilikan, &dtPerusahaan.Dokumen_perusahaan,
			&dtPerusahaan.Modal, &dtPerusahaan.Deskripsi, &dtPerusahaan.CreatedAt,
			&idParent, &idChild, &idArchiveAsset,
			&usr.Id, &usr.Username, &usr.Nama_lengkap, &usr.UserRole, &usr.Foto_profil)
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

	businessFieldsQuery := `
		SELECT b.id, b.nama 
		FROM perusahaan_business pb
		LEFT JOIN business_field b ON pb.id_business = b.id
		WHERE pb.id_perusahaan = ?
	`
	businessRows, err := con.Query(businessFieldsQuery, nId)
	if err != nil {
		res.Status = 401
		res.Message = "Gagal mengambil business fields"
		res.Data = err.Error()
		return res, err
	}
	defer businessRows.Close()

	var businessFields []BusinessField
	for businessRows.Next() {
		var field BusinessField
		if err := businessRows.Scan(&field.Id, &field.Nama); err != nil {
			res.Status = 401
			res.Message = "Failed to scan business field"
			res.Data = err.Error()
			return res, err
		}
		businessFields = append(businessFields, field)
	}
	dtPerusahaan.Field = businessFields

	if idChild.Valid && idChild.String != "" {
		childIds := strings.Split(idChild.String, ",")
		for _, childId := range childIds {
			var childPerusahaan TempPerusahaan
			childQuery := `
				SELECT perusahaan_id, status, name, username, lokasi, tipe, IFNULL(kelas,0), dokumen_kepemilikan, dokumen_perusahaan, 
				modal_awal, deskripsi, created_at 
				FROM perusahaan WHERE perusahaan_id = ?
			`
			err := con.QueryRow(childQuery, childId).Scan(&childPerusahaan.Id, &childPerusahaan.Status, &childPerusahaan.Nama,
				&childPerusahaan.Username, &childPerusahaan.Lokasi, &childPerusahaan.Tipe,
				&childPerusahaan.Kelas, &childPerusahaan.Dokumen_kepemilikan, &childPerusahaan.Dokumen_perusahaan,
				&childPerusahaan.Modal, &childPerusahaan.Deskripsi, &childPerusahaan.CreatedAt)
			if err != nil {
				res.Status = 401
				res.Message = "Gagal mengambil data perusahaan anak"
				res.Data = err.Error()
				return res, err
			}
			linkedPerusahaan = append(linkedPerusahaan, childPerusahaan)
		}
	}

	if idArchiveAsset.Valid && idArchiveAsset.String != "" {
		assetIds := strings.Split(idArchiveAsset.String, ",")
		for _, childId := range assetIds {
			var childAsset Asset
			childQuery := `
				SELECT id_asset, nama, alamat 
				FROM asset WHERE id_asset = ?
			`
			err := con.QueryRow(childQuery, childId).Scan(&childAsset.Id_asset, &childAsset.Nama, &childAsset.Alamat)
			if err != nil {
				res.Status = 401
				res.Message = "Gagal mengambil data perusahaan anak"
				res.Data = err.Error()
				return res, err
			}
			archivedAssets = append(archivedAssets, childAsset)
		}
	}

	dtPerusahaan.UserJoined = dtUser
	dtPerusahaan.LinkedPerusahaan = linkedPerusahaan
	dtPerusahaan.ArchivedAsset = archivedAssets

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
		WHERE up.id_user = ?
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

	query := `
	SELECT p.perusahaan_id, p.name, 
		COUNT(DISTINCT up.id_user) AS user_count, 
		COUNT(DISTINCT tr.id_asset) AS transaction_count 
	FROM perusahaan p 
	LEFT JOIN user_perusahaan up ON p.perusahaan_id = up.id_perusahaan 
	LEFT JOIN transaction_request tr ON p.perusahaan_id = tr.perusahaan_id 
	WHERE p.status = 'A'
    GROUP BY p.perusahaan_id;
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
		homeuser.SemuaAsset = append(homeuser.SemuaAsset, asset)
	}

	res.Status = http.StatusOK
	res.Message = "Berhasil mengambil data"
	res.Data = homeuser

	defer db.DbClose(con)
	return res, nil
}

func UpdatePerusahaanById(input string) (Response, error) {
	var res Response

	type TempUpdatePerusahaan struct {
		Id        int    `json:"id"`
		Username  string `json:"username"`
		Lokasi    string `json:"lokasi"`
		Tipe      string `json:"tipe"`
		Deskripsi string `json:"deskripsi"`
		Field     string `json:"field"`
		Kelas     string `json:"kelas"`
	}

	var dtPerusahaan TempUpdatePerusahaan
	err := json.Unmarshal([]byte(input), &dtPerusahaan)
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

	query := "UPDATE perusahaan SET username = ?, lokasi=?, tipe=?,kelas=?,deskripsi=?,updated_at=NOW() WHERE perusahaan_id = ?"
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	_, err = stmt.Exec(dtPerusahaan.Username, dtPerusahaan.Lokasi, dtPerusahaan.Tipe, dtPerusahaan.Kelas, dtPerusahaan.Deskripsi, dtPerusahaan.Id)
	if err != nil {
		res.Status = 401
		res.Message = "exec gagal"
		res.Data = err.Error()
		return res, err
	}

	// update field perusahaan
	deleteFieldQuery := `DELETE FROM perusahaan_business WHERE id_perusahaan = ?`
	stmtFieldQuery, err := con.Prepare(deleteFieldQuery)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmtFieldQuery.Close()
	_, err = stmtFieldQuery.Exec(dtPerusahaan.Id)
	if err != nil {
		res.Status = 401
		res.Message = "exec gagal"
		res.Data = err.Error()
		return res, err
	}

	insertFieldQuery := `INSERT INTO perusahaan_business (id_perusahaan,id_business) VALUES (?,?)`

	tempBusinessField := strings.Split(dtPerusahaan.Field, ",")
	for _, idField := range tempBusinessField {
		stmtInsertField, err := con.Prepare(insertFieldQuery)
		if err != nil {
			res.Status = 401
			res.Message = "stmt gagal"
			res.Data = err.Error()
			return res, err
		}

		defer stmtInsertField.Close()
		_, err = stmtInsertField.Exec(dtPerusahaan.Id, idField)
		if err != nil {
			res.Status = 401
			res.Message = "exec gagal"
			res.Data = err.Error()
			return res, err
		}
	}

	res.Status = http.StatusOK
	res.Message = "Berhasil mengupdate data"
	res.Data = dtPerusahaan

	defer db.DbClose(con)
	return res, nil
}

func AddUserCompany(input string) (Response, error) {
	var res Response

	type UserCompany struct {
		Id_user       int `json:"user_id"`
		Id_perusahaan int `json:"perusahaan_id"`
		Id_role       int `json:"role_id"`
	}
	var tempUserCompany UserCompany
	err := json.Unmarshal([]byte(input), &tempUserCompany)
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

	// Check if id_user exists
	var userExists bool
	err = con.QueryRow("SELECT EXISTS(SELECT 1 FROM user WHERE user_id = ?)", tempUserCompany.Id_user).Scan(&userExists)
	if err != nil {
		res.Status = 401
		res.Message = "gagal mengecek user"
		res.Data = err.Error()
		return res, err
	}
	if !userExists {
		res.Status = 404
		res.Message = "User tidak ditemukan"
		res.Data = nil
		return res, nil
	}

	// Check if id_perusahaan exists
	var perusahaanExists bool
	err = con.QueryRow("SELECT EXISTS(SELECT 1 FROM perusahaan WHERE perusahaan_id = ?)", tempUserCompany.Id_perusahaan).Scan(&perusahaanExists)
	if err != nil {
		res.Status = 401
		res.Message = "gagal mengecek perusahaan"
		res.Data = err.Error()
		return res, err
	}
	if !perusahaanExists {
		res.Status = 404
		res.Message = "Perusahaan tidak ditemukan"
		res.Data = nil
		return res, nil
	}

	// Check if id_role exists
	var roleExists bool
	err = con.QueryRow("SELECT EXISTS(SELECT 1 FROM role WHERE role_id = ?)", tempUserCompany.Id_role).Scan(&roleExists)
	if err != nil {
		res.Status = 401
		res.Message = "gagal mengecek role"
		res.Data = err.Error()
		return res, err
	}
	if !roleExists {
		res.Status = 404
		res.Message = "Role tidak ditemukan"
		res.Data = nil
		return res, nil
	}

	// cek sudah tergabung atau belum
	var alreadyJoin bool
	err = con.QueryRow("SELECT EXISTS(SELECT 1 FROM user_perusahaan WHERE id_user = ? AND id_perusahaan = ?)", tempUserCompany.Id_user, tempUserCompany.Id_perusahaan).Scan(&alreadyJoin)
	if err != nil {
		res.Status = 401
		res.Message = "gagal mengecek user dan perusahaan"
		res.Data = err.Error()
		return res, err
	}
	if alreadyJoin {
		res.Status = 404
		res.Message = "user sudah terdaftar pada perusahaan tersebut"
		res.Data = nil
		return res, nil
	}

	query := "INSERT INTO user_perusahaan (`id_user`, `id_perusahaan`, `id_role`) VALUES (?,?,?)"
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	_, err = stmt.Exec(tempUserCompany.Id_user, tempUserCompany.Id_perusahaan, tempUserCompany.Id_role)
	if err != nil {
		res.Status = 401
		res.Message = "exec gagal"
		res.Data = err.Error()
		return res, err
	}

	tempUserid := strconv.Itoa(tempUserCompany.Id_user)
	tempPerusahaanid := strconv.Itoa(tempUserCompany.Id_perusahaan)
	tempRoleid := strconv.Itoa(tempUserCompany.Id_role)

	res.Status = http.StatusOK
	res.Message = "Berhasil memasukkan user ke perusahaan"
	res.Data = map[string]string{
		"id_user":       tempUserid,
		"id_perusahaan": tempPerusahaanid,
		"id_role":       tempRoleid,
	}

	defer db.DbClose(con)
	return res, nil
}

func GetAllPerusahaanJoinedByUserId(user_id string) (Response, error) {
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
	SELECT p.perusahaan_id, p.name
	FROM user_perusahaan up 
	LEFT JOIN perusahaan p ON up.id_perusahaan = p.perusahaan_id
	WHERE up.id_user = ?
	`
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	result, err := stmt.Query(user_id)
	if err != nil {
		res.Status = 401
		res.Message = "exec gagal"
		res.Data = err.Error()
		return res, err
	}
	defer result.Close()

	for result.Next() {
		var dtPerusahaan Perusahaan
		err = result.Scan(&dtPerusahaan.Id, &dtPerusahaan.Nama)
		if err != nil {
			res.Status = 401
			res.Message = "rows scan"
			res.Data = err.Error()
			return res, err
		}
		arrPerusahaan = append(arrPerusahaan, dtPerusahaan)
	}

	if len(arrPerusahaan) == 0 {
		res.Status = 401
		res.Message = "Data tidak ditemukan"
		res.Data = "User tidak tergabung dalam perusahaan mana pun"
		return res, nil
	}

	res.Status = http.StatusOK
	res.Message = "Berhasil mengambil data"
	res.Data = arrPerusahaan

	defer db.DbClose(con)
	return res, nil
}

func JoinPerusahaan(id_perusahaan_1, id_perusahaan_2 string) (Response, error) {
	var res Response
	var perusahaan1Exists, perusahaan2Exists bool
	var currentIdChild string

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "Gagal membuka database"
		res.Data = err.Error()
		return res, err
	}
	defer db.DbClose(con)

	err = con.QueryRow("SELECT EXISTS(SELECT 1 FROM perusahaan WHERE perusahaan_id = ?)", id_perusahaan_1).Scan(&perusahaan1Exists)
	if err != nil {
		res.Status = 401
		res.Message = "Gagal mengecek perusahaan 1"
		res.Data = err.Error()
		return res, err
	}
	if !perusahaan1Exists {
		res.Status = 404
		res.Message = "Perusahaan 1 tidak ditemukan"
		return res, nil
	}

	err = con.QueryRow("SELECT EXISTS(SELECT 1 FROM perusahaan WHERE perusahaan_id = ?)", id_perusahaan_2).Scan(&perusahaan2Exists)
	if err != nil {
		res.Status = 401
		res.Message = "Gagal mengecek perusahaan 2"
		res.Data = err.Error()
		return res, err
	}
	if !perusahaan2Exists {
		res.Status = 404
		res.Message = "Perusahaan 2 tidak ditemukan"
		return res, nil
	}

	err = con.QueryRow("SELECT IFNULL(id_child, '') FROM perusahaan WHERE perusahaan_id = ?", id_perusahaan_1).Scan(&currentIdChild)
	if err != nil {
		res.Status = 500
		res.Message = "Gagal mengecek id_child perusahaan 1"
		res.Data = err.Error()
		return res, err
	}

	tx, err := con.Begin()
	if err != nil {
		res.Status = 500
		res.Message = "Gagal memulai transaksi"
		res.Data = err.Error()
		return res, err
	}

	var newIdChild string
	if currentIdChild == "" {
		newIdChild = id_perusahaan_2
	} else {
		newIdChild = currentIdChild + "," + id_perusahaan_2
	}

	updateQuery1 := `
		UPDATE perusahaan 
		SET id_child = ?, id_role = 'P' 
		WHERE perusahaan_id = ?`
	_, err = tx.Exec(updateQuery1, newIdChild, id_perusahaan_1)
	if err != nil {
		tx.Rollback()
		res.Status = 500
		res.Message = "Gagal memperbarui perusahaan 1"
		res.Data = err.Error()
		return res, err
	}

	updateQuery2 := `
		UPDATE perusahaan 
		SET id_parent = ?, id_role = 'C' 
		WHERE perusahaan_id = ?`
	_, err = tx.Exec(updateQuery2, id_perusahaan_1, id_perusahaan_2)
	if err != nil {
		tx.Rollback()
		res.Status = 500
		res.Message = "Gagal memperbarui perusahaan 2"
		res.Data = err.Error()
		return res, err
	}

	err = tx.Commit()
	if err != nil {
		res.Status = 500
		res.Message = "Gagal mengonfirmasi transaksi"
		res.Data = err.Error()
		return res, err
	}

	res.Status = http.StatusOK
	res.Message = "Berhasil menggabungkan dua perusahaan"
	return res, nil
}

func GetAssetArchiveByPerusahaanId(perusahaan_id string) (Response, error) {
	var res Response
	var dtArchiveAsset []Asset

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka koneksi"
		res.Data = err.Error()
		return res, err
	}

	query := `
		SELECT archive_asset
		FROM perusahaan 
		WHERE perusahaan_id = ? 
	`
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	var idArchiveAsset sql.NullString

	nId, _ := strconv.Atoi(perusahaan_id)
	err = stmt.QueryRow(nId).Scan(&idArchiveAsset)
	if err != nil {
		res.Status = 401
		res.Message = "rows gagal"
		res.Data = err.Error()
		return res, err
	}

	if idArchiveAsset.Valid && idArchiveAsset.String != "" {
		assetIds := strings.Split(idArchiveAsset.String, ",")
		for _, childId := range assetIds {
			var childAsset Asset
			childQuery := `
				SELECT id_asset, nama, alamat 
				FROM asset WHERE id_asset = ?
			`
			err := con.QueryRow(childQuery, childId).Scan(&childAsset.Id_asset, &childAsset.Nama, &childAsset.Alamat)
			if err != nil {
				res.Status = 401
				res.Message = "Gagal mengambil data perusahaan anak"
				res.Data = err.Error()
				return res, err
			}
			dtArchiveAsset = append(dtArchiveAsset, childAsset)
		}
	}

	res.Status = http.StatusOK
	res.Message = "Berhasil mengambil data"
	res.Data = dtArchiveAsset

	defer db.DbClose(con)

	return res, nil
}

func AddAssetArchive(input string) (Response, error) {
	var res Response
	var perusahaan1Exists bool
	var currentAsset string

	type TempUpdatePerusahaan struct {
		Perusahaan_id string `json:"perusahaan"`
		Asset_id      string `json:"asset"`
	}
	var dtPerusahaan TempUpdatePerusahaan
	err := json.Unmarshal([]byte(input), &dtPerusahaan)
	if err != nil {
		res.Status = 401
		res.Message = "gagal decode json"
		res.Data = err.Error()
		return res, err
	}

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "Gagal membuka database"
		res.Data = err.Error()
		return res, err
	}
	defer db.DbClose(con)

	err = con.QueryRow("SELECT EXISTS(SELECT 1 FROM perusahaan WHERE perusahaan_id = ?)", dtPerusahaan.Perusahaan_id).Scan(&perusahaan1Exists)
	if err != nil {
		res.Status = 401
		res.Message = "Gagal mengecek perusahaan 1"
		res.Data = err.Error()
		return res, err
	}
	if !perusahaan1Exists {
		res.Status = 404
		res.Message = "Perusahaan 1 tidak ditemukan"
		return res, nil
	}

	err = con.QueryRow("SELECT IFNULL(archive_asset, '') FROM perusahaan WHERE perusahaan_id = ?", dtPerusahaan.Perusahaan_id).Scan(&currentAsset)
	if err != nil {
		res.Status = 500
		res.Message = "Gagal mengecek archive asset perusahaan"
		res.Data = err.Error()
		return res, err
	}

	tx, err := con.Begin()
	if err != nil {
		res.Status = 500
		res.Message = "Gagal memulai transaksi"
		res.Data = err.Error()
		return res, err
	}

	var newIdChild string
	if currentAsset == "" {
		newIdChild = dtPerusahaan.Asset_id
	} else {
		newIdChild = currentAsset + "," + dtPerusahaan.Asset_id
	}

	updateQuery1 := `
		UPDATE perusahaan 
		SET archive_asset = ? 
		WHERE perusahaan_id = ?`
	_, err = tx.Exec(updateQuery1, newIdChild, dtPerusahaan.Perusahaan_id)
	if err != nil {
		tx.Rollback()
		res.Status = 500
		res.Message = "Gagal memperbarui archive asset"
		res.Data = err.Error()
		return res, err
	}

	err = tx.Commit()
	if err != nil {
		res.Status = 500
		res.Message = "Gagal mengonfirmasi transaksi"
		res.Data = err.Error()
		return res, err
	}

	res.Status = http.StatusOK
	res.Message = "Berhasil menambahkan aset archive"
	return res, nil
}

func RemoveAssetArchive(input string) (Response, error) {
	var res Response
	var perusahaan1Exists bool
	var currentAsset string

	type TempUpdatePerusahaan struct {
		Perusahaan_id string `json:"perusahaan"`
		Asset_id      string `json:"asset"`
	}
	var dtPerusahaan TempUpdatePerusahaan
	err := json.Unmarshal([]byte(input), &dtPerusahaan)
	if err != nil {
		res.Status = 401
		res.Message = "gagal decode json"
		res.Data = err.Error()
		return res, err
	}

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "Gagal membuka database"
		res.Data = err.Error()
		return res, err
	}
	defer db.DbClose(con)

	err = con.QueryRow("SELECT EXISTS(SELECT 1 FROM perusahaan WHERE perusahaan_id = ?)", dtPerusahaan.Perusahaan_id).Scan(&perusahaan1Exists)
	if err != nil {
		res.Status = 401
		res.Message = "Gagal mengecek perusahaan 1"
		res.Data = err.Error()
		return res, err
	}
	if !perusahaan1Exists {
		res.Status = 404
		res.Message = "Perusahaan 1 tidak ditemukan"
		return res, nil
	}

	err = con.QueryRow("SELECT IFNULL(archive_asset, '') FROM perusahaan WHERE perusahaan_id = ?", dtPerusahaan.Perusahaan_id).Scan(&currentAsset)
	if err != nil {
		res.Status = 500
		res.Message = "Gagal mengecek archive asset perusahaan"
		res.Data = err.Error()
		return res, err
	}

	tx, err := con.Begin()
	if err != nil {
		res.Status = 500
		res.Message = "Gagal memulai transaksi"
		res.Data = err.Error()
		return res, err
	}

	assets := strings.Split(currentAsset, ",")
	var newAssets []string
	for _, asset := range assets {
		if asset != dtPerusahaan.Asset_id {
			newAssets = append(newAssets, asset)
		}
	}

	newArchiveAsset := strings.Join(newAssets, ",")

	updateQuery1 := `
		UPDATE perusahaan 
		SET archive_asset = ? 
		WHERE perusahaan_id = ?`
	_, err = tx.Exec(updateQuery1, newArchiveAsset, dtPerusahaan.Perusahaan_id)
	if err != nil {
		tx.Rollback()
		res.Status = 500
		res.Message = "Gagal memperbarui archive asset"
		res.Data = err.Error()
		return res, err
	}

	err = tx.Commit()
	if err != nil {
		res.Status = 500
		res.Message = "Gagal mengonfirmasi transaksi"
		res.Data = err.Error()
		return res, err
	}

	res.Status = http.StatusOK
	res.Message = "Berhasil menambahkan aset archive"
	return res, nil
}

func LoginPerusahaan(id_perusahaan, id_user string) (Response, error) {
	var res Response
	type TempUser struct {
		Id               int    `json:"id"`
		Username         string `json:"username"`
		Password         string `json:"password"`
		Nama_lengkap     string `json:"nama_lengkap"`
		Alamat           string `json:"alamat"`
		Jenis_kelamin    string `json:"jenis_kelamin"`
		Tgl_lahir        string `json:"tgl_lahir"`
		Email            string `json:"email"`
		No_telp          string `json:"no_telp"`
		Foto_profil      string `json:"foto_profil"`
		Ktp              string `json:"ktp"`
		Kelas            int    `json:"kelas"`
		Status           string `json:"status"`
		Tipe             int    `json:"tipe"`
		First_login      string `json:"first_login"`
		Denied_by_admin  string `json:"denied_by_admin"`
		UserRole         string `json:"user_role"`
		PerusahaanJoined []Perusahaan
	}
	type TempPerusahaan struct {
		Id                  int     `json:"id_perusahaan"`
		Status              string  `json:"status"`
		Nama                string  `json:"nama"`
		Username            string  `json:"username"`
		Lokasi              string  `json:"lokasi"`
		Kelas               int     `json:"kelas"`
		Tipe                string  `json:"tipe"`
		Dokumen_kepemilikan string  `json:"dokumen_kepemilikan"`
		Dokumen_perusahaan  string  `json:"dokumen_perusahaan"`
		Modal               float64 `json:"modal"`
		Deskripsi           string  `json:"deskripsi"`
		CreatedAt           string  `json:"created_at"`
		Role_user           Role    `json:"role_user"`
		Field               []BusinessField
		UserJoined          []TempUser
		LinkedPerusahaan    []TempPerusahaan
		ArchivedAsset       []Asset
	}

	var dtPerusahaan TempPerusahaan
	var dtUser []TempUser
	var linkedPerusahaan []TempPerusahaan
	var archivedAssets []Asset

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka koneksi"
		res.Data = err.Error()
		return res, err
	}

	query := `
		SELECT p.perusahaan_id, p.status, p.name, p.username, p.lokasi, p.tipe, IFNULL(p.kelas,0), p.dokumen_kepemilikan, p.dokumen_perusahaan,
			p.modal_awal,p.deskripsi, p.created_at, p.id_parent, p.id_child, p.archive_asset,
			u.user_id, u.username, u.nama_lengkap, IFNULL(r.nama_role,""), IFNULL(up.id_role,0), IFNULL(r.nama_role,'')
		FROM perusahaan p
		LEFT JOIN user_perusahaan up ON p.perusahaan_id = up.id_perusahaan
		LEFT JOIN user u on up.id_user = u.user_id
		LEFT JOIN role r ON up.id_role = r.role_id 
		WHERE p.perusahaan_id = ? AND up.id_user = ?
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
	rows, err := stmt.Query(nId, id_user)
	if err != nil {
		res.Status = 401
		res.Message = "rows gagal"
		res.Data = err.Error()
		return res, err
	}
	defer rows.Close()

	var idParent, idChild, idArchiveAsset sql.NullString

	for rows.Next() {
		var usr TempUser
		var tempRole Role
		err = rows.Scan(
			&dtPerusahaan.Id, &dtPerusahaan.Status, &dtPerusahaan.Nama,
			&dtPerusahaan.Username, &dtPerusahaan.Lokasi, &dtPerusahaan.Tipe,
			&dtPerusahaan.Kelas, &dtPerusahaan.Dokumen_kepemilikan, &dtPerusahaan.Dokumen_perusahaan,
			&dtPerusahaan.Modal, &dtPerusahaan.Deskripsi, &dtPerusahaan.CreatedAt,
			&idParent, &idChild, &idArchiveAsset,
			&usr.Id, &usr.Username, &usr.Nama_lengkap, &usr.UserRole,
			&tempRole.Role_id, &tempRole.Nama_role)
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

	if idChild.Valid && idChild.String != "" {
		childIds := strings.Split(idChild.String, ",")
		for _, childId := range childIds {
			var childPerusahaan TempPerusahaan
			childQuery := `
				SELECT perusahaan_id, status, name, username, lokasi, tipe, IFNULL(kelas,0), dokumen_kepemilikan, dokumen_perusahaan, 
				modal_awal, deskripsi, created_at 
				FROM perusahaan WHERE perusahaan_id = ?
			`
			err := con.QueryRow(childQuery, childId).Scan(&childPerusahaan.Id, &childPerusahaan.Status, &childPerusahaan.Nama,
				&childPerusahaan.Username, &childPerusahaan.Lokasi, &childPerusahaan.Tipe,
				&childPerusahaan.Kelas, &childPerusahaan.Dokumen_kepemilikan, &childPerusahaan.Dokumen_perusahaan,
				&childPerusahaan.Modal, &childPerusahaan.Deskripsi, &childPerusahaan.CreatedAt)
			if err != nil {
				res.Status = 401
				res.Message = "Gagal mengambil data perusahaan anak"
				res.Data = err.Error()
				return res, err
			}
			linkedPerusahaan = append(linkedPerusahaan, childPerusahaan)
		}
	}

	if idArchiveAsset.Valid && idArchiveAsset.String != "" {
		assetIds := strings.Split(idArchiveAsset.String, ",")
		for _, childId := range assetIds {
			var childAsset Asset
			childQuery := `
				SELECT id_asset, nama, alamat 
				FROM asset WHERE id_asset = ?
			`
			err := con.QueryRow(childQuery, childId).Scan(&childAsset.Id_asset, &childAsset.Nama, &childAsset.Alamat)
			if err != nil {
				res.Status = 401
				res.Message = "Gagal mengambil data perusahaan anak"
				res.Data = err.Error()
				return res, err
			}
			archivedAssets = append(archivedAssets, childAsset)
		}
	}

	dtPerusahaan.UserJoined = dtUser
	dtPerusahaan.LinkedPerusahaan = linkedPerusahaan
	dtPerusahaan.ArchivedAsset = archivedAssets

	res.Status = http.StatusOK
	res.Message = "Berhasil mengambil data"
	res.Data = dtPerusahaan

	defer db.DbClose(con)

	return res, nil
}
