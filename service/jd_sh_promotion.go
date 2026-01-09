package service

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
	"wfc_jd_report/core"

	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
)

// JdShPromotionResponse 京东推广链接响应结构体
type JdShPromotionResponse struct {
	JdUnionOpenShPromotionGetResponce struct {
		Code      string `json:"code"`
		GetResult string `json:"getResult"` // JSON字符串，需要二次解析
	} `json:"jd_union_open_sh_promotion_get_responce"`
}

// JdShPromotionResult 京东推广链接查询结果
type JdShPromotionResult struct {
	Code      int               `json:"code"`
	Message   string            `json:"message"`
	Data      JdShPromotionData `json:"data"`
	RequestId string            `json:"requestId"`
}

// JdShPromotionData 推广链接数据
type JdShPromotionData struct {
	AppUrl               string `json:"appUrl"`               // App深度链接
	ClickMonitorUrl      string `json:"clickMonitorUrl"`      // 点击监测链接
	ClickUrl             string `json:"clickUrl"`             // 点击链接
	ImpressionMonitorUrl string `json:"impressionMonitorUrl"` // 曝光监测链接
}

func GetJdShPromotion(ctx *gin.Context) {
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

	promotion, err := fetchJdShPromotion(ctx, jdAppkey, jdSecretkey)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to fetch promotion: %v", err.Error())})
		return
	}

	// 返回京东 API 原始格式的数据
	ctx.JSON(http.StatusOK, gin.H{
		"jd_union_open_sh_promotion_get_responce": gin.H{
			"code":      "0",
			"getResult": promotion,
		},
	})
}

func fetchJdShPromotion(ctx *gin.Context, jdAppkey, jdSecretkey string) (*JdShPromotionResult, error) {
	// 从 ctx 获取所有请求参数
	method := ctx.DefaultQuery("method", "jd.union.open.sh.promotion.get")
	format := ctx.DefaultQuery("format", "json")
	version := ctx.DefaultQuery("v", "1.0")
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
	values.Set("sign", sign)
	values.Set("360buy_param_json", paramJson)

	requestUrl := apiUrl + "?" + values.Encode()

	log.Println("requestUrl: ", requestUrl)

	// 发送HTTP请求
	var response JdShPromotionResponse
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
	if response.JdUnionOpenShPromotionGetResponce.Code != "0" {
		return nil, fmt.Errorf("JD API error: response_code=%s", response.JdUnionOpenShPromotionGetResponce.Code)
	}

	// 二次解析getResult字符串
	var promotionResult JdShPromotionResult
	err = jsoniter.Unmarshal([]byte(response.JdUnionOpenShPromotionGetResponce.GetResult), &promotionResult)
	if err != nil {
		return nil, fmt.Errorf("failed to parse getResult: %w", err)
	}

	if promotionResult.Code != 200 {
		return nil, fmt.Errorf("JD API query error: code=%d, message=%s",
			promotionResult.Code,
			promotionResult.Message)
	}

	return &promotionResult, nil
}
