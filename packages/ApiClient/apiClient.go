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
	model.AddToDb(model.APIURL, "192.168.29.2:4000")
}

func AddKey(platform, profile, newKey string) {
	apiUrl := getUrl()
	if apiUrl == "" {
		errorhandling.HandleErrorPop(errors.New("no url set in db"))
		return
	}

	fullUrl := fmt.Sprintf("http://%s/addKey", apiUrl)
	parseUrl, err := url.Parse(fullUrl)
	if err != nil {
		errorhandling.HandleErrorPop(err)
		return
	}
	q := parseUrl.Query()
	q.Add("platform", platform)
	q.Add("profile", profile)
	q.Add("key", newKey)

	parseUrl.RawQuery = q.Encode()

	resp, err := http.Get(parseUrl.String())
	if err != nil {
		errorhandling.HandleErrorPop(err)
		return
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		errorhandling.HandleErrorPop(fmt.Errorf("Add key status code %d, message %s", resp.StatusCode, resp.Status))
	}

}

func AddProfile(platform, newProfile string) {
	apiUrl := getUrl()
	if apiUrl == "" {
		errorhandling.HandleErrorPop(errors.New("no url set in db"))
		return
	}

	fullUrl := fmt.Sprintf("http://%s/addProfile", apiUrl)
	parseUrl, err := url.Parse(fullUrl)
	if err != nil {
		errorhandling.HandleErrorPop(err)
		return
	}
	q := parseUrl.Query()
	q.Add("platform", platform)
	q.Add("profile", newProfile)

	parseUrl.RawQuery = q.Encode()

	resp, err := http.Get(parseUrl.String())
	if err != nil {
		errorhandling.HandleErrorPop(err)
		return
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		errorhandling.HandleErrorPop(fmt.Errorf("add status code %d, message %s", resp.StatusCode, resp.Status))
	}

}

func RotateKey(platform, profile string) {
	apiUrl := getUrl()
	if apiUrl == "" {
		errorhandling.HandleErrorPop(errors.New("no url set in db"))
		return
	}

	fullUrl := fmt.Sprintf("http://%s/rotateKey", apiUrl)
	parseUrl, err := url.Parse(fullUrl)
	if err != nil {
		errorhandling.HandleErrorPop(err)
		return
	}
	q := parseUrl.Query()
	q.Add("platform", platform)
	q.Add("profile", profile)

	parseUrl.RawQuery = q.Encode()

	resp, err := http.Get(parseUrl.String())
	if err != nil {
		errorhandling.HandleErrorPop(err)
		return
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		errorhandling.HandleErrorPop(fmt.Errorf("rotate key status code %d, message %s", resp.StatusCode, resp.Status))
	}

}

func GetKey(platform, profile string) string {
	apiUrl := getUrl()
	if apiUrl == "" {
		errorhandling.HandleErrorPop(errors.New("no url set in db"))
		return ""
	}

	fullUrl := fmt.Sprintf("http://%s/getKey", apiUrl)
	parseUrl, err := url.Parse(fullUrl)
	if err != nil {
		errorhandling.HandleErrorPop(err)
		return ""
	}
	q := parseUrl.Query()
	q.Add("platform", platform)
	q.Add("profile", profile)

	parseUrl.RawQuery = q.Encode()

	resp, err := http.Get(parseUrl.String())
	if err != nil {
		errorhandling.HandleErrorPop(err)
		return ""
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		errorhandling.HandleErrorPop(err)
		return ""
	}

	var mp = map[string]string{}
	json.Unmarshal(data, &mp)

	return mp["key"]

}

func ListProfileOnPlatform(platform string) []string {
	apiUrl := getUrl()
	if apiUrl == "" {
		errorhandling.HandleErrorPop(errors.New("no url set in db"))
		return nil
	}

	fullUrl := fmt.Sprintf("http://%s/listProfile", apiUrl)

	parseUrl, err := url.Parse(fullUrl)
	if err != nil {
		errorhandling.HandleErrorPop(err)
		return nil
	}
	q := parseUrl.Query()
	q.Add("platform", platform)

	parseUrl.RawQuery = q.Encode()

	apiUrl = parseUrl.String()

	resp, err := http.Get(apiUrl)
	if err != nil {
		errorhandling.HandleErrorPop(err)
		return nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		errorhandling.HandleErrorPop(err)
		return nil
	}

	var respUnmarshallJson = map[string][]string{}
	err = json.Unmarshal(body, &respUnmarshallJson)
	if err != nil {
		errorhandling.HandleErrorPop(err)
		return nil
	}
	return respUnmarshallJson["data"]

}

