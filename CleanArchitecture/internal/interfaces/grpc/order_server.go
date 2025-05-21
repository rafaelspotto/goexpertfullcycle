package grpc

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/goexpertfullcycle/cleanarchitecture/internal/core/usecase"
	"github.com/goexpertfullcycle/cleanarchitecture/proto"
	"google.golang.org/grpc"
)

type OrderServer struct {
	proto.UnimplementedOrderServiceServer
	OrderUseCase *usecase.OrderUseCase
}

func NewOrderServer(useCase *usecase.OrderUseCase) *OrderServer {
	return &OrderServer{
		OrderUseCase: useCase,
	}
}

func (s *OrderServer) CreateOrder(ctx context.Context, req *proto.CreateOrderRequest) (*proto.Order, error) {
	order, err := s.OrderUseCase.Create(req.Price, req.Tax)
	if err != nil {
		return nil, err
	}

	return &proto.Order{
		Id:         order.ID,
		Price:      order.Price,
		Tax:        order.Tax,
		FinalPrice: order.FinalPrice,
		CreatedAt:  order.CreatedAt.Format(time.RFC3339),
		UpdatedAt:  order.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func (s *OrderServer) ListOrders(ctx context.Context, req *proto.ListOrdersRequest) (*proto.ListOrdersResponse, error) {
	orders, err := s.OrderUseCase.List()
	if err != nil {
		return nil, err
	}

	var protoOrders []*proto.Order
	for _, order := range orders {
		protoOrders = append(protoOrders, &proto.Order{
			Id:         order.ID,
			Price:      order.Price,
			Tax:        order.Tax,
			FinalPrice: order.FinalPrice,
			CreatedAt:  order.CreatedAt.Format(time.RFC3339),
			UpdatedAt:  order.UpdatedAt.Format(time.RFC3339),
		})
	}

	return &proto.ListOrdersResponse{
		Orders: protoOrders,
	}, nil
}

func StartGRPCServer(useCase *usecase.OrderUseCase) error {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		return err
	}

	server := grpc.NewServer()
	proto.RegisterOrderServiceServer(server, NewOrderServer(useCase))

	log.Println("Starting gRPC server on :50051")
	return server.Serve(lis)
}
