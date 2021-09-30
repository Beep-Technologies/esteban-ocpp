create schema if not exists bb3;
create table bb3.ocpp_address (
    id serial primary key,
    country_code char(2) not null,
    /* ISO 3166 Alpha-2 Code */
    city varchar(255),
    line_1 varchar(255),
    line_2 varchar(255),
    zip_code varchar(10)
);
create table bb3.ocpp_charge_point (
    id serial primary key,
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
    /* user-set */
    ocpp_protocol varchar(20) not null,
    charge_point_identifier varchar(255) not null,
    description varchar(255) not null,
    location_latitude decimal(11, 8) not null,
    location_longitude decimal(11, 8) not null,
    address_id int not null,
    foreign key (address_id) references bb3.ocpp_address(id)
);
create table bb3.ocpp_transaction (
    id serial primary key,
    charge_point_id int not null,
    connector_id int not null,
    id_tag varchar(20) not null,
    ongoing bool not null,
    state varchar(50) not null,
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
insert into bb3.ocpp_address (
        country_code,
        city,
        line_1,
        line_2,
        zip_code
    )
values (
        'SG',
        'Singapore',
        '8 Somapah Road',
        '',
        '487372'
    );
insert into bb3.ocpp_charge_point (
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
        description,
        location_latitude,
        location_longitude,
        address_id
    )
values (
        'Schneider Electric',
        'EVlink Smart Wallbox',
        'EVB1A22PCRI3N192421401400',
        '3N192040721A1S1B755170001',
        '',
        '',
        '',
        '',
        '3.3.0.16',
        'ocpp1.6J',
        'SUTD_TEST',
        'EVLink @ SUTD for Testing',
        1.3413647,
        103.9611493,
        1
    );