# impress. Go GUI cross-platform library

## Getting Started

[![PkgGoDev](https://pkg.go.dev/badge/github.com/codeation/impress)](https://pkg.go.dev/github.com/codeation/impress)

See [project wiki](https://github.com/codeation/impress/wiki) 
for [a library overview](https://github.com/codeation/impress/wiki/Library-overview).

Some usage examples are located in [examples folder](https://github.com/codeation/impress/tree/master/examples).

## Basic Principles of Library Design

- Performance of the application as well as native applications.
- Simple and clean application code.
- Limited use of other libraries (graphic, low-level or native).
- Creating a GUI application without form designer or hard-coded widgets.
- Minimal library API to stay within the Go-way.

The basic idea is to avoid the event-driven programming paradigm. See
["Whats wrong with event-driven programming"](https://github.com/codeation/impress/wiki/Whats-wrong-with-event-driven-programming)
page.

## Hello World Example

Let's say hello:

!["Hello, world!" screenshot](https://codeation.github.io/pages/images/helloworld.png)

```
package main

import (
	"image"
	"image/color"

	"github.com/codeation/impress"
	"github.com/codeation/impress/event"

	_ "github.com/codeation/impress/duo"
)

var (
	background = color.RGBA{255, 255, 255, 0}
	foreground = color.RGBA{0, 0, 0, 0}
	underline  = color.RGBA{255, 0, 0, 0}
)

func main() {
	app := impress.NewApplication(image.Rect(0, 0, 640, 480), "Hello World Application")
	defer app.Close()

	font := impress.NewFont(15, map[string]string{"family": "Verdana"})
	defer font.Close()

	w := app.NewWindow(image.Rect(0, 0, 640, 480), background)
	defer w.Drop()

	w.Text("Hello, world!", font, image.Pt(280, 210), foreground)
	w.Line(image.Pt(270, 230), image.Pt(380, 230), underline)
	w.Show()
	app.Sync()

	for {
		action := <-app.Chan()
		if action == event.DestroyEvent || action == event.KeyExit {
			break
		}
	}
}
```

See an explanation of the source code in [a library overview](https://github.com/codeation/impress/wiki/Library-overview).

To run this example on Debian/ Ubuntu:

```
git clone https://github.com/codeation/impress.git
cd impress

wget https://github.com/codeation/it/releases/download/v0.2.4/it-linux.zip
unzip it-linux.zip
rm it-linux.zip

go run ./examples/simple/
```

## GUI driver

The library uses [a separate application](https://github.com/codeation/it) for drawing
instead of binding low-level library to a Golang.

Pros:
- To compile the application, it is not necessary to install low-level libraries.
- There are no additional restrictions on the application, such as GTK functions should only be called by the main thread.
- It will be possible to run several applications for a low-level drawing from different devices, thereby using the screens of different devices for the same application.
- Abnormal termination of low-level drawing application may not result in data loss.

Cons:
- Some loss of speed due to data transfer between applications.
- Additional complexity due to state synchronization between applications.

You can [download](https://github.com/codeation/it/releases)
the compiled binary `it` file or make it again from [the source](https://github.com/codeation/it).

You can specify the full path and name for the GUI driver via the environment variable, for example:

```
IMPRESS_TERMINAL_PATH=/path/it go run ./examples/simple/
```

or just copy the downloaded GUI driver to the working directory and run example:

```
go run ./examples/simple/
```

## Project State

Notes:

- The project in a beta stage.
- The project tested on Debian 11.6 and macOS Big Sur.
- The library may contain bugs.

A cross-platform [mind-map application](https://github.com/codeation/lineation/) is being developed to prove the value of the library.

## Contributing

First of all, welcome:

- any clue of the project importance (stars, etc)
- bug report, open [issue](https://github.com/codeation/impress/issues)
- any advice on the library design and principles
- help to correct grammatical and writing errors (PR or issue)
- any contribution to project, project wiki, examples

Stay tuned for more.
