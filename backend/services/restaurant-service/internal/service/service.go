package service

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"restaurant-service/internal/model"
	"restaurant-service/internal/repository"
)

var ErrUnauthorized = errors.New("unauthorized")

type RestaurantService struct {
	repo     *repository.RestaurantRepository
	menuRepo *repository.MenuRepository
}

func NewRestaurantService(repo *repository.RestaurantRepository, menuRepo *repository.MenuRepository) *RestaurantService {
	return &RestaurantService{repo: repo, menuRepo: menuRepo}
}

// CreateRestaurant for owners
func (s *RestaurantService) CreateRestaurant(ctx context.Context, ownerID uuid.UUID, role string, req model.CreateRestaurantReq) (*model.Restaurant, error) {
	if role != "restaurant_owner" && role != "admin" {
		return nil, ErrUnauthorized
	}

	rest := &model.Restaurant{
		ID:            uuid.New(),
		OwnerID:       ownerID,
		Name:          req.Name,
		MinOrderValue: req.MinOrderValue,
		DeliveryFee:   req.DeliveryFee,
	}
	if req.Description != "" {
		rest.Description = &req.Description
	}
	if req.Address != "" {
		rest.Address = &req.Address
	}
	if req.Phone != "" {
		rest.Phone = &req.Phone
	}

	if err := s.repo.Create(ctx, rest); err != nil {
		return nil, err
	}
	return rest, nil
}

func (s *RestaurantService) ListRestaurants(ctx context.Context) ([]model.Restaurant, error) {
	return s.repo.ListAll(ctx)
}

func (s *RestaurantService) GetRestaurantDetail(ctx context.Context, id uuid.UUID) (*model.Restaurant, error) {
	return s.repo.FindByID(ctx, id)
}

// ─── Menu Methods ────────────────────────────────────────────────────────────

func (s *RestaurantService) CreateMenuCategory(ctx context.Context, ownerID uuid.UUID, restID uuid.UUID, req model.CreateMenuCategoryReq) (*model.MenuCategory, error) {
	if err := s.checkOwnership(ctx, restID, ownerID); err != nil {
		return nil, err
	}

	cat := &model.MenuCategory{
		ID:           uuid.New(),
		RestaurantID: restID,
		Name:         req.Name,
		DisplayOrder: req.DisplayOrder,
	}
	if err := s.menuRepo.CreateCategory(ctx, cat); err != nil {
		return nil, err
	}
	return cat, nil
}

func (s *RestaurantService) CreateMenuItem(ctx context.Context, ownerID uuid.UUID, restID uuid.UUID, req model.CreateMenuItemReq) (*model.MenuItem, error) {
	if err := s.checkOwnership(ctx, restID, ownerID); err != nil {
		return nil, err
	}

	item := &model.MenuItem{
		ID:           uuid.New(),
		RestaurantID: restID,
		CategoryID:   req.CategoryID,
		Name:         req.Name,
		Price:        req.Price,
	}
	if req.Description != "" {
		item.Description = &req.Description
	}
	if req.ImageURL != "" {
		item.ImageURL = &req.ImageURL
	}

	if err := s.menuRepo.CreateItem(ctx, item); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *RestaurantService) GetFullMenu(ctx context.Context, restID uuid.UUID) (*model.MenuResponse, error) {
	categories, err := s.menuRepo.GetCategoriesByRestaurant(ctx, restID)
	if err != nil {
		return nil, err
	}
	items, err := s.menuRepo.GetItemsByRestaurant(ctx, restID)
	if err != nil {
		return nil, err
	}

	resp := &model.MenuResponse{
		Categories:    make([]model.CategoryWithItems, 0),
		Uncategorized: make([]model.MenuItem, 0),
	}

	// Grouping items by category
	catMap := make(map[uuid.UUID][]model.MenuItem)
	for _, it := range items {
		if it.CategoryID != nil {
			catMap[*it.CategoryID] = append(catMap[*it.CategoryID], it)
		} else {
			resp.Uncategorized = append(resp.Uncategorized, it)
		}
	}

	for _, cat := range categories {
		c := model.CategoryWithItems{
			MenuCategory: cat,
			Items:        catMap[cat.ID],
		}
		if c.Items == nil {
			c.Items = []model.MenuItem{}
		}
		resp.Categories = append(resp.Categories, c)
	}

	return resp, nil
}

func (s *RestaurantService) checkOwnership(ctx context.Context, restID, ownerID uuid.UUID) error {
	rest, err := s.repo.FindByID(ctx, restID)
	if err != nil {
		return err
	}
	if rest.OwnerID != ownerID {
		return ErrUnauthorized
	}
	return nil
}
