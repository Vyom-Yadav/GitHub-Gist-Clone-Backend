package controllers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Vyom-Yadav/GitHub-Gist-Clone-Backend/initializers"
	"github.com/Vyom-Yadav/GitHub-Gist-Clone-Backend/models"
	"github.com/Vyom-Yadav/GitHub-Gist-Clone-Backend/utils"
	"github.com/gin-gonic/gin"
	"github.com/thanhpk/randstr"
	"gorm.io/gorm"
)

type AuthController struct {
	DB *gorm.DB
}

func NetAuthController(DB *gorm.DB) AuthController {
	return AuthController{
		DB: DB,
	}
}

//	@Summary	Register a new user
//	@Tags		Authentication
//	@Accept		json
//	@Produce	json
//	@Param		SignUpInput	body		models.SignUpInput	true	"User object that needs to be added to the system"
//	@Success	201			{object}	map[string]string
//	@Failure	400			{object}	map[string]string
//	@Failure	409			{object}	map[string]string
//	@Failure	500			{object}	map[string]string
//	@Router		/auth/register [post]
func (ac *AuthController) SignUpUser(ctx *gin.Context) {
	var payload *models.SignUpInput

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	if payload.Password != payload.PasswordConfirm {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Passwords do not match"})
		return
	}

	hashedPassword, err := utils.HashPassword(payload.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	now := time.Now()
	newUser := models.User{
		Username:     payload.Username,
		FirstName:    payload.FirstName,
		LastName:     payload.LastName,
		Email:        strings.ToLower(payload.Email),
		Password:     hashedPassword,
		Role:         "user",
		Provider:     "local",
		Verified:     false,
		CreatedAt:    now,
		UpdatedAt:    now,
		UserMetadata: models.UserMetadata{},
		Gists:        nil,
	}

	result := ac.DB.Create(&newUser)

	if result.Error != nil && strings.Contains(result.Error.Error(), "duplicate key value violates unique") {
		ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": "User with that email already exists"})
		return
	} else if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Something bad happened"})
		return
	}

	config, err := initializers.LoadConfig("/app/env")
	if err != nil {
		ac.DB.Delete(&newUser)
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": "Something bad happened"})
		return
	}

	// Generate Verification Code
	code := randstr.String(20)

	verificationCode := utils.Encode(code)

	// Update User in Database
	newUser.VerificationCode = verificationCode
	ac.DB.Save(newUser)

	emailData := utils.EmailData{
		URL:       config.ClientOrigin + "/verifyemail/" + code,
		FirstName: newUser.FirstName,
		Subject:   "Your account verification code",
	}

	err = utils.SendEmail(&newUser, &emailData, "verificationCode.html")
	if err != nil {
		ac.DB.Delete(&newUser)
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": "Something bad happened"})
		return
	}

	message := "We sent an email with a verification code to " + newUser.Email
	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "message": message})
}

//	@Summary	Verify users email address
//	@Tags		Authentication
//	@Produce	json
//	@Param		verificationCode	path		string	true	"Verify the added user object to the database"
//	@Success	200					{object}	map[string]string
//	@Failure	400					{object}	map[string]string
//	@Failure	409					{object}	map[string]string
//	@Router		/auth/verifyemail/{verificationCode} [get]
func (ac *AuthController) VerifyEmail(ctx *gin.Context) {
	code := ctx.Params.ByName("verificationCode")
	verificationCode := utils.Encode(code)

	var updatedUser models.User
	result := ac.DB.First(&updatedUser, "verification_code = ?", verificationCode)
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid verification code or user doesn't exists"})
		return
	}

	if updatedUser.Verified {
		ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": "User already verified"})
		return
	}

	updatedUser.VerificationCode = ""
	updatedUser.Verified = true
	ac.DB.Save(&updatedUser)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "Email verified successfully"})
}

