package mikro

type Button int

const (
	// Browser Section

	ProjectButton Button = iota
	FavoritesButton
	BrowserButton

	// Edit section

	VolumeButton
	SwingButton
	TempoButton
	PluginButton
	SamplingButton

	ArrowLeftButton
	ArrowRightButton

	// Performance section

	PitchButton
	ModButton
	PerformButton
	NotesButton

	GroupButton
	AutoButton
	LockButton
	NoteRepeatButton

	// Transport section

	RestartButton
	EraseButton
	TapButton
	FollowButton

	PlayButton
	RecButton
	StopButton
	ShiftButton

	// Pads section

	FixedVelButton
	PadModeButton
	KeyboardButton
	ChordsButton
	StepButton

	SceneButton
	PatternButton
	EventsButton
	VariationButton
	DuplicateButton
	SelectButton
	SoloButton
	MuteButton

	EncoderButton
)
