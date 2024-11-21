-- +migrate Up
-- +migrate StatementBegin

CREATE TABLE IF NOT EXISTS user_type (
    id INTEGER PRIMARY KEY,
    user_type VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS user_base(
    id INTEGER PRIMARY KEY,
    username VARCHAR(100) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    user_type INTEGER REFERENCES user_type(id),
    pin       VARCHAR(255) NOT NULL,
    status    INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS user_info (
    id INTEGER PRIMARY KEY,
    user_id INTEGER REFERENCES user_base(id),
    name VARCHAR(255),
    gender VARCHAR(10),
    birth_date DATE,
    address VARCHAR(255),
    occupation VARCHAR(255),
    job_place VARCHAR(255),
    email_address VARCHAR(255),
    phone_number VARCHAR(20),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(100),
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_by VARCHAR(100)
);

CREATE TABLE IF NOT EXISTS account_type (
    id INTEGER PRIMARY KEY,
    account_name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(100),
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_by VARCHAR(100)
);

CREATE TABLE IF NOT EXISTS account_base (
    account_number VARCHAR(50) PRIMARY KEY,
    user_id INTEGER REFERENCES user_base(id),
    account_type INTEGER REFERENCES account_type(id),
    bs_balance DECIMAL(20, 2),
    current_balance DECIMAL(20, 2),
    joint_with INTEGER REFERENCES user_base(id),
    status     INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS account_code_base (
    account_code VARCHAR(50) PRIMARY KEY,
    contra_code VARCHAR(50),
    code_type VARCHAR(100),
    code_description VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS ledger_base (
    id INTEGER PRIMARY KEY,
    reference_no VARCHAR(100),
    account_code VARCHAR(50) REFERENCES account_code_base(account_code),
    trx_amount DECIMAL(20, 2),
    global_id VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(reference_no, account_code)
);

CREATE TABLE IF NOT EXISTS ledger_transaction (
    id INTEGER PRIMARY KEY,
    account_number VARCHAR(50) REFERENCES account_base(account_number),
    before_balance DECIMAL(20, 2),
    trx_amt DECIMAL(20, 2),
    after_balance DECIMAL(20, 2),
    global_id VARCHAR(255),
    remark VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(100)
);

CREATE TABLE IF NOT EXISTS ledger_master (
    id INTEGER PRIMARY KEY,
    account_code VARCHAR(50) REFERENCES account_code_base(account_code),
    account_type VARCHAR(50),
    credit_amt   DECIMAL(20,2),
    debit_amt    DECIMAL(20,2),
    total_amt    DECIMAL(20,2),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(account_code, created_at)
);

CREATE TABLE IF NOT EXISTS accrual_base (
    id INTEGER PRIMARY KEY,
    account_number VARCHAR(50) REFERENCES account_base(account_number),
    total_history INTEGER,
    base_amount DECIMAL(20, 2),
    interest_rate DECIMAL(5, 2),
    total_accrual DECIMAL(20, 2),
    current_accrued DECIMAL(20, 2),
    from_date DATE,
    to_date DATE,
    status VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(100),
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_by VARCHAR(100)
);

CREATE TABLE IF NOT EXISTS accrual_history (
    id INTEGER PRIMARY KEY,
    reference_no INTEGER REFERENCES accrual_base(id),
    history_number INTEGER,
    accrual_amount DECIMAL(20, 2),
    total_accrued DECIMAL(20, 2),
    to_be_accrued DECIMAL(20, 2),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(100)
);

CREATE TABLE IF NOT EXISTS user_mac_addresses (
    id INTEGER PRIMARY KEY,
    user_id INTEGER REFERENCES user_base(id),
    mac_address VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO account_code_base (account_code, code_description, code_type)
VALUES 
    ('10908213', 'Office Fund', 'B'),
    ('24033123', 'Saving of Public Fund', 'B'),
    ('24033124', 'Time Deposit of Public Fund', 'B'),
    ('42231554', 'Penalty of Time Deposit', 'P'),
    ('51264478', 'Interest of Time Deposit', 'P')
ON CONFLICT (account_code) 
DO NOTHING;

INSERT INTO user_type (id, user_type, created_at, updated_at)
VALUES 
    ( 1,'Common Customer', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ( 2,'Gold Customer', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ( 3,'Priority Customer', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT (id) DO NOTHING;

INSERT INTO account_type (id,account_name, created_by, updated_by)
VALUES 
    (1,'Basic Saving Account', 'admin', 'admin'),
    (2,'Time Deposit Account', 'admin', 'admin')
ON CONFLICT (id) DO NOTHING;

INSERT INTO ledger_master (id,account_code, account_type, debit_amt, credit_amt, total_amt, created_at, updated_at)
VALUES 
    (1
        ,'10908213', 
        CASE 
            WHEN LEFT('10908213', 1) IN ('1', '2') THEN 'B'
            WHEN LEFT('10908213', 1) IN ('3', '4', '5') THEN 'P'
        END, 0, 0, 0,
        CURRENT_DATE, CURRENT_DATE),
    (2
        ,'24033123', 
        CASE 
            WHEN LEFT('24033123', 1) IN ('1', '2') THEN 'B'
            WHEN LEFT('24033123', 1) IN ('3', '4', '5') THEN 'P'
        END, 0, 0, 0,
        CURRENT_DATE, CURRENT_DATE),
    (3
        ,'24033124', 
        CASE 
            WHEN LEFT('24033124', 1) IN ('1', '2') THEN 'B'
            WHEN LEFT('24033124', 1) IN ('3', '4', '5') THEN 'P'
        END, 0, 0, 0,
        CURRENT_DATE, CURRENT_DATE),
    (4
        ,'51264478', 
        CASE 
            WHEN LEFT('51264478', 1) IN ('1', '2') THEN 'B'
            WHEN LEFT('51264478', 1) IN ('3', '4', '5') THEN 'P'
        END, 0, 0, 0,
        CURRENT_DATE, CURRENT_DATE),
    (5
        ,'42231554', 
        CASE 
            WHEN LEFT('42231554', 1) IN ('1', '2') THEN 'B'
            WHEN LEFT('42231554', 1) IN ('3', '4', '5') THEN 'P'
        END, 0, 0, 0,
        CURRENT_DATE, CURRENT_DATE)
ON CONFLICT (account_code,created_at) 
DO NOTHING;


-- +migrate StatementEnd
