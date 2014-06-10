// A library for accessing Amazon S3 from Go
package s34go

//------------------------------------------------------------------------------------------------
//  HTTP Constants
//------------------------------------------------------------------------------------------------

const HTTP_GET = "GET"
const HTTP_PUT = "PUT"

const HTTP_HDR_AUTH = "Authorization"
const HTTP_HDR_DATE = "Date"
const HTTP_HDR_HOST = "Host"

//------------------------------------------------------------------------------------------------
//  S3-specific Constants
//------------------------------------------------------------------------------------------------

const DEFAULT_S3_ENDPOINT = "s3.amazonaws.com"
const DEFAULT_S3_REGION   = "us-east-1"

const S3_NAMESPACE =  "http://s3.amazonaws.com/doc/2006-03-01/";

const AMZN_HDR_SHA256 = "x-amz-content-sha256"
