CREATE TABLE IF NOT EXISTS quote_items (
    id BIGSERIAL PRIMARY KEY,

    quote_id BIGINT NOT NULL,

    description TEXT NOT NULL,

    quantity NUMERIC(10,2) NOT NULL DEFAULT 1,
    unit_price NUMERIC(12,2) NOT NULL DEFAULT 0,

    total NUMERIC(12,2) NOT NULL DEFAULT 0,
    total_overridden BOOLEAN NOT NULL DEFAULT FALSE,

    position INTEGER NOT NULL DEFAULT 0,

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_quote_items_quote
        FOREIGN KEY (quote_id)
        REFERENCES quotes(id)
        ON DELETE CASCADE,

    CONSTRAINT chk_quote_items_amounts
        CHECK (
            quantity > 0
            AND unit_price >= 0
            AND total >= 0
        ),

    CONSTRAINT chk_quote_items_position
        CHECK (
            position >= 0
        )
);