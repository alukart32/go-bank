ALTER TABLE "transfers" ADD COLUMN "from_entry_id" bigint NOT NULL UNIQUE;
ALTER TABLE "transfers" ADD FOREIGN KEY ("from_entry_id") REFERENCES "entries" ("id");

ALTER TABLE "transfers" ADD COLUMN "to_entry_id" bigint NOT NULL UNIQUE;
ALTER TABLE "transfers" ADD FOREIGN KEY ("to_entry_id") REFERENCES "entries" ("id");
