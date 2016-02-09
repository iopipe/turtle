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
var url = require('url')
var request = require("request")
var util = require('util')
var vm = require('vm')
var fs = require('fs')
var path = require('path')

var USERAGENT = "iopipe/0.0.5"

function funcCallback(call, done) {
  return function() {
    done(call.apply(this, [].slice.call(arguments)))
  }
}

function httpCallback(u, done) {
  return function() {
    if (arguments.length === 0) {
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
    } else {
      prevResult = arguments[0]
      request.post({url: url.format(u), body: prevResult, strictSSL: true,
                    headers: {
                      "User-Agent": "iopipe/0.0.3"
                    }
                   },
                    function(error, response, body) {
                      if (error || response.statusCode != 200) {
                        throw "HTTP response != 200"
                      }
                      done(body)
                    })
    }
  }
}

function pipescriptCallback(id, done) {
  // Pull from index (or use cached pipescripts)
  /* download script */
  var script = fs.readFileSync(path.join(".iopipe/filter_cache/", id))
  var input = ""

  return function() {
    var prevResult = ""
    if (arguments.length > 0) {
      prevResult = arguments[0]
    }
    var sandbox = { "module": { "exports": function () {} }, "msg": prevResult }
    var ctx = vm.createContext(sandbox)
    vm.runInContext(script, ctx)
    var result = vm.runInContext("module.exports(msg)", ctx)

    return done(result)
  }
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
exports.define = function() {
  var callbackList = []
  var nextCallback;
  var done = function(result) { };

  for (var i = arguments.length - 1; i > -1; i--) {
    var arg = arguments[i];

    if (typeof arg === "function") {
      nextCallback = funcCallback(arg, done)
    } else if (typeof(arg) === "string") {
      var u = url.parse(arg);

      if (u.protocol === 'http:' || u.protocol === 'https:') {
        var server = u.hostname
        nextCallback = httpCallback(u, done)
      } else {
        nextCallback = pipescriptCallback(arg, done)
      }
    } else {
      throw new Error("ERROR: unknown argument: " + arg)
    }
    done = nextCallback
  }
  return nextCallback
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
exports.exec = function() {
  var l = [].slice.call(arguments)
  return exports.define.apply(this, l)()
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
exports.property = function (index) {
  return function (obj) {
    return obj[index]
  }
}

exports.bind = function (method, arg) {
  return function (obj) {
    return obj[method].apply(obj, [].slice.call(arguments).slice(1))
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
exports.apply = function () {
  var l = [].slice.call(arguments)
  return function (input) {
    return input.apply(input, l)
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
exports.map = function (fun) {
  return function(input) {
    var result = []
    for (i in input) {
      result.push(
        _go(function() { return fun(input[i]) })
      )
    }
    return result
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
exports.tee = function() {
  var l = [].slice.call(arguments)
  return function() {
    var args = [].slice.call(arguments)
    var result = []
    for (f in l) {
      result.push(
        _go(function() { return l[f].apply(l[f], args)})
      )
    }
    return result
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
exports.reduce = function(fun) {
  return function(input) {
    return input.reduce(fun)
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
exports.fetch = function(u) {
  request.get({url: url.format(u), strictSSL: true,
               headers: {
                 "User-Agent": USERAGENT
               }
              }, function(error, response, body) {
    if (error || response.statusCode != 200) {
      throw "HTTP response != 200"
    }
    return body
  })
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

function _go(fun) {
  return fun()
}
