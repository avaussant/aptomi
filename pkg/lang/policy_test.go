package lang

import (
	"github.com/Aptomi/aptomi/pkg/runtime"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func TestPolicy_AddObjectAndGetObjectsByKind(t *testing.T) {
	namespace, policy := makePolicyWithObjects()

	// retrieve objects
	for _, kind := range []string{ServiceObject.Kind, ContractObject.Kind, RuleObject.Kind, DependencyObject.Kind} {
		assert.Equal(t, 10, len(policy.GetObjectsByKind(kind)), "Number of '%s' objects in the policy should be correct", kind)

		for i := 0; i < 10; i++ {
			name := kind + strconv.Itoa(i)
			getObject(t, policy, kind, name, namespace)
		}
	}

	for _, kind := range []string{ClusterObject.Kind} {
		assert.Equal(t, 10, len(policy.GetObjectsByKind(kind)), "Number of '%s' objects in the policy should be correct", kind)

		for i := 0; i < 10; i++ {
			name := kind + strconv.Itoa(i)
			getObject(t, policy, kind, name, runtime.SystemNS)
		}
	}
}

func TestPolicy_AddObjectIdempotent(t *testing.T) {
	// create two identical policies
	_, policy := makePolicyWithObjects()
	_, policyUpdated := makePolicyWithObjects()

	// add objects from one to another
	for _, pObjType := range PolicyObjects {
		objects := policy.GetObjectsByKind(pObjType.Kind)
		for _, obj := range objects {
			policyUpdated.AddObject(obj)
		}
	}

	// after addition, policy should stay the same
	for _, pObjType := range PolicyObjects {
		objects := policy.GetObjectsByKind(pObjType.Kind)
		objectsUpdated := policyUpdated.GetObjectsByKind(pObjType.Kind)
		assert.Equal(t, len(objects), len(objectsUpdated), "Policy should stay the same after calling AddObject() on the existing %s", pObjType.Kind)
	}
}

func TestPolicy_RemoveObject(t *testing.T) {
	// create two identical policies
	_, policy := makePolicyWithObjects()
	_, policyUpdated := makePolicyWithObjects()

	// delete objects from the updated policy
	for _, pObjType := range PolicyObjects {
		objects := policy.GetObjectsByKind(pObjType.Kind)
		for _, obj := range objects {
			assert.True(t, policyUpdated.RemoveObject(obj), "RemoveObject() should return true when removing an existing object")
		}
	}

	// after removal, policy should be empty
	for _, pObjType := range PolicyObjects {
		objectsUpdated := policyUpdated.GetObjectsByKind(pObjType.Kind)
		assert.Zero(t, len(objectsUpdated), "Policy should contain 0 %s objects after RemoveObject() is called", pObjType.Kind)
	}

	// try to delete objects once again from the empty policy
	for _, pObjType := range PolicyObjects {
		objects := policy.GetObjectsByKind(pObjType.Kind)
		for _, obj := range objects {
			assert.False(t, policyUpdated.RemoveObject(obj), "RemoveObject() should return false when removing a non-existing object")
		}
	}
}

func getObject(t *testing.T, policy *Policy, kind string, name string, namespace string) {
	// get within current namespace
	obj1, err := policy.GetObject(kind, name, namespace)
	assert.NoError(t, err, "Get object '%s/%s' should be successful", kind, name)
	assert.NotNil(t, obj1, "Get object '%s/%s' should return an object", kind, name)

	// get by absolute path
	obj2, err := policy.GetObject(kind, namespace+"/"+name, "")
	assert.NoError(t, err, "Get object '%s/%s/%s' should be successful", namespace, kind, name)
	assert.NotNil(t, obj2, "Get object '%s/%s/%s' should return an object", namespace, kind, name)

	// get by incorrect path (empty)
	obj3, err := policy.GetObject(kind, "", "")
	assert.Error(t, err, "Get object with incorrect locator (zero parts) should return an error")
	assert.Nil(t, obj3)

	// get by incorrect path (too many parts)
	obj4, err := policy.GetObject(kind, "extrapart"+"/"+namespace+"/"+name, "")
	assert.Error(t, err, "Get object with incorrect locator (too many parts) should return an error")
	assert.Nil(t, obj4)

	// get by incorrect namespace
	obj5, err := policy.GetObject(kind, name, "non-existing-namespace")
	assert.Error(t, err, "Get object with a non-existing namespace should return an error")
	assert.Nil(t, obj5)
}

func makePolicyWithObjects() (string, *Policy) {
	namespace := "main"
	policy := NewPolicy()
	for i := 0; i < 10; i++ {
		addObject(policy, &Service{
			TypeKind: ServiceObject.GetTypeKind(),
			Metadata: Metadata{
				Namespace: namespace,
				Name:      "service" + strconv.Itoa(i),
			},
		})
		addObject(policy, &Contract{
			TypeKind: ContractObject.GetTypeKind(),
			Metadata: Metadata{
				Namespace: namespace,
				Name:      "contract" + strconv.Itoa(i),
			},
		})
		addObject(policy, &Cluster{
			TypeKind: ClusterObject.GetTypeKind(),
			Metadata: Metadata{
				Namespace: runtime.SystemNS,
				Name:      "cluster" + strconv.Itoa(i),
			},
			Type: "kubernetes",
		})
		addObject(policy, &Rule{
			TypeKind: RuleObject.GetTypeKind(),
			Metadata: Metadata{
				Namespace: namespace,
				Name:      "rule" + strconv.Itoa(i),
			},
		})
		addObject(policy, &Dependency{
			TypeKind: DependencyObject.GetTypeKind(),
			Metadata: Metadata{
				Namespace: namespace,
				Name:      "dependency" + strconv.Itoa(i),
			},
			User:     "user" + strconv.Itoa(i),
			Contract: "contract" + strconv.Itoa(i),
		})
	}
	return namespace, policy
}

func addObject(policy *Policy, obj Base) {
	err := policy.AddObject(obj)
	if err != nil {
		panic(err)
	}
}
