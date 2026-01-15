-- name: AddToStorage :execresult
INSERT INTO locations (equipment_id,
                       user_id,
                       move_at,
                       move_code)
VALUES (@equipment_id,
        @user_id,
        @move_at,
        @move_code);

-- name: MoveToLocation :execresult
INSERT INTO locations (equipment_id,
                       user_id,
                       move_at,
                       move_code,
                       from_department_id,
                       from_employee_id,
                       from_contract_id,
                       to_department_id,
                       to_employee_id,
                       to_contract_id)
VALUES (@equipment_id,
        @user_id,
        @move_at,
        @move_code,
        @from_department_id,
        @from_employee_id,
        @from_contract_id,
        @to_department_id,
        @to_employee_id,
        @to_contract_id);

-- name: ListEquipmentFromLocation :many
select e.id,
       e.serial_number,
       e.company_title,
       e.profile_title,
       e.category_title,
       e.total
from (SELECT DISTINCT ON (l.equipment_id) eq.id,
                                          eq.serial_number,
                                          co.title         AS company_title,
                                          p.title          AS profile_title,
                                          ca.title         AS category_title,
                                          l.to_department_id,
                                          l.to_employee_id,
                                          l.to_contract_id,
                                          COUNT(*) OVER () AS total
      FROM locations l
               LEFT JOIN equipments eq ON eq.id = l.equipment_id
               LEFT JOIN companies co ON co.id = eq.company_id
               LEFT JOIN profiles p ON p.id = eq.profile_id
               LEFT JOIN categories ca ON ca.id = p.category_id
      WHERE eq.deleted_at IS NULL
      ORDER BY l.equipment_id, l.move_at DESC, l.id DESC) e
WHERE (
    @to_department_id::bigint = 0
        AND e.to_department_id IS NULL
        AND e.to_employee_id IS NULL
        AND e.to_contract_id IS NULL
    )
   OR (
    @to_department_id::bigint > 0
        AND e.to_department_id = @to_department_id
    )
ORDER BY e.profile_title, e.serial_number;