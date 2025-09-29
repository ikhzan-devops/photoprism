package vision

// Label defaults for remote classification providers.
// These defaults were validated with Ollama/OpenAI adapters and should align with entity source constants (SrcOllama, SrcOpenAI).
var (
	LabelSystemDefault = "You are a PhotoPrism vision model. Output concise JSON that matches the schema."
	LabelPromptDefault = "Analyze the image and return label objects with name, confidence (0-1), and topicality (0-1)."
	LabelSchemaDefault = "{\n  \"labels\": [{\n    \"name\": \"\",\n    \"confidence\": 0,\n    \"topicality\": 0\n  }]\n}"
)
