package controlles

import (
	"context"
	"fmt"
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

var invoicecollection *mongo.Collection = database.OpenCollection(database.Client, "invoice")

func GetInvoices() gin.HandlerFunc {

	return func(c *gin.Context) {

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		result, err := invoicecollection.Find(ctx, bson.M{})
		defer cancel()
		if err != nil {
			masg := fmt.Sprintf("Upable to get the data")
			c.JSON(http.StatusInternalServerError, gin.H{"Error": masg})
			return
		}

		var allinvoice []bson.M

		if err = result.All(ctx, &allinvoice); err != nil {

			log.Fatal(err)
		}

		c.JSON(http.StatusOK, allinvoice)
	}

}

func GetInvoice() gin.HandlerFunc {

	return func(c *gin.Context) {

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		invoice_id := c.Param("invoice_id")

		var invoice models.Invoice

		err := invoicecollection.FindOne(ctx, bson.M{"invoice_id": invoice_id}).Decode(&invoice)

		defer cancel()

		if err != nil {

			c.JSON(http.StatusBadRequest, gin.H{"error": "error occured while listing invoice item"})
			return
		}

		var invoice_view InvoiceViewFormat

		allorder_items, err := ItemsByOrder(invoice.Order_id)

		invoice_view.Order_id = invoice.Order_id
		invoice_view.Payment_due_date = invoice.Payment_due_date

		invoice_view.Payment_method = "null"

		if invoice.Payment_method != nil {

			invoice_view.Payment_method = *invoice.Payment_method
		}

		log.Println("check", allorder_items)
		log.Println("check 1", allorder_items[0])

		invoice_view.Invoice_id = invoice_id
		invoice_view.Payment_status = *&invoice.Payment_status
		invoice_view.Payment_due = allorder_items[0]["payment_due"]
		invoice_view.Table_number = allorder_items[0]["table_number"]
		invoice_view.Order_details = allorder_items[0]["order_items"]

		c.JSON(http.StatusOK, invoice_view)
	}
}

func CreateInvoice() gin.HandlerFunc {

	return func(c *gin.Context) {

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		var invoice models.Invoice

		if err := c.ShouldBindJSON(&invoice); err != nil {

			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var order models.Order
		err := orderItemCollection.FindOne(ctx, bson.M{"order_id": invoice.Order_id}).Decode(&order)

		if err != nil {

			msg := fmt.Sprintf("message: Order was not found")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
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

		validation := validate.Struct(&invoice)

		if validation != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validation.Error()})
			return
		}

		result, insertErr := invoicecollection.InsertOne(ctx, invoice)

		if insertErr != nil {
			msg := fmt.Sprintf("invoice item was not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		defer cancel()

		c.JSON(http.StatusOK, result)
	}
}

func UpdateInvoice() gin.HandlerFunc {

	return func(c *gin.Context) {

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var invoive models.Invoice

		invoice_id := c.Param("invoice_id")
		if err := c.ShouldBindJSON(&invoive); err != nil {

			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		myfilter := bson.M{"invoice_id": invoice_id}

		var updateobje primitive.D

		if invoive.Payment_method != nil {

			updateobje = append(updateobje, bson.E{"payment_method", invoive.Payment_method})
		}

		if invoive.Payment_status != nil {

			updateobje = append(updateobje, bson.E{"payment_status", invoive.Payment_status})
		}

		invoive.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		updateobje = append(updateobje, bson.E{"updated_at", invoive.Updated_at})

		uppsert := true

		opt := options.UpdateOptions{

			Upsert: &uppsert,
		}

		//status := "PENDING"

		//if invoive.Payment_status == nil {

		//	invoive.Payment_status = &status
		//}
		log.Println("check", updateobje)
		result, err := invoicecollection.UpdateOne(

			ctx,
			myfilter,
			bson.D{

				{"$set", updateobje},
			},
			&opt,
		)

		defer cancel()

		if err != nil {

			c.JSON(http.StatusInternalServerError, gin.H{"error": "Cound not update "})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}
