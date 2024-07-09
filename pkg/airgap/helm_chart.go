package airgap

type HelmChart struct {
	Version string
}

func (HelmChart) Download() (string, error) {
	//TODO implement me
	panic("implement me")
}

func (HelmChart) Verify() error {
	//TODO implement me
	panic("implement me")
}
