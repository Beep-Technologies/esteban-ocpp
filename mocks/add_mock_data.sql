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
        'http://cgw-adapter/callbacks/ocpp/StartTransaction'
    ),
    (
        'cgw',
        'StopTransaction',
        'http://cgw-adapter/callbacks/ocpp/StopTransaction'
    ),
    (
        'cgw',
        'StatusNotification',
        'http://cgw-adapter/callbacks/ocpp/StatusNotification'
    );
insert into bb3.ocpp_application_entity (application_id, entity_code)
values ('cgw', 'L0332137');
insert into bb3.ocpp_charge_point (
        /* user-set */
        entity_code,
        ocpp_protocol,
        charge_point_identifier,
        connector_count
    )
values (
        'L0332137',
        'ocpp1.6J',
        'TACW745020G0274',
        0
    );