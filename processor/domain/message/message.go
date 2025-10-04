package message

type Event struct {
	ID       string `json:"id"`
	Event    string `json:"event"`
	FileName string `json:"fileName"`
	Size     int64  `json:"size"`
	Path     string `json:"path"`
}

type MovieInfoEvent struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
}
