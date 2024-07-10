package rke2

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const (
	RKE2ReleaseURL = "https://github.com/rancher/rke2/releases/download/"
)

var (
	listRKE2Images = map[string]string{
		"RKE2ImagesLinux":   "rke2-images.linux-amd64.tar.zst",
		"RKE2ImagesCalico":  "rke2-images-calico.linux-amd64.tar.zst",
		"RKE2ImagesFlannel": "rke2-images-flannel.linux-amd64.tar.zst",
		"RKE2ImagesCilium":  "rke2-images-cilium.linux-amd64.tar.zst",
		"RKE2ImagesCanal":   "rke2-images-canal.linux-amd64.tar.zst",
		"RKE2ImagesMultus":  "rke2-images-multus.linux-amd64.tar.zst",
		"RKE2ImagesCore":    "rke2-images-core.linux-amd64.tar.zst",
		"RKE2Linux":         "rke2.linux-amd64.tar.gz",
		"RKE2SHA256":        "sha256sum-amd64.txt",
	}
)

type RKE2 struct {
	Version          string
	OutputDirTarball string
}

func New(version, outputDirTarball string) *RKE2 {
	return &RKE2{
		Version:          version,
		OutputDirTarball: outputDirTarball,
	}
}

func (r *RKE2) Download() error {
	// Create the destination directory if it doesn't exist
	if err := os.MkdirAll(r.OutputDirTarball, os.ModePerm); err != nil {
		log.Printf("failed to create destination directory: %v", err)
		return err
	}

	for _, image := range listRKE2Images {

		// Construct the file path
		filePath := filepath.Join(ensureTrailingSlash(r.OutputDirTarball), image)

		// Download the file
		resp, err := http.Get(RKE2ReleaseURL + replaceVersionLink(r.Version) + "/" + image)
		if err != nil {
			log.Printf("failed to download the images: %v", err)
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Printf("failed to download the images: HTTP status %s", resp.Status)
			return fmt.Errorf("failed to download the images: HTTP status %s", resp.Status)
		}

		// Create the file
		out, err := os.Create(filePath)
		if err != nil {
			log.Printf("failed to create file: %v", err)
			return err
		}
		defer out.Close()

		_, err = io.Copy(out, resp.Body)
		if err != nil {
			log.Printf("failed to save file: %v", err)
			return err
		}
		log.Printf("Image downloaded successfully to %s", filePath)
	}
	return nil
}

func (r *RKE2) Verify() error {
	// verify if all images has been downloaded successfully
	for _, image := range listRKE2Images {
		filePath := filepath.Join(ensureTrailingSlash(r.OutputDirTarball), image)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			log.Printf("file does not exist: %s", filePath)
			return err
		}
		log.Printf("Image verified successfully: %s", filePath)
	}
	return nil
}

func (r *RKE2) Upload() error {
	// Upload the tarball files to the registry
	// TODO: implement me if needed (prepared if we change the airgap with rke2-capi-provider to use registry instead of artifacts)
	return nil
}

func replaceVersionLink(version string) string {
	return strings.ReplaceAll(version, "+", "%2B")
}

func ensureTrailingSlash(dir string) string {
	if !strings.HasSuffix(dir, "/") {
		dir += "/"
	}
	return dir
}
