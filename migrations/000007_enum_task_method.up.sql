ALTER TABLE
    tasks
MODIFY
    COLUMN method ENUM(
        'GET',
        'POST',
        'PUT',
        'DELETE',
        'PATCH',
        'HEAD',
        'OPTIONS'
    );