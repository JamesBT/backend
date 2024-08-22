package controler

import (
	"TemplateProject/model"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// asset
func TambahAsset(c echo.Context) error {
	asetName := c.FormValue("nama")
	tipe := c.FormValue("tipe")
	nomorLegalitas := c.FormValue("nomor_legalitas")
	fileLegalitas, err := c.FormFile("file_legalitas")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Gagal membaca body request"})
	}
	status := c.FormValue("status")
	suratKuasa, err := c.FormFile("surat_kuasa")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Gagal membaca body request"})
	}
	alamat := c.FormValue("alamat")
	kondisi := c.FormValue("kondisi")
	koordinat := c.FormValue("titik_koordinat")
	batasKoordinat := c.FormValue("batas_koordinat")
	luas := c.FormValue("luas")
	nilai := c.FormValue("nilai")
	perusahaan_id := c.FormValue("perusahaan_id")
	result, err := model.CreateAsset(fileLegalitas, suratKuasa, asetName, perusahaan_id, tipe, nomorLegalitas, status, alamat, kondisi, koordinat, batasKoordinat, luas, nilai)
	if err != nil {
		return c.JSON(result.Status, map[string]string{"message": err.Error()})
	}

	ip := c.RealIP()
	model.InsertLog(ip, "UploadFoto", result.Data, 3)
	return c.JSON(http.StatusOK, result)
}

func TambahAssetChild(c echo.Context) error {
	parentId := c.FormValue("parentId")
	asetName := c.FormValue("nama")
	tipe := c.FormValue("tipe")
	nomorLegalitas := c.FormValue("nomor_legalitas")
	fileLegalitas, err := c.FormFile("file_legalitas")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Gagal membaca body request"})
	}
	status := c.FormValue("status")
	suratKuasa, err := c.FormFile("surat_kuasa")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Gagal membaca body request"})
	}
	alamat := c.FormValue("alamat")
	kondisi := c.FormValue("kondisi")
	koordinat := c.FormValue("titik_koordinat")
	batasKoordinat := c.FormValue("batas_koordinat")
	luas := c.FormValue("luas")
	nilai := c.FormValue("nilai")
	perusahaan_id := c.FormValue("perusahaan_id")
	result, err := model.CreateAssetChild(fileLegalitas, suratKuasa, parentId, asetName, perusahaan_id, tipe, nomorLegalitas, status, alamat, kondisi, koordinat, batasKoordinat, luas, nilai)
	if err != nil {
		return c.JSON(result.Status, map[string]string{"message": err.Error()})
	}

	ip := c.RealIP()
	model.InsertLog(ip, "UploadFoto", result.Data, 3)
	return c.JSON(http.StatusOK, result)
}

func GetAllAsset(c echo.Context) error {
	result, err := model.GetAllAsset()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}
	ip := c.RealIP()
	model.InsertLog(ip, "UploadFoto", result.Data, 3)
	return c.JSON(http.StatusOK, result)
}

func GetAssetById(c echo.Context) error {
	id := c.Param("id")
	result, err := model.GetAssetById(id)
	if err != nil {
		return c.JSON(result.Status, map[string]string{"message": err.Error()})
	}
	ip := c.RealIP()
	model.InsertLog(ip, "UploadFoto", result.Data, 3)
	return c.JSON(http.StatusOK, result)
}

func GetAssetByName(c echo.Context) error {
	nama := c.Param("nama")
	result, err := model.GetAssetByName(nama)
	if err != nil {
		return c.JSON(result.Status, map[string]string{"message": err.Error()})
	}
	ip := c.RealIP()
	model.InsertLog(ip, "UploadFoto", result.Data, 3)
	return c.JSON(http.StatusOK, result)
}

func GetAssetDetailedById(c echo.Context) error {
	id := c.Param("id")
	result, err := model.GetAssetDetailedById(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}
	ip := c.RealIP()
	model.InsertLog(ip, "UploadFoto", result.Data, 3)
	return c.JSON(http.StatusOK, result)
}

func UbahVisibilitasAset(c echo.Context) error {
	id := c.Param("id")
	input, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Gagal membaca body request"})
	}
	result, err := model.UbahVisibilitasAset(id, string(input))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}
	ip := c.RealIP()
	model.InsertLog(ip, "UploadFoto", result.Data, 3)
	return c.JSON(http.StatusOK, result)
}

