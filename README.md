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
    authenticator/ # Auth microservice, which follows the same structure as sensor, but the implementation is not polished.
deployments/       # Helm configuration files
scripts/           # Auxiliary shell scripts
test/              # Integration tests
    client/        # Generated client from swagger
```
The [go project layout](https://github.com/golang-standards/project-layout) standard was followed to a degree, 
with some modifications to allow multiple microservices in the same repository.

## Enviroment variables and run arguments

All environment variables and run arguments are already configured between the Dockerfile and Helm / Kubernetes scripts.
Nevertheless, here are the environment variables:

ENV Variable   | Used by       | Default | Description
---------------|---------------|:-------:| --------------------------------
ROOT_PASSWORD  | Authenticator |  1234   | Password for base user root
tls.crt        | Both          |         | Jwt Certificate public key
tls.key        | Authenticator |         | Jwt Certificate private key

To check the run time arguments you may run after `make`

`./out/authenticator -h`

and

`./out/sensor -h`


## Building, testing and running

`make help` will show the commands for building, running lint, tests and deployment.

### Quick Start

```
make
./scripts/kind.sh (ONLY IF YOU DON"T HAVE A LOCAL CLUSTER)
kubectl wait --namespace ingress-nginx \
  --for=condition=ready pod \
  --selector=app.kubernetes.io/component=controller \
  --timeout=90s  
make deploy
```

In order to deploy you need to have `kubectl` properly configured to the desired k8s cluster and helm installed.
If you want to run locally, you can use `scripts/kind.sh` to install kind and create a cluster with ingress enabled.
After that `make deploy` will deploy to your local cluster.

If you want to run outside the container you need to configure the enviroment variables and run the following commands:

```
make
./out/authenticator&
./out/sensor&
```
In this case you must have a mongo running on localhost or run with arguments. The default port for authenticator is 3000 and for sensor is 4000

## Basic tests

Here there is sequence of curl commands for the base use cases. We are assuming a local kubernetes deployed using `make deploy`.
If you are running the services outside the container you should replace localhost/sensor-metadata to localhost:4000
and localhost/authenticator to localhost:3000.

First Login:
```
curl --request POST http://localhost/authenticator/login --data-raw '{ "username" : "root", "password" : "1234"  } '
```
Copy the result and add this header in all POST, PUT and DELETE operations as follows:

```
--header 'Authorization: token [PASTE_TOKEN]'
```

Registering users:

```
curl --request POST http://localhost/authenticator/register \
--header 'Authorization:  token eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InJvb3QiLCJyb2xlIjoiQURNSU4ifQ.T0G4DKnjCIpUNQo0dRUH4gkzrcDzaXoGdIp_rjvRg06tRJpO2Etw1UAyOWgJFICCnkFU-tDFvceGZiM3jVQkjrNB-VFk8KiPgdLZttRatXEtk8a5NHfN3oyP0doBXqkw8s_olHzI0iwQK4yAVmsavTenOaQU9cpgk1U1Ite303x_hGy_S97NvC0RA_PIaHjdeNpJcTz9M1bRDxuKiUm_YAJuQu14wbw8chfNEaADJ7xoHjBRPtQqMH4R3Is0jJTd66thxfIUETxPpp7LXzVup0V-4bGuAjN06clXa-Klp8K8JhlTsEnl0BGUPSTBKQe0vWOA2WNu5mWQGtOcLSQhcA' \
--data-raw '{ "username" : "user", "password" : "1234", "role" : "USER" } '
```
The only role supported is ADMIN which allows POST, PUT and DELETE operations. GET do not require authentication.


Where the token should be replaced accordingly. For simplicity we will ommit the token on the following commands.


Creating a sensor:
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
There is no swagger for authenticator as it was not the focus of this work and it only has the two endpoints listed here.




