package helpers

import (
	"golang-authentication/prisma/db"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

var (
	Client    = CreateClient()
	SecretKey = GetSecret()
)

type Claims struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func CreateClient() *db.PrismaClient {
	client := db.NewClient()
	if err := client.Prisma.Connect(); err != nil {
		panic(err)
	}

	return client
}

func CloseClient(client *db.PrismaClient) {
	if err := client.Prisma.Disconnect(); err != nil {
		panic(err)
	}
}

func GetSecret() string {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	secret := os.Getenv("AUTH_SECRET")
	return secret
}

func GenerateToken(id, email string) (string, error) {
	expirationTime := time.Now().Add(1 * time.Hour)

	claims := &Claims{
		ID:    id,
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ParseToken(token string) (string, error) {
	claims := &Claims{}

	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})
	if err != nil {
		return "", err
	}

	if !tkn.Valid {
		return "", err
	}

	return claims.ID, nil
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("token")
		if err != nil {
			c.JSON(400, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		id, err := ParseToken(token)
		if err != nil {
			c.JSON(400, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		c.Set("id", id)

		c.Next()
	}
}
