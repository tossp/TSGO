package ipfs

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	config "github.com/ipfs/go-ipfs-config"
	icore "github.com/ipfs/interface-go-ipfs-core"
	"github.com/tossp/tsgo/pkg/setting"
)

var (
	basePath   = setting.UseDataPath("ipfs")
	GlobalPath = UseIpfsPath("global")
	defIpfs    *Ipfs
)

type Ipfs struct {
	ipfs icore.CoreAPI
	ctx  context.Context
	Stop context.CancelFunc
}

func (fs *Ipfs) Info() {
}

func init() {
	_ = os.MkdirAll(GlobalPath, 755)
}

func UseIpfsPath(elem ...string) string {
	return filepath.Join(basePath, filepath.Clean(filepath.Join(elem...)))
}

/// ------ Spawning the node

// Spawns a node on the default repo location, if the repo exists
func spawnDefault(ctx context.Context) (icore.CoreAPI, error) {
	_, err := config.PathRoot()
	if err != nil {
		// shouldn't be possible
		return nil, err
	}

	if err := setupPlugins(basePath); err != nil {
		return nil, err

	}

	return createNode(ctx, basePath)
}

// Spawns a node to be used just for this run (i.e. creates a tmp repo)
func spawnEphemeral(ctx context.Context) (icore.CoreAPI, error) {
	if err := setupPlugins(""); err != nil {
		return nil, err
	}

	// Create a Temporary Repo
	repoPath, err := createGlobalRepo(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create temp repo: %s", err)
	}

	// Spawning an ephemeral IPFS node
	return createNode(ctx, repoPath)
}

//

/// -------

