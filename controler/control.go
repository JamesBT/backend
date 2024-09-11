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
	provinsi := c.FormValue("provinsi")
	surat_legalitas := c.FormValue("surat_legalitas")
	gambar_asset, err := c.FormFile("gambar_asset")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Gagal membaca body request"})
	}
	usage := c.FormValue("usage")
	tags := c.FormValue("tags")
	result, err := model.CreateAsset(
		fileLegalitas, suratKuasa, gambar_asset, asetName, surat_legalitas, tipe, usage, tags, nomorLegalitas, status, alamat,
		kondisi, koordinat, batasKoordinat, luas, nilai, provinsi)
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
	surat_legalitas := c.FormValue("surat_legalitas")
	gambar_asset, err := c.FormFile("gambar_asset")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Gagal membaca body request"})
	}
	usage := c.FormValue("usage")
	tags := c.FormValue("tags")
	result, err := model.CreateAssetChild(
		fileLegalitas, suratKuasa, gambar_asset, parentId, asetName, surat_legalitas, tipe, usage, tags, nomorLegalitas, status, alamat,
		kondisi, koordinat, batasKoordinat, luas, nilai)
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

func GetAssetSurveyHistoryByAssetId(c echo.Context) error {
	id := c.Param("id")
	result, err := model.GetAssetSurveyHistoryByAssetId(id)
	if err != nil {
		return c.JSON(result.Status, map[string]string{"message": err.Error()})
	}
	ip := c.RealIP()
	model.InsertLog(ip, "UploadFoto", result.Data, 3)
	return c.JSON(http.StatusOK, result)
}

func GetAssetChildByParentId(c echo.Context) error {
	id := c.Param("id")
	result, err := model.GetAssetChildByParentId(id)
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

func GetAssetDetailedByPerusahaanId(c echo.Context) error {
	id := c.Param("id")
	result, err := model.GetAssetDetailedByPerusahaanId(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}
	ip := c.RealIP()
	model.InsertLog(ip, "UploadFoto", result.Data, 3)
	return c.JSON(http.StatusOK, result)
}

func GetAssetDetailedByUserId(c echo.Context) error {
	id := c.Param("id")
	result, err := model.GetAssetDetailedByUserId(id)
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

func JoinAsset(c echo.Context) error {
	input, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Gagal membaca body request"})
	}
	result, err := model.JoinAsset(string(input))
	if err != nil {
		return c.JSON(result.Status, map[string]string{"message": err.Error()})
	}
	ip := c.RealIP()
	model.InsertLog(ip, "UploadFoto", result.Data, 3)
	return c.JSON(http.StatusOK, result)
}

func GetAssetRentedByUserId(c echo.Context) error {
	id := c.Param("id")
	result, err := model.GetAssetRentedByUserId(id)
	if err != nil {
		return c.JSON(result.Status, map[string]string{"message": err.Error()})
	}
	ip := c.RealIP()
	model.InsertLog(ip, "UploadFoto", result.Data, 3)
	return c.JSON(http.StatusOK, result)
}

func UpdateAssetByIdWithoutGambar(c echo.Context) error {
	asetId := c.FormValue("id")
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
	surat_legalitas := c.FormValue("surat_legalitas")
	usage := c.FormValue("usage")
	tags := c.FormValue("tags")
	result, err := model.UpdateAssetByIdWithoutGambar(
		fileLegalitas, suratKuasa, asetId, asetName, surat_legalitas, tipe, usage, tags, nomorLegalitas, status, alamat,
		kondisi, koordinat, batasKoordinat, luas, nilai)
	if err != nil {
		return c.JSON(result.Status, map[string]string{"message": err.Error()})
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

func GetPerusahaanDetailById(c echo.Context) error {
	id := c.Param("id")
	result, err := model.GetPerusahaanDetailById(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}
	ip := c.RealIP()
	model.InsertLog(ip, "UploadFoto", result.Data, 3)
	return c.JSON(http.StatusOK, result)
}

func HomeUserPerusahaan(c echo.Context) error {
	id := c.Param("id")
	result, err := model.HomeUserPerusahaan(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}
	ip := c.RealIP()
	model.InsertLog(ip, "UploadFoto", result.Data, 3)
	return c.JSON(http.StatusOK, result)
}

func UpdatePerusahaanById(c echo.Context) error {
	akun, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Gagal membaca body request"})
	}
	result, err := model.UpdatePerusahaanById(string(akun))
	if err != nil {
		return c.JSON(result.Status, map[string]string{"message": err.Error()})
	}
	ip := c.RealIP()
	model.InsertLog(ip, "UploadFoto", result.Data, 3)
	return c.JSON(http.StatusOK, result)
}

func GetAllPerusahaanJoinedByUserId(c echo.Context) error {
	id := c.Param("id")
	result, err := model.GetAllPerusahaanJoinedByUserId(id)
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

func GetSurveyorByUserId(c echo.Context) error {
	id := c.Param("id")
	result, err := model.GetSurveyorByUserId(string(id))
	if err != nil {
		return c.JSON(result.Status, map[string]string{"message": err.Error()})
	}
	ip := c.RealIP()
	model.InsertLog(ip, "UploadFoto", result.Data, 3)
	return c.JSON(http.StatusOK, result)
}

func UpdateUserBySurveyorId(c echo.Context) error {
	input, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Gagal membaca body request"})
	}
	result, err := model.UpdateUserBySurveyorId(string(input))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}
	ip := c.RealIP()
	model.InsertLog(ip, "UploadFoto", result.Data, 3)
	return c.JSON(http.StatusOK, result)
}

func UpdateSurveyorByUserId(c echo.Context) error {
	input, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Gagal membaca body request"})
	}
	result, err := model.UpdateSurveyorByUserId(string(input))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
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

func GetAllSurveyReq(c echo.Context) error {
	result, err := model.GetAllSurveyReq()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}
	ip := c.RealIP()
	model.InsertLog(ip, "UploadFoto", result.Data, 3)
	return c.JSON(http.StatusOK, result)
}

