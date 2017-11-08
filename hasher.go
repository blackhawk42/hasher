package main

import(
	"flag"
	"os"
	"fmt"
	"hash"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/sha1"
	"crypto/md5"
	"hash/crc32"
	"hash/crc64"
	"hash/adler32"
	"hash/fnv"
	"strings"
	"sort"
)

const(
	DEFAULT_REPORT_CHANNEL_BUFFER int = 10
	DEFAULT_HASH_ALGORITHM string = "sha256"
	DEFAULT_REPORT_UPPER_FORMAT = "%X"
	DEFAULT_REPORT_LOWER_FORMAT = "%x"
)



var AvaiableHashes = [...]string{
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

//~ var AvaiableHashes = [...]{
	//~ "sha256": sha256.New,
	//~ "sha224": sha256.New224,
	//~ "sha512": sha512.New,
	//~ "sha384": sha512.New384,
	//~ "sha512/224": sha512.New512_224,
	//~ "sha512/256": sha512.New512_256,
	//~ "sha1": sha1.New,
	//~ "md5": md5.New,
	//~ "crc32": hash.Hash(crc32.NewIEEE),
	//~ "crc64": crc64.New,
	//~ "adler32": adler32.New,
	//~ "fnv1-32": fnv.New32,
	//~ "fnv1-64": fnv.New64,
	//~ "fnv1-128": fnv.New128,
	//~ "fnv1a-32": fnv.New32a,
	//~ "fnv1a-64": fnv.New64a,
	//~ "fnv1a-128": fnv.New128a,
//~ }

func main() {
	// Flag config
	
	flag.Usage = func () {
		fmt.Fprintf(os.Stderr, "use: %[1]s FILE1 [FILE2...]\n%[1]s -i\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nAvaiable hash algorithms: %s", strings.Join(AvaiableHashes[:], ", ") )
	}
	
	var useStdin = flag.Bool("stdin", false, "use stdin for input data")
	var sortingMode = flag.Bool("sort", false, "sort results, in contrasts to the inherent randomness of concurrency. May delay printing of results")
	var reportChannelBufferSize = flag.Int("b", DEFAULT_REPORT_CHANNEL_BUFFER, "`buffer size` of the channel to store results")
	var algorithm = flag.String("hash", DEFAULT_HASH_ALGORITHM, "`hash algorithm` to use from the avaiable listed")
	var upper = flag.Bool("U", false, "Report hashes in uppercase, instead of lowercase letters")
	flag.Parse()
	
	// Hasher config
	var hasher hash.Hash
	switch *algorithm {
		case "sha256":
			hasher = sha256.New()
		case "sha224":
			hasher = sha256.New224()
		case "sha512":
			hasher = sha512.New()
		case "sha384":
			hasher = sha512.New384()
		case "sha512/224":
			hasher = sha512.New512_224()
		case "sha512/256":
			hasher = sha512.New512_256()
		case "sha1":
			hasher = sha1.New()
		case "md5":
			hasher = md5.New()
		case "crc32":
			hasher = crc32.NewIEEE()
		case "crc64":
			hasher = crc64.New(crc64.MakeTable(crc64.ISO))
		case "adler32":
			hasher = adler32.New()
		case "fnv1-32":
			hasher = fnv.New32()
		case "fnv1-64":
			hasher = fnv.New64()
		case "fnv1-128":
			hasher = fnv.New128()
		case "fnv1a-32":
			hasher = fnv.New32a()
		case "fnv1a-64":
			hasher = fnv.New64a()
		case "fnv1a-128":
			hasher = fnv.New128a()
		default:
			fmt.Fprintf(os.Stderr, "hash not found: %s\n", *algorithm)
			os.Exit(1)
	}
	
	// Format config
	
	var reportFormat string
	if *upper {
		reportFormat = DEFAULT_REPORT_UPPER_FORMAT
	} else {
		reportFormat = DEFAULT_REPORT_LOWER_FORMAT
	}
	
	// Main logic
	
	if *useStdin {
		hash, err := getHash(hasher, os.Stdin)
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
			
			go goGetHash(&HashRequest{Hasher: hasher, Input: f, Number: currentNumber, Name: file}, reportChan)
			
			currentNumber++
		}
		
		if *sortingMode {
			reports := HashReportSlice( make([]*HashReport, *reportChannelBufferSize) )
			
			for range flag.Args() {
				reports = append(reports, <-reportChan)
			}
			
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
