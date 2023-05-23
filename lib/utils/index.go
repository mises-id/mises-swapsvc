package utils

import (
	"io"
	"io/ioutil"
	"math/bits"
	"os"
	"path"
	"strconv"
	"strings"
)

func FindStringArrayValueIndex(header []string, value string) int {
	index := -1
	for k, v := range header {
		if v == value {
			index = k
			break
		}
	}
	return index
}

func WirteLogAppend(path string, content string) error {

	arr := strings.Split(path, "/")
	filePath := strings.Join(arr[:len(arr)-1], "/")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		err := os.MkdirAll(filePath, os.ModePerm)
		if err != nil {
			return err
		}
	}
	fileObj, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer fileObj.Close()
	if _, err := io.WriteString(fileObj, content); err == nil {
		return err
	}
	return nil
}

func IntMin(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func SaveLocal(filePath, filename string, file io.Reader) (string, error) {
	var err error
	if _, err = os.Stat(filePath); os.IsNotExist(err) {
		err := os.MkdirAll(filePath, os.ModePerm)
		if err != nil {
			return "", err
		}
	}
	localFile := path.Join(filePath, filename)
	dst, err := os.Create(localFile)
	if err != nil {
		return "", err
	}
	defer dst.Close()
	if _, err = io.Copy(dst, file); err != nil {
		return "", err
	}
	return localFile, nil
}

func GetFolders(path string) (filenames []string, err error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return filenames, err
	}
	for _, file := range files {
		if !file.IsDir() {
			filenames = append(filenames, file.Name())
		}
	}
	return filenames, nil
}

func HammingDistanceString(hash1 string, hash2 string) (int, error) {
	h1, _ := strconv.ParseUint(hash1, 2, 64)
	h2, _ := strconv.ParseUint(hash2, 2, 64)
	return HammingDistance(h1, h2)
}

func HammingDistance(hash1 uint64, hash2 uint64) (int, error) {
	hamming := hash1 ^ hash2
	return popcnt(hamming), nil
}

func popcnt(x uint64) int { return bits.OnesCount64(x) }

func Convert2binary(n uint64) string {
	result := ""
	for ; n > 0; n /= 2 {
		lsb := n % 2
		result = strconv.Itoa(int(lsb)) + result
	}
	return result
}

func Hex2Dec(val string) uint64 {
	if val == "" {
		return 0
	}
	if strings.HasPrefix(val, "0x") {
		val = strings.TrimPrefix(val, "0x")
	}
	n, _ := strconv.ParseUint(val, 16, 64)
	return n
}
