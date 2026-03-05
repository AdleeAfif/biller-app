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

// MonthlySummaryResponse represents monthly summary response. It now
// includes totals for paid vs remaining commitments and separates the
// commitment list into paid/unpaid, merging defaults and month-specific
// entries.
type MonthlySummaryResponse struct {
	Salary                       float64             `json:"salary"`
	TotalPaidCommitment          float64             `json:"total_paid_commitment"`
	TotalRemainingCommitment     float64             `json:"total_remaining_commitment"`
	SalaryMinusPaidCommitment    float64             `json:"salary_minus_paid_commitment"`
	TotalOverallCommitment       float64             `json:"total_overall_commitment"`
	SalaryMinusOverallCommitment float64             `json:"salary_minus_overall_commitment"`
	PaidCommitments              []CommitmentSummary `json:"paid_commitments"`
	UnpaidCommitments            []CommitmentSummary `json:"unpaid_commitments"`
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
// (allows changing email, username, salary)
type UpdateUserRequest struct {
	Email         *string  `json:"email,omitempty" binding:"omitempty,email"`
	Username      *string  `json:"username,omitempty" binding:"omitempty,min=3"`
	DefaultSalary *float64 `json:"default_salary,omitempty" binding:"omitempty,gt=0"`
}

// ErrorResponse represents error response
type ErrorResponse struct {
	Error string `json:"error"`
}
