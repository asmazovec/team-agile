package closer

import (
	"context"
	"fmt"
	"sync"
)

// Dependency presents a service resource, could be released.
// User out of this package should be never able to release Dependency manually.
type Dependency struct {
	releaser Releaser
}

const (
	Ready status = iota
	Released
)

// Closer is a closer pattern for graceful shutdown.
// Closer walks on dependencies and emit release of each resource.
type Closer struct {
	Status status
	mu     sync.Mutex
	g      graph
}

type (
	Releaser func(context.Context) error
	graph    map[*Dependency]map[*Dependency]bool
	status   int
)

// Add instantiate a new dependant resource with releaser and dependencies.
func (c *Closer) Add(r Releaser, ds ...*Dependency) (*Dependency, error) {
	const op = "adding dependency"
	c.mu.Lock()
	defer c.mu.Unlock()

	// Verify dependencies.
	for _, d := range ds {
		if d == nil {
			return nil, fmt.Errorf("%s: dependency is nil", op)
		}
		if _, ok := c.g[d]; !ok {
			return nil, fmt.Errorf("%s: dependency not associated with current canceler", op)
		}
	}

	if c.Status != Ready {
		return nil, fmt.Errorf("%s: canceller actually not ready", op)
	}
	if c.g == nil {
		c.g = make(graph, len(ds))
	}
	from := &Dependency{releaser: r}
	c.g[from] = make(map[*Dependency]bool)
	for _, to := range ds {
		c.g[from][to] = true
	}

	return from, nil
}

func (c *Closer) Close(ctx context.Context) (errs []error, err error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.Status == Released {
		return nil, nil
	}
	c.Status = Released
	for top, ok := c.extractTop(); ok; top, ok = c.extractTop() {
		chError := make(chan error, len(top))
		for _, dep := range top {
			go func(dep *Dependency) {
				if dep.releaser == nil {
					chError <- nil
					return
				}
				chError <- dep.releaser(ctx)
			}(dep)
		}
		for i := 0; i < len(top); i++ {
			select {
			case e := <-chError:
				if e == nil {
					break
				}
				errs = append(errs, e)
			case <-ctx.Done():
				return nil, fmt.Errorf("shutdown cancelled: %v", ctx.Err())
			}
		}
	}
	return errs, nil
}

func (c *Closer) extractTop() (top []*Dependency, ok bool) {
	if len(c.g) == 0 {
		return nil, false
	}
	visited := make(map[*Dependency]bool)
	for from, ends := range c.g {
		for to := range ends {
			if !c.g[from][to] {
				continue
			}
			visited[to] = true
		}
	}
	for n := range c.g {
		if visited[n] {
			continue
		}
		top = append(top, n)
		delete(c.g, n)
	}
	return top, true
}
