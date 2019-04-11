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
	gotenv.Load(".env.gatewayuat")
	os.RemoveAll(os.Getenv("ENV_CLONE")+string(os.PathSeparator)+"RootNode/*")
	router := mux.NewRouter()
	router.HandleFunc("/dev", GetDev).Methods("GET")
	router.HandleFunc("/uat", GetUat).Methods("GET")
	fmt.Println("Starting the application...")
	http.ListenAndServe(":8081", router)
}
func GetUat(writer http.ResponseWriter, request *http.Request) {
	//os.RemoveAll(os.Getenv("ENV_CLONE")+string(os.PathSeparator)+"RootNode/*")

	url,env := getenviroment("UAT")
	recursivecall(url+"restman/1.0/folders/0000000000000000ffffffffffffec76/dependencies?level=1","",env)
	GitWorld("UAT")
	json.NewEncoder(writer).Encode("Trabajo realizado")
}

func GetDev(writer http.ResponseWriter, request *http.Request) {
	//os.RemoveAll(os.Getenv("ENV_CLONE")+string(os.PathSeparator)+"RootNode/*")

	url,env := getenviroment("DEV")
	recursivecall(url+"restman/1.0/folders/0000000000000000ffffffffffffec76/dependencies?level=1","",env)
	GitWorld("DEV")
	json.NewEncoder(writer).Encode("Trabajo realizado")
}
func recursivecall(url string,folder string,env string){
	var resp Entidades.TheContent
	jsonresp := calls(url)
	json.Unmarshal(jsonresp.Bytes(), &resp)
	carpeta :=cleanString(resp.Item.Resource.DependencyList.Reference.Name)
	folder = folder +string(os.PathSeparator) +carpeta
	for _, element := range resp.Item.Resource.DependencyList.Reference.Dependencies.Dependency {
		if element.Type == "FOLDER"{
			fmt.Println(makeurl(element.Type,element.ID,env))
			recursivecall(makeurl(element.Type,element.ID,env),folder,env)
		}else {
			if element.Type == "SERVICE" || element.Type == "POLICY" {
				fmt.Println(makeurl(element.Type,element.ID,env))
				thecallandsave(makeurl(element.Type,element.ID,env),element.Type,folder)
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
		err:=json.Unmarshal(jsonrespinfo.Bytes(), &respuestainfoservice)
		archivo :=cleanString(respuestainfoservice.Item.Name)
		createFile(os.Getenv("ENV_CLONE")+string(os.PathSeparator) +directory,string(os.PathSeparator) +"SERVICE-"+archivo+".xml",respuestainfoservice.Item.Resource.Service.Resources.ResourceSet.Resource.Content)
		}
	if typee == "POLICY"{
		var respuestainfopolicy Entidades.Policyinfo
		err := json.Unmarshal(jsonrespinfo.Bytes(), &respuestainfopolicy)
		archivo :=cleanString(respuestainfopolicy.Item.Name)
		createFile(os.Getenv("ENV_CLONE")+string(os.PathSeparator) +directory,string(os.PathSeparator) +"POLICY-"+archivo+".xml",respuestainfopolicy.Item.Resource.Policy.Resources.ResourceSet.Resource.Content)
		}
	}
func createFile(path string,name string,contenido string) {
	errdel := os.RemoveAll(path + name)
	errcrete := os.MkdirAll(path, os.ModePerm)
	f, err := os.OpenFile(path+name, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()
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
	f.Sync()
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
func makeurl(tipo string, id string,env string) string {
	var output string
	url,_ := getenviroment(env)
	switch tipo {
	case "FOLDER":
			output = url+"restman/1.0/folders/"+id+"/dependencies?level=1"
	case "SERVICE":
		output = url+"restman/1.0/services/"+id
	case "POLICY":
		output = url+"restman/1.0/policies/"+id
	}
	return output
}
func basicAuth() string {
	auth := os.Getenv("ENV_USER") + ":" + os.Getenv("ENV_PASS")
	return base64.StdEncoding.EncodeToString([]byte(auth))
}
func GitWorld(Enviroment string){
	if runtime.GOOS == "windows" {
		value, err:= exec.Command("git-work.bat",os.Getenv("ENV_CLONE"),Enviroment).Output()
		if err != nil{
			fmt.Println(err.Error())
		}
		println(string(value))
	}
	if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
		_, err:= exec.Command("git-work.sh",os.Getenv("ENV_CLONE"),Enviroment).Output()
		if err != nil{
			fmt.Println(err.Error())
		}
		//println(string(value))
	}
}
func getenviroment(s string)(url string,string2 string) {
	var output string
	switch s {
	case "DEV":
		output = "https://10.49.22.7:8443/"
	case "UAT":
		output = "https://10.49.22.14:8443/"

	}
	return output,s

}

