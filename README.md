# cloud-sql-proxy
Sample application with connecting with Cloud SQL.

## Create a DB with Cloud SQL
Follow this tutorial for creating a MySQL instance.

[Tutorial Link](https://cloud.google.com/sql/docs/mysql/quickstart)

## Create a Service Account
Follow the linked tutorial for creating a Service Account

[Service Account Tutorial](https://cloud.google.com/sql/docs/mysql/connect-external-app#4_if_required_by_your_authentication_method_create_a_service_account)
* Note: I only got the service account to connect to the DB by giving the account __Project > Editor__ role

## Reference
* https://medium.com/@DazWilkin/google-cloud-sql-6-ways-golang-a4aa497f3c67