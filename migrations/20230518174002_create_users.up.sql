CREATE TABLE users (
    id bigserial not null primary key,
    email varchar not null unique,
    name_user varchar not null, 
    encrypted_password varchar not null,
    role varchar not null,
    user_level integer not null default 1
);

Create TABLE tasks (
    id bigserial not null primary key,
    email_curator varchar not null,
    email_employee varchar not null,
    description varchar not null
);