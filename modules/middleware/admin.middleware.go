package middleware

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	boostrap "github.com/NVCLong/Alert-Server/bootstrap"
	"github.com/NVCLong/Alert-Server/common"
	redisService "github.com/NVCLong/Alert-Server/redis"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AbstractMiddleware interface {
	GetAdminHandlerFunc() gin.HandlerFunc
}
type AdminMiddleware struct {
	cacheService redisService.AbstractCacheService
	db           *gorm.DB
}

func NewAdminMiddleware(db *gorm.DB, cacheService redisService.AbstractCacheService) AbstractMiddleware {
	return &AdminMiddleware{
		cacheService: cacheService,
		db:           db,
	}
}
func (middleware *AdminMiddleware) GetAdminHandlerFunc() gin.HandlerFunc {
	return func(c *gin.Context) {

		if middleware.db == nil {
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

		if handleGetItemFromCahce(middleware.cacheService, strconv.FormatUint(uid, 10)) {
			c.Next()
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

		if err := middleware.db.Raw(query).Scan(&exists).Error; err != nil {
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

		handleSetItemToCache(middleware.cacheService, strconv.FormatUint(uid, 10), exists)
		log.Println("User is an admin")
		c.Next()
	}
}

func handleGetItemFromCahce(cacheService redisService.AbstractCacheService, key string) bool {
	fmt.Println("Start get from cache")
	item, error := cacheService.GetItem(key)
	if error != nil {
		return false
	}
	return item != ""
}

func handleSetItemToCache(cacheService redisService.AbstractCacheService, key string, item any) {
	ttlStr := boostrap.GetEnv(boostrap.EnvRedisTTL)
	var timeToLive int
	ttl, err := strconv.ParseInt(ttlStr, 10, 64)
	if err != nil {
		timeToLive = 0
	}
	timeToLive = int(ttl)
	cacheService.SetItem(key, item, time.Duration(timeToLive)*time.Second)
	fmt.Println("Set Item success")
}
