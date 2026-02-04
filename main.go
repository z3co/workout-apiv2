package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	db "github.com/z3co/workout-apiv2/db/gen"
)

type Server struct {
	store *db.Queries
}

func (server *Server) GetExercisesHandler(ctx *gin.Context) {
	exercises, err := server.store.ListExercises(context.Background())
	if err != nil {
		errString := fmt.Errorf("exercise not found: %s", err)
		ctx.Error(errString)
		ctx.String(http.StatusNotFound, errString.Error())
		return
	}

	ctx.JSON(200, exercises)
}

func (server *Server) GetExerciseByIdHandler(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		errString := fmt.Errorf("bad request %s", err)
		ctx.Error(errString)
		ctx.String(http.StatusBadRequest, errString.Error())
		return
	}
	exercise, err := server.store.GetExercise(ctx.Request.Context(), int64(id))
	if err != nil {
		errString := fmt.Errorf("exercise not found: %s", err)
		ctx.Error(errString)
		ctx.String(http.StatusNotFound, errString.Error())
		return
	}
	ctx.JSON(200, exercise)
}

func (server *Server) GetExerciseByEquipment(ctx *gin.Context) {
	name := ctx.Param("equipment")
	exercise, err := server.store.GetExerciseByEquipment(ctx.Request.Context(), pgtype.Text{
		String: name,
		Valid: true,
	})
	if err != nil {
		errString := fmt.Errorf("exercise not found: %s", err)
		ctx.Error(errString)
		ctx.String(http.StatusNotFound, errString.Error())
		return
	}
	ctx.JSON(200, exercise)
}

func (server *Server) GetExerciseByMuscle(ctx *gin.Context) {
	muscle := ctx.Param("muscle")
	exercise, err := server.store.GetExerciseByMuscle(ctx.Request.Context(), muscle)
	if err != nil {
		errString := fmt.Errorf("exercise not found: %s", err)
		ctx.Error(errString)
		ctx.String(http.StatusNotFound, errString.Error())
		return
	}
	fmt.Println(exercise)
	ctx.JSON(200, exercise)
}

func (server *Server) GetReplacmentByIdHandler(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		errString := fmt.Errorf("bad request %s", err)
		ctx.Error(errString)
		ctx.String(http.StatusBadRequest, errString.Error())
		return
	}
	exercise, err := server.store.GetReplaceMentExercise(ctx.Request.Context(), int64(id))
	if err != nil {
		errString := fmt.Errorf("exercise not found: %s", err)
		ctx.Error(errString)
		ctx.String(http.StatusNotFound, errString.Error())
		return 
	}
	ctx.JSON(200, exercise)
}

func (server *Server) GetSetsByExerciseId(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		errString := fmt.Errorf("bad request %s\n", err)
		ctx.Error(errString)
		ctx.String(http.StatusBadRequest, errString.Error())
		return
	}

	sets, err := server.store.GetSetsByExerciseId(ctx.Request.Context(), int64(id))
	if err != nil {
		errString := fmt.Errorf("exercise not found: %s", err)
		ctx.Error(errString)
		ctx.String(http.StatusNotFound, errString.Error())
		return
	}

	ctx.JSON(200, sets)
}

func (server *Server) CreateExercise(ctx *gin.Context) {
	type sets struct {
		Reps int32 `json:"reps"`
		Weight int32 `json:"weight"`
	}
	type request struct {
		Name string `json:"name"`
		Description pgtype.Text `json:"description"`
		Href pgtype.Text `json:"href"`
		Equipment pgtype.Text `json:"equipment"`
		TargetMuscle string `json:"target_muscle"`
		Replacement pgtype.Int8 `json:"replacement"`
		Sets []sets `json:"sets"`
	}

	var req request
	err := ctx.ShouldBindBodyWithJSON(req)
	if err != nil {
		errString := fmt.Errorf("bad request could not bind json: %s", err)
		ctx.Error(errString)
		ctx.String(http.StatusBadRequest, errString.Error())
		return
	}
	exercise := db.InsertExerciseParams{
		Name: req.Name,
		Description: req.Description,
		Sets: int32(len(req.Sets)),
		Href: req.Href,
		Equipment: req.Equipment,
		TargetMuscle: req.TargetMuscle,
		Replacement: req.Replacement,
	}
	exerciseId, err := server.store.InsertExercise(context.Background(), exercise)
	if err != nil {
		errString := fmt.Errorf("unable to insert exercise: %s", err)
		ctx.Error(errString)
		ctx.String(http.StatusInternalServerError, errString.Error())
		return
	}
	for _, set := range req.Sets {
		dbSet := db.InsertSetParams{
			ExerciseID: exerciseId,
			Reps: set.Reps,
		}
		err = server.store.InsertSet(context.Background(), dbSet)
		if err != nil {
			errString := fmt.Errorf("unable to insert set: %s", err)
			ctx.Error(errString)
			ctx.String(http.StatusInternalServerError, errString.Error())
			return
		}
	}
	ctx.String(http.StatusOK, fmt.Sprintf("Created new exercise with id: %v", exerciseId))
}

  

func RunApi(dbUrl string) *gin.Engine {
	conn, err := pgxpool.New(context.Background(), dbUrl)
	if err != nil {
		log.Fatalf("Unable to connect to db: %v\n", err)
	}

	server := &Server{
		store: db.New(conn),
	}

	router := gin.Default()

	router.GET("/exercises", server.GetExercisesHandler)
	router.GET("/exercise/:id", server.GetExerciseByIdHandler)
	router.GET("/equipment/:equipment", server.GetExerciseByEquipment)
	router.GET("/replacement/:id", server.GetReplacmentByIdHandler)
	router.GET("/muscle/:muscle", server.GetExerciseByMuscle)
	router.GET("/sets/:id", server.GetSetsByExerciseId)
	router.POST("/exercise", server.CreateExercise)

	return router
}

func main() {
	router := RunApi(os.Getenv("DATABASE_URL"))
	log.Fatal(router.Run())
}
