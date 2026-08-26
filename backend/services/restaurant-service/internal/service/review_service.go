package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"restaurant-service/internal/model"
	"restaurant-service/internal/repository"
)

type ReviewService struct {
	repo     *repository.ReviewRepository
	restRepo *repository.RestaurantRepository
}

func NewReviewService(repo *repository.ReviewRepository, restRepo *repository.RestaurantRepository) *ReviewService {
	return &ReviewService{repo: repo, restRepo: restRepo}
}

func (s *ReviewService) CreateReview(ctx context.Context, customerID, orderID uuid.UUID, req model.CreateReviewReq) error {
	// Check if already reviewed
	reviewed, err := s.repo.CheckOrderReviewed(ctx, orderID)
	if err != nil {
		return err
	}
	if reviewed {
		return errors.New("this order has already been reviewed")
	}

	review := &model.Review{
		ID:           uuid.New(),
		RestaurantID: req.RestaurantID,
		OrderID:      orderID,
		CustomerID:   customerID,
		Score:        req.Score,
	}

	if req.Comment != "" {
		review.Comment = &req.Comment
	}
	if req.ImageURL != "" {
		review.ImageURL = &req.ImageURL
	}

	for _, item := range req.Items {
		review.Items = append(review.Items, model.ReviewItem{
			ID:         uuid.New(),
			ReviewID:   review.ID,
			MenuItemID: item.MenuItemID,
			Score:      item.Score,
		})
	}

	return s.repo.CreateReview(ctx, review)
}

func (s *ReviewService) GetRestaurantReviews(ctx context.Context, restID uuid.UUID) ([]model.Review, error) {
	return s.repo.GetByRestaurant(ctx, restID)
}

func (s *ReviewService) ReplyToReview(ctx context.Context, ownerID, reviewID uuid.UUID, reply string) error {
	rv, err := s.repo.GetByID(ctx, reviewID)
	if err != nil {
		return errors.New("review not found")
	}
	
	// Verify owner
	rest, err := s.restRepo.FindByID(ctx, rv.RestaurantID)
	if err != nil || rest.OwnerID != ownerID {
		return errors.New("unauthorized to reply to this review")
	}

	return s.repo.Reply(ctx, reviewID, reply)
}
