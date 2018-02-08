Turtle
---------------------------------------

[![Greenkeeper badge](https://badges.greenkeeper.io/iopipe/turtle.svg)](https://greenkeeper.io/)
[![Gitter](https://img.shields.io/gitter/room/nwjs/nw.js.svg?maxAge=2592000)](https://gitter.im/iopipe/iopipe)

Apache 2.0 licensed.

Turtle is a toolkit for building and orchestrating event-driven and
serverless applications. These apps may run anywhere, either locally or,
via execution drivers, in the cloud. It's turtles all the way down.

Execution drivers exist for:

 - AWS Lambda

Drivers are planned (or in development) for:

 - Google Cloud Functions
 - Azure Functions
 - Docker (Engine & Swarm)

Turtle can:

 * Chain AWS Lambda Functions and local functions.
 * Convert NodeJS functions into serverless functions.
 * Compose applications with HTTP APIs.
 * Parallelize data into serverless workers (scatter & gather).

# CLI

Use the [Turtle CLI](https://github.com/iopipe/turtle-golang) to create and
export npm modules, share code, & provide runtime of magnetic functions.
The CLI is still in early development, with our NodeJS SDK being more
mature.

Find, download, and/or contribute to this tool in the [CLI repo](https://github.com/iopipe/turtle-golang).

# SDK

### NodeJS SDK:

The NodeJS SDK provides a generic callback chaining mechanism which allows
mixing HTTP(S) requests/POSTs, and function calls. Callbacks
receive the return of the previous function call or HTTP body.

The callback variable received by a function is also an AWS Lambda-compatible
"context" object. Because of this, you can chain standard callback-based NodeJS
functions, and functions written for AWS Lambda.

#### Installation

Installation of the NodeJS module is easy via npm:

```
$ npm install @iopipe/turtle
```

Our CLI is still in early development and may be found on the releases
page, with further instructions in the [CLI
repo](https://github.com/iopipe/iopipe-golang).

#### Basic usage:

```javascript
/* Create a Lambda function which returns event.key + 1. */
var turtle = require("@iopipe/turtle")()

exports.handle = turtle.define(
  (event, context) => {
    context.succeed(event.key + 1)
  }
)
```

#### Context argument

The context argument operates as both a callback and
an object with several methods, similar to the same
argument passed to AWS Lambda functions.

Developers may call `context()` directly, with its argument
passed as the event to the next function, or may call its
methods.

Context Methods:

 - context.done(err, data)
 - context.succeed(data)
 - context.fail(err)


Example of using context.fail to pass errors:

```javascript
var turtle = require("@iopipe/turtle")()

exports.handle = turtle.define(
  (event, context) => {
    try {
      throw "Ford, you're turning into a penguin. Stop it!"
    }
    catch (err) {
      context.fail(err)
    }
  }
)
```

#### Function Composition

Turtle supports the composition of functions, HTTP endpoints,
and modules, taken from functional-programming and flow-based
programming models. This simplifies code-reuse and works as
glue between algorithms.

There is (some) compatibility with [Rambda](http://ramdajs.com) for
function composition & developing functional applications.

By using function composition, you will gain additional insights
and increased granularity when utilizing (upcoming) telementry features.

Example:

```javascript
/* Return event.int + 1, square the result,
   print, then return the result. */
exports.handle = turtle.define(
  (event, context) => {
    context(event.int + 1)
	},
  (event, context) => {
    context(Math.pow(event, 2))
	},
  (event, context) => {
    console.log(event)
    context(event)
	}
)
```

#### HTTP endpoints as "functions"

The first argument to `define`, if a URL, is fetched via an HTTP get
request. Any URL string specified elsewhere in the argument list to
`define` is sent a POST rqeuest.

This first example fetches data from a URL, then performs a POST request
to another.

```javascript
exports.handle = turtle.define("http://localhost/get-data",
                               "http://localhost/post-data")
```

It's possible to use Turtle to fetch from a URL and perform data
transformations via composition as follows:

```javascript
exports.handle = turtle.define(
  "http://localhost/get-data",
  (data, callback) => {
    console.log("Fetched data: " + data)
  }
)
```

Often, users will need to use Turtle to fetch a URL somewhere in
the middle of a composition and will need to use functional tools
such as `turtle.fetch`. The following example also uses
`turtle.property`, which extracts a key from an ECMAscript `Object`:

```javascript
exports.handle = turtle.define(
  turtle.property("url"),
  turtle.fetch,
  (data, callback) => {
    console.log("Fetched data: " + data)
  }
)
```

#### Scatter & Gather

Turtle also acts as a client to serverless infrastructure allowing
the use of scatter & gather patterns such as map-reduce.

Below we initialize an AWS Lambda Client where a Lambda function may
be specified by its Amazon [URN](https://en.wikipedia.org/wiki/Uniform_Resource_Name)
and included in the execution chain:

```javascript
var turtle = require("@iopipe/turtle")()
var turtle_aws = require("@iopipe/turtle")(
  exec_driver: 'aws'
  exec_driver_opts: {
    region: 'us-west-1',
    access_key: 'itsasecrettoeverybody',
    secret_key: 'itsasecrettoeverybody'
  }
)
var crypto = require("crypto")

export.handler = turtle_aws.define("urn:someLambdaFunction",
                                   "urn:anotherLambdaFunction",
                                   turtle.property("property-of-result"),
                                   turtle.fetch, // fetch that as a URL
                                   (event, callback) => {
                                      callback(JSON.parse(event))
                                   },
                                   turtle.map(
                                     iopipe_aws.define(
                                       "urn:spawn_this_on_aws_for_each_value_in_parallel"
                                     )
                                   ))
```

For more information on using the NodeJS SDK, please refer to its documentation:
***https://github.com/iopipe/iopipe/blob/master/docs/nodejs.md***

### Go SDK:

Bundled with the [Turtle CLI](https://github.com/iopipe/turtle-golang) is
a Go SDK, still in early development.

---------
Security
---------

Applications are executed in individual virtual machines
whenever allowed by the executing environment.
The definition of a virtual machine here is lax,
such that it may describe a Javascript VM,
a Linux container, or a hardware-assisted x86
virtual machine. Users should exercise caution
when running community contributed code.

It is a project priority to make fetching, publishing,
and execution of functions secure for a
production-ready 1.0.0 release.

Modules are fetched and stored using sha256 hashes,
providing an advantage over module-hosting mechanisms
which are based simply on a name and version. Future
versions of Turtle will likely implement TUF for
state-of-the-art software assurance.

Contact security@iopipe.com for questions.

-------
LICENSE
-------

Apache 2.0. Copyright 2016. IOpipe, Inc.
