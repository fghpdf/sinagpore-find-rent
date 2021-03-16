package chrome

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"sync"

	"github.com/chromedp/chromedp"
	log "github.com/sirupsen/logrus"
)

func loopAnalyzeDirectory(ctx context.Context, tpls []string, maxConcurrent int) []string {
	concurrentGoroutines := make(chan struct{}, maxConcurrent)

	var wg sync.WaitGroup
	res := make([]string, 0)

	for _, tpl := range tpls {
		wg.Add(1)
		go func(tpl string) {
			defer wg.Done()
			concurrentGoroutines <- struct{}{}
			urls := analyzeDirectory(ctx, tpl)
			if len(urls) != 0 {
				res = append(res, urls...)
			}
			<-concurrentGoroutines
		}(tpl)
	}

	wg.Wait()

	return res
}

// tpl should like this https://condo.singaporeexpats.com/%sname/0-9
func analyzeDirectory(ctx context.Context, tpl string) []string {
	url := fmt.Sprintf(tpl, "")
	log.Infof("Now is %s\n", url)

	directoryUrls := []string{url}
	var totalStr string
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.Text("div.propertyfound > span.pageno", &totalStr),
	)
	if err != nil {
		log.Fatal(err)
	}

	total, err := strconv.ParseFloat(totalStr, 64)
	if err != nil {
		log.Fatal(err)
	}

	pages := int(math.Ceil(total / 50))

	log.Infof("pages is %d", pages)

	for i := 2; i <= pages; i++ {
		directoryUrls = append(directoryUrls, fmt.Sprintf(tpl, fmt.Sprintf("%d/", i)))
	}

	return directoryUrls
}
