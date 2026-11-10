INSERT INTO users (full_name, email, password_hash, role_id, is_active)
SELECT 'Employee', 'employee@example.com', 'REPLACE_WITH_BCRYPT_HASH', id, TRUE
FROM roles
WHERE code = 'employee';
