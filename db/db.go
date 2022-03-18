// forked from https://github.com/sdomino/scribble and adapted

package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var (
	ErrMissingKey    = errors.New("missing key - unable to save record")
	ErrMissingBucket = errors.New("missing bucket - no place to save record")
)

type (
	// Repo is what is used to interact with the scribble database. It runs
	// transactions, and provides log output
	Repo struct {
		mutex   sync.Mutex
		mutexes map[string]*sync.Mutex
		dir     string // the directory where scribble will create the database
	}
)

// Options uses for specification of working golang-scribble
type Options struct{}

// New creates a new scribble database at the desired directory location, and
// returns a *Repo to then use for interacting with the database
func New(dir string, options *Options) (*Repo, error) {
	dir = filepath.Clean(dir)
	repo := Repo{
		dir:     dir,
		mutexes: make(map[string]*sync.Mutex),
	}

	// if the database already exists, just use it
	if _, err := os.Stat(dir); err == nil {
		return &repo, nil
	}

	// if the database doesn't exist create it
	return &repo, os.MkdirAll(dir, 0755)
}

// Write locks the database and attempts to write the record to the database under
// the [bucket] specified with the [key] name given
func (d *Repo) Write(bucket, key string, v interface{}) error {
	// ensure there is a place to save record
	if bucket == "" {
		return ErrMissingBucket
	}

	// ensure there is a key (name) to save record as
	if key == "" {
		return ErrMissingKey
	}

	mutex := d.getOrCreateMutex(bucket)
	mutex.Lock()
	defer mutex.Unlock()

	dir := filepath.Join(d.dir, bucket)
	fnlPath := filepath.Join(dir, key+".json")
	tmpPath := fnlPath + ".tmp"

	return write(dir, tmpPath, fnlPath, v)
}

func write(dir, tmpPath, dstPath string, v interface{}) error {
	// create bucket directory
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// marshal the pointer to a non-struct and indent with tab
	b, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		return err
	}

	// write marshaled data to the temp file
	if err := ioutil.WriteFile(tmpPath, b, 0644); err != nil {
		return err
	}

	// move final file into place
	return os.Rename(tmpPath, dstPath)
}

// Read a record from the database
func (d *Repo) Read(bucket, key string, v interface{}) error {
	// ensure there is a place to save record
	if bucket == "" {
		return ErrMissingBucket
	}

	// ensure there is a key (name) to save record as
	if key == "" {
		return ErrMissingKey
	}

	record := filepath.Join(d.dir, bucket, key)
	// read record from database; if the file doesn't exist `read` will return an err
	return read(record, v)
}

func read(record string, v interface{}) error {
	b, err := ioutil.ReadFile(record + ".json")
	if err != nil {
		return err
	}

	// unmarshal data
	return json.Unmarshal(b, &v)
}

// ReadAll records from a bucket; this is returned as a slice of strings because
// there is no way of knowing what type the record is.
func (d *Repo) ReadAll(bucket string) (map[string][][]byte, error) {
	// ensure there is a bucket to read
	if bucket == "" {
		return nil, ErrMissingBucket
	}

	dir := filepath.Join(d.dir, bucket)
	// read all the files in the transaction.Bucket; an error here just means
	// the bucket is either empty or doesn't exist
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	return readAll(files, dir)
}

func readAll(files []os.FileInfo, dir string) (map[string][][]byte, error) {
	// the files read from the database
	records := make(map[string][][]byte)

	// iterate over each of the files, attempting to read the file. If successful
	// append the files to the bucket of read
	for _, file := range files {
		fileName := strings.Split(file.Name(), ".")[0]

		records[fileName] = [][]byte{}
		b, err := ioutil.ReadFile(filepath.Join(dir, file.Name()))

		if err != nil {
			return nil, err
		}

		// append read file
		records[fileName] = append(records[fileName], b)
	}
	// unmarhsal the read files as a comma delimeted byte array
	return records, nil
}

// Delete locks the database then attempts to remove the bucket/key
// specified by [path]
func (d *Repo) Delete(bucket, key string) error {
	path := filepath.Join(bucket, key)
	//
	mutex := d.getOrCreateMutex(bucket)
	mutex.Lock()
	defer mutex.Unlock()

	dir := filepath.Join(d.dir, path)
	switch fi, err := stat(dir); {
	// if fi is nil or error is not nil return
	case fi == nil, err != nil:
		return fmt.Errorf("Unable to find file or directory named %v\n", path)
	// remove directory and all contents
	case fi.Mode().IsDir():
		return os.RemoveAll(dir)
	// remove file
	case fi.Mode().IsRegular():
		return os.RemoveAll(dir + ".json")
	}
	return nil
}

//
func stat(path string) (fi os.FileInfo, err error) {
	// check for dir, if path isn't a directory check to see if it's a file
	if fi, err = os.Stat(path); os.IsNotExist(err) {
		fi, err = os.Stat(path + ".json")
	}
	return
}

// getOrCreateMutex creates a new bucket specific mutex any time a bucket
// is being modified to avoid unsafe operations
func (d *Repo) getOrCreateMutex(bucket string) *sync.Mutex {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	m, ok := d.mutexes[bucket]
	// if the mutex doesn't exist make it
	if !ok {
		m = &sync.Mutex{}
		d.mutexes[bucket] = m
	}
	return m
}
