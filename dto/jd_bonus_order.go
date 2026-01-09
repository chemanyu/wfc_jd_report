package dto

// OrderBonusRowResp 京东奖励订单信息
type OrderBonusRowResp struct {
	Account          string  `json:"account"`          // 账户ID
	ActivityId       int64   `json:"activityId"`       // 活动ID
	ActivityName     string  `json:"activityName"`     // 活动名称
	ActualBonusFee   float64 `json:"actualBonusFee"`   // 实际奖励金额
	ActualCosPrice   float64 `json:"actualCosPrice"`   // 实际计算佣金金额
	ActualFee        float64 `json:"actualFee"`        // 推客实际佣金
	BonusInvalidCode string  `json:"bonusInvalidCode"` // 奖励无效代码
	BonusInvalidText string  `json:"bonusInvalidText"` // 奖励无效文本
	BonusState       int     `json:"bonusState"`       // 奖励状态
	BonusText        string  `json:"bonusText"`        // 奖励状态文本
	ChannelId        int64   `json:"channelId"`        // 渠道关系ID
	CommissionRate   float64 `json:"commissionRate"`   // 佣金比例
	EstimateBonusFee float64 `json:"estimateBonusFee"` // 预估奖励金额
	EstimateCosPrice float64 `json:"estimateCosPrice"` // 预估计佣金额
	EstimateFee      float64 `json:"estimateFee"`      // 推客预估佣金
	Ext1             string  `json:"ext1"`             // 推广链接扩展字段
	FinalRate        float64 `json:"finalRate"`        // 最终分佣比例
	FinishTime       int64   `json:"finishTime"`       // 完成时间（时间戳）
	Id               int64   `json:"id"`               // 标记唯一订单行
	ItemId           string  `json:"itemId"`           // 联盟商品ID
	OrderId          int64   `json:"orderId"`          // 订单号
	OrderState       int     `json:"orderState"`       // 订单状态
	OrderText        string  `json:"orderText"`        // 订单状态文本
	OrderTime        int64   `json:"orderTime"`        // 下单时间（时间戳，毫秒）
	ParentId         int64   `json:"parentId"`         // 父单订单号
	PayPrice         float64 `json:"payPrice"`         // 支付价格
	Pid              string  `json:"pid"`              // 格式:子推客ID_子站长应用ID_子推客推广位ID
	PositionId       int64   `json:"positionId"`       // 推广位ID
	SkuId            int64   `json:"skuId"`            // 商品ID
	SkuName          string  `json:"skuName"`          // 商品名称
	SortValue        string  `json:"sortValue"`        // 排序值
	SubSideRate      float64 `json:"subSideRate"`      // 分成比例
	SubUnionId       string  `json:"subUnionId"`       // 子渠道标识
	SubsidyRate      float64 `json:"subsidyRate"`      // 补贴比例
	UnionAlias       string  `json:"unionAlias"`       // PID所属母账号平台名称
	UnionId          int64   `json:"unionId"`          // 推客ID
}
