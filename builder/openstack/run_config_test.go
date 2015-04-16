package openstack

import (
	"github.com/mitchellh/packer/packer"
	"io/ioutil"
	"os"
	"testing"
)

func init() {
	// Clear out the openstack env vars so they don't
	// affect our tests.
	os.Setenv("SDK_USERNAME", "")
	os.Setenv("SDK_PASSWORD", "")
	os.Setenv("SDK_PROVIDER", "")
}

func testRunConfig() *RunConfig {
	return &RunConfig{
		SourceImage: "abcd",
		Flavor:      "m1.small",
		SSHUsername: "root",
	}
}

func TestRunConfigPrepare(t *testing.T) {
	c := testRunConfig()
	err := c.Prepare(nil)
	if len(err) > 0 {
		t.Fatalf("err: %s", err)
	}
}

func TestRunConfigPrepare_InstanceType(t *testing.T) {
	c := testRunConfig()
	c.Flavor = ""
	if err := c.Prepare(nil); len(err) != 1 {
		t.Fatalf("err: %s", err)
	}
}

func TestRunConfigPrepare_SourceImage(t *testing.T) {
	c := testRunConfig()
	c.SourceImage = ""
	if err := c.Prepare(nil); len(err) != 1 {
		t.Fatalf("err: %s", err)
	}
}

func TestRunConfigPrepare_SSHPort(t *testing.T) {
	c := testRunConfig()
	c.SSHPort = 0
	if err := c.Prepare(nil); len(err) != 0 {
		t.Fatalf("err: %s", err)
	}

	if c.SSHPort != 22 {
		t.Fatalf("invalid value: %d", c.SSHPort)
	}

	c.SSHPort = 44
	if err := c.Prepare(nil); len(err) != 0 {
		t.Fatalf("err: %s", err)
	}

	if c.SSHPort != 44 {
		t.Fatalf("invalid value: %d", c.SSHPort)
	}
}

func TestRunConfigPrepare_SSHTimeout(t *testing.T) {
	c := testRunConfig()
	c.RawSSHTimeout = ""
	if err := c.Prepare(nil); len(err) != 0 {
		t.Fatalf("err: %s", err)
	}

	c.RawSSHTimeout = "bad"
	if err := c.Prepare(nil); len(err) != 1 {
		t.Fatalf("err: %s", err)
	}
}

func TestRunConfigPrepare_SSHUsername(t *testing.T) {
	c := testRunConfig()
	c.SSHUsername = ""
	if err := c.Prepare(nil); len(err) != 0 {
		t.Fatalf("err: %s", err)
	}
}

func TestRunConfigPrepare_Networks(t *testing.T) {
	c := testRunConfig()
	network_uuid_var := "7d83eb1e-76ed-4a35-9247-54380a5419ea"

	tpl, err := packer.NewConfigTemplate()
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	tpl.UserVars["network_uuid"] = network_uuid_var

	c.Networks = []string{`{{user "network_uuid"}}`}

	if err := c.Prepare(tpl); len(err) != 0 {
		t.Fatalf("err: %s", err)
	}
	expected := network_uuid_var
	if c.Networks[0] != expected {
		t.Fatalf("Networks was not templated. Value is: %s", c.Networks)
	}
}

func TestRunConfigPrepare_UserDataFile(t *testing.T) {
	c := testRunConfig()
	c.UserDataFile = "badfile"
	if err := c.Prepare(nil); len(err) != 1 {
		t.Fatalf("err: %s", err)
	}

	tf, err := ioutil.TempFile("", "packer")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	defer tf.Close()

	c.UserDataFile = tf.Name()
	if err := c.Prepare(nil); len(err) != 0 {
		t.Fatalf("err: %s", err)
	}
}
