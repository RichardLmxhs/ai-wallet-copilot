-- wallets: 保存钱包的基本信息与最近一次快照摘要
CREATE TABLE wallets
(
    address         TEXT PRIMARY KEY,                            -- checksummed address 建议在入库时统一格式
    chain           TEXT NOT NULL            DEFAULT 'ethereum', -- 'ethereum', 'arbitrum', etc.
    last_block      BIGINT,                                      -- 最后索引到的区块高度
    total_value_usd NUMERIC(30, 8),                              -- 快照计算的总市值（美元）
    token_count     INT                      DEFAULT 0,
    nft_count       INT                      DEFAULT 0,
    last_indexed_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    metadata        JSONB,                                       -- 可扩展：额外标签/备注
    created_at      TIMESTAMP WITH TIME ZONE DEFAULT now(),
    updated_at      TIMESTAMP WITH TIME ZONE DEFAULT now()
);

CREATE INDEX idx_wallets_chain_last_indexed ON wallets (chain, last_indexed_at);

-- wallet_assets: 每次快照的 token/nft 信息（一个 address + token 为一行）
CREATE TABLE wallet_assets
(
    id              BIGSERIAL PRIMARY KEY,
    wallet_address  TEXT NOT NULL REFERENCES wallets (address) ON DELETE CASCADE,
    chain           TEXT NOT NULL            DEFAULT 'ethereum',
    token_address   TEXT,            -- 对于 native token (ETH) 可为 NULL 或 'ETH'
    token_id        TEXT,            -- NFT 的 tokenId（string）
    token_type      TEXT,            -- 'ERC20' | 'ERC721' | 'ERC1155' | 'NATIVE'
    symbol          TEXT,
    name            TEXT,
    balance         NUMERIC(50, 18), -- token balance 原始精度
    balance_decimal INT,             -- token decimals
    balance_usd     NUMERIC(30, 8),  -- 估算 USD 值
    last_updated    TIMESTAMP WITH TIME ZONE DEFAULT now(),
    extra           JSONB            -- 存 metadata（image url, openSea data 等）
);

CREATE INDEX idx_wallet_assets_wallet ON wallet_assets (wallet_address);
CREATE INDEX idx_wallet_assets_tokenaddr ON wallet_assets (token_address);

-- wallet_transactions: 交易/tx 基本信息（按链 + txhash 唯一）
CREATE TABLE wallet_transactions
(
    id              BIGSERIAL PRIMARY KEY,
    chain           TEXT NOT NULL            DEFAULT 'ethereum',
    tx_hash         TEXT NOT NULL,
    block_number    BIGINT,
    block_timestamp TIMESTAMP WITH TIME ZONE,
    from_address    TEXT,
    to_address      TEXT,
    value           NUMERIC(50, 18), -- 原始 value (wei -> 转换时请统一单位)
    gas_used        NUMERIC(30, 0),
    gas_price       NUMERIC(30, 0),
    method          TEXT,            -- decode 后的方法名，如 "swapExactTokensForTokens"
    decoded         JSONB,           -- decode 后的 input 参数/结构
    status          SMALLINT,        -- 0/1 or null
    raw_tx          JSONB,           -- 原始 tx 结构（提供审计）
    wallet_address  TEXT,            -- 触发该 tx 与该 wallet 相关（索引查询用）
    created_at      TIMESTAMP WITH TIME ZONE DEFAULT now()
);

CREATE UNIQUE INDEX ux_wallet_transactions_chain_txhash ON wallet_transactions (chain, tx_hash);
CREATE INDEX idx_wallet_transactions_wallet ON wallet_transactions (wallet_address);
CREATE INDEX idx_wallet_transactions_block ON wallet_transactions (chain, block_number);

-- token_transfers: 交易内部的 token 转移（ERC20/ERC721 转账事件）
CREATE TABLE token_transfers
(
    id            BIGSERIAL PRIMARY KEY,
    chain         TEXT NOT NULL            DEFAULT 'ethereum',
    tx_hash       TEXT NOT NULL,
    log_index     INT,
    token_address TEXT,
    token_type    TEXT, -- ERC20/ERC721
    from_address  TEXT,
    to_address    TEXT,
    amount        NUMERIC(50, 18),
    token_id      TEXT,
    decimals      INT,
    metadata      JSONB,
    created_at    TIMESTAMP WITH TIME ZONE DEFAULT now()
);

