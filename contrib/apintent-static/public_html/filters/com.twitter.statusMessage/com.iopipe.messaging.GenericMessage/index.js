var obj = JSON.parse(input);
var tweet = obj["properties"];
var statusMessage = {
  "id":   "/objects/statusMessage/" + tweet["id_str"],
  "user": "/objects/user/" + tweet["user"]["id_str"],
  "text": tweet["text"]
};
JSON.stringify(statusMessage);
