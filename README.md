# impress. Go GUI cross-platform library

[![PkgGoDev](https://pkg.go.dev/badge/github.com/codeation/impress)](https://pkg.go.dev/github.com/codeation/impress)

See [project wiki](https://github.com/codeation/impress/wiki) 
for [a library overview](https://github.com/codeation/impress/wiki/Library-overview).

Some usage examples are in the [examples folder](https://github.com/codeation/impress/tree/master/examples).

## Basic Principles of Library Design

- Performance of the application as well as native applications.
- Facade pattern for masking complex components behind an API.
- Simple and clean application code.
- Minimal library API to follow the Go-way.
- Limited use of other libraries (graphic, low-level or native).
- Creating a GUI application without form designer or hard-coded widgets.

The basic idea is to avoid the event-driven programming paradigm. See
["What's wrong with event-driven programming"](https://github.com/codeation/impress/wiki/Whats-wrong-with-event-driven-programming)
page.

Go is a perfect language for developing desktop GUI applications. Compiled language is fast enough to spin an interactive application. Goroutines are helpful to handle the state of separate windows. The small runtime means that the application starts instantly. Go implement high level abstractions to complex application development.

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

See an explanation of the source code in [a library overview](https://github.com/codeation/impress/wiki/Library-overview).

### To run this example on Debian/ Ubuntu:

1. Install GTK+ 3 libraries, `pkg-config` and some libraries if you don't have them installed:

```
sudo apt-get install libgtk-3-dev
sudo apt-get install pkg-config libsystemd-dev libwebp-dev
```

2. Build impress terminal from source (see [impress terminal page](https://github.com/codeation/it)):

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

### To run this example on macOS:

0. Install [Homebrew](https://brew.sh/) if you don't have it installed.
1. Install GTK+ 3 (`brew install gtk+3`) and pkg-config (`brew install pkg-config`) if you don't have them installed.
2. Build impress terminal from source (see [impress terminal page](https://github.com/codeation/it)):

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

*The latest releases aren't tested on Apple machines. The earlier version tested on both Intel and Silicon platform and worked well. Please, use earlier release or open [issue](https://github.com/codeation/impress/issues) if some bugs have raised. PRs and MRs are welcome too.*

## GTK-3 driver

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

*On Debian 12, `libharfbuzz-gobject` should be installed. If you are getting `error while loading shared libraries: libharfbuzz-gobject.so.0: cannot open shared object file: No such file or directory`, run `sudo apt-get install libharfbuzz-gobject0`.*

## Project State

### Notes

- The project is in a beta stage. It's suitable for developing in-house applications.
- The project tested on Debian 12.0.
- The API is stable, but the details are subject to change.

A cross-platform [mind-map application](https://github.com/codeation/lineation/) is being developed to prove the underlying principles of the library.

### Short term roadmap

- Recommended Event Propagation module.
- Additional library with a set of widgets (text input, dialog, etc.).
- Developer version of the library with debugging and profiling.
- Auto build GTK-3 driver for a fixed set of Linux distributions.
- Distribution tool to pack the GUI application into Debian package.

### Long term roadmap

- [WebAssembly driver](https://github.com/codeation/canvas) for turning a browser into a client terminal.
- Two or more cooperate drivers for the same GUI application.
- GTK-4 driver version when Debian will contain GTK-4 apps.
- Apple Silicon native driver for macOS.
- iOS app used as GUI app remote client.

Stay tuned for more.

## Contributing

First, welcome:

- Any like to the project (star, post, link, etc.).
- Any recommendations about library design principles.
- Contribution to project, project wiki, examples.
- Bug report, open [issue](https://github.com/codeation/impress/issues).
- Or any help to fix grammatical or writing errors (PR or issue).
