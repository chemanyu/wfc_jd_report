package dto

// JdOrderInfo 京东订单信息（参考 Java 版本）
type JdOrderInfo struct {
	// 附加字段
	ParentOrderId string `json:"parentOrderId"` // 父单的订单号
	ClickId       string `json:"clickId"`       // 点击ID
	OrderStatus   int    `json:"orderStatus"`   // 订单状态

	// 基础字段
	Id                  string        `json:"id"`                  // 标记唯一订单行
	OrderId             string        `json:"orderId"`             // 订单号
	ParentId            string        `json:"parentId"`            // 父单订单号
	OrderTime           string        `json:"orderTime"`           // 下单时间
	FinishTime          string        `json:"finishTime"`          // 完成时间
	ModifyTime          string        `json:"modifyTime"`          // 更新时间
	OrderEmt            int           `json:"orderEmt"`            // 下单设备 1.pc 2.无线
	Plus                int           `json:"plus"`                // 是否PLUS会员
	UnionId             string        `json:"unionId"`             // 推客ID
	SkuId               string        `json:"skuId"`               // 商品ID
	SkuName             string        `json:"skuName"`             // 商品名称
	SkuNum              int           `json:"skuNum"`              // 商品数量
	SkuReturnNum        int           `json:"skuReturnNum"`        // 已退货数量
	SkuFrozenNum        int           `json:"skuFrozenNum"`        // 售后中数量
	Price               string        `json:"price"`               // 商品单价
	CommissionRate      string        `json:"commissionRate"`      // 佣金比例
	SubSideRate         string        `json:"subSideRate"`         // 分成比例
	SubsidyRate         string        `json:"subsidyRate"`         // 补贴比例
	FinalRate           string        `json:"finalRate"`           // 最终分佣比例
	EstimateCosPrice    string        `json:"estimateCosPrice"`    // 预估计佣金额
	EstimateFee         string        `json:"estimateFee"`         // 推客预估佣金
	ActualCosPrice      string        `json:"actualCosPrice"`      // 实际计算佣金金额
	ActualFee           string        `json:"actualFee"`           // 推客实际佣金
	ValidCode           int           `json:"validCode"`           // sku维度有效码
	TraceType           int           `json:"traceType"`           // 同跨店
	PositionId          string        `json:"positionId"`          // 推广位ID
	SiteId              string        `json:"siteId"`              // 应用id
	UnionAlias          string        `json:"unionAlias"`          // PID所属母账号平台名称
	Pid                 string        `json:"pid"`                 // 格式:子推客ID_子站长应用ID_子推客推广位ID
	Cid1                int64         `json:"cid1"`                // 一级类目id
	Cid2                int64         `json:"cid2"`                // 二级类目id
	Cid3                int64         `json:"cid3"`                // 三级类目id
	SubUnionId          string        `json:"subUnionId"`          // 子渠道标识
	UnionTag            string        `json:"unionTag"`            // 联盟标签数据
	PopId               string        `json:"popId"`               // 商家ID
	Ext1                string        `json:"ext1"`                // 推广链接扩展字段
	PayMonth            string        `json:"payMonth"`            // 预估结算时间
	CpActId             int64         `json:"cpActId"`             // 招商团活动id
	UnionRole           string        `json:"unionRole"`           // 站长角色
	GiftCouponOcsAmount string        `json:"giftCouponOcsAmount"` // 礼金分摊金额
	GiftCouponKey       string        `json:"giftCouponKey"`       // 礼金批次ID
	BalanceExt          string        `json:"balanceExt"`          // 计佣扩展信息
	Sign                string        `json:"sign"`                // 数据签名
	ProPriceAmount      string        `json:"proPriceAmount"`      // 价保赔付金额
	Rid                 string        `json:"rid"`                 // 团长渠道ID
	ExpressStatus       int           `json:"expressStatus"`       // 发货状态
	ChannelId           int64         `json:"channelId"`           // 渠道关系ID
	SkuTag              string        `json:"skuTag"`              // 64位标签字段
	ItemId              string        `json:"itemId"`              // 联盟商品ID
	GoodsInfo           *GoodsInfo    `json:"goodsInfo"`           // 商品信息
	CategoryInfo        *CategoryInfo `json:"categoryInfo"`        // 类目信息
}

// GoodsInfo 商品信息
type GoodsInfo struct {
	ImageUrl  string `json:"imageUrl"`  // sku主图链接
	Owner     string `json:"owner"`     // g=自营，p=pop
	MainSkuId string `json:"mainSkuId"` // 自营商品主Id
	ProductId string `json:"productId"` // 非自营商品主Id
	ShopName  string `json:"shopName"`  // 店铺名称
	ShopId    string `json:"shopId"`    // 店铺Id
}

// CategoryInfo 类目信息
type CategoryInfo struct {
	Cid1     int64  `json:"cid1"`     // 一级类目id
	Cid2     int64  `json:"cid2"`     // 二级类目id
	Cid3     int64  `json:"cid3"`     // 三级类目id
	Cid1Name string `json:"cid1Name"` // 一级类目名称
	Cid2Name string `json:"cid2Name"` // 二级类目名称
	Cid3Name string `json:"cid3Name"` // 三级类目名称
}
