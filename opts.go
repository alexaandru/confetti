package confetti

type optsLoader struct {
	errOnUnknown bool
}

type optsMockedSSMLoader struct {
	client SSMAPI
}

func (o optsLoader) Load(_ any, ownConfig *confetti) (err error) {
	ownConfig.errOnUnknown = o.errOnUnknown
	return
}

func (o optsMockedSSMLoader) Load(_ any, ownConfig *confetti) (err error) {
	ownConfig.mockedSSM = o.client
	return
}
