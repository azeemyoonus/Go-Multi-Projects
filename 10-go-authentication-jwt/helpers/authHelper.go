package helpers

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/amitamrutiya/10-go-authentication-jwt/database"
	"github.com/amitamrutiya/10-go-authentication-jwt/models"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type SignedDetails struct {
	ID        string
	Email     string
	Username  string
	User_type string
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
		ID:        user.ID.Hex(),
		Email:     *user.Email,
		Username:  *user.Username,
		User_type: *user.User_type,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}

	refreshClaims := &SignedDetails{
		ID:        user.ID.Hex(),
		Email:     *user.Email,
		Username:  *user.Username,
		User_type: *user.User_type,
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
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes)
}

func VerifyPassword(hashedPassword string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func UpdateAllTokens(signedToken string, signedRefreshToken string, userId string) error {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	_, err := userCollection.UpdateOne(ctx, bson.M{"_id": userId}, bson.M{"$set": bson.M{"token": signedToken, "refresh_token": signedRefreshToken, "updated_at": time.Now().Format(time.RFC3339)}})
	if err != nil {
		return err
	}
	return nil
}

func ValidateToken(signedToken string) (*SignedDetails, error) {
	claims := &SignedDetails{}
	token, err := jwt.ParseWithClaims(signedToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(SECRET_KEY), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("Invalid token")
	}
	// if claims.ExpiresAt < time.Now().Local().Unix() {
	// 	return nil, errors.New("Token expired")
	// }
	return claims, nil
}
