// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package policy

import (
	"context"
	"fmt"

	v1 "github.com/marmotedu/api/apiserver/v1"
	"github.com/marmotedu/component-base/pkg/json"
	metav1 "github.com/marmotedu/component-base/pkg/meta/v1"
	apiclientv1 "github.com/marmotedu/marmotedu-sdk-go/marmotedu/service/iam/apiserver/v1"
	"github.com/ory/ladon"
	"github.com/spf13/cobra"

	cmdutil "github.com/marmotedu/iam/internal/iamctl/cmd/util"
	"github.com/marmotedu/iam/internal/iamctl/util/templates"
	"github.com/marmotedu/iam/pkg/cli/genericclioptions"
)

const (
	createUsageStr = "create POLICY_NAME POLICY"
)

// CreateOptions is an options struct to support create subcommands.
type CreateOptions struct {
	Policy *v1.Policy

	Client apiclientv1.APIV1Interface
	genericclioptions.IOStreams
}

var (
	createExample = templates.Examples(`
		# Create a authorization policy
		iamctl policy create foo '{"description":"This is a updated policy","subjects":["users:<peter|ken>","users:maria","groups:admins"],"actions":["delete","<create|update>"],"effect":"allow","resources":["resources:articles:<.*>","resources:printer"],"conditions":{"remoteIPAddress":{"type":"CIDRCondition","options":{"cidr":"192.168.0.1/16"}}}}'`)

	createUsageErrStr = fmt.Sprintf(
		"expected '%s'.\nPOLICY_NAME and POLICY are required arguments for the create command",
		createUsageStr,
	)
)

// NewCreateOptions returns an initialized CreateOptions instance.
func NewCreateOptions(ioStreams genericclioptions.IOStreams) *CreateOptions {
	return &CreateOptions{
		IOStreams: ioStreams,
	}
}

// NewCmdCreate returns new initialized instance of create sub command.
func NewCmdCreate(f cmdutil.Factory, ioStreams genericclioptions.IOStreams) *cobra.Command {
	o := NewCreateOptions(ioStreams)

	cmd := &cobra.Command{
		Use:                   createUsageStr,
		DisableFlagsInUseLine: true,
		Aliases:               []string{},
		Short:                 "Create a authorization policy resource",
		TraverseChildren:      true,
		Long:                  "Create a authorization policy resource.",
		Example:               createExample,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(f, cmd, args))
			cmdutil.CheckErr(o.Validate(cmd, args))
			cmdutil.CheckErr(o.Run(args))
		},
		SuggestFor: []string{},
	}

	return cmd
}

// Complete completes all the required options.
func (o *CreateOptions) Complete(f cmdutil.Factory, cmd *cobra.Command, args []string) error {
	if len(args) < 2 {
		return cmdutil.UsageErrorf(cmd, createUsageErrStr)
	}

	var pol ladon.DefaultPolicy
	if err := json.Unmarshal([]byte(args[1]), &pol); err != nil {
		return err
	}

	o.Policy = &v1.Policy{
		ObjectMeta: metav1.ObjectMeta{
			Name: args[0],
		},
		Policy: v1.AuthzPolicy{
			DefaultPolicy: pol,
		},
	}

	clientConfig, err := f.ToRESTConfig() // 生成配置
	if err != nil {
		return err
	}
	o.Client, err = apiclientv1.NewForConfig(clientConfig) // 基于 rest.Config 类型的配置变量 创建RESTful API客户端和SDK的客户端
	if err != nil {
		return err
	}

	return nil
}

// Validate makes sure there is no discrepency in command options.
func (o *CreateOptions) Validate(cmd *cobra.Command, args []string) error {
	if errs := o.Policy.Validate(); len(errs) != 0 {
		return errs.ToAggregate()
	}

	return nil
}

// Run executes a create subcommand using the specified options.
func (o *CreateOptions) Run(args []string) error {
	ret, err := o.Client.Policies().Create(context.TODO(), o.Policy, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	fmt.Fprintf(o.Out, "policy/%s created\n", ret.Name)

	return nil
}
