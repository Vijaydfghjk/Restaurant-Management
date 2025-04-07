package controller

import (
	"context"
	"log"
	"net/http"
	database "restaurant-management/Database"
	model "restaurant-management/Model"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Orderitem_db struct {
	orderitem_collection *mongo.Collection
	validate             *validator.Validate
}

type OrderItemPack struct {
	Order_Items []model.OrderItem `json:"order_items"`
}

func Orderitemcontroll() *Orderitem_db {

	return &Orderitem_db{
		orderitem_collection: database.Open_collection(database.Client, "Orderitem"),
		validate:             validator.New(),
	}

}
func (a *Orderitem_db) CreateOrderItem(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	defer cancel()

	var myorders OrderItemPack

	if err := c.ShouldBindJSON(&myorders); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var orderItems []interface{}
	var wg sync.WaitGroup
	orderItemChan := make(chan interface{}, len(myorders.Order_Items))

	errorChan := make(chan error, len(myorders.Order_Items))

	for _, orderitem := range myorders.Order_Items {
		wg.Add(1)
		go func(item model.OrderItem) {

			defer wg.Done()

			if validationErr := a.validate.Struct(item); validationErr != nil {

				errorChan <- validationErr
				return
			}

			item.ID = primitive.NewObjectID()
			item.Order_item_id = item.ID.Hex()
			item.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
			item.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

			var qty = toFixed(*item.Quantity, 2)
			item.Quantity = &qty

			orderItemChan <- item

		}(orderitem)

	}

	wg.Wait()
	close(orderItemChan)
	close(errorChan)

	for err := range errorChan {
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
			return
		}
	}

	for values := range orderItemChan {

		orderItems = append(orderItems, values)
	}

	results, err := a.orderitem_collection.InsertMany(ctx, orderItems)

	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}

func (a *Orderitem_db) GetOrderItemsByOrder(c *gin.Context) {

	order_id := c.Param("order_id")

	all_item_orders, err := a.ItemsByOrder(order_id)

	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while listing order items by order ID"})
		return
	}

	c.JSON(http.StatusOK, all_item_orders)
}

