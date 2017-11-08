package main

import(
	"hash"
	"io"
)

// Simple fucntion to get a hash from a certain io.Reader of data
func getHash(hasher hash.Hash, input io.Reader) ([]byte, error) {
	_, err := io.Copy(hasher, input)
	if err != nil {
		return nil, err
	}
	
	return hasher.Sum(nil), nil
}

// Concurrent version of getHash, making use of apporpiate structs
func goGetHash(request *HashRequest, reportChan chan<- *HashReport) {
	report := &HashReport{HashRequest: request}
	
	report.Sum, report.Err = getHash(request.Hasher, request.Input)
	
	reportChan <- report
}
