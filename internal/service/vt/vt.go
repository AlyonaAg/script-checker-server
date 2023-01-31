package virustotal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

const (
	finishedStatus = "Scan finished, information embedded"
)

type VirusTotal struct {
	apikey string
}

type VirusTotalResponse struct {
	ResponseCode int    `json:"response_code"`
	Message      string `json:"verbose_msg"`
}

type ScanResponse struct {
	VirusTotalResponse

	ScanId    string `json:"scan_id"`
	Sha1      string `json:"sha1"`
	Resource  string `json:"resource"`
	Sha256    string `json:"sha256"`
	Permalink string `json:"permalink"`
	Md5       string `json:"md5"`
}

type FileScan struct {
	Detected bool   `json:"detected"`
	Version  string `json:"version"`
	Result   string `json:"result"`
	Update   string `json:"update"`
}

type ReportResponse struct {
	VirusTotalResponse
	Resource  string              `json:"resource"`
	ScanId    string              `json:"scan_id"`
	Sha1      string              `json:"sha1"`
	Sha256    string              `json:"sha256"`
	Md5       string              `json:"md5"`
	Scandate  string              `json:"scan_date"`
	Positives int                 `json:"positives"`
	Total     int                 `json:"total"`
	Permalink string              `json:"permalink"`
	Scans     map[string]FileScan `json:"scans"`
}

func (sr *ScanResponse) String() string {
	return fmt.Sprintf("scanid: %s, resource: %s, permalink: %s, md5: %s", sr.ScanId, sr.Resource, sr.Permalink, sr.Md5)
}

func NewVirusTotal(apikey string) (*VirusTotal, error) {
	vt := &VirusTotal{apikey: apikey}
	return vt, nil
}

func (vt *VirusTotal) GetDengerPercent(path string) (float64, string, error) {
	var percent float64
	var vtDanger string

	file, err := os.Open(path)
	if err != nil {
		return percent, vtDanger, err
	}
	defer file.Close()

	scanResp, err := vt.Scan(path, file)
	if err != nil {
		return percent, vtDanger, err
	}

	var status string
	var reposrtResp *ReportResponse

	for status != finishedStatus {
		fmt.Printf("[%s] status: %s\n", path, status)

		reposrtResp, err = vt.Report(scanResp.Md5)
		if err != nil {
			return percent, vtDanger, err
		}
		status = reposrtResp.Message
		time.Sleep(5 * time.Second)
	}
	fmt.Printf("[%s] path totalvirus result: %d/%d\n", path, reposrtResp.Positives, reposrtResp.Total)

	return float64(reposrtResp.Positives) / float64(reposrtResp.Total),
		fmt.Sprintf("%d/%d", reposrtResp.Positives, reposrtResp.Total),
		nil
}

func (vt *VirusTotal) Report(resource string) (*ReportResponse, error) {
	u, err := url.Parse("https://www.virustotal.com/vtapi/v2/file/report")

	params := url.Values{"apikey": {vt.apikey}, "resource": {resource}}

	resp, err := http.PostForm(u.String(), params)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if len(contents) == 0 {
		return &ReportResponse{}, nil
	}

	var reportResponse = &ReportResponse{}

	err = json.Unmarshal(contents, &reportResponse)

	return reportResponse, err
}

func (vt *VirusTotal) Scan(path string, file io.Reader) (*ScanResponse, error) {
	params := map[string]string{
		"apikey": vt.apikey,
	}

	request, err := newfileUploadRequest("http://www.virustotal.com/vtapi/v2/file/scan", params, path, file)

	if err != nil {
		return nil, err
	}

	client := &http.Client{}

	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var scanResponse = &ScanResponse{}
	err = json.Unmarshal(contents, &scanResponse)

	return scanResponse, err
}

// Creates a new file upload http request with optional extra params
func newfileUploadRequest(uri string, params map[string]string, path string, file io.Reader) (*http.Request, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}

	part, err := writer.CreateFormFile("file", filepath.Base(path))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)

	err = writer.Close()

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", uri, body)

	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, err
}
