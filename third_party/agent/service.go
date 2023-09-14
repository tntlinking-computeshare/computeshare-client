package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	transhttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/ipfs/kubo/core"
	"github.com/mohaijiang/computeshare-client/internal/conf"
	"io"
	"net/http"
	"net/url"
	"path"
)

type AgentService struct {
	confData *conf.Data
	ipfsNode *core.IpfsNode
	id       string
}

func NewAgentService(conn *transhttp.Client, ipfsNode *core.IpfsNode) *AgentService {

	//client := pb.New(conn)
	return &AgentService{
		//confData: confData,
		ipfsNode: ipfsNode,
	}
}

func (s *AgentService) Register() error {
	api, err := url.JoinPath(s.confData.GetComputerPowerApi(), "/v1/agent")
	if err != nil {
		return err
	}
	fmt.Println("HTTP JSON POST URL:", api)
	peerId := s.ipfsNode.Identity.String()
	var jsonData = []byte(`{
		"name": "` + peerId + `"
	}`)
	request, err := http.NewRequest("POST", api, bytes.NewBuffer(jsonData))
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	if err != nil {
		return err
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	fmt.Println("response Status:", response.Status)
	fmt.Println("response Headers:", response.Header)
	body, _ := io.ReadAll(response.Body)
	var m map[string]string
	err = json.Unmarshal(body, &m)
	if err != nil {
		return err
	}
	fmt.Println("response Body:", string(body))

	s.id = m["id"]

	return nil
}

func (s *AgentService) UnRegister() error {
	url := path.Join(s.confData.GetComputerPowerApi(), fmt.Sprintf("/v1/agent/%s", s.id))
	fmt.Println("HTTP JSON POST URL:", url)

	request, err := http.NewRequest("DELETE", url, nil)
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	fmt.Println("response Status:", response.Status)
	fmt.Println("response Headers:", response.Header)
	body, _ := io.ReadAll(response.Body)
	fmt.Println("response Body:", string(body))

	return nil
}
