// Copyright (C) 2019 Storj Labs, Inc.
// See LICENSE for copying information.

package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"storj.io/uplink"
)

const (
	myAPIKey = "13Yqgf2KjVB1mSUchWMVe2WNLW2Ynwue77nGtjeWEaunDcyHRC9cDwk98i7QEgV3z7Fn2Eav2EXsRBVS5oYss98pcY8au69TPk5GFUo"

	satellite    = "us-central-1.tardigrade.io:7777"
	myBucket     = "my-first-bucket"
	myUploadKey  = "foo/bar/baz"
	myData       = "one fish two fish red fish blue fish"
	myPassphrase = "you'll never guess this"
)

// UploadAndDownloadData uploads the specified data to the specified key in the
// specified bucket, using the specified Satellite, API key, and passphrase.
func UploadAndDownloadData(ctx context.Context,
	satelliteAddress, apiKey, passphrase, bucketName, uploadKey string,
	dataToUpload []byte) error {

	// Request access grant to the satellite with the API key and passphrase.
	access, err := uplink.RequestAccessWithPassphrase(ctx, satelliteAddress, apiKey, passphrase)
	if err != nil {
		return fmt.Errorf("could not request access grant: %v", err)
	}

	// Open up the Project we will be working with.
	project, err := uplink.OpenProject(ctx, access)
	if err != nil {
		return fmt.Errorf("could not open project: %v", err)
	}
	defer project.Close()

	// Ensure the desired Bucket within the Project is created.
	_, err = project.EnsureBucket(ctx, bucketName)
	if err != nil {
		return fmt.Errorf("could not ensure bucket: %v", err)
	}

	// Intitiate the upload of our Object to the specified bucket and key.
	upload, err := project.UploadObject(ctx, bucketName, uploadKey, nil)
	if err != nil {
		return fmt.Errorf("could not initiate upload: %v", err)
	}

	// Copy the data to the upload.
	buf := bytes.NewBuffer(dataToUpload)
	_, err = io.Copy(upload, buf)
	if err != nil {
		_ = upload.Abort()
		return fmt.Errorf("could not upload data: %v", err)
	}

	// Commit the uploaded object.
	err = upload.Commit()
	if err != nil {
		return fmt.Errorf("could not commit uploaded object: %v", err)
	}

	// Initiate a download of the same object again
	download, err := project.DownloadObject(ctx, bucketName, uploadKey, nil)
	if err != nil {
		return fmt.Errorf("could not open object: %v", err)
	}
	defer download.Close()

	// Read everything from the download stream
	receivedContents, err := ioutil.ReadAll(download)
	if err != nil {
		return fmt.Errorf("could not read data: %v", err)
	}

	// Check that the downloaded data is the same as the uploaded data.
	if !bytes.Equal(receivedContents, dataToUpload) {
		return fmt.Errorf("got different object back: %q != %q", dataToUpload, receivedContents)
	}

	return nil
}

func main() {
	err := UploadAndDownloadData(context.Background(),
		satellite, myAPIKey, myPassphrase, myBucket, myUploadKey, []byte(myData))
	if err != nil {
		log.Fatalln("error:", err)
	}

	fmt.Println("success!")
}
