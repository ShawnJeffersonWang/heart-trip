create table homestay_order
(
    id                   bigint auto_increment
        primary key,
    create_time          datetime      default CURRENT_TIMESTAMP not null,
    update_time          datetime      default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP,
    delete_time          datetime      default CURRENT_TIMESTAMP not null,
    del_state            tinyint       default 0                 not null,
    version              bigint        default 0                 not null comment '版本号',
    sn                   char(25)      default ''                not null comment '订单号',
    user_id              bigint        default 0                 not null comment '下单用户id',
    homestay_id          bigint        default 0                 not null comment '民宿id',
    title                varchar(32)   default ''                not null comment '标题',
    cover                varchar(1024) default ''                not null comment '封面',
    info                 varchar(4069) default ''                not null comment '介绍',
    homestay_price       bigint                                  not null comment '民宿价格(分)',
    homestay_business_id bigint        default 0                 not null comment '店铺id',
    homestay_user_id     bigint        default 0                 not null comment '店铺房东id',
    live_start_date      date                                    not null comment '开始入住日期',
    live_end_date        date                                    not null comment '结束入住日期',
    trade_state          tinyint(1)    default 0                 not null comment '-1: 已取消 0:待支付 1:未使用 2:已使用  3:已退款 4:已过期',
    trade_code           char(8)       default ''                not null comment '确认码',
    remark               varchar(64)   default ''                not null comment '用户下单备注',
    order_total_price    bigint        default 0                 not null comment '订单总价格（餐食总价格+民宿总价格）(分)',
    homestay_total_price bigint        default 0                 not null comment '民宿总价格(分)',
    constraint idx_sn
        unique (sn)
)
    comment '每一间民宿';

