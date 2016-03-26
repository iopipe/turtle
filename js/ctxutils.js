exports.callback = function (context) {
  return function(err, data) {
    if (err) {
      context.fail(err)
    } else { 
      context.succeed(data)
    }
  }
}
