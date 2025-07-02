package video

type Video struct {
	ID       string `json:"id" gorm:"column:id;primaryKey"`
	Event    string `json:"event" gorm:"column:event"`
	FileName string `json:"fileName" gorm:"column:file_name"`
	Size     int64  `json:"size" gorm:"column:size"`
	Path     string `json:"path" gorm:"column:path"`
}
