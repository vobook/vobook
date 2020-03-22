package requests

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"vobook/cmd/server/errors"
	"vobook/database/models"
	"vobook/services/fs"
)

type FileBase64 struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Filename    string `json:"filename"`
	Data        string `json:"data"`
}

func (r *FileBase64) Validate() error {
	errs := errors.Input{}

	if r.Filename == "" {
		errs.Add("filename", "Filename missing")
	}
	if r.Data == "" {
		errs.Add("data", "Data missing")
	}

	if errs.Has() {
		return errs
	}
	return nil
}

func (r *FileBase64) ToModel() (m *models.File, err error) {
	m = &models.File{
		Name:        r.Name,
		Description: r.Description,
		Filename:    r.Filename,
		Base64:      r.Data,
	}

	if m.Name == "" {
		m.Name = m.Filename
	}

	m.Bytes, err = base64.StdEncoding.DecodeString(r.Data)
	if err != nil {
		return
	}

	sha := sha256.New()
	sha.Write(m.Bytes)
	m.Hash = hex.EncodeToString(sha.Sum(nil))

	m.Size = int64(len(m.Bytes))
	m.Type = fs.Type(m.Filename)

	return
}
