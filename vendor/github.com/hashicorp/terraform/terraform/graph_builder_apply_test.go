package terraform

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/addrs"
	"github.com/hashicorp/terraform/plans"
	"github.com/hashicorp/terraform/states"
	"github.com/zclconf/go-cty/cty"
)

func TestApplyGraphBuilder_impl(t *testing.T) {
	var _ GraphBuilder = new(ApplyGraphBuilder)
}

func TestApplyGraphBuilder(t *testing.T) {
	changes := &plans.Changes{
		Resources: []*plans.ResourceInstanceChangeSrc{
			{
				Addr: mustResourceInstanceAddr("test_object.create"),
				ChangeSrc: plans.ChangeSrc{
					Action: plans.Create,
				},
			},
			{
				Addr: mustResourceInstanceAddr("test_object.other"),
				ChangeSrc: plans.ChangeSrc{
					Action: plans.Update,
				},
			},
			{
				Addr: mustResourceInstanceAddr("module.child.test_object.create"),
				ChangeSrc: plans.ChangeSrc{
					Action: plans.Create,
				},
			},
			{
				Addr: mustResourceInstanceAddr("module.child.test_object.other"),
				ChangeSrc: plans.ChangeSrc{
					Action: plans.Create,
				},
			},
		},
	}

	b := &ApplyGraphBuilder{
		Config:     testModule(t, "graph-builder-apply-basic"),
		Changes:    changes,
		Components: simpleMockComponentFactory(),
		Schemas:    simpleTestSchemas(),
	}

	g, err := b.Build(addrs.RootModuleInstance)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if g.Path.String() != addrs.RootModuleInstance.String() {
		t.Fatalf("wrong path %q", g.Path.String())
	}

	actual := strings.TrimSpace(g.String())

	expected := strings.TrimSpace(testApplyGraphBuilderStr)
	if actual != expected {
		t.Fatalf("wrong result\n\ngot:\n%s\n\nwant:\n%s", actual, expected)
	}
}

// This tests the ordering of two resources where a non-CBD depends
// on a CBD. GH-11349.
func TestApplyGraphBuilder_depCbd(t *testing.T) {
	changes := &plans.Changes{
		Resources: []*plans.ResourceInstanceChangeSrc{
			{
				Addr: mustResourceInstanceAddr("test_object.A"),
				ChangeSrc: plans.ChangeSrc{
					Action: plans.CreateThenDelete,
				},
			},
			{
				Addr: mustResourceInstanceAddr("test_object.B"),
				ChangeSrc: plans.ChangeSrc{
					Action: plans.Update,
				},
			},
		},
	}

	state := states.NewState()
	root := state.EnsureModule(addrs.RootModuleInstance)
	root.SetResourceInstanceCurrent(
		mustResourceInstanceAddr("test_object.A").Resource,
		&states.ResourceInstanceObjectSrc{
			Status:    states.ObjectReady,
			AttrsJSON: []byte(`{"id":"A"}`),
		},
		mustProviderConfig("provider.test"),
	)
	root.SetResourceInstanceCurrent(
		mustResourceInstanceAddr("test_object.B").Resource,
		&states.ResourceInstanceObjectSrc{
			Status:       states.ObjectReady,
			AttrsJSON:    []byte(`{"id":"B","test_list":["x"]}`),
			Dependencies: []addrs.AbsResource{mustResourceAddr("test_object.A")},
		},
		mustProviderConfig("provider.test"),
	)

	b := &ApplyGraphBuilder{
		Config:     testModule(t, "graph-builder-apply-dep-cbd"),
		Changes:    changes,
		Components: simpleMockComponentFactory(),
		Schemas:    simpleTestSchemas(),
		State:      state,
	}

	g, err := b.Build(addrs.RootModuleInstance)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if g.Path.String() != addrs.RootModuleInstance.String() {
		t.Fatalf("wrong path %q", g.Path.String())
	}

	// We're going to go hunting for our deposed instance node here, so we
	// can find out its key to use in the assertions below.
	var dk states.DeposedKey
	for _, v := range g.Vertices() {
		tv, ok := v.(*NodeDestroyDeposedResourceInstanceObject)
		if !ok {
			continue
		}
		if dk != states.NotDeposed {
			t.Fatalf("more than one deposed instance node in the graph; want only one")
		}
		dk = tv.DeposedKey
	}
	if dk == states.NotDeposed {
		t.Fatalf("no deposed instance node in the graph; want one")
	}

	destroyName := fmt.Sprintf("test_object.A (destroy deposed %s)", dk)

	// Create A, Modify B, Destroy A
	testGraphHappensBefore(
		t, g,
		"test_object.A",
		destroyName,
	)
	testGraphHappensBefore(
		t, g,
		"test_object.A",
		"test_object.B",
	)
	testGraphHappensBefore(
		t, g,
		"test_object.B",
		destroyName,
	)
}

