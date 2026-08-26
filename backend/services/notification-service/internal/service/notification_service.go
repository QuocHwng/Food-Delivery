package service

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"notification-service/internal/model"
	"notification-service/internal/repository"
	"notification-service/internal/websocket"
)

type NotificationService struct {
	repo *repository.NotificationRepository
	hub  *websocket.Hub
}

func NewNotificationService(repo *repository.NotificationRepository, hub *websocket.Hub) *NotificationService {
	return &NotificationService{
		repo: repo,
		hub:  hub,
	}
}

func (s *NotificationService) HandleEvent(event model.NotificationEvent) error {
	var dataStr *string
	if event.Data != nil {
		bytes, err := json.Marshal(event.Data)
		if err == nil {
			str := string(bytes)
			dataStr = &str
		}
	}

	notif := &model.Notification{
		ID:      uuid.New(),
		UserID:  event.UserID,
		Title:   event.Title,
		Message: event.Message,
		Type:    event.Type,
		Data:    dataStr,
	}

	// 1. Save to database
	if err := s.repo.Create(context.Background(), notif); err != nil {
		return err
	}

	// 2. Broadcast via WebSocket if user is online
	msgBytes, _ := json.Marshal(notif)
	s.hub.SendToUser(event.UserID.String(), msgBytes)

	return nil
}

func (s *NotificationService) GetMyNotifications(ctx context.Context, userID uuid.UUID) ([]model.Notification, error) {
	return s.repo.FindByUser(ctx, userID)
}

func (s *NotificationService) MarkAsRead(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	return s.repo.MarkAsRead(ctx, id, userID)
}

func (s *NotificationService) MarkAllAsRead(ctx context.Context, userID uuid.UUID) error {
	return s.repo.MarkAllAsRead(ctx, userID)
}
