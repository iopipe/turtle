var fs = require('fs')
var ctxutils = require('../../ctxutils')

module.exports = LocalDriver

function LocalDriver(opts) {
}

function get_filter_cache(id) {
  return path.join(".iopipe/filter_cache", id)
}

LocalDriver.prototype.invoke = function (event, context) {
  var id = event.id

  // Pull from index (or use cached pipescripts)
  /* download script */
  var script = fs.readFileSync(get_filter_cache(id))
  var input = ""

  return function(prevResult) {
    var sandbox = { "module": { "exports": function () {} }
                    ,"msg": prevResult
                    ,"context": context}
    var ctx = vm.createContext(sandbox)
    vm.runInContext(script, ctx)
    var result = vm.runInContext("module.exports(msg, context)", ctx)

    return context.done(result)
  }
}

LocalDriver.prototype.listFunctions = function(event, context) {
  fs.readdir(get_filter_cache("", cxtutils.callback(context)))
}
