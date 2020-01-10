create table "Files"(
    f_id bigserial not null primary key,
    f_parent int,
    f_name varchar,
    f_area varchar,
    f_type int,
    f_hash varchar,
    f_size int,
    f_date_create timestamp,
    f_date_create_from_source timestamp
);