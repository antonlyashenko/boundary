// Code generated by "make cli"; DO NOT EDIT.
package credentiallibrariescmd

import (
	"errors"
	"fmt"

	"github.com/hashicorp/boundary/api"
	"github.com/hashicorp/boundary/api/credentiallibraries"
	"github.com/hashicorp/boundary/internal/cmd/base"
	"github.com/hashicorp/boundary/internal/cmd/common"
	"github.com/hashicorp/go-secure-stdlib/strutil"
	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

func initVaultSshCertificateFlags() {
	flagsOnce.Do(func() {
		extraFlags := extraVaultSshCertificateActionsFlagsMapFunc()
		for k, v := range extraFlags {
			flagsVaultSshCertificateMap[k] = append(flagsVaultSshCertificateMap[k], v...)
		}
	})
}

var (
	_ cli.Command             = (*VaultSshCertificateCommand)(nil)
	_ cli.CommandAutocomplete = (*VaultSshCertificateCommand)(nil)
)

type VaultSshCertificateCommand struct {
	*base.Command

	Func string

	plural string

	extraVaultSshCertificateCmdVars
}

func (c *VaultSshCertificateCommand) AutocompleteArgs() complete.Predictor {
	initVaultSshCertificateFlags()
	return complete.PredictAnything
}

func (c *VaultSshCertificateCommand) AutocompleteFlags() complete.Flags {
	initVaultSshCertificateFlags()
	return c.Flags().Completions()
}

func (c *VaultSshCertificateCommand) Synopsis() string {
	if extra := extraVaultSshCertificateSynopsisFunc(c); extra != "" {
		return extra
	}

	synopsisStr := "credential library"

	synopsisStr = fmt.Sprintf("%s %s", "vault-ssh-certificate-type", synopsisStr)

	return common.SynopsisFunc(c.Func, synopsisStr)
}

func (c *VaultSshCertificateCommand) Help() string {
	initVaultSshCertificateFlags()

	var helpStr string
	helpMap := common.HelpMap("credential library")

	switch c.Func {

	default:

		helpStr = c.extraVaultSshCertificateHelpFunc(helpMap)

	}

	// Keep linter from complaining if we don't actually generate code using it
	_ = helpMap
	return helpStr
}

var flagsVaultSshCertificateMap = map[string][]string{

	"create": {"credential-store-id", "name", "description"},

	"update": {"id", "name", "description", "version"},
}

func (c *VaultSshCertificateCommand) Flags() *base.FlagSets {
	if len(flagsVaultSshCertificateMap[c.Func]) == 0 {
		return c.FlagSet(base.FlagSetNone)
	}

	set := c.FlagSet(base.FlagSetHTTP | base.FlagSetClient | base.FlagSetOutputFormat)
	f := set.NewFlagSet("Command Options")
	common.PopulateCommonFlags(c.Command, f, "vault-ssh-certificate-type credential library", flagsVaultSshCertificateMap, c.Func)

	extraVaultSshCertificateFlagsFunc(c, set, f)

	return set
}

func (c *VaultSshCertificateCommand) Run(args []string) int {
	initVaultSshCertificateFlags()

	switch c.Func {
	case "":
		return cli.RunResultHelp

	}

	c.plural = "vault-ssh-certificate-type credential library"
	switch c.Func {
	case "list":
		c.plural = "vault-ssh-certificate-type credential libraries"
	}

	f := c.Flags()

	if err := f.Parse(args); err != nil {
		c.PrintCliError(err)
		return base.CommandUserError
	}

	if strutil.StrListContains(flagsVaultSshCertificateMap[c.Func], "id") && c.FlagId == "" {
		c.PrintCliError(errors.New("ID is required but not passed in via -id"))
		return base.CommandUserError
	}

	var opts []credentiallibraries.Option

	if strutil.StrListContains(flagsVaultSshCertificateMap[c.Func], "credential-store-id") {
		switch c.Func {

		case "create":
			if c.FlagCredentialStoreId == "" {
				c.PrintCliError(errors.New("CredentialStore ID must be passed in via -credential-store-id or BOUNDARY_CREDENTIAL_STORE_ID"))
				return base.CommandUserError
			}

		}
	}

	client, err := c.Client()
	if c.WrapperCleanupFunc != nil {
		defer func() {
			if err := c.WrapperCleanupFunc(); err != nil {
				c.PrintCliError(fmt.Errorf("Error cleaning kms wrapper: %w", err))
			}
		}()
	}
	if err != nil {
		c.PrintCliError(fmt.Errorf("Error creating API client: %w", err))
		return base.CommandCliError
	}
	credentiallibrariesClient := credentiallibraries.NewClient(client)

	switch c.FlagName {
	case "":
	case "null":
		opts = append(opts, credentiallibraries.DefaultName())
	default:
		opts = append(opts, credentiallibraries.WithName(c.FlagName))
	}

	switch c.FlagDescription {
	case "":
	case "null":
		opts = append(opts, credentiallibraries.DefaultDescription())
	default:
		opts = append(opts, credentiallibraries.WithDescription(c.FlagDescription))
	}

	if c.FlagFilter != "" {
		opts = append(opts, credentiallibraries.WithFilter(c.FlagFilter))
	}

	var version uint32

	switch c.Func {

	case "update":
		switch c.FlagVersion {
		case 0:
			opts = append(opts, credentiallibraries.WithAutomaticVersioning(true))
		default:
			version = uint32(c.FlagVersion)
		}

	}

	if ok := extraVaultSshCertificateFlagsHandlingFunc(c, f, &opts); !ok {
		return base.CommandUserError
	}

	var resp *api.Response
	var item *credentiallibraries.CredentialLibrary

	var createResult *credentiallibraries.CredentialLibraryCreateResult

	var updateResult *credentiallibraries.CredentialLibraryUpdateResult

	switch c.Func {

	case "create":
		createResult, err = credentiallibrariesClient.Create(c.Context, "vault-ssh-certificate", c.FlagCredentialStoreId, opts...)
		if exitCode := c.checkFuncError(err); exitCode > 0 {
			return exitCode
		}
		resp = createResult.GetResponse()
		item = createResult.GetItem()

	case "update":
		updateResult, err = credentiallibrariesClient.Update(c.Context, c.FlagId, version, opts...)
		if exitCode := c.checkFuncError(err); exitCode > 0 {
			return exitCode
		}
		resp = updateResult.GetResponse()
		item = updateResult.GetItem()

	}

	resp, item, err = executeExtraVaultSshCertificateActions(c, resp, item, err, credentiallibrariesClient, version, opts)
	if exitCode := c.checkFuncError(err); exitCode > 0 {
		return exitCode
	}

	output, err := printCustomVaultSshCertificateActionOutput(c)
	if err != nil {
		c.PrintCliError(err)
		return base.CommandUserError
	}
	if output {
		return base.CommandSuccess
	}

	switch c.Func {

	}

	switch base.Format(c.UI) {
	case "table":
		c.UI.Output(printItemTable(item, resp))

	case "json":
		if ok := c.PrintJsonItem(resp); !ok {
			return base.CommandCliError
		}
	}

	return base.CommandSuccess
}

func (c *VaultSshCertificateCommand) checkFuncError(err error) int {
	if err == nil {
		return 0
	}
	if apiErr := api.AsServerError(err); apiErr != nil {
		c.PrintApiError(apiErr, fmt.Sprintf("Error from controller when performing %s on %s", c.Func, c.plural))
		return base.CommandApiError
	}
	c.PrintCliError(fmt.Errorf("Error trying to %s %s: %s", c.Func, c.plural, err.Error()))
	return base.CommandCliError
}

var (
	extraVaultSshCertificateActionsFlagsMapFunc = func() map[string][]string { return nil }
	extraVaultSshCertificateSynopsisFunc        = func(*VaultSshCertificateCommand) string { return "" }
	extraVaultSshCertificateFlagsFunc           = func(*VaultSshCertificateCommand, *base.FlagSets, *base.FlagSet) {}
	extraVaultSshCertificateFlagsHandlingFunc   = func(*VaultSshCertificateCommand, *base.FlagSets, *[]credentiallibraries.Option) bool { return true }
	executeExtraVaultSshCertificateActions      = func(_ *VaultSshCertificateCommand, inResp *api.Response, inItem *credentiallibraries.CredentialLibrary, inErr error, _ *credentiallibraries.Client, _ uint32, _ []credentiallibraries.Option) (*api.Response, *credentiallibraries.CredentialLibrary, error) {
		return inResp, inItem, inErr
	}
	printCustomVaultSshCertificateActionOutput = func(*VaultSshCertificateCommand) (bool, error) { return false, nil }
)
