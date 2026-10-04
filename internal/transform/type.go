package transform

type ResizeTransform struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

type CropTransform struct {
	Width  int `json:"width"`
	Height int `json:"height"`
	X      int `json:"x"`
	Y      int `json:"y"`
}

type FiltersTransform struct {
	Grayscale *bool `json:"grayscale"`
	Sepia     *bool `json:"sepia"`
}

type Transforms struct {
	Resize  *ResizeTransform  `json:"resize"`
	Crop    *CropTransform    `json:"crop"`
	Filters *FiltersTransform `json:"filters"`
	Format  *string           `json:"format"`
	Rotate  *int              `json:"rotate"`
}

type TransformRequest struct {
	Transformations Transforms `json:"transformations"`
}