func ListPlatforms() []string {
	apiUrl := getUrl()
	if apiUrl == "" {
		errorhandling.HandleErrorPop(errors.New("no url set in db"))
		return nil
	}

	fullUrl := fmt.Sprintf("http://%s/listPlatform", apiUrl)

	resp, err := http.Get(fullUrl)
	if err != nil {
		errorhandling.HandleErrorPop(err)
		return nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		errorhandling.HandleErrorPop(err)
		return nil
	}

	var respUnmarshallJson = map[string][]string{}
	err = json.Unmarshal(body, &respUnmarshallJson)
	if err != nil {
		errorhandling.HandleErrorPop(err)
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

	fmt.Println(string(payload))

	url := getUrl()
	if url == "" {
		errorhandling.HandleErrorPop(errors.New("no url set"))
		return
	}

	fullUrl := "http://" + url + "/createVideo"

	reader := bytes.NewReader(payload)

	req, err := http.NewRequest("POST", fullUrl, reader)
	if err != nil {
		errorhandling.HandleErrorPop(err)
		return
	}

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		errorhandling.HandleErrorPop(err)
		return
	}

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		body, _ := io.ReadAll(response.Body)
		errorhandling.HandleErrorPop(fmt.Errorf("error with create project status: %d msg :%s", response.StatusCode, string(body)))
	}

	fmt.Println("writing to the file")
	val := model.QueryDB(model.OUTPUTFOLDER)
	if val == nil {
		errorhandling.HandleErrorPop(errors.New("no output folder  in DB"))
		return
	}

	outputFolder, ok := val.(string)
	if !ok {
		errorhandling.HandleErrorPop(errors.New("output folder  in db not a string"))
		return
	}

	val = model.QueryDB(model.PROJECTNAME)
	if val == nil {
		errorhandling.HandleErrorPop(errors.New("project name not set"))
		return
	}
	projectName, ok := val.(string)
	if !ok {
		errorhandling.HandleErrorPop(errors.New("project name not string"))
		return
	}

	filename := filepath.Join(outputFolder, fmt.Sprintf("%s.xml", projectName))

	file, err := os.Create(filename)
	if err != nil {
		errorhandling.HandleErrorPop(err)
		return
	}

	file.ReadFrom(response.Body)
}

func createPayloadForVideoProject() []byte {
	val := model.QueryDB(model.PROJECTNAME)
	if val == nil {
		errorhandling.HandleErrorPop(errors.New("project name not set"))
		return nil
	}
	projectName, ok := val.(string)
	if !ok {
		errorhandling.HandleErrorPop(errors.New("project name not string"))
	}

	val = model.QueryDB(model.SENTENCEGAP)
	sentenceGap := 0.0
	if val != nil {
		sentenceGap, _ = val.(float64)
	}

	val = model.QueryDB(model.PARAGAP)
	paraGap := 0.0
	if val != nil {
		paraGap, _ = val.(float64)
	}

	val = model.QueryDB(model.IMAGEFOLDER)
	if val == nil {
		errorhandling.HandleErrorPop(errors.New("image folder name not set"))
		return nil
	}

	imageFolder, ok := val.(string)
	if !ok {
		errorhandling.HandleErrorPop(errors.New("image folder not string"))
		return nil

	}

	val = model.QueryDB(model.REUSEAUDIOFOLDER)
	if val == nil {
		errorhandling.HandleErrorPop(errors.New("no reuse audio found"))
		return nil
	}

	reuseAudio, ok := val.(string)
	if !ok {
		errorhandling.HandleErrorPop(errors.New("reuse audio folder not a string"))
		return nil
	}

	//read all images
	var images []string
	dirs, err := os.ReadDir(imageFolder)
	if err != nil {
		errorhandling.HandleErrorPop(fmt.Errorf("images folder error %w", err))
		return nil
	}

	for _, dir := range dirs {
		if !isImageExt(getExtension(dir.Name())) {
			continue
		}
		// fullFilePath := filepath.Join(imageFolder, )
		images = append(images, dir.Name())
	}

	val = model.QueryDB(model.IMAGESINORDER)
	if val == nil {
		errorhandling.HandleErrorPop(fmt.Errorf("images in order not set in DB"))
		return nil
	}
	imagesInOrder, _ := val.(bool)

	if imagesInOrder {
		sort.Slice(images, func(i, j int) bool {
			file1, err := strconv.ParseInt(strings.Split(images[i], ".")[0], 10, 64)
			if err != nil {
				errorhandling.HandleError(err)
				return false
			}
			file2, err := strconv.ParseInt(strings.Split(images[j], ".")[0], 10, 64)
			if err != nil {
				errorhandling.HandleError(err)
				return true
			}

			return file1 < file2

		})
	}

	for i, fileName := range images {
		images[i] = filepath.Join(imageFolder, fileName)
	}

	//get all audio folder sorted

	var audioFiles = [][]string{}
	dirs, err = os.ReadDir(reuseAudio)
	if err != nil {
		errorhandling.HandleErrorPop(fmt.Errorf("images folder error %w", err))
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
			errorhandling.HandleErrorPop(fmt.Errorf("audio file parse err %w", err))
			return nil
		}
		num2, err := strconv.ParseInt(file1[1], 10, 32)
		if err != nil {
			errorhandling.HandleErrorPop(fmt.Errorf("audio file parse err %w", err))
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
		for len(audioFiles) < fileInd[0] {
			audioFiles = append(audioFiles, []string{})
		}
		for len(audioFiles[fileInd[0]-1]) < fileInd[1] {
			audioFiles[fileInd[0]-1] = append(audioFiles[fileInd[0]-1], "")
		}

		audioFiles[fileInd[0]-1][fileInd[1]-1] = fullFileName
	}

	var audioTimings = make([][]float64, len(audioFiles))

	for i, filename := range audioFiles {
		duration, err := getWavDurationPara(filename)
		if err != nil {
			errorhandling.HandleErrorPop(err)
			return nil
		}
		audioTimings[i] = duration
	}

	var mp = map[string]any{}
	mp[model.JSONPROJECTNAME] = projectName
	mp[model.JSONIMAGES] = images
	mp[model.JSONAUDIOSTIMES] = audioTimings
	mp[model.JSONAUDIONAMES] = audioFiles
	mp[model.JSONPARAGAP] = paraGap
	mp[model.JSONSENTENCEGAP] = sentenceGap
	mp[model.JSONIMAGESINORDER] = imagesInOrder

	payload, err := json.Marshal(mp)
	if err != nil {
		errorhandling.HandleErrorPop(err)
		return nil
	}

	return payload

}

