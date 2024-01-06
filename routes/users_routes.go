package routes

import (
	"API_Rest/db"
	"API_Rest/models"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

// Create new user
func CreateUserHandler(c echo.Context) error {
	var user models.User
	//Binding data
	if err := c.Bind(&user); err != nil {
		return c.String(http.StatusUnprocessableEntity, err.Error())
	}

	//Validating data
	if err := c.Validate(user); err != nil {
		return c.String(http.StatusUnprocessableEntity, "Invalid data")
	}
	//Hasing password
	if err := user.HashSaltPassword(); err != nil {
		return c.String(http.StatusInternalServerError, "Error hasing the password")
	}

	///////////////////////////////////////////////////////////////////////////////////////////////////
	//Find a way to change the message and error code according to the error reported by the database//
	///////////////////////////////////////////////////////////////////////////////////////////////////
	//Creating user
	if err := db.DataBase().Create(&user).Error; err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.String(http.StatusCreated, "User created")
}

// Gets user all Users
func GetUsersHandler(c echo.Context) error {
	var users []models.User
	if err := db.DataBase().Select([]string{"email", "name", "funds", "role"}).Find(&users).Error; err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}
	return c.JSON(http.StatusOK, users)
}

// Get user using email like path parameter
func GetUserHandler(c echo.Context) error {
	var params struct {
		Email string `validate:"required,email"`
	}
	var user models.User
	params.Email = c.Param("email")
	if err := c.Validate(params); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, user)
	}

	if err := db.DataBase().Select([]string{"email", "name", "funds", "role"}).First(&user, "email = ?", params.Email).Error; err != nil {
		return c.JSON(http.StatusNotFound, user)
	}

	return c.JSON(http.StatusOK, user)
}

// Login user using request body parameter(email and password)
func LoginUserHandler(c echo.Context) error {
	var params struct {
		Email    string `validate:"required,email"`
		Password string `validate:"required,alphanum"`
	}

	if err := c.Bind(&params); err != nil {
		return c.String(http.StatusUnprocessableEntity, err.Error())
	}

	if err := c.Validate(params); err != nil {
		return c.String(http.StatusUnprocessableEntity, err.Error())
	}

	var user models.User
	if err := db.DataBase().First(&user, "email = ?", params.Email).Error; err != nil {
		return c.String(http.StatusUnauthorized, "Invalid login parameters")
	}

	if user.ValidatePassword(user.Password, params.Password) {
		return c.String(http.StatusUnauthorized, "Invalid login parameters")
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":    user.UserId,
		"user_role":  user.Role,
		"user_funds": user.Funds,
	})

	secretKey := []byte(os.Getenv("SECRET_KEY"))
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return c.String(http.StatusNotExtended, err.Error())
	}

	return c.JSON(http.StatusOK, tokenString)
}

// Delete users by id == you can lose all users
// Validar token
// Delete issuer ==  delete invoices == problems
func DeleteUserHandler(c echo.Context) error {
	var user models.User
	param := c.Param("id")
	if err := db.DataBase().First(&user, "user_id = ?", param).Error; err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}
	// Marc this user like a delete in DB
	if err := db.DataBase().Delete(&user).Error; err != nil {
		//Check this status
		return c.String(http.StatusNotModified, err.Error())
	}

	return c.String(http.StatusOK, "User deleted")
	// Real delete in DB
	// db.DB.Unscoped().Delete(&user)
}
