IOpipe
---------------------------------------
[![Gitter](https://img.shields.io/gitter/room/nwjs/nw.js.svg?maxAge=2592000)](https://gitter.im/iopipe/iopipe)

Apache 2.0 licensed.

IOpipe is a toolkit for building and orchestrating event-driven and
serverless applications. These apps may run *locally* or in the cloud
via AWS Lambda, Google Cloud Functions, or Azure Functions.

IOpipe can:

 * Chain AWS Lambda Functions and local functions
 * Convert NodeJS functions into serverless functions
 * Perform GET and POST HTTP requests
 * Parallelize data into serverless workers.

We call our serverless functions "kernels".  Kernels take and transform
input and communicate over the networking, operating in a fashion to
Unix pipes. A kernel may receive input or send output to/from
web service requests, functions, or local applications.


---------------------------------------
Usage
---------------------------------------

### NodeJS SDK:

The NodeJS SDK provides a generic callback chaining mechanism which allows
mixing HTTP(S) requests/POSTs, function calls, and kernels. Callbacks
receive the return of the previous function call or HTTP body.

The callback variable received by a function is *also* an AWS Lambda-compatible
"context" object. Because of this, you can chain standard callback-based NodeJS
functions, and functions written for AWS Lambda.

```javascript
var iopipe = require("iopipe")()

/* Get HTTP data, process it with SomeScript, and POST the results.
   Note that com.example.SomeScript would be present in .iopipe/filter_cache/ */
iopipe.exec("http://localhost/get-request",
            "com.example.SomeScript",
            "http://otherhost.post")

// Users may chain functions and HTTP requests.
iopipe.exec(function(_, callback) { callback("something") },
            function(arg, callback) { callback(arg) },
            "http://otherhost.post",
            your_callback)

// A function may also be returned then executed later.
var f = iopipe.define("http://fetch", "https://post")
f()

// A defined function also accepts parameters
var echo = require("iopipe-echo")
var f = iopipe.define(echo, console.log)
f("hello world")

/* Create an AWS Lambda function from any NodeJS function /w callback.
   The callback becomes the equivilent of a done or success call on AWS. */
export.handler = iopipe.define(function(event, callback) {
  console.log(event)
  callback()
})

/* Of course, this method chaining also works for creating AWS Lambda code.
   This example will fetch HTTP data from the URL in the event's 'url' key
   and return a SHA-256 of the retrieved content. */
var crypto = require("crypto")
export.handler = iopipe.define(iopipe.property("url"),
                               iopipe.fetch,
                               (event, callback) => {
                                  callback(crypto
                                           .createHash('sha256')
                                           .update(event)
                                           .digest('hex'))
                               })
```

### AWS Lambda Client

IOpipe also acts as an AWS Lambda Client where a Lambda function may
be specified by its URN and included in the execution chain:

```javascript
var iopipe = require("iopipe")()
var iopipe_aws = require("iopipe")(
  exec_driver: 'aws'
  exec_driver_opts: {
    region: 'us-west-1',
    access_key: 'itsasecrettoeverybody',
    secret_key: 'itsasecrettoeverybody'
  }
)
var crypto = require("crypto")

export.handler = iopipe_aws.define("urn:somefunction",
                                   "urn:anotherfunction",
                                   iopipe.property("property-of-result"),
                                   iopipe.fetch, # fetch that as a URL
                                   (event, callback) => {
                                      callback(JSON.parse(event))
                                   },
                                   iopipe.map(
                                     iopipe_aws.define(
                                       "urn:spawn_this_on_aws_for_each_value_in_parallel"
                                     )
                                   ))
```

For more information on using the NodeJS SDK, please refer to its documentation:
***https://github.com/iopipe/iopipe/blob/master/docs/nodejs.md***

---------------------------------------
Kernel functions
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

Functions should expect a "context" parameter which may be called
directly as a callback, but also offers the methods 'done', 'success',
and 'fail'. Users needing, for any reason, to create a context manually
may call iopipe.create_context(callback).

For more on writing filters see:
***https://github.com/iopipe/iopipe/blob/master/docs/kernels.md***

### CLI

A Go-based CLI exists to create and export npm modules, share code,
and provide runtime of magnetic kernels.

Find this tool in the [IOpipe-Golang repo](https://github.com/iopipe/iopipe-golang).

---------------------------------------
Security
---------------------------------------

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

Modules are fetched and stored using sha256 hashes,
providing an advantage over module-hosting mechanisms
which are based simply on a name and version. Future
versions of IOpipe will likely implement TUF for
state-of-the-art software assurance.

Contact security@iopipe.com for questions.

---------------------------------------
LICENSE
---------------------------------------

Apache 2.0
