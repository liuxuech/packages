package sms

type Sms interface {
	Send(opts *MessageOption) error // 单条短信发送
	SendBatch()                     // 批量发送
}

type MessageOption struct {
	Phones     string `validate:"required=true"` // 被推送短信的手机号，单个推送和批量推送都可以
	TemplateID string `validate:"required=true"` // 短信模板码

	Sign          string // 签名
	TemplateParam string // 模板参数
}
