var iopipe = require("./iopipe")

//var d = iopipe.define("http://192.241.174.50/twitter_status/", "7d6fd7774f0d87624da6dcf16d0d3d104c3191e771fbe2f39c86aed4b2bf1a0f")
//var d = iopipe.define(function() { return 1 }, "7d6fd7774f0d87624da6dcf16d0d3d104c3191e771fbe2f39c86aed4b2bf1a0f")
var d = iopipe.define("http://192.241.174.50/twitter_status/","8b4d14fd6343f8b722a804f279f389b78ac82e4ca61fdb158f2715ce2d7806ca", console.log, "https://api.twitter.com/1.1/statuses/update.json")
//var d = iopipe.define("http://192.241.174.50/twitter_status/", function(arg) { return arg }, "https://api.twitter.com/1.1/statuses/update.json")
//var d = iopipe.define(function() { return 1 }, function(one) { return one + 1 }, function(arg) { console.log("myapp-function:" + arg); return arg })

d()
