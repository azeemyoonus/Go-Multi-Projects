package helpers

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/amitamrutiya/16-restaurant-management/database"
	"github.com/amitamrutiya/16-restaurant-management/models"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

type SignedDetails struct {
	Email      string
	First_name string
	Last_name  string
	Uid        string
	jwt.StandardClaims
}

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")

var SECRET_KEY string = os.Getenv("JWT_SECRET_KEY")

func MatchUserTypeToUid(c *gin.Context, userId string) error {
	userType := c.GetString("user_type")
	uid := c.GetString("uid")
	if userType != "admin" && uid != userId {
		return errors.New("Unauthorized to access this resource")
	}
	err := CheckUserType(c, userType)
	return err
}

func CheckUserType(c *gin.Context, role string) error {
	userType := c.GetString("user_type")
	if userType != role {
		return errors.New("Unauthorized to access this resource")
	}
	return nil
}

func GenerateAllTokens(user *models.User) (signedToken string, signedRefreshToken string, err error) {
	claims := &SignedDetails{
		Email:      *user.Email,
		First_name: *user.First_name,
		Last_name:  *user.Last_name,
		Uid:        user.User_id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}

	refreshClaims := &SignedDetails{
		Email:      *user.Email,
		First_name: *user.First_name,
		Last_name:  *user.Last_name,
		Uid:        user.User_id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24*7)).Unix(),
		},
	}

	signedToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
	signedRefreshToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SECRET_KEY))

	if err != nil {
		return "", "", err
	}
	return signedToken, signedRefreshToken, nil
}

func HashPassowrd(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}

func VerifyPassword(hashedPassword string, password string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	check := true
	msg := ""
	if err != nil {
		msg = "Login password is incorrect"
		check = false
		return check, msg
	}
	return check, msg
}

func UpdateAllTokens(signedToken string, signedRefreshToken string, userId string) error {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var updateObj primitive.D
	updateObj = append(updateObj, bson.E{Key: "token", Value: signedToken})
	updateObj = append(updateObj, bson.E{Key: "refresh_token", Value: signedRefreshToken})
	Updated_at, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj = append(updateObj, bson.E{Key: "updated_at", Value: Updated_at})

	filter := bson.M{"user_id": userId}
	upsert := true

	opt := options.UpdateOptions{
		Upsert: &upsert,
	}

	_, err := userCollection.UpdateOne(ctx, filter, bson.D{{"$set", updateObj}}, &opt)

	if err != nil {
		return err
	}
	return nil
}

func ValidateToken(signedToken string) (*SignedDetails, error) {
	token, err := jwt.ParseWithClaims(signedToken, &SignedDetails{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SECRET_KEY), nil
	})

	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		return nil, errors.New("Invalid token")
	}

	if !token.Valid {
		return nil, errors.New("Invalid token")
	}
	if claims.ExpiresAt < time.Now().Local().Unix() {
		return nil, errors.New("Token expired")
	}
	return claims, nil
}
