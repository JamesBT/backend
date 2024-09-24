package model

import (
	"TemplateProject/db"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sort"
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
		Admin_role string `json:"admin_role"`
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

	query := "INSERT INTO role (nama_role,admin_role) VALUES (?,?)"
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(dtUserRole.Role_name, dtUserRole.Admin_role)
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
	insertPrivilegeQuery := "INSERT INTO role_privilege_user (id_role, id_privilege) VALUES (?, ?)"

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
		// id role bukan role_privilege_user
		Id_role    int    `json:"id"`
		Role_name  string `json:"role"`
		Admin_role string `json:"admin_role"`
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
	fmt.Println(dtUserRole.Id_role)
	fmt.Println(dtUserRole.Role_name)
	fmt.Println(dtUserRole.Privileges)

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	query := "UPDATE role SET `nama_role`=?,`admin_role`=? WHERE `role_id` = ?"
	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	_, err = stmt.Exec(dtUserRole.Role_name, dtUserRole.Admin_role, dtUserRole.Id_role)
	if err != nil {
		res.Status = 401
		res.Message = "exec gagal"
		res.Data = err.Error()
		return res, err
	}

	deletePrivilegesQuery := "DELETE FROM role_privilege_user WHERE `id_role` = ?"
	deletestmt, err := con.Prepare(deletePrivilegesQuery)
	if err != nil {
		res.Status = 401
		res.Message = "stmt delete gagal"
		res.Data = err.Error()
		return res, err
	}
	_, err = deletestmt.Exec(dtUserRole.Id_role)
	if err != nil {
		res.Status = 401
		res.Message = "exec delete privilege gagal"
		res.Data = err.Error()
		return res, err
	}
	deletestmt.Close()

	privileges := strings.Split(dtUserRole.Privileges, ",")
	insertPrivilegeQuery := "INSERT INTO role_privilege_user (id_role, id_privilege) VALUES (?, ?)"

	for _, privilege := range privileges {
		privilegeId := strings.TrimSpace(privilege)

		privstmt, err := con.Prepare(insertPrivilegeQuery)
		if err != nil {
			res.Status = 401
			res.Message = "stmt privilege gagal"
			res.Data = err.Error()
			return res, err
		}

		_, err = privstmt.Exec(dtUserRole.Id_role, privilegeId)
		privstmt.Close()
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

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	query := `
	SELECT role_id,nama_role
	FROM role
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
		var dtRole Role

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

func GetAllRoleAdmin() (Response, error) {
	var res Response
	var arrRole = []Role{}

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	// Modify the query to join role, role_privilege_user, and privilege tables
	query := `
	SELECT r.role_id, r.nama_role, p.nama_privilege
	FROM role r
	LEFT JOIN role_privilege_user rp ON r.role_id = rp.id_role
	LEFT JOIN privilege p ON rp.id_privilege = p.privilege_id
	WHERE r.admin_role = 'Y'
	ORDER BY r.role_id ASC
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

	roleMap := make(map[int]Role)
	var roleIds []int

	for result.Next() {
		var roleId int
		var namaRole, namaPrivilege sql.NullString

		err = result.Scan(&roleId, &namaRole, &namaPrivilege)
		if err != nil {
			res.Status = 401
			res.Message = "Rows scan gagal"
			res.Data = err.Error()
			return res, err
		}

		role, exists := roleMap[roleId]
		if !exists {
			role = Role{
				Role_id:   roleId,
				Nama_role: namaRole.String,
			}
			roleMap[roleId] = role
			roleIds = append(roleIds, roleId)
		}

		if namaPrivilege.Valid {
			role.Privilege = append(role.Privilege, namaPrivilege.String)
			roleMap[roleId] = role
		}
	}

	sort.Ints(roleIds)

	for _, id := range roleIds {
		arrRole = append(arrRole, roleMap[id])
	}

	res.Status = http.StatusOK
	res.Message = "Berhasil mengambil data"
	res.Data = arrRole

	defer db.DbClose(con)
	return res, nil
}

func GetAllPrivRole() (Response, error) {
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
	SELECT r.role_id, r.nama_role, IFNULL(p.nama_privilege,"")
	FROM role r
	LEFT JOIN role_privilege_user rp ON r.role_id = rp.id_role
	LEFT JOIN privilege p ON rp.id_privilege = p.privilege_id
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
		if privilege != "" {
			roleMap[roleId].Privilege = append(roleMap[roleId].Privilege, privilege)
		}
	}

	var roleIds []int
	for roleId := range roleMap {
		roleIds = append(roleIds, roleId)
	}
	sort.Ints(roleIds)
	for _, roleId := range roleIds {
		arrRole = append(arrRole, *roleMap[roleId])
	}

	// Prepare the response
	res.Status = http.StatusOK
	res.Message = "Berhasil mengambil data"
	res.Data = arrRole

	defer db.DbClose(con)
	return res, nil
}

