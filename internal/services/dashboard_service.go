package services

import (
	"lodge-system/internal/models"
	"lodge-system/internal/repository"
)

type DashboardService struct {
	repo *repository.DashboardRepository
}

func NewDashboardService(repo *repository.DashboardRepository) *DashboardService {
	return &DashboardService{repo: repo}
}

func (s *DashboardService) GetStaffStats() (*models.DashboardStats, error) {
	statCards, err := s.repo.StatCards()
	if err != nil {
		return nil, err
	}

	roomSummary, err := s.repo.RoomSummary()
	if err != nil {
		return nil, err
	}

	revenueByMonth, err := s.repo.RevenueByMonth(12)
	if err != nil {
		return nil, err
	}

	reservationsByDay, err := s.repo.ReservationsByDay(14)
	if err != nil {
		return nil, err
	}

	recentBookings, err := s.repo.RecentBookings(5)
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
