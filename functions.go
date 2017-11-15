package main

import(
	"hash"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/sha1"
	"crypto/md5"
	"hash/crc32"
	"hash/crc64"
	"hash/adler32"
	"hash/fnv"
	
	"io"
	"fmt"
)

// Simple fucntion to get a hash from a certain io.Reader of data
func getHash(hasherString string, input io.Reader) ([]byte, error) {
	var hasher hash.Hash
	
	switch hasherString {
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
	}
	
	_, err := io.Copy(hasher, input)
	if err != nil {
		return nil, err
	}
	
	return hasher.Sum(nil), nil
}

// Concurrent version of getHash, making use of apporpiate structs
func goGetHash(request *HashRequest, reportChan chan<- *HashReport) {
	report := &HashReport{HashRequest: request}
	
	report.Sum, report.Err = getHash(request.HasherString, request.Input)
	
	reportChan <- report
}

// Miscelaneous

// Is the string in the slice?
func stringInSlice(str string, slice []string) bool {
    for _, s := range slice {
        if s == str {
            return true
        }
    }
    return false
}

// Print a report, with a sum on the specified format, to the specified output device
func printReport(report *HashReport, reportFormat string, outDevice io.Writer) {
	fmt.Fprintf(outDevice, "%s %s\n", report.Report(reportFormat), report.Name)
}