func GetPrivRoleById(role_id string) (Response, error) {
	var res Response
	type TempRole struct {
		Role_id   int    `json:"role_id"`
		Nama_role string `json:"nama_role"`
		Privilege []int  `json:"privilege"`
	}
	var dtRole TempRole
	var privileges []int

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	query := `
	SELECT r.role_id, r.nama_role, IFNULL(p.privilege_id,0)
	FROM role r
	LEFT JOIN role_privilege_user rp ON r.role_id = rp.id_role
	LEFT JOIN privilege p ON rp.id_privilege = p.privilege_id
	WHERE r.role_id = ?
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

	rows, err := stmt.Query(role_id)
	if err != nil {
		res.Status = 401
		res.Message = "exec gagal"
		res.Data = err.Error()
		return res, err
	}
	defer rows.Close()

	// Variables to hold row data
	var privilegeId int
	isFirstRow := true

	for rows.Next() {
		err := rows.Scan(&dtRole.Role_id, &dtRole.Nama_role, &privilegeId)
		if err != nil {
			res.Status = 401
			res.Message = "gagal scan data"
			res.Data = err.Error()
			return res, err
		}

		if isFirstRow {
			isFirstRow = false
		}

		if privilegeId != 0 {
			privileges = append(privileges, privilegeId)
		}
	}
	dtRole.Privilege = privileges

	res.Status = http.StatusOK
	res.Message = "Berhasil mengambil data"
	res.Data = dtRole

	defer db.DbClose(con)
	return res, nil
}

func GetUserRoleByPerusahaanId(perusahaan_id, user_id int) (Response, error) {
	var res Response
	var dtRole Role
	var privileges []string

	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	query := `
	SELECT r.role_id, r.nama_role, IFNULL(rp.id_privilege,"")
	FROM user_perusahaan up
	LEFT JOIN role r ON up.id_role = r.role_id
	LEFT JOIN role_privilege_user rp ON r.role_id = rp.id_role
	LEFT JOIN privilege p ON rp.id_privilege = p.privilege_id
	WHERE up.id_perusahaan = ? AND up.id_user = ?
	`

	stmt, err := con.Prepare(query)
	if err != nil {
		res.Status = 401
		res.Message = "stmt gagal"
		res.Data = err.Error()
		return res, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(perusahaan_id, user_id)
	if err != nil {
		res.Status = 401
		res.Message = "exec gagal"
		res.Data = err.Error()
		return res, err
	}
	defer rows.Close()

	// Variables to hold row data
	var privilegeId string
	isFirstRow := true

	for rows.Next() {
		err := rows.Scan(&dtRole.Role_id, &dtRole.Nama_role, &privilegeId)
		if err != nil {
			res.Status = 401
			res.Message = "gagal scan data"
			res.Data = err.Error()
			return res, err
		}

		if isFirstRow {
			isFirstRow = false
		}

		if privilegeId != "" {
			privileges = append(privileges, privilegeId)
		}
	}
	dtRole.Privilege = privileges

	res.Status = http.StatusOK
	res.Message = "Berhasil mengambil data"
	res.Data = dtRole

	defer db.DbClose(con)
	return res, nil
}

func DeleteRoleById(id string) (Response, error) {
	var res Response
	con, err := db.DbConnection()
	if err != nil {
		res.Status = 401
		res.Message = "gagal membuka database"
		res.Data = err.Error()
		return res, err
	}

	checkAdminQuery := `
	SELECT role_id
	FROM user_role
	WHERE role_id = ?
	`
	var roleID string
	err = con.QueryRow(checkAdminQuery, id).Scan(&roleID)
	if err != nil {
		if err == sql.ErrNoRows {
		} else {
			res.Status = 401
			res.Message = "Gagal memeriksa role user"
			res.Data = err.Error()
			return res, err
		}
	} else {
		res.Status = 403
		res.Message = "Role masih terpakai di user lain"
		res.Data = nil
		return res, errors.New(res.Message)
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

	result, err := stmt.Exec(id)
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
