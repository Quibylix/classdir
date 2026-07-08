package hub

const (
	OpStroke = "stroke"
	OpClear  = "clear"
)

type AnnotationPoint struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type AnnotationPayload struct {
	Points    []AnnotationPoint `json:"points"`
	Color     string            `json:"color"`
	Thickness float64           `json:"thickness"`
}

type AnnotationOperation struct {
	Type    string             `json:"type"`
	ID      string             `json:"id"`
	Payload *AnnotationPayload `json:"payload,omitempty"`
}

type AnnotationParams struct {
	Type    string             `json:"type"`
	ID      string             `json:"id"`
	Payload *AnnotationPayload `json:"payload,omitempty"`
}

type annotationAddedData struct {
	Type    string             `json:"type"`
	ID      string             `json:"id"`
	Payload *AnnotationPayload `json:"payload,omitempty"`
}

type annotationAddedEvent struct {
	Event string              `json:"event"`
	Data  annotationAddedData `json:"data"`
}

type annotationsBatchData struct {
	OperationsBySlide map[int][]AnnotationOperation `json:"operations_by_slide"`
}

type annotationsBatchEvent struct {
	Event string               `json:"event"`
	Data  annotationsBatchData `json:"data"`
}
