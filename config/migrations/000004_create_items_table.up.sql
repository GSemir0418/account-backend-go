BEGIN;
CREATE TYPE kind AS ENUM ('expenses', 'in_come');
CREATE TABLE IF NOT EXISTS items (
  id   SERIAL PRIMARY KEY,
  user_id SERIAL NOT NULL,
  amount INTEGER NOT NULL,
  tag_ids INTEGER[] NOT NULL,
  kind kind NOT NULL,
  happened_at TIMESTAMP NOT NULL DEFAULT now(),
  created_at TIMESTAMP NOT NULL DEFAULT now(),
  updated_at TIMESTAMP NOT NULL DEFAULT now()
);
COMMIT;