CREATE TABLE processes (
    pid TEXT NOT NULL PRIMARY KEY,
    name TEXT NOT NULL,
    status INTEGER NOT NULL,
    type TEXT NOT NULL,
    listeners TEXT,
    root_directory TEXT,
    remote_server TEXT,
    tunnelled INTEGER,
    created_at NUMERIC NOT NULL
)