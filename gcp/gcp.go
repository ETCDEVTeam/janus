package gcp

import (
	"context"
	"io"
	"log"
	"os"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

// Use: go run gcs-deploy.go -bucket builds.etcdevteam.com -object go-ethereum/$(cat version-base.txt)/geth-classic-$TRAVIS_OS_NAME-$(cat version-app.txt).zip -file geth-classic-linux-14.0.zip -key ./.gcloud.key

// writeToGCP writes (uploads) a file or files to GCP Storage.
// 'object' is the path at which 'file' will be written,
// 'bucket' is the parent directory in which the object will be written.
func writeToGCP(client *storage.Client, bucket, object, file string) error {
	ctx := context.Background()
	// [START upload_file]
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	// Write object to storage, ensuring basename for file/object if exists.
	wc := client.Bucket(bucket).Object(object).NewWriter(ctx)
	if _, err = io.Copy(wc, f); err != nil {
		return err
	}
	if err := wc.Close(); err != nil {
		return err
	}
	// [END upload_file]
	log.Printf(`Successfully uploaded:
	bucket: %v
	object: %v
	file: %v`, bucket, object, file)
	return nil
}

// SendToGCP sends a file or files to Google Cloud Provider storage
// using a service account JSON key
func SendToGCP(bucket, object, file, key string) error {
	if _, e := os.Stat(file); e != nil {
		log.Fatal(file, e)
	}

	// Ensure key file exists.
	if _, e := os.Stat(key); e != nil {
		log.Fatal(file, e)
	}

	ctx := context.Background()

	client, err := storage.NewClient(ctx, option.WithServiceAccountFile(key))
	if err != nil {
		log.Fatal(err)
	}

	return writeToGCP(client, bucket, object, file)
}
