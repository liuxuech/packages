package sms

type Sms interface {
	Send(opts *SingleMsg) error // 单条短信发送
	SendBatch()                 // 批量发送
}

type SingleMsg struct {
	RegionId string
	// 签名，ep：普诺思博
	Sign string
	// 被推送短信的手机号
	TargetPhone string
	// 短信模板码
	TemplateCode string
	// 模板参数
	TemplateParam string
}

type BatchMsg struct {
	RegionId       string
	Sign           string   // 签名
	TargetPhones   []string // 目标手机号切片
	TemplateCode   string   // 短信模板码
	TemplateParams []string // 模板参数切片
}
