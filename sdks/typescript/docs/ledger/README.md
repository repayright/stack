# ledger

### Available Operations

* [addMetadataOnTransaction](#addmetadataontransaction) - Set the metadata of a transaction by its ID
* [addMetadataToAccount](#addmetadatatoaccount) - Add metadata to an account
* [countAccounts](#countaccounts) - Count the accounts from a ledger
* [countTransactions](#counttransactions) - Count the transactions from a ledger
* [createTransaction](#createtransaction) - Create a new transaction to a ledger
* [getAccount](#getaccount) - Get account by its address
* [getBalances](#getbalances) - Get the balances from a ledger's account
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

```typescript
import { SDK } from "@formance/formance-sdk";
import { AddMetadataOnTransactionResponse } from "@formance/formance-sdk/dist/sdk/models/operations";
import { ErrorsEnum } from "@formance/formance-sdk/dist/sdk/models/shared";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.ledger.addMetadataOnTransaction({
  idempotencyKey: "ipsam",
  requestBody: {
    "possimus": "aut",
    "quasi": "error",
    "temporibus": "laborum",
  },
  async: true,
  dryRun: true,
  ledger: "ledger001",
  txid: 1234,
}).then((res: AddMetadataOnTransactionResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## addMetadataToAccount

Add metadata to an account

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { AddMetadataToAccountResponse } from "@formance/formance-sdk/dist/sdk/models/operations";
import { ErrorsEnum } from "@formance/formance-sdk/dist/sdk/models/shared";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.ledger.addMetadataToAccount({
  idempotencyKey: "quasi",
  requestBody: {
    "voluptatibus": "vero",
    "nihil": "praesentium",
    "voluptatibus": "ipsa",
    "omnis": "voluptate",
  },
  address: "users:001",
  async: true,
  dryRun: true,
  ledger: "ledger001",
}).then((res: AddMetadataToAccountResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## countAccounts

Count the accounts from a ledger

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { CountAccountsResponse } from "@formance/formance-sdk/dist/sdk/models/operations";
import { ErrorsEnum } from "@formance/formance-sdk/dist/sdk/models/shared";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.ledger.countAccounts({
  address: "users:.+",
  ledger: "ledger001",
  metadata: {
    "perferendis": "doloremque",
    "reprehenderit": "ut",
    "maiores": "dicta",
  },
}).then((res: CountAccountsResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## countTransactions

Count the transactions from a ledger

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { CountTransactionsResponse } from "@formance/formance-sdk/dist/sdk/models/operations";
import { ErrorsEnum } from "@formance/formance-sdk/dist/sdk/models/shared";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.ledger.countTransactions({
  account: "users:001",
  destination: "users:001",
  endTime: new Date("2022-09-14T21:48:41.971Z"),
  ledger: "ledger001",
  metadata: {
    "dicta": "harum",
    "enim": "accusamus",
  },
  reference: "ref:001",
  source: "users:001",
  startTime: new Date("2022-01-30T20:15:26.045Z"),
}).then((res: CountTransactionsResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## createTransaction

Create a new transaction to a ledger

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { CreateTransactionResponse } from "@formance/formance-sdk/dist/sdk/models/operations";
import { ErrorsEnum } from "@formance/formance-sdk/dist/sdk/models/shared";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.ledger.createTransaction({
  idempotencyKey: "quae",
  postTransaction: {
    metadata: {
      "quidem": "molestias",
    },
    postings: [
      {
        amount: 100,
        asset: "COIN",
        destination: "users:002",
        source: "users:001",
      },
      {
        amount: 100,
        asset: "COIN",
        destination: "users:002",
        source: "users:001",
      },
      {
        amount: 100,
        asset: "COIN",
        destination: "users:002",
        source: "users:001",
      },
    ],
    reference: "ref:001",
    script: {
      plain: "vars {
    account $user
    }
    send [COIN 10] (
    	source = @world
    	destination = $user
    )
    ",
      vars: {
        "modi": "praesentium",
        "rem": "voluptates",
        "quasi": "repudiandae",
        "sint": "veritatis",
      },
    },
    timestamp: new Date("2022-03-02T21:33:21.372Z"),
  },
  async: true,
  dryRun: true,
  ledger: "ledger001",
}).then((res: CreateTransactionResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## getAccount

Get account by its address

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { GetAccountResponse } from "@formance/formance-sdk/dist/sdk/models/operations";
import { ErrorsEnum } from "@formance/formance-sdk/dist/sdk/models/shared";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.ledger.getAccount({
  address: "users:001",
  ledger: "ledger001",
}).then((res: GetAccountResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## getBalances

Get the balances from a ledger's account

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { GetBalancesResponse } from "@formance/formance-sdk/dist/sdk/models/operations";
import { ErrorsEnum } from "@formance/formance-sdk/dist/sdk/models/shared";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.ledger.getBalances({
  address: "users:001",
  cursor: "aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ==",
  ledger: "ledger001",
  pageSize: 318569,
}).then((res: GetBalancesResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## getBalancesAggregated

Get the aggregated balances from selected accounts

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { GetBalancesAggregatedResponse } from "@formance/formance-sdk/dist/sdk/models/operations";
import { ErrorsEnum } from "@formance/formance-sdk/dist/sdk/models/shared";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.ledger.getBalancesAggregated({
  address: "users:001",
  ledger: "ledger001",
}).then((res: GetBalancesAggregatedResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## getInfo

Show server information

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { GetInfoResponse } from "@formance/formance-sdk/dist/sdk/models/operations";
import { ErrorsEnum } from "@formance/formance-sdk/dist/sdk/models/shared";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.ledger.getInfo().then((res: GetInfoResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## getLedgerInfo

Get information about a ledger

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { GetLedgerInfoResponse } from "@formance/formance-sdk/dist/sdk/models/operations";
import { ErrorsEnum, MigrationInfoState } from "@formance/formance-sdk/dist/sdk/models/shared";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.ledger.getLedgerInfo({
  ledger: "ledger001",
}).then((res: GetLedgerInfoResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## getTransaction

Get transaction from a ledger by its ID

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { GetTransactionResponse } from "@formance/formance-sdk/dist/sdk/models/operations";
import { ErrorsEnum } from "@formance/formance-sdk/dist/sdk/models/shared";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.ledger.getTransaction({
  ledger: "ledger001",
  txid: 1234,
}).then((res: GetTransactionResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## listAccounts

List accounts from a ledger, sorted by address in descending order.

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { ListAccountsBalanceOperator, ListAccountsResponse } from "@formance/formance-sdk/dist/sdk/models/operations";
import { ErrorsEnum } from "@formance/formance-sdk/dist/sdk/models/shared";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.ledger.listAccounts({
  address: "users:.+",
  balance: 2400,
  balanceOperator: ListAccountsBalanceOperator.Gte,
  cursor: "aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ==",
  ledger: "ledger001",
  metadata: {
    "est": "quibusdam",
  },
  pageSize: 131797,
}).then((res: ListAccountsResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## listLogs

List the logs from a ledger, sorted by ID in descending order.

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { ListLogsResponse } from "@formance/formance-sdk/dist/sdk/models/operations";
import { ErrorsEnum, LogType } from "@formance/formance-sdk/dist/sdk/models/shared";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.ledger.listLogs({
  cursor: "aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ==",
  endTime: new Date("2021-07-27T01:56:50.693Z"),
  ledger: "ledger001",
  pageSize: 841386,
  startTime: new Date("2022-09-26T08:57:48.803Z"),
}).then((res: ListLogsResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## listTransactions

List transactions from a ledger, sorted by txid in descending order.

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { ListTransactionsResponse } from "@formance/formance-sdk/dist/sdk/models/operations";
import { ErrorsEnum } from "@formance/formance-sdk/dist/sdk/models/shared";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.ledger.listTransactions({
  account: "users:001",
  cursor: "aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ==",
  destination: "users:001",
  endTime: new Date("2022-08-08T19:05:24.174Z"),
  ledger: "ledger001",
  metadata: {
    "quos": "perferendis",
    "magni": "assumenda",
    "ipsam": "alias",
  },
  pageSize: 146441,
  reference: "ref:001",
  source: "users:001",
  startTime: new Date("2021-11-11T04:17:07.569Z"),
}).then((res: ListTransactionsResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## readStats

Get statistics from a ledger. (aggregate metrics on accounts and transactions)


### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { ReadStatsResponse } from "@formance/formance-sdk/dist/sdk/models/operations";
import { ErrorsEnum } from "@formance/formance-sdk/dist/sdk/models/shared";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.ledger.readStats({
  ledger: "ledger001",
}).then((res: ReadStatsResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## revertTransaction

Revert a ledger transaction by its ID

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { RevertTransactionResponse } from "@formance/formance-sdk/dist/sdk/models/operations";
import { ErrorsEnum } from "@formance/formance-sdk/dist/sdk/models/shared";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.ledger.revertTransaction({
  ledger: "ledger001",
  txid: 1234,
}).then((res: RevertTransactionResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```
