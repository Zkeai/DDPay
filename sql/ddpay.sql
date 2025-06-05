-- DDPay数据库初始化SQL
-- 合并自user.sql, membership.sql, subsite.sql和其他SQL文件

USE ddpay;

-- 用户表
CREATE TABLE IF NOT EXISTS `users` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `email` VARCHAR(255) NOT NULL,
    `password` VARCHAR(255) NOT NULL,
    `username` VARCHAR(255) NOT NULL,
    `avatar` VARCHAR(255) DEFAULT NULL,
    `role` VARCHAR(50) NOT NULL DEFAULT 'user',
    `status` TINYINT NOT NULL DEFAULT 1,
    `level` TINYINT NOT NULL DEFAULT 1 COMMENT '会员等级：1=青铜会员，2=白银会员，3=黄金会员，4=钻石会员',
    `email_verified` TINYINT NOT NULL DEFAULT 0,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `last_login_at` TIMESTAMP NULL DEFAULT NULL,
    `last_login_ip` VARCHAR(50) DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `email` (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 验证码表
CREATE TABLE IF NOT EXISTS `verification_codes` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `email` VARCHAR(255) NOT NULL,
    `code` VARCHAR(10) NOT NULL,
    `type` VARCHAR(50) NOT NULL,
    `expires_at` TIMESTAMP NOT NULL,
    `used` TINYINT NOT NULL DEFAULT 0,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    KEY `email_type_used` (`email`, `type`, `used`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- OAuth账号关联表
CREATE TABLE IF NOT EXISTS `oauth_accounts` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `user_id` BIGINT UNSIGNED NOT NULL,
    `provider` VARCHAR(50) NOT NULL,
    `provider_user_id` VARCHAR(255) NOT NULL,
    `provider_username` VARCHAR(255) DEFAULT NULL,
    `provider_email` VARCHAR(255) DEFAULT NULL,
    `provider_avatar` VARCHAR(255) DEFAULT NULL,
    `access_token` TEXT NOT NULL,
    `refresh_token` TEXT DEFAULT NULL,
    `token_expires_at` TIMESTAMP NULL DEFAULT NULL,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `provider_user_id` (`provider`, `provider_user_id`),
    KEY `user_id` (`user_id`),
    CONSTRAINT `oauth_accounts_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 登录日志表
CREATE TABLE IF NOT EXISTS `login_logs` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `user_id` BIGINT UNSIGNED DEFAULT NULL,
    `login_type` VARCHAR(50) NOT NULL,
    `ip` VARCHAR(50) NOT NULL,
    `user_agent` VARCHAR(255) NOT NULL,
    `status` TINYINT NOT NULL DEFAULT 0,
    `fail_reason` VARCHAR(255) DEFAULT NULL,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    KEY `user_id` (`user_id`),
    CONSTRAINT `login_logs_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 会员等级表
CREATE TABLE IF NOT EXISTS `membership_levels` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(50) NOT NULL COMMENT '等级名称',
    `level` INT NOT NULL COMMENT '等级数值',
    `icon` VARCHAR(255) NOT NULL COMMENT '等级图标URL',
    `price` DECIMAL(10,2) NOT NULL DEFAULT 0 COMMENT '升级价格',
    `description` VARCHAR(255) NOT NULL COMMENT '等级描述',
    `discount_rate` DECIMAL(3,2) NOT NULL DEFAULT 1.00 COMMENT '折扣率(0.1-1.0)',
    `max_subsites` INT NOT NULL DEFAULT 1 COMMENT '最大分站数量',
    `custom_service_access` TINYINT NOT NULL DEFAULT 0 COMMENT '专属客服权限',
    `vip_group_access` TINYINT NOT NULL DEFAULT 0 COMMENT 'VIP群权限',
    `priority` INT NOT NULL DEFAULT 1 COMMENT '优先级',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `level` (`level`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 会员权益表
CREATE TABLE IF NOT EXISTS `membership_benefits` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `level_id` BIGINT UNSIGNED NOT NULL COMMENT '关联的等级ID',
    `title` VARCHAR(100) NOT NULL COMMENT '权益标题',
    `description` VARCHAR(255) NOT NULL COMMENT '权益描述',
    `icon` VARCHAR(255) NOT NULL COMMENT '权益图标',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    KEY `level_id` (`level_id`),
    CONSTRAINT `membership_benefits_ibfk_1` FOREIGN KEY (`level_id`) REFERENCES `membership_levels` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 会员升级条件表
CREATE TABLE IF NOT EXISTS `membership_requirements` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `level_id` BIGINT UNSIGNED NOT NULL COMMENT '关联的等级ID',
    `type` VARCHAR(50) NOT NULL COMMENT '条件类型(充值金额/订单数/交易额/邀请人数)',
    `value` DECIMAL(10,2) NOT NULL COMMENT '条件值',
    `description` VARCHAR(255) NOT NULL COMMENT '条件描述',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    KEY `level_id` (`level_id`),
    CONSTRAINT `membership_requirements_ibfk_1` FOREIGN KEY (`level_id`) REFERENCES `membership_levels` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 用户会员记录表
CREATE TABLE IF NOT EXISTS `user_memberships` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    `level_id` BIGINT UNSIGNED NOT NULL COMMENT '会员等级ID',
    `start_date` TIMESTAMP NOT NULL COMMENT '开始日期',
    `end_date` TIMESTAMP NULL DEFAULT NULL COMMENT '结束日期(NULL表示永久)',
    `is_active` TINYINT NOT NULL DEFAULT 1 COMMENT '是否激活',
    `purchase_amount` DECIMAL(10,2) NOT NULL DEFAULT 0 COMMENT '购买金额',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    KEY `user_id` (`user_id`),
    KEY `level_id` (`level_id`),
    CONSTRAINT `user_memberships_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE,
    CONSTRAINT `user_memberships_ibfk_2` FOREIGN KEY (`level_id`) REFERENCES `membership_levels` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 会员交易记录表
CREATE TABLE IF NOT EXISTS `membership_transactions` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    `level_id` BIGINT UNSIGNED NOT NULL COMMENT '会员等级ID',
    `amount` DECIMAL(10,2) NOT NULL COMMENT '交易金额',
    `transaction_type` VARCHAR(20) NOT NULL COMMENT '交易类型(购买/续费/升级)',
    `payment_method` VARCHAR(50) NOT NULL COMMENT '支付方式',
    `status` VARCHAR(20) NOT NULL COMMENT '交易状态',
    `order_id` VARCHAR(64) NOT NULL COMMENT '订单号',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    KEY `user_id` (`user_id`),
    KEY `level_id` (`level_id`),
    UNIQUE KEY `order_id` (`order_id`),
    CONSTRAINT `membership_transactions_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE,
    CONSTRAINT `membership_transactions_ibfk_2` FOREIGN KEY (`level_id`) REFERENCES `membership_levels` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 分站表
CREATE TABLE IF NOT EXISTS `subsites` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `owner_id` BIGINT UNSIGNED NOT NULL COMMENT '所有者ID',
    `name` VARCHAR(100) NOT NULL COMMENT '分站名称',
    `subdomain` VARCHAR(50) NOT NULL COMMENT '子域名',
    `domain` VARCHAR(100) DEFAULT NULL COMMENT '自定义域名',
    `description` TEXT DEFAULT NULL COMMENT '分站描述',
    `logo` VARCHAR(255) DEFAULT NULL COMMENT '分站Logo',
    `theme` VARCHAR(50) DEFAULT 'default' COMMENT '分站主题',
    `status` TINYINT NOT NULL DEFAULT 1 COMMENT '状态：0=禁用，1=启用',
    `commission_rate` DECIMAL(5,2) NOT NULL DEFAULT 10.00 COMMENT '佣金比例(%)',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `subdomain` (`subdomain`),
    UNIQUE KEY `domain` (`domain`),
    KEY `owner_id` (`owner_id`),
    CONSTRAINT `subsites_ibfk_1` FOREIGN KEY (`owner_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 分站配置表
CREATE TABLE IF NOT EXISTS `subsite_configs` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `subsite_id` BIGINT UNSIGNED NOT NULL COMMENT '分站ID',
    `config` TEXT NOT NULL COMMENT 'JSON格式的配置',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `subsite_id` (`subsite_id`),
    CONSTRAINT `subsite_configs_ibfk_1` FOREIGN KEY (`subsite_id`) REFERENCES `subsites` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 分站商品表
CREATE TABLE IF NOT EXISTS `subsite_products` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `subsite_id` BIGINT UNSIGNED NOT NULL COMMENT '分站ID',
    `title` VARCHAR(100) NOT NULL COMMENT '商品标题',
    `description` TEXT DEFAULT NULL COMMENT '商品描述',
    `price` DECIMAL(10,2) NOT NULL COMMENT '商品价格',
    `image` VARCHAR(255) DEFAULT NULL COMMENT '商品图片',
    `stock` INT NOT NULL DEFAULT -1 COMMENT '库存，-1表示无限',
    `is_time_limited` TINYINT NOT NULL DEFAULT 0 COMMENT '是否限时：0=否，1=是',
    `start_time` TIMESTAMP NULL DEFAULT NULL COMMENT '开始时间',
    `end_time` TIMESTAMP NULL DEFAULT NULL COMMENT '结束时间',
    `status` TINYINT NOT NULL DEFAULT 1 COMMENT '状态：0=下架，1=上架',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    KEY `subsite_id` (`subsite_id`),
    CONSTRAINT `subsite_products_ibfk_1` FOREIGN KEY (`subsite_id`) REFERENCES `subsites` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 分站订单表
CREATE TABLE IF NOT EXISTS `subsite_orders` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `order_no` VARCHAR(32) NOT NULL COMMENT '订单号',
    `subsite_id` BIGINT UNSIGNED NOT NULL COMMENT '分站ID',
    `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    `product_id` BIGINT UNSIGNED NOT NULL COMMENT '商品ID',
    `quantity` INT NOT NULL DEFAULT 1 COMMENT '购买数量',
    `amount` DECIMAL(10,2) NOT NULL COMMENT '订单金额',
    `commission` DECIMAL(10,2) NOT NULL COMMENT '佣金金额',
    `status` TINYINT NOT NULL DEFAULT 0 COMMENT '状态：0=待支付，1=已支付，2=已完成，3=已取消',
    `pay_time` TIMESTAMP NULL DEFAULT NULL COMMENT '支付时间',
    `complete_time` TIMESTAMP NULL DEFAULT NULL COMMENT '完成时间',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `order_no` (`order_no`),
    KEY `subsite_id` (`subsite_id`),
    KEY `user_id` (`user_id`),
    KEY `product_id` (`product_id`),
    CONSTRAINT `subsite_orders_ibfk_1` FOREIGN KEY (`subsite_id`) REFERENCES `subsites` (`id`) ON DELETE CASCADE,
    CONSTRAINT `subsite_orders_ibfk_2` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE,
    CONSTRAINT `subsite_orders_ibfk_3` FOREIGN KEY (`product_id`) REFERENCES `subsite_products` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 分站余额表
CREATE TABLE IF NOT EXISTS `subsite_balances` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `owner_id` BIGINT UNSIGNED NOT NULL COMMENT '所有者ID',
    `amount` DECIMAL(10,2) NOT NULL DEFAULT 0.00 COMMENT '余额',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `owner_id` (`owner_id`),
    CONSTRAINT `subsite_balances_ibfk_1` FOREIGN KEY (`owner_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 分站余额变动记录表
CREATE TABLE IF NOT EXISTS `subsite_balance_logs` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `owner_id` BIGINT UNSIGNED NOT NULL COMMENT '所有者ID',
    `order_id` BIGINT UNSIGNED DEFAULT NULL COMMENT '关联订单ID',
    `amount` DECIMAL(10,2) NOT NULL COMMENT '变动金额',
    `before_balance` DECIMAL(10,2) NOT NULL COMMENT '变动前余额',
    `after_balance` DECIMAL(10,2) NOT NULL COMMENT '变动后余额',
    `type` VARCHAR(20) NOT NULL COMMENT '类型：commission=佣金收入，withdrawal=提现',
    `remark` VARCHAR(255) DEFAULT NULL COMMENT '备注',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    KEY `owner_id` (`owner_id`),
    CONSTRAINT `subsite_balance_logs_ibfk_1` FOREIGN KEY (`owner_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 分站提现申请表
CREATE TABLE IF NOT EXISTS `subsite_withdrawals` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `owner_id` BIGINT UNSIGNED NOT NULL COMMENT '所有者ID',
    `amount` DECIMAL(10,2) NOT NULL COMMENT '提现金额',
    `bank_name` VARCHAR(100) NOT NULL COMMENT '银行名称',
    `bank_account` VARCHAR(50) NOT NULL COMMENT '银行账号',
    `account_name` VARCHAR(50) NOT NULL COMMENT '开户名',
    `status` TINYINT NOT NULL DEFAULT 0 COMMENT '状态：0=待处理，1=已处理，2=已拒绝',
    `admin_remark` VARCHAR(255) DEFAULT NULL COMMENT '管理员备注',
    `processed_at` TIMESTAMP NULL DEFAULT NULL COMMENT '处理时间',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    KEY `owner_id` (`owner_id`),
    CONSTRAINT `subsite_withdrawals_ibfk_1` FOREIGN KEY (`owner_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 商户钱包表
CREATE TABLE IF NOT EXISTS `merchant_wallets` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `merchant_id` BIGINT NOT NULL COMMENT '商户ID',
    `chain` VARCHAR(20) NOT NULL COMMENT '链名称，如 eth、bsc、solana、tron',
    `address` VARCHAR(100) NOT NULL COMMENT '派生的钱包地址',
    `derivation_path` VARCHAR(100) NOT NULL COMMENT 'HD钱包派生路径',
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uniq_merchant_chain` (`merchant_id`, `chain`),
    INDEX `idx_address` (`address`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 插入默认会员等级数据
INSERT INTO `membership_levels` (`name`, `level`, `icon`, `price`, `description`, `discount_rate`, `max_subsites`, `custom_service_access`, `vip_group_access`, `priority`) VALUES
('青铜会员', 1, '/assets/membership/bronze.png', 0.00, '基础会员级别，享受平台基础服务', 1.00, 1, 0, 0, 1),
('白银会员', 2, '/assets/membership/silver.png', 99.00, '进阶会员级别，享受更多权益与折扣', 0.90, 3, 0, 1, 2),
('黄金会员', 3, '/assets/membership/gold.png', 299.00, '高级会员，享受VIP待遇与最大折扣', 0.80, 10, 1, 1, 3),
('钻石会员', 4, '/assets/membership/diamond.png', 999.00, '顶级会员待遇，无限制使用所有功能', 0.70, -1, 1, 1, 4);

-- 插入会员权益数据
INSERT INTO `membership_benefits` (`level_id`, `title`, `description`, `icon`) VALUES
(1, '创建1个分站', '可创建并管理1个分站', '/assets/membership/site.png'),
(1, '标准技术支持', '工作日9:00-18:00客服支持', '/assets/membership/support.png'),

(2, '创建3个分站', '可创建并管理最多3个分站', '/assets/membership/sites.png'),
(2, '9折优惠', '所有交易享受9折优惠', '/assets/membership/discount.png'),
(2, '会员专属群', '加入白银会员专属交流群', '/assets/membership/group.png'),

(3, '创建10个分站', '可创建并管理最多10个分站', '/assets/membership/sites.png'),
(3, '8折优惠', '所有交易享受8折优惠', '/assets/membership/discount.png'),
(3, '专属客服', '一对一专属客服7*12小时服务', '/assets/membership/vip-support.png'),
(3, 'VIP会员群', '加入金牌会员VIP交流群', '/assets/membership/vip-group.png'),
(3, '优先处理', '订单优先处理，技术问题优先解决', '/assets/membership/priority.png'),

(4, '无限分站', '可创建无限数量的分站', '/assets/membership/unlimited.png'),
(4, '7折优惠', '所有交易享受7折最大优惠', '/assets/membership/discount.png'),
(4, '24小时专属客服', '一对一专属客服全天候服务', '/assets/membership/vip-support.png'),
(4, '钻石VIP群', '加入钻石会员核心交流群', '/assets/membership/vip-group.png'),
(4, '最高优先级', '享受系统最高优先级处理所有请求', '/assets/membership/priority.png'),
(4, '专属定制', '享受专属定制开发服务', '/assets/membership/custom.png');

-- 插入会员升级条件数据
INSERT INTO `membership_requirements` (`level_id`, `type`, `value`, `description`) VALUES
(1, 'register', 0.00, '注册成为会员即可'),

(2, 'payment', 99.00, '一次性支付99元升级'),
(2, 'total_order', 5.00, '累计订单达到5笔'),

(3, 'payment', 299.00, '一次性支付299元升级'),
(3, 'total_payment', 1000.00, '累计充值金额达到1000元'),
(3, 'invitation', 3.00, '成功邀请3名用户注册'),

(4, 'payment', 999.00, '一次性支付999元升级'),
(4, 'total_payment', 5000.00, '累计充值金额达到5000元'),
(4, 'total_transaction', 10000.00, '累计交易额达到10000元'),
(4, 'invitation', 10.00, '成功邀请10名用户注册'); 