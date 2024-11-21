-- +migrate Down
-- +migrate StatementBegin

DROP TABLE IF EXISTS accrual_history CASCADE;
DROP TABLE IF EXISTS accrual_base CASCADE;
DROP TABLE IF EXISTS ledger_master CASCADE;
DROP TABLE IF EXISTS ledger_transaction CASCADE;
DROP TABLE IF EXISTS ledger_base CASCADE;
DROP TABLE IF EXISTS account_code_base CASCADE;
DROP TABLE IF EXISTS account_base CASCADE;
DROP TABLE IF EXISTS account_type CASCADE;
DROP TABLE IF EXISTS user_info CASCADE;
DROP TABLE IF EXISTS user_base CASCADE;
DROP TABLE IF EXISTS user_type CASCADE;
DROP TABLE IF EXISTS user_mac_addresses CASCADE;


-- +migrate StatementEnd