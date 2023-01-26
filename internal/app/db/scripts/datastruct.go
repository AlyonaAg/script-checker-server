package scripts

const ScriptsTableName = "scripts"

type ListScriptsFilter struct {
	Page 	int64
	Limit 	int64
}