// Package lang provides core constructs for describing Aptomi policy, as well as core structures for processing policy.
//
// Let's start with policy objects.
// Cluster - individual cluster where containers get deployed (e.g. k8s cluster).
// Contract - contract for a service (e.g. database).
// Context - a set of contexts, defining how contract can be fulfilled (e.g. MariaDB, MySQL, SQLite).
// Service - specific service implementation (set of containers to run, and dependencies on other services).
// Dependency - service use declaration, which triggers instantiation of a service .
// Rule - rules which constitute policy, allowing to change labels and perform actions during policy resolution.
// ACLRule - rules which define user roles for accessing Aptomi namespaces.
//
// Now, core structures:
// LabelSet - set of labels that get processed and transformed
// LabelOperations - how to transform labels
// Criteria - bool-based lists of expressions for defining complex criteria
package lang
