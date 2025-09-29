package vision

// CaptionPromptDefault is the default prompt used to generate captions.
var CaptionPromptDefault = "Create a caption with exactly one sentence in the active voice that describes the main visual content." +
	" Begin with the main subject and clear action. Avoid text formatting, meta-language, and filler words."

// CaptionModelDefault specifies the default model used to generate captions,
// see https://ollama.com/search?c=vision for a list of available models.
var CaptionModelDefault = "gemma3"
