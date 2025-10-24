package integration_test

import (
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"

	. "github.com/onsi/gomega"
)

var rubyBuildpack string

func TestIntegration(t *testing.T) {
	Expect := NewWithT(t).Expect

	output, err := exec.Command("bash", "-c", "../scripts/package.sh --version 1.2.3").CombinedOutput()
	Expect(err).NotTo(HaveOccurred(), string(output))

	rubyBuildpack, err = filepath.Abs("../build/buildpackage.cnb")
	Expect(err).NotTo(HaveOccurred())

	SetDefaultEventuallyTimeout(20 * time.Second)

	suite := spec.New("Integration", spec.Parallel(), spec.Report(report.Terminal{}))

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
