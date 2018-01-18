package main

import(
	"flag"
	"os"
	"fmt"
	"strings"
	"sort"
	"path/filepath"
)

const(
	DEFAULT_REPORT_CHANNEL_BUFFER int = 10
	DEFAULT_HASH_ALGORITHM string = "crc32"
	
	DEFAULT_REPORT_UPPER_FORMAT string = "%X"
	DEFAULT_REPORT_LOWER_FORMAT string = "%x"
)

var DEFAULT_OUTPUT_DEVICE *os.File = os.Stdout



var AvaiableHashes = []string{
	"sha256",
	"sha224",
	"sha512",
	"sha384",
	"sha512/224",
	"sha512/256",
	"sha1",
	"md5",
	"crc32",
	"crc64",
	"adler32",
	"fnv1-32",
	"fnv1-64",
	"fnv1-128",
	"fnv1a-32",
	"fnv1a-64",
	"fnv1a-128",
}

func main() {
	// Flag config
	
	flag.Usage = func () {
		fmt.Fprintf(os.Stderr, "use: %[1]s [OPTIONS] FILE1 [FILE2...]\n%[1]s [OPTIONS] -stdin\n\n", filepath.Base(os.Args[0]) )
		fmt.Printf("Concurrently calculate and print many hashes, mostly from Go's standard library.\n\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nAvaiable hash algorithms: %s\n", strings.Join(AvaiableHashes, ", ") )
	}
	
	var useStdin = flag.Bool("stdin", false, "use stdin for data input. Will calculate one hash")
	var sortingMode = flag.Bool("sort", false, "sort results in order of passed files, in contrast to the inherent randomness of concurrency. May delay printing of results")
	var reportChannelBufferSize = flag.Int("b", DEFAULT_REPORT_CHANNEL_BUFFER, "`buffer size` of the channel to store results")
	var algorithm = flag.String("hash", DEFAULT_HASH_ALGORITHM, "hash `algorithm` to use from the avaiable listed")
	var upper = flag.Bool("U", false, "report hashes in uppercase, instead of lowercase letters")
	var iterativeMode = flag.Bool("i", false, "iterative mode. No concurrency will be used. Useful for poor CPU or memory, to mitigitate I/O bottlenecks, etc.")
	var outputFile = flag.String("output", "", "output `file`. If not specified, all is printed to stdout")
	flag.Parse()
	
	// Make sure hash is avaiable
	
	if !stringInSlice(*algorithm, AvaiableHashes) {
		fmt.Fprintf(os.Stderr, "hash not avaiable: %s\n", *algorithm)
		flag.Usage()
		os.Exit(1)
	}
	
	// Set up output device
	
	var outputDevice *os.File = DEFAULT_OUTPUT_DEVICE
	
	if *outputFile != "" {
		var err error
		outputDevice, err = os.Create(*outputFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "opening output file %s: %v\n", *outputFile, err)
			os.Exit(1)
		}
		defer outputDevice.Close()
	}
	
	// Format config
	
	var reportFormat string
	if *upper {
		reportFormat = DEFAULT_REPORT_UPPER_FORMAT
	} else {
		reportFormat = DEFAULT_REPORT_LOWER_FORMAT
	}
	
	// Other exit points
	
	if len(flag.Args()) == 0 && !*useStdin {
		flag.Usage()
		os.Exit(0)
	}
	
	// Main logic
	
	// Number of reports actualy generated
	var currentNumber int = 0
	
	if *useStdin {
		hash, err := getHash(*algorithm, os.Stdin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "getting hash: %v\n", err)
			os.Exit(1)
		}
		
		fmt.Printf(reportFormat + "\n", hash)
		
	} else if *iterativeMode {
		
		// An "artificial" report. It facilitates reporting and provides
		// possibilities for debugging and extension
		currentHash := &HashReport{HashRequest: &HashRequest{} }
		
		for _, file := range flag.Args() {
			f, err := os.Open(file)
			if err != nil {
				fmt.Fprintf(os.Stderr, "opening %s: %v\n", file, err)
				continue
			}
			defer f.Close()
			
			currentNumber++
			
			currentHash.Input = f
			currentHash.Name = file
			currentHash.Number = currentNumber
			
			currentHash.Sum, currentHash.Err = getHash(*algorithm, f)
			printReport(currentHash, reportFormat, outputDevice)
		}
		
	} else {
		
		reportChan := make(chan *HashReport, *reportChannelBufferSize)
		
		for _, file := range flag.Args() {
			f, err := os.Open(file)
			if err != nil {
				fmt.Fprintf(os.Stderr, "opening %s: %v\n", file, err)
				continue
			}
			defer f.Close()
			
			currentNumber++
			
			go goGetHash(&HashRequest{HasherString: *algorithm, Input: f, Number: currentNumber, Name: file}, reportChan)
			
		}
		
		if *sortingMode {
			reports := HashReportSlice( make([]*HashReport, currentNumber ) )
			
			for i := 0; i < currentNumber; i++ {
				reports[i] = <-reportChan
			}
			
			sort.Sort(reports)
			
			for _, report := range reports {
				printReport(report, reportFormat, outputDevice)
			}
			
		} else {
			for i := 0; i < currentNumber; i++ {
				report := <- reportChan
				
				printReport(report, reportFormat, outputDevice)
			}
		}
	}
}
