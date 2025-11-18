-- name: CreateUser :one
INSERT INTO users (username, password_hash, email, role, employee_id)
VALUES (@username, @password_hash, @email, @role, @employee_id)
RETURNING *;

-- name: ReadUser :one
SELECT u.id,
       u.username,
       u.email,
       u.role,
       u.enabled,
       u.last_login_at,
       e.id          as employee_id,
       e.last_name   as employee_last_name,
       e.first_name  as employee_first_name,
       e.middle_name as employee_middle_name,
       e.phone       as employee_phone,
       d.id          as department_id,
       d.title       as department_title
FROM users u
         LEFT JOIN employees e ON e.id = u.employee_id
         LEFT JOIN departments d on d.id = e.department_id
WHERE u.id = @id;

-- name: UpdateUser :execresult
UPDATE users
SET username = @username,
    email    = @email
WHERE id = @id
  AND (username != @username OR email != @email);

-- name: DeleteUser :execresult
DELETE
FROM users
WHERE id = @id;

-- name: ListUser :many
SELECT u.id,
       u.username,
       u.email,
       u.role,
       u.enabled,
       u.last_login_at,
       e.id          as employee_id,
       e.last_name   as employee_last_name,
       e.first_name  as employee_first_name,
       e.middle_name as employee_middle_name,
       e.phone       as employee_phone,
       d.id          as department_id,
       d.title       as department_title
FROM users u
         LEFT JOIN employees e ON e.id = u.employee_id
         LEFT JOIN departments d on d.id = e.department_id
ORDER BY u.id;

-- name: GetPasswordHashUser :one
SELECT password_hash
FROM users
WHERE id = @id;

-- name: SetPasswordHashUser :execresult
UPDATE users
SET password_hash = @password_hash
WHERE id = @id;

-- name: SetRoleUser :execresult
UPDATE users
SET role = @role
WHERE id = @id;

-- name: SetEnabledUser :execresult
UPDATE users
SET enabled = @enabled
WHERE id = @id;

-- name: SetLastLoginAtUser :execresult
UPDATE users
SET last_login_at = now()
WHERE id = @id;

-- name: SetEmployeeUser :execresult
UPDATE users
SET employee_id = @employee_id
WHERE id = @id;

-- name: GetByUsernameUser :one
SELECT id, username, password_hash, email, role, enabled, last_login_at
FROM users
WHERE username = @id;