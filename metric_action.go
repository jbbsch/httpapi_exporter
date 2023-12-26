package main

import (
	//"bytes"
	"fmt"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

// ***************************************************************************************
// ***************************************************************************************
// metric fact
// ***************************************************************************************
// ***************************************************************************************

type MetricAction struct {
	Name    *Field              `yaml:"name,omitempty"`
	With    []any               `yaml:"with,omitempty"`
	When    []*exporterTemplate `yaml:"when,omitempty"`
	LoopVar string              `yaml:"loop_var,omitempty"`
	Vars    map[string]any      `yaml:"vars,omitempty"`
	Until   []*exporterTemplate `yaml:"until,omitempty"`

	mc           *MetricConfig
	metricFamily *MetricFamily
	vars         [][]any

	// Catches all undefined fields and must be empty after parsing.
	XXX map[string]interface{} `yaml:",inline" json:"-"`
}

func (a *MetricAction) Type() int {
	return metric_action
}

func (a *MetricAction) GetName(symtab map[string]any, logger log.Logger) string {
	str, err := a.Name.GetValueString(symtab, nil, false)
	if err != nil {
		level.Warn(logger).Log(
			"collid", CollectorId(symtab, logger),
			"script", ScriptName(symtab, logger),
			"msg", fmt.Sprintf("invalid action name: %v", err))
		return ""
	}
	return str
}

func (a *MetricAction) GetNameField() *Field {
	return a.Name
}
func (a *MetricAction) SetNameField(name *Field) {
	a.Name = name
}

func (a *MetricAction) GetWidh() []any {
	return a.With
}
func (a *MetricAction) SetWidth(with []any) {
	a.With = with
}

func (a *MetricAction) GetWhen() []*exporterTemplate {
	return a.When

}
func (a *MetricAction) SetWhen(when []*exporterTemplate) {
	a.When = when
}

func (a *MetricAction) GetLoopVar() string {
	return a.LoopVar
}
func (a *MetricAction) SetLoopVar(loopvar string) {
	a.LoopVar = loopvar
}

func (a *MetricAction) GetVars() [][]any {
	return a.vars
}
func (a *MetricAction) SetVars(vars [][]any) {
	a.vars = vars
}

func (a *MetricAction) GetUntil() []*exporterTemplate {
	return a.Until
}
func (a *MetricAction) SetUntil(until []*exporterTemplate) {
	a.Until = until
}

func (a *MetricAction) setBasicElement(
	nameField *Field,
	vars [][]any,
	with []any,
	loopVar string,
	when []*exporterTemplate,
	until []*exporterTemplate) error {
	return setBasicElement(a, nameField, vars, with, loopVar, when, until)
}

func (a *MetricAction) PlayAction(script *YAMLScript, symtab map[string]any, logger log.Logger) error {
	return PlayBaseAction(script, symtab, logger, a, a.CustomAction)
}

// only for MetricsAction
func (a *MetricAction) GetMetrics() []*GetMetricsRes {
	return nil
}

// only for MetricAction
func (a *MetricAction) GetMetric() *MetricConfig {
	return a.mc
}

func (a *MetricAction) SetMetricFamily(mf *MetricFamily) {
	a.metricFamily = mf
}

// only for PlayAction
func (a *MetricAction) SetPlayAction(scripts map[string]*YAMLScript) error {
	return nil
}

// specific behavior for the MetricAction

func (a *MetricAction) CustomAction(script *YAMLScript, symtab map[string]any, logger log.Logger) error {
	var (
		metric_channel chan<- Metric
		// mfs            []*MetricFamily
	)
	loop_var_idx := ""
	if raw_loop_var_idx, ok := symtab["loop_var_idx"].(int); ok {
		if raw_loop_var_idx > 0 {
			loop_var_idx = fmt.Sprintf(" %d", raw_loop_var_idx)
		}
	}
	level.Debug(logger).Log(
		"collid", CollectorId(symtab, logger),
		"script", ScriptName(symtab, logger),
		"name", a.GetName(symtab, logger),
		"msg", fmt.Sprintf("[Type: MetricAction] %s", loop_var_idx))

	if r_val, ok := symtab["__metric_channel"]; ok {
		if metric_channel, ok = r_val.(chan<- Metric); !ok {
			panic(fmt.Sprintf("collid=\"%s\" script=\"%s\" name=\"%s\" msg=\"invalid context (channel wrong type)\"",
				CollectorId(symtab, logger),
				ScriptName(symtab, logger),
				a.GetName(symtab, logger)))
		}
	} else {
		panic(fmt.Sprintf("collid=\"%s\" script=\"%s\" name=\"%s\" msg=\"invalid context (channel not set)\"",
			CollectorId(symtab, logger),
			ScriptName(symtab, logger),
			a.GetName(symtab, logger)))
	}

	// for _, mf := range mfs {
	level.Debug(logger).Log(
		"collid", CollectorId(symtab, logger),
		"script", ScriptName(symtab, logger),
		"name", a.GetName(symtab, logger),
		"msg", fmt.Sprintf("    metric_name: %s", a.metricFamily.Name()))
	a.metricFamily.Collect(symtab, logger, metric_channel)
	// }

	return nil
}

func (a *MetricAction) AddCustomTemplate(customTemplate *exporterTemplate) error {

	if err := AddCustomTemplate(a, customTemplate); err != nil {
		return err
	}

	return nil
}
