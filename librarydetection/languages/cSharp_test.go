package languages_test

import (
	. "github.com/onsi/ginkgo"
	"io/ioutil"

	"github.com/codersrank-org/repo_info_extractor/librarydetection/languages"
)

var _ = Describe("GoLibraryDetection", func() {
	fixture, err := ioutil.ReadFile("./fixtures/csharp.fixture")
	if err != nil {
		panic(err)
	}

	expectedLibraries := []string{
		"System",
		"System.Collections.Generic",
		"System.IO",
		"System.IO.Compression",
		"System.Linq",
		"System.Net",
		"System.Net.Http",
		"System.Net.Http.Headers",
		"System.Text",
		"System.Threading",
		"System.Threading.Tasks",
		"GitHub.DistributedTask.ObjectTemplating.Tokens",
		"GitHub.Runner.Common",
		"GitHub.Runner.Sdk",
		"GitHub.Runner.Worker.Container",
		"GitHub.Services.Common",
		"GitHub.DistributedTask.WebApi",
		"GitHub.DistributedTask.Pipelines",
		"GitHub.DistributedTask.Pipelines.ObjectTemplating.PipelineTemplateConstants",
	}

	analyzer := languages.NewCSharpAnalyzer()

	Describe("Extract C# Libraries", func() {
		It("Should be able to extract libraries", func() {
			libs, err := analyzer.ExtractLibraries(string(fixture))
			if err != nil {
				panic(err)
			}
			assertSameUnordered(libs, expectedLibraries)
		})
	})
})
