package models

import "time"

type Folder struct {
	Id           int       `json:"id"`
	FullPath     string    `json:"full_path"`
	ParentFolder string    `json:"parent_folder"`
	CreatedAt    time.Time `json:"created_at"`
}

type FolderContent struct {
	Key       string    `json:"key"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"created_at"`
}
