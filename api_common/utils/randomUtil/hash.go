package randomUtil

import (
	"hash/fnv"
	"strconv"
)

/*
Hash generates a hash of a string

ex.
"hello" -> 1542292725
*/
func Hash(s string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(s))
	return h.Sum64()
}

/*
generateHashString generates a hash of a string and returns the last 4 characters of the hash

ex.
"hello" -> "2725"
*/
func generateHashString(input string) (fullHash string) {
	fullHash = strconv.FormatUint(Hash(input), 16)
	return fullHash[len(fullHash)-4:]
}
