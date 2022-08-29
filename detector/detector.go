package detector

import (
	"istio-botdetector/config"
	"regexp"

	lru "github.com/hashicorp/golang-lru"
)

type Detector struct {
	Allow []*regexp.Regexp
	Deny  []*regexp.Regexp
}

type CacheBucket struct {
	cache *lru.Cache
}

func InitializeDetector(config config.PluginConfiguration) (Detector, error) {
	var detector Detector

	for _, ua := range config.Allow {
		reg, err := regexp.Compile(ua)
		if err != nil {
			return Detector{}, err
		}
		detector.Allow = append(detector.Allow, reg)
	}
	for _, ua := range config.Deny {
		reg, err := regexp.Compile(ua)
		if err != nil {
			return Detector{}, err
		}
		detector.Deny = append(detector.Deny, reg)
	}

	return detector, nil
}

func IsBot(userAgent string, d Detector, config config.PluginConfiguration, cacheBucket CacheBucket) bool {
	cached, ok := cacheBucket.cache.Get(userAgent)
	if ok {
		return cached.(bool)
	}

	for _, pattern := range d.Allow {
		if pattern.MatchString(userAgent) {
			cacheBucket.cache.Add(userAgent, false)
			return false
		}
	}

	if config.DenyAll {
		return true
	}
	for _, pattern := range d.Deny {
		if pattern.MatchString(userAgent) {
			cacheBucket.cache.Add(userAgent, true)
			return true
		}
	}

	cacheBucket.cache.Add(userAgent, false)
	return false
}

func InitializeCacheBucket(config config.PluginConfiguration) (CacheBucket, error) {
	cache, err := lru.New(config.CacheSize)
	if err != nil {
		return CacheBucket{}, err
	}
	return CacheBucket{
		cache: cache,
	}, nil
}
