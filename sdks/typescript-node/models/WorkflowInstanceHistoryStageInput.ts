/**
 * Formance Stack API
 * Open, modular foundation for unique payments flows  # Introduction This API is documented in **OpenAPI format**.  # Authentication Formance Stack offers one forms of authentication:   - OAuth2 OAuth2 - an open protocol to allow secure authorization in a simple and standard method from web, mobile and desktop applications. <SecurityDefinitions /> 
 *
 * OpenAPI spec version: develop
 * Contact: support@formance.com
 *
 * NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).
 * https://openapi-generator.tech
 * Do not edit the class manually.
 */

import { ActivityConfirmHold } from '../models/ActivityConfirmHold';
import { ActivityCreateTransaction } from '../models/ActivityCreateTransaction';
import { ActivityCreditWallet } from '../models/ActivityCreditWallet';
import { ActivityDebitWallet } from '../models/ActivityDebitWallet';
import { ActivityGetAccount } from '../models/ActivityGetAccount';
import { ActivityGetPayment } from '../models/ActivityGetPayment';
import { ActivityGetWallet } from '../models/ActivityGetWallet';
import { ActivityRevertTransaction } from '../models/ActivityRevertTransaction';
import { ActivityVoidHold } from '../models/ActivityVoidHold';
import { StripeTransferRequest } from '../models/StripeTransferRequest';
import { HttpFile } from '../http/http';

export class WorkflowInstanceHistoryStageInput {
    'getAccount'?: ActivityGetAccount;
    'createTransaction'?: ActivityCreateTransaction;
    'revertTransaction'?: ActivityRevertTransaction;
    'stripeTransfer'?: StripeTransferRequest;
    'getPayment'?: ActivityGetPayment;
    'confirmHold'?: ActivityConfirmHold;
    'creditWallet'?: ActivityCreditWallet;
    'debitWallet'?: ActivityDebitWallet;
    'getWallet'?: ActivityGetWallet;
    'voidHold'?: ActivityVoidHold;

    static readonly discriminator: string | undefined = undefined;

    static readonly attributeTypeMap: Array<{name: string, baseName: string, type: string, format: string}> = [
        {
            "name": "getAccount",
            "baseName": "GetAccount",
            "type": "ActivityGetAccount",
            "format": ""
        },
        {
            "name": "createTransaction",
            "baseName": "CreateTransaction",
            "type": "ActivityCreateTransaction",
            "format": ""
        },
        {
            "name": "revertTransaction",
            "baseName": "RevertTransaction",
            "type": "ActivityRevertTransaction",
            "format": ""
        },
        {
            "name": "stripeTransfer",
            "baseName": "StripeTransfer",
            "type": "StripeTransferRequest",
            "format": ""
        },
        {
            "name": "getPayment",
            "baseName": "GetPayment",
            "type": "ActivityGetPayment",
            "format": ""
        },
        {
            "name": "confirmHold",
            "baseName": "ConfirmHold",
            "type": "ActivityConfirmHold",
            "format": ""
        },
        {
            "name": "creditWallet",
            "baseName": "CreditWallet",
            "type": "ActivityCreditWallet",
            "format": ""
        },
        {
            "name": "debitWallet",
            "baseName": "DebitWallet",
            "type": "ActivityDebitWallet",
            "format": ""
        },
        {
            "name": "getWallet",
            "baseName": "GetWallet",
            "type": "ActivityGetWallet",
            "format": ""
        },
        {
            "name": "voidHold",
            "baseName": "VoidHold",
            "type": "ActivityVoidHold",
            "format": ""
        }    ];

    static getAttributeTypeMap() {
        return WorkflowInstanceHistoryStageInput.attributeTypeMap;
    }

    public constructor() {
    }
}

