/*
 * Copyright (c) 2023. yimincai(Neil) <bravc29229@gmail.com>.
 */

package minio

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/schollz/progressbar/v3"
	"github.com/yimincai/toolbox/pkg/logger"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// DestinationDir Bucket name and destination directory to save dumped files
var (
	Endpoint       = ""
	BucketName     = ""
	User           = ""
	Password       = ""
	DestinationDir = "./backup/minio"
	UseSSL         = false
)

func newMinioClient(endpoint, user, password string, ssl bool) *minio.Client {
	logger.Green(fmt.Sprintf("Initializing Minio client on %s using %s bucket", Endpoint, BucketName))

	// Initialize Minio client object
	minioClient, err := minio.New(Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(User, Password, ""),
		Secure: UseSSL,
	})
	if err != nil {
		log.Fatalln("Error initializing Minio client:", err)
	}

	return minioClient
}

// DumpBucket dumps all objects in a bucket to a directory
func DumpBucket(worker int) {
	minioClient := newMinioClient(Endpoint, User, Password, UseSSL)

	// List objects in the bucket
	ctx := context.Background()
	objectCh := minioClient.ListObjects(ctx, BucketName, minio.ListObjectsOptions{
		Recursive: true,
	})

	// Create destination directory if it does not exist
	if err := os.MkdirAll(DestinationDir, os.ModePerm); err != nil {
		log.Fatalf("Error creating destination directory: %v\n", err)
	}

	// Delete all files in the destination directory
	err := os.RemoveAll(DestinationDir)
	if err != nil {
		log.Fatal(err)
	}
	logger.Green("All files in the destination directory deleted.")

	// Create a channel to communicate with workers
	downloadCh := make(chan minio.ObjectInfo)

	var wg sync.WaitGroup

	// Start workers
	for i := 0; i < worker; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for object := range downloadCh {
				objectName := object.Key
				destinationPath := fmt.Sprintf("%s/%s", DestinationDir, objectName)

				// Create directories if the object has a directory structure
				if strings.Contains(objectName, "/") {
					dir := strings.Join(strings.Split(objectName, "/")[:len(strings.Split(objectName, "/"))-1], "/")
					if err := os.MkdirAll(fmt.Sprintf("%s/%s", DestinationDir, dir), os.ModePerm); err != nil {
						log.Fatalf("Error creating directory structure: %v\n", err)
					}
				}

				// Download object
				err := minioClient.FGetObject(ctx, BucketName, objectName, destinationPath, minio.GetObjectOptions{})
				if err != nil {
					logger.Red(fmt.Sprintf("Error downloading object %s: %v", objectName, err))
				}
			}
		}()
	}

	// Count the total number of objects for the progress bar
	totalObjects := 0
	for range objectCh {
		totalObjects++
	}
	objectCh = minioClient.ListObjects(ctx, BucketName, minio.ListObjectsOptions{
		Recursive: true,
	}) // Reset objectCh after counting

	// Create a progress bar
	bar := progressbar.Default(int64(totalObjects), "Downloading objects")

	// Send objects to the workers through the channel
	go func() {
		defer close(downloadCh)
		for object := range objectCh {
			if object.Err != nil {
				log.Fatalln("Error listing objects:", object.Err)
			}
			downloadCh <- object

			err = bar.Add(1) // Increment progress bar
			if err != nil {
				logger.Red(fmt.Sprint("Error incrementing progress bar:", err))
			}
		}
	}()

	// Wait for all workers to complete
	wg.Wait()

	bar.Finish() // Finish the progress bar
	log.Printf("Bucket %s dump completed.", BucketName)
}

