package main

import (
	"hash"
	"hash/crc32"
	"hash/crc64"
	"hash/adler32"
	"hash/fnv"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"

	"sort"

	"sync"

	"io"
	"fmt"
)

// Hashes that require individual preparation

func crc32IEEECreator() hash.Hash {
	return crc32.NewIEEE()
}

func fnv32Creator() hash.Hash {
	return fnv.New32()
}

func fnv32aCreator() hash.Hash {
	return fnv.New32a()
}

func fnv64Creator() hash.Hash {
	return fnv.New64()
}

func fnv64aCreator() hash.Hash {
	return fnv.New64a()
}

var crc64ISOTable *crc64.Table
var crc64ISOTableCreator sync.Once
func crc64ISOTableInit() {
	crc64ISOTable = crc64.MakeTable(crc64.ISO)
}

func crc64ISOCreator() hash.Hash {
	crc64ISOTableCreator.Do(crc64ISOTableInit)

	return crc64.New(crc64ISOTable)
}

func adler32Creator() hash.Hash {
	return adler32.New()
}

// Main hash map that defines what algorithms are avaiable and how they are created.
var hashStringToCreatorFuntion = map[string]func() hash.Hash{
	"crc32": crc32IEEECreator,
	"crc64-iso": crc64ISOCreator,
	"adler32": adler32Creator,
	"fnv32": fnv32Creator,
	"fnv32a": fnv32aCreator,
	"fnv64": fnv64Creator,
	"fnv64a": fnv64aCreator,
	"fnv128": fnv.New128,
	"fnv128a": fnv.New128a,
	"md5": md5.New,
	"sha1": sha1.New,
	"sha224": sha256.New224,
	"sha256": sha256.New,
	"sha384": sha512.New384,
	"sha512": sha512.New,
	"sha512/224": sha512.New512_224,
	"sha512/256": sha512.New512_256,
}

// GetAvaiableAlgorithms generates a sorted slice with all the avaiable hashes
// implemented.
func GetAvaiableAlgorithms() []string {
	algorithms := make([]string, len(hashStringToCreatorFuntion))
	i := 0

	for algorithm := range hashStringToCreatorFuntion {
		algorithms[i] = algorithm
		i++
	}

	sort.Strings(algorithms)

	return algorithms
}

// IsAnAvaiableAlgorithm checks if the named hashed is implemented.
func IsAnAvaiableAlgorithm(algorithm string) bool {
	_, ok := hashStringToCreatorFuntion[algorithm]

	return ok
}



// HashReader takes a simple io.Reader and hashes it. Mostly for internal use,
// but can be used as a stand-alone function.
func HashReader(algorithm string, reader io.Reader) ([]byte, error) {
	hasherCreator, ok := hashStringToCreatorFuntion[algorithm]
	if !ok {
		return nil, fmt.Errorf(HASH_NOT_SUPPORTED_ERROR_FORMAT, algorithm)
	}

	hasher := hasherCreator()

	_, err := io.Copy(hasher, reader)
	if err != nil {
		return nil, fmt.Errorf("hashing data: %v", err)
	}

	return hasher.Sum(nil), nil
}
