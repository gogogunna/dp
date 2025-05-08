--liquibase formatted sql
--changeset author:2

-- Универсальная функция для работы с токенами (получение или создание)
CREATE OR REPLACE FUNCTION get_or_create_user_id(
    p_token VARCHAR(64)
RETURNS INTEGER AS $$
DECLARE
    v_user_id INTEGER;
BEGIN
    -- SERIALIZABLE - максимальный уровень изоляции для гарантии консистентности
    SET TRANSACTION ISOLATION LEVEL SERIALIZABLE;

    -- Пытаемся найти существующий токен
    SELECT user_id INTO v_user_id
    FROM user_token
    WHERE token = p_token;

    -- Если не нашли - создаем новую запись
    IF v_user_id IS NULL THEN
        -- Генерируем новый user_id через последовательность
        INSERT INTO user_token (token, user_id)
        VALUES (
            p_token,
            nextval(pg_get_serial_sequence('user_token', 'user_id'))
        RETURNING user_id INTO v_user_id;
    END IF;

    RETURN v_user_id;
END;
$$ LANGUAGE plpgsql;

-- Батчевая вставка инструментов
CREATE OR REPLACE FUNCTION insert_instruments_batch(
    p_instruments JSONB
) RETURNS VOID AS $$
BEGIN
INSERT INTO instrument_data (figi, name, logo_path)
SELECT
    instrument->>'figi',
    instrument->>'name',
    instrument->>'logo_path'
FROM jsonb_array_elements(p_instruments) AS instrument
ON CONFLICT (figi) DO UPDATE SET
    name = EXCLUDED.name,
                          logo_path = EXCLUDED.logo_path;
END;
$$ LANGUAGE plpgsql;

-- Батчевое получение инструментов
CREATE OR REPLACE FUNCTION get_instruments_batch(
    p_figi_list JSONB
) RETURNS JSONB AS $$
BEGIN
RETURN (
    SELECT jsonb_agg(jsonb_build_object(
            'figi', figi,
            'name', name,
            'logo_path', logo_path
                     ))
    FROM instrument_data
    WHERE figi = ANY(ARRAY(SELECT jsonb_array_elements_text(p_figi_list)))
);
END;
$$ LANGUAGE plpgsql;

-- Батчевая вставка операций
CREATE OR REPLACE FUNCTION insert_operations_batch(
    p_operations JSONB
) RETURNS INTEGER AS $$
DECLARE
inserted_count INTEGER;
BEGIN
WITH inserted AS (
INSERT INTO user_operations (
    user_id,
    operation_id,
    operation_type,
    operation_time,
    payment,
    quantity,
    price,
    figi
)
SELECT
    (op->>'user_id')::INTEGER,
        (op->>'operation_id')::INTEGER,
        (op->>'operation_type')::INTEGER,
        (op->>'operation_time')::TIMESTAMP WITH TIME ZONE,
            (op->>'payment')::INTEGER,
            (op->>'quantity')::INTEGER,
            (op->>'price')::INTEGER,
            op->>'figi'
FROM jsonb_array_elements(p_operations) AS op
ON CONFLICT (user_id, operation_id) DO NOTHING
    RETURNING 1
    )
SELECT COUNT(*) INTO inserted_count FROM inserted;

RETURN inserted_count;
END;
$$ LANGUAGE plpgsql;

-- Батчевое получение операций с флагом новых операций
CREATE OR REPLACE FUNCTION get_operations_batch(
    p_user_id INTEGER,
    p_from_time TIMESTAMP WITH TIME ZONE,
    p_to_time TIMESTAMP WITH TIME ZONE
) RETURNS TABLE (
    operations JSONB,
    has_newer_operations BOOLEAN
) AS $$
BEGIN
RETURN QUERY
    WITH batch AS (
        SELECT jsonb_agg(jsonb_build_object(
            'user_id', user_id,
            'operation_id', operation_id,
            'operation_type', operation_type,
            'operation_time', operation_time,
            'payment', payment,
            'quantity', quantity,
            'price', price,
            'figi', figi
        )) AS data
        FROM user_operations
        WHERE
            user_id = p_user_id AND
            operation_time >= p_from_time AND
            operation_time <= p_to_time
        ORDER BY operation_time
    ),
    newer_check AS (
        SELECT EXISTS (
            SELECT 1 FROM user_operations
            WHERE
                user_id = p_user_id AND
                operation_time > p_to_time
            LIMIT 1
        ) AS has_newer
    )
SELECT
    COALESCE(batch.data, '[]'::jsonb),
    newer_check.has_newer
FROM batch, newer_check;
END;
$$ LANGUAGE plpgsql;