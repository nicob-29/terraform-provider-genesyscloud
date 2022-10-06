---
page_title: "genesyscloud_outbound_callanalysisresponseset Resource - terraform-provider-genesyscloud"
subcategory: ""
description: |-
  Genesys Cloud outbound Call Analysis Response Set
---
# genesyscloud_outbound_callanalysisresponseset (Resource)

Genesys Cloud outbound Call Analysis Response Set

## API Usage
The following Genesys Cloud APIs are used by this resource. Ensure your OAuth Client has been granted the necessary scopes and permissions to perform these operations:

- [GET /api/v2/outbound/callanalysisresponsesets](https://developer.genesys.cloud/devapps/api-explorer#get-api-v2-outbound-callanalysisresponsesets)
- [POST /api/v2/outbound/callanalysisresponsesets](https://developer.genesys.cloud/devapps/api-explorer#post-api-v2-outbound-callanalysisresponsesets)
- [GET /api/v2/outbound/callanalysisresponsesets/{callAnalysisSetId}](https://developer.genesys.cloud/devapps/api-explorer#get-api-v2-outbound-callanalysisresponsesets--callAnalysisSetId-)
- [PUT /api/v2/outbound/callanalysisresponsesets/{callAnalysisSetId}](https://developer.genesys.cloud/devapps/api-explorer#put-api-v2-outbound-callanalysisresponsesets--callAnalysisSetId-)
- [DELETE /api/v2/outbound/callanalysisresponsesets/{callAnalysisSetId}](https://developer.genesys.cloud/devapps/api-explorer#delete-api-v2-outbound-callanalysisresponsesets--callAnalysisSetId-)

## Example Usage

```terraform
resource "genesyscloud_outbound_callanalysisresponseset" "example-cars" {
  name                   = "Example Call Analysis Response Set"
  beep_detection_enabled = false
  responses {
    callable_person {
      name          = "Example Outbound Flow"
      data          = genesyscloud_flow.flow.id
      reaction_type = "transfer_flow"
    }
    callable_machine {
      reaction_type = "hangup"
    }
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The name of the Response Set.

### Optional

- `beep_detection_enabled` (Boolean) Whether to enable answering machine beep detection Defaults to `false`.
- `responses` (Block List, Max: 1) List of maps of disposition identifiers to reactions. Required if beep_detection_enabled = true. (see [below for nested schema](#nestedblock--responses))

### Read-Only

- `id` (String) The ID of this resource.

<a id="nestedblock--responses"></a>
### Nested Schema for `responses`

Optional:

- `callable_busy` (Block Set, Max: 1) (see [below for nested schema](#nestedblock--responses--callable_busy))
- `callable_disconnect` (Block Set, Max: 1) (see [below for nested schema](#nestedblock--responses--callable_disconnect))
- `callable_fax` (Block Set, Max: 1) (see [below for nested schema](#nestedblock--responses--callable_fax))
- `callable_lineconnected` (Block Set, Max: 1) (see [below for nested schema](#nestedblock--responses--callable_lineconnected))
- `callable_machine` (Block Set, Max: 1) (see [below for nested schema](#nestedblock--responses--callable_machine))
- `callable_noanswer` (Block Set, Max: 1) (see [below for nested schema](#nestedblock--responses--callable_noanswer))
- `callable_person` (Block Set, Max: 1) (see [below for nested schema](#nestedblock--responses--callable_person))
- `callable_sit` (Block Set, Max: 1) (see [below for nested schema](#nestedblock--responses--callable_sit))
- `uncallable_notfound` (Block Set, Max: 1) (see [below for nested schema](#nestedblock--responses--uncallable_notfound))
- `uncallable_sit` (Block Set, Max: 1) (see [below for nested schema](#nestedblock--responses--uncallable_sit))

<a id="nestedblock--responses--callable_busy"></a>
### Nested Schema for `responses.callable_busy`

Required:

- `reaction_type` (String) The reaction to take for a given call analysis result.

Optional:

- `data` (String) Parameter for this reaction. For transfer_flow, this would be the outbound flow id.
- `name` (String) Name of the parameter for this reaction. For transfer_flow, this would be the outbound flow name.


<a id="nestedblock--responses--callable_disconnect"></a>
### Nested Schema for `responses.callable_disconnect`

Required:

- `reaction_type` (String) The reaction to take for a given call analysis result.

Optional:

- `data` (String) Parameter for this reaction. For transfer_flow, this would be the outbound flow id.
- `name` (String) Name of the parameter for this reaction. For transfer_flow, this would be the outbound flow name.


<a id="nestedblock--responses--callable_fax"></a>
### Nested Schema for `responses.callable_fax`

Required:

- `reaction_type` (String) The reaction to take for a given call analysis result.

Optional:

- `data` (String) Parameter for this reaction. For transfer_flow, this would be the outbound flow id.
- `name` (String) Name of the parameter for this reaction. For transfer_flow, this would be the outbound flow name.


<a id="nestedblock--responses--callable_lineconnected"></a>
### Nested Schema for `responses.callable_lineconnected`

Required:

- `reaction_type` (String) The reaction to take for a given call analysis result.

Optional:

- `data` (String) Parameter for this reaction. For transfer_flow, this would be the outbound flow id.
- `name` (String) Name of the parameter for this reaction. For transfer_flow, this would be the outbound flow name.


<a id="nestedblock--responses--callable_machine"></a>
### Nested Schema for `responses.callable_machine`

Required:

- `reaction_type` (String) The reaction to take for a given call analysis result.

Optional:

- `data` (String) Parameter for this reaction. For transfer_flow, this would be the outbound flow id.
- `name` (String) Name of the parameter for this reaction. For transfer_flow, this would be the outbound flow name.


<a id="nestedblock--responses--callable_noanswer"></a>
### Nested Schema for `responses.callable_noanswer`

Required:

- `reaction_type` (String) The reaction to take for a given call analysis result.

Optional:

- `data` (String) Parameter for this reaction. For transfer_flow, this would be the outbound flow id.
- `name` (String) Name of the parameter for this reaction. For transfer_flow, this would be the outbound flow name.


<a id="nestedblock--responses--callable_person"></a>
### Nested Schema for `responses.callable_person`

Required:

- `reaction_type` (String) The reaction to take for a given call analysis result.

Optional:

- `data` (String) Parameter for this reaction. For transfer_flow, this would be the outbound flow id.
- `name` (String) Name of the parameter for this reaction. For transfer_flow, this would be the outbound flow name.


<a id="nestedblock--responses--callable_sit"></a>
### Nested Schema for `responses.callable_sit`

Required:

- `reaction_type` (String) The reaction to take for a given call analysis result.

Optional:

- `data` (String) Parameter for this reaction. For transfer_flow, this would be the outbound flow id.
- `name` (String) Name of the parameter for this reaction. For transfer_flow, this would be the outbound flow name.


<a id="nestedblock--responses--uncallable_notfound"></a>
### Nested Schema for `responses.uncallable_notfound`

Required:

- `reaction_type` (String) The reaction to take for a given call analysis result.

Optional:

- `data` (String) Parameter for this reaction. For transfer_flow, this would be the outbound flow id.
- `name` (String) Name of the parameter for this reaction. For transfer_flow, this would be the outbound flow name.


<a id="nestedblock--responses--uncallable_sit"></a>
### Nested Schema for `responses.uncallable_sit`

Required:

- `reaction_type` (String) The reaction to take for a given call analysis result.

Optional:

- `data` (String) Parameter for this reaction. For transfer_flow, this would be the outbound flow id.
- `name` (String) Name of the parameter for this reaction. For transfer_flow, this would be the outbound flow name.
