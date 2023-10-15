package model

type EditType int

const (
	EditTypeSlowed EditType = iota
	EditTypeSpedUp
	EditTypeReverb
	EditTypeReverse
	EditTypePitchShiftDown
	EditTypePitchShiftUp
)

func (e EditType) String() string {
	return [...]string{"slowed", "sped_up", "reverb", "reverse", "pitch_shift_down", "pitch_shift_up"}[e]
}

func (e EditType) Int() int {
	return int(e)
}

func GetEditType(i int) EditType {
	return [...]EditType{EditTypeSlowed, EditTypeSpedUp, EditTypeReverb, EditTypeReverse, EditTypePitchShiftDown, EditTypePitchShiftUp}[i]
}
