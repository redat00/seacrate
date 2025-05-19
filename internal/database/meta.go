package database

import (
	"context"

	"github.com/redat00/seacrate/internal/models"
)

func (d *DatabaseEngine) CreateMeta(metaKey string, metaValue string) error {
	_, err := d.c.Exec(context.Background(), "INSERT INTO meta VALUES ($1, $2)", metaKey, metaValue)
	if err != nil {
		return err
	}
	return nil
}

func (d *DatabaseEngine) GetMeta(metaKey string) (*models.Meta, error) {
	var meta models.Meta
	err := d.c.QueryRow(context.Background(), "SELECT key, value FROM meta WHERE key=$1", metaKey).Scan(&meta.Key, &meta.Value)
	if err != nil {
		return &meta, err
	}
	return &meta, nil
}

// TODO: UPDATE

func (d *DatabaseEngine) DeleteMeta(metaKey string) error {
	_, err := d.c.Exec(context.Background(), "DELETE FROM meta WHERE key=?", metaKey)
	if err != nil {
		return err
	}
	return nil
}
