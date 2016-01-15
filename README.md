IOpipe
---------------------------------------
Apache 2.0 licensed.

IOpipe simplifies the consumption and integration of web services through
the chaining of pipescripts, javascript-based microservices.

Pipescripts take and transform input, providing straight-forward output
in a fashion to Unix pipes. Scripts may receive input or send output to/from
web service requests, functions, or local applications.

Scripts may be embedded in applications, used from shell scripts, or run manually
via a CLI. Because pipescript is Javascript, embedding is possible with most
languages and made safe through sandboxing.

![Build Status](https://circleci.com/gh/iopipe/iopipe.png?circle-token=eae431abda6b19dbfca597af818bb01092211272)

---------------------------------------
Usage
---------------------------------------

### CLI binary installation using Docker:

```bash
$ docker pull iopipe/iopipe:trunk
$ eval $(echo "alias iopipe='docker run --rm -it iopipe/iopipe:trunk'" | tee -a ~/.bashrc)
$ iopipe --help
```

### Command-line

```sh
# Import a pipescript and name it com.example.SomeScript
$ iopipe import --name com.example.SomeScript - <<<'input'

# List pipescripts
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
```

### NodeJS SDK:

The NodeJS SDK provides a generic callback chaining mechanism which allows
mixing HTTP(S) requests/POSTs, function calls, and pipescripts. Callbacks
receive the return of the previous function call or HTTP body.

```javascript
var iopipe = require("iopipe")

// Where com.example.SomeScript is present in .iopipe/filter_cache/
iopipe.exec("http://localhost/get-request",
            "com.example.SomeScript",
            "http://otherhost.post")

// Users may chain functions and HTTP requests.
iopipe.exec(function() { return "something" },
            function(arg) { return arg },
            "http://otherhost.post",
            your_callback)

// A function may also be returned then executed later.
var f = iopipe.define("http://fetch", "https://post")
f()
```

For more information on using the NodeJS SDK, please refer to its documentation:
***https://github.com/iopipe/iopipe/blob/master/docs/nodejs.md***

---------------------------------------
Filters & Pipescript
---------------------------------------

Requests and responses and translated using filters written in
Pipescript (i.e. Javascript) or offered as web services.

All filters simply receive request or response data and output
translated request or response data. Pipescript is typically operated
upon locally in the client, whereas web-service based filters operate
server-side. Pipescript may also be used to build serverside filters
and applications.

Example:

```javascript
export.module = function(input) {
  return "I'm doing something with input: {0}".format(input)
}
```

For more on writing filters see:
***https://github.com/iopipe/iopipe/blob/master/docs/pipescript.md***

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
$ eval $(echo "alias iopipe='docker run --rm -it iopipe-dev'" | tee -a ~/.bashrc)
$ iopipe --help
```

---------------------------------------
Project goals
---------------------------------------

The principal goal of our project is to improve
human to machine and machine to machine communications
with a free, highly-distributed protocol.

Furthermore:

1. Support simple translation of existing APIs
2. Support native, greenfield applications
3. Design for an open and distributed web.
4. Permissive open source licensing.
5. Security. Security. Security.

---------------------------------------
Scaling
---------------------------------------

Filters are pulled from (de)centralized repositories
and scale should be considered when building and
deploying private filter repositories.

---------------------------------------
Security
---------------------------------------

Note that this tool communicates to 3rd-party
web services. Caution is advised when trusting
these services, as is standard practice with
web and cloud services.

Pipescripts are executed in individual
javascript VMs whenever allowed by the host
system. Users should still exercise caution
when running pipescripts provided by other
users.

Contact security@iopipe.com for questions.

---------------------------------------
LICENSE
---------------------------------------

Apache 2.0
