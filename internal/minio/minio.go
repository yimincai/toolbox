/*
 * Copyright (c) 2023. yimincai(Neil) <bravc29229@gmail.com>.
 */

package minio

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/spf13/viper"
)

// DestinationDir Bucket name and destination directory to save dumped files
var DestinationDir = "./backup/minio"

type Env struct {
	PrintEnv   bool   `mapstructure:"PRINT_ENV"`
	Endpoint   string `mapstructure:"MINIO_ENDPOINT"`
	BucketName string `mapstructure:"MINIO_BUCKET_NAME"`
	Port       string `mapstructure:"MINIO_PORT"`
	User       string `mapstructure:"MINIO_ROOT_USER"`
	Password   string `mapstructure:"MINIO_ROOT_PASSWORD"`
	SSL        bool   `mapstructure:"MINIO_SSL"`
}

func newEnv() *Env {
	// init env
	env := &Env{}
	envFile := ".env"
	viper.SetConfigFile(envFile)

	err := viper.ReadInConfig()
	if err == nil {
		log.Println("Using environment config file:", viper.ConfigFileUsed())
	}
	if err != nil {
		log.Println("Use environment variable")
		viper.AutomaticEnv()
		_ = viper.BindEnv("PRINT_ENV")
		_ = viper.BindEnv("MINIO_ENDPOINT")
		_ = viper.BindEnv("MINIO_BUCKET_NAME")
		_ = viper.BindEnv("MINIO_PORT")
		_ = viper.BindEnv("MINIO_ROOT_USER")
		_ = viper.BindEnv("MINIO_ROOT_PASSWORD")
		_ = viper.BindEnv("MINIO_SSL")
	}

	err = viper.Unmarshal(&env)
	if err != nil {
		log.Fatal("Environment can't be loaded:", err)
	}

	if env.PrintEnv {
		e := prettyPrint(&env)
		log.Printf("Env:%s\n", e)
	}

	return env
}

func prettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}

func newMinioClient(env *Env) *minio.Client {
	// Initialize Minio client object
	minioClient, err := minio.New(env.Endpoint+":"+env.Port, &minio.Options{
		Creds:  credentials.NewStaticV4(env.User, env.Password, ""),
		Secure: env.SSL,
	})
	if err != nil {
		log.Fatalln("Error initializing Minio client:", err)
	}

	return minioClient
}

func DumpBucket() {
	env := newEnv()
	minioClient := newMinioClient(env)

	// List objects in the bucket
	ctx := context.Background()
	objectCh := minioClient.ListObjects(ctx, env.BucketName, minio.ListObjectsOptions{
		Recursive: true,
	})

	// Create destination directory if it does not exist
	if err := os.MkdirAll(DestinationDir, os.ModePerm); err != nil {
		log.Fatalf("Error creating destination directory: %v\n", err)
	}

	// Delete all files in destination directory
	err := os.RemoveAll(DestinationDir)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("All files in destination directory deleted.")

	// Download objects from the bucket
	for object := range objectCh {
		if object.Err != nil {
			log.Fatalln("Error listing objects:", object.Err)
		}
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
		err := minioClient.FGetObject(ctx, env.BucketName, objectName, destinationPath, minio.GetObjectOptions{})
		if err != nil {
			log.Printf("Error downloading object %s: %v\n", objectName, err)
		} else {
			log.Printf("Downloaded: %s\n", objectName)
		}
	}

	log.Print("Bucket dump completed.")
}

func DeleteBucket() {
	env := newEnv()
	minioClient := newMinioClient(env)

	// List objects in the bucket
	ctx := context.Background()
	objectCh := minioClient.ListObjects(ctx, env.BucketName, minio.ListObjectsOptions{
		Recursive: true,
	})

	// Delete objects from the bucket
	for object := range objectCh {
		if object.Err != nil {
			log.Fatalln("Error listing objects:", object.Err)
		}
		objectName := object.Key
		err := minioClient.RemoveObject(ctx, env.BucketName, objectName, minio.RemoveObjectOptions{})
		if err != nil {
			log.Printf("Error deleting object %s: %v\n", objectName, err)
		} else {
			log.Printf("Deleted: %s\n", objectName)
		}
	}

	log.Print("Bucket cleanup completed.")
}

func UploadBucket() {
	env := newEnv()
	minioClient := newMinioClient(env)

	// Walk through the restore directory and upload files to the Minio bucket
	err := filepath.WalkDir(DestinationDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			objectName := filepath.ToSlash(path[len(DestinationDir)+1:])

			fileInfo, err := d.Info()
			if err != nil {
				return err
			}
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			// Upload the file to the Minio bucket
			_, err = minioClient.PutObject(context.Background(), env.BucketName, objectName, file, fileInfo.Size(), minio.PutObjectOptions{})
			if err != nil {
				return err
			}

			fmt.Printf("Uploaded: %s\n", objectName)
		}
		return nil
	})
	if err != nil {
		log.Fatalln("Error restoring files:", err)
	}

	fmt.Println("Upload completed.")
}

func RestoreBucket() {
	env := newEnv()
	minioClient := newMinioClient(env)

	// List objects in the bucket
	ctx := context.Background()
	objectCh := minioClient.ListObjects(ctx, env.BucketName, minio.ListObjectsOptions{
		Recursive: true,
	})

	// Delete objects from the bucket
	for object := range objectCh {
		if object.Err != nil {
			log.Fatalln("Error listing objects:", object.Err)
		}
		objectName := object.Key
		err := minioClient.RemoveObject(ctx, env.BucketName, objectName, minio.RemoveObjectOptions{})
		if err != nil {
			log.Printf("Error deleting object %s: %v\n", objectName, err)
		} else {
			log.Printf("Deleted: %s\n", objectName)
		}
	}

	log.Print("Bucket cleanup completed.")

	// Walk through the restore directory and upload files to the Minio bucket
	err := filepath.WalkDir(DestinationDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			objectName := filepath.ToSlash(path[len(DestinationDir)+1:])

			fileInfo, err := d.Info()
			if err != nil {
				return err
			}
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			// Upload the file to the Minio bucket
			_, err = minioClient.PutObject(context.Background(), env.BucketName, objectName, file, fileInfo.Size(), minio.PutObjectOptions{})
			if err != nil {
				return err
			}

			fmt.Printf("Restored: %s\n", objectName)
		}
		return nil
	})
	if err != nil {
		log.Fatalln("Error restoring files:", err)
	}

	fmt.Println("Restoration completed.")
}
