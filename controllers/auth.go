package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/snykk/beego-presence-api/constants"
	"github.com/snykk/beego-presence-api/dto"
	"github.com/snykk/beego-presence-api/helpers"
	"github.com/snykk/beego-presence-api/models"
	"golang.org/x/crypto/bcrypt"

	beego "github.com/beego/beego/v2/server/web"
)

type AuthController struct {
	beego.Controller
}

func (c *AuthController) URLMapping() {
	c.Mapping("Register", c.Register) // Maps POST /register to Register method
	c.Mapping("Login", c.Login)       // Maps POST /login to Login method
}

// @Title Register User
// @Description Register a new user with the provided data
// @Accept  json
// @Produce  json
// @Param registerRequest body dto.RegisterRequest true "Registration Data"
// @Success 201 {object} dto.UserResponse "User registered successfully"
// @Failure 400 Invalid input data
// @Failure 500 Failed to register user
// @router /regis [post]
func (c *AuthController) Register() {
	var req dto.RegisterRequest
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusBadRequest, "Invalid input", err)
		return
	}

	// Validate the request payload.
	if errorsMap, err := helpers.ValidatePayloads(req); err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusBadRequest, constants.ErrValidationMessage, errorsMap)
		return
	}

	// Fetch the department by ID.
	department, err := models.GetDepartmentById(req.DepartmentId, false, false)
	if err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusBadRequest, "Invalid department ID", err)
		return
	}

	// Hash the password before saving.
	hashedPassword, err := helpers.HashPassword(req.Password)
	if err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusInternalServerError, "Failed to hash password", err)
		return
	}

	// Convert the request data to a user model.
	user := req.ToUserModel(department)
	user.Password = hashedPassword

	// Create the user in the database.
	if err := models.CreateUser(user); err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusInternalServerError, "Failed to register user", err)
		return
	}

	// Return the registered user.
	helpers.SuccessResponse(c.Ctx.ResponseWriter, http.StatusCreated, "User registered successfully", map[string]interface{}{"users": dto.FromUserModelToRegisterResponse(user)})
}

// @Title User Login
// @Description Authenticate user and generate a JWT token
// @Accept  json
// @Produce  json
// @Param loginRequest body dto.LoginRequest true "Login Data"
// @Success 200 {object} dto.LoginResponse "Login successful"
// @Failure 400 Invalid credentials
// @Failure 500 Failed to generate token
// @router /login [post]
func (c *AuthController) Login() {
	var req dto.LoginRequest
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusBadRequest, "Invalid input", err)
		return
	}

	// Validate the login request payload.
	if errorsMap, err := helpers.ValidatePayloads(req); err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusBadRequest, constants.ErrValidationMessage, errorsMap)
		return
	}

	// Fetch the user by email.
	user, err := models.GetUserByEmail(req.Email)
	if err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusUnauthorized, "Invalid credentials", nil)
		return
	}

	// Verify the password.
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusUnauthorized, "Invalid credentials", nil)
		return
	}

	// Generate a JWT token for the user.
	token, err := helpers.GenerateJWT(user.Id, user.Email, user.Role)
	if err != nil {
		helpers.ErrorResponse(c.Ctx.ResponseWriter, http.StatusInternalServerError, "Failed to generate token", err)
		return
	}

	// Return the token.
	helpers.SuccessResponse(c.Ctx.ResponseWriter, http.StatusOK, "Login successful", dto.LoginResponse{Token: token})
}
