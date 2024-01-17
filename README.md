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

    font := impress.NewFont(15, map[string]string{"family": "Verdana"})
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

## Project State

### Notes

- The project is currently in its beta stage. It is highly suitable for the development of in-house applications.
- The project was tested on Debian 12.4 and macOS 14.2.1.
- While the API remains stable, please note that the specific details may be subject to change.

A cross-platform [mind-map application](https://codeation.github.io/lineation/) is being developed to showcase the core principles of the library.

### Key features

- The library uses [a separate application](https://codeation.github.io/impress/it-driver.html) for drawing instead of binding low-level library to a Golang.
- The library implements [lazy mode](https://codeation.github.io/impress/lazy-mode.html) drawing.

See [project site](https://codeation.github.io/impress/) for a basic principles and details.

## Contributing

First, welcome:

- Any like to the project (star, post, link, etc.).
- Any recommendations about library design principles.
- Contribution to project, documentation, examples.
- Bug report, open [issue](https://github.com/codeation/impress/issues).
- Or any help to fix grammatical or writing errors (PR or issue).
