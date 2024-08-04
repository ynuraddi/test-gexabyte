CREATE TABLE IF NOT EXISTS "currency" (
  "id" bigserial PRIMARY KEY,
  "symbol" varchar UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS "currency_price" (
  "id" bigserial PRIMARY KEY,
  "price" numeric(20,10) NOT NULL,
  "currency_id" bigint NOT NULL,
  "time" bigint NOT NULL,

  FOREIGN KEY(currency_id) REFERENCES currency(id) ON DELETE RESTRICT
);

INSERT INTO currency(symbol) VALUES('BTCUSDT');
INSERT INTO currency(symbol) VALUES('ETHUSDT');
INSERT INTO currency(symbol) VALUES('SOLUSDT');
INSERT INTO currency(symbol) VALUES('TRXUSDT');
