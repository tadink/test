create table if not exists actor
(
    id   int auto_increment
        primary key,
    name varchar(100) not null,
    constraint actor_unique
        unique (name)
);

create table if not exists actor_vod
(
    actor_id int null,
    vod_id   int null
);

create table if not exists director
(
    id   int auto_increment
        primary key,
    name varchar(50) not null,
    constraint director_unique
        unique (name)
);

create table if not exists director_vod
(
    director_id int null,
    vod_id      int null
);

create table if not exists tag
(
    id   int auto_increment
        primary key,
    name varchar(20) not null,
    constraint tag_unique
        unique (name)
);

create table if not exists tag_vod
(
    tag_id int null,
    vod_id int null
);

create table if not exists vod
(
    id           int auto_increment
        primary key,
    name         varchar(255)              not null,
    image        varchar(255)  default ''  not null,
    remark       varchar(100)  default ''  not null,
    rate         decimal(5, 1) default 0.0 not null,
    alias        varchar(255)  default ''  not null,
    area         varchar(30)   default ''  not null,
    language     varchar(20)   default ''  not null,
    publish_time int           default 0   not null,
    update_time  datetime                  null,
    source_url   varchar(255)  default ''  not null,
    description  text                      null,
    type_id      int           default 0   not null
);

create table if not exists vod_play_url
(
    id     int auto_increment
        primary key,
    name   varchar(50) default '' not null,
    url    varchar(255)           null,
    vod_id int                    not null
);

create table if not exists vod_type
(
    id        int auto_increment
        primary key,
    parent_id int default 0 not null,
    name      varchar(50)   not null
);

