CREATE TABLE accounts (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    date_created TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    date_edited TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TABLE permissions (
    id SERIAL PRIMARY KEY,
    name TEXT
);

CREATE TABLE permissions_members (
    permissions_id INT REFERENCES permissions(id) ON DELETE CASCADE,
    account_id INT REFERENCES accounts(id) ON DELETE CASCADE,
    PRIMARY KEY (permissions_id, account_id)
);

CREATE TABLE lists (
    id SERIAL PRIMARY KEY,
    title TEXT,
    description TEXT,
    account_id INT NOT NULL REFERENCES accounts(id), -- on delete cascade
    parent_list_id INT REFERENCES lists(id), -- on delete cascade
    permissions_id INT REFERENCES permissions(id), -- on delete cascade
    date_created TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    date_edited TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TABLE todos (
    id SERIAL PRIMARY KEY,
    body TEXT,
    list_id INT NOT NULL REFERENCES lists(id),
    completed BOOLEAN NOT NULL DEFAULT FALSE,
    account_id INT REFERENCES accounts(id), -- on delete cascade
    date_created TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    date_edited TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    permissions_id INT REFERENCES permissions(id)
);

CREATE TABLE notes (
    id SERIAL PRIMARY KEY,
    title TEXT, 
    body TEXT,
    account_id INT NOT NULL REFERENCES accounts(id), -- on delete cascade
    list_id INT NOT NULL REFERENCES lists(id), -- on delete cascade
    date_created TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    date_edited TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    permissions_id INT REFERENCES permissions(id) -- on delete cascade
);

CREATE INDEX idx_account_username ON accounts(username);
CREATE INDEX idx_account_email ON accounts(email);
CREATE INDEX idx_list_account ON lists(account_id);
CREATE INDEX idx_todo_list ON todos(list_id);
CREATE INDEX idx_note_list ON notes(list_id);

 