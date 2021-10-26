# beepbeep3-ocpp / bb3-ocpp / Esteban 🕴️

The BB3-OCPP microservice (a.k.a. Esteban) is intended to manage all the Charge Points (CPs) that are controlled by Beep. It interfaces with these Charge Points via the Open Charge Point Protocol (OCPP), which is essentially a RPC framework implemented over WebSockets.

## Setup

A working installation of docker and docker-compose is required.
- Copy and paste `.env.example` into a new `.env` file.

- Run the following

```
docker-compose up -d
go run cmd/bb3-ocpp-ws/main.go
```

This will set up a dockerized instance of a PostgreSQL (timescaledb) server and a web-based postgres ui (pgweb) at `localhost:8062`.

Swagger Docs can then be viewed at `localhost:8060/v2/ocpp/swagger/index.html`
