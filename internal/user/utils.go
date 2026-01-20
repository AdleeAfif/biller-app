package user

import (
	"context"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nkamil/biller-app/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

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

// recalculateMonthlyTotals recalculates totals for a monthly record
func recalculateMonthlyTotals(db *mongo.Database, userID primitive.ObjectID, year, month int) error {
	filter := bson.M{
		"user_id": userID,
		"year":    year,
		"month":   month,
	}

	var record models.MonthlyRecord
	err := db.Collection("monthly_records").FindOne(context.Background(), filter).Decode(&record)
	if err != nil {
		return err
	}

	// Calculate total commitment
	totalCommitment := 0.0
	for _, commitment := range record.Commitments {
		amount := commitment.Value
		if commitment.Type == models.CommitmentTypePercentage {
			amount = (commitment.Value / 100) * record.Salary
		}
		totalCommitment += amount
	}

	// Update record
	update := bson.M{"$set": bson.M{
		"total_commitment":  totalCommitment,
		"remaining_balance": record.Salary - totalCommitment,
	}}

	_, err = db.Collection("monthly_records").UpdateOne(context.Background(), filter, update)
	return err
}

// getOrCreateMonthlyRecord gets or creates a monthly record
func getOrCreateMonthlyRecord(db *mongo.Database, userID primitive.ObjectID, year, month int) (*models.MonthlyRecord, error) {
	filter := bson.M{
		"user_id": userID,
		"year":    year,
		"month":   month,
	}

	var record models.MonthlyRecord
	err := db.Collection("monthly_records").FindOne(context.Background(), filter).Decode(&record)

	if err == mongo.ErrNoDocuments {
		// Get user's default salary
		var user models.User
		err := db.Collection("users").FindOne(
			context.Background(),
			bson.M{"_id": userID},
		).Decode(&user)
		if err != nil {
			return nil, err
		}

		// Create new record with default salary
		record = models.MonthlyRecord{
			UserID:           userID,
			Year:             year,
			Month:            month,
			Salary:           user.DefaultSalary,
			Commitments:      []models.Commitment{},
			TotalCommitment:  0,
			RemainingBalance: user.DefaultSalary,
		}

		_, err = db.Collection("monthly_records").InsertOne(context.Background(), record)
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	return &record, nil
}
