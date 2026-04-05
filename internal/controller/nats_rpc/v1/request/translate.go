package request

// Translate -.
type Translate struct {
	Source      string `json:"source"      validate:"required"`
	Destination string `json:"destination" validate:"required"`
	Original    string `json:"original"    validate:"required"`
}
