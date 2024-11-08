package alipanSdk

const (
	driveBackup = 1 << iota
	driveResources
)

var driveMap = map[int]string{
	driveBackup:    "backup",
	driveResources: "resources",
}

const AllDrive = driveBackup | driveResources
