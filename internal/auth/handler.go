package auth

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nkamil/biller-app/internal/models"
	"github.com/nkamil/biller-app/pkg/jwt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// Handler handles authentication requests
type Handler struct {
	db         *mongo.Database
	jwtManager *jwt.JWTManager
}

// NewHandler creates a new auth handler
func NewHandler(db *mongo.Database, jwtManager *jwt.JWTManager) *Handler {
	return &Handler{
		db:         db,
		jwtManager: jwtManager,
	}
}

// Register handles user registration
func (h *Handler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	// Check if user already exists
	var existingUser models.User
	err := h.db.Collection("users").FindOne(
		context.Background(),
		bson.M{"$or": []bson.M{
			{"username": req.Username},
			{"email": req.Email},
		}},
	).Decode(&existingUser)

	if err == nil {
		c.JSON(http.StatusConflict, models.ErrorResponse{Error: "username or email already exists"})
		return
	} else if err != mongo.ErrNoDocuments {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "database error"})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "failed to hash password"})
		return
	}

	// Create user
	user := models.User{
		ID:            primitive.NewObjectID(),
		Email:         req.Email,
		Username:      req.Username,
		PasswordHash:  string(hashedPassword),
		Role:          models.RoleUser,
		DefaultSalary: 0,
		CreatedAt:     time.Now(),
	}

	_, err = h.db.Collection("users").InsertOne(context.Background(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "user created successfully",
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
	})
}

// Login handles user login
func (h *Handler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	// Find user
	var user models.User
	err := h.db.Collection("users").FindOne(
		context.Background(),
		bson.M{"username": req.Username, "deleted_at": nil},
	).Decode(&user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "invalid credentials"})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "database error"})
		return
	}

	// Compare password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "invalid credentials"})
		return
	}

	// Generate token
	token, err := h.jwtManager.GenerateToken(user.ID, user.Username, string(user.Role))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, models.LoginResponse{
		AccessToken: token,
	})
}
