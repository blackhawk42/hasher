package main

import(
	"fmt"
	"os"
)

// HashFileRequest represents a request to hash a file.
type HashFileRequest struct {
	// Number of request. Useful for sorting.
	Number int

	// File name of the file to be hashed
	FileName string

	// Hash of the request
	Hash string
}

func NewHashFileRequest(fileName string, number int, hash string) *HashFileRequest {
	return &HashFileRequest{
		FileName: fileName,
		Number: number,
		Hash: hash,
	}
}

// Execute the order to hash and generate a HashFileReport based on the results.
func (request *HashFileRequest) Execute() *HashFileReport {
	f, err := os.Open(request.FileName)
	if err != nil {
		return NewHashFileReport(request, err, nil)
	}
	defer f.Close()

	sum, err := HashReader(request.Hash, f)
	if err != nil {
		return NewHashFileReport(request, err, nil)
	}

	return NewHashFileReport(request, nil, sum)
}


// HashFileReport represents a post-hashing report after executing a HashFileRequest.
type HashFileReport struct {
	// The original hash request
	*HashFileRequest

	// Any error ocurred during hashing
	Err error

	// Sum of the hashing, if succesful
	Sum []byte
}

func NewHashFileReport(request *HashFileRequest, err error, sum []byte) *HashFileReport {
	return &HashFileReport{
		HashFileRequest: request,
		Err: err,
		Sum: sum,
	}
}

// Report results as a simple printable string. Must specify if the hex part of the
// report will be in upper (true) or lowercases (false).
func (rep *HashFileReport) Report(uppercase bool) string {
	if rep.Err != nil {
		return fmt.Sprintf("error in %s: %v", rep.FileName, rep.Err)
	}


	var format string
	if uppercase {
		format = "%X %s"
	} else {
		format = "%x %s"
	}

	return fmt.Sprintf(format, rep.Sum, rep.FileName)
}

// HashFileReportSlice gives you a simple slice of HashFileReport structs. It implements
// sort.Interface.
type HashFileReportSlice []*HashFileReport

func (rs HashFileReportSlice) PrintAllReports(uppercase bool) {
	for _, report := range rs {
		fmt.Printf("%s\n", report.Report(uppercase))
	}
}

func (s HashFileReportSlice) Len() int {
	return len(s)
}

func (s HashFileReportSlice) Less(i, j int) bool {
	return s[i].Number < s[j].Number
}

func (s HashFileReportSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
