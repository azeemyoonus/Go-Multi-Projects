package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/amitamrutiya/10-go-authentication-jwt/database"
	helper "github.com/amitamrutiya/10-go-authentication-jwt/helpers"
	"github.com/amitamrutiya/10-go-authentication-jwt/models"
	"github.com/go-playground/validator/v10"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")
var validate = validator.New()

func Singup() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := validate.Struct(user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if *user.Password != *user.ConfirmPassword {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Password and Confirm Password should be same"})
			return
		}

		count, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		if err != nil {
			log.Fatal(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occured while checking for the email"})
		}

		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email already exists"})
			return
		}

		user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()

		token, refreshToken, _ := helper.GenerateAllTokens(&user)
		user.Token = &token
		user.Refresh_token = &refreshToken

		hashedPassword := helper.HashPassowrd(*user.Password)
		user.Password = &hashedPassword
		user.ConfirmPassword = nil
		result, err := userCollection.InsertOne(ctx, user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while inserting user"})
			return
		}
		c.JSON(http.StatusOK, result)

	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var foundUser models.User
		err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while decoding user"})
			return
		}

		if foundUser.Email == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
			return
		}

		if !helper.VerifyPassword(*foundUser.Password, *user.Password) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid password"})
			return
		}

		token, refreshToken, _ := helper.GenerateAllTokens(&foundUser)
		err = helper.UpdateAllTokens(token, refreshToken, foundUser.ID.Hex())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while updating user"})
			return
		}
		user.Token = &token
		user.Refresh_token = &refreshToken

		user.Password = nil
		user.ConfirmPassword = nil

		c.JSON(http.StatusOK, user)
	}
}
