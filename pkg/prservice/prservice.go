package prservice

import (
    "net/http"  
    "io/ioutil"
    "sync"
    "time"
    "fmt"
    fm "serverlessgo/pkg/functionmanager"
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
        if _, ok := fm.GetFunctionManager().GetFunction(functionId); ok == false {
             makeFailedResponse(w, http.StatusBadRequest, "Function " + functionId + "  not exists!")
             return 
        }
                
        respData, _ := fm.GetFunctionManager().ExecuteFunction(functionId, evt)     
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
    info := fm.GetFunctionManager().GetAllFunctionsJSON()
    makeOKResponse(w, info)
}
func ServeHTTPHealthCheck(w http.ResponseWriter, req *http.Request) {
    t := time.Now()
    body := "OK " + t.Format("20060102150405")
    makeOKResponse(w, body)
}

// Singleton 
var instance *http.ServeMux
var once sync.Once

func GetPrserviceHttpHandler() (*http.ServeMux) {

    once.Do( func() {
        instance = http.NewServeMux()   
        instance.HandleFunc("/health", ServeHTTPHealthCheck)
        instance.HandleFunc("/config", ServeHTTPConfig)
        instance.HandleFunc("/info", ServeHTTPInfo)
        instance.HandleFunc("/invoke", ServeHTTPInvoke)
    })

    return instance
}