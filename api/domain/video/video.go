package video

type Video struct {
	ID   string `json:"id" gorm:"column:id;primaryKey"`
	Name string `json:"name" gorm:"column:file_name"`
}

type VideoFilter struct {
	Name string `json:"name"`
}
