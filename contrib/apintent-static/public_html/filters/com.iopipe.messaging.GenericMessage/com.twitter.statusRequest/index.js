var obj = JSON.parse(input);
var statusRequest = {
  "status": obj["properties"]["text"]
};
JSON.stringify(statusRequest);
