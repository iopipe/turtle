IOpipe
---------------------------------------
Apache 2.0 licensed.

IOpipe simplifies the consumption and integration of web services through
the chaining of kernels, single-function applications.

Kernels take and transform input, providing straight-forward output
in a fashion to Unix pipes. A kernel may receive input or send output to/from
web service requests, functions, or local applications.

IOpipe may be embedded in applications, used from shell scripts, or run manually
via a CLI to form complete applications. Kernels and pipelines may be run
within local processes, or dispatched to remote workers (i.e. "cloud").

![Build Status](https://circleci.com/gh/iopipe/iopipe.png?circle-token=eae431abda6b19dbfca597af818bb01092211272)
[![Coverage Status](https://coveralls.io/repos/github/iopipe/iopipe/badge.svg?branch=master&t=UYi1cn)](https://coveralls.io/github/iopipe/iopipe?branch=master)

---------------------------------------
Usage
---------------------------------------

### Installation:

Download the [latest binary release](https://github.com/iopipe/iopipe/releases) and chmod 755 the file.

Building from source? See [Build & Install from source](#build--install-from-source).

Alternatively, download & alias our Docker image:

```bash
$ docker pull iopipe/iopipe:trunk
$ docker run --name iopipe-data iopipe/iopipe:trunk
$ eval $(echo "alias iopipe='docker run --rm --volumes-from iopipe-data iopipe/iopipe:trunk'" | tee -a ~/.bashrc)
$ iopipe --help
```

OS-specific packages are forthcoming.

### Command-line

```sh
# Import a kernel and name it com.example.SomeScript
$ iopipe import --name com.example.SomeScript - <<<'input'

# List kernels
$ iopipe list

# Fetch response and process it with com.example.SomeScript
$ iopipe --debug exec http://localhost/some-request com.example.SomeScript

# Fetch response and convert it with SomeScript, sending the result to otherhost
$ iopipe --debug exec http://localhost/some-request com.example.SomeScript \
                      http://otherhost/request

# Fetch response and convert it with SomeScript, send that result to otherhost,
# & converting the response with the script ResponseScript
$ iopipe --debug exec http://localhost/some-request com.example.SomeScript \
                      http://otherhost/request some.example.ResponseScript

# Export an NPM module:
$ iopipe export --name my-module-name http://localhost/some-request com.example.SomeScript
```

### NodeJS SDK:

The NodeJS SDK provides a generic callback chaining mechanism which allows
mixing HTTP(S) requests/POSTs, function calls, and kernels. Callbacks
receive the return of the previous function call or HTTP body.

```javascript
var iopipe = require("iopipe")()

// Where com.example.SomeScript is present in .iopipe/filter_cache/
iopipe.exec("http://localhost/get-request",
            "com.example.SomeScript",
            "http://otherhost.post")

// Users may chain functions and HTTP requests.
iopipe.exec(function(_, ctx) { ctx.done("something") },
            function(arg, ctx) { ctx.done(arg) },
            "http://otherhost.post",
            your_callback)

// A function may also be returned then executed later.
var f = iopipe.define("http://fetch", "https://post")
f()

// A defined function also accepts parameters
var echo = require("iopipe-echo")
var f = iopipe.define(echo, console.log)
f("hello world")
```

For more information on using the NodeJS SDK, please refer to its documentation:
***https://github.com/iopipe/iopipe/blob/master/docs/nodejs.md***

---------------------------------------
Kernels
---------------------------------------

Requests and responses are translated using kernels, and
may pipe to other kernels, or to/from web service endpoints.

Kernels simply receive request or response data and output
translated request or response data.

Example:

```javascript
module.exports = function(input, context) {
  context.done("I'm doing something with input: {0}".format(input))
}
```

For more on writing filters see:
***https://github.com/iopipe/iopipe/blob/master/docs/kernels.md***

---------------------------------------
Build & Install from source
---------------------------------------

With a functioning golang 1.5 development environment:

```bash
$ go build
$ ./iopipe --help
```

Alternatively use Docker to build & deploy:

```bash
$ docker build -t iopipe-dev .
$ docker run --name iopipe-data iopipe-dev
$ eval $(echo "alias iopipe='docker run --rm --volumes-from iopipe-data iopipe-dev'" | tee -a ~/.bashrc)
$ iopipe --help
```

---------------------------------------
Project goals
---------------------------------------

The principal goal of our project is to improve
human to machine and machine to machine communication.
We believe this can be achieved without the creation
or use of new protocols, but through the use of
a flow-based programming model.

Furthermore:

1. Simplify the use and integration of existing APIs into
   user applications.
2. Support use by both existing and new applications.
3. Design for an open and distributed web.
4. Permissive open source licensing.
5. Secure sharing, execution, & communications.

---------------------------------------
Security
---------------------------------------

Note that this tool communicates to 3rd-party
web services. Caution is advised when trusting
these services, as is standard practice with
web and cloud services.

Kernels are executed in individual virtual machines
whenever allowed by the executing environment.
The definition of a virtual machine here is lax,
such that it may describe a Javascript VM,
a Linux container, or a hardware-assisted x86
virtual machine. Users should exercise caution
when running community created kernels.

It is a project priority to make fetching, publishing,
and execution of kernels secure for a
production-ready 1.0.0 release.

Contact security@iopipe.com for questions.

---------------------------------------
LICENSE
---------------------------------------

Apache 2.0
