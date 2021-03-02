package gen

import (
	"fmt"

	"github.com/dave/jennifer/jen"
	"github.com/google/uuid"
)

type SystemViewItem struct {
	PackagePath     string `toml:"package_path" json:"package_path"`
	PackageAlias    string `toml:"package_alias" json:"package_alias"`
	StructName      string `toml:"struct_name" json:"struct_name"`
	ComponentGetter string `toml:"component_getter" json:"component_getter"`
	FlagGetter      string `toml:"flag_getter" json:"flag_getter"`
}

type SystemMember struct {
	VarName         string `toml:"var_name" json:"var_name"`
	VarType         string `toml:"var_type" json:"var_type"`
	VarPackagePath  string `toml:"var_package_path" json:"var_package_path"`
	VarPackageAlias string `toml:"var_package_alias" json:"var_package_alias"`
	VarPrefix       string `toml:"var_prefix" json:"var_prefix"`
}

type SystemMatchFn struct {
	PackagePath  string `toml:"package_path" json:"package_path"`
	PackageAlias string `toml:"package_alias" json:"package_alias"`
	Name         string `toml:"name" json:"name"`
}

type SystemDef struct {
	UUID       string           `toml:"uuid" json:"uuid"`
	Priority   int64            `toml:"priority" json:"priority"`
	Name       string           `toml:"name" json:"name"`
	Async      bool             `toml:"async" json:"async"`
	Components []SystemViewItem `toml:"components" json:"components"`
	Members    []SystemMember   `toml:"members" json:"members"`
	// Custom function to add or remove an entity to this system
	AddRemoveMatchFn *SystemMatchFn `toml:"add_remove_match_fn" json:"add_remove_match_fn"`
	// Custom function to rescan all component references when a component slice changes capacity
	ResizeMatchFn         *SystemMatchFn `toml:"resize_match_fn" json:"resize_match_fn"`
	OnEntityAdded         string         `toml:"on_entity_added" json:"on_entity_added"`
	OnEntityRemoved       string         `toml:"on_entity_removed" json:"on_entity_removed"`
	OnComponentWillResize string         `toml:"on_component_will_resize" json:"on_component_will_resize"`
	OnComponentResized    string         `toml:"on_component_resized" json:"on_component_resized"`
	OnSetup               string         `toml:"on_setup" json:"on_setup"`
}

func (d SystemDef) sanitize() SystemDef {
	if d.UUID == "" {
		id, err := uuid.NewRandom()
		if err != nil {
			panic(err)
		}
		d.UUID = id.String()
	}
	for i, v := range d.Components {
		if v.ComponentGetter == "" {
			v.ComponentGetter = "Get" + v.StructName + "ComponentData"
		}
		if v.FlagGetter == "" {
			v.FlagGetter = "Get" + v.StructName + "ComponentFlag"
		}
		d.Components[i] = v
	}
	return d
}

type systemIDs struct {
	uuid           string
	nameconst      string
	view           string
	viewItem       string
	viewItemSorter string
}

