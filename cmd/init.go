package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/spf13/cobra"
)

const (
	wmsTemplateUrl = "https://github.com/devil-dwj/wms-template.git"
)

var Init = &cobra.Command{
	Use:   "init",
	Short: "Init project",
	Long:  "Init project using the repository template, eg: wmsctl init helloworld",
	Run:   run,
}

func run(cmd *cobra.Command, args []string) {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	if len(args) == 0 {
		fmt.Printf("args is nil")
		return
	}
	name := args[0]
	p := &Project{
		Name: path.Base(name),
		Path: name,
	}
	err = p.Init(wd)
	if err != nil {
		fmt.Printf("%+v", err)
	}
}

type Project struct {
	Name string // helloworld
	Path string // github.com/devil-dwj/helloworld
}

func (p *Project) Init(wd string) error {
	to := path.Join(wd, p.Name)

	fmt.Println("wmsctl creating project: ", p.Name)
	git := NewGit(wmsTemplateUrl, p.Path)
	err := git.Clone()
	if err != nil {
		return err
	}

	err = git.CopyToWork(to)
	if err != nil {
		return err
	}
	tree(to, wd)

	fmt.Println("Project created successfully")
	return nil
}
