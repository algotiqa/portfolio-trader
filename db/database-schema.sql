-- ======================================================================
-- ===   Sql Script for Database : Portfolio Trader
-- ===
-- === Build : 317
-- ======================================================================

CREATE TABLE portfolio
  (
    id        int,
    username  varchar(32)   not null,
    name      varchar(64)   not null,

    primary key(id)
  )
 ENGINE = InnoDB ;

CREATE INDEX portfolioIDX1 ON portfolio(username);

-- ======================================================================

CREATE TABLE trading_system
  (
    id                  int,
    username            varchar(32)     not null,
    name                varchar(64)     not null,
    timeframe           int             not null,
    data_product_id     int             not null,
    data_symbol         varchar(16)     not null,
    broker_product_id   int             not null,
    broker_symbol       varchar(16)     not null,
    point_value         double          not null,
    cost_per_operation  double          not null,
    margin_value        double          not null,
    increment           double          not null,
    market_type         char(2)         not null,
    currency_id         int             not null,
    currency_code       varchar(16)     not null,
    currency_symbol     varchar(4)      not null,
    trading_session_id  int             not null,
    session_name        varchar(64)     not null,
    session_config      varchar(1024)   not null,
    agent_profile_id    int,
    external_ref        varchar(64),
    strategy_type       char(2)         not null,
    overnight           tinyint         not null,
    tags                varchar(255),
    finalized           tinyint         not null,
    trading             tinyint         not null,
    running             tinyint         not null,
    auto_activation     tinyint         not null,
    active              tinyint         not null,
    status              tinyint         not null,
    suggested_action    tinyint         not null,
    first_trade         datetime,
    last_trade          datetime,
    last_net_profit     double,
    last_net_avg_trade  double,
    last_num_trades     int,
    portfolio_id        int,
    timezone            varchar(32)     not null,
    in_sample_from      int             not null,
    in_sample_to        int             not null,
    engine_code         varchar(16)     not null,

    primary key(id),

    foreign key(portfolio_id) references portfolio(id)
  )
 ENGINE = InnoDB ;

CREATE INDEX trading_systemIDX1 ON trading_system(username);

-- ======================================================================

CREATE TABLE trading_filter
  (
    trading_system_id  int,
    equ_avg_enabled    tinyint    not null,
    equ_avg_len        smallint   not null,
    pos_pro_enabled    tinyint    not null,
    pos_pro_len        smallint   not null,
    win_per_enabled    tinyint    not null,
    win_per_len        tinyint    not null,
    win_per_value      tinyint    not null,
    old_new_enabled    tinyint    not null,
    old_new_old_len    smallint   not null,
    old_new_old_perc   smallint   not null,
    old_new_new_len    smallint   not null,
    trendline_enabled  tinyint    not null,
    trendline_len      smallint   not null,
    trendline_value    smallint   not null,
    drawdown_enabled   tinyint    not null,
    drawdown_min       smallint   not null,
    drawdown_max       smallint   not null,

    primary key(trading_system_id),

    foreign key(trading_system_id) references trading_system(id)
  )
 ENGINE = InnoDB ;

-- ======================================================================

CREATE TABLE trading_position
  (
    trading_system_id   int,
    initial_capital     double          not null,
    max_tol_drawd_perc  double          not null,
    margin_override     double,
    max_units           int             not null,
    risk_per_unit       varchar(16)     not null,
    risk_value          double          not null,
    model               char(2)         not null,
    config              varchar(1024)   not null,

    primary key(trading_system_id),

    foreign key(trading_system_id) references trading_system(id)
  )
 ENGINE = InnoDB ;

-- ======================================================================

CREATE TABLE trade
  (
    id                     bigint        auto_increment,
    trading_system_id      int           not null,
    trade_type             char(2)       not null,
    entry_date             datetime      not null,
    entry_price            double        not null,
    entry_label            varchar(64),
    exit_date              datetime,
    exit_price             double,
    exit_label             varchar(64),
    gross_return           double,
    max_contracts          int           not null,
    entry_date_at_broker   datetime,
    entry_price_at_broker  double,
    exit_date_at_broker    datetime,
    exit_price_at_broker   double,

    primary key(id),

    foreign key(trading_system_id) references trading_system(id)
  )
 ENGINE = InnoDB ;

CREATE INDEX tradeIDX1 ON trade(trading_system_id);

-- ======================================================================

CREATE TABLE equity_bar
  (
    trade_id      bigint     not null,
    date          datetime   not null,
    gross_return  double     not null,
    contracts     int        not null
  )
 ENGINE = InnoDB partition by hash(trade_id) partitions 16;

CREATE INDEX equity_barIDX1 ON equity_bar(trade_id);

-- ======================================================================

CREATE TABLE live_period
  (
    id                 int        auto_increment,
    trading_system_id  int        not null,
    period             datetime   not null,
    active             tinyint    not null,

    primary key(id),

    foreign key(trading_system_id) references trading_system(id)
  )
 ENGINE = InnoDB ;

CREATE INDEX live_periodIDX1 ON live_period(trading_system_id);

-- ======================================================================

CREATE TABLE outbox
  (
    id         int           auto_increment,
    timestamp  datetime      not null,
    exchange   varchar(64)   not null,
    uuid       varchar(64)   not null,
    payload    mediumblob    not null,
    size       int           not null,

    primary key(id)
  )
 ENGINE = InnoDB ;

-- ======================================================================

