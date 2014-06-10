// A library for accessing Amazon S3 from Go
package s34go

import "crypto/hmac"
import "crypto/sha256"
import "encoding/hex"
import "io"
import "log"
import "net/http"
import "net/url"
import "sort"
import "strings"
import "time"


// Holds all information for a single request. Constructed using the "builder" pattern, with
// the various Set... and Add... methods.
type S3Request struct {
	method        string
	accessKey     string
	secretKey     string
	region        string
	hostName      string
	bucketName    string
	objectName    string
	queryParams   map[string][]string
	headers       map[string]string
	cookedHeaders []CookedHeader // filled by prepare()
	timestamp     time.Time      // filled by factory
}

//----------------------------------------------------------------------------------------------
// Methods to build an instance
//----------------------------------------------------------------------------------------------

func NewS3Request(accessKey string, secretKey string) *S3Request {
	queryParams := make(map[string][]string)
	headers := make(map[string]string)

	return &S3Request{HTTP_GET, accessKey, secretKey, DEFAULT_S3_REGION, DEFAULT_S3_ENDPOINT, "", "", queryParams, headers, nil, time.Now()}
}

func (request *S3Request) SetMethod(value string) *S3Request {
	// FIXME - validate, default to GET
	request.method = strings.ToUpper(strings.TrimSpace(value))
	return request
}

func (request *S3Request) SetHost(region string, hostName string) *S3Request {
	request.region = region
	request.hostName = hostName
	return request
}

func (request *S3Request) SetBucketName(value string) *S3Request {
	request.bucketName = value
	return request
}

func (request *S3Request) SetObjectName(value string) *S3Request {
	request.objectName = value
	return request
}

func (request *S3Request) AddHeader(name string, value string) *S3Request {
	request.headers[name] = value
	return request
}

func (request *S3Request) AddQueryParam(name string, value string) *S3Request {
	slice, exists := request.queryParams[name]
	if !exists {
		slice = make([]string, 0, 2)
	}
	request.queryParams[name] = append(slice, value)
	return request
}

func (request *S3Request) prepare() *S3Request {
	request.headers[HTTP_HDR_HOST] = request.hostName
	request.headers[HTTP_HDR_DATE] = request.timestamp.Format(time.RFC1123)
    request.headers[AMZN_HDR_SHA256] = request.payloadHash()
	request.cookedHeaders = request.cookHeaders()
    return request
}

//----------------------------------------------------------------------------------------------
// Other public methods
//----------------------------------------------------------------------------------------------

// Executes the request using the give client object, returning the raw response.
func (request *S3Request) execute(client *http.Client) (*S3Response, error) {
    method := request.method
    url := request.constructUrl()

    httpRequest,err := http.NewRequest(method, url, nil)
    if err != nil {
        return nil,err
    }

    for key,value := range request.headers {
        httpRequest.Header.Add(key, value)
    }
    httpRequest.Header.Add(HTTP_HDR_AUTH, request.authHeader())

    log.Println("executing", httpRequest)
    httpResponse,err := client.Do(httpRequest)
    if err != nil {
        return nil, err
    }

    return NewS3Response(httpResponse)
}

//----------------------------------------------------------------------------------------------
// Internals -- request-specific
//----------------------------------------------------------------------------------------------

type CookedHeader struct {
	rawName     string
	rawValue    string
	cookedName  string
	cookedValue string
}

type CookedHeaderSlice []CookedHeader

func (ch CookedHeaderSlice) Len() int           { return len(ch) }
func (ch CookedHeaderSlice) Less(i, j int) bool { return ch[i].cookedName < ch[j].cookedName }
func (ch CookedHeaderSlice) Swap(i, j int)      { ch[i], ch[j] = ch[j], ch[i] }

func (request *S3Request) cookHeaders() []CookedHeader {
	headers := make([]CookedHeader, 0, len(request.headers))
	for name, value := range request.headers {
		cookedName := strings.ToLower(strings.TrimSpace(name))
		cookedValue := strings.TrimSpace(value)
		headers = append(headers, CookedHeader{name, value, cookedName, cookedValue})
	}
	sort.Sort(CookedHeaderSlice(headers))
	return headers
}