// DeleteBucket deletes all objects in a bucket
func DeleteBucket(worker int) {
	minioClient := newMinioClient(Endpoint, User, Password, UseSSL)

	// List objects in the bucket
	ctx := context.Background()
	objectCh := minioClient.ListObjects(ctx, BucketName, minio.ListObjectsOptions{
		Recursive: true,
	})

	// Create a channel to communicate with workers
	deleteCh := make(chan minio.ObjectInfo)

	var wg sync.WaitGroup

	// Calculate the total number of objects for the progress bar
	totalObjects := 0
	for range objectCh {
		totalObjects++
	}
	objectCh = minioClient.ListObjects(ctx, BucketName, minio.ListObjectsOptions{
		Recursive: true,
	}) // Reset objectCh after counting

	// Start workers
	for i := 0; i < worker; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for object := range deleteCh {
				objectName := object.Key
				err := minioClient.RemoveObject(ctx, BucketName, objectName, minio.RemoveObjectOptions{})
				if err != nil {
					// Handle error if needed
					logger.Red(fmt.Sprintf("Error deleting object %s: %v", objectName, err))
				}
			}
		}()
	}

	// Create a progress bar
	bar := progressbar.Default(int64(totalObjects), "Deleting objects")

	// Send objects to the workers through the channel
	go func() {
		defer close(deleteCh)
		for object := range objectCh {
			if object.Err != nil {
				log.Fatalln("Error listing objects:", object.Err)
			}
			deleteCh <- object

			err := bar.Add(1) // Increment progress bar
			if err != nil {
				logger.Red(fmt.Sprint("Error incrementing progress bar:", err))
			}
		}
	}()

	// Wait for all workers to complete
	wg.Wait()

	bar.Finish() // Finish the progress bar
	log.Print("Bucket cleanup completed.")
}

// UploadBucket uploads all files in a directory to a bucket
func UploadBucket(worker int) {
	minioClient := newMinioClient(Endpoint, User, Password, UseSSL)

	// Create a channel to communicate with workers
	uploadCh := make(chan string)

	var wg sync.WaitGroup

	// Start workers
	for i := 0; i < worker; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for objectName := range uploadCh {
				func() {
					defer func() {
						if r := recover(); r != nil {
							logger.Red(fmt.Sprintf("Recovered panic: %v", r))
						}
					}()

					file, err := os.Open(filepath.Join(DestinationDir, objectName))
					if err != nil {
						logger.Red(fmt.Sprintf("Error opening file %s: %v", objectName, err))
						return
					}
					defer file.Close()

					// Get file information
					fileInfo, err := file.Stat()
					if err != nil {
						logger.Red(fmt.Sprintf("Error getting file info for %s: %v", objectName, err))
						return
					}

					// Upload the file to the Minio bucket
					_, err = minioClient.PutObject(context.Background(), BucketName, objectName, file, fileInfo.Size(), minio.PutObjectOptions{})
					if err != nil {
						logger.Red(fmt.Sprintf("Error uploading object %s: %v", objectName, err))
					}
				}()
			}
		}()
	}

	// Walk through the restore directory and send files to workers through the channel
	files, err := filepath.Glob(filepath.Join(DestinationDir, "*"))
	if err != nil {
		logger.Red(fmt.Sprintf("Error listing files: %v", err))
		return
	}

	// Calculate the total number of files for the progress bar
	totalFiles := len(files)

	// Create a progress bar
	bar := progressbar.Default(int64(totalFiles), "Uploading files")

	// Send objectName to workers for uploading
	for _, file := range files {
		objectName, err := filepath.Rel(DestinationDir, file)
		if err != nil {
			logger.Red(fmt.Sprintf("Error calculating object name: %v", err))
			continue
		}
		uploadCh <- filepath.ToSlash(objectName)

		err = bar.Add(1) // Increment progress bar
		if err != nil {
			logger.Red(fmt.Sprint("Error incrementing progress bar:", err))
		}
	}

	// Close the channel and wait for all workers to complete
	close(uploadCh)
	wg.Wait()

	bar.Finish() // Finish the progress bar
	log.Println("Upload completed.")
}

