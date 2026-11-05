package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/Kagami/go-face"
	"github.com/andre-ols/go-face-recognition/internal/entity"
	"github.com/andre-ols/go-face-recognition/internal/input"
	"github.com/andre-ols/go-face-recognition/internal/usecases"
)

var (
	modelsDir  = "models"
	personsDir = "persons"
	fontPath   = filepath.Join("fonts", "Roboto-Regular.ttf")
	threshold  = float32(0.3)
)

// This example shows the basic usage of the package: create an
// recognizer, recognize faces, classify them using few known ones.
func main() {
	initialTime := time.Now()

	defer func() {
		fmt.Printf("\x1b[34mTotal time: %s\x1b[0m\n", time.Since(initialTime))
	}()

	args := os.Args[1:]
	jobs, err := input.Resolve(args)
	if err != nil {
		log.Fatalf("Can't resolve inputs: %v", err)
	}

	loadPersonsUseCase := usecases.NewLoadPersonsUseCase()

	// Load persons from fs.
	persons, err := loadPersonsUseCase.Execute(personsDir)
	if err != nil {
		log.Fatalf("Can't load persons: %v", err)
	}

	initRecognizerTime := time.Now()

	// Init the recognizer.
	rec, err := face.NewRecognizer(modelsDir)
	if err != nil {
		log.Fatalf("Can't init face recognizer: %v", err)
	}
	// Free the resources when you're finished.
	defer rec.Close()

	fmt.Println("Time to init recognizer: ", time.Since(initRecognizerTime))

	// Recognize faces from a few images and assign them to persons.
	recognizePersonsUseCase := usecases.NewRecognizePersonsUseCase(rec)

	err = recognizePersonsUseCase.Execute(persons)
	if err != nil {
		log.Fatalf("Can't recognize persons: %v", err)
	}

	classifyPersonUseCase := usecases.NewClassifyPersonsUseCase(rec)

	// Explicit args or multi-job runs continue on per-image errors.
	continueOnError := len(args) > 0
	var failed int

	for _, job := range jobs {
		if err := processJob(job, persons, classifyPersonUseCase); err != nil {
			if continueOnError {
				log.Printf("Can't process %s: %v", job.In, err)
				failed++
				continue
			}
			log.Fatalf("Can't process %s: %v", job.In, err)
		}
	}

	if failed > 0 {
		log.Printf("%d of %d image(s) failed", failed, len(jobs))
		os.Exit(1)
	}
}

func processJob(job input.Job, persons map[int]*entity.Person, classify usecases.ClassifyPersonsUseCase) error {
	fmt.Printf("\033[0;36mProcessing %s\033[0m\n", job.In)

	recognizedFaces, err := classify.Execute(job.In, threshold)
	if err != nil {
		return fmt.Errorf("classify: %w", err)
	}

	fmt.Printf("\033[0;32mFound %d faces\033[0m\n", len(recognizedFaces))
	for _, recognizedFace := range recognizedFaces {
		fmt.Printf("\033[0;32mPerson: %s\033[0m\n", persons[recognizedFace.ID].Name)
	}

	drawer := entity.NewDrawer(job.In, fontPath)
	for _, recognizedFace := range recognizedFaces {
		drawer.DrawFace(recognizedFace.Face.Rectangle, persons[recognizedFace.ID].Name)
	}

	if err := drawer.SaveImage(job.Out); err != nil {
		return fmt.Errorf("save %s: %w", job.Out, err)
	}

	fmt.Printf("\033[0;32mWrote %s\033[0m\n", job.Out)
	return nil
}
