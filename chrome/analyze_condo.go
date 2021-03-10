package chrome

import (
	"context"
	"fmt"
	"sync"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	log "github.com/sirupsen/logrus"
)

func analyzeCondo(ctx context.Context, condoUrl string, wg *sync.WaitGroup, condos *[]*Condo) {

	defer wg.Done()

	if !isCondoUrl(condoUrl) {
		return
	}

	log.Infof("Now is %s\n", condoUrl)

	c := &Condo{url: condoUrl}
	facilities := make([]string, 0)
	var nodes []*cdp.Node

	ctxNew, _ := chromedp.NewContext(ctx)

	err := chromedp.Run(ctxNew,
		chromedp.Navigate(condoUrl),
		// name
		chromedp.Text(`div.propertytitlecontainer > div > h1`, &c.name),
		// address
		chromedp.Text(`div.propertytitlecontainer > div.propcol1 > :nth-child(2)`, &c.address),
		// district
		chromedp.Text(`div.propertytitlecontainer > div.propcol1 > :nth-child(3)`, &c.address),
		// tenure
		chromedp.Text(`div.propertytitlecontainer > div.propcol1 > :nth-child(4)`, &c.tenure),
		// developer
		chromedp.Text(`div.propertytitlecontainer > div.propcol1 > :nth-child(5)`, &c.developer),
		chromedp.Nodes(`div#tabs-property-facilities > :nth-child(2) > ul > li`, &nodes, chromedp.ByQueryAll),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var fac string
			var facStr string
			for i := 1; i <= len(nodes); i++ {
				sel := fmt.Sprintf("div#tabs-property-facilities > :nth-child(2) > ul > :nth-child(%d)", i)
				err := chromedp.Text(sel, &fac).Do(ctx)
				if err != nil {
					return err
				}
				facilities = append(facilities, fac)
				facStr += fac
			}

			c.facString = facStr

			return nil
		}),
	)
	if err != nil {
		log.Error(err)
		return
	}

	c.facility = analyzeFac(facilities)

	*condos = append(*condos, c)

	return
}
