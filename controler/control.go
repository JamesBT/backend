package controler

import (
	"TemplateProject/model"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

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
	akun, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Gagal membaca body request"})
	}
	result, err := model.UpdateUser(string(akun))
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
