# Data Processor

## Introduction

This Go microservice connects to a Mongo database with information about people, parses their names to find their middle names, and saved the result in a field in the collection.

## Requirements

- Go
- Docker

## Config

An environment variable for `MONGO_URI` can be set up to give access to a specific Mongo instance. Otherwise, a default sample instance will be used. If using a specific instance managed via Docker, you can ammend the `docker-compose.yml` if required to run using this script.
An environment variable for `PORT` can be set up to run the API on a specific port. Otherwise, a default port will be used.
There is an artifical sleep implemented between calling batches from the database. This has been used for testing purposes and can be configured/removed by altering the `sleepTime` variable in `handler.go`.

## How to Run

First, bring up the docker MongoDB container using the following command:

```
make database-up
```

To run the microservice, run the following command:

```
make run
```

If there is no data in the `persons` collection, the init function will populate the collection with some sample data.

To delete the `persons` collection, run the following command:

```
make database-drop
```

To stop the docker MongoDB container, run the following command:

```
make database-stop
```

## How to Use

Upon receiving a request on the `start` endpoint the app:

- checks if the process was already started, and return a 429 status, if that's the case
- if it wasn't started, mark it as started, then return 202 for REST
- each instance must then process batches of items in the DB, avoiding duplicate processing, until all the items have been processed

Upon receiving a request on the `stats` endpoint the app:

- checks if the process is in progress or finished and return the number of already processed items (with a 200 status for REST)
- if the process is not started, return 412
