package integration_test

import (
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/paketo-buildpacks/occam"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"

	. "github.com/onsi/gomega"
)

var rubyBuildpack string

func TestIntegration(t *testing.T) {
	pack := occam.NewPack()
	Expect := NewWithT(t).Expect

	output, err := exec.Command("bash", "-c", "../scripts/package.sh --version 1.2.3").CombinedOutput()
	Expect(err).NotTo(HaveOccurred(), string(output))

	rubyBuildpack, err = filepath.Abs("../build/buildpackage.cnb")
	Expect(err).NotTo(HaveOccurred())

	SetDefaultEventuallyTimeout(20 * time.Second)

	builder, err := pack.Builder.Inspect.Execute()
	Expect(err).NotTo(HaveOccurred())

	suite := spec.New("Integration", spec.Parallel(), spec.Report(report.Terminal{}))

	// This test will only run on the Bionic full stack, because stack upgrade
	// failures have only been observed when upgrading from the Bionic full stack.
	// All other tests will run against the Bionic base stack and Jammy base stack
	if builder.BuilderName == "paketobuildpacks/builder:buildpackless-full" {
		suite("StackUpgrades", testGracefulStackUpgrades)
	}

	suite("Passenger", testPassenger)
	suite("Puma", testPuma)
	suite("Rackup", testRackup)
	suite("RailsAssets", testRailsAssets)
	suite("Rake", testRake)
	suite("ReproducibleBuilds", testReproducibleBuilds)
	suite("Thin", testThin)
	suite("Unicorn", testUnicorn)

	suite.Run(t)
}
