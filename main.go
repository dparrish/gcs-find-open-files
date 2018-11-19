// Finds and lists any files in a Google Cloud Storage bucket that have any allUsers ACLs set on the file.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
)

func main() {
	flag.Parse()
	bucketName := flag.Arg(0)
	if bucketName == "" {
		fmt.Printf("Specify a bucket to list\n")
		os.Exit(1)
	}

	client, err := storage.NewClient(context.Background())
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	var matches int64
	st := time.Now()

	it := client.Bucket(bucketName).Objects(context.Background(), nil)
	for {
		attrs, err := it.Next()
		if err != nil {
			if err != iterator.Done {
				log.Printf("Error: %v", err)
			}
			break
		}
		for _, acl := range attrs.ACL {
			if acl.Entity == storage.AllUsers {
				fmt.Printf("gs://%s/%s %+v\n", bucketName, attrs.Name, acl)
				matches++
				break
			}
		}
	}

	fmt.Printf("Found %d total objects with open ACLs in %s", matches, time.Since(st))
}
