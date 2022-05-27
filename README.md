# Transaction handler
A web service that handles transactions


## Requirements
* Docker engine
* GNU make

## Configuration

Configuration is done via environment variables.

|Name|Type|Example|Description|
|-|-|-|-|
|QONTO_APP_LISTEN_ADDRESS|string|127.0.0.1:8080|Address that application will listen on|
|QONTO_DB_NAME|string|qonto|Database name to use|
|QONTO_DB_USER|string|root|User to access database|
|QONTO_DB_PASSWORD|string|root|Password to access database|
|QONTO_DB_ADDRESS|string|127.0.0.1:13306, server.example.com|Address of remote database server with or without port information|



## Running project locally

To manipulate local environment, `make` command is used.

You can run `make help` to see all available targets and short summary about each.
```sh
$ make help
clean - Remove project containers
down - Stop all project containers
logs - Follow logs from all containers in the project
mysql-enter - Run mysql client inside database container
run - Start project in background
run-database - Start only database container
test - Run short, non-integrational, tests
test-all - Run all available tests, including integration
test-integration - Run integration tests
```

If you prefer running/debugging code from host, you have to run at least database container using:
```sh
$ make run-database
```

It will run MySQL bound to `127.0.0.1:13306` by default (see [docker-compose.dev.yml](./docker/docker-compose.dev.yml))

To clean the project up, run:
```sh
$ make clean
```
it will drop MySQL container including the data, so you can start from scratch.

To run complete system, run:
```sh
$ make run
```
now the web service is available on `http://127.0.0.1:8080/v1/transfers` (see [docker-compose.dev.yml](./docker/docker-compose.dev.yml)) and you can try that out with

```sh
$ curl -x POST @sample1.json http://127.0.0.1:8080/v1/transfers
```

## Known issues and trade-offs
* Only credit operations are supported: the task states that transfer goes **from** Qonto account and amount in individual transfer is **always positive**
    at the same time, DB description, and actual data in example DB, contains also debit operations (negative amount)
* no audit features, even transaction timestamps are missing
* identification of customer is not implemented.
It can be solved in several ways:
    1. use trusted API gateway that sets proper headers/checks requests
    1. use mTLS with customer ID baked-in
* solution is not concurrent, so DB locks, sporadic errors may occur.
To solve it, we could:
    1. use a dispatcher in front of the service, so only one request is handled at the moment
    1. make the whole system asynchronously and return `future` object that can be polled later and checked for result

## Improvements to be done (business)
* if time frames for bulk operations are known in advance,
time-based scaling of infrastructure should be done to reduce ramp time
* route organizations with huge number of transfers to own set of servers to eliminate possibility of `noisy neighbor` situation

## Improvements to be done (technical)
* add linter
* storage layer should be refactored as pure interface layer, so business logic knows nothing about Querier, transactions (SQL-specific), etc.
* error subsystem is not optimal and messy
    * error responses on API level do not have solid structure
* introduce `build` target in Makefile

## Time spending

The time I spent could be 2-3h less, but I just can't send you incomplete solution or one that I don't like.

* 1h30m initial scaffolding of the project: 
    * whiteboard design
    * local developer setup (it was taken from [another project](https://github.com/maxim-nazarenko/fiskil-lms/) that was done in ~10h)
* 2h improving storage and write integration tests
* 1h30m implement outer layer of API handlers
* 1h refactoring local setup, adding more documentation
