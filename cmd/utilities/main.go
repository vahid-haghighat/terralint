package utilities

import (
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unsafe"
)

var src = rand.NewSource(time.Now().UnixNano())

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func AbsPath(targetPath string) (string, error) {
	if strings.Contains(targetPath, "~") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		targetPath = strings.ReplaceAll(targetPath, "~", home)
	}
	return filepath.Abs(targetPath)
}

func RandString(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return *(*string)(unsafe.Pointer(&b))
}

func Exists[T comparable](name T, array []T) bool {
	for _, item := range array {
		if name == item {
			return true
		}
	}
	return false
}

// ArrayDifference Returns a - b
func ArrayDifference[T comparable](a, b []T) []T {
	mb := make(map[T]bool, len(b))
	for _, x := range b {
		mb[x] = true
	}

	var diff []T
	for _, x := range a {
		if _, found := mb[x]; !found {
			diff = append(diff, x)
		}
	}

	return diff
}

// ArrayIntersection Returns the intersection between a and b
func ArrayIntersection[T comparable](a, b []T) []T {
	mb := make(map[T]bool, len(b))
	for _, x := range b {
		mb[x] = true
	}

	var intersection []T
	for _, x := range a {
		if _, found := mb[x]; found {
			intersection = append(intersection, x)
		}
	}

	return intersection
}

func ArrayUnion[T comparable](a, b []T) []T {
	ma := make(map[T]bool, len(a))
	mb := make(map[T]bool, len(b))

	result := make(map[T]bool)

	for key, _ := range ma {
		result[key] = true
	}

	for key, _ := range mb {
		result[key] = true
	}

	return MapKeys(result)
}

// MapDifference Returns a - b
func MapDifference[K comparable, V any](a, b map[K]V) map[K]V {
	result := make(map[K]V)

	for key, value := range a {
		if _, found := b[key]; !found {
			result[key] = value
		}
	}

	return result
}

func MapKeys[K comparable, V any](input map[K]V) []K {
	result := make([]K, 0, len(input))

	for key, _ := range input {
		result = append(result, key)
	}

	return result
}

func MapValues[K comparable, V any](input map[K]V) []V {
	result := make([]V, 0, len(input))

	for _, value := range input {
		result = append(result, value)
	}

	return result
}

func MergeMaps[K comparable, V any](first map[K]V, second map[K]V) map[K]V {
	result := make(map[K]V)

	for key, value := range first {
		result[key] = value
	}

	for key, value := range second {
		result[key] = value
	}

	return result
}

func GetPointer[T any](v T) *T {
	return &v
}
