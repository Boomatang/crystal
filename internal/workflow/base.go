package workflow

import (
	"fmt"
	"reflect"
	"runtime"

	"github.com/boomatang/crystal/internal/logger"
	"github.com/prometheus/client_golang/prometheus"
)

type Action interface {
	GetName() string
	RunAction(*NodeList) error
}

type ReconcileFunc func(*NodeList) error

type point struct {
	Name   string
	Action ReconcileFunc
}

func (p *point) GetName() string {
	return p.Name
}

func (p *point) RunAction(nodes *NodeList) error {
	return p.Action(nodes)
}

func NewPoint(f ReconcileFunc) *point {
	if reflect.TypeOf(f).Kind() != reflect.Func {
		panic("Not a function")
	}

	var name string
	fn := runtime.FuncForPC(reflect.ValueOf(f).Pointer())
	if fn != nil {
		name = fn.Name()
	} else {
		panic("Unknown function")
	}
	return &point{Name: name, Action: f}
}

type World struct {
	Name          string
	PreCondition  Action
	PostCondition Action
	Actions       []Action
	ErrorHandler  Action
}

func (w *World) RunAction(nodes *NodeList) error {
	timer := prometheus.NewTimer(ActionDuration.WithLabelValues(w.GetName()))
	defer timer.ObserveDuration()

	ActionTotal.WithLabelValues(w.GetName()).Inc()

	logger.Log.Info("workflow started", "workflow", w.GetName())
	if w.PreCondition != nil {
		logger.Log.Debug("executing precondition",
			"workflow", w.GetName(),
			"precondition", w.PreCondition.GetName())
		w.PreCondition.RunAction(nodes)
	}
	if w.Actions != nil {
		logger.Log.Debug("executing actions",
			"workflow", w.GetName(),
			"count", len(w.Actions))
		for _, action := range w.Actions {
			logger.Log.Debug("executing action",
				"workflow", w.GetName(),
				"action", action.GetName())
			action.RunAction(nodes)
		}
	}
	if w.PostCondition != nil {
		logger.Log.Debug("executing postcondition",
			"workflow", w.GetName(),
			"postcondition", w.PostCondition.GetName())
		w.PostCondition.RunAction(nodes)
	}
	if w.ErrorHandler != nil {
		logger.Log.Warn("executing error handler",
			"workflow", w.GetName(),
			"handler", w.ErrorHandler.GetName())
		w.ErrorHandler.RunAction(nodes)
	}
	return nil
}

func (w *World) GetName() string {
	return w.Name
}

func (w *World) AddAction(action Action) {
	w.Actions = append(w.Actions, action)
}

func NewWorld(name string) *World {
	return &World{Name: name}
}

type Graph interface {
	Render() string
}

type WorldGraph struct {
	World *World
}

var counter = 0

func (g *WorldGraph) Render() string {
	pre := make([]string, 0)
	next := make([]string, 0)
	s := ""

	switch v := g.World.PreCondition.(type) {
	case *World:
		world, ok := g.World.PreCondition.(*World)
		if !ok {
			panic("This should never happen")
		}

		subString, subNext := g.subRender(world, pre)
		s = fmt.Sprintf("%v\n%v", s, subString)
		pre = subNext

	case *point:
		p := name(g.World.PreCondition.GetName())
		s = fmt.Sprintf("%v\n%v", s, p)
		pre = append(pre, p)
	default:
		logger.Log.Warn("unknown precondition type", "type", fmt.Sprintf("%T", v))
	}

	for _, action := range g.World.Actions {
		name := name(action.GetName())
		switch v := action.(type) {
		case *World:
			world, ok := action.(*World)
			if !ok {
				panic("This should never happen")
			}

			subString, subNext := g.subRender(world, pre)
			s = fmt.Sprintf("%v\n%v", s, subString)
			next = append(next, subNext...)

		case *point:
			next = append(next, name)
			for _, node := range pre {
				s = fmt.Sprintf("%v\n%v -> %v", s, node, name)
			}

		default:
			logger.Log.Warn("unknown action type", "type", fmt.Sprintf("%T", v))
		}
	}

	if g.World.PostCondition != nil {
		switch v := g.World.PreCondition.(type) {
		case *World:
			world, ok := g.World.PostCondition.(*World)
			if !ok {
				panic("This should never happen")
			}

			subString, subNext := g.subRender(world, next)
			s = fmt.Sprintf("%v\n%v", s, subString)
			next = subNext

		case *point:
			pp := name(g.World.PostCondition.GetName())
			for _, node := range next {
				s = fmt.Sprintf("%v\n%v -> %v", s, node, pp)
			}
		default:
			logger.Log.Warn("unknown postcondition type", "type", fmt.Sprintf("%T", v))
		}
	}

	s = fmt.Sprintf("digraph {%v\n}", s)

	return s
}

func (g *WorldGraph) subRender(world *World, pre []string) (string, []string) {
	next := make([]string, 0)
	s := ""
	w := name(world.GetName())
	for _, node := range pre {
		s = fmt.Sprintf("%v\n%v -> %v", s, node, w)
	}
	pre = []string{w}

	if world.ErrorHandler != nil {
		s = fmt.Sprintf("%v\n%v -> %v", s, w, name(world.ErrorHandler.GetName()))

	}

	p := name(world.PreCondition.GetName())
	for _, node := range pre {
		s = fmt.Sprintf("%v\n%v -> %v", s, node, p)
	}
	pre = []string{p}
	for _, action := range world.Actions {
		name := name(action.GetName())
		next = append(next, name)
		for _, node := range pre {
			s = fmt.Sprintf("%v\n%v -> %v", s, node, name)
		}
	}
	if world.PostCondition != nil {
		pp := name(world.PostCondition.GetName())
		for _, node := range next {
			s = fmt.Sprintf("%v\n%v -> %v", s, node, pp)
		}
		next = []string{pp}
	}

	return s, next
}

func name(n string) string {
	counter++
	return fmt.Sprintf("\"%d.%v\"", counter, n)

}
