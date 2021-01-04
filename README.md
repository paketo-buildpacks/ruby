# Ruby Paketo Buildpack

## `gcr.io/paketo-buildpacks/ruby`

The Ruby Paketo Buildpack provides a set of collaborating buildpacks that
enable the building of a Ruby-based application. These buildpacks include:
- [Bundle Install](https://github.com/paketo-buildpacks/bundle-install)
- [Bundler](https://github.com/paketo-buildpacks/bundler)
- [MRI](https://github.com/paketo-buildpacks/mri)
- [Node Engine](https://github.com/paketo-buildpacks/node-engine)
- [Passenger](https://github.com/paketo-buildpacks/passenger)
- [Puma](https://github.com/paketo-buildpacks/puma)
- [Rackup](https://github.com/paketo-buildpacks/rackup)
- [Rails Assets](https://github.com/paketo-buildpacks/rails-assets)
- [Rake](https://github.com/paketo-buildpacks/rake)
- [Thin](https://github.com/paketo-buildpacks/thin)
- [Unicorn](https://github.com/paketo-buildpacks/unicorn)
- [Yarn Install](https://github.com/paketo-buildpacks/yarn-install)
- [Yarn](https://github.com/paketo-buildpacks/yarn)

The buildpack supports building simple Ruby applications or applications which
utilize [Bundler](https://bundler.io/) for managing their dependencies. Usage
examples can be found in the
[`samples` repository under the `ruby` directory](https://github.com/paketo-buildpacks/samples/tree/main/ruby).

#### The Ruby buildpack is compatible with the following builder(s):
- [Paketo Full Builder](https://github.com/paketo-buildpacks/full-builder)
- [Paketo Base Builder](https://github.com/paketo-buildpacks/base-builder)
