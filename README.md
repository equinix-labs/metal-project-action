# packet-project-action

Experimental Github Action for creating [Packet](https://packet.com) Projects.

Give a Packet User API Token, a new Project will be created, preconfigured with an SSH Key and API Key which can be used in subsequent actions.

Clean up the project with the [Packet Sweeper Action](https://github.com/displague/packet-sweeper-action).

See the [Packet Actions Example](https://github.com/displague/packet-actions-example) for usage examples.

## Input

With | Environment variable | Description
--- | --- | ---
`userToken` | `PACKET_AUTH_TOKEN` | (required) A Packet User API Token
`projectName` | - | Name for the project, API key, and SSH Key. A generated name will be used if not supplied.
`organizationID` | - | Organization ID that the Project will be created under. If not supplied, the default organization for the API User will be used.

## Output

Output Name | Environment Variable | Description
--- | --- | ---
`projectID` | `PACKET_PROJECT_ID` | a new Packet Project ID
`projectName` | `PACKET_PROJECT_NAME` | the generated (or supplied) name of the Packet Project
`projectToken` | `PACKET_PROJECT_TOKEN` | a new Packet Project API Token bound to this project
`projectSSHPrivateKey` | `PACKET_SSH_PRIVATE_KEY`  | a new SSH Private Key that can be used to authenticate to devices in this project
