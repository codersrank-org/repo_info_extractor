import os
from language.C import extract_libraries

def test_extract_libraries():
    dir_path = os.path.dirname(os.path.realpath(__file__))
    files = ['fixtures/Cs.cs']
    fq_files = [os.path.join(dir_path, f) for f in files]
    libs = extract_libraries(fq_files)
    assert len(libs) == 1

    assert libs['C#'] == [
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
        "GitHub.DistributedTask.Pipelines.ObjectTemplating.PipelineTemplateConstants"
        ]
