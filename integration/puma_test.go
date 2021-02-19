package integration_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/paketo-buildpacks/occam"
	"github.com/sclevine/spec"

	. "github.com/onsi/gomega"
	. "github.com/paketo-buildpacks/occam/matchers"
)

func testPuma(t *testing.T, context spec.G, it spec.S) {
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

	context("when building a puma app", func() {
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

			source, err = occam.Source(filepath.Join("testdata", "puma"))
			Expect(err).NotTo(HaveOccurred())
		})

		it.After(func() {
			Expect(docker.Container.Remove.Execute(container.ID)).To(Succeed())
			Expect(docker.Image.Remove.Execute(image.ID)).To(Succeed())
			Expect(docker.Volume.Remove.Execute(occam.CacheVolumeNames(name))).To(Succeed())
			Expect(os.RemoveAll(source)).To(Succeed())
		})

		it("creates a working OCI image with a puma start command", func() {
			var err error
			var logs fmt.Stringer
			image, logs, err = pack.WithNoColor().Build.
				WithBuildpacks(rubyBuildpack).
				WithPullPolicy("never").
				Execute(name, source)
			Expect(err).NotTo(HaveOccurred(), logs.String())

			container, err = docker.Container.Run.
				WithEnv(map[string]string{"PORT": "9292"}).
				WithPublish("9292").
				WithPublishAll().
				Execute(image.ID)
			Expect(err).NotTo(HaveOccurred())

			Eventually(container).Should(BeAvailable())
			Eventually(container).Should(Serve(ContainSubstring("Hello world!")).OnPort(9292))

			Expect(logs).To(ContainLines(ContainSubstring("MRI Buildpack")))
			Expect(logs).To(ContainLines(ContainSubstring("Bundler Buildpack")))
			Expect(logs).To(ContainLines(ContainSubstring("Bundle Install Buildpack")))
			Expect(logs).To(ContainLines(ContainSubstring("Puma Buildpack")))
			Expect(logs).NotTo(ContainLines(ContainSubstring("Procfile Buildpack")))
		})

		context("using optional utility buildpacks", func() {
			it.Before(func() {
				Expect(ioutil.WriteFile(filepath.Join(source, "Procfile"), []byte("web: bundle exec puma -b tcp://0.0.0.0:${PORT}"), 0644)).To(Succeed())
			})

			it("builds a working OCI image with start command from the Procfile and incorporating the utility buildpacks' effect", func() {
				var err error
				var logs fmt.Stringer
				image, logs, err = pack.WithNoColor().Build.
					WithEnv(map[string]string{"BPE_SOME_VARIABLE": "SOME_VALUE"}).
					WithBuildpacks(rubyBuildpack).
					WithPullPolicy("never").
					Execute(name, source)
				Expect(err).NotTo(HaveOccurred(), logs.String())

				container, err = docker.Container.Run.
					WithEnv(map[string]string{"PORT": "8080"}).
					WithPublish("8080").
					WithPublishAll().
					Execute(image.ID)
				Expect(err).NotTo(HaveOccurred())

				Eventually(container).Should(BeAvailable())

				Expect(image.Buildpacks[5].Key).To(Equal("paketo-buildpacks/environment-variables"))
				Expect(image.Buildpacks[5].Layers["environment-variables"].Metadata["variables"]).To(Equal(map[string]interface{}{"SOME_VARIABLE": "SOME_VALUE"}))

				Eventually(container).Should(Serve(ContainSubstring("Hello world!")).OnPort(8080))

				Expect(logs).To(ContainLines(ContainSubstring("MRI Buildpack")))
				Expect(logs).To(ContainLines(ContainSubstring("Bundler Buildpack")))
				Expect(logs).To(ContainLines(ContainSubstring("Bundle Install Buildpack")))
				Expect(logs).To(ContainLines(ContainSubstring("Puma Buildpack")))
				Expect(logs).To(ContainLines(ContainSubstring("Procfile Buildpack")))
				Expect(logs).To(ContainLines(ContainSubstring("bundle exec puma -b tcp://0.0.0.0:${PORT}")))
				Expect(logs).To(ContainLines(ContainSubstring("Environment Variables Buildpack")))
			})
		})
	})
}
