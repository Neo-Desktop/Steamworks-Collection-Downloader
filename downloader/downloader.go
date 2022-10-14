package downloader

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Neo-Desktop/Steamworks-Collection-Downloader/manifest"
	"github.com/Neo-Desktop/Steamworks-Collection-Downloader/webapi"

	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

type Downloader struct {
	apikey  string
	prefix  string
	Fetched map[string]*Downloadable
}

func (d *Downloader) Start(apiKey string, seeds []string, prefix string) error {
	d.apikey = apiKey
	d.prefix = prefix

	err, r := d.RequestFileDetails(apiKey, seeds)
	if err != nil {
		return err
	}

	for _, v := range seeds {
		d.AddSeed(v)
	}

	err = d.FetchCollections(r)
	return err
}

func (d *Downloader) AddSeed(id string) {
	if d.Fetched == nil {
		d.Fetched = make(map[string]*Downloadable)
	}
	d.Fetched[id] = new(Downloadable)
}

func (d *Downloader) FetchCollections(r *webapi.PublishedFileDetailsResponse) error {
	var err error

	if d.apikey == "" {
		return errors.New("API key not set")
	}

	for {
		for _, v := range r.Response.PublishedFileDetails {
			d.Fetched[v.PublishedFileID].fetched = true

			log.Printf("%s - has children (%t)\n", v.Title, v.NumChildren > 0)
			if v.NumChildren > 0 {
				for _, c := range v.Children {
					if _, ok := d.Fetched[c.PublishedFileID]; !ok {
						d.Fetched[c.PublishedFileID] = &Downloadable{fetched: false}
					}
				}
			}

			if len(v.Filename) > 3 && v.Filename[len(v.Filename)-3:] == "vpk" {
				d.Fetched[v.PublishedFileID] = &Downloadable{
					FileURL:    v.FileURL,
					PreviewURL: v.PreviewURL,
					Title:      v.Title,
					Entry: manifest.Entry{
						ID:                  v.PublishedFileID,
						TimeCreated:         fmt.Sprintf("%d", v.TimeCreated),
						TimeUpdated:         fmt.Sprintf("%d", v.TimeUpdated),
						HContentFile:        v.HContentFile,
						HContentFileSize:    v.FileSize,
						HContentPreview:     v.HContentPreview,
						HContentPreviewSize: v.PreviewFileSize,
					},
					fetched: true,
				}
			}
		}

		children := make([]string, 0)
		for k, v := range d.Fetched {
			if v.fetched == false {
				children = append(children, k)
			}
		}

		if len(children) == 0 {
			break
		}

		err, r = d.RequestFileDetails(d.apikey, children)
		if err != nil {
			return err
		}
	}

	for k, v := range d.Fetched {
		if v.FileURL == "" {
			delete(d.Fetched, k)
		}
	}

	return nil
}

func (d *Downloader) RequestFileDetails(key string, fileIDs []string) (error, *webapi.PublishedFileDetailsResponse) {
	client := new(http.Client)

	req, err := http.NewRequest(http.MethodGet, "", nil)
	if err != nil {
		return err, nil
	}

	req.URL.Scheme = "https"
	req.URL.Host = "api.steampowered.com"
	req.URL.Path = "/IPublishedFileService/GetDetails/v1"

	query := req.URL.Query()
	query.Add("includechildren", "true")
	query.Add("key", key)

	for i, v := range fileIDs {
		query.Add(fmt.Sprintf("publishedfileids[%d]", i), v)
	}

	req.URL.RawQuery = query.Encode()

	log.Printf("Requesting: %s - sleeping 5 seconds\n", req.URL)
	time.Sleep(5 * time.Second)

	res, err := client.Do(req)
	if err != nil {
		return err, nil
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err, nil
	}

	r := new(webapi.PublishedFileDetailsResponse)

	err = json.Unmarshal(body, r)
	if err != nil {
		return err, nil
	}

	return nil, r
}

func (d *Downloader) DisplayFetched() {
	keys := make(map[string]string, 0)
	for k, v := range d.Fetched {
		keys[v.Title] = k
	}

	k2 := maps.Keys(keys)
	slices.Sort(k2)

	for _, k := range k2 {
		v := d.Fetched[keys[k]]
		fmt.Printf("%s: %s (%s) [%s]\n", v.Title, v.FileURL, v.FileNamePrefixed(d.prefix), v.PreviewFileNamePrefixed(d.prefix))
	}
}