func (request *S3Request) canonicalPath() string {
	url := "/"
	if len(request.bucketName) > 0 {
		url = url + request.bucketName + "/"
	}
	if len(request.objectName) > 0 {
		url = url + request.objectName
	}
	return url
}

func (request *S3Request) canonicalQueryString() string {
	cookedParams := make(map[string][]string)
	cookedParamNames := make([]string, 0, len(request.queryParams))

	for param, values := range request.queryParams {
		cookedValues := make([]string, 0, len(values))
		for _, value := range values {
			cookedValues = append(cookedValues, uriEncode(value))
		}
		cookedParam := uriEncode(param)
		cookedParams[cookedParam] = cookedValues
		cookedParamNames = append(cookedParamNames, cookedParam)

	}

	sort.Sort(sort.StringSlice(cookedParamNames))

	result := ""
	for _, name := range cookedParamNames {
		for _, value := range cookedParams[name] {
			result = result + name + "=" + value + "&"
		}
	}
	return strings.TrimRight(result, "&")
}

func (request *S3Request) canonicalHeaders() string {
	result := ""
	for _, header := range request.cookedHeaders {
		result = result + header.cookedName + ":" + header.cookedValue + "\n"
	}
	return result
}

func (request *S3Request) signedHeaders() string {
	result := ""
	for _, header := range request.cookedHeaders {
		result = result + header.cookedName + ";"
	}
	return strings.TrimRight(result, ";")
}

func (request *S3Request) payloadHash() string {
	// FIXME - once we add payload to request
	return computeHash(strings.NewReader(""))
}

func (request *S3Request) scope() string {
	return request.timestamp.Format("20060102") + "/" + request.region + "/s3/aws4_request"
}

func (request *S3Request) canonicalRequestString() string {
	return request.method + "\n" +
		request.canonicalPath() + "\n" +
		request.canonicalQueryString() + "\n" +
		request.canonicalHeaders() + "\n" +
		request.signedHeaders() + "\n" +
		request.payloadHash()
}

func (request *S3Request) stringToSign() string {
log.Println("canonical request string:\n" + request.canonicalRequestString())
	return "AWS4-HMAC-SHA256" + "\n" +
		request.timestamp.UTC().Format("20060102T150405Z") + "\n" +
		request.scope() + "\n" +
		computeHash(strings.NewReader(request.canonicalRequestString()))
}

func (request *S3Request) signingKey() []byte {
	dateKey := sign([]byte("AWS4"+request.secretKey), []byte(request.timestamp.Format("20060102")))
	regionKey := sign(dateKey, []byte(request.region))
	serviceKey := sign(regionKey, []byte("s3"))
	return sign(serviceKey, []byte("aws4_request"))
}

func (request *S3Request) authHeader() string {
log.Println("string to sign:\n" + request.stringToSign())
	signedString := sign(request.signingKey(), []byte(request.stringToSign()))
	return "AWS4-HMAC-SHA256 " +
		"Credential=" + request.accessKey + "/" + request.scope() + ", " +
		"SignedHeaders=" + request.signedHeaders() + ", " +
		"Signature=" + hex.EncodeToString(signedString)
}

func (request *S3Request) constructUrl() string {
    url := "http://" + request.hostName + "/"
    if len(request.bucketName) > 0 {
        url = url + request.bucketName + "/"
        if len(request.objectName) > 0 {
            url = url + request.objectName
        }
    }

    return url
}

//----------------------------------------------------------------------------------------------
// Internals -- utility
//----------------------------------------------------------------------------------------------

// Amazon doesn't like the default Go encoding
func uriEncode(value string) string {
	initial := url.QueryEscape(value)
	return strings.Replace(initial, "+", "%20", -1)
}

// computes hexified hash over bytes from a Reader
// ... this is one place where I really miss exceptions
func computeHash(rdr io.Reader) string {
	sha256 := sha256.New()
	io.Copy(sha256, rdr)
	hash := sha256.Sum(make([]byte, 0))
	return hex.EncodeToString(hash)
}

// sign some content using a specific key
func sign(key []byte, content []byte) []byte {
	hmac := hmac.New(sha256.New, key)
	hmac.Write(content)
	return hmac.Sum(make([]byte, 0))
}

