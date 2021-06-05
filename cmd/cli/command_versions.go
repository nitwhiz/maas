package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/nitwhiz/maas/pkg/minecraft"
	"io"
	"net/http"
	"os"
	"path"
	"sort"
	"strings"
	"time"
)

const manifestCacheFile = "maas_cached_manifest.json"
const manifestCacheDuration = time.Hour * 4

type VersionsCmd struct {
	Latest         bool   `kong:"short='l',default='false',help='Just show the latest release and snapshot.'"`
	Type           string `kong:"short='t',default='all',help='Filter versions by type. Can be release, snapshot or all.'"`
	Search         string `kong:"short='s',help='Search by version id'"`
	Count          int    `kong:"short='c',default='10',help='Limit the results.'"`
	OrderBy        string `kong:"short='o',default='time',help='Order by specific field from the manifest JSON. Can be time or releaseTime.'"`
	OrderDirection string `kong:"short='d',default='desc',help='Change order direction. Can be asc or desc.'"`
	ForceDownload  bool   `kong:"short='f',default='false',help='Force re-download of manifest file. Otherwise cached manifest is used.'"`
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

func unmarshalManifest(bs []byte) (*minecraft.Manifest, error) {
	var manifest minecraft.Manifest

	err := json.Unmarshal(bs, &manifest)

	if err != nil {
		return nil, err
	}

	return &manifest, nil
}

func downloadManifest() (*minecraft.Manifest, error) {
	r, err := http.Get(minecraft.ManifestUrl)

	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(r.Body)

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(r.Body)

	if err != nil {
		return nil, err
	}

	bs := buf.Bytes()
	p := path.Join(os.TempDir(), manifestCacheFile)

	err = os.WriteFile(p, bs, 0666)

	if err != nil {
		return nil, err
	}

	return unmarshalManifest(bs)
}

func readManifestFromFile(path string) (*minecraft.Manifest, error) {
	bs, err := os.ReadFile(path)

	if err != nil {
		return nil, err
	}

	return unmarshalManifest(bs)
}

func readManifestFromCache() (*minecraft.Manifest, error) {
	p := path.Join(os.TempDir(), manifestCacheFile)

	cs, err := os.Stat(p)

	if os.IsNotExist(err) {
		return nil, err
	}

	if cs.ModTime().Before(time.Now().Add(-manifestCacheDuration)) {
		return nil, &CacheExpiredError{}
	}

	return readManifestFromFile(p)
}

func printLatest(m *minecraft.Manifest, t string) {
	if t == "release" {
		fmt.Println(m.Latest.Release)
	} else if t == "snapshot" {
		fmt.Println(m.Latest.Snapshot)
	} else {
		bs, err := json.MarshalIndent(m.Latest, "", "  ")

		if err != nil {
			fmt.Println("[]")
			return
		}

		fmt.Println(string(bs))
	}
}

func filterByType(versions []minecraft.Version, t string) []minecraft.Version {
	if t == "all" {
		return append([]minecraft.Version{}, versions...)
	}

	var result []minecraft.Version

	for _, v := range versions {
		if v.Type == t {
			result = append(result, v)
		}
	}

	return result
}

func filterByIdSearch(versions []minecraft.Version, sid string) []minecraft.Version {
	if sid == "" {
		return append([]minecraft.Version{}, versions...)
	}

	var result []minecraft.Version

	for _, v := range versions {
		if strings.Contains(v.Id, sid) {
			result = append(result, v)
		}
	}

	return result
}

func orderByFieldAndDirection(versions []minecraft.Version, field string, direction string) []minecraft.Version {
	v := append([]minecraft.Version{}, versions...)

	asc := direction == "asc"

	if field == "releaseTime" {
		sort.SliceStable(v, func(i int, j int) bool {
			if asc {
				return v[i].ReleaseTime < v[j].ReleaseTime
			} else {
				return v[i].ReleaseTime > v[j].ReleaseTime
			}
		})
	} else {
		sort.SliceStable(v, func(i int, j int) bool {
			if asc {
				return v[i].Time < v[j].Time
			} else {
				return v[i].Time > v[j].Time
			}
		})
	}

	return v
}

func getSlice(versions []minecraft.Version, size int) []minecraft.Version {
	if len(versions) >= size {
		return append([]minecraft.Version{}, versions[:size]...)
	}

	return append([]minecraft.Version{}, versions...)
}

func printFiltered(m *minecraft.Manifest, c *VersionsCmd) {
	versions := getSlice(
		orderByFieldAndDirection(
			filterByIdSearch(
				filterByType(m.Versions, c.Type),
				c.Search,
			),
			c.OrderBy,
			c.OrderDirection,
		),
		c.Count,
	)

	bs, err := json.MarshalIndent(versions, "", "  ")

	if err != nil {
		fmt.Println("[]")
		return
	}

	fmt.Println(string(bs))
}

func (c *VersionsCmd) Run() error {
	var manifest *minecraft.Manifest
	var err error

	if c.ForceDownload {
		manifest, err = downloadManifest()

		if err != nil {
			return err
		}
	} else {
		manifest, err = readManifestFromCache()

		if os.IsNotExist(err) || IsCacheExpired(err) {
			manifest, err = downloadManifest()

			if err != nil {
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
