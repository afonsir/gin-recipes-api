package handlers

import (
	"context"
	"crypto/sha256"
	"net/http"
	"os"
	"time"

	"github.com/auth0-community/go-auth0"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/rs/xid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/square/go-jose.v2"

	models "github.com/afonsir/gin-recipes-api/models"
)

type AuthHandler struct {
	collection    *mongo.Collection
	ctx           context.Context
	AuthMechanism string
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

type JWTOutput struct {
	Token   string    `json:"token"`
	Expires time.Time `json:"expires"`
}

func NewAuthHandler(
	ctx context.Context,
	collection *mongo.Collection,
	authCookieEnabled string,
) *AuthHandler {
	return &AuthHandler{
		collection:    collection,
		ctx:           ctx,
		AuthMechanism: authCookieEnabled,
	}
}

func (handler *AuthHandler) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		switch handler.AuthMechanism {
		case "COOKIE":
			session := sessions.Default(c)
			sessionToken := session.Get("token")

			if sessionToken == nil {
				c.JSON(http.StatusForbidden, gin.H{
					"error": "Not logged",
				})
				c.Abort()
			}

			c.Next()
		case "JWT":
			tokenValue := c.GetHeader("Authorization")
			claims := &Claims{}

			tkn, err := jwt.ParseWithClaims(tokenValue, claims,
				func(token *jwt.Token) (interface{}, error) {
					return []byte(os.Getenv("JWT_SECRET")), nil
				})

			if err != nil {
				c.AbortWithStatus(http.StatusUnauthorized)
			}

			if tkn == nil || !tkn.Valid {
				c.AbortWithStatus(http.StatusUnauthorized)
			}

			c.Next()
		case "AUTH0":
			var auth0Domain = "https://" + os.Getenv("AUTH0_DOMAIN") + "/"

			client := auth0.NewJWKClient(auth0.JWKClientOptions{
				URI: auth0Domain + ".well-known/jwks.json"},
				nil,
			)

			configuration := auth0.NewConfiguration(
				client,
				[]string{os.Getenv("AUTH0_API_IDENTIFIER")},
				auth0Domain,
				jose.RS256,
			)

			validator := auth0.NewValidator(configuration, nil)
			_, err := validator.ValidateRequest(c.Request)

			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{
					"message": "Invalid token",
				})
				c.Abort()
				return
			}

			c.Next()
		}
	}
}

// swagger:operation POST /signin users signIn
//
// Sign-in user, returns a token with 10 minutes of expiration time
//
// ---
// parameters:
// - name: user
//   in: body
//   description: user credentials
//   schema:
//     type: object
//     required:
//       - username
//       - password
//     properties:
//       username:
//         type: string
//       password:
//         type: string
//
// consumes:
// - application/json
//
// produces:
// - application/json
//
// responses:
//   '200':
//     description: Successful operation
//   '401':
//     description: Invalid input
//   '500':
//     description: Internal Server Error
func (handler *AuthHandler) SignInHandler(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	h := sha256.New()
	cur := handler.collection.FindOne(handler.ctx, bson.M{
		"username": user.Username,
		"password": string(h.Sum([]byte(user.Password))),
	})

	if cur.Err() != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid username or password",
		})
		return
	}

	expirationTime := time.Now().Add(10 * time.Minute)
	claims := &Claims{
		Username: user.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}

	jwtOutput := JWTOutput{
		Token:   tokenString,
		Expires: expirationTime,
	}

	c.JSON(http.StatusOK, jwtOutput)
}

func (handler *AuthHandler) SignInWithCookieHandler(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	h := sha256.New()
	cur := handler.collection.FindOne(handler.ctx, bson.M{
		"username": user.Username,
		"password": string(h.Sum([]byte(user.Password))),
	})

	if cur.Err() != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid username or password",
		})
		return
	}

	sessionToken := xid.New().String()
	session := sessions.Default(c)

	session.Set("username", user.Username)
	session.Set("token", sessionToken)
	session.Save()

	c.JSON(http.StatusOK, gin.H{
		"message": "User signed in",
	})
}

// swagger:operation POST /refresh users refreshToken
//
// Refresh a given valid token (+5 minutes)
//
// ---
// parameters:
// - name: user
//   in: body
//   description: user credentials
//   schema:
//     type: object
//     required:
//       - username
//       - password
//     properties:
//       username:
//         type: string
//       password:
//         type: string
//
// consumes:
// - application/json
//
// produces:
// - application/json
//
// responses:
//   '200':
//     description: Successful operation
//   '401':
//     description: Invalid input
//   '500':
//     description: Internal Server Error
func (handler *AuthHandler) RefreshHandler(c *gin.Context) {
	tokenValue := c.GetHeader("Authorization")
	claims := &Claims{}

	tkn, err := jwt.ParseWithClaims(tokenValue, claims,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	if tkn == nil || !tkn.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid token",
		})
		return
	}

	if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) > 30*time.Second {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Token is not expired yet",
		})
		return
	}

	expirationTime := time.Now().Add(5 * time.Minute)
	claims.ExpiresAt = expirationTime.Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(os.Getenv("JWT_SECRET"))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	jwtOutput := JWTOutput{
		Token:   tokenString,
		Expires: expirationTime,
	}

	c.JSON(http.StatusOK, jwtOutput)
}

// swagger:operation POST /signout users signOut
//
// Signout user (removes auth cookie)
//
// ---
// produces:
// - application/json
//
// responses:
//   '200':
//     description: Successful operation
func (handler *AuthHandler) SignOutHandler(c *gin.Context) {
	session := sessions.Default(c)

	session.Clear()
	session.Save()

	c.JSON(http.StatusOK, gin.H{
		"message": "Signed out...",
	})
}
