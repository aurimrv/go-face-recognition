package input

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

const (
	defaultImagesDir = "images"
	defaultInputName = "unknown.jpg"
	defaultOutputName = "result.jpg"
)

// Job is a single classify/draw/save work item.
type Job struct {
	In  string
	Out string
}

// Resolve maps CLI positional path arguments to analysis jobs.
// Empty args yield the historic default: images/unknown.jpg → images/result.jpg.
func Resolve(args []string) ([]Job, error) {
	if len(args) == 0 {
		in := filepath.Join(defaultImagesDir, defaultInputName)
		out := filepath.Join(defaultImagesDir, defaultOutputName)
		return []Job{{In: in, Out: out}}, nil
	}

	var jobs []Job
	for _, arg := range args {
		info, err := os.Stat(arg)
		if err != nil {
			if os.IsNotExist(err) {
				return nil, fmt.Errorf("path does not exist: %s", arg)
			}
			return nil, fmt.Errorf("cannot access path %s: %w", arg, err)
		}

		switch {
		case info.Mode().IsRegular():
			jobs = append(jobs, Job{In: arg, Out: resultBeside(arg)})
		case info.IsDir():
			dirJobs, err := walkDir(arg)
			if err != nil {
				return nil, err
			}
			jobs = append(jobs, dirJobs...)
		default:
			return nil, fmt.Errorf("path is neither a file nor a directory: %s", arg)
		}
	}

	if len(jobs) == 0 {
		return nil, fmt.Errorf("no eligible images found")
	}
	return jobs, nil
}

func resultBeside(in string) string {
	dir := filepath.Dir(in)
	stem := strings.TrimSuffix(filepath.Base(in), filepath.Ext(in))
	return filepath.Join(dir, "result_"+stem+".jpg")
}

func walkDir(root string) ([]Job, error) {
	var jobs []Job
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if !isJPEG(path) || isResultImage(filepath.Base(path)) {
			return nil
		}
		jobs = append(jobs, Job{In: path, Out: resultBeside(path)})
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("walk %s: %w", root, err)
	}
	return jobs, nil
}

func isJPEG(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return ext == ".jpg" || ext == ".jpeg"
}

func isResultImage(base string) bool {
	lower := strings.ToLower(base)
	if lower == "result.jpg" || lower == "result.jpeg" {
		return true
	}
	return strings.HasPrefix(lower, "result_") && isJPEG(base)
}
