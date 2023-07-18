package api

import (
	"context"
	"net/http"

	"github.com/fernandoglatz/home-management/models"
	"github.com/fernandoglatz/home-management/repositories"
	"github.com/fernandoglatz/home-management/services"
	"github.com/gin-gonic/gin"
)

// CreateUser handles the creation of a new user
func CreateUser(c *gin.Context) {
	// Bind request body to User struct
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// Insert user into MongoDB
	service := services.NewService[*models.User]()
	err := service.Save(context.Background(), &user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to create user: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// GetUser handles the retrieval of a user by ID
func GetUser(c *gin.Context) {
	// Get user ID from request parameters
	userID := c.Param("id")

	// Find user in MongoDB
	repository := repositories.NewRepository[*models.User]()
	user, err := repository.FindByID(context.Background(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// GetUser handles the retrieval of a user by ID
func GetAllUsers(c *gin.Context) {

	// Find user in MongoDB
	repository := repositories.NewRepository[*models.User]()
	users, err := repository.FindAll(context.Background())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}

// UpdateUser handles the update of a user
func UpdateUser(c *gin.Context) {
	// Get user ID from request parameters
	userID := c.Param("id")

	// Bind request body to User struct
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// Update user in MongoDB
	user.ID = userID // Set the user ID in the struct
	service := services.NewService[*models.User]()
	err := service.Save(context.Background(), &user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update user: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// DeleteUser handles the deletion of a user
func DeleteUser(c *gin.Context) {
	// Get user ID from request parameters
	userID := c.Param("id")

	repository := repositories.NewRepository[*models.User]()

	user, err := repository.FindByID(context.Background(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}

	err = repository.Delete(context.Background(), user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
}
