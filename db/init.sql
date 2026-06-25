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

-- Bersihkan system categories jika script di-run ulang
DELETE FROM categories WHERE user_id IS NULL;

DO $$
DECLARE
    -- Income Parents
    id_occ_earn UUID := gen_random_uuid();
    id_fin_inv UUID := gen_random_uuid();
    id_misc_inc UUID := gen_random_uuid();
    
    -- Expense Parents
    id_food UUID := gen_random_uuid();
    id_clothing UUID := gen_random_uuid();
    id_housing UUID := gen_random_uuid();
    id_transport UUID := gen_random_uuid();
    id_comm UUID := gen_random_uuid();
    id_ent UUID := gen_random_uuid();
    id_edu UUID := gen_random_uuid();
    id_gifts UUID := gen_random_uuid();
    id_med UUID := gen_random_uuid();
    id_fin_ins UUID := gen_random_uuid();
    id_misc_exp UUID := gen_random_uuid();

    -- Transfer Parents
    id_gen_trans UUID := gen_random_uuid();
    id_loan UUID := gen_random_uuid();
    id_misc_trans UUID := gen_random_uuid();
BEGIN
    -- ================= INCOME =================
    INSERT INTO categories (id, user_id, name, type, icon, color_class, sort_order) VALUES
    (id_occ_earn, NULL, 'Occupational Earnings', 'income', '💼', 'text-[#D68E5A] bg-[#D68E5A]/10', 1),
    (id_fin_inv, NULL, 'Finance & Investment', 'income', '📈', 'text-[#D68E5A] bg-[#D68E5A]/10', 2),
    (id_misc_inc, NULL, 'Miscellaneous', 'income', '📝', 'text-[#D68E5A] bg-[#D68E5A]/10', 3);

    INSERT INTO categories (user_id, parent_id, name, type, icon, color_class, sort_order) VALUES
    (NULL, id_occ_earn, 'Salary Income', 'income', '💼', 'text-[#D68E5A] bg-[#D68E5A]/10', 1),
    (NULL, id_occ_earn, 'Bonus Income', 'income', '💼', 'text-[#D68E5A] bg-[#D68E5A]/10', 2),
    (NULL, id_occ_earn, 'Overtime Pay', 'income', '💼', 'text-[#D68E5A] bg-[#D68E5A]/10', 3),
    (NULL, id_occ_earn, 'Side Job Income', 'income', '💼', 'text-[#D68E5A] bg-[#D68E5A]/10', 4),
    (NULL, id_occ_earn, 'Commission Income', 'income', '💼', 'text-[#D68E5A] bg-[#D68E5A]/10', 5),

    (NULL, id_fin_inv, 'Interest Income', 'income', '📈', 'text-[#D68E5A] bg-[#D68E5A]/10', 1),
    (NULL, id_fin_inv, 'Dividend Income', 'income', '📈', 'text-[#D68E5A] bg-[#D68E5A]/10', 2),
    (NULL, id_fin_inv, 'Capital Gains', 'income', '📈', 'text-[#D68E5A] bg-[#D68E5A]/10', 3),
    (NULL, id_fin_inv, 'Rental Income', 'income', '📈', 'text-[#D68E5A] bg-[#D68E5A]/10', 4),

    (NULL, id_misc_inc, 'Gift Received', 'income', '📝', 'text-[#D68E5A] bg-[#D68E5A]/10', 1),
    (NULL, id_misc_inc, 'Lottery/Prize', 'income', '📝', 'text-[#D68E5A] bg-[#D68E5A]/10', 2),
    (NULL, id_misc_inc, 'Refund', 'income', '📝', 'text-[#D68E5A] bg-[#D68E5A]/10', 3),
    (NULL, id_misc_inc, 'Other Income', 'income', '📝', 'text-[#D68E5A] bg-[#D68E5A]/10', 4);

    -- ================= EXPENSE =================
    INSERT INTO categories (id, user_id, name, type, icon, color_class, sort_order) VALUES
    (id_food, NULL, 'Food & Drink', 'expense', '🍔', 'text-orange-500 bg-orange-50', 1),
    (id_clothing, NULL, 'Clothing & Appearance', 'expense', '👔', 'text-purple-500 bg-purple-50', 2),
    (id_housing, NULL, 'Housing & Houseware', 'expense', '🏠', 'text-gray-600 bg-gray-100', 3),
    (id_transport, NULL, 'Transportation', 'expense', '🚗', 'text-teal-500 bg-teal-50', 4),
    (id_comm, NULL, 'Communication', 'expense', '📱', 'text-blue-500 bg-blue-50', 5),
    (id_ent, NULL, 'Entertainment', 'expense', '🎬', 'text-rose-500 bg-rose-50', 6),
    (id_edu, NULL, 'Education & Studying', 'expense', '📚', 'text-lime-500 bg-lime-50', 7),
    (id_gifts, NULL, 'Gifts & Donations', 'expense', '🎁', 'text-green-500 bg-green-50', 8),
    (id_med, NULL, 'Medical & Healthcare', 'expense', '🏥', 'text-red-500 bg-red-50', 9),
    (id_fin_ins, NULL, 'Finance & Insurance', 'expense', '💳', 'text-amber-500 bg-amber-50', 10),
    (id_misc_exp, NULL, 'Miscellaneous', 'expense', '📌', 'text-gray-500 bg-gray-50', 11);

    INSERT INTO categories (user_id, parent_id, name, type, icon, color_class, sort_order) VALUES
    (NULL, id_food, 'Restaurant/Dining', 'expense', '🍔', 'text-orange-500 bg-orange-50', 1),
    (NULL, id_food, 'Groceries', 'expense', '🍔', 'text-orange-500 bg-orange-50', 2),
    (NULL, id_food, 'Coffee/Tea', 'expense', '🍔', 'text-orange-500 bg-orange-50', 3),
    (NULL, id_food, 'Beverages', 'expense', '🍔', 'text-orange-500 bg-orange-50', 4),
    (NULL, id_food, 'Snacks', 'expense', '🍔', 'text-orange-500 bg-orange-50', 5),

    (NULL, id_clothing, 'Clothing', 'expense', '👔', 'text-purple-500 bg-purple-50', 1),
    (NULL, id_clothing, 'Shoes', 'expense', '👔', 'text-purple-500 bg-purple-50', 2),
    (NULL, id_clothing, 'Accessories', 'expense', '👔', 'text-purple-500 bg-purple-50', 3),
    (NULL, id_clothing, 'Laundry', 'expense', '👔', 'text-purple-500 bg-purple-50', 4),
    (NULL, id_clothing, 'Cosmetics', 'expense', '👔', 'text-purple-500 bg-purple-50', 5),

    (NULL, id_housing, 'Rent', 'expense', '🏠', 'text-gray-600 bg-gray-100', 1),
    (NULL, id_housing, 'Mortgage', 'expense', '🏠', 'text-gray-600 bg-gray-100', 2),
    (NULL, id_housing, 'Utilities', 'expense', '🏠', 'text-gray-600 bg-gray-100', 3),
    (NULL, id_housing, 'Furniture', 'expense', '🏠', 'text-gray-600 bg-gray-100', 4),
    (NULL, id_housing, 'Home Appliances', 'expense', '🏠', 'text-gray-600 bg-gray-100', 5),
    (NULL, id_housing, 'Repairs', 'expense', '🏠', 'text-gray-600 bg-gray-100', 6),

    (NULL, id_transport, 'Fuel', 'expense', '🚗', 'text-teal-500 bg-teal-50', 1),
    (NULL, id_transport, 'Public Transit', 'expense', '🚗', 'text-teal-500 bg-teal-50', 2),
    (NULL, id_transport, 'Taxi/Ride-hailing', 'expense', '🚗', 'text-teal-500 bg-teal-50', 3),
    (NULL, id_transport, 'Vehicle Maintenance', 'expense', '🚗', 'text-teal-500 bg-teal-50', 4),
    (NULL, id_transport, 'Parking/Toll', 'expense', '🚗', 'text-teal-500 bg-teal-50', 5),

    (NULL, id_comm, 'Mobile Plan', 'expense', '📱', 'text-blue-500 bg-blue-50', 1),
    (NULL, id_comm, 'Internet', 'expense', '📱', 'text-blue-500 bg-blue-50', 2),
    (NULL, id_comm, 'Cable TV', 'expense', '📱', 'text-blue-500 bg-blue-50', 3),
    (NULL, id_comm, 'Postage', 'expense', '📱', 'text-blue-500 bg-blue-50', 4),

    (NULL, id_ent, 'Movies/Theater', 'expense', '🎬', 'text-rose-500 bg-rose-50', 1),
    (NULL, id_ent, 'Music', 'expense', '🎬', 'text-rose-500 bg-rose-50', 2),
    (NULL, id_ent, 'Games', 'expense', '🎬', 'text-rose-500 bg-rose-50', 3),
    (NULL, id_ent, 'Sports', 'expense', '🎬', 'text-rose-500 bg-rose-50', 4),
    (NULL, id_ent, 'Travel/Vacation', 'expense', '🎬', 'text-rose-500 bg-rose-50', 5),
    (NULL, id_ent, 'Hobbies', 'expense', '🎬', 'text-rose-500 bg-rose-50', 6),

    (NULL, id_edu, 'Tuition', 'expense', '📚', 'text-lime-500 bg-lime-50', 1),
    (NULL, id_edu, 'Books', 'expense', '📚', 'text-lime-500 bg-lime-50', 2),
    (NULL, id_edu, 'Online Courses', 'expense', '📚', 'text-lime-500 bg-lime-50', 3),
    (NULL, id_edu, 'Stationery', 'expense', '📚', 'text-lime-500 bg-lime-50', 4),

    (NULL, id_gifts, 'Gifts', 'expense', '🎁', 'text-green-500 bg-green-50', 1),
    (NULL, id_gifts, 'Charity', 'expense', '🎁', 'text-green-500 bg-green-50', 2),
    (NULL, id_gifts, 'Church/Religious', 'expense', '🎁', 'text-green-500 bg-green-50', 3),

    (NULL, id_med, 'Doctor Visit', 'expense', '🏥', 'text-red-500 bg-red-50', 1),
    (NULL, id_med, 'Medicine', 'expense', '🏥', 'text-red-500 bg-red-50', 2),
    (NULL, id_med, 'Dental', 'expense', '🏥', 'text-red-500 bg-red-50', 3),
    (NULL, id_med, 'Vision', 'expense', '🏥', 'text-red-500 bg-red-50', 4),
    (NULL, id_med, 'Insurance Premium', 'expense', '🏥', 'text-red-500 bg-red-50', 5),

    (NULL, id_fin_ins, 'Bank Fees', 'expense', '💳', 'text-amber-500 bg-amber-50', 1),
    (NULL, id_fin_ins, 'Interest Payment', 'expense', '💳', 'text-amber-500 bg-amber-50', 2),
    (NULL, id_fin_ins, 'Investment Fees', 'expense', '💳', 'text-amber-500 bg-amber-50', 3),
    (NULL, id_fin_ins, 'Life Insurance', 'expense', '💳', 'text-amber-500 bg-amber-50', 4),

    (NULL, id_misc_exp, 'Other Expenses', 'expense', '📌', 'text-gray-500 bg-gray-50', 1);

    -- ================= TRANSFER =================
    INSERT INTO categories (id, user_id, name, type, icon, color_class, sort_order) VALUES
    (id_gen_trans, NULL, 'General Transfer', 'transfer', '🔄', 'text-orange-500 bg-orange-50', 1),
    (id_loan, NULL, 'Loan & Debt', 'transfer', '📄', 'text-amber-500 bg-amber-50', 2),
    (id_misc_trans, NULL, 'Miscellaneous', 'transfer', '📌', 'text-gray-500 bg-gray-50', 3);

    INSERT INTO categories (user_id, parent_id, name, type, icon, color_class, sort_order) VALUES
    (NULL, id_gen_trans, 'Bank Transfer', 'transfer', '🔄', 'text-orange-500 bg-orange-50', 1),
    (NULL, id_gen_trans, 'Savings Deposit', 'transfer', '🔄', 'text-orange-500 bg-orange-50', 2),
    (NULL, id_gen_trans, 'Savings Withdrawal', 'transfer', '🔄', 'text-orange-500 bg-orange-50', 3),

    (NULL, id_loan, 'Loan Payment', 'transfer', '📄', 'text-amber-500 bg-amber-50', 1),
    (NULL, id_loan, 'Credit Card Payment', 'transfer', '📄', 'text-amber-500 bg-amber-50', 2),
    (NULL, id_loan, 'Borrow Money', 'transfer', '📄', 'text-amber-500 bg-amber-50', 3),

    (NULL, id_misc_trans, 'Other Transfer', 'transfer', '📌', 'text-gray-500 bg-gray-50', 1);

END $$;

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
