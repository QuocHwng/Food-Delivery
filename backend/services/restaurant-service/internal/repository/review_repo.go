package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"restaurant-service/internal/model"
)

type ReviewRepository struct {
	db *sqlx.DB
}

func NewReviewRepository(db *sqlx.DB) *ReviewRepository {
	return &ReviewRepository{db: db}
}

// CreateReview inserts a review and its items transactionally, then updates averages
func (r *ReviewRepository) CreateReview(ctx context.Context, review *model.Review) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Insert Review
	qReview := `INSERT INTO restaurants.reviews (id, restaurant_id, order_id, customer_id, score, comment, image_url)
		VALUES (:id, :restaurant_id, :order_id, :customer_id, :score, :comment, :image_url)`
	if _, err := tx.NamedExecContext(ctx, qReview, review); err != nil {
		return err
	}

	// Insert Items
	if len(review.Items) > 0 {
		qItem := `INSERT INTO restaurants.review_items (id, review_id, menu_item_id, score)
			VALUES (:id, :review_id, :menu_item_id, :score)`
		for _, item := range review.Items {
			if _, err := tx.NamedExecContext(ctx, qItem, item); err != nil {
				return err
			}
		}
	}

	// Update Restaurant Avg
	qUpdateRest := `UPDATE restaurants.restaurants 
		SET total_ratings = total_ratings + 1,
		    avg_rating = ((avg_rating * total_ratings) + $1) / (total_ratings + 1)
		WHERE id = $2`
	if _, err := tx.ExecContext(ctx, qUpdateRest, review.Score, review.RestaurantID); err != nil {
		return err
	}

	// Update Menu Items Avg
	for _, item := range review.Items {
		qUpdateItem := `UPDATE restaurants.menu_items 
			SET total_ratings = total_ratings + 1,
			    avg_rating = ((avg_rating * total_ratings) + $1) / (total_ratings + 1)
			WHERE id = $2`
		if _, err := tx.ExecContext(ctx, qUpdateItem, item.Score, item.MenuItemID); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *ReviewRepository) GetByRestaurant(ctx context.Context, restID uuid.UUID) ([]model.Review, error) {
	var list []model.Review
	query := `SELECT * FROM restaurants.reviews WHERE restaurant_id = $1 ORDER BY created_at DESC LIMIT 50`
	err := r.db.SelectContext(ctx, &list, query, restID)
	return list, err
}

func (r *ReviewRepository) CheckOrderReviewed(ctx context.Context, orderID uuid.UUID) (bool, error) {
	var count int
	err := r.db.GetContext(ctx, &count, `SELECT COUNT(*) FROM restaurants.reviews WHERE order_id = $1`, orderID)
	return count > 0, err
}

func (r *ReviewRepository) Reply(ctx context.Context, reviewID uuid.UUID, reply string) error {
	query := `UPDATE restaurants.reviews SET owner_reply = $1, owner_reply_at = NOW(), updated_at = NOW() WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, reply, reviewID)
	return err
}

func (r *ReviewRepository) GetByID(ctx context.Context, reviewID uuid.UUID) (*model.Review, error) {
	var rv model.Review
	err := r.db.GetContext(ctx, &rv, `SELECT * FROM restaurants.reviews WHERE id = $1`, reviewID)
	return &rv, err
}
