package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"github.com/containers/storage/types"
)

// ParseIDMapping takes idmappings and subuid and subgid maps and returns a storage mapping
func ParseIDMapping(UIDMapSlice, GIDMapSlice []string, subUIDMap, subGIDMap string) (*types.IDMappingOptions, error) {
	return types.ParseIDMapping(UIDMapSlice, GIDMapSlice, subUIDMap, subGIDMap)
}

// DefaultStoreOptions returns the default storage options for containers
func DefaultStoreOptions() (types.StoreOptions, error) {
	return types.DefaultStoreOptions()
}

func validateMountOptions(mountOptions []string) error {
	var Empty struct{}
	// Add invalid options for ImageMount() here.
	invalidOptions := map[string]struct{}{
		"rw": Empty,
	}

	for _, opt := range mountOptions {
		if _, ok := invalidOptions[opt]; ok {
			return fmt.Errorf(" %q option not supported", opt)
		}
	}
	return nil
}

func applyNameOperation(oldNames []string, opParameters []string, op updateNameOperation) ([]string, error) {
	var result []string
	switch op {
	case setNames:
		// ignore all old names and just return new names
		result = opParameters
	case removeNames:
		// remove given names from old names
		result = make([]string, 0, len(oldNames))
		for _, name := range oldNames {
			if !slices.Contains(opParameters, name) {
				result = append(result, name)
			}
		}
	case addNames:
		result = slices.Concat(opParameters, oldNames)
	default:
		return result, errInvalidUpdateNameOperation
	}
	return dedupeStrings(result), nil
}

func moveToTrash(source, trashPath string) error {
	trashPath, err := os.MkdirTemp(trashPath, "")
	if err != nil {
		return fmt.Errorf("creating temp dir in %q: %w", trashPath, err)
	}
	if err := os.Rename(source, filepath.Join(trashPath, filepath.Base(source))); err != nil {
		return fmt.Errorf("moving %q to %q: %w", source, trashPath, err)
	}
	return nil
}

// containsIncompleteFlag returns true if map contains an incompleteFlag set to true
func containsIncompleteFlag(f map[string]interface{}) bool {
	if f == nil {
		return false
	}
	if flagValue, ok := f[incompleteFlag]; ok {
		if b, ok := flagValue.(bool); ok && b {
			return true
		}
	}
	return false
}
