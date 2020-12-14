package voice

type Voice interface {
	Call(*Options) error
}

type Options struct {
	TtsCode          string `validate:"required=true"`
	TtsParam         string
	CalledNumber     string `validate:"required=true"`
	CalledShowNumber string
}
