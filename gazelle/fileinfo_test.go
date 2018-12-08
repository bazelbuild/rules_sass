/* Copyright 2018 The Bazel Authors. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package gazelle

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestProtoRegexpGroupNames(t *testing.T) {
	names := sassRe.SubexpNames()
	nameMap := map[string]int{
		"import": importSubexpIndex,
	}
	for name, index := range nameMap {
		if names[index] != name {
			t.Errorf("sass regexp subexp %d is %s ; want %s", index, names[index], name)
		}
	}
	if len(names)-1 != len(nameMap) {
		t.Errorf("sass regexp has %d groups ; want %d", len(names), len(nameMap))
	}
}

func TestProtoFileInfo(t *testing.T) {
	for _, tc := range []struct {
		desc, name, sass string
		want             FileInfo
	}{
		{
			desc: "empty",
			name: "empty^file.sass",
			sass: "",
			want: FileInfo{},
		}, {
			desc: "import single quote",
			name: "single.sass",
			sass: `import 'single';`,
			want: FileInfo{
				Imports: []string{"single"},
			},
		}, {
			desc: "import double quote",
			name: "double.sass",
			sass: `import "double";`,
			want: FileInfo{
				Imports: []string{"double"},
			},
		}, {
			desc: "import two",
			name: "two.sass",
			sass: `import "first";
import "second";`,
			want: FileInfo{
				Imports: []string{"first", "second"},
			},
		}, {
			desc: "import depth",
			name: "deep.sass",
			sass: `import "from/internal/package";`,
			want: FileInfo{
				Imports: []string{"from/internal/package"},
			},
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			dir, err := ioutil.TempDir(os.Getenv("TEST_TEMPDIR"), "TestProtoFileinfo")
			if err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(dir)
			if err := ioutil.WriteFile(filepath.Join(dir, tc.name), []byte(tc.sass), 0600); err != nil {
				t.Fatal(err)
			}

			got := sassFileInfo(dir, tc.name)

			// Reexpose the fields we care bout for testing.
			got = FileInfo{
				Imports: got.Imports,
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("Inequalith.\ngot  %#v;\nwant %#v", got, tc.want)
			}
		})
	}
}
