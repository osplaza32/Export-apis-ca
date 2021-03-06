package main

import (
	"bytes"
	"crypto/tls"
	"github.com/gorilla/mux"
	"github.com/subosito/gotenv"
	"log"
	"os"
	"os/exec"
	"osplaza32/ExtractGolang/Entidades"
	"regexp"
	"runtime"

	//"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	xj "github.com/basgys/goxml2json"
	"net/http"
)
func init()  {
	cfg := &tls.Config{
		InsecureSkipVerify: true,
	}
	http.DefaultClient.Transport = &http.Transport{
		TLSClientConfig: cfg,
	}

}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/dev", GetDev).Methods("GET")
	router.HandleFunc("/uat", GetUat).Methods("GET")
	fmt.Println("Starting the application...")
	http.ListenAndServe(":8081", router)
}
func GetUat(writer http.ResponseWriter, request *http.Request) {
	gotenv.Load(".env.gatewayuat")
	recursivecall(os.Getenv("ENV_URL")+"restman/1.0/folders/0000000000000000ffffffffffffec76/dependencies?level=1","")
	GitWorld()
	json.NewEncoder(writer).Encode("Trabajo realizado")
}
func GetDev(writer http.ResponseWriter, request *http.Request) {
	gotenv.Load(".env.gatewaydev")
	recursivecall(os.Getenv("ENV_URL")+"restman/1.0/folders/0000000000000000ffffffffffffec76/dependencies?level=1","")
	GitWorld()
	json.NewEncoder(writer).Encode("Trabajo realizado")
}
func recursivecall(url string,folder string){
	var resp Entidades.TheContent
	jsonresp := calls(url)
	json.Unmarshal(jsonresp.Bytes(), &resp)
	carpeta :=cleanString(resp.Item.Resource.DependencyList.Reference.Name)
	folder = folder +string(os.PathSeparator) +carpeta
	for _, element := range resp.Item.Resource.DependencyList.Reference.Dependencies.Dependency {
		if element.Type == "FOLDER"{
			recursivecall(makeurl(element.Type,element.ID),folder)
		}else {
			if element.Type == "SERVICE" || element.Type == "POLICY" {
				thecallandsave(makeurl(element.Type,element.ID),element.Type,folder)
				}}}
}
func calls(url string) *bytes.Buffer {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Basic "+basicAuth())
	req.Header.Add("cache-control", "no-cache")
	req.Header.Add("Postman-Token", "621b450c-cc8c-4d83-b0fb-6f69fbd7987c")
	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()
	jsonresp, _ := xj.Convert(res.Body)
	return jsonresp
	}
func thecallandsave(url string,typee string,directory string) {
	jsonrespinfo:=calls(url)
	if typee == "SERVICE" {
		var respuestainfoservice Entidades.Serviceinfo
		json.Unmarshal(jsonrespinfo.Bytes(), &respuestainfoservice)
		archivo :=cleanString(respuestainfoservice.Item.Name)
		createFile(os.Getenv("ENV_CLONE")+string(os.PathSeparator) +directory,string(os.PathSeparator) +"SERVICE-"+archivo+".xml",respuestainfoservice.Item.Resource.Service.Resources.ResourceSet.Resource.Content)
		}
	if typee == "POLICY"{
		var respuestainfopolicy Entidades.Policyinfo
		json.Unmarshal(jsonrespinfo.Bytes(), &respuestainfopolicy)
		archivo :=cleanString(respuestainfopolicy.Item.Name)
		createFile(os.Getenv("ENV_CLONE")+string(os.PathSeparator) +directory,string(os.PathSeparator) +"POLICY-"+archivo+".xml",respuestainfopolicy.Item.Resource.Policy.Resources.ResourceSet.Resource.Content)
		}
	}
func createFile(path string,name string,contenido string) {
	os.MkdirAll(path,os.ModePerm)
	os.Remove(path+name)
	f, err := os.OpenFile(path+name, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
	l, err := f.WriteString(contenido)
	if err != nil {
		fmt.Println(err)
		f.Close()
		return
	}
	fmt.Println(l, "bytes escritos")
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	}
func cleanString(s string) string{
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		log.Fatal(err)
	}
	return reg.ReplaceAllString(s,"")


}
func makeurl(tipo string, id string) string {
	var output string
	switch tipo {
	case "FOLDER":
			output = os.Getenv("ENV_URL")+"restman/1.0/folders/"+id+"/dependencies?level=1"
	case "SERVICE":
		output = os.Getenv("ENV_URL")+"restman/1.0/services/"+id
	case "POLICY":
		output = os.Getenv("ENV_URL")+"restman/1.0/policies/"+id
	}
	return output
}
func basicAuth() string {
	auth := os.Getenv("ENV_USER") + ":" + os.Getenv("ENV_PASS")
	return base64.StdEncoding.EncodeToString([]byte(auth))
}
func GitWorld(){
	fmt.Println("Hola")
	if runtime.GOOS == "windows" {
		value, err:= exec.Command("git-work.bat",os.Getenv("ENV_CLONE")).Output()
		if err != nil{
			fmt.Println("aun no")
		}
		println(string(value))
	}
	if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
		value, err:= exec.Command("git-work.sh",os.Getenv("ENV_CLONE")).Output()
		if err != nil{
			fmt.Println("aun no")
		}
		println(string(value))
	}
}

