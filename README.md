# Sensor Metadata

This project contains a sensor metadata microservice. It exposes a json REST API for storing and querying sensor metadata. 
With this microservice you can:
* Store name, location (gps position), and a list of tags for each sensor.
* Retrieve metadata for an individual sensor by name or by id.
* Update a sensor’s metadata.
* Query to find the sensor nearest to a given location.

## Tech Stack

* [Go](https://golang.org/dl/) (1.19 or higher)
* [Mongo](https://www.mongodb.com/)
* [Docker](https://www.docker.com/)
* [Kubernetes](https://kubernetes.io/)
* [Helm](https://helm.sh/)

## Structure

```
api/               # Swagger files of the exposed APIs
cmd/               # Implementation of the microservices
    sensor/        # The root of the sensor meta-data microservice
        db/        # Database access layer
        handlers/  # End-point implementation using http
        service/   # Business logic layer 
        Dockerfile # Dockerfile for this microservice
        main.go    # Server entrypoint
    authenticator/ # Auth microservice (NOT IMPLEMENTED YET).
deployments/       # Helm configuration files
scripts/           # Auxiliary shell scripts
test/              # Integration tests
    client/        # Generated client from swagger
```
The [go project layout](https://github.com/golang-standards/project-layout) standard was followed to a degree, 
with some modifications to allow multiple microservices in the same repository.

## Enviroment variables and run arguments

All environment variables and run arguments are already configure between the Dockerfile and Helm / Kubernetes scripts.
Nevertheless, here they are:

TODO

## Building, testing and running

`make help` will show the commands for building, running lint, tests and deployment.

In order to deploy you need to have `kubectl` properly configured to the desired k8s cluster and helm installed.
If you want to run locally, you can use `scripts/kind.sh` to install kind and create a cluster with ingress enabled.
After that `make deploy` will deploy to your local cluster.

If you want to run separately you need to configure the enviroment variables and run the following commands:

TODO

You test the basic operations with the following commands:

``` 
curl --request POST http://localhost/sensor-metadata \
--data-raw '{ "name" : "Sensor 1", "tags" : [ "Tag1", "Tag2" ] , "location" : { "lat" : "35", "lon" : "45"  } } '
```
which should return something like:
`{"id":"63bcf00cf3ed6129b61c137b"}`

You can then try to find by id using:
`curl http://localhost/sensor-metadata/63bcf00cf3ed6129b61c137b`

Find by name using:
`curl http://localhost/sensor-metadata/by-name/Sensor%201`
(Note that spaces in the name should be replaced by %20)

Or find the nearest sensor using:
`curl http://localhost/sensor-metadata/nearest/35/45`

You may also update the sensor meta-data or delete it. Please check under `api/swagger.yml` for more information.




