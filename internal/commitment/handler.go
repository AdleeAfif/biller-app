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
	"go.mongodb.org/mongo-driver/mongo/options"
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

	opts := options.Update().SetUpsert(true)
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

	// Get or create monthly record (we need the existing data so we can append)
	rec, err := getOrCreateMonthlyRecord(h.db, userID, year, month)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "failed to get monthly record"})
		return
	}

	// Start with whatever commitments are already stored (this will include defaults
	// when the record was created) and append the incoming ones instead of replacing.
	commitments := rec.Commitments
	for _, input := range req.Commitments {
		commitments = append(commitments, models.Commitment{
			ID:     primitive.NewObjectID(),
			Name:   input.Name,
			Type:   input.Type,
			Value:  input.Value,
			IsPaid: false,
		})
	}

	// Update monthly record with merged commitments
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

	// ensure there is a monthly record for this month; if it didn't exist the
	// helper will create it and copy default commitments so the update below can
	// find the item by ID. We ignore the returned record as we don't need it here.
	if _, err := getOrCreateMonthlyRecord(h.db, userID, year, month); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "failed to get monthly record"})
		return
	}

	// Update the commitment's paid status
	filter := bson.M{
		"user_id":         userID,
		"year":            year,
		"month":           month,
		"commitments._id": commitmentID,
	}

	update := bson.M{
		"$set": bson.M{"commitments.$.is_paid": req.IsPaid},
	}

	monthlyResult, err := h.db.Collection("monthly_records").UpdateOne(context.Background(), filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "failed to update commitment"})
		return
	}

	defaultResult := &mongo.UpdateResult{}
	if monthlyResult.MatchedCount == 0 {
		// no match in monthly records; try default commitments
		defaultResult, err = h.db.Collection("default_commitments").UpdateOne(
			context.Background(),
			bson.M{"user_id": userID, "commitments._id": commitmentID},
			bson.M{"$set": bson.M{"commitments.$.is_paid": req.IsPaid}},
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "failed to update commitment"})
			return
		}
	} else {
		// also update default commitments entry if present (ignore error)
		_, _ = h.db.Collection("default_commitments").UpdateOne(
			context.Background(),
			bson.M{"user_id": userID, "commitments._id": commitmentID},
			bson.M{"$set": bson.M{"commitments.$.is_paid": req.IsPaid}},
		)
	}

	if monthlyResult.MatchedCount == 0 && defaultResult.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "commitment not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "commitment status updated successfully",
		"is_paid": req.IsPaid,
	})
}

// GetDefaultCommitments returns the user's saved default commitments
func (h *Handler) GetDefaultCommitments(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "unauthorized"})
		return
	}

	var def models.DefaultCommitment
	err = h.db.Collection("default_commitments").FindOne(context.Background(), bson.M{"user_id": userID}).Decode(&def)
	if err == mongo.ErrNoDocuments {
		// nothing set, use default salary and empty lists
		var user models.User
		salary := 0.0
		if err := h.db.Collection("users").FindOne(context.Background(), bson.M{"_id": userID}).Decode(&user); err == nil {
			salary = user.DefaultSalary
		}
		c.JSON(http.StatusOK, gin.H{
			"salary":                     salary,
			"paid_commitments":           []models.CommitmentSummary{},
			"unpaid_commitments":         []models.CommitmentSummary{},
			"total_paid_commitment":      0,
			"total_remaining_commitment": salary,
		})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "failed to fetch default commitments"})
		return
	}

	// get user's default salary for calculations
	var user models.User
	salary := 0.0
	if err := h.db.Collection("users").FindOne(context.Background(), bson.M{"_id": userID}).Decode(&user); err == nil {
		salary = user.DefaultSalary
	}

	// compute paid/unpaid totals
	totalPaid := 0.0
	paidList := []models.CommitmentSummary{}
	unpaidList := []models.CommitmentSummary{}
	for _, cmt := range def.Commitments {
		amount := cmt.Value
		if cmt.Type == models.CommitmentTypePercentage {
			amount = (cmt.Value / 100) * salary
		}
		item := models.CommitmentSummary{
			ID:     cmt.ID.Hex(),
			Name:   cmt.Name,
			Amount: amount,
			IsPaid: cmt.IsPaid,
		}
		if cmt.IsPaid {
			paidList = append(paidList, item)
			totalPaid += amount
		} else {
			unpaidList = append(unpaidList, item)
		}
	}
	totalRemaining := salary - totalPaid

	c.JSON(http.StatusOK, gin.H{
		"salary":                     salary,
		"paid_commitments":           paidList,
		"unpaid_commitments":         unpaidList,
		"total_paid_commitment":      totalPaid,
		"total_remaining_commitment": totalRemaining,
	})
}

// GetMonthlyCommitments returns the commitments for a given year/month
func (h *Handler) GetMonthlyCommitments(c *gin.Context) {
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

	// fetch commitments from dedicated collection (no defaults)
	var cmts []models.Commitment
	cursor, err := h.db.Collection("commitments").Find(
		context.Background(),
		bson.M{"user_id": userID, "year": year, "month": month},
	)
	if err != nil && err != mongo.ErrNoDocuments {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "failed to fetch commitments"})
		return
	}
	if cursor != nil {
		defer cursor.Close(context.Background())
		for cursor.Next(context.Background()) {
			var c models.Commitment
			if err := cursor.Decode(&c); err == nil {
				cmts = append(cmts, c)
			}
		}
	}

	// lookup salary separately (use monthly record if available, otherwise user default)
	salary := 0.0
	var rec models.MonthlyRecord
	err = h.db.Collection("monthly_records").FindOne(
		context.Background(),
		bson.M{"user_id": userID, "year": year, "month": month},
	).Decode(&rec)
	if err == nil {
		salary = rec.Salary
	} else if err == mongo.ErrNoDocuments {
		var user models.User
		if err := h.db.Collection("users").FindOne(context.Background(), bson.M{"_id": userID}).Decode(&user); err == nil {
			salary = user.DefaultSalary
		}
	} else {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "failed to fetch salary"})
		return
	}

	// calculate paid totals using cmts slice
	totalPaid := 0.0
	paidList := []models.CommitmentSummary{}
	unpaidList := []models.CommitmentSummary{}
	for _, cmt := range cmts {
		amount := cmt.Value
		if cmt.Type == models.CommitmentTypePercentage {
			amount = (cmt.Value / 100) * salary
		}
		tmp := models.CommitmentSummary{
			ID:     cmt.ID.Hex(),
			Name:   cmt.Name,
			Amount: amount,
			IsPaid: cmt.IsPaid,
		}
		if cmt.IsPaid {
			paidList = append(paidList, tmp)
			totalPaid += amount
		} else {
			unpaidList = append(unpaidList, tmp)
		}
	}
	totalRemaining := salary - totalPaid

	c.JSON(http.StatusOK, gin.H{
		"year":                       year,
		"month":                      month,
		"salary":                     salary,
		"paid_commitments":           paidList,
		"unpaid_commitments":         unpaidList,
		"total_paid_commitment":      totalPaid,
		"total_remaining_commitment": totalRemaining,
	})
}
