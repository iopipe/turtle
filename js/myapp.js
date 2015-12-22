var iopipe = require("./iopipe")

iopipe.exec("http://192.241.174.50/twitter_status/", function(arg) { console.log("Function:" + arg) })

console.log("End")
