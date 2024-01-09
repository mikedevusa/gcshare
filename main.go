package main

import (
    "fmt"
    "context"
	"time"
    "io"
    "os"
	"path/filepath"
	gomail "gopkg.in/gomail.v2"
    "cloud.google.com/go/storage"
    "gcshare/config"
)

var configuration = config.GetConfig()

func copyFileToGCS(localFilePath string, bucketName string, objectName string) error {
    ctx := context.Background()

    f, err := os.Open(localFilePath)
    if err != nil {
        return fmt.Errorf("failed to open local file: %v", err)
    }
    defer f.Close()

    client, err := storage.NewClient(ctx)
    if err != nil {
        return fmt.Errorf("failed to create storage client: %v", err)
    }
    defer client.Close()

    bucket := client.Bucket(bucketName)
    wc := bucket.Object(objectName).NewWriter(ctx)

    // Copy the contents of the local file to the object.
    if _, err := io.Copy(wc, f); err != nil {
        return fmt.Errorf("failed to copy file to GCS: %v", err)
    }

    // Close the writer to finalize the object.
    if err := wc.Close(); err != nil {
        return fmt.Errorf("failed to close writer: %v", err)
    }

    return nil
}

func generateSignedURL(bucketName, objectName string, expiration time.Duration) (string, error) {
    ctx := context.Background()
    client, err := storage.NewClient(ctx)
    if err != nil {
        return "", fmt.Errorf("storage.NewClient: %v", err)
    }

	opts := &storage.SignedURLOptions{
		Scheme:  storage.SigningSchemeV4,
		Method:  "GET",
		Expires: time.Now().Add(15 * time.Minute),
	}
	u, err := client.Bucket(bucketName).SignedURL(objectName, opts)
	if err != nil {
		return "", fmt.Errorf("Bucket(%q).SignedURL: %w", bucketName, err)
	}

    return u, nil
}

func sendUrl(recipientEmail string, signedUrl string, objectName string, password string) error {
    msg := gomail.NewMessage()
    msg.SetHeader("From", configuration.User)
    msg.SetHeader("To", recipientEmail)
    msg.SetHeader("Subject", objectName)
    msg.SetBody("text/html", signedUrl)

    n := gomail.NewDialer(configuration.Mailhost, configuration.Mailport, configuration.User, configuration.Password)

    if err := n.DialAndSend(msg); err != nil {
		return fmt.Errorf("Could not send email")
    }

	return nil
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Syntax: gcshare <filepath> <recipient_email>")
		os.Exit(1)
	}

	localFile := os.Args[1]
	recipient := os.Args[2]
	objectName := filepath.Base(localFile)

    if err := copyFileToGCS(localFile, configuration.Bucketname, objectName); err != nil {
        fmt.Println(err)
    } 

    expiration := 15 * time.Minute // Expire URL after 15 minutes

    url, err := generateSignedURL(configuration.Bucketname, objectName, expiration)
    if err != nil {
        fmt.Println("Error generating signed URL:", err)
    } else {
        fmt.Println("Generated signed URL:", url)
    }

	if err := sendUrl(recipient, url, objectName, configuration.Password); err != nil {
		fmt.Println("Error occurred sending email!")
	} else {
		fmt.Println("Send succcessful to ", recipient)
	}

}