package controllers

import (
	"fmt"
	"strconv"
	"time"

	"context"
	"net/http"

	"github.com/amitamrutiya/16-restaurant-management/database"
	helper "github.com/amitamrutiya/16-restaurant-management/helpers"
	"github.com/amitamrutiya/16-restaurant-management/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")
var validate = validator.New()

func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))
		if err != nil || recordPerPage < 1 {
			recordPerPage = 10
		}
		page, err1 := strconv.Atoi(c.Query("page"))
		if err1 != nil || page < 1 {
			page = 1
		}

		startIndex := (page - 1) * recordPerPage
		startIndex, err = strconv.Atoi(c.Query("startIndex"))
		if err != nil {
			startIndex = 0
		}

		matchStage := bson.D{{"$match", bson.M{}}}
		groupStage := bson.D{{"$group", bson.M{
			"_id":         nil,
			"total_count": bson.M{"$sum": 1},
			"data":        bson.M{"$push": "$$ROOT"}}}}
		projectStage := bson.D{{"$project", bson.M{
			"_id":         0,
			"total_count": 1,
			"data":        bson.M{"$slice": []interface{}{"$data", startIndex, recordPerPage}}}}}

		countCursor, err := userCollection.Aggregate(ctx, mongo.Pipeline{matchStage, groupStage, projectStage})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occured while fetching the users"})
		}
		var allUsers []bson.M
		if err = countCursor.All(ctx, &allUsers); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occured while fetching the users"})
		}

		c.JSON(http.StatusOK, allUsers[0])
	}
}

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User
		userId := c.Param("user_id")
		objId, err := primitive.ObjectIDFromHex(userId)
		filter := bson.M{"_id": objId}
		fmt.Println("userId" + userId)
		err = userCollection.FindOne(ctx, filter).Decode(&user)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found" + err.Error()})
			return
		}
		c.JSON(http.StatusOK, user)
	}
}

func DeleteUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		userId := c.Param("id")
		fmt.Println("userId" + userId)

		if err := helper.MatchUserTypeToUid(c, userId); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		objectID, err := primitive.ObjectIDFromHex(userId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Conversion error"})
			return
		}
		filter := bson.M{"_id": objectID}
		result, err := userCollection.DeleteOne(ctx, filter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while deleting user"})
			return
		}
		c.JSON(http.StatusOK, result)
	}
}
