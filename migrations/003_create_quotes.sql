CREATE TABLE IF NOT EXISTS quotes (
    id BIGSERIAL PRIMARY KEY,

    business_id BIGINT NOT NULL,
    customer_id BIGINT NOT NULL,

    quote_number TEXT NOT NULL,
    title TEXT,
    description TEXT,

    pricing_method TEXT NOT NULL DEFAULT 'items',

    items_subtotal NUMERIC(12,2) NOT NULL DEFAULT 0,
    manual_subtotal NUMERIC(12,2),
    additional_amount NUMERIC(12,2) NOT NULL DEFAULT 0,
    subtotal NUMERIC(12,2) NOT NULL DEFAULT 0,

    discount_type TEXT,
    discount_value NUMERIC(12,2) NOT NULL DEFAULT 0,
    discount_amount NUMERIC(12,2) NOT NULL DEFAULT 0,

    vat_rate NUMERIC(5,2) NOT NULL DEFAULT 0,
    vat_amount NUMERIC(12,2) NOT NULL DEFAULT 0,

    total NUMERIC(12,2) NOT NULL DEFAULT 0,

    status TEXT NOT NULL DEFAULT 'draft',

    valid_until DATE,
    notes TEXT,

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_quotes_business
        FOREIGN KEY (business_id)
        REFERENCES businesses(id)
        ON DELETE RESTRICT,

    CONSTRAINT fk_quotes_customer
        FOREIGN KEY (customer_id)
        REFERENCES customers(id)
        ON DELETE RESTRICT,

    CONSTRAINT uq_quotes_business_quote_number
        UNIQUE (business_id, quote_number),

    CONSTRAINT chk_quotes_pricing_method
        CHECK (
            pricing_method IN ('items', 'manual')
        ),

    CONSTRAINT chk_quotes_discount_type
        CHECK (
            discount_type IS NULL
            OR discount_type IN ('percent', 'fixed')
        ),

    CONSTRAINT chk_quotes_status
        CHECK (
            status IN (
                'draft',
                'sent',
                'viewed',
                'approved',
                'rejected',
                'expired'
            )
        ),

    CONSTRAINT chk_quotes_amounts
        CHECK (
            items_subtotal >= 0
            AND (manual_subtotal IS NULL OR manual_subtotal >= 0)
            AND additional_amount >= 0
            AND subtotal >= 0
            AND discount_value >= 0
            AND discount_amount >= 0
            AND vat_rate >= 0
            AND vat_amount >= 0
            AND total >= 0
        )
);