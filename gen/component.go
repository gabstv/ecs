package gen

import (
	"fmt"

	"github.com/dave/jennifer/jen"
	"github.com/google/uuid"
)

const (
	ecspkg = "github.com/gabstv/ecs/v2"
)

type ComponentDef struct {
	UUID          string
	StructName    string // if the ComponentName is "PositionComponent", this would be "Position"
	ComponentName string // if the StructName is "Position", this would be "PositionComponent"
	InitialCap    int    // initial slice capacity (default: 256)
	Async         bool   // Adds a mutex lock for accessing component data in parallel. Not recommended.
	NoInit        bool   // if true, the generated code does not define an init() to register this component
}

func (d ComponentDef) sanitize() ComponentDef {
	if d.UUID == "" {
		id, err := uuid.NewRandom()
		if err != nil {
			panic(err)
		}
		d.UUID = id.String()
	}
	if d.InitialCap == 0 {
		d.InitialCap = 256
	}
	if d.ComponentName == "" {
		d.ComponentName = d.StructName + "Component"
	}
	return d
}

type componentIDs struct {
	uuid        string
	cap         string
	drawer      string
	drawerSlice string
	watcher     string
	watcherImpl string
}

func Component(f *jen.File, def ComponentDef) {

	def = def.sanitize()

	ids := componentIDs{
		uuid:        "uuid" + def.ComponentName,
		cap:         "cap" + def.ComponentName,
		drawer:      "drawer" + def.ComponentName,
		drawerSlice: "slcdrawer" + def.ComponentName,
		watcher:     def.ComponentName + "Observer",
		watcherImpl: "m" + def.ComponentName + "Observer",
	}

	f.ImportName(ecspkg, "ecs")

	f.Const().Defs(
		jen.Id(ids.uuid).Op("=").Lit(def.UUID),
		jen.Id(ids.cap).Op("=").Lit(def.InitialCap),
	)

	f.Type().Id(ids.drawer).Struct(
		jen.Id("Entity").Qual(ecspkg, "Entity"),
		jen.Id("Data").Id(def.StructName),
	)

	f.Commentf("%s is a helper struct to access a valid pointer of %s", ids.watcher, def.StructName)
	f.Type().Id(ids.watcher).Interface(
		jen.Id("Entity").Params().Qual(ecspkg, "Entity"),
		jen.Id("Data").Params().Op("*").Id(def.StructName),
	)

	ijparams := func() []jen.Code {
		return []jen.Code{jen.Id("i").Int(), jen.Id("j").Int()}
	}

	f.Type().Id(ids.drawerSlice).Index().Id(ids.drawer)
	f.Comment("Len implementation of sort.Interface.Len")
	f.Func().Params(jen.Id("a").Id(ids.drawerSlice)).Id("Len").Params().Int().Block(jen.Return(jen.Len(jen.Id("a"))))
	f.Comment("Swap implementation of sort.Interface.Swap")
	f.Func().Params(jen.Id("a").Id(ids.drawerSlice)).Id("Swap").Params(ijparams()...).Block(
		jen.List(jen.Id("a").Index(jen.Id("i")), jen.Id("a").Index(jen.Id("j"))).Op("=").List(jen.Id("a").Index(jen.Id("j")), jen.Id("a").Index(jen.Id("i"))),
	)
	f.Comment("Less implementation of sort.Interface.Less")
	f.Func().Params(jen.Id("a").Id(ids.drawerSlice)).Id("Less").Params(ijparams()...).Bool().Block(
		jen.Return(jen.Id("a").Index(jen.Id("i")).Dot("Entity").Op("<").Id("a").Index(jen.Id("j")).Dot("Entity")),
	)

	f.Type().Id(ids.watcherImpl).Struct(
		jen.Id("c").Op("*").Id(def.ComponentName),
		jen.Id("entity").Qual(ecspkg, "Entity"),
	)

	f.Func().Params(jen.Id("w").Op("*").Id(ids.watcherImpl)).Id("Entity").Params().Qual(ecspkg, "Entity").Block(
		jen.Return(jen.Id("w").Dot("entity")),
	)

	f.Func().Params(jen.Id("w").Op("*").Id(ids.watcherImpl)).Id("Data").Params().Op("*").Id(def.StructName).Block(
		mutexLU(def.Async, "w.c.l.RLock", false),
		mutexLU(def.Async, "w.c.l.RUnlock", true),
		jen.Id("id").Op(":=").Id("w").Dot("c").Dot("indexof").Call(jen.Id("w").Dot("entity")),
		jen.If(jen.Id("id").Op("==").Lit(-1)).Block(jen.Return(jen.Nil())),
		jen.Return(jen.Op("&").Id("w").Dot("c").Dot("data").Index(jen.Id("id")).Dot("Data")),
	)

	f.Commentf("%s handles the component logic of the %s. Implements ecs.Component", def.ComponentName, def.StructName)
	f.Type().Id(def.ComponentName).Struct(
		jen.Id("initialized").Bool(),
		jen.Id("flag").Qual(ecspkg, "Flag"),
		jen.Id("world").Qual(ecspkg, "World"),
		jen.Id("wkey").Index(jen.Lit(4)).Byte(),
		jen.Id("data").Index().Id(ids.drawer),
		func() *jen.Statement {
			if def.Async {
				return jen.Id("l").Qual("sync", "RWMutex")
			}
			return jen.Empty()
		}(),
	)

	f.Commentf("Get%s returns the instance of the component in a World", def.ComponentName)
	f.Func().Id("Get" + def.ComponentName).Params(jen.Id("w").Qual(ecspkg, "World")).Op("*").Id(def.ComponentName).Block(
		jen.Return(jen.Id("w").Dot("C").Call(jen.Id(ids.uuid))).Assert(jen.Op("*").Id(def.ComponentName)),
	)

	f.Commentf("Get%sData returns the data of the component in a World", def.ComponentName)
	f.Func().Id("Get"+def.ComponentName+"Data").Params(jen.Id("w").Qual(ecspkg, "World"), jen.Id("e").Qual(ecspkg, "Entity")).Op("*").Id(def.StructName).Block(
		jen.Return(jen.Id("Get" + def.ComponentName).Call(jen.Id("w")).Dot("Data").Call(jen.Id("e"))),
	)

	f.Commentf("Set%sData updates/adds a %s to Entity e", def.ComponentName, def.StructName)
	f.Func().Id("Set"+def.ComponentName+"Data").Params(
		jen.Id("w").Qual(ecspkg, "World"),
		jen.Id("e").Qual(ecspkg, "Entity"),
		jen.Id("data").Id(def.StructName),
	).Block(
		jen.Id("Get"+def.ComponentName).Call(jen.Id("w")).Dot("Upsert").Call(jen.Id("e"), jen.Id("data")),
	)

	f.Commentf("Watch%sData gets a pointer getter of an entity's %s.", def.ComponentName, def.StructName)
	f.Comment("The pointer must not be stored because it may become invalid over time.")
	f.Func().Id("Watch"+def.ComponentName+"Data").Params(
		jen.Id("w").Qual(ecspkg, "World"),
		jen.Id("e").Qual(ecspkg, "Entity"),
	).Id(ids.watcher).Block(
		jen.Return(jen.Op("&").Id(ids.watcherImpl)).Block(
			jen.Id("c").Op(":").Id("Get"+def.ComponentName).Call(jen.Id("w")).Op(","),
			jen.Id("entity").Op(":").Id("e").Op(","),
		),
	)

	f.Comment("UUID implements ecs.Component")
	f.Func().Params(jen.Id(def.ComponentName)).Id("UUID").Params().String().Block(
		jen.Return(jen.Id(ids.uuid)),
	)

	f.Comment("Name implements ecs.Component")
	f.Func().Params(jen.Id(def.ComponentName)).Id("Name").Params().String().Block(
		jen.Return(jen.Lit(def.ComponentName)),
	)

	f.Comment("World implements ecs.Component")
	f.Func().Params(jen.Id("c").Op("*").Id(def.ComponentName)).Id("World").Params().Qual(ecspkg, "World").Block(
		jen.Return(jen.Id("c").Dot("world")),
	)

	f.Func().Params(jen.Id("c").Op("*").Id(def.ComponentName)).Id("indexof").Params(jen.Id("e").Qual(ecspkg, "Entity")).Int().Block(
		jen.Id("i").Op(":=").Qual("sort", "Search").Call(
			jen.Len(jen.Id("c").Dot("data")),
			jen.Func().Params(jen.Id("i").Int()).Bool().Block(jen.Return(jen.Id("c").Dot("data").Index(jen.Id("i")).Dot("Entity").Op(">=").Id("e"))),
		),
		jen.If(jen.Id("i").Op("<").Len(jen.Id("c").Dot("data")).Op("&&").Id("c").Dot("data").Index(jen.Id("i")).Dot("Entity").Op("==").Id("e")).Block(
			jen.Return(jen.Id("i")),
		),
		jen.Return(jen.Lit(-1)),
	)

	// Upsert()
	componentUpsert(f, def, ids)

	// Remove()
	componentRemove(f, def, ids)

	// Data(e)
	f.Commentf("Data retrieves the *%s of entity e. Warning: volatile pointer (do not store)", def.StructName)
	f.Func().Params(jen.Id("c").Op("*").Id(def.ComponentName)).Id("Data").Params(jen.Id("e").Qual(ecspkg, "Entity")).Op("*").Id(def.StructName).Block(
		mutexLU(def.Async, "c.l.RLock", false),
		mutexLU(def.Async, "c.l.RUnlock", true),
		jen.Id("index").Op(":=").Id("c").Dot("indexof").Call(jen.Id("e")),
		jen.If(jen.Id("index").Op(">").Lit(-1)).Block(
			jen.Return(jen.Op("&").Id("c").Dot("data").Index(jen.Id("index")).Dot("Data")),
		),
		jen.Return(jen.Nil()),
	)

	// Flag() ecs.Flag
	f.Comment("Flag returns the flag of this component")
	f.Func().Params(jen.Id("c").Op("*").Id(def.ComponentName)).Id("Flag").Params().Qual(ecspkg, "Flag").Block(
		jen.Return(jen.Id("c").Dot("flag")),
	)

	// Global Flag
	f.Commentf("Get%sFlag returns the flag of this component", def.ComponentName)
	// "Get" + def.ComponentName + "Flag"
	f.Func().Id("Get"+def.ComponentName+"Flag").Params(jen.Id("w").Qual(ecspkg, "World")).Qual(ecspkg, "Flag").Block(
		jen.Return(jen.Id("Get" + def.ComponentName).Call(jen.Id("w")).Dot("Flag").Call()),
	)

	// Setup()
	componentSetup(f, def, ids)

	// init()
	componentInit(f, def, ids)
}

