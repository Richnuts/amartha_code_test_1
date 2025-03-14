BEGIN;

-- mock user
CREATE TABLE IF NOT EXISTS "user" (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);

CREATE TYPE "duration_unit" AS ENUM('WEEK');

CREATE TABLE IF NOT EXISTS "loan" (
    id UUID PRIMARY KEY,
    user_id UUID,
    principal_amount INT NOT NULL,
    interest_rate DECIMAL(5,2) NOT NULL,
    outstanding_amount INT NOT NULL,
    duration INT NOT NULL DEFAULT 50,
    duration_unit duration_unit,
    is_deliquent BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    last_due_at TIMESTAMPTZ DEFAULT (NOW() + interval '1 week')
);

CREATE TABLE IF NOT EXISTS "payment" (
    id UUID PRIMARY KEY,
    loan_id UUID REFERENCES loan(id),
    payment_number INT NOT NULL, -- represents the number of duration
    amount_paid INT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_payment_loan_unique ON "public"."payment" (
    loan_id,
    payment_number
);

-- dummy data for testing
INSERT INTO "loan" (id, user_id, principal_amount, interest_rate, outstanding_amount, duration, duration_unit, created_at)
VALUES
    ('1b68dc3f-345c-4296-8b52-9fe56e9934cb', '1bd0d553-e770-4352-90ba-0d5cdda990bf', 10000, 5.00, 10500, 52, 'WEEK', NOW()),
    ('af21f244-d622-4356-b2fc-479011b49346', '1bd0d553-e770-4352-90ba-0d5cdda990bf', 10000, 5.00, 10500, 52, 'WEEK', NOW());

COMMIT;