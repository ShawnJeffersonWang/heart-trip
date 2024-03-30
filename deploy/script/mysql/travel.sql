create table guess
(
    id           bigint auto_increment
        primary key,
    homestay_id  bigint                                  not null,
    price_after  bigint                                  not null,
    price_before bigint                                  not null,
    cover        varchar(4096) default ''                not null,
    location     varchar(4096) default ''                not null,
    title        varchar(4096) default ''                null,
    is_collected tinyint       default 0                 not null,
    udate_time   datetime      default CURRENT_TIMESTAMP not null,
    create_time  datetime      default CURRENT_TIMESTAMP not null
);

create table history
(
    id                   bigint auto_increment
        primary key,
    create_time          datetime      default CURRENT_TIMESTAMP not null,
    last_browsing_time   datetime      default CURRENT_TIMESTAMP not null,
    title                varchar(32)   default ''                not null,
    cover                varchar(4096) default ''                not null,
    intro                varchar(4096) default ''                not null,
    location             varchar(2048) default ''                not null,
    price_before         bigint        default 0                 not null,
    price_after          bigint        default 0                 not null,
    row_state            bigint        default 0                 not null,
    homestay_business_id bigint                                  not null,
    rating_stars         float                                   not null,
    user_id              bigint                                  not null,
    homestay_id          bigint                                  not null,
    del_state            tinyint       default 0                 not null,
    version              bigint        default 0                 not null,
    delete_time          datetime      default CURRENT_TIMESTAMP not null
);

create table homestay
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

create table homestay_activity
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
    on homestay_activity (row_type, row_status, del_state);

create table homestay_business
(
    id           bigint auto_increment
        primary key,
    create_time  datetime      default CURRENT_TIMESTAMP not null,
    update_time  datetime      default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP,
    delete_time  datetime      default CURRENT_TIMESTAMP not null,
    del_state    tinyint       default 0                 not null,
    title        varchar(32)   default ''                not null comment '店铺名称',
    user_id      bigint        default 0                 not null comment '关联的用户id',
    info         varchar(128)  default ''                not null comment '店铺介绍',
    boss_info    varchar(128)  default ''                not null comment '房东介绍',
    license_fron varchar(256)  default ''                not null comment '营业执照正面',
    license_back varchar(256)  default ''                not null comment '营业执照背面',
    row_state    tinyint(1)    default 0                 not null comment '0:禁止营业 1:正常营业',
    star         double(2, 1)  default 0.0               not null comment '店铺整体评价，冗余',
    tags         varchar(32)   default ''                not null comment '每个店家一个标签，自己编辑',
    cover        varchar(1024) default ''                not null comment '封面图',
    header_img   varchar(1024) default ''                not null comment '店招门头图片',
    version      bigint        default 0                 not null comment '版本号',
    constraint idx_userId
        unique (user_id)
)
    comment '民宿店铺';

create table homestay_comment
(
    id              bigint auto_increment
        primary key,
    create_time     datetime      default CURRENT_TIMESTAMP not null,
    update_time     datetime      default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP,
    delete_time     datetime      default CURRENT_TIMESTAMP not null,
    del_state       tinyint       default 0                 not null,
    homestay_id     bigint        default 0                 not null comment '民宿id',
    user_id         bigint        default 0                 not null comment '用户id',
    content         varchar(1024) default ''                not null comment '评论内容',
    star            json                                    not null comment '星星数,多个维度',
    version         bigint        default 0                 not null comment '版本号',
    nickname        varchar(32)   default ''                not null,
    avatar          varchar(1024) default ''                not null,
    image_urls      varchar(4096) default ''                not null,
    like_count      bigint        default 0                 not null,
    comment_time    varchar(255)  default ''                not null,
    tidy_rating     varchar(255)  default ''                not null,
    traffic_rating  varchar(255)  default ''                not null,
    security_rating varchar(255)  default ''                not null,
    food_rating     varchar(255)  default ''                not null,
    cost_rating     varchar(255)  default ''                not null
)
    comment '民宿评价';

create table user_history
(
    id          bigint auto_increment
        primary key,
    history_id  bigint                             not null,
    user_id     bigint                             not null,
    del_state   tinyint  default 0                 not null,
    version     bigint   default 0                 not null,
    delete_time datetime default CURRENT_TIMESTAMP not null
);

create table user_homestay
(
    id          bigint auto_increment
        primary key,
    user_id     bigint                             not null,
    homestay_id bigint                             not null,
    del_state   tinyint  default 0                 not null,
    version     bigint   default 0                 not null,
    delete_time datetime default CURRENT_TIMESTAMP not null
);

