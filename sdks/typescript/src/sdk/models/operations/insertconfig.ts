/*
 * Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.
 */

import { SpeakeasyBase, SpeakeasyMetadata } from "../../../internal/utils";
import * as shared from "../shared";
import { AxiosResponse } from "axios";

export class InsertConfigResponse extends SpeakeasyBase {
  /**
   * Config created successfully.
   */
  @SpeakeasyMetadata()
  configResponse?: shared.ConfigResponse;

  @SpeakeasyMetadata()
  contentType: string;

  @SpeakeasyMetadata()
  statusCode: number;

  @SpeakeasyMetadata()
  rawResponse?: AxiosResponse;

  /**
   * Error
   */
  @SpeakeasyMetadata()
  webhooksErrorResponse?: shared.WebhooksErrorResponse;
}
