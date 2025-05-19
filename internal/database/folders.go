package database

import (
	"context"

	"github.com/redat00/seacrate/internal/models"
)

// Check if folder already exist in database
func (d *DatabaseEngine) checkIfFolderExist(path string) bool {
	var count int64
	d.c.QueryRow(context.Background(), "SELECT COUNT(id) FROM folders WHERE path = $1", path).Scan(&count)
	if count > 0 {
		return true
	}
	return false
}

// Check if folder is empty
func (d *DatabaseEngine) checkIfFolderIsEmpty(path string) bool {
	var fCount int64
	var sCount int64

	// Check for folders
	d.c.QueryRow(context.Background(), "SELECT COUNT(id) FROM folders WHERE parent_folder = $1", path).Scan(&fCount)
	if fCount > 0 {
		return false
	}

	// Check for secrets
	d.c.QueryRow(context.Background(), "SELECT COUNT(id) FROM secrets WHERE folder = $1", path).Scan(&sCount)
	if sCount > 0 {
		return false
	}

	return true
}

// Create a new folder within database
func (d *DatabaseEngine) createFolder(path string, parentFolder string) error {
	_, err := d.c.Exec(context.Background(), "INSERT INTO folders (path, parent_folder) VALUES ($1, $2)", path, parentFolder)
	return err
}

// Create folder that are succeding each other
//
// Note for future me, or future whoever :
//
// This could probably be improved by the use of either transactions or,
// since we're using PostgreSQL, a PostgreSQL function. Maybe even both.
// However at the time of writing this, I'm not quite up there with my
// PostgreSQL skills, so everything is done in the code.
func (d *DatabaseEngine) createMultipleFolders(folders []string) error {
	for i, folder := range folders {
		if !d.checkIfFolderExist(folder) {
			// The first folder is always at root
			if i == 0 {
				err := d.createFolder(folder, "/")
				if err != nil {
					return err
				}
			} else {
				err := d.createFolder(folder, folders[i-1])
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// Get folder from database by full path
func (d *DatabaseEngine) getFolder(path string) (models.Folder, error) {
	var folder models.Folder
	err := d.c.QueryRow(context.Background(), "SELECT id, path, parent_folder, created_at FROM folders WHERE path = $1", path).Scan(&folder.Id, &folder.FullPath, &folder.ParentFolder, &folder.CreatedAt)
	return folder, err
}

// Get a list of folders within a folder
func (d *DatabaseEngine) getFoldersInFolder(path string) ([]models.FolderContent, error) {
	var folders []models.FolderContent
	folderRows, err := d.c.Query(context.Background(), "SELECT path, created_at FROM folders WHERE parent_folder = $1", path)
	if err != nil {
		return folders, err
	}
	defer folderRows.Close()
	for folderRows.Next() {
		var fR models.FolderContent
		fR.Type = "folder"
		err := folderRows.Scan(&fR.Key, &fR.CreatedAt)
		if err != nil {
			return folders, err
		}
		folders = append(folders, fR)
	}
	return folders, nil
}

// Delete folder from database
func (d *DatabaseEngine) deleteFolder(path string) error {
	_, err := d.c.Exec(context.Background(), "DELETE FROM folders WHERE path = $1", path)
	return err
}

// Delete folder if it's empty and go recursively
//
// This function will only delete a folder if it's empty. On top of that, it will also
// recursevely delete the parent folder of the folder we're deleting if it's also empty.
// This will be done all the way to the root.
//
// This code might be improved by the use of PostgreSQL functions, for now that's what I
// felt was a good thing to do.
func (d *DatabaseEngine) deleteFolderIfEmptyRecursively(path string) error {
	for true {
		// Prevent root folder deletion
		if path == "/" {
			break
		}

		folder, err := d.getFolder(path)
		if err != nil {
			return err
		}
		if empty := d.checkIfFolderIsEmpty(folder.FullPath); empty {
			d.deleteFolder(folder.FullPath)
		} else {
			break
		}
		path = folder.ParentFolder
	}
	return nil
}
