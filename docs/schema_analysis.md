# Analisa Skema Database (DDL & DML) Berdasarkan Frontend (React)

Setelah melakukan analisa pada source code frontend di folder `finvera-app` (terutama di `components/desktop/DesktopProvider.jsx`, `lib/storage.js`, dan file halaman lainnya), berikut adalah struktur database (DDL) dan data awal (DML) yang dibutuhkan oleh backend (`finvera-be`).

---

## 1. DDL (Data Definition Language) - Struktur Tabel

Berdasarkan state yang disimpan di *localStorage* oleh frontend, kita membutuhkan tabel-tabel berikut di PostgreSQL:

### A. Tabel `users`
Dari `app/desktop/settings/user/page.jsx` dan fitur autentikasi.
- `id` (UUID/Int) - Primary Key
- `username` (Varchar)
- `email` (Varchar, Unique)
- `password` (Varchar) - Hashed
- `balance_visible` (Boolean) - Default: `true` (berdasarkan preferensi `BALANCE_VIS`)
- `created_at`, `updated_at`

### B. Tabel `accounts`
Dari `app/desktop/accounts/page.jsx`.
- `id` (UUID/Int) - Primary Key
- `user_id` (UUID/Int) - Foreign Key ke `users`
- `name` (Varchar) - Misal: "Cash", "Bank Account"
- `category` (Varchar) - Enum/String: 'Cash', 'Checking', 'Savings', 'Credit Card', 'Investment', 'Wallet', 'Other'
- `type` (Varchar) - Enum: 'asset', 'liability'
- `balance` (Decimal/Numeric) - Saldo akun (saat ini FE menggunakan float)
- `color` (Varchar) - Hex color code (misal: "#009E9E")
- `created_at`, `updated_at`, `deleted_at`

### C. Tabel `categories`
Dari `lib/presetCategories.js` dan manajemen kategori.
- `id` (UUID/Int) - Primary Key
- `user_id` (UUID/Int, Nullable) - Foreign Key ke `users` (atau `null` jika default system)
- `name` (Varchar)
- `type` (Varchar) - Enum: 'income', 'expense', 'transfer'
- `icon` (Varchar) - Emoji icon
- `color_class` (Varchar) - CSS class untuk warna
- `parent_id` (UUID/Int, Nullable) - Self-referencing FK untuk sub-kategori
- `created_at`, `updated_at`, `deleted_at`

### D. Tabel `tags`
- `id` (UUID/Int) - Primary Key
- `user_id` (UUID/Int) - Foreign Key ke `users`
- `name` (Varchar)
- `color` (Varchar) - Hex color code
- `created_at`, `updated_at`, `deleted_at`

### E. Tabel `transactions`
Dari `app/desktop/transactions/page.jsx`.
- `id` (UUID/Int) - Primary Key
- `user_id` (UUID/Int) - Foreign Key ke `users`
- `type` (Varchar) - Enum: 'income', 'expense', 'transfer'
- `amount` (Decimal/Numeric)
- `account_id` (UUID/Int) - Foreign Key ke `accounts`
- `target_account_id` (UUID/Int, Nullable) - FK ke `accounts` (hanya digunakan jika `type` = 'transfer')
- `category_id` (UUID/Int) - Foreign Key ke `categories`
- `date` (Timestamp) - Waktu transaksi dilakukan
- `note` (Varchar, max 200) - Opsional deskripsi
- `created_at`, `updated_at`, `deleted_at`

### F. Tabel `transaction_tags` (Pivot Table)
Karena satu transaksi bisa memiliki banyak tag (M:N).
- `transaction_id` (UUID/Int) - FK ke `transactions`
- `tag_id` (UUID/Int) - FK ke `tags`

### G. Tabel `scheduled_transactions`
Dari `app/desktop/scheduled/page.jsx`. Untuk transaksi otomatis berulang (Cron job).
- `id` (UUID/Int) - Primary Key
- `user_id` (UUID/Int) - FK ke `users`
- `name` (Varchar)
- `type` (Varchar) - Enum: 'income', 'expense', 'transfer'
- `amount` (Decimal/Numeric)
- `account_id` (UUID/Int) - FK ke `accounts`
- `target_account_id` (UUID/Int, Nullable) - FK ke `accounts`
- `category_id` (UUID/Int) - FK ke `categories`
- `note` (Varchar)
- `frequency` (Varchar) - Enum: 'daily', 'weekly', 'monthly', 'yearly'
- `next_run` (Timestamp) - Kapan transaksi ini harus di-execute selanjutnya
- `is_active` (Boolean) - Default: `true`
- `created_at`, `updated_at`, `deleted_at`

---

## 2. DML (Data Manipulation Language) - Initial Seeding

Ketika user baru saja register, FE mengharapkan ada data awal (default) agar tidak kosong. Kita perlu membuat logic "Seeder" di backend (atau men-trigger penciptaan data ini saat register user baru).

### A. Seeding Kategori Default (Berdasarkan `presetCategories.js`)
**Income:**
- Occupational Earnings (💼)
- Finance & Investment (📈)
- Miscellaneous (📝)

**Expense:**
- Food & Drink (🍔)
- Clothing & Appearance (👔)
- Housing & Houseware (🏠)
- Transportation (🚗)
- Communication (📱)
- Entertainment (🎬)
- Education & Studying (📚)
- Gifts & Donations (🎁)
- Medical & Healthcare (🏥)
- Finance & Insurance (💳)
- Miscellaneous (📌)

**Transfer:**
- General Transfer (🔄)
- Loan & Debt (📄)
- Miscellaneous (📌)

### B. Seeding Akun Default
Berdasarkan `DesktopProvider.jsx`, ketika user register, buat 3 akun awal:
1. **Cash** (Type: `asset`, Category: `Cash`, Color: `#009E9E`, Balance: 500.00 / 0)
2. **Bank Account** (Type: `asset`, Category: `Checking`, Color: `#E6923F`, Balance: 3500.00 / 0)
3. **Credit Card** (Type: `liability`, Category: `Credit Card`, Color: `#F14C4C`, Balance: -298.92 / 0)

### C. Seeding Tag Default
1. "Vacation 2026"
2. "Business Trip"
3. "Tax Deductible"

---

## 3. Catatan Khusus untuk Integrasi dengan Backend
1. **Database Transaction:** Di FE, kalkulasi saldo (*balance*) dihitung secara dinamis saat ada penambahan transaksi. Di Backend, jika ada mutasi/tambah transaksi baru (khususnya *transfer*), kita harus membungkusnya dalam **DB Transaction (`tx.Begin()`)** agar mutasi saldo pada `Account` ikut berubah secara konsisten, lalu baru commit/rollback.
2. **Cron Scheduler:** Untuk `scheduled_transactions`, di backend (Go) kita membutuhkan *Cron Job / Background Worker* (misal via library `robfig/cron`) yang berjalan periodik (contoh: setiap tengah malam). Worker ini mengecek field `next_run` <= `NOW()`. Jika ya, ia otomatis insert ke tabel `transactions`, meng-update `balance` akun, lalu set `next_run` maju sesuai `frequency`.
3. **Soft Delete vs Cascade Delete:** Banyak model memerlukan Soft Delete (`DeletedAt` gorm). Terutama *Category* dan *Account*, jika dihapus sebaiknya transaksi lama tetap valid atau Account tidak bisa dihapus kalau masih ada transaksinya (sesuai *business logic* ezBookkeeping).
