package store

import (
	"os"
	"path"
	"testing"
	"time"
)

// Testing for invariant behaviour that cannot be fully captured by the interface definition.

func Test_AllKvStoreImplementations_CaptureExpectedInvariantBehaviour(t *testing.T) {
	t.Run("InMemHashMapKVStorage", func(t *testing.T) {
		test_KvStoreImplementation_CapturesExpectedInvariantBehaviour(t, func() KvStore {
			return NewInMemHashMapKVStorage()
		})
	})
	
	t.Run("InMemSortedKVStorage", func(t *testing.T) {
		test_KvStoreImplementation_CapturesExpectedInvariantBehaviour(t, func() KvStore {
			return NewInMemSortedKVStorage()
		})
	})

	t.Run("FsAppendOnlyStorage", func(t *testing.T) {
		// Setup a uniqueish filename and remove it if necessary
		filename := path.Join(os.TempDir(), "kvstore_test_" + time.Now().Format(time.RFC3339Nano))
		defer os.Remove(filename)
		test_KvStoreImplementation_CapturesExpectedInvariantBehaviour(t, func() KvStore {
			return NewFsAppendOnlyStorage(filename)
		})
	})
}

func test_KvStoreImplementation_CapturesExpectedInvariantBehaviour(t *testing.T, storeFactory func() KvStore) {

	store := storeFactory()

	test_key := "key"
	test_value := "value"
	test_updated_value := "updated_value"

	t.Run("GetBeforeSet", func(t *testing.T) {
		_, exists, err := store.Get(test_key)
		if err != nil {
			t.Errorf("Retrieving a key before it has been set returned an error value : %v", err)
		} else if exists {
			t.Error("Retrieving a key before it has been set returned exists==true")
		}
	})

	store.Set(test_key, test_value)

	t.Run("GetAfterSet", func(t *testing.T) {
		result, exists, err := store.Get(test_key)
		if err != nil {
			t.Errorf("Retrieving a key after it has been set once returned an error value : %v", err)
		} else if !exists {
			t.Error("Retrieving a key after it has been set once returned exists==false")
		} else if result != test_value {
			t.Errorf("Retrieving a key after it has been set once returned an incorrect value '%s', expected '%s'", result, test_value)
		}
	})

	store.Set(test_key, test_updated_value)

	t.Run("GetAfterUpdate", func(t *testing.T) {
		result, exists, err := store.Get(test_key)
		if err != nil {
			t.Errorf("Retrieving a key after it has been updated returned an error value : %v", err)
		} else if !exists {
			t.Error("Retrieving a key after it has been updated returned exists==false")
		} else if result != test_updated_value {
			t.Errorf("Retrieving a key after it has been updated returned an incorrect value '%s', expected '%s'", result, test_value)
		}
	})
}
