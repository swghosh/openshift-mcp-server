package core

import (
	"fmt"

	"github.com/google/jsonschema-go/jsonschema"
	"k8s.io/utils/ptr"

	"github.com/containers/kubernetes-mcp-server/pkg/api"
)

func initMustGather() []api.ServerTool {
	return []api.ServerTool{
		{Tool: api.Tool{
			Name:        "plan_mustgather",
			Description: "Plan a must-gather operation for Kubernetes cluster debugging and troubleshooting",
			InputSchema: &jsonschema.Schema{
				Type: "object",
				Properties: map[string]*jsonschema.Schema{
					"node_name": {
						Type:        "string",
						Description: "Optional node to run the mustgather pod. If not provided, a random control-plane node will be selected automatically",
					},
					"node_selector": {
						Type:        "string",
						Description: "Optional node selector to use - only relevant when specifying a command and image which needs to capture data on a set of cluster nodes simultaneously",
					},
					"host_network": {
						Type:        "boolean",
						Description: "Run must-gather pods as hostNetwork: true - relevant if a specific command and image needs to capture host-level data",
					},
					"images": {
						Type:        "array",
						Description: "Specify must-gather plugin images to run. If not specified, OpenShift's default must-gather image will be used",
						Items: &jsonschema.Schema{
							Type: "string",
						},
					},
					"image_streams": {
						Type:        "array",
						Description: "Specify image streams (namespace/name:tag) containing must-gather plugin images to run",
						Items: &jsonschema.Schema{
							Type: "string",
						},
					},
					"all_images": {
						Type:        "boolean",
						Description: "Collect must-gather using the default image for all Operators on the cluster annotated with must-gather",
					},
					"dest_dir": {
						Type:        "string",
						Description: "Set a specific directory on the local machine to write gathered data to",
					},
					"source_dir": {
						Type:        "string",
						Description: "Set the specific directory on the pod copy the gathered data from",
					},
					"timeout": {
						Type:        "string",
						Description: "The length of time to gather data, like 5s, 2m, or 3h, higher than zero. Defaults to 10 minutes",
						Default:     "10m",
					},
					"run_namespace": {
						Type:        "string",
						Description: "An existing privileged namespace where must-gather pods should run. If not specified a temporary namespace will be generated",
					},
					"volume_percentage": {
						Type:        "integer",
						Description: "Specify maximum percentage of must-gather pod's allocated volume that can be used. If this limit is exceeded, must-gather will stop gathering, but still copy gathered data",
						Minimum:     ptr.To(0.0),
						Maximum:     ptr.To(100.0),
					},
					"keep": {
						Type:        "boolean",
						Description: "Do not delete temporary resources when the mustgather completes",
					},
					"since_time": {
						Type:        "string",
						Description: "Only return logs after a specific date (RFC3339). Defaults to all logs. Only one of since-time / since may be used",
					},
					"since": {
						Type:        "string",
						Description: "Only return logs newer than a relative duration like 5s, 2m, or 3h. Defaults to all logs. Only one of since-time / since may be used",
					},
				},
			},
			Annotations: api.ToolAnnotations{
				Title:           "MustGather: Plan",
				ReadOnlyHint:    ptr.To(true),
				DestructiveHint: ptr.To(false),
				IdempotentHint:  ptr.To(true),
				OpenWorldHint:   ptr.To(false),
			},
		}, Handler: planMustGather},
	}
}