//	@Summary	Resend verification email
//	@Tags		Authentication
//	@Accept		json
//	@Produce	json
//	@Param		ResendVerificationEmailInput	body		models.ResendVerificationEmailInput	true	"Resend verification email to the user with the given email"
//	@Success	200								{object}	map[string]string
//	@Failure	400								{object}	map[string]string
//	@Failure	409								{object}	map[string]string
//	@Failure	500								{object}	map[string]string
//	@Router		/auth/resendverificationemail [post]
func (ac *AuthController) ResendVerificationEmail(ctx *gin.Context) {
	var payload *models.ResendVerificationEmailInput

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	var user models.User
	result := ac.DB.First(&user, "email = ?", strings.ToLower(payload.Email))
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "User with that email doesn't exists"})
		return
	}

	if user.Verified {
		ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": "User already verified"})
		return
	}

	config, err := initializers.LoadConfig("/app/env")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": "Something bad happened"})
		return
	}

	// Generate Verification Code
	code := randstr.String(20)

	verificationCode := utils.Encode(code)

	// Update User in Database
	user.VerificationCode = verificationCode
	ac.DB.Save(user)

	emailData := utils.EmailData{
		URL:       config.ClientOrigin + "/verifyemail/" + code,
		FirstName: user.FirstName,
		Subject:   "Your account verification code",
	}

	err = utils.SendEmail(&user, &emailData, "verificationCode.html")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": "Something bad happened"})
		return
	}

	message := "We sent an email with a verification code to " + user.Email
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": message})
}

//	@Summary	Sign in a user
//	@Tags		Authentication
//	@Accept		json
//	@Produce	json
//	@Param		SignInInput	body		models.SignInInput	true	"Sign in a user"
//	@Success	200			{object}	map[string]string
//	@Failure	400			{object}	map[string]string
//	@Failure	409			{object}	map[string]string
//	@Failure	500			{object}	map[string]string
//	@Router		/auth/login [post]
func (ac *AuthController) SignInUser(ctx *gin.Context) {
	var payload *models.SignInInput

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	var user models.User

	result := ac.DB.First(&user, "email = ?", strings.ToLower(payload.Email))
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid email or Password"})
		return
	}

	if err := utils.VerifyPassword(user.Password, payload.Password); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid email or Password"})
		return
	}

	config, err := initializers.LoadConfig("/app/env")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": "Something bad happened"})
		return
	}

	// Generate Tokens
	accessToken, err := utils.CreateToken(config.AccessTokenExpiresIn, user.Username, config.AccessTokenPrivateKey)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	refreshToken, err := utils.CreateToken(config.RefreshTokenExpiresIn, user.Username, config.RefreshTokenPrivateKey)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.SetCookie("access_token", accessToken, config.AccessTokenMaxAge*60, "/", "localhost", false, true)
	ctx.SetCookie("refresh_token", refreshToken, config.RefreshTokenMaxAge*60, "/", "localhost", false, true)

	// TODO: See why this has httpOnly as 'false'
	ctx.SetCookie("logged_in", "true", config.AccessTokenMaxAge*60, "/", "localhost", false, false)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "access_token": accessToken})

}

//	@Summary	Refresh access token with refresh token
//	@Tags		Authentication
//	@Produce	json
//	@Success	200	{object}	map[string]string
//	@Failure	403	{object}	map[string]string
//	@Failure	500	{object}	map[string]string
//	@Router		/auth/refresh [get]
func (ac *AuthController) RefreshAccessToken(ctx *gin.Context) {
	message := "could not refresh access token"

	cookie, err := ctx.Cookie("refresh_token")

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": message})
		return
	}

	config, err := initializers.LoadConfig("/app/env")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": "Something bad happened"})
		return
	}

	sub, err := utils.ValidateToken(cookie, config.RefreshTokenPublicKey)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	var user models.User
	result := ac.DB.First(&user, "username = ?", fmt.Sprint(sub))
	if result.Error != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": "the user belonging to this token no logger exists"})
		return
	}

	accessToken, err := utils.CreateToken(config.AccessTokenExpiresIn, user.Username, config.AccessTokenPrivateKey)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.SetCookie("access_token", accessToken, config.AccessTokenMaxAge*60, "/", "localhost", false, true)
	ctx.SetCookie("logged_in", "true", config.AccessTokenMaxAge*60, "/", "localhost", false, false)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "access_token": accessToken})
}

