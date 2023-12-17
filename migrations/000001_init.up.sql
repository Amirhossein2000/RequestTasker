CREATE TABLE tasks (
    id INT AUTO_INCREMENT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    public_id CHAR(36) NOT NULL,
    url VARCHAR(255) NOT NULL,
    method VARCHAR(10) NOT NULL,
    headers JSON,
    body JSON
);

-- CREATE TABLE task_statuses (
--     id INT AUTO_INCREMENT PRIMARY KEY,
--     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
--     task_id INT,
--     status VARCHAR(255) NOT NULL,
--     FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE
-- );

-- CREATE TABLE task_results (
--     id INT AUTO_INCREMENT PRIMARY KEY,
--     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
--     task_id INT,
--     status_code INT,
--     headers JSON,
--     length INT,
--     FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE
-- );