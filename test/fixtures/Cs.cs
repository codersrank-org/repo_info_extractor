using System;
using System.Collections.Generic;
using System.IO;
using System.IO.Compression;
using System.Linq;
using System.Net;
using System.Net.Http;
using System.Net.Http.Headers;
using System.Text;
using System.Threading;
using System.Threading.Tasks;
using GitHub.DistributedTask.ObjectTemplating.Tokens;
using GitHub.Runner.Common;
using GitHub.Runner.Sdk;
using GitHub.Runner.Worker.Container;
using GitHub.Services.Common;
using WebApi = GitHub.DistributedTask.WebApi;
using Pipelines = GitHub.DistributedTask.Pipelines;
using PipelineTemplateConstants = GitHub.DistributedTask.Pipelines.ObjectTemplating.PipelineTemplateConstants;

namespace GitHub.Runner.Worker
{
    public class PrepareResult
    {
        public PrepareResult(List<JobExtensionRunner> containerSetupSteps, Dictionary<Guid, IActionRunner> preStepTracker)
        {
            this.ContainerSetupSteps = containerSetupSteps;
            this.PreStepTracker = preStepTracker;
        }

        public List<JobExtensionRunner> ContainerSetupSteps { get; set; }

        public Dictionary<Guid, IActionRunner> PreStepTracker { get; set; }
    }

    [ServiceLocator(Default = typeof(ActionManager))]
    public interface IActionManager : IRunnerService
    {
        Dictionary<Guid, ContainerInfo> CachedActionContainers { get; }
        Task<PrepareResult> PrepareActionsAsync(IExecutionContext executionContext, IEnumerable<Pipelines.JobStep> steps);
        Definition LoadAction(IExecutionContext executionContext, Pipelines.ActionStep action);
    }

    public sealed class ActionManager : RunnerService, IActionManager
    {
        private const int _defaultFileStreamBufferSize = 4096;

        //81920 is the default used by System.IO.Stream.CopyTo and is under the large object heap threshold (85k).
        private const int _defaultCopyBufferSize = 81920;
        private const string _dotcomApiUrl = "https://api.github.com";
        private readonly Dictionary<Guid, ContainerInfo> _cachedActionContainers = new Dictionary<Guid, ContainerInfo>();

        public Dictionary<Guid, ContainerInfo> CachedActionContainers => _cachedActionContainers;
        public async Task<PrepareResult> PrepareActionsAsync(IExecutionContext executionContext, IEnumerable<Pipelines.JobStep> steps)
        {
            ArgUtil.NotNull(executionContext, nameof(executionContext));
            ArgUtil.NotNull(steps, nameof(steps));

            executionContext.Output("Prepare all required actions");
            Dictionary<string, List<Guid>> imagesToPull = new Dictionary<string, List<Guid>>(StringComparer.OrdinalIgnoreCase);
            Dictionary<string, List<Guid>> imagesToBuild = new Dictionary<string, List<Guid>>(StringComparer.OrdinalIgnoreCase);
            Dictionary<string, ActionContainer> imagesToBuildInfo = new Dictionary<string, ActionContainer>(StringComparer.OrdinalIgnoreCase);
            List<JobExtensionRunner> containerSetupSteps = new List<JobExtensionRunner>();
            Dictionary<Guid, IActionRunner> preStepTracker = new Dictionary<Guid, IActionRunner>();
            IEnumerable<Pipelines.ActionStep> actions = steps.OfType<Pipelines.ActionStep>();

            // TODO: Deprecate the PREVIEW_ACTION_TOKEN
            // Log even if we aren't using it to ensure users know.
            if (!string.IsNullOrEmpty(executionContext.Global.Variables.Get("PREVIEW_ACTION_TOKEN")))
            {
                executionContext.Warning("The 'PREVIEW_ACTION_TOKEN' secret is deprecated. Please remove it from the repository's secrets");
            }

            // Clear the cache (for self-hosted runners)
            IOUtil.DeleteDirectory(HostContext.GetDirectory(WellKnownDirectory.Actions), executionContext.CancellationToken);

            // todo: Remove when feature flag DistributedTask.NewActionMetadata is removed
            var newActionMetadata = executionContext.Global.Variables.GetBoolean("DistributedTask.NewActionMetadata") ?? false;

            var repositoryActions = new List<Pipelines.ActionStep>();

            foreach (var action in actions)
            {
                if (action.Reference.Type == Pipelines.ActionSourceType.ContainerRegistry)
                {
                    ArgUtil.NotNull(action, nameof(action));
                    var containerReference = action.Reference as Pipelines.ContainerRegistryReference;
                    ArgUtil.NotNull(containerReference, nameof(containerReference));
                    ArgUtil.NotNullOrEmpty(containerReference.Image, nameof(containerReference.Image));

                    if (!imagesToPull.ContainsKey(containerReference.Image))
                    {
                        imagesToPull[containerReference.Image] = new List<Guid>();
                    }

                    Trace.Info($"Action {action.Name} ({action.Id}) needs to pull image '{containerReference.Image}'");
                    imagesToPull[containerReference.Image].Add(action.Id);
                }
                // todo: Remove when feature flag DistributedTask.NewActionMetadata is removed
                else if (action.Reference.Type == Pipelines.ActionSourceType.Repository && !newActionMetadata)
                {
                    // only download the repository archive
                    await DownloadRepositoryActionAsync(executionContext, action);

                    // more preparation base on content in the repository (action.yml)
                    var setupInfo = PrepareRepositoryActionAsync(executionContext, action);
                    if (setupInfo != null)
                    {
                        if (!string.IsNullOrEmpty(setupInfo.Image))
                        {
                            if (!imagesToPull.ContainsKey(setupInfo.Image))
                            {
                                imagesToPull[setupInfo.Image] = new List<Guid>();
                            }

                            Trace.Info($"Action {action.Name} ({action.Id}) from repository '{setupInfo.ActionRepository}' needs to pull image '{setupInfo.Image}'");
                            imagesToPull[setupInfo.Image].Add(action.Id);
                        }
                        else
                        {
                            ArgUtil.NotNullOrEmpty(setupInfo.ActionRepository, nameof(setupInfo.ActionRepository));

                            if (!imagesToBuild.ContainsKey(setupInfo.ActionRepository))
                            {
                                imagesToBuild[setupInfo.ActionRepository] = new List<Guid>();
                            }

                            Trace.Info($"Action {action.Name} ({action.Id}) from repository '{setupInfo.ActionRepository}' needs to build image '{setupInfo.Dockerfile}'");
                            imagesToBuild[setupInfo.ActionRepository].Add(action.Id);
                            imagesToBuildInfo[setupInfo.ActionRepository] = setupInfo;
                        }
                    }