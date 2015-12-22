/*
 pipe("http://api.twitter.com/blah/blah", function() {}, "sha256:DEADBEEF", "user/pipeline", "http://somedestination/")
*/
var url = require('url')
var http = require('http')

exports.exec = function() {
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

    if (typeof arg === "function") {
      console.log("Run function")
      if (i === 0) {
        callbackList.push(arg)
	continue
      }
      callbackList.push(function() { arg(prevResult) }
    } else if (typeof(arg) === "string") {
       var u = url.parse(arg);

       if (u.protocol === 'http:' || u.protocol === 'https:') {
         var server = u.hostname
         if (i > 0) {
           console.log("POST")
           u.method = 'POST';
           var req = http.request(u, function(res) {
             res.on('data', function(body) {
               console.log("Got response: " + body)
               prevResult = body
             })
           })
           req.write(prevResult)
           req.end()
           continue
         }
         console.log("GET")
         http.get(arg, function(res) {
           res.on('data', function(body) {
             console.log("Got response: " + body)
             prevResult = body
           })
           prevResult = res.read()
           console.log("Response: " + prevResult)
         })
      } else {
         // Pull from index (or use cached pipescripts)
      }
    } else {
      console.log("Unknown argument.")
    }
  }
}
