package main

import (
	"fmt"
	"slices"
	"sort"
	"strings"

	"github.com/mbordner/aoc2023/common"
	"github.com/mbordner/aoc2023/common/files"
)

func main() {
	modules, cables := getModules("../test2.txt")
	fmt.Println(cables)

	fmt.Println(modules)
	button := modules.Get(ButtonID).(*ButtonModule)
	button.Press()
	fmt.Println(modules)
	button.Press()
	fmt.Println(modules)
	button.Press()
	fmt.Println(modules)
	button.Press()
	fmt.Println(modules)

}

const (
	ButtonID      = "button"
	BroadcasterID = "broadcaster"
)

type PulseType int

func (t PulseType) String() string {
	if t == High {
		return "1"
	}
	return "0"
}

const (
	Low PulseType = iota
	High
)

type ModuleType int

const (
	Button ModuleType = iota
	BroadCaster
	FlipFlop
	Conjunction
	Output
)

type Cables struct {
	queue common.Queue[*Pulse]
}

func (c *Cables) Send(pulse *Pulse) {
	wasEmpty := c.queue.Empty()

	for _, t := range pulse.To {
		p := Pulse{From: pulse.From, Type: pulse.Type, To: []Module{t}}
		fmt.Println(&p)
		c.queue.Enqueue(&p)
	}

	if wasEmpty {
		for !c.queue.Empty() {
			p := *(c.queue.Dequeue())
			p.To[0].Receive(p)
		}
	}
}

func NewCables() *Cables {
	cables := new(Cables)
	cables.queue = make(common.Queue[*Pulse], 0, 100)
	return cables
}

type Pulse struct {
	From Module
	To   []Module
	Type PulseType
}

func (p *Pulse) String() string {
	to := make([]string, len(p.To))
	for i, m := range p.To {
		to[i] = m.ID()
	}
	return fmt.Sprintf("[%s -(%s)-> %s]", p.From.ID(), p.Type, strings.Join(to, ","))
}

type Modules struct {
	modulesMap map[string]Module
	modules    []Module
}

func (m *Modules) Set(id string, module Module) {
	if _, e := m.modulesMap[id]; !e {
		m.modulesMap[id] = module
		m.modules = append(m.modules, module)
		sort.Slice(m.modules, func(i, j int) bool {
			return m.modules[i].ID() < m.modules[j].ID()
		})
	}
}

func (m *Modules) Get(id string) Module {
	if module, exists := m.modulesMap[id]; exists {
		return module
	}
	return nil
}

func (m *Modules) GetAll() []Module {
	return m.modules
}

func (m *Modules) String() string {
	ms := make([]string, len(m.modules))
	for i, m := range m.modules {
		ms[i] = fmt.Sprintf("{%s:%s}", m.ID(), m.String())
	}
	return fmt.Sprintf("{%s}", strings.Join(ms, ","))
}

func NewModules() *Modules {
	modules := new(Modules)
	modules.modulesMap = make(map[string]Module)
	modules.modules = make([]Module, 0, 10)
	return modules
}

type Module interface {
	ID() string
	Receive(pulse *Pulse)
	Send(pulseType PulseType)
	Type() ModuleType
	DestinationModuleIDs() []string
	String() string
	AddInput(input Module)
	AddDestination(destination Module)
}

type ModuleBase struct {
	id                   string
	cables               *Cables
	moduleType           ModuleType
	destinationModuleIDs []string
	inputModules         []Module
	destinationModules   []Module
}

func (m *ModuleBase) InitBase(cables *Cables, moduleType ModuleType, id string, destinationModuleIDs []string) {
	m.cables = cables
	m.moduleType = moduleType
	m.id = id
	m.destinationModuleIDs = destinationModuleIDs
	m.destinationModules = make([]Module, 0, len(destinationModuleIDs))
	m.inputModules = make([]Module, 0, 10)
}

func (m *ModuleBase) ID() string {
	return m.id
}
func (m *ModuleBase) Type() ModuleType {
	return m.moduleType
}

func (m *ModuleBase) DestinationModuleIDs() []string {
	return m.destinationModuleIDs
}

func (m *ModuleBase) Send(pulseType PulseType) {
	m.cables.Send(&Pulse{From: m, Type: pulseType, To: m.destinationModules})
}

func (m *ModuleBase) Receive(pulse *Pulse) {

}

func (m *ModuleBase) String() string {
	return "{}"
}

func (m *ModuleBase) sortModules(modules []Module) {
	sort.Slice(modules, func(i, j int) bool {
		return modules[i].ID() < modules[j].ID()
	})
}

