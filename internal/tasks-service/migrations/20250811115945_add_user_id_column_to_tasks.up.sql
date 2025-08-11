-- ALTER TABLE tasks
--     ADD COLUMN user_id INTEGER REFERENCES users (id) ON DELETE CASCADE

ALTER TABLE tasks
    ADD COLUMN IF NOT EXISTS user_id INTEGER NOT NULL;

CREATE INDEX IF NOT EXISTS idx_tasks_user_id ON tasks (user_id);

-- Проверка на положительный идентификатор
ALTER TABLE tasks
    ADD CONSTRAINT tasks_user_id_positive CHECK (user_id > 0);
