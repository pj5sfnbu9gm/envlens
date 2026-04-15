// Package comparator orchestrates loading and diffing of multiple deployment
// targets against a single baseline environment file.
//
// It ties together the loader, resolver, and differ packages to provide a
// high-level API used by the CLI:
//
//	results, err := comparator.CompareAll("production", targets)
//	if err != nil {
//		log.Fatal(err)
//	}
//	for _, diff := range results {
//		fmt.Printf("--- %s vs %s ---\n", diff.From, diff.To)
//	}
package comparator
