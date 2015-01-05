package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/pmylund/go-cache"
)

var cacheProxies = cache.New(cache.NoExpiration, cache.NoExpiration)

func getProxies(domain string) []string {
	var proxies []string
	if x, found := cacheProxies.Get(domain); found {
		proxies = x.([]string)
	}
	return proxies
}

func setProxies(domain string, urls []string) {
	log.Info("Successfully cached: " + domain)
	cacheProxies.Set(domain, urls, cache.NoExpiration)
}

func addProxy(domain string, url string) {
	if domain == "" {
		return
	}
	if url == "" {
		return
	}

	urls := getProxies(domain)
	urls = append(urls, url)
	cacheProxies.Set(domain, urls, cache.NoExpiration)
}

func removeProxy(domain string, url string) {
	if domain == "" {
		return
	}
	if url == "" {
		return
	}

	urls := getProxies(domain)
	var newUrls []string
	for _, u := range urls {
		if u != url {
			newUrls = append(newUrls, u)
		}
	}
	cacheProxies.Set(domain, newUrls, cache.NoExpiration)
}
