package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/z3co/workout-apiv2/db/gen"
)

type Server struct {
	store *db.Queries
}

func (server *Server) getExercisesHandler(ctx *gin.Context) {
	exercises, err := server.store.ListExercises(ctx.Request.Context())
	if err != nil {
		errString := fmt.Errorf("Exercise not found: %s", err)
		ctx.Error(errString)
		ctx.String(http.StatusNotFound, errString.Error())
	}

	ctx.JSON(200, exercises)
}

func (server *Server) getExerciseByIdHandler(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		errString := fmt.Errorf("Bad request %s", err)
		ctx.Error(errString)
		ctx.String(http.StatusBadRequest, errString.Error())
	}
	exercise, err := server.store.GetExercise(ctx.Request.Context(), int64(id))
	if err != nil {
		errString := fmt.Errorf("Exercise not found: %s", err)
		ctx.Error(errString)
		ctx.String(http.StatusNotFound, errString.Error())
		return
	}
	ctx.JSON(200, exercise)
}

func (server *Server) getExerciseByEquipment(ctx *gin.Context) {
	name := ctx.Param("equipment")
	exercise, err := server.store.GetExerciseByEquipment(ctx.Request.Context(), pgtype.Text{
		String: name,
		Valid: true,
	})
	if err != nil {
		errString := fmt.Errorf("Exercise not found: %s", err)
		ctx.Error(errString)
		ctx.String(http.StatusNotFound, errString.Error())
	}
	ctx.JSON(200, exercise)
}

func (server *Server) getExerciseByMuscle(ctx *gin.Context) {
	muscle := ctx.Param("muscle")
	exercise, err := server.store.GetExerciseByMuscle(ctx.Request.Context(), muscle)
	if err != nil {
		errString := fmt.Errorf("Exercise not found: %s", err)
		ctx.Error(errString)
		ctx.String(http.StatusNotFound, errString.Error())
	}
	fmt.Println(exercise)
	ctx.JSON(200, exercise)
}

func (server *Server) getReplacmentByIdHandler(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		errString := fmt.Errorf("Bad request %s\n", err)
		ctx.Error(errString)
		ctx.String(http.StatusBadRequest, errString.Error())
	}
	exercise, err := server.store.GetReplaceMentExercise(ctx.Request.Context(), int64(id))
	if err != nil {
		errString := fmt.Errorf("Exercise not found: %s", err)
		ctx.Error(errString)
		ctx.String(http.StatusNotFound, errString.Error())
	}
	ctx.JSON(200, exercise)
}

func (server *Server) getSetsByExerciseId(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		errString := fmt.Errorf("Bad request %s\n", err)
		ctx.Error(errString)
		ctx.String(http.StatusBadRequest, errString.Error())
	}

	sets, err := server.store.GetSetsByExerciseId(ctx.Request.Context(), int64(id))
	if err != nil {
		errString := fmt.Errorf("Exercise not found: %s", err)
		ctx.Error(errString)
		ctx.String(http.StatusNotFound, errString.Error())
	}

	ctx.JSON(200, sets)
}

func main() {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Unable to connect to db: %v\n", err)
	}
	defer conn.Close(context.Background())

	server := &Server{
		store: db.New(conn),
	}

	router := gin.Default()

	router.GET("/exercises", server.getExercisesHandler)
	router.GET("/exercise/:id", server.getExerciseByIdHandler)
	router.GET("/equipment/:equipment", server.getExerciseByEquipment)
	router.GET("/replacement/:id", server.getReplacmentByIdHandler)
	router.GET("/muscle/:muscle", server.getExerciseByMuscle)
	router.GET("/sets/:id", server.getSetsByExerciseId)

	router.Run()
}
