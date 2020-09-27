# packet-project-action

Experimental Github Action for creating [Packet](https://packet.com) Projects.

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
