package handlers

import (
	"context"
	"golang-authentication/helpers"
	"golang-authentication/prisma/db"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type RegisterCred struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginCred struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Register(c *gin.Context) {
	var cred RegisterCred
	if err := c.ShouldBindJSON(&cred); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(cred.Password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	existingUser, err := helpers.Client.User.FindUnique(
		db.User.Email.Equals(cred.Email),
	).Exec(ctx)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if existingUser != nil {
		c.JSON(400, gin.H{"error": "Email already taken"})
		return
	}

	_, err = helpers.Client.User.CreateOne(
		db.User.Name.Set(cred.Name),
		db.User.Email.Set(cred.Email),
		db.User.Password.Set(string(hashedPassword)),
	).Exec(ctx)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"success": "User created successfully",
	})
}

func Login(c *gin.Context) {
	var cred LoginCred
	if err := c.ShouldBindJSON(&cred); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	ctx := context.Background()

	user, err := helpers.Client.User.FindUnique(
		db.User.Email.Equals(cred.Email),
	).Exec(ctx)
	if err != nil {
		if err.Error() == "ErrNotFound" {
			c.JSON(400, gin.H{"error": "User does not exists"})
			return
		}
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(cred.Password))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid password"})
		return
	}

	tokenString, err := helpers.GenerateToken(user.ID, user.Email)
	if err != nil {
		c.JSON(400, gin.H{"error": "Error generating token"})
		return
	}

	c.SetCookie("token", tokenString, 3600, "/", "localhost", false, true)

	c.JSON(200, gin.H{
		"success": "Logged in successfully",
		"token":   tokenString,
	})
}

func Logout(c *gin.Context) {
	c.SetCookie("token", "", -1, "/", "localhost", false, true)

	c.JSON(200, gin.H{
		"success": "Logged out successfully",
	})
}

func Me(c *gin.Context) {
	idAny, _ := c.Get("id")
	id, ok := idAny.(string)
	if !ok {
		c.JSON(400, gin.H{"error": "Not a valid User ID"})
		return
	}

	ctx := context.Background()

	user, err := helpers.Client.User.FindUnique(
		db.User.ID.Equals(id),
	).Omit(
		db.User.ID.Field(),
		db.User.Password.Field(),
	).Exec(ctx)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if user == nil {
		c.JSON(400, gin.H{"error": "User does not exists"})
		return
	}

	c.JSON(200, gin.H{
		"user": user,
	})
}
