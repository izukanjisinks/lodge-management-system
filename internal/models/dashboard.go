package models

type DashboardStatCards struct {
	NewBookingsThisMonth int `json:"new_bookings_this_month"`
	CheckInsToday        int `json:"checkins_today"`
	CheckOutsToday       int `json:"checkouts_today"`
}

type DashboardRoomSummary struct {
	Occupied  int `json:"occupied"`
	Reserved  int `json:"reserved"`
	Available int `json:"available"`
	NotReady  int `json:"not_ready"`
}

type DashboardRevenuePoint struct {
	Month   string  `json:"month"`   // "2026-01"
	Revenue float64 `json:"revenue"`
}

type DashboardReservationPoint struct {
	Day       string `json:"day"`       // "2026-03-13"
	Booked    int    `json:"booked"`
	Cancelled int    `json:"cancelled"`
}

type DashboardRecentBooking struct {
	ID         string `json:"id"`
	ClientName string `json:"client_name"`
	RoomName   string `json:"room_name"`
	RoomType   string `json:"room_type"`
	CheckIn    string `json:"check_in"`
	CheckOut   string `json:"check_out"`
	Status     string `json:"status"`
}

type DashboardStats struct {
	StatCards       DashboardStatCards          `json:"stat_cards"`
	RoomSummary     DashboardRoomSummary        `json:"room_summary"`
	RevenueByMonth  []DashboardRevenuePoint     `json:"revenue_by_month"`
	ReservationsByDay []DashboardReservationPoint `json:"reservations_by_day"`
	RecentBookings  []DashboardRecentBooking    `json:"recent_bookings"`
}