func GetAllSurveyReqDetailed(c echo.Context) error {
	result, err := model.GetAllSurveyReqDetailed()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}
	ip := c.RealIP()
	model.InsertLog(ip, "UploadFoto", result.Data, 3)
	return c.JSON(http.StatusOK, result)
}

func GetSurveyReqById(c echo.Context) error {
	id := c.Param("id")
	result, err := model.GetSurveyReqById(id)
	if err != nil {
		return c.JSON(result.Status, map[string]string{"message": err.Error()})
	}
	ip := c.RealIP()
	model.InsertLog(ip, "UploadFoto", result.Data, 3)
	return c.JSON(http.StatusOK, result)
}

func GetSurveyReqByAsetNama(c echo.Context) error {
	name := c.Param("nama")
	result, err := model.GetSurveyReqByAsetName(name)
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

func GetAllSurveyReqByUserId(c echo.Context) error {
	id := c.Param("id")
	result, err := model.GetAllSurveyReqByUserId(id)
	if err != nil {
		return c.JSON(result.Status, map[string]string{"message": err.Error()})
	}
	ip := c.RealIP()
	model.InsertLog(ip, "UploadFoto", result.Data, 3)
	return c.JSON(http.StatusOK, result)
}

func GetAllOngoingSurveyReqByUserId(c echo.Context) error {
	id := c.Param("id")
	result, err := model.GetAllOngoingSurveyReqByUserId(id)
	if err != nil {
		return c.JSON(result.Status, map[string]string{"message": err.Error()})
	}
	ip := c.RealIP()
	model.InsertLog(ip, "UploadFoto", result.Data, 3)
	return c.JSON(http.StatusOK, result)
}

func GetAllFinishedSurveyReqByUserId(c echo.Context) error {
	id := c.Param("id")
	result, err := model.GetAllFinishedSurveyReqByUserId(id)
	if err != nil {
		return c.JSON(result.Status, map[string]string{"message": err.Error()})
	}
	ip := c.RealIP()
	model.InsertLog(ip, "UploadFoto", result.Data, 3)
	return c.JSON(http.StatusOK, result)
}

func SubmitSurveyReqById(c echo.Context) error {
	surveyreq, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Gagal membaca body request"})
	}
	result, err := model.SubmitSurveyReqById(string(surveyreq))
	if err != nil {
		return c.JSON(result.Status, map[string]string{"message": err.Error()})
	}
	ip := c.RealIP()
	model.InsertLog(ip, "UploadFoto", result.Data, 3)
	return c.JSON(http.StatusOK, result)
}

