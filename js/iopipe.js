/*
 Usage:
 var iopipe = require('iopipe')
 iopipe.exec("http://api.twitter.com/blah/blah", function() {}, "sha256:DEADBEEF", "user/pipeline", "http://somedestination/")
*/
var url = require('url')
var request = require("request")
var util = require('util')
var vm = require('vm')

function funcCallback(call, done) {
  return function() {
    done(call.apply(this, [].slice.call(arguments)))
  }
}

function httpCallback(u, done) {
  return function() {
    if (arguments.length === 0) {
      request.get({url: url.format(u), strictSSL: true }, function(error, response, body) {
        done(body)
      })
    } else {
      prevResult = arguments[0]
      request.post({url: url.format(u), body: prevResult, strictSSL: true },
                    function(error, response, body) {
                      done(body)
                    })
    }
  }
}
/*  var lastCallback = arguments[0]
  if (arguments.length > 1) {
    lastCallback(arg(arguments[1]))
  } else {
    lastCallback(arg());
  }
}*/

/* I'm thinking that perhaps we should simply convert the args into
 a Node Stream .pipe(a).pipe(b).pipe(c) etc. This will be compatible with
 Node'isms and flexible for Node users */
exports.define = function() {
  var callbackList = []
  var nextCallback;
  var lastCallback = function(result) { console.log(result) };

  for (var i = arguments.length - 1; i > -1; i--) {
    var arg = arguments[i];
    console.log("Processing arg: " + arg)

    if (typeof arg === "function") {
      nextCallback = funcCallback(arg, lastCallback)
      console.log("Processed function: " + arg)
    } else if (typeof(arg) === "string") {
      var u = url.parse(arg);

      if (u.protocol === 'http:' || u.protocol === 'https:') {
        var server = u.hostname
        nextCallback = httpCallback(u, lastCallback)
      } else {
         // Pull from index (or use cached pipescripts)
         /* download script */
         var script = ""
         if (i === 0) {
           nextCallback = function() {
             return lastCallback(vm.runInNewContext(script))
           };
         } else {
           nextCallback = function(prevResult) {
             return lastCallback(vm.runInNewContext(script, {
               input: prevResult
             }))
           }
         }
         console.log("Processed pipescript")
      }
    } else {
      console.log("Unknown argument: " + arg)
    }
    lastCallback = nextCallback
  }
  return nextCallback
}

exports.exec = function() {
  var l = [].slice.call(arguments)
  exports.define.apply(this, l)()
}
