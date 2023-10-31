package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Vyom-Yadav/GitHub-Gist-Clone-Backend/initializers"
	"github.com/Vyom-Yadav/GitHub-Gist-Clone-Backend/models"
	"github.com/Vyom-Yadav/GitHub-Gist-Clone-Backend/utils"
	"github.com/gin-gonic/gin"
	"github.com/thanhpk/randstr"
	"go.uber.org/zap"
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
//	@Failure	502			{object}	map[string]string
//	@Router		/auth/register [post]
func (ac *AuthController) SignUpUser(ctx *gin.Context) {
	var payload *models.SignUpInput

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		zap.L().Error(err.Error())
		return
	}

	if payload.Password != payload.PasswordConfirm {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Passwords do not match"})
		return
	}

	hashedPassword, err := utils.HashItem(payload.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": err.Error()})
		zap.L().Error(err.Error())
		return
	}

	code := randstr.String(20)
	hashedVerificationCode, err := utils.HashItem(code)
	now := time.Now()
	newUser := models.User{
		Username:         payload.Username,
		FirstName:        payload.FirstName,
		LastName:         &payload.LastName,
		Email:            strings.ToLower(payload.Email),
		Password:         hashedPassword,
		Role:             "user",
		Provider:         "local",
		Verified:         false,
		CreatedAt:        now,
		UpdatedAt:        now,
		VerificationCode: &hashedVerificationCode,
		UserMetadata: models.UserMetadata{
			ProfilePicture: "default.png", // TODO: Change this to a default profile picture in some storage
		},
	}

	// Only non-zero associations will be upserted
	result := ac.DB.Create(&newUser)

	if result.Error != nil && strings.Contains(result.Error.Error(), "duplicate key value violates unique") {
		ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": "User with that email already exists"})
		return
	} else if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": "Something bad happened"})
		zap.L().Error(result.Error.Error())
		return
	}

	config, err := initializers.LoadConfig(os.Getenv("API_ENV_CONFIG_PATH"))
	if err != nil {
		txn := ac.DB.Delete(&models.User{}, "username = ?", newUser.Username)
		if txn.Error != nil {
			zap.L().Error("Error while deleting user", zap.Error(txn.Error))
		}
		utils.SomethingBadHappened(ctx)
		zap.L().Error(err.Error())
		return
	}

	verificationCodeErrorMessage := "error while sending verification code, please resend verification code"
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": verificationCodeErrorMessage})
		zap.L().Error(err.Error())
		return
	}

	emailData := utils.EmailData{
		URL:       config.ClientOrigin + "/verifyemail?verificationCode=" + code + "&username=" + newUser.Username,
		FirstName: newUser.FirstName,
		Subject:   "[GitHub-Gist-Clone] Verify your email address",
	}

	err = utils.SendEmail(newUser.Email, &emailData, "verificationCode.html")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": verificationCodeErrorMessage})
		zap.L().Error(err.Error())
		return
	}

	message := "We sent an email with a verification code to " + newUser.Email
	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "message": message})
}

//	@Summary	Get GitHub Client ID
//	@Tags		Authentication
//	@Produce	json
//	@Success	200	{object}	map[string]string
//	@Failure	500	{object}	map[string]string
//	@Router		/auth/github/clientid [get]
func (ac *AuthController) GetGitHubClientId(ctx *gin.Context) {
	config, err := initializers.LoadConfig(os.Getenv("API_ENV_CONFIG_PATH"))
	if err != nil {
		utils.SomethingBadHappened(ctx)
		zap.L().Error(err.Error())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "clientId": config.GitHubClientId})
}

