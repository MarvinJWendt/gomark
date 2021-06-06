# http

Package http provides HTTP client and server implementations.

Get, Head, Post, and PostForm make HTTP (or HTTPS) requests:

    resp, err := http.Get("http://example.com/")
    ...
    resp, err := http.Post("http://example.com/upload", "image/jpeg", &buf)
    ...
    resp, err := http.PostForm("http://example.com/form",
    	url.Values{"key": {"Value"}, "id": {"123"}})

The client must close the response body when finished with it:

    resp, err := http.Get("http://example.com/")
    if err != nil {
    	// handle error
    }
    defer resp.Body.Close()
    body, err := io.ReadAll(resp.Body)
    // ...

For control over HTTP client headers, redirect policy, and other settings,
create a Client:

    client := &http.Client{
    	CheckRedirect: redirectPolicyFunc,
    }

    resp, err := client.Get("http://example.com")
    // ...

    req, err := http.NewRequest("GET", "http://example.com", nil)
    // ...
    req.Header.Add("If-None-Match", `W/"wyzzy"`)
    resp, err := client.Do(req)
    // ...

For control over proxies, TLS configuration, keep-alives, compression, and
other settings, create a Transport:

    tr := &http.Transport{
    	MaxIdleConns:       10,
    	IdleConnTimeout:    30 * time.Second,
    	DisableCompression: true,
    }
    client := &http.Client{Transport: tr}
    resp, err := client.Get("https://example.com")

Clients and Transports are safe for concurrent use by multiple goroutines
and for efficiency should only be created once and re-used.

ListenAndServe starts an HTTP server with a given address and handler. The
handler is usually nil, which means to use DefaultServeMux. Handle and
HandleFunc add handlers to DefaultServeMux:

    http.Handle("/foo", fooHandler)

    http.HandleFunc("/bar", func(w http.ResponseWriter, r *http.Request) {
    	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
    })

    log.Fatal(http.ListenAndServe(":8080", nil))

More control over the server's behavior is available by creating a custom
Server:

    s := &http.Server{
    	Addr:           ":8080",
    	Handler:        myHandler,
    	ReadTimeout:    10 * time.Second,
    	WriteTimeout:   10 * time.Second,
    	MaxHeaderBytes: 1 << 20,
    }
    log.Fatal(s.ListenAndServe())

Starting with Go 1.6, the http package has transparent support for the
HTTP/2 protocol when using HTTPS. Programs that must disable HTTP/2 can do
so by setting Transport.TLSNextProto (for clients) or Server.TLSNextProto
(for servers) to a non-nil, empty map. Alternatively, the following GODEBUG
environment variables are currently supported:

    GODEBUG=http2client=0  # disable HTTP/2 client support
    GODEBUG=http2server=0  # disable HTTP/2 server support
    GODEBUG=http2debug=1   # enable verbose HTTP/2 debug logs
    GODEBUG=http2debug=2   # ... even more verbose, with frame dumps

The GODEBUG variables are not covered by Go's API compatibility promise.
Please report any issues before disabling HTTP/2 support:
https://golang.org/s/http2bug

The http package's Transport and Server both automatically enable HTTP/2
support for simple configurations. To enable HTTP/2 for more complex
configurations, to use lower-level HTTP/2 features, or to use a newer
version of Go's http2 package, import "golang.org/x/net/http2" directly and
use its ConfigureTransport and/or ConfigureServer functions. Manually
configuring HTTP/2 via the golang.org/x/net/http2 package takes precedence
over the net/http package's built-in HTTP/2 support.

## Constants

### DefaultMaxHeaderBytes

```go
const DefaultMaxHeaderBytes = 1 << 20 // 1 MB
```

DefaultMaxHeaderBytes is the maximum permitted size of the headers in an
HTTP request. This can be overridden by setting Server.MaxHeaderBytes.


### DefaultMaxIdleConnsPerHost

```go
const DefaultMaxIdleConnsPerHost = 2
```

DefaultMaxIdleConnsPerHost is the default value of Transport's
MaxIdleConnsPerHost.


### TimeFormat

```go
const TimeFormat = "Mon, 02 Jan 2006 15:04:05 GMT"
```

TimeFormat is the time format to use when generating times in HTTP headers.
It is like time.RFC1123 but hard-codes GMT as the time zone. The time being
formatted must be in UTC for Format to generate the correct format.
For parsing this time format, see ParseTime.


### TrailerPrefix

```go
const TrailerPrefix = "Trailer:"
```

TrailerPrefix is a magic prefix for ResponseWriter.Header map keys that, if
present, signals that the map entry is actually for the response trailers,
and not the response headers. The prefix is stripped after the ServeHTTP
call finishes and the values are sent in the trailers.
This mechanism is intended only for trailers that are not known prior to the
headers being written. If the set of trailers is fixed or known before the
header is written, the normal Go trailers mechanism is preferred:
https://golang.org/pkg/net/http/#ResponseWriter
https://golang.org/pkg/net/http/#example_ResponseWriter_trailers



## Constant Blocks

Common HTTP methods.
Unless otherwise noted, these are defined in RFC 7231 section 4.3.


```go
MethodGet     = "GET"
MethodHead    = "HEAD"
MethodPost    = "POST"
MethodPut     = "PUT"
MethodPatch   = "PATCH" // RFC 5789
MethodDelete  = "DELETE"
MethodConnect = "CONNECT"
MethodOptions = "OPTIONS"
MethodTrace   = "TRACE"
```

HTTP status codes as registered with IANA. See:
https://www.iana.org/assignments/http-status-codes/http-status-codes.xhtml


