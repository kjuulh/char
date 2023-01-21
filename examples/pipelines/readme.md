# Pipelines design docs

The goal of this experiment is to play with various code sharing features. The
pipelines in this experiment is supposed to model normal code sharing behavior
of libraries, with overriding capabilities of the downstream repository

The goal is to split the body of the work in three parts.

- libraries,
- compendiums
- articles

The terminology is as such:

## Libraries

Libraries provide raw functions, and is a general abstraction on an underlying
process, such as running a container, executing a shell script etc. Libraries
might define an API descripting how to interact with it. It is up to the
compendium, and article to use these in a sane manner.

Such that:

```rust
pub fn execute_shell(&self, input: &ShellOpts) -> Result<ShellOutput> {
   ...
}
```

These work similar to raw primitive functions, but should serve as an
opinionated flyweight. These may be extremely specific, such as building a go
binary, creating a github release etc. The details are left to the caller.

A version scheme of the library should follow semver, as that is the best model
at the moment for versioning. Libraries should be pulled using the native
package manager, or include as a submodule.

## Compendiums

A compendium is an opinionated collection of libraries, which consists of files,
configurations etc. A compendium is to be used by either other compendiums, or
articles. They are not to be used by libraries, this is to provide a natural
hierachy. The end result should be a directed acyclic graph.

A compendium should provide an API for either other compendiums, or articles.
These apis, need to remain flexible, and be open to mutations. As such all
primitive features, files, configurations need to be exposed as raw data or data
structures if suitable.

This is done using pipelines, or middleware for the specific parts. A data
object will pass from above, containing the data to implement the required
interfaces of the Compendium, these must be fulfilled for the construction of
the Compendium, else the compilation should fail.

The caller will have the ability to replace the specifics of the Compendium, by
replacing certain pipelines, or mutating the data. In case of mutations, the
data is only modified for that pipeline and below, if a fork occurs above, then
the divergent paths won't be affected.

Compendiums are driven by pipelines applied to them either from other
compendiums or articles. An article or compendium will only have access to their
direct dependent compendiums pipelines. An article will naturally expose
pipelines to be called by the user.

```rust
pub struct GoApplication;

impl Compendium for GoApplication {
  type Input = GoApplicationOpts;

  pub fn get_pipelines(&mut self) -> Result<Pipelines> {
    let pipelines = self.pipelines
      .clone()
      .add(self.get_application_pipelines())?
      .add(self.get_go_releaser_pipelines())?;

    Ok(pipelines)
  }
}
```

## Articles

An article is a specific implementation of a compendium, it is the end of the
chain, and is meant to be directly executed by the user, using a client
application.

It by default is supposed to be a golden path, I.e. it passes the defaults of
the Compendium, but on a case-by-case basis has the ability to modify its
pipelines to its needs. This may be changing certain default configurations,
mutate a dockerfile, add additional steps to a pipeline etc, and remove others.

We reason that once you stray from the golden path, you should be in control,
this may be done by forking the compendium's features.

It provides pipelines as actions, and implements a strict protocol for
communication.

```rust
pub struct MyGoService;

impl Article for MyGoService {

  fn get_pipelines(&mut self) -> Result<Pipelines> {
    let mut go_pipelines = GoApplication::new().get_pipelines()?;
    let api_pipeline = self.get_api_pipeline()?;
    let build_pipeline = self.get_docker_pipeline(go_pipelines.extract::<BuildPipeline>()?)?;;

    let pipelines = PipelineBuilder::new()
      .append(go_pipelines)
      .append(api_pipeline)
      .append(build_pipeline)
      .build()?

    Ok(pipelines)
  }
}
```

## Usage

A host app can now call these:

```bash
char ls
```

ls simply displays the information on what pipeline are available

```bash
char run build
```

run build will execute the pipeline build. It will validate input available from
char.toml, push these keys/values through to the pipeline, which will go through
all the steps.

- build
- MyGoService
- DockerBuildPipeline
  - (BuildPipeline)
- DockerLibrary
  - `fn docker_build(dockerfile_contents: string)`
- Native dependencies
  - "write dockerfile /tmp/abc/dockerfile"
  - "shell -> docker build -f /tmp/abc/dockerfile some-library-path"

## Common scenarios

### Replacing parts of a build script (dockerfile)

A compendium may embed/provide a dockerfile or any other resource, these are
provided through consistent interfaces.

```rust
pub struct GoDockerBuildPipeline;

impl GoDockerBuildPipeline {
  fn get_resources(&self) -> Result<(
    libraries::docker::DockerContents,
    libraries::docker::DockerBuildTags)> {
    return (self.contents, self.build_tags)
  }

  fn mutate_resources(
    &mut self,
    contents: libraries::docker::DockerContents,
    build_tags: libraries::docker::DockerBuildTags,
  ) -> Result<()> {
    self.contents = contents;
    self.build_tags = build_tags
  }
}

impl BuildPipeline for GoDockerBuildPipeline {
  fn execute(&mut self, config: Configuration) -> Result<()> {
    let (docker_contents, base_tags) = self.get_resources()
    libraries::docker::build(docker_contents, base_tags, config)
  }
}
```

In the article you can now replace the resources to fit your needs, that or
building your own pipeline.

```rust
pub struct MyGoService;

impl Article for MyGoService {

  fn get_pipelines(&mut self) -> Result<Pipelines> {
  let go_app = GoApplication::new();
  let go_pipelines = go_app.get_pipelines();

  let mut go_build_pipeline = go_pipelines.get_pipeline::<GoDockerBuildPipeline>()?;
  let go_resources = go_build_pipeline.get_resources()?;
  let go_resources = self.mutate_go_resources(&go_resources)?;
  go_build_pipeline.mutate_resources(go_resources)?;
  go_pipelines.replace::<GoDockerBuildPipeline>(go_build_pipeline)?;

  let pipelines = PipelineBuilder::new()
      .append(.get_pipelines())
      .build()?

    Ok(pipelines)
  }
}
```
