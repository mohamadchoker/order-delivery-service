package grpc

import (
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	" github.com/mohamadchoker/order-delivery-service/internal/domain"
	pb " github.com/mohamadchoker/order-delivery-service/proto"
)

// Proto to Domain conversions

func protoToAddress(p *pb.Address) domain.Address {
	if p == nil {
		return domain.Address{}
	}
	return domain.Address{
		Street:     p.Street,
		City:       p.City,
		State:      p.State,
		PostalCode: p.PostalCode,
		Country:    p.Country,
		Latitude:   p.Latitude,
		Longitude:  p.Longitude,
	}
}

func protoStatusToDomain(s pb.DeliveryStatus) domain.DeliveryStatus {
	switch s {
	case pb.DeliveryStatus_PENDING:
		return domain.DeliveryStatusPending
	case pb.DeliveryStatus_ASSIGNED:
		return domain.DeliveryStatusAssigned
	case pb.DeliveryStatus_PICKED_UP:
		return domain.DeliveryStatusPickedUp
	case pb.DeliveryStatus_IN_TRANSIT:
		return domain.DeliveryStatusInTransit
	case pb.DeliveryStatus_DELIVERED:
		return domain.DeliveryStatusDelivered
	case pb.DeliveryStatus_FAILED:
		return domain.DeliveryStatusFailed
	case pb.DeliveryStatus_CANCELLED:
		return domain.DeliveryStatusCancelled
	default:
		return domain.DeliveryStatusPending
	}
}

// Domain to Proto conversions

func addressToProto(a domain.Address) *pb.Address {
	return &pb.Address{
		Street:     a.Street,
		City:       a.City,
		State:      a.State,
		PostalCode: a.PostalCode,
		Country:    a.Country,
		Latitude:   a.Latitude,
		Longitude:  a.Longitude,
	}
}

func domainStatusToProto(s domain.DeliveryStatus) pb.DeliveryStatus {
	switch s {
	case domain.DeliveryStatusPending:
		return pb.DeliveryStatus_PENDING
	case domain.DeliveryStatusAssigned:
		return pb.DeliveryStatus_ASSIGNED
	case domain.DeliveryStatusPickedUp:
		return pb.DeliveryStatus_PICKED_UP
	case domain.DeliveryStatusInTransit:
		return pb.DeliveryStatus_IN_TRANSIT
	case domain.DeliveryStatusDelivered:
		return pb.DeliveryStatus_DELIVERED
	case domain.DeliveryStatusFailed:
		return pb.DeliveryStatus_FAILED
	case domain.DeliveryStatusCancelled:
		return pb.DeliveryStatus_CANCELLED
	default:
		return pb.DeliveryStatus_UNSPECIFIED
	}
}

func deliveryToProto(d *domain.DeliveryAssignment) *pb.DeliveryAssignment {
	proto := &pb.DeliveryAssignment{
		Id:                    d.ID.String(),
		OrderId:               d.OrderID,
		Status:                domainStatusToProto(d.Status),
		PickupAddress:         addressToProto(d.PickupAddress),
		DeliveryAddress:       addressToProto(d.DeliveryAddress),
		ScheduledPickupTime:   timestamppb.New(d.ScheduledPickupTime),
		EstimatedDeliveryTime: timestamppb.New(d.EstimatedDeliveryTime),
		Notes:                 d.Notes,
		CreatedAt:             timestamppb.New(d.CreatedAt),
		UpdatedAt:             timestamppb.New(d.UpdatedAt),
	}

	if d.DriverID != nil {
		proto.DriverId = *d.DriverID
	}

	if d.ActualPickupTime != nil {
		proto.ActualPickupTime = timestamppb.New(*d.ActualPickupTime)
	}

	if d.ActualDeliveryTime != nil {
		proto.ActualDeliveryTime = timestamppb.New(*d.ActualDeliveryTime)
	}

	return proto
}

// Error handling

func handleError(err error) error {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, domain.ErrInvalidInput):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, domain.ErrInvalidStatusTransition):
		return status.Error(codes.FailedPrecondition, err.Error())
	case errors.Is(err, domain.ErrAlreadyExists):
		return status.Error(codes.AlreadyExists, err.Error())
	default:
		return status.Error(codes.Internal, "internal server error")
	}
}
