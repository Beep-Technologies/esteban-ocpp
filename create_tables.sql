create schema if not exists bb3;
create table bb3.address (
    id int not null,
    country_code char(2) not null,
    /* ISO 3166 Alpha-2 Code */
    city varchar(255),
    line_1 varchar(255),
    line_2 varchar(255),
    zip_code varchar(10),
    primary key (id)
);
create table bb3.charge_point (
    id int not null,
    /* fixed to charge point */
    charge_point_vendor varchar(20) not null,
    charge_point_model varchar(20) not null,
    charge_point_serial_number varchar(25) not null,
    charge_box_serial_number varchar(20) not null,
    iccid varchar(20) not null,
    imsi varchar(20) not null,
    meter_type varchar(25) not null,
    meter_serial_number varchar(25) not null,
    firmware_version varchar(50) not null,
    /* user-set */
    ocpp_protocol varchar(20) not null,
    charge_point_identifier varchar(255) not null,
    description varchar(255) not null,
    location_latitude decimal(11, 8),
    location_longitude decimal(11, 8),
    address_id int not null,
    primary key (id),
    foreign key (address_id) references bb3.address(id)
);
create table bb3.transaction (
    id int not null,
    charge_point_id int not null,
    connector_id int not null,
    id_tag varchar(20) not null,
    started bool not null,
    stopped bool not null,
    start_timestamp timestamp without time zone not null,
    stop_timestamp timestamp without time zone not null,
    start_meter_value int not null,
    stop_meter_value int not null,
    stop_reason varchar(255) not null,
    primary key (id),
    foreign key (charge_point_id) references bb3.charge_point(id)
);
create table bb3.status_notification (
    id int not null,
    charge_point_id int not null,
    connector_id int not null,
    error_code varchar(255) not null,
    info varchar(255) not null,
    status varchar(255) not null,
    vendor_id varchar(255) not null,
    vendor_error_code varchar(255) not null,
    timestamp timestamp without time zone not null,
    reported_timestamp timestamp without time zone not null,
    primary key (id),
    foreign key (charge_point_id) references bb3.charge_point(id)
);