package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/dhowden/tag"
)

const (
	root     = `D:\music\takeout\Takeout\Google Play Music\トラック`
	distRoot = `D:\music\takeout\Takeout\Google Play Music\dist`
)

type Album struct {
	MaxTrack int
	Mp3s     []Mp3
}

type Mp3 struct {
	Name  string
	Track int
	Path  string
}

func main() {
	pattern := filepath.Join(root, `*.mp3`)
	paths, err := filepath.Glob(pattern)
	if err != nil {
		log.Fatal(err)
	}

	// paths = paths[1000:1200]

	albums := map[string]Album{}

	for _, path := range paths {
		f, err := os.Open(path)
		m, err := tag.ReadFrom(f)
		if err != nil {
			log.Fatal(err)
		}
		albumTitle := m.Album()
		title := m.Title()
		track, _ := m.Track()
		f.Close()

		mp3 := Mp3{
			Name:  title,
			Track: track,
			Path:  path,
		}

		album, ok := albums[albumTitle]
		if !ok {
			album := Album{
				MaxTrack: 0,
				Mp3s:     []Mp3{},
			}
			albums[albumTitle] = album
		}
		if album.MaxTrack < track {
			album.MaxTrack = track
		}
		album.Mp3s = append(album.Mp3s, mp3)
		albums[albumTitle] = album

		// log.Print(m.Title())
		// log.Print(m.Album())
		// log.Print(m.Track())

		// log.Print(m.Format())
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
	dists := map[string]bool{}

	for title, album := range albums {
		if title == `` {
			continue
		}
		// fmt.Printf("アルバム名：%s トラック数：%d\n", title, len(album.Mp3s))

		sort.Slice(album.Mp3s, func(i, j int) bool { return album.Mp3s[i].Track < album.Mp3s[j].Track })
		format := `%02d %s.mp3`
		if album.MaxTrack >= 100 {
			format = `%03d %s.mp3`
		}
		for _, mp3 := range album.Mp3s {
			name := fmt.Sprintf(format, mp3.Track, mp3.Name)
			name = strings.Replace(name, "?", "？", -1)
			name = strings.Replace(name, "*", "＊", -1)
			name = strings.Replace(name, ":", "：", -1)
			// fmt.Println(name)
			dir := filepath.Join(distRoot, Escape(title))
			MkDir(dir)

			dist := filepath.Join(dir, Escape(name))
			if _, ok := dists[dist]; ok {
				// fmt.Printf("重複：%s\n", dist)
			} else {
				dists[dist] = true
				if Exists(dist) {
					continue
				}
				err := os.Link(mp3.Path, dist)
				if err != nil {
					log.Println(err)
				}
			}

		}
	}

}

func Escape(str string) string {
	str = strings.Replace(str, `\`, "", -1)
	str = strings.Replace(str, "￥", "", -1)
	str = strings.Replace(str, "／", "", -1)
	str = strings.Replace(str, "<", "", -1)
	str = strings.Replace(str, ">", "", -1)
	str = strings.Replace(str, "|", "", -1)
	str = strings.Replace(str, `"`, "", -1)
	str = strings.Replace(str, "?", "？", -1)
	str = strings.Replace(str, "*", "＊", -1)
	str = strings.Replace(str, ":", "：", -1)
	str = strings.Replace(str, "!", "！", -1)
	str = strings.Replace(str, "#", "＃", -1)
	str = strings.Replace(str, "/", " ", -1)
	str = strings.Replace(str, "...", "…", -1)
	return str
}

// MkDir ディレクトリが存在しなければ作成する
func MkDir(path string) (err error) {
	fi, err := os.Stat(path)
	if err == nil {
		if fi.IsDir() {
			// Nothing to do. It already exists.
		} else {
			return errors.New("it already exists and isn't directory. Cannot recover")
		}
	} else {
		err = os.Mkdir(path, 0777)
		if err != nil {
			return
		}
	}
	return nil
}

// Exists ファイルまたはディレクトリの存在を確認する
func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