func (a *Orderitem_db) ItemsByOrder(id string) (OrterItems []primitive.M, err error) {

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	matchstage := bson.D{{"$match", bson.D{{"order_id", id}}}}

	foodlookup := bson.D{{"$lookup", bson.D{{"from", "Food"}, {"localField", "food_id"}, {"foreignField", "food_id"}, {"as", "food"}}}}
	unwindfood := bson.D{{"$unwind", bson.D{{"path", "$food"}, {"preserveNullAndEmptyArrays", true}}}}

	orderlookup := bson.D{{"$lookup", bson.D{{"from", "Order"}, {"localField", "order_id"}, {"foreignField", "order_id"}, {"as", "order"}}}}
	unwindorder := bson.D{{"$unwind", bson.D{{"path", "$order"}, {"preserveNullAndEmptyArrays", true}}}}

	tablelookup := bson.D{{"$lookup", bson.D{{"from", "Table"}, {"localField", "order.table_id"}, {"foreignField", "table_id"}, {"as", "table"}}}}
	unwindtable := bson.D{{"$unwind", bson.D{{"path", "$table"}, {"preserveNullAndEmptyArrays", true}}}}

	projectStage := bson.D{

		{"$project", bson.D{

			{"_id", 0},
			{"name", "$food.name"},
			{"quantity", 1},
			{"Unit_price", "$food.price"},
			{"food_image", "$food.food_image"},
			{"table_number", "$table.table_number"},
			{"table_id", "$table.table_id"},
			{"order_id", "$order.order_id"},
		},
		},
	}

	groupstage := bson.D{
		{"$group", bson.D{
			{"_id", bson.D{{"order_id", "$order_id"}, {"name", "$name"}, {"table_number", "$table_number"}, {"table_id", "$table_id"}}},
			{"payment_due", bson.D{{"$sum", bson.D{{"$multiply", bson.A{"$Unit_price", "$quantity"}}}}}},
			{"quantity", bson.D{{"$sum", "$quantity"}}},
			{"food_image", bson.D{{"$first", "$food_image"}}},
			//{"table_number", bson.D{{"$first", "$table_number"}}},
			{"Unit_price", bson.D{{"$first", "$Unit_price"}}},
			//{"table_id", bson.D{{"$first", "$table_id"}}},
		}},
	}

	projectStage2 := bson.D{
		{"$project", bson.D{
			{"_id", 0},
			{"OrderId", "$_id.order_id"},
			{"Table_number", "$_id.table_number"},
			{"Table_id", "$_id.table_id"},
			{"payment_due", 1},
			{"order_items", bson.D{
				{"$map", bson.D{
					{"input", bson.A{
						bson.D{
							{"foodname", "$_id.name"},
							{"Unit_price", "$Unit_price"},
							{"food_image", "$food_image"},
							{"quantity", "$quantity"},
							//{"table_number", "$table_number"},
							//{"table_id", "$table_id"},
						},
					}},
					{"as", "item"},
					{"in", "$$item"},
				}},
			}},
		}},
	}

	groupstageFinal := bson.D{
		{"$group", bson.D{
			{"_id", bson.D{{"OrderId", "$OrderId"}, {"Table_number", "$Table_number"}, {"Table_id", "$Table_id"}}},
			{"total_payment_due", bson.D{{"$sum", "$payment_due"}}},
			{"order_items", bson.D{{"$push", "$order_items"}}},
		}},
	}

	projectStageFinal := bson.D{
		{"$project", bson.D{
			{"_id", 0},
			{"OrderId", "$_id.OrderId"},
			{"Table_number", "$_id.Table_number"},
			{"TableId", "$_id.Table_id"},
			{"payment_due", "$total_payment_due"},
			{"order_items", 1},
		}},
	}

	result, err := a.orderitem_collection.Aggregate(ctx, mongo.Pipeline{

		matchstage,
		foodlookup,
		unwindfood,
		orderlookup,
		unwindorder,
		tablelookup,
		unwindtable,
		projectStage,
		groupstage,
		projectStage2,
		groupstageFinal,
		projectStageFinal,
	})

	if err != nil {

		return nil, err
	}

	if err = result.All(ctx, &OrterItems); err != nil {

		return nil, err
	}

	return OrterItems, nil
}

func (a *Orderitem_db) GetOrderItems(c *gin.Context) {

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

	defer cancel()

	results, err := a.orderitem_collection.Find(ctx, bson.M{})

	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var Allitems []bson.M
	for results.Next(ctx) {

		var item bson.M

		if err := results.Decode(&item); err != nil {

			log.Fatal(err)
		}

		Allitems = append(Allitems, item)
	}

	c.JSON(http.StatusOK, Allitems)
}

func (a *Orderitem_db) GetorderItem(c *gin.Context) {

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	order_itemId := c.Param("order_item_id")
	var orderitem model.OrderItem

	err := a.orderitem_collection.FindOne(ctx, bson.M{"order_item_id": order_itemId}).Decode(&orderitem)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, orderitem)

}

func (a *Orderitem_db) UpdateOrderitem(c *gin.Context) {

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	order_item_id := c.Param("order_item_id")
	var orderitem model.OrderItem

	if err := c.ShouldBindJSON(&orderitem); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if valerr := a.validate.Struct(orderitem); valerr != nil {

		c.JSON(http.StatusBadRequest, gin.H{"error": valerr.Error()})
		return
	}

	var update_order_item primitive.D
	if orderitem.Quantity != nil {

		update_order_item = append(update_order_item, bson.E{"quantity", orderitem.Quantity})
	}

	if orderitem.Food_id != nil {

		update_order_item = append(update_order_item, bson.E{"food_id", orderitem.Food_id})
	}

	orderitem.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	update_order_item = append(update_order_item, bson.E{"updated_at", orderitem.Updated_at})

	upsert := true
	opt := options.UpdateOptions{

		Upsert: &upsert,
	}

	filter := bson.M{"order_item_id": order_item_id}

	update_status, err := a.orderitem_collection.UpdateOne(

		ctx,
		filter,
		bson.D{

			{"$set", update_order_item},
		},
		&opt,
	)

	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, update_status)
}
