package relabel

import (
	"context"
	"sync"

	"github.com/grafana/alloy/internal/component"
	alloy_relabel "github.com/grafana/alloy/internal/component/common/relabel"
	"github.com/grafana/alloy/internal/component/discovery"
	"github.com/grafana/alloy/internal/featuregate"
	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/model/relabel"
)

func init() {
	component.Register(component.Registration{
		Name:      "discovery.relabel",
		Stability: featuregate.StabilityGenerallyAvailable,
		Args:      Arguments{},
		Exports:   Exports{},

		Build: func(opts component.Options, args component.Arguments) (component.Component, error) {
			return New(opts, args.(Arguments))
		},
	})
}

// Arguments holds values which are used to configure the discovery.relabel component.
type Arguments struct {
	// Targets contains the input 'targets' passed by a service discovery component.
	Targets []discovery.Target `alloy:"targets,attr"`

	TargetsArray *discovery.LabelArray `alloy:"targets_array,attr"`

	// The relabelling rules to apply to each target's label set.
	RelabelConfigs []*alloy_relabel.Config `alloy:"rule,block,optional"`
}

// Exports holds values which are exported by the discovery.relabel component.
type Exports struct {
	Output     []discovery.Target    `alloy:"output,attr"`
	Rules      alloy_relabel.Rules   `alloy:"rules,attr"`
	OutputArry *discovery.LabelArray `alloy:"output_array,attr"`
}

// Component implements the discovery.relabel component.
type Component struct {
	opts component.Options

	mut sync.RWMutex
	rcs []*relabel.Config
}

var _ component.Component = (*Component)(nil)

// New creates a new discovery.relabel component.
func New(o component.Options, args Arguments) (*Component, error) {
	c := &Component{opts: o}

	// Call to Update() to set the output once at the start
	if err := c.Update(args); err != nil {
		return nil, err
	}

	return c, nil
}

// Run implements component.Component.
func (c *Component) Run(ctx context.Context) error {
	<-ctx.Done()
	return nil
}

// Update implements component.Component.
func (c *Component) Update(args component.Arguments) error {
	c.mut.Lock()
	defer c.mut.Unlock()

	newArgs := args.(Arguments)
	if newArgs.TargetsArray == nil {
		c.handleTargets(newArgs)
	} else {
		c.handleTargetsArray(newArgs)
	}

	return nil
}

func (c *Component) handleTargetsArray(newArgs Arguments) {
	targets := &discovery.LabelArray{Lbls: make([]*discovery.Labels, 0)}
	relabelConfigs := alloy_relabel.ComponentToPromRelabelConfigs(newArgs.RelabelConfigs)
	c.rcs = relabelConfigs

	for _, t := range newArgs.TargetsArray.Lbls {
		lset, keep := relabel.Process(t.LabelsCopy(), relabelConfigs...)
		if keep {
			targets.Lbls = append(targets.Lbls, discovery.CapsulePool.FromLabels(lset))
		}
	}

	c.opts.OnStateChange(Exports{
		OutputArry: targets,
		Rules:      newArgs.RelabelConfigs,
	})
}

func (c *Component) handleTargets(newArgs Arguments) {
	targets := make([]discovery.Target, 0, len(newArgs.Targets))
	relabelConfigs := alloy_relabel.ComponentToPromRelabelConfigs(newArgs.RelabelConfigs)
	c.rcs = relabelConfigs

	for _, t := range newArgs.Targets {
		lset := componentMapToPromLabels(t)
		lset, keep := relabel.Process(lset, relabelConfigs...)
		if keep {
			targets = append(targets, promLabelsToComponent(lset))
		}
	}

	c.opts.OnStateChange(Exports{
		Output: targets,
		Rules:  newArgs.RelabelConfigs,
	})
}

func componentMapToPromLabels(ls discovery.Target) labels.Labels {
	res := make([]labels.Label, 0, len(ls))
	for k, v := range ls {
		res = append(res, labels.Label{Name: k, Value: v})
	}

	return res
}

func promLabelsToComponent(ls labels.Labels) discovery.Target {
	res := make(map[string]string, len(ls))
	for _, l := range ls {
		res[l.Name] = l.Value
	}

	return res
}
