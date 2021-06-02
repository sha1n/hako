package internal

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"runtime"

	"github.com/google/go-github/v35/github"
)

const (
	owner      = "sha1n"
	repo       = "hako"
	binaryName = "hako"
)

// Release a GitHub realease facade
type Release interface {
	TagName() string
	DownloadAsset() (io.ReadCloser, error)
}

type rel struct {
	delegate *github.RepositoryRelease
}

// GetLatestRelease returns the latest github non-draft release of this program.
func GetLatestRelease() (Release, error) {
	ctx := context.Background()
	client := github.NewClient(nil)

	release, _, err := client.Repositories.GetLatestRelease(ctx, owner, repo)

	return &rel{delegate: release}, err
}

func (r *rel) TagName() string {
	return *r.delegate.TagName
}

func (r *rel) DownloadAsset() (rc io.ReadCloser, err error) {
	ctx := context.Background()
	client := github.NewClient(nil)
	var assetID int64

	if assetID, err = findCompatibleAssetID(r.delegate); err == nil {
		rc, _, err = client.Repositories.DownloadReleaseAsset(ctx, owner, repo, assetID, http.DefaultClient)
	}

	return rc, err
}

func findCompatibleAssetID(release *github.RepositoryRelease) (int64, error) {
	requiredAssetName := getRequiredAssetName()
	// log.Debugf("Required asset name is %s. Looking for matching assets in latest release.", requiredAssetName)
	for _, asset := range (*release).Assets {
		if *asset.Name == requiredAssetName {
			// log.Debugf("Found asset ID = %d", *asset.ID)
			// log.Debugf("Found asset Name = %s", *asset.Name)
			return *asset.ID, nil
		}
	}
	return 0, fmt.Errorf("unable to find a compatible asset in the latest release (required=%s)", requiredAssetName)
}

func getRequiredAssetName() string {
	assertName := fmt.Sprintf("%s-%s-%s", binaryName, runtime.GOOS, runtime.GOARCH)
	if runtime.GOOS == "windows" {
		assertName += ".exe"
	}

	// log.Debugf("Required asset name is: %s", assertName)
	return assertName
}
