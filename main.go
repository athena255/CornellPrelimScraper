package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

var (
	// Example: http://www.cs.cornell.edu/courses/cs4410/2016fa
	searchString    = "http://www.cs.cornell.edu/courses/cs%s/%d%s"
	urlRegexStr     = `(?i)/courses/cs\d{3,4}/(\d{4})(fa|sp)/(.*(%s).*\.(.*))`
	fileNameRegex   = regexp.MustCompile(`.*/(.*)`)
	documentRegex   = regexp.MustCompile(`.*(pdf|txt|doc|ppt)`)
	defaultKeywords = `prelim\d*|final|sol|midterm|sample|review|exam|answer`
)

type Param struct {
	CourseID string
	InitYear string
	EndYear  string
	Verbose  bool
}

type Exam struct {
	FullURL    string
	Year       string
	Semester   string
	UniqueName string
	Category   string // prelim, final, homework solutions
	Extension  string // .pdf, .doc
}

func usage() {
	fmt.Println("Scraper scrapes for past exams, solutions, and homeworks.")
	fmt.Println("Usage: ./scraper <class number> <startYear> <endYear> <keywords> <options> ")
	fmt.Println(`Example: ./scraper 312 1999 2002 "prelim\d*|final|sol|midterm|sample|review|exam|answer"`)
	fmt.Println("Options: ")
	fmt.Println("v verbose mode prints all discovered links")
	fmt.Println("Keywords: Use regex to match with file names default is:" + `"prelim\d*|final|sol|midterm|sample|review|exam|answer"`)
	os.Exit(1)
}

func main() {
	nLen := len(os.Args)
	if nLen != 4 && nLen != 5 && nLen != 6 {
		usage()
	}

	param := Param{
		os.Args[1],
		os.Args[2],
		os.Args[3],
		false,
	}

	// See if user set verbose mode
	if nLen == 6 {
		param.Verbose = true
	}

	var urlRegex *regexp.Regexp
	if nLen == 5 {
		urlRegex = regexp.MustCompile(fmt.Sprintf(urlRegexStr, os.Args[4]))
	} else {
		urlRegex = regexp.MustCompile(fmt.Sprintf(urlRegexStr, defaultKeywords))
	}

	c := colly.NewCollector(
		colly.AllowedDomains("www.cs.cornell.edu"),
		// colly.CacheDir("./cache"), // cache makes us miss files
	)

	// Keep track of visited links
	setVisited := map[string]bool{}

	// Map unique documents to links
	setDocuments := map[string][]Exam{}

	// On every <a> with href, call this function
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Request.AbsoluteURL(e.Attr("href"))
		_, ok := setVisited[link]
		if ok {
			return
		}
		setVisited[link] = true // Mark this link as visited
		isExam := urlRegex.FindStringSubmatch(link)
		res := fileNameRegex.FindStringSubmatch(link)
		if isExam != nil { // If this links to a potential exam/homework/prelim
			fileName := res[1]
			// Create a new Exam record
			fmt.Println("======================")
			fmt.Println(isExam[3])
			newexam := Exam{
				link,
				isExam[1], // year
				isExam[2], // semester
				isExam[3], // documentName
				isExam[4], // documentCat
				isExam[5], // documentType
			}
			// Add Exam record to map
			setDocuments[fileName] = append(setDocuments[fileName], newexam)

		} else { // Else this links to some path that we have not visited
			if param.Verbose {
				fmt.Printf("\t[Discovered] %s\n", link)
			}
			isFile := documentRegex.FindStringSubmatch(link)
			// Only look at paths that are not files
			if !ok && isFile == nil && strings.Contains(link, "courses") {
				// Visit the unvisited path
				e.Request.Visit(link)
			}
		}
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("[Visiting]", r.URL)
		setVisited[r.URL.String()] = true
	})

	fmt.Println("***************************************************************")
	fmt.Printf("\nSearching for documents for cs%s from %s to %s...\n\n", param.CourseID, param.InitYear, param.EndYear)
	fmt.Println("***************************************************************")

	initYear, _ := strconv.Atoi(param.InitYear)
	endYear, _ := strconv.Atoi(param.EndYear)

	// From user specified InitYear to EndYear, search both fall and spring semesters
	for i := initYear; i <= endYear; i++ {
		fmt.Printf("\nNow searching for cs%s documents from fall of %d\n", param.CourseID, i)
		c.Visit(fmt.Sprintf(searchString, param.CourseID, i, "fa"))
		fmt.Println("===============================================================")
		fmt.Printf("\nNow searching for cs%s documents from spring of %d\n", param.CourseID, i)
		c.Visit(fmt.Sprintf(searchString, param.CourseID, i, "sp"))
		fmt.Println("===============================================================")
	}

	fmt.Printf("\nFound %d (unique) documents at these locations: \n\n", len(setDocuments))

	for k, v := range setDocuments {
		fmt.Printf("[%s]", k)
		for _, exam := range v {
			fmt.Printf("\n\t%s\n", exam.FullURL)
		}
		fmt.Printf("\n")
	}

}
