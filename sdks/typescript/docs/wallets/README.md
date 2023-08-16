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
  holdId: "aspernatur",
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
    expiresAt: new Date("2022-09-14T03:02:47.808Z"),
    name: "Miriam Hermann",
    priority: 574325,
  },
  id: "0afa563e-2516-4fe4-88b7-11e5b7fd2ed0",
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
    "praesentium": "natus",
  },
  name: "Joan Satterfield",
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
      amount: 807319,
      asset: "ea",
    },
    balance: "excepturi",
    metadata: {
      "ea": "accusantium",
    },
    reference: "ab",
    sources: [
      {
        balance: "ipsam",
        identifier: "voluptate",
        type: "autem",
      },
      {
        balance: "eaque",
        identifier: "pariatur",
        type: "nemo",
      },
      {
        balance: "perferendis",
        identifier: "fugiat",
        type: "amet",
      },
      {
        identifier: "cumque",
        type: "corporis",
      },
    ],
  },
  id: "fbb25870-5320-42c7-bd5f-e9b90c28909b",
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
      amount: 237893,
      asset: "asperiores",
    },
    balances: [
      "modi",
      "iste",
      "dolorum",
      "deleniti",
    ],
    description: "pariatur",
    destination: {
      balance: "nobis",
      identifier: "libero",
      type: "delectus",
    },
    metadata: {
      "quos": "aliquid",
      "dolorem": "dolorem",
    },
    pending: false,
  },
  id: "323f9b77-f3a4-4100-a74e-bf69280d1ba7",
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
  balanceName: "voluptate",
  id: "a89ebf73-7ae4-4203-8e5e-6a95d8a0d446",
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
  holdID: "quod",
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
    "qui": "dolorum",
    "a": "esse",
    "harum": "iusto",
    "ipsum": "quisquam",
  },
  pageSize: 947371,
  walletID: "amet",
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
  pageSize: 730856,
  walletID: "accusamus",
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
  id: "453f870b-326b-45a7-b429-cdb1a8422bb6",
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
  id: "79d23227-15bf-40cb-b1e3-1b8b90f3443a",
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
  id: "1108e0ad-cf4b-4921-879f-ce953f73ef7f",
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
    "quod": "odio",
    "similique": "facilis",
    "vero": "ducimus",
  },
  name: "Rosemarie Stanton",
  pageSize: 773326,
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
      "voluptatibus": "exercitationem",
    },
  },
  id: "d2cff7c7-0a45-4626-9436-813f16d9f5fc",
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
  holdId: "saepe",
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
