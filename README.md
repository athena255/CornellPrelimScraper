# CornellPrelimScraper

> Scrapes cs.cornell.edu for past prelims and generates a list of interesting files :)

The script searches for documents in course webpages from previous years. It uses regex strings and the fact that Cornell's CS course webpages are of the format: `cs.cornell.edu/cs[course number]/[year][semester]` to find a range of documents.

Given a course number, a starting year, an ending year, and an optional regex, it will a compile a list of documents reachable from the specified range of course webpages. 

The default regex is `prelim\d*|final|sol|midterm|sample|review|exam|answer`

Example outputs can be found in the `examples` folder. 

---

### Setup

> install go (brew)
```shell
$ brew install go
```
> install go (apt)
```shell
$ apt install golang-go
```
> clone repo and build
```shell
$ go get github.com/gocolly/colly
$ git clone https://github.com/athena255/CornellPrelimScraper.git
$ go build
```
---

## Usage
Note that in 2008 the course naming system went from using three-digit course numbers to four-digit course numbers. So use the three-digit version of the course number if searching for prelims in years prior to 2008. 
> Default search
```shell
$ ./scraper [course number] [start year] [end year]
$ ./scaper 4410 2010 2019
$ ./scraper 312 1999 2008
```
> Regex search
```shell
$ ./scraper [course number] [start year] [end year] [search terms]
$ ./scraper 3110 2012 2020 "prelim|exam|solution"
```
