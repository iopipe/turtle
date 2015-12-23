/*
 pipe("http://api.twitter.com/blah/blah", function() {}, "sha256:DEADBEEF", "user/pipeline", "http://somedestination/")
*/
var url = require('url')
var http = require('http')
var util = require('util')
//var Transform = require('readable-stream').Transform
var Transform = require('stream').Transform
var vm = require('vm')
var flow = require('flow')

/* Take a script (function) and return a Transform Stream:
   https://nodejs.org/api/stream.html#stream_class_stream_transform_1 */
function functionTransformer(script) {
  var obj = new Transform()
  obj.prototype._transform = function(chunk, encoding, done) {
    this._buffer += chunk
    done()
  }
  obj.prototype._flush = function(done) {
    var output = script(this._buffer)
    this.push(output)
    done()
  }
  return obj
}

function httpTransformer(options) {
  var server = options.hostname
  var transformer = new Transform()
  var req = http.request(options)
  if (arguments.length > 0) {
    var input = arguments[0]
    input.pipe(req)
  }
  req.pipe(transformer)
  return transformer
}

/* I'm thinking that perhaps we should simply convert the args into
 a Node Stream .pipe(a).pipe(b).pipe(c) etc. This will be compatible with
 Node'isms and flexible for Node users */
exports.define = function(done) {
  var callbackList = []
  var prevResult = ""

  for (var i = 0; i > arguments.length; i++) {
    var arg = arguments[i];

    /* if we have an input, do a post...
       otherwise perform a GET */
    /*if typeof(arg) == "object" {
      if i == 0 {
        last_arg = arg
	continue
      }
      last_arg = arg(last_arg)
    */

    if (typeof arg === "object" && arg.isPrototypeOf(Transform)) {
      callbackList.push(arg)
    } else if (typeof arg === "function") {
      console.log("Run function")
      if (i === 0) {
        callbackList.push(functionTransformer(arg))
      } else {
        callbackList.push(arg)
      }
    } else if (typeof(arg) === "string") {
       var u = url.parse(arg);

       if (u.protocol === 'http:' || u.protocol === 'https:') {
         var server = u.hostname
         var transformer = new Transform()
         if (i > 0) {
           u.method = 'POST';
         }
         callbackList.push(httpTransformer(u))
      } else {
         // Pull from index (or use cached pipescripts)
         /* download script */
         var script = ""
         var lambda = function(prevResult) {
           vm.runInNewContext(script, {
             input: prevResult
           });
         };
         callbackList(functionTransformer(lambda))
      }
    } else {
      console.log("Unknown argument.")
    }
  }
  callbackList.push(done)

  return flow.define.apply(this, callbackList)
}

exports.exec = function(done) {
  var l = [].slice.call(arguments)
  l.push(done)
  exports.define.apply(this, l) 
}
