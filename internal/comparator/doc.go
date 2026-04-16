// Package comparator loads multiple named environment files and compares
// each one against a designated baseline target, collecting all diff
// results into a map keyed by target name.
//
// Use CompareAll to run the full comparison and HasChanges to quickly
// determine whether any differences were found.
package comparator
