CREATE TABLE IF NOT EXISTS modules (
    id SERIAL PRIMARY KEY,
    course_id INT NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    "order" INT NOT NULL, -- "order" is a keyword in SQL
    pdf_path VARCHAR(255),
    video_path VARCHAR(255),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
    CONSTRAINT fk_modules_course FOREIGN KEY (course_id) REFERENCES courses(id) ON DELETE CASCADE
);

-- Add unique constraint to course_id and order, so no two modules in the same course have the same order
CREATE UNIQUE INDEX IF NOT EXISTS idx_course_id_order ON modules (course_id, "order") WHERE deleted_at IS NULL;