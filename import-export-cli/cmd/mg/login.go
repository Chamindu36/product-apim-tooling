/*
*  Copyright (c) WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
*
*  WSO2 Inc. licenses this file to you under the Apache License,
*  Version 2.0 (the "License"); you may not use this file except
*  in compliance with the License.
*  You may obtain a copy of the License at
*
*    http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing,
* software distributed under the License is distributed on an
* "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
* KIND, either express or implied.  See the License for the
* specific language governing permissions and limitations
* under the License.
 */

package mg

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/credentials"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
	"golang.org/x/crypto/ssh/terminal"
)

var loginUsername string
var loginPassword string
var loginPasswordStdin bool

const loginCmdLiteral = "login [environment]"
const loginCmdShortDesc = "Login to an Microgateway Adapter environment"
const loginCmdLongDesc = `Login to an Microgateway Adapter environment using username and password`
const loginCmdExamples = utils.ProjectName + " " + mgCmdLiteral + " login dev -u admin -p admin\n" +
	utils.ProjectName + " " + mgCmdLiteral + " login dev -u admin\n" +
	"cat ~/.mypassword | " + utils.ProjectName + " " + mgCmdLiteral + " login dev -u admin --password-stdin"

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:     loginCmdLiteral,
	Short:   loginCmdShortDesc,
	Long:    loginCmdLongDesc,
	Example: loginCmdExamples,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		environment := args[0]

		store, err := credentials.GetDefaultCredentialStore()
		if err != nil {
			utils.HandleErrorAndExit("Error occurred while loading credential store : ", err)
		}
		err = runLogin(store, environment)
		if err != nil {
			utils.HandleErrorAndExit("Error occurred while login : ", err)
		}
	},
}

func runLogin(store credentials.Store, environment string) error {
	mgwAdapterEndpoints, err := utils.GetEndpointsOfMgwAdapterEnv(environment, utils.MainConfigFilePath)
	if err != nil {
		return errors.New("Env " + environment + " does not exists. Add it using `apictl mg add env`")
	}

	if loginUsername == "" {
		fmt.Print("Username: ")
		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			loginUsername = scanner.Text()
		}
	}

	if loginPassword != "" {
		fmt.Println("Warning: Using --password in CLI is not secure. Use --password-stdin")
		if loginPasswordStdin {
			return errors.New("--password and --password-stdin are mutually exclusive")
		}
	}

	if loginPasswordStdin {
		data, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			return errors.New("Error reading password. Cause: " + err.Error())
		}
		loginPassword = strings.TrimRight(strings.TrimSuffix(string(data), "\n"), "\r")
	}

	if loginPassword == "" {
		fmt.Print("Enter Password: ")
		loginPasswordB, err := terminal.ReadPassword(0)
		loginPassword = string(loginPasswordB)
		fmt.Println()
		if err != nil {
			return errors.New("Error reading password. Cause: " + err.Error())
		}
	}

	accessToken, err := credentials.GetOAuthAccessTokenForMGAdapter(loginUsername, loginPassword,
		mgwAdapterEndpoints.AdapterEndpoint)
	if err != nil {
		utils.HandleErrorAndExit("Error getting access token from adapter", err)
	}

	err = store.SetMGToken(environment, accessToken)
	if err != nil {
		return err
	}
	fmt.Println("Successfully logged into Microgateway Adapter in environment: ", environment)
	return nil
}

// init using Cobra
func init() {
	MgCmd.AddCommand(loginCmd)

	loginCmd.Flags().StringVarP(&loginUsername, "username", "u", "", "Username for login")
	loginCmd.Flags().StringVarP(&loginPassword, "password", "p", "", "Password for login")
	loginCmd.Flags().BoolVarP(&loginPasswordStdin, "password-stdin", "", false, "Get password from stdin")
}