```go
StatusContinue           = 100 // RFC 7231, 6.2.1
StatusSwitchingProtocols = 101 // RFC 7231, 6.2.2
StatusProcessing         = 102 // RFC 2518, 10.1
StatusEarlyHints         = 103 // RFC 8297

StatusOK                   = 200 // RFC 7231, 6.3.1
StatusCreated              = 201 // RFC 7231, 6.3.2
StatusAccepted             = 202 // RFC 7231, 6.3.3
StatusNonAuthoritativeInfo = 203 // RFC 7231, 6.3.4
StatusNoContent            = 204 // RFC 7231, 6.3.5
StatusResetContent         = 205 // RFC 7231, 6.3.6
StatusPartialContent       = 206 // RFC 7233, 4.1
StatusMultiStatus          = 207 // RFC 4918, 11.1
StatusAlreadyReported      = 208 // RFC 5842, 7.1
StatusIMUsed               = 226 // RFC 3229, 10.4.1

StatusMultipleChoices  = 300 // RFC 7231, 6.4.1
StatusMovedPermanently = 301 // RFC 7231, 6.4.2
StatusFound            = 302 // RFC 7231, 6.4.3
StatusSeeOther         = 303 // RFC 7231, 6.4.4
StatusNotModified      = 304 // RFC 7232, 4.1
StatusUseProxy         = 305 // RFC 7231, 6.4.5

StatusTemporaryRedirect = 307 // RFC 7231, 6.4.7
StatusPermanentRedirect = 308 // RFC 7538, 3

StatusBadRequest                   = 400 // RFC 7231, 6.5.1
StatusUnauthorized                 = 401 // RFC 7235, 3.1
StatusPaymentRequired              = 402 // RFC 7231, 6.5.2
StatusForbidden                    = 403 // RFC 7231, 6.5.3
StatusNotFound                     = 404 // RFC 7231, 6.5.4
StatusMethodNotAllowed             = 405 // RFC 7231, 6.5.5
StatusNotAcceptable                = 406 // RFC 7231, 6.5.6
StatusProxyAuthRequired            = 407 // RFC 7235, 3.2
StatusRequestTimeout               = 408 // RFC 7231, 6.5.7
StatusConflict                     = 409 // RFC 7231, 6.5.8
StatusGone                         = 410 // RFC 7231, 6.5.9
StatusLengthRequired               = 411 // RFC 7231, 6.5.10
StatusPreconditionFailed           = 412 // RFC 7232, 4.2
StatusRequestEntityTooLarge        = 413 // RFC 7231, 6.5.11
StatusRequestURITooLong            = 414 // RFC 7231, 6.5.12
StatusUnsupportedMediaType         = 415 // RFC 7231, 6.5.13
StatusRequestedRangeNotSatisfiable = 416 // RFC 7233, 4.4
StatusExpectationFailed            = 417 // RFC 7231, 6.5.14
StatusTeapot                       = 418 // RFC 7168, 2.3.3
StatusMisdirectedRequest           = 421 // RFC 7540, 9.1.2
StatusUnprocessableEntity          = 422 // RFC 4918, 11.2
StatusLocked                       = 423 // RFC 4918, 11.3
StatusFailedDependency             = 424 // RFC 4918, 11.4
StatusTooEarly                     = 425 // RFC 8470, 5.2.
StatusUpgradeRequired              = 426 // RFC 7231, 6.5.15
StatusPreconditionRequired         = 428 // RFC 6585, 3
StatusTooManyRequests              = 429 // RFC 6585, 4
StatusRequestHeaderFieldsTooLarge  = 431 // RFC 6585, 5
StatusUnavailableForLegalReasons   = 451 // RFC 7725, 3

StatusInternalServerError           = 500 // RFC 7231, 6.6.1
StatusNotImplemented                = 501 // RFC 7231, 6.6.2
StatusBadGateway                    = 502 // RFC 7231, 6.6.3
StatusServiceUnavailable            = 503 // RFC 7231, 6.6.4
StatusGatewayTimeout                = 504 // RFC 7231, 6.6.5
StatusHTTPVersionNotSupported       = 505 // RFC 7231, 6.6.6
StatusVariantAlsoNegotiates         = 506 // RFC 2295, 8.1
StatusInsufficientStorage           = 507 // RFC 4918, 11.5
StatusLoopDetected                  = 508 // RFC 5842, 7.2
StatusNotExtended                   = 510 // RFC 2774, 7
StatusNetworkAuthenticationRequired = 511 // RFC 6585, 6
```



## Variables

### DefaultClient

```go
var DefaultClient = &Client{}
```

DefaultClient is the default Client and is used by Get, Head, and Post.


### DefaultServeMux

```go
var DefaultServeMux = &defaultServeMux
```

DefaultServeMux is the default ServeMux used by Serve.


### ErrAbortHandler

```go
var ErrAbortHandler = errors.New("net/http: abort Handler")
```

ErrAbortHandler is a sentinel panic value to abort a handler. While any
panic from ServeHTTP aborts the response to the client, panicking with
ErrAbortHandler also suppresses logging of a stack trace to the server's
error log.


### ErrBodyReadAfterClose

```go
var ErrBodyReadAfterClose = errors.New("http: invalid Read on closed Body")
```

ErrBodyReadAfterClose is returned when reading a Request or Response Body
after the body has been closed. This typically happens when the body is read
after an HTTP Handler calls WriteHeader or Write on its ResponseWriter.


### ErrHandlerTimeout

```go
var ErrHandlerTimeout = errors.New("http: Handler timeout")
```

ErrHandlerTimeout is returned on ResponseWriter Write calls in handlers
which have timed out.


### ErrLineTooLong

```go
var ErrLineTooLong = internal.ErrLineTooLong
```

ErrLineTooLong is returned when reading request or response bodies with
malformed chunked encoding.


### ErrMissingFile

```go
var ErrMissingFile = errors.New("http: no such file")
```

ErrMissingFile is returned by FormFile when the provided file field name is
either not present in the request or not a file field.


### ErrNoCookie

```go
var ErrNoCookie = errors.New("http: named cookie not present")
```

ErrNoCookie is returned by Request's Cookie method when a cookie is not
found.


### ErrNoLocation

```go
var ErrNoLocation = errors.New("http: no Location header in response")
```

ErrNoLocation is returned by Response's Location method when no Location
header is present.


### ErrServerClosed

```go
var ErrServerClosed = errors.New("http: Server closed")
```

ErrServerClosed is returned by the Server's Serve, ServeTLS, ListenAndServe,
and ListenAndServeTLS methods after a call to Shutdown or Close.


### ErrSkipAltProtocol

```go
var ErrSkipAltProtocol = errors.New("net/http: skip alternate protocol")
```

ErrSkipAltProtocol is a sentinel error value defined by
Transport.RegisterProtocol.


### ErrUseLastResponse

```go
var ErrUseLastResponse = errors.New("net/http: use last response")
```

ErrUseLastResponse can be returned by Client.CheckRedirect hooks to control
how redirects are processed. If returned, the next request is not sent and
the most recent response is returned with its body unclosed.


### NoBody

```go
var NoBody = noBody{}
```

