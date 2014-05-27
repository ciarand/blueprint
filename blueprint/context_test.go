package blueprint

import (
	"bytes"
	"testing"
)

type fooModule struct {
	properties struct {
		Foo string
	}
}

var FooModule = MakeModuleType("FooModule", newFooModule)

func newFooModule() (Module, interface{}) {
	m := &fooModule{}
	return m, &m.properties
}

func (f *fooModule) GenerateBuildActions(ModuleContext) {
}

func (f *fooModule) Foo() string {
	return f.properties.Foo
}

type barModule struct {
	properties struct {
		Bar bool
	}
}

var BarModule = MakeModuleType("BarModule", newBarModule)

func newBarModule() (Module, interface{}) {
	m := &barModule{}
	return m, &m.properties
}

func (b *barModule) GenerateBuildActions(ModuleContext) {
}

func (b *barModule) Bar() bool {
	return b.properties.Bar
}

func TestContextParse(t *testing.T) {
	ctx := NewContext()
	ctx.RegisterModuleType("foo_module", FooModule)
	ctx.RegisterModuleType("bar_module", BarModule)

	r := bytes.NewBufferString(`
		foo_module {
			name: "MyFooModule",
			deps: ["MyBarModule"],
		}

		bar_module {
			name: "MyBarModule",
		}
	`)

	_, errs := ctx.Parse(".", "Blueprint", r)
	if len(errs) > 0 {
		t.Errorf("unexpected parse errors:")
		for _, err := range errs {
			t.Errorf("  %s", err)
		}
		t.FailNow()
	}

	errs = ctx.resolveDependencies()
	if len(errs) > 0 {
		t.Errorf("unexpected dep errors:")
		for _, err := range errs {
			t.Errorf("  %s", err)
		}
		t.FailNow()
	}

	errs = ctx.checkForDependencyCycles()
	if len(errs) > 0 {
		t.Errorf("unexpected dep cycle errors:")
		for _, err := range errs {
			t.Errorf("  %s", err)
		}
		t.FailNow()
	}

}
