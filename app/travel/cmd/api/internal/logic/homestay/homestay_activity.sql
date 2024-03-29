create table looklook_travel.homestay_activity
(
    id          bigint auto_increment
        primary key,
    create_time datetime    default CURRENT_TIMESTAMP not null,
    update_time datetime    default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP,
    delete_time datetime    default CURRENT_TIMESTAMP not null,
    del_state   tinyint     default 0                 not null,
    row_type    varchar(32) default ''                not null comment '活动类型',
    data_id     bigint      default 0                 not null comment '业务表id（id跟随活动类型走）',
    row_status  tinyint(1)  default 0                 not null comment '0:下架 1:上架',
    version     bigint      default 0                 not null comment '版本号'
)
    comment '每一间民宿';

create index idx_rowType
    on looklook_travel.homestay_activity (row_type, row_status, del_state);

