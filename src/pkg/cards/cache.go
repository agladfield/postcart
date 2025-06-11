package cards

import (
	"embed"
	"path"

	"github.com/agladfield/postcart/pkg/shared/tools/img"
	"github.com/davidbyttow/govips/v2/vips"
	lru "github.com/hashicorp/golang-lru/v2"
)

type imageCache struct {
	*lru.Cache[string, *vips.ImageRef]
	fs       *embed.FS
	basePath string
}

func (cache *imageCache) Obtain(pathStr string) (*vips.ImageRef, error) {
	if cache.basePath != "" {
		pathStr = path.Join(cache.basePath, pathStr)
	}
	cachedImage, exists := cache.Get(pathStr)
	if exists {
		return cachedImage.Copy()
	} else {
		loadedImage, loadErr := img.LoadFromEmbed(cache.fs, pathStr)
		if loadErr != nil {
			return nil, loadErr
		}
		cache.Add(pathStr, loadedImage)

		return loadedImage.Copy()
	}
}

func evictImage(key string, image *vips.ImageRef) {
	image.Close()
}

func newImageCache(fs *embed.FS, size int, basePath ...string) (*imageCache, error) {
	lruCache, err := lru.NewWithEvict(size, evictImage)
	if err != nil {
		return nil, err
	}

	basePathString := ""
	if len(basePath) > 0 {
		basePathString = basePath[0]
	}

	return &imageCache{lruCache, fs, basePathString}, nil
}

func closeCache() {
	if sCache != nil {
		sCache.Purge()
	}
	if pcCache != nil {
		pcCache.Purge()
	}
}

// stamp cache
// postcad cache

var (
	sCache  *imageCache
	pcCache *imageCache
)

func createCaches() error {
	var cacheErr error
	sCache, cacheErr = newImageCache(&stampResources, 15)
	if cacheErr != nil {
		return cacheErr
	}
	pcCache, cacheErr = newImageCache(&postcardResources, 20)
	if cacheErr != nil {
		return cacheErr
	}

	return nil
}
