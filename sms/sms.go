package sms

type Sms interface {
	Send(opts *MessageOption) error // 单条短信发送
}

type MessageOption struct {
	Phones     string `validate:"required=true"` // 被推送短信的手机号，格式：15612345678,1341234567
	TemplateID string `validate:"required=true"` // 短信模板ID

	Sign          string // 签名
	TemplateParam string // 模板参数
}
