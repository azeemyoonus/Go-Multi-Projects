package controllers

import (
	"fmt"
	"strconv"
	"time"

	"context"
	"net/http"

	helper "github.com/amitamrutiya/10-go-authentication-jwt/helpers"
	"github.com/amitamrutiya/10-go-authentication-jwt/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {

		if err := helper.CheckUserType(c, "ADMIN"); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))
		if err != nil {
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
		fmt.Println("allUsers", allUsers)
		// if len(allUsers) > 0 {
		// 	totalCount := allUsers[0]["total_count"]
		// 	c.JSON(http.StatusOK, gin.H{"data": allUsers[0]["data"], "totalCount": totalCount})
		// 	return
		// }
		c.JSON(http.StatusOK, allUsers[0])
	}
}

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		userId := c.Param("id")
		fmt.Println("userId" + userId)

		if err := helper.MatchUserTypeToUid(c, userId); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		objectID, err := primitive.ObjectIDFromHex(userId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Conversion error"})
			return
		}
		filter := bson.M{"_id": objectID}
		err = userCollection.FindOne(ctx, filter).Decode(&user)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusOK, user)
	}
}