// transaction_request
func CreateTranReq(c echo.Context) error {
	// id_asset, user_id, perusahaan_id, nama_progress, tgl_meeting, lokasi_meeting, deskripsi
	idAsset := c.FormValue("idAsset")
	idUser := c.FormValue("idUser")
	idPerusahaan := c.FormValue("idPerusahaan")
	nama := c.FormValue("nama")
	tgl_meeting := c.FormValue("tgl_meeting")
	waktu_meeting := c.FormValue("waktu_meeting")
	lokasi_meeting := c.FormValue("lokasi_meeting")
	deskripsi := c.FormValue("deskripsi")
	proposal, err := c.FormFile("proposal")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Gagal membaca body request"})
	}
	result, err := model.CreateTranReq(proposal, idAsset, idUser, idPerusahaan, nama, tgl_meeting, waktu_meeting, lokasi_meeting, deskripsi)
	if err != nil {
		return c.JSON(result.Status, map[string]string{"message": err.Error()})
	}
	ip := c.RealIP()
	model.InsertLog(ip, "UploadFoto", result.Data, 3)
	return c.JSON(http.StatusOK, result)
}

func GetTranReqById(c echo.Context) error {
	id := c.Param("id")
	result, err := model.GetTranReqById(id)
	if err != nil {
		return c.JSON(result.Status, map[string]string{"message": err.Error()})
	}
	ip := c.RealIP()
	model.InsertLog(ip, "UploadFoto", result.Data, 3)
	return c.JSON(http.StatusOK, result)
}

func GetAllTranReq(c echo.Context) error {
	result, err := model.GetAllTranReq()
	if err != nil {
		return c.JSON(result.Status, map[string]string{"message": err.Error()})
	}
	ip := c.RealIP()
	model.InsertLog(ip, "UploadFoto", result.Data, 3)
	return c.JSON(http.StatusOK, result)
}

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

func GetTranReqByPerusahaanId(c echo.Context) error {
	id := c.Param("id")
	result, err := model.GetTranReqByPerusahaanId(id)
	if err != nil {
		return c.JSON(result.Status, map[string]string{"message": err.Error()})
	}
	ip := c.RealIP()
	model.InsertLog(ip, "UploadFoto", result.Data, 3)
	return c.JSON(http.StatusOK, result)
}

func GetAllUserTransaction(c echo.Context) error {
	result, err := model.GetAllUserTransaction()
	if err != nil {
		return c.JSON(result.Status, map[string]string{"message": err.Error()})
	}
	ip := c.RealIP()
	model.InsertLog(ip, "UploadFoto", result.Data, 3)
	return c.JSON(http.StatusOK, result)
}

func UserManagementGetMeetingByUserId(c echo.Context) error {
	id := c.Param("id")
	result, err := model.UserManagementGetMeetingByUserId(id)
	if err != nil {
		return c.JSON(result.Status, map[string]string{"message": err.Error()})
	}
	ip := c.RealIP()
	model.InsertLog(ip, "UploadFoto", result.Data, 3)
	return c.JSON(http.StatusOK, result)
}

func UserManagementGetMeetingByPerusahaanId(c echo.Context) error {
	id := c.Param("id")
	result, err := model.UserManagementGetMeetingByPerusahaanId(id)
	if err != nil {
		return c.JSON(result.Status, map[string]string{"message": err.Error()})
	}
	ip := c.RealIP()
	model.InsertLog(ip, "UploadFoto", result.Data, 3)
	return c.JSON(http.StatusOK, result)
}

func AcceptTransaction(c echo.Context) error {
	input, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Gagal membaca body request"})
	}
	result, err := model.AcceptTransaction(string(input))
	if err != nil {
		return c.JSON(result.Status, map[string]string{"message": err.Error()})
	}
	ip := c.RealIP()
	model.InsertLog(ip, "UploadFoto", result.Data, 3)
	return c.JSON(http.StatusOK, result)
}

func DeclineTransaction(c echo.Context) error {
	input, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Gagal membaca body request"})
	}
	result, err := model.DeclineTransaction(string(input))
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

