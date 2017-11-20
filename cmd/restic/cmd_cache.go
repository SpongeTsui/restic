package main

import (
	"path/filepath"

	"github.com/restic/restic/internal/cache"
	"github.com/restic/restic/internal/errors"
	"github.com/restic/restic/internal/fs"
	"github.com/spf13/cobra"
)

var cacheCmd = &cobra.Command{
	Use:   "cache",
	Short: "Cleanup the cache",
	Long: `
The "cache" cleans up old cache directories.
`,
	DisableAutoGenTag: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runCache(cacheOptions, globalOptions)
	},
}

// CacheOptions collects all options for the unlock command.
type CacheOptions struct {
	Cleanup bool
}

var cacheOptions CacheOptions

func init() {
	cmdRoot.AddCommand(cacheCmd)

	cacheCmd.Flags().BoolVar(&cacheOptions.Cleanup, "cleanup", false, "remove old cache directories")
}

func runCache(opts CacheOptions, gopts GlobalOptions) (err error) {
	if gopts.NoCache {
		return errors.Fatal("cannot clean up the cache without a cache")
	}

	if !opts.Cleanup {
		return errors.Fatal("no --cleanup passed, doing nothing")
	}

	cachedir := gopts.CacheDir
	if cachedir == "" {
		cachedir, err = cache.DefaultDir()
		if err != nil {
			return err
		}
	}

	oldCacheDirs, err := cache.Old(cachedir)
	if err != nil {
		return err
	}

	if len(oldCacheDirs) == 0 {
		Printf("no old cache dirs found, nothing to do\n")
		return nil
	}

	Printf("removing %d old cache dirs from %v\n", len(oldCacheDirs), cachedir)

	var foundErr bool
	for _, item := range oldCacheDirs {
		dir := filepath.Join(cachedir, item)
		err = fs.RemoveAll(dir)
		if err != nil {
			Warnf("unable to remove %v: %v\n", dir, err)
		}
	}

	if foundErr {
		return errors.Fatal("error(s) found while removing old cache directories")
	}

	return nil
}
