CREATE TABLE IF NOT EXISTS tasks
             (
                          id SERIAL PRIMARY KEY,
                          task    VARCHAR ( 255 ) NOT NULL,
                          is_done BOOLEAN DEFAULT false,
                          created_at TIMESTAMP NOT NULL DEFAULT Now ( ),
                          updated_at TIMESTAMP NOT NULL DEFAULT Now ( ),
                          deleted_at timestamp DEFAULT null
             );