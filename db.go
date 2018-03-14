package prometheusTelemetry

import (
	"time"

	"github.com/jinzhu/gorm"
)

var (
	beginTimer time.Time
)

func beginTransactionMetering(scope *gorm.Scope) {
	beginTimer = time.Now()
	IncrementCurrentDbQueries()
}

func endTransactionMetering(scope *gorm.Scope) {
	DecrementCurrentDbQueries()
	ObserveTransaction(scope.TableName(), beginTimer)
}

// SetupDatabaseTelemetry :
func SetupDatabaseTelemetry(database *gorm.DB) {
	beforeTransactionName := "gorm:begin_transaction"
	afterTransactionName := "gorm:commit_or_rollback_transaction"

	telemetryOnBeforeCallbackName := "telemetry:begin_transaction_metering"
	telemetryOnAfterCallbackName := "telemetry:end_transaction_metering"

	database.Callback().Create().Before(beforeTransactionName).Register(telemetryOnBeforeCallbackName, beginTransactionMetering)
	database.Callback().Create().After(afterTransactionName).Register(telemetryOnAfterCallbackName, endTransactionMetering)

	database.Callback().Delete().Before(beforeTransactionName).Register(telemetryOnBeforeCallbackName, beginTransactionMetering)
	database.Callback().Delete().After(afterTransactionName).Register(telemetryOnAfterCallbackName, endTransactionMetering)

	database.Callback().Query().Before(beforeTransactionName).Register(telemetryOnBeforeCallbackName, beginTransactionMetering)
	database.Callback().Query().After(afterTransactionName).Register(telemetryOnAfterCallbackName, endTransactionMetering)

	database.Callback().RowQuery().Before(beforeTransactionName).Register(telemetryOnBeforeCallbackName, beginTransactionMetering)
	database.Callback().RowQuery().After(afterTransactionName).Register(telemetryOnAfterCallbackName, endTransactionMetering)

	database.Callback().Update().Before(beforeTransactionName).Register(telemetryOnBeforeCallbackName, beginTransactionMetering)
	database.Callback().Update().After(afterTransactionName).Register(telemetryOnAfterCallbackName, endTransactionMetering)
}
