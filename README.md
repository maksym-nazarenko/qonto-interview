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
