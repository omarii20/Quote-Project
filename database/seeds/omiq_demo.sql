-- OMIQ Demo Database Seed
-- Development and demo environments ONLY.
--
-- WARNING:
-- This script deletes 11 existing quotes for business_id = 3.
-- It preserves existing customers and businesses.
-- Do NOT execute against a production database.
--
-- This script is intended for one-time execution.

BEGIN;

DO $$
DECLARE
    v_business_id BIGINT := 3;
    v_customer_1 BIGINT;
    v_customer_2 BIGINT;
    v_customer_3 BIGINT;
    v_customer_4 BIGINT;
    v_customer_5 BIGINT;
    v_quote_id BIGINT;
    v_old_count INTEGER;
BEGIN
    -- Safety checks
    IF NOT EXISTS (
        SELECT 1 FROM businesses WHERE id = v_business_id
    ) THEN
        RAISE EXCEPTION 'Business 3 does not exist';
    END IF;

    SELECT COUNT(*) INTO v_old_count
    FROM quotes
    WHERE business_id = v_business_id
      AND id IN (6, 13, 14, 15, 16, 17, 18, 19, 20, 22, 24);

    IF v_old_count <> 11 THEN
        RAISE EXCEPTION
            'Expected 11 original quotes, found %',
            v_old_count;
    END IF;

    IF EXISTS (
        SELECT 1 FROM quotes
        WHERE business_id = v_business_id
          AND quote_number IN (
              'Q-0014', 'Q-0015', 'Q-0016',
              'Q-0017', 'Q-0018'
          )
    ) THEN
        RAISE EXCEPTION 'Demo quote numbers already exist';
    END IF;

    -- Delete ONLY the 11 original quotes of business 3.
    -- Their quote_items are deleted by ON DELETE CASCADE.
    DELETE FROM quotes
    WHERE business_id = v_business_id
      AND id IN (6, 13, 14, 15, 16, 17, 18, 19, 20, 22, 24);

    -- Demo customers: fictional names and contact details
    INSERT INTO customers
        (business_id, name, phone, email, address, notes)
    VALUES
        (3, 'Nova Digital', '0500000101',
         'demo-nova@example.com', 'תל אביב',
         'לקוח הדגמה - פיתוח אפליקציות')
    RETURNING id INTO v_customer_1;

    INSERT INTO customers
        (business_id, name, phone, email, address, notes)
    VALUES
        (3, 'Pixel Studio', '0500000102',
         'demo-pixel@example.com', 'חיפה',
         'לקוח הדגמה - אתרי אינטרנט')
    RETURNING id INTO v_customer_2;

    INSERT INTO customers
        (business_id, name, phone, email, address, notes)
    VALUES
        (3, 'TechFlow Solutions', '0500000103',
         'demo-techflow@example.com', 'הרצליה',
         'לקוח הדגמה - מערכות ארגוניות')
    RETURNING id INTO v_customer_3;

    INSERT INTO customers
        (business_id, name, phone, email, address, notes)
    VALUES
        (3, 'CloudNine', '0500000104',
         'demo-cloudnine@example.com', 'נצרת',
         'לקוח הדגמה - שירותי ענן ותחזוקה')
    RETURNING id INTO v_customer_4;

    INSERT INTO customers
        (business_id, name, phone, email, address, notes)
    VALUES
        (3, 'Urban Market', '0500000105',
         'demo-urban@example.com', 'ירושלים',
         'לקוח הדגמה - מסחר אלקטרוני')
    RETURNING id INTO v_customer_5;

    -- Q-0014: React Native ordering application
    -- Items: 12,000 + 10,000 + 6,000 + 4,000 = 32,000
    -- Fixed discount: 2,000
    -- Taxable amount: 30,000
    -- VAT: 5,400
    -- Total: 35,400

    INSERT INTO quotes (
        business_id, customer_id, quote_number,
        title, description, pricing_method,
        items_subtotal, subtotal,
        discount_type, discount_value, discount_amount,
        vat_rate, vat_amount, total,
        status, valid_until, notes
    ) VALUES (
        3, v_customer_1, 'Q-0014',
        'פיתוח אפליקציית הזמנות',
        'פיתוח אפליקציית React Native לניהול הזמנות, לקוחות ותשלומים, כולל ממשק ניהול וחיבור לשרת.',
        'items',
        32000, 32000,
        'fixed', 2000, 2000,
        18, 5400, 35400,
        'sent', '2026-11-30',
        'המחיר כולל פיתוח, בדיקות והטמעה. תנאי תשלום: 30% בתחילת הפרויקט, 40% במהלך הפיתוח ו-30% במסירה.'
    )
    RETURNING id INTO v_quote_id;

    INSERT INTO quote_items
        (quote_id, description, quantity, unit_price, total, position)
    VALUES
        (v_quote_id, 'אפיון ועיצוב UI/UX לאפליקציה', 1, 12000, 12000, 1),
        (v_quote_id, 'פיתוח אפליקציית React Native', 1, 10000, 10000, 2),
        (v_quote_id, 'פיתוח Backend API ומסד נתונים', 1, 6000, 6000, 3),
        (v_quote_id, 'בדיקות QA והכנה לפרסום', 1, 4000, 4000, 4);

    -- Q-0015: Business website
    -- Items subtotal: 11,000
    -- VAT: 1,980
    -- Total: 12,980

    INSERT INTO quotes (
        business_id, customer_id, quote_number,
        title, description, pricing_method,
        items_subtotal, subtotal,
        vat_rate, vat_amount, total,
        status, valid_until, notes
    ) VALUES (
        3, v_customer_2, 'Q-0015',
        'פיתוח אתר תדמית מקצועי',
        'עיצוב ופיתוח אתר תדמית רספונסיבי עם ממשק ניהול תוכן, התאמה למובייל ואופטימיזציה למנועי חיפוש.',
        'items',
        11000, 11000,
        18, 1980, 12980,
        'approved', '2026-12-15',
        'ההצעה כוללת התאמה מלאה לעברית ולמכשירים ניידים. אחסון ודומיין אינם כלולים במחיר.'
    )
    RETURNING id INTO v_quote_id;

    INSERT INTO quote_items
        (quote_id, description, quantity, unit_price, total, position)
    VALUES
        (v_quote_id, 'אפיון ועיצוב ממשק האתר', 1, 3500, 3500, 1),
        (v_quote_id, 'פיתוח אתר רספונסיבי', 1, 4500, 4500, 2),
        (v_quote_id, 'מערכת ניהול תוכן וקידום SEO בסיסי', 1, 3000, 3000, 3);

    -- Q-0016: CRM system, manual pricing
    -- Manual subtotal: 42,000
    -- Fixed discount: 2,000
    -- Taxable amount: 40,000
    -- VAT: 7,200
    -- Total: 47,200

    INSERT INTO quotes (
        business_id, customer_id, quote_number,
        title, description, pricing_method,
        manual_subtotal, subtotal,
        discount_type, discount_value, discount_amount,
        vat_rate, vat_amount, total,
        status, valid_until, notes
    ) VALUES (
        3, v_customer_3, 'Q-0016',
        'פיתוח מערכת CRM לניהול לקוחות',
        'פיתוח מערכת CRM מותאמת אישית הכוללת ניהול לקוחות, משימות, הרשאות משתמשים, דוחות ו-Dashboard ניהולי.',
        'manual',
        42000, 42000,
        'fixed', 2000, 2000,
        18, 7200, 47200,
        'draft', '2026-12-31',
        'משך פיתוח משוער: 10-12 שבועות. ההצעה כוללת אפיון, פיתוח, בדיקות והדרכת משתמשים.'
    );

    -- Q-0017: Monthly maintenance, manual pricing
    -- Manual subtotal: 3,500
    -- VAT: 630
    -- Total: 4,130

    INSERT INTO quotes (
        business_id, customer_id, quote_number,
        title, description, pricing_method,
        manual_subtotal, subtotal,
        vat_rate, vat_amount, total,
        status, valid_until, notes
    ) VALUES (
        3, v_customer_4, 'Q-0017',
        'חבילת תחזוקה ותמיכה חודשית',
        'שירותי תחזוקה שוטפת לאתר ולשרת, ניטור תקלות, עדכוני אבטחה, גיבויים ותמיכה טכנית.',
        'manual',
        3500, 3500,
        18, 630, 4130,
        'viewed', '2026-11-30',
        'המחיר מתייחס לחודש שירות אחד. השירות כולל עד 12 שעות תמיכה בחודש, ללא שעות נוספות.'
    );

    -- Q-0018: E-commerce platform
    -- Items: 6,500 + 5,500 + 4,500 + 3,500 = 20,000
    -- Fixed discount: 1,000
    -- Taxable amount: 19,000
    -- VAT: 3,420
    -- Total: 22,420

    INSERT INTO quotes (
        business_id, customer_id, quote_number,
        title, description, pricing_method,
        items_subtotal, subtotal,
        discount_type, discount_value, discount_amount,
        vat_rate, vat_amount, total,
        status, valid_until, notes
    ) VALUES (
        3, v_customer_5, 'Q-0018',
        'פיתוח חנות אינטרנטית',
        'הקמת חנות E-commerce הכוללת קטלוג מוצרים, סל קניות, תהליך רכישה, ניהול הזמנות וחיבור לספק תשלומים.',
        'items',
        20000, 20000,
        'fixed', 1000, 1000,
        18, 3420, 22420,
        'sent', '2026-12-20',
        'ההצעה כוללת עיצוב מותאם למותג, מערכת ניהול מוצרים והדרכה. עמלות סליקה ושירותי צד שלישי אינם כלולים.'
    )
    RETURNING id INTO v_quote_id;

    INSERT INTO quote_items
        (quote_id, description, quantity, unit_price, total, position)
    VALUES
        (v_quote_id, 'עיצוב UI/UX לחנות', 1, 6500, 6500, 1),
        (v_quote_id, 'פיתוח קטלוג מוצרים וסל קניות', 1, 5500, 5500, 2),
        (v_quote_id, 'פיתוח מערכת ניהול הזמנות', 1, 4500, 4500, 3),
        (v_quote_id, 'אינטגרציית תשלומים ובדיקות', 1, 3500, 3500, 4);

    RAISE NOTICE 'OMIQ demo quotes created successfully';
END $$;

-- Verify all five new quotes and their calculations.
SELECT
    q.quote_number,
    c.name AS customer,
    q.title,
    q.pricing_method,
    q.status,
    q.subtotal,
    q.discount_amount,
    q.vat_amount,
    q.total
FROM quotes q
JOIN customers c ON c.id = q.customer_id
WHERE q.business_id = 3
ORDER BY q.quote_number;

COMMIT;
