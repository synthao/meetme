create table if not exists users
(
    id         bigint generated always as identity primary key,
    status     smallint  default 0,
    firstname  varchar(255)                        NOT NULL,
    lastname   varchar(255)                        NOT NULL,
    gender     SMALLINT                            NOT NULL,
    email      varchar(255)                        NOT NULL,
    birth_date date                                NOT NULL,
    created_at timestamp default current_timestamp NOT NULL,
    updated_at timestamp default current_timestamp NOT NULL
);

create table if not exists user_photos
(
    id         bigint generated always as identity primary key,
    user_id    bigint references users (id) ON DELETE CASCADE NOT NULL,
    name       varchar(255)                                   NOT NULL,
    created_at timestamp default current_timestamp            NOT NULL,
    updated_at timestamp default current_timestamp            NOT NULL
);