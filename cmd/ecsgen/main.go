package main

import (
	"bufio"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strconv"
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/gabstv/ecs/v2/gen"

	"github.com/urfave/cli/v2"
)

func main() {
	app := cli.NewApp()

	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name: "go-package",
			EnvVars: []string{
				"GOPACKAGE",
			},
		},
		&cli.StringFlag{
			Name:    "go-file",
			EnvVars: []string{"GOFILE"},
		},
		&cli.StringFlag{
			Name:    "out",
			EnvVars: []string{"OUTPUT"},
			Aliases: []string{"o", "output"},
		},
	}

	app.Action = run

	app.Run(os.Args)
}

type commentContext struct {
	isECS         bool
	isComponent   bool
	isSystem      bool
	isAsync       bool
	name          string
	uuid          string
	componentsraw []string
	members       []gen.SystemMember
	onSetup       string

	onEntityAdded         string
	onEntityRemoved       string
	onComponentWillResize string
	onComponentResized    string

	onAdd        string
	onRemove     string
	onWillRemove string
	onWillResize string
	onResized    string

	noInit bool

	initialCap int
	priority   int64

	sysAddRemoveFn *gen.SystemMatchFn
	sysResizeFn    *gen.SystemMatchFn
}

func doCommentGroup(ctx *commentContext, g *ast.CommentGroup) {
	for _, line := range strings.Split(g.Text(), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "ecs:") {
			ctx.isECS = true
			if line[4:] == "component" {
				ctx.isComponent = true
			} else if line[4:] == "system" {
				ctx.isSystem = true
			}
		}
		if strings.HasPrefix(line, "name:") {
			ctx.name = strings.TrimSpace(line[5:])
		}
		if strings.HasPrefix(line, "uuid:") {
			ctx.uuid = strings.TrimSpace(line[5:])
		}
		if strings.HasPrefix(line, "components:") {
			ctx.componentsraw = strings.Split(strings.TrimSpace(line[11:]), ",")
			for i, v := range ctx.componentsraw {
				ctx.componentsraw[i] = strings.TrimSpace(v)
			}
		}
		if strings.HasPrefix(line, "async:") {
			ctx.isAsync, _ = strconv.ParseBool(strings.TrimSpace(line[6:]))
		}
		if strings.HasPrefix(line, "member:") {
			memberdecl := strings.Split(strings.TrimSpace(line[7:]), " ")
			m := gen.SystemMember{
				VarName: memberdecl[0],
			}
			if len(memberdecl) > 1 {
				mprefix, mtype := gen.ParsePrefixType(memberdecl[1])
				m.VarType = mtype
				m.VarPrefix = mprefix
			}
			if len(memberdecl) > 2 {
				m.VarPackagePath = memberdecl[2]
			}
			if len(memberdecl) > 3 {
				m.VarPackageAlias = memberdecl[3]
			}
			ctx.members = append(ctx.members, m)
		}
		if strings.HasPrefix(line, "setup:") {
			ctx.onSetup = strings.TrimSpace(line[6:])
		}
		if rest, ok := gen.LinePrefixMatch(line, "entityadded:", "entity added:", "entity-added:", "on-entity-added:", "entity_added:", "on-entity_added:"); ok {
			ctx.onEntityAdded = strings.TrimSpace(rest)
		}
		if rest, ok := gen.LinePrefixMatch(line, "entityremoved:", "entity removed:", "entity-removed:", "on-entity-removed:", "entity_removed:", "on-entity_removed:"); ok {
			ctx.onEntityRemoved = strings.TrimSpace(rest)
		}
		if rest, ok := gen.LinePrefixMatch(line, "componentwillresize:", "component will resize:", "component-will-resize:", "on-component-will-resize:", "component_will_resize:", "on-component-will-resize:"); ok {
			ctx.onComponentWillResize = strings.TrimSpace(rest)
		}
		if rest, ok := gen.LinePrefixMatch(line, "componentresized:", "component resized:", "component-resized:", "on-component-resized:", "component_resized:", "on-component-resized:"); ok {
			ctx.onComponentResized = strings.TrimSpace(rest)
		}

		if rest, ok := gen.LinePrefixMatch(line, "add:", "on-add:", "on_add:"); ok {
			ctx.onAdd = strings.TrimSpace(rest)
		}
		if rest, ok := gen.LinePrefixMatch(line, "remove:", "on-remove:", "on_remove:"); ok {
			ctx.onRemove = strings.TrimSpace(rest)
		}
		if rest, ok := gen.LinePrefixMatch(line, "willremove:", "will-remove:", "will_remove:", "on-will-remove:", "on_will_remove:"); ok {
			ctx.onWillRemove = strings.TrimSpace(rest)
		}
		if rest, ok := gen.LinePrefixMatch(line, "willresize:", "will-resize:", "will_resize:", "on-will-resize:", "on_will_resize:"); ok {
			ctx.onWillResize = strings.TrimSpace(rest)
		}
		if rest, ok := gen.LinePrefixMatch(line, "resized:", "on-resized:", "on_resized:"); ok {
			ctx.onResized = strings.TrimSpace(rest)
		}
		if rest, ok := gen.LinePrefixMatch(line, "no-init:", "no_init:", "noinit:"); ok {
			if v, _ := strconv.ParseBool(strings.TrimSpace(rest)); v {
				ctx.noInit = true
			}
		}
		if rest, ok := gen.LinePrefixMatch(line, "init:"); ok {
			if v, _ := strconv.ParseBool(strings.TrimSpace(rest)); !v {
				ctx.noInit = true
			}
		}
		if rest, ok := gen.LinePrefixMatch(line, "cap:", "capacity:"); ok {
			if v, _ := strconv.ParseInt(strings.TrimSpace(rest), 10, 64); v > 0 {
				ctx.initialCap = int(v)
			}
		}
		if rest, ok := gen.LinePrefixMatch(line, "priority:"); ok {
			if v, _ := strconv.ParseInt(strings.TrimSpace(rest), 10, 64); v > 0 {
				ctx.priority = v
			}
		}
		if rest, ok := gen.LinePrefixMatch(line, "addremovematchfn:", "matchfn:", "add-remove-match-fn:", "match-fn:", "add_remove_match_fn:", "match_fn:"); ok {
			ctx.sysAddRemoveFn = &gen.SystemMatchFn{
				Name: rest,
			}
			if vs := strings.Split(rest, " "); len(vs) > 1 {
				ctx.sysAddRemoveFn.Name = vs[0]
				ctx.sysAddRemoveFn.PackagePath = vs[1]
				if len(vs) > 2 {
					ctx.sysAddRemoveFn.PackageAlias = vs[2]
				}
			}
		}
		if rest, ok := gen.LinePrefixMatch(line, "resizematchfn:", "resize-match-fn:", "resize_match_fn:"); ok {
			ctx.sysResizeFn = &gen.SystemMatchFn{
				Name: rest,
			}
			if vs := strings.Split(rest, " "); len(vs) > 1 {
				ctx.sysResizeFn.Name = vs[0]
				ctx.sysResizeFn.PackagePath = vs[1]
				if len(vs) > 2 {
					ctx.sysResizeFn.PackageAlias = vs[2]
				}
			}
		}
	}
}

