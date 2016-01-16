# Introduction

IOpipe simplifies the consumption and integration of
web services. The NodeJS SDK allows the chaining of 
web service / HTTP requests, local functions, and
pipescripts.

Pipescripts are portable javascript snippets that take
a single string argument and return a string. These
pipescripts may be written and used locally, shared
amongst the community, or shared privately within a team.

Because pipescripts define their input and output types,
it is easy for developers to discover and build
request/response workflows which automatically build
API requests and transform responses.

# Installation

IOpipe for NodeJS may be downloaded & installed using *npm*:

```bash
$ npm install -g iopipe
```

Typically, users will use the command-line tool to create, download,
share, and manage pipescripts. The tool may be used to seed a filter cache
for embedding into your NodeJS project or to export complete  npm-compatible
packages.

# Basic Usage:

The following example demonstrates the use of IOpipe as a simple function and callback management mechanism:

```javascript
var iopipe = require("iopipe")

var mypipe = iopipe.define(
  function() {
    return "hello world"
  }
  ,console.log
)

mypipe()
```

The *exec* function exists for those not needing a reference to the function:

```javascript
var iopipe = require("iopipe")
iopipe.exec(
  function() {
    return "hello world"
  }
  ,console.log
)
```

# Integrating HTTP(S) requests

HTTP(S) may be placed anywhere in a pipeline. If a URL is detected
as the first argument, then an HTTP GET is performed. Otherwise, the
POST method is sent.

The following performs an HTTP GET and prints the output to the console:

```javascript
var iopipe = require("iopipe")
iopipe.exec("http://127.0.0.1/my_request/", console.log)
```

Manipulating a response and forwarding it to another server is easily done:

```javascript
var iopipe = require("iopipe")
iopipe.exec(
  "http://127.0.0.1/my_request/"
  ,function(s) { var j = JSON.decode(s); return j["field"] },
  ,"http://127.127.127.127/update/"
)
```

# Leveraging pipescript

Functions need not be inlined, in fact the greatest value of IOpipe
is in using stored pipescript routines. These allow sharing of functional,
lambda-like methods to transform requests.

Modifying the previous example to convert the inline function to a pipescript:

```bash
$ # Write a pipescript file via the shell:
$ mkdir -p .iopipe/filter_cache/
$ cat <<EOF >.iopipe/filter_cache/myscript
module.exports = function(input) {
  var x = JSON.decode(input)
  return x["field"]
}
EOF
```

```javascript
var iopipe = require("iopipe")
iopipe.exec(
  "http://127.0.0.1/my_request/"
  ,"myscript"
  ,"http://127.127.127.127/update/"
)
```
