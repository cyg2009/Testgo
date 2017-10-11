package main

import (
    "net/http"  
    "io/ioutil"
    "sync"
    "time"
    "fmt"
)

func makeOKResponse(w http.ResponseWriter, body string){
    
    //w.Header().Set("Content-Type", "application/json")   
    w.Write([]byte(body))
}
func makeFailedResponse(w http.ResponseWriter, statusCode int, message string) {

   // bodyContent := slscommon.NewErrorMessage(statusCode, message).ToJsonString()    
    //w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(statusCode)
    w.Write([]byte(message))
}

func ServeHTTPInvoke(w http.ResponseWriter, req *http.Request) {
    if req.Method != "POST" {
        makeFailedResponse(w, http.StatusInternalServerError, "Please POST this request.")
        return 
    }
    evt, err := ioutil.ReadAll(req.Body)
    if err != nil {
        makeFailedResponse(w, http.StatusInternalServerError, err.Error())
        return 
    }

    functionId := req.Header.Get("function")
    if len(functionId) > 0 {        
        if _, ok := GetFunctionManager().GetFunction(functionId); ok == false {
             makeFailedResponse(w, http.StatusBadRequest, "Function " + functionId + "  not exists!")
             return 
        }
                
        respData, _ := GetFunctionManager().ExecuteFunction(functionId, evt)     
        fmt.Println("invoke result:" + respData)
        makeOKResponse(w, respData)
        return
    }
        
    makeFailedResponse(w,  http.StatusBadRequest, "Function header not spcified!")      
}

func ServeHTTPConfig(w http.ResponseWriter, req *http.Request) {
    makeOKResponse(w, "TBD:Configuration")
}

func ServeHTTPInfo(w http.ResponseWriter, req *http.Request) {
    info := GetFunctionManager().GetAllFunctionsJSON()
    makeOKResponse(w, info)
}
func ServeHTTPHealthCheck(w http.ResponseWriter, req *http.Request) {
    t := time.Now()
    body := "OK " + t.Format("20060102150405")
    makeOKResponse(w, body)
}

// Singleton 
var instanceHTTP *http.ServeMux
var oncehttp sync.Once

func GetPrserviceHttpHandler() (*http.ServeMux) {

    oncehttp.Do( func() {
        instanceHTTP = http.NewServeMux()   
        instanceHTTP.HandleFunc("/health", ServeHTTPHealthCheck)
        instanceHTTP.HandleFunc("/config", ServeHTTPConfig)
        instanceHTTP.HandleFunc("/info", ServeHTTPInfo)
        instanceHTTP.HandleFunc("/invoke", ServeHTTPInvoke)
    })

    return instanceHTTP
}