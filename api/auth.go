package api

import (
	"errors"
	"net/http"
	"os"
	"time"

	"cozeva.com/vault/interfaces"
	"cozeva.com/vault/pkg/security"
	"cozeva.com/vault/pkg/user"
	"cozeva.com/vault/pkg/validation"
	sqlquery "cozeva.com/vault/sql"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mattn/go-sqlite3"
)

// UserSignIn handles sign in logic for users
func UserSignup(c *gin.Context) {
	// Unmarshal JSON payload from request body
	var newUserSignupPayload interfaces.UserSignupPayload
	if err := c.BindJSON(&newUserSignupPayload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "Missing required parameter",
		})
		return
	}

	// Validate user inputs
	if newUserSignupPayload.FirstName == "" || newUserSignupPayload.LastName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "First name and last name must be a valid string",
		})
		return
	}

	// Email validation
	if !validation.IsValidEmail(newUserSignupPayload.Email) {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "Invalid email format",
		})
		return
	}

	// Password validation
	if err := validation.ValidatePassword(newUserSignupPayload.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})
		return
	}

	// Hash the password using Argon2
	hash, err := security.GenerateHash(newUserSignupPayload.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "fail",
			"message": "failed to hash password",
		})
		return
	}

	// Connect to the database
	db, err := sqlquery.GetSqlInstance()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "Failed to connect to db.",
		})
		return
	}
	defer db.Close()

	// Create a new user in the database
	_, err = db.Exec(sqlquery.InsertUserQuery, newUserSignupPayload.FirstName, newUserSignupPayload.LastName, newUserSignupPayload.Email, hash, newUserSignupPayload.ProfilePicture)
	if err != nil {
		// Check if it is a SQLite UNIQUE constraint error
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) && sqliteErr.Code == sqlite3.ErrConstraint {
			if sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
				c.JSON(http.StatusBadRequest, gin.H{
					"status":  "fail",
					"message": "email already exists",
				})
				return
			}
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "fail",
			"message": "Failed to create new user.",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Created new user with email address " + newUserSignupPayload.Email,
	})
}

func UserSignin(c *gin.Context) {
	var newUserSignInPayload interfaces.UserSignInPayload
	if err := c.BindJSON(&newUserSignInPayload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "Missing required credentials",
		})
		return
	}

	if newUserSignInPayload.Email == "" || newUserSignInPayload.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "Email and password must be a valid string",
		})
		return
	}

	db, err := sqlquery.GetSqlInstance()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "Failed to connect to db.",
		})
		return
	}
	defer db.Close()

	// Create a new user in the database
	row := db.QueryRow(sqlquery.CheckUserCreDentialsQuery, newUserSignInPayload.Email)
	var creds interfaces.UserCredentials
	dbErr := row.Scan(&creds.ID, &creds.PasswordHash)
	if dbErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "Please Enter Valid Credentials",
		})
		return
	}

	passwordMatch, _ := security.CompareHash(newUserSignInPayload.Password, creds.PasswordHash)
	if !passwordMatch {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "Password is incorrect",
		})
		return
	}

	sessionId := uuid.NewString()
	refreshToken := uuid.NewString()
	refreshTokenHash, err := security.GenerateHash(refreshToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "fail",
			"message": "failed to hash password",
		})
		return
	}

	ttlStr := os.Getenv("core.refresh.token.ttl")
	refreshTTL, err := time.ParseDuration(ttlStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "fail",
			"message": "failed to Parse Duration",
		})
		return
	}

	refreshTokenExpiry := time.Now().Add(refreshTTL).UTC().Format("2006-01-02 15:04:05")
	_, err = db.Exec(sqlquery.InsertUserSessionsQuery, creds.ID, sessionId, refreshTokenHash, refreshTokenExpiry)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "fail",
			"message": "Failed to insert user session in db",
		})
		return
	}
	jwtTTL := os.Getenv("core.jwt.token.ttl")
	jwtTokenTTL, err := time.ParseDuration(jwtTTL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "fail",
			"message": "failed to Parse Duration",
		})
		return
	}
	jwtToken, err := security.GenerateJWT(creds.ID, sessionId, jwtTokenTTL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "fail",
			"message": "Failed to generate jwt token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":        "success",
		"message":       "Logged in successfully",
		"access_token":  jwtToken,
		"refresh_token": refreshToken,
	})

}

