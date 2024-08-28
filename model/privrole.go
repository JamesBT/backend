package model

import (
	"TemplateProject/db"
	"encoding/json"
	"net/http"
	"strconv"
)

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
