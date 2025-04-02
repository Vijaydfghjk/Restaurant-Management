package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

/*
validate:"eq=CARD|eq=CASH|eq="
This rule means that Payment_method can only have one of these three values

"CARD"

"CASH"

Empty ("")

If the value is something else, like "CHECK" or "BANK", validation will fail.
*/

type Invoice struct {
	ID               primitive.ObjectID `bson:"_id"`
	Invoice_id       string             `json:"invoice_id"`
	Order_id         string             `json:"order_id" validate:"required"`
	Payment_method   *string            `json:"payment_method" validate:"eq=CARD|eq=CASH|eq="`
	Payment_status   *string            `json:"payment_status" validate:"required,eq=PENDING|eq=PAID"` // eq=PENDING|eq=PAID: This field can only have one of two values
	Payment_due_date time.Time          `json:"Payment_due_date"`
	Created_at       time.Time          `json:"created_at"`
	Updated_at       time.Time          `json:"updated_at"`
}
