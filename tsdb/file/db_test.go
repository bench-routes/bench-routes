package file

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestChainCreation(t *testing.T) {
	set := NewChainSet(FlushAsTime, time.Second)
	set.Run()

	// Create new chain.
	app, utils := set.NewChain("google", "https://google.co.in", false)

	// Add new blocks into the chain.
	app.Append(NewBlock("test", "1234|4321"))
	app.Append(NewBlock("test", "134|421"))
	app.Append(NewBlock("test", "124|431"))
	app.Append(NewBlock("test", "124|321"))

	// Sleep so that Run() commits the chain.
	time.Sleep(time.Second * 2)

	err := os.Remove(utils.Path())
	require.NoError(t, err)
	err = os.Remove(storagePrefix)
	require.NoError(t, err)
}

func TestLoadingExistingChain(t *testing.T) {
	set := NewChainSet(FlushAsTime, time.Second)
	set.Run()

	// Create new chain.
	_, utils := set.NewChain("google", "https://google.co.in", true)
	stream := utils.Stream()
	fmt.Println(stream)
	require.Equal(t, 4, len(stream))
}
