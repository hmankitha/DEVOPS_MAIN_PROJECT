package handlers

import (
	"net/http"
	"strconv"
	"user-management/models"
	"user-management/services"
	"user-management/utils"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService services.UserService
}

func NewUserHandler(userService services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) GetCurrentUser(c *gin.Context) {
	userID := c.GetString("user_id")

	user, err := h.userService.GetByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, utils.ErrorResponse("User not found"))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse(user, "User retrieved successfully"))
}

func (h *UserHandler) GetUserByID(c *gin.Context) {
	id := c.Param("id")

	user, err := h.userService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, utils.ErrorResponse("User not found"))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse(user, "User retrieved successfully"))
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := c.GetString("user_id")

	var req models.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(err.Error()))
		return
	}

	user, err := h.userService.UpdateProfile(userID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse(user, "Profile updated successfully"))
}

func (h *UserHandler) ChangePassword(c *gin.Context) {
	userID := c.GetString("user_id")

	var req models.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(err.Error()))
		return
	}

	if err := h.userService.ChangePassword(userID, &req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse(nil, "Password changed successfully"))
}

func (h *UserHandler) DeleteAccount(c *gin.Context) {
	userID := c.GetString("user_id")

	if err := h.userService.DeleteAccount(userID); err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to delete account"))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse(nil, "Account deleted successfully"))
}

func (h *UserHandler) ListUsers(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	offset := (page - 1) * limit

	users, err := h.userService.ListUsers(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to fetch users"))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse(users, "Users retrieved successfully"))
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")

	if err := h.userService.DeleteAccount(id); err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to delete user"))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse(nil, "User deleted successfully"))
}

func (h *UserHandler) UpdateUserRole(c *gin.Context) {
	id := c.Param("id")

	var req models.UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(err.Error()))
		return
	}

	if err := h.userService.UpdateUserRole(id, req.Role); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse(nil, "User role updated successfully"))
}

func (h *UserHandler) GetStats(c *gin.Context) {
	stats, err := h.userService.GetStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to fetch stats"))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse(stats, "Stats retrieved successfully"))
}
