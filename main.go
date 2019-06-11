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
	"strings"

	"encoding/base64"
	"encoding/json"
	"fmt"
	xj "github.com/basgys/goxml2json"
	"net/http"
)

func init() {
	cfg := &tls.Config{
		InsecureSkipVerify: true,
	}
	http.DefaultClient.Transport = &http.Transport{
		TLSClientConfig: cfg,
	}
}
func main() {
	err := gotenv.Load(".env.gateway")
	if err != nil {
		fmt.Println("error: main" + err.Error())

	}
	erro := os.RemoveAll(os.Getenv("ENV_CLONE") + string(os.PathSeparator) + "RootNode" + string(os.PathSeparator))
	if erro != nil {
		fmt.Println("error: main" + erro.Error())

	}
	router := mux.NewRouter()
	router.HandleFunc("/dev", GetDev).Methods("GET")
	router.HandleFunc("/uat", GetUat).Methods("GET")
	fmt.Println("Starting the application...")
	er := http.ListenAndServe(":8081", router)
	if er != nil {
		fmt.Println("error: main" + er.Error())

	}
}
func GetUat(writer http.ResponseWriter, request *http.Request) {
	//os.RemoveAll(os.Getenv("ENV_CLONE")+string(os.PathSeparator)+"RootNode/*")

	url, env := getenviroment("UAT")
	recursivecall(url+"restman/1.0/folders/0000000000000000ffffffffffffec76/dependencies?level=1", "", env)
	//GitWorld("UAT")
	encode := json.NewEncoder(writer).Encode("Trabajo realizado")
	if encode != nil {

		fmt.Println("error: GetUat" + encode.Error())

	}
}

func GetDev(writer http.ResponseWriter, request *http.Request) {
	//os.RemoveAll(os.Getenv("ENV_CLONE")+string(os.PathSeparator)+"RootNode/*")

	url, env := getenviroment("DEV")
	recursivecall(url+"restman/1.0/folders/0000000000000000ffffffffffffec76/dependencies?level=1", "", env)
	//GitWorld("DEV")
	json.NewEncoder(writer).Encode("Trabajo realizado")
}
func recursivecall(url string, folder string, env string) {
	var resp Entidades.TheContent
	var carpeta string
	jsonresp := calls(url)
	err := json.Unmarshal(jsonresp.Bytes(), &resp)
	if err != nil {
		var respaux Entidades.Thecontentespecial
		err := json.Unmarshal(jsonresp.Bytes(), &respaux)
		if err != nil {
			fmt.Println("error: recursivecall" + err.Error())
		}
		carpeta = cleanString(respaux.Resource.DependencyList.Reference.Name)

	} else {
		carpeta = cleanString(resp.Item.Resource.DependencyList.Reference.Name)
	}
	folder = folder + string(os.PathSeparator) + carpeta
	for _, element := range resp.Item.Resource.DependencyList.Reference.Dependencies.Dependency {
		if element.Type == "FOLDER" {
			fmt.Println("INFO: " + makeurl(element.Type, element.ID, env))
			recursivecall(makeurl(element.Type, element.ID, env), folder, env)
		} else {
			if element.Type == "SERVICE" || element.Type == "POLICY" {
				fmt.Println("INFO: " + makeurl(element.Type, element.ID, env))
				thecallandsave(makeurl(element.Type, element.ID, env), element.Type, folder,env)
			}
		}
	}
}
func calls(url string) *bytes.Buffer {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("error: " + err.Error())
	}
	//fmt.Println(url)
	//fmt.Println(req)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Basic "+basicAuth())
	req.Header.Add("cache-control", "no-cache")
	req.Header.Add("Postman-Token", "621b450c-cc8c-4d83-b0fb-6f69fbd7987c")
	res, erra := http.DefaultClient.Do(req)
	if erra != nil {
		fmt.Println("error: " + err.Error())
	}
	//fmt.Println(res)
	defer res.Body.Close()

	jsonresp, err := xj.Convert(res.Body)
	if err != nil {
		fmt.Println("error: " + err.Error())

	}
	return jsonresp
}
func thecallandsave(url,typee,directory,env string,) {
	var archivo,content string
	jsonrespinfo := calls(url)
	if typee == "SERVICE" {
		var respuestainfoservice Entidades.Serviceinfo
		err := json.Unmarshal(jsonrespinfo.Bytes(), &respuestainfoservice)
		if len(respuestainfoservice.Item.Link) > 1 {
			archivo = cleanStringandgetthename(respuestainfoservice.Item.Link[3].URI,typee,env)

		}
		content = respuestainfoservice.Item.Resource.PolicyVersion.XML
		if err != nil {
				fmt.Println("error: " + err.Error())
		}
	}
	if typee == "POLICY" {
		var respuestainfopolicy Entidades.Policyinfo
		err := json.Unmarshal(jsonrespinfo.Bytes(), &respuestainfopolicy)
		if len(respuestainfopolicy.Item.Link) > 1 {
			archivo = cleanStringandgetthename(respuestainfopolicy.Item.Link[3].URI,typee,env)

		}
		content = respuestainfopolicy.Item.Resource.PolicyVersion.XML

		if err != nil {
			fmt.Println("error: " + err.Error())
		}
	}
	createFile(os.Getenv("ENV_CLONE")+string(os.PathSeparator)+directory, string(os.PathSeparator)+typee+"-"+archivo+".xml", content)

}
func cleanStringandgetthename(url,typee,env string) string {
	var content string
	urlco, _ := getenviroment(env)


	jsonresp:=calls(remplaceip(url,urlco))
	if typee == "SERVICE" {
		var respuestainfoservice Entidades.ServiceinfoOLD
		err := json.Unmarshal(jsonresp.Bytes(), &respuestainfoservice)
		content = respuestainfoservice.Item.Name
		if err != nil {
			fmt.Println("error: cleanStringandgetthename" + err.Error())
		}
	}
	if typee == "POLICY" {
		var respuestainfopolicy Entidades.Policyinfo
		err := json.Unmarshal(jsonresp.Bytes(), &respuestainfopolicy)
		content = respuestainfopolicy.Item.Name

		if err != nil {
			fmt.Println("error: cleanStringandgetthename" + err.Error())
		}
	}
	return content
}

