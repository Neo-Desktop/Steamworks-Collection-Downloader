package manifest

type Entry struct {
	ID                  string
	TimeCreated         string
	TimeUpdated         string
	HContentFile        string
	HContentFileSize    string
	HContentPreview     string
	HContentPreviewSize string
}

func (w Entry) FileNamePrefixed(prefix string) string {
	return prefix + w.ID + ".vpk"
}

func (w Entry) PreviewFileNamePrefixed(prefix string) string {
	return prefix + w.ID + ".jpg"
}
