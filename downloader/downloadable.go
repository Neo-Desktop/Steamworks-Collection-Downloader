package downloader

import "github.com/Neo-Desktop/Steamworks-Collection-Downloader/manifest"

type Downloadable struct {
	PreviewFilename string
	FileURL         string
	PreviewURL      string
	Title           string
	Filesize        string
	manifest.Entry
	fetched bool
}
