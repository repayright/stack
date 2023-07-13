# wallets

### Available Operations

* [confirmHold](#confirmhold) - Confirm a hold
* [createBalance](#createbalance) - Create a balance
* [createWallet](#createwallet) - Create a new wallet
* [creditWallet](#creditwallet) - Credit a wallet
* [debitWallet](#debitwallet) - Debit a wallet
* [getBalance](#getbalance) - Get detailed balance
* [getHold](#gethold) - Get a hold
* [getHolds](#getholds) - Get all holds for a wallet
* [getTransactions](#gettransactions)
* [getWallet](#getwallet) - Get a wallet
* [getWalletSummary](#getwalletsummary) - Get wallet summary
* [listBalances](#listbalances) - List balances of a wallet
* [listWallets](#listwallets) - List all wallets
* [updateWallet](#updatewallet) - Update a wallet
* [voidHold](#voidhold) - Cancel a hold
* [walletsgetServerInfo](#walletsgetserverinfo) - Get server info

## confirmHold

Confirm a hold

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { ConfirmHoldResponse } from "@formance/formance-sdk/dist/sdk/models/operations";
import { WalletsErrorResponseErrorCode } from "@formance/formance-sdk/dist/sdk/models/shared";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.wallets.confirmHold({
  confirmHoldRequest: {
    amount: 100,
    final: true,
  },
  holdId: "nisi",
}).then((res: ConfirmHoldResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## createBalance

Create a balance

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { CreateBalanceResponse } from "@formance/formance-sdk/dist/sdk/models/operations";
import { WalletsErrorResponseErrorCode } from "@formance/formance-sdk/dist/sdk/models/shared";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.wallets.createBalance({
  createBalanceRequest: {
    expiresAt: new Date("2022-05-20T03:14:12.839Z"),
    name: "Fernando Aufderhar",
    priority: 716075,
  },
  id: "a4469b6e-2141-4959-890a-fa563e2516fe",
}).then((res: CreateBalanceResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## createWallet

Create a new wallet

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { CreateWalletResponse } from "@formance/formance-sdk/dist/sdk/models/operations";
import { WalletsErrorResponseErrorCode } from "@formance/formance-sdk/dist/sdk/models/shared";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.wallets.createWallet({
  metadata: {
    "maxime": "deleniti",
    "facilis": "in",
  },
  name: "Diane VonRueden",
}).then((res: CreateWalletResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## creditWallet

Credit a wallet

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { CreditWalletResponse } from "@formance/formance-sdk/dist/sdk/models/operations";
import { WalletsErrorResponseErrorCode } from "@formance/formance-sdk/dist/sdk/models/shared";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.wallets.creditWallet({
  creditWalletRequest: {
    amount: {
      amount: 469249,
      asset: "repellat",
    },
    balance: "quibusdam",
    metadata: {
      "saepe": "pariatur",
    },
    reference: "accusantium",
    sources: [
      {
        balance: "natus",
        identifier: "magni",
        type: "sunt",
      },
    ],
  },
  id: "cddc6926-01fb-4576-b0d5-f0d30c5fbb25",
}).then((res: CreditWalletResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## debitWallet

Debit a wallet

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { DebitWalletResponse } from "@formance/formance-sdk/dist/sdk/models/operations";
import { WalletsErrorResponseErrorCode } from "@formance/formance-sdk/dist/sdk/models/shared";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.wallets.debitWallet({
  debitWalletRequest: {
    amount: {
      amount: 521037,
      asset: "dignissimos",
    },
    balances: [
      "quis",
    ],
    description: "nesciunt",
    destination: {
      identifier: "perferendis",
      type: "dolores",
    },
    metadata: {
      "quam": "dolor",
      "vero": "nostrum",
      "hic": "recusandae",
      "omnis": "facilis",
    },
    pending: false,
  },
  id: "90c28909-b3fe-449a-8d9c-bf48633323f9",
}).then((res: DebitWalletResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## getBalance

Get detailed balance

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { GetBalanceResponse } from "@formance/formance-sdk/dist/sdk/models/operations";
import { WalletsErrorResponseErrorCode } from "@formance/formance-sdk/dist/sdk/models/shared";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.wallets.getBalance({
  balanceName: "cum",
  id: "77f3a410-0674-4ebf-a928-0d1ba77a89eb",
}).then((res: GetBalanceResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## getHold

Get a hold

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { GetHoldResponse } from "@formance/formance-sdk/dist/sdk/models/operations";
import { WalletsErrorResponseErrorCode } from "@formance/formance-sdk/dist/sdk/models/shared";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.wallets.getHold({
  holdID: "asperiores",
}).then((res: GetHoldResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## getHolds

Get all holds for a wallet

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { GetHoldsResponse } from "@formance/formance-sdk/dist/sdk/models/operations";
import { WalletsErrorResponseErrorCode } from "@formance/formance-sdk/dist/sdk/models/shared";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.wallets.getHolds({
  cursor: "aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ==",
  metadata: {
    "ipsum": "voluptate",
    "id": "saepe",
  },
  pageSize: 263322,
  walletID: "aspernatur",
}).then((res: GetHoldsResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## getTransactions

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { GetTransactionsResponse } from "@formance/formance-sdk/dist/sdk/models/operations";
import { WalletsErrorResponseErrorCode } from "@formance/formance-sdk/dist/sdk/models/shared";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.wallets.getTransactions({
  cursor: "aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ==",
  pageSize: 20651,
  walletID: "amet",
}).then((res: GetTransactionsResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## getWallet

Get a wallet

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { GetWalletResponse } from "@formance/formance-sdk/dist/sdk/models/operations";
import { WalletsErrorResponseErrorCode } from "@formance/formance-sdk/dist/sdk/models/shared";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.wallets.getWallet({
  id: "ce5e6a95-d8a0-4d44-ace2-af7a73cf3be4",
}).then((res: GetWalletResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## getWalletSummary

Get wallet summary

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { GetWalletSummaryResponse } from "@formance/formance-sdk/dist/sdk/models/operations";
import { WalletsErrorResponseErrorCode } from "@formance/formance-sdk/dist/sdk/models/shared";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.wallets.getWalletSummary({
  id: "53f870b3-26b5-4a73-829c-db1a8422bb67",
}).then((res: GetWalletSummaryResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## listBalances

List balances of a wallet

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { ListBalancesResponse } from "@formance/formance-sdk/dist/sdk/models/operations";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.wallets.listBalances({
  id: "9d232271-5bf0-4cbb-9e31-b8b90f3443a1",
}).then((res: ListBalancesResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## listWallets

List all wallets

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { ListWalletsResponse } from "@formance/formance-sdk/dist/sdk/models/operations";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.wallets.listWallets({
  cursor: "aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ==",
  metadata: {
    "aut": "quas",
  },
  name: "William Ortiz",
  pageSize: 984330,
}).then((res: ListWalletsResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## updateWallet

Update a wallet

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { UpdateWalletResponse } from "@formance/formance-sdk/dist/sdk/models/operations";
import { WalletsErrorResponseErrorCode } from "@formance/formance-sdk/dist/sdk/models/shared";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.wallets.updateWallet({
  requestBody: {
    metadata: {
      "facilis": "cupiditate",
      "qui": "quae",
    },
  },
  id: "879fce95-3f73-4ef7-bbc7-abd74dd39c0f",
}).then((res: UpdateWalletResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## voidHold

Cancel a hold

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { VoidHoldResponse } from "@formance/formance-sdk/dist/sdk/models/operations";
import { WalletsErrorResponseErrorCode } from "@formance/formance-sdk/dist/sdk/models/shared";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.wallets.voidHold({
  holdId: "exercitationem",
}).then((res: VoidHoldResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```

## walletsgetServerInfo

Get server info

### Example Usage

```typescript
import { SDK } from "@formance/formance-sdk";
import { WalletsgetServerInfoResponse } from "@formance/formance-sdk/dist/sdk/models/operations";
import { WalletsErrorResponseErrorCode } from "@formance/formance-sdk/dist/sdk/models/shared";

const sdk = new SDK({
  security: {
    authorization: "Bearer YOUR_ACCESS_TOKEN_HERE",
  },
});

sdk.wallets.walletsgetServerInfo().then((res: WalletsgetServerInfoResponse) => {
  if (res.statusCode == 200) {
    // handle response
  }
});
```
