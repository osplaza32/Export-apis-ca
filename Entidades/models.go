package Entidades

import "os/exec"

type TheContent struct {
		Item struct {
			Name      string `json:"Name"`
			Type      string `json:"Type"`
			TimeStamp string `json:"TimeStamp"`
			Link      struct {
				Rel string `json:"-rel"`
				URI string `json:"-uri"`
			} `json:"Link"`
			Resource struct {
				DependencyList struct {
					Reference struct {
						Name         string `json:"Name"`
						ID           string `json:"Id"`
						Type         string `json:"Type"`
						Dependencies struct {
							Dependency []struct {
								ID   string `json:"Id"`
								Type string `json:"Type"`
								Name string `json:"Name"`
							} `json:"Dependency"`
						} `json:"Dependencies"`
					} `json:"Reference"`
					Dependencies struct {
						Dependency []struct {
							Type string `json:"Type"`
							Name string `json:"Name"`
							ID   string `json:"Id"`
						} `json:"Dependency"`
					} `json:"Dependencies"`
					MissingDependencies string `json:"MissingDependencies"`
				} `json:"DependencyList"`
			} `json:"Resource"`
			L7 string `json:"-l7"`
		} `json:"Item"`
	}
	type Policyinfo struct {
		Item struct {
			Name      string `json:"Name"`
			ID        string `json:"Id"`
			Type      string `json:"Type"`
			TimeStamp string `json:"TimeStamp"`
			Link      []struct {
				Rel string `json:"-rel"`
				URI string `json:"-uri,omitempty"`
			} `json:"Link"`
			Resource struct {
				Policy struct {
					GUID         string `json:"-guid"`
					ID           string `json:"-id"`
					Version      string `json:"-version"`
					PolicyDetail struct {
						ID         string `json:"-id"`
						Version    string `json:"-version"`
						Name       string `json:"Name"`
						PolicyType string `json:"PolicyType"`
						Properties struct {
							Property []struct {
								Key          string `json:"-key"`
								LongValue    string `json:"LongValue,omitempty"`
								BooleanValue string `json:"BooleanValue,omitempty"`
							} `json:"Property"`
						} `json:"Properties"`
						FolderID string `json:"-folderId"`
						GUID     string `json:"-guid"`
					} `json:"PolicyDetail"`
					Resources struct {
						ResourceSet struct {
							Resource struct {
								Content string `json:"#content"`
								Type    string `json:"-type"`
							} `json:"Resource"`
							Tag string `json:"-tag"`
						} `json:"ResourceSet"`
					} `json:"Resources"`
				} `json:"Policy"`
			} `json:"Resource"`
			L7 string `json:"-l7"`
		} `json:"Item"`
	}
	type Serviceinfo struct {
		Item struct {
			Link []struct {
				URI string `json:"-uri,omitempty"`
				Rel string `json:"-rel,omitempty"`

			} `json:"Link"`
			Resource struct {
				Service struct {
					ID            string `json:"-id"`
					Version       string `json:"-version"`
					ServiceDetail struct {
						ServiceMappings struct {
							HTTPMapping struct {
								URLPattern string `json:"UrlPattern"`
								Verbs      struct {
									Verb []string `json:"Verb"`
								} `json:"Verbs"`
							} `json:"HttpMapping"`
						} `json:"ServiceMappings"`
						Properties struct {
							Property []struct {
								Key          string `json:"-key"`
								BooleanValue string `json:"BooleanValue,omitempty"`
								LongValue    string `json:"LongValue,omitempty"`
							} `json:"Property"`
						} `json:"Properties"`
						FolderID string `json:"-folderId"`
						ID       string `json:"-id"`
						Version  string `json:"-version"`
						Name     string `json:"Name"`
						Enabled  string `json:"Enabled"`
					} `json:"ServiceDetail"`
					Resources struct {
						ResourceSet struct {
							Resource struct {
								Content string `json:"#content"`
								Type    string `json:"-type"`
								Version string `json:"-version"`
							} `json:"Resource"`
							Tag string `json:"-tag"`
						} `json:"ResourceSet"`
					} `json:"Resources"`
				} `json:"Service"`
			} `json:"Resource"`
			L7        string `json:"-l7"`
			Name      string `json:"Name"`
			ID        string `json:"Id"`
			Type      string `json:"Type"`
			TimeStamp string `json:"TimeStamp"`
		} `json:"Item"`
	}


type Git struct {
	Cmd string
	cmd string
}

func (git *Git) Exists() bool {
	out, _ := exec.Command(git.cmd).Output()
	return len(out) > 0
}

func (git *Git) Init() ([]byte, error) {

	return exec.Command(git.cmd, "config user.email 'entel.bot@api.entel.com' && git config user.name  'entelito' && git init C:\\Users\\Administrador\\Desktop\\APIs_Chile").Output()
}

func (git *Git) GetConfig(file, key string) string {
	value, _ := exec.Command(git.cmd, "config", file, key).Output()
	return string(value)
}

func (git *Git) Add(file string) ([]byte, error) {
	return exec.Command(git.cmd, "add", file).Output()
}

func (git *Git) Commit(message string) ([]byte, error) {
	return exec.Command(git.cmd, "commit", "-m", message).Output()
}

func (git *Git) Clone(url, dir string) ([]byte, error) {
	return exec.Command(git.cmd, "clone", url, dir).Output()
}

func (git *Git) NewGit() *Git {
	return &Git{
		cmd: "git",
		Cmd:"git",
	}
}
