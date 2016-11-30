/**
  @module iopipe
  @description
  IOpipe helps developers build, connect, and scale code.

  Using flows, iopipe simplifies the consumption and integration
  of web services through the chaining of kernels,
  single-function applications.

  Kernels take and transform input, providing straight-forward output
  in a fashion to Unix pipes. A kernel may receive input or send output
  to/from web service requests, functions, or local applications.

  IOpipe may be embedded in applications, used from shell scripts,
  or run manually via a CLI to form complete applications. Kernels
  and pipelines may be run within local processes, or dispatched to
  remote workers (i.e. "cloud").

  Basic example usage:

  ```javascript
  var iopipe = require('iopipe')
  iopipe.exec("http://api.twitter.com/blah/blah"
              ,function() {}
              ,"sha256:DEADBEEF"
              ,"user/pipeline"
              ,"http://somedestination/")
  ```
*/
var events = require('events')
var url = require('url')
var request = require("request")
var vm = require('vm')
var path = require('path')

var local_driver = require('./exec_drivers/local')
var iopipe_agent = require('./lib/agent')

var USERAGENT = "iopipe/0.0.8"

/**
   @description
   Initalizes an IOpipe environment,
   accepting an argument specifying runtime options
   such as the execution driver ('local', 'aws', 'gcp', etc),
   and settings for those execution drivers.

   Without arguments, performs local execution and reads/writes
   kernels to the directory .iopipe_cache.

   ```javascript
   IOpipe({
     exec_driver: 'aws'
     exec_driver_opts: {
       region: 'us-west-1',
       access_key: 'itsasecrettoeverybody',
       secret_key: 'itsasecrettoeverybody'
     },
     metrics_opts: {
       client_id: "your_client_id"
     }
   })
   ```

   @param object options - Runtime options.
*/
function IOpipe(options) {
  var metrics_opts = {}
  this._metrics_agent = iopipe_agent(metrics_opts)

  var _exec_driver = 'local'
  if (options && "exec_driver" in options) {
    _exec_driver = options.exec_driver
  }
  driver_options = {}
  if (options &&
      "exec_driver_opts" in options) {
    driver_options = options.exec_driver_opts
  }
  this._exec_driver = require("./" + path.join('./exec_drivers/', _exec_driver, 'index.js'))(
    driver_options
  )
}

module.exports = function(options) {
  return new IOpipe(options)
}

function funcCallback(call, context) {
  return function() {
    var args = [].slice.call(arguments)
    if (args.length == 0) {
      args.push(undefined)
    }
    args.push(context)
    call.apply(this, args)
  }
}

function httpCallback(u, context) {
  return function() {
    if (arguments.length === 0) {
      request.get({url: url.format(u), strictSSL: true,
                   headers: {
                     "User-Agent": USERAGENT
                   }
                  }, function(error, response, body) {
        if (error || response.statusCode != 200) {
          context.fail("HTTP response != 200")
        }
        context.done(body)
      })
    } else {
      prevResult = arguments[0]
      request.post({url: url.format(u), body: prevResult, strictSSL: true,
                    headers: {
                      "User-Agent": USERAGENT
                    }
                   },
                    function(error, response, body) {
                      if (error || response.statusCode != 200) {
                        context.fail("HTTP response != 200")
                      }
                      context.done(body)
                    })
    }
  }
}

IOpipe.prototype.make_context = function(done) {
  var ctx = done
  ctx.done = done
  ctx.fail = function(failure) {
    throw failure
  }
  ctx.succeed = function(result) {
    var args = [].slice.call(arguments)
    return done.apply(this, args)
  }
  ctx.raw = done
  return ctx
}

/**
   @description
   Defines a pipeline, returning a function.
   Used for passing arguments to a  pipeline as this
   is not possible with 'exec', or for  reusing a
   pipeline. Users seeking a  method with callback
   should use exec (which actually wraps define),
   or call:

   ```javascript
   define(args...)(input)
   ```

   @param {...(string|function)} kernel - Kernels specified as functions, scripts, or HTTP endpoints.
*/
IOpipe.prototype.define = function() {
  /* We return a function that executes the pipeline,
     if arguments are supplied, the first is input, and the remainder
     are callbacks. */
  var iopipe = this
  var defargs = [].slice.call(arguments)
  var aws_context = null

  return function() {
    var done = function(e) { return e };

    /* support callback to function returned by define()
     * i.e. define(f1)(data, f2) ~= f2(f1(data)) */
    var largs = [].slice.call(arguments)
    if (largs.length > 1) {
      aws_context = largs[1]
      /* check for awsRequestId to detect Lambda */
      done = function (e) {
        if ('awsRequestId' in aws_context) {
          aws_context.succeed(e)
        } else {
          largs[1](e)
        }
      }
    }

    for (var i = defargs.length - 1; i > -1; i--) {
      var arg = defargs[i];

      var context = iopipe.make_context(done)

      if (typeof arg === "function") {
        done = funcCallback(arg, context)
      } else if (typeof(arg) === "string") {
        var u = url.parse(arg);

        if (u.protocol === 'http:' || u.protocol === 'https:') {
          var server = u.hostname
          done = httpCallback(u, context)
        } else {
          done = this._exec_driver.invoke({ id: arg }, context)
        }
      } else {
        throw new Error("ERROR: unknown argument: " + arg)
      }
    }

    // Call function with input data.
    done(largs[0])
  }
}

