package cmd

func ExampleVersion() {
	rootCmd.SetArgs([]string{"version"})
	_ = rootCmd.Execute()
	// output:
	// linux-utils version
	// ---
	// Version: unknown
	// Build User: unknown
	// Built On: unknown
}
