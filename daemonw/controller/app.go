package controller

import (
	"crypto/cipher"
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
	"github.com/pkg/errors"
	"github.com/shogo82148/androidbinary/apk"
	"image"
	"image/png"
	"io"
	"io/ioutil"
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
	t := c.Request.Header.Get("Content-Type")
	if !strings.HasPrefix(t, gin.MIMEMultipartPOSTForm) {
		fmt.Println(t)
		c.JSON(http.StatusBadRequest, entity.NewRespErr(xerr.CodeCrateApp, "unacceptable content type"))
		return
	}
	//b, _ := ioutil.ReadAll(c.Request.Body)
	//fmt.Println(string(b))
	err := c.Request.ParseMultipartForm(32 << 20)
	if err != nil {
		c.JSON(http.StatusInternalServerError, entity.NewRespErr(xerr.CodeCrateApp, "parse multipart body failed"))
		return
	}
	encrypted := c.PostForm("encrypted")
	name := c.PostForm("name")
	enc, _ := strconv.ParseBool(encrypted)
	r, _, err := c.Request.FormFile("apk")
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.NewRespErr(xerr.CodeCrateApp, "parse file data failed"))
		return
	}
	defer r.Close()
	stamp := time.Now().UnixNano()
	resFolder := filepath.Join(conf.Config.Data, "res")
	apkFolder := filepath.Join(conf.Config.Data, "apk")
	tempFile := filepath.Join(apkFolder, strconv.FormatInt(stamp, 10))
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
	if name != "" {
		app.Name = name
	}
	app.Encrypted = enc
	apkDir := filepath.Join(apkFolder, app.AppId, app.Version)
	if !util.ExistFile(apkDir) {
		err = os.MkdirAll(apkDir, os.ModePerm)
		util.PanicIfErr(err)
	}
	filePath := apkDir + "/" + app.Name + ".apk"
	//err = os.Rename(tempFile, filePath)
	err = encryptApp(tempFile, filePath)
	os.Remove(tempFile)
	if err != nil {
		os.Remove(filePath)
		panic(err)
	}
	iconDir := filepath.Join(resFolder, app.AppId, app.Version)
	if !util.ExistFile(iconDir) {
		err = os.MkdirAll(iconDir, os.ModePerm)
		util.PanicIfErr(err)
	}
	if icon != nil {
		iconFile := iconDir + "/icon.png"
		f, err := os.OpenFile(iconFile, os.O_CREATE|os.O_RDWR, os.ModePerm)
		if err == nil {
			defer f.Close()
			png.Encode(f, icon)
		}
	} else {
		r, _, err := c.Request.FormFile("icon")
		if err != nil {
			if err != http.ErrMissingFile {
				c.JSON(http.StatusBadRequest, entity.NewRespErr(xerr.CodeCrateApp, "parse icon data failed"))
				return
			}
		} else {
			iconFile := iconDir + "/icon.png"
			f, err := os.OpenFile(iconFile, os.O_CREATE|os.O_RDWR, os.ModePerm)
			if err == nil {
				defer f.Close()
				io.Copy(f, r)
			}
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
	AppInfoTaskLock.Lock()
	if _, ok := AppInfoSpiderTask[app.Id]; !ok {
		AppInfoSpiderChan <- *app
	}
	defer AppInfoTaskLock.Unlock()
	c.JSON(http.StatusOK, entity.NewResp().AddResult("app", app))
}

func encryptApp(originPath, newPath string) error {
	r, err := os.OpenFile(originPath, os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		return err
	}
	defer r.Close()
	w, err := os.OpenFile(newPath, os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		return err
	}
	defer w.Close()
	keyIv := crypto.RandomBytes(48)
	b, err := ioutil.ReadFile(conf.Config.RSAPublic)
	if err != nil {
		return err
	}
	b, err = util.HexStr2Bytes(string(b))
	if err != nil {
		panic(err)
	}
	publicKey, err := crypto.ParsePublicKey(b)
	if err != nil {
		return err
	}
	encKeyIv, err := crypto.RsaEnc(publicKey, keyIv)
	if err != nil {
		return err
	}
	_, err = w.Write(encKeyIv)
	if err != nil {
		return err
	}
	success := crypto.AesStream(r, w, keyIv[:32], keyIv[32:], "CTR", true)
	if !success {
		return errors.New("encrypt file failed")
	}
	return nil
}

func insertApp(app *entity.App) (exist bool, err error) {
	daoConn := dao.NewDao()
	err = daoConn.BeginTx()
	if err != nil {
		return false, err
	}
	existApp := &entity.App{}
	smt := `SELECT * from apps WHERE app_id=? AND version_code=?`
	err = daoConn.Get(existApp, smt, app.AppId, app.VersionCode)
	if err != nil {
		if err == sql.ErrNoRows {
			existApp = nil
		} else {
			daoConn.RollBack()
			return false, err
		}
	}
	//if exist,return
	if existApp != nil {
		daoConn.RollBack()
		return true, nil
	}

	//insert new app and version
	smt = `INSERT INTO apps(app_id,version,version_code,name,size,hash,encrypted,url,create_at)
						VALUES (:app_id,:version,:version_code,:name,:size,:hash,:encrypted,:url,:create_at) RETURNING id`
	app.CreateAt = time.Now()
	rows, err := daoConn.NamedQuery(smt, app)
	if err != nil {
		daoConn.RollBack()
		return false, err
	}
	for rows.Next() {
		err = rows.Scan(&app.Id)
		util.PanicIfErr(err)
	}

	var latest int32 = -1
	err = daoConn.Get(&latest, `SELECT COALESCE(MAX(version_code),-1) from apps WHERE app_id=?`, app.AppId)
	if err != nil {
		daoConn.RollBack()
		return false, err
	}

	if app.VersionCode > latest {
		latest = app.VersionCode
	}
	//need update latest app version
	smt = `INSERT INTO updates(app_id,latest) VALUES (?,?) ON CONFLICT(app_id) DO UPDATE SET latest=?`
	_, err = daoConn.Create(smt, app.AppId, latest, latest)

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
	protocol := "https"
	if !c.TLS {
		protocol = "http"
	}
	for i := 0; i < len(apps); i++ {
		apps[i].Url = fmt.Sprintf(`%s://%s:%d/api/download/app/%d?uuid=%s&c=%s`,
			protocol, c.Domain, c.Port, apps[i].Id, uuid, verifyCode)
		dir := filepath.Join(conf.Config.Data, "res", apps[i].AppId, apps[i].Version)
		if !util.ExistFile(dir + `/icon.png`) {
			apps[i].Icon = ""
		} else {
			apps[i].Icon = fmt.Sprintf(`%s://%s:%d/api/app/resources/%s/%s/icon.png?uuid=%s&c=%s`,
				protocol, c.Domain, c.Port, apps[i].AppId, apps[i].Version, uuid, verifyCode)
		}
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
	code := dao.Redis().Get("app:" + uuid).Val()
	if verifyCode != code {
		c.JSON(http.StatusBadRequest, entity.NewRespErr(xerr.CodeDownloadApp, "illegal request"))
		return
	}
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.NewRespErr(xerr.CodeQueryApp, "illegal request"))
		return
	}
	appDao := dao.NewAppDao()
	app, err := appDao.GetAppById(id)
	util.PanicIfErr(err)
	if app == nil {
		c.JSON(http.StatusNotFound, entity.NewRespErr(xerr.CodeDownloadApp, "app not found"))
		return
	}
	filePath := filepath.Join(conf.Config.Data, "apk", app.AppId, app.Version, app.Name+".apk")
	if !util.ExistFile(filePath) {
		c.JSON(http.StatusNotFound, entity.NewRespErr(xerr.CodeDownloadApp, "apk not found"))
		return
	}
	c.Writer.Header().Set("Content-Type", "application/vnd.android.package-archive")
	if app.Encrypted {
		c.File(filePath)
	} else {
		fr, err := os.Open(filePath)
		if err != nil {
			panic(err)
		}
		fi, _ := fr.Stat()
		c.Writer.Header().Set("Content-Length", strconv.FormatInt(fi.Size()-256, 10))
		defer fr.Close()
		lr := &io.LimitedReader{fr, 256}
		encKeyIv, err := ioutil.ReadAll(lr)
		if err != nil {
			panic(err)
		}
		pr, err := os.Open(conf.Config.RSAPrivate)
		if err != nil {
			panic(err)
		}
		defer pr.Close()
		b, err := ioutil.ReadAll(pr)
		if err != nil {
			panic(err)
		}
		b, err = util.HexStr2Bytes(string(b))
		if err != nil {
			panic(err)
		}
		privateKey, err := crypto.ParsePrivateKey(b)
		if err != nil {
			panic(err)
		}
		keyIv, err := crypto.RsaDec(privateKey, encKeyIv)
		if err != nil {
			panic(err)
		}
		stream, err := crypto.CtrCipher(keyIv[:32], keyIv[32:])
		if err != nil {
			panic(err)
		}
		cr := cipher.StreamReader{stream, fr}
		_, err = io.Copy(c.Writer, cr)
		if err != nil {
			panic(err)
		}
	}
}

