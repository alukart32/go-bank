CREATE TYPE "currency" AS ENUM (
  'RUB',
  'USD'
);

CREATE TABLE "accounts" (
  "id" uuid PRIMARY KEY,
  "owner" varchar NOT NULL,
  "balance" bigint NOT NULL DEFAULT 0,
  "currency" currency NOT NULL,
  "created_at" timestamptz DEFAULT (now())
);

CREATE TABLE "entries" (
  "id" bigserial PRIMARY KEY,
  "account_id" uuid NOT NULL,
  "amount" bigint NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "transfers" (
  "id" bigserial PRIMARY KEY,
  "from_account_id" uuid NOT NULL,
  "to_account_id" uuid NOT NULL,
  "amount" bigint NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE INDEX ON "accounts" ("owner");

CREATE INDEX ON "entries" ("account_id");

CREATE INDEX ON "transfers" ("from_account_id");

CREATE INDEX ON "transfers" ("to_account_id");

CREATE INDEX ON "transfers" ("from_account_id", "to_account_id");

COMMENT ON COLUMN "entries"."amount" IS 'can be negative or positive';

COMMENT ON COLUMN "transfers"."amount" IS 'must be positive';

ALTER TABLE "accounts" ADD CONSTRAINT valid_owner CHECK (char_length(owner) <= 70);

ALTER TABLE "accounts" ADD CONSTRAINT positive_balance CHECK (balance >= 0);

ALTER TABLE "entries" ADD FOREIGN KEY ("account_id") REFERENCES "accounts" ("id");

ALTER TABLE "transfers" ADD FOREIGN KEY ("from_account_id") REFERENCES "accounts" ("id");

ALTER TABLE "transfers" ADD FOREIGN KEY ("to_account_id") REFERENCES "accounts" ("id");