NoBody is an io.ReadCloser with no bytes. Read always returns EOF and Close
always returns nil. It can be used in an outgoing client request to
explicitly signal that a request has zero bytes. An alternative, however, is
to simply set Request.Body to nil.



## Variable Blocks


```go
// ErrNotSupported is returned by the Push method of Pusher
// implementations to indicate that HTTP/2 Push support is not
// available.
ErrNotSupported = &ProtocolError{"feature not supported"}

// Deprecated: ErrUnexpectedTrailer is no longer returned by
// anything in the net/http package. Callers should not
// compare errors against this variable.
ErrUnexpectedTrailer = &ProtocolError{"trailer header without chunked transfer encoding"}

// ErrMissingBoundary is returned by Request.MultipartReader when the
// request's Content-Type does not include a "boundary" parameter.
ErrMissingBoundary = &ProtocolError{"no multipart boundary param in Content-Type"}

// ErrNotMultipart is returned by Request.MultipartReader when the
// request's Content-Type is not multipart/form-data.
ErrNotMultipart = &ProtocolError{"request Content-Type isn't multipart/form-data"}

// Deprecated: ErrHeaderTooLong is no longer returned by
// anything in the net/http package. Callers should not
// compare errors against this variable.
ErrHeaderTooLong = &ProtocolError{"header too long"}

// Deprecated: ErrShortBody is no longer returned by
// anything in the net/http package. Callers should not
// compare errors against this variable.
ErrShortBody = &ProtocolError{"entity body too short"}

// Deprecated: ErrMissingContentLength is no longer returned by
// anything in the net/http package. Callers should not
// compare errors against this variable.
ErrMissingContentLength = &ProtocolError{"missing ContentLength in HEAD response"}
```

Errors used by the HTTP server.


```go
// ErrBodyNotAllowed is returned by ResponseWriter.Write calls
// when the HTTP method or response code does not permit a
// body.
ErrBodyNotAllowed = errors.New("http: request method or response status code does not allow body")

// ErrHijacked is returned by ResponseWriter.Write calls when
// the underlying connection has been hijacked using the
// Hijacker interface. A zero-byte write on a hijacked
// connection will return ErrHijacked without any other side
// effects.
ErrHijacked = errors.New("http: connection has been hijacked")

// ErrContentLength is returned by ResponseWriter.Write calls
// when a Handler set a Content-Length response header with a
// declared size and then attempted to write more bytes than
// declared.
ErrContentLength = errors.New("http: wrote more than the declared Content-Length")

// Deprecated: ErrWriteAfterFlush is no longer returned by
// anything in the net/http package. Callers should not
// compare errors against this variable.
ErrWriteAfterFlush = errors.New("unused")
```


```go
// ServerContextKey is a context key. It can be used in HTTP
// handlers with Context.Value to access the server that
// started the handler. The associated value will be of
// type *Server.
ServerContextKey = &contextKey{"http-server"}

// LocalAddrContextKey is a context key. It can be used in
// HTTP handlers with Context.Value to access the local
// address the connection arrived on.
// The associated value will be of type net.Addr.
LocalAddrContextKey = &contextKey{"local-addr"}
```



## Functions

### CanonicalHeaderKey

```go
func CanonicalHeaderKey(s string) string
```

CanonicalHeaderKey returns the canonical format of the header key s. The
canonicalization converts the first letter and any letter following a hyphen
to upper case; the rest are converted to lowercase. For example, the
canonical key for "accept-encoding" is "Accept-Encoding". If s contains a
space or invalid header field bytes, it is returned without modifications.


### DetectContentType

```go
func DetectContentType(data []byte) string
```

DetectContentType implements the algorithm described at
https://mimesniff.spec.whatwg.org/ to determine the Content-Type of the
given data. It considers at most the first 512 bytes of data.
DetectContentType always returns a valid MIME type: if it cannot determine a
more specific one, it returns "application/octet-stream".


### Error

```go
func Error(w ResponseWriter, error string, code int)
```

Error replies to the request with the specified error message and HTTP code.
It does not otherwise end the request; the caller should ensure no further
writes are done to w. The error message should be plain text.


### Get

```go
func Get(url string) (resp *Response, err error)
```

Get issues a GET to the specified URL. If the response is one of the
following redirect codes, Get follows the redirect, up to a maximum of 10
redirects:
301 (Moved Permanently)
302 (Found)
303 (See Other)
307 (Temporary Redirect)
308 (Permanent Redirect)
An error is returned if there were too many redirects or if there was an
HTTP protocol error. A non-2xx response doesn't cause an error. Any returned
error will be of type *url.Error. The url.Error value's Timeout method will
report true if request timed out or was canceled.
When err is nil, resp always contains a non-nil resp.Body. Caller should
close resp.Body when done reading from it.
Get is a wrapper around DefaultClient.Get.
To make a request with custom headers, use NewRequest and DefaultClient.Do.


### Handle

```go
func Handle(pattern string, handler Handler)
```

Handle registers the handler for the given pattern in the DefaultServeMux.
The documentation for ServeMux explains how patterns are matched.


### HandleFunc

```go
func HandleFunc(pattern string, handler func(ResponseWriter, *Request))
```

HandleFunc registers the handler function for the given pattern in the
DefaultServeMux. The documentation for ServeMux explains how patterns are
matched.


### Head

```go
func Head(url string) (resp *Response, err error)
```

Head issues a HEAD to the specified URL. If the response is one of the
following redirect codes, Head follows the redirect, up to a maximum of 10
redirects:
301 (Moved Permanently)
302 (Found)
303 (See Other)
307 (Temporary Redirect)
308 (Permanent Redirect)
Head is a wrapper around DefaultClient.Head


### ListenAndServe

```go
func ListenAndServe(addr string, handler Handler) error
```

ListenAndServe listens on the TCP network address addr and then calls Serve
with handler to handle requests on incoming connections. Accepted
connections are configured to enable TCP keep-alives.
The handler is typically nil, in which case the DefaultServeMux is used.
ListenAndServe always returns a non-nil error.


### ListenAndServeTLS

```go
func ListenAndServeTLS(addr, certFile, keyFile string, handler Handler) error
```

ListenAndServeTLS acts identically to ListenAndServe, except that it expects
HTTPS connections. Additionally, files containing a certificate and matching
private key for the server must be provided. If the certificate is signed by
a certificate authority, the certFile should be the concatenation of the
server's certificate, any intermediates, and the CA's certificate.


### MaxBytesReader

