// A library for accessing Amazon S3 from Go
package s34go

import "encoding/xml"
import "io"
import "io/ioutil"
import "log"
import "net/http"

// Holds the HTTP response and any transformations on it
type S3Response struct {
    httpResponse    *http.Response
    hasBeenClosed   bool
    *AWSError
}

func NewS3Response(httpResponse *http.Response) (*S3Response,error) {

    log.Println("status = ", httpResponse.StatusCode)
    log.Println("headers = ", httpResponse.Header)

    response := &S3Response{httpResponse, false, nil}

body,_ := ioutil.ReadAll(response.httpResponse.Body)
log.Println("response body:\n", string(body))

//    if httpResponse.StatusCode != 200 {
//        awsError := AWSError{}
//        err := response.UnmarshalBody(&awsError)
//        if err == nil {
//            err = S3Error{"Received error from AWS: " + awsError.Message}
//            response.AWSError = &awsError
//        }
//        return response, err
//    }

    return response, nil
}


//----------------------------------------------------------------------------------------------
// Public Methods
//----------------------------------------------------------------------------------------------

// Writes the HTTP response body to the passed Writer. Will close the response.
func (response *S3Response) WriteBody(dst io.Writer) error {
    defer response.Close()
    _, err := io.Copy(dst, response.httpResponse.Body)
    return err
}

// Unmarshals the HTTP response body into the given object. Will close the response.
func (response *S3Response) UnmarshalBody(obj interface{}) error {
    defer response.Close()
    unmarshaller := xml.NewDecoder(response.httpResponse.Body)
    return unmarshaller.Decode(obj)
}

// Close the HTTP response body, if it has not already been closed.
func (response *S3Response) Close() {
    if ! response.hasBeenClosed {
        response.httpResponse.Body.Close()
        response.hasBeenClosed = true
    }
}



//----------------------------------------------------------------------------------------------
// Internals
//----------------------------------------------------------------------------------------------
