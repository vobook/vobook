package file

import (
	"bytes"
	"image"
	"image/png"
	"path"
	"vobook/config"
	"vobook/database"
	"vobook/database/models"
	"vobook/services/fs"

	"github.com/nfnt/resize"

	"github.com/go-pg/pg/v9"
	log "github.com/sirupsen/logrus"
)

func Create(file *models.File) (err error) {
	existing, err := GetByHash(file.Hash)
	if err != nil && err != pg.ErrNoRows {
		return
	}

	// new file
	if err == pg.ErrNoRows {
		var storedIn string
		storedIn, err = fs.Save(file.Filename, file.Bytes)
		if err != nil {
			return
		}

		file.Path = storedIn
	} else {
		file.Path = existing.Path
		file.PreviewPath = existing.PreviewPath
	}

	return database.ORM().Insert(file)
}

func CreatePreview(file *models.File) (err error) {
	// only for pics for now
	// TODO vids, text, music, whatever
	if file.Type != fs.FileTypeImage {
		file.PreviewPath = file.Path
		return Update(file)
	}

	data, err := fs.Load(file.Path)
	if err != nil {
		return
	}
	pic, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return
	}

	buff := new(bytes.Buffer)
	dim := config.Get().Contacts.Preview
	previewPic := resize.Resize(dim.W, dim.H, pic, resize.Lanczos3)
	err = png.Encode(buff, previewPic)
	if err != nil {
		return
	}

	_, fname := path.Split(file.Path)
	file.PreviewPath, err = fs.Save(fname, buff.Bytes())
	if err != nil {
		return
	}

	return database.ORM().Update(file)
}

func GetByID(id string) (elem models.File, err error) {
	err = database.ORM().
		Model(&elem).
		Where("id = ?", id).
		First()

	return
}

func GetByHash(hash string) (elem models.File, err error) {
	err = database.ORM().
		Model(&elem).
		Where("hash = ?", hash).
		First()

	return
}

func Delete(id string) (err error) {
	file, err := GetByID(id)
	if err != nil {
		return
	}

	_, err = database.ORM().Model(&models.File{}).
		Where("id=?", file.ID).
		Delete()
	if err != nil {
		return
	}

	// delete file from FS if no more entries referencing to it
	go func() {
		_, err = GetByHash(file.Hash)
		if err == pg.ErrNoRows {
			err = fs.Delete(file.Path)
			if err != nil {
				log.Errorf("err deleting file from fs: %s", err)
			}
			if file.PreviewPath != "" {
				err = fs.Delete(file.PreviewPath)
				if err != nil {
					log.Errorf("err deleting preview file from fs: %s", err)
				}
			}
		}
	}()

	return
}

func Restore(id string) (err error) {
	_, err = database.ORM().
		Model(&models.File{}).
		Where("id = ?", id).
		Set("deleted_at = null").
		Update()

	return
}

func Update(m *models.File) (err error) {
	err = database.ORM().
		Update(m)

	return
}
