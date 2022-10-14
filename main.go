package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Neo-Desktop/Steamworks-Collection-Downloader/downloader"
	"github.com/Neo-Desktop/Steamworks-Collection-Downloader/manifest"
)

var (
	APIKEY       string
	SEEDS        []string
	PATH         string
	MANIFESTJSON string
)

const (
	REPOURL = "https://github.com/Neo-Desktop/Steamworks-Collection-Downloader"
)

func init() {
	flag.Usage = usage

	flagAPIKEY := flag.String("apikey", "", "Required: Steamworks Web API key [env: WSDL_APIKEY]")
	flagHELP := flag.Bool("help", false, "Displays this help message")
	flagPREFIX := flag.String("prefix", "."+string(os.PathSeparator)+"workshop", "path to the <gamedir>/addons/workshop folder [env: WSDL_PREFIX]")
	flagMANIFEST := flag.String("manifest", "workshopmanifest.json", "filename to store the current state of downloaded files (will be stored in prefix) [env: WSDL_MANIFEST]")

	envAPIKEY := os.Getenv("WSDL_APIKEY")
	envPREFIX := os.Getenv("WSDL_PREFIX")
	envSEEDS := os.Getenv("WSDL_SEEDS")
	envMANIFEST := os.Getenv("WSDL_MANIFEST")

	flag.Parse()

	if flagAPIKEY != nil && *flagAPIKEY != "" {
		APIKEY = *flagAPIKEY
	} else if envAPIKEY != "" {
		APIKEY = envAPIKEY
	} else {
		fmt.Fprintf(os.Stderr, "Please specify an APIKEY\n\n")
		usage()
		os.Exit(1)
	}

	if flagPREFIX != nil && *flagPREFIX != "" {
		PATH = *flagPREFIX
	} else if envPREFIX != "" {
		PATH = envPREFIX
	} else {
		fmt.Fprintf(os.Stderr, "Please specify a prefix\n\n")
		usage()
		os.Exit(1)
	}

	if flagMANIFEST != nil && *flagMANIFEST != "" {
		MANIFESTJSON = *flagMANIFEST
	} else if envMANIFEST != "" {
		MANIFESTJSON = envMANIFEST
	} else {
		fmt.Fprintf(os.Stderr, "Please specify a filename for the manifest\n\n")
		usage()
		os.Exit(1)
	}

	if len(flag.Args()) != 0 {
		SEEDS = flag.Args()
	} else if envSEEDS != "" {
		SEEDS = strings.Split(envSEEDS, " ")
	} else {
		fmt.Fprintf(os.Stderr, "Please specify one or more items to download (as command line arguments)\n\n")
		usage()
		os.Exit(1)
	}

	if PATH[len(PATH)-1] != os.PathSeparator {
		PATH = PATH + string(os.PathSeparator)
	}

	if *flagHELP {
		usage()
		os.Exit(1)
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [options] (-apikey <key>) <id> [...]\n", os.Args[0])

	flag.PrintDefaults()

	fmt.Fprintf(os.Stderr, "   <id> string\n        Required: one or more Steamworks collection(s) or file ID(s) to download (space seperated) [env: WSDL_SEEDS]\n")
	fmt.Fprintf(os.Stderr, "\nPlease visit %s for source code and more information\n\n", REPOURL)
}

func main() {
	d := new(downloader.Downloader)
	m := new(manifest.Manifest)

	err := d.Start(APIKEY, SEEDS, PATH)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(2)
	}

	d.DisplayFetched()

	err = m.Open(MANIFESTJSON, PATH)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(2)
	}

	d.UpdateFiles(m)

	m.Close()

	log.Println("Completed!")
}
