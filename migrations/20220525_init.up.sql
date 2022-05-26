-- ------------------------
-- Initial schema creation
-- ------------------------

CREATE TABLE IF NOT EXISTS `bank_accounts` (
    id INT NOT NULL AUTO_INCREMENT,
    organization_name TEXT NOT NULL,
    balance_cents INTEGER NOT NULL,
    iban VARCHAR(34) NOT NULL,
    bic TEXT NOT NULL,

    PRIMARY KEY(id),
    UNIQUE INDEX idx_iban (iban)
) ENGINE=InnoDB DEFAULT CHARACTER SET=utf8mb4;

CREATE TABLE IF NOT EXISTS `transactions` (
    id INT NOT NULL AUTO_INCREMENT,
    counterparty_name TEXT NOT NULL,
    counterparty_iban  VARCHAR(34) NOT NULL,
    counterparty_bic TEXT NOT NULL,
    amount_cents INTEGER NOT NULL,
    amount_currency TEXT NOT NULL,
    bank_account_id INTEGER NOT NULL,
    description TEXT,

    PRIMARY KEY(id),
    FOREIGN KEY (bank_account_id)
        REFERENCES bank_accounts (id)
) ENGINE=InnoDB DEFAULT CHARACTER SET=utf8mb4;

INSERT INTO `bank_accounts` (`organization_name`, `bic`, `iban`, `balance_cents`)
VALUES
    ( "ACME Corp", "OIVUSCLQXXX", "FR10474608000002006107XXXXX", 10000000);
