Set ENV variables:

            "TADO_USERNAME": "tadouser@example.com",
            "TADO_PASSWORD": "tadopassword",
            "TADO_HOME": "My Home",
            "BQ_PROJECT_ID": "bq-project-name",
            "BQ_DATASET": "bq_dataset_name",
            "BQ_TABLE": "bq_table_name"

Make BQ Table:

bq mk --schema tadodailyreport/tado_heating.schema -t $BQ_PROJECT_ID:$BQ_DATASET.$BQ_TABLE

Regenerate proto:

git clone https://github.com/googleapis/api-common-protos.git
cp bq-table.proto bq-field.proto /usr/local/include/

protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative tadodailyreport/tadoDailyReport.proto

protoc --bq-schema_out=. tadodailyreport/tadoDailyReport.proto