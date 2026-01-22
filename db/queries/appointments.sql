-- name: CreateAppointment :one
INSERT INTO appointments (id,  user_id, appointment_time, status, notes) VALUES (
    gen_random_uuid(),
    $1,
    $2,
    $3,
    $4
) RETURNING *;

-- name: GetAppointmentByID :one
SELECT * FROM appointments WHERE id = $1;

-- name: GetAppointmentsByUserID :many
SELECT * FROM appointments WHERE user_id = $1 ORDER BY appointment_time DESC LIMIT $2 OFFSET $3;

-- name: UpdateAppointmentStatus :one
UPDATE appointments
SET status = $2, updated_at = NOW()
WHERE id = $1
RETURNING *;