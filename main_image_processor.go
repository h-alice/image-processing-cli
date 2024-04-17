// Descritpion: This is the declaration of the main worker, and file processing subroutines.
package main

import (
	"context"
	"imagetools/config"
	"log"
	"path/filepath"
)

// Process image file with profile.
//
// profile: Image processing profile.
func processFile(profile config.ImageProcessingProfile) error {

	working_image, err := profile.CreateImageFile()
	if err != nil {
		log.Printf("[x] Error while creating image: %v", err)
		return err
	}

	// Create image processing pipeline.
	for _, pb := range profile.PipelineBlocks {
		// log.Printf("Processing Operation #%d: %s", index, pb.Operation)
		working_image = working_image.Then(config.PipelineBlockToOperation(pb))
		if working_image.LastError() != nil {
			log.Printf("[x] Error while processing image: %v", working_image.LastError())
			return working_image.LastError()
		}
	}

	return nil
}

// Main worker, this subroutine is designed to be run in a goroutine.
//
// ctx: Context.
// profile: Image processing profile.
// result_chan: Result channel, used to send result back to main thread.
func mainWorker(ctx context.Context, profile config.ImageProcessingProfile, result_chan chan<- error) {

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Get only file name from the path.
	input_image_name := profile.GetAssignedFilePath()
	input_image_name = filepath.Base(input_image_name)

	log.Printf("[.] Processing image [%s] with profile [%s]\n", input_image_name, profile.ProfileName)

	// Process image.
	for {
		select {
		case <-ctx.Done(): // Check if context is cancelled.
			log.Printf("[!] Context cancelled. Exiting...\n")
			result_chan <- nil
			return // Terminate goroutine.
		default:
			err := processFile(profile) // Process image.
			if err != nil {
				log.Printf("[x] Error while processing image: %v\n", err)
				result_chan <- err
				return // Terminate goroutine.
			}
			result_chan <- nil
			return // Goroutine finished.
		}
	}
}