func (d *Downloader) IsContentUpdateNeeded(v *manifest.Entry, v2 *manifest.Entry) bool {
	return v.TimeUpdated != v2.TimeUpdated || v.TimeCreated != v2.TimeCreated ||
		v.HContentFile != v2.HContentFile || v.HContentFileSize != v2.HContentFileSize
}

func (d *Downloader) IsPreviewUpdateNeeded(v *manifest.Entry, v2 *manifest.Entry) bool {
	return v.TimeUpdated != v2.TimeUpdated || v.TimeCreated != v2.TimeCreated ||
		v.HContentPreview != v2.HContentPreview || v.HContentPreviewSize != v2.HContentPreviewSize
}

func (d *Downloader) UpdateFiles(m *manifest.Manifest) {
	for k, v := range d.Fetched {
		contentUpdate := false
		previewUpdate := false
		contentSizeOnDisk := int64(0)
		previewSizeOnDisk := int64(0)

		if v2, ok := m.Entries[k]; ok {
			contentStat, err := os.Stat(v.FileNamePrefixed(d.prefix))
			if err == nil {
				contentSizeOnDisk = contentStat.Size()
			} else if errors.Is(err, os.ErrNotExist) {
				contentUpdate = true
				fmt.Fprintf(os.Stderr, "could not find file %s\n", v.FileNamePrefixed(d.prefix))
			} else {
				contentUpdate = true
				fmt.Fprintf(os.Stderr, "could not stat() %s: %s\n", v.FileNamePrefixed(d.prefix), err)
			}

			previewStat, err := os.Stat(v.PreviewFileNamePrefixed(d.prefix))
			if err == nil {
				previewSizeOnDisk = previewStat.Size()
			} else if errors.Is(err, os.ErrNotExist) {
				previewUpdate = true
				fmt.Fprintf(os.Stderr, "could not find file %s\n", v.PreviewFileNamePrefixed(d.prefix))
			} else {
				previewUpdate = true
				fmt.Fprintf(os.Stderr, "could not stat() %s: %s\n", v.PreviewFileNamePrefixed(d.prefix), err)
			}

			if !contentUpdate && !d.IsContentUpdateNeeded(&v.Entry, v2) &&
				fmt.Sprintf("%d", contentSizeOnDisk) == v.HContentFileSize {
				log.Printf("%s.vpk (%s): already up to date\n", v.ID, v.Title)
			} else {
				contentUpdate = true
				if fmt.Sprintf("%d", previewSizeOnDisk) != v.HContentPreviewSize {
					log.Printf("%s.vpk (%s): expected size %s, have %d\n", v.ID, v.Title, v.HContentPreviewSize, previewSizeOnDisk)
				}
			}

			if !previewUpdate && !d.IsPreviewUpdateNeeded(&v.Entry, v2) &&
				fmt.Sprintf("%d", previewSizeOnDisk) == v.HContentPreviewSize {
				log.Printf("%s.jpg (%s): already up to date\n", v.ID, v.Title)
			} else {
				previewUpdate = true
				if fmt.Sprintf("%d", previewSizeOnDisk) != v.HContentPreviewSize {
					log.Printf("%s.jpg (%s): expected size %s, have %d\n", v.ID, v.Title, v.HContentPreviewSize, previewSizeOnDisk)
				}
			}
		}

		if previewUpdate {
			err := d.DownloadFile(v.PreviewURL, v.PreviewFileNamePrefixed(d.prefix), v.HContentPreviewSize)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err)
				continue
			}
		}

		if contentUpdate {
			err := d.DownloadFile(v.FileURL, v.FileNamePrefixed(d.prefix), v.HContentFileSize)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err)
				continue
			}
		}

		err := m.AddEntry(&v.Entry)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error adding entry for %s (id %s) to manifest\n", v.Title, v.ID)
			fmt.Fprintf(os.Stderr, "%s\n", err)
		}
	}
}

func (d *Downloader) DownloadFile(url string, filename string, filesize string) error {
	out, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0664)
	defer out.Close()
	if err != nil {
		return fmt.Errorf("error creating file %s\n%w\n", filename, err)
	}

	log.Printf("downloading %s to %s", url, filename)

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("error downloading file %s\n%w\n", filename, err)
	}

	n, err := io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("error streaming file to %s\n%w\n", filename, err)
	}

	if fmt.Sprintf("%d", n) != filesize {
		return fmt.Errorf("expected %s bytes, only recieved %d bytes\n", filesize, n)
	}

	return nil
}
