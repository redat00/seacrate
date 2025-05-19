package database

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	seacrateErrors "github.com/redat00/seacrate/internal/errors"
	"github.com/redat00/seacrate/internal/models"
)

func splitSecretKey(key string) ([]string, string, string) {
	var folders []string
	var parentFolder string
	var secretKey string

	// Get the number of / in the string, and if one simply return the data
	if slashesCount := strings.Count(key, "/"); slashesCount == 1 {
		return folders, "/", strings.ReplaceAll(key, "/", "")
	}

	// Extract secret key, and get rid of it to only keep folders
	getLastSlashIndex := strings.LastIndex(key, "/")
	secretKey = key[getLastSlashIndex+1:]
	key = key[:getLastSlashIndex]
	parentFolder = key

	// Determine each folders that would need to be created
	currentPath := ""
	for _, folder := range strings.Split(key, "/") {
		currentPath = filepath.Join(currentPath, folder)
		if currentPath != "" {
			folders = append(folders, fmt.Sprintf("/%s", currentPath))
		}
	}

	return folders, parentFolder, secretKey
}

func (d *DatabaseEngine) checkIfSecretExist(key string, parentFolder string) bool {
	var count int64
	d.c.QueryRow(context.Background(), "SELECT COUNT(key) FROM secrets WHERE key = $1 AND folder = $2", key, parentFolder).Scan(&count)
	if count > 0 {
		return true
	}
	return false
}

func (d *DatabaseEngine) CreateSecret(key string, value string) error {
	folders, parentFolder, secretKey := splitSecretKey(key)

	// Make sure that no secret with similar key already exist
	if d.checkIfSecretExist(secretKey, parentFolder) {
		return seacrateErrors.ErrSecretDuplicateKey{Key: key}
	}

	// Make sure that we are not overriding any folder
	if d.checkIfFolderExist(key) {
		return seacrateErrors.ErrOverridingFolder{Key: key}
	}

	// We need to create each folder
	err := d.createMultipleFolders(folders)
	if err != nil {
		return err
	}

	// We can then create the secret
	_, err = d.c.Exec(context.Background(), "INSERT INTO secrets (key, value, folder) VALUES ($1, $2, $3)", secretKey, value, parentFolder)
	if err != nil {
		return err
	}

	return nil
}

// Get all secrets within a folder
// Beware: This is used to list only, and will not return the actual secrets values
func (d *DatabaseEngine) getSecretsInFolder(path string) ([]models.FolderContent, error) {
	var secrets []models.FolderContent
	secretRows, err := d.c.Query(context.Background(), "SELECT key, created_at FROM secrets WHERE folder = $1", path)
	if err != nil {
		return secrets, err
	}
	defer secretRows.Close()
	for secretRows.Next() {
		var sR models.FolderContent
		sR.Type = "secret"
		err := secretRows.Scan(&sR.Key, &sR.CreatedAt)
		if err != nil {
			return secrets, err
		}
		secrets = append(secrets, sR)
	}
	return secrets, nil
}

func (d *DatabaseEngine) GetSecret(key string) (bool, []models.FolderContent, *models.Secret, error) {
	var secret models.Secret
	var folderContent []models.FolderContent

	_, parentFolder, secretKey := splitSecretKey(key)

	if isSecret := d.checkIfSecretExist(secretKey, parentFolder); isSecret == true {
		var secret models.Secret
		err := d.c.QueryRow(context.Background(), "SELECT key, value, created_at, updated_at FROM secrets WHERE key = $1 AND folder = $2", secretKey, parentFolder).Scan(&secret.Key, &secret.Value, &secret.CreatedAt, &secret.UpdatedAt)
		return false, folderContent, &secret, err
	} else if isFolder := d.checkIfFolderExist(key); isFolder == true {
		// Get all secrets for given folder
		secrets, err := d.getSecretsInFolder(key)
		if err != nil {
			return true, folderContent, &secret, err
		}
		folderContent = append(folderContent, secrets...)

		// Get all folders for given folder
		folders, err := d.getFoldersInFolder(key)
		if err != nil {
			return true, folderContent, &secret, err
		}
		folderContent = append(folderContent, folders...)

		return true, folderContent, &secret, nil
	}

	return false, folderContent, &secret, seacrateErrors.ErrSecretNotFound{Key: key}
}

// TODO: WORK ON SECRET UPDATE

// Delete secret from database
//
// This will also handle the deletion of the parent folder, if it's empty
func (d *DatabaseEngine) DeleteSecret(key string) error {
	_, parentFolder, secretKey := splitSecretKey(key)
	if exist := d.checkIfSecretExist(secretKey, parentFolder); exist {
		_, err := d.c.Exec(context.Background(), "DELETE FROM secrets WHERE key = $1 AND folder = $2", secretKey, parentFolder)
		if err != nil {
			return err
		}
		err = d.deleteFolderIfEmptyRecursively(parentFolder)
		if err != nil {
			return err
		}
		return nil
	}
	return seacrateErrors.ErrSecretNotFound{Key: key}
}
