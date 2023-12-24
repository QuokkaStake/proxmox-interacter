package proxmox

import (
	"encoding/json"
	"fmt"
	"main/pkg/types"
	"main/pkg/utils"
)

func (c *Client) ParseContainersFromResponse(response *types.ProxmoxStatusResponse) ([]types.Container, error) {
	containers := make([]types.Container, 0)

	for _, rawData := range response.Data {
		resourceType, ok := rawData["type"]
		if !ok {
			return containers, fmt.Errorf("resource type is not present")
		}

		if resourceType != "lxc" && resourceType != "qemu" {
			continue
		}

		rawBytes, err := json.Marshal(rawData)
		if err != nil {
			return nil, err
		}

		container := types.Container{}
		if err := json.Unmarshal(rawBytes, &container); err != nil {
			return nil, err
		}

		container.Link = c.Config.GetResourceLink(container)

		containers = append(containers, container)
	}

	return containers, nil
}

func (c *Client) ParseNodesFromResponse(response *types.ProxmoxStatusResponse) ([]types.Node, error) {
	nodes := make([]types.Node, 0)

	for _, rawData := range response.Data {
		resourceType, ok := rawData["type"]
		if !ok {
			return nodes, fmt.Errorf("resource type is not present")
		}

		if resourceType != "node" {
			continue
		}

		rawBytes, err := json.Marshal(rawData)
		if err != nil {
			return nodes, err
		}

		node := types.Node{}
		if err := json.Unmarshal(rawBytes, &node); err != nil {
			return nodes, err
		}

		node.Link = c.Config.GetResourceLink(node)

		nodes = append(nodes, node)
	}

	return nodes, nil
}

func (c *Client) ParseStoragesFromResponse(response *types.ProxmoxStatusResponse) ([]types.Storage, error) {
	storages := make([]types.Storage, 0)

	for _, rawData := range response.Data {
		resourceType, ok := rawData["type"]
		if !ok {
			return storages, fmt.Errorf("resource type is not present")
		}

		if resourceType != "storage" {
			continue
		}

		rawBytes, err := json.Marshal(rawData)
		if err != nil {
			return storages, err
		}

		storage := types.Storage{}
		if err := json.Unmarshal(rawBytes, &storage); err != nil {
			return storages, err
		}

		storage.Link = c.Config.GetResourceLink(storage)

		storages = append(storages, storage)
	}

	return storages, nil
}

func (c *Client) ParseNodesWithAssetsFromResponse(response *types.ProxmoxStatusResponse) ([]types.NodeWithAssets, error) {
	nodes, err := c.ParseNodesFromResponse(response)
	if err != nil {
		return nil, err
	}

	containers, err := c.ParseContainersFromResponse(response)
	if err != nil {
		return nil, err
	}

	storages, err := c.ParseStoragesFromResponse(response)
	if err != nil {
		return nil, err
	}

	result := make([]types.NodeWithAssets, len(nodes))

	for index, node := range nodes {
		result[index] = types.NodeWithAssets{
			Node: node,
			Containers: utils.Filter(containers, func(c types.Container) bool {
				return c.Node == node.Node
			}),
			Storages: utils.Filter(storages, func(c types.Storage) bool {
				return c.Node == node.Node
			}),
		}
	}

	return result, nil
}
