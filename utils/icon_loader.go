package utils

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/storage"
)

func LoadIcon(name string) fyne.Resource {
	path := "assets/icons/" + name
	uri := storage.NewFileURI(path)
	read, err := storage.Reader(uri)
	if err != nil {
		return nil
	}
	defer read.Close()

	data := make([]byte, 0)
	buf := make([]byte, 1024)
	for {
		n, err := read.Read(buf)
		if n > 0 {
			data = append(data, buf[:n]...)
		}
		if err != nil {
			break
		}
	}

	return fyne.NewStaticResource(name, data)
}
