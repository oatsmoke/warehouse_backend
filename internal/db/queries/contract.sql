-- name: CreateContract :one
INSERT INTO contracts (number, address)
VALUES (@number, @address)
RETURNING *;

-- name: ReadContract :one
SELECT id, number, address, deleted_at
FROM contracts
WHERE id = @id;

-- name: UpdateContract :execresult
UPDATE contracts
SET number  = @number,
    address = @address
WHERE id = @id
  AND (number != @number OR address != @address);

-- name: DeleteContract :execresult
UPDATE contracts
SET deleted_at = now()
WHERE id = @id
  AND deleted_at IS NULL;

-- name: RestoreContract :execresult
UPDATE contracts
SET deleted_at = NULL
WHERE id = @id
  AND deleted_at IS NOT NULL;

-- name: ListContract :many
SELECT id, number, address, deleted_at, count(*) OVER () AS total
FROM contracts
WHERE (@with_deleted::bool = true OR deleted_at IS NULL)
  AND (@search::text = '' OR (number || ' ' || address) ILIKE '%' || @search || '%')
  AND (array_length(@ids::bigint[], 1) IS NULL OR id = ANY (@ids))
ORDER BY CASE WHEN @sort_column::text = 'id' AND @sort_order::text = 'asc' THEN id::text END,
         CASE WHEN @sort_column = 'id' AND @sort_order = 'desc' THEN id::text END DESC,
         CASE WHEN @sort_column = 'number' AND @sort_order = 'asc' THEN number END,
         CASE WHEN @sort_column = 'number' AND @sort_order = 'desc' THEN number END DESC,
         CASE WHEN @sort_column = 'address' AND @sort_order = 'asc' THEN address END,
         CASE WHEN @sort_column = 'address' AND @sort_order = 'desc' THEN address END DESC
LIMIT @pagination_limit OFFSET @pagination_offset;