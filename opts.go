package confetti

type optsLoader struct {
	errOnUnknown bool
}

type optsSSMClientLoader struct {
	client SSMAPI
}

func (o optsLoader) Load(_ any, ownConfig *confetti) (err error) {
	ownConfig.errOnUnknown = o.errOnUnknown
	return
}

func (o optsSSMClientLoader) Load(_ any, ownConfig *confetti) (err error) {
	ownConfig.ssmClient = o.client
	return
}
