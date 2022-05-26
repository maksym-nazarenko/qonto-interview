# Transaction handler
A web service that handles transactions


## Requirements
* Docker engine
* GNU make

## Configuration


## Running project locally

To manipulate local environment, `make` command is being used.

You can run `make help` to see all available targets and short summary about each.
```sh
$ make help

```

If you prefer running/debugging code from host, you have to run at least database container using:
```sh
$ make run
```
It will run MySQL bound to `127.0.0.1:13306` by default (see [docker-compose.dev.yml](./docker/docker-compose.dev.yml))


## Known issues and trade-offs
* the task states that transfer goes **from** Qonto account and amount in individual transfer is **always positive**
    at the same time, DB descriptin, and actual data in example DB, contains also debit operations (negative amount)
* no audit features, even transaction timestamps are missing
* storage layer should be refactored as pure interface layer, so business logic knows nothing about Querier, transactions (SQL-specific), etc.
* error subsystem is not optimal and messy
    * error responses on API level do not have solid structure
* add linter
* introduce `build` target in Makefile

## Time spending

* 1h30m initial scaffolding of the project: 
    * whiteboard design
    * local developer setup (it was taken from [another project](https://github.com/maxim-nazarenko/fiskil-lms/) that was done in ~10h)
* 2h improving storage and write integration tests
