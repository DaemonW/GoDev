package controller

import (
	"daemonw/conf"
	"daemonw/crypto"
	"daemonw/dao"
	"daemonw/entity"
	"daemonw/util"
	"daemonw/xerr"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/shogo82148/androidbinary/apk"
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

	app, err := ParseApkFromFile(tempFile)
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
	h, err := crypto.GetFileHash(filePath, "MD5")
	util.PanicIfErr(err)
	app.Hash = util.Bytes2HexStr(h)
	appDao := dao.NewAppDao()
	if err = appDao.CreateApp(app); err != nil {
		os.Remove(filePath)
		util.PanicIfErr(err)
	}
	c.JSON(http.StatusOK, entity.NewResp().AddResult("app", app))
}

func insertApp(app entity.App) error {
	daoConn := dao.NewDao()
	err := daoConn.BeginTx()
	if err != nil {
		return err
	}
	type AppVersion struct {
		AppId  string `db:"app_id"`
		Latest string `db:"latest"`
	}
	version := &AppVersion{}
	smt := `SELECT apps.app_id,updates.latest FROM apps RIGHT JOIN updates ON apps.app_id=updates.app_id WHERE apps.app_id=? AND apps.version=?`
	err = daoConn.Get(version, smt, app.AppId, app.Version)
	if err != nil && err != sql.ErrNoRows {
		daoConn.RollBack()
		return err
	}
	if version.AppId != "" {
		daoConn.RollBack()
		return errors.New("app already exists")
	}
	smt = `INSERT INTO apps(id,app_id,version,name,size,hash,encrypted,url)
						VALUES (:id,:app_id,:version,:name,:size,:hash,:encrypted,:url) RETURNING id`
	_, err = daoConn.CreateObj(smt, app)
	if err != nil {
		daoConn.RollBack()
		return err
	}
	return nil
}

func QueryApp(c *gin.Context) {
	uuid := c.Query("uuid")
	if uuid == "" {
		c.JSON(http.StatusBadRequest, entity.NewRespErr(xerr.CodeQueryApp, "illegal request"))
		return
	}
	appId := c.Query("app_id")
	appDao := dao.NewAppDao()
	if appId == "" {
		apps, err := appDao.GetAllApps()
		util.PanicIfErr(err)
		c.JSON(http.StatusOK, entity.NewResp().AddResult("apps", apps))
		return
	} else {
		version := c.Query("version")
		if version == "" {
			c.JSON(http.StatusBadRequest, entity.NewRespErr(xerr.CodeQueryApp, "illegal app version"))
			return
		}
		app, err := appDao.GetApp(appId, version)
		util.PanicIfErr(err)
		c.JSON(http.StatusOK, entity.NewResp().AddResult("app", app))
		return
	}
}

func DownloadApp(c *gin.Context) {
	uuid := c.Query("uuid")
	if uuid == "" {
		c.JSON(http.StatusBadRequest, entity.NewRespErr(xerr.CodeQueryApp, "illegal request"))
		return
	}
	verifyCode := c.Query("c")
	if verifyCode == "" {
		c.JSON(http.StatusBadRequest, entity.NewRespErr(xerr.CodeDownloadApp, "illegal request"))
		return
	}
	appId := c.Query("app_id")
	version := c.Query("version")
	name := c.Query("name")
	code := dao.Redis().Get(appId + ":" + version).String()
	if verifyCode != code {
		c.JSON(http.StatusBadRequest, entity.NewRespErr(xerr.CodeDownloadApp, "illegal request"))
		return
	}
	appDao := dao.NewAppDao()
	app, err := appDao.GetApp(appId, version)
	util.PanicIfErr(err)
	if app == nil {
		c.JSON(http.StatusNotFound, entity.NewRespErr(xerr.CodeDownloadApp, "app not found"))
		return
	}
	dir := filepath.Join(conf.Config.Data, appId, version)
	filePath := dir + "/" + name
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

func version2Code(version string) (int64, error) {
	fields := strings.Split(version, ".")
	if len(fields) < 3 {
		return 0, errors.New("illegal version code")
	}
	h, err := strconv.ParseInt(fields[0], 10, 32)
	if err != nil {
		return 0, err
	}
	m, err := strconv.ParseInt(fields[1], 10, 32)
	if err != nil {
		return 0, err
	}
	l, err := strconv.ParseInt(fields[2], 10, 32)
	if err != nil {
		return 0, err
	}
	return h*1e+08 + m*1e+04 + l, nil
}

func ParseApkFromReader(r io.ReaderAt, size int64) (*entity.App, error) {
	Apk, err := apk.OpenZipReader(r, size)
	if err != nil {
		return nil, err
	}
	app := &entity.App{}
	app.AppId = Apk.PackageName()
	manifest := Apk.Manifest()
	app.Name, _ = Apk.Label(nil)
	app.Version = manifest.VersionName.MustString()
	app.VersionCode = manifest.VersionCode.MustInt32()
	app.Size = size
	return app, nil
}

func ParseApkFromFile(filePath string) (*entity.App, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	fi, _ := f.Stat()
	size := fi.Size()
	return ParseApkFromReader(f, size)
}
