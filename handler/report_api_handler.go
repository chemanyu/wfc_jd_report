package handlers

import (
	"dmp_distribution/service"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ReportApiHandler struct {
	CommonHandler
}

var GetReportApiHandler = new(ReportApiHandler)

func init() {
	GetReportApiHandler.postMapping("order/query", getOrderQuery)
}

func getOrderQuery(ctx *gin.Context) {
	// 创建分发服务实例
	service.GetJdOrder(ctx)

	log.Printf("[Cron] Distribution service started successfully")

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "report_api retrieved successfully",
	})
}
