package main

import (
	"context"
	"fmt"
	"github.com/sunil206b/customer_api/src/customerpb"
	"google.golang.org/grpc"
	"log"
	"os"
	"strings"
	"time"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Inside Customer GRPC client...")
	port := os.Getenv("PORT")
	if port == "" {
		port = "50051"
	}
	cc, err := grpc.Dial("0.0.0.0:"+port, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to server %v\n", err)
		return
	}
	defer cc.Close()

	c := customerpb.NewCustomerServiceClient(cc)
	dateOfBirth := time.Date(
		1987, 10, 20, 20, 34, 58, 651387237, time.UTC)
	customer := &customerpb.CustomerRequest{
		Customer: &customerpb.Customer{
			CustomerId: "",
			FirstName: "Rahane",
			LastName: "Ajinkya",
			Title: "Mr.",
			Email: "rahane@bcc.com",
			HomePhone: "9848223344",
			BusinessPhone: "8000123123",
			GenderCode: "Male",
			DateOfBirth: dateOfBirth.Unix(),
		},
	}
	res, err := c.CreateCustomer(context.Background(), customer)
	if err != nil {
		log.Println(err)
	}
	fmt.Printf("Customer Created %v\n", res.GetCustomer())
	fmt.Printf("%s\n", strings.Repeat("=", 100))
	updateCust := &customerpb.Customer{
		CustomerId: "5f6c99c9667ca89f769c3534",
		FirstName: "Mahendhar Sing",
		LastName: "Dhoni",
		Title: "Mr.",
		Email: "dhoni@bcc.com",
		HomePhone: "9988223645",
		BusinessPhone: "8000123123",
		GenderCode: "Male",
		DateOfBirth: dateOfBirth.Unix(),
	}
	updateReq := &customerpb.UpdateCustomerRequest{
		Customer: updateCust,
	}
	updateRes, err := c.UpdateCustomer(context.Background(), updateReq)
	if err != nil {
		log.Println(err)
	}
	fmt.Printf("Updated Customer: %v\n", updateRes)

	/*fmt.Printf("%s\n", strings.Repeat("=", 100))
	deleteReq := &customerpb.DeleteCustomerRequest{
		CustomerId: "5f6c9bb9667ca89f769c3535",
	}
	deleteRes, err := c.DeleteCustomer(context.Background(), deleteReq)
	if err != nil {
		log.Println(err)
	}
	fmt.Printf("Deleted Customer Id: %v\n", deleteRes) */
	fmt.Printf("%s\n", strings.Repeat("=", 100))
	getCustReq := &customerpb.GetCustomerRequest{
		CustomerId: "5f6c9782667ca89f769c3533",
	}
	cust, err := c.GetCustomer(context.Background(), getCustReq)
	if err != nil {
		log.Println(err)
	}

	fmt.Printf("Customer From GetReq: %v\n", cust)

	fmt.Printf("%s\n", strings.Repeat("=", 100))
	listReq := &customerpb.ListCustomersRequest{}
	result, err := c.ListCustomers(context.Background(), listReq)
	if err != nil {
		log.Println(err)
	}
	for i, r := range result.GetCustomers() {
		fmt.Printf("Record %d: %v\n", i, r)
	}
}
