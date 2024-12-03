package apiclient

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	errorhandling "github.com/jala-R/VideoAutomatorGUI/packages/ErrorHandling"
	"github.com/jala-R/VideoAutomatorGUI/packages/GUI/model"
)

//enpints
//optimize script done
//createVideoProject done
//get platform  done
//get profile done
//get key done
//rotate key done
//Add profile done
//Add key

func init() {
	model.AddToDb(model.APIURL, "192.168.29.2:8000")
}

func AddKey(platform, profile, newKey string) {
	apiUrl := getUrl()
	if apiUrl == "" {
		errorhandling.HandleError(errors.New("no url set in db"))
		return
	}

	fullUrl := fmt.Sprintf("http://%s/addKey", apiUrl)
	parseUrl, err := url.Parse(fullUrl)
	if err != nil {
		errorhandling.HandleError(err)
		return
	}
	q := parseUrl.Query()
	q.Add("platform", platform)
	q.Add("profile", profile)
	q.Add("key", newKey)

	parseUrl.RawQuery = q.Encode()

	resp, err := http.Get(parseUrl.String())
	if err != nil {
		errorhandling.HandleError(err)
		return
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		errorhandling.HandleError(fmt.Errorf("rotate key status code %d, message %s", resp.StatusCode, resp.Status))
	}

}

func AddProfile(platform, newProfile string) {
	apiUrl := getUrl()
	if apiUrl == "" {
		errorhandling.HandleError(errors.New("no url set in db"))
		return
	}

	fullUrl := fmt.Sprintf("http://%s/addProfile", apiUrl)
	parseUrl, err := url.Parse(fullUrl)
	if err != nil {
		errorhandling.HandleError(err)
		return
	}
	q := parseUrl.Query()
	q.Add("platform", platform)
	q.Add("profile", newProfile)

	parseUrl.RawQuery = q.Encode()

	resp, err := http.Get(parseUrl.String())
	if err != nil {
		errorhandling.HandleError(err)
		return
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		errorhandling.HandleError(fmt.Errorf("rotate key status code %d, message %s", resp.StatusCode, resp.Status))
	}

}

func RotateKey(platform, profile string) {
	apiUrl := getUrl()
	if apiUrl == "" {
		errorhandling.HandleError(errors.New("no url set in db"))
		return
	}

	fullUrl := fmt.Sprintf("http://%s/rotateKey", apiUrl)
	parseUrl, err := url.Parse(fullUrl)
	if err != nil {
		errorhandling.HandleError(err)
		return
	}
	q := parseUrl.Query()
	q.Add("platform", platform)
	q.Add("profile", profile)

	parseUrl.RawQuery = q.Encode()

	resp, err := http.Get(parseUrl.String())
	if err != nil {
		errorhandling.HandleError(err)
		return
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		errorhandling.HandleError(fmt.Errorf("rotate key status code %d, message %s", resp.StatusCode, resp.Status))
	}

}

func GetKey(platform, profile string) string {
	apiUrl := getUrl()
	if apiUrl == "" {
		errorhandling.HandleError(errors.New("no url set in db"))
		return ""
	}

	fullUrl := fmt.Sprintf("http://%s/getKey", apiUrl)
	parseUrl, err := url.Parse(fullUrl)
	if err != nil {
		errorhandling.HandleError(err)
		return ""
	}
	q := parseUrl.Query()
	q.Add("platform", platform)
	q.Add("profile", profile)

	parseUrl.RawQuery = q.Encode()

	resp, err := http.Get(parseUrl.String())
	if err != nil {
		errorhandling.HandleError(err)
		return ""
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		errorhandling.HandleError(err)
		return ""
	}

	var mp = map[string]string{}
	json.Unmarshal(data, &mp)

	return mp["key"]

}

func ListProfileOnPlatform(platform string) []string {
	apiUrl := getUrl()
	if apiUrl == "" {
		errorhandling.HandleError(errors.New("no url set in db"))
		return nil
	}

	fullUrl := fmt.Sprintf("http://%s/listProfile", apiUrl)

	parseUrl, err := url.Parse(fullUrl)
	if err != nil {
		errorhandling.HandleError(err)
		return nil
	}
	q := parseUrl.Query()
	q.Add("platform", platform)

	parseUrl.RawQuery = q.Encode()

	apiUrl = parseUrl.String()

	resp, err := http.Get(apiUrl)
	if err != nil {
		errorhandling.HandleError(err)
		return nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		errorhandling.HandleError(err)
		return nil
	}

	var respUnmarshallJson = map[string][]string{}
	err = json.Unmarshal(body, &respUnmarshallJson)
	if err != nil {
		errorhandling.HandleError(err)
		return nil
	}

	return respUnmarshallJson["data"]

}

func ListPlatforms() []string {
	apiUrl := getUrl()
	if apiUrl == "" {
		errorhandling.HandleError(errors.New("no url set in db"))
		return nil
	}

	fullUrl := fmt.Sprintf("http://%s/listPlatform", apiUrl)

	resp, err := http.Get(fullUrl)
	if err != nil {
		errorhandling.HandleError(err)
		return nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		errorhandling.HandleError(err)
		return nil
	}

	var respUnmarshallJson = map[string][]string{}
	err = json.Unmarshal(body, &respUnmarshallJson)
	if err != nil {
		errorhandling.HandleError(err)
		return nil
	}

	return respUnmarshallJson["data"]
}

