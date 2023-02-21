package games

import (
	"net/http"
	"sync"
	"tic-tac-toe/pkg/websocket"

	"github.com/google/uuid"
)

type Game struct {
	ID       uuid.UUID `json:"id"`
	Capacity uint      `json:"capacity"`

	cli *websocket.Client
}

func New(capacity uint) *Game {
	g := &Game{
		ID:       uuid.New(),
		Capacity: capacity,
	}
	return g
}

func (g *Game) Close() {
	if g.cli == nil {
		return
	}

	g.cli.Close()
}

func (g *Game) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if g.cli == nil {
		g.cli = websocket.NewClient(g.Capacity)
	}
	g.cli.ServeHTTP(w, r)
}

// Broker is responsible of delegating the creation of a new Jam and the
// management of the Jam's websocket clients.
type Broker websocket.Broker[uuid.UUID, *Game]

type broker struct {
	m sync.Map
}

func NewBroker() Broker {
	b := &broker{}
	return b
}

// Delete deletes a jam from the broker.
func (b *broker) Delete(id uuid.UUID) {
	b.m.Delete(id)
}

func (b *broker) LoadAndDelete(id uuid.UUID) (value *Game, loaded bool) {
	actual, loaded := b.m.LoadAndDelete(id)
	if !loaded {
		return nil, loaded
	}

	return actual.(*Game), loaded
}

// Load loads an existing jam from the broker.
func (b *broker) Load(id uuid.UUID) (value *Game, ok bool) {
	v, ok := b.m.Load(id)
	if !ok {
		return nil, ok
	}
	return v.(*Game), ok
}

// Store stores the jam in the broker.
func (b *broker) Store(id uuid.UUID, g *Game) {
	b.m.Store(id, g)
}

// LoadOrStore returns the existing jam for the id if present.
// Otherwise, it stores and returns the given jam.
// The loaded result is true if the value was loaded, false if stored.
func (b *broker) LoadOrStore(id uuid.UUID, g *Game) (*Game, bool) {
	actual, loaded := b.m.LoadOrStore(id, g)
	if !loaded {
		actual = g
	}

	return actual.(*Game), loaded
}
