package controlles

import (
	"context"
	"log"
	"net/http"
	database "restaurant/Database"
	models "restaurant/Models"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type OrderItemPack struct {
	Table_id    *string            `json:"Table_id"`
	Order_Items []models.OrderItem `json:"order_items"`
}

var orderItemCollection *mongo.Collection = database.OpenCollection(database.Client, "orderItem")

func GetorderItems() gin.HandlerFunc {

	return func(c *gin.Context) {

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

		datas, err := orderItemCollection.Find(ctx, bson.M{})
		defer cancel()
		if err != nil {

			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while listing the order items"})
			return
		}

		Allitems := []bson.M{}
		for datas.Next(ctx) {

			var item bson.M

			if err := datas.Decode(&item); err != nil {

				log.Fatal(err)
			}

			Allitems = append(Allitems, item)
		}
		c.JSON(http.StatusOK, Allitems)
	}
}

func GetOrderItemsByOrder() gin.HandlerFunc {

	return func(c *gin.Context) {

		orderId := c.Param("order_id")

		all_item_orders, err := ItemsByOrder(orderId)

		log.Println("checking", all_item_orders)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while listing order items by order ID"})
			return
		}
		c.JSON(http.StatusOK, all_item_orders)
	}
}

func ItemsByOrder(id string) (OrterItems []primitive.M, err error) {

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	matchstage := bson.D{{"$match", bson.D{{"order_id", id}}}}
	lookupstage := bson.D{{"$lookup", bson.D{{"from", "food"}, {"localField", "food_id"}, {"foreignField", "food_id"}, {"as", "food"}}}}
	unwindstage := bson.D{{"$unwind", bson.D{{"path", "$food"}, {"preserveNullAndEmptyArrays", true}}}}

	lookupOrderStage := bson.D{{"$lookup", bson.D{{"from", "order"}, {"localField", "order_id"}, {"foreignField", "order_id"}, {"as", "order"}}}}
	unwindOrder := bson.D{{"$unwind", bson.D{{"path", "$order"}, {"preserveNullAndEmptyArrays", true}}}}

	lookupTableStage := bson.D{{"$lookup", bson.D{{"from", "table"}, {"localField", "order.table_id"}, {"foreignField", "table_id"}, {"as", "table"}}}}
	unwindtable := bson.D{{"$unwind", bson.D{{"path", "$table"}, {"preserveNullAndEmptyArrays", true}}}}

	// Project stage where amount is calculated
	project_stage := bson.D{
		{"$project", bson.D{
			{"_id", 0},
			{"amount", "$food.price"},
			{"name", "$food.name"},
			{"food_image", "$food.food_image"},
			{"table_number", "$table.table_number"},
			{"table_id", "$table.table_id"},
			{"order_id", "$order.order_id"},
			{"price", "$food.price"},
			{"quantity", 1},
		}},
	}

	groupStage := bson.D{
		{"$group", bson.D{ //{"table_id", "$table_id"}, {"table_number", "$table_number"}
			{"_id", bson.D{{"order_id", "$order_id"}}},
			{"payment_due", bson.D{{"$sum", "$amount"}}}, // Summing amount
			{"total_count", bson.D{{"$sum", 1}}},
			{"order_items", bson.D{{"$push", "$$ROOT"}}},
		}},
	}

	projectstage2 := bson.D{
		{"$project", bson.D{
			{"_id", 0},
			{"payment_due", 1},
			{"total_count", 1},
			{"table_number", "$_id.table_number"},
			{"order_items", 1},
		}},
	}

	// Aggregation pipeline execution
	result, err := orderItemCollection.Aggregate(ctx, mongo.Pipeline{
		matchstage,
		lookupstage,
		unwindstage,
		lookupOrderStage,
		unwindOrder,
		lookupTableStage,
		unwindtable,
		project_stage,
		groupStage,
		projectstage2,
	})
	if err != nil {
		return nil, err // Return error instead of panic
	}

	if err = result.All(ctx, &OrterItems); err != nil {
		return nil, err // Return error instead of panic
	}

	//log.Println("checking", OrterItems)
	return OrterItems, nil
}

func GetOrderItem() gin.HandlerFunc {

	return func(c *gin.Context) {

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		order_itemId := c.Param("order_item_id")
		var orderitem models.OrderItem
		err := orderItemCollection.FindOne(ctx, bson.M{"order_item_id": order_itemId}).Decode(&orderitem)

		defer cancel()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while listing ordered item"})
			return
		}
		c.JSON(http.StatusOK, orderitem)
	}
}

func UpdateOrderItem() gin.HandlerFunc {

	return func(c *gin.Context) {

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		orderitem_id := c.Param("order_item_id")

		filter := bson.M{"order_item_id": orderitem_id}

		var order_item models.OrderItem

		err := c.ShouldBindJSON(&order_item)

		if err != nil {

			c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
			return
		}

		var updatedobj primitive.D

		if order_item.Quantity != nil {

			updatedobj = append(updatedobj, bson.E{"quantity", *order_item.Quantity})
		}

		if order_item.Unit_price != nil {

			updatedobj = append(updatedobj, bson.E{"unit_price", *order_item.Unit_price})
		}

		if order_item.Food_id != nil {

			updatedobj = append(updatedobj, bson.E{"food_id", *order_item.Food_id})
		}

		order_item.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		updatedobj = append(updatedobj, bson.E{"updated_at", order_item.Updated_at})

		upsert := true

		opt := options.UpdateOptions{

			Upsert: &upsert,
		}
		result, err := orderItemCollection.UpdateOne(

			ctx,
			filter,
			bson.D{

				{"set", updatedobj},
			},

			&opt,
		)

		if err != nil {

			msg := "update is failed"
			c.JSON(http.StatusInternalServerError, gin.H{"Message": msg})
			return
		}

		defer cancel()
		c.JSON(http.StatusOK, result)
	}
}

func CreateOrderItem() gin.HandlerFunc {

	return func(c *gin.Context) {

		var ctx, cancel = context.WithTimeout(context.Background(), 200*time.Second)
		orderItemsToBeInserted := []interface{}{}

		var orderItam_pack OrderItemPack

		var order models.Order

		if err := c.ShouldBindJSON(&orderItam_pack); err != nil {

			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		order.Order_Date, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		order.Table_id = orderItam_pack.Table_id

		//	order_id := OrderItemOrderCreator(order)

		for _, orderItem := range orderItam_pack.Order_Items {

			//orderItem.Order_id = order_id

			validationErr := validate.Struct(orderItem)

			if validationErr != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
				return
			}

			orderItem.ID = primitive.NewObjectID()
			orderItem.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
			orderItem.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
			orderItem.Order_item_id = orderItem.ID.Hex()
			var num = toFixed(*orderItem.Unit_price, 2)
			orderItem.Unit_price = &num
			orderItemsToBeInserted = append(orderItemsToBeInserted, orderItem)

		}

		Insertedorder_items, err := orderItemCollection.InsertMany(ctx, orderItemsToBeInserted)

		if err != nil {
			log.Fatal(err)
		}
		defer cancel()
		c.JSON(http.StatusOK, Insertedorder_items)
	}
}
