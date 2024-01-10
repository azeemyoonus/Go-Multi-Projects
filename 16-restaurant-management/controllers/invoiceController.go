package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/amitamrutiya/16-restaurant-management/database"
	"github.com/amitamrutiya/16-restaurant-management/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type InvoiceViewFormat struct {
	Invoice_id       string
	Payment_method   string
	Order_id         string
	Payment_status   *string
	Payment_due      interface{}
	Table_number     interface{}
	Payment_due_date time.Time
	Order_details    interface{}
}

var invoiceCollection *mongo.Collection = database.OpenCollection(database.Client, "invoice")

func GetInvoices() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("GetInvoices")
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		result, err := invoiceCollection.Find(context.TODO(), bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while fetching invoice items"})
			return
		}

		var allInvoices []bson.M
		if err = result.All(ctx, &allInvoices); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while decoding invoice items"})
			return
		}

		c.JSON(200, allInvoices)
	}
}

func GetInvoice() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("GetInvoice")
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		InvoiceId := c.Param("invoice_id")
		objID,_ := primitive.ObjectIDFromHex(InvoiceId)
		filter := bson.M{"_id": objID}
		var invoice models.Invoice

		err := invoiceCollection.FindOne(ctx, filter).Decode(&invoice)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while fetching invoice item" + err.Error()})
			return
		}
		var invoiceView InvoiceViewFormat

		allOrderItems, err := ItemsByOrder(invoice.Order_id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while fetching order items"})
			return
		}
		fmt.Println(allOrderItems)
		invoiceView.Order_id = invoice.Order_id
		invoiceView.Payment_due_date = invoice.Payment_due_date
		invoiceView.Payment_method = "null"
		if invoice.Payment_method != nil {
			invoiceView.Payment_method = *invoice.Payment_method
		}
		invoiceView.Invoice_id = invoice.Invoice_id
		invoiceView.Payment_status = *&invoice.Payment_status
		invoiceView.Payment_due = allOrderItems[0]["payment_due"]
		invoiceView.Table_number = allOrderItems[0]["table_number"]
		invoiceView.Order_details = allOrderItems[0]["order_details"]

		c.JSON(200, invoiceView)
	}
}

func CreateInvoice() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("CreateInvoice")
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var order models.Order
		var invoice models.Invoice
		if err := c.BindJSON(&invoice); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		err := orderCollection.FindOne(ctx, bson.M{"order_id": invoice.Order_id}).Decode(&order)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while fetching order item"})
			return
		}

		status := "PENDING"
		if invoice.Payment_status != nil {
			invoice.Payment_status = &status
		}

		invoice.Payment_due_date, _ = time.Parse(time.RFC3339, time.Now().AddDate(0, 0, 1).Format(time.RFC3339))
		invoice.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		invoice.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		invoice.ID = primitive.NewObjectID()
		invoice.Invoice_id = invoice.ID.Hex()

		validationErr := validate.Struct(invoice)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error while validating request body " + validationErr.Error()})
			return
		}

		result, insertErr := invoiceCollection.InsertOne(ctx, invoice)
		if insertErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while inserting new invoice " + insertErr.Error()})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

func UpdateInvoice() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("UpdateInvoice")
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		InvoiceId := c.Param("invoice_id")
		filter := bson.M{"invoice_id": InvoiceId}

		var invoice models.Invoice
		if err := c.BindJSON(&invoice); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var updateObj primitive.D

		if invoice.Payment_method != nil {
			updateObj = append(updateObj, bson.E{Key: "payment_method", Value: invoice.Payment_method})
		}

		if invoice.Payment_status != nil {
			updateObj = append(updateObj, bson.E{Key: "payment_status", Value: invoice.Payment_status})
		}

		invoice.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObj = append(updateObj, bson.E{Key: "updated_at", Value: invoice.Updated_at})

		upsert := true
		opt := options.UpdateOptions{
			Upsert: &upsert,
		}

		status := "PENDING"
		if invoice.Payment_status != nil {
			invoice.Payment_status = &status
		}

		result, err := invoiceCollection.UpdateOne(
			ctx,
			filter,
			bson.D{{"$set", updateObj}},
			&opt,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while updating invoice item"})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

func DeleteInvoice() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("DeleteInvoice")
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		InvoiceId := c.Param("invoice_id")
		filter := bson.M{"invoice_id": InvoiceId}

		result, err := invoiceCollection.DeleteOne(ctx, filter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while deleting invoice item"})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}
