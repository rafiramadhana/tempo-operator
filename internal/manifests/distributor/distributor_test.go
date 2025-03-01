package distributor

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8slabels "k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/ptr"

	configv1alpha1 "github.com/grafana/tempo-operator/apis/config/v1alpha1"
	"github.com/grafana/tempo-operator/apis/tempo/v1alpha1"
	"github.com/grafana/tempo-operator/internal/manifests/manifestutils"
)

func TestBuildDistributor(t *testing.T) {

	tests := []struct {
		name                   string
		enableGateway          bool
		enableReceiverTLS      bool
		expectedContainerPorts []corev1.ContainerPort
		expectedServicePorts   []corev1.ServicePort
		expectedResources      corev1.ResourceRequirements
		expectedVolumes        []corev1.Volume
		expectedVolumeMounts   []corev1.VolumeMount
	}{
		{
			name:          "Gateway disabled",
			enableGateway: false,
			expectedContainerPorts: []corev1.ContainerPort{
				{
					Name:          manifestutils.OtlpGrpcPortName,
					ContainerPort: manifestutils.PortOtlpGrpcServer,
					Protocol:      corev1.ProtocolTCP,
				},
				{
					Name:          manifestutils.HttpPortName,
					ContainerPort: manifestutils.PortHTTPServer,
					Protocol:      corev1.ProtocolTCP,
				},
				{
					Name:          manifestutils.HttpMemberlistPortName,
					ContainerPort: manifestutils.PortMemberlist,
					Protocol:      corev1.ProtocolTCP,
				},
				{
					Name:          manifestutils.PortOtlpHttpName,
					ContainerPort: manifestutils.PortOtlpHttp,
					Protocol:      corev1.ProtocolTCP,
				},
				{
					Name:          manifestutils.PortJaegerThriftHTTPName,
					ContainerPort: manifestutils.PortJaegerThriftHTTP,
					Protocol:      corev1.ProtocolTCP,
				},
				{
					Name:          manifestutils.PortJaegerThriftCompactName,
					ContainerPort: manifestutils.PortJaegerThriftCompact,
					Protocol:      corev1.ProtocolUDP,
				},
				{
					Name:          manifestutils.PortJaegerThriftBinaryName,
					ContainerPort: manifestutils.PortJaegerThriftBinary,
					Protocol:      corev1.ProtocolUDP,
				},
				{
					Name:          manifestutils.PortJaegerGrpcName,
					ContainerPort: manifestutils.PortJaegerGrpc,
					Protocol:      corev1.ProtocolTCP,
				},
				{
					Name:          manifestutils.PortZipkinName,
					ContainerPort: manifestutils.PortZipkin,
					Protocol:      corev1.ProtocolTCP,
				},
			},
			expectedServicePorts: []corev1.ServicePort{
				{
					Name:       manifestutils.OtlpGrpcPortName,
					Protocol:   corev1.ProtocolTCP,
					Port:       manifestutils.PortOtlpGrpcServer,
					TargetPort: intstr.FromString(manifestutils.OtlpGrpcPortName),
				},
				{
					Name:       manifestutils.HttpPortName,
					Protocol:   corev1.ProtocolTCP,
					Port:       manifestutils.PortHTTPServer,
					TargetPort: intstr.FromString(manifestutils.HttpPortName),
				},
				{
					Name:       manifestutils.PortOtlpHttpName,
					Port:       manifestutils.PortOtlpHttp,
					TargetPort: intstr.FromString(manifestutils.PortOtlpHttpName),
					Protocol:   corev1.ProtocolTCP,
				},
				{
					Name:       manifestutils.PortJaegerThriftHTTPName,
					Port:       manifestutils.PortJaegerThriftHTTP,
					TargetPort: intstr.FromString(manifestutils.PortJaegerThriftHTTPName),
					Protocol:   corev1.ProtocolTCP,
				},
				{
					Name:       manifestutils.PortJaegerThriftCompactName,
					Port:       manifestutils.PortJaegerThriftCompact,
					TargetPort: intstr.FromString(manifestutils.PortJaegerThriftCompactName),
					Protocol:   corev1.ProtocolUDP,
				},
				{
					Name:       manifestutils.PortJaegerThriftBinaryName,
					Port:       manifestutils.PortJaegerThriftBinary,
					TargetPort: intstr.FromString(manifestutils.PortJaegerThriftBinaryName),
					Protocol:   corev1.ProtocolUDP,
				},
				{
					Name:       manifestutils.PortJaegerGrpcName,
					Port:       manifestutils.PortJaegerGrpc,
					TargetPort: intstr.FromString(manifestutils.PortJaegerGrpcName),
					Protocol:   corev1.ProtocolTCP,
				},
				{
					Name:       manifestutils.PortZipkinName,
					Port:       manifestutils.PortZipkin,
					TargetPort: intstr.FromString(manifestutils.PortZipkinName),
					Protocol:   corev1.ProtocolTCP,
				},
			},
			expectedResources: corev1.ResourceRequirements{
				Limits: corev1.ResourceList{
					corev1.ResourceCPU:    *resource.NewMilliQuantity(270, resource.BinarySI),
					corev1.ResourceMemory: *resource.NewQuantity(257698032, resource.BinarySI),
				},
				Requests: corev1.ResourceList{
					corev1.ResourceCPU:    *resource.NewMilliQuantity(81, resource.BinarySI),
					corev1.ResourceMemory: *resource.NewQuantity(77309416, resource.BinarySI),
				},
			},
			expectedVolumes: []corev1.Volume{
				{
					Name: manifestutils.ConfigVolumeName,
					VolumeSource: corev1.VolumeSource{
						ConfigMap: &corev1.ConfigMapVolumeSource{
							LocalObjectReference: corev1.LocalObjectReference{
								Name: "tempo-test",
							},
						},
					},
				},
				{
					Name: manifestutils.TmpStorageVolumeName,
					VolumeSource: corev1.VolumeSource{
						EmptyDir: &corev1.EmptyDirVolumeSource{},
					},
				},
			},
			expectedVolumeMounts: []corev1.VolumeMount{
				{
					Name:      manifestutils.ConfigVolumeName,
					MountPath: "/conf",
					ReadOnly:  true,
				},
				{
					Name:      manifestutils.TmpStorageVolumeName,
					MountPath: manifestutils.TmpStoragePath,
				},
			},
		},
		{
			name:              "Receiver TLS enable",
			enableGateway:     false,
			enableReceiverTLS: true,
			expectedContainerPorts: []corev1.ContainerPort{
				{
					Name:          manifestutils.OtlpGrpcPortName,
					ContainerPort: manifestutils.PortOtlpGrpcServer,
					Protocol:      corev1.ProtocolTCP,
				},
				{
					Name:          manifestutils.HttpPortName,
					ContainerPort: manifestutils.PortHTTPServer,
					Protocol:      corev1.ProtocolTCP,
				},
				{
					Name:          manifestutils.HttpMemberlistPortName,
					ContainerPort: manifestutils.PortMemberlist,
					Protocol:      corev1.ProtocolTCP,
				},
				{
					Name:          manifestutils.PortOtlpHttpName,
					ContainerPort: manifestutils.PortOtlpHttp,
					Protocol:      corev1.ProtocolTCP,
				},
				{
					Name:          manifestutils.PortJaegerThriftHTTPName,
					ContainerPort: manifestutils.PortJaegerThriftHTTP,
					Protocol:      corev1.ProtocolTCP,
				},
				{
					Name:          manifestutils.PortJaegerThriftCompactName,
					ContainerPort: manifestutils.PortJaegerThriftCompact,
					Protocol:      corev1.ProtocolUDP,
				},
				{
					Name:          manifestutils.PortJaegerThriftBinaryName,
					ContainerPort: manifestutils.PortJaegerThriftBinary,
					Protocol:      corev1.ProtocolUDP,
				},
				{
					Name:          manifestutils.PortJaegerGrpcName,
					ContainerPort: manifestutils.PortJaegerGrpc,
					Protocol:      corev1.ProtocolTCP,
				},
				{
					Name:          manifestutils.PortZipkinName,
					ContainerPort: manifestutils.PortZipkin,
					Protocol:      corev1.ProtocolTCP,
				},
			},
			expectedServicePorts: []corev1.ServicePort{
				{
					Name:       manifestutils.OtlpGrpcPortName,
					Protocol:   corev1.ProtocolTCP,
					Port:       manifestutils.PortOtlpGrpcServer,
					TargetPort: intstr.FromString(manifestutils.OtlpGrpcPortName),
				},
				{
					Name:       manifestutils.HttpPortName,
					Protocol:   corev1.ProtocolTCP,
					Port:       manifestutils.PortHTTPServer,
					TargetPort: intstr.FromString(manifestutils.HttpPortName),
				},
				{
					Name:       manifestutils.PortOtlpHttpName,
					Port:       manifestutils.PortOtlpHttp,
					TargetPort: intstr.FromString(manifestutils.PortOtlpHttpName),
					Protocol:   corev1.ProtocolTCP,
				},
				{
					Name:       manifestutils.PortJaegerThriftHTTPName,
					Port:       manifestutils.PortJaegerThriftHTTP,
					TargetPort: intstr.FromString(manifestutils.PortJaegerThriftHTTPName),
					Protocol:   corev1.ProtocolTCP,
				},
				{
					Name:       manifestutils.PortJaegerThriftCompactName,
					Port:       manifestutils.PortJaegerThriftCompact,
					TargetPort: intstr.FromString(manifestutils.PortJaegerThriftCompactName),
					Protocol:   corev1.ProtocolUDP,
				},
				{
					Name:       manifestutils.PortJaegerThriftBinaryName,
					Port:       manifestutils.PortJaegerThriftBinary,
					TargetPort: intstr.FromString(manifestutils.PortJaegerThriftBinaryName),
					Protocol:   corev1.ProtocolUDP,
				},
				{
					Name:       manifestutils.PortJaegerGrpcName,
					Port:       manifestutils.PortJaegerGrpc,
					TargetPort: intstr.FromString(manifestutils.PortJaegerGrpcName),
					Protocol:   corev1.ProtocolTCP,
				},
				{
					Name:       manifestutils.PortZipkinName,
					Port:       manifestutils.PortZipkin,
					TargetPort: intstr.FromString(manifestutils.PortZipkinName),
					Protocol:   corev1.ProtocolTCP,
				},
			},
			expectedResources: corev1.ResourceRequirements{
				Limits: corev1.ResourceList{
					corev1.ResourceCPU:    *resource.NewMilliQuantity(270, resource.BinarySI),
					corev1.ResourceMemory: *resource.NewQuantity(257698032, resource.BinarySI),
				},
				Requests: corev1.ResourceList{
					corev1.ResourceCPU:    *resource.NewMilliQuantity(81, resource.BinarySI),
					corev1.ResourceMemory: *resource.NewQuantity(77309416, resource.BinarySI),
				},
			},
			expectedVolumes: []corev1.Volume{
				{
					Name: manifestutils.ConfigVolumeName,
					VolumeSource: corev1.VolumeSource{
						ConfigMap: &corev1.ConfigMapVolumeSource{
							LocalObjectReference: corev1.LocalObjectReference{
								Name: "tempo-test",
							},
						},
					},
				},
				{
					Name: manifestutils.TmpStorageVolumeName,
					VolumeSource: corev1.VolumeSource{
						EmptyDir: &corev1.EmptyDirVolumeSource{},
					},
				},
				{
					Name: "ca-custom",
					VolumeSource: corev1.VolumeSource{
						ConfigMap: &corev1.ConfigMapVolumeSource{
							LocalObjectReference: corev1.LocalObjectReference{
								Name: "ca-custom",
							},
						},
					},
				},
				{
					Name: "cert-custom",
					VolumeSource: corev1.VolumeSource{
						Secret: &corev1.SecretVolumeSource{
							SecretName: "cert-custom",
						},
					},
				},
			},
			expectedVolumeMounts: []corev1.VolumeMount{
				{
					Name:      manifestutils.ConfigVolumeName,
					MountPath: "/conf",
					ReadOnly:  true,
				},
				{
					Name:      manifestutils.TmpStorageVolumeName,
					MountPath: manifestutils.TmpStoragePath,
				},
				{
					Name:      "ca-custom",
					MountPath: manifestutils.CAReceiver,
					ReadOnly:  true,
				},
				{
					Name:      "cert-custom",
					MountPath: manifestutils.TempoReceiverTLSDir(),
					ReadOnly:  true,
				},
			},
		},
		{
			name:          "Gateway enable",
			enableGateway: true,
			expectedContainerPorts: []corev1.ContainerPort{
				{
					Name:          manifestutils.OtlpGrpcPortName,
					ContainerPort: manifestutils.PortOtlpGrpcServer,
					Protocol:      corev1.ProtocolTCP,
				},
				{
					Name:          manifestutils.HttpPortName,
					ContainerPort: manifestutils.PortHTTPServer,
					Protocol:      corev1.ProtocolTCP,
				},
				{
					Name:          manifestutils.HttpMemberlistPortName,
					ContainerPort: manifestutils.PortMemberlist,
					Protocol:      corev1.ProtocolTCP,
				},
			},
			expectedServicePorts: []corev1.ServicePort{
				{
					Name:       manifestutils.OtlpGrpcPortName,
					Protocol:   corev1.ProtocolTCP,
					Port:       manifestutils.PortOtlpGrpcServer,
					TargetPort: intstr.FromString(manifestutils.OtlpGrpcPortName),
				},
				{
					Name:       manifestutils.HttpPortName,
					Protocol:   corev1.ProtocolTCP,
					Port:       manifestutils.PortHTTPServer,
					TargetPort: intstr.FromString(manifestutils.HttpPortName),
				},
			},
			expectedResources: corev1.ResourceRequirements{
				Limits: corev1.ResourceList{
					corev1.ResourceCPU:    *resource.NewMilliQuantity(260, resource.BinarySI),
					corev1.ResourceMemory: *resource.NewQuantity(236223200, resource.BinarySI),
				},
				Requests: corev1.ResourceList{
					corev1.ResourceCPU:    *resource.NewMilliQuantity(78, resource.BinarySI),
					corev1.ResourceMemory: *resource.NewQuantity(70866960, resource.BinarySI),
				},
			},
			expectedVolumes: []corev1.Volume{
				{
					Name: manifestutils.ConfigVolumeName,
					VolumeSource: corev1.VolumeSource{
						ConfigMap: &corev1.ConfigMapVolumeSource{
							LocalObjectReference: corev1.LocalObjectReference{
								Name: "tempo-test",
							},
						},
					},
				},
				{
					Name: manifestutils.TmpStorageVolumeName,
					VolumeSource: corev1.VolumeSource{
						EmptyDir: &corev1.EmptyDirVolumeSource{},
					},
				},
			},
			expectedVolumeMounts: []corev1.VolumeMount{
				{
					Name:      manifestutils.ConfigVolumeName,
					MountPath: "/conf",
					ReadOnly:  true,
				},
				{
					Name:      manifestutils.TmpStorageVolumeName,
					MountPath: manifestutils.TmpStoragePath,
				},
			},
		},
	}

	for _, ts := range tests {
		t.Run(ts.name, func(t *testing.T) {
			objects, err := BuildDistributor(manifestutils.Params{Tempo: v1alpha1.TempoStack{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test",
					Namespace: "project1",
				},
				Spec: v1alpha1.TempoStackSpec{
					Images: configv1alpha1.ImagesSpec{
						Tempo: "docker.io/grafana/tempo:1.5.0",
					},
					ServiceAccount: "tempo-test-serviceaccount",
					Template: v1alpha1.TempoTemplateSpec{
						Distributor: v1alpha1.TempoDistributorSpec{
							TLS: v1alpha1.ReceiversTLSSpec{
								Enabled: ts.enableReceiverTLS,
								CA:      "ca-custom",
								Cert:    "cert-custom",
							},
							TempoComponentSpec: v1alpha1.TempoComponentSpec{
								Replicas:     ptr.To(int32(1)),
								NodeSelector: map[string]string{"a": "b"},
								Tolerations: []corev1.Toleration{
									{
										Key: "c",
									},
								},
							},
						},
						Gateway: v1alpha1.TempoGatewaySpec{
							Enabled: ts.enableGateway,
						},
					},
					Resources: v1alpha1.Resources{
						Total: &corev1.ResourceRequirements{
							Limits: corev1.ResourceList{
								corev1.ResourceCPU:    resource.MustParse("1000m"),
								corev1.ResourceMemory: resource.MustParse("2Gi"),
							},
						},
					},
				},
			}})
			require.NoError(t, err)

			labels := manifestutils.ComponentLabels("distributor", "test")
			annotations := manifestutils.CommonAnnotations("")
			assert.Equal(t, 2, len(objects))
			assert.Equal(t, &v1.Deployment{
				TypeMeta: metav1.TypeMeta{
					APIVersion: v1.SchemeGroupVersion.String(),
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "tempo-test-distributor",
					Namespace: "project1",
					Labels:    labels,
				},
				Spec: v1.DeploymentSpec{
					Replicas: ptr.To(int32(1)),
					Selector: &metav1.LabelSelector{
						MatchLabels: labels,
					},
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Labels:      k8slabels.Merge(labels, map[string]string{"tempo-gossip-member": "true"}),
							Annotations: annotations,
						},
						Spec: corev1.PodSpec{
							ServiceAccountName: "tempo-test-serviceaccount",
							NodeSelector:       map[string]string{"a": "b"},
							Tolerations: []corev1.Toleration{
								{
									Key: "c",
								},
							},
							Affinity: manifestutils.DefaultAffinity(labels),
							Containers: []corev1.Container{
								{
									Name:  "tempo",
									Image: "docker.io/grafana/tempo:1.5.0",
									Args: []string{
										"-target=distributor",
										"-config.file=/conf/tempo.yaml",
										"-log.level=info",
									},
									Env:          []corev1.EnvVar{},
									VolumeMounts: ts.expectedVolumeMounts,
									Ports:        ts.expectedContainerPorts,
									ReadinessProbe: &corev1.Probe{
										ProbeHandler: corev1.ProbeHandler{
											HTTPGet: &corev1.HTTPGetAction{
												Scheme: corev1.URISchemeHTTP,
												Path:   manifestutils.TempoReadinessPath,
												Port:   intstr.FromString(manifestutils.HttpPortName),
											},
										},
										InitialDelaySeconds: 15,
										TimeoutSeconds:      1,
									},
									Resources:       ts.expectedResources,
									SecurityContext: manifestutils.TempoContainerSecurityContext(),
								},
							},
							Volumes: ts.expectedVolumes,
						},
					},
				},
			}, objects[0])

			assert.NoError(t, err)
			assert.Equal(t, &corev1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "tempo-test-distributor",
					Namespace: "project1",
					Labels:    labels,
				},
				Spec: corev1.ServiceSpec{
					Ports:    ts.expectedServicePorts,
					Selector: labels,
				},
			}, objects[1])
		})
	}
}