//	@Summary	Log out a user
//	@Tags		Authentication
//	@Produce	json
//	@Success	200	{object}	map[string]string
//	@Router		/auth/logout [get]
func (ac *AuthController) LogoutUser(ctx *gin.Context) {
	ctx.SetCookie("access_token", "", -1, "/", "localhost", false, true)
	ctx.SetCookie("refresh_token", "", -1, "/", "localhost", false, true)
	ctx.SetCookie("logged_in", "", -1, "/", "localhost", false, false)

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}

//	@Summary	Send reset code for password reset
//	@Tags		Authentication
//	@Accept		json
//	@Produce	json
//	@Param		ForgotPasswordInput	body		models.ForgotPasswordInput	true	"The Input for sending password reset code"
//	@Success	200					{object}	map[string]string
//	@Failure	400					{object}	map[string]string
//	@Failure	401					{object}	map[string]string
//	@Failure	500					{object}	map[string]string
//	@Router		/auth/forgotpassword [post]
func (ac *AuthController) ForgotPassword(ctx *gin.Context) {
	var payload *models.ForgotPasswordInput

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	message := "You will receive a reset email if user with that email exists"

	var user models.User
	result := ac.DB.First(&user, "email = ?", strings.ToLower(payload.Email))
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid email or Password"})
		return
	}

	if !user.Verified {
		ctx.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Account not verified"})
		return
	}

	config, err := initializers.LoadConfig("/app/env")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": "Something bad happened"})
		return
	}

	// Generate Verification Code
	resetToken := randstr.String(20)

	passwordResetToken := utils.Encode(resetToken)
	user.PasswordResetToken = passwordResetToken

	user.PasswordResetAt = time.Now().Add(time.Minute * 15)
	ac.DB.Save(&user)

	emailData := utils.EmailData{
		URL:       config.ClientOrigin + "/resetpassword/" + resetToken,
		FirstName: user.FirstName,
		Subject:   "Your password reset token (valid for 15 min)",
	}

	err = utils.SendEmail(&user, &emailData, "resetPassword.html")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": "Something bad happened"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": message})
}

//	@Summary	Reset password
//	@Tags		Authentication
//	@Accept		json
//	@Produce	json
//	@Param		ResetPasswordInput	body		models.ResetPasswordInput	true	"The input required to reset the password"
//	@Param		resetToken			path		string						true	"The token required to reset the password"
//	@Success	200					{object}	map[string]string
//	@Failure	400					{object}	map[string]string
//	@Router		/auth/resetpassword/{resetToken} [patch]
func (ac *AuthController) ResetPassword(ctx *gin.Context) {
	var payload *models.ResetPasswordInput
	resetToken := ctx.Params.ByName("resetToken")

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	if payload.Password != payload.PasswordConfirm {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Passwords do not match"})
		return
	}

	hashedPassword, err := utils.HashPassword(payload.Password)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": err.Error()})
		return
	}
	passwordResetToken := utils.Encode(resetToken)

	var updatedUser models.User
	result := ac.DB.First(&updatedUser, "password_reset_token = ? AND password_reset_at > ?", passwordResetToken, time.Now())
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "The reset token is invalid or has expired"})
		return
	}

	updatedUser.Password = hashedPassword
	updatedUser.PasswordResetToken = ""
	ac.DB.Save(&updatedUser)

	// Invalidate user session
	ctx.SetCookie("access_token", "", -1, "/", "localhost", false, true)
	ctx.SetCookie("refresh_token", "", -1, "/", "localhost", false, true)
	ctx.SetCookie("logged_in", "", -1, "/", "localhost", false, false)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "Password data updated successfully"})
}

//	@Summary	Check if username is available
//	@Tags		Authentication
//	@Produce	json
//	@Param		username	path		string	true	"The username to check availability"
//	@Success	200			{object}	map[string]string
//	@Failure	400			{object}	map[string]string
//	@Router		/auth/usernameavailable/{username} [get]
func (ac *AuthController) CheckUsernameAvailability(ctx *gin.Context) {
	username := ctx.Params.ByName("username")

	var user models.User
	result := ac.DB.First(&user, "username = ?", username)
	if result.Error != nil {
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "Username is available"})
		return
	}

	ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Username is not available"})
}
