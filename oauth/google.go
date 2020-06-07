package oauth

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/matac42/LiveShare/database"
	"golang.org/x/oauth2"
)

// Conf wraped oauth2.Config.
type Conf struct {
	oauth2.Config
}

// CredentialInfo store oauth2 access token etc...
type CredentialInfo struct {
	gorm.Model
	Token oauth2.Token
}

// APIInfo include api info from google.
type APIInfo struct {
	Web struct {
		ClientID                string   `json:"client_id"`
		ProjectID               string   `json:"project_id"`
		AuthURI                 string   `json:"auth_uri"`
		TokenURI                string   `json:"token_uri"`
		AuthProviderX509CertURL string   `json:"auth_provider_x509_cert_url"`
		ClientSecret            string   `json:"client_secret"`
		RedirectURIs            []string `json:"redirect_uris"`
		JSOrigins               []string `json:"javascript_origins"`
	}
}

// CreateAPIInfo create APIInfo based on apiinfo.json.
func CreateAPIInfo() APIInfo {
	apiInfo := APIInfo{}
	raw, err := ioutil.ReadFile("~/Downloads/apiinfo.json")
	if err != nil {
		fmt.Println(err.Error())
	}

	json.Unmarshal(raw, &apiInfo)

	return apiInfo
}

// CreateConf create oauth2 config structure.
func CreateConf() Conf {
	apiInfo := CreateAPIInfo()
	conf := Conf{
		oauth2.Config{
			ClientID:     apiInfo.Web.ClientID,
			ClientSecret: apiInfo.Web.ClientSecret,
			Scopes:       []string{"https://www.googleapis.com/auth/youtube.readonly"},
			Endpoint: oauth2.Endpoint{
				AuthURL:  apiInfo.Web.AuthURI,
				TokenURL: apiInfo.Web.TokenURI,
			},
			RedirectURL: apiInfo.Web.RedirectURIs[0],
		},
	}
	return conf
}

// Google get user email.
func (conf Conf) Google(c *gin.Context) {

	url := conf.AuthCodeURL("state", oauth2.AccessTypeOffline)
	c.Redirect(http.StatusMovedPermanently, url)

}

// CallBack fires fn when there was a connection to /callback.
func (conf Conf) CallBack(c *gin.Context) {
	code := GetAuthCode(c)
	tok := conf.GetAccessToken(c, code)
	// まずaccess tokenをdbに保存することから．
	cre := CredentialInfo{}
	cre.Token = *tok
	SaveCredentialInfo(cre)

	// info := conf.GetChannelInfo(c, tok)
	// // とりあえずprintしてから
	// fmt.Println(info)
	// infoから帰って来たインスタンスをmysqlに保存までやる
	// チャンネル名とチャンネルidくらいは先にとって一緒に保存しておきたい．

}

// GetAuthCode get credential code from google.
func GetAuthCode(c *gin.Context) string {
	code := c.Request.URL.Query().Get("code")
	err := c.Request.URL.Query().Get("error")
	if err != "" {
		fmt.Println(err)
	}
	// if _, err := fmt.Scan(&code); err != nil {
	// 	log.Fatal(err)
	// }

	return code
}

// GetAccessToken get accesstoken from google.
func (conf Conf) GetAccessToken(c *gin.Context, code string) *oauth2.Token {
	tok, err := conf.Exchange(c, code)
	if err != nil {
		log.Fatal(err)
	}

	return tok
}

// GetChannelInfo get ...
func (conf Conf) GetChannelInfo(c *gin.Context, tok *oauth2.Token) {
	client := conf.Client(c, tok)
	client.Get("...")

	// db保存用の構造体の初期化をここでして，いれてそれを返す感じ．
}

// SaveCredentialInfo is routines for storing credential info to the database.
func SaveCredentialInfo(cre CredentialInfo) {
	db, err := database.SQLConnect()
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	db.AutoMigrate(cre)
	error := db.Create(&cre).Error
	if error != nil {
		fmt.Println(error)
	} else {
		fmt.Println("success addition access token to db!!!")
	}
}
