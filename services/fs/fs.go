package fs

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"vobook/config"
)

type FileType uint8

const (
	_ FileType = iota
	FileTypeImage
	FileTypeVideo
	FileTypeAudio
	FileTypeOther
)

var ImageTypes = []string{"png", "jpg", "jpeg", "svg", "gif"}
var VideoTypes = []string{"mpeg", "mp4", "mov", "webm", "mkv", "flv", "gifv", "avi", "wmv"}
var AudioTypes = []string{"mp3", "ogg", "aac", "wav"}

// Type determines type of file from filename by its extension
func Type(filename string) FileType {
	ext := filepath.Ext(filename)

	// compares two extensions
	// first one is predefined of type, dot at beginning should be added
	// second is from filename and should be lowercased
	isExtEqual := func(ext1, ext2 string) bool {
		ext1 = "." + ext1
		ext2 = strings.ToLower(ext2)
		return ext1 == ext2
	}

	if ext == "" {
		return FileTypeOther
	}

	for _, imageExt := range ImageTypes {
		if isExtEqual(imageExt, ext) {
			return FileTypeImage
		}
	}

	for _, videoExt := range VideoTypes {
		if isExtEqual(videoExt, ext) {
			return FileTypeVideo
		}
	}

	for _, audioExt := range AudioTypes {
		if isExtEqual(audioExt, ext) {
			return FileTypeAudio
		}
	}

	return FileTypeOther
}

// Store writes given file bytes in storage
func Store(name string, b []byte) (filename string, err error) {
	t := time.Now().UTC()
	dateFolder := fmt.Sprintf("%d-%02d", t.Year(), t.Month())
	filename, err = UniqueFilename(path.Join(dateFolder, name))
	if err != nil {
		return
	}

	err = os.MkdirAll(path.Join(config.Get().FileStorage.Dir, dateFolder), 0755)
	if err != nil {
		err = fmt.Errorf("fs.Store: MkdirAll: %s", err.Error())
		return
	}
	err = ioutil.WriteFile(FullPath(filename), b, 0755)
	if err != nil {
		err = fmt.Errorf("fs.Store: WriteFile: %s", err.Error())
	}
	return
}

// Delete deletes file
func Delete(filename string) (err error) {
	return os.Remove(FullPath(filename))
}

// UniqueFilename creates unique filename using incremental numeric suffix
func UniqueFilename(filename string) (string, error) {
	dir, fname := path.Split(filename)
	ext := filepath.Ext(fname)
	name := strings.TrimSuffix(fname, ext)
	next := 1
	for {
		_, err := os.Stat(FullPath(filename))
		if os.IsNotExist(err) {
			return filename, nil
		}
		if err != nil {
			return "", fmt.Errorf("fs.UniqueFilename: %s", err.Error())
		}

		nextFilename := fmt.Sprintf("%s-%d%s", name, next, ext)
		filename = path.Join(dir, nextFilename)
		next++
	}
}

func FullPath(p string) string {
	return path.Join(config.Get().FileStorage.Dir, p)
}

// Base64 returns file bytes encoded in base64
func Base64(filename string) (b64 string, err error) {
	f, err := os.Open(FullPath(filename))
	if err != nil {
		return
	}

	reader := bufio.NewReader(f)
	content, err := ioutil.ReadAll(reader)
	if err != nil {
		return
	}

	b64 = base64.StdEncoding.EncodeToString(content)
	return
}