func planMustGather(params api.ToolHandlerParams) (*api.ToolCallResult, error) {
	args := params.GetArguments()

	// Extract arguments
	nodeName := getStringArg(args, "node_name")
	nodeSelector := getStringArg(args, "node_selector")
	hostNetwork := getBoolArg(args, "host_network")
	images := getStringSliceArg(args, "images")
	imageStreams := getStringSliceArg(args, "image_streams")
	allImages := getBoolArg(args, "all_images")
	destDir := getStringArg(args, "dest_dir")
	sourceDir := getStringArg(args, "source_dir")
	timeout := getStringArg(args, "timeout")
	runNamespace := getStringArg(args, "run_namespace")
	volumePercentage := getIntArg(args, "volume_percentage")
	keep := getBoolArg(args, "keep")
	sinceTime := getStringArg(args, "since_time")
	since := getStringArg(args, "since")

	if timeout == "" {
		timeout = "10m"
	}

	plan := "# Must-Gather Plan\n\n"
	plan += "This is a dummy implementation of the plan_mustgather tool.\n\n"
	plan += "## Planned Collection Configuration:\n\n"

	// Node configuration
	if nodeName != "" {
		plan += "- **Target Node**: " + nodeName + "\n"
	}
	if nodeSelector != "" {
		plan += "- **Node Selector**: " + nodeSelector + "\n"
	}
	if hostNetwork {
		plan += "- **Host Network**: Enabled (pods will run with hostNetwork: true)\n"
	}

	// Image configuration
	if allImages {
		plan += "- **Images**: All operator images (annotated with must-gather)\n"
	} else if len(images) > 0 {
		plan += "- **Images**:\n"
		for _, img := range images {
			plan += "  - " + img + "\n"
		}
	} else {
		plan += "- **Images**: Default OpenShift must-gather image\n"
	}

	if len(imageStreams) > 0 {
		plan += "- **Image Streams**:\n"
		for _, is := range imageStreams {
			plan += "  - " + is + "\n"
		}
	}

	// Runtime configuration
	plan += "- **Timeout**: " + timeout + "\n"
	if runNamespace != "" {
		plan += "- **Run Namespace**: " + runNamespace + "\n"
	} else {
		plan += "- **Run Namespace**: Temporary namespace (will be auto-generated)\n"
	}

	// Storage configuration
	if destDir != "" {
		plan += "- **Destination Directory**: " + destDir + "\n"
	}
	if sourceDir != "" {
		plan += "- **Source Directory**: " + sourceDir + "\n"
	}
	if volumePercentage > 0 {
		plan += "- **Volume Percentage Limit**: " + fmt.Sprintf("%d%%", volumePercentage) + "\n"
	}

	// Time filtering
	if sinceTime != "" {
		plan += "- **Log Filter (Since Time)**: " + sinceTime + "\n"
	}
	if since != "" {
		plan += "- **Log Filter (Since Duration)**: " + since + "\n"
	}

	// Cleanup configuration
	if keep {
		plan += "- **Resource Cleanup**: Disabled (resources will be kept after completion)\n"
	} else {
		plan += "- **Resource Cleanup**: Enabled (temporary resources will be deleted)\n"
	}

	plan += "\n## Standard Collection Includes:\n"
	plan += "- Cluster info and version\n"
	plan += "- Node information and status\n"
	plan += "- Pod logs and descriptions\n"
	plan += "- Events across all namespaces\n"
	plan += "- ConfigMaps and Secrets (metadata only)\n"
	plan += "- Network policies and configurations\n"
	plan += "- Storage information and PVCs\n"
	plan += "- Operator-specific diagnostics (if applicable)\n"

	plan += "\n## Next Steps:\n"
	plan += "1. Review the planned collection configuration\n"
	plan += "2. Execute the must-gather operation with specified parameters\n"
	plan += "3. Monitor collection progress and resource usage\n"
	plan += "4. Analyze collected data for troubleshooting\n"
	plan += "\n*Note: This is a dummy tool for demonstration purposes.*"

	return api.NewToolCallResult(plan, nil), nil
}

// Helper functions to extract typed arguments
func getStringArg(args map[string]any, key string) string {
	if val, ok := args[key]; ok && val != nil {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

func getBoolArg(args map[string]any, key string) bool {
	if val, ok := args[key]; ok && val != nil {
		if b, ok := val.(bool); ok {
			return b
		}
	}
	return false
}

func getIntArg(args map[string]any, key string) int {
	if val, ok := args[key]; ok && val != nil {
		if i, ok := val.(float64); ok {
			return int(i)
		}
	}
	return 0
}

func getStringSliceArg(args map[string]any, key string) []string {
	if val, ok := args[key]; ok && val != nil {
		if slice, ok := val.([]interface{}); ok {
			var result []string
			for _, item := range slice {
				if str, ok := item.(string); ok {
					result = append(result, str)
				}
			}
			return result
		}
	}
	return nil
}
