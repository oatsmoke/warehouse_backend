-- name: CreateEmployee :one
INSERT INTO employees (last_name, first_name, middle_name, phone)
VALUES (@last_name, @first_name, @middle_name, @phone)
RETURNING *;

-- name: ReadEmployee :one
SELECT e.id,
       e.last_name,
       e.first_name,
       e.middle_name,
       e.phone,
       e.deleted_at,
       d.id    as department_id,
       d.title as department_title
FROM employees e
         LEFT JOIN departments d ON d.id = e.department_id
WHERE e.id = @id;

-- name: UpdateEmployee :execresult
UPDATE employees
SET last_name   = @last_name,
    first_name  = @first_name,
    middle_name = @middle_name,
    phone       = @phone
WHERE id = @id
  AND (last_name != @last_name OR first_name != @first_name OR middle_name != @middle_name OR phone != @phone);

-- name: DeleteEmployee :execresult
UPDATE employees
SET deleted_at = now()
WHERE id = @id
  AND deleted_at IS NULL;

-- name: RestoreEmployee :execresult
UPDATE employees
SET deleted_at = NULL
WHERE id = @id
  AND deleted_at IS NOT NULL;

-- name: ListEmployee :many
SELECT e.id,
       e.last_name,
       e.first_name,
       e.middle_name,
       e.phone,
       e.deleted_at,
       d.id             as department_id,
       d.title          as department_title,
       COUNT(*) OVER () AS total
FROM employees e
         LEFT JOIN public.departments d ON d.id = e.department_id
WHERE (@with_deleted::bool = true OR e.deleted_at IS NULL)
  AND (@search::text = '' OR
       (e.last_name || ' ' || e.first_name || ' ' || e.middle_name || ' ' || e.phone || ' ' || d.title) ILIKE
       '%' || @search || '%')
  AND (array_length(@ids::bigint[], 1) IS NULL OR e.id = ANY (@ids))
ORDER BY CASE WHEN @sort_column::text = 'id' AND @sort_order::text = 'asc' THEN e.id::text END,
         CASE WHEN @sort_column = 'id' AND @sort_order = 'desc' THEN e.id::text END DESC,
         CASE WHEN @sort_column = 'last_name' AND @sort_order = 'asc' THEN e.last_name END,
         CASE WHEN @sort_column = 'last_name' AND @sort_order = 'desc' THEN e.last_name END DESC,
         CASE WHEN @sort_column = 'first_name' AND @sort_order = 'asc' THEN e.first_name END,
         CASE WHEN @sort_column = 'first_name' AND @sort_order = 'desc' THEN e.first_name END DESC,
         CASE WHEN @sort_column = 'middle_name' AND @sort_order = 'asc' THEN e.middle_name END,
         CASE WHEN @sort_column = 'middle_name' AND @sort_order = 'desc' THEN e.middle_name END DESC,
         CASE WHEN @sort_column = 'phone' AND @sort_order = 'asc' THEN e.phone END,
         CASE WHEN @sort_column = 'phone' AND @sort_order = 'desc' THEN e.phone END DESC,
         CASE WHEN @sort_column = 'department_title' AND @sort_order = 'asc' THEN d.title END,
         CASE WHEN @sort_column = 'department_title' AND @sort_order = 'desc' THEN d.title END DESC
LIMIT @pagination_limit OFFSET @pagination_offset;

-- name: SetDepartmentEmployee :execresult
UPDATE employees
SET department_id = @department_id
WHERE id = @id;
