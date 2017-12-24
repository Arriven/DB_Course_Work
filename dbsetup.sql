CREATE TABLE IF NOT EXISTS phonebook("id" SERIAL PRIMARY KEY, "name" varchar(50), "phone" varchar(100));
DELETE FROM phonebook;
INSERT INTO phonebook VALUES (default, "SomeName", "0123456789");