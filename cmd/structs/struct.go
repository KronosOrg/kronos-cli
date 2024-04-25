package structs

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Holiday struct {
	Name string `json:"name"`
	Date string `json:"date"`
}

type IncludedObject struct {
	ApiVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	Namespace  string `json:"namespace"`
	IncludeRef string `json:"includeRef"`
	ExcludeRef string `json:"excludeRef"`
}

type KronosAppSpec struct {
	StartSleep      string           `json:"startSleep"`
	EndSleep        string           `json:"endSleep"`
	WeekDays        string           `json:"weekdays"`
	TimeZone        string           `json:"timezone,omitempty"`
	Holidays        []Holiday        `json:"holidays,omitempty"`
	IncludedObjects []IncludedObject `json:"includedObjects"`
	ForceWake       bool             `json:"forceWake,omitempty"`
	ForceSleep      bool             `json:"forceSleep,omitempty"`
}

type KronosAppStatus struct {
	CreatedSecrets []string `json:"secretCreated,omitempty"`
}

type KronosApp struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KronosAppSpec   `json:"spec,omitempty"`
	Status KronosAppStatus `json:"status,omitempty"`
}
type KronosAppList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KronosApp `json:"items"`
}