func System(f *jen.File, def SystemDef) {

	def = def.sanitize()

	if len(def.Components) < 1 {
		panic("a system needs at least one component")
	}

	// 	{{$type := (printf "%sSystem" .Name)}}
	// {{$uuid := (or .Vars.UUID $type)}}
	// {{$matchfn := (or .Vars.Match (printf "match%s" $type))}}
	// {{$resizematchfn := (or .Vars.ResizeMatch (printf "resizematch%s" $type))}}
	// {{$view := (printf "view%s" $type)}}
	// {{$viewitem := (printf "VI%s" $type)}}
	// {{$priority := (or .Vars.Priority "0")}}

	ids := systemIDs{
		uuid:           "uuid" + def.Name,
		nameconst:      "name" + def.Name,
		view:           "view" + def.Name,
		viewItem:       "VI" + def.Name,
		viewItemSorter: "sortedVI" + def.Name + "s",
	}

	f.Const().Defs(
		jen.Id(ids.uuid).Op("=").Lit(def.UUID),
		jen.Id(ids.nameconst).Op("=").Lit(def.Name),
	)

	systemView(f, def, ids)

	members := make([]jen.Code, 0)
	for _, v := range def.Members {
		members = append(members, jen.Id(v.VarName).Op(v.VarPrefix).Qual(v.VarPackagePath, v.VarType))
	}

	// system struct def
	f.Commentf("%s implements ecs.System", def.Name)
	f.Type().Id(def.Name).Struct(
		jen.Qual(ecspkg, "BaseSystem"),
		jen.Id("initialized").Bool(),
		jen.Id("world").Qual(ecspkg, "World"),
		jen.Id("view").Op("*").Id(ids.view),
		jen.Add(members...),
	)

	f.Commentf("Get%s returns the instance of the system in a World", def.Name)
	f.Func().Id("Get" + def.Name).Params(jen.Id("w").Qual(ecspkg, "World")).Op("*").Id(def.Name).Block(
		jen.Return(jen.Id("w").Dot("S").Call(jen.Id(ids.uuid)).Assert(jen.Op("*").Id(def.Name))),
	)

	f.Comment("UUID implements ecs.System")
	f.Func().Params(jen.Id(def.Name)).Id("UUID").Params().String().Block(
		jen.Return(jen.Id(ids.uuid)),
	)

	f.Comment("Name implements ecs.System")
	f.Func().Params(jen.Id(def.Name)).Id("Name").Params().String().Block(
		jen.Return(jen.Id(ids.nameconst)),
	)

	// if no custom match function is defined, we use the standard logic to create one
	if def.AddRemoveMatchFn == nil {
		def.AddRemoveMatchFn = &SystemMatchFn{
			Name: "addRemoveMatchFn" + def.Name,
		}
		f.Func().Id(def.AddRemoveMatchFn.Name).Params(jen.Id("f").Qual(ecspkg, "Flag"), jen.Id("w").Qual(ecspkg, "World")).Bool().Block(func() []jen.Code {
			out := make([]jen.Code, 0)
			for _, v := range def.Components {
				blk := jen.If(jen.Op("!").Id("f").Dot("Contains").Call(jen.Qual(v.PackagePath, v.FlagGetter).Call(jen.Id("w")))).Block(
					jen.Return(jen.Lit(false)),
				)
				out = append(out, blk)
			}
			out = append(out, jen.Return(jen.Lit(true)))
			return out
		}()...)
	}
	if def.ResizeMatchFn == nil {
		def.ResizeMatchFn = &SystemMatchFn{
			Name: "resizeMatchFn" + def.Name,
		}
		f.Func().Id(def.ResizeMatchFn.Name).Params(jen.Id("f").Qual(ecspkg, "Flag"), jen.Id("w").Qual(ecspkg, "World")).Bool().Block(
			jen.Return(jen.Qual(def.AddRemoveMatchFn.PackagePath, def.AddRemoveMatchFn.Name).Call(jen.Id("f"), jen.Id("w"))),
		)
	}

	f.Comment("ensure matchfn")
	f.Var().Id("_").Qual(ecspkg, "MatchFn").Op("=").Qual(def.AddRemoveMatchFn.PackagePath, def.AddRemoveMatchFn.Name)

	f.Comment("ensure resizematchfn")
	f.Var().Id("_").Qual(ecspkg, "MatchFn").Op("=").Qual(def.ResizeMatchFn.PackagePath, def.ResizeMatchFn.Name)

	sysPtr := jen.Id("s").Op("*").Id(def.Name)

	f.Func().Params(sysPtr).Id("match").Params(jen.Id("eflag").Qual(ecspkg, "Flag")).Bool().Block(
		jen.Return(jen.Qual(def.AddRemoveMatchFn.PackagePath, def.AddRemoveMatchFn.Name).Call(jen.Id("eflag"), jen.Id("s").Dot("world"))),
	)

	f.Func().Params(sysPtr).Id("resizematch").Params(jen.Id("eflag").Qual(ecspkg, "Flag")).Bool().Block(
		jen.Return(jen.Qual(def.ResizeMatchFn.PackagePath, def.ResizeMatchFn.Name).Call(jen.Id("eflag"), jen.Id("s").Dot("world"))),
	)

	enAdded := jen.Empty()
	enRemoved := jen.Empty()

	if def.OnEntityAdded != "" {
		enAdded = jen.Id("s").Dot(def.OnEntityAdded).Call(jen.Id("e"))
	}
	if def.OnEntityRemoved != "" {
		enRemoved = jen.Id("s").Dot(def.OnEntityRemoved).Call(jen.Id("e"))
	}

	f.Func().Params(sysPtr).Id("ComponentAdded").Params(jen.Id("e").Qual(ecspkg, "Entity"), jen.Id("eflag").Qual(ecspkg, "Flag")).Block(
		jen.If(jen.Id("s").Dot("match").Call(jen.Id("eflag"))).Block(
			jen.If(jen.Id("s").Dot("view").Dot("Add").Call(jen.Id("e"))).Block(
				// TODO: dispatch event that this entity was added to this system
				enAdded,
			),
		).Else().Block(
			jen.If(jen.Id("s").Dot("view").Dot("Remove").Call(jen.Id("e"))).Block(
				// TODO: dispatch event that this entity was removed from this system
				enRemoved,
			),
		),
	)

	f.Func().Params(sysPtr).Id("ComponentRemoved").Params(jen.Id("e").Qual(ecspkg, "Entity"), jen.Id("eflag").Qual(ecspkg, "Flag")).Block(
		jen.If(jen.Id("s").Dot("match").Call(jen.Id("eflag"))).Block(
			jen.If(jen.Id("s").Dot("view").Dot("Add").Call(jen.Id("e"))).Block(
				// TODO: dispatch event that this entity was added to this system
				enAdded,
			),
		).Else().Block(
			jen.If(jen.Id("s").Dot("view").Dot("Remove").Call(jen.Id("e"))).Block(
				// TODO: dispatch event that this entity was removed from this system
				enRemoved,
			),
		),
	)

	cnResized := jen.Empty()
	cnWillResize := jen.Empty()

	if def.OnComponentResized != "" {
		cnResized = jen.Id("s").Dot(def.OnComponentResized).Call(jen.Id("cflag"))
	}
	if def.OnComponentWillResize != "" {
		cnWillResize = jen.Id("s").Dot(def.OnComponentWillResize).Call(jen.Id("cflag"))
	}

	f.Func().Params(sysPtr).Id("ComponentResized").Params(jen.Id("cflag").Qual(ecspkg, "Flag")).Block(
		jen.If(jen.Id("s").Dot("resizematch").Call(jen.Id("cflag"))).Block(
			jen.Id("s").Dot("view").Dot("rescan").Call(),
			cnResized,
		),
	)

	f.Func().Params(sysPtr).Id("ComponentWillResize").Params(jen.Id("cflag").Qual(ecspkg, "Flag")).Block(
		jen.If(jen.Id("s").Dot("resizematch").Call(jen.Id("cflag"))).Block(
			cnWillResize,
			jen.Id("s").Dot("view").Dot("clearpointers").Call(),
		),
	)

	f.Func().Params(sysPtr).Id("V").Params().Op("*").Id(ids.view).Block(
		jen.Return(jen.Id("s").Dot("view")),
	)

	f.Func().Params(sysPtr).Id("View").Params().Op("*").Id(ids.view).Block(
		jen.Return(jen.Id("s").Dot("view")),
	)

	f.Func().Params(sysPtr).Id("Priority").Params().Int64().Block(
		jen.Return(jen.Lit(def.Priority)),
	)

	onSetup := jen.Empty()
	if def.OnSetup != "" {
		onSetup = jen.Id("s").Dot(def.OnSetup).Call(jen.Id("w"))
	}

	f.Func().Params(sysPtr).Id("Setup").Params(jen.Id("w").Qual(ecspkg, "World")).Block(
		jen.If(jen.Id("s").Dot("initialized")).Block(jen.Panic(jen.Lit(fmt.Sprintf("%s called Setup() more than once", def.Name)))),
		jen.Id("s").Dot("view").Op("=").Id("new"+ids.view).Call(jen.Id("w")),
		jen.Id("s").Dot("world").Op("=").Id("w"),
		jen.Id("s").Dot("BaseSystem").Dot("Enable").Call(),
		jen.Id("s").Dot("initialized").Op("=").Lit(true),
		onSetup,
	)

	// {{if not .SkipRegister}}
	f.Func().Id("init").Params().Block(
		jen.Qual(ecspkg, "RegisterSystem").Call(
			jen.Func().Params().Qual(ecspkg, "System").Block(
				jen.Return(jen.Op("&").Id(def.Name).Block()),
			),
		),
	)
	// {{end}}
}

