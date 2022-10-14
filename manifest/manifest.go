package manifest

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Manifest struct {
	filename string
	prefix   string

	Entries map[string]*Entry

	*os.File
}

func (m *Manifest) Open(filename string, prefix string) (err error) {
	filename = prefix + filename
	m.File, err = os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0664)
	if err != nil {
		return err
	}

	buffer, err := ioutil.ReadAll(m)
	if err != nil {
		return err
	}

	err = m.Parse(buffer)
	if err != nil {
		return err
	}

	return nil
}

func (m *Manifest) Parse(manifest []byte) error {
	m.Entries = make(map[string]*Entry)

	if len(manifest) == 0 {
		return nil
	}

	err := json.Unmarshal(manifest, &m.Entries)
	if err != nil {
		return err
	}

	return nil
}

func (m *Manifest) AddEntry(entry *Entry) error {
	m.Entries[entry.ID] = entry

	_, err := m.Seek(0, 0)
	if err != nil {
		return err
	}

	err = m.Truncate(0)
	if err != nil {
		return err
	}

	b, err := json.Marshal(m.Entries)
	if err != nil {
		return err
	}

	_, err = m.Write(b)
	if err != nil {
		return err
	}

	return nil
}
