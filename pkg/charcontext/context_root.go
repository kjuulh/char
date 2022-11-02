package charcontext

import (
	"context"
	"errors"
	"os"
	"path"
)

var ErrNoContextFound = errors.New("could not find project root")

const CharFileName = ".char.yml"

func FindLocalRoot(ctx context.Context) (string, error) {
	curdir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return recursiveFindLocalRoot(ctx, curdir)

	//output, err := exec.Command("git", "rev-parse", "--show-toplevel").CombinedOutput()
	//if err != nil {
	//	return "", err
	//}
	//if len(output) == 0 {
	//	return "", errors.New("could not find absolute path")
	//}
	//if _, err := os.Stat(string(output)); errors.Is(err, os.ErrNotExist) {
	//	return "", fmt.Errorf("path does not exist %s", string(output))
	//}

	//return string(output), nil
}

func recursiveFindLocalRoot(ctx context.Context, localpath string) (string, error) {
	entries, err := os.ReadDir(localpath)
	if err != nil {
		return "", err
	}

	for _, entry := range entries {
		if entry.Name() == CharFileName {
			return localpath, nil
		}
	}

	if localpath == "/" {
		return "", ErrNoContextFound
	}

	return recursiveFindLocalRoot(ctx, path.Dir(localpath))
}

func ChangeToPath(_ context.Context, path string) error {
	return os.Chdir(path)
}
