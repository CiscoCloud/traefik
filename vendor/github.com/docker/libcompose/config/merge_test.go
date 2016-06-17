package config

import "testing"

type NullLookup struct {
}

func (n *NullLookup) Lookup(file, relativeTo string) ([]byte, string, error) {
	return nil, "", nil
}

func (n *NullLookup) ResolvePath(path, inFile string) string {
	return ""
}

func TestExtendsInheritImage(t *testing.T) {
	configV1, _, _, err := Merge(NewServiceConfigs(), nil, &NullLookup{}, "", []byte(`
parent:
  image: foo
child:
  extends:
    service: parent
`))
	if err != nil {
		t.Fatal(err)
	}

	configV2, _, _, err := Merge(NewServiceConfigs(), nil, &NullLookup{}, "", []byte(`
version: '2'
services:
  parent:
    image: foo
  child:
    extends:
      service: parent
`))
	if err != nil {
		t.Fatal(err)
	}

	for _, config := range []map[string]*ServiceConfig{configV1, configV2} {
		parent := config["parent"]
		child := config["child"]

		if parent.Image != "foo" {
			t.Fatal("Invalid image", parent.Image)
		}

		if child.Build.Context != "" {
			t.Fatal("Invalid build", child.Build)
		}

		if child.Image != "foo" {
			t.Fatal("Invalid image", child.Image)
		}
	}
}

func TestExtendsInheritBuild(t *testing.T) {
	configV1, _, _, err := Merge(NewServiceConfigs(), nil, &NullLookup{}, "", []byte(`
parent:
  build: .
child:
  extends:
    service: parent
`))
	if err != nil {
		t.Fatal(err)
	}

	configV2, _, _, err := Merge(NewServiceConfigs(), nil, &NullLookup{}, "", []byte(`
version: '2'
services:
  parent:
    build:
      context: .
  child:
    extends:
      service: parent
`))
	if err != nil {
		t.Fatal(err)
	}

	for _, config := range []map[string]*ServiceConfig{configV1, configV2} {
		parent := config["parent"]
		child := config["child"]

		if parent.Build.Context != "." {
			t.Fatal("Invalid build", parent.Build)
		}

		if child.Build.Context != "." {
			t.Fatal("Invalid build", child.Build)
		}

		if child.Image != "" {
			t.Fatal("Invalid image", child.Image)
		}
	}
}

func TestExtendBuildOverImage(t *testing.T) {
	configV1, _, _, err := Merge(NewServiceConfigs(), nil, &NullLookup{}, "", []byte(`
parent:
  image: foo
child:
  build: .
  extends:
    service: parent
`))
	if err != nil {
		t.Fatal(err)
	}

	configV2, _, _, err := Merge(NewServiceConfigs(), nil, &NullLookup{}, "", []byte(`
version: '2'
services:
  parent:
    image: foo
  child:
    build:
      context: .
    extends:
      service: parent
`))
	if err != nil {
		t.Fatal(err)
	}

	for _, config := range []map[string]*ServiceConfig{configV1, configV2} {
		parent := config["parent"]
		child := config["child"]

		if parent.Image != "foo" {
			t.Fatal("Invalid image", parent.Image)
		}

		if child.Build.Context != "." {
			t.Fatal("Invalid build", child.Build)
		}

		if child.Image != "" {
			t.Fatal("Invalid image", child.Image)
		}
	}
}

func TestExtendImageOverBuild(t *testing.T) {
	configV1, _, _, err := Merge(NewServiceConfigs(), nil, &NullLookup{}, "", []byte(`
parent:
  build: .
child:
  image: foo
  extends:
    service: parent
`))
	if err != nil {
		t.Fatal(err)
	}

	configV2, _, _, err := Merge(NewServiceConfigs(), nil, &NullLookup{}, "", []byte(`
version: '2'
services:
  parent:
    build:
      context: .
  child:
    image: foo
    extends:
      service: parent
`))
	if err != nil {
		t.Fatal(err)
	}

	for _, config := range []map[string]*ServiceConfig{configV1, configV2} {
		parent := config["parent"]
		child := config["child"]

		if parent.Image != "" {
			t.Fatal("Invalid image", parent.Image)
		}

		if parent.Build.Context != "." {
			t.Fatal("Invalid build", parent.Build)
		}

		if child.Build.Context != "" {
			t.Fatal("Invalid build", child.Build)
		}

		if child.Image != "foo" {
			t.Fatal("Invalid image", child.Image)
		}
	}
}

func TestRestartNo(t *testing.T) {
	configV1, _, _, err := Merge(NewServiceConfigs(), nil, &NullLookup{}, "", []byte(`
test:
  restart: "no"
  image: foo
`))
	if err != nil {
		t.Fatal(err)
	}

	configV2, _, _, err := Merge(NewServiceConfigs(), nil, &NullLookup{}, "", []byte(`
version: '2'
services:
  test:
    restart: "no"
    image: foo
`))
	if err != nil {
		t.Fatal(err)
	}

	for _, config := range []map[string]*ServiceConfig{configV1, configV2} {
		test := config["test"]

		if test.Restart != "no" {
			t.Fatal("Invalid restart policy", test.Restart)
		}
	}
}

func TestRestartAlways(t *testing.T) {
	configV1, _, _, err := Merge(NewServiceConfigs(), nil, &NullLookup{}, "", []byte(`
test:
  restart: always
  image: foo
`))
	if err != nil {
		t.Fatal(err)
	}

	configV2, _, _, err := Merge(NewServiceConfigs(), nil, &NullLookup{}, "", []byte(`
version: '2'
services:
  test:
    restart: always
    image: foo
`))
	if err != nil {
		t.Fatal(err)
	}

	for _, config := range []map[string]*ServiceConfig{configV1, configV2} {
		test := config["test"]

		if test.Restart != "always" {
			t.Fatal("Invalid restart policy", test.Restart)
		}
	}
}

func TestIsValidRemote(t *testing.T) {
	gitUrls := []string{
		"git://github.com/docker/docker",
		"git@github.com:docker/docker.git",
		"git@bitbucket.org:atlassianlabs/atlassian-docker.git",
		"https://github.com/docker/docker.git",
		"http://github.com/docker/docker.git",
		"http://github.com/docker/docker.git#branch",
		"http://github.com/docker/docker.git#:dir",
	}
	incompleteGitUrls := []string{
		"github.com/docker/docker",
	}
	invalidGitUrls := []string{
		"http://github.com/docker/docker.git:#branch",
	}
	for _, url := range gitUrls {
		if !IsValidRemote(url) {
			t.Fatalf("%q should have been a valid remote", url)
		}
	}
	for _, url := range incompleteGitUrls {
		if !IsValidRemote(url) {
			t.Fatalf("%q should have been a valid remote", url)
		}
	}
	for _, url := range invalidGitUrls {
		if !IsValidRemote(url) {
			t.Fatalf("%q should have been a valid remote", url)
		}
	}
}
