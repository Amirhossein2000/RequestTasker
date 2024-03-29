CREATE TABLE tasks (
    id INT AUTO_INCREMENT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    public_id CHAR(36) NOT NULL,
    url VARCHAR(255) NOT NULL,
    method VARCHAR(10) NOT NULL,
    headers JSON,
    body JSON
);