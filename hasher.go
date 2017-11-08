package main

import(
	"flag"
	"os"
	"fmt"
	"strings"
	"sort"
)

const(
	DEFAULT_REPORT_CHANNEL_BUFFER int = 10
	DEFAULT_HASH_ALGORITHM string = "sha256"
	DEFAULT_REPORT_UPPER_FORMAT = "%X"
	DEFAULT_REPORT_LOWER_FORMAT = "%x"
)



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
		fmt.Fprintf(os.Stderr, "use: %[1]s FILE1 [FILE2...]\n%[1]s -i\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nAvaiable hash algorithms: %s", strings.Join(AvaiableHashes, ", ") )
	}
	
	var useStdin = flag.Bool("stdin", false, "use stdin for input data")
	var sortingMode = flag.Bool("sort", false, "sort results, in contrasts to the inherent randomness of concurrency. May delay printing of results")
	var reportChannelBufferSize = flag.Int("b", DEFAULT_REPORT_CHANNEL_BUFFER, "`buffer size` of the channel to store results")
	var algorithm = flag.String("hash", DEFAULT_HASH_ALGORITHM, "`hash algorithm` to use from the avaiable listed")
	var upper = flag.Bool("U", false, "Report hashes in uppercase, instead of lowercase letters")
	flag.Parse()
	
	// Make sure hash is avaiable
	
	if !stringInSlice(*algorithm, AvaiableHashes) {
		fmt.Fprintf(os.Stderr, "hash not avaiable: %s\n", *algorithm)
		flag.Usage()
		os.Exit(1)
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
	
	if *useStdin {
		hash, err := getHash(*algorithm, os.Stdin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "getting hash: %v\n", err)
			os.Exit(1)
		}
		
		fmt.Printf("%x\n", hash)
		
	} else {
		
		reportChan := make(chan *HashReport, *reportChannelBufferSize)
		var currentNumber int = 0
		
		for _, file := range flag.Args() {
			f, err := os.Open(file)
			if err != nil {
				fmt.Fprintf(os.Stderr, "opening %s: %v\n", file, err)
				continue
			}
			defer f.Close()
			
			go goGetHash(&HashRequest{HasherString: *algorithm, Input: f, Number: currentNumber, Name: file}, reportChan)
			
			currentNumber++
		}
		
		if *sortingMode {
			reports := HashReportSlice( make([]*HashReport, len(flag.Args()) ) )
			
			for i := range flag.Args() {
				reports[i] = <-reportChan
			}
			fmt.Println(reports.Len())
			
			sort.Sort(reports)
			
			for _, report := range reports {
				fmt.Printf("%s: %s\n", report.Name, report.Report(reportFormat) )
			}
			
		} else {
			for range flag.Args() {
				report := <- reportChan
				
				fmt.Printf("%s: %s\n", report.Name, report.Report(reportFormat) )
			}
		}
	}
}
