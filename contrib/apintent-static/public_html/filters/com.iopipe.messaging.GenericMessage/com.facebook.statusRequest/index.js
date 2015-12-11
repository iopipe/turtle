var obj = JSON.parse(input);
var statusRequest = {
  "message": obj["properties"]["text"]
};
JSON.stringify(statusRequest);