// This tests the ordering of two resources that are both CBD that
// require destroy/create.
func TestApplyGraphBuilder_doubleCBD(t *testing.T) {
	changes := &plans.Changes{
		Resources: []*plans.ResourceInstanceChangeSrc{
			{
				Addr: mustResourceInstanceAddr("test_object.A"),
				ChangeSrc: plans.ChangeSrc{
					Action: plans.CreateThenDelete,
				},
			},
			{
				Addr: mustResourceInstanceAddr("test_object.B"),
				ChangeSrc: plans.ChangeSrc{
					Action: plans.CreateThenDelete,
				},
			},
		},
	}

	b := &ApplyGraphBuilder{
		Config:        testModule(t, "graph-builder-apply-double-cbd"),
		Changes:       changes,
		Components:    simpleMockComponentFactory(),
		Schemas:       simpleTestSchemas(),
		DisableReduce: true,
	}

	g, err := b.Build(addrs.RootModuleInstance)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if g.Path.String() != addrs.RootModuleInstance.String() {
		t.Fatalf("wrong path %q", g.Path.String())
	}

	// We're going to go hunting for our deposed instance node here, so we
	// can find out its key to use in the assertions below.
	var destroyA, destroyB string
	for _, v := range g.Vertices() {
		tv, ok := v.(*NodeDestroyDeposedResourceInstanceObject)
		if !ok {
			continue
		}

		switch tv.Addr.Resource.Name {
		case "A":
			destroyA = fmt.Sprintf("test_object.A (destroy deposed %s)", tv.DeposedKey)
		case "B":
			destroyB = fmt.Sprintf("test_object.B (destroy deposed %s)", tv.DeposedKey)
		default:
			t.Fatalf("unknown instance: %s", tv.Addr)
		}
	}

	// Create A, Modify B, Destroy A
	testGraphHappensBefore(
		t, g,
		"test_object.A",
		destroyA,
	)
	testGraphHappensBefore(
		t, g,
		"test_object.A",
		"test_object.B",
	)
	testGraphHappensBefore(
		t, g,
		"test_object.B",
		destroyB,
	)

	// actual := strings.TrimSpace(g.String())
	// expected := strings.TrimSpace(testApplyGraphBuilderDoubleCBDStr)
	// if actual != expected {
	// 	t.Fatalf("wrong result\n\ngot:\n%s\n\nwant:\n%s", actual, expected)
	// }
}

// This tests the ordering of two resources being destroyed that depend
// on each other from only state. GH-11749
func TestApplyGraphBuilder_destroyStateOnly(t *testing.T) {
	changes := &plans.Changes{
		Resources: []*plans.ResourceInstanceChangeSrc{
			{
				Addr: mustResourceInstanceAddr("module.child.test_object.A"),
				ChangeSrc: plans.ChangeSrc{
					Action: plans.Delete,
				},
			},
			{
				Addr: mustResourceInstanceAddr("module.child.test_object.B"),
				ChangeSrc: plans.ChangeSrc{
					Action: plans.Delete,
				},
			},
		},
	}

	state := MustShimLegacyState(&State{
		Modules: []*ModuleState{
			&ModuleState{
				Path: []string{"root", "child"},
				Resources: map[string]*ResourceState{
					"test_object.A": &ResourceState{
						Type: "test_object",
						Primary: &InstanceState{
							ID:         "foo",
							Attributes: map[string]string{},
						},
						Provider: "provider.test",
					},

					"test_object.B": &ResourceState{
						Type: "test_object",
						Primary: &InstanceState{
							ID:         "bar",
							Attributes: map[string]string{},
						},
						Dependencies: []string{"test_object.A"},
						Provider:     "provider.test",
					},
				},
			},
		},
	})

	b := &ApplyGraphBuilder{
		Config:        testModule(t, "empty"),
		Changes:       changes,
		State:         state,
		Components:    simpleMockComponentFactory(),
		Schemas:       simpleTestSchemas(),
		DisableReduce: true,
	}

	g, diags := b.Build(addrs.RootModuleInstance)
	if diags.HasErrors() {
		t.Fatalf("err: %s", diags.Err())
	}
	t.Logf("Graph:\n%s", g.String())

	if g.Path.String() != addrs.RootModuleInstance.String() {
		t.Fatalf("wrong path %q", g.Path.String())
	}

	testGraphHappensBefore(
		t, g,
		"module.child.test_object.B (destroy)",
		"module.child.test_object.A (destroy)")
}

