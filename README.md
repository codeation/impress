# impress. Go GUI cross-platform library

[![PkgGoDev](https://pkg.go.dev/badge/github.com/codeation/impress)](https://pkg.go.dev/github.com/codeation/impress)

See [project site](https://codeation.github.io/impress/) for a technical details
and [a library overview](https://codeation.github.io/impress/library-overview.html).

Some usage examples are in the [examples folder](https://github.com/codeation/impress/tree/master/examples).

## Hello World Example

Let's say hello:

<img src="https://codeation.github.io/images/hello_small.png" width="545" height="350" />

```
package main

import (
    "image"
    "image/color"

    "github.com/codeation/impress"
    "github.com/codeation/impress/event"

    _ "github.com/codeation/impress/duo"
)

func main() {
    app := impress.NewApplication(image.Rect(0, 0, 480, 240), "Hello World Application")
    defer app.Close()

    font := app.NewFont(15, map[string]string{"family": "Verdana"})
    defer font.Close()

    w := app.NewWindow(image.Rect(0, 0, 480, 240), color.RGBA{255, 255, 255, 255})
    defer w.Drop()

    w.Text("Hello, world!", font, image.Pt(200, 100), color.RGBA{0, 0, 0, 255})
    w.Line(image.Pt(200, 120), image.Pt(300, 120), color.RGBA{255, 0, 0, 255})
    w.Show()
    app.Sync()

    for {
        e := <-app.Chan()
        if e == event.DestroyEvent || e == event.KeyExit {
            break
        }
    }
}
```

See an explanation of the source code in [a library overview](https://codeation.github.io/impress/library-overview.html).

### To run this example on Debian/ Ubuntu:

0. Install `gcc`, `make`, `pkg-config` if you don't have them installed.

1. Install GTK+ 3 libraries if you don't have them installed:

```
sudo apt-get install libgtk-3-dev
```

2. Build impress terminal from source:

```
git clone https://github.com/codeation/it.git
cd it
make
cd ..
```

3. Then run example:

```
git clone https://github.com/codeation/impress.git
cd impress
IMPRESS_TERMINAL_PATH=../it/it go run ./examples/simple/
```

Steps 0-2 are needed to build a impress terminal binary. See [impress terminal page](https://codeation.github.io/impress/it-driver.html) for other options for downloading or building impress terminal app.

## Technical details

Basic Principles of Library Design:

- Performance of the application as well as native applications.
- Facade pattern for masking complex components behind an API.
- Simple and clean application code.
- Minimal library API to follow the Go-way.
- Limited use of other libraries (graphic, low-level or native).
- Creating a GUI application without form designer or hard-coded widgets.

The main idea is to stay away from the event driven programming paradigm. See
["What's wrong with event-driven programming"](https://codeation.github.io/impress/what-is-wrong-with-event-oriented-programming.html)
page.

The library uses 
[a separate application (GTK+ 3 terminal)](https://codeation.github.io/impress/it-driver.html) for drawing instead of binding low-level library to a Golang.

## Project State

- The project is currently in its beta stage. It is highly suitable for the development of in-house applications.
- The project was tested on Debian 12.9 and macOS 15.3.
- While the API remains stable, please note that the specific details may be subject to change.

[The project roadmap](https://codeation.github.io/impress/roadmap.html)
includes both short-term and long-term project stages. 

A cross-platform [mind-map application](https://codeation.github.io/lineation/) is being developed to showcase the core principles of the library.

Up to 2 ms from the moment you press the key, the mind-map application window is redrawn. Measured on a computer with an N200 processor without a video card.

## Contributing

First, welcome:

- Any like to the project (star, post, link, etc.).
- Any recommendations about library design principles.
- Contribution to project, documentation, examples.
- Bug report, open [issue](https://github.com/codeation/impress/issues).
- Or any help to fix grammatical or writing errors (PR or issue).
