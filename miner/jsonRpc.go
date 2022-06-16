package miner

import (
	"bytes"
	entity2 "chia-miner/miner/entity"
	"chia-miner/pkg/config"
	"compress/gzip"
	"encoding/hex"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"sync/atomic"
	"time"
)

var singleJsonRpc *JsonRpc

func InitJsonRpc(cfg *config.Config) {
	singleJsonRpc = &JsonRpc{
		miningInfoClient: &http.Client{},
		submitClient:     &http.Client{},
		cfg:              cfg,
	}
}

func GetJsonRpc() *JsonRpc {
	return singleJsonRpc
}

type JsonRpc struct {
	cfg              *config.Config
	jsonRpcId        int64
	miningInfoClient *http.Client
	submitClient     *http.Client
}

func (j *JsonRpc) GetMiningInfo() (*entity2.MiningInfo, error) {
	raw, err := j.call("pos_getMiningInfo", []interface{}{}, nil, j.miningInfoClient, nil)
	if err != nil {
		return nil, err
	}
	ret := gjson.Parse(raw)
	challenge, _ := hex.DecodeString(ret.Get("result.challenge").String())
	if len(challenge) != 32 {
		return nil, ErrBadMiningInfo
	}
	miningInfo := &entity2.MiningInfo{
		Height:         uint32(ret.Get("result.height").Int()),
		Challenge:      challenge,
		ReceiveTime:    time.Now().Unix(),
		Difficulty:     ret.Get("result.difficulty").Uint(),
		Epoch:          ret.Get("result.epoch").Int(),
		FilterBits:     int(ret.Get("result.filter_bits").Int()),
		ServerTime:     ret.Get("result.now").Int(),
		ScanIterations: ret.Get("result.scan_iterations").Int(),
		BestQuality:    math.MaxInt64,
	}
	return miningInfo, nil
}

func (j *JsonRpc) Submit(info *entity2.SubmitProof) {
	raw, err := j.call("pos_submitProof", info, nil, j.submitClient, nil)
	if err != nil {
		logrus.Errorf("submitProof fail, error %v", err)
		return
	}
	logrus.Infof("submitProof %v", raw)
}

// CallJsonRpc call json rpc method
func (j *JsonRpc) call(method string, params interface{}, result interface{}, rpcClient *http.Client, headers map[string]string) (raw string, err error) {
	var body io.Reader
	if data, err := json.Marshal(map[string]interface{}{"jsonrpc": "2.0", "method": method, "params": params, "id": atomic.AddInt64(&j.jsonRpcId, 1)}); err != nil {
		return "", err
	} else {
		logrus.Tracef("Call %s by %s", method, string(data))
		body = bytes.NewReader(data)
	}

	request, err := http.NewRequest(http.MethodPost, j.cfg.Rpc.Url, body)
	if err != nil {
		return "", err
	}
	request.Header.Set("Accept-Encoding", "gzip, deflate, br")
	for k, v := range headers {
		request.Header.Set(k, v)
	}
	request.Header.Set("Content-Type", "application/json; charset=utf-8")
	request.Header.Set("Authorization", j.cfg.GetAuthorizationToken())
	if response, err := rpcClient.Do(request); err != nil {
		return "", err
	} else if response.StatusCode == http.StatusUnauthorized {
		return "", ErrUnauthenticated
	} else if response.StatusCode == http.StatusForbidden {
		return "", ErrPermissionDenied
	} else {
		if body, err := readResponseBody(response); err != nil {
			return "", errors.Wrapf(err, "status code %d", response.StatusCode)
		} else {
			resp := struct {
				Error *struct {
					Code    int
					Message string
				}
			}{}
			if err := json.Unmarshal(body, &resp); err != nil {
				return string(body), err
			}
			if resp.Error != nil && (resp.Error.Code != 0 || resp.Error.Message != "") {
				return string(body), RawRpcError{
					code:    resp.Error.Code,
					message: resp.Error.Message,
				}
			}

			if response.StatusCode != http.StatusOK {
				return string(body), errors.Wrapf(ErrBadRequest, "status code %d", response.StatusCode)
			}

			// success
			if result != nil {
				jr := struct{ Result interface{} }{Result: result}
				if err := json.Unmarshal(body, &jr); err != nil {
					return string(body), err
				}
				if jr.Result == nil {
					return string(body), ErrNotFoundData
				}
			}
			return string(body), nil
		}
	}
}
func readResponseBody(response *http.Response) ([]byte, error) {
	defer response.Body.Close()

	// check gzip
	reader := response.Body
	if response.Header.Get("Content-Encoding") == "gzip" {
		r, err := gzip.NewReader(response.Body)
		if err != nil {
			return nil, err
		}
		defer r.Close()
		reader = r
	}

	// read
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	return data, nil
}
