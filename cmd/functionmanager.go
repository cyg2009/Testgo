package main

import (
    "errors"
    "sync"
    "fmt"
	"encoding/json"
	"plugin"
	"strings"
	 "os"

)

type ServerlessFunction struct{
    id string
    functionpath string
    started bool
    trigger func(string) string
}

func (sf *ServerlessFunction) Start(){

   if !sf.started {
        fmt.Println("Load plugin:" + sf.functionpath)
        p, err := plugin.Open(sf.functionpath)
        if err != nil {
        	fmt.Println(err.Error())
            return 
        }

        greetSymbol, err := p.Lookup("Handler")
        if err != nil {
        	fmt.Println(err.Error())
            return 
        }

        sf.trigger = greetSymbol.(func(string) string)
        sf.started = true            
   } 
}

func (sf *ServerlessFunction) Stop(){

   if sf.started {
        sf.started = false
   } 
}

func (sf *ServerlessFunction) Trigger (event []byte) (string, error) {
    
    if !sf.started {
        sf.Start()
    }
   
    
    return sf.trigger(string(event[:])), nil
}


type ServerlessFunctionManager struct{
    functionStore map[string]*ServerlessFunction
}


func (mgr *ServerlessFunctionManager) LoadFunction(functionId string) *ServerlessFunction {
	temp := strings.Split(functionId, ":")
    funcName := ""
    funcTag := "latest"
    if len(temp) > 1 {
        l := len(temp)
        funcTag = temp[l - 1]
        funcName = temp[l - 2]
    }

    dest := os.Getenv("RUNTIME_ROOT")
    if len(dest) == 0 {
        dest = "/var/runtime"
    }   
    dest += "/func/" + functionId
    fmt.Println("function path:" + dest)
    if _, err := os.Stat(dest); os.IsNotExist(err) {
    	fmt.Println(err.Error())
        return nil
    } 
   
    funcpath := dest + "/" + funcName + ".so." + funcTag
    fmt.Println("funcpath:" + funcpath)
    
    sf := &ServerlessFunction{
        id: functionId,   
        functionpath: funcpath,
        started: false, 
        trigger: nil,
    }


    return sf
}

func (mgr *ServerlessFunctionManager) GetFunction (functionId string) (*ServerlessFunction, bool){
	
    ff, ok := mgr.functionStore[functionId]   
    if ok {
        return ff, ok
    }
    fmt.Println("load function:" + functionId)
    ff = mgr.LoadFunction(functionId)
    
    if ff == nil {
        return nil, false
    }
    
    mgr.functionStore[functionId]  = ff
    
    return ff, true
}


func (mgr *ServerlessFunctionManager) ExecuteFunction (functionId string, event []byte) (string, error) {
    
    ff, ok := mgr.functionStore[functionId]
    if ok {
        fmt.Println("trigger function:" + functionId)	
        fmt.Println("trigger event:" + string(event))
        return ff.Trigger(event)
    
    }

    return  "", errors.New(functionId + " not exists!")
}

func (mgr *ServerlessFunctionManager) GetAllFunctionsJSON () (string){
    ret, err := json.Marshal(mgr.functionStore)

    if err != nil {
        return err.Error()
    }

    return string(ret)
}


// Singleton 
var instanceFM *ServerlessFunctionManager
var oncefm sync.Once

func GetFunctionManager() (*ServerlessFunctionManager) {

    oncefm.Do( func() {
       
            instanceFM = &ServerlessFunctionManager {
                functionStore: make(map[string]*ServerlessFunction),
            }        
    })

    return instanceFM
}


