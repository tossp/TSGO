package ipfs

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"

	config "github.com/ipfs/go-ipfs-config"
	"github.com/ipfs/go-ipfs/repo/fsrepo"
	icore "github.com/ipfs/interface-go-ipfs-core"
)

//func createTempRepo(ctx context.Context) (string, error) {
//    repoPath, err := ioutil.TempDir("", "ipfs-shell")
//    if err != nil {
//        return "", fmt.Errorf("failed to get temp dir: %s", err)
//    }
//
//    // Create a config with default options and a 2048 bit key
//    cfg, err := config.Init(ioutil.Discard, 2048)
//    if err != nil {
//        return "", err
//    }
//
//    // Create the repo with the config
//    err = fsrepo.Init(repoPath, cfg)
//    if err != nil {
//        return "", fmt.Errorf("failed to init ephemeral node: %s", err)
//    }
//
//    return repoPath, nil
//}

func createGlobalRepo(ctx context.Context) (string, error) {
	return GlobalPath, createRepo(ctx, GlobalPath)
}

func createRepo(ctx context.Context, repoPath string) error {
	// Create a config with default options and a 2048 bit key
	cfg, err := config.Init(ioutil.Discard, 2048)
	if err != nil {
		return err
	}

	// Create the repo with the config
	err = fsrepo.Init(repoPath, cfg)
	if err != nil {
		return fmt.Errorf("failed to init ephemeral node: %s", err)
	}

	return nil
}

func loadRepo(ctx context.Context, repoPath string) (ipfs icore.CoreAPI, err error) {
	if err = os.MkdirAll(repoPath, 0700); err != nil {
		return
	}
	if err = setupPlugins(repoPath); err != nil {
		return nil, err
	}
	return createNode(ctx, repoPath)
}
