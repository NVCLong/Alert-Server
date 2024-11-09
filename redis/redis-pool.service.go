package redis

import (
	"encoding/json"
	"fmt"
	"github.com/NVCLong/Alert-Server/common"
	"github.com/NVCLong/Alert-Server/dto"
	conditionbatch "github.com/NVCLong/Alert-Server/modules/condition-batch"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"sync"
)

var conditionBatchService conditionbatch.AbstractService

func StartWorkerPool(client *redis.Client, wg *sync.WaitGroup, ctx *gin.Context, numberWorkers int, db *gorm.DB) {
	conditionBatchService = conditionbatch.NewBatchService(db)
	for w := 1; w <= numberWorkers; w++ {
		wg.Add(1)
		go ProcessJob(w, client, wg, ctx, conditionBatchService)
	}
}
func PushJobToQueue(context *gin.Context, job string, logger common.AbstractLogger) {
	redisClient := NewRedisConnection()
	defer redisClient.Close()
	err := redisClient.RPush(context, "jobQueue", job).Err()
	if err != nil {
		logger.Debug("Fail to push job into queue")
		return
	}
}
func ProcessJob(id int, client *redis.Client, wg *sync.WaitGroup, ctx *gin.Context, batchService conditionbatch.AbstractService) {
	defer wg.Done()
	for {
		job, err := client.BLPop(ctx, 0, "jobQueue").Result()
		if err != nil {
			fmt.Printf("Worker %d error fetching job: %v\n", id, err)
			continue
		}
		fmt.Printf("Worker %d processing job: %s\n", id, job[1])

		var jobCreation dto.JobCreation
		err = json.Unmarshal([]byte(job[1]), &jobCreation)
		if err != nil {
			fmt.Printf("Worker %d failed to parse job details: %v\n", id, err)
			return
		}
		batchService.TriggerCondition(jobCreation)
		return
	}
}
