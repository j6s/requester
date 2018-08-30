# Requester

Simple command line utility that crawls a website and calls every internal
link on the site. This is useful in order to warm caches of application and
for finding invalid links.

## Usage

```
Requester
===========
Simple website crawler that tries to follow every link that points
back to the same domain.

Usage:
$ ./requester {OPTIONAL ARGUMENTS} http://first-url.com http://second-url.com

Optional arguments:
 -list-success
       If specified successful (200) requests will also be printed in the summary
 -workers int
       The number of workers to use (default 8)
```

## Example

```
$ requester --list-success thej6s.com
2018/08/30 23:24:39 [200 OK] GET http://thej6s.com
2018/08/30 23:24:39 [200 OK] GET https://thej6s.com
2018/08/30 23:24:39 [200 OK] GET https://thej6s.com/projects
2018/08/30 23:24:39 [200 OK] GET https://thej6s.com/about
2018/08/30 23:24:39 [200 OK] GET http://thej6s.com/rss
2018/08/30 23:24:39 [200 OK] GET https://thej6s.com/articles/things-i-learned-2018-34
2018/08/30 23:24:39 [200 OK] GET https://thej6s.com/articles/things-i-learned-2018-33
2018/08/30 23:24:39 [200 OK] GET https://thej6s.com/articles/things-i-learned-2018-24
2018/08/30 23:24:39 [200 OK] GET https://thej6s.com/articles/things-i-learned-2018-23
2018/08/30 23:24:39 [200 OK] GET https://thej6s.com/articles/things-i-learned-2018-22
2018/08/30 23:24:39 [200 OK] GET https://thej6s.com/page:2
2018/08/30 23:24:39 [200 OK] GET http://thej6s.com/privacy
2018/08/30 23:24:39 [200 OK] GET https://thej6s.com/rss
2018/08/30 23:24:39 [200 OK] GET https://thej6s.com/privacy
2018/08/30 23:24:39 [200 OK] GET https://thej6s.com/articles/things-i-learned-2018-21
2018/08/30 23:24:39 [200 OK] GET https://thej6s.com/articles/things-i-learned-2018-20
2018/08/30 23:24:40 [200 OK] GET https://thej6s.com/articles/things-i-learned-2018-19
2018/08/30 23:24:40 [200 OK] GET https://thej6s.com/articles/things-i-learned-2018-18
2018/08/30 23:24:40 [200 OK] GET https://thej6s.com/articles/things-i-learned-2018-17
2018/08/30 23:24:40 [200 OK] GET https://thej6s.com/
2018/08/30 23:24:40 [200 OK] GET https://thej6s.com/page:3
2018/08/30 23:24:40 [200 OK] GET https://thej6s.com/articles/network-bridging-for-pi-zero
2018/08/30 23:24:40 [200 OK] GET https://thej6s.com/articles/thinkpad-t430-fan-control
2018/08/30 23:24:40 [200 OK] GET https://thej6s.com/articles/the-new-me
2018/08/30 23:24:40 ## 200 OK: 24 Requests
2018/08/30 23:24:40 	 -> http://thej6s.com
2018/08/30 23:24:40 	 -> https://thej6s.com
2018/08/30 23:24:40 	 -> https://thej6s.com/projects
2018/08/30 23:24:40 	 -> https://thej6s.com/about
2018/08/30 23:24:40 	 -> http://thej6s.com/rss
2018/08/30 23:24:40 	 -> https://thej6s.com/articles/things-i-learned-2018-34
2018/08/30 23:24:40 	 -> https://thej6s.com/articles/things-i-learned-2018-33
2018/08/30 23:24:40 	 -> https://thej6s.com/articles/things-i-learned-2018-24
2018/08/30 23:24:40 	 -> https://thej6s.com/articles/things-i-learned-2018-23
2018/08/30 23:24:40 	 -> https://thej6s.com/articles/things-i-learned-2018-22
2018/08/30 23:24:40 	 -> https://thej6s.com/page:2
2018/08/30 23:24:40 	 -> http://thej6s.com/privacy
2018/08/30 23:24:40 	 -> https://thej6s.com/rss
2018/08/30 23:24:40 	 -> https://thej6s.com/privacy
2018/08/30 23:24:40 	 -> https://thej6s.com/articles/things-i-learned-2018-21
2018/08/30 23:24:40 	 -> https://thej6s.com/articles/things-i-learned-2018-20
2018/08/30 23:24:40 	 -> https://thej6s.com/articles/things-i-learned-2018-19
2018/08/30 23:24:40 	 -> https://thej6s.com/articles/things-i-learned-2018-18
2018/08/30 23:24:40 	 -> https://thej6s.com/articles/things-i-learned-2018-17
2018/08/30 23:24:40 	 -> https://thej6s.com/
2018/08/30 23:24:40 	 -> https://thej6s.com/page:3
2018/08/30 23:24:40 	 -> https://thej6s.com/articles/network-bridging-for-pi-zero
2018/08/30 23:24:40 	 -> https://thej6s.com/articles/thinkpad-t430-fan-control
2018/08/30 23:24:40 	 -> https://thej6s.com/articles/the-new-me
```

## TODO

Write an actual README, this one sucks.
