# proluxe-cxp-api

API for Proluxe MFG 

## BigQuery

For local development you must authenticate with Google Cloud. Run the following command and follow the prompts:


```
yay -S google-cloud-cli
source ~/.zshrc

gcloud init
gcloud auth application-default login
cp /home/scott/.config/gcloud/application_default_credentials.json .

``` 

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