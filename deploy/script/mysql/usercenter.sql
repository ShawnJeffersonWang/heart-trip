create table message
(
    id           bigint auto_increment
        primary key,
    from_user_id bigint                             not null,
    to_user_id   bigint                             not null,
    content      text                               not null,
    del_state    tinyint  default 0                 not null,
    create_time  datetime default CURRENT_TIMESTAMP not null
);

create index idx_to_user_id_create_time on message (to_user_id, create_time);
create index idx_from_to_create_time on message (from_user_id, to_user_id, create_time);

create table user
(
    id          bigint auto_increment
        primary key,
    create_time datetime     default CURRENT_TIMESTAMP not null,
    update_time datetime     default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP,
    delete_time datetime     default CURRENT_TIMESTAMP not null,
    del_state   tinyint      default 0                 not null,
    version     bigint       default 0                 not null comment '版本号',
    mobile      char(11)     default ''                not null,
    password    varchar(255) default ''                not null,
    nickname    varchar(255) default ''                not null,
    sex         tinyint(1)   default 0                 not null comment '性别 0:男 1:女',
    avatar      varchar(255) default ''                not null,
    info        varchar(255) default ''                not null,
    constraint idx_mobile
        unique (mobile)
)
    comment '用户表';

create table user_auth
(
    id          bigint auto_increment
        primary key,
    create_time datetime    default CURRENT_TIMESTAMP not null,
    update_time datetime    default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP,
    delete_time datetime    default CURRENT_TIMESTAMP not null,
    del_state   tinyint     default 0                 not null,
    version     bigint      default 0                 not null comment '版本号',
    user_id     bigint      default 0                 not null,
    auth_key    varchar(64) default ''                not null comment '平台唯一id',
    auth_type   varchar(12) default ''                not null comment '平台类型',
    constraint idx_type_key
        unique (auth_type, auth_key),
    constraint idx_userId_key
        unique (user_id, auth_type)
)
    comment '用户授权表';

