package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/dhowden/tag"
)

const (
	root = `D:\music\takeout\Takeout\Google Play Music\トラック`
)

func main() {
	pattern := filepath.Join(root, `*.mp3`)
	paths, err := filepath.Glob(pattern)
	if err != nil {
		log.Fatal(err)
	}
	for _, path := range paths {
		// path := filepath.Join(root, `FIELD OF VIEW - 1995 BEST 100 - 君がいたから.mp3`)
		f, err := os.Open(path)
		m, err := tag.ReadFrom(f)
		if err != nil {
			log.Fatal(err)
		}
		// log.Print(m.Format())
		log.Print(m.Title())
		log.Print(m.Album())
		log.Print(m.Track())
		// log.Print(m.AlbumArtist())
		// log.Print(m.Artist())
		// log.Print(m.Composer())
		// log.Print(m.Comment())
		// log.Print(m.Genre())
		// log.Print(m.Year())
		// log.Print(m.Disc())
		// log.Print(m.Picture())
		// log.Print(m.Lyrics())
	}

}