func componentUpsert(f *jen.File, def ComponentDef, ids componentIDs) {

	f.Comment("Upsert creates or updates a component data of an entity.")
	f.Commentf("Not recommended to be used directly. Use Set%sData to change component", def.ComponentName)
	f.Comment("data outside of a system loop.")

	f.Func().Params(jen.Id("c").Op("*").Id(def.ComponentName)).Id("Upsert").Params(
		jen.Id("e").Qual(ecspkg, "Entity"),
		jen.Id("data").Interface(),
	).Params(jen.Id("added").Bool()).Block(
		jen.Id("added").Op("=").Lit(false),
		jen.List(jen.Id("v"), jen.Id("ok")).Op(":=").Id("data").Assert(jen.Id(def.StructName)),
		jen.If(jen.Op("!").Id("ok")).Block(jen.Panic(jen.Lit(fmt.Sprintf("data must be %s", def.StructName)))),
		mutexLU(def.Async, "c.l.RLock", false),
		jen.Id("id").Op(":=").Id("c").Dot("indexof").Call(jen.Id("e")),
		mutexLU(def.Async, "c.l.RUnlock", false),
		jen.If(jen.Id("id").Op(">").Lit(-1)).Block(
			mutexLU(def.Async, "c.l.Lock", false),
			jen.Id("dwr").Op(":=").Op("&").Id("c").Dot("data").Index(jen.Id("id")),
			jen.Id("dwr").Dot("Data").Op("=").Id("v"),
			mutexLU(def.Async, "c.l.Unlock", false),
		).Else().Block(
			jen.Id("added").Op("=").Lit(true),
		),
		mutexLU(def.Async, "c.l.Lock", false),
		jen.Id("rsz").Op(":=").Lit(false),
		jen.If(jen.Cap(jen.Id("c").Dot("data")).Op("==").Len(jen.Id("c").Dot("data"))).Block(
			jen.Id("rsz").Op("=").Lit(true),
			jen.Id("c").Dot("world").Dot("CWillResize").Call(jen.Id("c"), jen.Id("c").Dot("wkey")),
			// TODO: {{if .Vars.OnWillResize}}{{.Vars.OnWillResize}}{{end}}
		),
		jen.Id("newindex").Op(":=").Len(jen.Id("c").Dot("data")),
		jen.Id("c").Dot("data").Op("=").Append(jen.Id("c").Dot("data"), jen.Id(ids.drawer).Block(
			jen.Id("Entity").Op(":").Id("e").Op(","),
			jen.Id("Data").Op(":").Id("v").Op(","),
		)),
		jen.If(jen.Len(jen.Id("c").Dot("data")).Op(">").Lit(1)).Block(
			jen.If(jen.Id("c").Dot("data").Index(jen.Id("newindex")).Dot("Entity").Op("<").Id("c").Dot("data").Index(jen.Id("newindex").Op("-").Lit(1)).Dot("Entity")).Block(
				jen.Id("c").Dot("world").Dot("CWillResize").Call(jen.Id("c"), jen.Id("c").Dot("wkey")),
				// TODO: {{if .Vars.OnWillResize}}{{.Vars.OnWillResize}}{{end}}
				jen.Qual("sort", "Sort").Call(jen.Id(ids.drawerSlice).Call(jen.Id("c").Dot("data"))),
				jen.Id("rsz").Op("=").Lit(true),
			),
		),
		mutexLU(def.Async, "c.l.Unlock", false),
		jen.If(jen.Id("rsz")).Block(
			// TODO: {{if .Vars.OnResize}}{{.Vars.OnResize}}{{end}}
			jen.Id("c").Dot("world").Dot("CResized").Call(jen.Id("c"), jen.Id("c").Dot("wkey")),
			jen.Qual(ecspkg, "DispatchComponentEvent").Call(jen.Id("c"), jen.Qual(ecspkg, "EvtComponentsResized"), jen.Lit(0)),
		),
		// TODO: {{if .Vars.OnAdd}}{{.Vars.OnAdd}}{{end}}
		jen.Id("c").Dot("world").Dot("CAdded").Call(jen.Id("e"), jen.Id("c"), jen.Id("c").Dot("wkey")),
		jen.Qual(ecspkg, "DispatchComponentEvent").Call(jen.Id("c"), jen.Qual(ecspkg, "EvtComponentAdded"), jen.Id("e")),
		jen.Return(),
	)
}

