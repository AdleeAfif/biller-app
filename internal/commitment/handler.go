package commitment

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nkamil/biller-app/internal/middleware"
	"github.com/nkamil/biller-app/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Handler handles commitment-related requests
type Handler struct {
	db *mongo.Database
}

// NewHandler creates a new commitment handler
func NewHandler(db *mongo.Database) *Handler {
	return &Handler{db: db}
}

// SetDefaultCommitments sets default monthly commitments
func (h *Handler) SetDefaultCommitments(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "unauthorized"})
		return
	}

	var req models.CommitmentsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	// Convert input to commitments
	commitments := make([]models.Commitment, len(req.Commitments))
	for i, input := range req.Commitments {
		commitments[i] = models.Commitment{
			ID:     primitive.NewObjectID(),
			Name:   input.Name,
			Type:   input.Type,
			Value:  input.Value,
			IsPaid: false,
		}
	}

	// Upsert default commitments
	filter := bson.M{"user_id": userID}
	update := bson.M{
		"$set": bson.M{
			"user_id":     userID,
			"commitments": commitments,
			"updated_at":  time.Now(),
		},
	}

	opts := mongo.NewUpdateOptions().SetUpsert(true)
	_, err = h.db.Collection("default_commitments").UpdateOne(
		context.Background(),
		filter,
		update,
		opts,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "failed to set default commitments"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "default commitments set successfully",
		"commitments": commitments,
	})
}

// SetMonthlyCommitments sets commitments for a specific month
func (h *Handler) SetMonthlyCommitments(c *gin.Context) {
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

	var req models.CommitmentsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	// Get or create monthly record
	record, err := getOrCreateMonthlyRecord(h.db, userID, year, month)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "failed to get monthly record"})
		return
	}

	// Convert input to commitments
	commitments := make([]models.Commitment, len(req.Commitments))
	for i, input := range req.Commitments {
		commitments[i] = models.Commitment{
			ID:     primitive.NewObjectID(),
			Name:   input.Name,
			Type:   input.Type,
			Value:  input.Value,
			IsPaid: false,
		}
	}

	// Update monthly record with new commitments
	filter := bson.M{
		"user_id": userID,
		"year":    year,
		"month":   month,
	}

	update := bson.M{"$set": bson.M{"commitments": commitments}}
	_, err = h.db.Collection("monthly_records").UpdateOne(context.Background(), filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "failed to update commitments"})
		return
	}

	// Recalculate totals
	if err := recalculateMonthlyTotals(h.db, userID, year, month); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "failed to recalculate totals"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "monthly commitments set successfully",
		"year":        year,
		"month":       month,
		"commitments": commitments,
	})
}

// UpdateCommitmentPaidStatus updates the paid status of a commitment
func (h *Handler) UpdateCommitmentPaidStatus(c *gin.Context) {
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

	commitmentIDStr := c.Param("commitment_id")
	commitmentID, err := primitive.ObjectIDFromHex(commitmentIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid commitment ID"})
		return
	}

	var req models.CommitmentPaidRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	// Update the commitment's paid status
	filter := bson.M{
		"user_id":          userID,
		"year":             year,
		"month":            month,
		"commitments._id":  commitmentID,
	}

	update := bson.M{
		"$set": bson.M{"commitments.$.is_paid": req.IsPaid},
	}

	result, err := h.db.Collection("monthly_records").UpdateOne(context.Background(), filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "failed to update commitment"})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "commitment not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "commitment status updated successfully",
		"is_paid": req.IsPaid,
	})
}
