package lang

import "reflect"

// LabelCluster is a special label name where cluster should be stored. It's required by the engine during policy processing
const LabelCluster = "cluster"

// LabelSet defines the set of labels that will be manipulated throughout policy execution. All labels are stored
// in a 'key' -> 'value' map
type LabelSet struct {
	Labels map[string]string
}

// NewLabelSet creates a new LabelSet from a given map of text labels
func NewLabelSet(labels map[string]string) *LabelSet {
	result := &LabelSet{Labels: make(map[string]string, len(labels))}
	result.AddLabels(labels)
	return result
}

// AddLabels adds new labels to the current set of labels
func (src *LabelSet) AddLabels(addMap map[string]string) {
	for k, v := range addMap {
		src.Labels[k] = v
	}
}

// ApplyTransform applies a given set of label transformations to the current set of labels.
// The method teturns true if changes have been made to the current set
func (src *LabelSet) ApplyTransform(ops LabelOperations) bool {
	changed := false
	if ops != nil {
		// set labels
		for k, v := range ops["set"] {
			if src.Labels[k] != v {
				src.Labels[k] = v
				changed = true
			}
		}

		// remove labels
		for k := range ops["remove"] {
			if _, exists := src.Labels[k]; exists {
				delete(src.Labels, k)
				changed = true
			}
		}
	}
	return changed
}

// Equal compares two labels sets. If one is nil and another one is empty, it will return true as well
func (src *LabelSet) Equal(dst *LabelSet) bool {
	if len(src.Labels) == 0 && len(dst.Labels) == 0 {
		return true
	}
	return reflect.DeepEqual(src.Labels, dst.Labels)
}
