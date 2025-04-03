package framework

// experimental features on branch dev

type Fwm struct {
	// Fwm is a struct that contains the framework's configuration and state.
	State int
	Conf  string
}

func (fmw Fwm) GetState() int {
	return fmw.State
}
