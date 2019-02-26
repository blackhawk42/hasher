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
	DEFAULT_USE_STDIN bool = false
	DEFAULT_SORTING_MODE bool = false
	DEFAULT_HASH_ALGORITHM string = "crc32"
	DEFAULT_REPORT_IN_UPPER bool = false
	DEFAULT_WORKERS int = 0
)

const (
	HASH_NOT_SUPPORTED_ERROR_FORMAT = "hash algorithm not supported: %s"
)

// Flag config

var useStdin = flag.Bool("stdin", DEFAULT_USE_STDIN, "Use stdin for data input. Will calculate one hash.")
var sortingMode = flag.Bool("sort", DEFAULT_SORTING_MODE, "Sort results in order of passed files, in contrast to the inherent randomness of concurrency. May delay printing of results. Final results are similar to specifying a single worker, but sorting mode will do the processing concurrently using all specified workers and then sort.")
var algorithm = flag.String("hash", DEFAULT_HASH_ALGORITHM, "Hash `algorithm` to use from the avaiable listed.")
var upper = flag.Bool("U", DEFAULT_REPORT_IN_UPPER, "Report hashes in uppercase, instead of lowercase letters.")
var workers = flag.Int("workers", DEFAULT_WORKERS, "The `number` of workers to use for cuncurrency. if <= 0, workers default to number of files to process. A single (1) worker (basically-but-not-exactly sequential mode) can be of help with I/O bottlenecks, and like sorting mode results are printed in the given order, without the overhead of actual sorting (if that flag is not used), but without the benefits of concurrent processing.")

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "use: %[1]s [OPTIONS] FILE1 [FILE2...]\nuse: %[1]s -stdin [OPTIONS]\n\n", filepath.Base(os.Args[0]) )
		fmt.Printf("Concurrently calculate and print many hashes, mostly from Go's standard library.\n\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nAvaiable hash algorithms: %s\n", strings.Join(GetAvaiableAlgorithms(), ", ") )
	}
	flag.Parse()

	// Set up worker length
	if *workers <= 0 {
		*workers = len(flag.Args())
	}

	// Other exit points

	// If no args and not using stdin, let's consider it another (sucessful) way to ask for help
	if len(flag.Args()) == 0 && !*useStdin {
		flag.Usage()
		os.Exit(0)
	}

	// Preemtively check against not valid algorithms
	if !IsAnAvaiableAlgorithm(*algorithm) {
		fmt.Fprintf(os.Stderr, HASH_NOT_SUPPORTED_ERROR_FORMAT + "\n", *algorithm)
		flag.Usage()
		os.Exit(1)
	}

	// Main logic

	// The special case of wanting to hash stdin
	if *useStdin {
		sum, err := HashReader(*algorithm, os.Stdin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		} else {
			if *upper {
				fmt.Printf("%X\n", sum)
			} else {
				fmt.Printf("%x\n", sum)
			}

			os.Exit(0)
		}
	}

	// In theory, we've already dealt with the possibility of the algorithm not
	// being supported
	requestsChan, _ := GenHashingPipeline(flag.Args(), *algorithm)

	workersSlice := make([]<-chan *HashFileReport, *workers)
	for i := range workersSlice {
		workersSlice[i] = HashPipeline(requestsChan)
	}

	reportsChan := MergePipelines(workersSlice)

	if *sortingMode {
		reportsSlice := make( HashFileReportSlice, len(flag.Args()) )
		i := 0
		for report := range reportsChan {
			reportsSlice[i] = report
			i++
		}

		sort.Sort(reportsSlice)

		reportsSlice.PrintAllReports(*upper)

	} else { // We don't bother to create slices for non-sorting mode
		for report := range reportsChan {
			fmt.Printf("%s\n", report.Report(*upper))
		}
	}

}
