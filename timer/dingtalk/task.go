package dingtalk

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"io"
	"net"
	"net/http"
	"net/url"
	"poem-bot/global"
	"poem-bot/model"
	"time"
)

type DingTalk struct {
	robotToken string
	secret     string
}

const (
	BaseUrl            = "https://oapi.dingtalk.com/robot/send?"
	defaultDialTimeout = 2 * time.Second
	defaultKeepAlive   = 2 * time.Second
)

var (
	myHTTPClient *http.Client
	globalD      *DingTalk
)

func init() {
	myHTTPClient = initDefaultHTTPClient()
}

func InitDingTalk() *DingTalk {
	global.LOG.Info("初始化DingTalk配置")
	robotToken := global.VP.GetString("bot.dingtalk.robot-token")
	secret := global.VP.GetString("bot.dingtalk.secret")
	dt := &DingTalk{
		robotToken: robotToken,
		secret:     secret,
	}
	return dt
}

func initDefaultHTTPClient() *http.Client {
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   defaultDialTimeout,
				KeepAlive: defaultKeepAlive,
			}).DialContext,
		},
		Timeout: defaultDialTimeout,
	}
	return client
}

func SendPoem() {
	global.LOG.Info("钉钉定时机器人启动")
	if globalD == nil {
		globalD = InitDingTalk()
	}

	poem := model.GetRandomPoem()
	msg := model.NewTextMsg(poem)

	global.LOG.Info("向机器人发送内容, ", zap.String("内容", poem))
	err := globalD.sendTextMsg(msg)
	if err != nil {
		global.LOG.Error("SendPoem error", zap.Error(err))
	}
}

func (d *DingTalk) sendTextMsg(msg model.IDingMsg) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return d.sendMessage(ctx, msg)
}

func (d *DingTalk) sendMessage(ctx context.Context, msg model.IDingMsg) error {
	var (
		uri  string
		resp *http.Response
		err  error
	)

	value := url.Values{}
	value.Set("access_token", d.robotToken)
	if d.secret != "" {
		t := time.Now().UnixNano() / 1e6
		value.Set("timestamp", fmt.Sprintf("%d", t))
		value.Set("sign", d.sign(t, d.secret))

	}

	uri = BaseUrl + value.Encode()
	header := map[string]string{
		"Content-type": "application/json",
	}
	resp, err = doRequest(ctx, "POST", uri, header, msg.Parse())

	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("send msg err. http code: %d, token: %s, msg: %s",
			resp.StatusCode, d.robotToken, msg.Parse())
	}
	body, _ := io.ReadAll(resp.Body)
	global.LOG.Info("钉钉返回: " + string(body))

	var respMsg model.ResponseMsg
	err = json.Unmarshal(body, &respMsg)
	if err != nil {
		return err
	}
	if respMsg.ErrCode != 0 {
		return fmt.Errorf("send msg err. err msg: %s", respMsg.ErrMsg)
	}
	return nil
}

func (d *DingTalk) sign(t int64, secret string) string {
	strToHash := fmt.Sprintf("%d\n%s", t, secret)
	hmac256 := hmac.New(sha256.New, []byte(secret))
	hmac256.Write([]byte(strToHash))
	data := hmac256.Sum(nil)
	return base64.StdEncoding.EncodeToString(data)
}

func doRequest(ctx context.Context, callMethod string, endPoint string, header map[string]string, body []byte) (*http.Response, error) {
	req, err := http.NewRequest(callMethod, endPoint, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	if header != nil && len(header) > 0 {
		for k, v := range header {
			req.Header.Set(k, v)
		}
	}
	req = req.WithContext(ctx)
	response, err := myHTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	if response == nil {
		return nil, fmt.Errorf("reponse is nil, please check it")
	}
	return response, nil
}