func Run() {
	/// --- Part I: Getting a IPFS node running

	fmt.Println("-- Getting an IPFS node running -- ")

	ctx, cancel := context.WithCancel(context.Background())
	//defer cancel()
	/*
		// Spawn a node using the default path (~/.ipfs), assuming that a repo exists there already
		fmt.Println("Spawning node on default repo")
		ipfs, err := spawnDefault(ctx)
		if err != nil {
			fmt.Println("No IPFS repo available on the default path")
		}
	*/

	// Spawn a node using a temporary path, creating a temporary repo for the run
	fmt.Println("Spawning node on a temporary repo")
	ipfs, err := spawnEphemeral(ctx)
	if err != nil {
		panic(fmt.Errorf("failed to spawn ephemeral node: %s", err))
	}
	defIpfs = &Ipfs{
		ipfs: ipfs,
		ctx:  ctx,
		Stop: cancel,
	}

	fmt.Println("IPFS node is running")

	///// --- Part II: Adding a file and a directory to IPFS
	//
	//fmt.Println("\n-- Adding and getting back files & directories --")
	//
	//inputBasePath := "./example-folder/"
	//inputPathFile := inputBasePath + "ipfs.paper.draft3.pdf"
	//inputPathDirectory := inputBasePath + "test-dir"
	//
	//someFile, err := getUnixfsNode(inputPathFile)
	//if err != nil {
	//    panic(fmt.Errorf("Could not get File: %s", err))
	//}
	//
	//cidFile, err := ipfs.Unixfs().Add(ctx, someFile)
	//if err != nil {
	//    panic(fmt.Errorf("Could not add File: %s", err))
	//}
	//
	//fmt.Printf("Added file to IPFS with CID %s\n", cidFile.String())
	//
	//someDirectory, err := getUnixfsNode(inputPathDirectory)
	//if err != nil {
	//    panic(fmt.Errorf("Could not get File: %s", err))
	//}
	//
	//cidDirectory, err := ipfs.Unixfs().Add(ctx, someDirectory)
	//if err != nil {
	//    panic(fmt.Errorf("Could not add Directory: %s", err))
	//}
	//
	//fmt.Printf("Added directory to IPFS with CID %s\n", cidDirectory.String())
	//
	///// --- Part III: Getting the file and directory you added back
	//
	//outputBasePath := "./example-folder/"
	//outputPathFile := outputBasePath + strings.Split(cidFile.String(), "/")[2]
	//outputPathDirectory := outputBasePath + strings.Split(cidDirectory.String(), "/")[2]
	//
	//rootNodeFile, err := ipfs.Unixfs().Get(ctx, cidFile)
	//if err != nil {
	//    panic(fmt.Errorf("Could not get file with CID: %s", err))
	//}
	//
	//err = files.WriteTo(rootNodeFile, outputPathFile)
	//if err != nil {
	//    panic(fmt.Errorf("Could not write out the fetched CID: %s", err))
	//}
	//
	//fmt.Printf("Got file back from IPFS (IPFS path: %s) and wrote it to %s\n", cidFile.String(), outputPathFile)
	//
	//rootNodeDirectory, err := ipfs.Unixfs().Get(ctx, cidDirectory)
	//if err != nil {
	//    panic(fmt.Errorf("Could not get file with CID: %s", err))
	//}
	//
	//err = files.WriteTo(rootNodeDirectory, outputPathDirectory)
	//if err != nil {
	//    panic(fmt.Errorf("Could not write out the fetched CID: %s", err))
	//}
	//
	//fmt.Printf("Got directory back from IPFS (IPFS path: %s) and wrote it to %s\n", cidDirectory.String(), outputPathDirectory)
	//
	///// --- Part IV: Getting a file from the IPFS Network
	//
	//fmt.Println("\n-- Going to connect to a few nodes in the Network as bootstrappers --")
	//
	//bootstrapNodes := []string{
	//    // IPFS Bootstrapper nodes.
	//    "/dnsaddr/bootstrap.libp2p.io/ipfs/QmNnooDu7bfjPFoTZYxMNLWUQJyrVwtbZg5gBMjTezGAJN",
	//    "/dnsaddr/bootstrap.libp2p.io/ipfs/QmQCU2EcMqAqQPR2i9bChDtGNJchTbq5TbXJJ16u19uLTa",
	//    "/dnsaddr/bootstrap.libp2p.io/ipfs/QmbLHAnMoJPWSCR5Zhtx6BHJX9KiKNN6tpvbUcqanj75Nb",
	//    "/dnsaddr/bootstrap.libp2p.io/ipfs/QmcZf59bWwK5XFi76CZX8cbJ4BhTzzA3gU1ZjYZcYW3dwt",
	//
	//    // IPFS Cluster Pinning nodes
	//    "/ip4/138.201.67.219/tcp/4001/ipfs/QmUd6zHcbkbcs7SMxwLs48qZVX3vpcM8errYS7xEczwRMA",
	//    "/ip4/138.201.67.220/tcp/4001/ipfs/QmNSYxZAiJHeLdkBg38roksAR9So7Y5eojks1yjEcUtZ7i",
	//    "/ip4/138.201.68.74/tcp/4001/ipfs/QmdnXwLrC8p1ueiq2Qya8joNvk3TVVDAut7PrikmZwubtR",
	//    "/ip4/94.130.135.167/tcp/4001/ipfs/QmUEMvxS2e7iDrereVYc5SWPauXPyNwxcy9BXZrC1QTcHE",
	//
	//    // You can add more nodes here, for example, another IPFS node you might have running locally, mine was:
	//    // "/ip4/127.0.0.1/tcp/4010/ipfs/QmZp2fhDLxjYue2RiUvLwT9MWdnbDxam32qYFnGmxZDh5L",
	//}
	//
	//go connectToPeers(ctx, ipfs, bootstrapNodes)
	//
	//exampleCIDStr := "QmUaoioqU7bxezBQZkUcgcSyokatMY71sxsALxQmRRrHrj"
	//
	//fmt.Printf("Fetching a file from the network with CID %s\n", exampleCIDStr)
	//outputPath := outputBasePath + exampleCIDStr
	//testCID := icorepath.New(exampleCIDStr)
	//
	//rootNode, err := ipfs.Unixfs().Get(ctx, testCID)
	//if err != nil {
	//    panic(fmt.Errorf("Could not get file with CID: %s", err))
	//}
	//
	//err = files.WriteTo(rootNode, outputPath)
	//if err != nil {
	//    panic(fmt.Errorf("Could not write out the fetched CID: %s", err))
	//}
	//
	//fmt.Printf("Wrote the file to %s\n", outputPath)
	//
	//fmt.Println("\nAll done! You just finalized your first tutorial on how to use go-ipfs as a library")
}
