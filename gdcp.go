// Copyright 2011 Richard Larocque. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v2"
)

func usage(progname string) {
	fmt.Fprintf(os.Stderr, "Usage: %s [options] SOURCE DEST\n", progname)
	flag.PrintDefaults()
	os.Exit(2)
}

var (
	srcFile, dstFile                string
	update, allowDupes, keepHistory bool
)

func main() {
	flag.BoolVar(&update, "update", false,
		"If a file of the same name already exists, update it.")
	flag.BoolVar(&allowDupes, "allowDupes", false,
		"Allow this client to create files whose names shadow existing files.")
	flag.BoolVar(&keepHistory, "keepHistory", true,
		"Whether or not file history is preserved.")
	flag.Parse()

	if flag.NArg() != 2 {
		usage(path.Base(os.Args[0]))
	}

	if (allowDupes && update) {
		log.Fatal("--allowDupes and --update flags are incompatible");
	}

	srcFile = flag.Arg(0)
	dstFile = flag.Arg(1)

	ctx := context.Background()
	// b, err := ioutil.ReadFile("client_secret.json")
	// if err != nil {
	// 	log.Fatalf("Unable to read client secret file: %v", err)
	//}
	// It's more convenient to embed secrets at compile time.
	config, err := google.ConfigFromJSON(ClientSecretJson(), drive.DriveScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(ctx, config)

	service, err := drive.New(client)
	if err != nil {
		log.Fatal(err)
	}

	fromFile, err := os.Open(srcFile)
	if err != nil {
		log.Fatalf("Error opening %s: %v", srcFile, err)
	}

	fileList, err := service.Files.
		List().
		Q(fmt.Sprintf("title='%s'", dstFile)).
		MaxResults(2).Do()
	if err != nil {
		log.Fatal(err)
	} else if len(fileList.Items) > 1 && update {
		log.Fatalf(
			"Many files match name '%s'. Aborting update.",
			dstFile)
	} else if len(fileList.Items) == 1 && !(allowDupes || update) {
		log.Fatalf(
			"File named '%s' exists. Will not upload dupe.",
			dstFile)
	}

	if len(fileList.Items) == 0 || !update {
		_, err := service.Files.
			Insert(&drive.File{Title: dstFile}).
			Media(fromFile).Do()
		if err != nil {
			log.Fatal(err)
		}
	}

	if len(fileList.Items) == 1 && update {
		fileId := fileList.Items[0].Id
		_, err := service.Files.
			Update(fileId, fileList.Items[0]).
			NewRevision(keepHistory).
			Media(fromFile).Do()
		if err != nil {
			log.Fatal(err)
		}
	}
}
