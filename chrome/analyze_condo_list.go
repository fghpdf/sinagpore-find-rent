package chrome

import (
	"context"
	"sync"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	log "github.com/sirupsen/logrus"
)

func analyzeCondoList(ctx context.Context, directoryUrl string, wg *sync.WaitGroup, condoUrls *[]string) {
	defer wg.Done()

	log.Infof("Now is %s\n", directoryUrl)
	ctx2, _ := chromedp.NewContext(ctx)

	var nodes []*cdp.Node
	err := chromedp.Run(ctx2,
		chromedp.Navigate(directoryUrl),
		chromedp.Nodes(`a.title_link`, &nodes, chromedp.ByQueryAll),
	)
	if err != nil {
		log.Error(err)
	}

	for _, node := range nodes {
		if url, ok := node.Attribute("href"); ok {
			*condoUrls = append(*condoUrls, url)
		}
	}
}