func GetAllUser(c echo.Context) error {
	result, err := model.GetAllUser()
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

func GetUserDetailedById(c echo.Context) error {
	id := c.Param("id")
	result, err := model.GetUserDetailedById(id)
	if err != nil {
		return c.JSON(result.Status, map[string]string{"message": err.Error()})
	}
	ip := c.RealIP()
	model.InsertLog(ip, "UploadFoto", result.Data, 3)
	return c.JSON(http.StatusOK, result)
}

// admin - management user
func UpdateUserById(c echo.Context) error {
	datauser, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Gagal membaca body request"})
	}
	result, err := model.UpdateUserById(string(datauser))
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

func GetAllUserByPerusahaanId(c echo.Context) error {
	id := c.Param("id")
	result, err := model.GetAllUserByPerusahaanId(id)
	if err != nil {
		return c.JSON(result.Status, map[string]string{"message": err.Error()})
	}
	ip := c.RealIP()
	model.InsertLog(ip, "UploadFoto", result.Data, 3)
	return c.JSON(http.StatusOK, result)
}

func AdminUserManagement(c echo.Context) error {
	result, err := model.AdminUserManagement()
	if err != nil {
		return c.JSON(result.Status, map[string]string{"message": err.Error()})
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

func GetAllVerify(c echo.Context) error {
	result, err := model.GetAllVerify()
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

func GetVerifyPerusahaanDetailedById(c echo.Context) error {
	id := c.Param("id")
	result, err := model.GetVerifyPerusahaanDetailedById(id)
	if err != nil {
		return c.JSON(result.Status, map[string]string{"message": err.Error()})
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

func VerifyOTP(c echo.Context) error {
	input, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Gagal membaca body request"})
	}
	result, err := model.VerifyOTP(string(input))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}
	ip := c.RealIP()
	model.InsertLog(ip, "UploadFoto", result.Data, 3)
	return c.JSON(http.StatusOK, result)
}

func CreateNotification(c echo.Context) error {
	notif, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Gagal membaca body request"})
	}
	result, err := model.CreateNotification(string(notif))
	if err != nil {
		return c.JSON(result.Status, map[string]string{"message": err.Error()})
	}
	ip := c.RealIP()
	model.InsertLog(ip, "UploadFoto", result.Data, 3)
	return c.JSON(http.StatusOK, result)
}

func GetNotificationById(c echo.Context) error {
	id := c.Param("id")
	result, err := model.GetNotificationById(id)
	if err != nil {
		return c.JSON(result.Status, map[string]string{"message": err.Error()})
	}
	ip := c.RealIP()
	model.InsertLog(ip, "UploadFoto", result.Data, 3)
	return c.JSON(http.StatusOK, result)
}

func GetNotificationByUserIdReceiver(c echo.Context) error {
	id := c.Param("id")
	result, err := model.GetNotificationByUserIdReceiver(id)
	if err != nil {
		return c.JSON(result.Status, map[string]string{"message": err.Error()})
	}
	ip := c.RealIP()
	model.InsertLog(ip, "UploadFoto", result.Data, 3)
	return c.JSON(http.StatusOK, result)
}

func GetNotificationByPerusahaanIdReceiver(c echo.Context) error {
	id := c.Param("id")
	result, err := model.GetNotificationByPerusahaanIdReceiver(id)
	if err != nil {
		return c.JSON(result.Status, map[string]string{"message": err.Error()})
	}
	ip := c.RealIP()
	model.InsertLog(ip, "UploadFoto", result.Data, 3)
	return c.JSON(http.StatusOK, result)
}

// kelas
func GetAllKelas(c echo.Context) error {
	result, err := model.GetAllKelas()
	if err != nil {
		return c.JSON(result.Status, map[string]string{"message": err.Error()})
	}
	ip := c.RealIP()
	model.InsertLog(ip, "UploadFoto", result.Data, 3)
	return c.JSON(http.StatusOK, result)
}

// business field
func GetAllBusinessField(c echo.Context) error {
	result, err := model.GetAllBusinessField()
	if err != nil {
		return c.JSON(result.Status, map[string]string{"message": err.Error()})
	}
	ip := c.RealIP()
	model.InsertLog(ip, "UploadFoto", result.Data, 3)
	return c.JSON(http.StatusOK, result)
}

func GetFile(c echo.Context) error {
	path := c.FormValue("path")
	return c.File(path)
}

func CreateMeeting(c echo.Context) error {
	id := c.FormValue("id")
	tanggal_meeting := c.FormValue("tanggal_meeting")
	waktu_meeting := c.FormValue("waktu_meeting")
	tempat_meeting := c.FormValue("tempat_meeting")
	waktu_mulai := c.FormValue("waktu_mulai")
	waktu_selesai := c.FormValue("waktu_selesai")
	notes := c.FormValue("notes")
	tipe_dokumen := c.FormValue("tipe_dokumen")
	dokumen, err := c.FormFile("dokumen")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Gagal membaca body request"})
	}
	result, err := model.CreateMeeting(dokumen, id, tanggal_meeting, waktu_meeting, tempat_meeting, waktu_mulai, waktu_selesai, notes, tipe_dokumen)
	if err != nil {
		return c.JSON(result.Status, map[string]string{"message": err.Error()})
	}

	ip := c.RealIP()
	model.InsertLog(ip, "UploadFoto", result.Data, 3)
	return c.JSON(http.StatusOK, result)
}

func CreateMeetingWithoutDocument(c echo.Context) error {
	id := c.FormValue("id")
	tanggal_meeting := c.FormValue("tanggal_meeting")
	waktu_meeting := c.FormValue("waktu_meeting")
	tempat_meeting := c.FormValue("tempat_meeting")
	waktu_mulai := c.FormValue("waktu_mulai")
	waktu_selesai := c.FormValue("waktu_selesai")
	notes := c.FormValue("notes")
	result, err := model.CreateMeetingWithoutDocument(id, tanggal_meeting, waktu_meeting, tempat_meeting, waktu_mulai, waktu_selesai, notes)
	if err != nil {
		return c.JSON(result.Status, map[string]string{"message": err.Error()})
	}

	ip := c.RealIP()
	model.InsertLog(ip, "UploadFoto", result.Data, 3)
	return c.JSON(http.StatusOK, result)
}

func GetAllProgress(c echo.Context) error {
	result, err := model.GetAllProgress()
	if err != nil {
		return c.JSON(result.Status, map[string]string{"message": err.Error()})
	}
	ip := c.RealIP()
	model.InsertLog(ip, "UploadFoto", result.Data, 3)
	return c.JSON(http.StatusOK, result)
}

func GetProgressByUserId(c echo.Context) error {
	id := c.Param("id")
	result, err := model.GetProgressByUserId(id)
	if err != nil {
		return c.JSON(result.Status, map[string]string{"message": err.Error()})
	}
	ip := c.RealIP()
	model.InsertLog(ip, "UploadFoto", result.Data, 3)
	return c.JSON(http.StatusOK, result)
}

func GetProgressNotDoneByUserId(c echo.Context) error {
	id := c.Param("id")
	result, err := model.GetProgressNotDoneByUserId(id)
	if err != nil {
		return c.JSON(result.Status, map[string]string{"message": err.Error()})
	}
	ip := c.RealIP()
	model.InsertLog(ip, "UploadFoto", result.Data, 3)
	return c.JSON(http.StatusOK, result)
}

func GetProgressById(c echo.Context) error {
	id := c.Param("id")
	result, err := model.GetProgressById(id)
	if err != nil {
		return c.JSON(result.Status, map[string]string{"message": err.Error()})
	}
	ip := c.RealIP()
	model.InsertLog(ip, "UploadFoto", result.Data, 3)
	return c.JSON(http.StatusOK, result)
}

func GetProgressByUserAsetId(c echo.Context) error {
	id := c.Param("id")
	aset := c.Param("aset")
	result, err := model.GetProgressByUserAsetId(id, aset)
	if err != nil {
		return c.JSON(result.Status, map[string]string{"message": err.Error()})
	}
	ip := c.RealIP()
	model.InsertLog(ip, "UploadFoto", result.Data, 3)
	return c.JSON(http.StatusOK, result)
}

func GetAllUsage(c echo.Context) error {
	result, err := model.GetAllUsage()
	if err != nil {
		return c.JSON(result.Status, map[string]string{"message": err.Error()})
	}
	ip := c.RealIP()
	model.InsertLog(ip, "UploadFoto", result.Data, 3)
	return c.JSON(http.StatusOK, result)
}

func GetAllTags(c echo.Context) error {
	result, err := model.GetAllTags()
	if err != nil {
		return c.JSON(result.Status, map[string]string{"message": err.Error()})
	}
	ip := c.RealIP()
	model.InsertLog(ip, "UploadFoto", result.Data, 3)
	return c.JSON(http.StatusOK, result)
}

func GetAllProvinsi(c echo.Context) error {
	result, err := model.GetAllProvinsi()
	if err != nil {
		return c.JSON(result.Status, map[string]string{"message": err.Error()})
	}
	ip := c.RealIP()
	model.InsertLog(ip, "UploadFoto", result.Data, 3)
	return c.JSON(http.StatusOK, result)
}
