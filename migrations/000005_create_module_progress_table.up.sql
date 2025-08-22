CREATE TABLE IF NOT EXISTS module_progresses (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    module_id INT NOT NULL,
    is_completed BOOLEAN DEFAULT FALSE NOT NULL,
    completed_at TIMESTAMPTZ, -- Nullable, set when is_completed is true
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
    CONSTRAINT fk_module_progresses_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_module_progresses_module FOREIGN KEY (module_id) REFERENCES modules(id) ON DELETE CASCADE,
    CONSTRAINT uq_user_module CHECK (deleted_at IS NOT NULL OR NOT EXISTS (
        SELECT 1 FROM module_progresses mp
        WHERE mp.user_id = user_id AND mp.module_id = module_id AND mp.deleted_at IS NULL
    ))
);