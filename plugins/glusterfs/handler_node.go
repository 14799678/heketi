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

package glusterfs

import (
	"errors"
	"github.com/heketi/heketi/requests"
)

func (m *GlusterFSPlugin) NodeAddDevice(id string, req *requests.DeviceAddRequest) error {

	m.rwlock.Lock()
	defer m.rwlock.Unlock()

	if node, ok := m.db.nodes[id]; ok {

		for device := range req.Devices {
			err := node.DeviceAdd(&req.Devices[device])
			if err != nil {
				return err
			}
		}

	} else {
		return errors.New("Node not found")
	}

	// Create a new ring
	err := m.ring.CreateRing()
	if err != nil {
		return nil
	}

	// Save db to persistent storage
	m.db.Commit()

	return nil
}

func (m *GlusterFSPlugin) NodeAdd(v *requests.NodeAddRequest) (*requests.NodeInfoResp, error) {

	node := NewNodeDB(v)

	m.rwlock.Lock()
	defer m.rwlock.Unlock()

	// Save to the db
	m.db.nodes[node.Info.Id] = node

	// Save db to persistent storage
	m.db.Commit()

	return &node.Info, nil
}

func (m *GlusterFSPlugin) NodeList() (*requests.NodeListResponse, error) {

	m.rwlock.RLock()
	defer m.rwlock.RUnlock()

	list := &requests.NodeListResponse{}
	list.Nodes = make([]requests.NodeInfoResp, 0)

	for _, info := range m.db.nodes {
		list.Nodes = append(list.Nodes, info.Info)
	}

	return list, nil
}

func (m *GlusterFSPlugin) NodeRemove(id string) error {

	m.rwlock.Lock()
	defer m.rwlock.Unlock()

	// :TODO: What happens when we remove a node that has
	// brick in use?

	if _, ok := m.db.nodes[id]; ok {
		delete(m.db.nodes, id)
	} else {
		return errors.New("Id not found")
	}

	// Create a new ring
	m.ring.CreateRing()

	// Save db to persistent storage
	m.db.Commit()
	return nil

}

func (m *GlusterFSPlugin) NodeInfo(id string) (*requests.NodeInfoResp, error) {

	m.rwlock.RLock()
	defer m.rwlock.RUnlock()

	if node, ok := m.db.nodes[id]; ok {
		return &node.Info, nil
	} else {
		return nil, errors.New("Id not found")
	}

}