//	@Summary	Register a new user through GitHub
//	@Tags		Authentication
//	@Produce	json
//	@Param		code		query		string	true	"code received from GitHub API after user authorizes application"
//	@Param		newUsername	query		string	false	"new username to be used for the user"
//	@Success	201			{object}	map[string]string
//	@Failure	400			{object}	map[string]string
//	@Failure	409			{object}	map[string]string
//	@Failure	500			{object}	map[string]string
//	@Failure	502			{object}	map[string]string
//	@Router		/auth/github/callback [get]
func (ac *AuthController) GitHubCallback(ctx *gin.Context) {
	// TODO: Make sure frontend sets scope=read:user user:email
	// TODO: Check the state parameter to prevent CSRF attacks
	code, exists := ctx.GetQuery("code")
	if !exists {
		utils.SomethingBadHappened(ctx)
		zap.L().Error("code parameter not found in query")
		return
	}
	newUsername, _ := ctx.GetQuery("newUsername")
	accessToken, err := getAccessToken(code)
	if err != nil {
		utils.SomethingBadHappened(ctx)
		zap.L().Error(err.Error())
		return
	}

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	gitHubUserUrl := "https://api.github.com/user"
	req, err := http.NewRequest("GET", gitHubUserUrl, nil)
	if err != nil {
		utils.SomethingBadHappened(ctx)
		zap.L().Error(err.Error())
		return
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := client.Do(req)
	if err != nil {
		utils.SomethingBadHappened(ctx)
		zap.L().Error(err.Error())
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		utils.SomethingBadHappened(ctx)
		zap.L().Error("GitHub API returned status code " + strconv.Itoa(resp.StatusCode))
		return
	}

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		utils.SomethingBadHappened(ctx)
		zap.L().Error(err.Error())
		return
	}

	// Print keys of result map
	for k := range result {
		zap.L().Info(k)
	}

	id, ok := result["id"]
	if !ok {
		utils.SomethingBadHappened(ctx)
		zap.L().Error("id not found in GitHub response")
		return
	}
	stringId := strconv.FormatFloat(id.(float64), 'f', 0, 64)
	var user models.User
	queryResult := ac.DB.First(&user, "github_user_id = ?", stringId)
	if errors.Is(queryResult.Error, gorm.ErrRecordNotFound) {
		signUpThroughGitHub(ctx, result, accessToken, newUsername, ac)
	} else if queryResult.Error != nil {
		utils.SomethingBadHappened(ctx)
		zap.L().Error(queryResult.Error.Error())
	} else if queryResult.Error == nil && queryResult.RowsAffected == 1 {
		// Changes in email and username through GitHub will not be reflected in our database
		// Reason: Too much redirection and user experience will be bad, so we will not allow it
		// Instead we use github_user_id to uniquely identify a user from GitHub
		githubGistAccessToken, encounteredError := registerLoginCookies(ctx, user.Username)
		if !encounteredError {
			ctx.JSON(http.StatusOK, gin.H{"status": "success", "access_token": githubGistAccessToken})
		}
	}
}

func signUpThroughGitHub(ctx *gin.Context, result map[string]interface{}, accessToken, newUsername string, ac *AuthController) {
	username, ok := result["login"].(string)
	if !ok {
		utils.SomethingBadHappened(ctx)
		zap.L().Error("Key login does not exists in result map")
		return
	}

	// Check if username already exists
	user := models.User{}
	queryResult := ac.DB.First(&user, "username = ?", username)

	if queryResult.Error == nil && queryResult.RowsAffected == 1 {
		if newUsername == "" {
			// User would be redirected to another page to choose a new username, GitHub OAuth is still used
			ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": "User with that username already exists, please choose a new username"})
			zap.L().Warn("User with that username already exists, redirect", zap.String("username", username))
			return
		}
	}

	if queryResult.Error != nil {
		if errors.Is(queryResult.Error, gorm.ErrRecordNotFound) {
			newUser, err := createUserFromGitHubOAuth(result, accessToken)
			if err != nil {
				utils.SomethingBadHappened(ctx)
				zap.L().Error(err.Error())
				return
			}

			if newUsername != "" {
				newUser.Username = newUsername
			}

			creationResult := ac.DB.Create(&newUser)

			if creationResult.Error != nil && strings.Contains(creationResult.Error.Error(), "duplicate key value violates unique") {
				ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": "User with that email already exists"})
				return
			} else if creationResult.Error != nil {
				ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": "Something bad happened"})
				zap.L().Error(creationResult.Error.Error())
				return
			}
			zap.L().Info("New user created", zap.String("username", newUser.Username))
			githubGistAccessToken, encounteredError := registerLoginCookies(ctx, newUser.Username)
			if !encounteredError {
				ctx.JSON(http.StatusOK, gin.H{"status": "success", "access_token": githubGistAccessToken})
			}
		} else {
			utils.SomethingBadHappened(ctx)
			zap.L().Error(queryResult.Error.Error())
		}
	}
}

func createUserFromGitHubOAuth(result map[string]interface{}, accessToken string) (*models.User, error) {
	username, ok := result["login"]
	if !ok {
		return nil, errors.New("key login does not exists in result map")
	}
	id, ok := result["id"]
	if !ok {
		return nil, errors.New("key id does not exists in result map")
	}

	email, err := getUsersGitHubEmail(accessToken)
	if err != nil {
		return nil, err
	}

	password, _ := utils.HashItem(randstr.String(32))
	now := time.Now()
	githubUserId := strconv.FormatFloat(id.(float64), 'f', 0, 64)
	user := models.User{
		Username:     username.(string),
		Email:        email,
		Password:     password,
		GithubUserId: &githubUserId,
		Role:         "user",
		Provider:     "local",
		Verified:     true,
		CreatedAt:    now,
		UpdatedAt:    now,
		UserMetadata: models.UserMetadata{
			ProfilePicture: "default.png", // TODO: Change this to a default profile picture in some storage
		},
	}

	var firstName = username.(string)
	var lastName *string = nil
	if avatarUrl, ok := result["avatar_url"]; ok {
		user.UserMetadata.ProfilePicture = avatarUrl.(string)
	}
	if name, ok := result["name"]; ok {
		if name.(string) != "" {
			splitN := strings.SplitN(name.(string), " ", 2)
			firstName = splitN[0]
			if splitN[1] != "" {
				lastName = &splitN[1]
			}
		}
	}
	if location, ok := result["location"]; ok {
		location := location.(string)
		user.UserMetadata.Location = &location
	}
	if bio, ok := result["bio"]; ok {
		bio := bio.(string)
		user.UserMetadata.Tagline = &bio
	}
	if twitterUsername, ok := result["twitter_username"]; ok {
		twitterUsername := twitterUsername.(string)
		user.UserMetadata.Twitter = &twitterUsername
	}
	user.FirstName = firstName
	user.LastName = lastName

	return &user, nil
}

