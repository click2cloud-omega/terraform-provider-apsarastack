package apsarastack

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccApsaraStackSecurityGroupsDataSourceBasic(t *testing.T) {
	rand := acctest.RandIntRange(1000, 9999)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackSecurityGroupsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${apsarastack_security_group.default.name}"`,
		}),
		fakeConfig: testAccCheckApsaraStackSecurityGroupsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${apsarastack_security_group.default.name}_fake"`,
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackSecurityGroupsDataSourceConfig(rand, map[string]string{
			"ids": `[ "${apsarastack_security_group.default.id}" ]`,
		}),
		fakeConfig: testAccCheckApsaraStackSecurityGroupsDataSourceConfig(rand, map[string]string{
			"ids": `[ "${apsarastack_security_group.default.id}_fake" ]`,
		}),
	}

	vpcIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackSecurityGroupsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${apsarastack_security_group.default.name}"`,
			"vpc_id":     `"${apsarastack_security_group.default.vpc_id}"`,
		}),
		fakeConfig: testAccCheckApsaraStackSecurityGroupsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${apsarastack_security_group.default.name}"`,
			"vpc_id":     `"${apsarastack_security_group.default.vpc_id}_fake"`,
		}),
	}

	tagsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackSecurityGroupsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${apsarastack_security_group.default.name}"`,
			"tags": `{
                         from = "datasource"
                         usage1 = "test"
                         usage2 = "test"
                        }`,
		}),
		fakeConfig: testAccCheckApsaraStackSecurityGroupsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${apsarastack_security_group.default.name}"`,
			"tags": `{
                         from = "datasource"
                         usage1 = "test"
                         usage2 = "test_fake"
                        }`,
		}),
	}

	resourceGroupIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackSecurityGroupsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${apsarastack_security_group.default.name}"`,
			//"resource_group_id": fmt.Sprintf(`"%s"`, os.Getenv("APSARASTACK_RESOURCE_GROUP_ID")),
		}),
		fakeConfig: testAccCheckApsaraStackSecurityGroupsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${apsarastack_security_group.default.name}"`,
			//"resource_group_id": fmt.Sprintf(`"%s_fake"`, os.Getenv("APSARASTACK_RESOURCE_GROUP_ID")),
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackSecurityGroupsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${apsarastack_security_group.default.name}"`,
			"ids":        `[ "${apsarastack_security_group.default.id}" ]`,
			"vpc_id":     `"${apsarastack_security_group.default.vpc_id}"`,
			//"resource_group_id": fmt.Sprintf(`"%s"`, os.Getenv("APSARASTACK_RESOURCE_GROUP_ID")),
			"tags": `{
                         from = "datasource"
                         usage1 = "test"
                         usage2 = "test"
                        }`,
		}),
		fakeConfig: testAccCheckApsaraStackSecurityGroupsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${apsarastack_security_group.default.name}_fake"`,
			"ids":        `[ "${apsarastack_security_group.default.id}" ]`,
			"vpc_id":     `"${apsarastack_security_group.default.vpc_id}"`,
			//"resource_group_id": fmt.Sprintf(`"%s"`, os.Getenv("APSARASTACK_RESOURCE_GROUP_ID")),
			"tags": `{
                         from = "datasource"
                         usage1 = "test"
                         usage2 = "test"
                        }`,
		}),
	}

	securityGroupsCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf, idsConf, vpcIdConf, tagsConf, resourceGroupIdConf, allConf)
}

func testAccCheckApsaraStackSecurityGroupsDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccCheckApsaraStackSecurityGroupsDataSourceConfig%d"
}
resource "apsarastack_vpc" "default" {
  cidr_block = "172.16.0.0/12"
  name = "${var.name}"
}

resource "apsarastack_security_group" "default" {
  name        = "${var.name}"
  description = "test security group"
  vpc_id      = "${apsarastack_vpc.default.id}"
  //resource_group_id = "%s"
  tags = {
		from = "datasource"
		usage1 = "test"
		usage2 = "test"
  }
}

data "apsarastack_security_groups" "default" {
  %s
}`, rand, os.Getenv("APSARASTACK_RESOURCE_GROUP_ID"), strings.Join(pairs, "\n  "))
	return config
}

var existSecurityGroupsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":           "1",
		"names.#":         "1",
		"groups.#":        "1",
		"groups.0.vpc_id": CHECKSET,
		//"groups.0.resource_group_id":   os.Getenv("APSARASTACK_RESOURCE_GROUP_ID"),
		//"groups.0.security_group_type": "normal",
		"groups.0.name":        fmt.Sprintf("tf-testAccCheckApsaraStackSecurityGroupsDataSourceConfig%d", rand),
		"groups.0.tags.from":   "datasource",
		"groups.0.tags.usage1": "test",
		"groups.0.tags.usage2": "test",
		//"groups.0.inner_access":        "true",
		"groups.0.creation_time": CHECKSET,
		"groups.0.description":   "test security group",
		"groups.0.id":            CHECKSET,
	}
}

var fakeSecurityGroupsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":    "0",
		"names.#":  "0",
		"groups.#": "0",
	}
}

var securityGroupsCheckInfo = dataSourceAttr{
	resourceId:   "data.apsarastack_security_groups.default",
	existMapFunc: existSecurityGroupsMapFunc,
	fakeMapFunc:  fakeSecurityGroupsMapFunc,
}
