-- если когда-то всё-таки был FK — снимем безопасно
ALTER TABLE IF EXISTS tasks DROP CONSTRAINT IF EXISTS tasks_user_id_fkey;

-- снимем наш CHECK, если он есть
ALTER TABLE IF EXISTS tasks DROP CONSTRAINT IF EXISTS tasks_user_id_positive;

-- снимем индекс, если он есть
DROP INDEX IF EXISTS idx_tasks_user_id;

-- и уберём колонку
ALTER TABLE IF EXISTS tasks DROP COLUMN IF EXISTS user_id;
