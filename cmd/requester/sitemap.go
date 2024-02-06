package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type SitemapUrl struct {
	Loc        string     `xml:"loc"`
	Lastmod    *time.Time `xml:"lastmod"`
	Changefreq string     `xml:"changefreq"`
	Priority   int        `xml:"priority"`
}

type SitemapReference struct {
	Loc     string     `xml:"loc"`
	Lastmod *time.Time `xml:"lastmod"`
}

type Sitemap struct {
	Urls       []SitemapUrl       `xml:"url"`
	References []SitemapReference `xml:"sitemap"`
}

func GetUrlsInSitemap(sitemapUrl string) ([]url.URL, error) {

	response, err := http.Get(sitemapUrl)
	if err != nil {
		return []url.URL{}, fmt.Errorf("cannot request sitemap %s: %s", sitemapUrl, err)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return []url.URL{}, fmt.Errorf("cannot read response of %s: %s", sitemapUrl, err)
	}

	var sitemap Sitemap
	err = xml.Unmarshal(body, &sitemap)
	if err != nil {
		return []url.URL{}, fmt.Errorf("cannot parse XML of sitemap %s: %s", sitemapUrl, err)
	}

	urls := make([]url.URL, len(sitemap.Urls))
	for i, urlInSitemap := range sitemap.Urls {
		parsed, err := url.Parse(urlInSitemap.Loc)
		if err != nil {
			return urls, fmt.Errorf("cannot parse URL %s in sitemap %s: %s", urlInSitemap, sitemapUrl, err)
		}
		urls[i] = *parsed
	}

	for _, ref := range sitemap.References {
		refUrls, err := GetUrlsInSitemap(ref.Loc)
		if err != nil {
			return urls, fmt.Errorf("error requesting sitemap %s referenced by %s: %s", ref.Loc, sitemapUrl, err)
		}

		urls = append(urls, refUrls...)
	}

	return urls, nil
}
