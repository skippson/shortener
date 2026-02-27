create table if not exists urls (
    id serial primary key,
    original text not null unique,
    shortened varchar(10) not null unique,
    created_at timestamp not null default now()
);