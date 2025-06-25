package workflow

import (
	"fmt"
	"reflect"
	"runtime"
)

type Action interface {
	GetName() string
	RunAction() error
}

type ReconcileFunc func() error

type point struct {
	Name   string
	Action ReconcileFunc
}

func (p *point) GetName() string {
	return p.Name
}

func (p *point) RunAction() error {
	return p.Action()
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

func (w *World) RunAction() error {
	fmt.Println(w.GetName())
	if w.PreCondition != nil {
		fmt.Printf("Precondition is running: %v \n", w.PreCondition.GetName())
		w.PreCondition.RunAction()
	}
	if w.Actions != nil {
		fmt.Println("Start to run actions")
		for _, action := range w.Actions {
			fmt.Println(action.GetName())
			action.RunAction()
		}
	}
	if w.PostCondition != nil {
		fmt.Printf("PostCondition is running: %v \n", w.PostCondition.GetName())
		w.PostCondition.RunAction()
	}
	if w.ErrorHandler != nil {
		fmt.Printf("ErrorHandler is running: %v \n", w.ErrorHandler.GetName())
		w.ErrorHandler.RunAction()
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
		fmt.Println("we not good")
		fmt.Println(v)
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
			fmt.Println("we not good")
			fmt.Println(v)
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
			fmt.Println("we not good")
			fmt.Println(v)
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
