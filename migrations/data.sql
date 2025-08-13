-- =================================================================
-- Tabel Inti & Referensi
-- =================================================================

-- Tabel untuk menyimpan informasi pengguna (jika sistem multi-user)
-- 1. USERS
CREATE TABLE users (
                       id BIGSERIAL PRIMARY KEY,
                       name VARCHAR(100) NOT NULL,
                       email VARCHAR(255) NOT NULL UNIQUE,
                       username VARCHAR(50) UNIQUE,
                       phone VARCHAR(20),
                       password VARCHAR(255) NOT NULL,
                       status TEXT CHECK (status IN ('inactive','active','banned')) DEFAULT 'inactive',
                       email_verified_at TIMESTAMP NULL,
                       phone_verified_at TIMESTAMP NULL,
                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                       updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 2. USER TOKENS (Refresh Token / Device Login)
CREATE TABLE user_tokens (
                             id BIGSERIAL PRIMARY KEY,
                             user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                             refresh_token TEXT NOT NULL,
                             ip_address VARCHAR(45),
                             user_agent TEXT,
                             expires_at TIMESTAMP NOT NULL,
                             created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 3. ROLES
CREATE TABLE roles (
                       id BIGSERIAL PRIMARY KEY,
                       name VARCHAR(50) NOT NULL,
                       slug VARCHAR(50) NOT NULL UNIQUE,
                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 4. PERMISSIONS
CREATE TABLE permissions (
                             id BIGSERIAL PRIMARY KEY,
                             name VARCHAR(50) NOT NULL,
                             slug VARCHAR(50) NOT NULL UNIQUE,
                             created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 5. ROLE_USER (Many-to-Many)
CREATE TABLE role_user (
                           user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                           role_id BIGINT NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
                           PRIMARY KEY (user_id, role_id)
);

-- 6. PERMISSION_ROLE (Many-to-Many)
CREATE TABLE permission_role (
                                 role_id BIGINT NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
                                 permission_id BIGINT NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
                                 PRIMARY KEY (role_id, permission_id)
);

-- 7. USER VERIFICATIONS (OTP / Email Token)
CREATE TABLE user_verifications (
                                    id BIGSERIAL PRIMARY KEY,
                                    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                                    type TEXT CHECK (type IN ('email','phone','password_reset')) NOT NULL,
                                    token VARCHAR(255) NOT NULL,
                                    otp_code VARCHAR(6),
                                    expires_at TIMESTAMP NOT NULL,
                                    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


-- Tabel untuk menyimpan semua akun keuangan pengguna
-- Contoh: Rekening bank, e-wallet, dompet tunai, kartu kredit.
CREATE TABLE Akun (
    id_akun INT PRIMARY KEY AUTO_INCREMENT,
    id_pengguna INT NOT NULL,
    nama_akun VARCHAR(50) NOT NULL,
    tipe_akun ENUM('Tabungan', 'E-Wallet', 'Tunai', 'Kartu Kredit', 'Investasi') NOT NULL,
    saldo_awal DECIMAL(15, 2) DEFAULT 0.00,
    FOREIGN KEY (id_pengguna) REFERENCES Pengguna(id_pengguna) ON DELETE CASCADE
);

-- Tabel untuk mengelola kategori pemasukan dan pengeluaran secara dinamis
CREATE TABLE Kategori (
    id_kategori INT PRIMARY KEY AUTO_INCREMENT,
    id_pengguna INT NOT NULL,
    nama_kategori VARCHAR(50) NOT NULL,
    tipe_kategori ENUM('Pemasukan', 'Pengeluaran') NOT NULL,
    FOREIGN KEY (id_pengguna) REFERENCES Pengguna(id_pengguna) ON DELETE CASCADE
);

-- =================================================================
-- Tabel Transaksi (Pusat dari Semua Aktivitas)
-- =================================================================

-- Tabel tunggal untuk mencatat SEMUA pergerakan uang: pemasukan, pengeluaran, dan transfer.
CREATE TABLE Transaksi (
    id_transaksi INT PRIMARY KEY AUTO_INCREMENT,
    id_pengguna INT NOT NULL,
    deskripsi VARCHAR(255) NOT NULL,
    jumlah DECIMAL(15, 2) NOT NULL,
    tanggal DATE NOT NULL,
    
    -- Tipe transaksi untuk membedakan jenis pergerakan uang
    tipe_transaksi ENUM('Pemasukan', 'Pengeluaran', 'Transfer') NOT NULL,
    
    -- Foreign Keys
    id_akun INT NOT NULL,             -- Akun yang terlibat (sumber untuk pengeluaran/transfer, tujuan untuk pemasukan)
    id_akun_tujuan INT NULL,          -- Diisi HANYA jika tipe_transaksi adalah 'Transfer'
    id_kategori INT NULL,             -- Diisi HANYA jika tipe_transaksi adalah 'Pemasukan' atau 'Pengeluaran'
    
    FOREIGN KEY (id_pengguna) REFERENCES Pengguna(id_pengguna) ON DELETE CASCADE,
    FOREIGN KEY (id_akun) REFERENCES Akun(id_akun) ON DELETE CASCADE,
    FOREIGN KEY (id_akun_tujuan) REFERENCES Akun(id_akun) ON DELETE SET NULL,
    FOREIGN KEY (id_kategori) REFERENCES Kategori(id_kategori) ON DELETE SET NULL
);


-- =================================================================
-- Tabel Spesialisasi (Hutang, Piutang, Aset, Budget)
-- =================================================================

-- Tabel untuk mencatat hutang
CREATE TABLE Hutang (
    id_hutang INT PRIMARY KEY AUTO_INCREMENT,
    id_pengguna INT NOT NULL,
    pemberi_hutang VARCHAR(100) NOT NULL,
    deskripsi VARCHAR(255),
    jumlah_pokok DECIMAL(15, 2) NOT NULL,
    sisa_hutang DECIMAL(15, 2) NOT NULL,
    tanggal_hutang DATE NOT NULL,
    tanggal_jatuh_tempo DATE NULL,
    status ENUM('Belum Lunas', 'Lunas') DEFAULT 'Belum Lunas',
    FOREIGN KEY (id_pengguna) REFERENCES Pengguna(id_pengguna) ON DELETE CASCADE
);

-- Tabel untuk mencatat piutang (orang lain berhutang ke kita)
CREATE TABLE Piutang (
    id_piutang INT PRIMARY KEY AUTO_INCREMENT,
    id_pengguna INT NOT NULL,
    peminjam VARCHAR(100) NOT NULL,
    deskripsi VARCHAR(255),
    jumlah_pokok DECIMAL(15, 2) NOT NULL,
    sisa_piutang DECIMAL(15, 2) NOT NULL,
    tanggal_piutang DATE NOT NULL,
    tanggal_jatuh_tempo DATE NULL,
    status ENUM('Belum Dibayar', 'Lunas') DEFAULT 'Belum Dibayar',
    FOREIGN KEY (id_pengguna) REFERENCES Pengguna(id_pengguna) ON DELETE CASCADE
);

-- Tabel penghubung untuk melacak pembayaran hutang. Setiap pembayaran adalah sebuah transaksi.
CREATE TABLE PembayaranHutang (
    id_pembayaran_hutang INT PRIMARY KEY AUTO_INCREMENT,
    id_hutang INT NOT NULL,
    id_transaksi INT NOT NULL,
    jumlah_pembayaran DECIMAL(15, 2) NOT NULL,
    tanggal_pembayaran DATE NOT NULL,
    FOREIGN KEY (id_hutang) REFERENCES Hutang(id_hutang) ON DELETE CASCADE,
    FOREIGN KEY (id_transaksi) REFERENCES Transaksi(id_transaksi) ON DELETE CASCADE
);

-- Tabel untuk mencatat aset yang dimiliki
CREATE TABLE Aset (
    id_aset INT PRIMARY KEY AUTO_INCREMENT,
    id_pengguna INT NOT NULL,
    nama_aset VARCHAR(100) NOT NULL,
    kategori_aset VARCHAR(50), -- Contoh: Properti, Kendaraan, Saham
    nilai_perolehan DECIMAL(15, 2) NOT NULL,
    nilai_saat_ini DECIMAL(15, 2),
    tanggal_perolehan DATE NOT NULL,
    FOREIGN KEY (id_pengguna) REFERENCES Pengguna(id_pengguna) ON DELETE CASCADE
);

-- Tabel untuk rencana anggaran (budgeting)
CREATE TABLE RencanaBudget (
    id_budget INT PRIMARY KEY AUTO_INCREMENT,
    id_pengguna INT NOT NULL,
    id_kategori INT NOT NULL,
    periode_bulan INT NOT NULL, -- (1-12)
    periode_tahun INT NOT NULL,
    jumlah_anggaran DECIMAL(15, 2) NOT NULL,
    FOREIGN KEY (id_pengguna) REFERENCES Pengguna(id_pengguna) ON DELETE CASCADE,
    FOREIGN KEY (id_kategori) REFERENCES Kategori(id_kategori) ON DELETE CASCADE
);

-- Tabel untuk melacak progres dana darurat atau tujuan keuangan lainnya
CREATE TABLE TujuanKeuangan (
    id_tujuan INT PRIMARY KEY AUTO_INCREMENT,
    id_pengguna INT NOT NULL,
    nama_tujuan VARCHAR(100) NOT NULL, -- Contoh: 'Dana Darurat', 'DP Rumah'
    target_jumlah DECIMAL(15, 2) NOT NULL,
    target_tanggal DATE NULL,
    terkumpul DECIMAL(15, 2) DEFAULT 0.00,
    FOREIGN KEY (id_pengguna) REFERENCES Pengguna(id_pengguna) ON DELETE CASCADE
);