-- add mock data
insert into bb3.ocpp_application (id, name)
values ('cgw', 'chargegowhere');
insert into bb3.ocpp_application_callback (
        application_id,
        callback_event,
        callback_url
    )
values (
        'cgw',
        'StartTransaction',
        'http://localhost:8080/connectivity/bb3-ocpp/callbacks/StartTransaction'
    ),
    (
        'cgw',
        'StopTransaction',
        'http://localhost:8080/connectivity/bb3-ocpp/callbacks/StopTransaction'
    ),
    (
        'cgw',
        'StatusNotification',
        'http://localhost:8080/connectivity/bb3-ocpp/callbacks/StatusNotification'
    ),
    (
        'cgw',
        'MeterValues',
        'http://localhost:8080/connectivity/bb3-ocpp/callbacks/MeterValues'
    );
insert into bb3.ocpp_application_entity (application_id, entity_code)
values ('cgw', 'B2070807');
insert into bb3.ocpp_charge_point (
        /* user-set */
        entity_code,
        ocpp_protocol,
        charge_point_identifier,
        connector_count
    )
values (
        'B2070807',
        'ocpp1.6J',
        'SUTD_TEST',
        0
    );