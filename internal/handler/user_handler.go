package handler

import (
	"net/http"
	"strconv"

	"go-base-project/internal/model"
	"go-base-project/internal/service"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// Register godoc
// @Summary      Daftarkan user baru
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body body model.RegisterRequest true "Data registrasi"
// @Success      201  {object} model.APIResponse
// @Router       /api/v1/auth/register [post]
func (h *UserHandler) Register(c *gin.Context) {
	var req model.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse("Input tidak valid", err.Error()))
		return
	}

	user, err := h.userService.Register(c.Request.Context(), req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "email sudah terdaftar" {
			statusCode = http.StatusConflict
		}
		c.JSON(statusCode, model.ErrorResponse(err.Error(), nil))
		return
	}

	c.JSON(http.StatusCreated, model.SuccessResponse("Registrasi berhasil", user))
}

// Login godoc
// @Summary      Login user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body body model.LoginRequest true "Kredensial login"
// @Success      200  {object} model.APIResponse
// @Router       /api/v1/auth/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req model.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse("Input tidak valid", err.Error()))
		return
	}

	result, err := h.userService.Login(c.Request.Context(), req)
	if err != nil {
		statusCode := http.StatusUnauthorized 
		c.JSON(statusCode, model.ErrorResponse(err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse("Login berhasil", result))
}

// GetAll godoc
// @Summary      Ambil semua user (dengan pagination)
// @Tags         users
// @Security     BearerAuth
// @Produce      json
// @Param        page   query int    false "Halaman (default: 1)"
// @Param        limit  query int    false "Jumlah per halaman (default: 10)"
// @Param        search query string false "Kata kunci pencarian"
// @Success      200 {object} model.APIResponse
// @Router       /api/v1/users [get]
func (h *UserHandler) GetAll(c *gin.Context) {
	var query model.PaginationQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse("Parameter tidak valid", err.Error()))
		return
	}

	users, meta, err := h.userService.GetAll(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse("Gagal ambil data", err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.SuccessListResponse("Berhasil", users, meta))
}

// GetByID godoc
// @Summary      Ambil user berdasarkan ID
// @Tags         users
// @Security     BearerAuth
// @Produce      json
// @Param        id path int true "User ID"
// @Success      200 {object} model.APIResponse
// @Router       /api/v1/users/{id} [get]
func (h *UserHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse("ID tidak valid", nil))
		return
	}

	user, err := h.userService.GetByID(c.Request.Context(), id)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "user tidak ditemukan" {
			statusCode = http.StatusNotFound 
		}
		c.JSON(statusCode, model.ErrorResponse(err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse("Berhasil", user))
}

// Update godoc
// @Summary      Update data user
// @Tags         users
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id   path int                      true "User ID"
// @Param        body body model.UpdateUserRequest  true "Data yang diupdate"
// @Success      200  {object} model.APIResponse
// @Router       /api/v1/users/{id} [put]
func (h *UserHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse("ID tidak valid", nil))
		return
	}

	var req model.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse("Input tidak valid", err.Error()))
		return
	}

	updated, err := h.userService.Update(c.Request.Context(), id, req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "user tidak ditemukan" {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, model.ErrorResponse(err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse("User berhasil diupdate", updated))
}

// Delete godoc
// @Summary      Hapus user
// @Tags         users
// @Security     BearerAuth
// @Produce      json
// @Param        id path int true "User ID"
// @Success      200 {object} model.APIResponse
// @Router       /api/v1/users/{id} [delete]
func (h *UserHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse("ID tidak valid", nil))
		return
	}

	if err := h.userService.Delete(c.Request.Context(), id); err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "user tidak ditemukan" {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, model.ErrorResponse(err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse("User berhasil dihapus", nil))
}
