package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"text/template"

	"github.com/gabstv/ecs/rx"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "name, n",
			Usage: "Template: {{.Name}}",
		},
		cli.StringFlag{
			Name:  "package, p",
			Usage: "Package: {{.Package}}",
		},
		cli.StringSliceFlag{
			Name:  "template, t",
			Usage: "Template file(s) to use",
		},
		cli.BoolFlag{
			Name: "async",
		},
		cli.BoolFlag{
			Name: "skip-register",
		},
		cli.StringSliceFlag{
			Name: "vars",
		},
		cli.StringFlag{
			Name: "output, o",
		},
		cli.StringSliceFlag{
			Name: "components",
		},
		cli.StringFlag{
			Name:  "split",
			Value: ";",
		},
		cli.BoolFlag{
			Name: "system-tpl",
		},
		cli.BoolFlag{
			Name: "component-tpl",
		},
	}
	app.Action = run
	if err := app.Run(os.Args); err != nil {
		println(err.Error())
		os.Exit(1)
	}
}

func run(c *cli.Context) error {
	name := c.String("name")
	packagen := c.String("package")
	templatep := c.StringSlice("template")
	async := c.Bool("async")
	rawvars := c.StringSlice("vars")
	rawviewitems := c.StringSlice("components")
	vars := make(map[string]string)
	viewitems := make([]map[string]string, 0)
	for _, v := range rawvars {
		vs := strings.SplitN(v, "=", 2)
		if len(vs) == 2 {
			vars[vs[0]] = vs[1]
			println(vs[0])
			println(vs[1])
		} else {
			vars[vs[0]] = "1"
		}
	}
	for _, v := range rawviewitems {
		vsplit := strings.Split(v, c.String("split"))
		item := make(map[string]string)
		item["Name"] = vsplit[0]
		if len(vsplit) > 1 {
			if vsplit[1] != "" {
				item["Type"] = vsplit[1]
			} else {
				item["Type"] = "*" + vsplit[0]
			}
		} else {
			item["Type"] = "*" + vsplit[0]
		}
		if len(vsplit) > 2 {
			item["Getter"] = vsplit[2]
		} else {
			item["Getter"] = fmt.Sprintf("Get%sComponent(v.world).Data(e)", vsplit[0])
		}
		viewitems = append(viewitems, item)
	}
	var tpl *template.Template
	if len(templatep) > 0 {
		var err error
		tpl, err = template.ParseFiles(templatep...)
		if err != nil {
			return err
		}
	} else {
		if c.Bool("system-tpl") {
			f, err := rx.FS().Open("templates/system.tmpl")
			if err != nil {
				return err
			}
			d, _ := ioutil.ReadAll(f)
			f.Close()
			tpl, err = template.New("").Parse(string(d))
			if err != nil {
				return err
			}
		} else if c.Bool("component-tpl") {
			f, err := rx.FS().Open("templates/component.tmpl")
			if err != nil {
				return err
			}
			d, _ := ioutil.ReadAll(f)
			f.Close()
			tpl, err = template.New("").Parse(string(d))
			if err != nil {
				return err
			}
		}
	}
	tpld := struct {
		Package      string
		Name         string
		Async        bool
		Vars         map[string]string
		SkipRegister bool
		ViewItems    []map[string]string
	}{
		Package:      packagen,
		Name:         name,
		Async:        async,
		Vars:         vars,
		SkipRegister: c.Bool("skip-register"),
		ViewItems:    viewitems,
	}
	f, err := os.Create(c.String("output"))
	if err != nil {
		return err
	}
	defer f.Close()
	if err := tpl.Execute(f, tpld); err != nil {
		return err
	}
	return nil
}
