/*
 Usage:
 var iopipe = require('iopipe')
 iopipe.exec("http://api.twitter.com/blah/blah", function() {}, "sha256:DEADBEEF", "user/pipeline", "http://somedestination/")
*/
var url = require('url')
var request = require("request")
var util = require('util')
var vm = require('vm')
var fs = require('fs')

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
  var script = fs.readFileSync(".iopipe/filter_cache/" + id)
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

exports.exec = function() {
  var l = [].slice.call(arguments)
  return exports.define.apply(this, l)()
}

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

/* Return a function that accepts a function parameter,
   currying any parameters passed to apply() itself.
   for instance, the following is a "hello world" for apply:
     apply("Hello world")(function(x) { console.log(x) }))

   This is useful with iopipe where a function returns another
   function and the developer wishes to call this with an iopipe
   pipeline:

     iopipe.exec(function() { return function (x) { console.log(x) } }
                 ,iopipe.apply("hello world"))
*/
exports.apply = function () {
  var l = [].slice.call(arguments)
  return function (input) {
    return input.apply(input, l)
  }
}

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

exports.reduce = function(fun) {
  return function(input) {
    return input.reduce(fun)
  }
}

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

/* Monkey-patch Object.values function,
   this makes it easier to map assoc arrays using tee:
    > iopipe.tee(Object.keys, Object.values)({"hello": "world")
    ("hello", "world")
*/
if (!Object.hasOwnProperty("values")) {
  Object.values = function (arr) {
    return Object.keys(arr).map(function(y) {return arr[y]})
  }
}

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

function _go(fun) {
  return fun()
}
