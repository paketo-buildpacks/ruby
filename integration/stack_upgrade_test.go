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

func testGracefulStackUpgrades(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect     = NewWithT(t).Expect
		Eventually = NewWithT(t).Eventually

		imageIDs     map[string]struct{}
		containerIDs map[string]struct{}

		pack   occam.Pack
		docker occam.Docker

		name   string
		source string
	)

	it.Before(func() {
		var err error
		name, err = occam.RandomName()
		Expect(err).NotTo(HaveOccurred())

		pack = occam.NewPack()
		docker = occam.NewDocker()

		imageIDs = map[string]struct{}{}
		containerIDs = map[string]struct{}{}

		Expect(docker.Pull.Execute("paketobuildpacks/builder-jammy-buildpackless-full")).To(Succeed())
		Expect(docker.Pull.Execute("paketobuildpacks/run-jammy-full")).To(Succeed())
	})

	it.After(func() {
		for id := range containerIDs {
			Expect(docker.Container.Remove.Execute(id)).To(Succeed())
		}

		for id := range imageIDs {
			Expect(docker.Image.Remove.Execute(id)).To(Succeed())
		}

		Expect(docker.Image.Remove.Execute("paketobuildpacks/builder-jammy-buildpackless-full")).To(Succeed())
		Expect(docker.Image.Remove.Execute("paketobuildpacks/run-jammy-full")).To(Succeed())

		Expect(docker.Volume.Remove.Execute(occam.CacheVolumeNames(name))).To(Succeed())
		Expect(os.RemoveAll(source)).To(Succeed())
	})

	context("when building a rails app and the stack changes between builds", func() {
		var (
			err         error
			logs        fmt.Stringer
			firstImage  occam.Image
			secondImage occam.Image

			firstContainer  occam.Container
			secondContainer occam.Container
		)

		it.Before(func() {
			source, err = occam.Source(filepath.Join("testdata", "rails"))
			Expect(err).NotTo(HaveOccurred())
		})

		it("creates a working OCI image on rebuild", func() {
			build := pack.WithNoColor().Build.
				WithPullPolicy("never").
				WithBuildpacks(rubyBuildpack)

			firstImage, logs, err = build.Execute(name, source)
			Expect(err).NotTo(HaveOccurred(), logs.String())

			imageIDs[firstImage.ID] = struct{}{}

			firstContainer, err = docker.Container.Run.
				WithEnv(map[string]string{
					"PORT":            "8080",
					"SECRET_KEY_BASE": "some-secret",
				}).
				WithPublish("8080").
				WithPublishAll().
				Execute(firstImage.ID)
			Expect(err).NotTo(HaveOccurred())

			containerIDs[firstContainer.ID] = struct{}{}
			Eventually(firstContainer).Should(Serve(ContainSubstring("Hello World!")).OnPort(8080))

			// Second pack build, upgrade stack image
			secondImage, logs, err = build.WithBuilder("paketobuildpacks/builder-jammy-buildpackless-full").Execute(name, source)
			Expect(err).NotTo(HaveOccurred(), logs.String())

			imageIDs[secondImage.ID] = struct{}{}

			secondContainer, err = docker.Container.Run.
				WithEnv(map[string]string{
					"PORT":            "8080",
					"SECRET_KEY_BASE": "some-secret",
				}).
				WithPublish("8080").
				WithPublishAll().
				Execute(secondImage.ID)
			Expect(err).NotTo(HaveOccurred())

			containerIDs[secondContainer.ID] = struct{}{}
			Eventually(secondContainer).Should(Serve(ContainSubstring("Hello World!")).OnPort(8080))
		})
	})
}
