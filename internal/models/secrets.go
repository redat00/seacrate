package models

import (
	"time"
)

type Secret struct {
	Key          string    `json:"key"`
	Folder       string    `json:"folder,omitempty"`
	ParentFolder string    `json:"parent_folder,omitempty"`
	Value        string    `json:"value"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
