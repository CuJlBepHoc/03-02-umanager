package linkgrpc

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"gitlab.com/robotomize/gb-golang/homework/03-02-umanager/internal/database"
	"gitlab.com/robotomize/gb-golang/homework/03-02-umanager/pkg/pb"
)

var _ pb.LinkServiceServer = (*Handler)(nil)

func New(linksRepository linksRepository, timeout time.Duration) *Handler {
	return &Handler{linksRepository: linksRepository, timeout: timeout}
}

type Handler struct {
	pb.UnimplementedLinkServiceServer
	linksRepository linksRepository
	timeout         time.Duration
}

func (h Handler) GetLinkByUserID(ctx context.Context, id *pb.GetLinksByUserId) (*pb.ListLinkResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, h.timeout)
	defer cancel()

	links, err := h.linksRepository.FindByUserID(ctx, id.UserId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var pbLinks []*pb.Link
	for _, link := range links {
		pbLinks = append(pbLinks, &pb.Link{
			Id:     link.ID.Hex(),
			Title:  link.Title,
			Url:    link.URL,
			Images: link.Images,
		})
	}

	return &pb.ListLinkResponse{Links: pbLinks}, nil
}

// func (h Handler) mustEmbedUnimplementedLinkServiceServer() {
// 	// Do nothing since we are embedding the unimplemented server.
// }

func (h Handler) CreateLink(ctx context.Context, request *pb.CreateLinkRequest) (*pb.Empty, error) {
	ctx, cancel := context.WithTimeout(ctx, h.timeout)
	defer cancel()

	_, err := h.linksRepository.Create(ctx, database.CreateLinkReq{
		UserID: request.UserId,
		Title:  request.Title,
		URL:    request.Url,
		Images: request.Images,
		Tags:   request.Tags,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.Empty{}, nil
}

func (h Handler) GetLink(ctx context.Context, request *pb.GetLinkRequest) (*pb.Link, error) {
	ctx, cancel := context.WithTimeout(ctx, h.timeout)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(request.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	l, err := h.linksRepository.FindByID(ctx, objectID)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.Link{
		Id:        l.ID.Hex(),
		Title:     l.Title,
		Url:       l.URL,
		Images:    l.Images,
		Tags:      l.Tags,
		UserId:    l.UserID,
		CreatedAt: l.CreatedAt.Format(time.RFC3339),
		UpdatedAt: l.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func (h Handler) UpdateLink(ctx context.Context, request *pb.UpdateLinkRequest) (*pb.Empty, error) {
	ctx, cancel := context.WithTimeout(ctx, h.timeout)
	defer cancel()

	id, err := primitive.ObjectIDFromHex(request.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid link ID")
	}

	_, err = h.linksRepository.Update(ctx, database.UpdateLinkReq{
		ID:     id,
		UserID: request.UserId,
		Title:  request.Title,
		URL:    request.Url,
		Images: request.Images,
		Tags:   request.Tags,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.Empty{}, nil
}

func (h Handler) DeleteLink(ctx context.Context, request *pb.DeleteLinkRequest) (*pb.Empty, error) {
	ctx, cancel := context.WithTimeout(ctx, h.timeout)
	defer cancel()

	id, err := primitive.ObjectIDFromHex(request.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid link ID")
	}

	err = h.linksRepository.Delete(ctx, id)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.Empty{}, nil
}

func (h Handler) ListLinks(ctx context.Context, request *pb.Empty) (*pb.ListLinkResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, h.timeout)
	defer cancel()

	links, err := h.linksRepository.FindAll(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var pbLinks []*pb.Link
	for _, link := range links {
		pbLinks = append(pbLinks, &pb.Link{
			Id:        link.ID.Hex(),
			Title:     link.Title,
			Url:       link.URL,
			Images:    link.Images,
			Tags:      link.Tags,
			UserId:    link.UserID,
			CreatedAt: link.CreatedAt.Format(time.RFC3339),
			UpdatedAt: link.UpdatedAt.Format(time.RFC3339),
		})
	}

	return &pb.ListLinkResponse{Links: pbLinks}, nil
}