// perusahaan
func TambahPerusahaan(c echo.Context) error {
	userid := c.FormValue("userid")
	nama := c.FormValue("nama")
	username := c.FormValue("username")
	lokasi := c.FormValue("lokasi")
	tipe := c.FormValue("tipe")
	dokumen_kepemilikan, err := c.FormFile("dokumen_kepemilikan")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Gagal membaca body request"})
	}
	dokumen_perusahaan, err := c.FormFile("dokumen_perusahaan")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Gagal membaca body request"})
	}
	modal := c.FormValue("modal")
	deskripsi := c.FormValue("deskripsi")
	result, err := model.CreatePerusahaan(dokumen_kepemilikan, dokumen_perusahaan, userid, nama, username, lokasi, tipe, modal, deskripsi)
	if err != nil {
		return c.JSON(result.Status, map[string]string{"message": err.Error()})
	}

	ip := c.RealIP()
	model.InsertLog(ip, "UploadFoto", result.Data, 3)
	return c.JSON(http.StatusOK, result)
}

func GetAllPerusahaanUnverified(c echo.Context) error {
	result, err := model.GetAllPerusahaanUnverified()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}
	ip := c.RealIP()
	model.InsertLog(ip, "UploadFoto", result.Data, 3)
	return c.JSON(http.StatusOK, result)
}

func GetAllPerusahaanDetailed(c echo.Context) error {
	result, err := model.GetAllPerusahaanDetailed()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}
	ip := c.RealIP()
	model.InsertLog(ip, "UploadFoto", result.Data, 3)
	return c.JSON(http.StatusOK, result)
}

func GetPerusahaanByUserId(c echo.Context) error {
	id := c.Param("id")
	result, err := model.GetPerusahaanByUserId(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}
	ip := c.RealIP()
	model.InsertLog(ip, "UploadFoto", result.Data, 3)
	return c.JSON(http.StatusOK, result)
}

// privilege

// role

// surveyor
func LoginSurveyor(c echo.Context) error {
	akun, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Gagal membaca body request"})
	}
	result, err := model.LoginSurveyor(string(akun))
	if err != nil {
		return c.JSON(result.Status, map[string]string{"message": err.Error()})
	}
	ip := c.RealIP()
	model.InsertLog(ip, "UploadFoto", result.Data, 3)
	return c.JSON(http.StatusOK, result)
}

func SignUpSurveyor(c echo.Context) error {
	akun, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Gagal membaca body request"})
	}
	result, err := model.SignUpSurveyor(string(akun))
	if err != nil {
		return c.JSON(result.Status, map[string]string{"message": err.Error()})
	}
	ip := c.RealIP()
	model.InsertLog(ip, "UploadFoto", result.Data, 3)
	return c.JSON(http.StatusOK, result)
}

func GetSurveyorByName(c echo.Context) error {
	nama := c.Param("nama")
	result, err := model.GetSurveyorByName(nama)
	if err != nil {
		return c.JSON(result.Status, map[string]string{"message": err.Error()})
	}
	ip := c.RealIP()
	model.InsertLog(ip, "UploadFoto", result.Data, 3)
	return c.JSON(http.StatusOK, result)
}

func GetAllSurveyor(c echo.Context) error {
	result, err := model.GetAllSurveyor()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}
	ip := c.RealIP()
	model.InsertLog(ip, "UploadFoto", result.Data, 3)
	return c.JSON(http.StatusOK, result)
}

func GetAllSurveyorDetailed(c echo.Context) error {
	result, err := model.GetAllSurveyorDetailed()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}
	ip := c.RealIP()
	model.InsertLog(ip, "UploadFoto", result.Data, 3)
	return c.JSON(http.StatusOK, result)
}

func UpdateSurveyorById(c echo.Context) error {
	akun, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Gagal membaca body request"})
	}
	result, err := model.UpdateSurveyorById(string(akun))
	if err != nil {
		return c.JSON(result.Status, map[string]string{"message": err.Error()})
	}
	ip := c.RealIP()
	model.InsertLog(ip, "UploadFoto", result.Data, 3)
	return c.JSON(http.StatusOK, result)
}

