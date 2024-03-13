package handler

type ToDoRequest struct {
	Kegiatan string `json:"kegiatan" form:"kegiatan"`
}
