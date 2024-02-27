package handlers

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenValue := c.GetHeader("Authorization")
		claims := &jwt.MapClaims{}
		tkn, err := jwt.ParseWithClaims(tokenValue, claims,
			func(t *jwt.Token) (interface{}, error) {
				return []byte(os.Getenv("JWT_SECRET")), nil
			},
		)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		if tkn == nil || !tkn.Valid {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		c.Next()
	}
}

type AuthHandler struct {
	ctx        context.Context
	collection *mongo.Collection
}

type JwtOutput struct {
	Token   string    `json:"token,omitempty"`
	Expires time.Time `json:"expires,omitempty"`
}

func NewAuthHandler(ctx context.Context, collection *mongo.Collection) *AuthHandler {
	return &AuthHandler{
		ctx:        ctx,
		collection: collection,
	}
}

func (r *AuthHandler) SignInHandler(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h := sha256.New()
	pswd := base64.StdEncoding.EncodeToString(h.Sum([]byte(user.Password)))

	log.Println("pasw", pswd)
	cur := r.collection.FindOne(r.ctx,
		bson.M{
			"username": user.Username,
			"password": pswd,
		},
	)
	if err := cur.Err(); err != nil {
		log.Println("err", err.Error())
		c.JSON(http.StatusUnauthorized,
			gin.H{"error": "Invalid username or password"})
		return
	}

	expireTime := time.Now().Add(10 * time.Minute)
	claims := jwt.MapClaims{
		"username": user.Username,
		"exp":      expireTime.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		panic(err)
	}
	c.JSON(http.StatusOK, JwtOutput{
		Token:   tokenString,
		Expires: expireTime,
	})
}

func (r *AuthHandler) RefreshHandler(c *gin.Context) {
	tokenValue := c.GetHeader("Authorization")
	claims := &jwt.MapClaims{}
	tkn, err := jwt.ParseWithClaims(tokenValue, claims,
		func(t *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		},
	)
	if err != nil || tkn == nil || !tkn.Valid {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
	jwtDate, err := claims.GetExpirationTime()
	if err != nil {
		panic(err)
	}

	if jwtDate.Time.Sub(time.Now()) > time.Second*30 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token is not expired yet"})
		return
	}

	expireTime := time.Now().Add(5 * time.Minute)
	(*claims)["exp"] = expireTime.Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, JwtOutput{
		Token:   tokenString,
		Expires: expireTime,
	})
}
