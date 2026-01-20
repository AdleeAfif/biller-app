package user

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nkamil/biller-app/internal/middleware"
	"github.com/nkamil/biller-app/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Handler handles user-related requests
type Handler struct {
	db *mongo.Database
}

// NewHandler creates a new user handler
func NewHandler(db *mongo.Database) *Handler {
	return &Handler{db: db}
}

// SetDefaultSalary sets the default salary for a user
func (h *Handler) SetDefaultSalary(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "unauthorized"})
		return
	}

	var req models.SalaryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	// Update user's default salary
	result, err := h.db.Collection("users").UpdateOne(
		context.Background(),
		bson.M{"_id": userID},
		bson.M{"$set": bson.M{"default_salary": req.Salary}},
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "failed to update salary"})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "user not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "default salary updated successfully",
		"salary":  req.Salary,
	})
}

// SetMonthlySalary sets the salary for a specific month
func (h *Handler) SetMonthlySalary(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "unauthorized"})
		return
	}

	year, month, err := parseYearMonth(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid year or month"})
		return
	}

	var req models.SalaryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	// Find or create monthly record
	filter := bson.M{
		"user_id": userID,
		"year":    year,
		"month":   month,
	}

	var record models.MonthlyRecord
	err = h.db.Collection("monthly_records").FindOne(context.Background(), filter).Decode(&record)

	if err == mongo.ErrNoDocuments {
		// Create new record
		record = models.MonthlyRecord{
			UserID:      userID,
			Year:        year,
			Month:       month,
			Salary:      req.Salary,
			Commitments: []models.Commitment{},
		}
		_, err = h.db.Collection("monthly_records").InsertOne(context.Background(), record)
	} else if err == nil {
		// Update existing record
		update := bson.M{"$set": bson.M{"salary": req.Salary}}
		_, err = h.db.Collection("monthly_records").UpdateOne(context.Background(), filter, update)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "failed to update monthly salary"})
		return
	}

	// Recalculate totals
	if err := recalculateMonthlyTotals(h.db, userID, year, month); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "failed to recalculate totals"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "monthly salary updated successfully",
		"year":    year,
		"month":   month,
		"salary":  req.Salary,
	})
}
