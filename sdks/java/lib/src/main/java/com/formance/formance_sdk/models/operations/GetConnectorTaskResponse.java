/* 
 * Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.
 */

package com.formance.formance_sdk.models.operations;

import com.fasterxml.jackson.annotation.JsonProperty;
import java.net.http.HttpResponse;

public class GetConnectorTaskResponse {
    
    public String contentType;

    public GetConnectorTaskResponse withContentType(String contentType) {
        this.contentType = contentType;
        return this;
    }
    
    
    public Integer statusCode;

    public GetConnectorTaskResponse withStatusCode(Integer statusCode) {
        this.statusCode = statusCode;
        return this;
    }
    
    
    public HttpResponse<byte[]> rawResponse;

    public GetConnectorTaskResponse withRawResponse(HttpResponse<byte[]> rawResponse) {
        this.rawResponse = rawResponse;
        return this;
    }
    
    /**
     * OK
     */
    
    public com.formance.formance_sdk.models.shared.TaskResponse taskResponse;

    public GetConnectorTaskResponse withTaskResponse(com.formance.formance_sdk.models.shared.TaskResponse taskResponse) {
        this.taskResponse = taskResponse;
        return this;
    }
    
    public GetConnectorTaskResponse(@JsonProperty("ContentType") String contentType, @JsonProperty("StatusCode") Integer statusCode) {
        this.contentType = contentType;
        this.statusCode = statusCode;
  }
}