```go
func MaxBytesReader(w ResponseWriter, r io.ReadCloser, n int64) io.ReadCloser
```

MaxBytesReader is similar to io.LimitReader but is intended for limiting the
size of incoming request bodies. In contrast to io.LimitReader,
MaxBytesReader's result is a ReadCloser, returns a non-EOF error for a Read
beyond the limit, and closes the underlying reader when its Close method is
called.
MaxBytesReader prevents clients from accidentally or maliciously sending a
large request and wasting server resources.


### NewRequest

```go
func NewRequest(method, url string, body io.Reader) (*Request, error)
```

NewRequest wraps NewRequestWithContext using the background context.


### NewRequestWithContext

```go
func NewRequestWithContext(ctx context.Context, method, url string, body io.Reader) (*Request, error)
```

NewRequestWithContext returns a new Request given a method, URL, and
optional body.
If the provided body is also an io.Closer, the returned Request.Body is set
to body and will be closed by the Client methods Do, Post, and PostForm, and
Transport.RoundTrip.
NewRequestWithContext returns a Request suitable for use with Client.Do or
Transport.RoundTrip. To create a request for use with testing a Server
Handler, either use the NewRequest function in the net/http/httptest
package, use ReadRequest, or manually update the Request fields. For an
outgoing client request, the context controls the entire lifetime of a
request and its response: obtaining a connection, sending the request, and
reading the response headers and body. See the Request type's documentation
for the difference between inbound and outbound request fields.
If body is of type *bytes.Buffer, *bytes.Reader, or *strings.Reader, the
returned request's ContentLength is set to its exact value (instead of -1),
GetBody is populated (so 307 and 308 redirects can replay the body), and
Body is set to NoBody if the ContentLength is 0.


### NotFound

```go
func NotFound(w ResponseWriter, r *Request)
```

NotFound replies to the request with an HTTP 404 not found error.


### ParseHTTPVersion

```go
func ParseHTTPVersion(vers string) (major, minor int, ok bool)
```

ParseHTTPVersion parses an HTTP version string. "HTTP/1.0" returns (1, 0,
true).


### ParseTime

```go
func ParseTime(text string) (t time.Time, err error)
```

ParseTime parses a time header (such as the Date: header), trying each of
the three formats allowed by HTTP/1.1: TimeFormat, time.RFC850, and
time.ANSIC.


### Post

```go
func Post(url, contentType string, body io.Reader) (resp *Response, err error)
```

Post issues a POST to the specified URL.
Caller should close resp.Body when done reading from it.
If the provided body is an io.Closer, it is closed after the request.
Post is a wrapper around DefaultClient.Post.
To set custom headers, use NewRequest and DefaultClient.Do.
See the Client.Do method documentation for details on how redirects are
handled.


### PostForm

```go
func PostForm(url string, data url.Values) (resp *Response, err error)
```

PostForm issues a POST to the specified URL, with data's keys and values
URL-encoded as the request body.
The Content-Type header is set to application/x-www-form-urlencoded. To set
other headers, use NewRequest and DefaultClient.Do.
When err is nil, resp always contains a non-nil resp.Body. Caller should
close resp.Body when done reading from it.
PostForm is a wrapper around DefaultClient.PostForm.
See the Client.Do method documentation for details on how redirects are
handled.


### ProxyFromEnvironment

```go
func ProxyFromEnvironment(req *Request) (*url.URL, error)
```

ProxyFromEnvironment returns the URL of the proxy to use for a given
request, as indicated by the environment variables HTTP_PROXY, HTTPS_PROXY
and NO_PROXY (or the lowercase versions thereof). HTTPS_PROXY takes
precedence over HTTP_PROXY for https requests.
The environment values may be either a complete URL or a "host[:port]", in
which case the "http" scheme is assumed. An error is returned if the value
is a different form.
A nil URL and nil error are returned if no proxy is defined in the
environment, or a proxy should not be used for the given request, as defined
by NO_PROXY.
As a special case, if req.URL.Host is "localhost" (with or without a port
number), then a nil URL and nil error will be returned.


### ProxyURL

```go
func ProxyURL(fixedURL *url.URL) func(*Request) (*url.URL, error)
```

ProxyURL returns a proxy function (for use in a Transport) that always
returns the same URL.


### ReadRequest

```go
func ReadRequest(b *bufio.Reader) (*Request, error)
```

ReadRequest reads and parses an incoming request from b.
ReadRequest is a low-level function and should only be used for specialized
applications; most code should use the Server to read requests and handle
them via the Handler interface. ReadRequest only supports HTTP/1.x requests.
For HTTP/2, use golang.org/x/net/http2.


### ReadResponse

```go
func ReadResponse(r *bufio.Reader, req *Request) (*Response, error)
```

ReadResponse reads and returns an HTTP response from r. The req parameter
optionally specifies the Request that corresponds to this Response. If nil,
a GET request is assumed. Clients must call resp.Body.Close when finished
reading resp.Body. After that call, clients can inspect resp.Trailer to find
key/value pairs included in the response trailer.


### Redirect

```go
func Redirect(w ResponseWriter, r *Request, url string, code int)
```

Redirect replies to the request with a redirect to url, which may be a path
relative to the request path.
The provided code should be in the 3xx range and is usually
StatusMovedPermanently, StatusFound or StatusSeeOther.
If the Content-Type header has not been set, Redirect sets it to "text/html;
charset=utf-8" and writes a small HTML body. Setting the Content-Type header
to any value, including nil, disables that behavior.


### Serve

```go
func Serve(l net.Listener, handler Handler) error
```

Serve accepts incoming HTTP connections on the listener l, creating a new
service goroutine for each. The service goroutines read requests and then
call handler to reply to them.
The handler is typically nil, in which case the DefaultServeMux is used.
HTTP/2 support is only enabled if the Listener returns *tls.Conn connections
and they were configured with "h2" in the TLS Config.NextProtos.
Serve always returns a non-nil error.


### ServeContent

```go
func ServeContent(w ResponseWriter, req *Request, name string, modtime time.Time, content io.ReadSeeker)
```

