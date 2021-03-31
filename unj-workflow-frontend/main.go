package main

import (
	"strings"
	v1beta1 "unj-workflow-frontend/virtualservice/networking/v1beta1"

	appsv1 "github.com/pulumi/pulumi-kubernetes/sdk/v2/go/kubernetes/apps/v1"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v2/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v2/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi/config"
)

type VirtualService struct {
	Enabled  bool
	ApiPath  string
	Gateways []string
	Hosts    []string
}

type Service struct {
	Apply      bool
	Type       string
	Port       int
	TargetPort int
	Protocol   string
}

type Config struct {
	Replicas   int
	PullSecret string
}

type Resources struct {
	LimitCpu       string
	LimitMemory    string
	RequestsCpu    string
	RequestsMemory string
}

type Image struct {
	ImageTag   string
	Repository string
	PullPolicy string
}

type Env struct {
	DateHour            int
	KeycloakBaseUrl     string
	AppName             string
	KeycloakClientId    string
	ApplicationInitials string
	AppTokenName        string
	ApmUrl              string
	ApmToken            string
	SajKeycloakRealm    string
	SajAppEnvironment   string
	SajLangugae         string
	Tz                  string
	SajApolloUrl        string
	SajThemeUrl         string
}

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		if !strings.Contains(ctx.Stack(), "template") {

			conf := config.New(ctx, "")
			//config aplicação
			var configs Config
			conf.RequireObject("configs", &configs)
			var env Env
			conf.RequireObject("envs", &env)
			//config kubernetes
			conf_k8s := config.New(ctx, "kubernetes")
			//config istio
			conf_istio := config.New(ctx, "istio")
			//config relacionadas a imagem
			conf_image := config.New(ctx, "image")
			var image Image
			conf_image.RequireObject("configs", &image)
			//config resources
			conf_resources := config.New(ctx, "resources")
			var resource Resources
			conf_resources.RequireObject("configs", &resource)
			// config service
			conf_service := config.New(ctx, "service")
			var service Service
			conf_service.RequireObject("service", &service)

			appLabels := pulumi.StringMap{
				"app":     pulumi.String(ctx.Project()),
				"release": pulumi.String(ctx.Project() + "-" + conf_k8s.Require("namespace")),
			}
			annotations := pulumi.StringMap{
				"sidecar.istio.io/inject": pulumi.String(conf_istio.Require("istio_enabled")),
			}

			metadata := &metav1.ObjectMetaArgs{
				Name:   pulumi.String(ctx.Project() + "-" + conf_k8s.Require("namespace")),
				Labels: appLabels,
			}

			deployment, err := appsv1.NewDeployment(ctx, ctx.Project(), &appsv1.DeploymentArgs{
				Metadata: metadata,
				Spec: appsv1.DeploymentSpecArgs{
					Selector: &metav1.LabelSelectorArgs{
						MatchLabels: appLabels,
					},
					Replicas: pulumi.Int(configs.Replicas),
					Template: &corev1.PodTemplateSpecArgs{
						Metadata: &metav1.ObjectMetaArgs{
							Annotations: annotations,
							Labels:      appLabels,
						},
						Spec: &corev1.PodSpecArgs{
							ImagePullSecrets: corev1.LocalObjectReferenceArray{
								corev1.LocalObjectReferenceArgs{
									Name: pulumi.String(configs.PullSecret),
								},
							},
							Containers: corev1.ContainerArray{
								corev1.ContainerArgs{
									Name:            pulumi.String(ctx.Project()),
									Image:           pulumi.String(image.Repository + ":" + image.ImageTag),
									ImagePullPolicy: pulumi.String(image.PullPolicy),
									Env: corev1.EnvVarArray{

										corev1.EnvVarArgs{
											Name:  pulumi.String("SAJ_KEYCLOAK_URL"),
											Value: pulumi.String(env.KeycloakBaseUrl),
										},
										corev1.EnvVarArgs{
											Name:  pulumi.String("SAJ_KEYCLOAK_CLIENT_ID"),
											Value: pulumi.String(env.KeycloakClientId),
										},
										corev1.EnvVarArgs{
											Name:  pulumi.String("SAJ_KEYCLOAK_REALM"),
											Value: pulumi.String(env.SajKeycloakRealm),
										},
										corev1.EnvVarArgs{
											Name:  pulumi.String("SAJ_HEADER_APP_NAME"),
											Value: pulumi.String(env.AppName),
										},
										corev1.EnvVarArgs{
											Name:  pulumi.String("SAJ_HEADER_APP_INITIALS"),
											Value: pulumi.String(env.ApplicationInitials),
										},
										corev1.EnvVarArgs{
											Name:  pulumi.String("SAJ_APOLLO_URL"),
											Value: pulumi.String(env.SajApolloUrl),
										},
										corev1.EnvVarArgs{
											Name:  pulumi.String("SAJ_APOLLO_APM_SERVER_URL"),
											Value: pulumi.String(env.ApmUrl),
										},
										corev1.EnvVarArgs{
											Name:  pulumi.String("SAJ_APOLLO_APM_TOKEN"),
											Value: pulumi.String(env.ApmToken),
										},
										corev1.EnvVarArgs{
											Name:  pulumi.String("SAJ_APOLLO_APM_ENVIRONMENT"),
											Value: pulumi.String(env.SajAppEnvironment),
										},
										corev1.EnvVarArgs{
											Name:  pulumi.String("SAJ_LANGUAGE"),
											Value: pulumi.String(env.SajLangugae),
										},
										corev1.EnvVarArgs{
											Name:  pulumi.String("TZ"),
											Value: pulumi.String(env.Tz),
										},
										corev1.EnvVarArgs{
											Name:  pulumi.String("SAJ_THEME_URL"),
											Value: pulumi.String(env.SajThemeUrl),
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
											"cpu":    pulumi.String(resource.RequestsCpu),
											"memory": pulumi.String(resource.RequestsMemory),
										},
										Limits: pulumi.StringMap{
											"cpu":    pulumi.String(resource.LimitCpu),
											"memory": pulumi.String(resource.LimitMemory),
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
			svc, err := corev1.NewService(ctx, ctx.Project()+"-"+conf_k8s.Require("namespace"), &corev1.ServiceArgs{
				Metadata: metadata,
				Spec: &corev1.ServiceSpecArgs{
					Type: pulumi.String(service.Type),
					Ports: &corev1.ServicePortArray{
						&corev1.ServicePortArgs{
							Port:       pulumi.Int(service.Port),
							TargetPort: pulumi.Int(service.TargetPort),
							Protocol:   pulumi.String(service.Protocol),
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
			var virtualservices VirtualService
			conf_istio.RequireObject("virtualservice", &virtualservices)
			var gateway pulumi.StringArray
			for idx := 0; idx < len(virtualservices.Gateways); idx++ {
				gateway = append(gateway, pulumi.String(virtualservices.Gateways[idx]))
			}
			var host pulumi.StringArray
			for idx := 0; idx < len(virtualservices.Gateways); idx++ {
				host = append(host, pulumi.String(virtualservices.Hosts[idx]))
			}

			if virtualservices.Enabled {

				vs, err := v1beta1.NewVirtualService(ctx, ctx.Project()+"-"+conf_k8s.Require("namespace"), &v1beta1.VirtualServiceArgs{

					Metadata: metadata,
					Spec: v1beta1.VirtualServiceSpecArgs{
						Gateways: pulumi.StringArray(gateway),
						Hosts:    pulumi.StringArray(host),
						Http: v1beta1.VirtualServiceSpecHttpArray{
							v1beta1.VirtualServiceSpecHttpArgs{
								Match: v1beta1.VirtualServiceSpecHttpMatchArray{
									v1beta1.VirtualServiceSpecHttpMatchArgs{
										Uri: pulumi.StringMap{
											"prefix": pulumi.String(virtualservices.ApiPath),
										},
									},
								},
								Route: v1beta1.VirtualServiceSpecHttpRouteArray{
									v1beta1.VirtualServiceSpecHttpRouteArgs{
										Destination: v1beta1.VirtualServiceSpecHttpRouteDestinationArgs{
											Host: pulumi.String(ctx.Project() + "-" + conf_k8s.Require("namespace")),
											Port: v1beta1.VirtualServiceSpecHttpRouteDestinationPortArgs{
												Number: pulumi.Int(service.Port),
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

		}
		return nil
	})
}
