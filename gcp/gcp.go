package gcp

import (
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

var decryptedKeyFileName = "./" + strconv.Itoa(os.Getpid()) + "-gcloud.json"

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

	var arbitraryMap = make(map[string]interface{})

	// Read key file
	bRead, eRead := ioutil.ReadFile(key)
	if eRead != nil {
		return eRead
	}

	// Attempt to unmarshal key file, checks for encryption
	e := json.Unmarshal(bRead, &arbitraryMap)
	if e != nil {
		log.Println("key is possibly encryped, attempting to decrypt with $GCP_PASSWD")

		passwd := os.Getenv("GCP_PASSWD")
		if passwd == "" {
			log.Fatalln("env GCP_PASSWD not set, cannot decrypt")
		}
		// Assume reading for decoding error is it's encrypted... attempt to decrypt
		if e := exec.Command("openssl", "aes-256-cbc", "-k", passwd, "-in", key, "-out", decryptedKeyFileName, "-d").Run(); e != nil {
			log.Fatal("could not parse nor decrypt given key file (please ensure env var GCP_PASSWD is set)", key, e)
		}

		log.Println("decrypted key file to ", decryptedKeyFileName)
		key = decryptedKeyFileName

		defer os.Remove(key) // Only remove *unecrypted* key file
	}

	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithServiceAccountFile(key))
	if err != nil {
		log.Fatal(err)
	}

	deployError := writeToGCP(client, bucket, object, file)
	return deployError
}
