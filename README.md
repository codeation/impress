# impress

**Cross platform GUI Library for Go**

## Proof of Concept Version

Notes:

- This project is still in the early stages of development and is not yet in a usable state.
- The project tested on Debian 9 and MacOS 10.13.

## Basic Principles of Library Building Design

- Performance of the application as well as native applications
- A convenient functions interface to simplify application program structure
- Using Go Project libraries to draw graphics ("image", "golang.org/x/image", etc)
- Limited use of other libraries (graphic, low-level or native)
- Creating a application without form designer or hard-coded widgets (it can be separate libraries)

## Examples

Let's say hello:

!["Hello, world!" screenshot](https://codeation.github.io/pages/images/helloworld.png)

```
package main

import (
	"github.com/codeation/impress"
	"github.com/codeation/impress/bitmap"
	"golang.org/x/image/font/gofont/goregular"
	"log"
)

func main() {
	a := impress.NewApplication()
	a.Title("Hello World Application")
	a.Size(bitmap.NewRect(0, 0, 640, 480))

	font, err := bitmap.NewFont(goregular.TTF, 15)
	if err != nil {
		log.Println(err)
		return
	}
	defer font.Close()

	w := a.NewWindow(bitmap.NewRect(0, 0, 640, 480), bitmap.NewColor(255, 255, 255))
	w.Text("Hello, world!", font, bitmap.NewPoint(280, 225), bitmap.NewColor(0, 0, 0))
	w.Line(bitmap.NewPoint(270, 230), bitmap.NewPoint(380, 230), bitmap.NewColor(255, 0, 0))
	w.Show()

	go func() {
		for {
			switch a.Event() {
			case impress.DestroyEvent:
				a.Quit()
			}
		}
	}()

	a.Main()
}
```

## Installation

The dependency Go libraries can be installed with the following commands:

```
go get github.com/codeation/impress
```

Currently, the library uses [GTK+ 3](https://www.gtk.org) for rendering, event collecting, etc. You should install `libgtk+-3.0` and packages that depend on GTK.

On Debian/ Ubuntu you can run:

```
sudo apt-get install libgtk-3-dev
```

Also [pkg-config](https://www.freedesktop.org/wiki/Software/pkg-config/) must be installed.

## Contributing

First of all, welcome:

- any advice on the library design and principles
- help to correct grammatical and writing errors
- contribution in the near future

## Using

A cross-platform [mind-map application](https://github.com/codeation/lineation/) is being developed to check the library's applicability.

Stay tuned for more.
