create table looklook_travel.homestay
(
    id                   bigint auto_increment
        primary key,
    create_time          datetime      default CURRENT_TIMESTAMP not null,
    update_time          datetime      default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP,
    delete_time          datetime      default CURRENT_TIMESTAMP not null,
    del_state            tinyint       default 0                 not null,
    version              bigint        default 0                 not null comment '版本号',
    title                varchar(32)   default ''                not null comment '标题',
    cover                varchar(4096) default ''                not null comment '轮播图，第一张封面',
    intro                varchar(4069) default ''                not null comment '介绍',
    location             varchar(2048) default ''                not null comment '位置',
    homestay_business_id bigint        default 0                 not null comment '民宿店铺id',
    user_id              bigint        default 0                 not null comment '房东id，冗余字段',
    row_state            tinyint(1)    default 0                 not null comment '0:下架 1:上架',
    rating_stars         float         default 0                 not null comment '评分',
    price_before         bigint        default 0                 not null comment '民宿价格（分）',
    price_after          bigint        default 0                 not null,
    clean_video          varchar(1024) default ''                not null,
    image_urls           varchar(4096) default ''                not null
)
    comment '每一间民宿';

