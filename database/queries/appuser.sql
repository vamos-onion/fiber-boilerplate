-- name: GetAllAppusers :many
SELECT *
FROM appuser;

-- name: GetAppusersByName :one
SELECT *
FROM appuser
WHERE name = $1;

-- name: SearchAppusers :many
SELECT *
FROM appuser
WHERE (@uuid::varchar IS NULL OR uuid = CAST(@uuid AS UUID))
  AND (@name::varchar IS NULL OR name = @name)
  AND (@gender::enum_gender IS NULL OR gender = @gender)
  AND (@withdraw::boolean IS NULL OR withdraw = @withdraw)
          ? 1 = @options::text;

-- name: CreateAppuser :one
INSERT INTO appuser (name, birthday, gender, withdraw)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpdateAppuser :one
UPDATE appuser
SET name     = $2,
    birthday = $3,
    gender   = $4,
    withdraw = $5
WHERE appuser.uuid = $1
RETURNING *;