package grpc

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	" github.com/mohamadchoker/order-delivery-service/internal/service"
	pb " github.com/mohamadchoker/order-delivery-service/proto"
)

// Handler implements the gRPC DeliveryService
type Handler struct {
	pb.UnimplementedDeliveryServiceServer
	useCase service.DeliveryUseCase
	logger  *zap.Logger
}

// NewHandler creates a new gRPC handler
func NewHandler(useCase service.DeliveryUseCase, logger *zap.Logger) *Handler {
	return &Handler{
		useCase: useCase,
		logger:  logger,
	}
}

// CreateDeliveryAssignment creates a new delivery assignment
func (h *Handler) CreateDeliveryAssignment(ctx context.Context, req *pb.CreateDeliveryAssignmentRequest) (*pb.DeliveryAssignment, error) {
	// Validate request
	if req.OrderId == "" {
		return nil, status.Error(codes.InvalidArgument, "order_id is required")
	}
	if req.PickupAddress == nil || req.DeliveryAddress == nil {
		return nil, status.Error(codes.InvalidArgument, "pickup_address and delivery_address are required")
	}

	// Convert proto to domain
	input := service.CreateDeliveryInput{
		OrderID:               req.OrderId,
		PickupAddress:         protoToAddress(req.PickupAddress),
		DeliveryAddress:       protoToAddress(req.DeliveryAddress),
		ScheduledPickupTime:   req.ScheduledPickupTime.AsTime(),
		EstimatedDeliveryTime: req.EstimatedDeliveryTime.AsTime(),
		Notes:                 req.Notes,
	}

	// Create delivery assignment
	assignment, err := h.useCase.CreateDeliveryAssignment(ctx, input)
	if err != nil {
		return nil, handleError(err)
	}

	return deliveryToProto(assignment), nil
}

// GetDeliveryAssignment retrieves a delivery assignment by ID
func (h *Handler) GetDeliveryAssignment(ctx context.Context, req *pb.GetDeliveryAssignmentRequest) (*pb.DeliveryAssignment, error) {
	// Parse UUID
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid id format")
	}

	// Get delivery assignment
	assignment, err := h.useCase.GetDeliveryAssignment(ctx, id)
	if err != nil {
		return nil, handleError(err)
	}

	return deliveryToProto(assignment), nil
}

// UpdateDeliveryStatus updates the status of a delivery
func (h *Handler) UpdateDeliveryStatus(ctx context.Context, req *pb.UpdateDeliveryStatusRequest) (*pb.DeliveryAssignment, error) {
	// Parse UUID
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid id format")
	}

	// Convert proto status to domain status
	domainStatus := protoStatusToDomain(req.Status)

	// Update status
	assignment, err := h.useCase.UpdateDeliveryStatus(ctx, id, domainStatus, req.Notes)
	if err != nil {
		return nil, handleError(err)
	}

	return deliveryToProto(assignment), nil
}

// ListDeliveryAssignments lists delivery assignments with pagination
func (h *Handler) ListDeliveryAssignments(ctx context.Context, req *pb.ListDeliveryAssignmentsRequest) (*pb.ListDeliveryAssignmentsResponse, error) {
	// Prepare input
	input := service.ListDeliveryInput{
		Page:     int(req.Page),
		PageSize: int(req.PageSize),
	}

	if req.Status != pb.DeliveryStatus_UNSPECIFIED {
		domainStatus := protoStatusToDomain(req.Status)
		input.Status = &domainStatus
	}

	if req.DriverId != "" {
		input.DriverID = &req.DriverId
	}

	// List assignments
	assignments, totalCount, err := h.useCase.ListDeliveryAssignments(ctx, input)
	if err != nil {
		return nil, handleError(err)
	}

	// Convert to proto
	protoAssignments := make([]*pb.DeliveryAssignment, len(assignments))
	for i, assignment := range assignments {
		protoAssignments[i] = deliveryToProto(assignment)
	}

	return &pb.ListDeliveryAssignmentsResponse{
		Assignments: protoAssignments,
		TotalCount:  int32(totalCount),
		Page:        req.Page,
		PageSize:    req.PageSize,
	}, nil
}

// AssignDriver assigns a driver to a delivery
func (h *Handler) AssignDriver(ctx context.Context, req *pb.AssignDriverRequest) (*pb.DeliveryAssignment, error) {
	// Parse UUID
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid id format")
	}

	if req.DriverId == "" {
		return nil, status.Error(codes.InvalidArgument, "driver_id is required")
	}

	// Assign driver
	assignment, err := h.useCase.AssignDriver(ctx, id, req.DriverId)
	if err != nil {
		return nil, handleError(err)
	}

	return deliveryToProto(assignment), nil
}

// GetDeliveryMetrics retrieves delivery metrics
func (h *Handler) GetDeliveryMetrics(ctx context.Context, req *pb.GetDeliveryMetricsRequest) (*pb.DeliveryMetrics, error) {
	startTime := req.StartTime.AsTime()
	endTime := req.EndTime.AsTime()

	var driverID *string
	if req.DriverId != "" {
		driverID = &req.DriverId
	}

	metrics, err := h.useCase.GetDeliveryMetrics(
		ctx,
		startTime,
		endTime,
		driverID,
	)
	if err != nil {
		return nil, handleError(err)
	}

	return &pb.DeliveryMetrics{
		TotalDeliveries:            metrics.TotalDeliveries,
		CompletedDeliveries:        metrics.CompletedDeliveries,
		FailedDeliveries:           metrics.FailedDeliveries,
		CancelledDeliveries:        metrics.CancelledDeliveries,
		AverageDeliveryTimeMinutes: metrics.AverageDeliveryTimeMinutes,
		OnTimeDeliveryRate:         metrics.OnTimeDeliveryRate,
	}, nil
}

func (h *Handler) DeleteDeliveryAssignment(ctx context.Context, req *pb.DeleteDeliveryAssignmentRequest) (*empty.Empty, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid delivery ID")
	}

	if err := h.useCase.DeleteDeliveryAssignment(ctx, id); err != nil {
		return nil, handleError(err)
	}

	return &empty.Empty{}, nil
}
