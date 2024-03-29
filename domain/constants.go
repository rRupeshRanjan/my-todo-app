package domain

const (
	AppServerPort        = "app.server.port"
	AppLogLocation       = "app.log.location"
	AppAccessLogLocation = "app.access.log.location"
	FiberLogFormat       = "fiber.log.format"
	FiberLogTimeFormat   = "fiber.log.timeFormat"
	CorsAllowedOrigin    = "app.cors.allowOrigins"
	CorsAllowedHeaders   = "app.cors.allowHeaders"
	SqlDriver            = "sql.driver"
	SqlDatabaseName      = "sql.database.name"
)

var SupportedSearchParams = map[string]string{
	"page":        "0",
	"perPage":     "10",
	"dueByFrom":   "-1",
	"addedOnFrom": "-1",
	"dueByTo":     "9999999999999",
	"addedOnTo":   "9999999999999",
	"id":          "",
	"status":      "",
}
