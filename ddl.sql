CREATE TABLE consumer (
      id SERIAL PRIMARY KEY,
      account_limit INTEGER NOT NULL,
      balance INTEGER NOT null
);

CREATE TABLE transactions (
      id SERIAL PRIMARY KEY,
      consumer_id INTEGER NOT NULL,
      value INTEGER NOT NULL,
      type CHAR(1) NOT NULL,
      description VARCHAR(10) NOT NULL,
      created_at TIMESTAMP NOT NULL DEFAULT NOW(),
      CONSTRAINT fk_consumer_transactions_id FOREIGN KEY (consumer_id) REFERENCES consumer(id)
);

INSERT INTO public.consumer (id,account_limit,balance) VALUES (1,100000,0);
INSERT INTO public.consumer (id,account_limit,balance) VALUES (2,80000,0);
INSERT INTO public.consumer (id,account_limit,balance) VALUES (3,1000000,0);
INSERT INTO public.consumer (id,account_limit,balance) VALUES (4,10000000,0);
INSERT INTO public.consumer (id,account_limit,balance) VALUES (5,500000,0);
