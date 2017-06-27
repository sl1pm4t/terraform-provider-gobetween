package gobetween

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"

	"github.com/dustinkirkland/golang-petname"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	gb "github.com/yyyar/gobetween/src/config"
)

func TestAccServer_basic(t *testing.T) {
	var s gb.Server
	r := rand.New(rand.NewSource(99))

	name := strings.ToLower(petname.Generate(2, "-"))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccServer_basic(name, r.Int(), r.Int()),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServerExists(t, "gobetween_server.foo", &s),
				),
			},
		},
	})
}

func testAccCheckServerExists(t *testing.T, name string, server *gb.Server) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found in state: %s", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		id := rs.Primary.ID
		client := testAccProvider.Meta().(GBProvider).Client

		sv, err := client.Api.GetServer(id)
		if err != nil {
			return err
		}

		if s != nil {
			*server = *sv
			return nil
		}

		return fmt.Errorf("Server not found: %s", rs.Primary.ID)
	}
}

func testAccServer_basic(name string, port1, port2 int) string {
	return fmt.Sprintf(`
resource "gobetween_server" "foo" {
  name = "%s"
  bind = "0.0.0.0:9999"
  balance = "weight"

  discovery {
	  static_list = ["127.0.0.1:%d weight=1", "127.0.0.1:%d weight=2"]
  }

}`, name, port1, port2)
}
