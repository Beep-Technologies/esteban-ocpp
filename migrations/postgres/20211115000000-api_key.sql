-- +migrate Up
create table bb3.ocpp_application_api_key (
    api_key_hash varchar(255) primary key,
    application_id varchar(255) not null,
    description varchar(255) not null,
    is_active boolean not null,
    foreign key (application_id) references bb3.ocpp_application(id)
);
-- +migrate Down
