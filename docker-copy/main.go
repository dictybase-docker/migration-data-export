package main

import (
	"archive/tar"
	"log"

	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/memfs"
	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/urfave/cli/v2"
	"golang.org/x/exp/slices"
)

const ExitError = 2

type extracDataFromContainerProperties struct {
	container   types.Container
	client      *client.Client
	target, src string
}

func main() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "src",
				Usage:    "folder in the docker container that will be copied out",
				Required: true,
				Aliases:  []string{"s"},
			},
			&cli.StringFlag{
				Name:    "output",
				Usage:   "folder where the content will be copied",
				Aliases: []string{"o"},
			},
			&cli.StringFlag{
				Name:     "container",
				Usage:    "name of the container from where the content will be copied",
				Required: true,
				Aliases:  []string{"c"},
			},
			&cli.StringFlag{
				Name:    "data-repo",
				Usage:   "git repository from here additional data will be copied",
				Aliases: []string{"r"},
				Value:   "https://github.com/dictyBase/migration-data.git",
			},
		},
		Action: CopyFromContainer,
	}

	if err := app.Run(os.Args); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

// CopyFromContainer is the action associated with the CLI application.
// It is responsible for copying files from a Docker container to a local directory.
// It initializes a new Docker client, retrieves the specified container,
// and calls the extracDataFromContainer function to perform the file extraction.
func CopyFromContainer(cltx *cli.Context) error {
	client, err := client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)
	if err != nil {
		return cli.Exit(
			fmt.Sprintf("error in connecting to docker host %s", err),
			ExitError,
		)
	}
	defer client.Close()
	cont, err := getContainer(cltx.String("container"), client)
	if err != nil {
		return cli.Exit(err.Error(), ExitError)
	}
	target := cltx.String("output")
	err = extracDataFromContainer(&extracDataFromContainerProperties{
		client:    client,
		container: cont,
		src:       normalizeSrc(cltx.String("src")),
		target:    target,
	})
	if err != nil {
		return cli.Exit(err.Error(), ExitError)
	}
	fs, err := cloneDataRepo(cltx.String("data-repo"))
	if err != nil {
		return cli.Exit(err.Error(), ExitError)
	}
	err = CopyFile(
		"import/data/stockcenter/gwdi_strain.csv",
		filepath.Join(
			target,
			cltx.String("src"),
			"stockcenter/gwdi_strain.csv",
		),
		fs,
	)
	if err != nil {
		return cli.Exit(err.Error(), ExitError)
	}

	stckTarget := filepath.Join(target, cltx.String("src"), "stockcenter")
	for _, folder := range []string{"formatted_sequence", "images", "raw_sequence"} {
		err = copyFolder(
			filepath.Join("plasmid", folder),
			filepath.Join(stckTarget, folder),
			fs,
		)
		if err != nil {
			return cli.Exit(err.Error(), ExitError)
		}
	}

	return nil
}

func CopyFile(src, dst string, blfs billy.Filesystem) error {
	// Open the source file
	srcFile, err := blfs.Open(src)
	if err != nil {
		return fmt.Errorf("error in opening source file %s", err)
	}
	defer srcFile.Close()

	// Create the destination file
	dstFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("error in creating destination file %s", err)
	}
	defer dstFile.Close()

	// Copy the contents from source to destination
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return fmt.Errorf("error in copying file %s", err)
	}

	return nil
}

func copyFolder(src, dest string, blfs billy.Filesystem) error {
	// Read the source directory
	srcDirInfos, err := blfs.ReadDir(src)
	if err != nil {
		return fmt.Errorf("error in opening source dir %s", err)
	}
	// Create destination directory
	if err := os.MkdirAll(dest, 0755); err != nil {
		return fmt.Errorf("error in creating destination folder %s", err)
	}
	// Copy each file/folder to the destination directory
	for _, fileInfo := range srcDirInfos {
		srcPath := filepath.Join(src, fileInfo.Name())
		destPath := filepath.Join(dest, fileInfo.Name())
		if !fileInfo.IsDir() {
			if err := CopyFile(srcPath, destPath, blfs); err != nil {
				return fmt.Errorf("Failed to copy file: %s", srcPath)
			}
		} else {
			if err := copyFolder(srcPath, destPath, blfs); err != nil {
				return err
			}
		}
	}

	return nil
}

// getContainer takes a container name and a Docker client as parameters
// and returns the corresponding container object.
// It retrieves a list of all containers and searches for a container with a matching name.
func getContainer(name string, client *client.Client) (types.Container, error) {
	containers, err := client.ContainerList(
		context.Background(),
		types.ContainerListOptions{All: true},
	)

	if err != nil {
		return types.Container{}, fmt.Errorf(
			"error in getting list of container %s",
			err,
		)
	}
	idx := slices.IndexFunc(containers, func(cont types.Container) bool {
		return strings.Contains(cont.Names[0], name)
	})
	if idx == -1 {
		return types.Container{}, fmt.Errorf(
			"unable to find container %s",
			"dsc-orders",
		)
	}

	return containers[idx], nil
}

// extracDataFromContainer takes an extracDataFromContainerProperties struct as a parameter
// and performs the extraction of files from the container.
// It removes the target folder, copies files from the container to the target folder,
// and creates the necessary directories and files.
func extracDataFromContainer(args *extracDataFromContainerProperties) error {
	target := args.target
	if _, err := os.Stat(target); os.IsNotExist(err) {
		return fmt.Errorf(
			"given folder %s does not exist, has to be created before extracing data",
			target,
		)
	}
	readCloser, _, err := args.client.CopyFromContainer(
		context.Background(),
		args.container.ID,
		args.src,
	)
	if err != nil {
		return fmt.Errorf("error in copying from container %s", err)
	}
	defer readCloser.Close()
	trd := tar.NewReader(readCloser)
	for {
		hdr, err := trd.Next()
		if err == io.EOF {
			break
		}
		switch hdr.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(filepath.Join(target, hdr.Name), 0750); err != nil {
				return fmt.Errorf("error in making dir %s %s\n", hdr.Name, err)
			}
		case tar.TypeReg:
			owr, err := os.Create(filepath.Join(target, hdr.Name))
			if err != nil {
				return fmt.Errorf(
					"error in creating file %s %s\n",
					err,
					hdr.Name,
				)
			}
			defer owr.Close()
			if _, err := io.Copy(owr, trd); err != nil {
				return fmt.Errorf(
					"error in writing to file %s %s",
					err,
					hdr.Name,
				)
			}
		default:
			continue
		}
	}
	return nil
}

// normalizeSrc takes a source path as a parameter and returns the normalized source path.
// It ensures that the source path ends with a forward slash ("/") by appending it if necessary.
func normalizeSrc(src string) string {
	if strings.HasSuffix(src, "/") {
		return src
	}
	return fmt.Sprintf("%s/", src)
}

// cloneDataRepo clones a git repository from the specified URL and returns
// a billy.Filesystem representing the cloned repository.
//
// Example:
//
//	fs, err := cloneDataRepo("https://github.com/example/repo.git")
//	if err != nil {
//	  log.Fatal(err)
//	}
//
// Use the cloned repository filesystem (fs) for further operations.
func cloneDataRepo(url string) (billy.Filesystem, error) {
	fs := memfs.New()
	_, err := git.Clone(memory.NewStorage(), fs, &git.CloneOptions{
		URL: url,
	})
	if err != nil {
		return fs, fmt.Errorf("error in cloning repo %s", err)
	}
	return fs, nil
}
