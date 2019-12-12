package controller

import (
	"daemonw/conf"
	"daemonw/dao"
	"daemonw/entity"
	"daemonw/util"
	"daemonw/xlog"
	"github.com/PuerkitoBio/goquery"
	"io"
	"net/http"
	"os"
	"path/filepath"
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
				info.Version = app.Version
				db := dao.NewAppDao()
				err := db.CreateAppInfo(info)
				if err != nil {
					xlog.Error().Msgf("err: %s", err.Error())
				}
				saveIcon(app, info.Icon)
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

	doc.Find("div.app-intro.cf").Each(func(i int, selection *goquery.Selection) {
		selection.Find("img.yellow-flower").Each(func(i int, selection *goquery.Selection) {
			src, exist := selection.Attr("src")
			if exist {
				appInfo.Icon = src
			}
		})
	})

	doc.Find("div.app-intro.cf").Each(func(i int, selection *goquery.Selection) {
		selection.Find("div.intro-titles").Each(func(i int, selection *goquery.Selection) {
			children := selection.Children()
			children.Each(func(i int, selection *goquery.Selection) {
				if (i == 0) {
					appInfo.Vendor = selection.Text()
				} else if (i == 1) {
					appInfo.Name = selection.Text()
				} else if (i == 2) {
					text := selection.Text()
					s := strings.Split(text, "|")
					s = strings.Split(s[0], "：")
					if len(s) >= 2 {
						appInfo.Category = s[1]
					}
				}
			})
		})
	})

	urls := make([]string, 0)
	doc.Find("div#J_thumbnail_wrap").Each(func(i int, selection *goquery.Selection) {
		selection.Find("img").Each(func(i int, selection *goquery.Selection) {
			src, exist := selection.Attr("src")
			if exist {
				urls = append(urls, src)
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
	if len(urls) != 0 {
		appInfo.ImageDetail = strings.Join(urls, `,`)
	}
	return appInfo, nil
}

func saveIcon(app entity.App, iconUrl string) {
	iconFile := filepath.Join(conf.Config.Data, app.AppId, app.Version, "icon.png")
	if util.ExistFile(iconFile) {
		return
	}
	if iconUrl == "" {
		return
	}
	resp, err := http.Get(iconUrl)
	if err != nil {
		xlog.Error().Msgf(`err: save icon failed for "%s“`, app.Name)
		return
	}
	if resp.StatusCode == http.StatusOK {
		defer resp.Body.Close()
		f, err := os.OpenFile(iconFile, os.O_CREATE|os.O_RDWR, os.ModePerm)
		if err == nil {
			defer f.Close()
			io.Copy(f, resp.Body)
		}
	}
}

type ApkPureSpider struct {
}

func (spider *ApkPureSpider) FetchApkInfo(pkg string) (info *entity.AppInfo, err error) {
	return nil, nil
}
