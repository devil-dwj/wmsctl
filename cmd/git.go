package cmd

import (
	"context"
	"fmt"
	stdurl "net/url"
	"os"
	"os/exec"
	"path"
	"strings"
)

type Git struct {
	templateUrl string // https://github.com/devil-dwj/wms-template.git
	home        string // "home\\user/.wmsctl/repo/github.com/devil-dwj"
	homePath    string // "home\\user/.wmsctl/repo/github.com/devil-dwj/wms-template@main"
	path        string // "github/devil_dwj/hello"
}

func NewGit(templateUrl string, path string) *Git {
	home := withHomeDir("repo/" + repoDir(templateUrl))
	return &Git{
		templateUrl: templateUrl,
		home:        home,
		homePath:    gpath(home, templateUrl),
		path:        path,
	}
}

func (r *Git) Clone() error {
	if _, err := os.Stat(r.homePath); !os.IsNotExist(err) {
		fmt.Printf("git pull %s ...\n", r.templateUrl)
		return r.Pull()
	}

	fmt.Printf("git clone %s ...\n", r.templateUrl)
	cmd := exec.CommandContext(context.Background(), "git", "clone", r.templateUrl, r.homePath)
	out, err := cmd.CombinedOutput()
	fmt.Println(string(out))
	if err != nil {
		return err
	}
	return nil
}

func (r *Git) Pull() error {
	ctx := context.Background()
	cmd := exec.CommandContext(ctx, "git", "symbolic-ref", "HEAD")
	cmd.Dir = r.homePath
	_, err := cmd.CombinedOutput()
	if err != nil {
		return nil
	}
	cmd = exec.CommandContext(ctx, "git", "pull")
	cmd.Dir = r.homePath
	out, err := cmd.CombinedOutput()
	fmt.Println(string(out))
	if err != nil {
		return err
	}
	return err
}

func (p *Git) CopyToWork(to string) error {
	fmt.Println("copy to work:", to)
	mod, err := modulePath(path.Join(p.homePath, "go.mod"))
	if err != nil {
		return err
	}
	err = copyDir(p.homePath, to, []string{mod, p.path}, []string{".git", ".github"})
	return err
}

func repoDir(url string) string {
	if !strings.Contains(url, "//") {
		url = "//" + url
	}
	if strings.HasPrefix(url, "//git@") {
		url = "ssh:" + url
	} else if strings.HasPrefix(url, "//") {
		url = "https:" + url
	}
	u, err := stdurl.Parse(url)
	if err == nil {
		url = fmt.Sprintf("%s://%s%s", u.Scheme, u.Hostname(), u.Path)
	}
	var start int
	start = strings.Index(url, "//")
	if start == -1 {
		start = strings.Index(url, ":") + 1
	} else {
		start += 2
	}
	end := strings.LastIndex(url, "/")
	return url[start:end]
}
