package model

import (
	"TemplateProject/db"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
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

// CRUD role ============================================================================

func CreateRole(userRole string) (Response, error) {
	var res Response

	type InputRole struct {
		Role_name  string `json:"role"`
		Privileges string `json:"privilege"`
	}
	var dtUserRole InputRole
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

	query := "INSERT INTO role (nama_role) VALUES (?)"
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(dtUserRole.Role_name)
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

	privileges := strings.Split(dtUserRole.Privileges, ",")
	insertPrivilegeQuery := "INSERT INTO role_privilege (id_role, id_privilege) VALUES (?, ?)"

	for _, privilege := range privileges {
		privilegeId := strings.TrimSpace(privilege)

		stmt, err := con.Prepare(insertPrivilegeQuery)
		if err != nil {
			res.Status = 401
			res.Message = "stmt privilege gagal"
			res.Data = err.Error()
			return res, err
		}

		_, err = stmt.Exec(lastId, privilegeId)
		stmt.Close()
		if err != nil {
			res.Status = 401
			res.Message = "exec privilege gagal"
			res.Data = err.Error()
			return res, err
		}
	}

	res.Status = http.StatusOK
	res.Message = "Berhasil membuat role privilege"
	res.Data = map[string]interface{}{
		"id_role":    lastId,
		"role":       dtUserRole.Role_name,
		"privileges": dtUserRole.Privileges,
	}

	defer db.DbClose(con)
	return res, nil
}

func EditRole(userRole string) (Response, error) {
	var res Response

	type InputRole struct {
		// id role bukan role_privilege
		Id_role    string `json:"id"`
		Role_name  string `json:"role"`
		Privileges string `json:"privilege"`
	}
	var dtUserRole InputRole
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

	query := "UPDATE role SET nama_role=? WHERE role_id = ?"
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	_, err = stmt.Exec(dtUserRole.Role_name, dtUserRole.Id_role)
	if err != nil {
		res.Status = 401
		res.Message = "exec gagal"
		res.Data = err.Error()
		return res, err
	}

	deletePrivilegesQuery := "DELETE FROM role_privilege WHERE id_role = ?"
	stmt, err = con.Prepare(deletePrivilegesQuery)
	if err != nil {
		res.Status = 401
		res.Message = "stmt delete gagal"
		res.Data = err.Error()
		return res, err
	}
	_, err = stmt.Exec(dtUserRole.Id_role)
	if err != nil {
		res.Status = 401
		res.Message = "exec delete privilege gagal"
		res.Data = err.Error()
		return res, err
	}
	stmt.Close()

	privileges := strings.Split(dtUserRole.Privileges, ",")
	insertPrivilegeQuery := "INSERT INTO role_privilege (id_role, id_privilege) VALUES (?, ?)"

	for _, privilege := range privileges {
		privilegeId := strings.TrimSpace(privilege)

		stmt, err := con.Prepare(insertPrivilegeQuery)
		if err != nil {
			res.Status = 401
			res.Message = "stmt privilege gagal"
			res.Data = err.Error()
			return res, err
		}

		_, err = stmt.Exec(dtUserRole.Id_role, privilegeId)
		stmt.Close()
		if err != nil {
			res.Status = 401
			res.Message = "exec privilege gagal"
			res.Data = err.Error()
			return res, err
		}
	}

	res.Status = http.StatusOK
	res.Message = "Berhasil membuat role privilege"
	res.Data = map[string]interface{}{
		"id_role":    dtUserRole.Id_role,
		"role":       dtUserRole.Role_name,
		"privileges": dtUserRole.Privileges,
	}

	defer db.DbClose(con)
	return res, nil
}

func GetAllRole() (Response, error) {
	var res Response
	var arrRole = []Role{}

	// Map to track roles and their associated privileges
	roleMap := make(map[int]*Role)

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	query := `
	SELECT r.role_id, r.nama_role, p.nama_privilege
	FROM role_privilege rp
	LEFT JOIN privilege p ON rp.id_privilege = p.privilege_id
	LEFT JOIN role r ON rp.id_role = r.role_id
	ORDER BY r.role_id
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

	// Iterate over the results and populate the roleMap
	for result.Next() {
		var roleId int
		var roleName, privilege string

		err = result.Scan(&roleId, &roleName, &privilege)
		if err != nil {
			res.Status = 401
			res.Message = "rows scan"
			res.Data = err.Error()
			return res, err
		}

		// Check if the role already exists in the map
		if _, exists := roleMap[roleId]; !exists {
			// If not, create a new Role entry and add it to the map
			roleMap[roleId] = &Role{
				Role_id:   roleId,
				Nama_role: roleName,
				Privilege: []string{},
			}
		}

		// Append the privilege to the role's Privileges list
		roleMap[roleId].Privilege = append(roleMap[roleId].Privilege, privilege)
	}

	// Convert the map to a slice for the response
	for _, role := range roleMap {
		arrRole = append(arrRole, *role)
	}

	// Prepare the response
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
