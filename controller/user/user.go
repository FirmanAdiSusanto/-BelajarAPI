package user

import (
	"BelajarAPIi/helper"
	"BelajarAPIi/middlewares"
	"BelajarAPIi/model"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type UserController struct {
	Model model.UserModel
}

func (us *UserController) Register() echo.HandlerFunc {
	return func(c echo.Context) error {
		var input RegisterRequest
		err := c.Bind(&input)
		if err != nil {
			if strings.Contains(err.Error(), "unsupport") {
				return c.JSON(http.StatusUnsupportedMediaType,
					helper.ResponseFormat(http.StatusUnsupportedMediaType, "format data tidak didukung", nil))
			}
			return c.JSON(http.StatusBadRequest,
				helper.ResponseFormat(http.StatusBadRequest, "data yang dikirimkan tidak sesuai", nil))
		}

		// Validasi panjang password
		if len(input.Password) < 10 {
			return c.JSON(http.StatusBadRequest,
				helper.ResponseFormat(http.StatusBadRequest, "Panjang password harus minimal 10 karakter", nil))
		}

		// Validasi setidaknya satu huruf besar
		hasUppercase := false
		for _, char := range input.Password {
			if char >= 'A' && char <= 'Z' {
				hasUppercase = true
				break
			}
		}
		if !hasUppercase {
			return c.JSON(http.StatusBadRequest,
				helper.ResponseFormat(http.StatusBadRequest, "Password harus mengandung setidaknya satu huruf besar", nil))
		}

		validate := validator.New(validator.WithRequiredStructEnabled())
		err = validate.Struct(input)
		if err != nil {
			return c.JSON(http.StatusBadRequest,
				helper.ResponseFormat(http.StatusBadRequest, "Data yang dikirim tidak sesuai", nil))
		}

		var processInput model.User
		processInput.Hp = input.Hp
		processInput.Nama = input.Nama
		processInput.Password = input.Password

		// Memeriksa apakah nomor sudah terdaftar
		if us.Model.CekUser(processInput.Hp) {
			return c.JSON(http.StatusConflict,
				helper.ResponseFormat(http.StatusConflict, "Nomor sudah terdaftar", nil))
		}

		err = us.Model.AddUser(processInput) // Fungsi yang Anda buat sendiri
		if err != nil {
			return c.JSON(http.StatusInternalServerError,
				helper.ResponseFormat(http.StatusInternalServerError, "Terjadi kesalahan pada sistem", nil))
		}
		return c.JSON(http.StatusCreated,
			helper.ResponseFormat(http.StatusCreated, "Selamat data sudah terdaftar", nil))
	}
}

func (us *UserController) Login() echo.HandlerFunc {
	return func(c echo.Context) error {
		var input LoginRequest
		err := c.Bind(&input)
		if err != nil {
			if strings.Contains(err.Error(), "unsupport") {
				return c.JSON(http.StatusUnsupportedMediaType,
					helper.ResponseFormat(http.StatusUnsupportedMediaType, "format data tidak didukung", nil))
			}
			return c.JSON(http.StatusBadRequest,
				helper.ResponseFormat(http.StatusBadRequest, "data yang dikirmkan tidak sesuai", nil))
		}

		validate := validator.New(validator.WithRequiredStructEnabled())
		err = validate.Struct(input)

		if err != nil {
			for _, val := range err.(validator.ValidationErrors) {
				fmt.Println(val.Error())
			}
		}

		result, err := us.Model.Login(input.Hp, input.Password)
		if err != nil {
			return c.JSON(http.StatusInternalServerError,
				helper.ResponseFormat(http.StatusInternalServerError, "terjadi kesalahan pada sistem", nil))
		}
		token, err := middlewares.GenerateJWT(result.Hp)
		if err != nil {
			return c.JSON(http.StatusInternalServerError,
				helper.ResponseFormat(http.StatusInternalServerError, "terjadi kesalahan pada sistem, gagal memproses data", nil))
		}

		var responseData LoginResponse
		responseData.Hp = result.Hp
		responseData.Nama = result.Nama
		responseData.Token = token

		return c.JSON(http.StatusOK,
			helper.ResponseFormat(http.StatusOK, "selamat anda berhasil login", responseData))

	}
}

func (us *UserController) Update() echo.HandlerFunc {
	return func(c echo.Context) error {
		var hp = c.Param("hp")
		var input model.User
		err := c.Bind(&input)
		if err != nil {
			log.Println("masalah baca input:", err.Error())
			if strings.Contains(err.Error(), "unsupport") {
				return c.JSON(http.StatusUnsupportedMediaType,
					helper.ResponseFormat(http.StatusUnsupportedMediaType, "format data tidak didukung", nil))
			}
			return c.JSON(http.StatusBadRequest,
				helper.ResponseFormat(http.StatusBadRequest, "data yang dikirmkan tidak sesuai", nil))
		}

		isFound := us.Model.CekUser(hp)

		if !isFound {
			return c.JSON(http.StatusNotFound,
				helper.ResponseFormat(http.StatusNotFound, "data tidak ditemukan", nil))
		}

		err = us.Model.Update(hp, input)

		if err != nil {
			log.Println("masalah database :", err.Error())
			return c.JSON(http.StatusInternalServerError,
				helper.ResponseFormat(http.StatusInternalServerError, "terjadi kesalahan saat update data", nil))
		}

		return c.JSON(http.StatusOK,
			helper.ResponseFormat(http.StatusOK, "data berhasil di update", nil))
	}
}

func (us *UserController) ListUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		listUser, err := us.Model.GetAllUser()
		if err != nil {
			return c.JSON(http.StatusInternalServerError,
				helper.ResponseFormat(http.StatusInternalServerError, "terjadi kesalahan pada sistem", nil))
		}
		return c.JSON(http.StatusOK,
			helper.ResponseFormat(http.StatusOK, "berhasil mendapatkan data", listUser))
	}
}

