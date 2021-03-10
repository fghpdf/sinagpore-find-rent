package chrome

import (
	"context"
	"sync"
	"time"

	"github.com/chromedp/chromedp"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	log "github.com/sirupsen/logrus"
)

type Condo struct {
	name     string
	address  string
	district string
	// 资产期限
	tenure string
	// 开发商
	developer string
	url       string
	facility  *Facility
	facString string
}

type Facility struct {
	// 泳池
	pool bool
	// 网球场
	tennisCourt bool
	// 读书角
	readingCorner bool
	// 屋顶花园
	rooftopGarden bool
	// 健身区
	fitnessArea bool
	// 俱乐部
	clubHouse bool
	// 健身房
	gymnasium bool
	// 烧烤设备
	bbqPit bool
	// 秘密花园
	secretGarden bool
	// 慢跑道
	joggingTrack bool
	// 蒸汽室
	steamRoom bool
	// 停车场
	carPark bool
	// 安保
	security bool
}

func Run() {
	// create chrome instance
	ctx, cancel := chromedp.NewContext(
		context.Background(),
		chromedp.WithLogf(log.Infof),
	)
	defer cancel()

	run(ctx)
}

func RunWithRemote() {
	debugUrl := getDebugURL()

	ctx, cancel := chromedp.NewRemoteAllocator(context.Background(), debugUrl)
	defer cancel()

	run(ctx)
}

func run(ctx context.Context) {
	// create a timeout
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	condoUrls := analyzeCondoList(ctx, "https://condo.singaporeexpats.com/name/0-9")

	log.WithField("condo number", len(condoUrls)).Info("need to analyze")

	var wg sync.WaitGroup

	condos := make([]*Condo, 0)
	for _, condoUrl := range condoUrls {
		wg.Add(1)
		go analyzeCondo(ctx, condoUrl, &wg, &condos)
	}

	wg.Wait()

	mongoCtx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	client, err := mongo.Connect(mongoCtx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Panic(err)
	}
	collection := client.Database("condo").Collection("condos")

	for _, condo := range condos {
		log.Info(condo)
		err := condo.insert(mongoCtx, collection)
		if err != nil {
			log.Errorf("mongo %v", err)
			continue
		}
	}

}
