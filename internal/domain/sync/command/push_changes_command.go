package command

type PushChangesCommand struct {
	AccountUID string
	DeviceUID  string
	BaseCursor string
	Mutations  []Mutation
}
