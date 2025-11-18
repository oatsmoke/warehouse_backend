-- name: CreateDepartment :one
INSERT INTO departments (title)
VALUES (@title)
RETURNING *;

-- name: ReadDepartment :one
SELECT id, title, deleted_at
FROM departments
WHERE id = @id;

-- name: UpdateDepartment :execresult
UPDATE departments
SET title = @title
WHERE id = @id
  AND title != @title;

-- name: DeleteDepartment :execresult
UPDATE departments
SET deleted_at = now()
WHERE id = @id
  AND deleted_at IS NULL;

-- name: RestoreDepartment :execresult
UPDATE departments
SET deleted_at = NULL
WHERE id = @id
  AND deleted_at IS NOT NULL;

-- name: ListDepartment :many
SELECT id, title, deleted_at, count(*) OVER () AS total
FROM departments
WHERE (@with_deleted::bool = true OR deleted_at IS NULL)
  AND (@search::text = '' OR title ILIKE '%' || @search || '%')
  AND (array_length(@ids::bigint[], 1) IS NULL OR id = ANY (@ids))
ORDER BY CASE WHEN @sort_column::text = 'id' AND @sort_order::text = 'asc' THEN id::text END,
         CASE WHEN @sort_column = 'id' AND @sort_order = 'desc' THEN id::text END DESC,
         CASE WHEN @sort_column = 'title' AND @sort_order = 'asc' THEN title END,
         CASE WHEN @sort_column = 'title' AND @sort_order = 'desc' THEN title END DESC
LIMIT @pagination_limit OFFSET @pagination_offset;