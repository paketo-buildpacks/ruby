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

func testRake(t *testing.T, context spec.G, it spec.S) {
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

	context("when building a rake container", func() {
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
		})

		it.After(func() {
			Expect(docker.Container.Remove.Execute(container.ID)).To(Succeed())
			Expect(docker.Image.Remove.Execute(image.ID)).To(Succeed())
			Expect(docker.Volume.Remove.Execute(occam.CacheVolumeNames(name))).To(Succeed())
			Expect(os.RemoveAll(source)).To(Succeed())
		})

		context("uses the rake gem", func() {
			it("creates a working OCI image", func() {
				var err error
				source, err = occam.Source(filepath.Join("testdata", "rake"))
				Expect(err).NotTo(HaveOccurred())

				var logs fmt.Stringer
				image, logs, err = pack.WithNoColor().Build.
					WithBuildpacks(rubyBuildpack).
					WithPullPolicy("never").
					Execute(name, source)
				Expect(err).NotTo(HaveOccurred(), logs.String())

				container, err = docker.Container.Run.Execute(image.ID)
				Expect(err).NotTo(HaveOccurred())

				rLogs := func() fmt.Stringer {
					rakeLogs, err := docker.Container.Logs.Execute(container.ID)
					Expect(err).NotTo(HaveOccurred())
					return rakeLogs
				}

				Eventually(rLogs).Should(ContainSubstring("I am a rake task"))

				Expect(logs).To(ContainLines(ContainSubstring("MRI Buildpack")))
				Expect(logs).To(ContainLines(ContainSubstring("Bundler Buildpack")))
				Expect(logs).To(ContainLines(ContainSubstring("Bundle Install Buildpack")))
				Expect(logs).To(ContainLines(ContainSubstring("Rake Buildpack")))
				Expect(logs).To(ContainLines(ContainSubstring("bundle exec rake")))
				Expect(logs).NotTo(ContainLines(ContainSubstring("Procfile Buildpack")))
				Expect(logs).NotTo(ContainLines(ContainSubstring("Environment Variables Buildpack")))
			})
		})

		context("does not use rake gem", func() {
			it("creates a working OCI image", func() {
				var err error
				source, err = occam.Source(filepath.Join("testdata", "rake_no_gem"))
				Expect(err).NotTo(HaveOccurred())

				var logs fmt.Stringer
				image, logs, err = pack.WithNoColor().Build.
					WithBuildpacks(rubyBuildpack).
					WithPullPolicy("never").
					Execute(name, source)
				Expect(err).NotTo(HaveOccurred(), logs.String())

				container, err = docker.Container.Run.Execute(image.ID)
				Expect(err).NotTo(HaveOccurred())

				rLogs := func() fmt.Stringer {
					rakeLogs, err := docker.Container.Logs.Execute(container.ID)
					Expect(err).NotTo(HaveOccurred())
					return rakeLogs
				}

				Eventually(rLogs).Should(ContainSubstring("I am a rake task"))

				Expect(logs).To(ContainLines(ContainSubstring("MRI Buildpack")))
				Expect(logs).NotTo(ContainLines(ContainSubstring("Bundler Buildpack")))
				Expect(logs).NotTo(ContainLines(ContainSubstring("Bundle Install Buildpack")))
				Expect(logs).To(ContainLines(ContainSubstring("Rake Buildpack")))
				Expect(logs).To(ContainLines(ContainSubstring("rake")))
				Expect(logs).NotTo(ContainLines(ContainSubstring("Procfile Buildpack")))
				Expect(logs).NotTo(ContainLines(ContainSubstring("Environment Variables Buildpack")))
			})
		})

		context("using optional utility buildpacks", func() {
			it("creates a working OCI image that uses the start command from and includes environment buildpack functionality", func() {
				var err error
				source, err = occam.Source(filepath.Join("testdata", "rake"))
				Expect(err).NotTo(HaveOccurred())
				Expect(ioutil.WriteFile(filepath.Join(source, "Procfile"), []byte("web: bundle exec rake proc"), 0644)).To(Succeed())

				var logs fmt.Stringer
				image, logs, err = pack.WithNoColor().Build.
					WithBuildpacks(rubyBuildpack).
					WithEnv(map[string]string{"BPE_SOME_VARIABLE": "SOME_VALUE"}).
					WithPullPolicy("never").
					Execute(name, source)
				Expect(err).NotTo(HaveOccurred(), logs.String())

				container, err = docker.Container.Run.Execute(image.ID)
				Expect(err).NotTo(HaveOccurred())

				Expect(image.Buildpacks[5].Key).To(Equal("paketo-buildpacks/environment-variables"))
				Expect(image.Buildpacks[5].Layers["environment-variables"].Metadata["variables"]).To(Equal(map[string]interface{}{"SOME_VARIABLE": "SOME_VALUE"}))

				rLogs := func() fmt.Stringer {
					rakeLogs, err := docker.Container.Logs.Execute(container.ID)
					Expect(err).NotTo(HaveOccurred())
					return rakeLogs
				}

				Eventually(rLogs).Should(ContainSubstring("I am the proc rake task"))

				Expect(logs).To(ContainLines(ContainSubstring("MRI Buildpack")))
				Expect(logs).To(ContainLines(ContainSubstring("Bundler Buildpack")))
				Expect(logs).To(ContainLines(ContainSubstring("Bundle Install Buildpack")))
				Expect(logs).To(ContainLines(ContainSubstring("Rake Buildpack")))
				Expect(logs).To(ContainLines(ContainSubstring("Procfile Buildpack")))
				Expect(logs).To(ContainLines(ContainSubstring("bundle exec rake proc")))
				Expect(logs).To(ContainLines(ContainSubstring("Environment Variables Buildpack")))
			})
		})
	})
}
