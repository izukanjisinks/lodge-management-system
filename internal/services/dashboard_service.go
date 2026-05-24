package services

import (
	"lodge-system/internal/models"
	"lodge-system/internal/repository"

	"github.com/google/uuid"
)

type DashboardService struct {
	repo *repository.DashboardRepository
}

func NewDashboardService(repo *repository.DashboardRepository) *DashboardService {
	return &DashboardService{repo: repo}
}

func (s *DashboardService) GetStaffStats(orgID uuid.UUID, branchID *uuid.UUID) (*models.DashboardStats, error) {
	statCards, err := s.repo.StatCards(orgID, branchID)
	if err != nil {
		return nil, err
	}

	roomSummary, err := s.repo.RoomSummary(orgID, branchID)
	if err != nil {
		return nil, err
	}

	revenueByMonth, err := s.repo.RevenueByMonth(orgID, branchID, 12)
	if err != nil {
		return nil, err
	}

	reservationsByDay, err := s.repo.ReservationsByDay(orgID, branchID, 14)
	if err != nil {
		return nil, err
	}

	recentBookings, err := s.repo.RecentBookings(orgID, branchID, 5)
	if err != nil {
		return nil, err
	}

	return &models.DashboardStats{
		StatCards:         statCards,
		RoomSummary:       roomSummary,
		RevenueByMonth:    revenueByMonth,
		ReservationsByDay: reservationsByDay,
		RecentBookings:    recentBookings,
	}, nil
}
