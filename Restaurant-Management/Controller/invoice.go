package controller

import (
	"context"
	"log"
	"net/http"
	database "restaurant-management/Database"
	model "restaurant-management/Model"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
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

type Invoice_db struct {
	invoice_collection *mongo.Collection
	validate           *validator.Validate
}

func Invoiceccontrll() *Invoice_db {

	return &Invoice_db{
		invoice_collection: database.Open_collection(database.Client, "Invoice"),
		validate:           validator.New(),
	}
}

func (a *Invoice_db) Create_inoice(c *gin.Context) {

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var invoice model.Invoice

	if err := c.ShouldBindJSON(&invoice); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if vallerr := a.validate.Struct(invoice); vallerr != nil {

		c.JSON(http.StatusBadRequest, gin.H{"error": vallerr.Error()})
		return
	}

	invoice.ID = primitive.NewObjectID()
	invoice.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	invoice.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	invoice.Payment_due_date, _ = time.Parse(time.RFC3339, time.Now().AddDate(0, 0, 1).Format(time.RFC3339))
	invoice.Invoice_id = invoice.ID.Hex()

	status := "PENDING"

	invoice.Payment_status = &status

	inserted_data, err := a.invoice_collection.InsertOne(ctx, invoice)

	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, inserted_data)
}

func (a *Invoice_db) GetInvoices(c *gin.Context) {

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	defer cancel()

	results, err := a.invoice_collection.Find(ctx, bson.M{})

	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var invoices []bson.M
	for results.Next(ctx) {

		var invoice bson.M

		if err := results.Decode(&invoice); err != nil {

			log.Fatal(err)
		}
		invoices = append(invoices, invoice)
	}
	c.JSON(http.StatusOK, invoices)
}

func (a *Invoice_db) Getinvoice(c *gin.Context) {

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	defer cancel()

	invoice_id := c.Param("invoice_id")

	var invoice model.Invoice

	err := a.invoice_collection.FindOne(ctx, bson.M{"invoice_id": invoice_id}).Decode(&invoice)

	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	orderItemobject := Orderitemcontroll()

	mycollection, err := orderItemobject.ItemsByOrder(invoice.Order_id)
	//log.Println("value3", mycollection[0]["Table_number"])
	//log.Println("value3", mycollection)

	if err != nil {

		log.Fatal("Error is ", err)

	}

	var invice_view InvoiceViewFormat

	invice_view.Invoice_id = invoice.Invoice_id
	invice_view.Payment_method = *invoice.Payment_method
	invice_view.Order_id = invoice.Order_id
	invice_view.Payment_status = invoice.Payment_status
	invice_view.Payment_due = mycollection[0]["payment_due"]
	invice_view.Table_number = mycollection[0]["Table_number"]
	invice_view.Order_details = mycollection[0]["order_items"]

	c.JSON(http.StatusOK, invice_view)
}

func (a *Invoice_db) Update_invoice(c *gin.Context) {

	invoice_id := c.Param("invoice_id")

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	defer cancel()

	var invoice model.Invoice

	if err := c.ShouldBindJSON(&invoice); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var updated_invoice primitive.D

	if invoice.Payment_method != nil {
		updated_invoice = append(updated_invoice, bson.E{"payment_method", invoice.Payment_method})
	}

	if invoice.Payment_status != nil {

		updated_invoice = append(updated_invoice, bson.E{"payment_status", invoice.Payment_status})
	}

	invoice.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	updated_invoice = append(updated_invoice, bson.E{"updated_at", invoice.Updated_at})

	upsert := true
	opt := options.UpdateOptions{

		Upsert: &upsert,
	}

	filter := bson.M{"invoice_id": invoice_id}
	updated, err := a.invoice_collection.UpdateOne(

		ctx,
		filter,
		bson.D{

			{"$set", updated_invoice},
		},
		&opt,
	)

	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, updated)
}
