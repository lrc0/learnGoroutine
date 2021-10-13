package main

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"gopkg.in/logger.v1"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"path"
	"regexp"
	// "strconv"
	"strings"
	"sync"
	"time"
)

const (
	// imageURLRegexp = `\S+\d+\.html`   //匹配图片网页的正则
	imageURLRegexp = `\d{5}\w+\.html` //匹配图片网页的正则
	// imageURLRegexp2 = `\d{5}+\.html`   //匹配图片网页的正则

	maxRetry      = 2
	maxPagePoller = 20 //总共处理page的goroutine数量
	maxDownloader = 20 //总共下载图片的goroutine数量
)

const (
	ready = iota
	done
	fail
)

//一个页面包含的要素
type page struct {
	url   string
	body  *[]byte
	retry int
}

//一张图片包含的要素
type image struct {
	imageURL string
	filename string //图片文件名
	folder   string //存放图片的文件夹名称
	retry    int
}

type context struct {
	locker        *sync.RWMutex
	imgMap        map[string]int
	imageChan     chan *image
	imgListMap    map[string]int
	imageListChan chan *page
	pageMap       map[string]int
	pageChan      chan *page
	parseChan     chan *page //待解析的网页channel
	rootPath      string
}

func main() {
	host := "http://www.27baola.net"
	ctx := initContext(host)
	start(ctx)
	// monitor(ctx)
	select {}
}

// func monitor(ctx *context) {
// 	for {
// 		if len(ctx.imageChan) == 0 && len(ctx.pageChan) == 0 && len(ctx.parseChan) == 0 {
// 			os.Exit(-1)
// 		}
// 	}
// }

func initContext(host string) *context {
	filename := strings.TrimLeft(host, "http://")
	creatfolder(filename)

	return &context{
		locker:    new(sync.RWMutex),
		imgMap:    make(map[string]int),
		imageChan: make(chan *image, maxDownloader),
		pageMap:   make(map[string]int),
		pageChan:  make(chan *page, maxPagePoller),
		parseChan: make(chan *page, maxPagePoller),
		rootPath:  host + "/gif",
	}
}

func start(ctx *context) {
	client := httpClient(120)
	//获取html
	for i := 0; i < maxPagePoller; i++ {
		go func() {
			for {
				page := <-ctx.pageChan //不停地从pageChan里面读取page,直到读完阻塞
				page.pollPage(ctx, client)
			}
		}()
	}

	//下载图片
	for i := 0; i < maxDownloader; i++ {
		go func() {
			for {
				image := <-ctx.imageChan
				image.downloadImage(ctx, client)
			}
		}()
	}

	go func() {
		for {
			page := <-ctx.parseChan
			page.parsePage(ctx)
		}
	}()

	ctx.pageChan <- &page{ctx.rootPath, nil, 0}
}

func httpClient(tout int) *http.Client {
	timeout := time.Duration(tout) * time.Second

	// 增加timeout
	hClient := &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				c, err := net.DialTimeout(netw, addr, timeout)
				if err != nil {
					log.Errorf("dial timeout: %v, message: %s", timeout, err.Error())
					return nil, err
				}
				c.SetDeadline(time.Now().Add(timeout))
				return c, nil
			},
		},
	}

	return hClient
}

//抓取html
func (p *page) pollPage(ctx *context, client *http.Client) {

	ctx.locker.RLock()
	if ctx.pageMap[p.url] == done {
		ctx.locker.RUnlock()
		return
	}
	ctx.locker.RUnlock()
	defer p.pageRetry(ctx)

	req, err := http.NewRequest("GET", p.url, nil)
	if err != nil {
		log.Error(err)
		return
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Error(err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
		return
	}
	ctx.locker.Lock()
	ctx.pageMap[p.url] = done
	ctx.locker.Unlock()

	p.body = &body
	ctx.parseChan <- p
}

func (p *page) pageRetry(ctx *context) {
	ctx.locker.RLock()
	if ctx.pageMap[p.url] == done {
		ctx.locker.RUnlock()
		return
	}
	ctx.locker.RUnlock()

	if p.retry++; p.retry < maxRetry { //重试两次
		go func() {
			ctx.pageChan <- p
		}()
	} else {
		ctx.locker.Lock()
		ctx.pageMap[p.url] = fail
		ctx.locker.Unlock()
	}
}

func (p *page) parsePage(ctx *context) {
	b := matchImageURL(p.url)
	if b {
		p.findImage(ctx)
	} else {
		p.findPage(ctx)
	}
}

func matchImageURL(url string) bool {
	b, err := regexp.MatchString(imageURLRegexp, url)
	if err != nil {
		log.Error(err)
		return false
	}

	return b
}

//找到html里面的链接
func (p *page) findPage(ctx *context) {

	pageURL, err := url.Parse(p.url)
	if err != nil {
		log.Error(err)
		return
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(*p.body))
	if err != nil {
		log.Error(err)
		return
	}
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		attr, e := s.Attr("href")
		if !e || attr == "" {
			return
		}

		absURL := toAbs(pageURL, attr)
		if absURL == nil {
			return
		}

		absURL.Fragment = ""

		url := absURL.String()

		ctx.locker.RLock()
		_, exist := ctx.pageMap[url]
		ctx.locker.RUnlock()
		if !exist {
			ctx.locker.Lock()
			ctx.pageMap[url] = ready
			ctx.locker.Unlock()
			go func() {
				ctx.pageChan <- &page{url: url}
			}()
		}
	})
}

