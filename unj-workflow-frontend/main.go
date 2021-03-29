package main

import (
	"strconv"
	"strings"

	apiext "github.com/pulumi/pulumi-kubernetes/sdk/v2/go/kubernetes/apiextensions"
	appsv1 "github.com/pulumi/pulumi-kubernetes/sdk/v2/go/kubernetes/apps/v1"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v2/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v2/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi/config"
)

type VirtualServiceSpecHttpMatchRouteDestinationPort struct {
	Number int
}

type VirtualServiceSpecHttpMatchRouteDestination struct {
	Host string
	Port VirtualServiceSpecHttpMatchRouteDestinationPort
}

type VirtualServiceSpecHttpMatchRoute struct {
	Destination VirtualServiceSpecHttpMatchRouteDestination
}

type VirtualServiceSpecHttpMatchUri struct {
	Prefix string
}

type VirtualServiceSpecHttpMatch struct {
	Name  string
	Uri   VirtualServiceSpecHttpMatchUri
	Route VirtualServiceSpecHttpMatchRoute
}

type VirtualServiceSpecHttp struct {
	Match VirtualServiceSpecHttpMatch
}

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		stack := ctx.Stack()
		if !strings.Contains(stack, "template") {

			conf := config.New(ctx, "")
			conf_k8s := config.New(ctx, "kubernetes")
			conf_istio := config.New(ctx, "istio")
			conf_image := config.New(ctx, "image")
			conf_resources := config.New(ctx, "resources")
			name := ctx.Project()

			replicas, err := strconv.ParseInt(conf.Require("replicas"), 10, 64)
			if err != nil {
				return err
			}

			namespace := conf_k8s.Require("namespace")

			istio_enabled := conf_istio.Require("istio_enabled")

			pullSecret := conf.Require("pullSecret")

			repository := conf_image.Require("repository")
			image_tag := conf_image.Require("image_tag")
			pullPolicy := conf_image.Require("pullPolicy")

			appLabels := pulumi.StringMap{
				"app":     pulumi.String(name),
				"release": pulumi.String(name + "-" + namespace),
			}
			annotations := pulumi.StringMap{
				"sidecar.istio.io/inject": pulumi.String(istio_enabled),
			}

			deployment, err := appsv1.NewDeployment(ctx, name, &appsv1.DeploymentArgs{
				Metadata: &metav1.ObjectMetaArgs{
					Name:   pulumi.String(name + "-" + namespace),
					Labels: appLabels,
				},
				Spec: appsv1.DeploymentSpecArgs{
					Selector: &metav1.LabelSelectorArgs{
						MatchLabels: appLabels,
					},
					Replicas: pulumi.Int(replicas),
					Template: &corev1.PodTemplateSpecArgs{
						Metadata: &metav1.ObjectMetaArgs{
							Annotations: annotations,
							Labels:      appLabels,
						},
						Spec: &corev1.PodSpecArgs{
							ImagePullSecrets: corev1.LocalObjectReferenceArray{
								corev1.LocalObjectReferenceArgs{
									Name: pulumi.String(pullSecret),
								},
							},
							Containers: corev1.ContainerArray{
								corev1.ContainerArgs{
									Name:            pulumi.String(name),
									Image:           pulumi.String(repository + ":" + image_tag),
									ImagePullPolicy: pulumi.String(pullPolicy),
									Env: corev1.EnvVarArray{
										corev1.EnvVarArgs{
											Name:  pulumi.String("SAJ_KEYCLOAK_URL"),
											Value: pulumi.String("https://identity-platform.softplan.io/auth"),
										},
									},
									Ports: corev1.ContainerPortArray{
										&corev1.ContainerPortArgs{
											Name:          pulumi.String("http"),
											ContainerPort: pulumi.Int(80),
											Protocol:      pulumi.String("TCP"),
										},
									},
									Resources: &corev1.ResourceRequirementsArgs{
										Requests: pulumi.StringMap{
											"cpu":    pulumi.String(conf_resources.Require("requests_cpu")),
											"memory": pulumi.String(conf_resources.Require("requests_memory")),
										},
										Limits: pulumi.StringMap{
											"cpu":    pulumi.String(conf_resources.Require("limit_cpu")),
											"memory": pulumi.String(conf_resources.Require("limit_memory")),
										},
									},
									ReadinessProbe: corev1.ProbeArgs{
										InitialDelaySeconds: pulumi.Int(30),
										TcpSocket: corev1.TCPSocketActionArgs{
											Port: pulumi.String("http"),
										},
									},
								}},
						},
					},
				},
			})
			if err != nil {
				return err
			}

			ctx.Export("name", deployment.Metadata.Elem().Name())

			//service
			svc, err := corev1.NewService(ctx, name+"-"+namespace, &corev1.ServiceArgs{
				Metadata: &metav1.ObjectMetaArgs{
					Name:   pulumi.String(name + "-" + namespace),
					Labels: appLabels,
				},
				Spec: &corev1.ServiceSpecArgs{
					Type: pulumi.String("ClusterIP"),
					Ports: &corev1.ServicePortArray{
						&corev1.ServicePortArgs{
							Port:       pulumi.Int(80),
							TargetPort: pulumi.Int(80),
							Protocol:   pulumi.String("TCP"),
						},
					},
					Selector: appLabels,
				},
			})
			if err != nil {
				return err
			}
			ctx.Export("name", svc.Metadata.Elem().Name())

			//virtualservice
			vs, err := apiext.NewCustomResource(ctx, name+"-"+namespace, &apiext.CustomResourceArgs{
				ApiVersion: pulumi.String("networking.istio.io/v1beta1"),
				Kind:       pulumi.String("VirtualService"),
				Metadata: &metav1.ObjectMetaArgs{
					Name:   pulumi.String(name + "-" + namespace),
					Labels: appLabels,
				},
				OtherFields: map[string]interface{}{
					"spec": map[string]interface{}{
						"gateways": []string{"teste.com.br", "dfasfdf"},
						"hosts":    []string{"teste.com.br", "dfasfdf"},
						"http": []VirtualServiceSpecHttp{
							Match: []VirtualServiceSpecHttpMatch{
								Name: "default",
								Uri: VirtualServiceSpecHttpMatchUri{
									Prefix: "/",
								},
								Route: VirtualServiceSpecHttpMatchRoute{
									Destination: VirtualServiceSpecHttpMatchRouteDestination{
										Host: "unjdfdf.com.br",
										Port: VirtualServiceSpecHttpMatchRouteDestinationPort{
											Number: 80,
										},
									},
								},
							},
						},
					},
				},
			})
			if err != nil {
				return err
			}
			ctx.Export("name", vs.Metadata.Elem().Name())

		}
		return nil
	})
}
