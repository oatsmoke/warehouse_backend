-- name: CreateProfile :one
INSERT INTO profiles (title, category_id)
VALUES (@title, @category_id)
RETURNING *;

-- name: ReadProfile :one
SELECT p.id,
       p.title,
       p.deleted_at,
       c.id    as category_id,
       c.title as category_title
FROM profiles p
         INNER JOIN categories c ON c.id = p.category_id
WHERE p.id = @id;

-- name: UpdateProfile :execresult
UPDATE profiles
SET title       = @title,
    category_id = @category_id
WHERE id = @id
  AND (title != @title OR category_id != @category_id);

-- name: DeleteProfile :execresult
UPDATE profiles
SET deleted_at = now()
WHERE id = @id
  AND deleted_at IS NULL;

-- name: RestoreProfile :execresult
UPDATE profiles
SET deleted_at = NULL
WHERE id = @id
  AND deleted_at IS NOT NULL;

-- name: ListProfile :many
SELECT p.id,
       p.title,
       p.deleted_at,
       c.id             as category_id,
       c.title          as category_title,
       COUNT(*) OVER () AS total
FROM profiles p
         INNER JOIN categories c ON c.id = p.category_id
WHERE (@with_deleted::bool = true OR p.deleted_at IS NULL)
  AND (@search::text = '' OR (p.title || ' ' || c.title) ILIKE '%' || @search || '%')
  AND (array_length(@ids::bigint[], 1) IS NULL OR p.id = ANY (@ids))
ORDER BY CASE WHEN @sort_column::text = 'id' AND @sort_order::text = 'asc' THEN p.id::text END,
         CASE WHEN @sort_column = 'id' AND @sort_order = 'desc' THEN p.id::text END DESC,
         CASE WHEN @sort_column = 'title' AND @sort_order = 'asc' THEN p.title END,
         CASE WHEN @sort_column = 'title' AND @sort_order = 'desc' THEN p.title END DESC,
         CASE WHEN @sort_column = 'category_title' AND @sort_order = 'asc' THEN c.title END,
         CASE WHEN @sort_column = 'category_title' AND @sort_order = 'desc' THEN c.title END DESC
LIMIT @pagination_limit OFFSET @pagination_offset;