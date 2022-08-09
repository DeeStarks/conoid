CREATE TABLE processes (
    pid TEXT NOT NULL PRIMARY KEY,
    name TEXT NOT NULL,
    status TEXT NOT NULL,
    type TEXT NOT NULL,
    listeners TEXT,
    root_directory TEXT,
    client_address TEXT,
    tunnelled INTEGER,
    created_at NUMERIC NOT NULL
)