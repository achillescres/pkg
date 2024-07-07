package batch

const PostgresMaxParameters = 65535

func PostgresMaxRowsFor(parametersNumber uint) int {
	return (PostgresMaxParameters - 1) / int(parametersNumber)
}
