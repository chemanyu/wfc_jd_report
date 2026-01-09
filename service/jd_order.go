package service

import (
	"crypto/md5"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
	"wfc_jd_report/core"
	"wfc_jd_report/dto"
	"wfc_jd_report/module"

	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
)

var (
	AdToolWfcTokenMapper = module.AdToolWfcToken{}
)

// JdOrderResponse JD API 响应结构体
type JdOrderResponse struct {
	JdUnionOpenOrderRowQueryResponse struct {
		Code        string `json:"code"`
		QueryResult string `json:"queryResult"` // 注意：这是一个JSON字符串，需要二次解析
	} `json:"jd_union_open_order_row_query_responce"`
}

// JdQueryResult JD API 查询结果
type JdQueryResult struct {
	Code      int               `json:"code"`
	Message   string            `json:"message"`
	RequestId string            `json:"requestId"`
	HasMore   bool              `json:"hasMore"`
	Data      []dto.JdOrderInfo `json:"data"`
}

func GetJdOrder(ctx *gin.Context) {
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

	orders, err := fetchJdOrdersPage(ctx, jdAppkey, jdSecretkey)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to fetch orders: %v", err)})
		return
	}

	// 返回订单数据
	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": orders,
		"msg":  "success",
	})
}

// fetchJdOrdersPage 获取单页JD订单数据
func fetchJdOrdersPage(ctx *gin.Context, jdAppkey, jdSecretkey string) (*JdOrderResponse, error) {
	// 从 ctx 获取所有请求参数
	method := ctx.DefaultQuery("method", "jd.union.open.order.row.query")
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
	apiUrl := "https://美数域名/wfc/jd/order/query"
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

	fmt.Println("Request URL:", requestUrl)

	// 发送HTTP请求
	var response JdOrderResponse
	err := core.Fast_Http("", requestUrl, "", "GET", func(content []byte) error {
		resp := string(content)
		// 解析响应
		if err := jsoniter.Unmarshal(content, &response); err != nil {
			return fmt.Errorf("failed to unmarshal response: %w", err)
		}
		if strings.Contains(resp, `success":false`) {
			return errors.New(requestUrl + "\n" + resp)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	// 检查响应状态
	if response.JdUnionOpenOrderRowQueryResponse.Code != "0" {
		return nil, fmt.Errorf("JD API error: code=%s", response.JdUnionOpenOrderRowQueryResponse.Code)
	}

	// 二次解析queryResult字符串
	var queryResult JdQueryResult
	err = jsoniter.Unmarshal([]byte(response.JdUnionOpenOrderRowQueryResponse.QueryResult), &queryResult)
	if err != nil {
		return nil, fmt.Errorf("failed to parse queryResult: %w", err)
	}

	if queryResult.Code != 200 {
		return nil, fmt.Errorf("JD API query error: code=%d, message=%s",
			queryResult.Code,
			queryResult.Message)
	}

	// 处理订单数据，设置附加字段
	for i := range queryResult.Data {
		queryResult.Data[i].ParentOrderId = queryResult.Data[i].ParentId
		queryResult.Data[i].ClickId = queryResult.Data[i].Ext1
		queryResult.Data[i].OrderStatus = queryResult.Data[i].ValidCode
	}

	return &response, nil
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
