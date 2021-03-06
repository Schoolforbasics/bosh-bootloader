package templates_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/cloudfoundry/bosh-bootloader/aws/cloudformation/templates"
)

var _ = Describe("InternalSubnetTemplateBuilder", func() {
	var builder templates.InternalSubnetTemplateBuilder

	BeforeEach(func() {
		builder = templates.NewInternalSubnetTemplateBuilder()
	})

	Describe("InternalSubnet", func() {
		It("returns a template with parameters for the internal subnet", func() {
			subnet := builder.InternalSubnet(0, "1", "10.0.16.0/20")

			Expect(subnet.Parameters).To(HaveLen(1))
			Expect(subnet.Parameters).To(HaveKeyWithValue("InternalSubnet1CIDR", templates.Parameter{
				Description: "CIDR block for InternalSubnet1.",
				Type:        "String",
				Default:     "10.0.16.0/20",
			}))
		})

		It("returns a template with resources for the internal subnet", func() {
			subnet := builder.InternalSubnet(0, "1", "10.0.16.0/20")

			Expect(subnet.Resources).To(HaveLen(4))
			Expect(subnet.Resources).To(HaveKeyWithValue("InternalSubnet1", templates.Resource{
				Type: "AWS::EC2::Subnet",
				Properties: templates.Subnet{
					AvailabilityZone: map[string]interface{}{
						"Fn::Select": []interface{}{
							"0",
							map[string]templates.Ref{
								"Fn::GetAZs": templates.Ref{"AWS::Region"},
							},
						},
					},
					CidrBlock: templates.Ref{"InternalSubnet1CIDR"},
					VpcId:     templates.Ref{"VPC"},
					Tags: []templates.Tag{
						{
							Key:   "Name",
							Value: "Internal1",
						},
					},
				},
			}))

			Expect(subnet.Resources).To(HaveKeyWithValue("InternalRouteTable", templates.Resource{
				Type: "AWS::EC2::RouteTable",
				Properties: templates.RouteTable{
					VpcId: templates.Ref{"VPC"},
				},
			}))

			Expect(subnet.Resources).To(HaveKeyWithValue("InternalRoute", templates.Resource{
				Type:      "AWS::EC2::Route",
				DependsOn: "NATInstance",
				Properties: templates.Route{
					DestinationCidrBlock: "0.0.0.0/0",
					RouteTableId:         templates.Ref{"InternalRouteTable"},
					InstanceId:           templates.Ref{"NATInstance"},
				},
			}))

			Expect(subnet.Resources).To(HaveKeyWithValue("InternalSubnet1RouteTableAssociation", templates.Resource{
				Type: "AWS::EC2::SubnetRouteTableAssociation",
				Properties: templates.SubnetRouteTableAssociation{
					RouteTableId: templates.Ref{"InternalRouteTable"},
					SubnetId:     templates.Ref{"InternalSubnet1"},
				},
			}))
		})

		It("returns a template with outputs for the internal subnet", func() {
			subnet := builder.InternalSubnet(0, "1", "10.0.16.0/20")

			Expect(subnet.Outputs).To(HaveLen(3))
			Expect(subnet.Outputs).To(HaveKeyWithValue("InternalSubnet1CIDR", templates.Output{
				Value: templates.Ref{"InternalSubnet1CIDR"},
			}))
			Expect(subnet.Outputs).To(HaveKeyWithValue("InternalSubnet1AZ", templates.Output{
				Value: templates.FnGetAtt{
					[]string{
						"InternalSubnet1",
						"AvailabilityZone",
					},
				},
			}))
			Expect(subnet.Outputs).To(HaveKeyWithValue("InternalSubnet1Name", templates.Output{
				Value: templates.Ref{"InternalSubnet1"},
			}))
		})
	})

})