func CreateVideoProject() {
	//create payload
	payload := createPayloadForVideoProject()
	if payload == nil {
		return
	}

	url := getUrl()
	if url == "" {
		errorhandling.HandleError(errors.New("no url set"))
		return
	}

	fullUrl := "http://" + url + "/createVideo"

	reader := bytes.NewReader(payload)

	req, err := http.NewRequest("POST", fullUrl, reader)
	if err != nil {
		errorhandling.HandleError(err)
		return
	}

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		errorhandling.HandleError(err)
	}

	fmt.Println("writing to the file")
	val := model.QueryDB(model.OUTPUTFOLDER)
	if val == nil {
		errorhandling.HandleError(errors.New("no output folder  in DB"))
		return
	}

	outputFolder, ok := val.(string)
	if !ok {
		errorhandling.HandleError(errors.New("output folder  in db not a string"))
		return
	}

	val = model.QueryDB(model.PROJECTNAME)
	if val == nil {
		errorhandling.HandleError(errors.New("project name not set"))
		return
	}
	projectName, ok := val.(string)
	if !ok {
		errorhandling.HandleError(errors.New("project name not string"))
		return
	}

	filename := filepath.Join(outputFolder, fmt.Sprintf("%s.xml", projectName))

	file, err := os.Create(filename)
	if err != nil {
		errorhandling.HandleError(err)
		return
	}

	file.ReadFrom(response.Body)
}

func createPayloadForVideoProject() []byte {
	val := model.QueryDB(model.PROJECTNAME)
	if val == nil {
		errorhandling.HandleError(errors.New("project name not set"))
		return nil
	}
	projectName, ok := val.(string)
	if !ok {
		errorhandling.HandleError(errors.New("project name not string"))
	}

	val = model.QueryDB(model.IMAGEFOLDER)
	if val == nil {
		errorhandling.HandleError(errors.New("image folder name not set"))
		return nil
	}

	imageFolder, ok := val.(string)
	if !ok {
		errorhandling.HandleError(errors.New("image folder not string"))
		return nil

	}

	val = model.QueryDB(model.REUSEAUDIOFOLDER)
	if val == nil {
		errorhandling.HandleError(errors.New("no reuse audio found"))
		return nil
	}

	reuseAudio, ok := val.(string)
	if !ok {
		errorhandling.HandleError(errors.New("reuse audio folder not a string"))
		return nil
	}

	//read all images
	var images []string
	dirs, err := os.ReadDir(imageFolder)
	if err != nil {
		errorhandling.HandleError(fmt.Errorf("images folder error %w", err))
		return nil
	}

	for _, dir := range dirs {
		if !isImageExt(getExtension(dir.Name())) {
			continue
		}
		fullFilePath := filepath.Join(imageFolder, dir.Name())
		images = append(images, fullFilePath)
	}

	//get all audio folder sorted

	var audioFiles []string
	dirs, err = os.ReadDir(reuseAudio)
	if err != nil {
		errorhandling.HandleError(fmt.Errorf("images folder error %w", err))
		return nil
	}

	var audioFileInInt [][]int
	for _, dir := range dirs {
		fileName := dir.Name()
		if !strings.Contains(fileName, ".wav") {
			continue
		}
		file1 := strings.Split(fileName, ".")

		num1, err := strconv.ParseInt(file1[0], 10, 32)
		if err != nil {
			errorhandling.HandleError(fmt.Errorf("audio file parse err %w", err))
			return nil
		}
		num2, err := strconv.ParseInt(file1[1], 10, 32)
		if err != nil {
			errorhandling.HandleError(fmt.Errorf("audio file parse err %w", err))
			return nil
		}

		audioFileInInt = append(audioFileInInt, []int{int(num1), int(num2)})
	}

	sort.Slice(audioFileInInt, func(i, j int) bool {
		if audioFileInInt[i][0] < audioFileInInt[j][0] {
			return true
		} else if audioFileInInt[i][0] > audioFileInInt[j][0] {
			return false
		} else {
			return audioFileInInt[i][1] < audioFileInInt[j][1]
		}
	})

	for _, fileInd := range audioFileInInt {
		fileName := fmt.Sprintf("%d.%d.wav", fileInd[0], fileInd[1])
		fullFileName := filepath.Join(reuseAudio, fileName)

		audioFiles = append(audioFiles, fullFileName)
	}

	var audioTimings = make([]float64, len(audioFiles))

	for i, filename := range audioFiles {
		duration, err := getWavDuration(filename)
		if err != nil {
			errorhandling.HandleError(err)
			return nil
		}
		audioTimings[i] = duration
	}

	var mp = map[string]any{}
	mp[model.JSONPROJECTNAME] = projectName
	mp[model.JSONIMAGES] = images
	mp[model.JSONAUDIOS] = audioTimings

	payload, err := json.Marshal(mp)
	if err != nil {
		errorhandling.HandleError(err)
		return nil
	}

	return payload

}

func OptimizeScript(script string) string {

	url := getUrl()
	if url == "" {
		errorhandling.HandleError(errors.New("no url from db"))
		return ""
	}

	fullUrl := fmt.Sprintf("http://%s/optimizeScript", url)

	req, err := http.NewRequest("POST", fullUrl, strings.NewReader(script))
	if err != nil {
		errorhandling.HandleError(err)
		return ""
	}

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		errorhandling.HandleError(err)
		return ""
	}

	data, err := io.ReadAll(response.Body)
	if err != nil {
		errorhandling.HandleError(err)
		return ""
	}

	return string(data)

}