//找到html里面的gif图片地址列表
func (p *page) findImage(ctx *context) {

	filename := strings.TrimLeft(ctx.rootPath, "http://")
	folder := creatfolder(filename)

	pageURL, err := url.Parse(p.url)
	if err != nil {
		log.Error(err)
		return
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(*p.body))
	if err != nil {
		log.Error(err)
		return
	}
	doc.Find("img").Each(func(i int, s *goquery.Selection) {
		attr, e := s.Attr("src")
		if !e || attr == "" {
			return
		}

		absURL := toAbs(pageURL, attr)
		absURL.Fragment = ""

		url := absURL.String()

		ctx.locker.RLock()
		_, exist := ctx.imgMap[url]
		ctx.locker.RUnlock()
		if !exist {
			ctx.locker.Lock()
			ctx.imgMap[url] = ready
			ctx.locker.Unlock()
			if strings.Contains(url, "%20") {
				url = strings.Trim(url, "%20")
			}
			filename := path.Base(url)
			ctx.imageChan <- &image{imageURL: url, filename: filename, folder: folder, retry: 0}
		}
	})
}

func (i *image) downloadImage(ctx *context, client *http.Client) {

	ctx.locker.RLock()
	if ctx.imgMap[i.imageURL] == done {
		ctx.locker.RUnlock()
		return
	}
	ctx.locker.RUnlock()
	defer i.imageRetry(ctx)

	if path.Ext(i.imageURL) == ".jpg" {
		ctx.locker.Lock()
		ctx.imgMap[i.imageURL] = done
		ctx.locker.Unlock()
		return
	}

	err := i.recordURL(i.imageURL + "\n")
	if err != nil {
		log.Error(err)
		return
	}
	fmt.Println("Image url: ", i.imageURL)

	req, err := http.NewRequest("GET", i.imageURL, nil)
	if err != nil {
		log.Error(err)
		return
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Error(err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		ctx.locker.Lock()
		ctx.imgMap[i.imageURL] = done
		ctx.locker.Unlock()
		return
	}

	image := <-ctx.imageChan
	filename := image.folder + "/" + image.filename

	f, err := os.Create(filename)
	if err != nil {
		log.Error(err)
		return
	}
	defer f.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Info("url: ", i.imageURL)
		log.Error(err)
		return
	}

	err = ioutil.WriteFile(filename, body, 0666)
	if err != nil {
		log.Error(err)
		return
	}

	ctx.locker.Lock()
	ctx.imgMap[i.imageURL] = done
	ctx.locker.Unlock()
}

func (i *image) imageRetry(ctx *context) {
	ctx.locker.RLock()
	if ctx.pageMap[i.imageURL] == done {
		ctx.locker.RUnlock()
		return
	}
	ctx.locker.RUnlock()

	if i.retry++; i.retry < maxRetry { //重试两次
		go func() {
			ctx.imageChan <- i
		}()
	} else {
		ctx.locker.Lock()
		ctx.imgMap[i.imageURL] = fail
		ctx.locker.Unlock()
	}
}

func toAbs(u *url.URL, href string) *url.URL {
	buf := new(bytes.Buffer)
	if h := strings.ToLower(href); strings.Index(h, "http://") == 0 || strings.Index(h, "https://") == 0 {
		buf.WriteString(href)
	} else {
		buf.WriteString(u.Scheme)
		buf.WriteString("://")
		buf.WriteString(u.Host)

		switch href[0] {
		case '?':
			if len(u.Path) == 0 {
				buf.WriteString("/")
			} else {
				buf.WriteString(u.Path)
			}
			buf.WriteString(href)
		case '/':
			buf.WriteString(href)
		default:
			p := "/" + path.Dir(u.Path) + "/" + href
			buf.WriteString(path.Clean(p))
		}
	}

	rurl, err := url.Parse(buf.String())
	if err != nil {
		log.Error(err)
		return nil
	}
	return rurl
}

func creatfolder(filename string) string {
	err := os.Mkdir(filename, 0777)
	if os.IsExist(err) {
		return filename
	}
	if err != nil {
		log.Error(err)
		return ""
	}
	return filename
}

func (i *image) recordURL(url string) error {
	filename := i.folder + "/" + "record"
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0666)
	if os.IsNotExist(err) {
		f, err = os.Create(filename)
		if err != nil {
			log.Error(err)
			return err
		}
	}
	if err != nil {
		log.Error(err)
		return err
	}

	_, err = io.WriteString(f, url)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}
