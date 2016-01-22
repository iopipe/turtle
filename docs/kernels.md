Kernels must currently be written in a subset of ECMAScript / Javascript.
Future versions of IOpipe may allow kernels to be developed in
other languages, or to run binary kernels.

A CommonJS module format is employed, expecting a function defined
as 'module.exports'. This function should return a result intended
for the next kernel, function, or HTTP endpoint.

Currently, a single argument and single string-type return variable
are supported.

------------------
Example kernel
------------------

The following converts a JSON document representing a "GenericMessage"
into a Twitter status update request (as expected by the Twitter API).

```javascript
module.exports = function(input) {
  var obj = JSON.parse(input)
  var statusRequest = {
    "status": obj["properties"]["text"]
  }
  return JSON.stringify(statusRequest)
}
```
