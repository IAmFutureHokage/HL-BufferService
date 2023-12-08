package types

type IcePhenomeniaState byte

const (
	None IcePhenomeniaState = iota
	True
	Empty
)

func (ct IcePhenomeniaState) ToByte() byte {
	return byte(ct)
}

func (ct *IcePhenomeniaState) FromByte(b byte) {
	*ct = IcePhenomeniaState(b)
}

type PrecipitationDuration byte

const (
	No PrecipitationDuration = iota
	Less1
	From1To3
	From3To6
	From6To12
	More12
)

func (ct PrecipitationDuration) ToByte() byte {
	return byte(ct)
}

func (ct *PrecipitationDuration) FromByte(b byte) {
	*ct = PrecipitationDuration(b)
}

type SnowHeight byte

const (
	NoSnow SnowHeight = iota
	Less5
	From5to10
	From11to15
	From16to20
	From21to25
	From26to35
	From36to50
	From51to70
	More70
)

func (ct SnowHeight) ToByte() byte {
	return byte(ct)
}

func (ct *SnowHeight) FromByte(b byte) {
	*ct = SnowHeight(b)
}
