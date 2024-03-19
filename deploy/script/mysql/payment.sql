create table third_payment
(
    id               bigint auto_increment
        primary key,
    sn               char(25)     default ''                    not null comment '流水单号',
    create_time      datetime     default CURRENT_TIMESTAMP     not null,
    update_time      datetime     default CURRENT_TIMESTAMP     not null on update CURRENT_TIMESTAMP,
    delete_time      datetime     default CURRENT_TIMESTAMP     not null,
    del_state        tinyint(1)   default 0                     not null,
    version          bigint       default 0                     not null comment '乐观锁版本号',
    user_id          bigint       default 0                     not null comment '用户id',
    pay_mode         varchar(20)  default ''                    not null comment '支付方式 1:微信支付',
    trade_type       varchar(20)  default ''                    not null comment '第三方支付类型',
    trade_state      varchar(20)  default ''                    not null comment '第三方交易状态',
    pay_total        bigint       default 0                     not null comment '支付总金额(分)',
    transaction_id   char(32)     default ''                    not null comment '第三方支付单号',
    trade_state_desc varchar(256) default ''                    not null comment '支付状态描述',
    order_sn         char(25)     default ''                    not null comment '业务单号',
    service_type     varchar(32)  default ''                    not null comment '业务类型 ',
    pay_status       tinyint(1)   default 0                     not null comment '平台内交易状态   -1:支付失败 0:未支付 1:支付成功 2:已退款',
    pay_time         datetime     default '1970-01-01 08:00:00' not null comment '支付成功时间',
    constraint idx_sn
        unique (sn)
)
    comment '第三方支付流水记录';

