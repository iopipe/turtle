Pipescript is a subset of ECMAScript / Javascript.

Scripts initiate with a variable predefined, 'input',
containing input bytes to be parsed. Each script is associated
with a pre-defined input and output type.

------------------
Example pipescript
------------------

The following converts a GenericMessage into a Twitter status update request:

```javascript
var obj = JSON.parse(input);
var statusRequest = {
  "status": obj["properties"]["text"]
};
JSON.stringify(statusRequest);
```
