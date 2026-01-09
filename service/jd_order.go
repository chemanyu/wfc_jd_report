package service

import (
	"crypto/md5"
	"dmp_distribution/core"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
)

// JD API 配置
const (
	JD_APP_KEY      = "133a2a74cab37ba1a2ee7cdbeb0cc479"
	JD_APP_SECRET   = "c772a9a6f34a400097a21860ed0830f4" // 请替换为您的JD应用密钥
	JD_ACCESS_TOKEN = ""                                 // 请替换为您的JD访问令牌
)

// JD API 响应结构体
type JdOrderResponse struct {
	JdUnionOpenOrderRowQueryResponse struct {
		Code        string `json:"code"`
		QueryResult string `json:"queryResult"` // 注意：这是一个JSON字符串，需要二次解析
	} `json:"jd_union_open_order_row_query_responce"`
}

// JD API 查询结果结构体（用于解析queryResult字符串）
type JdQueryResult struct {
	Code      int    `json:"code"`
	Message   string `json:"message"`
	RequestId string `json:"requestId"`
	HasMore   bool   `json:"hasMore"`
	Data      []struct {
		OrderId          int64   `json:"orderId"`
		OrderTime        string  `json:"orderTime"`
		FinishTime       string  `json:"finishTime"`
		ModifyTime       string  `json:"modifyTime"`
		UnionId          int64   `json:"unionId"`
		SkuId            int64   `json:"skuId"`
		SkuName          string  `json:"skuName"`
		Price            float64 `json:"price"`
		FinalRate        float64 `json:"finalRate"`
		EstimateCosPrice float64 `json:"estimateCosPrice"`
		EstimateFee      float64 `json:"estimateFee"`
		ActualCosPrice   float64 `json:"actualCosPrice"`
		ActualFee        float64 `json:"actualFee"`
		ValidCode        int     `json:"validCode"`
		PositionId       int64   `json:"positionId"`
		Pid              string  `json:"pid"`
		Account          string  `json:"account"`
	} `json:"data"`
}

// JD订单数据结构用于Excel导出
type JdOrderData struct {
	OrderId          int64   `json:"orderId"`
	OrderTime        string  `json:"orderTime"`
	FinishTime       string  `json:"finishTime"`
	ModifyTime       string  `json:"modifyTime"`
	UnionId          int64   `json:"unionId"`
	SkuId            int64   `json:"skuId"`
	SkuName          string  `json:"skuName"`
	Price            float64 `json:"price"`
	FinalRate        float64 `json:"finalRate"`
	EstimateCosPrice float64 `json:"estimateCosPrice"`
	EstimateFee      float64 `json:"estimateFee"`
	ActualCosPrice   float64 `json:"actualCosPrice"`
	ActualFee        float64 `json:"actualFee"`
	ValidCode        int     `json:"validCode"`
	PositionId       int64   `json:"positionId"`
	Pid              string  `json:"pid"`
	Account          string  `json:"account"`
}

func GetJdOrder(ctx *gin.Context) {
	// 获取查询参数
	startTimeStr := ctx.Query("startTime")
	endTimeStr := ctx.Query("endTime")
	pageSize := ctx.DefaultQuery("pageSize", "200")
	orderType := ctx.Query("type") // 1-3

	if startTimeStr == "" || endTimeStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "startTime and endTime are required"})
		return
	}
	// startTimeStr, _ = url.QueryUnescape(startTimeStr)
	// endTimeStr, _ = url.QueryUnescape(endTimeStr)

	// 处理特殊的24:00:00时间格式，将其转换为23:59:59
	endTimeStr = strings.Replace(endTimeStr, "24:00:00", "23:59:59", -1)

	// 解析时间
	startTime, err := time.Parse("2006-01-02 15:04:05", startTimeStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Invalid startTime format: %s, expected: 2006-01-02 15:04:05 (e.g., 2025-09-23 14:00:00)", startTimeStr),
		})
		return
	}

	endTime, err := time.Parse("2006-01-02 15:04:05", endTimeStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Invalid endTime format: %s, expected: 2006-01-02 15:04:05 (e.g., 2025-09-23 23:59:59). Note: use 23:59:59 instead of 24:00:00", endTimeStr),
		})
		return
	}

	var allOrders []JdOrderData

	// 按小时循环时间段
	currentTime := startTime
	for currentTime.Before(endTime) {
		// 计算当前小时的结束时间
		hourEndTime := currentTime.Add(time.Hour)
		if hourEndTime.After(endTime) {
			hourEndTime = endTime
		}

		currentStartStr := currentTime.Format("2006-01-02 15:04:05")
		currentEndStr := hourEndTime.Format("2006-01-02 15:04:05")

		// 循环调用API，type分别为1、2、3
		//for orderType := 1; orderType <= 3; orderType++ {
		size, _ := strconv.Atoi(pageSize)
		typeInt, _ := strconv.Atoi(orderType)
		orders, err := fetchJdOrders(currentStartStr, currentEndStr, size, typeInt)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to fetch orders for type %d, time %s-%s: %v", 1, currentStartStr, currentEndStr, err)})
			return
		}
		allOrders = append(allOrders, orders...)

		// 在不同类型的API调用之间添加延迟
		// if orderType < 3 {
		// 	time.Sleep(1 * time.Second)
		// }
		//}

		// 移动到下一个小时
		currentTime = hourEndTime
	}

	if len(allOrders) == 0 {
		ctx.JSON(http.StatusOK, gin.H{"message": "No orders found", "count": 0})
		return
	}

}

