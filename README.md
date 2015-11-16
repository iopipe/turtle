---------------------------------------
IOPIPE
---------------------------------------
Apache 2.0 licensed.

The IOPIPE cli tool uses the open PipeAPI to
provide compatibility and discovery between web services.

All web sevices and IoT devices are expected to provide
HTTP services with an API. To simplify automatic integration
with these services, it's necessary to have a lingua de franca,
an abstract language for which to communicate. For this reason,
PipeAPI is designed to provide a simplified and unified interface
to arbitrary data, sensors, and machine operations.

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
API Brokers (porting legacy APIs)
---------------------------------------

A reference broker will be published, TBA.

---------------------------------------
Scaling
---------------------------------------

Scaling of individual APIs is performed as is
typically done for web services. Globally, scaling is
performed through the distributed nature of the web,
as in there is no centralized server hosting all
IOPIPE web services.

---------------------------------------
Security
---------------------------------------

All API endpoints MUST be exposed via TLS 1.2.
This requirement is based on flaws in earlier versions
of TLS and the fact we have the opportunity to specify
a greenfield requirement.

Note that API Brokers may communicate to backend web
services through less secure channels. Not all API
brokers and API endpoints are operated by IO PIPE,
and caution is advised (to the same degree caution
may be warranted when accessing a mobile app or
web service).

Contact security@iopipe.com for questions.

---------------------------------------
LICENSE
---------------------------------------

Apache 2.0