/**
  @description
  Executes a pipeline, a la waterfall async pattern.
  Each argument is a callback for the result of the previous
  function. The final function may be seen as being the penultimate
  callback for triggering events.

  Usage:

  ```javascript
  iopipe.exec("http://127.0.0.1"
              ,"my_pipescript"
              ,function(i) { return i }
              ,"http://127.0.0.2/post"
              ,callback)
  ```

  @param {...(string|function)} kernel - Kernels specified as functions, scripts, or HTTP endpoints.
*/
IOpipe.prototype.exec = function() {
  var l = [].slice.call(arguments)
  return this.define.apply(this, l)()
}

/**
  Returns a function to access a property/index in an input array.

  Example:

  ```javascript
  iopipe.define(iopipe.property(0))(["hello", "world"])
  //=> "hello"
  ```

  @param {*} property - Property to access in input to returned function.
*/
IOpipe.prototype.property = function (index) {
  return function (obj, done) {
    done(obj[index])
  }
}

IOpipe.prototype.bind = function (method, arg) {
  return function (obj, done) {
    done(obj[method].apply(obj, [].slice.call(arguments).slice(1)))
  }
}

/**
   Return a function that accepts a function parameter,
   currying any parameters passed to apply() itself.
   for instance, the following is a "hello world" for apply:
     apply("Hello world")(function(x) { console.log(x) }))

   This is useful with iopipe where a function returns another
   function and the developer wishes to call this with an iopipe
   pipeline:

  ```javascript
  iopipe.exec(function() { return function (x) { console.log(x) } }
              ,iopipe.apply("hello world"))
  ```

  @param {...*} arguments - Arguments to pass to input of returned function.
*/
IOpipe.prototype.apply = function () {
  var l = [].slice.call(arguments)
  return function (input, done) {
    done(input.apply(input, l))
  }
}

/**
  Returns a map function for executing pipelines for each value
  in an input array. This is how one loops over elements and performs
  transformations of multiple elements with iopipe.

  Example (adds 1 to each array value):

  ```javascript
  iopipe.map(function(i) { return i + 1 })([0, 1, 2])
  //=> [1, 2, 3]
  ```

  @param function function - Function to call against each input provided to output function.
*/
IOpipe.prototype.map = function(fun) {
  var iopipe = this
  return function(input, done) {
    var result = []
    var waiter = new events.EventEmitter()
    var eventid = 'map-callback'
    waiter.setMaxListeners(1)
    waiter.on(eventid, function(msg) {
      setImmediate(function() {
        result.push(msg)
        if (input.length === result.length) {
          done(result)
          waiter.removeAllListeners(eventid)
        }
      })
    })
    for (i in input) {
      fun(input[i], iopipe.make_context(function(msg) {
        waiter.emit(eventid, msg)
      }))
    }
  }
}

/**
  Returns a function which executes each argument
  function/pipeline against a single input. That is,
  each passed argument (function) is called with
  the given input. The effective opposite of map(),
  although equally parallelizable.

  Example:

  ```javascript
  function echo(i) {
    return i
  }
  iopipe.tee(echo, echo)("hello world")
  //=> ["hello world", "hello world"]
  ```

  @param {...function} function - Functions to call against the input to the output function.
*/
IOpipe.prototype.tee = function() {
  var iopipe = this
  var tfuncs = [].slice.call(arguments)
  return function(input, context) {
    var result = []
    var waiter = new events.EventEmitter()
    var eventid = 'tee-callback'
    waiter.setMaxListeners(1)
    waiter.on(eventid, function(msg) {
      setImmediate(function() {
        result.push(msg)
        if (result.length === tfuncs.length) {
          context(result)
          waiter.removeAllListeners(eventid)
        }
      })
    })
    for (f in tfuncs) {
      tfuncs[f].apply(tfuncs[f], [input, iopipe.make_context(function(msg) {
        waiter.emit(eventid, msg)
      })])
    }
  }
}

/**
  Returns a reduce function for consolidating results or
  "squeezing" an array into a single value output.

  Example (sum):

  ```javascript
  iopipe.reduce(function(prev, curr) { return prev + curr })([2, 2])
  //=> 4
  ```

  @param function function - Function to reduce params to returned function.
*/
IOpipe.prototype.reduce = function(fun) {
  return function(input, done) {
    done(input.reduce(fun))
  }
}

/**
  Returns a function which fetches the input URL via HTTP(s).
  This is useful if using iopipe to create a URL, as may happen
  if transforming some data or some API result into a new API request.

  Example:

  ```javascript
  var getHNitem = iopipe.define(
                    "https://hacker-news.firebaseio.com/v0/items/".concat,
                    ,iopipe.fetch)
  getHNitem(1000, function(data) {
    console.log("Got HackerNews story:")
    console.log(data)
  })
  ```
*/
IOpipe.prototype.fetch = function(u) {
  return function(input, done) {
    request.get({url: url.format(u), strictSSL: true,
                 headers: {
                   "User-Agent": USERAGENT
                 }
    }, function(error, response, body) {
      if (error || response.statusCode != 200) {
        throw "HTTP response != 200"
      }
      done(body)
    })
  }
}


/**
  Creates a new function accepting a callback around a function
  which does not accept a callback parameter as its last argument.

  @param {...function} function - Function to wrap a callback around.
*/
IOpipe.prototype.callback = function(fun) {
  return function(input, done) {
    done(fun(input))
  }
}

/**
   We monkey-patch the Object.values function,
   this makes it easier to map assoc arrays using tee.
   Some Javascript implementations already offer this
   function with the same interface. Monkey-patching this is
   ugly, but should be mostly-safe(?)

   Example:

   ```javascript
   iopipe.tee(Object.keys, Object.values)({"hello": "world"})
   //=> ["hello", "world"]
   ```

   @param array array - Associative array to return values of.
*/
if (!Object.hasOwnProperty("values")) {
  Object.values = function (arr) {
    return Object.keys(arr).map(function(y) {return arr[y]})
  }
}
