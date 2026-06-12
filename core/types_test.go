package core

import (
	"reflect"
	"testing"
	"time"
)

func TestInterfaceMetric(t *testing.T) {

	t.Run("Metric doit contenir exactement 5 champs", func(t *testing.T) {
		nbChamps := reflect.TypeOf(Metric{}).NumField()
		if nbChamps != 5 {
			t.Errorf("Le contrat a changé ! Attendu 5 champs, obtenu %d", nbChamps)
		}
	})

	t.Run("Metric.Timestamp doit être de type 'time.Time'", func(t *testing.T) {
		typeField := reflect.TypeOf(Metric{}.Timestamp)
		if typeField != reflect.TypeOf(time.Time{}) {
			t.Errorf("Attendu time.Time, obtenu %v", typeField)
		}
	})

	t.Run("Metric.Timestamp doit être de type 'time.Time'", func(t *testing.T) {
		typeField := reflect.TypeOf(Metric{}.Timestamp)
		if typeField != reflect.TypeOf(time.Time{}) {
			t.Errorf("Attendu time.Time, obtenu %v", typeField)
		}
	})

	t.Run("Metric.ID doit être de type 'string'", func(t *testing.T) {
		typeField := reflect.TypeOf(Metric{}.ID)
		if typeField != reflect.TypeOf("") {
			t.Errorf("Attendu string, obtenu %v", typeField)
		}
	})

	t.Run("Metric.Value doit être de type 'string'", func(t *testing.T) {
		typeField := reflect.TypeOf(Metric{}.Value)
		if typeField != reflect.TypeOf("") {
			t.Errorf("Attendu string, obtenu %v", typeField)
		}
	})

	t.Run("Metric.Format doit être de type 'string'", func(t *testing.T) {
		typeField := reflect.TypeOf(Metric{}.Format)
		if typeField != reflect.TypeOf("") {
			t.Errorf("Attendu string, obtenu %v", typeField)
		}
	})

	t.Run("Metric.Unit doit être de type string", func(t *testing.T) {
		typeField := reflect.TypeOf(Metric{}.Unit)
		if typeField != reflect.TypeOf("") {
			t.Errorf("Attendu string, obtenu %v", typeField)
		}
	})

}
