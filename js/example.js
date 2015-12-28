var iopipe = require("./iopipe")

var d = iopipe.define("http://192.241.174.50/twitter_status/", function(arg) { return arg }, "https://api.twitter.com/1.1/statuses/update.json")
//var d = iopipe.define(function() { return 1 }, function(one) { return one + 1 }, function(arg) { console.log("myapp-function:" + arg); return arg })

d()