func getUsersGitHubEmail(accessToken string) (string, error) {
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	gitHubEmailUrl := "https://api.github.com/user/emails"
	req, err := http.NewRequest("GET", gitHubEmailUrl, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("could not get email from GitHub")
	}
	var result []map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "", err
	}

	for _, emailMap := range result {
		if emailMap["primary"] == true {
			return emailMap["email"].(string), nil
		}
	}

	return "", errors.New("primary email not found")
}

func getAccessToken(code string) (string, error) {
	config, err := initializers.LoadConfig(os.Getenv("API_ENV_CONFIG_PATH"))
	if err != nil {
		return "", err
	}

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	gitHubAccessTokenUrl := "https://github.com/login/oauth/access_token"
	body := map[string]string{
		"client_id":     config.GitHubClientId,
		"client_secret": config.GitHubClientSecret,
		"code":          code,
	}
	postBody, _ := json.Marshal(body)
	req, err := http.NewRequest("POST", gitHubAccessTokenUrl, bytes.NewBuffer(postBody))
	if err != nil {
		return "", err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")


	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("could not get access token from GitHub")
	}

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "", err
	}

	zap.L().Info("GitHub access token result", zap.Any("result", result))

	accessToken, ok := result["access_token"].(string)
	if !ok {
		return "", errors.New("could not get access token")
	}

	return accessToken, nil
}

//	@Summary	Verify users email address
//	@Tags		Authentication
//	@Produce	json
//	@Param		username			query		string	true	"Username of the user to be verified"
//	@Param		verificationCode	query		string	true	"Verify the added user object to the database"
//	@Success	200					{object}	map[string]string
//	@Failure	400					{object}	map[string]string
//	@Failure	409					{object}	map[string]string
//	@Router		/auth/verifyemail [get]
func (ac *AuthController) VerifyEmail(ctx *gin.Context) {
	code, exists := ctx.GetQuery("verificationCode")
	if !exists {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Missing verification code"})
		return
	}
	username, exists := ctx.GetQuery("username")
	if !exists {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Missing username"})
		return
	}

	var updatedUser models.User
	result := ac.DB.First(&updatedUser, "username = ?", username)
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "User doesn't exists"})
		return
	}

	if updatedUser.Verified {
		ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": "User already verified"})
		return
	}

	equalityError := utils.VerifyItem(*updatedUser.VerificationCode, code)
	if equalityError != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid verification code"})
		return
	}

	updatedUser.VerificationCode = nil
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

	config, err := initializers.LoadConfig(os.Getenv("API_ENV_CONFIG_PATH"))
	if err != nil {
		utils.SomethingBadHappened(ctx)
		return
	}

	// Generate Verification Code
	code := randstr.String(20)

	hashedVerificationCode, err := utils.HashItem(code)

	// Update User in Database
	user.VerificationCode = &hashedVerificationCode
	ac.DB.Save(user)

	emailData := utils.EmailData{
		URL:       config.ClientOrigin + "/verifyemail?verificationCode=" + code + "&username=" + user.Username,
		FirstName: user.FirstName,
		Subject:   "[GitHub-Gist-Clone] Verify your email address",
	}

	err = utils.SendEmail(user.Email, &emailData, "verificationCode.html")
	if err != nil {
		utils.SomethingBadHappened(ctx)
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

	if err := utils.VerifyItem(user.Password, payload.Password); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid email or Password"})
		return
	}

	if !user.Verified {
		ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": "Please verify your email address"})
		return
	}

	accessToken, encounteredError := registerLoginCookies(ctx, user.Username)
	if !encounteredError {
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "access_token": accessToken})
	}
}

