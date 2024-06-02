CREATE TABLE account (
    id SERIAL PRIMARY KEY,
    username TEXT UNIQUE NOT NULL,
    email TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    date_created TIMESTAMPTZ NOT NULL,
    date_edited TIMESTAMPTZ
);

CREATE TABLE todo (
    id SERIAL PRIMARY KEY,
    title TEXT,
    completed BOOLEAN,
    account_id INT REFERENCES account(id)
);