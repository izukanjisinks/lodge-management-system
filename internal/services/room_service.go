package services

import (
	"errors"
	"time"

	"lodge-system/internal/models"
	"lodge-system/internal/repository"

	"github.com/google/uuid"
)

type RoomService struct {
	repo *repository.RoomRepository
}

func NewRoomService(repo *repository.RoomRepository) *RoomService {
	return &RoomService{repo: repo}
}

func (s *RoomService) Create(room *models.Room) error {
	if room.Name == "" {
		return errors.New("room name is required")
	}
	if !models.ValidRoomTypes[room.Type] {
		return errors.New("invalid room type: must be single, double, suite, cabin, or conference")
	}
	if room.Capacity <= 0 {
		return errors.New("capacity must be greater than 0")
	}
	if room.PricePerNight < 0 {
		return errors.New("price per night cannot be negative")
	}
	if room.Amenities == nil {
		room.Amenities = []string{}
	}
	room.IsAvailable = true
	return s.repo.Create(room)
}

func (s *RoomService) GetByID(id uuid.UUID) (*models.Room, error) {
	return s.repo.GetByID(id)
}

func (s *RoomService) List(roomType string, isAvailable *bool, page, pageSize int) ([]models.Room, int, error) {
	if roomType != "" && !models.ValidRoomTypes[roomType] {
		return nil, 0, errors.New("invalid room type filter")
	}
	return s.repo.List(roomType, isAvailable, page, pageSize)
}

func (s *RoomService) ListAvailable(checkIn, checkOut time.Time, roomType string) ([]models.Room, error) {
	if checkIn.IsZero() || checkOut.IsZero() {
		return nil, errors.New("check_in and check_out are required")
	}
	if !checkOut.After(checkIn) {
		return nil, errors.New("check_out must be after check_in")
	}
	if roomType != "" && !models.ValidRoomTypes[roomType] {
		return nil, errors.New("invalid room type filter")
	}
	return s.repo.ListAvailable(checkIn, checkOut, roomType)
}

func (s *RoomService) Update(id uuid.UUID, updates *models.Room) (*models.Room, error) {
	room, err := s.repo.GetByID(id)
	if err != nil {
		return nil, errors.New("room not found")
	}
	if updates.Name != "" {
		room.Name = updates.Name
	}
	if updates.Type != "" {
		if !models.ValidRoomTypes[updates.Type] {
			return nil, errors.New("invalid room type")
		}
		room.Type = updates.Type
	}
	if updates.Capacity > 0 {
		room.Capacity = updates.Capacity
	}
	if updates.PricePerNight >= 0 {
		room.PricePerNight = updates.PricePerNight
	}
	if updates.Amenities != nil {
		room.Amenities = updates.Amenities
	}
	if updates.Description != "" {
		room.Description = updates.Description
	}
	room.IsAvailable = updates.IsAvailable
	if err := s.repo.Update(room); err != nil {
		return nil, err
	}
	return s.repo.GetByID(id)
}

func (s *RoomService) UpdateImages(id uuid.UUID, images []string) (*models.Room, error) {
	_, err := s.repo.GetByID(id)
	if err != nil {
		return nil, errors.New("room not found")
	}
	if images == nil {
		images = []string{}
	}
	if err := s.repo.UpdateImages(id, images); err != nil {
		return nil, err
	}
	return s.repo.GetByID(id)
}

func (s *RoomService) SetAvailability(id uuid.UUID, available bool) error {
	_, err := s.repo.GetByID(id)
	if err != nil {
		return errors.New("room not found")
	}
	return s.repo.SetAvailability(id, available)
}

func (s *RoomService) Delete(id uuid.UUID) error {
	return s.repo.Delete(id)
}
