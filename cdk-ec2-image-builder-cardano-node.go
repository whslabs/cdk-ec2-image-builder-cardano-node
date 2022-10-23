package main

import (
	"log"
	"os"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	// "github.com/aws/aws-cdk-go/awscdk/v2/awssqs"
	"github.com/aws/aws-cdk-go/awscdk/v2/awscloudtrail"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsimagebuilder"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslogs"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslogsdestinations"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3assets"
	"github.com/aws/aws-cdk-go/awscdk/v2/awssns"
	"github.com/aws/aws-cdk-go/awscdk/v2/awssnssubscriptions"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type CdkEc2ImageBuilderCardanoNodeStackProps struct {
	awscdk.StackProps
}

func NewCdkEc2ImageBuilderCardanoNodeStack(scope constructs.Construct, id string, props *CdkEc2ImageBuilderCardanoNodeStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	content, err := os.ReadFile("component.yaml")
	if err != nil {
		log.Fatal(err)
	}

	// The code that defines your stack goes here
	component := awsimagebuilder.NewCfnComponent(stack, jsii.String("Component"), &awsimagebuilder.CfnComponentProps{
		Data:     jsii.String(string(content)),
		Name:     jsii.String("whslabs-cardano-node"),
		Platform: jsii.String("Linux"),
		Version:  jsii.String("1.0.0"),
	})

	imageRecipe := awsimagebuilder.NewCfnImageRecipe(stack, jsii.String("ImageRecipe"), &awsimagebuilder.CfnImageRecipeProps{
		Name:        jsii.String("whslabs-cardano-node"),
		ParentImage: awscdk.Fn_Sub(jsii.String("arn:${AWS::Partition}:imagebuilder:${AWS::Region}:aws:image/amazon-linux-2-x86/x.x.x"), nil),
		Version:     jsii.String("1.0.0"),
		BlockDeviceMappings: []interface{}{
			&awsimagebuilder.CfnImageRecipe_InstanceBlockDeviceMappingProperty{
				DeviceName: jsii.String("/dev/xvda"),
				Ebs: &awsimagebuilder.CfnImageRecipe_EbsInstanceBlockDeviceSpecificationProperty{
					DeleteOnTermination: jsii.Bool(true),
					VolumeSize:          jsii.Number(20),
				},
			},
		},
		Components: []interface{}{
			&awsimagebuilder.CfnImageRecipe_ComponentConfigurationProperty{
				ComponentArn: component.AttrArn(),
			},
		},
	})

	bucket := awss3.NewBucket(stack, jsii.String("BucketImageBuilder"), &awss3.BucketProps{
		BlockPublicAccess: awss3.BlockPublicAccess_BLOCK_ALL(),
		RemovalPolicy:     awscdk.RemovalPolicy_DESTROY,
	})

	roleImageBuilder := awsiam.NewRole(stack, jsii.String("RoleImageBuilder"), &awsiam.RoleProps{
		AssumedBy: awsiam.NewServicePrincipal(jsii.String("ec2.amazonaws.com"), nil),
		ManagedPolicies: &[]awsiam.IManagedPolicy{
			awsiam.ManagedPolicy_FromAwsManagedPolicyName(jsii.String("AmazonSSMManagedInstanceCore")),
			awsiam.ManagedPolicy_FromAwsManagedPolicyName(jsii.String("EC2InstanceProfileForImageBuilder")),
		},
	})

	roleImageBuilder.AttachInlinePolicy(awsiam.NewPolicy(stack, jsii.String("PolicyImageBuilder"), &awsiam.PolicyProps{
		Statements: &[]awsiam.PolicyStatement{
			awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
				Actions: &[]*string{
					jsii.String("s3:PutObject"),
				},
				Resources: &[]*string{
					bucket.ArnForObjects(jsii.String("*")),
				},
			}),
		},
	}))

	instanceProfile := awsiam.NewCfnInstanceProfile(stack, jsii.String("InstanceProfile"), &awsiam.CfnInstanceProfileProps{
		Roles: &[]*string{
			roleImageBuilder.RoleName(),
		},
	})

	infrastructureConfiguration := awsimagebuilder.NewCfnInfrastructureConfiguration(stack, jsii.String("InfrastructureConfiguration"), &awsimagebuilder.CfnInfrastructureConfigurationProps{
		InstanceProfileName: instanceProfile.Ref(),
		Name:                jsii.String("whslabs-cardano-node"),
		Logging: &awsimagebuilder.CfnInfrastructureConfiguration_LoggingProperty{
			S3Logs: &awsimagebuilder.CfnInfrastructureConfiguration_S3LogsProperty{
				S3BucketName: bucket.BucketName(),
			},
		},
	})

	awsimagebuilder.NewCfnImagePipeline(stack, jsii.String("ImagePipline"), &awsimagebuilder.CfnImagePipelineProps{
		ImageRecipeArn:                 imageRecipe.AttrArn(),
		InfrastructureConfigurationArn: infrastructureConfiguration.AttrArn(),
		Name:                           jsii.String("whslabs-cardano-node"),
		Status:                         jsii.String("DISABLED"),
	})

	cloudtrail := awscloudtrail.NewTrail(stack, jsii.String("CloudTrail"), &awscloudtrail.TrailProps{
		Bucket: awss3.NewBucket(stack, jsii.String("BucketCloudTrail"), &awss3.BucketProps{
			BlockPublicAccess: awss3.BlockPublicAccess_BLOCK_ALL(),
			RemovalPolicy:     awscdk.RemovalPolicy_DESTROY,
		}),
		SendToCloudWatchLogs: jsii.Bool(true),
	})

	topic := awssns.NewTopic(stack, jsii.String("Topic"), nil)

	topic.AddSubscription(awssnssubscriptions.NewEmailSubscription(jsii.String("hswongac@gmail.com"), nil))

	roleFunction := awsiam.NewRole(stack, jsii.String("RoleFunction"), &awsiam.RoleProps{
		AssumedBy: awsiam.NewServicePrincipal(jsii.String("lambda.amazonaws.com"), nil),
	})

	roleFunction.AttachInlinePolicy(awsiam.NewPolicy(stack, jsii.String("PolicyFunction"), &awsiam.PolicyProps{
		Statements: &[]awsiam.PolicyStatement{
			awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
				Actions: &[]*string{
					jsii.String("sns:Publish"),
				},
				Resources: &[]*string{
					topic.TopicArn(),
				},
			}),
			awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
				Actions: &[]*string{
					jsii.String("logs:CreateLogGroup"),
					jsii.String("logs:CreateLogStream"),
					jsii.String("logs:PutLogEvents"),
				},
				Resources: &[]*string{
					jsii.String("*"),
				},
			}),
		},
	}))

	lambda := awslambda.NewFunction(stack, jsii.String("Function"), &awslambda.FunctionProps{
		Code: awslambda.Code_FromAsset(jsii.String("./"), &awss3assets.AssetOptions{
			Bundling: &awscdk.BundlingOptions{
				Image: awscdk.DockerImage_FromRegistry(jsii.String("nixos/nix")),
				Command: &[]*string{jsii.String(`nix build \
--extra-experimental-features flakes \
--extra-experimental-features nix-command \
&& cp result/bin/rust-lambda-cloudtrail /asset-output/bootstrap`)},
				Entrypoint: &[]*string{jsii.String("/bin/sh"), jsii.String("-c")},
				User:       jsii.String("root:root"),
			},
		}),
		Handler: jsii.String("dummy"),
		Runtime: awslambda.Runtime_PROVIDED_AL2(),
		Environment: &map[string]*string{
			"CLOUDWATCH_BASE_URL": jsii.String("https://us-east-1.console.aws.amazon.com/cloudwatch/home?region=us-east-1#logsV2:log-groups/log-group/$252Faws$252Fimagebuilder$252F"),
			"TOPIC_ARN":           topic.TopicArn(),
		},
		Role: roleFunction,
	})

	lambda.AddPermission(jsii.String("Permission"), &awslambda.Permission{
		Principal: awsiam.NewServicePrincipal(jsii.String("logs.amazonaws.com"), nil),
		Action:    jsii.String("lambda:InvokeFunction"),
		SourceArn: cloudtrail.LogGroup().LogGroupArn(),
	})

	awslogs.NewSubscriptionFilter(stack, jsii.String("SubscriptionFilter"), &awslogs.SubscriptionFilterProps{
		Destination:   awslogsdestinations.NewLambdaDestination(lambda, nil),
		FilterPattern: awslogs.FilterPattern_All(awslogs.FilterPattern_StringValue(jsii.String("$.userAgent"), jsii.String("="), jsii.String("imagebuilder.amazonaws.com")), awslogs.FilterPattern_StringValue(jsii.String("$.eventName"), jsii.String("="), jsii.String("CreateImage"))),
		LogGroup:      cloudtrail.LogGroup(),
	})

	// example resource
	// queue := awssqs.NewQueue(stack, jsii.String("CdkEc2ImageBuilderCardanoNodeQueue"), &awssqs.QueueProps{
	// 	VisibilityTimeout: awscdk.Duration_Seconds(jsii.Number(300)),
	// })

	return stack
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	NewCdkEc2ImageBuilderCardanoNodeStack(app, "CdkEc2ImageBuilderCardanoNodeStack", &CdkEc2ImageBuilderCardanoNodeStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)
}

// env determines the AWS environment (account+region) in which our stack is to
// be deployed. For more information see: https://docs.aws.amazon.com/cdk/latest/guide/environments.html
func env() *awscdk.Environment {
	// If unspecified, this stack will be "environment-agnostic".
	// Account/Region-dependent features and context lookups will not work, but a
	// single synthesized template can be deployed anywhere.
	//---------------------------------------------------------------------------
	return nil

	// Uncomment if you know exactly what account and region you want to deploy
	// the stack to. This is the recommendation for production stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String("123456789012"),
	//  Region:  jsii.String("us-east-1"),
	// }

	// Uncomment to specialize this stack for the AWS Account and Region that are
	// implied by the current CLI configuration. This is recommended for dev
	// stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String(os.Getenv("CDK_DEFAULT_ACCOUNT")),
	//  Region:  jsii.String(os.Getenv("CDK_DEFAULT_REGION")),
	// }
}
