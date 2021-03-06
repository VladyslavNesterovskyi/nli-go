package central

// information about a name in a database, for disambiguation by the user
type NameInformation struct {
	Name         string
	DatabaseName string
	EntityType   string
	SharedId     string
	Information  string
}

func (nameInformation NameInformation) GetIdentifier() string {
	return nameInformation.DatabaseName + "/" + nameInformation.SharedId
}