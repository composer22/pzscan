# pzscan
[![License MIT](https://img.shields.io/npm/l/express.svg)](http://opensource.org/licenses/MIT)
[![Build Status](https://travis-ci.org/composer22/pzscan.svg?branch=master)](http://travis-ci.org/composer22/pzscan)
[![Current Release](https://img.shields.io/badge/release-none-blue.svg)](https://github.com/composer22/pzscan/releases/tag/none)
[![Coverage Status](https://coveralls.io/repos/composer22/pzscan/badge.svg?branch=master)](https://coveralls.io/r/composer22/pzscan?branch=master)

A simple site scanner to validate links and content is SEO compliant.

## Description

This simple application will transverse a given URL and report back the following confirmations:

* All resources can be loaded on the page (CSS, js, images), and all links point to a working URL (non-4xx/non-5xx response).
* Pages have a canonical link tag.
* Pages have meta descriptions and each description is between 131 and 154 characters.
* Pages have title tags between 57 and 68 characters.
* Images have "alt" attributes.
* Pages are allowed only one "h1" tag.

Each element result is written to the log as INFO with a json encoded structure of the statistics of the scan.

## Usage

```

Usage: pzscan [options...]

Server options:
    -H, --hostname HOSTNAME          HOSTNAME to scan (default: example.com).
    -X, --procs MAX                  MAX processor cores to use from the
	                                 machine (default 1).
    -m, --minutes MAX                MAX minutes to live (default: 5).
    -W, --workers MAX                MAX running workers allowed (default: 4).

Common options:
    -h, --help                       Show this message

Example:

    # Scan craigslist; 1 processor; 2 min max; 10 worker go routines.

    ./pzscan -H "craigslist.com" -X 1 -m 2 -W 10

```

## Configuration

```
command line flags only

```

## Building

This code currently requires version 1.42 or higher of Go, but we encourage the use of the latest stable release.

Information on Golang installation, including pre-built binaries, is available at
<http://golang.org/doc/install>.  Stable branches of operating system packagers provided by
your OS vendor may not be sufficient.

Run `go version` to see the version of Go which you have installed.

Run `go build` inside the directory to build.

Run `go test ./...` to run the unit regression tests.

A successful build run produces no messages and creates an executable called `clidemo` in this
directory.  You can invoke that binary, with no options to start a server with acceptable standalone defaults.

Run `go help` for more guidance, and visit <http://golang.org/> for tutorials, presentations, references and more.


## License

(The MIT License)

Copyright (c) 2015 Pyxxel Inc.

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to
deal in the Software without restriction, including without limitation the
rights to use, copy, modify, merge, publish, distribute, sublicense, and/or
sell copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS
IN THE SOFTWARE.
