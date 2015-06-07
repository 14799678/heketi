//
// Copyright (c) 2014 The heketi Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package mock

import (
	"errors"
	"github.com/lpabon/heketi/models"
)

type Node struct {
	node *models.NodeInfoResp
}

type MockDB struct {
	nodes      map[uint64]*Node
	current_id uint64
}

type MockNode struct {
	db MockDB
}

func NewMockNode() *MockNode {
	m := &MockNode{}
	m.db.nodes = make(map[uint64]*Node)

	return m
}

// Not concurrent safe
func (m *MockNode) Add(v *models.NodeAddRequest) (*models.NodeInfoResp, error) {
	m.db.current_id++

	info := &models.NodeInfoResp{}
	info.Name = v.Name
	info.Zone = v.Zone
	info.Id = m.db.current_id

	node := &Node{
		node: info,
	}

	m.db.nodes[m.db.current_id] = node

	return m.Info(info.Id)
}

func (m *MockNode) List() (*models.NodeListResponse, error) {

	list := &models.NodeListResponse{}
	list.Nodes = make([]models.NodeInfoResp, 0)

	for id, _ := range m.db.nodes {
		info, err := m.Info(id)
		if err != nil {
			return nil, err
		}
		list.Nodes = append(list.Nodes, *info)
	}

	return list, nil
}

func (m *MockNode) Remove(id uint64) error {

	if _, ok := m.db.nodes[id]; ok {
		delete(m.db.nodes, id)
		return nil
	} else {
		return errors.New("Id not found")
	}

}

func (m *MockNode) Info(id uint64) (*models.NodeInfoResp, error) {

	if node, ok := m.db.nodes[id]; ok {
		info := &models.NodeInfoResp{}
		*info = *node.node
		return info, nil
	} else {
		return nil, errors.New("Id not found")
	}

}
