package client

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	dac "github.com/xinsnake/go-http-digest-auth-client"
)

type User struct {
	Country      string        `json:"country"`
	EmailAddress string        `json:"emailAddress"`
	FirstName    string        `json:"firstName"`
	ID           string        `json:"id"`
	LastName     string        `json:"lastName"`
	Roles        []Role        `json:"roles"`
	TeamIds      []interface{} `json:"teamIds"`
	Username     string        `json:"username"`
}

type NewUser struct {
	Roles    []string `json:"roles"`
	Username string   `json:"username"`
}

type NewReturnUser struct {
	InviterUsername string        `json:"inviterUsername"`
	OrgID           string        `json:"orgId"`
	OrgName         string        `json:"orgName"`
	TeamIds         []interface{} `json:"teamIds"`
	Username        string        `json:"username"`
}
type UpdateUser struct {
	Roles []Role `json:"roles"`
}
type Role struct {
	GroupID  string `json:"groupId,omitempty"`
	RoleName string `json:"roleName"`
	OrgID    string `json:"orgId,omitempty"`
}

var (
	Errors = make(map[int]string)
)

func init() {
	Errors[400] = "Bad Request, StatusCode = 400"
	Errors[404] = "User Does Not Exist , StatusCode = 404"
	Errors[409] = "User Already Exist, StatusCode = 409"
	Errors[401] = "Unautharized Access, StatusCode = 401"
	Errors[429] = "User Has Sent Too Many Request, StatusCode = 429"
}

type Client struct {
	publickey  string
	privateKey string
	orgid      string
	httpClient *http.Client
}

func NewClient(publickey string, privateKey string, orgid string) *Client {
	// publickey = os.Getenv("PUBLIC_KEY")
	// privateKey = os.Getenv("PRIVATE_KEY")
	// orgid = os.Getenv("ORGID")
	return &Client{
		publickey:  publickey,
		privateKey: privateKey,
		orgid:      orgid,
		httpClient: &http.Client{},
	}
}

func (c *Client) GetUser(username string) (*User, error) {
	fmt.Println("Get user")
	body, err := c.gethttpRequest(fmt.Sprintf("%v", username), "GET", bytes.Buffer{})
	if err != nil {
		log.Println("[READ ERROR]: ", err)
		return nil, err
	}
	user := &User{}
	err = json.NewDecoder(body).Decode(user)
	if err != nil {
		log.Println("[READ ERROR]: ", err)
		return nil, err
	}
	return user, nil
}
func (c *Client) gethttpRequest(username, method string, body bytes.Buffer) (closer io.ReadCloser, err error) {
	t := dac.NewTransport(c.publickey, c.privateKey)
	uri := "https://cloud.mongodb.com/api/atlas/v1.0/users/byName/" + username
	req, err := http.NewRequest(method, uri, nil)
	if err != nil {
		log.Fatalln(err)
	}
	res, err := t.RoundTrip(req)
	if err != nil {
		log.Println("[ERROR]: ", err)
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		respBody := new(bytes.Buffer)
		_, err := respBody.ReadFrom(res.Body)
		if err != nil {
			return nil, fmt.Errorf("Error : %v", Errors[res.StatusCode])
		}
		return nil, fmt.Errorf("Error : %v ", Errors[res.StatusCode])
	}
	return res.Body, nil
}

func (c *Client) AddNewUser(item *NewUser) (*NewReturnUser, error) {
	fmt.Println("New user")
	userjson := NewUser{
		Roles:    item.Roles,
		Username: item.Username,
	}
	reqjson, _ := json.Marshal(userjson)
	payload := strings.NewReader(string(reqjson))
	url := "https://cloud.mongodb.com/api/atlas/v1.0/orgs/" + c.orgid + "/invites?pretty=true"
	method := "POST"
	req, err := http.NewRequest(method, url, nil)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	digestParts := digestParts(resp)
	digestParts["uri"] = url
	digestParts["method"] = method
	digestParts["username"] = c.publickey
	digestParts["password"] = c.privateKey
	req, err = http.NewRequest(method, url, payload)
	req.Header.Set("Authorization", getDigestAuthrization(digestParts))
	req.Header.Set("Content-Type", "application/json")

	resp, err = client.Do(req)
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}
	statuscode := (int)(resp.StatusCode)
	if statuscode >= 200 && statuscode <= 400 {
		newbody := resp.Body
		newItemUser := &NewReturnUser{}
		err = json.NewDecoder(newbody).Decode(newItemUser)
		if err != nil {
			return nil, fmt.Errorf(err.Error())
		}
		return newItemUser, nil
	}
	return nil, fmt.Errorf("Error : %v ", Errors[statuscode])
}

func (c *Client) UpdateUser(updatevalue *UpdateUser, userId string) (*User, error) {
	fmt.Println("Update user")
	reqjson, _ := json.Marshal(updatevalue)
	payload := strings.NewReader(string(reqjson))
	url := "https://cloud.mongodb.com/api/public/v1.0/users/" + userId
	method := "PATCH"
	req, err := http.NewRequest(method, url, nil)
	req.Header.Add("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error client do ,", err)
	}
	defer resp.Body.Close()
	digestParts := digestParts(resp)
	digestParts["uri"] = url
	digestParts["method"] = method
	digestParts["username"] = c.publickey
	digestParts["password"] = c.privateKey
	req, err = http.NewRequest(method, url, payload)
	req.Header.Set("Authorization", getDigestAuthrization(digestParts))
	req.Header.Set("Content-Type", "application/json")
	resp, err = client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		fmt.Println("eror :", err)
	}
	statuscode := (int)(resp.StatusCode)
	if statuscode >= 200 && statuscode < 400 {
		newbody := resp.Body
		newItemUser := &User{}
		err = json.NewDecoder(newbody).Decode(newItemUser)
		if err != nil {
			return nil, fmt.Errorf(err.Error())
		}
		return newItemUser, nil
	}
	return nil, fmt.Errorf("Error : %v ", Errors[statuscode])
}

func digestParts(resp *http.Response) map[string]string {
	result := map[string]string{}
	if len(resp.Header["Www-Authenticate"]) > 0 {
		wantedHeaders := []string{"nonce", "realm", "qop"}
		responseHeaders := strings.Split(resp.Header["Www-Authenticate"][0], ",")
		for _, r := range responseHeaders {
			for _, w := range wantedHeaders {
				if strings.Contains(r, w) {
					result[w] = strings.Split(r, `"`)[1]
				}
			}
		}
	}
	return result
}

func getMD5(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func getCnonce() string {
	b := make([]byte, 8)
	io.ReadFull(rand.Reader, b)
	return fmt.Sprintf("%x", b)[:16]
}

func getDigestAuthrization(digestParts map[string]string) string {
	d := digestParts
	ha1 := getMD5(d["username"] + ":" + d["realm"] + ":" + d["password"])
	ha2 := getMD5(d["method"] + ":" + d["uri"])
	nonceCount := 00000001
	cnonce := getCnonce()
	response := getMD5(fmt.Sprintf("%s:%s:%v:%s:%s:%s", ha1, d["nonce"], nonceCount, cnonce, d["qop"], ha2))
	authorization := fmt.Sprintf(`Digest username="%s", realm="%s", nonce="%s", uri="%s", cnonce="%s", nc="%v", qop="%s", response="%s"`,
		d["username"], d["realm"], d["nonce"], d["uri"], cnonce, nonceCount, d["qop"], response)
	return authorization
}

func (c *Client) IsRetry(err error) bool {
	if err != nil {
		if strings.Contains(err.Error(), "\"responseCode\":503") == true {
			return true
		}
	}
	return false
}
