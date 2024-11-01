CREATE TABLE IF NOT EXISTS accounts (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    date_created TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    date_edited TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS permissions (
    id SERIAL PRIMARY KEY,
    name TEXT
);

CREATE TABLE IF NOT EXISTS permissions_members (
    permissions_id INT REFERENCES permissions(id) ON DELETE CASCADE,
    account_id INT REFERENCES accounts(id) ON DELETE CASCADE,
    PRIMARY KEY (permissions_id, account_id)
);

CREATE TABLE IF NOT EXISTS lists (
    id SERIAL PRIMARY KEY,
    title TEXT,
    description TEXT,
    account_id INT NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
    parent_list_id INT REFERENCES lists(id) ON DELETE CASCADE,
    permissions_id INT REFERENCES permissions(id) ON DELETE CASCADE,
    date_created TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    date_edited TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS todos (
    id SERIAL PRIMARY KEY,
    body TEXT,
    list_id INT NOT NULL REFERENCES lists(id),
    completed BOOLEAN NOT NULL DEFAULT FALSE,
    account_id INT REFERENCES accounts(id) ON DELETE CASCADE,
    date_created TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    date_edited TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    permissions_id INT REFERENCES permissions(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS notes (
    id SERIAL PRIMARY KEY,
    title TEXT, 
    body TEXT,
    account_id INT NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
    list_id INT NOT NULL REFERENCES lists(id) ON DELETE CASCADE,
    date_created TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    date_edited TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    permissions_id INT REFERENCES permissions(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS refreshtokens (
    id SERIAL PRIMARY KEY,
    token VARCHAR(256)
);

CREATE INDEX IF NOT EXISTS idx_account_username ON accounts(username);
CREATE INDEX IF NOT EXISTS idx_account_email ON accounts(email);
CREATE INDEX IF NOT EXISTS idx_list_account ON lists(account_id);
CREATE INDEX IF NOT EXISTS idx_todo_list ON todos(list_id);
CREATE INDEX IF NOT EXISTS idx_note_list ON notes(list_id);
CREATE INDEX IF NOT EXISTS idx_refresh_token ON refreshtokens(token);

 