// survey_request
func CreateSurveyReq(c echo.Context) error {
	surveyreq, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Gagal membaca body request"})
	}
	result, err := model.CreateSurveyReq(string(surveyreq))
	if err != nil {
		return c.JSON(result.Status, map[string]string{"message": err.Error()})
	}
	ip := c.RealIP()
	model.InsertLog(ip, "UploadFoto", result.Data, 3)
	return c.JSON(http.StatusOK, result)
}

func GetSurveyReqByAsetId(c echo.Context) error {
	id := c.Param("id")
	result, err := model.GetSurveyReqByAsetId(id)
	if err != nil {
		return c.JSON(result.Status, map[string]string{"message": err.Error()})
	}
	ip := c.RealIP()
	model.InsertLog(ip, "UploadFoto", result.Data, 3)
	return c.JSON(http.StatusOK, result)
}

// transaction_request
func GetTranReqByUserId(c echo.Context) error {
	id := c.Param("id")
	result, err := model.GetTranReqByUserId(id)
	if err != nil {
		return c.JSON(result.Status, map[string]string{"message": err.Error()})
	}
	ip := c.RealIP()
	model.InsertLog(ip, "UploadFoto", result.Data, 3)
	return c.JSON(http.StatusOK, result)
}

// user
func Login(c echo.Context) error {
	akun, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Gagal membaca body request"})
	}
	result, err := model.Login(string(akun))
	if err != nil {
		return c.JSON(result.Status, map[string]string{"message": err.Error()})
	}
	ip := c.RealIP()
	model.InsertLog(ip, "UploadFoto", result.Data, 3)
	return c.JSON(http.StatusOK, result)
}

func SignUp(c echo.Context) error {
	akun, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Gagal membaca body request"})
	}
	result, err := model.SignUp(string(akun))
	if err != nil {
		return c.JSON(result.Status, map[string]string{"message": err.Error()})
	}
	ip := c.RealIP()
	model.InsertLog(ip, "UploadFoto", result.Data, 3)
	return c.JSON(http.StatusOK, result)
}

func GetUserById(c echo.Context) error {
	id := c.Param("id")
	result, err := model.GetUserById(id)
	if err != nil {
		return c.JSON(result.Status, map[string]string{"message": err.Error()})
	}
	ip := c.RealIP()
	model.InsertLog(ip, "UploadFoto", result.Data, 3)
	return c.JSON(http.StatusOK, result)
}

func UpdateUser(c echo.Context) error {
	userId := c.FormValue("id")
	username := c.FormValue("username")
	nama_lengkap := c.FormValue("nama_lengkap")
	alamat := c.FormValue("alamat")
	jenis_kelamin := c.FormValue("jenis_kelamin")
	tanggal_lahir := c.FormValue("tanggal_lahir")
	email := c.FormValue("email")
	no_telp := c.FormValue("no_telp")
	fileFoto, err := c.FormFile("fileFoto")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Gagal membaca body request"})
	}
	result, err := model.UpdateUser(fileFoto, userId, username, nama_lengkap, alamat, jenis_kelamin, tanggal_lahir, email, no_telp)
	if err != nil {
		return c.JSON(result.Status, map[string]string{"message": err.Error()})
	}

	ip := c.RealIP()
	model.InsertLog(ip, "UploadFoto", result.Data, 3)
	return c.JSON(http.StatusOK, result)
}

func UpdateUserFull(c echo.Context) error {
	userId := c.FormValue("id")
	username := c.FormValue("username")
	nama_lengkap := c.FormValue("nama_lengkap")
	alamat := c.FormValue("alamat")
	jenis_kelamin := c.FormValue("jenis_kelamin")
	tanggal_lahir := c.FormValue("tanggal_lahir")
	email := c.FormValue("email")
	no_telp := c.FormValue("no_telp")
	fileFoto, err := c.FormFile("fileFoto")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Gagal membaca body request"})
	}
	fileKTP, err := c.FormFile("fileKTP")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Gagal membaca body request"})
	}
	result, err := model.UpdateUserFull(fileFoto, fileKTP, userId, username, nama_lengkap, alamat, jenis_kelamin, tanggal_lahir, email, no_telp)
	if err != nil {
		return c.JSON(result.Status, map[string]string{"message": err.Error()})
	}

	ip := c.RealIP()
	model.InsertLog(ip, "UploadFoto", result.Data, 3)
	return c.JSON(http.StatusOK, result)
}

