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

type OrderItemPack struct {
	Table_id    *string
	Order_items []models.OrderItem
}

var orderItemCollection *mongo.Collection = database.OpenCollection(database.Client, "orderItem")

// GetOrderItems : Get all orderItems
func GetOrderItems() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("GetOrderItems")
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		result, err := orderItemCollection.Find(context.TODO(), bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while fetching orderItem items"})
			return
		}

		var allOrderItems []bson.M
		if err = result.All(ctx, &allOrderItems); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while decoding orderItem items"})
			return
		}

		c.JSON(http.StatusOK, allOrderItems)
	}
}

// GetOrderItem : Get orderItem by id
func GetOrderItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("GetOrderItem")
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		orderItemId := c.Param("orderItem_id")
		objID, err := primitive.ObjectIDFromHex(orderItemId)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid orderItem ID format"})
			return
		}
		filter := bson.M{"_id": objID}

		var orderItem bson.M
		if err := orderItemCollection.FindOne(ctx, filter).Decode(&orderItem); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while fetching orderItem" + err.Error()})
			return
		}

		c.JSON(http.StatusOK, orderItem)
	}
}

// GetOrderItemsByOrder : Get orderItem by order
func GetOrderItemsByOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("GetOrderItemsByOrder")

		orderId := c.Param("order_id")
		allOrderItems, err := ItemsByOrder(orderId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while fetching order items"})
			return
		}

		c.JSON(http.StatusOK, allOrderItems)
	}
}

// ItemsByOrder : Get orderItem by order id
func ItemsByOrder(id string) (OrderItems []primitive.M, err error) {
	fmt.Println("ItemsByOrder")
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	matchStage := bson.D{{"$match", bson.D{{"order_id", id}}}}

	lookupStage := bson.D{{"$lookup", bson.D{{"from", "food"}, {"localField", "food_id"}, {"foreignField", "food_id"}, {"as", "food"}}}}
	unwindStage := bson.D{{"$unwind", bson.D{{"path", "$food"}, {"preserveNullAndEmptyArrays", true}}}}

	lookupOrderStage := bson.D{{"$lookup", bson.D{{"from", "order"}, {"localField", "order_id"}, {"foreignField", "order_id"}, {"as", "order"}}}}
	unwindOrderStage := bson.D{{"$unwind", bson.D{{"path", "$order"}, {"preserveNullAndEmptyArrays", true}}}}

	lookupTableStage := bson.D{{"$lookup", bson.D{{"from", "table"}, {"localField", "order.table_id"}, {"foreignField", "table_id"}, {"as", "table"}}}}
	unwindTableStage := bson.D{{"$unwind", bson.D{{"path", "$table"}, {"preserveNullAndEmptyArrays", true}}}}

	projectStage := bson.D{
		{
			"$project", bson.D{
				{"_id", 0},
				{"amount", "$food.price"},
				{"total_count", 1},
				{"food_name", "$food.food_name"},
				{"food_image", "$food.food_image"},
				{"table_number", "$table.table_number"},
				{"table_id", "$table.table_id"},
				{"order_id", "$_id"},
				{"price", "$food.price"},
				{"quantity", 1},
			},
		},
	}

	groupStage := bson.D{
		{
			"$group", bson.D{
				{"_id", bson.D{{"order_id", "$order_id"}, {"table_id", "$table_id"}, {"table_number", "$table_number"}}},
				{"payment_due", bson.D{{"$sum", "$amount"}}},
				{"total_count", bson.D{{"$sum", 1}}},
				{"order_items", bson.D{{"$push", "$$ROOT"}}},
			},
		},
	}

	projectStage2 := bson.D{
		{
			"$project", bson.D{
				{"_id", 0},
				{"order_id", "$_id.order_id"},
				{"table_id", "$_id.table_id"},
				{"table_number", "$_id.table_number"},
				{"payment_due", 1},
				{"total_count", 1},
				{"order_items", 1},
			},
		},
	}

	cursor, err := orderItemCollection.Aggregate(ctx, mongo.Pipeline{matchStage, lookupStage, unwindStage, lookupOrderStage, unwindOrderStage, lookupTableStage, unwindTableStage, projectStage, groupStage, projectStage2})
	if err != nil {
		return nil, err
	}

	if err = cursor.All(ctx, &OrderItems); err != nil {
		return nil, err
	}

	return OrderItems, nil
}

// CreateOrderItem : Create orderItem
func CreateOrderItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("CreateOrderItem")
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var orderItemPack OrderItemPack
		var order models.Order

		if err := c.BindJSON(&orderItemPack); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		order.Order_Date, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		order.Table_id = orderItemPack.Table_id
		order_id := OrderItemOrderCreator(order)
		orderItemsToBeInserted := []interface{}{}
		for _, orderItem := range orderItemPack.Order_items {
			orderItem.Order_id = order_id
			validationErr := validate.Struct(orderItem)
			if validationErr != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Error while validating request body " + validationErr.Error()})
				return
			}

			orderItem.ID = primitive.NewObjectID()
			orderItem.OrderItem_id = orderItem.ID.Hex()
			orderItem.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
			orderItem.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
			var num = toFixed(*orderItem.Unit_price, 2)
			orderItem.Unit_price = &num
			orderItemsToBeInserted = append(orderItemsToBeInserted, orderItem)
		}
		insertedOrderItems, err := orderItemCollection.InsertMany(ctx, orderItemsToBeInserted)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while inserting new orderItem " + err.Error()})
			return
		}

		c.JSON(http.StatusOK, insertedOrderItems)
	}
}

// UpdateOrderItem : Update orderItem
func UpdateOrderItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("UpdateOrderItem")
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var orderItem models.OrderItem

		orderItemId := c.Param("orderItem_id")
		filter := bson.M{"orderItem_id": orderItemId}

		var updateObj primitive.D

		if orderItem.Unit_price != nil {
			updateObj = append(updateObj, bson.E{Key: "unit_price", Value: *&orderItem.Unit_price})
		}
		if orderItem.Quantity != nil {
			updateObj = append(updateObj, bson.E{Key: "quantity", Value: *orderItem.Quantity})
		}
		if orderItem.Food_id != nil {
			updateObj = append(updateObj, bson.E{Key: "food_id", Value: *orderItem.Food_id})
		}

		orderItem.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObj = append(updateObj, bson.E{Key: "updated_at", Value: orderItem.Updated_at})

		upsert := true
		opt := options.UpdateOptions{
			Upsert: &upsert,
		}

		result, err := orderItemCollection.UpdateOne(ctx, filter, bson.D{
			{"$set", updateObj},
		}, &opt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while updating orderItem " + err.Error()})
			return
		}

		c.JSON(http.StatusOK, result)

	}
}

// DeleteOrderItem : Delete orderItem
func DeleteOrderItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("DeleteOrderItem")
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		orderItemId := c.Param("orderItem_id")
		filter := bson.M{"orderItem_id": orderItemId}

		result, err := orderItemCollection.DeleteOne(ctx, filter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while deleting orderItem " + err.Error()})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}
