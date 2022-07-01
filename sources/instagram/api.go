package instagram

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const (
	graphUrl            = "https://graph.instagram.com/v14.0/"
	IMAGE_TYPE          = "IMAGE"
	VIDEO_TYPE          = "VIDEO"
	CAROUSEL_ALBUM_TYPE = "CAROUSEL_ALBUM"
)

type Cursors struct {
	After  string `json:"after"`
	Before string `json:"before"`
}

type Paging struct {
	Next    string   `json:"next"`
	Cursors *Cursors `json:"cursors"`
}

type UserResponse struct {
	AccountType string `json:"account_type"`
	ID          string `json:"id"`
	MediaCount  int    `json:"media_count"`
	Username    string `json:"username"`
}

type MediaItem struct {
	Caption      string `json:"caption"`
	ID           string `json:"id"`
	MediaType    string `json:"media_type"`
	MediaUrl     string `json:"media_url"`
	Permalink    string `json:"permalink"`
	ThumbnailUrl string `json:"thumbnail_url"`
	Timestamp    string `json:"timestamp"`
	Username     string `json:"username"`
}

type PagingResponse struct {
	Data   []*MediaItem `json:"data"`
	Paging *Paging      `json:"paging"`
	cur    int
	next   int
	api    *InstagramApi
}

func (p *PagingResponse) Item() *MediaItem {
	return p.Data[p.cur]
}

func (p *PagingResponse) Next() bool {
	p.cur = p.next
	if len(p.Data) == p.cur {
		if p.Paging.Next == "" {
			return false
		}
		p.cur = 0
		p.next = 0
		if err := p.api.next(p.Paging.Next, p); err != nil {
			return false
		}
	}
	p.next++
	return true
}

type apiResponse map[string]interface{}

type InstagramApi struct {
	access_token string
}

func NewAPI(token string) *InstagramApi {
	return &InstagramApi{
		access_token: token,
	}
}

func (api *InstagramApi) Me(fields ...string) *UserResponse {
	params := url.Values{}
	if len(fields) > 0 {
		params.Set("fields", strings.Join(fields, ","))
	}
	r := &UserResponse{}
	api.get("me", params, r)
	return r
}

func (api *InstagramApi) MeMedia(fields ...string) (*PagingResponse, error) {
	return api.UserMedia("me", fields...)
}

func (api *InstagramApi) UserMedia(userID string, fields ...string) (*PagingResponse, error) {
	params := url.Values{}
	if len(fields) > 0 {
		params.Set("fields", strings.Join(fields, ","))
	}
	r := &PagingResponse{}
	err := api.get(userID+"/media", params, r)
	return r, err
}

func buildGetRequest(urlStr string, params url.Values) (*http.Request, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	// If we are getting, then we can't merge query params
	if params != nil {
		if u.RawQuery != "" {
			return nil, fmt.Errorf("Cannot merge query params in urlStr and params")
		}
		u.RawQuery = params.Encode()
	}

	return http.NewRequest(http.MethodGet, u.String(), nil)
}

func (api *InstagramApi) next(urlStr string, r interface{}) error {
	req, err := buildGetRequest(urlStr, nil)
	if err != nil {
		return err
	}
	return api.do(req, r)
}

func (api *InstagramApi) get(path string, params url.Values, r interface{}) error {
	u := graphUrl + path
	params.Set("access_token", api.access_token)
	req, err := buildGetRequest(u, params)
	if err != nil {
		return err
	}
	return api.do(req, r)
}

func (api *InstagramApi) do(req *http.Request, r interface{}) error {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		io.CopyN(ioutil.Discard, resp.Body, 512)
		resp.Body.Close()
	}()

	if resp.StatusCode == 401 {
		return fmt.Errorf("auth error %d", resp.StatusCode)
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("access error %d", resp.StatusCode)
	}

	return decodeResponse(resp.Body, r)
}

func decodeResponse(body io.Reader, to interface{}) error {
	err := json.NewDecoder(body).Decode(to)

	if err != nil {
		return fmt.Errorf("instagram: error decoding body; %s", err.Error())
	}
	return nil
}
