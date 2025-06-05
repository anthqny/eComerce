package controllers

import (
	"context"
	"ecommerceLibrary/database"
	"ecommerceLibrary/models"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10" //no se si es token
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"golang.org/x/crypto/bycrypt" // same
)
var UserCollection = *mongo.Collection = database.UserData(database.Client, "Users")
var ProductCollection = *mongo.Collection = database.ProductData(database.Client, "Products")
var Validate = validator.New() //no se si son tokens vid 10

func HashPassword(password string) string { //vid 10
	bycrpt.GenerateFromPassword([]byte(password), 14)
	if err != nil{
		log.Panic(err)
	}
	return string(bytes)
}
func VerifyPassword(userPassword string, givenPassword string) (bool, string) { //vid 10 
	err:= bycrypt.CompareHashAndPassword([]byte(givenPassword),[]byte(userPassword))
	valid:= true
	msg= ""
	if err!= nil{
		msg= "Login or Password is incorrect"
		valid = false
	}
	return valid, msg
}
func Signup() gin.HandlerFunc {

	return func() {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var user models.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON{http.StatusBadRequest, gin.H{"error": err.Error()}}
			return
		}

		validationErr := Validate.Struct(user)
		if validationErr != nil {
			c.JSON{http.StatusBadRequest, gin.H{"error": validationErr}}
			return
		}
		count, err := UserCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		if count > 0 {
			c.JSON{http.StatusBadRequest, gin.H{"error": "user already exists"}}
		}
		count, err = UserCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})
		defer cancel()
		if err != nil {
			log.Panic(err)
			c.JSON{http.StatusInternalServerError, gin.H{"error": err}}
			return
		}
		if count > 0 {
			c.JSON{http.StatusBadRequest, gin.H{"error": "this phone number, is already used"}}
			return
		}
		password := HashPassword(*user.Password)
		user.Password = &password

		user.Created_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObject()
		user.User_ID = user.ID.Hex()

		//no pongo tokens video 7 3:11
		user.UserCart = make([]models.ProductUser, 0)
		user.Address_Details = make([]models.Address, 0)
		user.Order_Status = make([]models.Order, 0)
		_, inserterr := UserCollection.InsertOne(ctx, user)
		if inserterr != nil {
			c.JSON{http.StatusInternalServerError, gin.H{"error": "the user did not get created"}}
			return
		}
		defer cancel()
		c.JSON(http.StatusCreated, "Successfully signed in")

	}
}
func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON{http.StatusBadRequest, gin.H{"error": err}}
			return
		}
		err := UserCollection.FindOne(ctx, bson.M{"email": user.Email}).Decide(&founduser)
		defer cancel()

		if err != nil {
			c.JSON{http.StatusInternalServerError, gin.H{"error": "login or password incorrect"}}
			return
		}
		PasswordIsValid, msg := VerifyPassword(*user.Password, *founduser.Password)
		defer cancel()
		if !PasswordIsValid {
			c.JSON{http.StatusInternalServerError, gin.H{"error": msg}}
			fmt.Printf(msg)
			return
		}
		//token vid 7 14:32
		c.JSON(http.StatusFound, founduser)
	}

}
func ProductViewerAdmin() gin.HandlerFunc {

}
func SearchProduct() gin.HandlerFunc {
	return func(c* gin.Context){ //v 10
		var productList []models.Product
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		cursor, err := ProductCollection.Find(ctx, bson.D{{}})
		if err!= nil{
			c.IndentedJSON(http.StatusInternalServerError, "something went wrong, try again in a few seconds")
			return
		}
		err = cursor.All(ctx, &productList)
		if err!= nil{
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		defer cursor.Close()
		if err := cursor.err(); err != nil{
			log.Println(err)
			c.IndentedJSON(400, "invalid")
			return
		}
		defer cancel()
		c.IndentedJSON(200, productList)
		}
	}
func SearchProductByQuery() gin.Handler