func DeleteUserById(c echo.Context) error {
	id := c.Param("id")
	result, err := model.DeleteUserById(id)
	if err != nil {
		return c.JSON(result.Status, map[string]string{"message": err.Error()})
	}
	ip := c.RealIP()
	model.InsertLog(ip, "UploadFoto", result.Data, 3)
	return c.JSON(http.StatusOK, result)
}

func GetUserKTP(c echo.Context) error {
	id := c.Param("id")

	fmt.Println(id)
	result, err := model.GetUserKTP(id)
	if err != nil {
		return c.JSON(result.Status, map[string]string{"message": err.Error()})
	}

	fmt.Println(result.Data)

	path, ok := result.Data.(string)
	if !ok {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Data is not a valid file path"})
	}

	ip := c.RealIP()
	model.InsertLog(ip, "UploadFoto", result.Data, 3)

	return c.File(path)
}

func GetUserFoto(c echo.Context) error {
	id := c.Param("id")

	fmt.Println(id)
	result, err := model.GetUserFoto(id)
	if err != nil {
		return c.JSON(result.Status, map[string]string{"message": err.Error()})
	}

	fmt.Println(result.Data)

	path, ok := result.Data.(string)
	if !ok {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Data is not a valid file path"})
	}

	ip := c.RealIP()
	model.InsertLog(ip, "UploadFoto", result.Data, 3)

	return c.File(path)
}

func GetAllUserUnverified(c echo.Context) error {
	result, err := model.GetAllUserUnverified()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}
	ip := c.RealIP()
	model.InsertLog(ip, "UploadFoto", result.Data, 3)
	return c.JSON(http.StatusOK, result)
}

func VerifyUserAccept(c echo.Context) error {
	input, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Gagal membaca body request"})
	}
	result, err := model.VerifyUserAccept(string(input))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}
	ip := c.RealIP()
	model.InsertLog(ip, "UploadFoto", result.Data, 3)
	return c.JSON(http.StatusOK, result)
}

func VerifyUserDecline(c echo.Context) error {
	input, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Gagal membaca body request"})
	}
	result, err := model.VerifyUserDecline(string(input))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}
	ip := c.RealIP()
	model.InsertLog(ip, "UploadFoto", result.Data, 3)
	return c.JSON(http.StatusOK, result)
}

func VerifyPerusahaanAccept(c echo.Context) error {
	input, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Gagal membaca body request"})
	}
	result, err := model.VerifyPerusahaanAccept(string(input))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}
	ip := c.RealIP()
	model.InsertLog(ip, "UploadFoto", result.Data, 3)
	return c.JSON(http.StatusOK, result)
}

func VerifyPerusahaanDecline(c echo.Context) error {
	input, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Gagal membaca body request"})
	}
	result, err := model.VerifyPerusahaanDecline(string(input))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}
	ip := c.RealIP()
	model.InsertLog(ip, "UploadFoto", result.Data, 3)
	return c.JSON(http.StatusOK, result)
}

func VerifyAssetAccept(c echo.Context) error {
	input, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Gagal membaca body request"})
	}
	result, err := model.VerifyAssetAccept(string(input))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}
	ip := c.RealIP()
	model.InsertLog(ip, "UploadFoto", result.Data, 3)
	return c.JSON(http.StatusOK, result)
}

func VerifyAssetReassign(c echo.Context) error {
	input, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Gagal membaca body request"})
	}
	result, err := model.ReassignAsset(string(input))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}
	ip := c.RealIP()
	model.InsertLog(ip, "UploadFoto", result.Data, 3)
	return c.JSON(http.StatusOK, result)
}

// user_privilege

// user_role

// fungsi tambahan
func GetFoto(c echo.Context) error {
	id := c.FormValue("id")
	folder := c.FormValue("folder")
	result := model.GetPhoto(folder, id)
	ip := c.RealIP()
	model.InsertLog(ip, "GetPhotoFolder", "getfoto("+id+")", 1)
	return c.File(result)
}

func UploadFoto(c echo.Context) error {
	folder := c.FormValue("folder")
	id := c.FormValue("id")
	fotoFile, err := c.FormFile("photo")
	nId, _ := strconv.Atoi(id)
	tId := int64(nId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}
	result, err := model.UploadFotoFolder(fotoFile, tId, folder)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}
	ip := c.RealIP()
	model.InsertLog(ip, "UploadFoto", result.Data, 3)
	return c.JSON(http.StatusOK, result)
}