ServeContent replies to the request using the content in the provided
ReadSeeker. The main benefit of ServeContent over io.Copy is that it handles
Range requests properly, sets the MIME type, and handles If-Match,
If-Unmodified-Since, If-None-Match, If-Modified-Since, and If-Range
requests.
If the response's Content-Type header is not set, ServeContent first tries
to deduce the type from name's file extension and, if that fails, falls back
to reading the first block of the content and passing it to
DetectContentType. The name is otherwise unused; in particular it can be
empty and is never sent in the response.
If modtime is not the zero time or Unix epoch, ServeContent includes it in a
Last-Modified header in the response. If the request includes an
If-Modified-Since header, ServeContent uses modtime to decide whether the
content needs to be sent at all.
The content's Seek method must work: ServeContent uses a seek to the end of
the content to determine its size.
If the caller has set w's ETag header formatted per RFC 7232, section 2.3,
ServeContent uses it to handle requests using If-Match, If-None-Match, or
If-Range.
Note that *os.File implements the io.ReadSeeker interface.


### ServeFile

```go
func ServeFile(w ResponseWriter, r *Request, name string)
```

ServeFile replies to the request with the contents of the named file or
directory.
If the provided file or directory name is a relative path, it is interpreted
relative to the current directory and may ascend to parent directories. If
the provided name is constructed from user input, it should be sanitized
before calling ServeFile.
As a precaution, ServeFile will reject requests where r.URL.Path contains a
".." path element; this protects against callers who might unsafely use
filepath.Join on r.URL.Path without sanitizing it and then use that
filepath.Join result as the name argument.
As another special case, ServeFile redirects any request where r.URL.Path
ends in "/index.html" to the same path, without the final "index.html". To
avoid such redirects either modify the path or use ServeContent.
Outside of those two special cases, ServeFile does not use r.URL.Path for
selecting the file or directory to serve; only the file or directory
provided in the name argument is used.


### ServeTLS

```go
func ServeTLS(l net.Listener, handler Handler, certFile, keyFile string) error
```

ServeTLS accepts incoming HTTPS connections on the listener l, creating a
new service goroutine for each. The service goroutines read requests and
then call handler to reply to them.
The handler is typically nil, in which case the DefaultServeMux is used.
Additionally, files containing a certificate and matching private key for
the server must be provided. If the certificate is signed by a certificate
authority, the certFile should be the concatenation of the server's
certificate, any intermediates, and the CA's certificate.
ServeTLS always returns a non-nil error.


### SetCookie

```go
func SetCookie(w ResponseWriter, cookie *Cookie)
```

SetCookie adds a Set-Cookie header to the provided ResponseWriter's headers.
The provided cookie must have a valid Name. Invalid cookies may be silently
dropped.


### StatusText

```go
func StatusText(code int) string
```

StatusText returns a text for the HTTP status code. It returns the empty
string if the code is unknown.



## Types

### ConnState

```go
type ConnState int
```

A ConnState represents the state of a client connection to a server. It's
used by the optional Server.ConnState hook.



#### ConnState.String

```go
func (c ConnState) String() string
```



### Dir

```go
type Dir string
```

A Dir implements FileSystem using the native file system restricted to a
specific directory tree.
While the FileSystem.Open method takes '/'-separated paths, a Dir's string
value is a filename on the native file system, not a URL, so it is separated
by filepath.Separator, which isn't necessarily '/'.
Note that Dir could expose sensitive files and directories. Dir will follow
symlinks pointing out of the directory tree, which can be especially
dangerous if serving from a directory in which users are able to create
arbitrary symlinks. Dir will also allow access to files and directories
starting with a period, which could expose sensitive directories like .git
or sensitive files like .htpasswd. To exclude files with a leading period,
remove the files/directories from the server or create a custom FileSystem
implementation.
An empty Dir is treated as ".".



#### Dir.Open

```go
func (d Dir) Open(name string) (File, error)
```

Open implements FileSystem using os.Open, opening files for reading rooted
and relative to the directory d.


### HandlerFunc

```go
type HandlerFunc func(ResponseWriter, *Request)
```

The HandlerFunc type is an adapter to allow the use of ordinary functions as
HTTP handlers. If f is a function with the appropriate signature,
HandlerFunc(f) is a Handler that calls f.



#### HandlerFunc.ServeHTTP

```go
func (f HandlerFunc) ServeHTTP(w ResponseWriter, r *Request)
```

ServeHTTP calls f(w, r).


### Header

```go
type Header map[string][]string
```

A Header represents the key-value pairs in an HTTP header.
The keys should be in canonical form, as returned by CanonicalHeaderKey.



#### Header.Add

```go
func (h Header) Add(key, value string)
```

Add adds the key, value pair to the header. It appends to any existing
values associated with key. The key is case insensitive; it is canonicalized
by CanonicalHeaderKey.


#### Header.Clone

```go
func (h Header) Clone() Header
```

Clone returns a copy of h or nil if h is nil.


#### Header.Del

```go
func (h Header) Del(key string)
```

Del deletes the values associated with key. The key is case insensitive; it
is canonicalized by CanonicalHeaderKey.


#### Header.Get

```go
func (h Header) Get(key string) string
```

Get gets the first value associated with the given key. If there are no
values associated with the key, Get returns "". It is case insensitive;
textproto.CanonicalMIMEHeaderKey is used to canonicalize the provided key.
To use non-canonical keys, access the map directly.


#### Header.Set

```go
func (h Header) Set(key, value string)
```

Set sets the header entries associated with key to the single element value.
It replaces any existing values associated with key. The key is case
insensitive; it is canonicalized by textproto.CanonicalMIMEHeaderKey. To use
non-canonical keys, assign to the map directly.


#### Header.Values

```go
func (h Header) Values(key string) []string
```

Values returns all values associated with the given key. It is case
insensitive; textproto.CanonicalMIMEHeaderKey is used to canonicalize the
provided key. To use non-canonical keys, access the map directly. The
returned slice is not a copy.


#### Header.Write

```go
func (h Header) Write(w io.Writer) error
```

Write writes a header in wire format.


#### Header.WriteSubset

```go
func (h Header) WriteSubset(w io.Writer, exclude map[string]bool) error
```

WriteSubset writes a header in wire format. If exclude is not nil, keys
where exclude[key] == true are not written. Keys are not canonicalized
before checking the exclude map.



## Structs

### Client

```go
func (c *Client) PostForm(url string, data url.Values) (resp *Response, err error)
```

