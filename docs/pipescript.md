Pipescript is a subset of ECMAScript / Javascript.

Scripts initiate with a variable predefined, 'input',
containing input bytes to be parsed. Each script is associated
with a pre-defined input and output type.

Scripts are expected to output an Object matching the following
Object Schema. Objects should adhere to their Class Schema, which
is their JSON Schema definition.

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

-------------
Object Schema
-------------

```json
{
	"classid": "string",
	"classvers": "string",
	"properties": {}
}
```

JSON Schema references SHOULD be supported as an experimental
extension to the APIntents schema, allowing properties to point
to local and remote resources. Clients MAY refuse multi-host
traversal, either entirely, or limit the number of redirections.
Support for JSON Schema  references MUST be a compile or
run-time option while this feature remains experimental.

-------------
Class Schema
-------------

Objects implement a class. Each class is defined
by JSON Hyper-Schema which describes the properties of an Object.

```json
{
  "title": "ArticleText",
  "type": "object",
  "properties": {
    "title": {
      "title": "Article Title",
      "type": "string"
    },
    "text": {
      "title": "Article Text",
      "type": "string"
    },
  },
  "required" : ["title", "text"],
  "links": [
    {
      "rel": "author",
      "href": "/Users/{authorId}"
    }
  ]
}

```

