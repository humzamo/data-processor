# Data Processor

## Introduction

This Go microservice connects to a Mongo database with information about people, parses their names to find their middle names, and saves the result in a field in the collection.

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

Upon receiving a request on the `pause` endpoint the app:

- checks if the process is finished and if so, returns a 200 status
- if the process is not started, return 412
- if the process is not finished, pauses the process and returns a 200 status

## Assumptions/Further Improvements

- The authentication for the sample DB in the container has been kept very simple for testing purposes; this could (and should!) be made more secure (e.g. by storing the password in an environment variable or secret)
- Some assumptions have been made for populating data (i.e. if there is no data in the collection, add the sample data)
- The processing function logic has been kept simple and only parses a name for middle names. This could be made more detailed/complicated with additional requirements.
- The sleeping logic has been used for testing and to simulate what could happen there the database is much larger and processing needs to be more segmented,
- For simplicity, the processing status has been restricted to a variable within the API instance, but this could be expanded to allow for syncronisation across multiple instances (for example by having an additional endpoint to manage, control, and lock the status across instances). At present, each instance would separately manage the processing status (e.g. if processing was started for one instance, it wouldn't start for all instances). However, since each instance connects to the same MongoDB, and the processed flag is updated for each record, there would not be any duplicate processing and the instances would still align and ultimately finish together.
