package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserRole represents user role type
type UserRole string

const (
	RoleUser  UserRole = "user"
	RoleAdmin UserRole = "admin"
)

// User represents a user in the system
type User struct {
	ID            primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Email         string             `json:"email" bson:"email"`
	Username      string             `json:"username" bson:"username"`
	PasswordHash  string             `json:"-" bson:"password_hash"`
	Role          UserRole           `json:"role" bson:"role"`
	DefaultSalary float64            `json:"default_salary" bson:"default_salary"`
	CreatedAt     time.Time          `json:"created_at" bson:"created_at"`
	DeletedAt     *time.Time         `json:"deleted_at,omitempty" bson:"deleted_at,omitempty"`
}

// CommitmentType represents the type of commitment
type CommitmentType string

const (
	CommitmentTypeDecimal    CommitmentType = "decimal"
	CommitmentTypePercentage CommitmentType = "percentage"
)

// Commitment represents a financial commitment
type Commitment struct {
	ID     primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name   string             `json:"name" bson:"name"`
	Type   CommitmentType     `json:"type" bson:"type"`
	Value  float64            `json:"value" bson:"value"`
	IsPaid bool               `json:"is_paid" bson:"is_paid"`
}

// MonthlyRecord represents a user's monthly financial record
type MonthlyRecord struct {
	ID               primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID           primitive.ObjectID `json:"user_id" bson:"user_id"`
	Year             int                `json:"year" bson:"year"`
	Month            int                `json:"month" bson:"month"`
	Salary           float64            `json:"salary" bson:"salary"`
	Commitments      []Commitment       `json:"commitments" bson:"commitments"`
	TotalCommitment  float64            `json:"total_commitment" bson:"total_commitment"`
	RemainingBalance float64            `json:"remaining_balance" bson:"remaining_balance"`
}

// YearlySummary represents a user's yearly summary
type YearlySummary struct {
	ID              primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID          primitive.ObjectID `json:"user_id" bson:"user_id"`
	Year            int                `json:"year" bson:"year"`
	TotalSalary     float64            `json:"total_salary" bson:"total_salary"`
	TotalCommitment float64            `json:"total_commitment" bson:"total_commitment"`
	TotalRemaining  float64            `json:"total_remaining" bson:"total_remaining"`
}

// DefaultCommitment represents a user's default monthly commitments
type DefaultCommitment struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID      primitive.ObjectID `json:"user_id" bson:"user_id"`
	Commitments []Commitment       `json:"commitments" bson:"commitments"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}
