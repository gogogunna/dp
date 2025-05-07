--liquibase formatted sql
--changeset author:1

-- Таблица instrument_data
CREATE TABLE instrument_data (
                                 figi VARCHAR(20) PRIMARY KEY,
                                 name VARCHAR(255) NOT NULL,
                                 logo_path VARCHAR(512)
);

CREATE INDEX idx_instrument_data_figi ON instrument_data(figi);

-- Таблица user_token
CREATE TABLE user_token (
                            token VARCHAR(64) PRIMARY KEY,
                            user_id INTEGER NOT NULL,
);

CREATE INDEX idx_user_token_user_id ON user_token(user_id);

-- Таблица user_operations
CREATE TABLE user_operations (
                                 user_id INTEGER NOT NULL,
                                 operation_id INTEGER NOT NULL,
                                 operation_type INTEGER NOT NULL,
                                 operation_time TIMESTAMP WITH TIME ZONE NOT NULL,
                                 payment INTEGER NOT NULL,
                                 quantity INTEGER,
                                 price INTEGER,
                                 figi VARCHAR(12),
                                 PRIMARY KEY (user_id, operation_id)
);

CREATE INDEX idx_user_operations_composite ON user_operations(user_id, operation_id);
CREATE INDEX idx_user_operations_time ON user_operations(operation_time);
CREATE INDEX idx_user_operations_figi ON user_operations(figi) WHERE figi IS NOT NULL;