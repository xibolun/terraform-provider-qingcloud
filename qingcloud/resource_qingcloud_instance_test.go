package qingcloud

import (
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	qc "github.com/yunify/qingcloud-sdk-go/service"
)

func TestAccQingcloudInstance_basic(t *testing.T) {
	var instance qc.DescribeInstancesOutput
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: "qingcloud_instance.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccInstanceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists("qingcloud_instance.foo", &instance),
					resource.TestCheckResourceAttr(
						"qingcloud_instance.foo", "image_id", "centos7x64d"),
					resource.TestCheckResourceAttr(
						"qingcloud_instance.foo", "managed_vxnet_id", "vxnet-0"),
					resource.TestCheckResourceAttr(
						"qingcloud_instance.foo", "cpu", "1"),
					resource.TestCheckResourceAttr(
						"qingcloud_instance.foo", "memory", "1024"),
					resource.TestCheckResourceAttr(
						"qingcloud_instance.foo", "instance_class", "0"),
				),
			},
			resource.TestStep{
				Config: testAccInstanceConfigTwo,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists("qingcloud_instance.foo", &instance),
					resource.TestCheckResourceAttr(
						"qingcloud_instance.foo", resourceName, "instance"),
					resource.TestCheckResourceAttr(
						"qingcloud_instance.foo", resourceDescription, "instance"),
					resource.TestCheckResourceAttr(
						"qingcloud_instance.foo", "image_id", "centos7x64d"),
					resource.TestCheckResourceAttr(
						"qingcloud_instance.foo", "managed_vxnet_id", "vxnet-0"),
					resource.TestCheckResourceAttr(
						"qingcloud_instance.foo", "cpu", "2"),
					resource.TestCheckResourceAttr(
						"qingcloud_instance.foo", "memory", "2048"),
					resource.TestCheckResourceAttr(
						"qingcloud_instance.foo", "instance_class", "0"),
				),
			},
		},
	})
}

func testAccCheckInstanceExists(n string, i *qc.DescribeInstancesOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Instance ID is set")
		}

		client := testAccProvider.Meta().(*QingCloudClient)
		input := new(qc.DescribeInstancesInput)
		input.Instances = []*string{qc.String(rs.Primary.ID)}
		d, err := client.instance.DescribeInstances(input)

		log.Printf("[WARN] instance id %#v", rs.Primary.ID)

		if err != nil {
			return err
		}

		if d == nil || len(d.InstanceSet) == 0 {
			return fmt.Errorf("Instance not found ")
		}

		*i = *d
		return nil
	}
}

func testAccCheckInstanceDestroy(s *terraform.State) error {
	return testAccCheckInstanceWithProvider(s, testAccProvider)
}

func testAccCheckInstanceWithProvider(s *terraform.State, provider *schema.Provider) error {
	client := provider.Meta().(*QingCloudClient)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "qingcloud_instance" {
			continue
		}
		// Try to find the resource
		input := new(qc.DescribeInstancesInput)
		input.Instances = []*string{qc.String(rs.Primary.ID)}
		output, err := client.instance.DescribeInstances(input)
		if err == nil {
			if !isInstanceDeleted(output.InstanceSet) {
				return fmt.Errorf("Found  Instance: %s", rs.Primary.ID)
			}
		}
	}
	return nil
}

const testAccInstanceConfig = `
resource "qingcloud_keypair" "foo"{
	public_key = "    ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCyLSPqVIdXGH0QlGeWcPwa1fjTRKl6WtMiaSsP8/GnwjakDSKILUCoNe1yIpiK8F0/gmL71xaDQyfl7k6aE+gn6lSLUjpDmucAF1luGg6l7CIN+6hCqY3YqlAI05Tqwu0PdLAwCbGwdHcaWfECcbROJk5D0zpCTHmissrrAxdOv72g9Ple8KJ6C7F1tz6wmG0zUeineguGjW/PvfZiBDWZ/CyXGPeMDJxv3lrIiLa/ShgnQOxFTdHJPCw+F0/XlSzlIzP3gfni1vXxJWvYjdE9ULo7Z1DLWgZ73FCbeAvX/0e9C9jwT21Qa5RUy4pSP8m4WXSJgw2f9IpY1vIJFSZP root@centos1    "
}
resource "qingcloud_instance" "foo" {
	image_id = "centos7x64d"
	keypair_ids = ["${qingcloud_keypair.foo.id}"]
}
`
const testAccInstanceConfigTwo = `
resource "qingcloud_keypair" "foo"{
	public_key = "    ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCyLSPqVIdXGH0QlGeWcPwa1fjTRKl6WtMiaSsP8/GnwjakDSKILUCoNe1yIpiK8F0/gmL71xaDQyfl7k6aE+gn6lSLUjpDmucAF1luGg6l7CIN+6hCqY3YqlAI05Tqwu0PdLAwCbGwdHcaWfECcbROJk5D0zpCTHmissrrAxdOv72g9Ple8KJ6C7F1tz6wmG0zUeineguGjW/PvfZiBDWZ/CyXGPeMDJxv3lrIiLa/ShgnQOxFTdHJPCw+F0/XlSzlIzP3gfni1vXxJWvYjdE9ULo7Z1DLWgZ73FCbeAvX/0e9C9jwT21Qa5RUy4pSP8m4WXSJgw2f9IpY1vIJFSZP root@centos1    "
}
resource "qingcloud_instance" "foo" {
	image_id = "centos7x64d"
	keypair_ids = ["${qingcloud_keypair.foo.id}"]
	cpu = 2
    memory = 2048
	name = "instance"
	description = "instance"
}
`
