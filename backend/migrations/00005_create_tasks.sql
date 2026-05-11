-- +goose Up
CREATE TABLE tasks (
    id               UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id          UUID        NOT NULL REFERENCES users    (id) ON DELETE CASCADE,
    project_id       UUID                    REFERENCES projects (id) ON DELETE SET NULL,
    parent_task_id   UUID                    REFERENCES tasks    (id) ON DELETE CASCADE,
    title            TEXT        NOT NULL CHECK (char_length(title) BETWEEN 1 AND 500),
    description_md   TEXT        NOT NULL DEFAULT '',
    description_html TEXT        NOT NULL DEFAULT '',
    links            JSONB       NOT NULL DEFAULT '[]',
    priority         SMALLINT    NOT NULL DEFAULT 0 CHECK (priority BETWEEN 0 AND 4),
    due_at           TIMESTAMPTZ,
    completed_at     TIMESTAMPTZ,
    sort_order       INT         NOT NULL DEFAULT 0,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_tasks_user_id    ON tasks (user_id);
CREATE INDEX idx_tasks_project_id ON tasks (project_id);
CREATE INDEX idx_tasks_parent_id  ON tasks (parent_task_id);
CREATE INDEX idx_tasks_due_at     ON tasks (due_at) WHERE due_at IS NOT NULL;
CREATE INDEX idx_tasks_completed  ON tasks (completed_at) WHERE completed_at IS NOT NULL;

-- +goose Down
DROP TABLE IF EXISTS tasks;