func registerLoginCookies(ctx *gin.Context, username string) (string, bool) {
	config, err := initializers.LoadConfig(os.Getenv("API_ENV_CONFIG_PATH"))
	if err != nil {
		utils.SomethingBadHappened(ctx)
		zap.L().Error("Error loading config", zap.Error(err))
		return "", true
	}

	// Generate Tokens
	accessToken, err := utils.CreateToken(config.AccessTokenExpiresIn, username, config.AccessTokenPrivateKey)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		zap.L().Error("Error creating access token", zap.Error(err))
		return "", true
	}

	refreshToken, err := utils.CreateToken(config.RefreshTokenExpiresIn, username, config.RefreshTokenPrivateKey)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		zap.L().Error("Error creating refresh token", zap.Error(err))
		return "", true
	}

	ctx.SetCookie("access_token", accessToken, config.AccessTokenMaxAge*60, "/", "localhost", false, true)
	ctx.SetCookie("refresh_token", refreshToken, config.RefreshTokenMaxAge*60, "/", "localhost", false, true)
	ctx.SetCookie("logged_in", "true", config.AccessTokenMaxAge*60, "/", "localhost", false, true)
	return accessToken, false
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

	config, err := initializers.LoadConfig(os.Getenv("API_ENV_CONFIG_PATH"))
	if err != nil {
		utils.SomethingBadHappened(ctx)
		return
	}

	sub, err := utils.ValidateToken(cookie, config.RefreshTokenPublicKey)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	var user models.User
	result := ac.DB.First(&user, "username = ?", sub.(string))
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
	ctx.SetCookie("logged_in", "true", config.AccessTokenMaxAge*60, "/", "localhost", false, true)

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
	ctx.SetCookie("logged_in", "", -1, "/", "localhost", false, true)

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
		ctx.JSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "Account not verified"})
		return
	}

	config, err := initializers.LoadConfig(os.Getenv("API_ENV_CONFIG_PATH"))
	if err != nil {
		utils.SomethingBadHappened(ctx)
		return
	}

	// Generate Verification Code
	resetToken := randstr.String(20)

	passwordResetToken, err := utils.HashItem(resetToken)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	user.PasswordResetToken = &passwordResetToken

	user.PasswordResetAt = time.Now().Add(time.Minute * 15)
	ac.DB.Save(&user)

	emailData := utils.EmailData{
		URL:       config.ClientOrigin + "/resetpassword?resetToken=" + resetToken + "&username=" + user.Username,
		FirstName: user.FirstName,
		Subject:   "Your password reset token (valid for 15 min)",
	}

	err = utils.SendEmail(user.Email, &emailData, "resetPassword.html")
	if err != nil {
		utils.SomethingBadHappened(ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": message})
}

//	@Summary	Reset password
//	@Tags		Authentication
//	@Accept		json
//	@Produce	json
//	@Param		ResetPasswordInput	body		models.ResetPasswordInput	true	"The input required to reset the password"
//	@Param		username			query		string						true	"The username of the user"
//	@Param		resetToken			query		string						true	"The token required to reset the password"
//	@Success	200					{object}	map[string]string
//	@Failure	400					{object}	map[string]string
//	@Router		/auth/resetpassword [patch]
func (ac *AuthController) ResetPassword(ctx *gin.Context) {
	var payload *models.ResetPasswordInput
	resetToken, exists := ctx.GetQuery("resetToken")
	if !exists {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Reset token not provided"})
		return
	}

	username, exists := ctx.GetQuery("username")
	if !exists {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Username not provided"})
		return
	}

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	if payload.Password != payload.PasswordConfirm {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Passwords do not match"})
		return
	}

	hashedPassword, err := utils.HashItem(payload.Password)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	var updatedUser models.User
	result := ac.DB.First(&updatedUser, "username = ?", username)
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "The user does not exists"})
		return
	}

	if !updatedUser.PasswordResetAt.After(time.Now()) {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Reset token expired"})
		return
	}

	validResetTokenError := utils.VerifyItem(*updatedUser.PasswordResetToken, resetToken)
	if validResetTokenError != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid reset token"})
		return
	}

	updatedUser.Password = hashedPassword
	updatedUser.PasswordResetToken = nil
	ac.DB.Save(&updatedUser)

	// Invalidate user session
	ctx.SetCookie("access_token", "", -1, "/", "localhost", false, true)
	ctx.SetCookie("refresh_token", "", -1, "/", "localhost", false, true)
	ctx.SetCookie("logged_in", "", -1, "/", "localhost", false, true)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "Password data updated successfully"})
}

//	@Summary	Check if username is available
//	@Tags		Authentication
//	@Produce	json
//	@Param		username	path		string	true	"The username to check availability"
//	@Success	200			{object}	map[string]string
//	@Router		/auth/usernameavailable/{username} [get]
func (ac *AuthController) CheckUsernameAvailability(ctx *gin.Context) {
	username := ctx.Params.ByName("username")

	var user models.User
	result := ac.DB.First(&user, "username = ?", username)
	if result.Error != nil {
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "Username is available"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "fail", "message": "Username is not available"})
}
