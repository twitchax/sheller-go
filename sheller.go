package sheller

// TODO:
//   * Allow user to "suppress" copy of local environment.

// UseExecutable creates a new Executable instance.
func UseExecutable(exe string) *Executable {
	return CreateExecutable(exe)
}
