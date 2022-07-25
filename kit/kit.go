package kit

import (
	"path/filepath"
	"strings"
)

func init() {

}

type FileType int

const (
	FT_gif FileType = 1
	FT_png FileType = 2
	FT_mp4 FileType = 3
	FT_mp3 FileType = 4
	FT_o   FileType = 0
)

func GetFileType(path string) FileType {
	ext := filepath.Ext(path)
	ext = strings.ToLower(ext)
	switch ext {
	case ".mp4":
		return FT_mp4
	case ".mp3":
		return FT_mp3
	case ".gif":
		return FT_gif
	case ".png", ".jpg":
		return FT_png
	}
	return FT_o
}

func GetStr(str *string) string {
	if str == nil {
		return ""
	} else {
		return *str
	}
}
