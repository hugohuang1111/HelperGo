package utils

import (
	"errors"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/cavaliercoder/grab"
)

func DownloadFile(url string, localPath string, fb func(totalLen, downLen, incLen int64)) error {
	client := grab.NewClient()
	req, _ := grab.NewRequest(localPath, url)
	resp := client.Do(req)

	t := time.NewTicker(500 * time.Millisecond)
	defer t.Stop()

	var lastBytes int64 = 0
END:
	for {
		select {
		case <-t.C:
			fb(resp.Size, resp.BytesComplete(), resp.BytesComplete()-lastBytes)
			lastBytes = resp.BytesComplete()

		case <-resp.Done:
			break END
		}
	}

	return resp.Err()
}

func DownloadFile1(url string, localPath string, fb func(totalLen, downLen, incLen int64)) error {
	var (
		fsize   int64
		buf     = make([]byte, 32*1024)
		written int64
	)
	tmpFilePath := localPath + ".download"
	client := new(http.Client)
	client.Timeout = time.Second * 60 //设置超时时间

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:84.0) Gecko/20100101 Firefox/84.0")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	//读取服务器返回的文件大小
	lenStr := resp.Header.Get("Content-Length")
	fsize = -1
	if "" != lenStr {
		fsize, err = strconv.ParseInt(lenStr, 10, 32)
		if err != nil {
			return err
		}
	}
	//创建文件
	MakeSureDirExist(tmpFilePath)
	file, err := os.Create(tmpFilePath)
	if err != nil {
		return err
	}
	defer file.Close()
	if resp.Body == nil {
		return errors.New("body is null")
	}
	defer resp.Body.Close()

	for {
		//读取bytes
		nr, er := resp.Body.Read(buf)
		if nr > 0 {
			//写入bytes
			nw, ew := file.Write(buf[0:nr])
			//数据长度大于0
			if nw > 0 {
				written += int64(nw)
			}
			//写入出错
			if ew != nil {
				err = ew
				break
			}
			//读取是数据长度不等于写入的数据长度
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er != nil {
			if er != io.EOF {
				err = er
			}
			break
		}
		fb(fsize, written, int64(nr))
	}

	if nil != err {
		return err
	}
	if fsize != written {
		return errors.New("receive data length is not equal to fsize")
	}
	file.Close()
	err = os.Rename(tmpFilePath, localPath)

	return err
}

func IsHttp(s string) bool {
	if strings.HasPrefix(s, "http://") {
		return true
	}
	if strings.HasPrefix(s, "https://") {
		return true
	}

	return false
}
