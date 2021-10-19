-- add mock data
insert into bb3.ocpp_application (
    uuid, 
    name
    )
values (
        'cde0496a-bcd8-408e-9d07-65d07d841487',
        'busways'
    );
insert into bb3.ocpp_charge_point (
        application_id,
        charge_point_vendor,
        charge_point_model,
        charge_point_serial_number,
        charge_box_serial_number,
        iccid,
        imsi,
        meter_type,
        meter_serial_number,
        firmware_version,
        /* user-set */
        ocpp_protocol,
        charge_point_identifier,
        connector_count
    ) 
values (
        1,
        '',
        '',
        '',
        '',
        '',
        '',
        '',
        '',
        '',
        'ocpp1.6J',
        'SUTD_TEST',
        0
    );