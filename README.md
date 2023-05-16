# P.I.R.S.

> Process Integrated Register Systems

P.I.R.S. is a cloud-first business process (BPM) register. It can be deployed for internal use or for public access.
PIRS consists of several microservices to manage register of business processes. PIRS deployments can be connected
together to make large network of accessible business processes.

---

## requiments

* go 1.20
* docker (optional)
* docker-compose (optional)

## code organization

### docs

technical and business documentation

### k8s

k8s deployment scripts

### pkg

all go microservices

### gateway

java gateway microservice

## running solution

this solution is intended to run in cloud and consist of more than serval microservices so getting it up and running can
be a bit tricky but
docker-compose_dev.yml can be used to run services that are defined there. also some changes in respective microservices
config files may be needed.
also, no guarantee that this compose file will be always upo to date.

for development purposes of individual microservices usage of mocks and high unit test coverage are highly encouraged so
each microservice can
be developed truly independently of whole solution.

## releasing

```docker-compose_release.yml``` is used for creating new image releases, when creating release change version number of
release tag in current versions section
and then run ```docker-compose -f docker-compose_release.yml build``` and then push images to registry.

releasing is done manually in this development stage and will be automated with RC CI jobs




