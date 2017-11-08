package main

import(
	"fmt"
	"io"
)

// Represents a request for hash generation
type HashRequest struct {
	// Hash to use, represented as a string. The function will deal with
	// actually creting the hash-generating interface
	HasherString string
	
	// Input to get data to hash
	Input io.Reader
	
	// Name to indentify this input (e.g., filename)
	Name string
	
	// Number of request
	Number int
}

// Represents a post-hashing report
type HashReport struct {
	// The request used to generate this report. Note that this includes Number,
	// useful for comparing, sorting, etc.
	*HashRequest
	
	// Hash sum
	Sum []byte
	
	// Error during hashing, if any
	Err error
}

// Utility method to simplify reporting of hashing results.
func (r *HashReport) Report(format string) string {
	if r.Err != nil {
		return fmt.Sprintf("%v", r.Err)
	} else {
		return fmt.Sprintf(format, r.Sum)
	}
}

// Slice of hash reports. Implements sort.Interface
type HashReportSlice []*HashReport

func (s HashReportSlice) Len() int {
	return len(s)
}

func (s HashReportSlice) Less(i, j int) bool {
	if s[i].Number < s[j].Number {
		return true
	} else {
		return false
	}
}

func (s HashReportSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