A Client is an HTTP client. Its zero value (DefaultClient) is a usable
client that uses DefaultTransport.
The Client's Transport typically has internal state (cached TCP
connections), so Clients should be reused instead of created as needed.
Clients are safe for concurrent use by multiple goroutines.
A Client is higher-level than a RoundTripper (such as Transport) and
additionally handles HTTP details such as cookies and redirects.
When following redirects, the Client will forward all headers set on the
initial Request except:
• when forwarding sensitive headers like "Authorization",
"WWW-Authenticate", and "Cookie" to untrusted targets. These headers will be
ignored when following a redirect to a domain that is not a subdomain match
or exact match of the initial domain. For example, a redirect from "foo.com"
to either "foo.com" or "sub.foo.com" will forward the sensitive headers, but
a redirect to "bar.com" will not.
• when forwarding the "Cookie" header with a non-nil cookie Jar. Since each
redirect may mutate the state of the cookie jar, a redirect may possibly
alter a cookie set in the initial request. When forwarding the "Cookie"
header, any mutated cookies will be omitted, with the expectation that the
Jar will insert those mutated cookies with the updated values (assuming the
origin matches). If Jar is nil, the initial cookies are forwarded without
change.



#### Client.CloseIdleConnections

```go
func (c *Client) CloseIdleConnections()
```

CloseIdleConnections closes any connections on its Transport which were
previously connected from previous requests but are now sitting idle in a
"keep-alive" state. It does not interrupt any connections currently in use.
If the Client's Transport does not have a CloseIdleConnections method then
this method does nothing.


#### Client.Do

```go
func (c *Client) Do(req *Request) (*Response, error)
```

Do sends an HTTP request and returns an HTTP response, following policy
(such as redirects, cookies, auth) as configured on the client.
An error is returned if caused by client policy (such as CheckRedirect), or
failure to speak HTTP (such as a network connectivity problem). A non-2xx
status code doesn't cause an error.
If the returned error is nil, the Response will contain a non-nil Body which
the user is expected to close. If the Body is not both read to EOF and
closed, the Client's underlying RoundTripper (typically Transport) may not
be able to re-use a persistent TCP connection to the server for a subsequent
"keep-alive" request.
The request Body, if non-nil, will be closed by the underlying Transport,
even on errors.
On error, any Response can be ignored. A non-nil Response with a non-nil
error only occurs when CheckRedirect fails, and even then the returned
Response.Body is already closed.
Generally Get, Post, or PostForm will be used instead of Do.
If the server replies with a redirect, the Client first uses the
CheckRedirect function to determine whether the redirect should be followed.
If permitted, a 301, 302, or 303 redirect causes subsequent requests to use
HTTP method GET (or HEAD if the original request was HEAD), with no body. A
307 or 308 redirect preserves the original HTTP method and body, provided
that the Request.GetBody function is defined. The NewRequest function
automatically sets GetBody for common standard library body types.
Any returned error will be of type *url.Error. The url.Error value's Timeout
method will report true if request timed out or was canceled.


#### Client.Get

```go
func (c *Client) Get(url string) (resp *Response, err error)
```

Get issues a GET to the specified URL. If the response is one of the
following redirect codes, Get follows the redirect after calling the
Client's CheckRedirect function:
301 (Moved Permanently)
302 (Found)
303 (See Other)
307 (Temporary Redirect)
308 (Permanent Redirect)
An error is returned if the Client's CheckRedirect function fails or if
there was an HTTP protocol error. A non-2xx response doesn't cause an error.
Any returned error will be of type *url.Error. The url.Error value's Timeout
method will report true if the request timed out.
When err is nil, resp always contains a non-nil resp.Body. Caller should
close resp.Body when done reading from it.
To make a request with custom headers, use NewRequest and Client.Do.


#### Client.Head

```go
func (c *Client) Head(url string) (resp *Response, err error)
```

Head issues a HEAD to the specified URL. If the response is one of the
following redirect codes, Head follows the redirect after calling the
Client's CheckRedirect function:
301 (Moved Permanently)
302 (Found)
303 (See Other)
307 (Temporary Redirect)
308 (Permanent Redirect)


#### Client.Post

```go
func (c *Client) Post(url, contentType string, body io.Reader) (resp *Response, err error)
```

Post issues a POST to the specified URL.
Caller should close resp.Body when done reading from it.
If the provided body is an io.Closer, it is closed after the request.
To set custom headers, use NewRequest and Client.Do.
See the Client.Do method documentation for details on how redirects are
handled.


#### Client.PostForm

```go
func (c *Client) PostForm(url string, data url.Values) (resp *Response, err error)
```

PostForm issues a POST to the specified URL, with data's keys and values
URL-encoded as the request body.
The Content-Type header is set to application/x-www-form-urlencoded. To set
other headers, use NewRequest and Client.Do.
When err is nil, resp always contains a non-nil resp.Body. Caller should
close resp.Body when done reading from it.
See the Client.Do method documentation for details on how redirects are
handled.


### Cookie

```go
func (c *Cookie) String() string
```

A Cookie represents an HTTP cookie as sent in the Set-Cookie header of an
HTTP response or the Cookie header of an HTTP request.
See https://tools.ietf.org/html/rfc6265 for details.



#### Cookie.String

```go
func (c *Cookie) String() string
```

String returns the serialization of the cookie for use in a Cookie header
(if only Name and Value are set) or a Set-Cookie response header (if other
fields are set). If c is nil or c.Name is invalid, the empty string is
returned.


### ProtocolError

```go
func (pe *ProtocolError) Error() string
```

ProtocolError represents an HTTP protocol error.
Deprecated: Not all errors in the http package related to protocol errors
are of type ProtocolError.



#### ProtocolError.Error

```go
func (pe *ProtocolError) Error() string
```



### PushOptions

```go
type PushOptions struct {
	// Method specifies the HTTP method for the promised request.
	// If set, it must be "GET" or "HEAD". Empty means "GET".
	Method string

	// Header specifies additional promised request headers. This cannot
	// include HTTP/2 pseudo header fields like ":path" and ":scheme",
	// which will be added automatically.
	Header Header
}
```

PushOptions describes options for Pusher.Push.



### Request

```go
func (r *Request) Write(w io.Writer) error
```

A Request represents an HTTP request received by a server or to be sent by a
client.
The field semantics differ slightly between client and server usage. In
addition to the notes on the fields below, see the documentation for
Request.Write and RoundTripper.



#### Request.AddCookie

```go
func (r *Request) AddCookie(c *Cookie)
```

