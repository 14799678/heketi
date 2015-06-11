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
	"github.com/lpabon/heketi/requests"
	"github.com/lpabon/heketi/tests"
	"os"
	"testing"
)

func TestGlusterFSDBFileLoad(t *testing.T) {
	tmpfile := tests.Tempfile()
	defer os.Remove(tmpfile)

	db := NewGlusterFSDB(tmpfile)

	db.nodes["one"] = &NodeDB{
		Node: &requests.NodeInfoResp{
			Name: "nodetest",
		},
	}
	db.volumes["a"] = &VolumeDB{
		Volume: &requests.VolumeInfoResp{
			Name: "volumetest",
		},
	}

	err := db.Commit()
	tests.Assert(t, err == nil)

	newdb := NewGlusterFSDB(tmpfile)
	tests.Assert(t, newdb != nil)

	tests.Assert(t, newdb.nodes["one"].Node.Name == db.nodes["one"].Node.Name)
	tests.Assert(t, newdb.volumes["a"].Volume.Name == db.volumes["a"].Volume.Name)
}
