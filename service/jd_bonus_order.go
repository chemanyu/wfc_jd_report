package service

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
	"wfc_jd_report/core"
	"wfc_jd_report/dto"

	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
)

// JdOrderResponse JD API 响应结构体
type JdOrderBonusResponse struct {
	JdUnionOpenOrderBonusQueryResponce struct {
		Code        string `json:"code"`
		QueryResult string `json:"queryResult"` // 注意：这是一个JSON字符串，需要二次解析
	} `json:"jd_union_open_order_bonus_query_responce"`
}

// JdQueryResult JD API 查询结果
type JdQueryBonusResult struct {
	Code        int                     `json:"code"`
	Message     string                  `json:"message"`
	HasMore     bool                    `json:"hasMore"`
	Data        []dto.OrderBonusRowResp `json:"data"`
	RequestId   string                  `json:requestId`
	ApiCodeEnum string                  `json:apiCodeEnum`
}

func GetJdBonusOrder(ctx *gin.Context) {
	// 获取请求参数中的 app_key
	appKey := ctx.Query("app_key")
	if appKey == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "app_key参数不能为空"})
		return
	}

	// 通过美数appkey获取京东凭证
	jdAppkey, jdSecretkey, err := AdToolWfcTokenMapper.GetJdCredentialsByMsAppkey(appKey)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("获取京东凭证失败: %v", err)})
		return
	}

	orders, err := fetchJdBonusOrdersPage(ctx, jdAppkey, jdSecretkey)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to fetch orders: %v", err.Error())})
		return
	}

	// 返回京东 API 原始格式的数据
	ctx.JSON(http.StatusOK, gin.H{
		"jd_union_open_order_row_query_responce": gin.H{
			"code":        "0",
			"queryResult": orders,
		},
	})
}

// fetchJdOrdersPage 获取单页JD订单数据
func fetchJdBonusOrdersPage(ctx *gin.Context, jdAppkey, jdSecretkey string) (*JdQueryBonusResult, error) {
	// 从 ctx 获取所有请求参数
	method := ctx.DefaultQuery("method", "jd.union.open.order.bonus.query")
	format := ctx.DefaultQuery("format", "json")
	version := ctx.DefaultQuery("v", "1.0")
	signMethod := ctx.DefaultQuery("sign_method", "md5")
	timestamp := ctx.Query("timestamp")
	if timestamp == "" {
		// 如果没有传 timestamp，使用当前时间减3分钟
		timestamp = time.Now().Add(-3 * time.Minute).Format("2006-01-02 15:04:05")
	}

	// 获取业务参数 360buy_param_json
	paramJson := ctx.Query("360buy_param_json")
	if paramJson == "" {
		return nil, fmt.Errorf("360buy_param_json 参数不能为空")
	}

	// 构建签名参数（按字母排序）
	params := map[string]string{
		"method":            method,
		"app_key":           jdAppkey,
		"timestamp":         timestamp,
		"format":            format,
		"v":                 version,
		"sign_method":       signMethod,
		"360buy_param_json": paramJson,
	}

	// 生成签名
	sign := generateJdSign(params, jdSecretkey)

	// 构建请求URL，所有参数拼接在URL上
	apiUrl := "https://api.jd.com/routerjson"
	values := url.Values{}
	values.Set("method", method)
	values.Set("app_key", jdAppkey)
	values.Set("timestamp", timestamp)
	values.Set("format", format)
	values.Set("v", version)
	values.Set("sign_method", signMethod)
	values.Set("sign", sign)
	values.Set("360buy_param_json", paramJson)

	requestUrl := apiUrl + "?" + values.Encode()

	log.Println("requestUrl: ", requestUrl)

	// 发送HTTP请求
	var response JdOrderBonusResponse
	err := core.Fast_Http("", requestUrl, "", "GET", func(content []byte) error {
		resp := string(content)
		log.Println("resp: ", resp)
		if strings.Contains(resp, "error_response") {
			return fmt.Errorf("errInfo: %s", resp)
		}
		// 解析响应
		if err := jsoniter.Unmarshal(content, &response); err != nil {
			return fmt.Errorf("failed to unmarshal response: %w", err)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	// 检查响应状态
	if response.JdUnionOpenOrderBonusQueryResponce.Code != "0" {
		return nil, fmt.Errorf("JD API error: response_code=%s", response)
	}

	// 二次解析queryResult字符串
	var queryResult JdQueryBonusResult
	err = jsoniter.Unmarshal([]byte(response.JdUnionOpenOrderBonusQueryResponce.QueryResult), &queryResult)
	if err != nil {
		return nil, fmt.Errorf("failed to parse queryResult: %w", err)
	}

	if queryResult.Code != 200 {
		return nil, fmt.Errorf("JD API query error: code=%d, message=%s",
			queryResult.Code,
			queryResult.Message)
	}
	if len(queryResult.Data) != 0 {
		queryResult.HasMore = true
	}

	// 处理订单数据，设置附加字段并过滤账户
	var filterResult []dto.OrderBonusRowResp
	for i := range queryResult.Data {
		_, exists := WfcAccountMapping[queryResult.Data[i].Account]
		log.Println("account_exist: ", exists, queryResult.Data[i].Account)
		if !exists { // 账户不存在
			continue
		}
		filterResult = append(filterResult, queryResult.Data[i])
	}
	queryResult.Data = filterResult

	// 将解析后的 queryResult 重新序列化为 JSON 对象（而不是字符串）
	// 但这里我们直接返回 queryResult，让上层处理
	return &queryResult, nil
}