func (m *ModuleBase) AddInput(input Module) {
	if !slices.Contains(m.inputModules, input) {
		m.inputModules = append(m.inputModules, input)
		m.sortModules(m.inputModules)
	}
}

func (m *ModuleBase) AddDestination(destination Module) {
	if !slices.Contains(m.destinationModules, destination) {
		m.destinationModules = append(m.destinationModules, destination)
	}
}

type OutputModule struct {
	ModuleBase
}

func NewOutputModule(cables *Cables, id string) *OutputModule {
	output := new(OutputModule)
	output.InitBase(cables, Output, id, []string{})
	return output
}

type ButtonModule struct {
	ModuleBase
}

func (b *ButtonModule) Press() {
	b.Send(Low)
}

func NewButtonModule(cables *Cables) *ButtonModule {
	button := new(ButtonModule)
	button.InitBase(cables, Button, ButtonID, []string{BroadcasterID})
	return button
}

type BroadCasterModule struct {
	ModuleBase
}

func (b *BroadCasterModule) Receive(pulse *Pulse) {
	b.Send(pulse.Type)
}

func NewBroadCasterModule(cables *Cables, destinations []string) *BroadCasterModule {
	broadCaster := new(BroadCasterModule)
	broadCaster.InitBase(cables, BroadCaster, BroadcasterID, destinations)
	return broadCaster
}

type FlipFlopModule struct {
	ModuleBase
	on bool
}

func (f *FlipFlopModule) Receive(pulse *Pulse) {
	if pulse.Type == Low {
		f.on = !f.on
		if f.on {
			f.Send(High)
		} else {
			f.Send(Low)
		}
	}
}

func (f *FlipFlopModule) String() string {
	if f.on {
		return "{1}"
	}
	return "{0}"
}

func NewFlipFlopModule(cables *Cables, id string, destinations []string) *FlipFlopModule {
	flipflop := new(FlipFlopModule)
	flipflop.InitBase(cables, FlipFlop, id, destinations)
	return flipflop
}

type ConjunctionModule struct {
	ModuleBase
	lastPulses map[string]PulseType
}

func (c *ConjunctionModule) AddInput(input Module) {
	c.ModuleBase.AddInput(input)
	c.lastPulses[input.ID()] = Low
}

func (c *ConjunctionModule) Receive(pulse *Pulse) {
	c.lastPulses[pulse.From.ID()] = pulse.Type
	allHigh := true
	for _, pt := range c.lastPulses {
		if pt == Low {
			allHigh = false
			break
		}
	}
	if allHigh {
		c.Send(Low)
	} else {
		c.Send(High)
	}
}

func (c *ConjunctionModule) String() string {
	lps := make([]string, 0, len(c.lastPulses))
	for _, input := range c.ModuleBase.inputModules {
		id := input.ID()
		lps = append(lps, fmt.Sprintf("%s:%s", id, c.lastPulses[id]))
	}
	return fmt.Sprintf("{%s}", strings.Join(lps, ","))
}

func NewConjunctionModule(cables *Cables, id string, destinations []string) *ConjunctionModule {
	conjunction := new(ConjunctionModule)
	conjunction.InitBase(cables, Conjunction, id, destinations)
	conjunction.lastPulses = make(map[string]PulseType)
	return conjunction
}

func getModules(filename string) (*Modules, *Cables) {
	cables := NewCables()
	modules := NewModules()
	button := NewButtonModule(cables)
	modules.Set(button.ID(), button)
	lines := files.MustGetLines(filename)
	for _, line := range lines {
		tokens := strings.Split(line, " -> ")
		destinations := strings.Split(tokens[1], ", ")
		if tokens[0] == BroadcasterID {
			broadcaster := NewBroadCasterModule(cables, destinations)
			modules.Set(broadcaster.ID(), broadcaster)
		} else {
			var module Module
			if tokens[0][0] == '%' {
				module = NewFlipFlopModule(cables, tokens[0][1:], destinations)
			} else if tokens[0][0] == '&' {
				module = NewConjunctionModule(cables, tokens[0][1:], destinations)
			} else {
				panic("invalid module id")
			}
			modules.Set(module.ID(), module)
		}
	}
	for _, module := range modules.GetAll() {
		for _, destination := range module.DestinationModuleIDs() {
			d := modules.Get(destination)
			if d == nil {
				d = NewOutputModule(cables, destination)
				modules.Set(d.ID(), d)
			}
			d.AddInput(module)
			module.AddDestination(d)
		}
	}
	return modules, cables
}
