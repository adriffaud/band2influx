package band2influx

// Datapoint is sensor data at a given point in time
type Datapoint struct {
	Timestamp    int64 `db:"TIMESTAMP" json:"timestamp"`
	RawIntensity int   `db:"RAW_INTENSITY" json:"rawIntensity"`
	Steps        int   `db:"STEPS" json:"steps"`
	RawKind      int   `db:"RAW_KIND" json:"rawKind"`
	HeartRate    int   `db:"HEART_RATE" json:"heartRate"`
}
