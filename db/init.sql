-- ========================================================================================
-- FINVERA DDL & DML Script
-- Bisa langsung di-run di PostgreSQL (misal via pgAdmin, DBeaver, atau ekstensi VSCode)
-- Pastikan sudah connect ke database `finvera` sebelum menjalankan script ini.
-- ========================================================================================

-- Pastikan extention UUID terinstall
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- ========================================================================================
-- DDL (Data Definition Language)
-- ========================================================================================

-- 1. Table: users
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(100) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    default_currency VARCHAR(10) DEFAULT 'IDR',
    first_day_of_week INT DEFAULT 1,
    fiscal_year_start INT DEFAULT 1,
    theme VARCHAR(20) DEFAULT 'light',
    totp_secret VARCHAR(255),
    totp_enabled BOOLEAN DEFAULT FALSE,
    email_verified_at TIMESTAMP,
    balance_visible BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- 2. Table: accounts
CREATE TABLE IF NOT EXISTS accounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE,
    parent_id UUID REFERENCES accounts(id) ON DELETE SET NULL,
    name VARCHAR(100) NOT NULL,
    type VARCHAR(50) NOT NULL, -- 'asset', 'liability'
    category VARCHAR(50), -- 'Checking', 'Cash', 'Credit Card', dll
    currency VARCHAR(10) DEFAULT 'IDR',
    icon VARCHAR(100),
    color VARCHAR(20),
    balance NUMERIC(15,2) DEFAULT 0,
    initial_balance NUMERIC(15,2) DEFAULT 0,
    statement_day INT,
    is_hidden BOOLEAN DEFAULT FALSE,
    sort_order INT DEFAULT 0,
    note VARCHAR(500),
    last_reconciled_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- 3. Table: categories
CREATE TABLE IF NOT EXISTS categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE, -- Nullable untuk default system categories
    parent_id UUID REFERENCES categories(id) ON DELETE SET NULL,
    name VARCHAR(100) NOT NULL,
    type VARCHAR(50) NOT NULL, -- 'income', 'expense', 'transfer'
    icon VARCHAR(100),
    color VARCHAR(20),
    color_class VARCHAR(100),
    sort_order INT DEFAULT 0,
    note VARCHAR(500),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- 4. Table: tags
CREATE TABLE IF NOT EXISTS tags (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE,
    name VARCHAR(100) NOT NULL,
    color VARCHAR(20),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- 5. Table: transactions
CREATE TABLE IF NOT EXISTS transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE,
    type VARCHAR(50) NOT NULL, -- 'income', 'expense', 'transfer'
    amount NUMERIC(15,2) NOT NULL,
    account_id UUID NOT NULL REFERENCES accounts(id) ON DELETE RESTRICT,
    target_account_id UUID REFERENCES accounts(id) ON DELETE RESTRICT,
    category_id UUID NOT NULL REFERENCES categories(id) ON DELETE RESTRICT,
    date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    note VARCHAR(500),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- 6. Table: transaction_tags
CREATE TABLE IF NOT EXISTS transaction_tags (
    transaction_id UUID NOT NULL REFERENCES transactions(id) ON DELETE CASCADE,
    tag_id UUID NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
    PRIMARY KEY (transaction_id, tag_id)
);

-- 7. Table: scheduled_transactions
CREATE TABLE IF NOT EXISTS scheduled_transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE,
    name VARCHAR(100) NOT NULL,
    type VARCHAR(50) NOT NULL, -- 'income', 'expense', 'transfer'
    amount NUMERIC(15,2) NOT NULL,
    account_id UUID NOT NULL REFERENCES accounts(id) ON DELETE RESTRICT,
    target_account_id UUID REFERENCES accounts(id) ON DELETE RESTRICT,
    category_id UUID NOT NULL REFERENCES categories(id) ON DELETE RESTRICT,
    note VARCHAR(500),
    frequency VARCHAR(20) NOT NULL, -- 'daily', 'weekly', 'monthly', 'yearly'
    next_run TIMESTAMP NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);


-- ========================================================================================
-- DML (Data Manipulation Language) - SEEDING
-- ========================================================================================

-- Karena data kategori ini "default system", kita insert dengan user_id = NULL
-- (sehingga berlaku global/default untuk semua user yang baru mendaftar)
INSERT INTO categories (id, user_id, name, type, icon, color_class) VALUES
    -- INCOME
    (gen_random_uuid(), NULL, 'Occupational Earnings', 'income', '💼', 'text-[#D68E5A] bg-[#D68E5A]/10'),
    (gen_random_uuid(), NULL, 'Finance & Investment', 'income', '📈', 'text-[#D68E5A] bg-[#D68E5A]/10'),
    (gen_random_uuid(), NULL, 'Miscellaneous', 'income', '📝', 'text-[#D68E5A] bg-[#D68E5A]/10'),
    
    -- EXPENSE
    (gen_random_uuid(), NULL, 'Food & Drink', 'expense', '🍔', 'text-orange-500 bg-orange-50'),
    (gen_random_uuid(), NULL, 'Clothing & Appearance', 'expense', '👔', 'text-purple-500 bg-purple-50'),
    (gen_random_uuid(), NULL, 'Housing & Houseware', 'expense', '🏠', 'text-gray-600 bg-gray-100'),
    (gen_random_uuid(), NULL, 'Transportation', 'expense', '🚗', 'text-teal-500 bg-teal-50'),
    (gen_random_uuid(), NULL, 'Communication', 'expense', '📱', 'text-blue-500 bg-blue-50'),
    (gen_random_uuid(), NULL, 'Entertainment', 'expense', '🎬', 'text-rose-500 bg-rose-50'),
    (gen_random_uuid(), NULL, 'Education & Studying', 'expense', '📚', 'text-lime-500 bg-lime-50'),
    (gen_random_uuid(), NULL, 'Gifts & Donations', 'expense', '🎁', 'text-green-500 bg-green-50'),
    (gen_random_uuid(), NULL, 'Medical & Healthcare', 'expense', '🏥', 'text-red-500 bg-red-50'),
    (gen_random_uuid(), NULL, 'Finance & Insurance', 'expense', '💳', 'text-amber-500 bg-amber-50'),
    (gen_random_uuid(), NULL, 'Miscellaneous', 'expense', '📌', 'text-gray-500 bg-gray-50'),
    
    -- TRANSFER
    (gen_random_uuid(), NULL, 'General Transfer', 'transfer', '🔄', 'text-orange-500 bg-orange-50'),
    (gen_random_uuid(), NULL, 'Loan & Debt', 'transfer', '📄', 'text-amber-500 bg-amber-50'),
    (gen_random_uuid(), NULL, 'Miscellaneous', 'transfer', '📌', 'text-gray-500 bg-gray-50')
ON CONFLICT DO NOTHING;

-- Catatan:
-- Untuk DML (Seeding) Accounts dan Tags, mereka bergantung pada user_id.
-- Sehingga Seeding untuk Account dan Tag paling ideal dieksekusi via kode Go
-- di dalam service `RegisterUser`, contoh:
/*
INSERT INTO accounts (user_id, name, category, type, color, balance) VALUES
    ('id-user-baru', 'Cash', 'Cash', 'asset', '#009E9E', 0),
    ('id-user-baru', 'Bank Account', 'Checking', 'asset', '#E6923F', 0),
    ('id-user-baru', 'Credit Card', 'Credit Card', 'liability', '#F14C4C', 0);
*/
