-- name: ListExercises :many
SELECT 
	name,
	description,
	href,
	equipment,
	target_muscle
FROM exercises;

-- name: GetExercise :one
SELECT 
	name,
	description,
	href,
	equipment,
	target_muscle
FROM exercises
WHERE id = $1 LIMIT 1;

-- name: GetExerciseByEquipment :many
SELECT 
	name,
	description,
	href,
	equipment,
	target_muscle
FROM exercises
WHERE equipment = $1;

-- name: GetExerciseByMuscle :many
SELECT 
	name,
	description,
	href,
	equipment,
	target_muscle
FROM exercises
WHERE target_muscle = $1;

-- name: GetReplaceMentExercise :one
SELECT 
	e.name,
	e.description,
	e.href,
	e.equipment,
	e.target_muscle
FROM exercises o JOIN exercises e ON o.replacement = e.id
WHERE o.id = $1 LIMIT 1;

-- name: GetSetsByExerciseId :many
SELECT 
	reps
FROM sets
WHERE exercise_id = $1;

-- name: InsertExercise :one
INSERT INTO exercises (
	name, description, sets, href, equipment, target_muscle, replacement
) VALUES (
	$1, $2, $3, $4, $5, $6, $7
) RETURNING id;

-- name: InsertSet :exec
INSERT INTO sets (
	exercise_id, reps
) VALUES (
	$1, $2
); 
