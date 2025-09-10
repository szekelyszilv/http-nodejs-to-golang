package main

/*
#include <stdint.h>
*/
import "C"

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"runtime/cgo"
	"slices"
	"strings"
	"unsafe"

	"http-nodejs-to-golang/internal/handlers"
)

var handler http.Handler

func main() {}

func init() {
	fmt.Println("Go init")
	router := handlers.CreateRouter()
	handler = router
}

type HttpPair struct {
	Request  *http.Request
	Response *httptest.ResponseRecorder
}

//export GoHandlerNew
func GoHandlerNew(methodPtr *C.char, urlPtr *C.char, headersPtr *C.char, bodyLen C.size_t, bodyPtr unsafe.Pointer, errPtr **C.char) C.uintptr_t {
	method := C.GoString(methodPtr)
	url := C.GoString(urlPtr)
	headersStr := C.GoString(headersPtr)
	body := slices.Clone(unsafe.Slice((*byte)(bodyPtr), bodyLen))

	req, err := http.NewRequest(method, url, bytes.NewReader(body))
	if err != nil {
		writeErrorToC(err, errPtr)
		return 0
	}

	err = json.NewDecoder(strings.NewReader(headersStr)).Decode(&req.Header)
	if err != nil {
		writeErrorToC(err, errPtr)
		return 0
	}

	httpPair := &HttpPair{
		Request:  req,
		Response: httptest.NewRecorder(),
	}
	handle := cgo.NewHandle(httpPair)

	fmt.Printf("Created new handle: %v\n", handle)

	return C.uintptr_t(handle)
}

//export GoHandlerFree
func GoHandlerFree(handle C.uintptr_t) {
	fmt.Printf("Freeing handle: %v\n", handle)
	(cgo.Handle)(handle).Delete()
}

//export GoHandlerRun
func GoHandlerRun(handle C.uintptr_t) {
	ref := (cgo.Handle(handle).Value()).(*HttpPair)
	handler.ServeHTTP(ref.Response, ref.Request)
}

//export GoHandlerGetResponseStatusCode
func GoHandlerGetResponseStatusCode(handle C.uintptr_t) C.int {
	ref := ((cgo.Handle)(handle).Value()).(*HttpPair)
	return C.int(ref.Response.Code)
}

//export GoHandlerGetResponseHeaders
func GoHandlerGetResponseHeaders(handle C.uintptr_t) *C.char {
	ref := (cgo.Handle(handle).Value()).(*HttpPair)
	headerBytes, _ := json.Marshal(ref.Response.Header())
	return C.CString(string(headerBytes))
}

//export GoHandlerGetResponseBodySize
func GoHandlerGetResponseBodySize(handle C.uintptr_t) C.size_t {
	ref := (cgo.Handle(handle).Value()).(*HttpPair)
	if ref.Response.Body == nil {
		return 0
	}
	return C.size_t(ref.Response.Body.Len())
}

//export GoHandlerGetResponseBodyBytes
func GoHandlerGetResponseBodyBytes(handle C.uintptr_t, buf unsafe.Pointer, buflen C.size_t) C.int {
	ref := (cgo.Handle(handle).Value()).(*HttpPair)
	if ref.Response.Body == nil {
		return 0
	}
	bodyBytes := ref.Response.Body.Bytes()
	if len(bodyBytes) > int(buflen) {
		return 1
	}
	copy(unsafe.Slice((*byte)(buf), buflen), bodyBytes)
	return 0
}

func writeErrorToC(err error, errPtr **C.char) {
	if err == nil {
		return
	}
	errMsg := C.CString(err.Error())
	*errPtr = errMsg
}
