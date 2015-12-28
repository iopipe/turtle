---------------------------------------
IOPIPE
---------------------------------------
Apache 2.0 licensed.

The IOPIPE cli tool uses pipescripts, javascript-based microservices,
to simplify the consumption and integration of web services.

Pipescripts can transform data into web service requests and
can transform responses into JSON documents.

Web services may be placed anywhere within a pipeline, allowing
transformations (and code execution) to happen on both client & server.

---------------------------------------
Usage
---------------------------------------

```sh
# Fetch a web service response (Curl-like usage)
$ iopipe fetch http://localhost/some-request

# Fetch response and convert to an object (to stdout)
$ iopipe --debug exec http://localhost/some-request com.example.SomeObject

# Fetch response and convert to an object, sending SomeObject to otherhost.
$ iopipe --debug exec http://localhost/some-request com.example.SomeObject \
                      http://otherhost/request

# Fetch response and convert to an object, sending SomeObject to otherhost,
# & converting the response into a ResponseObject
$ iopipe --debug exec http://localhost/some-request com.example.SomeObject \
                      http://otherhost/request some.example.ResponseObject
```

NodeJS SDK:

```javascript
var iopipe = require("iopipe")

// Note that pipescript objects such as SomeObject are not *yet* supported.
iopipe.exec("http://localhost/some-request", "com.example.SomeObject",
            "http://otherhost.request")

// Users may chain functions and HTTP requests.
iopipe.exec(function() { return "something" }, function(arg) { return arg },
            "http://otherhost.request")

// A function may also be returned then executed later.
var f = iopipe.define("http://fetch", "https://post")
f()
```

---------------------------------------
Filters & Pipescript
---------------------------------------

Requests and responses and translated using filters written in
Pipescript (i.e. Javascript) or offered as web services.

All filters simply receive request or response data and output
translated request or response data. Pipescript is typically operated
upon locally in the client, whereas web-service based filters operate
server-side. Pipescript may also be used to build serverside filters
and applications.

For writing filters see:
***https://github.com/ewindisch/iopipe/blob/master/docs/pipescript.md***

---------------------------------------
Build & Install from source
---------------------------------------

With a functioning golang 1.5 development environment:

```bash
$ go build
$ ./iopipe --help
```

Alternatively use Docker to build & deploy:

```bash
$ docker build -t iopipe-dev .
$ eval $(echo "alias iopipe='docker run --rm -it iopipe-dev'" | tee -a ~/.bashrc)
$ iopipe --help
```

---------------------------------------
Project goals
---------------------------------------

The principal goal of our project is to improve
human to machine and machine to machine communications
with a free, highly-distributed protocol.

Furthermore:

1. Support simple translation of existing APIs
2. Support native, greenfield applications
3. Design for an open and distributed web.
4. Permissive open source licensing.
5. Security. Security. Security.

---------------------------------------
Scaling
---------------------------------------

Scaling of individual APIs is performed as is
typically done for web services. Globally, scaling is
performed through the distributed nature of the web,
as in there is no centralized server hosting all
IOPIPE web services.

Filters are pulled from (de)centralized repositories
and scale should be considered when building and
deploying filter repositories.

---------------------------------------
Security
---------------------------------------

All API endpoints MUST be exposed via TLS 1.2.
This requirement is based on flaws in earlier versions
of TLS and the fact we have the opportunity to specify
a greenfield requirement.

Note that this tool communicates to 3rd-party
web services. Caution is advised when trusting
these services, as is standard practice with
web and cloud services.

Contact security@iopipe.com for questions.

---------------------------------------
LICENSE
---------------------------------------

Apache 2.0
