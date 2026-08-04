-- sqlite3
PRAGMA foreign_keys = ON;

DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS todos;
DROP TABLE IF EXISTS statuses;
DROP TABLE IF EXISTS priorities;

CREATE TABLE IF NOT EXISTS users (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  username VARCHAR(32) UNIQUE NOT NULL,
  email VARCHAR(254) NOT NULL UNIQUE,
  password VARCHAR(64) NOT NULL
);

CREATE TABLE IF NOT EXISTS statuses (
  id TINYINT PRIMARY KEY,
  emoji VARCHAR(6) NOT NULL,
  name VARCHAR(15) UNIQUE NOT NULL
);

INSERT INTO statuses (id, emoji, name) VALUES
  (1, '📝', 'Planning'),
  (2, '🟢', 'Active'),
  (3, '🚧', 'In Progress'),
  (4, '❌', 'Cancelled'),
  (5, '🗃️', 'Archive'),
  (6, '🗑️', 'Trash');

CREATE TABLE IF NOT EXISTS priorities (
  id TINYINT PRIMARY KEY,
  emoji VARCHAR(6) NOT NULL,
  name VARCHAR(15) UNIQUE NOT NULL
);

INSERT INTO priorities (id, emoji, name) VALUES
  (1, '🔴', 'Important'),
  (2, '🟠', 'Highest'),
  (3, '🟡', 'High'),
  (4, '🟢', 'Medium'),
  (5, '🔵', 'Low'),
  (6, '⚫', 'Lowest');

CREATE TABLE IF NOT EXISTS todos (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  task VARCHAR(255) NOT NULL,
  done BOOLEAN DEFAULT 0,
  due TIMESTAMP,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

  user_id INTEGER NOT NULL,
  status_id TINYINT NOT NULL DEFAULT 1,
  priority_id TINYINT NOT NULL DEFAULT 4,

  CONSTRAINT fk_todos_user FOREIGN KEY (user_id)
    REFERENCES users(id)
    ON UPDATE CASCADE ON DELETE CASCADE,
  CONSTRAINT fk_todos_status FOREIGN KEY (status_id)
    REFERENCES statuses(id)
    ON UPDATE CASCADE ON DELETE RESTRICT,
  CONSTRAINT fk_todos_priority FOREIGN KEY (priority_id)
    REFERENCES priorities(id)
    ON UPDATE CASCADE ON DELETE RESTRICT
);

