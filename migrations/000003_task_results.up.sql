CREATE TABLE task_results (
    id INT AUTO_INCREMENT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    task_id INT,
    status_code INT,
    headers JSON,
    length INT,
    FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE
);