// 获取JD订单数据
func fetchJdOrders(startTime, endTime string, pageSize, orderType int) ([]JdOrderData, error) {
	var allOrders []JdOrderData
	pageIndex := 1
	hasMore := true

	for hasMore {
		orders, hasMoreData, err := fetchJdOrdersPage(startTime, endTime, pageSize, pageIndex, orderType)
		if err != nil {
			return nil, err
		}

		allOrders = append(allOrders, orders...)
		hasMore = hasMoreData
		pageIndex++

		// 防止无限循环，最多获取100页
		if pageIndex > 100 {
			break
		}

		// 在分页请求之间添加短暂延迟
		if hasMore {
			time.Sleep(500 * time.Millisecond)
		}
	}

	return allOrders, nil
}

// 获取单页JD订单数据
func fetchJdOrdersPage(startTime, endTime string, pageSize, pageIndex, orderType int) ([]JdOrderData, bool, error) {
	// JD API 参数
	method := "jd.union.open.order.row.query"
	version := "1.0"
	timestamp := time.Now().Format("2006-01-02 15:04:05.000-0700")

	// 构建请求参数
	orderReq := map[string]interface{}{
		"orderReq": map[string]interface{}{
			"pageIndex": pageIndex,
			"pageSize":  pageSize,
			"startTime": startTime,
			"endTime":   endTime,
			"type":      orderType,
		},
	}

	paramJsonBytes, _ := jsoniter.Marshal(orderReq)
	paramJson := string(paramJsonBytes)

	// 构建签名参数
	params := map[string]string{
		"access_token":      JD_ACCESS_TOKEN,
		"app_key":           JD_APP_KEY,
		"method":            method,
		"v":                 version,
		"timestamp":         timestamp,
		"360buy_param_json": paramJson,
	}

	// 生成签名
	sign := generateJdSign(params, JD_APP_SECRET)

	// 构建请求URL
	apiUrl := "https://api.jd.com/routerjson"
	values := url.Values{}
	values.Set("access_token", JD_ACCESS_TOKEN)
	values.Set("app_key", JD_APP_KEY)
	values.Set("method", method)
	values.Set("v", version)
	values.Set("sign", sign)
	values.Set("360buy_param_json", paramJson)
	values.Set("timestamp", timestamp)

	requestUrl := apiUrl + "?" + values.Encode()

	fmt.Println("Request URL:", requestUrl)

	// 发送HTTP请求
	var response JdOrderResponse
	err := core.Fast_Http("", requestUrl, "", "GET", func(content []byte) error {
		resp := string(content)
		//log.Println("tzfb-impression : ", uri+callbackurl, resp)
		if strings.Contains(resp, `success":false`) {
			return errors.New(requestUrl + "\n" + resp)
		}
		return nil
	})

	if err != nil {
		return nil, false, err
	}

	// 检查响应状态
	if response.JdUnionOpenOrderRowQueryResponse.Code != "0" {
		return nil, false, fmt.Errorf("JD API error: code=%s", response.JdUnionOpenOrderRowQueryResponse.Code)
	}

	// 二次解析queryResult字符串
	var queryResult JdQueryResult
	err = jsoniter.Unmarshal([]byte(response.JdUnionOpenOrderRowQueryResponse.QueryResult), &queryResult)
	if err != nil {
		return nil, false, fmt.Errorf("failed to parse queryResult: %w", err)
	}

	if queryResult.Code != 200 {
		return nil, false, fmt.Errorf("JD API query error: code=%d, message=%s",
			queryResult.Code,
			queryResult.Message)
	}

	// 转换数据格式，只保留指定的PositionId
	var orders []JdOrderData
	for _, item := range queryResult.Data {
		order := JdOrderData{
			OrderId:          item.OrderId,
			OrderTime:        item.OrderTime,
			FinishTime:       item.FinishTime,
			ModifyTime:       item.ModifyTime,
			UnionId:          item.UnionId,
			SkuId:            item.SkuId,
			SkuName:          item.SkuName,
			Price:            item.Price,
			FinalRate:        item.FinalRate,
			EstimateCosPrice: item.EstimateCosPrice,
			EstimateFee:      item.EstimateFee,
			ActualCosPrice:   item.ActualCosPrice,
			ActualFee:        item.ActualFee,
			ValidCode:        item.ValidCode,
			PositionId:       item.PositionId,
			Pid:              item.Pid,
			Account:          item.Account,
		}
		orders = append(orders, order)
	}

	return orders, queryResult.HasMore, nil
}

// 生成JD API签名
func generateJdSign(params map[string]string, appSecret string) string {
	// 1. 将所有请求参数按照字母先后顺序排列
	var keys []string
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// 2. 把所有参数名和参数值进行拼接
	var signStr strings.Builder
	for _, k := range keys {
		signStr.WriteString(k)
		signStr.WriteString(params[k])
	}

	// 3. 把appSecret夹在字符串的两端
	finalStr := appSecret + signStr.String() + appSecret
	fmt.Println("String to Sign:", finalStr)

	// 4. 使用MD5进行加密，再转化成大写
	hash := md5.Sum([]byte(finalStr))
	return strings.ToUpper(fmt.Sprintf("%x", hash))
}
