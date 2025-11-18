-- name: CreateCompany :one
INSERT INTO companies (title)
VALUES (@title)
RETURNING *;

-- name: ReadCompany :one
SELECT id, title, deleted_at
FROM companies
WHERE id = @id;

-- name: UpdateCompany :execresult
UPDATE companies
SET title = @title
WHERE id = @id
  AND title != @title;

-- name: DeleteCompany :execresult
UPDATE companies
SET deleted_at = now()
WHERE id = @id
  AND deleted_at IS NULL;

-- name: RestoreCompany :execresult
UPDATE companies
SET deleted_at = NULL
WHERE id = @id
  AND deleted_at IS NOT NULL;

-- name: ListCompany :many
SELECT id, title, deleted_at, count(*) OVER () AS total
FROM companies
WHERE (@with_deleted::bool = true OR deleted_at IS NULL)
  AND (@search::text = '' OR title ILIKE '%' || @search || '%')
  AND (array_length(@ids::bigint[], 1) IS NULL OR id = ANY (@ids))
ORDER BY CASE WHEN @sort_column::text = 'id' AND @sort_order::text = 'asc' THEN id::text END,
         CASE WHEN @sort_column = 'id' AND @sort_order = 'desc' THEN id::text END DESC,
         CASE WHEN @sort_column = 'title' AND @sort_order = 'asc' THEN title END,
         CASE WHEN @sort_column = 'title' AND @sort_order = 'desc' THEN title END DESC
LIMIT @pagination_limit OFFSET @pagination_offset;