func GetAppInfo(c *gin.Context) {
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
	appDao := dao.NewAppDao()
	info, err := appDao.GetAppInfoById(id)
	util.PanicIfErr(err)
	if info == nil {
		c.JSON(http.StatusNotFound, entity.NewRespErr(xerr.CodeDownloadApp, "app not found"))
		return
	}
	detailDir := filepath.Join(conf.Config.Data, "res", info.Package, info.Version, "details")
	fis := util.ListFilesInfo(detailDir)
	if fis != nil {
		verifyCode := dao.Redis().Get("app:" + uuid).Val()
		if verifyCode == "" {
			verifyCode = util.RandomCharacters(32)
			success, err := dao.Redis().SetNX("app:"+uuid, verifyCode, time.Minute*30).Result()
			util.PanicIfErr(err)
			if !success {
				xlog.Error().Msg("generate app resources verify code failed")
				return
			}
		}
		c := conf.Config
		protocol := "https"
		if !c.TLS {
			protocol = "http"
		}
		n := len(fis)
		details := make([]string, n)
		for i := 0; i < n; i++ {
			img := fis[i].Name()
			details[i] = fmt.Sprintf(`%s://%s:%d/api/app/resources/%s/%s/details/%s?uuid=%s&c=%s`,
				protocol, c.Domain, c.Port, info.Package, info.Version, img, uuid, verifyCode)
		}
		info.ImageDetail = strings.Join(details, ",")
	} else {
		app, err := appDao.GetAppById(id)
		util.PanicIfErr(err)
		if app != nil {
			AppInfoTaskLock.Lock()
			if _, ok := AppInfoSpiderTask[info.Id]; !ok {
				AppInfoSpiderChan <- *app
			}
			defer AppInfoTaskLock.Unlock()
		}
	}
	c.JSON(http.StatusOK, entity.NewResp().AddResult("info", info))
}

