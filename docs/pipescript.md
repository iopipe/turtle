Pipescript is a subset of ECMAScript / Javascript.

Scripts are written in CommonJS module format in that they expect
a function defined as 'module.exports'. This function should return
a result intended for the next pipescript, function, or HTTP endpoint.

Currently, a single argument and single string-type return variable
are supported.

------------------
Example pipescript
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