CREATE UNIQUE INDEX ux_token_transfers_txhash_logindex ON token_transfers (tx_hash, log_index);
CREATE INDEX idx_token_transfers_to ON token_transfers (to_address);
CREATE INDEX idx_token_transfers_from ON token_transfers (from_address);
CREATE INDEX idx_token_transfers_token ON token_transfers (token_address);

-- contract_calls: 更细粒度的合约交互（如果你想按方法保留细节）
CREATE TABLE contract_calls
(
    id               BIGSERIAL PRIMARY KEY,
    chain            TEXT NOT NULL            DEFAULT 'ethereum',
    tx_hash          TEXT NOT NULL,
    contract_address TEXT,
    method           TEXT,
    args             JSONB,
    return_data      JSONB,
    created_at       TIMESTAMP WITH TIME ZONE DEFAULT now()
);

CREATE INDEX idx_contract_calls_contract ON contract_calls (contract_address);
CREATE INDEX idx_contract_calls_method ON contract_calls (method);

-- analysis_results: AI 分析的结果快照（每次请求或周期分析产生一条）
CREATE TABLE analysis_results
(
    id               BIGSERIAL PRIMARY KEY,
    wallet_address   TEXT NOT NULL REFERENCES wallets (address) ON DELETE CASCADE,
    chain            TEXT NOT NULL            DEFAULT 'ethereum',
    request_id       TEXT,  -- 如果有追踪 id（关联任务）
    model_name       TEXT,  -- e.g., "gpt-4o-mini"
    prompt_hash      TEXT,  -- 用于去重/版本化
    behavior_summary TEXT,
    risk_score       INT,   -- 0-100
    risk_details     JSONB, -- 结构化风险分析
    suggestions      JSONB, -- 建议数组
    raw_response     JSONB, -- LLM 原始输出
    created_at       TIMESTAMP WITH TIME ZONE DEFAULT now()
);

CREATE INDEX idx_analysis_wallet_created ON analysis_results (wallet_address, created_at DESC);

-- risk_flags: 规则引擎产生的结构化标记（便于搜索）
CREATE TABLE risk_flags
(
    id             BIGSERIAL PRIMARY KEY,
    wallet_address TEXT NOT NULL REFERENCES wallets (address) ON DELETE CASCADE,
    chain          TEXT NOT NULL            DEFAULT 'ethereum',
    flag_type      TEXT,  -- e.g., "scam", "mixer", "high_freq"
    score          INT,
    evidence       JSONB, -- 证据：tx_hash、addresses、timestamps
    created_at     TIMESTAMP WITH TIME ZONE DEFAULT now()
);

CREATE INDEX idx_risk_flags_wallet ON risk_flags (wallet_address);
CREATE INDEX idx_risk_flags_type ON risk_flags (flag_type);

-- price_cache: 保存 token 价格快照（用于估算 USD）
CREATE TABLE price_cache
(
    id            BIGSERIAL PRIMARY KEY,
    token_address TEXT,
    chain         TEXT NOT NULL            DEFAULT 'ethereum',
    price_usd     NUMERIC(30, 8),
    source        TEXT,
    ts            TIMESTAMP WITH TIME ZONE DEFAULT now()
);

CREATE INDEX idx_price_token_time ON price_cache (token_address, ts DESC);

-- jobs: 异步任务/批量作业队列（简单实现）
CREATE TABLE jobs
(
    id           BIGSERIAL PRIMARY KEY,
    job_type     TEXT NOT NULL,                              -- e.g., 'index_wallet', 'recalculate_risk', 'ai_analysis'
    payload      JSONB,
    status       TEXT NOT NULL            DEFAULT 'pending', -- pending/running/done/failed
    attempts     INT                      DEFAULT 0,
    last_error   TEXT,
    scheduled_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    created_at   TIMESTAMP WITH TIME ZONE DEFAULT now(),
    updated_at   TIMESTAMP WITH TIME ZONE DEFAULT now()
);

CREATE INDEX idx_jobs_status_schedule ON jobs (status, scheduled_at);

-- users: 如果你需要用户机制（可选）
CREATE TABLE users
(
    id            BIGSERIAL PRIMARY KEY,
    email         TEXT UNIQUE,
    password_hash TEXT,
    name          TEXT,
    created_at    TIMESTAMP WITH TIME ZONE DEFAULT now(),
    metadata      JSONB
);
