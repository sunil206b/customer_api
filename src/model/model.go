package model

import (
	"github.com/sunil206b/customer_api/src/customerpb"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// Customer will hold the customer information and passes it to the MongoDB
type Customer struct {
	ID            primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	CustomerID    string             `json:"customerId" bson:"customer_id,omitempty"`
	FirstName     string             `json:"firstName" bson:"first_name"`
	LastName      string             `json:"lastName" bson:"last_name"`
	MiddleName    string             `json:"middleName" bson:"middle_name"`
	Title         string             `json:"title" bson:"title"`
	Email         string             `json:"email" bson:"email"`
	HomePhone     string             `json:"homePhone" bson:"home_phone"`
	BusinessPhone string             `json:"businessPhone" bson:"business_phone"`
	GenderCode    string             `json:"genderCode" bson:"gender_code"`
	DateOfBirth   time.Time          `json:"dateOfBirth" bson:"date_of_birth"`
	CreatedAt     time.Time          `json:"createAt" bson:"created_at"`
	UpdatedAt     time.Time          `json:"updateAt" bson:"updated_at"`
}

// CopyToModelCustomer method will copy the Customer details from model Customer to protobuf Customer
func (customer *Customer) CopyToModelCustomer(cust *customerpb.Customer) {
	customer.CustomerID = cust.CustomerId
	customer.FirstName = cust.FirstName
	customer.LastName = cust.LastName
	customer.MiddleName = cust.MiddleName
	customer.Title = cust.Title
	customer.Email = cust.Email
	customer.HomePhone = cust.HomePhone
	customer.BusinessPhone = cust.BusinessPhone
	customer.GenderCode = cust.GenderCode
	if cust.DateOfBirth != 0 {
		customer.DateOfBirth = time.Unix(int64(cust.DateOfBirth), 0)
	}
	if customer.CreatedAt.IsZero() {
		customer.CreatedAt = time.Now()
	}
	customer.UpdatedAt = time.Now()
}

// CopyFromModelCustomer method will copy the Customer details from protobuf Customer to model Customer
func (customer *Customer) CopyFromModelCustomer(cust *customerpb.Customer) {
	cust.CustomerId = customer.CustomerID
	cust.FirstName = customer.FirstName
	cust.LastName = customer.LastName
	cust.MiddleName = customer.MiddleName
	cust.Title = customer.Title
	cust.Email = customer.Email
	cust.HomePhone = customer.HomePhone
	cust.BusinessPhone = customer.BusinessPhone
	cust.GenderCode = customer.GenderCode
	if !customer.DateOfBirth.IsZero() {
		cust.DateOfBirth = customer.DateOfBirth.Unix()
	}
	cust.CreatedAt = customer.CreatedAt.Unix()
	cust.UpdatedAt = customer.UpdatedAt.Unix()
}
