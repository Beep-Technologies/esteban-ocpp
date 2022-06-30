# Esteban-OCPP 🕴️⚡


Esteban-OCPP was developed at Beep Technologies in Singapore. The aim of Esteban is to support the deployment and popularity of electric mobility, so it is easy to install and to use. Estebann provides basic functions for the administration of charge points, user data, simple REST APIs for extensibility and was tested successfully in operation.

Esteban-OCPP considered as an open platform to implement, test and evaluate novel ideas for electric mobility, like authentication protocols, reservation mechanisms for charge points (future), and business models for electric mobility. Esteban is distributed under [GPL](LICENSE.md) and is free to use. If you are going to deploy Esteban we are happy to see the Esteban Logo on a Charge Point.

Esteban-OCPP was developed from scratch in Golang, and was highly inspired by [SteVe](https://github.com/RWTH-i5-IDSG/steve) built by RWTH Aachen University. 

Our team created Esteban-OCPP as we needed a highly-performannt and robust OCPP implementation that is highly portable, and can run either on the cloud, or on a local device running on the edge at a charger site deployment. 

Esteban-OCPP works out of the box with VoltNow & VoltPOS solutions from Beep, which enable transient customers and users to make EV Charger payments without the need for downloading any 3rd party applications.


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

## Future Works & Roadmap

Our team plans to add in support for OCPP2.0 as the protocol matures, and implement the Open InterCharge Protocol (OICP), to enable CPOs and eMSPs to easily achieve eRoaming through Esteban-OCPP.

## GDPR

If you are in the EU and offer vehicle charging to other people using SteVe, keep in mind that you have to comply to the General Data Protection Regulation (GDPR) as Esteban processes charging transactions, which can be considered personal data.

## Enterprise Support

As an Esteban Enterprise Editionn customer, you will have professional support included. Get reliable, high-touch support from senior support engineers.

#### How to receive support

Please write us an email to hello@beepbeep.tech to get help with your Esteban Enterprise Editionn customer. Include a detailed description as well as screen-shots, where necessary.

#### Support ticket (only for Premier Support and Corporate Support)

You can create a support ticket on our community platform. Please contact us for details. You first need to register on the OpenProject Community edition to create a support ticket.

#### Contact us

If you have more questions, please contact us at hello@beepbeep.tech