// RestoreBucket restores a bucket from a directory
func RestoreBucket(worker int) {
	minioClient := newMinioClient(Endpoint, User, Password, UseSSL)

	// List objects in the bucket
	objectCh := minioClient.ListObjects(context.Background(), BucketName, minio.ListObjectsOptions{
		Recursive: true,
	})

	var totalDeleteObjects int
	for range objectCh {
		totalDeleteObjects++
	}

	// Reset objectCh channel after counting total objects
	objectCh = minioClient.ListObjects(context.Background(), BucketName, minio.ListObjectsOptions{
		Recursive: true,
	})

	deleteCh := make(chan string, worker)
	uploadCh := make(chan string, worker)

	var dwg sync.WaitGroup

	// Start delete workers
	for i := 0; i < worker; i++ {
		dwg.Add(1)
		go func() {
			defer dwg.Done()
			for objectName := range deleteCh {
				err := minioClient.RemoveObject(context.Background(), BucketName, objectName, minio.RemoveObjectOptions{})
				if err != nil {
					logger.Red(fmt.Sprintf("Error deleting object %s: %v", objectName, err))
				}
			}
		}()
	}

	// Create progress bars
	deleteBar := progressbar.Default(int64(totalDeleteObjects), "Deleting files")

	// Send object names to delete and upload workers and update progress bars
	for object := range objectCh {
		if object.Err != nil {
			logger.Red(fmt.Sprintf("Error listing objects: %v", object.Err))
			continue
		}
		objectName := object.Key
		deleteCh <- objectName

		err := deleteBar.Add(1)
		if err != nil {
			logger.Red(fmt.Sprint("Error incrementing progress bar:", err))
		}
	}

	// Close the channels and wait for all workers to complete
	close(deleteCh)
	dwg.Wait()
	deleteBar.Finish()

	var uwg sync.WaitGroup
	// Start upload workers
	// Start workers
	for i := 0; i < worker; i++ {
		uwg.Add(1)
		go func() {
			defer uwg.Done()
			for objectName := range uploadCh {
				func() {
					defer func() {
						if r := recover(); r != nil {
							logger.Red(fmt.Sprintf("Recovered panic: %v", r))
						}
					}()

					file, err := os.Open(filepath.Join(DestinationDir, objectName))
					if err != nil {
						logger.Red(fmt.Sprintf("Error opening file %s: %v", objectName, err))
						return
					}
					defer file.Close()

					// Get file information
					fileInfo, err := file.Stat()
					if err != nil {
						logger.Red(fmt.Sprintf("Error getting file info for %s: %v", objectName, err))
						return
					}

					// Upload the file to the Minio bucket
					_, err = minioClient.PutObject(context.Background(), BucketName, objectName, file, fileInfo.Size(), minio.PutObjectOptions{})
					if err != nil {
						logger.Red(fmt.Sprintf("Error uploading object %s: %v", objectName, err))
					}
				}()
			}
		}()
	}

	// Walk through the restore directory and send files to workers through the channel
	files, err := filepath.Glob(filepath.Join(DestinationDir, "*"))
	if err != nil {
		logger.Red(fmt.Sprintf("Error listing files: %v", err))
		return
	}

	// Calculate the total number of files for the progress bar
	totalFiles := len(files)

	// Create a progress restoreBar
	restoreBar := progressbar.Default(int64(totalFiles), "Uploading files")

	// Send objectName to workers for uploading
	for _, file := range files {
		objectName, err := filepath.Rel(DestinationDir, file)
		if err != nil {
			logger.Red(fmt.Sprintf("Error calculating object name: %v", err))
			continue
		}
		uploadCh <- filepath.ToSlash(objectName)

		err = restoreBar.Add(1) // Increment progress bar
		if err != nil {
			logger.Red(fmt.Sprint("Error incrementing progress bar:", err))
		}
	}

	// Close the channel and wait for all workers to complete
	close(uploadCh)
	uwg.Wait()
	restoreBar.Finish() // Finish the progress bar

	log.Println("Bucket cleanup and restoration completed.")
}
