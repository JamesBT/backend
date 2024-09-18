package model

import (
	"TemplateProject/db"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
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
	var dtPerusahaan Perusahaan
	var dtUser []User

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka koneksi"
		res.Data = err.Error()
		return res, err
	}

	query := `
		SELECT p.perusahaan_id, p.status, p.name, p.username, p.lokasi, p.tipe, IFNULL(p.kelas,0), p.dokumen_kepemilikan, p.dokumen_perusahaan,
			p.modal_awal,p.deskripsi, p.created_at, u.user_id, u.username, u.nama_lengkap 
		FROM perusahaan p
		LEFT JOIN user_perusahaan up ON p.perusahaan_id = up.id_perusahaan
		LEFT JOIN user u on up.id_user = u.user_id
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

	for rows.Next() {
		var usr User
		err = rows.Scan(&dtPerusahaan.Id, &dtPerusahaan.Status, &dtPerusahaan.Nama, &dtPerusahaan.Username, &dtPerusahaan.Lokasi, &dtPerusahaan.Tipe,
			&dtPerusahaan.Kelas,
			&dtPerusahaan.Dokumen_kepemilikan, &dtPerusahaan.Dokumen_perusahaan, &dtPerusahaan.Modal,
			&dtPerusahaan.Deskripsi, &dtPerusahaan.CreatedAt, &usr.Id, &usr.Username, &usr.Nama_lengkap)
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

	dtPerusahaan.UserJoined = dtUser

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
		COUNT(DISTINCT tr.id_transaksi_jual_sewa) AS transaction_count 
	FROM perusahaan p 
	LEFT JOIN user_perusahaan up ON p.perusahaan_id = up.id_perusahaan 
	LEFT JOIN transaction_request tr ON p.perusahaan_id = tr.perusahaan_id 
	WHERE p.status = 'A'
	GROUP BY p.perusahaan_id
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
		fmt.Println("1")
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
		Modal     string `json:"modal"`
		Deskripsi string `json:"deskripsi"`
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

	query := "UPDATE perusahaan SET username = ?, lokasi=?, tipe=?,modal_awal=?,deskripsi=?,updated_at=NOW() WHERE perusahaan_id = ?"
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	_, err = stmt.Exec(dtPerusahaan.Username, dtPerusahaan.Lokasi, dtPerusahaan.Tipe, dtPerusahaan.Modal, dtPerusahaan.Deskripsi, dtPerusahaan.Id)
	if err != nil {
		res.Status = 401
		res.Message = "exec gagal"
		res.Data = err.Error()
		return res, err
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
		Id_user       string `user_id`
		Id_perusahaan string `perusahaan_id`
		Id_role       string `role_id`
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
	err = con.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE user_id = ?)", tempUserCompany.Id_user).Scan(&userExists)
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
	err = con.QueryRow("SELECT EXISTS(SELECT 1 FROM roles WHERE role_id = ?)", tempUserCompany.Id_role).Scan(&roleExists)
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

	res.Status = http.StatusOK
	res.Message = "Berhasil memasukkan user ke perusahaan"
	res.Data = map[string]string{
		"id_user":       tempUserCompany.Id_user,
		"id_perusahaan": tempUserCompany.Id_perusahaan,
		"id_role":       tempUserCompany.Id_role,
	}

	defer db.DbClose(con)
	return res, nil
}