// This tests the ordering of destroying a single count of a resource.
func TestApplyGraphBuilder_destroyCount(t *testing.T) {
	changes := &plans.Changes{
		Resources: []*plans.ResourceInstanceChangeSrc{
			{
				Addr: mustResourceInstanceAddr("test_object.A[1]"),
				ChangeSrc: plans.ChangeSrc{
					Action: plans.Delete,
				},
			},
			{
				Addr: mustResourceInstanceAddr("test_object.B"),
				ChangeSrc: plans.ChangeSrc{
					Action: plans.Update,
				},
			},
		},
	}

	b := &ApplyGraphBuilder{
		Config:     testModule(t, "graph-builder-apply-count"),
		Changes:    changes,
		Components: simpleMockComponentFactory(),
		Schemas:    simpleTestSchemas(),
	}

	g, err := b.Build(addrs.RootModuleInstance)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if g.Path.String() != addrs.RootModuleInstance.String() {
		t.Fatalf("wrong module path %q", g.Path)
	}

	actual := strings.TrimSpace(g.String())
	expected := strings.TrimSpace(testApplyGraphBuilderDestroyCountStr)
	if actual != expected {
		t.Fatalf("wrong result\n\ngot:\n%s\n\nwant:\n%s", actual, expected)
	}
}

func TestApplyGraphBuilder_moduleDestroy(t *testing.T) {
	changes := &plans.Changes{
		Resources: []*plans.ResourceInstanceChangeSrc{
			{
				Addr: mustResourceInstanceAddr("module.A.test_object.foo"),
				ChangeSrc: plans.ChangeSrc{
					Action: plans.Delete,
				},
			},
			{
				Addr: mustResourceInstanceAddr("module.B.test_object.foo"),
				ChangeSrc: plans.ChangeSrc{
					Action: plans.Delete,
				},
			},
		},
	}

	b := &ApplyGraphBuilder{
		Config:     testModule(t, "graph-builder-apply-module-destroy"),
		Changes:    changes,
		Components: simpleMockComponentFactory(),
		Schemas:    simpleTestSchemas(),
	}

	g, err := b.Build(addrs.RootModuleInstance)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	testGraphHappensBefore(
		t, g,
		"module.B.test_object.foo (destroy)",
		"module.A.test_object.foo (destroy)",
	)
}

func TestApplyGraphBuilder_provisioner(t *testing.T) {
	changes := &plans.Changes{
		Resources: []*plans.ResourceInstanceChangeSrc{
			{
				Addr: mustResourceInstanceAddr("test_object.foo"),
				ChangeSrc: plans.ChangeSrc{
					Action: plans.Create,
				},
			},
		},
	}

	b := &ApplyGraphBuilder{
		Config:     testModule(t, "graph-builder-apply-provisioner"),
		Changes:    changes,
		Components: simpleMockComponentFactory(),
		Schemas:    simpleTestSchemas(),
	}

	g, err := b.Build(addrs.RootModuleInstance)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	testGraphContains(t, g, "provisioner.test")
	testGraphHappensBefore(
		t, g,
		"provisioner.test",
		"test_object.foo",
	)
}

func TestApplyGraphBuilder_provisionerDestroy(t *testing.T) {
	changes := &plans.Changes{
		Resources: []*plans.ResourceInstanceChangeSrc{
			{
				Addr: mustResourceInstanceAddr("test_object.foo"),
				ChangeSrc: plans.ChangeSrc{
					Action: plans.Delete,
				},
			},
		},
	}

	b := &ApplyGraphBuilder{
		Destroy:    true,
		Config:     testModule(t, "graph-builder-apply-provisioner"),
		Changes:    changes,
		Components: simpleMockComponentFactory(),
		Schemas:    simpleTestSchemas(),
	}

	g, err := b.Build(addrs.RootModuleInstance)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	testGraphContains(t, g, "provisioner.test")
	testGraphHappensBefore(
		t, g,
		"provisioner.test",
		"test_object.foo (destroy)",
	)
}

func TestApplyGraphBuilder_targetModule(t *testing.T) {
	changes := &plans.Changes{
		Resources: []*plans.ResourceInstanceChangeSrc{
			{
				Addr: mustResourceInstanceAddr("test_object.foo"),
				ChangeSrc: plans.ChangeSrc{
					Action: plans.Update,
				},
			},
			{
				Addr: mustResourceInstanceAddr("module.child2.test_object.foo"),
				ChangeSrc: plans.ChangeSrc{
					Action: plans.Update,
				},
			},
		},
	}

	b := &ApplyGraphBuilder{
		Config:     testModule(t, "graph-builder-apply-target-module"),
		Changes:    changes,
		Components: simpleMockComponentFactory(),
		Schemas:    simpleTestSchemas(),
		Targets: []addrs.Targetable{
			addrs.RootModuleInstance.Child("child2", addrs.NoKey),
		},
	}

	g, err := b.Build(addrs.RootModuleInstance)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	testGraphNotContains(t, g, "module.child1.output.instance_id")
}

