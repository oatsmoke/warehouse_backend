-- name: CreateEquipment :one
INSERT INTO equipments (serial_number, profile_id, company_id)
VALUES (@serial_number, @profile_id, @company_id)
RETURNING *;

-- name: ReadEquipment :one
SELECT e.id,
       e.serial_number,
       e.deleted_at,
       co.id    as company_id,
       co.title as company_title,
       p.id     as profile_id,
       p.title  as profile_title,
       ca.id    as category_id,
       ca.title as category_title
FROM equipments e
         INNER JOIN companies co on co.id = e.company_id
         INNER JOIN profiles p ON p.id = e.profile_id
         INNER JOIN categories ca ON ca.id = p.category_id
WHERE e.id = @id;

-- name: UpdateEquipment :execresult
UPDATE equipments
SET company_id    = @company_id,
    profile_id    = @profile_id,
    serial_number = @serial_number

WHERE id = @id
  AND (company_id != @company_id OR
       profile_id != @profile_id OR
       serial_number != @serial_number);

-- name: DeleteEquipment :execresult
UPDATE equipments
SET deleted_at = now()
WHERE id = @id
  AND deleted_at IS NULL;

-- name: RestoreEquipment :execresult
UPDATE equipments
SET deleted_at = NULL
WHERE id = @id
  AND deleted_at IS NOT NULL;

-- name: ListEquipment :many
SELECT e.id,
       e.serial_number,
       e.deleted_at,
       p.id             as profile_id,
       p.title          as profile_title,
       c.id             as category_id,
       c.title          as category_title,
       COUNT(*) OVER () AS total
FROM equipments e
         INNER JOIN profiles p ON p.id = e.profile_id
         INNER JOIN categories c ON c.id = p.category_id
WHERE (@with_deleted::bool = true OR e.deleted_at IS NULL)
  AND (@search::text = '' OR (e.serial_number || ' ' || p.title || ' ' || c.title) ILIKE '%' || @search || '%')
  AND (array_length(@ids::bigint[], 1) IS NULL OR e.id = ANY (@ids))
ORDER BY CASE WHEN @sort_column::text = 'id' AND @sort_order::text = 'asc' THEN e.id::text END,
         CASE WHEN @sort_column = 'id' AND @sort_order = 'desc' THEN e.id::text END DESC,
         CASE WHEN @sort_column = 'serial_number' AND @sort_order = 'asc' THEN e.serial_number END,
         CASE WHEN @sort_column = 'serial_number' AND @sort_order = 'desc' THEN e.serial_number END DESC,
         CASE WHEN @sort_column = 'profile_title' AND @sort_order = 'asc' THEN p.title END,
         CASE WHEN @sort_column = 'profile_title' AND @sort_order = 'desc' THEN p.title END DESC,
         CASE WHEN @sort_column = 'category_title' AND @sort_order = 'asc' THEN c.title END,
         CASE WHEN @sort_column = 'category_title' AND @sort_order = 'desc' THEN c.title END DESC
LIMIT @pagination_limit OFFSET @pagination_offset;