func componentRemove(f *jen.File, def ComponentDef, ids componentIDs) {

	lock := func() *jen.Statement {
		if def.Async {
			return jen.Id("c").Dot("l").Dot("Lock").Call()
		}
		return jen.Empty()
	}
	unlock := func() *jen.Statement {
		if def.Async {
			return jen.Defer().Id("c").Dot("l").Dot("Unlock").Call()
		}
		return jen.Empty()
	}

	f.Commentf("Remove a %s data from entity e", def.StructName)
	if def.Async {
		f.Comment("Warning: DO NOT call remove inside the system entities loop")
	}
	f.Comment("Returns false if the component was not present in Entity e")
	f.Func().Params(jen.Id("c").Op("*").Id(def.ComponentName)).Id("Remove").Params(jen.Id("e").Qual(ecspkg, "Entity")).Bool().Block(
		lock(),
		unlock(),
		jen.Id("i").Op(":=").Id("c").Dot("indexof").Call(jen.Id("e")),
		jen.If(jen.Id("i").Op("==").Lit(-1)).Block(jen.Return(jen.Lit(false))),
		// TODO: {{if .Vars.BeforeRemove}}{{.Vars.BeforeRemove}}{{end}}
		jen.Id("c").Dot("data").Op("=").Id("c").Dot("data").Index(
			jen.Empty(),
			jen.Id("i").Op("+").Copy(
				jen.Id("c").Dot("data").Index(jen.Id("i"), jen.Empty()),
				jen.Id("c").Dot("data").Index(jen.Id("i").Op("+").Lit(1), jen.Empty()),
			),
		),
		jen.Id("c").Dot("world").Dot("CRemoved").Call(jen.Id("e"), jen.Id("c"), jen.Id("c").Dot("wkey")),
		// TODO: {{if .Vars.OnRemove}}{{.Vars.OnRemove}}{{end}}
		jen.Qual(ecspkg, "DispatchComponentEvent").Call(jen.Id("c"), jen.Qual(ecspkg, "EvtComponentRemoved"), jen.Id("e")),
		jen.Return(jen.Lit(true)),
	)
}

