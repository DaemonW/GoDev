package controller

import (
	"daemonw/dao"
	"daemonw/entity"
	"daemonw/xlog"
	"github.com/PuerkitoBio/goquery"
	"strings"
)

var (
	AppInfoSpiderChan = make(chan entity.App, 10)
)

func init() {
	go func() {
		spider := &MiStoreSpider{}
		for {
			app := <-AppInfoSpiderChan
			info, err := spider.FetchApkInfo(app.AppId)
			if err == nil && info != nil {
				info.Id = app.Id
				db := dao.NewAppDao()
				err := db.CreateAppInfo(info)
				if err != nil {
					xlog.Error().Msgf("err: %s", err.Error())
				}
			}
		}
	}()
}

type PkgSpider interface {
	FetchApkInfo(pkg string) (string, string, error)
}

const (
	MiStoreUrl = "http://app.mi.com/details?id="
)

type MiStoreSpider struct {
}

func (spider *MiStoreSpider) FetchApkInfo(pkg string) (info *entity.AppInfo, err error) {
	appUrl := MiStoreUrl + pkg
	doc, err := goquery.NewDocument(appUrl)
	if err != nil {
		return nil, err
	}
	appInfo := &entity.AppInfo{}
	doc.Find("div.app-text").Each(func(i int, selection *goquery.Selection) {
		selection.Find("p.pslide").Each(func(i int, selection *goquery.Selection) {
			switch i {
			case 0:
				appInfo.Description = selection.Text()
			case 1:
				appInfo.ChangeLog = selection.Text()
			}
		})
	})
	urls := strings.Builder{}
	doc.Find("div#J_thumbnail_wrap").Each(func(i int, selection *goquery.Selection) {
		selection.Find("img").Each(func(i int, selection *goquery.Selection) {
			src, exist := selection.Attr("src")
			if exist {
				urls.WriteString(src + ",")
			}
		})
	})

	doc.Find("div.details.preventDefault").Each(func(i int, selection *goquery.Selection) {
		selection.Find("ul.cf").Each(func(i int, selection *goquery.Selection) {
			selection.Find("li.weight-font").Each(func(i int, selection *goquery.Selection) {
				if i == 1 {
					appInfo.Version = selection.Next().Text()
				}
				if i == 3 {
					appInfo.Package = selection.Next().Text()
				}
			})
		})
	})
	appInfo.ImageDetail = urls.String()
	return appInfo, nil
}

type ApkPureSpider struct {
}

func (spider *ApkPureSpider) FetchApkInfo(pkg string) (info *entity.AppInfo, err error) {
	return nil, nil
}
