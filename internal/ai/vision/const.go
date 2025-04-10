package vision

type ModelType = string

const (
	ModelTypeLabels         ModelType = "labels"
	ModelTypeNsfw           ModelType = "nsfw"
	ModelTypeFaceEmbeddings ModelType = "face/embeddings"
	ModelTypeCaption        ModelType = "caption"
)
