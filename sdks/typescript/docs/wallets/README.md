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
  holdId: "maxime",
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
    expiresAt: new Date("2021-08-05T03:52:18.835Z"),
    name: "Alice Bradtke",
    priority: 714242,
  },
  id: "7fd2ed02-8921-4cdd-8692-601fb576b0d5",
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
    "perferendis": "fugiat",
    "amet": "aut",
    "cumque": "corporis",
    "hic": "libero",
  },
  name: "Ernest Hayes",
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
      amount: 54338,
      asset: "quis",
    },
    balance: "nesciunt",
    metadata: {
      "perferendis": "dolores",
    },
    reference: "minus",
    sources: [
      {
        identifier: "vero",
        type: "nostrum",
      },
      {
        balance: "recusandae",
        identifier: "omnis",
        type: "facilis",
      },
    ],
  },
  id: "90c28909-b3fe-449a-8d9c-bf48633323f9",
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
      amount: 739551,
      asset: "voluptate",
    },
    balances: [
      "reiciendis",
      "amet",
    ],
    description: "dolorum",
    destination: {
      identifier: "veritatis",
      type: "ipsa",
    },
    metadata: {
      "iure": "odio",
    },
    pending: false,
  },
  id: "4ebf6928-0d1b-4a77-a89e-bf737ae4203c",
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
  balanceName: "accusamus",
  id: "5e6a95d8-a0d4-446c-a2af-7a73cf3be453",
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
  holdID: "sapiente",
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
    "nihil": "sit",
    "expedita": "neque",
    "sed": "vel",
  },
  pageSize: 730442,
  walletID: "voluptas",
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
  pageSize: 646265,
  walletID: "quam",
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
  id: "3429cdb1-a842-42bb-a79d-2322715bf0cb",
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
  id: "b1e31b8b-90f3-4443-a110-8e0adcf4b921",
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
  id: "879fce95-3f73-4ef7-bbc7-abd74dd39c0f",
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
    "nulla": "fugit",
    "porro": "maiores",
  },
  name: "Ted Romaguera MD",
  pageSize: 269479,
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
      "ea": "aspernatur",
      "vel": "possimus",
    },
  },
  id: "436813f1-6d9f-45fc-a6c5-56146c3e250f",
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
  holdId: "libero",
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
