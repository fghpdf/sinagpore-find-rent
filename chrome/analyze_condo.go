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

	c := &Condo{Url: condoUrl}
	facilities := make([]string, 0)
	var leftNodes []*cdp.Node
	var rightNodes []*cdp.Node

	ctxNew, _ := chromedp.NewContext(ctx)

	err := chromedp.Run(ctxNew,
		chromedp.Navigate(condoUrl),
		// Name
		chromedp.Text(`div.propertytitlecontainer > div > h1`, &c.Name),
		// address
		chromedp.Text(`div.propertytitlecontainer > div.propcol1 > :nth-child(2)`, &c.Address),
		// district
		chromedp.Text(`div.propertytitlecontainer > div.propcol1 > :nth-child(3)`, &c.Address),
		// tenure
		chromedp.Text(`div.propertytitlecontainer > div.propcol1 > :nth-child(4)`, &c.Tenure),
		// developer
		chromedp.Text(`div.propertytitlecontainer > div.propcol1 > :nth-child(5)`, &c.Developer),
		chromedp.Nodes(`div#tabs-property-facilities > :nth-child(2) > ul > li`, &leftNodes, chromedp.ByQueryAll, chromedp.AtLeast(0)),
		chromedp.Nodes(`div#tabs-property-facilities > :nth-child(3) > ul > li`, &rightNodes, chromedp.ByQueryAll, chromedp.AtLeast(0)),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var fac string
			var facStr string
			for i := 1; i <= len(leftNodes); i++ {
				sel := fmt.Sprintf("div#tabs-property-facilities > :nth-child(2) > ul > :nth-child(%d)", i)
				err := chromedp.Text(sel, &fac).Do(ctx)
				if err != nil {
					return err
				}
				facilities = append(facilities, fac)
				facStr += fmt.Sprintf("| %s", fac)
			}

			for i := 1; i < len(rightNodes); i++ {
				sel := fmt.Sprintf("div#tabs-property-facilities > :nth-child(3) > ul > :nth-child(%d)", i)
				err := chromedp.Text(sel, &fac).Do(ctx)
				if err != nil {
					return err
				}
				facilities = append(facilities, fac)
				facStr += fmt.Sprintf("| %s", fac)
			}

			c.FacString = facStr

			return nil
		}),
	)
	if err != nil {
		log.Errorf("analyze condo %v", err)
		return
	}

	c.Facility = analyzeFac(facilities)

	*condos = append(*condos, c)

	return
}