AddCookie adds a cookie to the request. Per RFC 6265 section 5.4, AddCookie
does not attach more than one Cookie header field. That means all cookies,
if any, are written into the same line, separated by semicolon. AddCookie
only sanitizes c's name and value, and does not sanitize a Cookie header
already present in the request.


#### Request.BasicAuth

```go
func (r *Request) BasicAuth() (username, password string, ok bool)
```

BasicAuth returns the username and password provided in the request's
Authorization header, if the request uses HTTP Basic Authentication. See RFC
2617, Section 2.


#### Request.Clone

```go
func (r *Request) Clone(ctx context.Context) *Request
```

Clone returns a deep copy of r with its context changed to ctx. The provided
ctx must be non-nil.
For an outgoing client request, the context controls the entire lifetime of
a request and its response: obtaining a connection, sending the request, and
reading the response headers and body.


#### Request.Context

```go
func (r *Request) Context() context.Context
```

Context returns the request's context. To change the context, use
WithContext.
The returned context is always non-nil; it defaults to the background
context.
For outgoing client requests, the context controls cancellation.
For incoming server requests, the context is canceled when the client's
connection closes, the request is canceled (with HTTP/2), or when the
ServeHTTP method returns.


#### Request.Cookie

```go
func (r *Request) Cookie(name string) (*Cookie, error)
```

Cookie returns the named cookie provided in the request or ErrNoCookie if
not found. If multiple cookies match the given name, only one cookie will be
returned.


#### Request.Cookies

```go
func (r *Request) Cookies() []*Cookie
```

Cookies parses and returns the HTTP cookies sent with the request.


#### Request.FormFile

```go
func (r *Request) FormFile(key string) (multipart.File, *multipart.FileHeader, error)
```

FormFile returns the first file for the provided form key. FormFile calls
ParseMultipartForm and ParseForm if necessary.


#### Request.FormValue

```go
func (r *Request) FormValue(key string) string
```

FormValue returns the first value for the named component of the query. POST
and PUT body parameters take precedence over URL query string values.
FormValue calls ParseMultipartForm and ParseForm if necessary and ignores
any errors returned by these functions. If key is not present, FormValue
returns the empty string. To access multiple values of the same key, call
ParseForm and then inspect Request.Form directly.


#### Request.MultipartReader

```go
func (r *Request) MultipartReader() (*multipart.Reader, error)
```

MultipartReader returns a MIME multipart reader if this is a
multipart/form-data or a multipart/mixed POST request, else returns nil and
an error. Use this function instead of ParseMultipartForm to process the
request body as a stream.


#### Request.ParseForm

```go
func (r *Request) ParseForm() error
```

ParseForm populates r.Form and r.PostForm.
For all requests, ParseForm parses the raw query from the URL and updates
r.Form.
For POST, PUT, and PATCH requests, it also reads the request body, parses it
as a form and puts the results into both r.PostForm and r.Form. Request body
parameters take precedence over URL query string values in r.Form.
If the request Body's size has not already been limited by MaxBytesReader,
the size is capped at 10MB.
For other HTTP methods, or when the Content-Type is not
application/x-www-form-urlencoded, the request Body is not read, and
r.PostForm is initialized to a non-nil, empty value.
ParseMultipartForm calls ParseForm automatically. ParseForm is idempotent.


#### Request.ParseMultipartForm

```go
func (r *Request) ParseMultipartForm(maxMemory int64) error
```

ParseMultipartForm parses a request body as multipart/form-data. The whole
request body is parsed and up to a total of maxMemory bytes of its file
parts are stored in memory, with the remainder stored on disk in temporary
files. ParseMultipartForm calls ParseForm if necessary. After one call to
ParseMultipartForm, subsequent calls have no effect.


#### Request.PostFormValue

```go
func (r *Request) PostFormValue(key string) string
```

PostFormValue returns the first value for the named component of the POST,
PATCH, or PUT request body. URL query parameters are ignored. PostFormValue
calls ParseMultipartForm and ParseForm if necessary and ignores any errors
returned by these functions. If key is not present, PostFormValue returns
the empty string.


#### Request.ProtoAtLeast

```go
func (r *Request) ProtoAtLeast(major, minor int) bool
```

ProtoAtLeast reports whether the HTTP protocol used in the request is at
least major.minor.


#### Request.Referer

```go
func (r *Request) Referer() string
```

Referer returns the referring URL, if sent in the request.
Referer is misspelled as in the request itself, a mistake from the earliest
days of HTTP. This value can also be fetched from the Header map as
Header["Referer"]; the benefit of making it available as a method is that
the compiler can diagnose programs that use the alternate (correct English)
spelling req.Referrer() but cannot diagnose programs that use
Header["Referrer"].


#### Request.SetBasicAuth

```go
func (r *Request) SetBasicAuth(username, password string)
```

SetBasicAuth sets the request's Authorization header to use HTTP Basic
Authentication with the provided username and password.
With HTTP Basic Authentication the provided username and password are not
encrypted.
Some protocols may impose additional requirements on pre-escaping the
username and password. For instance, when used with OAuth2, both arguments
must be URL encoded first with url.QueryEscape.


#### Request.UserAgent

```go
func (r *Request) UserAgent() string
```

UserAgent returns the client's User-Agent, if sent in the request.


#### Request.WithContext

```go
func (r *Request) WithContext(ctx context.Context) *Request
```

WithContext returns a shallow copy of r with its context changed to ctx. The
provided ctx must be non-nil.
For outgoing client request, the context controls the entire lifetime of a
request and its response: obtaining a connection, sending the request, and
reading the response headers and body.
To create a new request with a context, use NewRequestWithContext. To change
the context of a request, such as an incoming request you want to modify
before sending back out, use Request.Clone. Between those two uses, it's
rare to need WithContext.


#### Request.Write

```go
func (r *Request) Write(w io.Writer) error
```

Write writes an HTTP/1.1 request, which is the header and body, in wire
format. This method consults the following fields of the request:
Host



## Interfaces

### CloseNotifier

```go
func (c ConnState) String() string
```

The CloseNotifier interface is implemented by ResponseWriters which allow
detecting when the underlying connection has gone away.
This mechanism can be used to cancel long operations on the server if the
client has disconnected before the response is ready.
Deprecated: the CloseNotifier interface predates Go's context package. New
code should use Request.Context instead.


### CookieJar

```go
func (d Dir) Open(name string) (File, error)
```

A CookieJar manages storage and use of cookies in HTTP requests.
Implementations of CookieJar must be safe for concurrent use by multiple
goroutines.
The net/http/cookiejar package provides a CookieJar implementation.


