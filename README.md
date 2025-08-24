# proluxe-cxp-api

API for Proluxe MFG 

## Authentication & User Context

### `GetCurrentUser` Helper

The `GetCurrentUser` function is used throughout the API to extract the current user's email address from the request context.

**Behavior (as of August 2025):**
- The function now reads the user's email from the `X-User-Email` HTTP header.
- It returns a single string: the user's email address.
- If the header is missing, it returns an empty string and logs a warning.

**Usage Example:**
```go
email := GetCurrentUser(c)
```

**Breaking Change:**  
Previous versions of this function returned three values (name, email, avatar) extracted from a JWT. All usages must now expect a single return value (email). Update any code that previously used name or avatar accordingly.


## BigQuery

For local development you must authenticate with Google Cloud. Run the following command and follow the prompts:


```
yay -S google-cloud-cli
source ~/.zshrc

gcloud init
gcloud auth application-default login
cp /home/scott/.config/gcloud/application_default_credentials.json .

``` 

## AlgoliaIndex Fetching and Flexible SOQL Mapping

The `model/algolia.go` helpers provide a flexible way to fetch and map Salesforce records for Algolia indexing.

### Usage

You can now specify both the field to use as the Algolia `objectID` and the field to use as the display `Name` when calling `FetchRecords`:

```go
// Use default fields ("Id" for objectID, "Name" for Name)
records := FetchRecords(client, soqlQuery)

// Use custom fields (e.g., "Model_Number__c" as objectID, "Name" as Name)
records := FetchRecords(client, soqlQuery, "Model_Number__c", "Name")
```

All `FetchAlgolia*` helpers now use this pattern for clarity and flexibility.

### Example

```go
// Fetch catalog models using Model_Number__c as the Algolia objectID
catalog := FetchAlgoliaCatalog(client)
// catalog[0].ObjectId == Model_Number__c
// catalog[0].Name == Name
```

**Note:**  
If you pass only one field, it is used as the objectID and "Name" is used for the name. If you pass two fields, the first is the objectID, the second is the name.

### Testing

Unit tests are not yet implemented for these helpers.  
**TODO:** Add tests for:
- Mapping of idField and nameField
- Handling of missing/invalid fields
- SOQL query error handling
- End-to-end mapping from Salesforce SObject to AlgoliaIndex

```
WITH RECURSIVE BomHierarchy AS (
    -- Anchor part: Initial selection from bom_costs, incorporating QtyOrdered
    SELECT
        bc.Name,
        bc.ExternalId,
        bc.SrcPl,
        bc.CompItemC,
        bc.ItemC,
        bc.Labor AS LaborUnitCost,
        bc.Material AS MaterialUnitCost,
        bc.QOH,
        bc.UOM,
        bc.QtyPer,
        bc.POQtyRequired,
        bc.POQtyReceived,
        bc.POQtyOutstanding,
        bc.WOQtyRequired,
        bc.WOQtyAccepted,
        bc.WOQtyWIP,
        bc.WOQtyShipped,
        bc.WOQtyScrapped,
        o.QtyOrdered,
        (bc.QtyPer * o.QtyOrdered) AS QtyOnOrder,
        0 AS Level
    FROM
        `proluxe-portal.mrp_nightly.bom_costs` bc
    JOIN
        `proluxe-portal.mrp_nightly.orders` o ON bc.ParentExternalId = o.ExternalId
    WHERE
        o.DueDate <= DATE_ADD(CURRENT_DATE(), INTERVAL 4 MONTH)

    UNION ALL

    -- Recursive part: Propagate QtyOrdered through the BOM hierarchy
    SELECT
        bc.Name,
        bc.ExternalId,
        bc.SrcPl,
        bc.CompItemC,
        bc.ItemC,
        bc.Labor AS LaborUnitCost,
        bc.Material AS MaterialUnitCost,
        bc.QOH,
        bc.UOM,
        bc.QtyPer,
        bc.POQtyRequired,
        bc.POQtyReceived,
        bc.POQtyOutstanding,
        bc.WOQtyRequired,
        bc.WOQtyAccepted,
        bc.WOQtyWIP,
        bc.WOQtyShipped,
        bc.WOQtyScrapped,
        h.QtyOrdered, -- Maintaining QtyOrdered from the hierarchy
        (bc.QtyPer * h.QtyOrdered) AS QtyOnOrder,
        h.Level + 1
    FROM
        `proluxe-portal.mrp_nightly.bom_costs` bc
    JOIN
        BomHierarchy h ON bc.ItemC = h.CompItemC
)
SELECT 
    ExternalId, 
    ANY_VALUE(Name) as Name, 
    ANY_VALUE(SrcPl) as SrcPl, 
    ANY_VALUE(UOM) as UOM, 
    ANY_VALUE(QOH) as QOH, 
    ANY_VALUE(MaterialUnitCost) as MaterialUnitCost, 
    ANY_VALUE(LaborUnitCost) as LaborUnitCost, 
    ANY_VALUE(POQtyRequired) as POQtyRequired, 
    ANY_VALUE(POQtyReceived) as POQtyReceived, 
    ANY_VALUE(POQtyOutstanding) as POQtyOutstanding, 
    ANY_VALUE(WOQtyRequired) as WOQtyRequired, 
    ANY_VALUE(WOQtyAccepted) as WOQtyAccepted, 
    ANY_VALUE(WOQtyWIP) as WOQtyWIP, 
    ANY_VALUE(WOQtyShipped) as WOQtyShipped, 
    ANY_VALUE(WOQtyScrapped) as WOQtyScrapped, 
    SUM(QtyOnOrder) as QtyOnOrder,
    SUM(MaterialUnitCost) as MaterialCost,
    SUM(LaborUnitCost) as LaborCost
FROM BomHierarchy
GROUP BY ExternalId
HAVING ExternalId = "Main_110115131C1400"
```