func systemView(f *jen.File, def SystemDef, ids systemIDs) {
	f.Type().Id(ids.view).Struct(
		jen.Id("entities").Index().Id(ids.viewItem),
		jen.Id("world").Qual(ecspkg, "World"),
		func() *jen.Statement {
			if def.Async {
				return jen.Id("l").Qual("sync", "RWMutex")
			}
			return jen.Empty()
		}(),
	)

	vitypes := make([]jen.Code, 1)
	vitypes[0] = jen.Id("Entity").Qual(ecspkg, "Entity")

	for _, v := range def.Components {
		if v.PackageAlias != "" {
			f.ImportName(v.PackagePath, v.PackageAlias)
		}
		vitypes = append(vitypes, jen.Id(v.StructName).Op("*").Qual(v.PackagePath, v.StructName))
	}

	f.Type().Id(ids.viewItem).Struct(vitypes...)

	ijparams := func() []jen.Code {
		return []jen.Code{jen.Id("i").Int(), jen.Id("j").Int()}
	}

	f.Type().Id(ids.viewItemSorter).Index().Id(ids.viewItem)
	f.Comment("Len implementation of sort.Interface.Len")
	f.Func().Params(jen.Id("a").Id(ids.viewItemSorter)).Id("Len").Params().Int().Block(jen.Return(jen.Len(jen.Id("a"))))
	f.Comment("Swap implementation of sort.Interface.Swap")
	f.Func().Params(jen.Id("a").Id(ids.viewItemSorter)).Id("Swap").Params(ijparams()...).Block(
		jen.List(jen.Id("a").Index(jen.Id("i")), jen.Id("a").Index(jen.Id("j"))).Op("=").List(jen.Id("a").Index(jen.Id("j")), jen.Id("a").Index(jen.Id("i"))),
	)
	f.Comment("Less implementation of sort.Interface.Less")
	f.Func().Params(jen.Id("a").Id(ids.viewItemSorter)).Id("Less").Params(ijparams()...).Bool().Block(
		jen.Return(jen.Id("a").Index(jen.Id("i")).Dot("Entity").Op("<").Id("a").Index(jen.Id("j")).Dot("Entity")),
	)

	f.Func().Id("new" + ids.view).Params(jen.Id("w").Qual(ecspkg, "World")).Op("*").Id(ids.view).Block(
		jen.Return(jen.Op("&").Id(ids.view).Block(
			jen.Id("entities").Op(":").Make(jen.Index().Id(ids.viewItem), jen.Lit(0)).Op(","),
			jen.Id("world").Op(":").Id("w").Op(","),
		)),
	)

	viewx := func() *jen.Statement {
		return jen.Id("v").Op("*").Id(ids.view)
	}

	f.Func().Params(viewx()).Id("Matches").Params().Index().Id(ids.viewItem).Block(
		func() []jen.Code {
			if def.Async {
				return []jen.Code{
					jen.Id("v").Dot("l").Dot("RLock").Call(),
					jen.Defer().Id("v").Dot("l").Dot("RUnlock").Call(),
					jen.Id("eclone").Op(":=").Make(
						jen.Index().Qual(ecspkg, "Entity"),
						jen.Len(jen.Id("v").Dot("entities")),
					),
					jen.Copy(jen.Id("eclone"), jen.Id("v").Dot("entities")),
					jen.Return(jen.Id("eclone")),
				}
			}
			return []jen.Code{jen.Return(jen.Id("v").Dot("entities"))}
		}()...,
	)

	f.Func().Params(viewx()).Id("indexof").Params(jen.Id("e").Qual(ecspkg, "Entity")).Int().Block(
		jen.Id("i").Op(":=").Qual("sort", "Search").Call(
			jen.Len(jen.Id("v").Dot("entities")),
			jen.Func().Params(jen.Id("i").Int()).Bool().Block(jen.Return(jen.Id("v").Dot("entities").Index(jen.Id("i")).Dot("Entity").Op(">=").Id("e"))),
		),
		jen.If(jen.Id("i").Op("<").Len(jen.Id("v").Dot("entities")).Op("&&").Id("v").Dot("entities").Index(jen.Id("i")).Dot("Entity").Op("==").Id("e")).Block(
			jen.Return(jen.Id("i")),
		),
		jen.Return(jen.Lit(-1)),
	)

	f.Comment("Fetch a specific entity")
	f.Func().Params(viewx()).Id("Fetch").Params(jen.Id("e").Qual(ecspkg, "Entity")).Params(
		jen.Id("data").Id(ids.viewItem),
		jen.Id("ok").Bool(),
	).Block(
		mutexLU(def.Async, "v.l.RLock", false),
		mutexLU(def.Async, "v.l.RUnlock", true),
		jen.Id("i").Op(":=").Id("v").Dot("indexof").Call(jen.Id("e")),
		jen.If(jen.Id("i").Op("==").Lit(-1)).Block(
			jen.Return(jen.Id(ids.viewItem).Block(), jen.Lit(false)),
		),
		jen.Return(jen.Id("v").Dot("entities").Index(jen.Id("i")), jen.Lit(true)),
	)

	getvitypes := make([]jen.Code, 1)
	getvitypes[0] = jen.Id("Entity").Op(":").Id("e").Op(",")
	for _, v := range def.Components {
		getvitypes = append(getvitypes, jen.Id(v.StructName).Op(":").Qual(v.PackagePath, v.ComponentGetter).Call(jen.Id("v").Dot("world"), jen.Id("e")).Op(","))
	}

	f.Comment("Add a new entity")
	f.Func().Params(viewx()).Id("Add").Params(jen.Id("e").Qual(ecspkg, "Entity")).Bool().Block(
		mutexLU(def.Async, "v.l.RLock", false),
		mutexLU(def.Async, "v.l.RUnlock", true),
		jen.Comment("MUST NOT add an Entity twice:"),
		jen.If(jen.Id("i").Op(":=").Id("v").Dot("indexof").Call(jen.Id("e")), jen.Id("i").Op(">").Lit(-1)).Block(
			jen.Return(jen.Lit(false)),
		),
		jen.Id("v").Dot("entities").Op("=").Append(jen.Id("v").Dot("entities"), jen.Id(ids.viewItem).Block(getvitypes...)),
		jen.If(jen.Len(jen.Id("v").Dot("entities")).Op(">").Lit(1)).Block(
			jen.If(jen.Id("v").Dot("entities").Index(jen.Len(jen.Id("v").Dot("entities")).Op("-").Lit(1)).Dot("Entity").Op("<").Id("v").Dot("entities").Index(jen.Len(jen.Id("v").Dot("entities")).Op("-").Lit(2)).Dot("Entity")).Block(
				jen.Qual("sort", "Sort").Call(jen.Id(ids.viewItemSorter).Call(jen.Id("v").Dot("entities"))),
			),
		),
		jen.Return(jen.Lit(true)),
	)

	f.Comment("Remove an entity")
	f.Func().Params(viewx()).Id("Remove").Params(jen.Id("e").Qual(ecspkg, "Entity")).Bool().Block(
		mutexLU(def.Async, "v.l.Lock", false),
		mutexLU(def.Async, "v.l.Unlock", true),
		jen.If(jen.Id("i").Op(":=").Id("v").Dot("indexof").Call(jen.Id("e")), jen.Id("i").Op("!=").Lit(-1)).Block(
			jen.Id("v").Dot("entities").Op("=").Append(jen.Id("v").Dot("entities").Index(jen.Empty(), jen.Id("i")), jen.Id("v").Dot("entities").Index(jen.Id("i").Op("+").Lit(1), jen.Empty()).Op("...")),
			jen.Return(jen.Lit(true)),
		),
		jen.Return(jen.Lit(false)),
	)

	clearptrsvi := make([]jen.Code, 1)
	clearptrsvi[0] = jen.Id("e").Op(":=").Id("v").Dot("entities").Index(jen.Id("i")).Dot("Entity")
	for _, v := range def.Components {
		clearptrsvi = append(clearptrsvi, jen.Id("v").Dot("entities").Index(jen.Id("i")).Dot(v.StructName).Op("=").Nil())
	}
	clearptrsvi = append(clearptrsvi, jen.Id("_").Op("=").Id("e"))

	f.Func().Params(viewx()).Id("clearpointers").Params().Block(
		mutexLU(def.Async, "v.l.Lock", false),
		mutexLU(def.Async, "v.l.Unlock", true),
		jen.For(jen.Id("i").Op(":=").Range().Id("v").Dot("entities")).Block(clearptrsvi...),
	)

	rescanvi := make([]jen.Code, 1)
	rescanvi[0] = jen.Id("e").Op(":=").Id("v").Dot("entities").Index(jen.Id("i")).Dot("Entity")
	for _, v := range def.Components {
		rescanvi = append(rescanvi, jen.Id("v").Dot("entities").Index(jen.Id("i")).Dot(v.StructName).Op("=").Qual(v.PackagePath, v.ComponentGetter).Call(jen.Id("v").Dot("world"), jen.Id("e")))
	}
	rescanvi = append(rescanvi, jen.Id("_").Op("=").Id("e"))
	//TODO: {{if .Vars.ViewRescan}}
	// rescanvi = append(rescanvi, call of whatev
	// {{//}}

	f.Func().Params(viewx()).Id("rescan").Params().Block(
		mutexLU(def.Async, "v.l.Lock", false),
		mutexLU(def.Async, "v.l.Unlock", true),
		jen.For(jen.Id("i").Op(":=").Range().Id("v").Dot("entities")).Block(rescanvi...),
	)
}

// func systemView(f *jen.File, def SystemDef, ids systemIDs) {