### File

```go
type File interface {
	io.Closer
	io.Reader
	io.Seeker
	Readdir(count int) ([]fs.FileInfo, error)
	Stat() (fs.FileInfo, error)
}
```

A File is returned by a FileSystem's Open method and can be served by the
FileServer implementation.
The methods should behave the same as those on an *os.File.


### FileSystem

```go
type FileSystem interface {
	Open(name string) (File, error)
}
```

A FileSystem implements access to a collection of named files. The elements
in a file path are separated by slash ('/', U+002F) characters, regardless
of host operating system convention. See the FileServer function to convert
a FileSystem to a Handler.
This interface predates the fs.FS interface, which can be used instead: the
FS adapter function converts an fs.FS to a FileSystem.
FS converts fsys to a FileSystem implementation, for use with FileServer and
NewFileTransport.


### Flusher

```go
type Flusher interface {
	// Flush sends any buffered data to the client.
	Flush()
}
```

The Flusher interface is implemented by ResponseWriters that allow an HTTP
handler to flush buffered data to the client.
The default HTTP/1.x and HTTP/2 ResponseWriter implementations support
Flusher, but ResponseWriter wrappers may not. Handlers should always test
for this ability at runtime.
Note that even for ResponseWriters that support Flush, if the client is
connected through an HTTP proxy, the buffered data may not reach the client
until the response completes.


### Handler

```go
func (h Header) WriteSubset(w io.Writer, exclude map[string]bool) error
```

A Handler responds to an HTTP request.
ServeHTTP should write reply headers and data to the ResponseWriter and then
return. Returning signals that the request is finished; it is not valid to
use the ResponseWriter or read from the Request.Body after or concurrently
with the completion of the ServeHTTP call.
Depending on the HTTP client software, HTTP protocol version, and any
intermediaries between the client and the Go server, it may not be possible
to read from the Request.Body after writing to the ResponseWriter. Cautious
handlers should read the Request.Body first, and then reply.
Except for reading the body, handlers should not modify the provided
Request.
If ServeHTTP panics, the server (the caller of ServeHTTP) assumes that the
effect of the panic was isolated to the active request. It recovers the
panic, logs a stack trace to the server error log, and either closes the
network connection or sends an HTTP/2 RST_STREAM, depending on the HTTP
protocol. To abort a handler so the client sees an interrupted response but
the server doesn't log an error, panic with the value ErrAbortHandler.
FileServer returns a handler that serves HTTP requests with the contents of
the file system rooted at root.
As a special case, the returned file server redirects any request ending in
"/index.html" to the same path, without the final "index.html".
To use the operating system's file system implementation, use http.Dir:
http.Handle("/", http.FileServer(http.Dir("/tmp")))
To use an fs.FS implementation, use http.FS to convert it:
http.Handle("/", http.FileServer(http.FS(fsys)))
NotFoundHandler returns a simple request handler that replies to each
request with a “404 page not found” reply.
RedirectHandler returns a request handler that redirects each request it
receives to the given url using the given status code.
The provided code should be in the 3xx range and is usually
StatusMovedPermanently, StatusFound or StatusSeeOther.
StripPrefix returns a handler that serves HTTP requests by removing the
given prefix from the request URL's Path (and RawPath if set) and invoking
the handler h. StripPrefix handles a request for a path that doesn't begin
with prefix by replying with an HTTP 404 not found error. The prefix must
match exactly: if the prefix in the request contains escaped characters the
reply is also an HTTP 404 not found error.
TimeoutHandler returns a Handler that runs h with the given time limit.
The new Handler calls h.ServeHTTP to handle each request, but if a call runs
for longer than its time limit, the handler responds with a 503 Service
Unavailable error and the given message in its body. (If msg is empty, a
suitable default message will be sent.) After such a timeout, writes by h to
its ResponseWriter will return ErrHandlerTimeout.
TimeoutHandler supports the Pusher interface but does not support the
Hijacker or Flusher interfaces.


### Hijacker

```go
type Hijacker interface {
	// Hijack lets the caller take over the connection.
	// After a call to Hijack the HTTP server library
	// will not do anything else with the connection.
	//
	// It becomes the caller's responsibility to manage
	// and close the connection.
	//
	// The returned net.Conn may have read or write deadlines
	// already set, depending on the configuration of the
	// Server. It is the caller's responsibility to set
	// or clear those deadlines as needed.
	//
	// The returned bufio.Reader may contain unprocessed buffered
	// data from the client.
	//
	// After a call to Hijack, the original Request.Body must not
	// be used. The original Request's Context remains valid and
	// is not canceled until the Request's ServeHTTP method
	// returns.
	Hijack() (net.Conn, *bufio.ReadWriter, error)
}
```

The Hijacker interface is implemented by ResponseWriters that allow an HTTP
handler to take over the connection.
The default ResponseWriter for HTTP/1.x connections supports Hijacker, but
HTTP/2 connections intentionally do not. ResponseWriter wrappers may also
not support Hijacker. Handlers should always test for this ability at
runtime.


### Pusher

```go
type Pusher interface {
	// Push initiates an HTTP/2 server push. This constructs a synthetic
	// request using the given target and options, serializes that request
	// into a PUSH_PROMISE frame, then dispatches that request using the
	// server's request handler. If opts is nil, default options are used.
	//
	// The target must either be an absolute path (like "/path") or an absolute
	// URL that contains a valid host and the same scheme as the parent request.
	// If the target is a path, it will inherit the scheme and host of the
	// parent request.
	//
	// The HTTP/2 spec disallows recursive pushes and cross-authority pushes.
	// Push may or may not detect these invalid pushes; however, invalid
	// pushes will be detected and canceled by conforming clients.
	//
	// Handlers that wish to push URL X should call Push before sending any
	// data that may trigger a request for URL X. This avoids a race where the
	// client issues requests for X before receiving the PUSH_PROMISE for X.
	//
	// Push will run in a separate goroutine making the order of arrival
	// non-deterministic. Any required synchronization needs to be implemented
	// by the caller.
	//
	// Push returns ErrNotSupported if the client has disabled push or if push
	// is not supported on the underlying connection.
	Push(target string, opts *PushOptions) error
}
```

Pusher is the interface implemented by ResponseWriters that support HTTP/2
server push. For more background, see
https://tools.ietf.org/html/rfc7540#section-8.2.



