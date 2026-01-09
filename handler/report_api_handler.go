package handlers

import (
	"wfc_jd_report/service"

	"github.com/gin-gonic/gin"
)

type ReportApiHandler struct {
	CommonHandler
}

var GetReportApiHandler = new(ReportApiHandler)

func init() {
	GetReportApiHandler.getMapping("order/query", getOrderQuery)
	GetReportApiHandler.getMapping("order/bonus/query", getOrderBonusQuery)
}

func getOrderQuery(ctx *gin.Context) {
	// 调用服务层，服务层会直接写入响应
	service.GetJdOrder(ctx)
}

func getOrderBonusQuery(ctx *gin.Context) {
	// 调用服务层，服务层会直接写入响应
	service.GetJdBonusOrder(ctx)
}
