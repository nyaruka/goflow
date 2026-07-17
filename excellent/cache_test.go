package excellent_test

import (
	"fmt"
	"sync"
	"testing"

	"github.com/nyaruka/goflow/excellent"
	"github.com/stretchr/testify/assert"
)

func TestTemplateCache(t *testing.T) {
	cache := excellent.NewTemplateCache()

	t1 := cache.Get(`Hello @contact.name`)
	t2 := cache.Get(`Hello @contact.name`)
	assert.Same(t, t1, t2)
	assert.Equal(t, `Hello @contact.name`, t1.String())

	t3 := cache.Get(`Goodbye @contact.name`)
	assert.NotSame(t, t1, t3)

	// safe for concurrent use (run with -race)
	var wg sync.WaitGroup
	for range 10 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := range 100 {
				parsed := cache.Get(fmt.Sprintf("template @(%d)", i%10))
				assert.NotNil(t, parsed)
			}
		}()
	}
	wg.Wait()
}
