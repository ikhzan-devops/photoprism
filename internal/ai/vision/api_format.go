package vision

type ApiFormat = string

const (
	ApiFormatUrl    ApiFormat = "url"
	ApiFormatImages ApiFormat = "images"
	ApiFormatVision ApiFormat = "vision"
	ApiFormatOllama ApiFormat = "ollama"
	ApiFormatOpenAI ApiFormat = "openai"
)