func UserLogout(c *gin.Context) {
	currentUser, err := user.GetCurrentUser(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "Failed to parse current user info.",
		})
		return
	}

	db, err := sqlquery.GetSqlInstance()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "fail",
			"message": "Failed to connect to db.",
		})
		return
	}
	defer db.Close()

	_, err = db.Exec(sqlquery.DeleteSessionIDFromUserSessionsTable, currentUser.Uid, currentUser.SessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "fail",
			"message": "Failed to delete user session from db.",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Logged out successfully",
	})

}

func RefreshTokenHandler(c *gin.Context) {
	var newRefreshTokenHandlerPayload interfaces.RefreshTokenHandlerPayload
	if err := c.BindJSON(&newRefreshTokenHandlerPayload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "Missing required credentials",
		})
		return
	}

	currentUser, err := user.GetCurrentUser(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "Failed to parse current user info.",
		})
		return
	}

	db, err := sqlquery.GetSqlInstance()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "fail",
			"message": "Failed to connect to db.",
		})
		return
	}
	defer db.Close()

	row := db.QueryRow(sqlquery.SelectUserSessionQuery, currentUser.Uid, currentUser.SessionID)
	var refreshTokenInfo interfaces.UserRefreshTokenInfo
	dbErr := row.Scan(&refreshTokenInfo.ID, &refreshTokenInfo.RefreshTokenHash, &refreshTokenInfo.RefreshTokenExpiry)
	if dbErr != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "Invalid Session",
		})
		return
	}
	// Check expiry
	if time.Now().After(refreshTokenInfo.RefreshTokenExpiry) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "refresh token expired",
		})
		return
	}
	// Check if refresh token is valid
	hash, _ := security.CompareHash(newRefreshTokenHandlerPayload.RefreshToken, refreshTokenInfo.RefreshTokenHash)
	if !hash {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "invalid refresh token",
		})
		return
	}

	jwtTTL := os.Getenv("core.jwt.token.ttl")
	jwtTokenTTL, err := time.ParseDuration(jwtTTL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "fail",
			"message": "failed to Parse Duration",
		})
		return
	}
	newSessionId := uuid.New().String() // Generate new session
	jwtToken, err := security.GenerateJWT(currentUser.Uid, newSessionId, jwtTokenTTL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "fail",
			"message": "Failed to generate jwt token",
		})
		return
	}

	// Refresh token will expire soon — generating new refresh token
	timeLeft := time.Until(refreshTokenInfo.RefreshTokenExpiry)
	formattedRefreshTokenExpiry := refreshTokenInfo.RefreshTokenExpiry.Format("2006-01-02 15:04:05")
	if timeLeft <= jwtTokenTTL {
		refreshTokenInfo.RefreshTokenHash, _ = security.GenerateHash(uuid.NewString())
		newRefreshTokenTTL, _ := time.ParseDuration(os.Getenv("core.jwt.refresh.token.ttl"))
		formattedRefreshTokenExpiry = time.Now().Add(newRefreshTokenTTL).UTC().Format("2006-01-02 15:04:05")
	}

	_, err = db.Exec(sqlquery.UpdateUserSessionIDForUserQuery, newSessionId, refreshTokenInfo.RefreshTokenHash, formattedRefreshTokenExpiry, time.Now().UTC().Format("2006-01-02 15:04:05"), refreshTokenInfo.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "fail",
			"message": "unknown error occured",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":        "success",
		"message":       "New acess token generated",
		"access_token":  jwtToken,
		"refresh_token": newRefreshTokenHandlerPayload.RefreshToken,
	})

}
