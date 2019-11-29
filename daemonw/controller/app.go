package controller

import (
	"daemonw/conf"
	"daemonw/crypto"
	"daemonw/dao"
	"daemonw/entity"
	"daemonw/util"
	"daemonw/xerr"
	"daemonw/xlog"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/shogo82148/androidbinary/apk"
	"image"
	"image/png"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	counter = time.Now().Unix()
	locker  = sync.RWMutex{}
)

func CreateApp(c *gin.Context) {
	if t := c.Request.Header.Get("Content-Type"); !strings.HasPrefix(t, gin.MIMEMultipartPOSTForm) {
		fmt.Println(t)
		c.JSON(http.StatusBadRequest, entity.NewRespErr(xerr.CodeCrateApp, "unacceptable content type"))
		return
	}
	err := c.Request.ParseMultipartForm(32 << 20)
	if err != nil {
		c.JSON(http.StatusInternalServerError, entity.NewRespErr(xerr.CodeCrateApp, "parse multipart body failed"))
		return
	}
	encrypted := c.PostForm("encrypted")
	enc, _ := strconv.ParseBool(encrypted)
	r, _, err := c.Request.FormFile("data")
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.NewRespErr(xerr.CodeCrateApp, "parse file data failed"))
		return
	}
	defer r.Close()
	stamp := time.Now().UnixNano()
	tempFile := filepath.Join(conf.Config.Data, strconv.FormatInt(stamp, 10))
	f, err := os.OpenFile(tempFile, os.O_CREATE|os.O_RDWR, os.ModePerm)
	defer f.Close()
	_, err = io.Copy(f, r)
	util.PanicIfErr(err)

	app, icon, err := ParseApkFromFile(tempFile)
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.NewRespErr(xerr.CodeCrateApp, "parse apk file failed"))
		os.Remove(tempFile)
		return
	}
	app.Encrypted = enc
	dir := filepath.Join(conf.Config.Data, app.AppId, app.Version)
	if !util.ExistFile(dir) {
		err = os.MkdirAll(dir, os.ModePerm)
		util.PanicIfErr(err)
	}
	filePath := dir + "/" + app.Name + ".apk"
	err = os.Rename(tempFile, filePath)
	util.PanicIfErr(err)
	if icon != nil {
		iconFile := dir + "/icon.png"
		f, err := os.OpenFile(iconFile, os.O_CREATE|os.O_RDWR, os.ModePerm)
		if err == nil {
			defer f.Close()
			png.Encode(f, icon)
		}
	}
	h, err := crypto.GetFileHash(filePath, "MD5")
	util.PanicIfErr(err)
	app.Hash = util.Bytes2HexStr(h)
	exist, err := insertApp(app)
	if err != nil {
		os.Remove(filePath)
		util.PanicIfErr(err)
	}
	if exist {
		c.JSON(http.StatusBadRequest, entity.NewRespErr(xerr.CodeCrateApp, "apk already exist"))
		return
	}
	c.JSON(http.StatusOK, entity.NewResp().AddResult("app", app))
}

func insertApp(app *entity.App) (exist bool, err error) {
	daoConn := dao.NewDao()
	err = daoConn.BeginTx()
	if err != nil {
		return false, err
	}
	existApp := &entity.Update{}
	smt := `SELECT apps.app_id, updates.latest FROM apps LEFT JOIN updates ON apps.app_id=updates.app_id WHERE updates.app_id=?`
	err = daoConn.Get(existApp, smt, app.AppId)
	if err != nil && err != sql.ErrNoRows {
		daoConn.RollBack()
		return false, err
	}
	if existApp != nil && existApp.AppId != "" {
		daoConn.RollBack()
		return true, nil
	}
	smt = `INSERT INTO apps(app_id,version,version_code,name,size,hash,encrypted,url)
						VALUES (:app_id,:version,:version_code,:name,:size,:hash,:encrypted,:url)`
	_, err = daoConn.CreateObj(smt, app)
	if err != nil {
		daoConn.RollBack()
		return false, err
	}

	//need update latest app version
	latestApp := &entity.Update{app.AppId, app.VersionCode}
	if existApp.Latest <= 0 {
		smt := `INSERT INTO updates(app_id,latest) VALUES (:app_id,:latest)`
		_, err = daoConn.CreateObj(smt, latestApp)
	} else if app.VersionCode > existApp.Latest {
		smt := `UPDATE updates SET latest=:latest WHERE updates.app_id=':app_id'`
		_, err = daoConn.UpdateObj(smt, latestApp)
	}
	if err != nil {
		daoConn.RollBack()
		return false, err
	} else {
		err = daoConn.Commit()
		if err != nil {
			return false, err
		}
	}
	return false, nil
}

