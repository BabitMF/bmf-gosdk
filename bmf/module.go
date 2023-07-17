package bmf

type Module interface {
	Init() error
	Reset() error
	Process(task *Task) error
	Close() error
	GetModuleInfo() (interface{}, error)
	NeedHungryCheck(istreamId int32) (bool, error)
	IsHungry(istreamId int32) (bool, error)
	IsInfinity() (bool, error)
}
