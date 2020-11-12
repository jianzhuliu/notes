/*

 */
package gocache

import (
	"fmt"
	"gocache/consistenthash"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/golang/protobuf/proto"
	pb "gocache/gocachepb"
)

const (
	defaultBasePath = "/gocache/"
	defaultReplics  = 50 //默认虚拟节点数
)

type HTTPPool struct {
	baseUrl  string //域名(ip)+端口, http://127.0.0.1:9003
	basePath string //url path

	//添加节点选择
	mu          sync.Mutex
	hashObj     *consistenthash.HashObj
	httpGetters map[string]*httpGetter //节点 对应的 httpGetter
}

func NewHTTPPool(baseUrl string) *HTTPPool {
	return &HTTPPool{
		baseUrl:  baseUrl,
		basePath: defaultBasePath,
	}
}

func (p *HTTPPool) Log(format string, args ...interface{}) {
	log.Printf("[Server %s] %s", p.baseUrl, fmt.Sprintf(format, args...))
}

func (p *HTTPPool) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	//判断 url
	if !strings.HasPrefix(path, p.basePath) {
		p.Log("unexpected path : %s", path)
		http.Error(rw, "bad request", http.StatusBadRequest)
		return
	}

	p.Log("%s %s", r.Method, path)

	//基本参数校验
	// /<basePath>/<groupName>/<key>
	parts := strings.SplitN(path[len(p.basePath):], "/", 2)
	if len(parts) != 2 {
		http.Error(rw, "bad request", http.StatusBadRequest)
		return
	}

	groupName := strings.TrimSpace(parts[0])
	key := strings.TrimSpace(parts[1])

	if len(groupName) == 0 || len(key) == 0 {
		http.Error(rw, "group or key should be empty", http.StatusNotFound)
		return
	}

	group := GetGroup(groupName)
	if group == nil {
		http.Error(rw, "no such group", http.StatusNotFound)
		return
	}

	view, err := group.Get(key)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	//body := view.ByteSlice()
	//使用 protobuf 协议
	body, err := proto.Marshal(&pb.Response{Value: view.ByteSlice()})
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/octet-stream")
	rw.WriteHeader(http.StatusOK)
	rw.Write(body)
}

//设置分布式节点
func (p *HTTPPool) SetNodes(nodes ...string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.hashObj = consistenthash.New(defaultReplics, nil)
	p.hashObj.Add(nodes...)
	p.httpGetters = make(map[string]*httpGetter, len(nodes))

	for _, node := range nodes {
		p.httpGetters[node] = NewHttpGetter(node + p.basePath)
	}
}

//根据 key， 返回对应的 NodeGetter
func (p *HTTPPool) PickNode(key string) (NodeGetter, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()

	//本机节点除外
	if node := p.hashObj.Get(key); node != "" && node != p.baseUrl {
		p.Log("pick node %s", node)
		return p.httpGetters[node], true
	}

	return nil, false
}

//确保 *HTTPPool 实现了 NodePicker 接口
var _ NodePicker = (*HTTPPool)(nil)

//缓存数据来源，通过 http 方式远程获取
type httpGetter struct {
	baseUrl string
}

func NewHttpGetter(baseUrl string) *httpGetter {
	return &httpGetter{baseUrl: baseUrl}
}

func (h *httpGetter) Get(group string, key string) ([]byte, error) {
	//构造 url 格式， /<basePath>/<groupName>/<key>
	targetUrl := fmt.Sprintf("%v%v/%v", h.baseUrl, url.QueryEscape(group), url.QueryEscape(key))
	log.Println("going to get targetUrl ", targetUrl)

	resp, err := http.Get(targetUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	//判断状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server return %d", resp.StatusCode)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body %v", err)
	}

	//使用 protobuf 协议
	out := &pb.Response{}

	if err = proto.Unmarshal(data, out); err != nil {
		return nil, fmt.Errorf("decoding response body: %v", err)
	}

	return out.Value, nil
}
