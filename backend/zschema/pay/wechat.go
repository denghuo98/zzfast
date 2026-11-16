package pay

import "time"

// WechatPrepayInput 微信支付预下单入参（金额单位：分）。
type WechatPrepayInput struct {
	OutTradeNo  string     `json:"outTradeNo"  dc:"商户订单号"`
	Description string     `json:"description" dc:"商品描述"`
	Total       int64      `json:"total"       dc:"订单总金额（分）"`
	Attach      string     `json:"attach"      dc:"附加数据"`
	ExpireAt    *time.Time `json:"expireAt"    dc:"订单失效时间"`
	OpenId      string     `json:"openId"      dc:"JSAPI 支付用户 OpenId"`
}

// WechatNativeOutput Native 支付预下单输出。
type WechatNativeOutput struct {
	CodeURL string `json:"codeUrl" dc:"二维码链接"`
}

// WechatJSAPIOutput JSAPI 支付预下单输出。
type WechatJSAPIOutput struct {
	PrepayId string            `json:"prepayId" dc:"预支付交易会话标识"`
	Params   map[string]string `json:"params"   dc:"前端调起支付所需参数"`
}

// WechatQueryOutput 支付订单查询输出。
type WechatQueryOutput struct {
	OutTradeNo     string `json:"outTradeNo"     dc:"商户订单号"`
	TransactionId  string `json:"transactionId"  dc:"微信支付订单号"`
	TradeState     string `json:"tradeState"     dc:"交易状态"`
	TradeStateDesc string `json:"tradeStateDesc" dc:"交易状态描述"`
	TradeType      string `json:"tradeType"      dc:"交易类型"`
	SuccessTime    string `json:"successTime"    dc:"支付完成时间"`
	Total          int64  `json:"total"          dc:"订单总金额（分）"`
	PayerTotal     int64  `json:"payerTotal"     dc:"用户实付金额（分）"`
	OpenId         string `json:"openId"         dc:"支付用户 OpenId"`
}

// WechatRefundInput 微信退款入参（金额单位：分）。
type WechatRefundInput struct {
	OutTradeNo  string `json:"outTradeNo"  dc:"商户订单号"`
	OutRefundNo string `json:"outRefundNo" dc:"商户退款单号"`
	Reason      string `json:"reason"      dc:"退款原因"`
	Refund      int64  `json:"refund"      dc:"退款金额（分）"`
	Total       int64  `json:"total"       dc:"原订单总金额（分）"`
}

// WechatRefundOutput 微信退款输出。
type WechatRefundOutput struct {
	RefundId    string `json:"refundId"    dc:"微信支付退款号"`
	OutRefundNo string `json:"outRefundNo" dc:"商户退款单号"`
	Status      string `json:"status"      dc:"退款状态"`
	Total       int64  `json:"total"       dc:"订单总金额（分）"`
	Refund      int64  `json:"refund"      dc:"退款金额（分）"`
}

// WechatNotifyInput 微信支付/退款回调入参。
type WechatNotifyInput struct {
	Id           string               `json:"id"           dc:"通知 ID"`
	EventType    string               `json:"eventType"    dc:"事件类型"`
	ResourceType string               `json:"resourceType" dc:"资源类型"`
	Summary      string               `json:"summary"      dc:"通知摘要"`
	Payment      *WechatNotifyPayment `json:"payment"      dc:"支付通知内容"`
	Refund       *WechatNotifyRefund  `json:"refund"       dc:"退款通知内容"`
}

// WechatNotifyPayment 支付成功通知。
type WechatNotifyPayment struct {
	OutTradeNo     string `json:"outTradeNo"`
	TransactionId  string `json:"transactionId"`
	TradeState     string `json:"tradeState"`
	TradeStateDesc string `json:"tradeStateDesc"`
	TradeType      string `json:"tradeType"`
	SuccessTime    string `json:"successTime"`
	Total          int64  `json:"total"`
	PayerTotal     int64  `json:"payerTotal"`
	OpenId         string `json:"openId"`
	Attach         string `json:"attach"`
}

// WechatNotifyRefund 退款结果通知。
type WechatNotifyRefund struct {
	OutRefundNo   string `json:"outRefundNo"`
	RefundId      string `json:"refundId"`
	OutTradeNo    string `json:"outTradeNo"`
	TransactionId string `json:"transactionId"`
	RefundStatus  string `json:"refundStatus"`
	SuccessTime   string `json:"successTime"`
	Total         int64  `json:"total"`
	Refund        int64  `json:"refund"`
	PayerTotal    int64  `json:"payerTotal"`
	PayerRefund   int64  `json:"payerRefund"`
}
