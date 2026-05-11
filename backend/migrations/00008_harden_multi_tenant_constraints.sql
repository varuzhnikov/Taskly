-- +goose Up
ALTER TABLE projects
    ADD CONSTRAINT projects_id_user_id_unique UNIQUE (id, user_id);

ALTER TABLE labels
    ADD CONSTRAINT labels_id_user_id_unique UNIQUE (id, user_id);

ALTER TABLE tasks
    ADD CONSTRAINT tasks_id_user_id_unique UNIQUE (id, user_id),
    ADD CONSTRAINT tasks_parent_task_owner_fk
        FOREIGN KEY (parent_task_id, user_id)
        REFERENCES tasks (id, user_id)
        ON DELETE CASCADE;

ALTER TABLE task_labels
    ADD COLUMN user_id UUID;

UPDATE task_labels tl
SET user_id = t.user_id
FROM tasks t
WHERE t.id = tl.task_id;

ALTER TABLE task_labels
    ALTER COLUMN user_id SET NOT NULL,
    ADD CONSTRAINT task_labels_user_fk
        FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    ADD CONSTRAINT task_labels_task_owner_fk
        FOREIGN KEY (task_id, user_id)
        REFERENCES tasks (id, user_id)
        ON DELETE CASCADE,
    ADD CONSTRAINT task_labels_label_owner_fk
        FOREIGN KEY (label_id, user_id)
        REFERENCES labels (id, user_id)
        ON DELETE CASCADE;

ALTER TABLE comments
    ADD CONSTRAINT comments_task_owner_fk
        FOREIGN KEY (task_id, user_id)
        REFERENCES tasks (id, user_id)
        ON DELETE CASCADE;

-- +goose Down
ALTER TABLE comments
    DROP CONSTRAINT IF EXISTS comments_task_owner_fk;

ALTER TABLE task_labels
    DROP CONSTRAINT IF EXISTS task_labels_label_owner_fk,
    DROP CONSTRAINT IF EXISTS task_labels_task_owner_fk,
    DROP CONSTRAINT IF EXISTS task_labels_user_fk,
    DROP COLUMN IF EXISTS user_id;

ALTER TABLE tasks
    DROP CONSTRAINT IF EXISTS tasks_parent_task_owner_fk,
    DROP CONSTRAINT IF EXISTS tasks_id_user_id_unique;

ALTER TABLE labels
    DROP CONSTRAINT IF EXISTS labels_id_user_id_unique;

ALTER TABLE projects
    DROP CONSTRAINT IF EXISTS projects_id_user_id_unique;
