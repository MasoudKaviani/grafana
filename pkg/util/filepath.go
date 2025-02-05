package util

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// ErrWalkSkipDir is the Error returned when we want to skip descending into a directory
var ErrWalkSkipDir = errors.New("skip this directory")

// WalkFunc is a callback function called for each path as a directory is walked
// If resolvedPath != "", then we are following symbolic links.
type WalkFunc func(resolvedPath string, info os.FileInfo, err error) error

// Walk walks a path, optionally following symbolic links, and for each path,
// it calls the walkFn passed.
//
// It is similar to filepath.Walk, except that it supports symbolic links and
// can detect infinite loops while following sym links.
// It solves the issue where your WalkFunc needs a path relative to the symbolic link
// (resolving links within walkfunc loses the path to the symbolic link for each traversal).
func Walk(path string, followSymlinks bool, detectSymlinkInfiniteLoop bool, followDistFolder bool, walkFn WalkFunc) error {
	info, err := os.Lstat(path)
	if err != nil {
		return err
	}
	var symlinkPathsFollowed map[string]bool
	var resolvedPath string
	if followSymlinks {
		resolvedPath = path
		if detectSymlinkInfiniteLoop {
			symlinkPathsFollowed = make(map[string]bool, 8)
		}
	}
	return walk(path, info, resolvedPath, symlinkPathsFollowed, followDistFolder, walkFn)
}

// walk walks the path. It is a helper/sibling function to Walk.
// It takes a resolvedPath into consideration. This way, paths being walked are
// always relative to the path argument, even if symbolic links were resolved).
//
// If resolvedPath is "", then we are not following symbolic links.
// If symlinkPathsFollowed is not nil, then we need to detect infinite loop.
func walk(path string, info os.FileInfo, resolvedPath string, symlinkPathsFollowed map[string]bool, followDistFolder bool, walkFn WalkFunc) error {
	if info == nil {
		return errors.New("walk: Nil FileInfo passed")
	}
	err := walkFn(resolvedPath, info, nil)
	if err != nil {
		if info.IsDir() && errors.Is(err, ErrWalkSkipDir) {
			err = nil
		}
		return err
	}

	if resolvedPath != "" && info.Mode()&os.ModeSymlink == os.ModeSymlink {
		// We only want to lstat on directories. If this entry is a symbolic link to a file, no need to recurse.
		statInfo, err := os.Stat(resolvedPath)
		if err != nil {
			return err
		}
		if !statInfo.IsDir() {
			return nil
		}

		path2, err := filepath.EvalSymlinks(resolvedPath)
		if err != nil {
			return err
		}
		// vout("SymLink Path: %v, links to: %v", resolvedPath, path2)
		if symlinkPathsFollowed != nil {
			if _, ok := symlinkPathsFollowed[path2]; ok {
				errMsg := "potential symLink infinite loop, path: %v, link to: %v"
				return fmt.Errorf(errMsg, resolvedPath, path2)
			}
			symlinkPathsFollowed[path2] = true
		}
		info2, err := os.Lstat(path2)
		if err != nil {
			return err
		}
		return walk(path, info2, path2, symlinkPathsFollowed, followDistFolder, walkFn)
	} else if info.IsDir() {
		list, err := os.ReadDir(path)
		if err != nil {
			return walkFn(resolvedPath, info, err)
		}
		var subFiles = make([]subFile, 0)
		for _, file := range list {
			path2 := filepath.Join(path, file.Name())
			var resolvedPath2 string
			if resolvedPath != "" {
				resolvedPath2 = filepath.Join(resolvedPath, file.Name())
			}
			fileInfo, err := file.Info()
			if err != nil {
				return fmt.Errorf("unable to read file info: %v, path: %v", file.Name(), path2)
			}

			subFiles = append(subFiles, subFile{path: path2, resolvedPath: resolvedPath2, fileInfo: fileInfo})
		}

		if containsDistFolder(subFiles) && followDistFolder {
			err := walk(
				filepath.Join(path, "dist"),
				info,
				filepath.Join(resolvedPath, "dist"),
				symlinkPathsFollowed,
				followDistFolder,
				walkFn)

			if err != nil {
				return err
			}
		} else {
			for _, p := range subFiles {
				err = walk(p.path, p.fileInfo, p.resolvedPath, symlinkPathsFollowed, followDistFolder, walkFn)

				if err != nil {
					return err
				}
			}
		}

		return nil
	}
	return nil
}

type subFile struct {
	path, resolvedPath string
	fileInfo           os.FileInfo
}

func containsDistFolder(subFiles []subFile) bool {
	for _, p := range subFiles {
		if p.fileInfo.IsDir() && p.fileInfo.Name() == "dist" {
			return true
		}
	}

	return false
}

// CleanRelativePath returns the shortest path name equivalent to path
// by purely lexical processing. It make sure the provided path is rooted
// and then uses filepath.Clean and filepath.Rel to make sure the path
// doesn't include any separators or elements that shouldn't be there
// like ., .., //.
func CleanRelativePath(path string) (string, error) {
	cleanPath := filepath.Clean(filepath.Join("/", path))
	rel, err := filepath.Rel("/", cleanPath)
	if err != nil {
		// slash is prepended above therefore this is not expected to fail
		return "", err
	}

	return rel, nil
}
