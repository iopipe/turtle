---------------------------------------
IOPIPE
---------------------------------------
Apache 2.0 licensed.

The IOPIPE cli tool uses APIntents to provide compatibility
and discovery between web services.

All web sevices and IoT devices are expected to provide
HTTP services with an API. To simplify automatic integration
with these services, it's necessary to have a lingua de franca,
an abstract language for which to communicate. For this reason,
APIntents are designed to provide a simplified and unified interface
to arbitrary data, sensors, and machine operations.


---------------------------------------
Usage
---------------------------------------

```sh
# Fetch a web service response (Curl-like usage)
$ iopipe fetch http://localhost/some-request

# Fetch response and convert to an object (to stdout)
$ iopipe convert http://localhost/some-request com.example.SomeObject

# Fetch response and pipe into another web service
$ iopipe copy http://localhost/dogs/spot http://otherhost/dogs/
```

The above commands may be temporary as we create functioning
piping mechanisms and object conversion tooling. It's intended,
for example, that 'iopipe copy' will be able to transform a
response not only into a request, but will be able to transform them.

For instance, copy should be able to turn a "Dog" object into a
generic "Animal" one, for a receiving service which does not
know about Dogs.

Pipes from and into the command-line's STDIN/STDOUT are
also intended, but are not yet integrated.

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
