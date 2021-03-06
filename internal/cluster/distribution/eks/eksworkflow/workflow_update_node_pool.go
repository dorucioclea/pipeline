// Copyright © 2020 Banzai Cloud
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package eksworkflow

import (
	"fmt"
	"time"

	"go.uber.org/cadence"
	"go.uber.org/cadence/workflow"

	"github.com/banzaicloud/pipeline/internal/cluster"
	"github.com/banzaicloud/pipeline/pkg/sdk/process"
)

const UpdateNodePoolWorkflowName = "eks-update-node-pool"

type UpdateNodePoolWorkflowInput struct {
	ProviderSecretID string
	Region           string

	StackName string

	OrganizationID  uint
	ClusterID       uint
	ClusterSecretID string
	ClusterName     string
	NodePoolName    string

	NodeImage string
}

func UpdateNodePoolWorkflow(ctx workflow.Context, input UpdateNodePoolWorkflowInput) (err error) {
	activityOptions := workflow.ActivityOptions{
		ScheduleToStartTimeout: time.Duration(workflow.GetInfo(ctx).ExecutionStartToCloseTimeoutSeconds) * time.Second,
	}

	ctx = workflow.WithActivityOptions(ctx, activityOptions)

	proc := process.Start(
		workflow.WithStartToCloseTimeout(ctx, 10*time.Minute),
		input.OrganizationID,
		fmt.Sprint(input.ClusterID),
	)
	defer proc.RecordEnd(err)
	defer func() {
		status := cluster.Running
		statusMessage := cluster.RunningMessage

		if err != nil {
			status = cluster.Warning
			statusMessage = fmt.Sprintf("failed to update node pool: %s", err.Error())
		}

		_ = setClusterStatus(ctx, input.ClusterID, status, statusMessage)
	}()

	{
		activityInput := UpdateNodeGroupActivityInput{
			SecretID:     input.ProviderSecretID,
			Region:       input.Region,
			ClusterName:  input.ClusterName,
			StackName:    input.StackName,
			NodePoolName: input.NodePoolName,
			NodeImage:    input.NodeImage,
		}

		activityOptions := activityOptions
		activityOptions.StartToCloseTimeout = 5 * time.Minute
		activityOptions.RetryPolicy = &cadence.RetryPolicy{
			InitialInterval:          20 * time.Second,
			BackoffCoefficient:       1.1,
			MaximumAttempts:          10,
			NonRetriableErrorReasons: []string{"cadenceInternal:Panic", ErrReasonStackFailed},
		}

		var output UpdateNodeGroupActivityOutput

		processEvent := process.NewEvent(workflow.WithStartToCloseTimeout(ctx, 10*time.Minute), UpdateNodeGroupActivityName)
		err = workflow.ExecuteActivity(
			workflow.WithActivityOptions(ctx, activityOptions),
			UpdateNodeGroupActivityName,
			activityInput,
		).Get(ctx, &output)
		processEvent.RecordEnd(err)
		if err != nil || !output.NodePoolChanged {
			return
		}
	}

	{
		activityInput := WaitCloudFormationStackUpdateActivityInput{
			SecretID:  input.ProviderSecretID,
			Region:    input.Region,
			StackName: input.StackName,
		}

		activityOptions := activityOptions
		activityOptions.StartToCloseTimeout = 100 * 10 * time.Minute // TODO: calculate based on desired node count (limited to around 100 nodes now)
		activityOptions.HeartbeatTimeout = time.Minute
		activityOptions.RetryPolicy = &cadence.RetryPolicy{
			InitialInterval:          20 * time.Second,
			BackoffCoefficient:       1.1,
			MaximumAttempts:          20,
			NonRetriableErrorReasons: []string{"cadenceInternal:Panic"},
		}

		processEvent := process.NewEvent(workflow.WithStartToCloseTimeout(ctx, 10*time.Minute), WaitCloudFormationStackUpdateActivityName)
		err = workflow.ExecuteActivity(
			workflow.WithActivityOptions(ctx, activityOptions),
			WaitCloudFormationStackUpdateActivityName,
			activityInput,
		).Get(ctx, nil)
		processEvent.RecordEnd(err)
		if err != nil {
			return err
		}
	}

	return nil
}
