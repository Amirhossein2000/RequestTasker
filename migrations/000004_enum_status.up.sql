ALTER TABLE
    task_statuses
MODIFY
    COLUMN status ENUM('done', 'in_process', 'error', 'new') NOT NULL;