func (us *UserController) Profile() echo.HandlerFunc {
	return func(c echo.Context) error {
		var hp = c.Param("hp")
		result, err := us.Model.GetProfile(hp)

		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				return c.JSON(http.StatusNotFound,
					helper.ResponseFormat(http.StatusNotFound, "data tidak ditemukan", nil))
			}
			return c.JSON(http.StatusInternalServerError,
				helper.ResponseFormat(http.StatusInternalServerError, "terjadi kesalahan pada sistem", nil))
		}
		return c.JSON(http.StatusOK,
			helper.ResponseFormat(http.StatusOK, "berhasil mendapatkan data", result))
	}
}

// Fungsi Controller untuk menambah AddActivity
func (us *UserController) AddActivity() echo.HandlerFunc {
	return func(c echo.Context) error {
		// Mengambil token JWT dari konteks
		token := c.Get("user").(*jwt.Token)

		// Mengambil claims dari token JWT
		claims := token.Claims.(jwt.MapClaims)

		// Mengambil nomor HP dari claims
		hp := claims["hp"].(string)

		var newActivity model.Activity
		if err := c.Bind(&newActivity); err != nil {
			return c.JSON(http.StatusBadRequest, helper.ResponseFormat(http.StatusBadRequest, "Gagal memproses permintaan", nil))
		}

		if err := us.Model.AddActivity(hp, newActivity); err != nil {
			return c.JSON(http.StatusInternalServerError, helper.ResponseFormat(http.StatusInternalServerError, "Gagal menambahkan", nil))
		}

		return c.JSON(http.StatusCreated, helper.ResponseFormat(http.StatusCreated, "Berhasil Ditambahkan", nil))
	}
}

func (us *UserController) UpdateActivity() echo.HandlerFunc {
	return func(c echo.Context) error {
		// Mengambil token JWT dari konteks
		token := c.Get("user").(*jwt.Token)

		// Mengambil claims dari token JWT
		claims := token.Claims.(jwt.MapClaims)

		// Mengambil nomor HP dari claims
		hp := claims["hp"].(string)

		// Mendapatkan ID kegiatan dari parameter URL
		id := c.Param("id")

		var updatedActivity model.Activity
		if err := c.Bind(&updatedActivity); err != nil {
			return c.JSON(http.StatusBadRequest, helper.ResponseFormat(http.StatusBadRequest, "Gagal memproses permintaan", nil))
		}

		// Update the activity only if it belongs to the current user (hp)
		if err := us.Model.UpdateActivity(hp, id, updatedActivity); err != nil {
			return c.JSON(http.StatusInternalServerError, helper.ResponseFormat(http.StatusInternalServerError, "Gagal mengupdate kegiatan", nil))
		}

		return c.JSON(http.StatusOK, helper.ResponseFormat(http.StatusOK, "Kegiatan berhasil diupdate", nil))
	}
}

func (us *UserController) GetActivities() echo.HandlerFunc {
	return func(c echo.Context) error {
		// Mengambil token JWT dari konteks
		token := c.Get("user").(*jwt.Token)

		// Mengambil claims dari token JWT
		claims := token.Claims.(jwt.MapClaims)

		// Mengambil nomor HP dari claims
		hp := claims["hp"].(string)

		// Get activities for the current user (hp)
		activities, err := us.Model.GetActivities(hp)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, helper.ResponseFormat(http.StatusInternalServerError, "Gagal mendapatkan kegiatan", nil))
		}

		return c.JSON(http.StatusOK, helper.ResponseFormat(http.StatusOK, "Berhasil mendapatkan kegiatan", activities))
	}
}
