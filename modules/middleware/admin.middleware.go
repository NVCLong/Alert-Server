package middleware

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/NVCLong/Alert-Server/common"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AdminMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		if db == nil {
			log.Println("Database connection is nil")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			c.Abort()
			return
		}

		adminHeader := c.Request.Header.Get("UserId")
		uidString := strings.Split(adminHeader, " ")[0]
		if uidString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Do not have user id"})
			c.Abort()
			return
		}
		uid, err := strconv.ParseUint(uidString, 10, 64) // Change to ParseInt if you need a signed integer
		if err != nil {
			log.Println("Error parsing UID:", err)
			c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
			c.Abort()
			return
		}
		var exists bool
		query := fmt.Sprintf(`
		    SELECT EXISTS (
			SELECT 1 
			FROM users 
			JOIN user_setting_info usi on usi.user_id = users.id
			WHERE users.id = %d 
			AND usi.role = '%s'
		);`, uid, common.ADMIN_ROLE)

		if err := db.Raw(query).Scan(&exists).Error; err != nil {
			log.Println("Error in query user with id ", uid)
			c.JSON(http.StatusUnauthorized, gin.H{"message": "User is not admin"})
			c.Abort()
			return
		}

		if !exists {
			log.Println("Do not found user in admin role")
			c.JSON(http.StatusUnauthorized, gin.H{"message": "User is not admin"})
			c.Abort()
			return
		}

		log.Println("User is an admin")
		c.Next()
	}
}
