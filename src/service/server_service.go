package service

import (
	"context"
	"github.com/sunil206b/customer_api/src/customerpb"
	"github.com/sunil206b/customer_api/src/model"
	"go.mongodb.org/mongo-driver/bson"
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
		return nil, status.Errorf(codes.Internal, "Error occurred when inserting record in DB %v\n", err)
	}
	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, status.Errorf(codes.Internal, "Error occurred when connecting mongo id to ObjectID %v\n", err)
	}
	modelCustomer.ID = oid
	modelCustomer.CustomerID = oid.Hex()
	filter := bson.M{"_id": oid}
	_, err = s.collection.ReplaceOne(ctx, filter, modelCustomer)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Cannot update the customer with specified Id: %v\n", err)
	}
	modelCustomer.CopyFromModelCustomer(customer)
	return &customerpb.CustomerResponse{
		Customer: customer,
	}, nil
}

// GetCustomer GRPC method implemented on Server to Get the customer from DB
func (s *Server) GetCustomer(ctx context.Context, req *customerpb.GetCustomerRequest) (*customerpb.GetCustomerResponse, error) {
	log.Println("Inside GetCustomer GRPC server service...")
	customerId := req.GetCustomerId()
	oid, err := primitive.ObjectIDFromHex(customerId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Error occurred while converting Hex customer id to mongo object id %v\n", err)
	}
	customer := &model.Customer{}
	filter := bson.M{"_id": oid}
	result := s.collection.FindOne(ctx, filter)
	if err = result.Decode(customer); err != nil {
		return nil, status.Errorf(codes.NotFound, "Cannot find the customer with specified Id: %v\n", err)
	}
	pbCustomer := &customerpb.Customer{}
	customer.CopyFromModelCustomer(pbCustomer)
	return &customerpb.GetCustomerResponse{
		Customer: pbCustomer,
	}, nil
}

// UpdateCustomer GRPC method implemented on Server to Update the customer in DB
func (s *Server) UpdateCustomer(ctx context.Context, req *customerpb.UpdateCustomerRequest) (*customerpb.UpdateCustomerResponse, error) {
	log.Println("Inside UpdateCustomer GRPC server service...")
	customer := req.GetCustomer()
	oid, err := primitive.ObjectIDFromHex(customer.CustomerId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Error occurred while converting Hex customer id to mongo object id %v\n", err)
	}
	filter := bson.M{"_id": oid}
	result := s.collection.FindOne(ctx, filter)
	modelCustomer := &model.Customer{}
	if err = result.Decode(modelCustomer); err != nil {
		return nil, status.Errorf(codes.NotFound, "Cannot find the customer with specified Id: %v\n", err)
	}
	modelCustomer.CopyToModelCustomer(customer)
	_, err = s.collection.ReplaceOne(ctx, filter, modelCustomer)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Cannot update the customer with specified Id: %v\n", err)
	}
	return &customerpb.UpdateCustomerResponse{
		Customer: customer,
	}, nil
}

// DeleteCustomer GRPC method implemented on Server to Delete the customer from DB
func (s *Server) DeleteCustomer(ctx context.Context, req *customerpb.DeleteCustomerRequest) (*customerpb.DeleteCustomerResponse, error) {
	log.Println("Inside DeleteCustomer GRPC server service...")
	customerId := req.GetCustomerId()
	oid, err := primitive.ObjectIDFromHex(customerId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Error occurred while converting Hex customer id to mongo object id %v\n", err)
	}
	filter := bson.M{"_id": oid}
	res, err := s.collection.DeleteOne(ctx, filter)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error occurred while deleting customer with id %v\n", err)
	}
	if res.DeletedCount == 0 {
		return nil, status.Errorf(codes.NotFound, "Customer not found with id %v\n", err)
	}
	return &customerpb.DeleteCustomerResponse{
		CustomerId: customerId,
	}, nil
}

// ListCustomers GRPC method implemented on Server to Get all the customer in DB
func (s *Server) ListCustomers(ctx context.Context, req *customerpb.ListCustomersRequest) (*customerpb.ListCustomersResponse, error) {
	log.Println("Inside ListCustomers GRPC server service...")
	cur, err := s.collection.Find(ctx, primitive.D{})
	if err != nil || cur.Err() != nil {
		return nil, status.Errorf(codes.Internal, "Error occurred while getting all customers from DB %v\n", err)
	}
	defer cur.Close(ctx)
	customersList := make([]*customerpb.Customer, 0)
	for cur.Next(ctx) {
		customer := &model.Customer{}
		err = cur.Decode(customer)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "Error occurred while decoding to customer %v\n", err)
		}
		pbCustomer := &customerpb.Customer{}
		customer.CopyFromModelCustomer(pbCustomer)
		customersList = append(customersList, pbCustomer)
	}
	return &customerpb.ListCustomersResponse{
		Customers: customersList,
	}, nil
}