func DeleteApp(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.NewRespErr(xerr.CodeDeleteApp, "illegal request"))
		return
	}
	appDao := dao.NewAppDao()
	err = appDao.BeginTx()
	util.PanicIfErr(err)
	app, err := appDao.GetAppById(id)
	if err != nil {
		appDao.RollBack()
		panic(err)
	}
	if app == nil {
		return
	}
	err = appDao.DeleteApp(app.AppId, app.Version)
	if err != nil {
		appDao.RollBack()
		panic(err)
	}

	err = appDao.DeleteAppInfoById(id)
	if err != nil {
		appDao.RollBack()
		panic(err)
	}
	latest := -1
	err = appDao.Get(&latest, `SELECT COALESCE(MAX(version_code),-1) FROM apps WHERE app_id=?`, app.AppId)
	if err != nil {
		appDao.RollBack()
		panic(err)
	}
	if latest == -1 {
		_, err = appDao.Delete(`DELETE FROM updates WHERE app_id=?`, app.AppId)
		if err != nil {
			appDao.RollBack()
			panic(err)
		}
	} else {
		_, err = appDao.Exec(`UPDATE updates SET latest=? WHERE app_id=?`, latest, app.AppId)
		if err != nil {
			appDao.RollBack()
			panic(err)
		}
	}
	apkFolder := filepath.Join(conf.Config.Data, "apk", app.AppId, app.Version)
	err = os.RemoveAll(apkFolder)
	if err != nil {
		appDao.RollBack()
		panic(err)
	}

	resFolder := filepath.Join(conf.Config.Data, "res", app.AppId, app.Version)
	err = os.RemoveAll(resFolder)
	if err != nil {
		appDao.RollBack()
		panic(err)
	}

	err = appDao.Commit()
	if err != nil {
		appDao.RollBack()
		panic(err)
	}
	c.JSON(http.StatusOK, entity.NewResp().AddResult("msg", "delete app success"))
}

