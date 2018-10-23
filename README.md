# impress

**Cross platform GUI Library for Go**

## Proof of Concept Version

Notes:

- This project is still in the early stages of development and is not yet in a usable state.
- The project tested on Debian 9 and MacOS 10.13.

## Basic Principles of Library Building Design

- Performance of the application as well as native applications
- A convenient functions interface to simplify application program structure
- Limited use of other libraries (graphic, low-level or native)
- Creating a application without form designer or hard-coded widgets (it can be separate libraries)

## Examples

Let's say hello:

!["Hello, world!" screenshot](https://codeation.github.io/pages/images/helloworld.png)

```
package main

import (
	"log"

	"github.com/codeation/impress"
	_ "github.com/codeation/impress/duo"
)

func main() {
	a := impress.NewApplication()
	a.Title("Hello World Application")
	a.Size(impress.NewRect(0, 0, 640, 480))

	font, err := impress.NewFont(`{"family":"Verdana"}`, 15)
	if err != nil {
		log.Fatal(err)
	}
	defer font.Close()

	w := a.NewWindow(impress.NewRect(0, 0, 640, 480), impress.NewColor(255, 255, 255))
	w.Text("Hello, world!", font, impress.NewPoint(280, 210), impress.NewColor(0, 0, 0))
	w.Line(impress.NewPoint(270, 230), impress.NewPoint(380, 230), impress.NewColor(255, 0, 0))
	w.Show()

	for a.Event() != impress.DestroyEvent {
	}
	a.Quit()
}
```

## Installation

The dependency Go libraries can be installed with the following commands:

```
go get github.com/codeation/impress
```

Currently, the library uses [a separate application](https://github.com/codeation/it) as a GUI driver
for rendering, event collecting, etc. You can [download](https://github.com/codeation/it/releases)
the compiled binary file or make it again from [the source](https://github.com/codeation/it).

You can specify the full path and name for the GUI driver via the environment variable, for example:

```
IMPRESS_TERMINAL_PATH=/path/it
```

or just copy the downloaded GUI driver to the working directory.

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

## Contributing

First of all, welcome:

- any advice on the library design and principles
- help to correct grammatical and writing errors
- contribution in the near future

## Using

A cross-platform [mind-map application](https://github.com/codeation/lineation/) is being developed to check the library's applicability.

Stay tuned for more.
