package confetti

type optsLoader struct {
	errOnUnknown bool
}

func (o optsLoader) Load(_ any, ownConfig *confetti) (err error) {
	ownConfig.errOnUnknown = o.errOnUnknown
	return
}
