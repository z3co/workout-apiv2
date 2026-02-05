package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	db "github.com/z3co/workout-apiv2/db/gen"
)

func RunDBContainer(ctx context.Context) (*postgres.PostgresContainer, string, error) {
	pgContainer, err := postgres.Run(ctx,
		"postgres:alpine",
		postgres.WithDatabase("test"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		postgres.WithInitScripts(filepath.Join("db", "sql", "workout.sql")),
		postgres.BasicWaitStrategies(),
	)
	if err != nil {
		return nil, "", err
	}
	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return nil, "", err
	}
	return pgContainer, connStr, err
}

func TestDBInsert(t *testing.T) {
	ctx := context.Background()
	pgContainer, connStr, err := RunDBContainer(ctx)
	if err != nil {
		t.Fatalf("could not start postgres container: %s", err)
	}
	t.Cleanup(func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			t.Fatalf("could not stop postgres container: %s", err)
		}
	})
	testCases := []db.InsertExerciseParams{
		{Name: "chest press", Description: pgtype.Text{String: "exercies", Valid: true}, Sets: int32(3), Href: pgtype.Text{String: "hello", Valid: true}, Equipment: pgtype.Text{String: "barbell", Valid: true}, TargetMuscle: "chest", Replacement: pgtype.Int8{Int64: 1, Valid: true}},
		{Name: "bench press", Description: pgtype.Text{String: "exercies", Valid: true}, Sets: int32(3), Href: pgtype.Text{String: "hello", Valid: true}, Equipment: pgtype.Text{String: "barbell", Valid: true}, TargetMuscle: "chest", Replacement: pgtype.Int8{Int64: 1, Valid: true}},
		{Name: "leg press", Description: pgtype.Text{String: "exercies", Valid: true}, Sets: int32(3), Href: pgtype.Text{String: "hello", Valid: true}, Equipment: pgtype.Text{String: "machine", Valid: true}, TargetMuscle: "legs", Replacement: pgtype.Int8{Int64: 1, Valid: true}},
	}
	testCasesSets := []db.InsertSetParams{
		{Reps: int32(10)},
		{Reps: int32(8)},
		{Reps: int32(6)},
	}
	t.Run("proper insert and get", func(t *testing.T) {
		conn, err := pgx.Connect(ctx, connStr)
		assert.NoError(t, err)
		store := db.New(conn)
		var returnIds []int64
		for _, exercise := range testCases {
			id, err := store.InsertExercise(ctx, exercise)
			assert.NoError(t, err)
			assert.NotNil(t, id)
			returnIds = append(returnIds, id)
			for _, set := range testCasesSets {
				set.ExerciseID = id
				err := store.InsertSet(ctx, set)
				assert.NoError(t, err)
			}
		}
		var returnCases []db.GetExerciseRow
		expectedReturnCases := []db.GetExerciseRow{
			{Name: "chest press", Sets: int32(3), Description: pgtype.Text{String: "exercies", Valid: true}, Href: pgtype.Text{String: "hello", Valid: true}, Equipment: pgtype.Text{String: "barbell", Valid: true}, TargetMuscle: "chest", Replacement: pgtype.Int8{Int64: int64(1), Valid: true}},
			{Name: "bench press", Sets: int32(3), Description: pgtype.Text{String: "exercies", Valid: true}, Href: pgtype.Text{String: "hello", Valid: true}, Equipment: pgtype.Text{String: "barbell", Valid: true}, TargetMuscle: "chest", Replacement: pgtype.Int8{Int64: int64(1), Valid: true}},
			{Name: "leg press", Sets: int32(3), Description: pgtype.Text{String: "exercies", Valid: true}, Href: pgtype.Text{String: "hello", Valid: true}, Equipment: pgtype.Text{String: "machine", Valid: true}, TargetMuscle: "legs", Replacement: pgtype.Int8{Int64: int64(1), Valid: true}},
		}
		for _, id := range returnIds {
			exercise, err := store.GetExercise(ctx, id)
			assert.NoError(t, err)
			returnCases = append(returnCases, exercise)
		}
		assert.Equal(t, expectedReturnCases, returnCases)

	})
}

func toPgText(str string) pgtype.Text {
	pgText := pgtype.Text{
		String: str,
		Valid:  true,
	}
	return pgText
}

func TestRunApi(t *testing.T) {
	ctx := context.Background()
	pgContainer, connStr, err := RunDBContainer(ctx)
	if err != nil {
		t.Fatalf("could not run db: %s", err)
	}
	t.Cleanup(func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			t.Fatalf("could not stop db: %s", err)
		}
	})

	type testCase struct {
		// Inputs
		db.InsertExerciseParams
		Reps []int `json:"-"`
		// Expected values
		code int
	}
	testSets := []int{
		10, 8, 6,
	}

	testCases := []testCase{
		{InsertExerciseParams: db.InsertExerciseParams{Name: "chest press", Description: toPgText("press with your chest"), Href: toPgText("chest.press"),
			Equipment: toPgText("barbell"), TargetMuscle: "chest", Replacement: pgtype.Int8{Int64: int64(1), Valid: true}}, code: 200, Reps: testSets},
		{InsertExerciseParams: db.InsertExerciseParams{Name: "squat", Description: toPgText("squat"), Href: toPgText("chest.press"), Equipment: toPgText("barbell"), TargetMuscle: "legs", Replacement: pgtype.Int8{Int64: int64(1), Valid: true}}, code: 200, Reps: testSets},
		{InsertExerciseParams: db.InsertExerciseParams{Name: "skull chrushers", Description: toPgText("crush your skull"), Href: toPgText("chest.press"), Equipment: toPgText("ezbar"), TargetMuscle: "triceps", Replacement: pgtype.Int8{Int64: int64(1), Valid: true}}, code: 200, Reps: testSets},
	}
	conn, err := pgx.Connect(ctx, connStr)
	if err != nil {
		t.Fatalf("error connect to db: %s", err)
	}
	store := db.New(conn)
	var returnIds []int64
	for _, exercise := range testCases {
		id, err := store.InsertExercise(ctx, exercise.InsertExerciseParams)
		if err != nil {
			t.Fatalf("could not insert exercises: %s", err)
		}
		returnIds = append(returnIds, id)
	}
	for _, id := range returnIds {
		for _, reps := range testSets {
			err := store.InsertSet(ctx, db.InsertSetParams{
				Reps:       int32(reps),
				ExerciseID: id,
			})
			if err != nil {
				t.Fatalf("could not insert into sets: %s", err)
			}
		}
	}
	t.Run("test get exercises", func(t *testing.T) {
		router := RunApi(connStr)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/exercises", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		var response []testCase
		err := json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Fatalf("Response body could not be parsed: %v\n%v", err, w.Body)
		}
		var returnCases []testCase
		for _, m := range testCases {
			m.Reps = nil
			m.code = 0
			returnCases = append(returnCases, m)
		}
		assert.Equal(t, returnCases, response)
	})
}
