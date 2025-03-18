package google

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/bigquery"
	u "github.com/scottraio/go-utils"
	"google.golang.org/api/iterator"
)

func BigQuery() (*bigquery.Client, context.Context) {
	ctx := context.Background()
	client, err := bigquery.NewClient(ctx, u.GetDotEnvVariable("GCP_PROJECT_ID"))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()
	return client, ctx
}

func BigQueryInsert(ctx context.Context, client *bigquery.Client, data []interface{}, table string) error {
	inserter := client.Dataset(u.GetDotEnvVariable("BIGQUERY_DATASET_ID")).Table(table).Inserter()
	err := inserter.Put(ctx, data)

	// syncLog := []SyncLog{
	// 	{
	// 		CreatedAt:   util.PSTDateToTime(time.Now()),
	// 		ProcessedAt: util.PSTDateToCivil(time.Now()),
	// 		Success:     err == nil,
	// 		Message:     SyncLogMessage(err),
	// 		Table:       table,
	// 	},
	// }

	// syncLogErr := LogSyncLog(ctx, client, syncLog)

	// if syncLogErr != nil {
	// 	return syncLogErr
	// }

	if err != nil {
		return err
	}

	return nil
}

func BigQueryDelete(ctx context.Context, client *bigquery.Client, table string) (*bigquery.Job, error) {
	query := fmt.Sprintf(`DELETE FROM %s WHERE Id IS NOT NULL`, table)

	job, err := client.Query(query).Run(ctx)
	if err != nil {
		// handle error
		fmt.Println("Error running data:", err)
	}

	_, err = job.Wait(ctx)
	// if status.Err() != nil {
	// 	fmt.Println(status.Err())
	// }
	// if err != nil || status.Err() != nil {
	// 	syncLog := []SyncLog{
	// 		{
	// 			CreatedAt:   util.PSTDateToTime(time.Now()),
	// 			ProcessedAt: util.PSTDateToCivil(time.Now()),
	// 			Success:     err == nil,
	// 			Message:     SyncLogMessage(err),
	// 			Table:       table,
	// 		},
	// 	}

	// 	syncLogErr := LogSyncLog(ctx, client, syncLog)

	// 	if syncLogErr != nil {
	// 		return job, syncLogErr
	// 	}
	// }

	return job, err
}

func BigQueryFind(ctx context.Context, client *bigquery.Client, query string) (*bigquery.RowIterator, error) {
	it, err := client.Query(query).Read(ctx)
	if err != nil {
		// handle error
		fmt.Println("Error querying data:", err)
	}
	if it == nil {
		err = fmt.Errorf("query did not initialize the iterator")
		fmt.Println("Query did not initialize the iterator")
		return nil, err // or handle accordingly
	}

	return it, nil
}

// BigQueryIsTableEmpty checks if a given BigQuery table is empty.
func BigQueryIsTableEmpty(ctx context.Context, client *bigquery.Client, tableName string) (bool, error) {
	query := client.Query(fmt.Sprintf("SELECT COUNT(*) as count FROM `%s`", tableName))
	it, err := query.Read(ctx)
	if err != nil {
		return false, fmt.Errorf("query execution failed: %w", err)
	}

	var data struct {
		Count int64
	}
	err = it.Next(&data)
	if err == iterator.Done {
		return true, nil // If no rows are returned, the table is empty
	}
	if err != nil {
		return false, fmt.Errorf("error reading query result: %w", err)
	}

	return data.Count == 0, nil
}

func FullTableName(table string) string {
	return u.GetDotEnvVariable("BIGQUERY_DATASET_ID") + "." + table
}

func HandleClientError(client *bigquery.Client, ctx context.Context) {
	if client == nil || ctx == nil {
		fmt.Println("Failed to initialize BigQuery client or context")
	}
}
