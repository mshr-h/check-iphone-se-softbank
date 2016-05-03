package main

import (
	"fmt"
	"image"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/disintegration/imaging"

	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/paint"
)

const (
	sb_url   = "http://yoyakuget.com/img/sb40.jpg"
	filename = "sb40.jpg"
)

func downloadImage(url string) {
	resp, err := http.Get(url)

	if err != nil {
		log.Fatal(err)
	}

	contents, err := ioutil.ReadAll(resp.Body)

	file, err := os.Create(filename)
	defer file.Close()

	if err != nil {
		log.Fatal(err)
	}

	file.Write(contents)
}

func main() {

	driver.Main(func(s screen.Screen) {
		downloadImage(sb_url)
		img, err := imaging.Open("sb40.jpg")
		if err != nil {
			panic(err)
		}
		bufSize := img.Bounds().Size()

		winSize := screen.NewWindowOptions{bufSize.X + 15, bufSize.Y + 35}
		w, err := s.NewWindow(&winSize)
		if err != nil {
			log.Fatal(err)
		}
		defer w.Release()

		b, err := s.NewBuffer(bufSize)
		if err != nil {
			log.Fatal(err)
		}
		defer b.Release()
		copyImage(b.RGBA(), img)

		for {
			e := w.NextEvent()
			switch e := e.(type) {
			case lifecycle.Event:
				if e.To == lifecycle.StageDead {
					return
				}

			case key.Event:
				if e.Code == key.CodeF5 {
					downloadImage(sb_url)
					img, err := imaging.Open("sb40.jpg")
					if err != nil {
						panic(err)
					}
					copyImage(b.RGBA(), img)
				} else if e.Code == key.CodeEscape {
					return
				}

			case paint.Event:
				w.Upload(image.Point{}, b, b.Bounds())
				w.Publish()

			case error:
				log.Print(e)
			}
		}
	})
}

func copyImage(dst *image.RGBA, src image.Image) error {
	b := dst.Bounds()
	if b != src.Bounds() {
		return fmt.Errorf("image size not match")
	}

	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			dst.Set(x, y, src.At(x, y))
		}
	}
	return nil
}
