create schema if not exists bb3;
create table bb3.ocpp_application (
    id serial primary key,
    uuid varchar(36) not null,
    name varchar(255) not null,
    unique (uuid)
);
create table bb3.ocpp_application_callback (
    id serial primary key,
    application_id int not null,
    callback_event varchar(255) not null,
    callback_url varchar(2048) not null,
    unique (application_id, callback_event)
);
create table bb3.ocpp_charge_point (
    id serial primary key,
    application_id int not null,
    /* fixed to charge point */
    charge_point_vendor varchar(20) not null,
    charge_point_model varchar(20) not null,
    charge_point_serial_number varchar(25) not null,
    charge_box_serial_number varchar(25) not null,
    iccid varchar(20) not null,
    imsi varchar(20) not null,
    meter_type varchar(25) not null,
    meter_serial_number varchar(25) not null,
    firmware_version varchar(50) not null,
    connector_count int not null,
    /* user-set */
    charge_point_identifier varchar(255) not null unique,
    ocpp_protocol varchar(20) not null,
    foreign key (application_id) references bb3.ocpp_application(id)
);
create table bb3.ocpp_charge_point_id_tag (
    id serial primary key,
    charge_point_id int not null,
    id_tag varchar(20) not null,
    foreign key (charge_point_id) references bb3.ocpp_charge_point(id)
);
create table bb3.ocpp_transaction (
    id serial primary key,
    charge_point_id int not null,
    connector_id int not null,
    id_tag varchar(20) not null,
    state varchar(50) not null,
    remote_initiated bool not null,
    start_timestamp timestamp without time zone not null,
    stop_timestamp timestamp without time zone not null,
    start_meter_value int not null,
    stop_meter_value int not null,
    stop_reason varchar(255) not null,
    foreign key (charge_point_id) references bb3.ocpp_charge_point(id)
);
create table bb3.ocpp_status_notification (
    id serial primary key,
    charge_point_id int not null,
    connector_id int not null,
    error_code varchar(255) not null,
    info varchar(255) not null,
    status varchar(255) not null,
    vendor_id varchar(255) not null,
    vendor_error_code varchar(255) not null,
    timestamp timestamp without time zone not null,
    reported_timestamp timestamp without time zone not null,
    foreign key (charge_point_id) references bb3.ocpp_charge_point(id)
);
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