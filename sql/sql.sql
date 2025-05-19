USE ddpay;

CREATE TABLE merchant_wallets (
                                  id BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键ID',
                                  merchant_id BIGINT NOT NULL COMMENT '商户ID',
                                  chain VARCHAR(20) NOT NULL COMMENT '链名称，如 eth、bsc、solana、tron',
                                  address VARCHAR(100) NOT NULL COMMENT '派生的钱包地址',
                                  derivation_path VARCHAR(100) NOT NULL COMMENT 'HD钱包派生路径',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',

    PRIMARY KEY (id),
    UNIQUE KEY uniq_merchant_chain (merchant_id, chain),
    INDEX idx_address (address)


);


