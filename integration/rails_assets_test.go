package integration_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/paketo-buildpacks/occam"
	"github.com/sclevine/spec"

	. "github.com/onsi/gomega"
	. "github.com/paketo-buildpacks/occam/matchers"
)

func testRailsAssets(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect     = NewWithT(t).Expect
		Eventually = NewWithT(t).Eventually

		pack   occam.Pack
		docker occam.Docker
	)

	it.Before(func() {
		pack = occam.NewPack()
		docker = occam.NewDocker()
	})

	context("when building a rails app", func() {
		var (
			image     occam.Image
			container occam.Container

			name   string
			source string
		)

		it.Before(func() {
			var err error
			name, err = occam.RandomName()
			Expect(err).NotTo(HaveOccurred())

			source, err = occam.Source(filepath.Join("testdata", "rails"))
			Expect(err).NotTo(HaveOccurred())
		})

		it.After(func() {
			Expect(docker.Container.Remove.Execute(container.ID)).To(Succeed())
			Expect(docker.Image.Remove.Execute(image.ID)).To(Succeed())
			Expect(docker.Volume.Remove.Execute(occam.CacheVolumeNames(name))).To(Succeed())
			Expect(os.RemoveAll(source)).To(Succeed())
		})

		it("creates a working OCI image with rails assets precompiled", func() {
			var err error
			var logs fmt.Stringer
			image, logs, err = pack.WithNoColor().Build.
				WithBuildpacks(rubyBuildpack).
				WithPullPolicy("never").
				Execute(name, source)
			Expect(err).NotTo(HaveOccurred(), logs.String())

			container, err = docker.Container.Run.
				WithEnv(map[string]string{
					"PORT":            "8080",
					"SECRET_KEY_BASE": "some-secret",
				}).
				WithPublish("8080").
				WithPublishAll().
				Execute(image.ID)
			Expect(err).NotTo(HaveOccurred())

			Eventually(container).Should(BeAvailable())

			Eventually(container).Should(Serve(ContainSubstring("Hello World!")).OnPort(8080))

			Expect(logs).To(ContainLines(ContainSubstring("Buildpack for MRI")))
			Expect(logs).To(ContainLines(ContainSubstring("Buildpack for Bundler")))
			Expect(logs).To(ContainLines(ContainSubstring("Buildpack for Bundle Install")))
			Expect(logs).To(ContainLines(ContainSubstring("Buildpack for Node Engine")))
			Expect(logs).To(ContainLines(ContainSubstring("Buildpack for Yarn")))
			Expect(logs).To(ContainLines(ContainSubstring("Buildpack for Yarn Install")))
			Expect(logs).To(ContainLines(ContainSubstring("Buildpack for Rails Assets")))
			Expect(logs).To(ContainLines(ContainSubstring("Buildpack for Puma")))
			Expect(logs).NotTo(ContainLines(ContainSubstring("Buildpack for Procfile")))
			Expect(logs).NotTo(ContainLines(ContainSubstring("Buildpack for Environment Variables")))
		})
	})
}
