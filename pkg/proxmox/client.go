package proxmox

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"main/pkg/types"
	"net/http"

	"github.com/rs/zerolog"
)

type Client struct {
	Config types.ProxmoxConfig
	Logger zerolog.Logger
}

func NewClient(config types.ProxmoxConfig, logger *zerolog.Logger) *Client {
	return &Client{
		Config: config,
		Logger: logger.With().
			Str("component", "proxmox_client").
			Str("host", config.URL).
			Logger(),
	}
}

func (c *Client) RelativeLink(url string) string {
	return fmt.Sprintf("%s%s", c.Config.URL, url)
}

func (c *Client) GetResources() (*types.ProxmoxStatusResponse, error) {
	var response *types.ProxmoxStatusResponse
	url := c.RelativeLink("/api2/json/cluster/resources")
	err := c.QueryAndDecode(url, &response)

	return response, err
}

func (c *Client) GetNodesWithAssets() ([]types.NodeWithAssets, error) {
	resources, err := c.GetResources()
	if err != nil {
		return []types.NodeWithAssets{}, err
	}

	nodes, err := c.ParseNodesWithAssetsFromResponse(resources)
	if err != nil {
		return []types.NodeWithAssets{}, err
	}

	return nodes, nil
}

func (c *Client) StartContainer(container types.Container) (*types.ProxmoxActionResponse, error) {
	return c.ContainerAction(container, "start")
}

func (c *Client) StopContainer(container types.Container) (*types.ProxmoxActionResponse, error) {
	return c.ContainerAction(container, "stop")
}

func (c *Client) RebootContainer(container types.Container) (*types.ProxmoxActionResponse, error) {
	return c.ContainerAction(container, "reboot")
}

func (c *Client) GetLxcContainerConfig(container types.Container) (*types.ProxmoxLxcConfigResponse, error) {
	var response *types.ProxmoxLxcConfigResponse
	url := c.RelativeLink(fmt.Sprintf(
		"/api2/extjs/nodes/%s/lxc/%d/config",
		container.Node,
		container.VMID,
	))
	err := c.QueryAndDecode(url, &response)

	return response, err
}

func (c *Client) GetQemuContainerConfig(container types.Container) (*types.ProxmoxQemuConfigResponse, error) {
	var response *types.ProxmoxQemuConfigResponse
	url := c.RelativeLink(fmt.Sprintf(
		"/api2/extjs/nodes/%s/qemu/%d/config",
		container.Node,
		container.VMID,
	))
	err := c.QueryAndDecode(url, &response)

	return response, err
}

func (c *Client) ContainerAction(container types.Container, action string) (*types.ProxmoxActionResponse, error) {
	var response *types.ProxmoxActionResponse
	url := c.RelativeLink(fmt.Sprintf(
		"/api2/extjs/nodes/%s/%s/status/%s",
		container.Node,
		container.ID,
		action,
	))
	err := c.QueryAndDecodePost(url, nil, &response)
	if err != nil {
		return nil, err
	}

	if response.Success != 1 {
		return nil, fmt.Errorf(response.Message)
	}

	return response, err
}

/* Query functions */

func (c *Client) Query(url string) (io.ReadCloser, error) {
	return c.DoQuery("GET", url, nil)
}

func (c *Client) QueryPost(url string, body interface{}) (io.ReadCloser, error) {
	return c.DoQuery("POST", url, body)
}

func (c *Client) QueryDelete(url string) error {
	_, err := c.DoQuery("DELETE", url, nil)
	return err
}

func (c *Client) QueryAndDecode(url string, output interface{}) error {
	body, err := c.Query(url)
	if err != nil {
		return err
	}

	defer body.Close()
	return json.NewDecoder(body).Decode(&output)
}

func (c *Client) QueryAndDecodePost(url string, postBody interface{}, output interface{}) error {
	body, err := c.QueryPost(url, postBody)
	if err != nil {
		return err
	}

	defer body.Close()

	return json.NewDecoder(body).Decode(&output)
}

func (c *Client) DoQuery(method string, url string, body interface{}) (io.ReadCloser, error) {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	var req *http.Request
	var err error

	if body != nil {
		buffer := new(bytes.Buffer)

		if err := json.NewEncoder(buffer).Encode(body); err != nil {
			return nil, err
		}

		req, err = http.NewRequest(method, url, buffer)
	} else {
		req, err = http.NewRequest(method, url, nil)
	}

	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("PVEAPIToken=%s=%s", c.Config.User, c.Config.Token))

	c.Logger.Trace().
		Str("url", url).
		Str("method", method).
		Msg("Doing a Proxmox API query")

	resp, err := client.Do(req)
	if err != nil {
		c.Logger.Error().
			Str("url", url).
			Str("method", method).
			Err(err).
			Msg("Error querying Proxmox")
		return nil, err
	}

	if resp.StatusCode >= http.StatusBadRequest {
		c.Logger.Error().
			Str("url", url).
			Str("method", method).
			Int("status", resp.StatusCode).
			Msg("Got error code from Proxmox")
		return nil, fmt.Errorf("Could not fetch request. Status code: %d", resp.StatusCode)
	}

	return resp.Body, nil
}
