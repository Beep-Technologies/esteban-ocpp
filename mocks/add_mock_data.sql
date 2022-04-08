-- add mock data
insert into bb3.ocpp_application (id, name)
values ('cgw', 'chargegowhere'), ('busways', 'busways');
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
    ),
    (
        'busways',
        'StartTransaction',
        'http://localhost:8080/callbacks/bb3-ocpp-callback/start-transaction'
    ),
    (
        'busways',
        'StopTransaction',
        'http://localhost:8080/callbacks/bb3-ocpp-callback/stop-transaction'
    ),
    (
        'busways',
        'StatusNotification',
        'http://localhost:8080/callbacks/bb3-ocpp-callback/status-notification'
    ),
    (
        'busways',
        'MeterValues',
        'http://localhost:8080/callbacks/bb3-ocpp-callback/meter-values'
    );
insert into bb3.ocpp_application_entity (application_id, entity_code)
values ('cgw', 'B2070807'), ('busways', 'B7940015');
insert into bb3.ocpp_charge_point (
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
    ),(
        'B7940015',
        'ocpp1.6J',
        'SUTD_TEST',
        0
    );