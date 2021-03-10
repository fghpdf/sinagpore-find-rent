package chrome

import (
	"context"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	log "github.com/sirupsen/logrus"
)

func analyzeCondoList(ctx context.Context, directoryUrl string) []string {
	log.Infof("Now is %s\n", directoryUrl)

	condoUrls := make([]string, 0)

	var nodes []*cdp.Node
	err := chromedp.Run(ctx,
		chromedp.Navigate(directoryUrl),
		chromedp.Nodes(`a.title_link`, &nodes, chromedp.ByQueryAll),
	)
	if err != nil {
		log.Error(err)
	}

	for _, node := range nodes {
		if url, ok := node.Attribute("href"); ok {
			condoUrls = append(condoUrls, url)
		}
	}

	return condoUrls
}