// Ensure that an update resulting from the removal of a resource happens after
// that resource is destroyed.
func TestApplyGraphBuilder_updateFromOrphan(t *testing.T) {
	schemas := simpleTestSchemas()
	instanceSchema := schemas.Providers["test"].ResourceTypes["test_object"]

	bBefore, _ := plans.NewDynamicValue(
		cty.ObjectVal(map[string]cty.Value{
			"id":          cty.StringVal("b_id"),
			"test_string": cty.StringVal("a_id"),
		}), instanceSchema.ImpliedType())
	bAfter, _ := plans.NewDynamicValue(
		cty.ObjectVal(map[string]cty.Value{
			"id":          cty.StringVal("b_id"),
			"test_string": cty.StringVal("changed"),
		}), instanceSchema.ImpliedType())

	changes := &plans.Changes{
		Resources: []*plans.ResourceInstanceChangeSrc{
			{
				Addr: mustResourceInstanceAddr("test_object.a"),
				ChangeSrc: plans.ChangeSrc{
					Action: plans.Delete,
				},
			},
			{
				Addr: mustResourceInstanceAddr("test_object.b"),
				ChangeSrc: plans.ChangeSrc{
					Action: plans.Update,
					Before: bBefore,
					After:  bAfter,
				},
			},
		},
	}

	state := states.NewState()
	root := state.EnsureModule(addrs.RootModuleInstance)
	root.SetResourceInstanceCurrent(
		addrs.Resource{
			Mode: addrs.ManagedResourceMode,
			Type: "test_object",
			Name: "a",
		}.Instance(addrs.NoKey),
		&states.ResourceInstanceObjectSrc{
			Status:    states.ObjectReady,
			AttrsJSON: []byte(`{"id":"a_id"}`),
		},
		addrs.ProviderConfig{
			Type: "test",
		}.Absolute(addrs.RootModuleInstance),
	)
	root.SetResourceInstanceCurrent(
		addrs.Resource{
			Mode: addrs.ManagedResourceMode,
			Type: "test_object",
			Name: "b",
		}.Instance(addrs.NoKey),
		&states.ResourceInstanceObjectSrc{
			Status:    states.ObjectReady,
			AttrsJSON: []byte(`{"id":"b_id","test_string":"a_id"}`),
			Dependencies: []addrs.AbsResource{
				addrs.AbsResource{
					Resource: addrs.Resource{
						Mode: addrs.ManagedResourceMode,
						Type: "test_object",
						Name: "a",
					},
					Module: root.Addr,
				},
			},
		},
		addrs.ProviderConfig{
			Type: "test",
		}.Absolute(addrs.RootModuleInstance),
	)

	b := &ApplyGraphBuilder{
		Config:     testModule(t, "graph-builder-apply-orphan-update"),
		Changes:    changes,
		Components: simpleMockComponentFactory(),
		Schemas:    schemas,
		State:      state,
	}

	g, err := b.Build(addrs.RootModuleInstance)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	expected := strings.TrimSpace(`
test_object.a (destroy)
test_object.b
  test_object.a (destroy)
`)

	instanceGraph := filterInstances(g)
	got := strings.TrimSpace(instanceGraph.String())

	if got != expected {
		t.Fatalf("expected:\n%s\ngot:\n%s", expected, got)
	}
}

const testApplyGraphBuilderStr = `
meta.count-boundary (EachMode fixup)
  module.child.test_object.other
  test_object.other
module.child.test_object.create
  module.child.test_object.create (prepare state)
module.child.test_object.create (prepare state)
  provider.test
  provisioner.test
module.child.test_object.other
  module.child.test_object.create
  module.child.test_object.other (prepare state)
module.child.test_object.other (prepare state)
  provider.test
provider.test
provider.test (close)
  module.child.test_object.other
  test_object.other
provisioner.test
provisioner.test (close)
  module.child.test_object.create
root
  meta.count-boundary (EachMode fixup)
  provider.test (close)
  provisioner.test (close)
test_object.create
  test_object.create (prepare state)
test_object.create (prepare state)
  provider.test
test_object.other
  test_object.create
  test_object.other (prepare state)
test_object.other (prepare state)
  provider.test
`

const testApplyGraphBuilderDestroyCountStr = `
meta.count-boundary (EachMode fixup)
  test_object.B
provider.test
provider.test (close)
  test_object.B
root
  meta.count-boundary (EachMode fixup)
  provider.test (close)
test_object.A (prepare state)
  provider.test
test_object.A[1] (destroy)
  provider.test
test_object.B
  test_object.A (prepare state)
  test_object.A[1] (destroy)
  test_object.B (prepare state)
test_object.B (prepare state)
  provider.test
`