package framework

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	configv1 "github.com/openshift/api/config/v1"
	samplesapi "github.com/openshift/api/samples/v1"
	configv1client "github.com/openshift/client-go/config/clientset/versioned/typed/config/v1"
)

const (
	clusterVersionName = "version"
)

func addCompomentOverride(overrides []configv1.ComponentOverride, override configv1.ComponentOverride) ([]configv1.ComponentOverride, bool) {
	for i, o := range overrides {
		if o.Group == override.Group && o.Kind == override.Kind &&
			o.Namespace == override.Namespace && o.Name == override.Name {
			if overrides[i].Unmanaged == override.Unmanaged {
				return overrides, false
			}
			overrides[i].Unmanaged = override.Unmanaged
			return overrides, true
		}
	}
	return append(overrides, override), true
}

// DisableCVOForOperator disables the samples operator deployment so we can modify the version env
func DisableCVOForOperator(operatorClient *configv1client.ConfigV1Client) error {
	cv, err := operatorClient.ClusterVersions().Get(context.TODO(), clusterVersionName, metav1.GetOptions{})
	if err != nil {
		return err
	}

	var changed bool
	cv.Spec.Overrides, changed = addCompomentOverride(cv.Spec.Overrides, configv1.ComponentOverride{
		Kind:      "Deployment",
		Namespace: samplesapi.OperatorNamespace,
		Name:      "cluster-samples-operator",
		Unmanaged: true,
	})
	if !changed {
		return nil
	}
	if _, err := operatorClient.ClusterVersions().Update(context.TODO(), cv, metav1.UpdateOptions{}); err != nil {
		return err
	}

	return nil
}
