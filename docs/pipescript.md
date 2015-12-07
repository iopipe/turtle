Pipescript is a subset of ECMAScript / Javascript.

Scripts initiate with a variable predefined, 'input',
containing a JSON string to be parsed.

Scripts are expected to return a JSON string.

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

