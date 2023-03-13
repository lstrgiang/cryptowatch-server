BEGIN;

CREATE TABLE IF NOT EXISTS users (
    id integer PRIMARY KEY,
    email text NOT NULL, 
    balance float DEFAULT 0,
    created_at timestamp default current_timestamp,
    updated_at timestamp default current_timestamp,
    deleted_at timestamp default null
);


CREATE TYPE order_type AS ENUM ('sell', 'buy');
CREATE TABLE IF NOT EXISTS orders (
    id uuid PRIMARY KEY,
    user_id integer NOT NULL constraint order_user_fk references users,
    usd_price float NOT NULL,
    amount float NOT NULL,
    type order_type NOT NULL,
    created_at timestamp default current_timestamp
);

CREATE INDEX user_order_type_identifier ON orders
    (
        user_id,
        type
    );

CREATE INDEX user_identifier ON orders
    (
        user_id
    );

COMMIT;
