package handlers

import (
	"net/http"

	"mini-payment-system/internal/delivery/http/dto"
	useruc "mini-payment-system/internal/usecase/user"
	"mini-payment-system/pkg/response"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service *useruc.Service
}

func NewUserHandler(service *useruc.Service) *UserHandler {
	return &UserHandler{service: service}
}

// Create godoc
// @Summary Create user
// @Description Create a new user record.
// @Tags users
// @Accept json
// @Produce json
// @Param request body dto.CreateUserRequest true "Create user payload"
// @Success 201 {object} response.UserResponse
// @Failure 400 {object} response.ErrorBody
// @Failure 409 {object} response.ErrorBody
// @Failure 500 {object} response.ErrorBody
// @Router /users [post]
func (h *UserHandler) Create(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	user, err := h.service.Create(c.Request.Context(), req.Name, req.Email)
	if err != nil {
		apiErr := mapError(err)
		response.Error(c, apiErr.Status, apiErr.Code, apiErr.Message)
		return
	}
	response.JSON(c, http.StatusCreated, user)
}

// List godoc
// @Summary List users
// @Description Return all users ordered by latest created first.
// @Tags users
// @Produce json
// @Success 200 {object} response.UserListResponse
// @Failure 500 {object} response.ErrorBody
// @Router /users [get]
func (h *UserHandler) List(c *gin.Context) {
	users, err := h.service.List(c.Request.Context())
	if err != nil {
		apiErr := mapError(err)
		response.Error(c, apiErr.Status, apiErr.Code, apiErr.Message)
		return
	}
	response.List(c, http.StatusOK, users, len(users))
}

// GetByID godoc
// @Summary Get user by ID
// @Tags users
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} response.UserResponse
// @Failure 404 {object} response.ErrorBody
// @Failure 500 {object} response.ErrorBody
// @Router /users/{id} [get]
func (h *UserHandler) GetByID(c *gin.Context) {
	user, err := h.service.GetByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		apiErr := mapError(err)
		response.Error(c, apiErr.Status, apiErr.Code, apiErr.Message)
		return
	}
	response.JSON(c, http.StatusOK, user)
}

// Update godoc
// @Summary Update user
// @Description Update name and/or email of a user.
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param request body dto.UpdateUserRequest true "Update user payload"
// @Success 200 {object} response.UserResponse
// @Failure 400 {object} response.ErrorBody
// @Failure 404 {object} response.ErrorBody
// @Failure 409 {object} response.ErrorBody
// @Failure 500 {object} response.ErrorBody
// @Router /users/{id} [put]
func (h *UserHandler) Update(c *gin.Context) {
	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	user, err := h.service.Update(c.Request.Context(), c.Param("id"), req.Name, req.Email)
	if err != nil {
		apiErr := mapError(err)
		response.Error(c, apiErr.Status, apiErr.Code, apiErr.Message)
		return
	}
	response.JSON(c, http.StatusOK, user)
}

// Delete godoc
// @Summary Delete user
// @Description Delete a user by ID.
// @Tags users
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} response.MessageResponse
// @Failure 404 {object} response.ErrorBody
// @Failure 500 {object} response.ErrorBody
// @Router /users/{id} [delete]
func (h *UserHandler) Delete(c *gin.Context) {
	if err := h.service.Delete(c.Request.Context(), c.Param("id")); err != nil {
		apiErr := mapError(err)
		response.Error(c, apiErr.Status, apiErr.Code, apiErr.Message)
		return
	}
	response.Message(c, http.StatusOK, "user deleted successfully")
}
