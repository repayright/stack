/* 
 * Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.
 */

package com.formance.formance_sdk.models.operations;

import com.fasterxml.jackson.annotation.JsonProperty;
import com.formance.formance_sdk.utils.SpeakeasyMetadata;
import java.time.OffsetDateTime;

public class ListLogsRequest {
    /**
     * Pagination cursor, will return the logs after a given ID. (in descending order).
     */
    @SpeakeasyMetadata("queryParam:style=form,explode=true,name=after")
    public String after;

    public ListLogsRequest withAfter(String after) {
        this.after = after;
        return this;
    }
    
    /**
     * Parameter used in pagination requests. Maximum page size is set to 15.
     * Set to the value of next for the next page of results.
     * Set to the value of previous for the previous page of results.
     * No other parameters can be set when this parameter is set.
     * 
     */
    @SpeakeasyMetadata("queryParam:style=form,explode=true,name=cursor")
    public String cursor;

    public ListLogsRequest withCursor(String cursor) {
        this.cursor = cursor;
        return this;
    }
    
    /**
     * Filter transactions that occurred before this timestamp.
     * The format is RFC3339 and is exclusive (for example, "2023-01-02T15:04:01Z" excludes the first second of 4th minute).
     * 
     */
    @SpeakeasyMetadata("queryParam:style=form,explode=true,name=endTime")
    public OffsetDateTime endTime;

    public ListLogsRequest withEndTime(OffsetDateTime endTime) {
        this.endTime = endTime;
        return this;
    }
    
    /**
     * Name of the ledger.
     */
    @SpeakeasyMetadata("pathParam:style=simple,explode=false,name=ledger")
    public String ledger;

    public ListLogsRequest withLedger(String ledger) {
        this.ledger = ledger;
        return this;
    }
    
    /**
     * The maximum number of results to return per page.
     * 
     */
    @SpeakeasyMetadata("queryParam:style=form,explode=true,name=pageSize")
    public Long pageSize;

    public ListLogsRequest withPageSize(Long pageSize) {
        this.pageSize = pageSize;
        return this;
    }
    
    /**
     * Filter transactions that occurred after this timestamp.
     * The format is RFC3339 and is inclusive (for example, "2023-01-02T15:04:01Z" includes the first second of 4th minute).
     * 
     */
    @SpeakeasyMetadata("queryParam:style=form,explode=true,name=startTime")
    public OffsetDateTime startTime;

    public ListLogsRequest withStartTime(OffsetDateTime startTime) {
        this.startTime = startTime;
        return this;
    }
    
    public ListLogsRequest(@JsonProperty("ledger") String ledger) {
        this.ledger = ledger;
  }
}
