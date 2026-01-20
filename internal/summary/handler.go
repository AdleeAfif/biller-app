package summary

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nkamil/biller-app/internal/middleware"
	"github.com/nkamil/biller-app/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Handler handles summary-related requests
type Handler struct {
	db *mongo.Database
}

// NewHandler creates a new summary handler
func NewHandler(db *mongo.Database) *Handler {
	return &Handler{db: db}
}

// GetMonthlySummary returns monthly summary
func (h *Handler) GetMonthlySummary(c *gin.Context) {
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

	// Get monthly record
	var record models.MonthlyRecord
	filter := bson.M{
		"user_id": userID,
		"year":    year,
		"month":   month,
	}

	err = h.db.Collection("monthly_records").FindOne(context.Background(), filter).Decode(&record)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// No record exists, return default values
			var user models.User
			err := h.db.Collection("users").FindOne(
				context.Background(),
				bson.M{"_id": userID},
			).Decode(&user)
			if err != nil {
				c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "failed to get user"})
				return
			}

			c.JSON(http.StatusOK, models.MonthlySummaryResponse{
				Salary:           user.DefaultSalary,
				TotalCommitment:  0,
				RemainingBalance: user.DefaultSalary,
				Commitments:      []models.CommitmentSummary{},
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "failed to get monthly record"})
		return
	}

	// Build commitment summaries
	commitmentSummaries := make([]models.CommitmentSummary, len(record.Commitments))
	for i, commitment := range record.Commitments {
		amount := commitment.Value
		if commitment.Type == models.CommitmentTypePercentage {
			amount = (commitment.Value / 100) * record.Salary
		}

		commitmentSummaries[i] = models.CommitmentSummary{
			ID:     commitment.ID.Hex(),
			Name:   commitment.Name,
			Amount: amount,
			IsPaid: commitment.IsPaid,
		}
	}

	response := models.MonthlySummaryResponse{
		Salary:           record.Salary,
		TotalCommitment:  record.TotalCommitment,
		RemainingBalance: record.RemainingBalance,
		Commitments:      commitmentSummaries,
	}

	c.JSON(http.StatusOK, response)
}

// GetYearlySummary returns yearly summary
func (h *Handler) GetYearlySummary(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "unauthorized"})
		return
	}

	yearStr := c.Param("year")
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid year"})
		return
	}

	// Get all monthly records for the year
	cursor, err := h.db.Collection("monthly_records").Find(
		context.Background(),
		bson.M{
			"user_id": userID,
			"year":    year,
		},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "failed to get records"})
		return
	}
	defer cursor.Close(context.Background())

	var records []models.MonthlyRecord
	if err := cursor.All(context.Background(), &records); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "failed to decode records"})
		return
	}

	// Calculate totals
	totalSalary := 0.0
	totalCommitment := 0.0
	totalRemaining := 0.0
	monthlyBreakdown := make([]models.MonthlyBreakdown, 0)

	for _, record := range records {
		totalSalary += record.Salary
		totalCommitment += record.TotalCommitment
		totalRemaining += record.RemainingBalance

		monthlyBreakdown = append(monthlyBreakdown, models.MonthlyBreakdown{
			Month:     record.Month,
			Remaining: record.RemainingBalance,
		})
	}

	response := models.YearlySummaryResponse{
		Year:             year,
		TotalSalary:      totalSalary,
		TotalCommitment:  totalCommitment,
		TotalRemaining:   totalRemaining,
		MonthlyBreakdown: monthlyBreakdown,
	}

	c.JSON(http.StatusOK, response)
}

// parseYearMonth parses year and month from URL params
func parseYearMonth(c *gin.Context) (int, int, error) {
	yearStr := c.Param("year")
	monthStr := c.Param("month")

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return 0, 0, err
	}

	month, err := strconv.Atoi(monthStr)
	if err != nil || month < 1 || month > 12 {
		return 0, 0, err
	}

	return year, month, nil
}