func remplaceip(s string, s2 string) string {
	myText := strings.Replace(s, "https://CDV1UTAPIGWIN01:8443/", s2, -1)
	fmt.Println("te quiero ver camviar"+myText)
	return myText

}
func createFile(path string, name string, contenido string) {
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		fmt.Println("error: createFile" + err.Error())
	}
	f, err := os.OpenFile(path+name, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()
	if err != nil {
		fmt.Println("error: createFile" + err.Error())

	}
	l, err := f.WriteString(contenido)
	if err != nil {
		fmt.Println("error: createFile" + err.Error())
		f.Close()

	}
	fmt.Println(l, "bytes escritos SUCESS")
	ERRA := f.Sync()
	if ERRA != nil {
		fmt.Println("error: createFile" + err.Error())
	}
	err = f.Close()
	if err != nil {
		fmt.Println("error: createFile" + err.Error())
	}
}
func cleanString(s string) string {
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		log.Fatal(err)
	}
	return reg.ReplaceAllString(s, "")

}
func makeurl(tipo string, id string, env string) string {
	var output string
	url, _ := getenviroment(env)
	switch tipo {
	case "FOLDER":
		output = url + "restman/1.0/folders/" + id + "/dependencies?level=1"
	case "SERVICE":
		output = url + "restman/1.0/services/" + id + "/versions/active"
	case "POLICY":
		output = url + "restman/1.0/policies/" + id + "/versions/active"
	}
	return output
}
func basicAuth() string {
	auth := os.Getenv("ENV_USER") + ":" + os.Getenv("ENV_PASS")
	return base64.StdEncoding.EncodeToString([]byte(auth))
}
func GitWorld(Enviroment string) {
	if runtime.GOOS == "windows" {
		value, err := exec.Command("git-work.bat", os.Getenv("ENV_CLONE"), Enviroment).Output()
		if err != nil {
			fmt.Println("error: " + err.Error())
		}
		println(string(value))
	}
	if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
		_, err := exec.Command("git-work.sh", os.Getenv("ENV_CLONE"), Enviroment).Output()
		if err != nil {
			fmt.Println("error: " + err.Error())
		}
		//println(string(value))
	}
}
func getenviroment(s string) (url string, string2 string) {
	var output string
	switch s {
	case "DEV":
		output = "https://10.49.22.7:8443/"
	case "UAT":
		output = "https://10.49.22.14:8443/"

	}
	return output, s

}
