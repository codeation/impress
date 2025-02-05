package main

import (
	"image"
	"image/color"
	"log"

	"github.com/codeation/impress"
	"github.com/codeation/impress/duo/duodriver"
	"github.com/codeation/impress/event"
)

func main() {
	d1, err := duodriver.New()
	if err != nil {
		log.Fatal(err)
	}

	app1 := impress.MakeApplication(d1, image.Rect(0, 0, 480, 240), "Application 1")
	defer app1.Close()

	w1 := app1.NewWindow(image.Rect(0, 0, 480, 240), color.RGBA{255, 255, 255, 255})
	defer w1.Drop()

	w1.Fill(image.Rect(100, 50, 380, 190), color.RGBA{255, 0, 0, 255})
	w1.Show()
	app1.Sync()

	d2, err := duodriver.New()
	if err != nil {
		log.Fatal(err)
	}

	app2 := impress.MakeApplication(d2, image.Rect(0, 0, 480, 240), "Application 2")
	defer app2.Close()

	w2 := app2.NewWindow(image.Rect(0, 0, 480, 240), color.RGBA{255, 255, 255, 255})
	defer w2.Drop()

	w2.Fill(image.Rect(100, 50, 380, 190), color.RGBA{0, 0, 255, 255})
	w2.Show()
	app2.Sync()

	for {
		var e event.Eventer
		select {
		case e = <-app1.Chan():
		case e = <-app2.Chan():
		}

		if e == event.DestroyEvent || e == event.KeyExit {
			break
		}
	}
}
