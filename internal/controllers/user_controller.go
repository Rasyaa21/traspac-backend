package controllers

import (
	dto "gin-backend-app/internal/dto/common"
	"gin-backend-app/internal/dto/request"
	"gin-backend-app/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController struct {
    UserService *services.UserService
}

func NewUserController(userService *services.UserService) *UserController {
    return &UserController{
        UserService: userService,
    }
}

// LoginUser godoc
// @Summary User login
// @Description Authenticate user with email and password
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body request.LoginUserRequest true "Login credentials"
// @Success 201 {object} dto.Response "User created successfully"
// @Failure 400 {object} dto.ErrorResponse "Invalid request data"
// @Failure 401 {object} dto.ErrorResponse "Invalid credentials"
// @Router /auth/login [post]
func (uc *UserController) LoginUser(c *gin.Context) {
	var req request.LoginUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		dto.SendError(c, http.StatusBadRequest, "Invalid Request Data")
		return
	}

    loginResponse, err := uc.UserService.LoginUser(req)
    if err != nil {
        dto.SendError(c, http.StatusUnauthorized, err.Error())
        return
    }

    dto.SendResponse(c, http.StatusOK, loginResponse, "Login successful")
}

// RegisterUser godoc
// @Summary Register new user  
// @Description Create a new user account
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body request.CreateUserRequest true "User registration data"
// @Success 201 {object} dto.Response "User created successfully"
// @Failure 409 {object} dto.ErrorResponse "Invalid request data or user already exists"
// @Router /auth/register [post]
func (uc *UserController) RegisterUser(c *gin.Context) {
	var req request.CreateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		dto.SendError(c, http.StatusBadRequest, "Invalid Request Data")
		return
	}

    registerResponse, err := uc.UserService.CreateUser(req)
    if err != nil {
        dto.SendError(c, http.StatusBadRequest, err.Error())
        return
    }

    dto.SendResponse(c, http.StatusCreated, registerResponse, "User created successfully")
}