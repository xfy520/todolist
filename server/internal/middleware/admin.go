package middleware

import (
	"net/http"

	"todo-server/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func AdminRequired(userRepo *repository.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
			c.Abort()
			return
		}

		user, err := userRepo.FindByID(userID.(uuid.UUID))
		if err != nil || !user.IsAdmin {
			c.JSON(http.StatusForbidden, gin.H{"error": "需要管理员权限"})
			c.Abort()
			return
		}

		c.Next()
	}
}
