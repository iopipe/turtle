Kernels must currently be written in a subset of ECMAScript / Javascript.
Future versions of IOpipe may allow kernels to be developed in
other languages, or to run binary kernels.

A CommonJS module format is employed, expecting a function defined
as 'module.exports'. This function should accept two parameters,
an input variable, and a "context" object providing callbacks.
The context object has the properties 'succeed', 'done', and 'fail'.

A function should pass its output as intended for the next kernel,
function, or HTTP endpoint via context.done() or context.succeed()
callbacks.

------------------
Example kernel
------------------

The following converts a JSON document representing a "GenericMessage"
into a Twitter status update request (as expected by the Twitter API).

```javascript
module.exports = function(input, context) {
  var obj = JSON.parse(input)
  var statusRequest = {
    "status": obj["properties"]["text"]
  }
  context.done(JSON.stringify(statusRequest))
}
```
