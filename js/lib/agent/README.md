IOpipe Analytics & Distributed Tracing Agent
--------------------------------------------
[![Gitter](https://img.shields.io/gitter/room/nwjs/nw.js.svg?maxAge=2592000)](https://gitter.im/iopipe/iopipe)

This package provides analytics and distributed tracing for
event-driven, "serverless" applications.

# Installation & usage

Installation is simple. Just require this module with your client id
([contact us to get set up](mailto:support@iopipe.com)) and it will
automatically monitor and collect metrics from your application
running on AWS Lambda, Google Cloud Functions, Azure, or any
other "serverless" environment.

Example:

```javascript
var iopipe_metrics = require("iopipe-agent")({ clientId: "YOUR_ID"})

exports.handle = iopipe_metrics(
  function (event, context) {
    context.succeed("This is my serverless function!")
  }
)
```

# Configuration

When requiring the metrics agent, it accepts a config object where
you define the URL for the collector service and a client id. By default,
telemetry will be reported to the [IOpipe Telemetry Service](https://www.iopipe.com/).

```javascript
var config = {
    url: "https://metrics-api.iopipe.com", // This is the default value
    clientId: "YOUR_ID"
}
var iopipe_metrics = require("iopipe-agent")(config)
```

To _override_ the collector service, specify a URL (ex. "127.0.0.1")
running the [IOpipe Collector](https://github.com/iopipe/iopipe-collector)
as follows:

```javascript
var iopipe_metrics = require("iopipe-agent")({ url: "https://127.0.0.1", clientId: "YOUR_ID" })

exports.handle = iopipe_metrics(
  function (event, context) {
    context.succeed("Reporting these metrics to my own collector!")
  }
)
```

# Data reported

The following is provided to the collector service,
either the IOpipe Telemetry Service or the [open serverless
collector](https://github.com/iopipe/iopipe-collector).

 - function_id (a hashed identifier of the function)
 - client_id (your client id, common among all your functions)
 - environment  (restricted view of Node's process var)
 - errors
 - events       (custom events sent via `.emit('event', [data])`
 - time_sec_nanosec  (execution time: [secs, nanosecs])
 - time_sec          (execution time: secs)
 - time_nanosec      (execution time: nanosecs)

# License

Apache 2.0
