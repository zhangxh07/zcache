package cache

import (
	"fmt"
	"google.golang.org/protobuf/proto"
	"io/ioutil"
	"net/http"
	"net/url"
	pb "zcache/zcachepb"
)

type httpGetter struct {
	bashURL string
}

// 向其他节点获取缓存值
//func (h *httpGetter) Get(group string, key string) ([]byte, error) {
//	u := fmt.Sprintf("%v%v/%v", h.bashURL, url.QueryEscape(group), url.QueryEscape(key))
//	res, err := http.Get(u)
//	if err != nil {
//		return nil, err
//	}
//	defer res.Body.Close()
//
//	if res.StatusCode != http.StatusOK {
//		return nil, fmt.Errorf("server returned: %v", res.Status)
//	}
//
//	bytes, err := ioutil.ReadAll(res.Body)
//	if err != nil {
//		return nil, fmt.Errorf("reading response body: %v", err)
//	}
//	return bytes, nil
//}

func (h *httpGetter) Get(in *pb.Request, out *pb.Response) error {
	u := fmt.Sprintf(
		"%v%v/%v",
		h.bashURL,
		url.QueryEscape(in.GetGroup()),
		url.QueryEscape(in.GetKey()),
	)
	res, err := http.Get(u)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned: %v", res.Status)
	}

	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("reading response body: %v", err)
	}

	if err = proto.Unmarshal(bytes, out); err != nil {
		return fmt.Errorf("decoding response body: %v", err)
	}

	return nil
}

var _ PeerGetter = (*httpGetter)(nil)