const (
	ScopeIcon      = "icon"
	ScopeName      = "name"
	ScopeEncrypted = "encrypted"
)

func UpdateApp(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.NewRespErr(xerr.CodeUpdateApp, "illegal request"))
		return
	}
	appDao := dao.NewAppDao()
	err = appDao.BeginTx()
	util.PanicIfErr(err)
	app, err := appDao.GetAppById(id)
	if err != nil {
		appDao.RollBack()
		panic(err)
	}
	if app == nil {
		appDao.Commit()
		return
	}
	_s := c.Query("scope")
	scopes := strings.Split(_s, ",")
	var name string
	var encrypted bool

	if hasScope(scopes, ScopeIcon) {
		err = c.Request.ParseMultipartForm(32 << 20)
		if err != nil {
			appDao.RollBack()
			xlog.Error().Msgf("err: %s", err.Error())
			c.JSON(http.StatusBadRequest, entity.NewRespErr(xerr.CodeUpdateApp, "parse icon file failed"))
			return
		}
		r, _, err := c.Request.FormFile("icon")
		if err != nil {
			appDao.RollBack()
			xlog.Error().Msgf("err: %s", err.Error())
			c.JSON(http.StatusBadRequest, entity.NewRespErr(xerr.CodeUpdateApp, "parse icon file failed"))
			return
		}
		defer r.Close()
		dir := filepath.Join(conf.Config.Data, "res", app.AppId, app.Version)
		iconFile := dir + "/icon.png"
		f, err := os.OpenFile(iconFile, os.O_RDWR|os.O_CREATE, os.ModePerm)
		if err != nil {
			appDao.RollBack()
			xlog.Error().Msgf("err: %s", err.Error())
			c.JSON(http.StatusBadRequest, entity.NewRespErr(xerr.CodeUpdateApp, "update icon failed"))
			return
		}
		defer f.Close()
		_, err = io.Copy(f, r)
		if err != nil {
			appDao.RollBack()
			os.Remove(iconFile)
			panic(err)
		}
	}

	if hasScope(scopes, ScopeName) {
		name = c.PostForm("name")
		_, err := appDao.Exec(`UPDATE apps SET name=? WHERE id=?`, name, app.Id)
		if err != nil {
			appDao.RollBack()
			panic(err)
		}
	}

	if hasScope(scopes, ScopeEncrypted) {
		encrypted, _ = strconv.ParseBool(c.PostForm("encrypted"))
		_, err := appDao.Exec(`UPDATE apps SET encrypted=? WHERE id=?`, encrypted, app.Id)
		if err != nil {
			appDao.RollBack()
			panic(err)
		}
	}
	err = appDao.Commit()
	if err != nil {
		appDao.RollBack()
		panic(err)
	}
	c.JSON(http.StatusOK, entity.NewResp().AddResult("msg", "update app success"))

}

func hasScope(scopes []string, scope string) bool {
	if scopes == nil || len(scopes) == 0 {
		return false
	}
	for _, s := range scopes {
		if s == scope {
			return true
		}
	}
	return false
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
