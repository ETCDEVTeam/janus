package gcp

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

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
func SendToGCP(to, files, key string) error {

	to = filepath.Clean(to)
	files = filepath.Clean(files)
	key = filepath.Clean(key)

	// Ensure key file exists.
	if _, e := os.Stat(key); e != nil {
		return e
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
			return errors.New("env GCP_PASSWD not set, cannot decrypt")
		}
		// Assume reading for decoding error is it's encrypted... attempt to decrypt
		if decryptError := exec.Command("openssl", "aes-256-cbc", "-k", passwd, "-in", key, "-out", decryptedKeyFileName, "-d").Run(); decryptError != nil {
			return decryptError
		}

		log.Println("decrypted key file to ", decryptedKeyFileName)
		key = decryptedKeyFileName

		// Only remove *unecrypted* key file
		defer func() {
			log.Printf("removing key: %v", key)
			if errRm := os.Remove(key); errRm != nil {
				log.Println(errRm)
			}
		}()
	}

	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithServiceAccountFile(key))
	if err != nil {
		return err
	}

	// Use glob to get matching file paths.
	globs, e := filepath.Glob(files)
	if e != nil {
		return e
	}
	// Ensure there is something to upload
	if len(globs) == 0 {
		return errors.New("no files matching '-to' pattern were found")
	}
	// Upload each file
	for _, f := range globs {
		fi, e := os.Stat(f)
		if e != nil {
			return e
		}
		if fi.IsDir() {
			log.Printf("%s is a directory, continuing", fi.Name())
			continue
		}
		// eg.
		// to: builds.etcdevteam.com/go-ethereum/3.5.x
		// file: ./dist/geth.zip
		//
		// Set bucket as first in separator-split path
		// eg. builds.etcdevteam.com
		bucket := strings.Split(filepath.ToSlash(to), "/")[0]

		// Get relative path of 'to' based on 'bucket'
		// eg. go-ethereum/3.5.x
		object, relError := filepath.Rel(bucket, to)
		if relError != nil {
			return relError
		}

		// Append file to 'to' path.
		// eg. go-ethereum/3.5.x/geth.zip
		deployObject := filepath.Join(object, filepath.Base(f))
		// Ensure actual '/' [slash]es are used, because Windows will make them '\' [backslashes]
		// and google won't folder-ize the path
		deployObject = filepath.ToSlash(deployObject)

		// Send it.
		if deployError := writeToGCP(client, bucket, deployObject, f); deployError != nil {
			return deployError
		}
	}
	return nil
}