func QueryApps(c *gin.Context) {
	uuid := c.Query("uuid")
	if uuid == "" {
		c.JSON(http.StatusBadRequest, entity.NewRespErr(xerr.CodeQueryApp, "illegal request"))
		return
	}
	latest, _ := strconv.ParseBool(c.Query("latest"))
	app_id := c.Query("app_id")
	appDao := dao.NewAppDao()
	//查询所有app
	if !latest {
		var apps []entity.App
		var err error
		if app_id != "" {
			apps, err = appDao.GetAllApp(app_id)
		} else {
			apps, err = appDao.GetAllApps()
		}
		util.PanicIfErr(err)
		fillAppUrl(uuid, apps)
		c.JSON(http.StatusOK, entity.NewResp().AddResult("apps", apps))
		return
	} else {
		if app_id != "" {
			app, err := appDao.GetLatestApp(app_id)
			util.PanicIfErr(err)
			fillAppUrl(uuid, app)
			c.JSON(http.StatusOK, entity.NewResp().AddResult("apps", app))
		} else {
			apps, err := appDao.GetLatestApps()
			util.PanicIfErr(err)
			fillAppUrl(uuid, apps)
			c.JSON(http.StatusOK, entity.NewResp().AddResult("apps", apps))
		}
		return
	}
}

func fillAppUrl(uuid string, apps []entity.App) {
	if apps == nil {
		return
	}
	verifyCode := dao.Redis().Get("app:" + uuid).Val()
	if verifyCode == "" {
		verifyCode = util.RandomCharacters(32)
		success, err := dao.Redis().SetNX("app:"+uuid, verifyCode, time.Minute*30).Result()
		util.PanicIfErr(err)
		if !success {
			xlog.Error().Msg("generate app download verify code failed")
			return
		}
	}
	c := conf.Config
	protol := "https"
	if !c.TLS {
		protol = "http"
	}
	for i := 0; i < len(apps); i++ {
		apps[i].Url = fmt.Sprintf(`%s://%s:%d/api/resource/app/downloads/%s/%s/%s.apk?uuid=%s&c=%s`,
			protol, c.Domain, c.Port, apps[i].AppId, apps[i].Version, apps[i].Name, uuid, verifyCode)
		apps[i].Icon = fmt.Sprintf(`%s://%s:%d/api/resource/app/downloads/%s/%s/icon.png?uuid=%s&c=%s`,
			protol, c.Domain, c.Port, apps[i].AppId, apps[i].Version, uuid, verifyCode)
	}
}

func DownloadApp(c *gin.Context) {
	uuid := c.Query("uuid")
	if uuid == "" {
		c.JSON(http.StatusBadRequest, entity.NewRespErr(xerr.CodeQueryApp, "illegal request"))
		return
	}
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.NewRespErr(xerr.CodeQueryApp, "illegal request"))
		return
	}
	verifyCode := c.Query("c")
	if verifyCode == "" {
		c.JSON(http.StatusBadRequest, entity.NewRespErr(xerr.CodeDownloadApp, "illegal request"))
		return
	}
	code := dao.Redis().Get("app:" + uuid).Val()
	if verifyCode != code {
		c.JSON(http.StatusBadRequest, entity.NewRespErr(xerr.CodeDownloadApp, "illegal request"))
		return
	}
	appDao := dao.NewAppDao()
	app, err := appDao.GetAppById(id)
	util.PanicIfErr(err)
	if app == nil {
		c.JSON(http.StatusNotFound, entity.NewRespErr(xerr.CodeDownloadApp, "app not found"))
		return
	}
	dir := filepath.Join(conf.Config.Data, app.AppId, app.Version)
	filePath := dir + "/" + app.Name + ".apk"
	if !util.ExistFile(filePath) {
		c.JSON(http.StatusNotFound, entity.NewRespErr(xerr.CodeDownloadApp, "app not found"))
		return
	}
	if app.Encrypted {
		c.File(filePath)
	} else {
		c.File(filePath)
	}
}

func ParseApkFromReader(r io.ReaderAt, size int64) (*entity.App, image.Image, error) {
	Apk, err := apk.OpenZipReader(r, size)
	if err != nil {
		return nil, nil, err
	}
	var icon image.Image
	app := &entity.App{}
	app.AppId = Apk.PackageName()
	icon, err = Apk.Icon(nil)
	if err != nil {
		icon = nil
	}
	manifest := Apk.Manifest()
	app.Name, _ = Apk.Label(nil)
	app.Version = manifest.VersionName.MustString()
	app.VersionCode = manifest.VersionCode.MustInt32()
	app.Size = size
	return app, icon, nil
}

func ParseApkFromFile(filePath string) (*entity.App, image.Image, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()
	fi, _ := f.Stat()
	size := fi.Size()
	return ParseApkFromReader(f, size)
}
