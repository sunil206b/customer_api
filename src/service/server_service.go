package service

import (
	"context"
	"github.com/sunil206b/customer_api/src/customerpb"
	"github.com/sunil206b/customer_api/src/model"
	"log"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// NewServer function will return the server instance
func NewServer(col *mongo.Collection) *Server {
	return &Server{
		collection: col,
	}
}

// Server will implement all the GRPC server methods
type Server struct {
	collection *mongo.Collection
}

// CreateCustomer GRPC method implemented on Server to Create new customer in DB
func (s *Server) CreateCustomer(ctx context.Context, req *customerpb.CustomerRequest) (*customerpb.CustomerResponse, error) {
	log.Println("Inside CreateCustomer GRPC service...")
	customer := req.GetCustomer()
	modelCustomer := &model.Customer{}
	modelCustomer.CopyToModelCustomer(customer)
	res, err := s.collection.InsertOne(ctx, modelCustomer)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error occured when inserting record in DB %v", err)
	}
	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, status.Errorf(codes.Internal, "Error occured when conerting mongo id to ObjectID %v", err)
	}
	modelCustomer.ID = oid
	modelCustomer.CustomerID = oid.Hex()
	modelCustomer.CopyFromModelCustomer(customer)
	return &customerpb.CustomerResponse{
		Customer: customer,
	}, nil
}

// GetCustomer GRPC method implemented on Server to Get the customer from DB
func (s *Server) GetCustomer(ctx context.Context, req *customerpb.GetCustomerRequest) (*customerpb.GetCustomerResponse, error) {
	return nil, nil
}

// UpdateCustomer GRPC method implemented on Server to Update the customer in DB
func (s *Server) UpdateCustomer(ctx context.Context, req *customerpb.UpdateCustomerRequest) (*customerpb.UpdateCustomerResponse, error) {
	return nil, nil
}

// DeleteCustomer GRPC method implemented on Server to Delete the customer from DB
func (s *Server) DeleteCustomer(ctx context.Context, req *customerpb.DeleteCustomerRequest) (*customerpb.DeleteCustomerResponse, error) {
	return nil, nil
}

// ListCustomers GRPC method implemented on Server to Get all the customer in DB
func (s *Server) ListCustomers(ctx context.Context, req *customerpb.ListCustomersRequest) (*customerpb.ListCustomersResponse, error) {
	return nil, nil
}
