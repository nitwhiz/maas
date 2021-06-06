package main

import (
	"encoding/json"
	"fmt"
	"github.com/nitwhiz/maas/pkg/minecraft"
	"os"
	"text/tabwriter"
	"time"
)

const manifestCacheFile = "maas_cached_manifest.json"
const manifestCacheDuration = time.Hour * 4

type VersionsCmd struct {
	Latest         bool   `kong:"short='l',default='false',help='Just show the latest release and snapshot'"`
	Type           string `kong:"short='t',default='all',help='Filter versions by type. Can be release, snapshot or all'"`
	Search         string `kong:"short='s',help='Search by version id'"`
	Count          int    `kong:"short='c',default='10',help='Limit the results'"`
	OrderBy        string `kong:"short='o',default='time',help='Order by specific field from the manifest JSON. Can be time or releaseTime'"`
	OrderDirection string `kong:"short='d',default='desc',help='Change order direction. Can be asc or desc'"`
	ForceDownload  bool   `kong:"short='f',default='false',help='Force re-download of manifest file. Otherwise cached manifest is used'"`
}

type CacheExpiredError struct {
}

func (e *CacheExpiredError) Error() string {
	return "cache is expired"
}

func IsCacheExpired(e error) bool {
	if _, ok := e.(*CacheExpiredError); ok {
		return true
	}

	return false
}

func readManifestFromCache(cachePath string, cacheExpiry time.Duration) (*minecraft.Manifest, error) {
	cs, err := os.Stat(cachePath)

	if os.IsNotExist(err) {
		return nil, err
	}

	if cs.ModTime().Before(time.Now().Add(-cacheExpiry)) {
		return nil, &CacheExpiredError{}
	}

	return minecraft.ReadManifestFromFile(cachePath)
}

func cacheManifest(m *minecraft.Manifest, file string) error {
	bs, err := json.Marshal(m)

	if err != nil {
		return err
	}

	err = os.WriteFile(file, bs, 0666)

	return err
}

func printLatest(m *minecraft.Manifest, t string) {
	if t == "all" || t == "release" {
		fmt.Printf("release: %s\n", m.Latest.Release)
	}

	if t == "all" || t == "snapshot" {
		fmt.Printf("snapshot: %s\n", m.Latest.Snapshot)
	}
}

func printFiltered(m *minecraft.Manifest, c *VersionsCmd) {
	m.FilterVersionsByType(c.Type).
		FilterVersionsByIdSubstring(c.Search).
		OrderVersionsByFieldAndDirection(c.OrderBy, c.OrderDirection)

	sum := len(m.Versions)

	m.SliceVersions(c.Count)

	writer := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	_, _ = fmt.Fprintln(writer, "ID\tTYPE\tRELEASE TIME\tSHA1")

	for _, v := range m.Versions {
		_, _ = fmt.Fprintf(writer, "%s\t%s\t%s\t%s\n", v.Id, v.Type, v.ReleaseTime, v.SHA1)
	}

	_ = writer.Flush()

	if sum > c.Count {
		fmt.Println("(more)")
	}
}

func (c *VersionsCmd) Run() error {
	var manifest *minecraft.Manifest
	var err error

	if c.ForceDownload {
		manifest, err = minecraft.DownloadManifest()

		if err != nil {
			return err
		}

		if err := cacheManifest(manifest, manifestCacheFile); err != nil {
			return err
		}
	} else {
		manifest, err = readManifestFromCache(manifestCacheFile, manifestCacheDuration)

		if os.IsNotExist(err) || IsCacheExpired(err) {
			manifest, err = minecraft.DownloadManifest()

			if err != nil {
				return err
			}

			if err := cacheManifest(manifest, manifestCacheFile); err != nil {
				return err
			}
		}
	}

	if c.Latest {
		printLatest(manifest, c.Type)
	} else {
		printFiltered(manifest, c)
	}

	return nil
}
