package main

import (
	"sync"
	"fmt"
)

// GenHashingPipeline generates the initial channel for the pipeline.
//
// It receives the list of filenames that will be hashed, a hashing algorithm from the
// avaiabale algorithms and returns a channel with all the HashFileRequests, ready
// to be consumed by the next stage.
//
// Will return error if the the named hash is not avaiable.
func GenHashingPipeline(filenames []string, hash string) (<-chan *HashFileRequest, error) {
	if !IsAnAvaiableAlgorithm(hash) {
		return nil, fmt.Errorf(HASH_NOT_SUPPORTED_ERROR_FORMAT, hash)
	}

	outChannel := make(chan *HashFileRequest)

	go func() {
		for i, fileName := range filenames {
			outChannel <- NewHashFileRequest(fileName, i, hash)
		}

		close(outChannel)
	}()

	return outChannel, nil
}

// HashPipeline hashes all the jobs on the pipeline. Returns a channel with
// all the generated HashFileReports.
func HashPipeline(requestsChan <-chan *HashFileRequest) <-chan *HashFileReport{
	outChannel := make(chan *HashFileReport)

	go func() {
		for request := range requestsChan {
			outChannel <- request.Execute()
		}

		close(outChannel)
	}()

	return outChannel
}

// MergePipelines takes a slice of channels with all the final HashFileReports and
// merges them into a single channel.
func MergePipelines(channels []<-chan *HashFileReport) <-chan *HashFileReport {
	var wg sync.WaitGroup
	outChannel := make(chan *HashFileReport)

	output := func(channel <-chan *HashFileReport) {
		for report := range channel {
			outChannel <- report
		}
		wg.Done()
	}

	wg.Add(len(channels))
	for _, channel := range channels {
		go output(channel)
	}

	go func() {
		wg.Wait()
		close(outChannel)
	}()

	return outChannel
}
