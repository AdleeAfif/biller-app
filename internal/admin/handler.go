package admin

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nkamil/biller-app/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Handler handles admin-related requests
type Handler struct {
	db *mongo.Database
}

// NewHandler creates a new admin handler
func NewHandler(db *mongo.Database) *Handler {
	return &Handler{db: db}
}

// ListUsers returns all users
func (h *Handler) ListUsers(c *gin.Context) {
	cursor, err := h.db.Collection("users").Find(
		context.Background(),
		bson.M{"deleted_at": nil},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "failed to get users"})
		return
	}
	defer cursor.Close(context.Background())

	var users []models.User
	if err := cursor.All(context.Background(), &users); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "failed to decode users"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"users": users,
		"count": len(users),
	})
}

// UpdateUser updates a user
func (h *Handler) UpdateUser(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid user ID"})
		return
	}

	var req models.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	// Build update document
	updateFields := bson.M{}
	if req.Email != nil {
		updateFields["email"] = *req.Email
	}
	if req.DefaultSalary != nil {
		updateFields["default_salary"] = *req.DefaultSalary
	}

	if len(updateFields) == 0 {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "no fields to update"})
		return
	}

	update := bson.M{"$set": updateFields}
	result, err := h.db.Collection("users").UpdateOne(
		context.Background(),
		bson.M{"_id": userID, "deleted_at": nil},
		update,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "failed to update user"})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "user not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "user updated successfully",
	})
}

// DeleteUser soft deletes a user
func (h *Handler) DeleteUser(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid user ID"})
		return
	}

	// Soft delete
	now := time.Now()
	update := bson.M{"$set": bson.M{"deleted_at": now}}
	result, err := h.db.Collection("users").UpdateOne(
		context.Background(),
		bson.M{"_id": userID, "deleted_at": nil},
		update,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "failed to delete user"})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "user not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "user deleted successfully",
	})
}
