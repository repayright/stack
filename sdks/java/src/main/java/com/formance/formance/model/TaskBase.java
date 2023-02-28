/*
 * Formance Stack API
 * Open, modular foundation for unique payments flows  # Introduction This API is documented in **OpenAPI format**.  # Authentication Formance Stack offers one forms of authentication:   - OAuth2 OAuth2 - an open protocol to allow secure authorization in a simple and standard method from web, mobile and desktop applications. <SecurityDefinitions /> 
 *
 * The version of the OpenAPI document: v1.0.20230228
 * Contact: support@formance.com
 *
 * NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).
 * https://openapi-generator.tech
 * Do not edit the class manually.
 */


package com.formance.formance.model;

import java.util.Objects;
import java.util.Arrays;
import com.formance.formance.model.PaymentStatus;
import com.google.gson.TypeAdapter;
import com.google.gson.annotations.JsonAdapter;
import com.google.gson.annotations.SerializedName;
import com.google.gson.stream.JsonReader;
import com.google.gson.stream.JsonWriter;
import java.io.IOException;
import java.time.OffsetDateTime;
import java.util.UUID;

/**
 * TaskBase
 */
@javax.annotation.Generated(value = "org.openapitools.codegen.languages.JavaClientCodegen")
public class TaskBase {
  public static final String SERIALIZED_NAME_ID = "id";
  @SerializedName(SERIALIZED_NAME_ID)
  private UUID id;

  public static final String SERIALIZED_NAME_CONNECTOR_ID = "connectorId";
  @SerializedName(SERIALIZED_NAME_CONNECTOR_ID)
  private UUID connectorId;

  public static final String SERIALIZED_NAME_CREATED_AT = "createdAt";
  @SerializedName(SERIALIZED_NAME_CREATED_AT)
  private OffsetDateTime createdAt;

  public static final String SERIALIZED_NAME_UPDATED_AT = "updatedAt";
  @SerializedName(SERIALIZED_NAME_UPDATED_AT)
  private OffsetDateTime updatedAt;

  public static final String SERIALIZED_NAME_STATUS = "status";
  @SerializedName(SERIALIZED_NAME_STATUS)
  private PaymentStatus status;

  public static final String SERIALIZED_NAME_STATE = "state";
  @SerializedName(SERIALIZED_NAME_STATE)
  private Object state;

  public static final String SERIALIZED_NAME_ERROR = "error";
  @SerializedName(SERIALIZED_NAME_ERROR)
  private String error;

  public TaskBase() {
  }

  public TaskBase id(UUID id) {
    
    this.id = id;
    return this;
  }

   /**
   * Get id
   * @return id
  **/
  @javax.annotation.Nonnull

  public UUID getId() {
    return id;
  }


  public void setId(UUID id) {
    this.id = id;
  }


  public TaskBase connectorId(UUID connectorId) {
    
    this.connectorId = connectorId;
    return this;
  }

   /**
   * Get connectorId
   * @return connectorId
  **/
  @javax.annotation.Nonnull

  public UUID getConnectorId() {
    return connectorId;
  }


  public void setConnectorId(UUID connectorId) {
    this.connectorId = connectorId;
  }


  public TaskBase createdAt(OffsetDateTime createdAt) {
    
    this.createdAt = createdAt;
    return this;
  }

   /**
   * Get createdAt
   * @return createdAt
  **/
  @javax.annotation.Nonnull

  public OffsetDateTime getCreatedAt() {
    return createdAt;
  }


  public void setCreatedAt(OffsetDateTime createdAt) {
    this.createdAt = createdAt;
  }


  public TaskBase updatedAt(OffsetDateTime updatedAt) {
    
    this.updatedAt = updatedAt;
    return this;
  }

   /**
   * Get updatedAt
   * @return updatedAt
  **/
  @javax.annotation.Nonnull

  public OffsetDateTime getUpdatedAt() {
    return updatedAt;
  }


  public void setUpdatedAt(OffsetDateTime updatedAt) {
    this.updatedAt = updatedAt;
  }


  public TaskBase status(PaymentStatus status) {
    
    this.status = status;
    return this;
  }

   /**
   * Get status
   * @return status
  **/
  @javax.annotation.Nonnull

  public PaymentStatus getStatus() {
    return status;
  }


  public void setStatus(PaymentStatus status) {
    this.status = status;
  }


  public TaskBase state(Object state) {
    
    this.state = state;
    return this;
  }

   /**
   * Get state
   * @return state
  **/
  @javax.annotation.Nonnull

  public Object getState() {
    return state;
  }


  public void setState(Object state) {
    this.state = state;
  }


  public TaskBase error(String error) {
    
    this.error = error;
    return this;
  }

   /**
   * Get error
   * @return error
  **/
  @javax.annotation.Nullable

  public String getError() {
    return error;
  }


  public void setError(String error) {
    this.error = error;
  }


  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (o == null || getClass() != o.getClass()) {
      return false;
    }
    TaskBase taskBase = (TaskBase) o;
    return Objects.equals(this.id, taskBase.id) &&
        Objects.equals(this.connectorId, taskBase.connectorId) &&
        Objects.equals(this.createdAt, taskBase.createdAt) &&
        Objects.equals(this.updatedAt, taskBase.updatedAt) &&
        Objects.equals(this.status, taskBase.status) &&
        Objects.equals(this.state, taskBase.state) &&
        Objects.equals(this.error, taskBase.error);
  }

  @Override
  public int hashCode() {
    return Objects.hash(id, connectorId, createdAt, updatedAt, status, state, error);
  }

  @Override
  public String toString() {
    StringBuilder sb = new StringBuilder();
    sb.append("class TaskBase {\n");
    sb.append("    id: ").append(toIndentedString(id)).append("\n");
    sb.append("    connectorId: ").append(toIndentedString(connectorId)).append("\n");
    sb.append("    createdAt: ").append(toIndentedString(createdAt)).append("\n");
    sb.append("    updatedAt: ").append(toIndentedString(updatedAt)).append("\n");
    sb.append("    status: ").append(toIndentedString(status)).append("\n");
    sb.append("    state: ").append(toIndentedString(state)).append("\n");
    sb.append("    error: ").append(toIndentedString(error)).append("\n");
    sb.append("}");
    return sb.toString();
  }

  /**
   * Convert the given object to string with each line indented by 4 spaces
   * (except the first line).
   */
  private String toIndentedString(Object o) {
    if (o == null) {
      return "null";
    }
    return o.toString().replace("\n", "\n    ");
  }

}

