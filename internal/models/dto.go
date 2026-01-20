package models

// RegisterRequest represents registration request body
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Username string `json:"username" binding:"required,min=3"`
	Password string `json:"password" binding:"required,min=6"`
}

// LoginRequest represents login request body
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse represents login response
type LoginResponse struct {
	AccessToken string `json:"access_token"`
}

// SalaryRequest represents salary update request
type SalaryRequest struct {
	Salary float64 `json:"salary" binding:"required,gt=0"`
}

// CommitmentsRequest represents commitments request body
type CommitmentsRequest struct {
	Commitments []CommitmentInput `json:"commitments" binding:"required,dive"`
}

// CommitmentInput represents commitment input
type CommitmentInput struct {
	Name  string         `json:"name" binding:"required"`
	Type  CommitmentType `json:"type" binding:"required,oneof=decimal percentage"`
	Value float64        `json:"value" binding:"required,gt=0"`
}

// CommitmentPaidRequest represents commitment paid status update
type CommitmentPaidRequest struct {
	IsPaid bool `json:"is_paid"`
}

// MonthlySummaryResponse represents monthly summary response
type MonthlySummaryResponse struct {
	Salary           float64                  `json:"salary"`
	TotalCommitment  float64                  `json:"total_commitment"`
	RemainingBalance float64                  `json:"remaining_balance"`
	Commitments      []CommitmentSummary      `json:"commitments"`
}

// CommitmentSummary represents commitment in summary
type CommitmentSummary struct {
	ID     string  `json:"id"`
	Name   string  `json:"name"`
	Amount float64 `json:"amount"`
	IsPaid bool    `json:"is_paid"`
}

// MonthlyBreakdown represents monthly breakdown in yearly summary
type MonthlyBreakdown struct {
	Month     int     `json:"month"`
	Remaining float64 `json:"remaining"`
}

// YearlySummaryResponse represents yearly summary response
type YearlySummaryResponse struct {
	Year             int                `json:"year"`
	TotalSalary      float64            `json:"total_salary"`
	TotalCommitment  float64            `json:"total_commitment"`
	TotalRemaining   float64            `json:"total_remaining"`
	MonthlyBreakdown []MonthlyBreakdown `json:"monthly_breakdown"`
}

// UpdateUserRequest represents admin user update request
type UpdateUserRequest struct {
	Email         *string  `json:"email,omitempty" binding:"omitempty,email"`
	DefaultSalary *float64 `json:"default_salary,omitempty" binding:"omitempty,gt=0"`
}

// ErrorResponse represents error response
type ErrorResponse struct {
	Error string `json:"error"`
}
