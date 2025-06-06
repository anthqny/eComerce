package controllers

import (
	"context"
	"github.com/anthqny/eComerce/database"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)
	type Aplication struct{
		prodCollection *mongo.Collection
		userCollection *mongo.Collection
	}
	func NewApplicaction(prodCollection, userCollection *mongo.Collection) *Aplication{
		return &Application{
			prodCollection: prodCollection,
			userCollection: userCollection
		}
	}
func (app *Applicaction) AddToCart() gin.Handler {
	return func(c *gin.Context){
		productQueryID := c.Query("id")
		if productQueryID == "" {
			log.Println("product is empty")

			_=	c.AbortWithError(http.StatusBadRequest, errors.New("product id is empty"))
			return
		}
		userQueryID:= c.Query("userID")
		if userQueryID == ""{
			log.Println("user id is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("user id is empty"))
			return
		}
		productID, err := primitive.ObjectIDFromHex(productQueryID)
		if err != nil{
			log.Println(err)
			c.AbortWithError(http.StatusInternalServerError)
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		
	err= database.AddProductToCart(ctx, app.prodCollection, app.userCollection,productID, userQueryID)
		if err != nil{
			c.IndentedJSON(http.StatusInternalServerError, err)
			
		}
		c.IndentedJSON(200, "Successfully added to the cart")
	}
}
func RemoveItem() gin.HandlerFunc {
	return func(c *gin.Context){
		productQueryID := c.Query("id")
		if productQueryID == "" {
			log.Println("product is empty")

			_=	c.AbortWithError(http.StatusBadRequest, errors.New("product id is empty"))
			return
		}
		userQueryID:= c.Query("userID")
		if userQueryID == ""{
			log.Println("user id is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("user id is empty"))
			return
		}
		productID, err := primitive.ObjectIDFromHex(productQueryID)
		if err != nil{
			log.Println(err)
			c.AbortWithError(http.StatusInternalServerError)
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err = database.RemoveCartItem(ctx, app.prodCollection, app.userCollection, ProductID, userQueryID)
		if err != nil{
			c.IndentedJSON(http.StatusInternalServerError, err)
			return
		}
		c.IndentedJSON(200, "Successfully removed item from the cart")

	}
}
func GetItemFrom() gin.HandlerFunc {
	return func(c *gin.Context){
		user_id:= c.Query("id")

		if user_id == ""{
			c.Header("Content-Type","application/json")
			c.JSON(http.StatusNotFound, gin.H{"error": "invalid id"})
			c.Abort()
			return
		}
		usert_id, _ := primitive.ObjectIDFromHex(user_id)

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var filledcar models.User
		err := UserCollection.FindOne(ctx, bson.D{primitive.E{Key:"_id", Value: usert_id}}).Decode(&filledcar)
		
		if err != nil{
			log.Println(err)
			c.IndentedJSON(500, "not found")
			return
		}
		filter_match := bson.D{{Key:"$match", Value: bson.D{primitive.E{Key:"_id", Value: usert_id}}}}
		unwind := bson.D{{Key: "$unwind", Value: bson.D{primitive.E{Key:"path", Value:"$usercart"}}}}
		grouping := bson.D{{Key: "$group", Value:bson.D{primitive.E{Key:"_id", Value: "$_id"},{Key:"total", Value:bson.D{primitive.E{Key:"$sum", Value: "$usercart.price"}}}}}}
		pointcursor, err:= UserCollection.Aggregate(ctx, mongo.Pipeline{filter_match, uwind, grouping})
		fi err != nil{
			log.Println(err)
		}
		var listing []bson.M
		err= pointcursor.All(ctx, &listing); err != nil{
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)	
		}
		for _, json := range listing{
			c.IndentedJSON(200,json["total"])
			c.IndentedJSON(200, filledcar.UserCart)
		}
		ctx.Done()
	}

	
}
func (app *Application) BuyFromCart() gin.HandlerFunc {
	return func(c *gin.Context){
		userQueryID:= c.Query("id")
		if userQueryID == ""{
			log.Println("user id is empty")
			_= c.AbortWithError(http.StatusBadRequest, errors.New("UserID is empty"))	
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		err :=database.BuyItemFromCart(ctx, app.userCollection, userQueryID)
		if err != nil{
			c.IndentedJSON(http.StatusInternalServerError, err)
		}

		c.IndentedJSON("succesfully placed the order ")
	}
}
func (app *Application) InstantBuy() gin.HandlerFunc {
	return func(c *gin.Context){
		productQueryID := c.Query("id")
		if productQueryID == "" {
			log.Println("product is empty")

			_=	c.AbortWithError(http.StatusBadRequest, errors.New("product id is empty"))
			return
		}
		userQueryID:= c.Query("userID")
		if userQueryID == ""{
			log.Println("user id is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("user id is empty"))
			return
		}
		productID, err := primitive.ObjectIDFromHex(productQueryID)
		if err != nil{
			log.Println(err)
			c.AbortWithError(http.StatusInternalServerError)
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err = database.InstantBuyer(ctx. app.prodCollection, app.userCollection, productID, UserQueryID)
		if err!= nil{
			c.IndentedJSON(htto.StatusInternalServerError, err)
		}

		c.IndentedJSON(200,"Successfullt removed item from the cart")
	}
}