func run(c *cli.Context) error {
	println(c.String("go-package"))
	packageName := c.String("go-package")
	gofile := c.String("go-file")
	out := c.String("out")

	if out == "" {
		out = gofile[:len(gofile)-2] + "ecs.go"
	}

	println(gofile)
	// println(c.Args().First())

	fset := token.NewFileSet() // positions are relative to fset
	f, err := parser.ParseFile(fset, gofile, nil, parser.ParseComments)
	if err != nil {
		return err
	}
	//spew.Dump(f.Comments)

	jenf := jen.NewFile(packageName)
	jenf.PackageComment("Code generated by ecs https://github.com/gabstv/ecs; DO NOT EDIT.")

	for _, v := range f.Comments {
		end := v.End()
		ctx := &commentContext{}

		doCommentGroup(ctx, v)

		if ctx.isECS {

			if ctx.uuid == "" {
				//TODO: generate uuid on the fly
			}

			pos := fset.Position(end)
			println(pos.String())
			if ctx.name == "" {
				lstr, err := readFileLine(gofile, pos.Line+1)
				if err != nil {
					// comp/sys cannot be generated (no name)
				} else {
					l1 := strings.TrimSpace(lstr[5:])
					ctx.name = strings.Split(l1, " ")[0]
					println("name: " + ctx.name)
				}
			}
			println("")
			if ctx.isComponent {
				gen.Component(jenf, gen.ComponentDef{
					UUID:          ctx.uuid,
					StructName:    ctx.name,
					ComponentName: ctx.name + "Component",
					InitialCap:    ctx.initialCap,
					Async:         ctx.isAsync,
					NoInit:        ctx.noInit,
					OnSetup:       ctx.onSetup,
					OnWillResize:  ctx.onWillResize,
					OnResized:     ctx.onResized,
					OnAdd:         ctx.onAdd,
					OnRemove:      ctx.onRemove,
					OnWillRemove:  ctx.onWillRemove,
				})
			} else if ctx.isSystem {

				items := make([]gen.SystemViewItem, 0, len(ctx.componentsraw))
				for _, v := range ctx.componentsraw {
					items = append(items, gen.SystemViewItem{
						StructName: v,
						// PackagePath: ,
						// PackageAlias: ,
						// ComponentGetter: ,
						// FlagGetter: ,
					})
				}

				gen.System(jenf, gen.SystemDef{
					UUID:                  ctx.uuid,
					Name:                  ctx.name,
					Priority:              0,
					Async:                 ctx.isAsync,
					Components:            items,
					Members:               ctx.members,
					AddRemoveMatchFn:      ctx.sysAddRemoveFn,
					ResizeMatchFn:         ctx.sysResizeFn,
					OnSetup:               ctx.onSetup,
					OnEntityAdded:         ctx.onEntityAdded,
					OnEntityRemoved:       ctx.onEntityRemoved,
					OnComponentWillResize: ctx.onComponentWillResize,
					OnComponentResized:    ctx.onComponentResized,
				})
			}

		}
	}

	if err := jenf.Save(out); err != nil {
		println("ERRRRR: " + err.Error())
		return err
	}

	return nil
}

func readFileLine(file string, line int) (string, error) {
	f, err := os.Open(file)
	if err != nil {
		return "", err
	}
	defer f.Close()
	r := bufio.NewReader(f)
	l := 0
	for {
		l++
		sline, err := r.ReadString('\n')
		if err != nil {
			break
		}
		if l == line {
			return strings.TrimSpace(sline), nil
		}
	}
	return "", nil
}
