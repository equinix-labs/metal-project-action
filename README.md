# metal-project-action

Experimental Github Action for creating [Equinix Metal](https://metal.equinix.com) Projects.

> :bulb: See also:
>
> - [equinix-metal-sweeper](https://github.com/equinix-labs/metal-sweeper-action) action
> - [equinix-metal-examples](https://github.com/equinix-labs/metal-actions-example) examples

Given a Equinix Metal User API Token, a new Project will be created, preconfigured with an SSH Key and API Key which can be used in subsequent actions.

Clean up the project with the [Equinix Metal Sweeper Action](https://github.com/equinix-labs/metal-sweeper-action).

See the [Equinix Metal Actions Example](https://github.com/equinix-labs/metal-actions-example) for usage examples.

## Input

| With             | Description                                                                                                                      |
| ---------------- | -------------------------------------------------------------------------------------------------------------------------------- |
| `userToken`      | (required) A Equinix Metal User API Token                                                                                        |
| `projectName`    | Name for the project, API key, and SSH Key. A generated name will be used if not supplied.                                       |
| `organizationID` | Organization ID that the Project will be created under. If not supplied, the default organization for the API User will be used. |

## Output

| Output Name                  | Description                                                                                        |
| ---------------------------- | -------------------------------------------------------------------------------------------------- |
| `projectID`                  | a new Equinix Metal Project ID                                                                     |
| `projectName`                | the generated (or supplied) name of the Equinix Metal Project                                      |
| `projectToken`               | a new Equinix Metal Project API Token bound to this project                                        |
| `projectSSHPrivateKeyBase64` | a new SSH Private Key (base64 encoded) that can be used to authenticate to devices in this project |
| `projectSSHPublicKey`        | a new SSH Public Key that can be used to authenticate to devices in this project                   |
| `organizationID`             | ID of the Organization responsible for the project.                                                |
