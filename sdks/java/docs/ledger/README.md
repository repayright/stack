# ledger

### Available Operations

* [addMetadataOnTransaction](#addmetadataontransaction) - Set the metadata of a transaction by its ID
* [addMetadataToAccount](#addmetadatatoaccount) - Add metadata to an account
* [countAccounts](#countaccounts) - Count the accounts from a ledger
* [countTransactions](#counttransactions) - Count the transactions from a ledger
* [createTransaction](#createtransaction) - Create a new transaction to a ledger
* [getAccount](#getaccount) - Get account by its address
* [getBalancesAggregated](#getbalancesaggregated) - Get the aggregated balances from selected accounts
* [getInfo](#getinfo) - Show server information
* [getLedgerInfo](#getledgerinfo) - Get information about a ledger
* [getTransaction](#gettransaction) - Get transaction from a ledger by its ID
* [listAccounts](#listaccounts) - List accounts from a ledger
* [listLogs](#listlogs) - List the logs from a ledger
* [listTransactions](#listtransactions) - List transactions from a ledger
* [readStats](#readstats) - Get statistics from a ledger
* [revertTransaction](#reverttransaction) - Revert a ledger transaction by its ID

## addMetadataOnTransaction

Set the metadata of a transaction by its ID

### Example Usage

```java
package hello.world;

import com.formance.formance_sdk.SDK;
import com.formance.formance_sdk.models.operations.AddMetadataOnTransactionRequest;
import com.formance.formance_sdk.models.operations.AddMetadataOnTransactionResponse;
import com.formance.formance_sdk.models.shared.Security;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("nobis") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
                }})
                .build();

            AddMetadataOnTransactionRequest req = new AddMetadataOnTransactionRequest(1234L, "ledger001") {{
                idempotencyKey = "enim";
                requestBody = new java.util.HashMap<String, String>() {{
                    put("nemo", "minima");
                    put("excepturi", "accusantium");
                    put("iure", "culpa");
                }};
                dryRun = true;
            }};            

            AddMetadataOnTransactionResponse res = sdk.ledger.addMetadataOnTransaction(req);

            if (res.statusCode == 200) {
                // handle response
            }
        } catch (Exception e) {
            // handle exception
        }
    }
}
```

## addMetadataToAccount

Add metadata to an account

### Example Usage

```java
package hello.world;

import com.formance.formance_sdk.SDK;
import com.formance.formance_sdk.models.operations.AddMetadataToAccountRequest;
import com.formance.formance_sdk.models.operations.AddMetadataToAccountResponse;
import com.formance.formance_sdk.models.shared.Security;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("doloribus") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
                }})
                .build();

            AddMetadataToAccountRequest req = new AddMetadataToAccountRequest(                new java.util.HashMap<String, String>() {{
                                put("architecto", "mollitia");
                                put("dolorem", "culpa");
                                put("consequuntur", "repellat");
                                put("mollitia", "occaecati");
                            }}, "users:001", "ledger001") {{
                idempotencyKey = "numquam";
                dryRun = true;
            }};            

            AddMetadataToAccountResponse res = sdk.ledger.addMetadataToAccount(req);

            if (res.statusCode == 200) {
                // handle response
            }
        } catch (Exception e) {
            // handle exception
        }
    }
}
```

## countAccounts

Count the accounts from a ledger

### Example Usage

```java
package hello.world;

import com.formance.formance_sdk.SDK;
import com.formance.formance_sdk.models.operations.CountAccountsRequest;
import com.formance.formance_sdk.models.operations.CountAccountsResponse;
import com.formance.formance_sdk.models.shared.Security;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("commodi") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
                }})
                .build();

            CountAccountsRequest req = new CountAccountsRequest("ledger001") {{
                address = "users:.+";
                metadata = new java.util.HashMap<String, Object>() {{
                    put("molestiae", "velit");
                    put("error", "quia");
                }};
            }};            

            CountAccountsResponse res = sdk.ledger.countAccounts(req);

            if (res.statusCode == 200) {
                // handle response
            }
        } catch (Exception e) {
            // handle exception
        }
    }
}
```

## countTransactions

Count the transactions from a ledger

### Example Usage

```java
package hello.world;

import com.formance.formance_sdk.SDK;
import com.formance.formance_sdk.models.operations.CountTransactionsRequest;
import com.formance.formance_sdk.models.operations.CountTransactionsResponse;
import com.formance.formance_sdk.models.shared.Security;
import java.time.OffsetDateTime;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("quis") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
                }})
                .build();

            CountTransactionsRequest req = new CountTransactionsRequest("ledger001") {{
                account = "users:001";
                destination = "users:001";
                endTime = OffsetDateTime.parse("2022-04-29T17:10:10.440Z");
                metadata = new java.util.HashMap<String, String>() {{
                    put("enim", "odit");
                    put("quo", "sequi");
                    put("tenetur", "ipsam");
                }};
                reference = "ref:001";
                source = "users:001";
                startTime = OffsetDateTime.parse("2021-05-11T16:11:54.761Z");
            }};            

            CountTransactionsResponse res = sdk.ledger.countTransactions(req);

            if (res.statusCode == 200) {
                // handle response
            }
        } catch (Exception e) {
            // handle exception
        }
    }
}
```

## createTransaction

Create a new transaction to a ledger

### Example Usage

```java
package hello.world;

import com.formance.formance_sdk.SDK;
import com.formance.formance_sdk.models.operations.CreateTransactionRequest;
import com.formance.formance_sdk.models.operations.CreateTransactionResponse;
import com.formance.formance_sdk.models.shared.PostTransaction;
import com.formance.formance_sdk.models.shared.PostTransactionScript;
import com.formance.formance_sdk.models.shared.Posting;
import com.formance.formance_sdk.models.shared.Security;
import java.time.OffsetDateTime;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("aut") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
                }})
                .build();

            CreateTransactionRequest req = new CreateTransactionRequest(                new PostTransaction(                new java.util.HashMap<String, String>() {{
                                                put("error", "temporibus");
                                            }}) {{
                                postings = new com.formance.formance_sdk.models.shared.Posting[]{{
                                    add(new Posting(100L, "COIN", "users:002", "users:001") {{
                                        amount = 100L;
                                        asset = "COIN";
                                        destination = "users:002";
                                        source = "users:001";
                                    }}),
                                    add(new Posting(100L, "COIN", "users:002", "users:001") {{
                                        amount = 100L;
                                        asset = "COIN";
                                        destination = "users:002";
                                        source = "users:001";
                                    }}),
                                    add(new Posting(100L, "COIN", "users:002", "users:001") {{
                                        amount = 100L;
                                        asset = "COIN";
                                        destination = "users:002";
                                        source = "users:001";
                                    }}),
                                }};
                                reference = "ref:001";
                                script = new PostTransactionScript("vars {
                                account $user
                                }
                                send [COIN 10] (
                                	source = @world
                                	destination = $user
                                )
                                ") {{
                                    vars = new java.util.HashMap<String, Object>() {{
                                        put("reiciendis", "voluptatibus");
                                    }};
                                }};;
                                timestamp = OffsetDateTime.parse("2021-08-05T19:50:46.898Z");
                            }};, "ledger001") {{
                idempotencyKey = "praesentium";
                dryRun = true;
            }};            

            CreateTransactionResponse res = sdk.ledger.createTransaction(req);

            if (res.createTransactionResponse != null) {
                // handle response
            }
        } catch (Exception e) {
            // handle exception
        }
    }
}
```

## getAccount

Get account by its address

### Example Usage

```java
package hello.world;

import com.formance.formance_sdk.SDK;
import com.formance.formance_sdk.models.operations.GetAccountRequest;
import com.formance.formance_sdk.models.operations.GetAccountResponse;
import com.formance.formance_sdk.models.shared.Security;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("voluptatibus") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
                }})
                .build();

            GetAccountRequest req = new GetAccountRequest("users:001", "ledger001") {{
                expand = "ipsa";
            }};            

            GetAccountResponse res = sdk.ledger.getAccount(req);

            if (res.accountResponse != null) {
                // handle response
            }
        } catch (Exception e) {
            // handle exception
        }
    }
}
```

## getBalancesAggregated

Get the aggregated balances from selected accounts

### Example Usage

```java
package hello.world;

import com.formance.formance_sdk.SDK;
import com.formance.formance_sdk.models.operations.GetBalancesAggregatedRequest;
import com.formance.formance_sdk.models.operations.GetBalancesAggregatedResponse;
import com.formance.formance_sdk.models.shared.Security;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("omnis") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
                }})
                .build();

            GetBalancesAggregatedRequest req = new GetBalancesAggregatedRequest("ledger001") {{
                address = "users:001";
            }};            

            GetBalancesAggregatedResponse res = sdk.ledger.getBalancesAggregated(req);

            if (res.aggregateBalancesResponse != null) {
                // handle response
            }
        } catch (Exception e) {
            // handle exception
        }
    }
}
```

## getInfo

Show server information

### Example Usage

```java
package hello.world;

import com.formance.formance_sdk.SDK;
import com.formance.formance_sdk.models.operations.GetInfoResponse;
import com.formance.formance_sdk.models.shared.Security;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("voluptate") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
                }})
                .build();

            GetInfoResponse res = sdk.ledger.getInfo();

            if (res.configInfoResponse != null) {
                // handle response
            }
        } catch (Exception e) {
            // handle exception
        }
    }
}
```

## getLedgerInfo

Get information about a ledger

### Example Usage

```java
package hello.world;

import com.formance.formance_sdk.SDK;
import com.formance.formance_sdk.models.operations.GetLedgerInfoRequest;
import com.formance.formance_sdk.models.operations.GetLedgerInfoResponse;
import com.formance.formance_sdk.models.shared.Security;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("cum") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
                }})
                .build();

            GetLedgerInfoRequest req = new GetLedgerInfoRequest("ledger001");            

            GetLedgerInfoResponse res = sdk.ledger.getLedgerInfo(req);

            if (res.ledgerInfoResponse != null) {
                // handle response
            }
        } catch (Exception e) {
            // handle exception
        }
    }
}
```

## getTransaction

Get transaction from a ledger by its ID

### Example Usage

```java
package hello.world;

import com.formance.formance_sdk.SDK;
import com.formance.formance_sdk.models.operations.GetTransactionRequest;
import com.formance.formance_sdk.models.operations.GetTransactionResponse;
import com.formance.formance_sdk.models.shared.Security;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("perferendis") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
                }})
                .build();

            GetTransactionRequest req = new GetTransactionRequest(1234L, "ledger001") {{
                expand = "doloremque";
            }};            

            GetTransactionResponse res = sdk.ledger.getTransaction(req);

            if (res.getTransactionResponse != null) {
                // handle response
            }
        } catch (Exception e) {
            // handle exception
        }
    }
}
```

## listAccounts

List accounts from a ledger, sorted by address in descending order.

### Example Usage

```java
package hello.world;

import com.formance.formance_sdk.SDK;
import com.formance.formance_sdk.models.operations.ListAccountsRequest;
import com.formance.formance_sdk.models.operations.ListAccountsResponse;
import com.formance.formance_sdk.models.shared.Security;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("reprehenderit") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
                }})
                .build();

            ListAccountsRequest req = new ListAccountsRequest("ledger001") {{
                address = "users:.+";
                cursor = "aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ==";
                expand = "ut";
                metadata = new java.util.HashMap<String, String>() {{
                    put("dicta", "corporis");
                    put("dolore", "iusto");
                    put("dicta", "harum");
                    put("enim", "accusamus");
                }};
                pageSize = 414263L;
            }};            

            ListAccountsResponse res = sdk.ledger.listAccounts(req);

            if (res.accountsCursorResponse != null) {
                // handle response
            }
        } catch (Exception e) {
            // handle exception
        }
    }
}
```

## listLogs

List the logs from a ledger, sorted by ID in descending order.

### Example Usage

```java
package hello.world;

import com.formance.formance_sdk.SDK;
import com.formance.formance_sdk.models.operations.ListLogsRequest;
import com.formance.formance_sdk.models.operations.ListLogsResponse;
import com.formance.formance_sdk.models.shared.Security;
import java.time.OffsetDateTime;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("repudiandae") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
                }})
                .build();

            ListLogsRequest req = new ListLogsRequest("ledger001") {{
                cursor = "aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ==";
                endTime = OffsetDateTime.parse("2022-10-13T20:38:16.426Z");
                pageSize = 692472L;
                startTime = OffsetDateTime.parse("2021-11-13T09:08:33.009Z");
            }};            

            ListLogsResponse res = sdk.ledger.listLogs(req);

            if (res.logsCursorResponse != null) {
                // handle response
            }
        } catch (Exception e) {
            // handle exception
        }
    }
}
```

## listTransactions

List transactions from a ledger, sorted by id in descending order.

### Example Usage

```java
package hello.world;

import com.formance.formance_sdk.SDK;
import com.formance.formance_sdk.models.operations.ListTransactionsRequest;
import com.formance.formance_sdk.models.operations.ListTransactionsResponse;
import com.formance.formance_sdk.models.shared.Security;
import java.time.OffsetDateTime;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("pariatur") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
                }})
                .build();

            ListTransactionsRequest req = new ListTransactionsRequest("ledger001") {{
                account = "users:001";
                cursor = "aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ==";
                destination = "users:001";
                endTime = OffsetDateTime.parse("2022-06-29T05:25:54.356Z");
                expand = "rem";
                metadata = new java.util.HashMap<String, String>() {{
                    put("quasi", "repudiandae");
                    put("sint", "veritatis");
                    put("itaque", "incidunt");
                    put("enim", "consequatur");
                }};
                pageSize = 667411L;
                reference = "ref:001";
                source = "users:001";
                startTime = OffsetDateTime.parse("2022-08-09T16:21:07.003Z");
            }};            

            ListTransactionsResponse res = sdk.ledger.listTransactions(req);

            if (res.transactionsCursorResponse != null) {
                // handle response
            }
        } catch (Exception e) {
            // handle exception
        }
    }
}
```

## readStats

Get statistics from a ledger. (aggregate metrics on accounts and transactions)


### Example Usage

```java
package hello.world;

import com.formance.formance_sdk.SDK;
import com.formance.formance_sdk.models.operations.ReadStatsRequest;
import com.formance.formance_sdk.models.operations.ReadStatsResponse;
import com.formance.formance_sdk.models.shared.Security;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("deserunt") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
                }})
                .build();

            ReadStatsRequest req = new ReadStatsRequest("ledger001");            

            ReadStatsResponse res = sdk.ledger.readStats(req);

            if (res.statsResponse != null) {
                // handle response
            }
        } catch (Exception e) {
            // handle exception
        }
    }
}
```

## revertTransaction

Revert a ledger transaction by its ID

### Example Usage

```java
package hello.world;

import com.formance.formance_sdk.SDK;
import com.formance.formance_sdk.models.operations.RevertTransactionRequest;
import com.formance.formance_sdk.models.operations.RevertTransactionResponse;
import com.formance.formance_sdk.models.shared.Security;

public class Application {
    public static void main(String[] args) {
        try {
            SDK sdk = SDK.builder()
                .setSecurity(new Security("distinctio") {{
                    authorization = "Bearer YOUR_ACCESS_TOKEN_HERE";
                }})
                .build();

            RevertTransactionRequest req = new RevertTransactionRequest(1234L, "ledger001");            

            RevertTransactionResponse res = sdk.ledger.revertTransaction(req);

            if (res.revertTransactionResponse != null) {
                // handle response
            }
        } catch (Exception e) {
            // handle exception
        }
    }
}
```