func componentSetup(f *jen.File, def ComponentDef, ids componentIDs) {
	f.Comment("Setup is called by ecs.World")
	f.Comment("Do not call this by yourself")
	f.Func().Params(jen.Id("c").Op("*").Id(def.ComponentName)).Id("Setup").Params(
		jen.Id("w").Qual(ecspkg, "World"),
		jen.Id("f").Qual(ecspkg, "Flag"),
		jen.Id("key").Index(jen.Lit(4)).Byte(),
	).Block(
		jen.If(jen.Id("c").Dot("initialized")).Block(jen.Panic(jen.Lit(fmt.Sprintf("%s called Setup() more than once", def.ComponentName)))),
		jen.Id("c").Dot("flag").Op("=").Id("f"),
		jen.Id("c").Dot("world").Op("=").Id("w"),
		jen.Id("c").Dot("wkey").Op("=").Id("key"),
		jen.Id("c").Dot("data").Op("=").Make(jen.Index().Id(ids.drawer), jen.Lit(0), jen.Id(ids.cap)),
		jen.Id("c").Dot("initialized").Op("=").Lit(true),
		// TODO: {{if .Vars.Setup}}{{.Vars.Setup}}{{end}}
	)
}

func componentInit(f *jen.File, def ComponentDef, ids componentIDs) {
	if def.NoInit {
		return
	}
	f.Func().Id("init").Params().Block(
		jen.Qual(ecspkg, "RegisterComponent").Call(
			jen.Func().Params().Qual(ecspkg, "Component").Block(
				jen.Return(jen.Op("&").Id(def.ComponentName).Block()),
			),
		),
	)
}