func getWavDurationPara(sentences []string) ([]float64, error) {
	var ans = []float64{}
	for _, filename := range sentences {
		duration := 10.0
		var err error
		if filename != "" {
			duration, err = getWavDuration(filename)
			if err != nil {
				return nil, err
			}
		}
		ans = append(ans, duration)
	}

	return ans, nil
}

func OptimizeScript(script string, strict16 bool) string {

	url := getUrl()
	if url == "" {
		errorhandling.HandleErrorPop(errors.New("no url from db"))
		return ""
	}

	var bodyMp = map[string]any{}
	bodyMp["script"] = script
	bodyMp["strict16"] = strict16

	body, err := json.Marshal(bodyMp)
	if err != nil {
		errorhandling.HandleErrorPop(err)
		return ""
	}

	fullUrl := fmt.Sprintf("http://%s/optimizeScript", url)

	req, err := http.NewRequest("POST", fullUrl, bytes.NewReader(body))
	if err != nil {
		errorhandling.HandleErrorPop(err)
		return ""
	}

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		errorhandling.HandleErrorPop(err)
		return ""
	}

	data, err := io.ReadAll(response.Body)
	if err != nil {
		errorhandling.HandleErrorPop(err)
		return ""
	}

	return string(data)

}

func GetKeyDetailsForPlatform(platform string) map[string]int {
	apiUrl := getUrl()
	if apiUrl == "" {
		errorhandling.HandleErrorPop(errors.New("no url set in db"))
		return nil
	}

	fullUrl := fmt.Sprintf("http://%s/listProfileKeys", apiUrl)

	parseUrl, err := url.Parse(fullUrl)
	if err != nil {
		errorhandling.HandleErrorPop(err)
		return nil
	}
	q := parseUrl.Query()
	q.Add("platform", platform)

	parseUrl.RawQuery = q.Encode()

	apiFullUrl := parseUrl.String()

	resp, err := http.Get(apiFullUrl)
	if err != nil {
		errorhandling.HandleError(err)
		return nil
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		errorhandling.HandleError(err)
		return nil
	}

	var result = map[string]int{}

	json.Unmarshal(data, &result)

	return result

}

func PresistProfiles(profiles []string) {
	val := model.QueryDB(model.PLATFORMFORKEYSVIEW)
	if val == nil {
		errorhandling.HandleError(fmt.Errorf("platform not set to presist keys"))
		return
	}
	platform, _ := val.(string)

	payload := map[string][]string{
		"profile": profiles,
	}

	data, _ := json.Marshal(payload)

	apiUrl := getUrl()
	if apiUrl == "" {
		errorhandling.HandleErrorPop(errors.New("no url set in db"))
		return
	}

	fullUrl := fmt.Sprintf("http://%s/profileSubmit", apiUrl)

	parseUrl, err := url.Parse(fullUrl)
	if err != nil {
		errorhandling.HandleErrorPop(err)
		return
	}
	q := parseUrl.Query()
	q.Add("platform", platform)

	parseUrl.RawQuery = q.Encode()

	apiFullUrl := parseUrl.String()

	req, err := http.NewRequest("POST", apiFullUrl, bytes.NewReader(data))
	if err != nil {
		errorhandling.HandleError(err)
		return
	}

	_, err = http.DefaultClient.Do(req)
	if err != nil {
		errorhandling.HandleError(err)
		return
	}

}
